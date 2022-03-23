package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/mitchellh/mapstructure"
	auth "github.com/savannahghi/clinical/pkg/clinical/application/authorization"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/scalarutils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
)

// constants and defaults
const (
	// LimitedProfileEncounterCount is the number of encounters to show when a
	// patient has approved limited access to their health record
	LimitedProfileEncounterCount = 5

	// MaxClinicalRecordPageSize is the maximum number of encounters we can show on a timeline
	MaxClinicalRecordPageSize = 50

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
	fullAccessLevel          = "FULL_ACCESS"
	partialAccessLevel       = "PROFILE_AND_RECENT_VISITS_ACCESS"
	timeFormatStr            = "2006-01-02T15:04:05+03:00"
	baseFHIRURL              = "https://healthcare.googleapis.com/v1"
	cloudhealthEmail         = "cloudhealth@healthcloud.co.ke"
	defaultTimeoutSeconds    = 10
)

// FHIRUseCase represents all the FHIR business logic
type FHIRUseCase interface {
	CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error)
	CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	OpenOrganizationEpisodes(
		ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error)
	GetORCreateOrganization(ctx context.Context, providerSladeCode string) (*string, error)
	GetOrganization(ctx context.Context, providerSladeCode string) (*string, error)
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	CreateOrganization(ctx context.Context, providerSladeCode string) (*string, error)
	SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error)
	POSTRequest(
		resourceName string, path string, params url.Values, body io.Reader) ([]byte, error)
	SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error)
	HasOpenEpisode(
		ctx context.Context,
		patient domain.FHIRPatient,
	) (bool, error)
	OpenEpisodes(
		ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error)
	CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error)
	GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error)
	Encounters(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum) ([]*domain.FHIREncounter, error)
	SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error)
	StartEncounter(ctx context.Context, episodeID string) (string, error)
	StartEpisodeByOtp(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error)
	UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error)
	SearchEpisodeEncounter(
		ctx context.Context,
		episodeReference string,
	) (*domain.FHIREncounterRelayConnection, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	EndEpisode(ctx context.Context, episodeID string) (bool, error)
	GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error)
	SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error)
	CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error)
	SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error)
	CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error)
	CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	DeleteFHIRComposition(ctx context.Context, id string) (bool, error)
	SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error)
	UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error)
	SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error)
	SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error)
	CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error)
	SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error)
	CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error)
	DeleteFHIRObservation(ctx context.Context, id string) (bool, error)
	GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error)
	DeleteFHIRPatient(ctx context.Context, id string) (bool, error)
	DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error)
	DeleteFHIRResourceType(results []map[string]string) error
}

// FHIRUseCaseImpl represents the FHIR usecase implementation
type FHIRUseCaseImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewFHIRUseCaseImpl initializes the new FHIR implementation
func NewFHIRUseCaseImpl(infra infrastructure.Infrastructure) FHIRUseCase {
	return &FHIRUseCaseImpl{
		infrastructure: infra,
	}
}

// Encounters returns encounters that belong to the indicated patient.
//
// The patientReference should be a [string] in the format "Patient/<patient resource ID>".
func (fh FHIRUseCaseImpl) Encounters(
	ctx context.Context,
	patientReference string,
	status *domain.EncounterStatusEnum,
) ([]*domain.FHIREncounter, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	searchParams := url.Values{}
	if status != nil {
		searchParams.Add("status:exact", status.String())
	}
	searchParams.Add("patient", patientReference)

	bs, err := fh.POSTRequest("Encounter", "_search", searchParams, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to search for encounter: %v", err)
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal FHIR encounter search response: %v", err)
	}

	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
	for _, k := range mandatoryKeys {
		_, found := respMap[k]
		if !found {
			return nil, fmt.Errorf("search response does not have key '%s'", k)
		}
	}
	resourceType, ok := respMap["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("search: the resourceType is not a string")
	}
	if resourceType != "Bundle" {
		return nil, fmt.Errorf("search: the resourceType value is not 'Bundle' as expected")
	}
	resultType, ok := respMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("search: the type is not a string")
	}
	if resultType != "searchset" {
		return nil, fmt.Errorf("search: the type value is not 'searchset' as expected")
	}

	output := []*domain.FHIREncounter{}
	respEntries := respMap["entry"]
	if respEntries == nil {
		return output, nil
	}
	entries, ok := respEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf("search: entries is not a list of maps, it is: %T", respEntries)
	}

	for _, en := range entries {
		entry, ok := en.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected each entry to be map, they are %T instead", en)
		}
		expectedKeys := []string{"fullUrl", "resource", "search"}
		for _, k := range expectedKeys {
			_, found := entry[k]
			if !found {
				return nil, fmt.Errorf("search entry does not have key '%s'", k)
			}
		}
		resource := entry["resource"]
		var encounter domain.FHIREncounter
		resourceBs, err := json.Marshal(resource)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %v", err)
		}
		err = json.Unmarshal(resourceBs, &encounter)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %v", err)
		}
		output = append(output, &encounter)
	}
	return output, nil
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (fh FHIRUseCaseImpl) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := fh.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "EpisodeOfCare"
	path := "_search"
	output := domain.FHIREpisodeOfCareRelayConnection{}

	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIREpisodeOfCare

		resourceBs, err := json.Marshal(result)
		if err != nil {
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &domain.FHIREpisodeOfCareRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateEpisodeOfCare is the final common pathway for creation of episodes of
// care.
func (fh FHIRUseCaseImpl) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	payload, err := converterandformatter.StructToMap(episode)
	if err != nil {
		return nil, fmt.Errorf("unable to turn episode of care input into a map: %v", err)
	}

	// search for the episode of care before creating new one.
	episodeOfCareSearchParams := map[string]interface{}{
		"patient":      fmt.Sprintf(*episode.Patient.Reference),
		"status":       string(domain.EpisodeOfCareStatusEnumActive),
		"organization": *episode.ManagingOrganization.Reference,
		"_sort":        "date",
		"_count":       "1",
	}

	episodeOfCarePayload, err := fh.SearchFHIREpisodeOfCare(ctx, episodeOfCareSearchParams)
	if err != nil {
		return nil, fmt.Errorf("unable to get patients episodes of care: %v", err)
	}

	// don't create a new episode if there is an ongoing one
	if len(episodeOfCarePayload.Edges) >= 1 {
		episodeOfCare := *episodeOfCarePayload.Edges[0].Node
		encounters, err := fh.Encounters(ctx, *episodeOfCare.Patient.Reference, nil)
		if err == nil {
			output := &domain.EpisodeOfCarePayload{
				EpisodeOfCare: &episodeOfCare,
				TotalVisits:   len(encounters),
			}
			return output, nil
		}
	}

	// create a new episode if none has been found
	data, err := fh.infrastructure.FHIRRepo.CreateFHIRResource("EpisodeOfCare", payload)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to create episode of care resource: %v", err)
	}
	fhirEpisode := &domain.FHIREpisodeOfCare{}
	err = json.Unmarshal(data, fhirEpisode)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal episode of care response JSON: data: %v\n, error: %v",
			string(data), err)
	}

	encounters, err := fh.Encounters(ctx, *episode.Patient.Reference, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get encounters for episode %s: %v",
			*episode.ID, err,
		)
	}
	output := &domain.EpisodeOfCarePayload{
		EpisodeOfCare: fhirEpisode,
		TotalVisits:   len(encounters),
	}
	return output, nil
}

