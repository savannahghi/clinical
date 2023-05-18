package clinical

import (
	"context"
	"fmt"
	"sync"

	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"

	linq "github.com/ahmetb/go-linq/v3"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/scalarutils"
	log "github.com/sirupsen/logrus"
)

// PatientTimeline return's the patient's historical timeline sorted in descending order i.e when it was first recorded
// The timeline consists of Allergies, Observations, Medication statement and Test results
func (c *UseCasesClinicalImpl) PatientTimeline(ctx context.Context, patientID string) ([]dto.TimelineResource, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	timeline := []dto.TimelineResource{}
	wg := &sync.WaitGroup{}
	mut := &sync.Mutex{}

	patientFilterParams := map[string]interface{}{
		"patient": fmt.Sprintf("Patient/%v", patientID),
	}

	// timelineResourceFunc is a go routine that fetches particular FHIR resource and
	// adds it to the timeline
	type timelineResourceFunc func(wg *sync.WaitGroup, mut *sync.Mutex)

	allergyIntoleranceResourceFunc := func(wg *sync.WaitGroup, mut *sync.Mutex) {
		defer wg.Done()

		conn, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, patientFilterParams, *identifiers, dto.Pagination{Skip: true})
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("AllergyIntolerance search error: %v", err)

			return
		}

		for _, edge := range conn.Allergies {
			if edge.ID == nil {
				continue
			}

			if edge.Code == nil {
				continue
			}

			if edge.Reaction == nil {
				continue
			}

			if len(edge.Reaction) < 1 {
				continue
			}

			if edge.Reaction[0].Manifestation == nil {
				continue
			}

			if len(edge.Reaction[0].Manifestation) < 1 {
				continue
			}

			if edge.RecordedDate == nil {
				continue
			}

			timelineResource := dto.TimelineResource{
				ID:           *edge.ID,
				ResourceType: dto.ResourceTypeAllergyIntolerance,
				Name:         edge.Code.Text,
				Value:        edge.Reaction[0].Manifestation[0].Text,
				Status:       edge.ClinicalStatus.Text,
				Date:         *edge.RecordedDate,
			}

			mut.Lock()
			timeline = append(timeline, timelineResource)
			mut.Unlock()
		}
	}

	observationResourceFunc := func(wg *sync.WaitGroup, mut *sync.Mutex) {
		defer wg.Done()

		conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, patientFilterParams, *identifiers, dto.Pagination{Skip: true})
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("Observation search error: %v", err)

			return
		}

		for _, edge := range conn.Observations {
			if edge.ID == nil {
				continue
			}

			if edge.Code.Coding == nil {
				continue
			}

			if len(edge.Code.Coding) < 1 {
				continue
			}

			if edge.Status == nil {
				continue
			}

			if edge.EffectiveInstant == nil {
				continue
			}

			instant := helpers.ParseDate(string(*edge.EffectiveInstant))
			date, err := scalarutils.NewDate(instant.Day(), int(instant.Month()), instant.Year())

			if err != nil {
				utils.ReportErrorToSentry(err)
				log.Errorf("date conversion error: %v", err)

				return
			}

			timelineResource := dto.TimelineResource{
				ID:           *edge.ID,
				ResourceType: dto.ResourceTypeObservation,
				Name:         edge.Code.Text,
				Value:        *edge.ValueString,
				Status:       string(*edge.Status),
				Date:         *date,
			}

			mut.Lock()
			timeline = append(timeline, timelineResource)
			mut.Unlock()
		}
	}

	medicationStatementResourceFunc := func(wg *sync.WaitGroup, mut *sync.Mutex) {
		defer wg.Done()

		conn, err := c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, patientFilterParams, *identifiers, dto.Pagination{Skip: true})
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("MedicationStatement search error: %v", err)

			return
		}

		for _, edge := range conn.Edges {
			if edge.Node == nil {
				continue
			}

			if edge.Node.ID == nil {
				continue
			}

			if edge.Node.MedicationCodeableConcept == nil {
				continue
			}

			if edge.Node.MedicationCodeableConcept.Coding == nil {
				continue
			}

			if len(edge.Node.MedicationCodeableConcept.Coding) < 1 {
				continue
			}

			if edge.Node.Status == nil {
				continue
			}

			if edge.Node.Category == nil {
				continue
			}

			if edge.Node.EffectiveDateTime == nil {
				continue
			}

			instant := edge.Node.EffectiveDateTime.AsTime()
			date, err := scalarutils.NewDate(instant.Day(), int(instant.Month()), instant.Year())

			if err != nil {
				utils.ReportErrorToSentry(err)
				log.Errorf("date conversion error: %v", err)

				return
			}

			timelineResource := dto.TimelineResource{
				ID:           *edge.Node.ID,
				ResourceType: dto.ResourceTypeMedicationStatement,
				Name:         edge.Node.Category.Text,
				Value:        edge.Node.MedicationCodeableConcept.Coding[0].Display,
				Status:       string(*edge.Node.Status),
				Date:         *date,
			}

			mut.Lock()
			timeline = append(timeline, timelineResource)
			mut.Unlock()
		}
	}

	resources := []timelineResourceFunc{
		allergyIntoleranceResourceFunc,
		observationResourceFunc,
		medicationStatementResourceFunc,
	}

	for _, resource := range resources {
		wg.Add(1)

		go resource(wg, mut)
	}

	wg.Wait()

	return timeline, nil
}

// PatientHealthTimeline return's the patient's historical timeline sorted in descending order i.e when it was first recorded
// The timeline consists of Allergies, Observations, Medication statement and Test results
func (c *UseCasesClinicalImpl) PatientHealthTimeline(ctx context.Context, input dto.HealthTimelineInput) (*dto.HealthTimeline, error) {
	records, err := c.PatientTimeline(ctx, input.PatientID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("cannot retrieve patient timeline error: %w", err)
	}

	data := &dto.HealthTimeline{}
	timeline := []dto.TimelineResource{}

	sortFunc := func(i, j interface{}) bool {
		itemI := i.(dto.TimelineResource)
		timeI := itemI.Date.AsTime()

		itemJ := j.(dto.TimelineResource)
		timeJ := itemJ.Date.AsTime()

		return timeI.After(timeJ)
	}

	linq.From(records).Sort(sortFunc).Skip(input.Offset).Take(input.Limit).ToSlice(&timeline)

	data.TotalCount = len(records)
	data.Timeline = timeline

	return data, nil
}
