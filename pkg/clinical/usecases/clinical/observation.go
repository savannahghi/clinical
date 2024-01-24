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

// RecordMuac is used to record a patient's Muac
func (c *UseCasesClinicalImpl) RecordMuac(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	muacObservation, err := c.RecordObservation(ctx, input, common.MuacCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return muacObservation, nil
}

// RecordOxygenSaturation is used to record a patient's oxygen saturation
func (c *UseCasesClinicalImpl) RecordOxygenSaturation(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	oxygenSaturationObservation, err := c.RecordObservation(ctx, input, common.OxygenSaturationCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return oxygenSaturationObservation, nil
}

// GetPatientTemperatureEntries returns all the temperature entries for a patient, they are automatically sorted in chronological order
func (c *UseCasesClinicalImpl) GetPatientTemperatureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.TemperatureCIELTerminologyCode, pagination)
}

// RecordHeight records a patient's height and saves it to fhir
func (c *UseCasesClinicalImpl) RecordHeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	heightObservation, err := c.RecordObservation(ctx, input, common.HeightCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return heightObservation, nil
}

// GetPatientHeightEntries gets the height records of a patient
func (c *UseCasesClinicalImpl) GetPatientHeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.HeightCIELTerminologyCode, pagination)
}

// PatchPatientHeight patches the height record of a patient
func (c *UseCasesClinicalImpl) PatchPatientHeight(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientWeight patches the weight record of a patient
func (c *UseCasesClinicalImpl) PatchPatientWeight(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientBMI patches the BMI record of a patient
func (c *UseCasesClinicalImpl) PatchPatientBMI(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientTemperature patches the temperature record of a patient
func (c *UseCasesClinicalImpl) PatchPatientTemperature(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientDiastolicBloodPressure patches the diastolic blood pressure record of a patient
func (c *UseCasesClinicalImpl) PatchPatientDiastolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientSystolicBloodPressure patches the Systolic blood pressure record of a patient
func (c *UseCasesClinicalImpl) PatchPatientSystolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientRespiratoryRate patches the respiration rate record of a patient
func (c *UseCasesClinicalImpl) PatchPatientRespiratoryRate(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// RecordWeight records a patient's weight
func (c *UseCasesClinicalImpl) RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	weightObservation, err := c.RecordObservation(ctx, input, common.WeightCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return weightObservation, nil
}

// RecordViralLoad records the patient viral load
func (c *UseCasesClinicalImpl) RecordViralLoad(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	viralLoadObservation, err := c.RecordObservation(ctx, input, common.ViralLoadCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return viralLoadObservation, nil
}

// GetPatientWeightEntries gets the weight records of a patient
func (c *UseCasesClinicalImpl) GetPatientWeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.WeightCIELTerminologyCode, pagination)
}

// GetPatientMuacEntries gets the patient's muac
func (c *UseCasesClinicalImpl) GetPatientMuacEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.MuacCIELTerminologyCode, pagination)
}

// GetPatientOxygenSaturationEntries gets the patient's oxygen saturation
func (c *UseCasesClinicalImpl) GetPatientOxygenSaturationEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.OxygenSaturationCIELTerminologyCode, pagination)
}

// GetPatientViralLoad gets the patient's viral load
func (c *UseCasesClinicalImpl) GetPatientViralLoad(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.ViralLoadCIELTerminologyCode, pagination)
}

// RecordRespiratoryRate records a patient's respiratory rate
func (c *UseCasesClinicalImpl) RecordRespiratoryRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	respiratoryRateObservation, err := c.RecordObservation(ctx, input, common.RespiratoryRateCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return respiratoryRateObservation, nil
}

// GetPatientRespiratoryRateEntries gets a patient's respiratory rate entries
func (c *UseCasesClinicalImpl) GetPatientRespiratoryRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.RespiratoryRateCIELTerminologyCode, pagination)
}

// RecordPulseRate records a patient's pulse rate
func (c *UseCasesClinicalImpl) RecordPulseRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	pulseRateObservation, err := c.RecordObservation(ctx, input, common.PulseCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return pulseRateObservation, nil
}

// GetPatientPulseRateEntries gets the pulse rate records of a patient
func (c *UseCasesClinicalImpl) GetPatientPulseRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.PulseCIELTerminologyCode, pagination)
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
func (c *UseCasesClinicalImpl) GetPatientBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.BloodPressureCIELTerminologyCode, pagination)
}