// CreateFHIRCondition creates a FHIRCondition instance
func (fh FHIRUseCaseImpl) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	resourceType := "Condition"
	resource := domain.FHIRCondition{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRConditionRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// CreateFHIROrganization creates a FHIROrganization instance
func (fh FHIRUseCaseImpl) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	resourceType := "Organization"
	resource := domain.FHIROrganization{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIROrganizationRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// OpenOrganizationEpisodes return all organization specific open episodes
func (fh FHIRUseCaseImpl) OpenOrganizationEpisodes(
	ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	organizationID, err := fh.GetORCreateOrganization(ctx, providerSladeCode)
	if err != nil {
		return nil, fmt.Errorf(
			"internal server error in retrieving service provider : %v", err)
	}
	organizationReference := fmt.Sprintf("Organization/%s", *organizationID)
	searchParams := url.Values{}
	searchParams.Add("status", domain.EpisodeOfCareStatusEnumActive.String())
	searchParams.Add("organization", organizationReference)
	return fh.SearchEpisodesByParam(ctx, searchParams)
}

// GetORCreateOrganization retrieve an organisation via its code if not found create a new one.
func (fh FHIRUseCaseImpl) GetORCreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	retrievedOrg, err := fh.GetOrganization(ctx, providerSladeCode)
	if err != nil {
		return nil, fmt.Errorf(
			"internal server error in getting organisation : %v", err)
	}
	if retrievedOrg == nil {
		createdOrg, err := fh.CreateOrganization(ctx, providerSladeCode)
		if err != nil {
			return nil, fmt.Errorf(
				"internal server error in creating organisation : %v", err)
		}
		return createdOrg, nil
	}
	return retrievedOrg, nil
}

// CreateOrganization creates an organization given ist provider code
func (fh FHIRUseCaseImpl) CreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	identifier := []*domain.FHIRIdentifierInput{
		{
			Use:   "official",
			Value: providerSladeCode,
		},
	}
	organizationInput := domain.FHIROrganizationInput{
		Identifier: identifier,
		Name:       &providerSladeCode,
	}
	createdOrganization, err := fh.CreateFHIROrganization(ctx, organizationInput)
	if err != nil {
		return nil, err
	}
	organisationID := createdOrganization.Resource.ID
	return organisationID, nil
}

// SearchFHIROrganization provides a search API for FHIROrganization
func (fh FHIRUseCaseImpl) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := fh.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Organization"
	path := "_search"
	output := domain.FHIROrganizationRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIROrganization

		resourceBs, err := json.Marshal(result)
		if err != nil {
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &domain.FHIROrganizationRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// GetOrganization retrieves an organization given its code
func (fh FHIRUseCaseImpl) GetOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	searchParam := map[string]interface{}{
		"identifier": providerSladeCode,
	}
	organization, err := fh.SearchFHIROrganization(ctx, searchParam)
	if err != nil {
		return nil, err
	}
	if organization.Edges == nil {
		return nil, nil
	}
	ORGID := organization.Edges[0].Node.ID
	return ORGID, nil
}

