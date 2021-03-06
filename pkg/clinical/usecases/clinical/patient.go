package clinical

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"sync"

	linq "github.com/ahmetb/go-linq/v3"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/fhir"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/scalarutils"
	log "github.com/sirupsen/logrus"
)

var isc interserviceclient.InterServiceClient

// constants and defaults
const (
	// LimitedProfileEncounterCount is the number of encounters to show when a
	// patient has approved limited access to their health record
	LimitedProfileEncounterCount = 5

	NHIFImageFrontPicName            = "nhif_front_photo"
	NHIFImageRearPicName             = "nhif_rear_photo"
	RelationshipSystem               = "http://terminology.hl7.org/CodeSystem/v2-0131"
	RelationshipVersion              = "2.9"
	StringTimeParseMonthNameLayout   = "2006-Jan-02"
	StringTimeParseMonthNumberLayout = "2006-01-02"
	SavannahAdminEmail               = "SAVANNAH_ADMIN_EMAIL"
	TwilioSMSNumberEnvVarName        = "TWILIO_SMS_NUMBER"

	notFoundWithSearchParams = "could not find a patient with the provided parameters"
	internalError            = "an error occurred on our end. Please try again later"
	timeFormatStr            = "2006-01-02T15:04:05+03:00"
	//defaultTimeoutSeconds    = 10
)

// UseCasesClinical represents all the patient business logic
type UseCasesClinical interface {
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
	AddNHIF(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error)
	CreateUpdatePatientExtraInformation(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error)
	AllergySummary(ctx context.Context, patientID string) ([]string, error)
	DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error)
	StartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error)
	FindPatientsByMSISDN(ctx context.Context, msisdn string) (*domain.PatientConnection, error)
	PatientTimeline(ctx context.Context, patientID string, count int) ([]map[string]interface{}, error)
	PatientHealthTimeline(ctx context.Context, input domain.HealthTimelineInput) (*domain.HealthTimeline, error)
	GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error)
}

// UseCasesClinicalImpl represents the patient usecase implementation
type UseCasesClinicalImpl struct {
	infrastructure infrastructure.Infrastructure
	fhirUsecase    fhir.UseCasesFHIR
}

// NewUseCasesClinicalImpl initializes new Clinical/Patient implementation
func NewUseCasesClinicalImpl(infra infrastructure.Infrastructure) UseCasesClinical {
	return &UseCasesClinicalImpl{
		infrastructure: infra,
	}
}

