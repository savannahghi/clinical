package clinical

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

var (
	fhirAllergyIntoleranceClinicalStatusURL     = "http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical"
	fhirAllergyIntoleranceVerificationStatusURL = "http://terminology.hl7.org/CodeSystem/allergyintolerance-verification"
	unknownConceptID                            = "1067"
)

// CreateFHIRPatient creates a patient on FHIR store
func (c *UseCasesClinicalImpl) CreatePubsubPatient(ctx context.Context, payload dto.PatientPubSubMessage) error {
	year, month, day := payload.DateOfBirth.Date()
	patientName := strings.Split(payload.Name, " ")
	registrationInput := domain.SimplePatientRegistrationInput{
		Names: []*domain.NameInput{{FirstName: patientName[0], LastName: patientName[1]}},
		BirthDate: scalarutils.Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		},
		PhoneNumbers: []*domain.PhoneNumberInput{{Msisdn: payload.PhoneNumber, CommunicationOptIn: false}},
		Gender:       string(payload.Gender),
		Active:       payload.Active,
	}

	patientInput, err := c.SimplePatientRegistrationInputToPatientInput(ctx, registrationInput)
	if err != nil {
		return err
	}

	newID := uuid.New().String()
	patientInput.ID = &newID

	patientInput.Identifier = append(patientInput.Identifier, common.DefaultIdentifier())

	clientSystem := scalarutils.URI("mycarehub.client.id")
	userSelected := false

	clientIdentifier := &domain.FHIRIdentifierInput{
		Use:   domain.IdentifierUseEnumOfficial,
		Value: payload.ClientID,
		Type: domain.FHIRCodeableConceptInput{
			Text: payload.ClientID,
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &clientSystem,
					Code:         scalarutils.Code(payload.ClientID),
					Display:      payload.ClientID,
					UserSelected: &userSelected,
				},
			},
		},
		System: &clientSystem,
		Period: common.DefaultPeriodInput(),
	}

	patientInput.Identifier = append(patientInput.Identifier, clientIdentifier)

	userSystem := scalarutils.URI("mycarehub.user.id")

	userIdentifier := &domain.FHIRIdentifierInput{
		Use:   domain.IdentifierUseEnumOfficial,
		Value: payload.UserID,
		Type: domain.FHIRCodeableConceptInput{
			Text: payload.UserID,
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &userSystem,
					Code:         scalarutils.Code(payload.UserID),
					Display:      payload.UserID,
					UserSelected: &userSelected,
				},
			},
		},
		System: &userSystem,
		Period: common.DefaultPeriodInput(),
	}

	patientInput.Identifier = append(patientInput.Identifier, userIdentifier)

	tags, err := c.CreateTenantMetaTags(ctx, payload.OrganizationID, payload.FacilityID)
	if err != nil {
		return err
	}

	patientInput.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	patient, err := c.infrastructure.FHIR.CreateFHIRPatient(ctx, *patientInput)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return err
	}

	err = c.infrastructure.Pubsub.NotifyPatientFHIRIDUpdate(ctx, dto.UpdatePatientFHIRID{
		FhirID:   *patient.PatientRecord.ID,
		ClientID: payload.ClientID,
	})
	if err != nil {
		return err
	}

	return nil
}

// CreatePubsubOrganization creates a FHIR organisation resource
func (c *UseCasesClinicalImpl) CreatePubsubOrganization(ctx context.Context, data dto.FacilityPubSubMessage) error {
	use := domain.ContactPointUseEnumWork
	rank := int64(1)
	phoneSystem := domain.ContactPointSystemEnumPhone
	input := domain.FHIROrganizationInput{
		ID:     data.ID,
		Active: &data.Active,
		Name:   &data.Name,
		Telecom: []*domain.FHIRContactPointInput{
			{
				System: &phoneSystem,
				Value:  &data.Phone,
				Use:    &use,
				Rank:   &rank,
				Period: common.DefaultPeriodInput(),
			},
		},
	}

	response, err := c.infrastructure.FHIR.CreateFHIROrganization(ctx, input)
	if err != nil {
		return err
	}

	err = c.infrastructure.Pubsub.NotifyFacilityFHIRIDUpdate(ctx, dto.UpdateFacilityFHIRID{
		FacilityID: *data.ID,
		FhirID:     *response.Resource.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

// CreatePubsubVitals creates FHIR observation vitals.
func (c *UseCasesClinicalImpl) CreatePubsubVitals(ctx context.Context, data dto.VitalSignPubSubMessage) error {
	input, err := c.ComposeVitalsInput(ctx, data)
	if err != nil {
		return err
	}

	tags, err := c.CreateTenantMetaTags(ctx, data.OrganizationID, data.FacilityID)
	if err != nil {
		return err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	_, err = c.infrastructure.FHIR.CreateFHIRObservation(ctx, *input)
	if err != nil {
		return err
	}

	return nil
}

// CreatePubsubAllergyIntolerance creates FHIR allergy intolerance
func (c *UseCasesClinicalImpl) CreatePubsubAllergyIntolerance(ctx context.Context, data dto.PatientAllergyPubSubMessage) error {
	input, err := c.ComposeAllergyIntoleranceInput(ctx, data)
	if err != nil {
		return err
	}

	tags, err := c.CreateTenantMetaTags(ctx, data.OrganizationID, data.FacilityID)
	if err != nil {
		return err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	_, err = c.infrastructure.FHIR.CreateFHIRAllergyIntolerance(ctx, *input)
	if err != nil {
		return err
	}

	return nil
}

// CreatePubsubTestResult creates a test result as an observation
func (c *UseCasesClinicalImpl) CreatePubsubTestResult(ctx context.Context, data dto.PatientTestResultPubSubMessage) error {
	input, err := c.ComposeTestResultInput(ctx, data)
	if err != nil {
		return err
	}

	tags, err := c.CreateTenantMetaTags(ctx, data.OrganizationID, data.FacilityID)
	if err != nil {
		return err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	_, err = c.infrastructure.FHIR.CreateFHIRObservation(ctx, *input)
	if err != nil {
		return err
	}

	return nil
}

// CreatePubsubMedicationStatement creates a FHIR medication statement
func (c *UseCasesClinicalImpl) CreatePubsubMedicationStatement(ctx context.Context, data dto.MedicationPubSubMessage) error {
	input, err := c.ComposeMedicationStatementInput(ctx, data)
	if err != nil {
		return err
	}

	tags, err := c.CreateTenantMetaTags(ctx, data.OrganizationID, data.FacilityID)
	if err != nil {
		return err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	_, err = c.infrastructure.FHIR.CreateFHIRMedicationStatement(ctx, *input)
	if err != nil {
		return err
	}

	return nil
}

// CreatePubsubTenant registers a tenant in this service
func (c *UseCasesClinicalImpl) CreatePubsubTenant(ctx context.Context, data dto.OrganizationInput) error {
	if data.Identifiers[0].Type != dto.OrganizationIdentifierType("MCHProgram") {
		return fmt.Errorf("invalid identifier type %v", data.Identifiers[0].Type)
	}

	organization, err := c.RegisterTenant(ctx, data)
	if err != nil {
		return err
	}

	err = c.infrastructure.Pubsub.NotifyProgramFHIRIDUpdate(ctx, dto.UpdateProgramFHIRID{
		ProgramID:    data.Identifiers[0].Value,
		FHIRTenantID: organization.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