// SearchEpisodesByParam search episodes by params
func (fh FHIRUseCaseImpl) SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {

	bs, err := fh.POSTRequest(
		"EpisodeOfCare", "_search", searchParams, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to search for episode of care: %v", err)
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal FHIR episode of care search response: %v", err)
	}

	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
	for _, k := range mandatoryKeys {
		_, found := respMap[k]
		if !found {
			return nil, fmt.Errorf("search response does not have key '%s'", k)
		}
	}
	resourceType, ok := respMap["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("search: the resourceType is not a string")
	}
	if resourceType != "Bundle" {
		return nil, fmt.Errorf("search: the resourceType value is not 'Bundle' as expected")
	}
	resultType, ok := respMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("search: the type is not a string")
	}
	if resultType != "searchset" {
		return nil, fmt.Errorf("search: the type value is not 'searchset' as expected")
	}

	output := []*domain.FHIREpisodeOfCare{}
	respEntries := respMap["entry"]
	if respEntries == nil {
		return output, nil
	}
	entries, ok := respEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf("search: entries is not a list of maps, it is: %T", respEntries)
	}

	for _, en := range entries {
		entry, ok := en.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected each entry to be map, they are %T instead", en)
		}
		expectedKeys := []string{"fullUrl", "resource", "search"}
		for _, k := range expectedKeys {
			_, found := entry[k]
			if !found {
				return nil, fmt.Errorf("search entry does not have key '%s'", k)
			}
		}
		resource := entry["resource"].(map[string]interface{})

		period := resource["period"].(map[string]interface{})

		// parse period->start as map[string]interface{}
		resStart := helpers.ParseDate(period["start"].(string))
		startDateAsMap := make(map[string]int)
		startDateAsMap["year"] = resStart.Year()
		startDateAsMap["month"] = int(resStart.Month())
		startDateAsMap["day"] = resStart.Day()
		period["start"] = scalarutils.DateTime(resStart.Format(timeFormatStr))

		// parse period->end as map[string]interface{}
		resEnd := helpers.ParseDate(period["end"].(string))
		endDateAsMap := make(map[string]int)
		endDateAsMap["year"] = resEnd.Year()
		endDateAsMap["month"] = int(resEnd.Month())
		endDateAsMap["day"] = resEnd.Day()
		period["end"] = scalarutils.DateTime(resEnd.Format(timeFormatStr))

		//update the original period resource
		resource["period"] = period

		var episode domain.FHIREpisodeOfCare

		err := mapstructure.Decode(resource, &episode)
		if err != nil {
			log.Errorf("unable to map decode resource: %v", err)
			return nil, fmt.Errorf(internalError)
		}

		output = append(output, &episode)
	}
	return output, nil
}

// POSTRequest is used to manually compose POST requests to the FHIR service
//
// - `resourceName` is a FHIR resource name e.g "Patient"
// - `path` is a sub-path e.g `_search` under a resource
// - `params` should be query params, sent as `url.Values`
func (fh FHIRUseCaseImpl) POSTRequest(
	resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
	fhirHeaders, err := fh.FHIRHeaders()
	if err != nil {
		return nil, fmt.Errorf("unable to get FHIR headers: %v", err)
	}

	url := fmt.Sprintf(
		"%s/%s/%s?%s", fh.infrastructure.FHIRRepo.FHIRRestURL(), resourceName, path, params.Encode())
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to compose FHIR POST request: %v", err)
	}
	for k, v := range fhirHeaders {
		for _, h := range v {
			req.Header.Add(k, h)
		}
	}

	httpClient := &http.Client{Timeout: time.Second * defaultTimeoutSeconds}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP response error: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(
			"search: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// OpenEpisodes returns the IDs of a patient's open episodes
func (fh FHIRUseCaseImpl) OpenEpisodes(
	ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	searchParams := url.Values{}
	searchParams.Add("status:exact", domain.EpisodeOfCareStatusEnumActive.String())
	searchParams.Add("patient", patientReference)
	return fh.SearchEpisodesByParam(ctx, searchParams)
}

// HasOpenEpisode determines if a patient has an open episode
func (fh FHIRUseCaseImpl) HasOpenEpisode(
	ctx context.Context,
	patient domain.FHIRPatient,
) (bool, error) {
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	episodes, err := fh.OpenEpisodes(ctx, patientReference)
	if err != nil {
		return false, err
	}
	return len(episodes) > 0, nil
}

// GetBearerToken logs in and gets a Google bearer auth token.
// The user referred to by `cloudhealthEmail` needs to have IAM permissions
// that allow them to read and write from the project's Cloud Healthcare base.
func (fh FHIRUseCaseImpl) GetBearerToken() (string, error) {
	ctx := context.Background()
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
	creds, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		return "", fmt.Errorf("default creds error: %v", err)
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("oauth token error: %v", err)
	}
	return fmt.Sprintf("Bearer %s", token.AccessToken), nil
}

