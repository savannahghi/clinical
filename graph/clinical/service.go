// Package clinical implements a simplified GraphQL interface to a FHIR server
// that acts as a clinical data repository.
package clinical

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"text/template"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/cloudhealth/cloudhealth"
)

// constants and defaults
const (
	// LimitedProfileEncounterCount is the number of encounters to show when a
	// patient has approved limited access to their health record
	LimitedProfileEncounterCount = 5

	// MaxClinicalRecordPageSize is the maximum number of encounters we can show on a timeline
	MaxClinicalRecordPageSize = 50

	// CenturyHours is the number of hours in a (fictional) century of leap years
	CenturyHours                     = 878400
	BreakGlassCollectionName         = "break_glass"
	DefaultLanguage                  = "English"
	NHIFImageFrontPicName            = "nhif_front_photo"
	NHIFImageRearPicName             = "nhif_rear_photo"
	RelationshipSystem               = "http://terminology.hl7.org/CodeSystem/v2-0131"
	RelationshipVersion              = "2.9"
	CallCenterNumber                 = "0790 360 360"
	StringTimeParseMonthNameLayout   = "2006-Jan-02"
	StringTimeParseMonthNumberLayout = "2006-01-02"
	SavannahAdminEmail               = "SAVANNAH_ADMIN_EMAIL"
	EmailWelcomeSubject              = "Welcome to Be.Well"
	TwilioSMSNumberEnvVarName        = "TWILIO_SMS_NUMBER"

	notFoundWithSearchParams = "could not find a patient with the provided parameters"
	internalError            = "an error occured on our end. Please try again later"
	fullAccessLevel          = "FULL_ACCESS"
	partialAccessLevel       = "PROFILE_AND_RECENT_VISITS_ACCESS"
	timeFormatStr            = "2006-01-02T15:04:05+03:00"
)

// dependencies names. Should match the names in the yaml file
const (
	mailgunService = "mailgun"
	smsService     = "sms"
	twilioService  = "twilio"
)

// specific endpoint paths for ISC
const (
	// mailgun isc paths
	sendEmail = "internal/send_email"

	// twilio isc paths
	sendTwilioSMS = "internal/send_sms"

	// sms isc paths
	sendSMS = "internal/send_sms"
)

// NewService initializes a new clinical service
func NewService() *Service {
	var config *base.DepsConfig
	config, err := base.LoadDepsFromYAML()
	if err != nil {
		log.Panicf("unable to load dependencies from YAML: %s", err)
	}

	twilioClient, err := base.SetupISCclient(*config, twilioService)
	if err != nil {
		log.Panicf("unable to set up Twilio ISC client: %v", err)
	}

	mailgunClient, err := base.SetupISCclient(*config, mailgunService)
	if err != nil {
		log.Panicf("unable to set up Mailgun ISC client: %v", err)
	}

	smsClient, err := base.SetupISCclient(*config, smsService)
	if err != nil {
		log.Panicf("unable to set up SMS ISC client: %v", err)
	}

	smsISC := &base.SmsISC{
		Isc:      smsClient,
		EndPoint: sendSMS,
	}
	twilioISC := &base.SmsISC{
		Isc:      twilioClient,
		EndPoint: sendTwilioSMS,
	}

	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase app: %s", err)
	}

	ctx := context.Background()
	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore client: %s", err)
	}
	clinicalRepository := cloudhealth.NewService()
	return &Service{
		clinicalRepository: clinicalRepository,
		firestoreClient:    firestoreClient,
		mailgun:            mailgunClient,
		twilio:             twilioISC,
		sms:                smsISC,
	}
}

// Service is a clinical service
type Service struct {
	clinicalRepository *cloudhealth.Service
	mailgun            *base.InterServiceClient
	twilio             *base.SmsISC
	sms                *base.SmsISC
	firestoreClient    *firestore.Client
}

func (s Service) checkPreconditions() {
	if s.clinicalRepository == nil {
		log.Panicf("*cloudhealth.Service is nil")
	}

	if s.firestoreClient == nil {
		log.Panicf("nil firestore client in health passport service")
	}

	if s.mailgun == nil {
		log.Panicf("nil mailgun ISC in health passport service")
	}

	if s.twilio == nil {
		log.Panicf("nil twilio ISC in health passport service")
	}

	if s.sms == nil {
		log.Panicf("nil SMS service in health passport service")
	}
}

func (s Service) validateSearchParams(params map[string]interface{}) (url.Values, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	output := url.Values{}
	for k, v := range params {
		val, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("the search/filter params should all be sent as strings")
		}
		output.Add(k, val)
	}
	return output, nil
}

// searchFilterHelper helps with the composition of FHIR REST search and filter requests.
//
// - the `resourceName` is a FHIR resource name e.g "Patient", "Appointment" etc
// - the `path` is a resource sub-path e.g "_search". If there is no sub-path, send a blank string
// - `params` should contain the filter parameters e.g
//
//    params := url.Values{}
//    params.Add("_content", search)
func (s Service) searchFilterHelper(
	ctx context.Context, resourceName string, path string, params url.Values,
) ([]map[string]interface{}, error) {

	s.checkPreconditions()

	bs, err := s.clinicalRepository.POSTRequest(resourceName, path, params, nil)
	if err != nil {
		log.Errorf("unable to search: %v", err)
		return nil, fmt.Errorf("unable to search: %v", err)
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		log.Errorf("%s could not be found with search params %v: %s", resourceName, params, err)
		return nil, fmt.Errorf(
			"%s could not be found with search params %v: %s", resourceName, params, err)
	}

	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
	for _, k := range mandatoryKeys {
		_, found := respMap[k]
		if !found {
			log.Errorf("Search: mandatory search result key %s not found", k)
			return nil, fmt.Errorf("server error: mandatory search result key %s not found", k)
		}
	}
	resourceType, ok := respMap["resourceType"].(string)
	if !ok {
		log.Errorf("Search: the resourceType is not a string")
		return nil, fmt.Errorf("server error: the resourceType is not a string")
	}
	if resourceType != "Bundle" {
		log.Errorf("Search: the resourceType value is not 'Bundle' as expected")
		return nil, fmt.Errorf(
			"server error: the resourceType value is not 'Bundle' as expected")
	}

	resultType, ok := respMap["type"].(string)
	if !ok {
		log.Errorf("Search: the search result type value is not a string")
		return nil, fmt.Errorf("server error: the search result type value is not a string")
	}
	if resultType != "searchset" {
		log.Errorf("Search: the type value is not 'searchset' as expected")
		return nil, fmt.Errorf("server error: the type value is not 'searchset' as expected")
	}

	respEntries := respMap["entry"]
	if respEntries == nil {
		return []map[string]interface{}{}, nil
	}
	entries, ok := respEntries.([]interface{})
	if !ok {
		log.Errorf("Search: entries is not a list of maps, it is: %T", respEntries)
		return nil, fmt.Errorf(
			"server error: entries is not a list of maps, it is: %T", respEntries)
	}

	results := []map[string]interface{}{}
	for _, en := range entries {
		entry, ok := en.(map[string]interface{})
		if !ok {
			log.Errorf("Search: expected each entry to be map, they are %T instead", en)
			return nil, fmt.Errorf(
				"server error: expected each entry to be map, they are %T instead", en)
		}
		expectedKeys := []string{"fullUrl", "resource", "search"}
		for _, k := range expectedKeys {
			_, found := entry[k]
			if !found {
				log.Errorf("Search: FHIR search entry does not have key '%s'", k)
				return nil, fmt.Errorf("server error: FHIR search entry does not have key '%s'", k)
			}
		}

		resource, ok := entry["resource"].(map[string]interface{})
		if !ok {
			{
				log.Errorf("Search: result entry %#v is not a map", entry["resource"])
				return nil, fmt.Errorf("server error: result entry %#v is not a map", entry["resource"])
			}
		}
		results = append(results, resource)
	}
	return results, nil
}

