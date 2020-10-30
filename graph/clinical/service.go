// Package clinical implements a simplified GraphQL interface to a FHIR server
// that acts as a clinical data repository.
package clinical

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/cloudhealth/cloudhealth"
	"gitlab.slade360emr.com/go/sms/graph/sms"
)

const (
	// LimitedProfileEncounterCount is the number of encounters to show when a
	// patient has approved limited access to their health record
	LimitedProfileEncounterCount = 5

	// MaxClinicalRecordPageSize is the maximum number of encounters we can show on a timeline
	MaxClinicalRecordPageSize = 50
)

// NewService initializes a new clinical service
func NewService() *Service {
	clinicalRepository := cloudhealth.NewService()
	smsRepository := sms.NewService()
	return &Service{clinicalRepository: clinicalRepository, smsRepository: smsRepository}
}

// Service is a clinical service
type Service struct {
	clinicalRepository *cloudhealth.Service
	smsRepository      *sms.Service
}

func (s Service) checkPreconditions() {
	if s.clinicalRepository == nil {
		log.Panicf("clinical.Service *cloudhealth.Service is nil")
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
		logrus.Errorf("unable to search: %v", err)
		return nil, fmt.Errorf("unable to search: %v", err)
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		logrus.Errorf("%s could not be found with search params %v: %s", resourceName, params, err)
		return nil, fmt.Errorf(
			"%s could not be found with search params %v: %s", resourceName, params, err)
	}

	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
	for _, k := range mandatoryKeys {
		_, found := respMap[k]
		if !found {
			logrus.Errorf("Search: mandatory search result key %s not found", k)
			return nil, fmt.Errorf("server error: mandatory search result key %s not found", k)
		}
	}
	resourceType, ok := respMap["resourceType"].(string)
	if !ok {
		logrus.Errorf("Search: the resourceType is not a string")
		return nil, fmt.Errorf("server error: the resourceType is not a string")
	}
	if resourceType != "Bundle" {
		logrus.Errorf("Search: the resourceType value is not 'Bundle' as expected")
		return nil, fmt.Errorf(
			"server error: the resourceType value is not 'Bundle' as expected")
	}

	resultType, ok := respMap["type"].(string)
	if !ok {
		logrus.Errorf("Search: the search result type value is not a string")
		return nil, fmt.Errorf("server error: the search result type value is not a string")
	}
	if resultType != "searchset" {
		logrus.Errorf("Search: the type value is not 'searchset' as expected")
		return nil, fmt.Errorf("server error: the type value is not 'searchset' as expected")
	}

	respEntries := respMap["entry"]
	if respEntries == nil {
		return []map[string]interface{}{}, nil
	}
	entries, ok := respEntries.([]interface{})
	if !ok {
		logrus.Errorf("Search: entries is not a list of maps, it is: %T", respEntries)
		return nil, fmt.Errorf(
			"server error: entries is not a list of maps, it is: %T", respEntries)
	}

	results := []map[string]interface{}{}
	for _, en := range entries {
		entry, ok := en.(map[string]interface{})
		if !ok {
			logrus.Errorf("Search: expected each entry to be map, they are %T instead", en)
			return nil, fmt.Errorf(
				"server error: expected each entry to be map, they are %T instead", en)
		}
		expectedKeys := []string{"fullUrl", "resource", "search"}
		for _, k := range expectedKeys {
			_, found := entry[k]
			if !found {
				logrus.Errorf("Search: FHIR search entry does not have key '%s'", k)
				return nil, fmt.Errorf("server error: FHIR search entry does not have key '%s'", k)
			}
		}

		resource, ok := entry["resource"].(map[string]interface{})
		if !ok {
			{
				logrus.Errorf("Search: result entry %#v is not a map", entry["resource"])
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

// AllergySummary returns a short list of the patient's active and confirmed
// allergies (by name)
func (s Service) AllergySummary(
	ctx context.Context, patientID string) ([]string, error) {
	s.checkPreconditions()

	params := map[string]interface{}{
		"clinical-status":     "active",
		"verification-status": "confirmed",
		"type":                "allergy",
		"criticality":         "high",
		"patient":             fmt.Sprintf("Patient/%s", patientID),
	}
	results, err := s.SearchFHIRAllergyIntolerance(ctx, params)
	if err != nil {
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

// PatientTimeline returns the patient's visit note timeline (a list of
// narratives that are sorted with the most recent one first), while
// respecting the approval level
func (s Service) PatientTimeline(
	ctx context.Context, episodeID string) ([]map[string]interface{}, error) {
	episode, accessLevel, err := s.getTimelineEpisode(ctx, episodeID)
	if err != nil {
		return nil, err
	}
	encounterSearchParams := map[string]interface{}{
		"patient": *episode.Patient.Reference,
		"sort":    "-date", // reverse chronological order
	}
	count := MaxClinicalRecordPageSize
	if accessLevel == "PROFILE_AND_RECENT_VISITS_ACCESS" {
		count = LimitedProfileEncounterCount
	}
	encounterSearchParams["_count"] = strconv.Itoa(count)
	return s.getTimelineVisitSummaries(ctx, encounterSearchParams, count)
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
