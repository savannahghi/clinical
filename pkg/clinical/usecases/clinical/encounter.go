package clinical

import (
	"context"
	"fmt"
	"time"

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
	encounterClassDisplay := string(dto.EncounterClassAmbulatory)
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

	encounterPayload.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	encounter, err := c.infrastructure.FHIR.CreateFHIREncounter(ctx, encounterPayload)
	if err != nil {
		return "", err
	}

	return *encounter.Resource.ID, nil
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

	encounters := mapFHIREncounterToEncounterDTO(encounterResponses.Encounters)

	pagedInfo := dto.PageInfo{
		HasNextPage:     encounterResponses.HasNextPage,
		EndCursor:       &encounterResponses.NextCursor,
		HasPreviousPage: encounterResponses.HasPreviousPage,
		StartCursor:     &encounterResponses.PreviousCursor,
	}

	connection := dto.CreateEncounterConnection(encounters, pagedInfo, encounterResponses.TotalCount)

	return &connection, nil
}

func mapFHIREncounterToEncounterDTO(fhirEncounters []domain.FHIREncounter) []*dto.Encounter {
	var encounters []*dto.Encounter

	for _, fhirEncounter := range fhirEncounters {
		encounter := &dto.Encounter{
			ID:              *fhirEncounter.ID,
			Status:          dto.EncounterStatusEnum(fhirEncounter.Status),
			Class:           dto.EncounterClass(fhirEncounter.Class.Display),
			EpisodeOfCareID: *fhirEncounter.EpisodeOfCare[0].ID,
		}

		if fhirEncounter.Subject.ID != nil {
			encounter.PatientID = *fhirEncounter.Subject.ID
		}

		encounters = append(encounters, encounter)
	}

	return encounters
}