// ProblemSummary returns a short list of the patient's active and confirmed
// problems (by name).
func (s Service) ProblemSummary(
	ctx context.Context, patientID string) ([]string, error) {
	s.checkPreconditions()

	params := map[string]interface{}{
		"clinical-status":     "active",
		"verification-status": "confirmed",
		"category":            "problem-list-item",
		"subject":             fmt.Sprintf("Patient/%s", patientID),
	}
	results, err := s.SearchFHIRCondition(ctx, params)
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

// VisitSummary returns a narrative friendly display of the data that has
// been associated with a single visit
func (s Service) VisitSummary(
	ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
	s.checkPreconditions()
	encounterPayload, err := s.GetFHIREncounter(ctx, encounterID)
	if err != nil {
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
			conn, err := s.SearchFHIRAllergyIntolerance(ctx, patientFilterParams)
			if err != nil {
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := base.StructToMap(edge.Node)
				if err != nil {
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "Encounter":
			conn, err := s.SearchFHIREncounter(ctx, encounterInstanceFilterParams)
			if err != nil {
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := base.StructToMap(edge.Node)
				if err != nil {
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "Condition":
			conn, err := s.SearchFHIRCondition(ctx, encounterFilterParams)
			if err != nil {
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := base.StructToMap(edge.Node)
				if err != nil {
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "Observation":
			conn, err := s.SearchFHIRObservation(ctx, encounterFilterParams)
			if err != nil {
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := base.StructToMap(edge.Node)
				if err != nil {
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "Composition":
			conn, err := s.SearchFHIRComposition(ctx, encounterFilterParams)
			if err != nil {
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := base.StructToMap(edge.Node)
				if err != nil {
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "MedicationRequest":
			conn, err := s.SearchFHIRMedicationRequest(ctx, encounterFilterParams)
			if err != nil {
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := base.StructToMap(edge.Node)
				if err != nil {
					return nil, fmt.Errorf("%s edge struct to map error: %w", resourceName, err)
				}
				if rMap == nil {
					continue
				}
				nodes[resourceName] = append(nodes[resourceName], rMap)
			}
		case "ServiceRequest":
			conn, err := s.SearchFHIRServiceRequest(ctx, encounterFilterParams)
			if err != nil {
				return nil, fmt.Errorf("%s search error: %w", resourceName, err)
			}
			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}
				rMap, err := base.StructToMap(edge.Node)
				if err != nil {
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
func (s Service) PatientTimelineWithCount(
	ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	episode, _, err := s.getTimelineEpisode(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	encounterSearchParams := map[string]interface{}{
		"patient": *episode.Patient.Reference,
		"sort":    "-date", // reverse chronological order
		"_count":  strconv.Itoa(count),
	}
	return s.getTimelineVisitSummaries(ctx, encounterSearchParams, count)
}

func (s Service) getTimelineEpisode(ctx context.Context, episodeID string) (*FHIREpisodeOfCare, string, error) {
	s.checkPreconditions()
	episodePayload, err := s.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}
	episode := episodePayload.Resource
	activeEpisodeStatus := EpisodeOfCareStatusEnumActive
	if episode.Patient == nil || episode.Patient.Reference == nil {
		return nil, "", fmt.Errorf("the episode with ID %s has no patient reference", episodeID)
	}
	if episodePayload.Resource.Status.String() != activeEpisodeStatus.String() {
		return nil, "", fmt.Errorf("the episode with ID %s is not active", episodeID)
	}
	if episode.Type == nil {
		return nil, "", fmt.Errorf("the episode with ID %s has a nil type", episodeID)
	}
	if len(episode.Type) != 1 {
		return nil, "", fmt.Errorf("expected the episode type to have just one entry")
	}
	accessLevel := episode.Type[0].Text
	if accessLevel != "FULL_ACCESS" && accessLevel != "PROFILE_AND_RECENT_VISITS_ACCESS" {
		return nil, "", fmt.Errorf("unknown episode access level: %s", accessLevel)
	}
	return episode, accessLevel, nil
}

func (s Service) getTimelineVisitSummaries(
	ctx context.Context,
	encounterSearchParams map[string]interface{},
	count int,
) ([]map[string]interface{}, error) {
	encounterConn, err := s.SearchFHIREncounter(ctx, encounterSearchParams)
	if err != nil {
		return nil, fmt.Errorf("unable to search for encounters for the timeline: %w", err)
	}
	visitSummaries := []map[string]interface{}{}
	if encounterConn == nil || encounterConn.Edges == nil {
		return visitSummaries, nil
	}
	for _, edge := range encounterConn.Edges {
		if edge.Node == nil || edge.Node.ID == nil {
			continue
		}
		summary, err := s.VisitSummary(ctx, *edge.Node.ID, count)
		if err != nil {
			return nil, err
		}
		if summary != nil {
			visitSummaries = append(visitSummaries, summary)
		}
	}
	return visitSummaries, nil
}

// CreateEpisodeOfCare is the final common pathway for creation of episodes of
// care.
func (s Service) CreateEpisodeOfCare(
	ctx context.Context,
	ep FHIREpisodeOfCare,
) (*EpisodeOfCarePayload, error) {
	s.checkPreconditions()
	payload, err := base.StructToMap(ep)
	if err != nil {
		return nil, fmt.Errorf("unable to turn episode of care input into a map: %v", err)
	}

	// search for the episode of care before creating new one.
	episodeOfCareSearchParams := map[string]interface{}{
		"patient":      fmt.Sprintf(*ep.Patient.Reference),
		"status":       string(EpisodeOfCareStatusActive),
		"organization": ep.ManagingOrganization.Reference,
		"_sort":        "date",
		"_count":       "1",
	}

	episodeOfCarePayload, err := s.SearchFHIREpisodeOfCare(ctx, episodeOfCareSearchParams)

	if err == nil {
		if len(episodeOfCarePayload.Edges) == 1 {
			episodeOfCare := *episodeOfCarePayload.Edges[0].Node
			encounters, err := s.Encounters(ctx, *episodeOfCare.Patient.Reference, nil)
			if err == nil {
				output := &EpisodeOfCarePayload{
					EpisodeOfCare: &episodeOfCare,
					TotalVisits:   len(encounters),
				}
				return output, nil
			}
		}

		cloudhealthService := cloudhealth.NewService()
		data, err := cloudhealthService.CreateFHIRResource("EpisodeOfCare", payload)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to create episode of care resource: %v", err)
		}
		episode := &FHIREpisodeOfCare{}
		err = json.Unmarshal(data, episode)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to unmarshal episode of care response JSON: data: %v\n, error: %v",
				string(data), err)
		}

		encounters, err := s.Encounters(ctx, *episode.Patient.Reference, nil)
		if err != nil {
			return nil, fmt.Errorf(
				"unable to get encounters for episode %s: %v",
				*episode.ID, err,
			)
		}
		output := &EpisodeOfCarePayload{
			EpisodeOfCare: episode,
			TotalVisits:   len(encounters),
		}
		return output, nil

	}

	return nil, fmt.Errorf("unable to get patients episodes of care: %v", err)
}

// Encounters returns encounters that belong to the indicated patient.
//
// The patientReference should be a [string] in the format "Patient/<patient resource ID>".
func (s Service) Encounters(
	ctx context.Context,
	patientReference string,
	status *EncounterStatusEnum,
) ([]*FHIREncounter, error) {
	s.checkPreconditions()

	searchParams := url.Values{}
	if status != nil {
		searchParams.Add("status:exact", status.String())
	}
	searchParams.Add("patient", patientReference)
	cloudhealthService := cloudhealth.NewService()
	bs, err := cloudhealthService.POSTRequest("Encounter", "_search", searchParams, nil)
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

	output := []*FHIREncounter{}
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
		var encounter FHIREncounter
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

// StartEpisodeByOtp starts a patient OTP verified episode
func (s Service) StartEpisodeByOtp(
	ctx context.Context, input OTPEpisodeCreationInput) (*EpisodeOfCarePayload, error) {
	s.checkPreconditions()
	validPhone, err := base.ValidateMSISDN(
		input.Msisdn,
		input.Otp,
		false,
		s.firestoreClient,
	)
	if err != nil {
		log.Printf(
			"invalid phone: \nPhone: %s\nOTP: %s\n", input.Msisdn, input.Otp)
		return nil, fmt.Errorf(
			"invalid phone: got %s when validating %s", err, input.Msisdn)
	}
	organizationID, err := s.GetORCreateOrganization(ctx, input.ProviderCode)
	if err != nil {
		return nil, fmt.Errorf(
			"internal server error in retrieving service provider : %v", err)
	}
	ep := ComposeOneHealthEpisodeOfCare(
		validPhone, input.FullAccess, *organizationID, input.ProviderCode,
		input.PatientID)
	return s.CreateEpisodeOfCare(ctx, ep)
}

// UpgradeEpisode starts a patient OTP verified episode
func (s Service) UpgradeEpisode(
	ctx context.Context, input OTPEpisodeUpgradeInput) (*EpisodeOfCarePayload, error) {
	s.checkPreconditions()

	// retrieve and validate the episode
	episode, err := s.GetActiveEpisode(ctx, input.EpisodeID)
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
	encounters, err := s.Encounters(ctx, *episode.Patient.Reference, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get encounters for episode %s: %v",
			*episode.ID, err,
		)
	}

	// if it already has the correct access level, return early
	if episodeTypes[0].Text == fullAccessLevel {
		return &EpisodeOfCarePayload{
			EpisodeOfCare: episode,
			TotalVisits:   len(encounters),
		}, nil
	}

	// validate the MSISDN and OTP
	_, err = base.ValidateMSISDN(
		input.Msisdn,
		input.Otp,
		false,
		s.firestoreClient,
	)
	if err != nil {
		log.Printf(
			"invalid phone: \nPhone: %s\nOTP: %s\n", input.Msisdn, input.Otp)
		return nil, fmt.Errorf(
			"invalid phone: got %s when validating %s", err, input.Msisdn)
	}

	// patch the episode status
	episode.Type = []*FHIRCodeableConcept{
		{Text: fullAccessLevel},
	}
	payload, err := base.StructToMap(episode)
	if err != nil {
		return nil, fmt.Errorf("unable to turn episode of care input into a map: %v", err)
	}

	cloudhealthService := cloudhealth.NewService()
	_, err = cloudhealthService.UpdateFHIRResource(
		"EpisodeOfCare", *episode.ID, payload)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to update episode of care resource: %v", err)
	}
	return &EpisodeOfCarePayload{
		EpisodeOfCare: episode,
		TotalVisits:   len(encounters),
	}, nil
}

// StartEpisodeByBreakGlass starts an emergency episode
func (s Service) StartEpisodeByBreakGlass(
	ctx context.Context, input BreakGlassEpisodeCreationInput) (*EpisodeOfCarePayload, error) {
	s.checkPreconditions()
	validPhone, err := base.ValidateMSISDN(
		input.Msisdn,
		input.Otp,
		false,
		s.firestoreClient,
	)
	if err != nil {
		log.Printf(
			"invalid phone: \nPhone: %s\nOTP: %s\n", input.Msisdn, input.Otp)
		return nil, fmt.Errorf("invalid phone number/OTP: %w", err)
	}
	_, err = base.SaveDataToFirestore(
		s.firestoreClient, s.getBreakGlassCollectionName(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to log break glass operation: %v", err)
	}

	// alert patient
	err = s.sendAlertToPatient(ctx, validPhone, input.PatientID)
	if err != nil {
		log.Printf("failed to send alert message during StartEpisodeByBreakGlass login: %s", err)
	}

	// alert next-of-kin
	err = s.sendAlertToNextOfKin(ctx, input.PatientID)
	if err != nil {
		log.Printf("failed to send alert message to next of kin during StartEpisodeByBreakGlass login: %s", err)
	}

	// alert admin
	pp, err := s.FindPatientByID(ctx, input.PatientID)
	if err == nil {
		patientName := pp.PatientRecord.Name[0].Text
		err = s.sendAlertToAdmin(patientName, validPhone)
		if err != nil {
			log.Printf("failed to send alert message to admin during StartEpisodeByBreakGlass login: %s", err)
		}
	}
	organizationID, err := s.GetORCreateOrganization(ctx, input.ProviderCode)
	if err != nil {
		return nil, fmt.Errorf(
			"internal server error in retrieving service provider : %v", err)
	}
	ep := ComposeOneHealthEpisodeOfCare(
		validPhone, input.FullAccess, *organizationID, input.ProviderCode, input.PatientID)
	return s.CreateEpisodeOfCare(ctx, ep)
}

// GetOrganization retrieves an organization given its code
func (s Service) GetOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	s.checkPreconditions()
	searchParam := map[string]interface{}{
		"identifier": providerSladeCode,
	}
	organization, err := s.SearchFHIROrganization(ctx, searchParam)
	if err != nil {
		return nil, err
	}
	if organization.Edges == nil {
		return nil, nil
	}
	ORGID := organization.Edges[0].Node.ID
	return ORGID, nil
}

// CreateOrganization creates an organization given ist provider code
func (s Service) CreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	s.checkPreconditions()
	identifier := []*FHIRIdentifierInput{
		{
			Use:   "official",
			Value: providerSladeCode,
		},
	}
	organizationInput := FHIROrganizationInput{
		Identifier: identifier,
		Name:       &providerSladeCode,
	}
	createdOrganization, err := s.CreateFHIROrganization(ctx, organizationInput)
	if err != nil {
		return nil, err
	}
	organisationID := createdOrganization.Resource.ID
	return organisationID, nil
}

// GetORCreateOrganization retrieve an organisation via its code if not found create a new one.
func (s Service) GetORCreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	s.checkPreconditions()
	retrievedOrg, err := s.GetOrganization(ctx, providerSladeCode)
	if err != nil {
		return nil, fmt.Errorf(
			"internal server error in getting organisation : %v", err)
	}
	if retrievedOrg == nil {
		createdOrg, err := s.CreateOrganization(ctx, providerSladeCode)
		if err != nil {
			return nil, fmt.Errorf(
				"internal server error in creating organisation : %v", err)
		}
		return createdOrg, nil
	}
	return retrievedOrg, nil
}

// OpenOrganizationEpisodes return all organization specific open episodes
func (s Service) OpenOrganizationEpisodes(
	ctx context.Context, providerSladeCode string) ([]*FHIREpisodeOfCare, error) {
	organizationID, err := s.GetORCreateOrganization(ctx, providerSladeCode)
	if err != nil {
		return nil, fmt.Errorf(
			"internal server error in retrieving service provider : %v", err)
	}
	organizationReference := fmt.Sprintf("Organization/%s", *organizationID)
	searchParams := url.Values{}
	searchParams.Add("status", EpisodeOfCareStatusActive.String())
	searchParams.Add("organization", organizationReference)
	return s.SearchEpisodesByParam(ctx, searchParams)
}

// SearchEpisodeEncounter returns all encounters in a visit
func (s Service) SearchEpisodeEncounter(
	ctx context.Context,
	episodeReference string,
) (*FHIREncounterRelayConnection, error) {
	s.checkPreconditions()

	episodeRef := fmt.Sprintf("Episode/%s", episodeReference)
	encounterFilterParams := map[string]interface{}{
		"episodeOfCare": episodeRef,
		"status":        "in_progress",
	}
	encounterConn, err := s.SearchFHIREncounter(ctx, encounterFilterParams)

	if err != nil {
		return nil, fmt.Errorf("unable to search encounter: %w", err)
	}

	return encounterConn, nil
}

// StartEncounter starts an encounter within an episode of care
func (s Service) StartEncounter(
	ctx context.Context, episodeID string) (string, error) {
	s.checkPreconditions()

	episodePayload, err := s.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return "", fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}
	activeEpisodeStatus := EpisodeOfCareStatusEnumActive
	activeEncounterStatus := EncounterStatusEnumInProgress
	if episodePayload.Resource.Status.String() != activeEpisodeStatus.String() {
		return "", fmt.Errorf("an encounter can only be started for an active episode")
	}
	episodeRef := fmt.Sprintf("EpisodeOfCare/%s", *episodePayload.Resource.ID)

	now := time.Now()
	startTime := base.DateTime(now.Format("2006-01-02T15:04:05+03:00"))

	encounterClassCode := base.Code("AMB")
	encounterClassSystem := base.URI("http://terminology.hl7.org/CodeSystem/v3-ActCode")
	encounterClassVersion := "2018-08-12"
	encounterClassDisplay := "ambulatory"
	encounterClassUserSelected := false

	encounterInput := FHIREncounterInput{
		Status: activeEncounterStatus,
		Class: FHIRCodingInput{
			System:       &encounterClassSystem,
			Version:      &encounterClassVersion,
			Code:         encounterClassCode,
			Display:      encounterClassDisplay,
			UserSelected: &encounterClassUserSelected,
		},
		Subject: &FHIRReferenceInput{
			Reference: episodePayload.Resource.Patient.Reference,
			Display:   episodePayload.Resource.Patient.Display,
			Type:      episodePayload.Resource.Patient.Type,
		},
		EpisodeOfCare: []*FHIRReferenceInput{
			{
				Reference: &episodeRef,
			},
		},
		ServiceProvider: &FHIRReferenceInput{
			Display: episodePayload.Resource.ManagingOrganization.Display,
			Type:    episodePayload.Resource.ManagingOrganization.Type,
		},
		Period: &FHIRPeriodInput{
			Start: startTime,
		},
	}
	encPl, err := s.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		return "", fmt.Errorf("unable to start encounter: %w", err)
	}
	return *encPl.Resource.ID, nil
}

// EndEncounter ends an encounter
func (s Service) EndEncounter(
	ctx context.Context, encounterID string) (bool, error) {
	s.checkPreconditions()

	resourceType := "Encounter"
	encounterPayload, err := s.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return false, fmt.Errorf("unable to get encounter with ID %s: %w", encounterID, err)
	}
	updatedStatus := EncounterStatusEnumFinished
	encounterPayload.Resource.Status = updatedStatus

	// workaround for odd date comparison behavior on the Google Cloud Healthcare API
	// the end time must be at least 24 hours after the start time
	// so: if the time now is less than 24 hours after start, set the end to be
	// 24 hours after the start of the visit. If the time now is more than 24 hours
	// after the start, use the current time as the end of the visit
	end := time.Now().Add(time.Hour * 24)
	endTime := base.DateTime(end.Format(timeFormatStr))
	encounterPayload.Resource.Period.End = endTime

	payload, err := base.StructToMap(encounterPayload.Resource)
	if err != nil {
		return false, fmt.Errorf("unable to turn the updated episode of care into a map: %v", err)
	}
	cloudhealthService := cloudhealth.NewService()
	_, err = cloudhealthService.UpdateFHIRResource(resourceType, encounterID, payload)
	if err != nil {
		return false, fmt.Errorf("unable to create/update %s resource: %w", resourceType, err)
	}
	return true, nil
}

// EndEpisode ends an episode of care by patching it's status to "finished"
func (s Service) EndEpisode(
	ctx context.Context, episodeID string) (bool, error) {
	s.checkPreconditions()
	resourceType := "EpisodeOfCare"
	episodePayload, err := s.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return false, fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}
	startTime := base.DateTime(time.Now().Format(timeFormatStr))
	if episodePayload.Resource.Period != nil {
		startTime = episodePayload.Resource.Period.Start
	}

	// Close all encounters in this visit
	encounterConn, err := s.SearchEpisodeEncounter(ctx, episodeID)
	if err != nil {
		return false, fmt.Errorf("unable to search episode encounter %w", err)
	}
	for _, edge := range encounterConn.Edges {
		_, err = s.EndEncounter(ctx, *edge.Node.ID)
		if err != nil {
			// !TODO track or log the ones we have failed to close
			continue
		}
	}
	// // workaround for odd date comparison behavior on the Google Cloud Healthcare API
	// the end time must be at least 24 hours after the start time
	// so: if the time now is less than 24 hours after start, set the end to be
	// 24 hours after the start of the visit. If the time now is more than 24 hours
	// after the start, use the current time as the end of the visit
	end := time.Now().Add(time.Hour * 24)
	endTime := base.DateTime(end.Format(timeFormatStr))

	updatedStatus := EpisodeOfCareStatusEnumFinished
	episodePayload.Resource.Status = &updatedStatus
	episodePayload.Resource.Period.Start = startTime
	episodePayload.Resource.Period.End = endTime
	payload, err := base.StructToMap(episodePayload.Resource)
	if err != nil {
		return false, fmt.Errorf("unable to turn the updated episode of care into a map: %v", err)
	}

	cloudhealthService := cloudhealth.NewService()
	_, err = cloudhealthService.UpdateFHIRResource(resourceType, episodeID, payload)
	if err != nil {
		return false, fmt.Errorf("unable to create/update %s resource: %w", resourceType, err)
	}
	patientReference := *episodePayload.Resource.Patient.Reference
	patientID, err := GetPatientIDFromEpisode(patientReference)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve patientID from episodeOfCare: %w", err)
	}
	err = s.sendEpisodeEndAlert(ctx, patientID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// sendEpisodeEndAlert send an alert to the patient informing a visit has been closed
func (s Service) sendEpisodeEndAlert(ctx context.Context, patientID string) error {
	patientPayload, err := s.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return err
	}
	patientName := patientPayload.Resource.Name
	patientContacts := patientPayload.Resource.Telecom

	for _, contact := range patientContacts {
		if *contact.System == ContactPointSystemEnumPhone {

			message := composeAlertMessage(patientName)
			phone := *contact.Value

			if phone == "" {
				continue
			}
			err := base.SendSMS(
				[]string{*contact.Value}, message, *s.sms, *s.twilio)
			if err != nil {
				return err
			}

			return nil

		}
	}
	return nil
}

// CreatePatient creates or updates a patient record on FHIR
func (s Service) CreatePatient(
	ctx context.Context, input PatientInput) (*PatientPayload, error) {
	s.checkPreconditions()

	// set the record ID if not set
	if input.ID == nil {
		newID := uuid.New().String()
		input.ID = &newID
	}

	// set the default communication preferences if not set
	if input.Communication == nil {
		input.Communication = DefaultCommunication()
	}

	// set or add the default record identifier
	if input.Identifier == nil {
		input.Identifier = []*FHIRIdentifierInput{DefaultIdentifier()}
	}
	if input.Identifier != nil {
		input.Identifier = append(input.Identifier, DefaultIdentifier())
	}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn patient input into a map: %v", err)
	}

	cloudhealthService := cloudhealth.NewService()
	data, err := cloudhealthService.CreateFHIRResource("Patient", payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update patient resource: %v", err)
	}
	patient := &FHIRPatient{}
	err = json.Unmarshal(data, patient)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal patient response JSON: data: %v\n, error: %v",
			string(data), err)
	}
	output := &PatientPayload{
		PatientRecord:   patient,
		HasOpenEpisodes: false, // the patient is newly created so we can safely assume this
		OpenEpisodes:    []*FHIREpisodeOfCare{},
	}
	return output, nil
}

