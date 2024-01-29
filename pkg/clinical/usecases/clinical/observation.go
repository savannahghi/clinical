package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
)

var (
	// Vitals
	VitalsCodeableConcept = []*domain.FHIRCodeableConceptInput{
		{
			Coding: []*domain.FHIRCodingInput{
				{
					Code:    "vital-signs",
					Display: "Vital Signs",
				},
			},
			Text: "Vital Signs",
		},
	}

	// Lab
	LabCodeableConcept = []*domain.FHIRCodeableConceptInput{
		{
			Coding: []*domain.FHIRCodingInput{
				{
					Code:    "laboratory",
					Display: "Laboratory",
				},
			},
			Text: "Laboratory",
		},
	}
)

// ObservationPayload models a common data class re-usable for use with any type of observation
type ObservationPayload struct {
	Input dto.ObservationInput        `json:"input,omitempty"`
	Obs   domain.FHIRObservationInput `json:"obs,omitempty"`
}

// RecordTemperature is used to record a patient's temperature and saves it as a FHIR observation
func (c *UseCasesClinicalImpl) RecordTemperature(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	temperatureObservation, err := c.recordObservation(ctx, payload, common.TemperatureCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return temperatureObservation, nil
}

// RecordMuac is used to record a patient's Muac
func (c *UseCasesClinicalImpl) RecordMuac(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	muacObservation, err := c.recordObservation(ctx, payload, common.MuacCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return muacObservation, nil
}

// RecordOxygenSaturation is used to record a patient's oxygen saturation
func (c *UseCasesClinicalImpl) RecordOxygenSaturation(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	oxygenSaturationObservation, err := c.recordObservation(ctx, payload, common.OxygenSaturationCIELTerminologyCode)
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
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	heightObservation, err := c.recordObservation(ctx, payload, common.HeightCIELTerminologyCode)
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

// PatchPatientOxygenSaturation patches the oxygen saturation record of a patient
func (c *UseCasesClinicalImpl) PatchPatientOxygenSaturation(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientPulseRate patches the pulse rate record of a patient
func (c *UseCasesClinicalImpl) PatchPatientPulseRate(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientViralLoad patches the viral load record of a patient
func (c *UseCasesClinicalImpl) PatchPatientViralLoad(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// PatchPatientMuac patches the muac record of a patient
func (c *UseCasesClinicalImpl) PatchPatientMuac(ctx context.Context, id string, value string) (*dto.Observation, error) {
	return c.PatchPatientObservations(ctx, id, value)
}

// RecordWeight records a patient's weight
func (c *UseCasesClinicalImpl) RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	weightObservation, err := c.recordObservation(ctx, payload, common.WeightCIELTerminologyCode)
	if err != nil {
		return nil, err
	}

	return weightObservation, nil
}

// RecordViralLoad records the patient viral load
func (c *UseCasesClinicalImpl) RecordViralLoad(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	viralLoadObservation, err := c.recordObservation(ctx, payload, common.ViralLoadCIELTerminologyCode)
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
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	respiratoryRateObservation, err := c.recordObservation(ctx, payload, common.RespiratoryRateCIELTerminologyCode)
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
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	pulseRateObservation, err := c.recordObservation(ctx, payload, common.PulseCIELTerminologyCode)
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
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	bloodPressureObservation, err := c.recordObservation(ctx, payload, common.BloodPressureCIELTerminologyCode)
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
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	bmiObservation, err := c.recordObservation(ctx, payload, common.BMICIELTerminologyCode)
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
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	return c.recordObservation(ctx, payload, common.BloodSugarCIELTerminologyCode)
}

// GetPatientBloodSugarEntries retrieves all blood sugar entries for a patient
func (c *UseCasesClinicalImpl) GetPatientBloodSugarEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.BloodSugarCIELTerminologyCode, pagination)
}

// RecordLastMenstrualPeriod records last menstrual period
func (c *UseCasesClinicalImpl) RecordLastMenstrualPeriod(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	return c.recordObservation(ctx, payload, common.LastMenstrualPeriodCIELTerminologyCode)
}

// GetPatientLastMenstrualPeriodEntries retrieves all blood sugar entries for a patient
func (c *UseCasesClinicalImpl) GetPatientLastMenstrualPeriodEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.LastMenstrualPeriodCIELTerminologyCode, pagination)
}

// RecordDiastolicBloodPressure records diastolic blood pressure
func (c *UseCasesClinicalImpl) RecordDiastolicBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	return c.recordObservation(ctx, payload, common.DiastolicBloodPressureCIELTerminologyCode)
}

// RecordColposcopy records colposcopy findings
func (c *UseCasesClinicalImpl) RecordColposcopy(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: VitalsCodeableConcept,
		},
	}

	return c.recordObservation(ctx, payload, common.ColposcopyCIELTerminologyCode)
}

// GetPatientDiastolicBloodPressureEntries retrieves all diastolic blood pressure entries for a patient
func (c *UseCasesClinicalImpl) GetPatientDiastolicBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error) {
	return c.GetPatientObservations(ctx, patientID, encounterID, date, common.DiastolicBloodPressureCIELTerminologyCode, pagination)
}

// recordObservation is an extracted function that takes any observation input and saves it to FHIR.
// A concept ID is also passed so that we can get the concept code of the passed observation
func (c *UseCasesClinicalImpl) recordObservation(ctx context.Context, observation ObservationPayload, conceptID string) (*dto.Observation, error) {
	err := observation.Input.Validate()
	if err != nil {
		return nil, err
	}

	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, observation.Input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot record an observation in a finished encounter")
	}

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	vitalsConcept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, conceptID)
	if err != nil {
		return nil, err
	}

	system := "http://terminology.hl7.org/CodeSystem/observation-category"
	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))
	observation.Obs.Category[0].Coding[0].System = (*scalarutils.URI)(&system)

	obs := domain.FHIRObservationInput{
		Status:           (*domain.ObservationStatusEnum)(&observation.Input.Status),
		Category:         observation.Obs.Category,
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
		ValueString: &observation.Input.Value,
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

	obs.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	fhirObservation, err := c.infrastructure.FHIR.CreateFHIRObservation(ctx, obs)
	if err != nil {
		return nil, err
	}

	return mapFHIRObservationToObservationDTO(*fhirObservation), nil
}

// RecordHPV is used to record HPV test results. We record it as observations as specified in https://build.fhir.org/ig/HL7/cqf-measures/Measure-EXM124-FHIR.html
// Check whether the gender of the patient is valid and that the patient is within the acceptable age range
func (c *UseCasesClinicalImpl) RecordHPV(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
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

	patientRecord, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, *encounter.Resource.Subject.ID)
	if err != nil {
		return nil, err
	}

	// Check if patient gender is female
	if enumutils.Gender(*patientRecord.Resource.Gender) != "female" {
		return nil, fmt.Errorf("cannot record HPV results for a male patients")
	}

	currentTime := time.Now()
	age := currentTime.Year() - patientRecord.Resource.BirthDate.AsTime().Year()

	// Check if the birthday has occurred this year
	if currentTime.YearDay() < patientRecord.Resource.BirthDate.AsTime().YearDay() {
		age--
	}

	if age < 35 || age > 65 {
		return nil, fmt.Errorf("cannot record HPV results for an age of %d. The patient has to be between 25 - 65 years", age)
	}

	payload := ObservationPayload{
		Input: input,
		Obs: domain.FHIRObservationInput{
			Category: LabCodeableConcept,
		},
	}

	return c.recordObservation(ctx, payload, common.HPVCIELTerminologyCode)
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
