package clinical

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// constants and defaults
const (
	timeFormatStr = "2006-01-02T15:04:05+03:00"
	dateFormatStr = "2006-01-02"
)

// CreateEpisodeOfCare creates an episode of care. An Episode of Care represents a period of time during which a patient is under the care of a particular provider/facility.
// An Episode of Care includes one or more encounters.
func (c *UseCasesClinicalImpl) CreateEpisodeOfCare(ctx context.Context, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error) {
	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	active := domain.EpisodeOfCareStatusEnum(string(input.Status))

	episodeOfCare := &domain.FHIREpisodeOfCareInput{
		Status: &active,
		Period: common.DefaultPeriodInput(),
	}

	facility, err := c.infrastructure.FHIR.GetFHIROrganization(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	orgRef := fmt.Sprintf("Organization/%s", *facility.Resource.ID)
	orgType := scalarutils.URI("Organization")

	episodeOfCare.ManagingOrganization = &domain.FHIRReferenceInput{
		ID:        facility.Resource.ID,
		Reference: &orgRef,
		Display:   *facility.Resource.Name,
		Type:      &orgType,
	}

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	patientType := scalarutils.URI("Patient")

	episodeOfCare.Patient = &domain.FHIRReferenceInput{
		ID:        patient.Resource.ID,
		Reference: &patientRef,
		Display:   patient.Resource.Name[0].Text,
		Type:      &patientType,
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	// search for the episode of care before creating new one.
	episodeOfCareSearchParams := map[string]interface{}{
		"patient":      patientRef,
		"status":       string(domain.EpisodeOfCareStatusEnumActive),
		"organization": orgRef,
		"_sort":        "date",
		"_count":       "1",
	}

	episodeOfCarePayload, err := c.infrastructure.FHIR.SearchFHIREpisodeOfCare(ctx, episodeOfCareSearchParams, *identifiers, dto.Pagination{})
	if err != nil {
		return nil, fmt.Errorf("unable to get patients episodes of care: %w", err)
	}

	// don't create a new episode if there is an ongoing one
	if len(episodeOfCarePayload.Edges) >= 1 {
		return nil, fmt.Errorf("an active episode of care already exists")
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	episodeOfCare.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	episode, err := c.infrastructure.FHIR.CreateEpisodeOfCare(ctx, *episodeOfCare)
	if err != nil {
		return nil, err
	}

	return mapFHIREpisodeToEpisodeDTO(*episode.EpisodeOfCare), nil
}

func (c *UseCasesClinicalImpl) PatchEpisodeOfCare(ctx context.Context, id string, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid episode of care id: %s", id)
	}

	status := domain.EpisodeOfCareStatusEnum(strings.ToLower(string(input.Status)))
	episodeOfCareInput := &domain.FHIREpisodeOfCareInput{
		Status: &status,
	}

	if status.IsFinal() {
		episode, err := c.infrastructure.FHIR.GetFHIREpisodeOfCare(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("unable to get episode with ID %s: %w", id, err)
		}

		var startTime scalarutils.DateTime
		if episode.Resource.Period == nil {
			startTime = scalarutils.DateTime(time.Now().Format(timeFormatStr))
		} else {
			startTime = episode.Resource.Period.Start
		}

		// workaround for odd date comparison behavior on the Google Cloud Healthcare API
		end := startTime.Time().Add(time.Hour * 24)
		endTime := scalarutils.DateTime(end.Format(timeFormatStr))

		episodeOfCareInput.Period = &domain.FHIRPeriodInput{Start: startTime, End: endTime}
	}

	episode, err := c.infrastructure.FHIR.PatchFHIREpisodeOfCare(ctx, id, *episodeOfCareInput)
	if err != nil {
		return nil, err
	}

	return mapFHIREpisodeToEpisodeDTO(*episode), nil
}

func (c *UseCasesClinicalImpl) EndEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid episode of care id: %s", id)
	}

	episode, err := c.infrastructure.FHIR.GetFHIREpisodeOfCare(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get episode of care: %w", err)
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	// Close all encounters in this visit
	encounterConn, err := c.infrastructure.FHIR.SearchEpisodeEncounter(ctx, id, *identifiers, dto.Pagination{})
	if err != nil {
		return nil, fmt.Errorf("unable to search episode encounter %w", err)
	}

	for _, edge := range encounterConn.Encounters {
		_, err = c.infrastructure.FHIR.EndEncounter(ctx, *edge.ID)
		if err != nil {
			return nil, fmt.Errorf("unable to end encounter %s: err: %w", *edge.ID, err)
		}
	}

	_, err = c.infrastructure.FHIR.EndEpisode(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to end episode of care: %w", err)
	}

	return mapFHIREpisodeToEpisodeDTO(*episode.Resource), nil
}

func (c *UseCasesClinicalImpl) GetEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid episode of care id: %s", id)
	}

	episode, err := c.infrastructure.FHIR.GetFHIREpisodeOfCare(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get episode of care: %w", err)
	}

	return mapFHIREpisodeToEpisodeDTO(*episode.Resource), nil
}

func mapFHIREpisodeToEpisodeDTO(episode domain.FHIREpisodeOfCare) *dto.EpisodeOfCare {
	return &dto.EpisodeOfCare{
		ID:        *episode.ID,
		Status:    dto.EpisodeOfCareStatusEnum(string(*episode.Status)),
		PatientID: *episode.Patient.ID,
	}
}