// DeletePatient is used to inactivate / retire a health record.
// It is used to implement "forget me".
//
// The patient's identifiers and name are removed. Their date of birth is set to the
// start of the year. They are then marked as inactive.
//
// This is done because there may already be patient records linked to this patient.
func (s Service) DeletePatient(
	ctx context.Context, input RetirePatientInput) (bool, error) {
	s.checkPreconditions()

	patientPayload, err := s.FindPatientByID(ctx, input.ID)
	if err != nil {
		return false, fmt.Errorf("no patient with ID %s: %w", input.ID, err)
	}

	dob := patientPayload.PatientRecord.BirthDate
	dob.Month = 1 // anonymize
	dob.Day = 1   // anonymize

	payload := []map[string]interface{}{
		{
			"op":    "replace",
			"path":  "/active",
			"value": false,
		},
		{
			"op":    "replace",
			"path":  "/identifier",
			"value": nil,
		},
		{
			"op":    "replace",
			"path":  "/name",
			"value": nil,
		},
		{
			"op":    "replace",
			"path":  "/birthDate",
			"value": dob,
		},
	}

	cloudhealthService := cloudhealth.NewService()
	resp, err := cloudhealthService.PatchFHIRResource(
		"Patient",
		*patientPayload.PatientRecord.ID,
		payload,
	)
	if err != nil {
		return false, fmt.Errorf(
			"unable to anonymize (forget) patient, response %s, error: %v",
			string(resp), err,
		)
	}
	return true, nil
}