// FHIRHeaders composes suitable FHIR headers, with authentication and content
// type already set
func (fh FHIRUseCaseImpl) FHIRHeaders() (http.Header, error) {
	headers := make(map[string][]string)
	bearerHeader, err := fh.GetBearerToken()
	if err != nil {
		return nil, fmt.Errorf("can't get bearer token: %v", err)
	}
	headers["Content-Type"] = []string{"application/fhir+json; charset=utf-8"}
	headers["Accept"] = []string{"application/fhir+json; charset=utf-8"}
	headers["Authorization"] = []string{bearerHeader}
	return headers, nil
}

// CreateFHIREncounter creates a FHIREncounter instance
func (fh FHIRUseCaseImpl) CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	resourceType := "Encounter"
	resource := domain.FHIREncounter{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIREncounterRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// GetFHIREpisodeOfCare retrieves instances of FHIREpisodeOfCare by ID
func (fh FHIRUseCaseImpl) GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {

	resourceType := "EpisodeOfCare"
	var resource domain.FHIREpisodeOfCare

	data, err := fh.infrastructure.FHIRRepo.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &domain.FHIREpisodeOfCareRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// StartEncounter starts an encounter within an episode of care
func (c *ClinicalUseCaseImpl) StartEncounter(
	ctx context.Context, episodeID string) (string, error) {
	episodePayload, err := c.fhir.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return "", fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}
	activeEpisodeStatus := domain.EpisodeOfCareStatusEnumActive
	activeEncounterStatus := domain.EncounterStatusEnumInProgress
	if episodePayload.Resource.Status.String() != activeEpisodeStatus.String() {
		return "", fmt.Errorf("an encounter can only be started for an active episode")
	}
	episodeRef := fmt.Sprintf("EpisodeOfCare/%s", *episodePayload.Resource.ID)

	now := time.Now()
	startTime := scalarutils.DateTime(now.Format("2006-01-02T15:04:05+03:00"))

	encounterClassCode := scalarutils.Code("AMB")
	encounterClassSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/v3-ActCode")
	encounterClassVersion := "2018-08-12"
	encounterClassDisplay := "ambulatory"
	encounterClassUserSelected := false

	encounterInput := domain.FHIREncounterInput{
		Status: activeEncounterStatus,
		Class: domain.FHIRCodingInput{
			System:       &encounterClassSystem,
			Version:      &encounterClassVersion,
			Code:         encounterClassCode,
			Display:      encounterClassDisplay,
			UserSelected: &encounterClassUserSelected,
		},
		Subject: &domain.FHIRReferenceInput{
			Reference: episodePayload.Resource.Patient.Reference,
			Display:   episodePayload.Resource.Patient.Display,
			Type:      episodePayload.Resource.Patient.Type,
		},
		EpisodeOfCare: []*domain.FHIRReferenceInput{
			{
				Reference: &episodeRef,
			},
		},
		ServiceProvider: &domain.FHIRReferenceInput{
			Display: episodePayload.Resource.ManagingOrganization.Display,
			Type:    episodePayload.Resource.ManagingOrganization.Type,
		},
		Period: &domain.FHIRPeriodInput{
			Start: startTime,
		},
	}
	encPl, err := c.fhir.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		return "", fmt.Errorf("unable to start encounter: %w", err)
	}
	return *encPl.Resource.ID, nil
}

// StartEncounter starts an encounter within an episode of care
func (fh *FHIRUseCaseImpl) StartEncounter(
	ctx context.Context, episodeID string) (string, error) {
	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return "", fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}
	activeEpisodeStatus := domain.EpisodeOfCareStatusEnumActive
	activeEncounterStatus := domain.EncounterStatusEnumInProgress
	if episodePayload.Resource.Status.String() != activeEpisodeStatus.String() {
		return "", fmt.Errorf("an encounter can only be started for an active episode")
	}
	episodeRef := fmt.Sprintf("EpisodeOfCare/%s", *episodePayload.Resource.ID)

	now := time.Now()
	startTime := scalarutils.DateTime(now.Format("2006-01-02T15:04:05+03:00"))

	encounterClassCode := scalarutils.Code("AMB")
	encounterClassSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/v3-ActCode")
	encounterClassVersion := "2018-08-12"
	encounterClassDisplay := "ambulatory"
	encounterClassUserSelected := false

	encounterInput := domain.FHIREncounterInput{
		Status: activeEncounterStatus,
		Class: domain.FHIRCodingInput{
			System:       &encounterClassSystem,
			Version:      &encounterClassVersion,
			Code:         encounterClassCode,
			Display:      encounterClassDisplay,
			UserSelected: &encounterClassUserSelected,
		},
		Subject: &domain.FHIRReferenceInput{
			Reference: episodePayload.Resource.Patient.Reference,
			Display:   episodePayload.Resource.Patient.Display,
			Type:      episodePayload.Resource.Patient.Type,
		},
		EpisodeOfCare: []*domain.FHIRReferenceInput{
			{
				Reference: &episodeRef,
			},
		},
		ServiceProvider: &domain.FHIRReferenceInput{
			Display: episodePayload.Resource.ManagingOrganization.Display,
			Type:    episodePayload.Resource.ManagingOrganization.Type,
		},
		Period: &domain.FHIRPeriodInput{
			Start: startTime,
		},
	}
	encPl, err := fh.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		return "", fmt.Errorf("unable to start encounter: %w", err)
	}
	return *encPl.Resource.ID, nil
}