// RecordBMI records a patient's BMI
func (c *UseCasesClinicalImpl) RecordBMI(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	bmiObservation, err := c.RecordObservation(ctx, input, common.BMICIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return bmiObservation, nil
}

// GetPatientBMIEntries retrieves all BMI entries for a patient
func (c *UseCasesClinicalImpl) GetPatientBMIEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.BMICIELTerminologyCode, pagination)
}

// RecordBloodSugar records a patient's blood sugar level (Serum glucose)
func (c *UseCasesClinicalImpl) RecordBloodSugar(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	return c.RecordObservation(ctx, input, common.BloodSugarCIELTerminologyCode)
}

// GetPatientBloodSugarEntries retrieves all blood sugar entries for a patient
func (c *UseCasesClinicalImpl) GetPatientBloodSugarEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.BloodSugarCIELTerminologyCode, pagination)
}

// RecordLastMenstrualPeriod records last menstrual period
func (c *UseCasesClinicalImpl) RecordLastMenstrualPeriod(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	return c.RecordObservation(ctx, input, common.LastMenstrualPeriodCIELTerminologyCode)
}

// GetPatientLastMenstrualPeriodEntries retrieves all blood sugar entries for a patient
func (c *UseCasesClinicalImpl) GetPatientLastMenstrualPeriodEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.LastMenstrualPeriodCIELTerminologyCode, pagination)
}

// RecordDiastolicBloodPressure records diastolic blood pressure
func (c *UseCasesClinicalImpl) RecordDiastolicBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	return c.RecordObservation(ctx, input, common.DiastolicBloodPressureCIELTerminologyCode)
}

// RecordColposcopy records colposcopy findings
func (c *UseCasesClinicalImpl) RecordColposcopy(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	return c.RecordObservation(ctx, input, common.ColposcopyCIELTerminologyCode)
}

// GetPatientDiastolicBloodPressureEntries retrieves all diastolic blood pressure entries for a patient
func (c *UseCasesClinicalImpl) GetPatientDiastolicBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.DiastolicBloodPressureCIELTerminologyCode, pagination)
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

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	vitalsConcept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, vitalSignConceptID)
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
		Code: &domain.FHIRCodeableConceptInput{
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
			ID:        encounter.Resource.ID,
			Reference: &encounterReference,
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	observation.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	fhirObservation, err := c.infrastructure.FHIR.CreateFHIRObservation(ctx, observation)
	if err != nil {
		return nil, err
	}

	return mapFHIRObservationToObservationDTO(*fhirObservation), nil
}

// GetPatientObservations is a helper function used to fetch patient's observations based on the passed CIEL
// terminology code. The observations will be sorted in a chronological order
func (c *UseCasesClinicalImpl) GetPatientObservations(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, observationCode string, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
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

	searchParams := map[string]interface{}{
		"patient": patientReference,
		"code":    observationCode,
	}

	if encounterID != nil {
		encounterReference := fmt.Sprintf("Encounter/%s", *encounterID)
		searchParams["encounter"] = encounterReference
	}

	if date != nil {
		searchParams["date"] = date.AsTime().Format(dateFormatStr)
	}

	patientObs, err := c.infrastructure.FHIR.SearchPatientObservations(ctx, searchParams, *identifiers, *pagination)
	if err != nil {
		return nil, err
	}

	for _, obs := range patientObs.Observations {
		if obs.Subject == nil {
			continue
		}

		if obs.Subject.ID == nil {
			continue
		}

		observations = append(observations, mapFHIRObservationToObservationDTO(obs))
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     patientObs.HasNextPage,
		EndCursor:       &patientObs.NextCursor,
		HasPreviousPage: patientObs.HasPreviousPage,
		StartCursor:     &patientObs.PreviousCursor,
	}

	connection := dto.CreateObservationConnection(observations, pageInfo, patientObs.TotalCount)

	return &connection, nil
}

// PatchPatientObservations update a patient's observation resource
func (c *UseCasesClinicalImpl) PatchPatientObservations(ctx context.Context, id string, value string) (*dto.Observation, error) {
	if value == "" {
		return nil, fmt.Errorf("observation value required")
	}

	if id == "" {
		return nil, fmt.Errorf("an observation id is required")
	}

	observation, err := c.infrastructure.FHIR.GetFHIRObservation(ctx, id)
	if err != nil {
		return nil, err
	}

	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, *observation.Resource.Encounter.ID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot patch an observation in a finished encounter")
	}

	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	observationInput := &domain.FHIRObservationInput{
		EffectiveInstant: &instant,
		ValueString:      &value,
	}

	output, err := c.infrastructure.FHIR.PatchFHIRObservation(ctx, id, *observationInput)
	if err != nil {
		return nil, err
	}

	result := mapFHIRObservationToObservationDTO(*output)

	return result, nil
}