// ProblemSummary returns a short list of the patient's active and confirmed problems (by name).
func (c *UseCasesClinicalImpl) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	if patientID == "" {
		return nil, fmt.Errorf("patient ID cannot be empty")
	}

	params := map[string]interface{}{
		"clinical-status":     "active",
		"verification-status": "confirmed",
		"category":            "problem-list-item",
		"subject":             fmt.Sprintf("Patient/%s", patientID),
	}
	results, err := c.infrastructure.FHIR.SearchFHIRCondition(ctx, params)
	if err != nil {
		utils.ReportErrorToSentry(err)
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

// VisitSummary returns a narrative friendly display of the data that has been associated with a single visit
func (c *UseCasesClinicalImpl) VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
	encounterPayload, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"encounter with ID %s not found: %w", encounterID, err)
	}
	encounter := encounterPayload.Resource
	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.ID)
	encounterFilterParams := map[string]interface{}{
		"encounter": encounterRef,
		"_count":    strconv.Itoa(count),
	}
	encounterInstanceFilterParams := map[string]interface{}{
		"_id": encounterID,
	}

	if encounterPayload.Resource.Subject == nil || encounterPayload.Resource.Subject.Reference == nil {
		return nil, fmt.Errorf("invalid: Encounter/%s has no patient reference", *encounterPayload.Resource.ID)
	}
	patientFilterParams := map[string]interface{}{
		"patient": *encounterPayload.Resource.Subject.Reference,
		"_count":  strconv.Itoa(count),
	}

	resources := []string{
		"Condition",
		"AllergyIntolerance",
		"Observation",
		"Composition",
		"MedicationRequest",
		"ServiceRequest",
		"Encounter",
	}
	nodes := make(map[string][]map[string]interface{})
	for _, resourceName := range resources {
		nodes[resourceName] = []map[string]interface{}{}
		switch resourceName {
		case "AllergyIntolerance":
			conn, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, patientFilterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := converterandformatter.StructToMap(edge.Node)
				if err != nil {
					utils.ReportErrorToSentry(err)
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "Encounter":
			conn, err := c.infrastructure.FHIR.SearchFHIREncounter(ctx, encounterInstanceFilterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := converterandformatter.StructToMap(edge.Node)
				if err != nil {
					utils.ReportErrorToSentry(err)
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "Condition":
			conn, err := c.infrastructure.FHIR.SearchFHIRCondition(ctx, encounterFilterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := converterandformatter.StructToMap(edge.Node)
				if err != nil {
					utils.ReportErrorToSentry(err)
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "Observation":
			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, encounterFilterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := converterandformatter.StructToMap(edge.Node)
				if err != nil {
					utils.ReportErrorToSentry(err)
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "Composition":
			conn, err := c.infrastructure.FHIR.SearchFHIRComposition(ctx, encounterFilterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := converterandformatter.StructToMap(edge.Node)
				if err != nil {
					utils.ReportErrorToSentry(err)
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "MedicationRequest":
			conn, err := c.infrastructure.FHIR.SearchFHIRMedicationRequest(ctx, encounterFilterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := converterandformatter.StructToMap(edge.Node)
				if err != nil {
					utils.ReportErrorToSentry(err)
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "ServiceRequest":
			conn, err := c.infrastructure.FHIR.SearchFHIRServiceRequest(ctx, encounterFilterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := converterandformatter.StructToMap(edge.Node)
				if err != nil {
					utils.ReportErrorToSentry(err)
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		default:
			// did you forget to add a case for the resource?
			return nil, fmt.Errorf(
				"server error: unknown resource %s when composing visit summary", resourceName)
		}
	}
	output := make(map[string]interface{})
	for k, v := range nodes {
		if v != nil {
			output[k] = v
		}
	}
	return output, nil
}

// PatientTimelineWithCount returns the patient's visit note timeline (a list of
// narratives that are sorted with the most recent one first), while
// respecting the approval level AND limiting the number
func (c *UseCasesClinicalImpl) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	episode, _, err := c.getTimelineEpisode(ctx, episodeID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}
	encounterSearchParams := map[string]interface{}{
		"patient": *episode.Patient.Reference,
		"sort":    "-date", // reverse chronological order
		"_count":  strconv.Itoa(count),
	}
	return c.getTimelineVisitSummaries(ctx, encounterSearchParams, count)
}

// PatientSearch searches for a patient by identifiers and names
func (c *UseCasesClinicalImpl) PatientSearch(ctx context.Context, search string) (*domain.PatientConnection, error) {

	params := url.Values{}
	params.Add("_content", search) // entire doc

	bs, err := c.infrastructure.FHIRRepo.POSTRequest("Patient", "_search", params, nil)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to search: %v", err)
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to map decode resource: %v", err)
			return nil, fmt.Errorf(internalError)
		}

		hasOpenEpisodes, err := c.infrastructure.FHIR.HasOpenEpisode(ctx, patient)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
func (c *UseCasesClinicalImpl) FindPatientsByMSISDN(ctx context.Context, msisdn string) (*domain.PatientConnection, error) {

	search, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't normalize contact: %w", err)
	}
	return c.PatientSearch(ctx, *search)
}

// CheckPatientExistenceUsingPhoneNumber checks whether a patient with the phone number they're trying to register with exists
func (c *UseCasesClinicalImpl) CheckPatientExistenceUsingPhoneNumber(ctx context.Context, patientInput domain.SimplePatientRegistrationInput) (bool, error) {
	exists := false
	for _, phone := range patientInput.PhoneNumbers {
		phoneNumber := &phone.Msisdn
		patient, err := c.FindPatientsByMSISDN(ctx, *phoneNumber)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
			return nil, fmt.Errorf("invalid email: %v", emailErr)
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
func (c *UseCasesClinicalImpl) SimplePatientRegistrationInputToPatientInput(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.FHIRPatientInput, error) {
	_, err := c.CheckPatientExistenceUsingPhoneNumber(ctx, input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to check patient existence")
	}

	contacts, err := c.ContactsToContactPointInput(ctx, input.PhoneNumbers, input.Emails)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't register patient with invalid contacts: %v", err)
	}

	ids, err := helpers.IDToIdentifier(input.IdentificationDocuments, input.PhoneNumbers)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't register patient with invalid identifiers: %v", err)
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

// RegisterPatient implements simple patient registration
func (c *UseCasesClinicalImpl) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	patientInput, err := c.SimplePatientRegistrationInputToPatientInput(ctx, input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	output, err := c.CreatePatient(ctx, *patientInput)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create patient: %v", err)
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

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn patient input into a map: %v", err)
	}

	data, err := c.infrastructure.FHIRRepo.CreateFHIRResource("Patient", payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update patient resource: %v", err)
	}
	patient := &domain.FHIRPatient{}
	err = json.Unmarshal(data, patient)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (c *UseCasesClinicalImpl) FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error) {
	if id == "" {
		return nil, fmt.Errorf("patient ID cannot be empty")
	}

	data, err := c.infrastructure.FHIRRepo.GetFHIRResource("Patient", id)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to get patient with ID %s, err: %v", id, err)
	}
	var patient domain.FHIRPatient
	err = json.Unmarshal(data, &patient)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to unmarshal patient data from JSON, err: %v", err)
	}
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	openEpisodes, err := c.infrastructure.FHIR.OpenEpisodes(ctx, patientReference)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to get open episodes for %s, err: %v", patientReference, err)
	}
	return &domain.PatientPayload{
		PatientRecord:   &patient,
		OpenEpisodes:    openEpisodes,
		HasOpenEpisodes: len(openEpisodes) > 0,
	}, nil
}

// UpdatePatient patches a patient record with fresh data.
// It updates elements that are set and ignores the ones that are nil.
func (c *UseCasesClinicalImpl) UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	op := "add" // this method replaces data at the indicated paths

	if input.ID == "" {
		return nil, fmt.Errorf("can't update with blank ID")
	}

	patientInput, err := c.SimplePatientRegistrationInputToPatientInput(ctx, input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't compose patient input: %v", err)
	}

	patientPayload, err := c.FindPatientByID(ctx, input.ID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't get patient with ID %s: %v", input.ID, err)
	}

	patches := []map[string]interface{}{}
	if patientInput.Identifier != nil {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/identifier"
		patch["value"] = patientInput.Identifier
		patches = append(patches, patch)
	}
	if patientInput.Active != patientPayload.PatientRecord.Active {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/active"
		patch["value"] = patientInput.Active
		patches = append(patches, patch)
	}
	if patientInput.Name != nil {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/name"
		patch["value"] = patientInput.Name
		patches = append(patches, patch)
	}
	if patientInput.Telecom != nil {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/telecom"
		patch["value"] = patientInput.Telecom
		patches = append(patches, patch)
	}
	if patientInput.Gender != patientPayload.PatientRecord.Gender {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/gender"
		patch["value"] = patientInput.Gender
		patches = append(patches, patch)
	}
	if patientInput.BirthDate != patientPayload.PatientRecord.BirthDate {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/birthDate"
		patch["value"] = patientInput.BirthDate
		patches = append(patches, patch)
	}
	if patientInput.Address != nil {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/address"
		patch["value"] = patientInput.Address
		patches = append(patches, patch)
	}
	if patientInput.MaritalStatus != nil {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/maritalStatus"
		patch["value"] = patientInput.MaritalStatus
		patches = append(patches, patch)
	}
	if patientInput.Photo != nil {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/photo"
		patch["value"] = patientInput.Photo
		patches = append(patches, patch)
	}
	if patientInput.Contact != nil {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/contact"
		patch["value"] = patientInput.Contact
		patches = append(patches, patch)
	}
	if patientInput.Communication != nil {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/communication"
		patch["value"] = patientInput.Communication
		patches = append(patches, patch)
	}

	data, err := c.infrastructure.FHIRRepo.PatchFHIRResource("Patient", input.ID, patches)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("UpdatePatient: %v", err)
	}
	patient := domain.FHIRPatient{}
	err = json.Unmarshal(data, &patient)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to unmarshal patient response JSON: data: %v\n, error: %v",
			string(data), err)
	}
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	openEpisodes, err := c.infrastructure.FHIR.OpenEpisodes(ctx, patientReference)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to get open episodes for %s, err: %v", patientReference, err)
	}
	return &domain.PatientPayload{
		PatientRecord:   &patient,
		OpenEpisodes:    openEpisodes,
		HasOpenEpisodes: len(openEpisodes) > 0,
	}, nil
}

// AddNextOfKin patches a patient with next of kin
func (c *UseCasesClinicalImpl) AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
	if input.PatientID == "" {
		return nil, fmt.Errorf("a patient ID must be specified")
	}

	_, err := c.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"can't get patient with ID %s: %v", input.PatientID, err)
	}

	updatedContacts := []*domain.FHIRPatientContactInput{}

	names := helpers.NameToHumanName(input.Names)
	if len(names) == 0 {
		return nil, fmt.Errorf("the contact must have a name")
	}

	contacts, err := c.ContactsToContactPointInput(
		ctx, input.PhoneNumbers, input.Emails)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("invalid contacts: %v", err)
	}

	if len(names) != 1 {
		return nil, fmt.Errorf("the contact should have one name")
	}

	addresses := helpers.PhysicalPostalAddressesToCombinedFHIRAddress(
		input.PhysicalAddresses,
		input.PostalAddresses,
	)
	userSelected := true
	relationshipSystem := scalarutils.URI(RelationshipSystem)
	relationshipVersion := RelationshipVersion
	gender := domain.PatientContactGenderEnum(input.Gender)
	if !gender.IsValid() {
		return nil, fmt.Errorf(
			"'%s' is not a valid gender; valid values are %s",
			input.Gender,
			domain.AllPatientContactGenderEnum,
		)
	}
	updatedContacts = append(updatedContacts, &domain.FHIRPatientContactInput{
		Relationship: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						Display:      domain.RelationshipTypeDisplay(input.Relationship),
						System:       &relationshipSystem,
						Version:      &relationshipVersion,
						Code:         scalarutils.Code(input.Relationship.String()),
						UserSelected: &userSelected,
					},
				},
				Text: domain.RelationshipTypeDisplay(input.Relationship),
			},
		},
		Name:    names[0],
		Telecom: contacts,
		Address: addresses,
		Gender:  &gender,
		Period:  common.DefaultPeriodInput(),
	})
	patches := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/contact",
			"value": updatedContacts,
		},
	}

	data, err := c.infrastructure.FHIRRepo.PatchFHIRResource(
		"Patient", input.PatientID, patches)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("UpdatePatient: %v", err)
	}
	patient := domain.FHIRPatient{}
	err = json.Unmarshal(data, &patient)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to unmarshal patient response JSON: data: %v\n, error: %v",
			string(data), err)
	}
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	openEpisodes, err := c.infrastructure.FHIR.OpenEpisodes(ctx, patientReference)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to get open episodes for patient %s: %w", patientReference, err)
	}

	return &domain.PatientPayload{
		PatientRecord:   &patient,
		OpenEpisodes:    openEpisodes,
		HasOpenEpisodes: len(openEpisodes) > 0,
	}, nil
}