// StartEpisodeByOtp starts a patient OTP verified episode
func (fh *FHIRUseCaseImpl) StartEpisodeByOtp(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	isVerified, normalized, err := fh.infrastructure.Engagement.VerifyOTP(ctx, input.Msisdn, input.Otp)
	if err != nil {
		log.Printf(
			"invalid phone: \nPhone: %s\nOTP: %s\n", input.Msisdn, input.Otp)
		return nil, fmt.Errorf(
			"invalid phone: got %s when validating %s", err, input.Msisdn)
	}
	if !isVerified {
		return nil, fmt.Errorf("invalid OTP")
	}

	organizationID, err := fh.GetORCreateOrganization(ctx, input.ProviderCode)
	if err != nil {
		return nil, fmt.Errorf(
			"internal server error in retrieving service provider : %v", err)
	}
	ep := helpers.ComposeOneHealthEpisodeOfCare(
		normalized,
		input.FullAccess,
		*organizationID,
		input.ProviderCode,
		input.PatientID,
	)
	return fh.CreateEpisodeOfCare(ctx, ep)
}

// UpgradeEpisode starts a patient OTP verified episode
func (fh *FHIRUseCaseImpl) UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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
	// retrieve and validate the episode
	episode, err := fh.GetActiveEpisode(ctx, input.EpisodeID)
	if err != nil {
		return nil, fmt.Errorf("can't get active episode to upgrade: %w", err)
	}
	if episode == nil {
		return nil, fmt.Errorf("system error: nil episode of care")
	}
	episodeTypes := episode.Type
	if episodeTypes == nil {
		return nil, fmt.Errorf("system error: nil episode type")
	}
	if len(episodeTypes) != 1 {
		return nil, fmt.Errorf(
			"system error: expected episode type to have exactly one entry, got %d", len(episodeTypes))
	}
	if episodeTypes[0] == nil {
		return nil, fmt.Errorf("system error: nil episode")
	}
	encounters, err := fh.Encounters(ctx, *episode.Patient.Reference, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get encounters for episode %s: %v",
			*episode.ID, err,
		)
	}

	// if it already has the correct access level, return early
	if episodeTypes[0].Text == fullAccessLevel {
		return &domain.EpisodeOfCarePayload{
			EpisodeOfCare: episode,
			TotalVisits:   len(encounters),
		}, nil
	}

	// validate the MSISDN and OTP
	isVerified, _, err := fh.infrastructure.Engagement.VerifyOTP(ctx, input.Msisdn, input.Otp)
	if err != nil {
		log.Printf(
			"invalid phone: \nPhone: %s\nOTP: %s\n", input.Msisdn, input.Otp)
		return nil, fmt.Errorf(
			"invalid phone: got %s when validating %s", err, input.Msisdn)
	}
	if !isVerified {
		return nil, fmt.Errorf("invalid OTP")
	}

	// patch the episode status
	episode.Type = []*domain.FHIRCodeableConcept{
		{Text: fullAccessLevel},
	}
	payload, err := converterandformatter.StructToMap(episode)
	if err != nil {
		return nil, fmt.Errorf("unable to turn episode of care input into a map: %v", err)
	}

	_, err = fh.infrastructure.FHIRRepo.UpdateFHIRResource(
		"EpisodeOfCare", *episode.ID, payload)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to update episode of care resource: %v", err)
	}
	return &domain.EpisodeOfCarePayload{
		EpisodeOfCare: episode,
		TotalVisits:   len(encounters),
	}, nil
}

// SearchEpisodeEncounter returns all encounters in a visit
func (fh *FHIRUseCaseImpl) SearchEpisodeEncounter(
	ctx context.Context,
	episodeReference string,
) (*domain.FHIREncounterRelayConnection, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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
	episodeRef := fmt.Sprintf("Episode/%s", episodeReference)
	encounterFilterParams := map[string]interface{}{
		"episodeOfCare": episodeRef,
		"status":        "in_progress",
	}
	encounterConn, err := fh.SearchFHIREncounter(ctx, encounterFilterParams)

	if err != nil {
		return nil, fmt.Errorf("unable to search encounter: %w", err)
	}

	return encounterConn, nil
}

// EndEncounter ends an encounter
func (fh *FHIRUseCaseImpl) EndEncounter(
	ctx context.Context, encounterID string) (bool, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get user: %w", err)
	}
	isAuthorized, err := auth.IsAuthorized(user, auth.ProblemSummaryView)
	if err != nil {
		return false, err
	}
	if !isAuthorized {
		return false, fmt.Errorf("user not authorized to access this resource")
	}
	resourceType := "Encounter"
	encounterPayload, err := fh.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return false, fmt.Errorf("unable to get encounter with ID %s: %w", encounterID, err)
	}
	updatedStatus := domain.EncounterStatusEnumFinished
	encounterPayload.Resource.Status = updatedStatus

	// workaround for odd date comparison behavior on the Google Cloud Healthcare API
	// the end time must be at least 24 hours after the start time
	// so: if the time now is less than 24 hours after start, set the end to be
	// 24 hours after the start of the visit. If the time now is more than 24 hours
	// after the start, use the current time as the end of the visit
	end := time.Now().Add(time.Hour * 24)
	endTime := scalarutils.DateTime(end.Format(timeFormatStr))
	encounterPayload.Resource.Period.End = endTime

	payload, err := converterandformatter.StructToMap(encounterPayload.Resource)
	if err != nil {
		return false, fmt.Errorf("unable to turn the updated episode of care into a map: %v", err)
	}

	_, err = fh.infrastructure.FHIRRepo.UpdateFHIRResource(resourceType, encounterID, payload)
	if err != nil {
		return false, fmt.Errorf("unable to create/update %s resource: %w", resourceType, err)
	}
	return true, nil
}

