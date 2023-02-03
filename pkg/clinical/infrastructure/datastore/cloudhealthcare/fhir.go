package fhir

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset"
	"github.com/savannahghi/clinical/pkg/clinical/repository"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// constants and defaults
const (
	internalError            = "an error occurred on our end. Please try again later"
	timeFormatStr            = "2006-01-02T15:04:05+03:00"
	notFoundWithSearchParams = "could not find a patient with the provided parameters"
)

var (
	patientResourceType       = "Patient"
	episodeOfCareResourceType = "EpisodeOfCare"
)

// StoreImpl represents the FHIR infrastructure implementation
type StoreImpl struct {
	Dataset dataset.FHIRRepository
}

// NewFHIRStoreImpl initializes the new FHIR implementation
func NewFHIRStoreImpl(
	dataset dataset.FHIRRepository,
) repository.FHIR {
	return &StoreImpl{
		Dataset: dataset,
	}
}

// Encounters returns encounters that belong to the indicated patient.
//
// The patientReference should be a [string] in the format "Patient/<patient resource ID>".
func (fh StoreImpl) Encounters(
	ctx context.Context,
	patientReference string,
	status *domain.EncounterStatusEnum,
) ([]*domain.FHIREncounter, error) {
	searchParams := url.Values{}
	if status != nil {
		searchParams.Add("status:exact", status.String())
	}
	searchParams.Add("patient", patientReference)

	bs, err := fh.Dataset.POSTRequest("Encounter", "_search", searchParams, nil)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to search for encounter: %v", err)
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
			utils.ReportErrorToSentry(err)
			return nil, fmt.Errorf("unable to marshal resource to JSON: %v", err)
		}
		err = json.Unmarshal(resourceBs, &encounter)
		if err != nil {
			utils.ReportErrorToSentry(err)
			return nil, fmt.Errorf("unable to unmarshal resource: %v", err)
		}
		output = append(output, &encounter)
	}
	return output, nil
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (fh StoreImpl) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "EpisodeOfCare"
	path := "_search"
	output := domain.FHIREpisodeOfCareRelayConnection{}

	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}
	for _, result := range resources {
		var resource domain.FHIREpisodeOfCare

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
func (fh StoreImpl) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
	payload, err := converterandformatter.StructToMap(episode)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
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
	data, err := fh.Dataset.CreateFHIRResource("EpisodeOfCare", payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to create episode of care resource: %v", err)
	}
	fhirEpisode := &domain.FHIREpisodeOfCare{}
	err = json.Unmarshal(data, fhirEpisode)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to unmarshal episode of care response JSON: data: %v\n, error: %v",
			string(data), err)
	}
	encounters, err := fh.Encounters(ctx, *episode.Patient.Reference, nil)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh StoreImpl) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	resourceType := "Condition"
	resource := domain.FHIRCondition{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh StoreImpl) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	resourceType := "Organization"
	resource := domain.FHIROrganization{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %v", resourceType, err)
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

// SearchFHIROrganization provides a search API for FHIROrganization
func (fh StoreImpl) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "Organization"
	path := "_search"
	output := domain.FHIROrganizationRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIROrganization

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
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

// FindOrganizationByID finds and retrieves organization details using the specified organization ID
func (fh StoreImpl) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("organization ID is required")
	}

	data, err := fh.Dataset.GetFHIRResource("Organization", organizationID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to retrieve organization: %w", err)
	}
	var organization domain.FHIROrganization
	err = json.Unmarshal(data, &organization)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to unmarshal organization data from JSON to the target struct: %w", err)
	}
	return &domain.FHIROrganizationRelayPayload{
		Resource: &organization,
	}, nil

}

// SearchEpisodesByParam search episodes by params
func (fh StoreImpl) SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
	bs, err := fh.Dataset.POSTRequest(
		"EpisodeOfCare", "_search", searchParams, nil)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to search for episode of care: %v", err)
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to map decode resource: %v", err)
			return nil, fmt.Errorf(internalError)
		}

		output = append(output, &episode)
	}
	return output, nil
}

// OpenEpisodes returns the IDs of a patient's open episodes
func (fh StoreImpl) OpenEpisodes(
	ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	searchParams := url.Values{}
	searchParams.Add("status:exact", domain.EpisodeOfCareStatusEnumActive.String())
	searchParams.Add("patient", patientReference)
	return fh.SearchEpisodesByParam(ctx, searchParams)
}

