package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	auth "github.com/savannahghi/clinical/pkg/clinical/application/authorization"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/scalarutils"
	log "github.com/sirupsen/logrus"
)

// ClinicalUseCase ...
type ClinicalUseCase interface {
	ProblemSummary(ctx context.Context, patientID string) ([]string, error)
	VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error)
	PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error)
	ContactsToContactPointInput(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error)
	RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	CreatePatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error)
	FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error)
	PatientSearch(ctx context.Context, search string) (*domain.PatientConnection, error)
	UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error)
	AddNhif(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error)
	RegisterUser(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	CreateUpdatePatientExtraInformation(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error)
	AllergySummary(ctx context.Context, patientID string) ([]string, error)
	DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error)
	StartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error)
}

// ClinicalUseCaseImpl ...
type ClinicalUseCaseImpl struct {
	infrastructure infrastructure.Infrastructure
	fhir           FHIRUseCase
}

// NewClinicalUseCaseImpl ...
func NewClinicalUseCaseImpl(infra infrastructure.Infrastructure, fhir FHIRUseCase) ClinicalUseCase {
	return &ClinicalUseCaseImpl{
		infrastructure: infra,
		fhir:           fhir,
	}
}

// ProblemSummary ...
func (c *ClinicalUseCaseImpl) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	user, err := profileutils.GetLoggedInUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get user: %w", err)
	}
	isAuthorized, err := auth.IsAuthorized(user, auth.ProblemSummaryView)
	if err != nil {
		return nil, err
	}
	if !isAuthorized {
		return nil, fmt.Errorf("user not authorized to access this resource")
	}

	params := map[string]interface{}{
		"clinical-status":     "active",
		"verification-status": "confirmed",
		"category":            "problem-list-item",
		"subject":             fmt.Sprintf("Patient/%s", patientID),
	}
	results, err := c.fhir.SearchFHIRCondition(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error when searching for patient conditions: %w", err)
	}
	output := []string{}
	for _, conditionEdge := range results.Edges {
		condition := conditionEdge.Node
		if condition.Code == nil {
			return nil, fmt.Errorf("server error: every condition must have a code")
		}
		if condition.Code.Text == "" {
			return nil, fmt.Errorf("server error: every condition code must have it's text set")
		}
		output = append(output, condition.Code.Text)
	}
	return output, nil
}

// VisitSummary ...
func (c *ClinicalUseCaseImpl) VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
	return nil, nil
}

// PatientTimelineWithCount ...
func (c *ClinicalUseCaseImpl) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	return nil, nil
}

func (c *ClinicalUseCaseImpl) birthdateMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	parsedDate := helpers.ParseDate(resourceCopy["birthDate"].(string))

	dateMap := make(map[string]interface{})

	dateMap["year"] = parsedDate.Year()
	dateMap["month"] = parsedDate.Month()
	dateMap["day"] = parsedDate.Day()

	resourceCopy["birthDate"] = dateMap

	return resourceCopy

}

func (c *ClinicalUseCaseImpl) periodMapper(period map[string]interface{}) map[string]interface{} {

	periodCopy := period

	parsedStartDate := helpers.ParseDate(periodCopy["start"].(string))

	periodCopy["start"] = scalarutils.DateTime(parsedStartDate.Format(timeFormatStr))

	parsedEndDate := helpers.ParseDate(periodCopy["end"].(string))

	periodCopy["end"] = scalarutils.DateTime(parsedEndDate.Format(timeFormatStr))

	return periodCopy
}

func (c *ClinicalUseCaseImpl) identifierMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	if _, ok := resource["identifier"]; ok {

		newIdentifiers := []map[string]interface{}{}

		for _, identifier := range resource["identifier"].([]interface{}) {

			identifier := identifier.(map[string]interface{})

			if _, ok := identifier["period"]; ok {

				period := identifier["period"].(map[string]interface{})
				newPeriod := c.periodMapper(period)

				identifier["period"] = newPeriod
			}

			newIdentifiers = append(newIdentifiers, identifier)
		}

		resourceCopy["identifier"] = newIdentifiers
	}

	return resourceCopy
}

func (c *ClinicalUseCaseImpl) nameMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newNames := []map[string]interface{}{}

	if _, ok := resource["name"]; ok {

		for _, name := range resource["name"].([]interface{}) {

			name := name.(map[string]interface{})

			if _, ok := name["period"]; ok {

				period := name["period"].(map[string]interface{})
				newPeriod := c.periodMapper(period)

				name["period"] = newPeriod
			}

			newNames = append(newNames, name)
		}

	}

	resourceCopy["name"] = newNames

	return resourceCopy
}

