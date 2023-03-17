package clinical

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
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

	facility, err := c.infrastructure.FHIR.FindOrganizationByID(ctx, facilityID)
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

func mapFHIREpisodeToEpisodeDTO(episode domain.FHIREpisodeOfCare) *dto.EpisodeOfCare {
	return &dto.EpisodeOfCare{
		ID:        *episode.ID,
		Status:    dto.EpisodeOfCareStatusEnum(string(*episode.Status)),
		PatientID: *episode.Patient.ID,
	}
}