// EndEpisode ends an episode of care by patching it's status to "finished"
func (fh *FHIRUseCaseImpl) EndEpisode(
	ctx context.Context, episodeID string) (bool, error) {
	resourceType := "EpisodeOfCare"
	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return false, fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}
	startTime := scalarutils.DateTime(time.Now().Format(timeFormatStr))
	if episodePayload.Resource.Period != nil {
		startTime = episodePayload.Resource.Period.Start
	}

	// Close all encounters in this visit
	encounterConn, err := fh.SearchEpisodeEncounter(ctx, episodeID)
	if err != nil {
		return false, fmt.Errorf("unable to search episode encounter %w", err)
	}
	for _, edge := range encounterConn.Edges {
		_, err = fh.EndEncounter(ctx, *edge.Node.ID)
		if err != nil {
			log.Printf("unable to end encounter %s", *edge.Node.ID)
			continue
		}
	}
	// // workaround for odd date comparison behavior on the Google Cloud Healthcare API
	// the end time must be at least 24 hours after the start time
	// so: if the time now is less than 24 hours after start, set the end to be
	// 24 hours after the start of the visit. If the time now is more than 24 hours
	// after the start, use the current time as the end of the visit
	end := time.Now().Add(time.Hour * 24)
	endTime := scalarutils.DateTime(end.Format(timeFormatStr))

	updatedStatus := domain.EpisodeOfCareStatusEnumFinished
	episodePayload.Resource.Status = &updatedStatus
	episodePayload.Resource.Period.Start = startTime
	episodePayload.Resource.Period.End = endTime
	payload, err := converterandformatter.StructToMap(episodePayload.Resource)
	if err != nil {
		return false, fmt.Errorf("unable to turn the updated episode of care into a map: %v", err)
	}

	_, err = fh.infrastructure.FHIRRepo.UpdateFHIRResource(resourceType, episodeID, payload)
	if err != nil {
		return false, fmt.Errorf("unable to create/update %s resource: %w", resourceType, err)
	}
	return true, nil
}

// GetActiveEpisode returns any ACTIVE episode that has to the indicated ID
func (fh *FHIRUseCaseImpl) GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {

	searchParams := url.Values{}
	searchParams.Add("status:exact", domain.EpisodeOfCareStatusEnumActive.String())
	searchParams.Add("_id", episodeID) // logical ID of the resource

	bs, err := fh.POSTRequest(
		"EpisodeOfCare", "_search", searchParams, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to search for episode of care: %v", err)
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal FHIR episode of care search response: %v", err)
	}

	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
	for _, k := range mandatoryKeys {
		_, found := respMap[k]
		if !found {
			return nil, fmt.Errorf("search response does not have key '%s'", k)
		}
	}
	resourceType, ok := respMap["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("search: the resourceType is not a string")
	}
	if resourceType != "Bundle" {
		return nil, fmt.Errorf("search: the resourceType value is not 'Bundle' as expected")
	}
	resultType, ok := respMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("search: the type is not a string")
	}
	if resultType != "searchset" {
		return nil, fmt.Errorf("search: the type value is not 'searchset' as expected")
	}

	respEntries := respMap["entry"]
	if respEntries == nil {
		return nil, fmt.Errorf("there is no ACTIVE episode with the ID %s", episodeID)
	}
	entries, ok := respEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf("search: entries is not a list of maps, it is: %T", respEntries)
	}
	if len(entries) != 1 {
		return nil, fmt.Errorf(
			"expected exactly one ACTIVE episode for episode ID %s, got %d", episodeID, len(entries))
	}

	entry, ok := entries[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected each entry to be map, they are %T instead", entry)
	}
	expectedKeys := []string{"fullUrl", "resource", "search"}
	for _, k := range expectedKeys {
		_, found := entry[k]
		if !found {
			return nil, fmt.Errorf("search entry does not have key '%s'", k)
		}
	}
	resource := entry["resource"]
	var episode domain.FHIREpisodeOfCare
	resourceBs, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal resource to JSON: %v", err)
	}
	err = json.Unmarshal(resourceBs, &episode)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal resource: %v", err)
	}
	return &episode, nil
}