func (c *ClinicalUseCaseImpl) telecomMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newTelecoms := []map[string]interface{}{}

	if _, ok := resource["telecom"]; ok {

		for _, telecom := range resource["telecom"].([]interface{}) {

			telecom := telecom.(map[string]interface{})

			if _, ok := telecom["period"]; ok {

				period := telecom["period"].(map[string]interface{})
				newPeriod := c.periodMapper(period)

				telecom["period"] = newPeriod
			}

			newTelecoms = append(newTelecoms, telecom)
		}

	}

	resourceCopy["telecom"] = newTelecoms

	return resourceCopy
}

func (c *ClinicalUseCaseImpl) addressMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newAddresses := []map[string]interface{}{}

	if _, ok := resource["address"]; ok {

		for _, address := range resource["address"].([]interface{}) {

			address := address.(map[string]interface{})

			if _, ok := address["period"]; ok {

				period := address["period"].(map[string]interface{})
				newPeriod := c.periodMapper(period)

				address["period"] = newPeriod
			}

			newAddresses = append(newAddresses, address)
		}
	}

	resourceCopy["address"] = newAddresses

	return resourceCopy
}

func (c *ClinicalUseCaseImpl) photoMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newPhotos := []map[string]interface{}{}

	if _, ok := resource["photo"]; ok {

		for _, photo := range resource["photo"].([]interface{}) {

			photo := photo.(map[string]interface{})

			parsedDate := helpers.ParseDate(photo["creation"].(string))

			photo["creation"] = scalarutils.DateTime(parsedDate.Format(timeFormatStr))

			newPhotos = append(newPhotos, photo)
		}
	}

	resourceCopy["photo"] = newPhotos

	return resourceCopy
}

func (c *ClinicalUseCaseImpl) contactMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newContacts := []map[string]interface{}{}

	if _, ok := resource["contact"]; ok {

		for _, contact := range resource["contact"].([]interface{}) {

			contact := contact.(map[string]interface{})

			if _, ok := contact["name"]; ok {

				name := contact["name"].(map[string]interface{})
				if _, ok := name["period"]; ok {

					period := name["period"].(map[string]interface{})
					newPeriod := c.periodMapper(period)

					name["period"] = newPeriod
				}

				contact["name"] = name
			}

			if _, ok := contact["telecom"]; ok {

				newTelecoms := []map[string]interface{}{}

				for _, telecom := range contact["telecom"].([]interface{}) {

					telecom := telecom.(map[string]interface{})

					if _, ok := telecom["period"]; ok {

						period := telecom["period"].(map[string]interface{})
						newPeriod := c.periodMapper(period)

						telecom["period"] = newPeriod
					}

					newTelecoms = append(newTelecoms, telecom)
				}

				contact["telecom"] = newTelecoms
			}

			if _, ok := contact["period"]; ok {

				period := contact["period"].(map[string]interface{})
				newPeriod := c.periodMapper(period)

				contact["period"] = newPeriod
			}

			newContacts = append(newContacts, contact)
		}
	}

	resourceCopy["contact"] = newContacts

	return resourceCopy
}