// FindPatientByID retrieves a single patient by their ID
func (s Service) FindPatientByID(
	ctx context.Context, id string) (*PatientPayload, error) {
	s.checkPreconditions()
	cloudhealthService := cloudhealth.NewService()
	data, err := cloudhealthService.GetFHIRResource("Patient", id)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get patient with ID %s, err: %v", id, err)
	}
	var patient FHIRPatient
	err = json.Unmarshal(data, &patient)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal patient data from JSON, err: %v", err)
	}
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	openEpisodes, err := s.OpenEpisodes(ctx, patientReference)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get open episodes for %s, err: %v", patientReference, err)
	}
	return &PatientPayload{
		PatientRecord:   &patient,
		OpenEpisodes:    openEpisodes,
		HasOpenEpisodes: len(openEpisodes) > 0,
	}, nil
}

// PatientSearch searches for a patient by identifiers and names
func (s Service) PatientSearch(
	ctx context.Context, search string) (*PatientConnection, error) {
	s.checkPreconditions()

	params := url.Values{}
	params.Add("_content", search) // entire doc
	cloudhealthService := cloudhealth.NewService()
	bs, err := cloudhealthService.POSTRequest(
		"Patient", "_search", params, nil)
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
		return &PatientConnection{
			Edges:    []*PatientEdge{},
			PageInfo: &base.PageInfo{},
		}, nil
	}
	entries, ok := respEntries.([]interface{})
	if !ok {
		log.Errorf("Search: entries is not a list of maps, it is: %T", respEntries)
		return nil, fmt.Errorf(notFoundWithSearchParams)
	}

	output := PatientConnection{}
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

		resource = s.birthdateMapper(resource)
		resource = s.identifierMapper(resource)
		resource = s.nameMapper(resource)
		resource = s.telecomMapper(resource)
		resource = s.addressMapper(resource)
		resource = s.photoMapper(resource)
		resource = s.contactMapper(resource)

		var patient FHIRPatient

		err := mapstructure.Decode(resource, &patient)
		if err != nil {
			log.Errorf("unable to map decode resource: %v", err)
			return nil, fmt.Errorf(internalError)
		}

		hasOpenEpisodes, err := s.HasOpenEpisode(ctx, patient)
		if err != nil {
			log.Errorf("error while checking if hasOpenEpisodes: %v", err)
			return nil, fmt.Errorf(internalError)
		}
		output.Edges = append(output.Edges, &PatientEdge{
			Node:            &patient,
			HasOpenEpisodes: hasOpenEpisodes,
		})
	}
	return &output, nil
}