// AddNHIF patches a patient with NHIF details
func (c *UseCasesClinicalImpl) AddNHIF(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
	if input == nil {
		return nil, fmt.Errorf("AddNHIF: nil input")
	}

	if input.PatientID == "" {
		return nil, fmt.Errorf("a patient ID must be specified")
	}

	patientPayload, err := c.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"can't get patient with ID %s: %v", input.PatientID, err)
	}

	existingIdentifiers := patientPayload.PatientRecord.Identifier
	updatedIdentifierInputs := []*domain.FHIRIdentifierInput{}
	for _, existing := range existingIdentifiers {
		updatedTypeCoding := []*domain.FHIRCodingInput{}
		for _, coding := range existing.Type.Coding {
			updatedTypeCoding = append(updatedTypeCoding, &domain.FHIRCodingInput{
				System:       coding.System,
				Version:      coding.Version,
				Code:         coding.Code,
				Display:      coding.Display,
				UserSelected: coding.UserSelected,
			})
		}
		updatedIdentifierInputs = append(updatedIdentifierInputs, &domain.FHIRIdentifierInput{
			ID:  existing.ID,
			Use: existing.Use,
			Type: domain.FHIRCodeableConceptInput{
				ID:     existing.Type.ID,
				Text:   existing.Type.Text,
				Coding: updatedTypeCoding,
			},
			System: existing.System,
			Value:  existing.Value,
			Period: &domain.FHIRPeriodInput{
				ID:    existing.Period.ID,
				Start: existing.Period.Start,
				End:   existing.Period.End,
			},
		})
	}
	patches := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/identifier",
			"value": updatedIdentifierInputs,
		},
	}

	data, err := c.infrastructure.FHIRRepo.PatchFHIRResource(
		"Patient", input.PatientID, patches)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("UpdatePatient: %v", err)
	}
	patient := domain.FHIRPatient{}
	err = json.Unmarshal(data, &patient)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to unmarshal patient response JSON: data: %v\n, error: %v",
			string(data), err)
	}
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	openEpisodes, err := c.infrastructure.FHIR.OpenEpisodes(ctx, patientReference)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to get open episodes for %s, err: %v", patientReference, err)
	}
	return &domain.PatientPayload{
		PatientRecord:   &patient,
		OpenEpisodes:    openEpisodes,
		HasOpenEpisodes: len(openEpisodes) > 0,
	}, nil
}