// PatientSearch searches for a patient by identifiers and names
func (c *ClinicalUseCaseImpl) PatientSearch(ctx context.Context, search string) (*domain.PatientConnection, error) {

	params := url.Values{}
	params.Add("_content", search) // entire doc

	bs, err := c.fhir.POSTRequest("Patient", "_search", params, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to search: %v", err)
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		log.Errorf("unable to unmarshal FHIR search response: %v", err)
		return nil, fmt.Errorf(notFoundWithSearchParams)
	}

	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
	for _, k := range mandatoryKeys {
		_, found := respMap[k]
		if !found {
			log.Errorf("search response does not have key '%s'", k)
			return nil, fmt.Errorf(notFoundWithSearchParams)
		}
	}
	resourceType, ok := respMap["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("search: the resourceType is not a string")
	}
	if resourceType != "Bundle" {
		log.Errorf("Search: the resourceType value is not 'Bundle' as expected")
		return nil, fmt.Errorf(notFoundWithSearchParams)
	}

	resultType, ok := respMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("search: the type is not a string")
	}
	if resultType != "searchset" {
		log.Errorf("Search: the type value is not 'searchset' as expected")
		return nil, fmt.Errorf(notFoundWithSearchParams)
	}

	respEntries := respMap["entry"]
	if respEntries == nil {
		return &domain.PatientConnection{
			Edges:    []*domain.PatientEdge{},
			PageInfo: &firebasetools.PageInfo{},
		}, nil
	}
	entries, ok := respEntries.([]interface{})
	if !ok {
		log.Errorf("Search: entries is not a list of maps, it is: %T", respEntries)
		return nil, fmt.Errorf(notFoundWithSearchParams)
	}

	output := domain.PatientConnection{}
	for _, en := range entries {
		entry, ok := en.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected each entry to be map, they are %T instead", en)
		}
		expectedKeys := []string{"fullUrl", "resource", "search"}
		for _, k := range expectedKeys {
			_, found := entry[k]
			if !found {
				log.Errorf("search entry does not have key '%s'", k)
				return nil, fmt.Errorf(notFoundWithSearchParams)
			}
		}

		resource := entry["resource"].(map[string]interface{})

		resource = c.birthdateMapper(resource)
		resource = c.identifierMapper(resource)
		resource = c.nameMapper(resource)
		resource = c.telecomMapper(resource)
		resource = c.addressMapper(resource)
		resource = c.photoMapper(resource)
		resource = c.contactMapper(resource)

		var patient domain.FHIRPatient

		err := mapstructure.Decode(resource, &patient)
		if err != nil {
			log.Errorf("unable to map decode resource: %v", err)
			return nil, fmt.Errorf(internalError)
		}

		hasOpenEpisodes, err := c.fhir.HasOpenEpisode(ctx, patient)
		if err != nil {
			log.Errorf("error while checking if hasOpenEpisodes: %v", err)
			return nil, fmt.Errorf(internalError)
		}
		output.Edges = append(output.Edges, &domain.PatientEdge{
			Node:            &patient,
			HasOpenEpisodes: hasOpenEpisodes,
		})
	}
	return &output, nil
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
func (c *ClinicalUseCaseImpl) FindPatientsByMSISDN(ctx context.Context, msisdn string) (*domain.PatientConnection, error) {

	search, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		return nil, fmt.Errorf("can't normalize contact: %w", err)
	}
	return c.PatientSearch(ctx, *search)
}

// CheckPatientExistenceUsingPhoneNumber checks whether a patient with the phone number they're trying to register with exists
func (c *ClinicalUseCaseImpl) CheckPatientExistenceUsingPhoneNumber(ctx context.Context, patientInput domain.SimplePatientRegistrationInput) (bool, error) {
	exists := false
	for _, phone := range patientInput.PhoneNumbers {
		phoneNumber := &phone.Msisdn
		patient, err := c.FindPatientsByMSISDN(ctx, *phoneNumber)
		if err != nil {
			return false, fmt.Errorf("unable to find patient")
		}
		if len(patient.Edges) > 1 {
			exists = true
			break
		}
	}
	return exists, nil
}

