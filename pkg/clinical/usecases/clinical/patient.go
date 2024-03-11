package clinical

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/scalarutils"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
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

// GetMedicalData returns a limited subset of specific medical data that for a specific patient
// These include: Allergies, Viral Load, Body Mass Index, Weight, CD4 Count using their respective OCL CIEL Terminology
// For each category the latest three records are fetched
func (c *UseCasesClinicalImpl) GetMedicalData(ctx context.Context, patientID string) (*dto.MedicalData, error) {
	data := &dto.MedicalData{}

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

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	for _, field := range fields {
		switch field {
		case "Regimen":
			conn, err := c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				if edge.Node.ID == nil {
					continue
				}

				if edge.Node.Status == nil {
					continue
				}

				if edge.Node.MedicationCodeableConcept.Coding == nil {
					continue
				}

				if len(edge.Node.MedicationCodeableConcept.Coding) < 1 {
					continue
				}

				if edge.Node.Subject.ID == nil {
					continue
				}

				data.Regimen = append(data.Regimen, mapFHIRMedicationStatementToMedicationStatementDTO(edge.Node))
			}
		case "AllergyIntolerance":
			conn, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Allergies {
				if edge.ID == nil {
					continue
				}

				if edge.Patient == nil {
					continue
				}

				if edge.Patient.ID == nil {
					continue
				}

				if edge.Code == nil {
					continue
				}

				if edge.Code.Coding == nil {
					continue
				}

				if len(edge.Code.Coding) < 1 {
					continue
				}

				if edge.Code.Coding[0].ID == nil {
					continue
				}

				data.Allergies = append(data.Allergies, mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(edge))
			}

		case "Weight":
			filterParams["code"] = common.WeightCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, observation := range conn.Observations {
				if !hasNilInObservation(observation) {
					data.Weight = append(data.Weight, mapFHIRObservationToObservationDTO(observation))
				}
			}

		case "BMI":
			filterParams["code"] = common.BMICIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, observation := range conn.Observations {
				if !hasNilInObservation(observation) {
					data.BMI = append(data.BMI, mapFHIRObservationToObservationDTO(observation))
				}
			}

		case "ViralLoad":
			filterParams["code"] = common.ViralLoadCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, observation := range conn.Observations {
				if !hasNilInObservation(observation) {
					data.ViralLoad = append(data.ViralLoad, mapFHIRObservationToObservationDTO(observation))
				}
			}

		case "CD4Count":
			filterParams["code"] = common.CD4CountCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers, dto.Pagination{Skip: true})
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, observation := range conn.Observations {
				if !hasNilInObservation(observation) {
					data.CD4Count = append(data.CD4Count, mapFHIRObservationToObservationDTO(observation))
				}
			}
		}
	}

	return data, nil
}

func hasNilInObservation(observation domain.FHIRObservation) bool {
	if observation.ID == nil {
		return true
	}

	if observation.Code == nil {
		return true
	}

	if observation.Code.Coding == nil {
		return true
	}

	if len(observation.Code.Coding) < 1 {
		return true
	}

	if observation.Subject == nil {
		return true
	}

	if observation.Subject.ID == nil {
		return true
	}

	return false
}

func mapFHIRMedicationStatementToMedicationStatementDTO(fhirAllergyIntolerance *domain.FHIRMedicationStatement) *dto.MedicationStatement {
	return &dto.MedicationStatement{
		ID:     *fhirAllergyIntolerance.ID,
		Status: dto.MedicationStatementStatusEnum(*fhirAllergyIntolerance.Status),
		Medication: dto.Medication{
			Name: fhirAllergyIntolerance.MedicationCodeableConcept.Coding[0].Display,
			Code: fhirAllergyIntoleranceClinicalStatusURL,
		},
		PatientID: *fhirAllergyIntolerance.Subject.ID,
	}
}

func mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(fhirAllergyIntolerance domain.FHIRAllergyIntolerance) *dto.Allergy {
	allergyIntolerance := &dto.Allergy{
		ID:        *fhirAllergyIntolerance.ID,
		PatientID: *fhirAllergyIntolerance.Patient.ID,
		Code:      string(*fhirAllergyIntolerance.Code.Coding[0].Code),
		Name:      string(fhirAllergyIntolerance.Code.Coding[0].Display),
		System:    string(fhirAllergyIntolerance.Code.Text),
	}

	if fhirAllergyIntolerance.Encounter != nil && fhirAllergyIntolerance.Encounter.ID != nil {
		allergyIntolerance.EncounterID = *fhirAllergyIntolerance.Encounter.ID
	}

	if fhirAllergyIntolerance.OnsetPeriod != nil {
		allergyIntolerance.OnsetDateTime = fhirAllergyIntolerance.OnsetPeriod.Start
	}

	if len(fhirAllergyIntolerance.Reaction) > 0 {
		reaction := fhirAllergyIntolerance.Reaction[0]
		if reaction.Severity != nil {
			allergyIntolerance.Reaction.Severity = dto.AllergyIntoleranceReactionSeverityEnum(*reaction.Severity)
		}

		if len(reaction.Manifestation) > 0 {
			manifestation := reaction.Manifestation[0]
			if len(manifestation.Coding) > 0 {
				coding := manifestation.Coding[0]
				if coding.System != nil {
					allergyIntolerance.Reaction.System = string(*coding.System)
				}

				allergyIntolerance.Reaction.Code = string(*coding.Code)
				allergyIntolerance.Reaction.Name = string(coding.Display)
			}
		}
	}

	return allergyIntolerance
}