// CreateUpdatePatientExtraInformation updates a patient's extra info
func (c *UseCasesClinicalImpl) CreateUpdatePatientExtraInformation(
	ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
	if input.PatientID == "" {
		return false, fmt.Errorf("patient ID cannot empty: %v", input.PatientID)
	}

	patientPayload, err := c.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to get patient with ID %s: %w", input.PatientID, err)
	}
	patient := patientPayload.PatientRecord

	patches := []map[string]interface{}{}
	op := "add" // the content will be appended to the element identified in the path

	if input.MaritalStatus != nil {
		inputMaritalStatus := helpers.MaritalStatusEnumToCodeableConcept(*input.MaritalStatus)
		if patient.MaritalStatus != inputMaritalStatus {
			patch := make(map[string]interface{})
			patch["op"] = op
			patch["path"] = "/maritalStatus"
			patch["value"] = inputMaritalStatus
			patches = append(patches, patch)
		}
	}

	if input.Languages != nil {
		langs := []enumutils.Language{}
		for _, l := range input.Languages {
			langs = append(langs, *l)
		}
		communicationInput := helpers.LanguagesToCommunicationInputs(langs)
		if len(input.Languages) > 0 {
			patch := make(map[string]interface{})
			patch["op"] = op
			patch["path"] = "/communication"
			patch["value"] = communicationInput
			patches = append(patches, patch)
		}
	}

	_, err = c.ContactsToContactPointInput(ctx, nil, input.Emails)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("an error occurred: %v", err)
	}

	_, err = c.infrastructure.FHIRRepo.PatchFHIRResource("Patient", input.PatientID, patches)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("UpdatePatient: %v", err)
	}
	return true, nil
}