// SearchFHIRServiceRequest provides a search API for FHIRServiceRequest
func (fh *FHIRUseCaseImpl) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := fh.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "ServiceRequest"
	path := "_search"
	output := domain.FHIRServiceRequestRelayConnection{}

	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRServiceRequest

		resourceBs, err := json.Marshal(result)
		if err != nil {
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &domain.FHIRServiceRequestRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRServiceRequest creates a FHIRServiceRequest instance
func (fh *FHIRUseCaseImpl) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	resourceType := "ServiceRequest"
	resource := domain.FHIRServiceRequest{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRServiceRequestRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// SearchFHIRAllergyIntolerance provides a search API for FHIRAllergyIntolerance
func (fh *FHIRUseCaseImpl) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := fh.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "AllergyIntolerance"
	path := "_search"
	output := domain.FHIRAllergyIntoleranceRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRAllergyIntolerance

		resourceBs, err := json.Marshal(result)
		if err != nil {
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &domain.FHIRAllergyIntoleranceRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRAllergyIntolerance creates a FHIRAllergyIntolerance instance
func (fh *FHIRUseCaseImpl) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	resourceType := "AllergyIntolerance"
	resource := domain.FHIRAllergyIntolerance{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRAllergyIntoleranceRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIRAllergyIntolerance updates a FHIRAllergyIntolerance instance
// The resource must have it's ID set.
func (fh *FHIRUseCaseImpl) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	resourceType := "AllergyIntolerance"
	resource := domain.FHIRAllergyIntolerance{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRAllergyIntoleranceRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// SearchFHIRComposition provides a search API for FHIRComposition
func (fh *FHIRUseCaseImpl) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := fh.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Composition"
	path := "_search"
	output := domain.FHIRCompositionRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRComposition

		resourceBs, err := json.Marshal(result)
		if err != nil {
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &domain.FHIRCompositionRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRComposition creates a FHIRComposition instance
func (fh *FHIRUseCaseImpl) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	resourceType := "Composition"
	resource := domain.FHIRComposition{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRCompositionRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIRComposition updates a FHIRComposition instance
// The resource must have it's ID set.
func (fh *FHIRUseCaseImpl) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	resourceType := "Composition"
	resource := domain.FHIRComposition{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRCompositionRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIRComposition deletes the FHIRComposition identified by the supplied ID
func (fh *FHIRUseCaseImpl) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get user: %w", err)
	}
	isAuthorized, err := auth.IsAuthorized(user, auth.ProblemSummaryView)
	if err != nil {
		return false, err
	}
	if !isAuthorized {
		return false, fmt.Errorf("user not authorized to access this resource")
	}

	resourceType := "Composition"
	resp, err := fh.infrastructure.FHIRRepo.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// SearchFHIRCondition provides a search API for FHIRCondition
func (fh *FHIRUseCaseImpl) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := fh.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Condition"
	path := "_search"
	output := domain.FHIRConditionRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRCondition

		resourceBs, err := json.Marshal(result)
		if err != nil {
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &domain.FHIRConditionRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// UpdateFHIRCondition updates a FHIRCondition instance
// The resource must have it's ID set.
func (fh *FHIRUseCaseImpl) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	resourceType := "Condition"
	resource := domain.FHIRCondition{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRConditionRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// GetFHIREncounter retrieves instances of FHIREncounter by ID
func (fh *FHIRUseCaseImpl) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	resourceType := "Encounter"
	var resource domain.FHIREncounter

	data, err := fh.infrastructure.FHIRRepo.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &domain.FHIREncounterRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// SearchFHIREncounter provides a search API for FHIREncounter
func (fh *FHIRUseCaseImpl) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := fh.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Encounter"
	path := "_search"
	output := domain.FHIREncounterRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIREncounter

		resourceBs, err := json.Marshal(result)
		if err != nil {
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &domain.FHIREncounterRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// SearchFHIRMedicationRequest provides a search API for FHIRMedicationRequest
func (fh *FHIRUseCaseImpl) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := fh.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "MedicationRequest"
	path := "_search"
	output := domain.FHIRMedicationRequestRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRMedicationRequest

		resourceBs, err := json.Marshal(result)
		if err != nil {
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &domain.FHIRMedicationRequestRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRMedicationRequest creates a FHIRMedicationRequest instance
func (fh *FHIRUseCaseImpl) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	resourceType := "MedicationRequest"
	resource := domain.FHIRMedicationRequest{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRMedicationRequestRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIRMedicationRequest updates a FHIRMedicationRequest instance
// The resource must have it's ID set.
func (fh *FHIRUseCaseImpl) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	resourceType := "MedicationRequest"
	resource := domain.FHIRMedicationRequest{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRMedicationRequestRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIRMedicationRequest deletes the FHIRMedicationRequest identified by the supplied ID
func (fh *FHIRUseCaseImpl) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get user: %w", err)
	}
	isAuthorized, err := auth.IsAuthorized(user, auth.ProblemSummaryView)
	if err != nil {
		return false, err
	}
	if !isAuthorized {
		return false, fmt.Errorf("user not authorized to access this resource")
	}
	resourceType := "MedicationRequest"
	resp, err := fh.infrastructure.FHIRRepo.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// SearchFHIRObservation provides a search API for FHIRObservation
func (fh *FHIRUseCaseImpl) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := fh.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Observation"
	path := "_search"
	output := domain.FHIRObservationRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRObservation

		resourceBs, err := json.Marshal(result)
		if err != nil {
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &domain.FHIRObservationRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRObservation creates a FHIRObservation instance
func (fh *FHIRUseCaseImpl) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
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

	resourceType := "Observation"
	resource := domain.FHIRObservation{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.infrastructure.FHIRRepo.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRObservationRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIRObservation deletes the FHIRObservation identified by the passed ID
func (fh *FHIRUseCaseImpl) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	user, err := fh.infrastructure.BaseExtension.GetLoggedInUser(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get user: %w", err)
	}
	isAuthorized, err := auth.IsAuthorized(user, auth.ProblemSummaryView)
	if err != nil {
		return false, err
	}
	if !isAuthorized {
		return false, fmt.Errorf("user not authorized to access this resource")
	}

	resourceType := "Observation"
	resp, err := fh.infrastructure.FHIRRepo.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// GetFHIRPatient retrieves instances of FHIRPatient by ID
func (fh *FHIRUseCaseImpl) GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {

	resourceType := "Patient"
	var resource domain.FHIRPatient

	data, err := fh.infrastructure.FHIRRepo.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	hasOpenEpisodes, err := fh.HasOpenEpisode(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get open episodes for patient %#v: %w", resource, err)
	}
	payload := &domain.FHIRPatientRelayPayload{
		Resource:        &resource,
		HasOpenEpisodes: hasOpenEpisodes,
	}
	return payload, nil
}

// DeleteFHIRPatient deletes the FHIRPatient identified by the supplied ID
func (fh *FHIRUseCaseImpl) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	patientEverythingBs, err := fh.infrastructure.FHIRRepo.GetFHIRPatientEverything(id)
	if err != nil {
		return false, fmt.Errorf("unable to get patient's compartment: %v", err)
	}

	var patientEverything map[string]interface{}
	err = json.Unmarshal(patientEverythingBs, &patientEverything)
	if err != nil {
		return false, fmt.Errorf("unable to unmarshal patient everything")
	}

	respEntries := patientEverything["entry"]
	if respEntries == nil {
		return false, nil
	}

	entries, ok := respEntries.([]interface{})
	if !ok {
		return false, fmt.Errorf(
			"server error: entries is not a list of maps, it is: %T", respEntries)
	}

	// This list stores assorted ResourceTypes and ResourceIDs found in an Encounter
	// i.e Observations, Medication Request etc
	assortedResourceTypes := []map[string]string{}

	// This list stores all the encounters ResourceType and ResourceID in a patient's compartment
	encounters := []map[string]string{}

	// This list stores all the Episodesofcare ResourceType and ResourceIDs in a patient's compartment
	episodesOfCare := []map[string]string{}

	// This list stores the patient ResourceType and ResourceID
	patient := []map[string]string{}

	for _, en := range entries {
		entry, ok := en.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf(
				"server error: expected each entry to be map, they are %T instead",
				en,
			)
		}
		resource, ok := entry["resource"]
		if !ok {
			{
				return false, fmt.Errorf(
					"server error: result entry %#v is not a map",
					entry["resource"],
				)
			}
		}
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf(
				"server error: expected each entry to be map, they are %T instead",
				resource,
			)
		}

		resourceType := resourceMap["resourceType"].(string)

		resourceTypeIDMap := map[string]string{
			"resourceType": resourceType,
			"resourceID":   resourceMap["id"].(string),
		}

		switch resourceType {
		case "Encounter":
			encounters = append(
				encounters,
				resourceTypeIDMap,
			)
			continue

		case "EpisodeOfCare":
			episodesOfCare = append(
				episodesOfCare,
				resourceTypeIDMap,
			)
			continue
		case "Patient":
			patient = append(
				patient,
				resourceTypeIDMap,
			)
			continue
		}

		assortedResourceTypes = append(assortedResourceTypes, resourceTypeIDMap)
	}

	// Order of deletion matters to avoid conflicts
	// First delete the ResourceTypes found in an encounter
	if err = fh.DeleteFHIRResourceType(assortedResourceTypes); err != nil {
		return false, err
	}

	// Secondly, delete the encounters. This will bring no conflict
	// as it ensures the any ResourceType that refers to the encounter is not found
	if err = fh.DeleteFHIRResourceType(encounters); err != nil {
		return false, err
	}

	// Thirdly, delete the episodes of care. This will bring no conflict
	// as it ensures the any Encounter that refers to the EpisodeOfCare is not found
	if err = fh.DeleteFHIRResourceType(episodesOfCare); err != nil {
		return false, err
	}

	// Finally delete the patient ResourceType
	if err = fh.DeleteFHIRResourceType(patient); err != nil {
		return false, err
	}

	return true, nil
}

// DeleteFHIRResourceType takes a ResourceType and ID and deletes them from FHIR
func (fh *FHIRUseCaseImpl) DeleteFHIRResourceType(results []map[string]string) error {
	for _, result := range results {
		resourceType := result["resourceType"]
		resourceID := result["resourceID"]

		resp, err := fh.infrastructure.FHIRRepo.DeleteFHIRResource(
			resourceType,
			resourceID,
		)
		if err != nil {
			return fmt.Errorf(
				"unable to delete %s, response %s, error: %v",
				resourceType, string(resp), err,
			)
		}
	}
	return nil
}

// DeleteFHIRServiceRequest deletes the FHIRServiceRequest identified by the supplied ID
func (fh *FHIRUseCaseImpl) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	resourceType := "ServiceRequest"
	resp, err := fh.infrastructure.FHIRRepo.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}
