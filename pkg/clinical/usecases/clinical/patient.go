package clinical

import (
	"context"
	"fmt"
	"sync"

	linq "github.com/ahmetb/go-linq/v3"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/scalarutils"
	log "github.com/sirupsen/logrus"
)

// UseCasesClinicalImpl represents the patient usecase implementation
type UseCasesClinicalImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewUseCasesClinicalImpl initializes new Clinical/Patient implementation
func NewUseCasesClinicalImpl(infra infrastructure.Infrastructure) *UseCasesClinicalImpl {
	return &UseCasesClinicalImpl{
		infrastructure: infra,
	}
}

// PatientTimeline return's the patient's historical timeline sorted in descending order i.e when it was first recorded
// The timeline consists of Allergies, Observations, Medication statement and Test results
func (c *UseCasesClinicalImpl) PatientTimeline(ctx context.Context, patientID string) ([]dto.TimelineResource, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
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

		conn, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, patientFilterParams)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("AllergyIntolerance search error: %v", err)

			return
		}

		for _, edge := range conn.Edges {
			if edge.Node == nil {
				continue
			}

			if edge.Node.ID == nil {
				continue
			}

			if edge.Node.Code == nil {
				continue
			}

			if edge.Node.Reaction == nil && len(edge.Node.Reaction) < 1 {
				continue
			}

			if edge.Node.Reaction[0].Manifestation == nil && len(edge.Node.Reaction[0].Manifestation) < 1 {
				continue
			}

			if edge.Node.RecordedDate == nil {
				continue
			}

			timelineResource := dto.TimelineResource{
				ID:           *edge.Node.ID,
				ResourceType: dto.ResourceTypeAllergyIntolerance,
				Name:         edge.Node.Code.Text,
				Value:        edge.Node.Reaction[0].Manifestation[0].Text,
				Status:       edge.Node.ClinicalStatus.Text,
				Date:         *edge.Node.RecordedDate,
			}

			mut.Lock()
			timeline = append(timeline, timelineResource)
			mut.Unlock()
		}
	}

	observationResourceFunc := func(wg *sync.WaitGroup, mut *sync.Mutex) {
		defer wg.Done()

		conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, patientFilterParams)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("Observation search error: %v", err)

			return
		}

		for _, edge := range conn.Edges {
			if edge.Node == nil {
				continue
			}

			if edge.Node.ID == nil {
				continue
			}

			if edge.Node.Code.Coding == nil && len(edge.Node.Code.Coding) < 1 {
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
				ResourceType: dto.ResourceTypeObservation,
				Name:         edge.Node.Code.Text,
				Value:        edge.Node.Code.Coding[0].Display,
				Status:       string(*edge.Node.Status),
				Date:         *date,
			}

			mut.Lock()
			timeline = append(timeline, timelineResource)
			mut.Unlock()
		}
	}

	medicationStatementResourceFunc := func(wg *sync.WaitGroup, mut *sync.Mutex) {
		defer wg.Done()

		conn, err := c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, patientFilterParams)
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

			if edge.Node.MedicationCodeableConcept.Coding == nil && len(edge.Node.MedicationCodeableConcept.Coding) < 1 {
				continue
			}

			if edge.Node.Status == nil {
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
				Name:         edge.Node.Subject.Display,
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

// GetMedicalData returns a limited subset of specific medical data that for a specific patient
// These include: Allergies, Viral Load, Body Mass Index, Weight, CD4 Count using their respective OCL CIEL Terminology
// For each category the latest three records are fetched
func (c *UseCasesClinicalImpl) GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error) {
	data := &domain.MedicalData{}

	filterParams := map[string]interface{}{
		"patient": fmt.Sprintf("Patient/%v", patientID),
		"_count":  common.MedicalDataCount,
		"_sort":   "-date",
	}

	fields := []string{
		"Regimen",
		"AllergyIntolerance",
		"Weight",
		"BMI",
		"ViralLoad",
		"CD4Count",
	}

	for _, field := range fields {
		switch field {
		case "Regimen":
			conn, err := c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				data.Regimen = append(data.Regimen, edge.Node)
			}
		case "AllergyIntolerance":
			conn, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				data.Allergies = append(data.Allergies, edge.Node)
			}

		case "Weight":
			filterParams["code"] = common.WeightCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				data.Weight = append(data.Weight, edge.Node)
			}

		case "BMI":
			filterParams["code"] = common.BMICIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				data.BMI = append(data.BMI, edge.Node)
			}

		case "ViralLoad":
			filterParams["code"] = common.ViralLoadCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				data.ViralLoad = append(data.ViralLoad, edge.Node)
			}

		case "CD4Count":
			filterParams["code"] = common.CD4CountCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				data.CD4Count = append(data.CD4Count, edge.Node)
			}
		}
	}

	return data, nil
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

// CreateFHIROrganization creates a FHIROrganization instance
func (c *UseCasesClinicalImpl) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	organizationRelayPayload, err := c.infrastructure.FHIR.CreateFHIROrganization(ctx, input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return organizationRelayPayload, nil
}