// UpdatePatient patches a patient record with fresh data.
// It updates elements that are set and ignores the ones that are nil.
func (s Service) UpdatePatient(
	ctx context.Context, input SimplePatientRegistrationInput) (*PatientPayload, error) {
	s.checkPreconditions()

	op := "add" // this method replaces data at the indicated paths

	if input.ID == "" {
		return nil, fmt.Errorf("can't update with blank ID")
	}

	patientInput, err := s.SimplePatientRegistrationInputToPatientInput(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("can't compose patient input: %v", err)
	}

	patientPayload, err := s.FindPatientByID(ctx, input.ID)
	if err != nil {
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
	if patientInput.Active != *patientPayload.PatientRecord.Active {
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
	if patientInput.Gender != patientPayload.PatientRecord.Gender.String() {
		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/gender"
		patch["value"] = patientInput.Gender
		patches = append(patches, patch)
	}
	if patientInput.BirthDate != *patientPayload.PatientRecord.BirthDate {
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

	cloudhealthService := cloudhealth.NewService()
	data, err := cloudhealthService.PatchFHIRResource("Patient", input.ID, patches)
	if err != nil {
		return nil, fmt.Errorf("UpdatePatient: %v", err)
	}
	patient := FHIRPatient{}
	err = json.Unmarshal(data, &patient)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal patient response JSON: data: %v\n, error: %v",
			string(data), err)
	}
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	openEpisodes, err := s.OpenEpisodes(ctx, patientReference)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get open episodes for %s, err: %v", patientReference, err)
	}
	return &PatientPayload{
		PatientRecord:   &patient,
		OpenEpisodes:    openEpisodes,
		HasOpenEpisodes: len(openEpisodes) > 0,
	}, nil
}

// AddNextOfKin patches a patient with next of kin
func (s Service) AddNextOfKin(
	ctx context.Context, input SimpleNextOfKinInput) (*PatientPayload, error) {
	s.checkPreconditions()
	if input.PatientID == "" {
		return nil, fmt.Errorf("a patient ID must be specified")
	}

	_, err := s.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get patient with ID %s: %v", input.PatientID, err)
	}

	updatedContacts := []*FHIRPatientContactInput{}

	names := NameToHumanName(input.Names)
	if len(names) == 0 {
		return nil, fmt.Errorf("the contact must have a name")
	}

	contacts, err := ContactsToContactPointInput(
		input.PhoneNumbers, input.Emails, s.firestoreClient)
	if err != nil {
		return nil, fmt.Errorf("invalid contacts: %v", err)
	}

	if len(names) != 1 {
		return nil, fmt.Errorf("the contact should have one name")
	}

	addresses := PhysicalPostalAddressesToCombinedFHIRAddress(
		input.PhysicalAddresses,
		input.PostalAddresses,
	)
	userSelected := true
	relationshipSystem := base.URI(RelationshipSystem)
	relationshipVersion := RelationshipVersion
	gender := PatientContactGenderEnum(input.Gender)
	if !gender.IsValid() {
		return nil, fmt.Errorf(
			"'%s' is not a valid gender; valid values are %s",
			input.Gender,
			AllPatientContactGenderEnum,
		)
	}
	updatedContacts = append(updatedContacts, &FHIRPatientContactInput{
		Relationship: []*FHIRCodeableConceptInput{
			{
				Coding: []*FHIRCodingInput{
					{
						Display:      RelationshipTypeDisplay(input.Relationship),
						System:       &relationshipSystem,
						Version:      &relationshipVersion,
						Code:         base.Code(input.Relationship.String()),
						UserSelected: &userSelected,
					},
				},
				Text: RelationshipTypeDisplay(input.Relationship),
			},
		},
		Name:    names[0],
		Telecom: contacts,
		Address: addresses,
		Gender:  &gender,
		Period:  DefaultPeriodInput(),
	})
	patches := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/contact",
			"value": updatedContacts,
		},
	}
	cloudhealthService := cloudhealth.NewService()
	data, err := cloudhealthService.PatchFHIRResource(
		"Patient", input.PatientID, patches)
	if err != nil {
		return nil, fmt.Errorf("UpdatePatient: %v", err)
	}
	patient := FHIRPatient{}
	err = json.Unmarshal(data, &patient)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal patient response JSON: data: %v\n, error: %v",
			string(data), err)
	}
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	openEpisodes, err := s.OpenEpisodes(ctx, patientReference)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get open episodes for patient %s: %w", patientReference, err)
	}

	return &PatientPayload{
		PatientRecord:   &patient,
		OpenEpisodes:    openEpisodes,
		HasOpenEpisodes: len(openEpisodes) > 0,
	}, nil
}