// ContactsToContactPointInput translates phone and email contacts to
// FHIR contact points
func (c *ClinicalUseCaseImpl) ContactsToContactPointInput(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
	if phones == nil && emails == nil {
		return nil, nil
	}
	output := []*domain.FHIRContactPointInput{}
	rank := int64(1)
	phoneSystem := domain.ContactPointSystemEnumPhone
	use := domain.ContactPointUseEnumHome

	for _, phone := range phones {
		if phone.IsUssd {
			continue // don't verify USSD
		}
		isVerified, normalized, err := c.infrastructure.Engagement.VerifyOTP(
			ctx, phone.Msisdn, phone.VerificationCode)
		if err != nil {
			return nil, fmt.Errorf("invalid phone: %w", err)
		}
		if !isVerified {
			return nil, fmt.Errorf("invalid OTP")
		}
		phoneContact := &domain.FHIRContactPointInput{
			System: &phoneSystem,
			Use:    &use,
			Rank:   &rank,
			Period: common.DefaultPeriodInput(),
			Value:  &normalized,
		}
		output = append(output, phoneContact)
		rank++
	}

	emailSystem := domain.ContactPointSystemEnumEmail
	for _, email := range emails {
		err := c.infrastructure.FirestoreRepo.ValidateEmail(ctx, email.Email, email.CommunicationOptIn)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %v", err)
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

// SimplePatientRegistrationInputToPatientInput transforms a patient input into
// a
func (c *ClinicalUseCaseImpl) SimplePatientRegistrationInputToPatientInput(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.FHIRPatientInput, error) {
	exists, err := c.CheckPatientExistenceUsingPhoneNumber(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to check patient existence")
	}
	if exists {
		return nil, fmt.Errorf("a patient registered with that phone number already exists")
	}

	contacts, err := c.ContactsToContactPointInput(ctx, input.PhoneNumbers, input.Emails)
	if err != nil {
		return nil, fmt.Errorf("can't register patient with invalid contacts: %v", err)
	}

	ids, err := helpers.IDToIdentifier(input.IdentificationDocuments, input.PhoneNumbers)
	if err != nil {
		return nil, fmt.Errorf("can't register patient with invalid identifiers: %v", err)
	}

	photos, err := c.infrastructure.Engagement.PhotosToAttachments(ctx, input.Photos)
	if err != nil {
		return nil, fmt.Errorf("can't process patient photos: %v", err)
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
	patientInput.Photo = photos
	patientInput.Address = helpers.PhysicalPostalAddressesToFHIRAddresses(
		input.PhysicalAddresses, input.PostalAddresses)
	patientInput.MaritalStatus = helpers.MaritalStatusEnumToCodeableConceptInput(
		input.MaritalStatus)
	patientInput.Communication = helpers.LanguagesToCommunicationInputs(input.Languages)
	return &patientInput, nil
}

// RegisterPatient implements simple patient registration
func (c *ClinicalUseCaseImpl) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	patientInput, err := c.SimplePatientRegistrationInputToPatientInput(ctx, input)
	if err != nil {
		return nil, err
	}
	output, err := c.CreatePatient(ctx, *patientInput)
	if err != nil {
		return nil, fmt.Errorf("unable to create patient: %v", err)
	}
	for _, patientEmail := range input.Emails {
		err = c.infrastructure.Engagement.SendPatientWelcomeEmail(ctx, patientEmail.Email)
		if err != nil {
			return nil, fmt.Errorf("unable to send welcome email: %w", err)
		}
	}

	return output, nil
}

// CreatePatient creates or updates a patient record on FHIR
func (c *ClinicalUseCaseImpl) CreatePatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	// set the record ID if not set
	if input.ID == nil {
		newID := uuid.New().String()
		input.ID = &newID
	}

	// set or add the default record identifier
	if input.Identifier == nil {
		input.Identifier = []*domain.FHIRIdentifierInput{common.DefaultIdentifier()}
	}
	if input.Identifier != nil {
		input.Identifier = append(input.Identifier, common.DefaultIdentifier())
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn patient input into a map: %v", err)
	}

	data, err := c.infrastructure.FHIRRepo.CreateFHIRResource("Patient", payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update patient resource: %v", err)
	}
	patient := &domain.FHIRPatient{}
	err = json.Unmarshal(data, patient)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal patient response JSON: data: %v\n, error: %v",
			string(data), err)
	}
	output := &domain.PatientPayload{
		PatientRecord:   patient,
		HasOpenEpisodes: false, // the patient is newly created so we can safely assume this
		OpenEpisodes:    []*domain.FHIREpisodeOfCare{},
	}
	return output, nil
}

// FindPatientByID retrieves a single patient by their ID
func (c *ClinicalUseCaseImpl) FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error) {
	return nil, nil
}

// UpdatePatient patches a patient record with fresh data.
// It updates elements that are set and ignores the ones that are nil.
func (c *ClinicalUseCaseImpl) UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return nil, nil
}

// AddNextOfKin patches a patient with next of kin
func (c *ClinicalUseCaseImpl) AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
	return nil, nil
}

// AddNhif patches a patient with NHIF details
func (c *ClinicalUseCaseImpl) AddNhif(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
	return nil, nil
}

// RegisterUser implements creates a user profile and simple patient registration
func (c *ClinicalUseCaseImpl) RegisterUser(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return nil, nil
}

// CreateUpdatePatientExtraInformation updates a patient's extra info
func (c *ClinicalUseCaseImpl) CreateUpdatePatientExtraInformation(
	ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
	return false, nil
}

// AllergySummary returns a short list of the patient's active and confirmed
// allergies (by name)
func (c *ClinicalUseCaseImpl) AllergySummary(
	ctx context.Context, patientID string) ([]string, error) {
	return nil, nil
}

// DeleteFHIRPatientByPhone delete's a patient's FHIR compartment
// given their phone number
func (c *ClinicalUseCaseImpl) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	return false, nil
}

//StartEpisodeByBreakGlass starts an emergency episode
func (c *ClinicalUseCaseImpl) StartEpisodeByBreakGlass(
	ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	return nil, nil
}