// HasOpenEpisode determines if a patient has an open episode
func (fh StoreImpl) HasOpenEpisode(
	ctx context.Context,
	patient domain.FHIRPatient,
) (bool, error) {
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	episodes, err := fh.OpenEpisodes(ctx, patientReference)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, err
	}
	return len(episodes) > 0, nil
}

// CreateFHIREncounter creates a FHIREncounter instance
func (fh StoreImpl) CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
	resourceType := "Encounter"
	resource := domain.FHIREncounter{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}
	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh StoreImpl) GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
	resourceType := "EpisodeOfCare"
	var resource domain.FHIREpisodeOfCare

	data, err := fh.Dataset.GetFHIRResource(resourceType, id)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}
	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &domain.FHIREpisodeOfCareRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// StartEncounter starts an encounter within an episode of care
func (fh *StoreImpl) StartEncounter(
	ctx context.Context, episodeID string) (string, error) {
	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
		return "", fmt.Errorf("unable to start encounter: %w", err)
	}
	return *encPl.Resource.ID, nil
}

// SearchEpisodeEncounter returns all encounters in a visit
func (fh *StoreImpl) SearchEpisodeEncounter(
	ctx context.Context,
	episodeReference string,
) (*domain.FHIREncounterRelayConnection, error) {
	episodeRef := fmt.Sprintf("Episode/%s", episodeReference)
	encounterFilterParams := map[string]interface{}{
		"episodeOfCare": episodeRef,
		"status":        "in_progress",
	}
	encounterConn, err := fh.SearchFHIREncounter(ctx, encounterFilterParams)

	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to search encounter: %w", err)
	}
	return encounterConn, nil
}

// EndEncounter ends an encounter
func (fh *StoreImpl) EndEncounter(
	ctx context.Context, encounterID string) (bool, error) {
	resourceType := "Encounter"
	encounterPayload, err := fh.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to turn the updated episode of care into a map: %v", err)
	}

	_, err = fh.Dataset.UpdateFHIRResource(resourceType, encounterID, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to create/update %s resource: %w", resourceType, err)
	}
	return true, nil
}

// EndEpisode ends an episode of care by patching it's status to "finished"
func (fh *StoreImpl) EndEpisode(
	ctx context.Context, episodeID string) (bool, error) {
	resourceType := "EpisodeOfCare"
	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}
	startTime := scalarutils.DateTime(time.Now().Format(timeFormatStr))
	if episodePayload.Resource.Period != nil {
		startTime = episodePayload.Resource.Period.Start
	}

	// Close all encounters in this visit
	encounterConn, err := fh.SearchEpisodeEncounter(ctx, episodeID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to search episode encounter %w", err)
	}
	for _, edge := range encounterConn.Edges {
		_, err = fh.EndEncounter(ctx, *edge.Node.ID)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to turn the updated episode of care into a map: %v", err)
	}

	_, err = fh.Dataset.UpdateFHIRResource(resourceType, episodeID, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to create/update %s resource: %w", resourceType, err)
	}
	return true, nil
}

// GetActiveEpisode returns any ACTIVE episode that has to the indicated ID
func (fh *StoreImpl) GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {

	searchParams := url.Values{}
	searchParams.Add("status:exact", domain.EpisodeOfCareStatusEnumActive.String())
	searchParams.Add("_id", episodeID) // logical ID of the resource

	bs, err := fh.Dataset.POSTRequest(
		"EpisodeOfCare", "_search", searchParams, nil)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to search for episode of care: %v", err)
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to marshal resource to JSON: %v", err)
	}
	err = json.Unmarshal(resourceBs, &episode)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to unmarshal resource: %v", err)
	}
	return &episode, nil
}