// AddNhif patches a patient with NHIF details
func (s Service) AddNhif(
	ctx context.Context, input *SimpleNHIFInput) (*PatientPayload, error) {
	s.checkPreconditions()
	if input == nil {
		return nil, fmt.Errorf("AddNhif: nil input")
	}

	if input.PatientID == "" {
		return nil, fmt.Errorf("a patient ID must be specified")
	}

	patientPayload, err := s.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get patient with ID %s: %v", input.PatientID, err)
	}

	existingIdentifers := patientPayload.PatientRecord.Identifier
	updatedIdentifierInputs := []*FHIRIdentifierInput{}
	for _, existing := range existingIdentifers {
		updatedTypeCoding := []*FHIRCodingInput{}
		for _, coding := range existing.Type.Coding {
			updatedTypeCoding = append(updatedTypeCoding, &FHIRCodingInput{
				System:       coding.System,
				Version:      coding.Version,
				Code:         coding.Code,
				Display:      coding.Display,
				UserSelected: coding.UserSelected,
			})
		}
		updatedIdentifierInputs = append(updatedIdentifierInputs, &FHIRIdentifierInput{
			ID:  existing.ID,
			Use: existing.Use,
			Type: FHIRCodeableConceptInput{
				ID:     existing.Type.ID,
				Text:   existing.Type.Text,
				Coding: updatedTypeCoding,
			},
			System: existing.System,
			Value:  existing.Value,
			Period: &FHIRPeriodInput{
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
	cloudhealthService := cloudhealth.NewService()
	data, err := cloudhealthService.PatchFHIRResource(
		"Patient", input.PatientID, patches)
	if err != nil {
		return nil, fmt.Errorf("UpdatePatient: %v", err)
	}
	patient := FHIRPatient{}
	err = json.Unmarshal(data, &patient)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal patient response JSON: data: %v\n, error: %v",
			string(data), err)
	}
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	openEpisodes, err := s.OpenEpisodes(ctx, patientReference)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get open episodes for %s, err: %v", patientReference, err)
	}
	return &PatientPayload{
		PatientRecord:   &patient,
		OpenEpisodes:    openEpisodes,
		HasOpenEpisodes: len(openEpisodes) > 0,
	}, nil
}