// AllergySummary returns a short list of the patient's active and confirmed
// allergies (by name)
func (c *UseCasesClinicalImpl) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	if patientID == "" {
		return nil, fmt.Errorf("patient ID cannot be empty")
	}
	params := map[string]interface{}{
		"clinical-status":     "active",
		"verification-status": "confirmed",
		"type":                "allergy",
		"criticality":         "high",
		"patient":             fmt.Sprintf("Patient/%s", patientID),
	}
	results, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("error when searching for patient allergies: %w", err)
	}
	output := []string{}
	for _, edge := range results.Edges {
		allergy := edge.Node
		if allergy.Code == nil {
			return nil, fmt.Errorf("server error: every allergy must have a code")
		}
		if allergy.Code.Text == "" {
			return nil, fmt.Errorf("server error: every allergy code must have it's text set")
		}
		output = append(output, allergy.Code.Text)
	}
	return output, nil
}

// DeleteFHIRPatientByPhone delete's a patient's FHIR compartment
// given their phone number
func (c *UseCasesClinicalImpl) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	patient, err := c.FindPatientsByMSISDN(ctx, phoneNumber)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to find patient by phone number")
	}

	edges := patient.Edges
	var patientID string
	for _, edge := range edges {
		patientID = *edge.Node.ID
	}

	return c.infrastructure.FHIR.DeleteFHIRPatient(ctx, patientID)
}

