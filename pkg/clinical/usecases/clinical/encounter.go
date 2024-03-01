package clinical

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// StartEncounter starts an encounter within an episode of care
func (c *UseCasesClinicalImpl) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	if episodeID == "" {
		return "", fmt.Errorf("an episode of care ID is required")
	}

	episodeOfCare, err := c.infrastructure.FHIR.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return "", err
	}

	encounterClassCode := scalarutils.Code("AMB")
	encounterClassSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/v3-ActCode")
	encounterClassVersion := "2018-08-12"
	encounterClassDisplay := string(dto.EncounterClassEnumAmbulatory)
	encounterClassUserSelected := false
	now := time.Now()
	startTime := scalarutils.DateTime(now.Format("2006-01-02T15:04:05+03:00"))

	episodeReference := fmt.Sprintf("EpisodeOfCare/%s", *episodeOfCare.Resource.ID)
	encounterPayload := domain.FHIREncounterInput{
		Status: domain.EncounterStatusEnum(dto.EncounterStatusEnumInProgress),
		Class: domain.FHIRCodingInput{
			System:       &encounterClassSystem,
			Version:      &encounterClassVersion,
			Code:         encounterClassCode,
			Display:      encounterClassDisplay,
			UserSelected: &encounterClassUserSelected,
		},
		Subject: &domain.FHIRReferenceInput{
			ID:        episodeOfCare.Resource.Patient.ID,
			Reference: episodeOfCare.Resource.Patient.Reference,
			Type:      episodeOfCare.Resource.Patient.Type,
			Display:   episodeOfCare.Resource.Patient.Display,
		},
		EpisodeOfCare: []*domain.FHIRReferenceInput{
			{
				ID:        episodeOfCare.Resource.ID,
				Reference: &episodeReference,
			},
		},
		ServiceProvider: &domain.FHIRReferenceInput{
			Reference: episodeOfCare.Resource.ManagingOrganization.Reference,
			Type:      episodeOfCare.Resource.ManagingOrganization.Type,
			Display:   episodeOfCare.Resource.ManagingOrganization.Display,
		},
		Period: &domain.FHIRPeriodInput{
			Start: startTime,
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return "", err
	}

	encounterPayload.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	encounter, err := c.infrastructure.FHIR.CreateFHIREncounter(ctx, encounterPayload)
	if err != nil {
		return "", err
	}

	return *encounter.Resource.ID, nil
}

func (c *UseCasesClinicalImpl) PatchEncounter(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error) {
	if encounterID == "" {
		return nil, fmt.Errorf("an encounterID is required")
	}

	status := domain.EncounterStatusEnum(strings.ToLower(string(input.Status)))
	encounterInput := domain.FHIREncounterInput{
		Status: status,
	}

	if status.IsFinal() {
		encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, encounterID)
		if err != nil {
			return nil, fmt.Errorf("unable to get encounter with ID %s: %w", encounterID, err)
		}

		var startTime scalarutils.DateTime
		if encounter.Resource.Period == nil {
			startTime = scalarutils.DateTime(time.Now().Format(timeFormatStr))
		} else {
			startTime = encounter.Resource.Period.Start
		}

		// workaround for odd date comparison behavior on the Google Cloud Healthcare API
		end := startTime.Time().Add(time.Hour * 24)
		endTime := scalarutils.DateTime(end.Format(timeFormatStr))

		encounterInput.Period = &domain.FHIRPeriodInput{Start: startTime, End: endTime}
	}

	fhirEncounter, err := c.infrastructure.FHIR.PatchFHIREncounter(ctx, encounterID, encounterInput)

	if err != nil {
		return nil, err
	}

	encounters := []*dto.Encounter{}

	err = mapstructure.Decode([]domain.FHIREncounter{*fhirEncounter}, &encounters)

	if err != nil {
		return nil, err
	}

	return encounters[0], nil
}

// EndEncounter marks an encounter as finished and updates the endtime field
func (c *UseCasesClinicalImpl) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	if encounterID == "" {
		return false, fmt.Errorf("an encounterID is required")
	}

	ok, err := c.infrastructure.FHIR.EndEncounter(ctx, encounterID)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// ListPatientEncounters lists all the encounters that a patient has been part of
func (c *UseCasesClinicalImpl) ListPatientEncounters(ctx context.Context, patientID string, pagination *dto.Pagination) (*dto.EncounterConnection, error) {
	if patientID == "" {
		return nil, fmt.Errorf("a patient ID is required")
	}

	err := pagination.Validate()
	if err != nil {
		return nil, err
	}

	_, err = c.infrastructure.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", patientID)

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	encounterResponses, err := c.infrastructure.FHIR.SearchPatientEncounters(ctx, patientReference, nil, *identifiers, *pagination)
	if err != nil {
		return nil, err
	}

	encounters := []*dto.Encounter{}

	err = mapstructure.Decode(encounterResponses.Encounters, &encounters)
	if err != nil {
		return nil, err
	}

	pagedInfo := dto.PageInfo{
		HasNextPage:     encounterResponses.HasNextPage,
		EndCursor:       &encounterResponses.NextCursor,
		HasPreviousPage: encounterResponses.HasPreviousPage,
		StartCursor:     &encounterResponses.PreviousCursor,
	}

	connection := dto.CreateEncounterConnection(encounters, pagedInfo, encounterResponses.TotalCount)

	return &connection, nil
}

// GetEncounterAssociatedResources get all resources assocuated with an encounter
func (c *UseCasesClinicalImpl) GetEncounterAssociatedResources(ctx context.Context, encounterID string) (*dto.EncounterAssociatedResources, error) {
	if encounterID == "" {
		return nil, fmt.Errorf("an encounterID is required")
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	encounterSearchParams := map[string]interface{}{
		"_id":         encounterID,
		"_sort":       "date",
		"_count":      "1",
		"_revinclude": []string{"RiskAssessment:encounter", "Consent:data"},
	}

	encounterAllData, err := c.infrastructure.FHIR.SearchFHIREncounterAllData(ctx, encounterSearchParams, *identifiers, dto.Pagination{})
	if err != nil {
		return nil, err
	}

	result := dto.EncounterAssociatedResources{}

	for _, encounterData := range encounterAllData.Resources {
		switch encounterData["resourceType"] {
		case "RiskAssessment":
			var riskAssessment dto.RiskAssessment

			riskAssessmentBytes, err := json.Marshal(encounterData)

			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(riskAssessmentBytes, &riskAssessment); err != nil {
				return nil, err
			}

			result.RiskAssessment = &riskAssessment
		case "Consent":
			var consent dto.Consent

			consentBytes, err := json.Marshal(encounterData)

			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(consentBytes, &consent); err != nil {
				return nil, err
			}

			result.Consent = &consent
		}
	}

	return &result, nil
}