func mapFHIRObservationToObservationDTO(fhirObservation domain.FHIRObservation) *dto.Observation {
	var value string

	if fhirObservation.ValueQuantity != nil {
		value = fmt.Sprintf("%v %v", fhirObservation.ValueQuantity.Value, fhirObservation.ValueQuantity.Unit)
	}

	if fhirObservation.ValueCodeableConcept != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueCodeableConcept)
	}

	if fhirObservation.ValueString != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueString)
	}

	if fhirObservation.ValueBoolean != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueBoolean)
	}

	if fhirObservation.ValueInteger != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueInteger)
	}

	if fhirObservation.ValueRange != nil {
		value = fmt.Sprintf("%v %v - %v %v", fhirObservation.ValueRange.High.Value, fhirObservation.ValueRange.High.Unit, fhirObservation.ValueRange.Low.Value, fhirObservation.ValueRange.Low.Unit)
	}

	if fhirObservation.ValueRatio != nil {
		value = fmt.Sprintf("%v %v : %v %v", fhirObservation.ValueRatio.Numerator.Value, fhirObservation.ValueRatio.Numerator.Unit, fhirObservation.ValueRatio.Denominator.Value, fhirObservation.ValueRatio.Denominator.Unit)
	}

	if fhirObservation.ValueSampledData != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueSampledData.ID)
	}

	if fhirObservation.ValueTime != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueTime)
	}

	if fhirObservation.ValueDateTime != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueDateTime)
	}

	if fhirObservation.ValuePeriod != nil {
		value = fmt.Sprintf("%v - %v", fhirObservation.ValuePeriod.Start, fhirObservation.ValuePeriod.End)
	}

	obs := &dto.Observation{
		ID:           *fhirObservation.ID,
		Status:       dto.ObservationStatus(*fhirObservation.Status),
		Name:         fhirObservation.Code.Coding[0].Display,
		Value:        value,
		PatientID:    *fhirObservation.Subject.ID,
		TimeRecorded: string(*fhirObservation.EffectiveInstant),
	}

	if fhirObservation.Encounter != nil && fhirObservation.Encounter.ID != nil {
		obs.EncounterID = *fhirObservation.Encounter.ID
	}

	if fhirObservation.Note != nil {
		obs.Note = string(*fhirObservation.Note[0].Text)
	}

	for _, interpretation := range fhirObservation.Interpretation {
		obs.Interpretation = append(obs.Interpretation, interpretation.Text)
	}

	return obs
}