//StartEpisodeByBreakGlass starts an emergency episode
func (c *UseCasesClinicalImpl) StartEpisodeByBreakGlass(
	ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	normalized, err := converterandformatter.NormalizeMSISDN(input.ProviderPhone)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to normalize phone number: %w", err)
	}

	// validatePhone patient phone number
	validatePhone, err := converterandformatter.NormalizeMSISDN(input.PatientPhone)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("invalid patient phone number: %v", err)
	}
	// alert patient
	err = c.sendAlertToPatient(ctx, *validatePhone, input.PatientID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		log.Printf("failed to send alert message during StartEpisodeByBreakGlass login: %s", err)
	}
	// alert next-of-kin
	err = c.sendAlertToNextOfKin(ctx, input.PatientID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		log.Printf("failed to send alert message to next of kin during StartEpisodeByBreakGlass login: %s", err)
	}

	// alert admin
	pp, err := c.FindPatientByID(ctx, input.PatientID)
	if err == nil {
		patientName := pp.PatientRecord.Name[0].Text
		err = c.sendAlertToAdmin(ctx, patientName, *normalized)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Printf("failed to send alert message to admin during StartEpisodeByBreakGlass login: %s", err)
		}
	}
	organizationID, err := c.fhirUsecase.GetORCreateOrganization(ctx, input.ProviderCode)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"internal server error in retrieving service provider : %v", err)
	}
	ep := helpers.ComposeOneHealthEpisodeOfCare(
		*normalized,
		input.FullAccess,
		*organizationID,
		input.ProviderCode,
		input.PatientID,
	)
	return c.infrastructure.FHIR.CreateEpisodeOfCare(ctx, ep)
}

// PatientTimeline return's the patient's historical timeline sorted in descending order i.e when it was first recorded
// The timeline consists of Allergies, Observations, Medication statement and Test results
func (c *UseCasesClinicalImpl) PatientTimeline(ctx context.Context, patientID string, count int) ([]map[string]interface{}, error) {

	timeline := []map[string]interface{}{}
	wg := &sync.WaitGroup{}
	mut := &sync.Mutex{}

	patientFilterParams := map[string]interface{}{
		"patient": fmt.Sprintf("Patient/%v", patientID),
	}

	// timelineResourceFunc is a go routine that fetches particular FHIR resource and
	// adds it to the timeline
	type timelineResourceFunc func(wg *sync.WaitGroup, mut *sync.Mutex)

	allergyIntoleranceResourceFunc := func(wg *sync.WaitGroup, mut *sync.Mutex) {
		defer wg.Done()

		conn, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, patientFilterParams)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("AllergyIntolerance search error: %v", err)
			return
		}
		for _, edge := range conn.Edges {
			if edge.Node == nil {
				continue
			}

			rMap, err := converterandformatter.StructToMap(edge.Node)
			if err != nil {
				utils.ReportErrorToSentry(err)
				log.Errorf("AllergyIntolerance edge struct to map error: %v", err)
				return
			}
			if rMap == nil {
				continue
			}

			rMap["resourceType"] = "AllergyIntolerance"
			rMap["timelineDate"] = rMap["recordedDate"]

			mut.Lock()
			timeline = append(timeline, rMap)
			mut.Unlock()
		}
	}

	observationResourceFunc := func(wg *sync.WaitGroup, mut *sync.Mutex) {
		defer wg.Done()

		conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, patientFilterParams)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("Observation search error: %v", err)
			return
		}
		for _, edge := range conn.Edges {
			if edge.Node == nil {
				continue
			}

			rMap, err := converterandformatter.StructToMap(edge.Node)
			if err != nil {
				utils.ReportErrorToSentry(err)
				log.Errorf("Observation edge struct to map error: %v", err)
				return
			}
			if rMap == nil {
				continue
			}

			rMap["resourceType"] = "Observation"
			rMap["timelineDate"] = rMap["effectiveInstant"]

			mut.Lock()
			timeline = append(timeline, rMap)
			mut.Unlock()
		}
	}

	medicationStatementResourceFunc := func(wg *sync.WaitGroup, mut *sync.Mutex) {
		defer wg.Done()

		conn, err := c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, patientFilterParams)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("MedicationStatement search error: %v", err)
			return
		}
		for _, edge := range conn.Edges {
			if edge.Node == nil {
				continue
			}

			rMap, err := converterandformatter.StructToMap(edge.Node)
			if err != nil {
				utils.ReportErrorToSentry(err)
				log.Errorf("MedicationStatement edge struct to map error: %v", err)
				return
			}
			if rMap == nil {
				continue
			}

			rMap["resourceType"] = "MedicationStatement"
			rMap["timelineDate"] = rMap["effectiveDateTime"]

			mut.Lock()
			timeline = append(timeline, rMap)
			mut.Unlock()
		}
	}

	resources := []timelineResourceFunc{
		allergyIntoleranceResourceFunc,
		observationResourceFunc,
		medicationStatementResourceFunc,
	}

	for _, resource := range resources {
		wg.Add(1)
		go resource(wg, mut)
	}

	wg.Wait()

	return timeline, nil
}

