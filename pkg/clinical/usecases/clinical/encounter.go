package clinical

import (
	"context"
	"fmt"

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