// CreatePatient creates a new patient
func (c *UseCasesClinicalImpl) CreatePatient(ctx context.Context, input dto.PatientInput) (*dto.Patient, error) {
	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	facility, err := c.infrastructure.FHIR.GetFHIROrganization(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	nameInput := &domain.NameInput{
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		OtherNames: input.OtherNames,
	}

	phoneNumbers := []*domain.PhoneNumberInput{}

	for _, contact := range input.Contacts {
		if contact.Type == dto.ContactTypePhoneNumber {
			number := &domain.PhoneNumberInput{
				Msisdn: contact.Value,
			}
			phoneNumbers = append(phoneNumbers, number)
		}
	}

	documents := []*domain.IdentificationDocument{}

	for _, identifier := range input.Identifiers {
		doc := &domain.IdentificationDocument{
			DocumentType:   domain.IDDocumentType(identifier.Type),
			DocumentNumber: identifier.Value,
		}

		documents = append(documents, doc)
	}

	registrationInput := domain.SimplePatientRegistrationInput{
		Names:                   []*domain.NameInput{nameInput},
		BirthDate:               input.BirthDate,
		PhoneNumbers:            phoneNumbers,
		Gender:                  string(input.Gender),
		Active:                  true,
		IdentificationDocuments: documents,
	}

	patientInput, err := c.SimplePatientRegistrationInputToPatientInput(ctx, registrationInput)
	if err != nil {
		return nil, err
	}

	orgRef := fmt.Sprintf("Organization/%s", *facility.Resource.ID)
	orgType := scalarutils.URI("Organization")

	patientInput.ManagingOrganization = &domain.FHIRReferenceInput{
		ID:        facility.Resource.ID,
		Reference: &orgRef,
		Display:   *facility.Resource.Name,
		Type:      &orgType,
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	patientInput.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	patient, err := c.infrastructure.FHIR.CreateFHIRPatient(ctx, *patientInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRPatientToPatientDTO(patient.PatientRecord), nil
}

func (c *UseCasesClinicalImpl) PatchPatient(ctx context.Context, id string, input dto.PatientInput) (*dto.Patient, error) {
	if id == "" {
		return nil, fmt.Errorf("a patient ID is required")
	}

	registrationInput := domain.SimplePatientRegistrationInput{
		Gender:    string(input.Gender),
		BirthDate: input.BirthDate,
	}

	if input.FirstName != "" {
		nameInput := &domain.NameInput{
			FirstName:  input.FirstName,
			LastName:   input.LastName,
			OtherNames: input.OtherNames,
		}
		registrationInput.Names = append(registrationInput.Names, nameInput)
	}

	for _, contact := range input.Contacts {
		if contact.Type == dto.ContactTypePhoneNumber {
			number := &domain.PhoneNumberInput{
				Msisdn: contact.Value,
			}
			registrationInput.PhoneNumbers = append(
				registrationInput.PhoneNumbers,
				number,
			)
		}
	}

	for _, identifier := range input.Identifiers {
		doc := &domain.IdentificationDocument{
			DocumentType:   domain.IDDocumentType(identifier.Type),
			DocumentNumber: identifier.Value,
		}

		registrationInput.IdentificationDocuments = append(
			registrationInput.IdentificationDocuments,
			doc,
		)
	}

	patientInput, err := c.SimplePatientRegistrationInputToPatientInput(ctx, registrationInput)
	if err != nil {
		return nil, err
	}

	patient, err := c.infrastructure.FHIR.PatchFHIRPatient(ctx, id, *patientInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRPatientToPatientDTO(patient), nil
}

func (c *UseCasesClinicalImpl) DeletePatient(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("a patient ID is required")
	}

	ok, err := c.infrastructure.FHIR.DeleteFHIRPatient(ctx, id)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func mapFHIRPatientToPatientDTO(patient *domain.FHIRPatient) *dto.Patient {
	numbers := []string{}

	for _, phone := range patient.Telecom {
		if *phone.System == domain.ContactPointSystemEnumPhone {
			numbers = append(numbers, *phone.Value)
		}
	}

	return &dto.Patient{
		ID:          *patient.ID,
		Active:      *patient.Active,
		Name:        patient.Name[0].Text,
		PhoneNumber: numbers,
		Gender:      dto.Gender(patient.Gender.String()),
		BirthDate:   *patient.BirthDate,
	}
}

// PatientMedicationHistory is used to retrieve all the patient clinical information. From Observations, Conditions,
// MedicationStatement etc. relevant to referral from one facility to another.
func (c *UseCasesClinicalImpl) PatientMedicationHistory(ctx context.Context, patientID string, pagination *dto.Pagination) (*dto.PatientMedicationHistoryOutput, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
	}

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", *patient.Resource.ID)

	searchParams := map[string]interface{}{
		"patient": patientReference,
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	observationsChan := make(chan *domain.PagedFHIRObservations)
	conditionsChan := make(chan *domain.PagedFHIRCondition)
	medicationStatementsChan := make(chan *domain.FHIRMedicationStatementRelayConnection)
	errChan := make(chan error, 3)

	go func() {
		result, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, searchParams, *identifiers, *pagination)
		if err != nil {
			errChan <- err
			return
		}
		observationsChan <- result
	}()

	go func() {
		result, err := c.infrastructure.FHIR.SearchFHIRCondition(ctx, searchParams, *identifiers, *pagination)
		if err != nil {
			errChan <- err
			return
		}
		conditionsChan <- result
	}()

	go func() {
		result, err := c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, searchParams, *identifiers, *pagination)
		if err != nil {
			errChan <- err
			return
		}
		medicationStatementsChan <- result
	}()

	var observations *domain.PagedFHIRObservations

	var conditions *domain.PagedFHIRCondition

	var medicationStatements *domain.FHIRMedicationStatementRelayConnection

	// ensure whichever goroutine finishes first is processed, whether it returns a result or an error.
	for i := 0; i < 3; i++ {
		select {
		case obs := <-observationsChan:
			observations = obs
		case conds := <-conditionsChan:
			conditions = conds
		case meds := <-medicationStatementsChan:
			medicationStatements = meds
		case err := <-errChan:
			return nil, err
		}
	}

	var conditionsResult []*dto.Condition
	for _, item := range conditions.Conditions {
		conditionsResult = append(conditionsResult, mapFHIRConditionToConditionDTO(item))
	}

	var observationList []*dto.Observation
	for _, item := range observations.Observations {
		observationList = append(observationList, mapFHIRObservationToObservationDTO(item))
	}

	var medicationStatement []*dto.MedicationStatement
	for _, item := range medicationStatements.Edges {
		medicationStatement = append(medicationStatement, mapFHIRMedicationStatementToMedicationStatementDTO(item.Node))
	}

	output := &dto.PatientMedicationHistoryOutput{
		Conditions:   conditionsResult,
		Observations: observationList,
		Patient:      *mapFHIRPatientToPatientDTO(patient.Resource),
		Medications:  medicationStatement,
	}

	return output, nil
}
