package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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

// GetPatientTemperatureEntries returns all the temperature entries for a patient, they are automatically sorted in chronological order
func (c *UseCasesClinicalImpl) GetPatientTemperatureEntries(ctx context.Context, patientID string) ([]*dto.Observation, error) {
	return c.GetPatientObservations(ctx, patientID, common.TemperatureCIELTerminologyCode)
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

// RecordPulseRate records a patient's pulse rate
func (c *UseCasesClinicalImpl) RecordPulseRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	pulseRateObservation, err := c.RecordObservation(ctx, input, common.PulseCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return pulseRateObservation, nil
}

// RecordBloodPressure records a patient's blood pressure
func (c *UseCasesClinicalImpl) RecordBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	bloodPressureObservation, err := c.RecordObservation(ctx, input, common.BloodPressureCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return bloodPressureObservation, nil
}

// GetPatientBloodPressureEntries retrieves all blood pressure entries for a patient
func (c *UseCasesClinicalImpl) GetPatientBloodPressureEntries(ctx context.Context, patientID string) ([]*dto.Observation, error) {
	return c.GetPatientObservations(ctx, patientID, common.BloodPressureCIELTerminologyCode)
}

// RecordBMI records a patient's BMI
func (c *UseCasesClinicalImpl) RecordBMI(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	bmiObservation, err := c.RecordObservation(ctx, input, common.BMICIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return bmiObservation, nil
}

// RecordObservation is an extracted function that takes any observation input and saves it to FHIR.
// A concept ID is also passed so that we can get the concept code of the passed observation
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

// GetPatientObservations is a helper function used to fetch patient's observations based off the passed CIEL
// terminology code. The observations will be sorted in a chronological error
func (c *UseCasesClinicalImpl) GetPatientObservations(ctx context.Context, patientID string, observationCode string) ([]*dto.Observation, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
	}

	_, err = c.infrastructure.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", patientID)
	observations := []*dto.Observation{}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	patientObs, err := c.infrastructure.FHIR.SearchPatientObservations(ctx, patientReference, observationCode, *identifiers)
	if err != nil {
		return nil, err
	}

	for _, obs := range patientObs {
		observations = append(observations, mapFHIRObservationToObservationDTO(obs))
	}

	return observations, nil
}
