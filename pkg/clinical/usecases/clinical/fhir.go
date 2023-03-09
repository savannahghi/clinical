package clinical

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/scalarutils"
)

// FindOrganizationByID finds and retrieves organization details using the specified organization ID
func (c *UseCasesClinicalImpl) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("organization ID is required")
	}
	return c.infrastructure.FHIR.FindOrganizationByID(ctx, organizationID)
}

// GetTenantMetaTags is a helper to create tags that are used to identify which tenant a resource belongs to
// and are saved in a resources `Meta` attribute
func (c *UseCasesClinicalImpl) GetTenantMetaTags(ctx context.Context) ([]domain.FHIRCodingInput, error) {

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	organisation, err := c.FindOrganizationByID(ctx, identifiers.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tenant organisation: %w", err)
	}

	userSelected := false
	organisationTagVersion := "1.0"
	organisationTagSystem := scalarutils.URI("http://mycarehub/tenant-identification/organisation")

	tags := []domain.FHIRCodingInput{
		{
			System:       &organisationTagSystem,
			Version:      &organisationTagVersion,
			Code:         scalarutils.Code(identifiers.OrganizationID),
			Display:      *organisation.Resource.Name,
			UserSelected: &userSelected,
		},
	}

	return tags, nil
}

// RegisterPatient implements simple patient registration
func (c *UseCasesClinicalImpl) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	exist, err := c.CheckPatientExistenceUsingPhoneNumber(ctx, input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to check patient existence")
	}

	if exist {
		return nil, fmt.Errorf("patient with phone number already exists")
	}

	patientInput, err := c.SimplePatientRegistrationInputToPatientInput(ctx, input)
	if err != nil {
		return nil, err
	}

	output, err := c.CreatePatient(ctx, *patientInput)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return output, nil
}

// CheckPatientExistenceUsingPhoneNumber checks whether a patient with the phone number they're trying to register with exists
func (c *UseCasesClinicalImpl) CheckPatientExistenceUsingPhoneNumber(ctx context.Context, patientInput domain.SimplePatientRegistrationInput) (bool, error) {
	exists := false
	for _, phone := range patientInput.PhoneNumbers {
		phoneNumber := &phone.Msisdn
		patient, err := c.FindPatientsByMSISDN(ctx, *phoneNumber)
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

// FindPatientByID retrieves a single patient by their ID
func (c *UseCasesClinicalImpl) FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error) {
	if id == "" {
		return nil, fmt.Errorf("patient ID cannot be empty")
	}

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, id)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get patient with ID %s, err: %w", id, err)
	}

	patientReference := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	openEpisodes, err := c.infrastructure.FHIR.OpenEpisodes(ctx, patientReference)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to get open episodes for %s, err: %w", patientReference, err)
	}
	return &domain.PatientPayload{
		PatientRecord:   patient.Resource,
		OpenEpisodes:    openEpisodes,
		HasOpenEpisodes: len(openEpisodes) > 0,
	}, nil
}

// FindPatientsByMSISDN finds a patient's record(s), given a search term
// e.g their phone number.
//
// It intentionally does NOT have the following:
//
// 1. Pagination - if we need to paginate this data, something has gone seriously wrong
// 2. Filtering - the MSISDN is enough of a filter
// 3. Sorting - the API will take sensible choices by default
//
// Known limitations:
//
// 1. The normalization of phone number assumes Kenyan (+254) numbers only
func (c *UseCasesClinicalImpl) FindPatientsByMSISDN(ctx context.Context, msisdn string) (*domain.PatientConnection, error) {

	search, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		return nil, fmt.Errorf("can't normalize contact: %w", err)
	}
	return c.infrastructure.FHIR.SearchFHIRPatient(ctx, *search)
}

// ContactsToContactPointInput translates phone and email contacts to
// FHIR contact points
func (c *UseCasesClinicalImpl) ContactsToContactPointInput(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
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

// CreatePatient creates or updates a patient record on FHIR
func (c *UseCasesClinicalImpl) CreatePatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	// set the record ID if not set
	if input.ID == nil {
		newID := uuid.New().String()
		input.ID = &newID
	}

	if input.Gender == nil {
		return nil, fmt.Errorf("please provide the patients gender")
	}

	// set or add the default record identifier
	if input.Identifier == nil {
		input.Identifier = []*domain.FHIRIdentifierInput{common.DefaultIdentifier()}
	}
	if input.Identifier != nil {
		input.Identifier = append(input.Identifier, common.DefaultIdentifier())
	}

	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	patientRecord, err := c.infrastructure.FHIR.CreateFHIRPatient(ctx, input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return patientRecord, nil
}
