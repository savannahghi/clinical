package clinical

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/scalarutils"
)

// GetTenantMetaTags is a helper to create tags that are used to identify which tenant a resource belongs to
// and are saved in a resources `Meta` attribute
func (c *UseCasesClinicalImpl) GetTenantMetaTags(ctx context.Context) ([]domain.FHIRCodingInput, error) {
	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	organisation, err := c.infrastructure.FHIR.FindOrganizationByID(ctx, identifiers.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tenant organisation: %w", err)
	}

	facility, err := c.infrastructure.FHIR.FindOrganizationByID(ctx, identifiers.FacilityID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tenant organisation: %w", err)
	}

	userSelected := false

	organisationTagVersion := "1.0"
	organisationTagSystem := scalarutils.URI("http://mycarehub/tenant-identification/organisation")

	facilityTagVersion := "1.0"
	facilityTagSystem := scalarutils.URI("http://mycarehub/tenant-identification/facility")

	tags := []domain.FHIRCodingInput{
		{
			System:       &organisationTagSystem,
			Version:      &organisationTagVersion,
			Code:         scalarutils.Code(identifiers.OrganizationID),
			Display:      *organisation.Resource.Name,
			UserSelected: &userSelected,
		},
		{
			System:       &facilityTagSystem,
			Version:      &facilityTagVersion,
			Code:         scalarutils.Code(identifiers.FacilityID),
			Display:      *facility.Resource.Name,
			UserSelected: &userSelected,
		},
	}

	return tags, nil
}

// CheckPatientExistenceUsingPhoneNumber checks whether a patient with the phone number they're trying to register with exists
func (c *UseCasesClinicalImpl) CheckPatientExistenceUsingPhoneNumber(ctx context.Context, patientInput domain.SimplePatientRegistrationInput) (bool, error) {
	exists := false

	for _, phone := range patientInput.PhoneNumbers {
		phoneNumber := &phone.Msisdn

		search, err := converterandformatter.NormalizeMSISDN(*phoneNumber)
		if err != nil {
			return false, fmt.Errorf("can't normalize contact: %w", err)
		}

		patient, err := c.infrastructure.FHIR.SearchFHIRPatient(ctx, *search)
		if err != nil {
			return false, fmt.Errorf("unable to find patient by phonenumber: %s", *phoneNumber)
		}

		if len(patient.Edges) > 1 {
			exists = true
			break
		}
	}

	return exists, nil
}

// SimplePatientRegistrationInputToPatientInput transforms a patient input into
// a
func (c *UseCasesClinicalImpl) SimplePatientRegistrationInputToPatientInput(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.FHIRPatientInput, error) {
	contacts, err := c.ContactsToContactPointInput(ctx, input.PhoneNumbers, input.Emails)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't register patient with invalid contacts: %w", err)
	}

	ids, err := helpers.IDToIdentifier(input.IdentificationDocuments, input.PhoneNumbers)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't register patient with invalid identifiers: %w", err)
	}

	// fullPatientInput is to be filled up by processing the simple patient input
	gender := domain.PatientGenderEnum(input.Gender)
	patientInput := domain.FHIRPatientInput{
		BirthDate: &input.BirthDate,
		Gender:    &gender,
		Active:    &input.Active,
	}
	patientInput.Identifier = ids
	patientInput.Telecom = contacts
	patientInput.Name = helpers.NameToHumanName(input.Names)
	// patientInput.Photo = photos
	patientInput.Address = helpers.PhysicalPostalAddressesToFHIRAddresses(
		input.PhysicalAddresses, input.PostalAddresses)
	patientInput.MaritalStatus = helpers.MaritalStatusEnumToCodeableConceptInput(
		input.MaritalStatus)
	patientInput.Communication = helpers.LanguagesToCommunicationInputs(input.Languages)

	return &patientInput, nil
}

// ContactsToContactPointInput translates phone and email contacts to
// FHIR contact points
func (c *UseCasesClinicalImpl) ContactsToContactPointInput(_ context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
	if phones == nil && emails == nil {
		return nil, nil
	}

	output := []*domain.FHIRContactPointInput{}
	rank := int64(1)
	phoneSystem := domain.ContactPointSystemEnumPhone
	use := domain.ContactPointUseEnumHome

	for _, phone := range phones {
		normalized, err := converterandformatter.NormalizeMSISDN(phone.Msisdn)
		if err != nil {
			utils.ReportErrorToSentry(err)
			return nil, fmt.Errorf("failed to normalize phonenumber")
		}

		phoneContact := &domain.FHIRContactPointInput{
			System: &phoneSystem,
			Use:    &use,
			Rank:   &rank,
			Period: common.DefaultPeriodInput(),
			Value:  normalized,
		}
		output = append(output, phoneContact)
		rank++
	}

	emailSystem := domain.ContactPointSystemEnumEmail

	for _, email := range emails {
		emailErr := utils.ValidateEmail(email.Email)
		if emailErr != nil {
			return nil, fmt.Errorf("invalid email: %w", emailErr)
		}

		emailContact := &domain.FHIRContactPointInput{
			System: &emailSystem,
			Use:    &use,
			Rank:   &rank,
			Period: common.DefaultPeriodInput(),
			Value:  &email.Email,
		}
		output = append(output, emailContact)
		rank++
	}

	return output, nil
}
