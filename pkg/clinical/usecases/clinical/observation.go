package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// RecordTemperature is used to record a patient's temperature and saves it as a FHIR observation
func (c *UseCasesClinicalImpl) RecordTemperature(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	temperatureObservation, err := c.RecordObservation(ctx, input, common.TemperatureCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return temperatureObservation, nil
}

// RecordHeight records a patient's height and saves it to fhir
func (c *UseCasesClinicalImpl) RecordHeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	heightObservation, err := c.RecordObservation(ctx, input, common.HeightCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return heightObservation, nil
}

// RecordWeight records a patient's weight
func (c *UseCasesClinicalImpl) RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	weightObservation, err := c.RecordObservation(ctx, input, common.WeightCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return weightObservation, nil
}

// RecordRespiratoryRate records a patient's respiratory rate
func (c *UseCasesClinicalImpl) RecordRespiratoryRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	respiratoryRateObservation, err := c.RecordObservation(ctx, input, common.RespiratoryRateCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return respiratoryRateObservation, nil
}

func (c *UseCasesClinicalImpl) RecordObservation(ctx context.Context, input dto.ObservationInput, vitalSignConceptID string) (*dto.Observation, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot record an observation in a finished encounter")
	}

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)

	vitalsConcept, err := c.getCIELConcept(ctx, vitalSignConceptID)
	if err != nil {
		return nil, err
	}

	system := "http://terminology.hl7.org/CodeSystem/observation-category"
	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))
	observation := domain.FHIRObservationInput{
		Status: (*domain.ObservationStatusEnum)(&input.Status),
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&system),
						Code:    "vital-signs",
						Display: "Vital Signs",
					},
				},
				Text: "Vital Signs",
			},
		},
		EffectiveInstant: &instant,
		Code: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&vitalsConcept.URL),
					Code:    scalarutils.Code(vitalsConcept.ID),
					Display: vitalsConcept.DisplayName,
				},
			},
			Text: vitalsConcept.DisplayName,
		},
		ValueString: &input.Value,
		Subject: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.Subject.ID,
			Reference: &patientReference,
			Display:   encounter.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID: encounter.Resource.ID,
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	observation.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	fhirObservation, err := c.infrastructure.FHIR.CreateFHIRObservation(ctx, observation)
	if err != nil {
		return nil, err
	}

	return mapFHIRObservationToObservationDTO(fhirObservation.Resource), nil
}