// SearchFHIRServiceRequest provides a search API for FHIRServiceRequest
func (fh *StoreImpl) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "ServiceRequest"
	path := "_search"
	output := domain.FHIRServiceRequestRelayConnection{}

	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRServiceRequest

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	resourceType := "ServiceRequest"
	resource := domain.FHIRServiceRequest{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "AllergyIntolerance"
	path := "_search"
	output := domain.FHIRAllergyIntoleranceRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRAllergyIntolerance

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	resourceType := "AllergyIntolerance"
	resource := domain.FHIRAllergyIntolerance{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	resourceType := "AllergyIntolerance"
	resource := domain.FHIRAllergyIntolerance{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "Composition"
	path := "_search"
	output := domain.FHIRCompositionRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRComposition

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	resourceType := "Composition"
	resource := domain.FHIRComposition{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {

	resourceType := "Composition"
	resource := domain.FHIRComposition{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	resourceType := "Composition"
	resp, err := fh.Dataset.DeleteFHIRResource(resourceType, id)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// SearchFHIRCondition provides a search API for FHIRCondition
func (fh *StoreImpl) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "Condition"
	path := "_search"
	output := domain.FHIRConditionRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRCondition

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {

	resourceType := "Condition"
	resource := domain.FHIRCondition{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {

	resourceType := "Encounter"
	var resource domain.FHIREncounter

	data, err := fh.Dataset.GetFHIRResource(resourceType, id)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &domain.FHIREncounterRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// SearchFHIREncounter provides a search API for FHIREncounter
func (fh *StoreImpl) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "Encounter"
	path := "_search"
	output := domain.FHIREncounterRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIREncounter

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "MedicationRequest"
	path := "_search"
	output := domain.FHIRMedicationRequestRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRMedicationRequest

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {

	resourceType := "MedicationRequest"
	resource := domain.FHIRMedicationRequest{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {

	resourceType := "MedicationRequest"
	resource := domain.FHIRMedicationRequest{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	resourceType := "MedicationRequest"
	resp, err := fh.Dataset.DeleteFHIRResource(resourceType, id)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// SearchFHIRObservation provides a search API for FHIRObservation
func (fh *StoreImpl) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "Observation"
	path := "_search"
	output := domain.FHIRObservationRelayConnection{}
	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRObservation

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	resourceType := "Observation"
	resource := domain.FHIRObservation{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
func (fh *StoreImpl) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {

	resourceType := "Observation"
	resp, err := fh.Dataset.DeleteFHIRResource(resourceType, id)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// GetFHIRPatient retrieves instances of FHIRPatient by ID
func (fh *StoreImpl) GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
	var resource domain.FHIRPatient

	data, err := fh.Dataset.GetFHIRResource(patientResourceType, id)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", patientResourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", patientResourceType, err)
	}

	hasOpenEpisodes, err := fh.HasOpenEpisode(ctx, resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get open episodes for patient %#v: %w", resource, err)
	}
	payload := &domain.FHIRPatientRelayPayload{
		Resource:        &resource,
		HasOpenEpisodes: hasOpenEpisodes,
	}
	return payload, nil
}

// DeleteFHIRPatient deletes the FHIRPatient identified by the supplied ID
func (fh *StoreImpl) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	patientEverythingBs, err := fh.Dataset.GetFHIRPatientEverything(id)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to get patient's compartment: %v", err)
	}

	var patientEverything map[string]interface{}
	err = json.Unmarshal(patientEverythingBs, &patientEverything)
	if err != nil {
		utils.ReportErrorToSentry(err)
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

	// This list stores all the observations resource types
	observations := []map[string]string{}

	medicationRequests := []map[string]string{}

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

		case "Observation":
			observations = append(
				observations,
				resourceTypeIDMap,
			)
			continue

		case "MedicationRequest":
			medicationRequests = append(
				medicationRequests,
				resourceTypeIDMap,
			)
			continue
		}

		assortedResourceTypes = append(assortedResourceTypes, resourceTypeIDMap)
	}

	// Special case, a medication request causes the failure for deleting a FHIR Condition
	if err = fh.DeleteFHIRResourceType(medicationRequests); err != nil {
		return false, err
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

	if err = fh.DeleteFHIRResourceType(observations); err != nil {
		return false, err
	}

	// Finally delete the patient ResourceType
	if err = fh.DeleteFHIRResourceType(patient); err != nil {
		return false, err
	}
	return true, nil
}

// DeleteFHIRResourceType takes a ResourceType and ID and deletes them from FHIR
func (fh *StoreImpl) DeleteFHIRResourceType(results []map[string]string) error {
	for _, result := range results {
		resourceType := result["resourceType"]
		resourceID := result["resourceID"]

		resp, err := fh.Dataset.DeleteFHIRResource(
			resourceType,
			resourceID,
		)
		if err != nil {
			utils.ReportErrorToSentry(err)
			return fmt.Errorf(
				"unable to delete %s:%s, response %s, error: %v",
				resourceType, resourceID, string(resp), err,
			)
		}
	}
	return nil
}

// DeleteFHIRServiceRequest deletes the FHIRServiceRequest identified by the supplied ID
func (fh *StoreImpl) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	resourceType := "ServiceRequest"
	resp, err := fh.Dataset.DeleteFHIRResource(resourceType, id)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// CreateFHIRMedicationStatement creates a new FHIR Medication statement instance
func (fh *StoreImpl) CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
	resourceType := "MedicationStatement"

	resource := domain.FHIRMedicationStatement{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRMedicationStatementRelayPayload{
		Resource: &resource,
	}

	return output, nil
}

// CreateFHIRMedication creates a new FHIR Medication instance
func (fh *StoreImpl) CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
	resourceType := "Medication"

	resource := domain.FHIRMedication{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(resourceType, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRMedicationRelayPayload{
		Resource: &resource,
	}

	return output, nil
}

// SearchFHIRMedicationStatement used to search for a fhir medication statement
func (fh *StoreImpl) SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}

	urlParams, err := validateSearchParams(params)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	resourceName := "MedicationStatement"
	path := "_search"
	output := domain.FHIRMedicationStatementRelayConnection{}

	resources, err := fh.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRMedicationStatement

		resourceBs, err := json.Marshal(result)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRMedicationStatementRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRPatient creates a patient on FHIR
func (fh *StoreImpl) CreateFHIRPatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	resource := domain.FHIRPatient{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", patientResourceType, err)
	}

	data, err := fh.Dataset.CreateFHIRResource(patientResourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", patientResourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			patientResourceType, string(data), err)
	}

	output := &domain.PatientPayload{
		PatientRecord: &resource,
		// The patient is newly created so we can assume they have no open episodes
		HasOpenEpisodes: false,
		OpenEpisodes:    []*domain.FHIREpisodeOfCare{},
	}
	return output, nil
}

// PatchFHIRPatient is used to patch a patient resource
func (fh *StoreImpl) PatchFHIRPatient(ctx context.Context, id string, params []map[string]interface{}) (*domain.FHIRPatient, error) {
	resource := domain.FHIRPatient{}

	data, err := fh.Dataset.PatchFHIRResource(patientResourceType, id, params)
	if err != nil {
		return nil, fmt.Errorf("unable to patch %s resource: %v", patientResourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			patientResourceType, string(data), err)
	}
	return &resource, nil
}

// UpdateFHIREpisodeOfCare updates a fhir episode of care
func (fh *StoreImpl) UpdateFHIREpisodeOfCare(ctx context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error) {
	resource := domain.FHIREpisodeOfCare{}

	if fhirResourceID == "" {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	data, err := fh.Dataset.UpdateFHIRResource(episodeOfCareResourceType, fhirResourceID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to update %s resource: %w", episodeOfCareResourceType, err)
	}
	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			episodeOfCareResourceType, string(data), err)
	}
	return &resource, nil
}

// CreateFHIRResource creates a FHIR resource
func (fh *StoreImpl) CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error) {
	return fh.Dataset.CreateFHIRResource(resourceType, payload)
}

// PatchFHIRResource patches a FHIR resource
func (fh *StoreImpl) PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error) {
	return fh.Dataset.PatchFHIRResource(resourceType, fhirResourceID, payload)
}

// UpdateFHIRResource updates a FHIR resource
func (fh *StoreImpl) UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
	return fh.Dataset.UpdateFHIRResource(resourceType, fhirResourceID, payload)
}

// SearchFHIRPatient searches for a FHIR patient
func (fh *StoreImpl) SearchFHIRPatient(ctx context.Context, searchParams string) (*domain.PatientConnection, error) {
	params := url.Values{}
	params.Add("_content", searchParams)

	bs, err := fh.Dataset.POSTRequest(patientResourceType, "_search", params, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to search for patient: %w", err)
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

		resource = birthdateMapper(resource)
		resource = identifierMapper(resource)
		resource = nameMapper(resource)
		resource = telecomMapper(resource)
		resource = addressMapper(resource)
		resource = photoMapper(resource)
		resource = contactMapper(resource)

		var patient domain.FHIRPatient

		err := mapstructure.Decode(resource, &patient)
		if err != nil {
			utils.ReportErrorToSentry(err)
			log.Errorf("unable to map decode resource: %v", err)
			return nil, fmt.Errorf(internalError)
		}

		hasOpenEpisodes, err := fh.HasOpenEpisode(ctx, patient)
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

// POSTRequest is used to manually compose POST requests to the FHIR service
func (fh *StoreImpl) POSTRequest(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
	return fh.Dataset.POSTRequest(resourceName, path, params, body)
}

// GetFHIRResource gets a FHIR resource.
func (fh *StoreImpl) GetFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return fh.Dataset.GetFHIRResource(resourceType, fhirResourceID)
}

// FHIRHeaders composes suitable FHIR headers, with authentication and content
// type already set
func (fh *StoreImpl) FHIRHeaders() (http.Header, error) {
	return fh.Dataset.FHIRHeaders()
}