// GetMedicalData returns a limited subset of specific medical data that for a specific patient
// These include: Allergies, Viral Load, Body Mass Index, Weight, CD4 Count using their respective OCL CIEL Terminology
// For each category the latest three records are fetched
func (c *UseCasesClinicalImpl) GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error) {
	data := &domain.MedicalData{}

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

	for _, field := range fields {
		switch field {
		case "Regimen":
			conn, err := c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				data.Regimen = append(data.Regimen, edge.Node)
			}
		case "AllergyIntolerance":
			conn, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				data.Allergies = append(data.Allergies, edge.Node)
			}

		case "Weight":
			filterParams["code"] = common.WeightCIELTerminologyCode
			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				data.Weight = append(data.Weight, edge.Node)
			}

		case "BMI":
			filterParams["code"] = common.BMICIELTerminologyCode
			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				data.BMI = append(data.BMI, edge.Node)
			}

		case "ViralLoad":
			filterParams["code"] = common.ViralLoadCIELTerminologyCode
			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				data.ViralLoad = append(data.ViralLoad, edge.Node)
			}

		case "CD4Count":
			filterParams["code"] = common.CD4CountCIELTerminologyCode
			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				data.CD4Count = append(data.CD4Count, edge.Node)
			}

		}

	}

	return data, nil
}

// PatientHealthTimeline return's the patient's historical timeline sorted in descending order i.e when it was first recorded
// The timeline consists of Allergies, Observations, Medication statement and Test results
func (c *UseCasesClinicalImpl) PatientHealthTimeline(ctx context.Context, input domain.HealthTimelineInput) (*domain.HealthTimeline, error) {
	records, err := c.PatientTimeline(ctx, input.PatientID, 0)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("cannot retrieve patient timeline error: %w", err)
	}

	data := &domain.HealthTimeline{}
	timeline := []map[string]interface{}{}

	sortFunc := func(i, j interface{}) bool {
		itemI := i.(map[string]interface{})
		dateStringI, ok := itemI["timelineDate"].(string)
		if !ok {
			return false
		}
		timeI := helpers.ParseDate(dateStringI)

		itemJ := j.(map[string]interface{})
		dateStringJ, _ := itemJ["timelineDate"].(string)
		if !ok {
			return false
		}
		timeJ := helpers.ParseDate(dateStringJ)

		return timeI.After(timeJ)
	}

	linq.From(records).Sort(sortFunc).Skip(input.Offset).Take(input.Limit).ToSlice(&timeline)

	data.TotalCount = len(records)
	data.Timeline = timeline

	return data, nil
}
