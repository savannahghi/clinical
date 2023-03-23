package clinical

import (
	"context"
	"fmt"

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
			conn, err := c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, filterParams, *identifiers)
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
			conn, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, filterParams, *identifiers)
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

				if edge.Node.Patient == nil {
					continue
				}

				if edge.Node.Patient.ID == nil {
					continue
				}

				if edge.Node.Encounter == nil {
					continue
				}

				if edge.Node.Encounter.ID == nil {
					continue
				}

				if edge.Node.Code == nil {
					continue
				}

				if edge.Node.Code.Coding == nil {
					continue
				}

				if len(edge.Node.Code.Coding) < 1 {
					continue
				}

				data.Allergies = append(data.Allergies, mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(edge.Node))
			}

		case "Weight":
			filterParams["code"] = common.WeightCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if !hasNilInObservation(edge) {
					data.Weight = append(data.Weight, mapFHIRObservationToObservationDTO(edge.Node))
				}
			}

		case "BMI":
			filterParams["code"] = common.BMICIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if !hasNilInObservation(edge) {
					data.BMI = append(data.BMI, mapFHIRObservationToObservationDTO(edge.Node))
				}
			}

		case "ViralLoad":
			filterParams["code"] = common.ViralLoadCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if !hasNilInObservation(edge) {
					data.ViralLoad = append(data.ViralLoad, mapFHIRObservationToObservationDTO(edge.Node))
				}
			}

		case "CD4Count":
			filterParams["code"] = common.CD4CountCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if !hasNilInObservation(edge) {
					data.CD4Count = append(data.CD4Count, mapFHIRObservationToObservationDTO(edge.Node))
				}
			}
		}
	}

	return data, nil
}

func hasNilInObservation(observation *domain.FHIRObservationRelayEdge) bool {
	if observation.Node == nil {
		return true
	}

	if observation.Node.Code.Coding == nil {
		return true
	}

	if len(observation.Node.Code.Coding) < 1 {
		return true
	}

	if observation.Node.Subject == nil {
		return true
	}

	if observation.Node.Subject.ID == nil {
		return true
	}

	if observation.Node.Encounter == nil {
		return true
	}

	if observation.Node.Encounter.ID == nil {
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
			Code: string(fhirAllergyIntolerance.MedicationCodeableConcept.Coding[0].Code),
		},
		PatientID: *fhirAllergyIntolerance.Subject.ID,
	}
}

func mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(fhirAllergyIntolerance *domain.FHIRAllergyIntolerance) *dto.Allergy {
	allergyIntolerance := &dto.Allergy{
		ID:          *fhirAllergyIntolerance.ID,
		PatientID:   *fhirAllergyIntolerance.Patient.ID,
		EncounterID: *fhirAllergyIntolerance.Encounter.ID,
		Code:        string(fhirAllergyIntolerance.Code.Coding[0].Code),
		System:      string(*fhirAllergyIntolerance.Code.Coding[0].System),
	}

	if fhirAllergyIntolerance.OnsetPeriod != nil {
		allergyIntolerance.OnsetDateTime = fhirAllergyIntolerance.OnsetPeriod.Start
	}

	if fhirAllergyIntolerance.Reaction != nil {
		allergyIntolerance.Reaction.Severity = dto.AllergyIntoleranceReactionSeverityEnum(*fhirAllergyIntolerance.Reaction[0].Severity)
	}

	return allergyIntolerance
}

func mapFHIRObservationToObservationDTO(fhirObservation *domain.FHIRObservation) *dto.Observation {
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

	return &dto.Observation{
		ID:          *fhirObservation.ID,
		Status:      dto.ObservationStatus(*fhirObservation.Status),
		Name:        fhirObservation.Code.Coding[0].Display,
		Value:       value,
		PatientID:   *fhirObservation.Subject.ID,
		EncounterID: *fhirObservation.Encounter.ID,
	}
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

	exists, err := c.CheckPatientExistenceUsingPhoneNumber(ctx, registrationInput)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to check patient existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("patient with phone number already exists")
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

	patientInput.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	patient, err := c.infrastructure.FHIR.CreateFHIRPatient(ctx, *patientInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRPatientToPatientDTO(patient.PatientRecord), nil
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