// GetActiveEpisode returns any ACTIVE episode that has to the indicated ID
func (s Service) GetActiveEpisode(ctx context.Context, episodeID string) (*FHIREpisodeOfCare, error) {
	s.checkPreconditions()

	searchParams := url.Values{}
	searchParams.Add("status:exact", EpisodeOfCareStatusActive.String())
	searchParams.Add("_id", episodeID) // logical ID of the resource

	cloudhealthService := cloudhealth.NewService()
	bs, err := cloudhealthService.POSTRequest(
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
	var episode FHIREpisodeOfCare
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

// SearchEpisodesByParam search episodes by params
func (s Service) SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*FHIREpisodeOfCare, error) {
	s.checkPreconditions()

	cloudhealthService := cloudhealth.NewService()
	bs, err := cloudhealthService.POSTRequest(
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

	output := []*FHIREpisodeOfCare{}
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
		resStart := ParseDate(period["start"].(string))
		startDateAsMap := make(map[string]int)
		startDateAsMap["year"] = resStart.Year()
		startDateAsMap["month"] = int(resStart.Month())
		startDateAsMap["day"] = resStart.Day()
		period["start"] = resStart

		// parse period->end as map[string]interface{}
		resEnd := ParseDate(period["end"].(string))
		endDateAsMap := make(map[string]int)
		endDateAsMap["year"] = resEnd.Year()
		endDateAsMap["month"] = int(resEnd.Month())
		endDateAsMap["day"] = resEnd.Day()
		period["end"] = resEnd

		//update the original period resource
		resource["period"] = period

		var episode FHIREpisodeOfCare

		err := mapstructure.Decode(resource, &episode)
		if err != nil {
			log.Errorf("unable to map decode resource: %v", err)
			return nil, fmt.Errorf(internalError)
		}

		output = append(output, &episode)
	}
	return output, nil
}

// OpenEpisodes returns the IDs of a patient's open episodes
func (s Service) OpenEpisodes(
	ctx context.Context, patientReference string) ([]*FHIREpisodeOfCare, error) {
	searchParams := url.Values{}
	searchParams.Add("status:exact", EpisodeOfCareStatusActive.String())
	searchParams.Add("patient", patientReference)
	return s.SearchEpisodesByParam(ctx, searchParams)
}

// HasOpenEpisode determines if a patient has an open episode
func (s Service) HasOpenEpisode(
	ctx context.Context,
	patient FHIRPatient,
) (bool, error) {
	s.checkPreconditions()
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	episodes, err := s.OpenEpisodes(ctx, patientReference)
	if err != nil {
		return false, err
	}
	return len(episodes) > 0, nil
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
func (s Service) FindPatientsByMSISDN(
	ctx context.Context, msisdn string) (*PatientConnection, error) {
	s.checkPreconditions()

	search, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return nil, fmt.Errorf("can't normalize contact: %w", err)
	}
	return s.PatientSearch(ctx, search)
}

// RegisterPatient implements simple patient registration
func (s Service) RegisterPatient(
	ctx context.Context, input SimplePatientRegistrationInput) (*PatientPayload, error) {
	s.checkPreconditions()
	patientInput, err := s.SimplePatientRegistrationInputToPatientInput(ctx, input)
	if err != nil {
		return nil, err
	}
	log.Printf("patient registration input: %#v", patientInput)
	output, err := s.CreatePatient(ctx, *patientInput)
	if err != nil {
		return nil, fmt.Errorf("unable to create patient: %v", err)
	}
	for _, patientEmail := range input.Emails {
		err = s.SendPatientWelcomeEmail(ctx, patientEmail.Email)
		if err != nil {
			return nil, fmt.Errorf("unable to send welcome email: %w", err)
		}
	}
	return output, nil
}

func (s Service) getBreakGlassCollectionName() string {
	suffixed := base.SuffixCollection(BreakGlassCollectionName)
	return suffixed
}

// SimplePatientRegistrationInputToPatientInput transforms a patient input into
// a
func (s Service) SimplePatientRegistrationInputToPatientInput(
	ctx context.Context, input SimplePatientRegistrationInput) (*PatientInput, error) {
	s.checkPreconditions()

	contacts, err := ContactsToContactPointInput(
		input.PhoneNumbers, input.Emails, s.firestoreClient)
	if err != nil {
		return nil, fmt.Errorf("can't register patient with invalid contacts: %v", err)
	}

	ids, err := IDToIdentifier(
		input.IdentificationDocuments, input.PhoneNumbers)
	if err != nil {
		return nil, fmt.Errorf("can't register patient with invalid identifiers: %v", err)
	}

	photos, err := PhotosToAttachments(ctx, input.Photos)
	if err != nil {
		return nil, fmt.Errorf("can't process patient photos: %v", err)
	}

	// fullPatientInput is to be filled up by processing the simple patient input
	patientInput := PatientInput{
		BirthDate: input.BirthDate,
		Gender:    input.Gender,
		Active:    input.Active,
	}
	patientInput.Identifier = ids
	patientInput.Telecom = contacts
	patientInput.Name = NameToHumanName(input.Names)
	patientInput.Photo = photos
	patientInput.Address = PhysicalPostalAddressesToFHIRAddresses(
		input.PhysicalAddresses, input.PostalAddresses)
	patientInput.MaritalStatus = MaritalStatusEnumToCodeableConceptInput(
		input.MaritalStatus)
	patientInput.Communication = LanguagesToCommunicationInputs(input.Languages)
	return &patientInput, nil
}

// SendPatientWelcomeEmail will send a welcome email to the practitioner
func (s Service) SendPatientWelcomeEmail(ctx context.Context, emailaddress string) error {
	s.checkPreconditions()

	text := generatePatientWelcomeEmailTemplate()
	if !govalidator.IsEmail(emailaddress) {
		return nil
	}
	body := map[string]interface{}{
		"to":      []string{emailaddress},
		"text":    text,
		"subject": EmailWelcomeSubject,
	}

	resp, err := s.mailgun.MakeRequest(http.MethodPost, sendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send welcome email: %w", err)
	}
	if base.IsDebug() {
		b, _ := httputil.DumpResponse(resp, true)
		log.Println(string(b))
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got error status %s from email service", resp.Status)
	}

	return nil
}

// sendAlertToPatient to send notification to patient when break glass request is made
func (s Service) sendAlertToPatient(ctx context.Context, phoneNumber string, patientID string) error {
	patientPayload, err := s.FindPatientByID(ctx, patientID)
	if err != nil {
		return err
	}

	name := patientPayload.PatientRecord.Name[0].Given[0]
	if name == nil {
		return fmt.Errorf("nil patient name")
	}
	text := createAlertMessage(*name)

	err = base.SendSMS(
		[]string{phoneNumber}, text, *s.sms, *s.twilio)
	if err != nil {
		return err
	}
	return nil
}

// sendAlertToAdmin send email to admin notifying them of the access.
func (s Service) sendAlertToAdmin(patientName string, patientContact string) error {
	adminEmail, err := base.GetEnvVar(SavannahAdminEmail)
	if err != nil {
		return err
	}

	var writer bytes.Buffer
	t := template.Must(template.New("profile").Parse(adminEmailMessage))
	_ = t.Execute(&writer, struct {
		name   string
		number string
	}{
		name:   patientName,
		number: patientContact,
	})
	subject := "Breaking Glass Access notice"

	body := map[string]interface{}{
		"to":      []string{adminEmail},
		"text":    writer.String(),
		"subject": subject,
	}

	resp, err := s.mailgun.MakeRequest(http.MethodPost, sendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send Alert to admin email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("unable to send Alert to admin email : %v, with status code %v", err, resp.StatusCode)
		return fmt.Errorf("got error status %s from email service", resp.Status)
	}

	return nil
}

// SearchFHIRServiceRequest provides a search API for FHIRServiceRequest
func (s Service) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*FHIRServiceRequestRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "ServiceRequest"
	path := "_search"
	output := FHIRServiceRequestRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIRServiceRequest

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
		output.Edges = append(output.Edges, &FHIRServiceRequestRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRServiceRequest creates a FHIRServiceRequest instance
func (s Service) CreateFHIRServiceRequest(ctx context.Context, input FHIRServiceRequestInput) (*FHIRServiceRequestRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "ServiceRequest"
	resource := FHIRServiceRequest{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRServiceRequestRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIRServiceRequest deletes the FHIRServiceRequest identified by the supplied ID
func (s Service) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	resourceType := "ServiceRequest"
	resp, err := s.clinicalRepository.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

//sendAlertToNextOfKin send an alert message to the patient's next of kin.
func (s Service) sendAlertToNextOfKin(ctx context.Context, patientID string) error {
	patientPayload, err := s.FindPatientByID(ctx, patientID)
	if err != nil {
		return err
	}
	patientContacts := patientPayload.PatientRecord.Contact

	for _, patientRelation := range patientContacts {
		if patientRelation.Name == nil {
			continue
		}
		if len(patientRelation.Name.Given) == 0 {
			continue
		}
		for _, codeableConcept := range patientRelation.Relationship {
			for _, coding := range codeableConcept.Coding {
				if coding.Code == "N" {
					// this is the next of kin
					phoneNextOfKin := patientRelation.Telecom
					for _, number := range phoneNextOfKin {
						if number == nil {
							continue
						}
						if number.Value == nil {
							continue
						}
						if number.System.String() == ContactPointSystemPhone.String() {
							text := createAlertMessage(*patientRelation.Name.Given[0])
							err := base.SendSMS(
								[]string{*number.Value}, text, *s.sms, *s.twilio)

							if err != nil {
								return err
							}
							return nil
						}
					}
					break
				}
			}

		}
	}
	err = fmt.Errorf("failed to send message")
	return err
}

// SearchFHIRAllergyIntolerance provides a search API for FHIRAllergyIntolerance
func (s Service) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*FHIRAllergyIntoleranceRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "AllergyIntolerance"
	path := "_search"
	output := FHIRAllergyIntoleranceRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIRAllergyIntolerance

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
		output.Edges = append(output.Edges, &FHIRAllergyIntoleranceRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRAllergyIntolerance creates a FHIRAllergyIntolerance instance
func (s Service) CreateFHIRAllergyIntolerance(ctx context.Context, input FHIRAllergyIntoleranceInput) (*FHIRAllergyIntoleranceRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "AllergyIntolerance"
	resource := FHIRAllergyIntolerance{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRAllergyIntoleranceRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIRAllergyIntolerance updates a FHIRAllergyIntolerance instance
// The resource must have it's ID set.
func (s Service) UpdateFHIRAllergyIntolerance(ctx context.Context, input FHIRAllergyIntoleranceInput) (*FHIRAllergyIntoleranceRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "AllergyIntolerance"
	resource := FHIRAllergyIntolerance{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRAllergyIntoleranceRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// SearchFHIRComposition provides a search API for FHIRComposition
func (s Service) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*FHIRCompositionRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Composition"
	path := "_search"
	output := FHIRCompositionRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIRComposition

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
		output.Edges = append(output.Edges, &FHIRCompositionRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRComposition creates a FHIRComposition instance
func (s Service) CreateFHIRComposition(ctx context.Context, input FHIRCompositionInput) (*FHIRCompositionRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Composition"
	resource := FHIRComposition{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRCompositionRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIRComposition updates a FHIRComposition instance
// The resource must have it's ID set.
func (s Service) UpdateFHIRComposition(ctx context.Context, input FHIRCompositionInput) (*FHIRCompositionRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Composition"
	resource := FHIRComposition{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRCompositionRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIRComposition deletes the FHIRComposition identified by the supplied ID
func (s Service) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	resourceType := "Composition"
	resp, err := s.clinicalRepository.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// SearchFHIRCondition provides a search API for FHIRCondition
func (s Service) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*FHIRConditionRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Condition"
	path := "_search"
	output := FHIRConditionRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIRCondition

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
		output.Edges = append(output.Edges, &FHIRConditionRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRCondition creates a FHIRCondition instance
func (s Service) CreateFHIRCondition(ctx context.Context, input FHIRConditionInput) (*FHIRConditionRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Condition"
	resource := FHIRCondition{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRConditionRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIRCondition updates a FHIRCondition instance
// The resource must have it's ID set.
func (s Service) UpdateFHIRCondition(ctx context.Context, input FHIRConditionInput) (*FHIRConditionRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Condition"
	resource := FHIRCondition{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRConditionRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// GetFHIREncounter retrieves instances of FHIREncounter by ID
func (s Service) GetFHIREncounter(ctx context.Context, id string) (*FHIREncounterRelayPayload, error) {
	s.checkPreconditions()

	resourceType := "Encounter"
	var resource FHIREncounter

	data, err := s.clinicalRepository.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &FHIREncounterRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// SearchFHIREncounter provides a search API for FHIREncounter
func (s Service) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*FHIREncounterRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Encounter"
	path := "_search"
	output := FHIREncounterRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIREncounter

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
		output.Edges = append(output.Edges, &FHIREncounterRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIREncounter creates a FHIREncounter instance
func (s Service) CreateFHIREncounter(ctx context.Context, input FHIREncounterInput) (*FHIREncounterRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Encounter"
	resource := FHIREncounter{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIREncounterRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// GetFHIREpisodeOfCare retrieves instances of FHIREpisodeOfCare by ID
func (s Service) GetFHIREpisodeOfCare(ctx context.Context, id string) (*FHIREpisodeOfCareRelayPayload, error) {
	s.checkPreconditions()

	resourceType := "EpisodeOfCare"
	var resource FHIREpisodeOfCare

	data, err := s.clinicalRepository.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &FHIREpisodeOfCareRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (s Service) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*FHIREpisodeOfCareRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "EpisodeOfCare"
	path := "_search"
	output := FHIREpisodeOfCareRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIREpisodeOfCare

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
		output.Edges = append(output.Edges, &FHIREpisodeOfCareRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// SearchFHIRMedicationRequest provides a search API for FHIRMedicationRequest
func (s Service) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*FHIRMedicationRequestRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "MedicationRequest"
	path := "_search"
	output := FHIRMedicationRequestRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIRMedicationRequest

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
		output.Edges = append(output.Edges, &FHIRMedicationRequestRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRMedicationRequest creates a FHIRMedicationRequest instance
func (s Service) CreateFHIRMedicationRequest(ctx context.Context, input FHIRMedicationRequestInput) (*FHIRMedicationRequestRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "MedicationRequest"
	resource := FHIRMedicationRequest{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRMedicationRequestRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIRMedicationRequest updates a FHIRMedicationRequest instance
// The resource must have it's ID set.
func (s Service) UpdateFHIRMedicationRequest(ctx context.Context, input FHIRMedicationRequestInput) (*FHIRMedicationRequestRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "MedicationRequest"
	resource := FHIRMedicationRequest{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRMedicationRequestRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIRMedicationRequest deletes the FHIRMedicationRequest identified by the supplied ID
func (s Service) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	resourceType := "MedicationRequest"
	resp, err := s.clinicalRepository.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// SearchFHIRObservation provides a search API for FHIRObservation
func (s Service) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*FHIRObservationRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Observation"
	path := "_search"
	output := FHIRObservationRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIRObservation

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
		output.Edges = append(output.Edges, &FHIRObservationRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRObservation creates a FHIRObservation instance
func (s Service) CreateFHIRObservation(ctx context.Context, input FHIRObservationInput) (*FHIRObservationRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Observation"
	resource := FHIRObservation{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRObservationRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// SearchFHIROrganization provides a search API for FHIROrganization
func (s Service) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*FHIROrganizationRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Organization"
	path := "_search"
	output := FHIROrganizationRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIROrganization

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
		output.Edges = append(output.Edges, &FHIROrganizationRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIROrganization creates a FHIROrganization instance
func (s Service) CreateFHIROrganization(ctx context.Context, input FHIROrganizationInput) (*FHIROrganizationRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Organization"
	resource := FHIROrganization{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIROrganizationRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// GetFHIRPatient retrieves instances of FHIRPatient by ID
func (s Service) GetFHIRPatient(ctx context.Context, id string) (*FHIRPatientRelayPayload, error) {
	s.checkPreconditions()

	resourceType := "Patient"
	var resource FHIRPatient

	data, err := s.clinicalRepository.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	hasOpenEpisodes, err := s.HasOpenEpisode(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get open episodes for patient %#v: %w", resource, err)
	}
	payload := &FHIRPatientRelayPayload{
		Resource:        &resource,
		HasOpenEpisodes: hasOpenEpisodes,
	}
	return payload, nil
}

func (s Service) birthdateMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	parsedDate := ParseDate(resourceCopy["birthDate"].(string))

	dateMap := make(map[string]interface{})

	dateMap["year"] = parsedDate.Year()
	dateMap["month"] = parsedDate.Month()
	dateMap["day"] = parsedDate.Day()

	resourceCopy["birthDate"] = dateMap

	return resourceCopy

}

func (s Service) periodMapper(period map[string]interface{}) map[string]interface{} {

	periodCopy := period

	parsedStartDate := ParseDate(periodCopy["start"].(string))

	periodCopy["start"] = parsedStartDate

	parsedEndDate := ParseDate(periodCopy["end"].(string))

	periodCopy["end"] = parsedEndDate

	return periodCopy
}

func (s Service) identifierMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	if _, ok := resource["identifier"]; ok {

		newIdentifiers := []map[string]interface{}{}

		for _, identifier := range resource["identifier"].([]interface{}) {

			identifier := identifier.(map[string]interface{})

			if _, ok := identifier["period"]; ok {

				period := identifier["period"].(map[string]interface{})
				newPeriod := s.periodMapper(period)

				identifier["period"] = newPeriod
			}

			newIdentifiers = append(newIdentifiers, identifier)
		}

		resourceCopy["identifier"] = newIdentifiers
	}

	return resourceCopy
}

func (s Service) nameMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newNames := []map[string]interface{}{}

	if _, ok := resource["name"]; ok {

		for _, name := range resource["name"].([]interface{}) {

			name := name.(map[string]interface{})

			if _, ok := name["period"]; ok {

				period := name["period"].(map[string]interface{})
				newPeriod := s.periodMapper(period)

				name["period"] = newPeriod
			}

			newNames = append(newNames, name)
		}

	}

	resourceCopy["name"] = newNames

	return resourceCopy
}

func (s Service) telecomMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newTelecoms := []map[string]interface{}{}

	if _, ok := resource["telecom"]; ok {

		for _, telecom := range resource["telecom"].([]interface{}) {

			telecom := telecom.(map[string]interface{})

			if _, ok := telecom["period"]; ok {

				period := telecom["period"].(map[string]interface{})
				newPeriod := s.periodMapper(period)

				telecom["period"] = newPeriod
			}

			newTelecoms = append(newTelecoms, telecom)
		}

	}

	resourceCopy["telecom"] = newTelecoms

	return resourceCopy
}

func (s Service) addressMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newAddresses := []map[string]interface{}{}

	if _, ok := resource["address"]; ok {

		for _, address := range resource["address"].([]interface{}) {

			address := address.(map[string]interface{})

			if _, ok := address["period"]; ok {

				period := address["period"].(map[string]interface{})
				newPeriod := s.periodMapper(period)

				address["period"] = newPeriod
			}

			newAddresses = append(newAddresses, address)
		}
	}

	resourceCopy["address"] = newAddresses

	return resourceCopy
}

func (s Service) photoMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newPhotos := []map[string]interface{}{}

	if _, ok := resource["photo"]; ok {

		for _, photo := range resource["photo"].([]interface{}) {

			photo := photo.(map[string]interface{})

			parsedDate := ParseDate(photo["creation"].(string))

			photo["creation"] = parsedDate

			newPhotos = append(newPhotos, photo)
		}
	}

	resourceCopy["photo"] = newPhotos

	return resourceCopy
}

func (s Service) contactMapper(resource map[string]interface{}) map[string]interface{} {

	resourceCopy := resource

	newContacts := []map[string]interface{}{}

	if _, ok := resource["contact"]; ok {

		for _, contact := range resource["contact"].([]interface{}) {

			contact := contact.(map[string]interface{})

			if _, ok := contact["name"]; ok {

				name := contact["name"].(map[string]interface{})
				if _, ok := name["period"]; ok {

					period := name["period"].(map[string]interface{})
					newPeriod := s.periodMapper(period)

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
						newPeriod := s.periodMapper(period)

						telecom["period"] = newPeriod
					}

					newTelecoms = append(newTelecoms, telecom)
				}

				contact["telecom"] = newTelecoms
			}

			if _, ok := contact["period"]; ok {

				period := contact["period"].(map[string]interface{})
				newPeriod := s.periodMapper(period)

				contact["period"] = newPeriod
			}

			newContacts = append(newContacts, contact)
		}
	}

	resourceCopy["contact"] = newContacts

	return resourceCopy
}

// CreateUpdatePatientExtraInformation updates a patient's extra info
func (s Service) CreateUpdatePatientExtraInformation(
	ctx context.Context, input PatientExtraInformationInput) (bool, error) {
	s.checkPreconditions()

	patientPayload, err := s.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		return false, fmt.Errorf("unable to get patient with ID %s: %w", input.PatientID, err)
	}
	patient := patientPayload.PatientRecord

	patches := []map[string]interface{}{}
	op := "add" // the content will be appended to the element identified in the path

	if input.MaritalStatus != nil {
		inputMaritalStatus := MaritalStatusEnumToCodeableConcept(*input.MaritalStatus)
		if patient.MaritalStatus != inputMaritalStatus {
			patch := make(map[string]interface{})
			patch["op"] = op
			patch["path"] = "/maritalStatus"
			patch["value"] = inputMaritalStatus
			patches = append(patches, patch)
		}
	}

	if input.Languages != nil {
		langs := []base.Language{}
		for _, l := range input.Languages {
			langs = append(langs, *l)
		}
		communicationInput := LanguagesToCommunicationInputs(langs)
		if len(input.Languages) > 0 {
			patch := make(map[string]interface{})
			patch["op"] = op
			patch["path"] = "/communication"
			patch["value"] = communicationInput
			patches = append(patches, patch)
		}
	}

	if len(input.Emails) > 0 {
		emailInput, err := ContactsToContactPoint(nil, input.Emails, s.firestoreClient)
		if err != nil {
			return false, fmt.Errorf("unable to process email addresses")
		}
		telecom := patient.Telecom
		if telecom == nil {
			telecom = []*FHIRContactPoint{}
		}
		telecom = append(telecom, emailInput...)

		patch := make(map[string]interface{})
		patch["op"] = op
		patch["path"] = "/telecom"
		patch["value"] = telecom
		patches = append(patches, patch)
	}

	cloudhealthService := cloudhealth.NewService()
	_, err = cloudhealthService.PatchFHIRResource("Patient", input.PatientID, patches)
	if err != nil {
		return false, fmt.Errorf("UpdatePatient: %v", err)
	}
	return true, nil
}
