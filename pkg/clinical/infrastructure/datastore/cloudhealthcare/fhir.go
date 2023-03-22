package fhir

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/scalarutils"
)

// constants and defaults
const (
	internalError = "an error occurred on our end. Please try again later"
	timeFormatStr = "2006-01-02T15:04:05+03:00"
)

// resource types
const (
	organizationResource            = "Organization"
	patientResourceType             = "Patient"
	episodeOfCareResourceType       = "EpisodeOfCare"
	observationResourceType         = "Observation"
	allergyIntoleranceResourceType  = "AllergyIntolerance"
	serviceRequestResourceType      = "ServiceRequest"
	medicationRequestResourceType   = "MedicationRequest"
	conditionResourceType           = "Condition"
	encounterResourceType           = "Encounter"
	compositionResourceType         = "Composition"
	medicationStatementResourceType = "MedicationStatement"
	medicationResourceType          = "Medication"
)

// Dataset ...
type Dataset interface {
	GetFHIRResource(resourceType, fhirResourceID string, resource interface{}) error
	CreateFHIRResource(resourceType string, payload map[string]interface{}, resource interface{}) error
	DeleteFHIRResource(resourceType, fhirResourceID string) error
	PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}, resource interface{}) error
	UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error
	SearchFHIRResource(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers) ([]map[string]interface{}, error)

	GetFHIRPatientAllData(fhirResourceID string) ([]byte, error)
}

// StoreImpl represents the FHIR infrastructure implementation
type StoreImpl struct {
	Dataset Dataset
}

// NewFHIRStoreImpl initializes the new FHIR implementation
func NewFHIRStoreImpl(
	dataset Dataset,
) *StoreImpl {
	return &StoreImpl{
		Dataset: dataset,
	}
}

// Encounters returns encounters that belong to the indicated patient.
//
// The patientReference should be a [string] in the format "Patient/<patient resource ID>".
func (fh StoreImpl) SearchPatientEncounters(
	_ context.Context,
	patientReference string,
	status *domain.EncounterStatusEnum,
	tenant dto.TenantIdentifiers,
) ([]*domain.FHIREncounter, error) {
	params := map[string]interface{}{
		"patient": patientReference,
	}
	if status != nil {
		params["status:exact"] = status.String()
	}

	resources, err := fh.Dataset.SearchFHIRResource(encounterResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	output := []*domain.FHIREncounter{}

	for _, resource := range resources {
		var encounter domain.FHIREncounter

		resourceBs, err := json.Marshal(resource)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &encounter)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		output = append(output, &encounter)
	}

	return output, nil
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (fh StoreImpl) SearchFHIREpisodeOfCare(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	output := domain.FHIREpisodeOfCareRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(episodeOfCareResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIREpisodeOfCare

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", episodeOfCareResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIREpisodeOfCareRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// CreateEpisodeOfCare is the final common pathway for creation of episodes of
// care.
func (fh StoreImpl) CreateEpisodeOfCare(_ context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error) {
	payload, err := converterandformatter.StructToMap(episode)
	if err != nil {
		return nil, fmt.Errorf("unable to turn episode of care input into a map: %w", err)
	}

	fhirEpisode := &domain.FHIREpisodeOfCare{}
	// create a new episode if none has been found

	err = fh.Dataset.CreateFHIRResource(episodeOfCareResourceType, payload, fhirEpisode)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to create episode of care resource: %w", err)
	}

	output := &domain.EpisodeOfCarePayload{
		EpisodeOfCare: fhirEpisode,
	}

	return output, nil
}

// CreateFHIRCondition creates a FHIRCondition instance
func (fh StoreImpl) CreateFHIRCondition(_ context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", conditionResourceType, err)
	}

	resource := &domain.FHIRCondition{}

	err = fh.Dataset.CreateFHIRResource(conditionResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", conditionResourceType, err)
	}

	output := &domain.FHIRConditionRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// CreateFHIROrganization creates a FHIROrganization instance
func (fh StoreImpl) CreateFHIROrganization(_ context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", organizationResource, err)
	}

	resource := &domain.FHIROrganization{}

	err = fh.Dataset.CreateFHIRResource(organizationResource, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", organizationResource, err)
	}

	output := &domain.FHIROrganizationRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// SearchFHIROrganization provides a search API for FHIROrganization
func (fh StoreImpl) SearchFHIROrganization(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIROrganizationRelayConnection, error) {
	output := domain.FHIROrganizationRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(organizationResource, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIROrganization

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", organizationResource, err)
		}

		output.Edges = append(output.Edges, &domain.FHIROrganizationRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// GetFHIROrganization finds and retrieves organization details using the specified organization ID
func (fh StoreImpl) GetFHIROrganization(_ context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("organization ID is required")
	}

	organization := &domain.FHIROrganization{}

	err := fh.Dataset.GetFHIRResource(organizationResource, organizationID, organization)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve organization: %w", err)
	}

	return &domain.FHIROrganizationRelayPayload{
		Resource: organization,
	}, nil
}

// SearchEpisodesByParam search episodes by params
func (fh StoreImpl) SearchEpisodesByParam(_ context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers) ([]*domain.FHIREpisodeOfCare, error) {
	resources, err := fh.Dataset.SearchFHIRResource(episodeOfCareResourceType, searchParams, tenant)
	if err != nil {
		return nil, err
	}

	output := []*domain.FHIREpisodeOfCare{}

	for _, resource := range resources {
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

		// update the original period resource
		resource["period"] = period

		var episode domain.FHIREpisodeOfCare

		err := mapstructure.Decode(resource, &episode)
		if err != nil {
			return nil, fmt.Errorf(internalError)
		}

		output = append(output, &episode)
	}

	return output, nil
}

// OpenEpisodes returns the IDs of a patient's open episodes
func (fh StoreImpl) OpenEpisodes(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers) ([]*domain.FHIREpisodeOfCare, error) {
	params := map[string]interface{}{
		"status:exact": domain.EpisodeOfCareStatusEnumActive.String(),
		"patient":      patientReference,
	}

	return fh.SearchEpisodesByParam(ctx, params, tenant)
}

// HasOpenEpisode determines if a patient has an open episode
func (fh StoreImpl) HasOpenEpisode(
	ctx context.Context,
	patient domain.FHIRPatient,
	tenant dto.TenantIdentifiers,
) (bool, error) {
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)

	episodes, err := fh.OpenEpisodes(ctx, patientReference, tenant)
	if err != nil {
		return false, err
	}

	return len(episodes) > 0, nil
}

// CreateFHIREncounter creates a FHIREncounter instance
func (fh StoreImpl) CreateFHIREncounter(_ context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", encounterResourceType, err)
	}

	resource := &domain.FHIREncounter{}

	err = fh.Dataset.CreateFHIRResource(encounterResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", encounterResourceType, err)
	}

	output := &domain.FHIREncounterRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// GetFHIREpisodeOfCare retrieves instances of FHIREpisodeOfCare by ID
func (fh StoreImpl) GetFHIREpisodeOfCare(_ context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
	resource := &domain.FHIREpisodeOfCare{}

	err := fh.Dataset.GetFHIRResource(episodeOfCareResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", episodeOfCareResourceType, id, err)
	}

	payload := &domain.FHIREpisodeOfCareRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// StartEncounter starts an encounter within an episode of care
func (fh *StoreImpl) StartEncounter(
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

// SearchEpisodeEncounter returns all encounters in a visit
func (fh *StoreImpl) SearchEpisodeEncounter(
	ctx context.Context,
	episodeReference string,
	tenant dto.TenantIdentifiers,
) (*domain.FHIREncounterRelayConnection, error) {
	episodeRef := fmt.Sprintf("Episode/%s", episodeReference)
	encounterFilterParams := map[string]interface{}{
		"episodeOfCare": episodeRef,
		"status":        "in_progress",
	}
	encounterConn, err := fh.SearchFHIREncounter(ctx, encounterFilterParams, tenant)

	if err != nil {
		return nil, fmt.Errorf("unable to search encounter: %w", err)
	}

	return encounterConn, nil
}

// EndEncounter ends an encounter
func (fh *StoreImpl) EndEncounter(
	ctx context.Context, encounterID string) (bool, error) {
	encounterPayload, err := fh.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return false, err
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
		return false, fmt.Errorf("unable to turn the updated episode of care into a map: %w", err)
	}

	encounter := &domain.FHIREncounter{}

	err = fh.Dataset.UpdateFHIRResource(encounterResourceType, encounterID, payload, encounter)
	if err != nil {
		return false, fmt.Errorf("unable to create/update %s resource: %w", encounterResourceType, err)
	}

	return true, nil
}

// EndEpisode ends an episode of care by patching it's status to "finished"
func (fh *StoreImpl) EndEpisode(
	ctx context.Context, episodeID string) (bool, error) {
	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return false, fmt.Errorf("unable to get episode with ID %s: %w", episodeID, err)
	}

	startTime := scalarutils.DateTime(time.Now().Format(timeFormatStr))
	if episodePayload.Resource.Period != nil {
		startTime = episodePayload.Resource.Period.Start
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
		return false, fmt.Errorf("unable to turn the updated episode of care into a map: %w", err)
	}

	episode := &domain.FHIREpisodeOfCare{}

	err = fh.Dataset.UpdateFHIRResource(episodeOfCareResourceType, episodeID, payload, episode)
	if err != nil {
		return false, fmt.Errorf("unable to create/update %s resource: %w", episodeOfCareResourceType, err)
	}

	return true, nil
}

// GetActiveEpisode returns any ACTIVE episode that has to the indicated ID
func (fh *StoreImpl) GetActiveEpisode(_ context.Context, episodeID string, tenant dto.TenantIdentifiers) (*domain.FHIREpisodeOfCare, error) {
	params := map[string]interface{}{
		"status:exact": domain.EpisodeOfCareStatusEnumActive.String(),
		"_id":          episodeID,
	}

	resources, err := fh.Dataset.SearchFHIRResource(episodeOfCareResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	if len(resources) != 1 {
		return nil, fmt.Errorf(
			"expected exactly one ACTIVE episode for episode ID %s, got %d", episodeID, len(resources))
	}

	var episode domain.FHIREpisodeOfCare

	resourceBs, err := json.Marshal(resources[0])
	if err != nil {
		return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
	}

	err = json.Unmarshal(resourceBs, &episode)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
	}

	return &episode, nil
}

// SearchFHIRServiceRequest provides a search API for FHIRServiceRequest
func (fh *StoreImpl) SearchFHIRServiceRequest(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIRServiceRequestRelayConnection, error) {
	output := domain.FHIRServiceRequestRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(serviceRequestResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRServiceRequest

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", serviceRequestResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRServiceRequestRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// CreateFHIRServiceRequest creates a FHIRServiceRequest instance
func (fh *StoreImpl) CreateFHIRServiceRequest(_ context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", serviceRequestResourceType, err)
	}

	resource := &domain.FHIRServiceRequest{}

	err = fh.Dataset.CreateFHIRResource(serviceRequestResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", serviceRequestResourceType, err)
	}

	output := &domain.FHIRServiceRequestRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// SearchFHIRAllergyIntolerance provides a search API for FHIRAllergyIntolerance
func (fh *StoreImpl) SearchFHIRAllergyIntolerance(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	output := domain.FHIRAllergyIntoleranceRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(allergyIntoleranceResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRAllergyIntolerance

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", allergyIntoleranceResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRAllergyIntoleranceRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// CreateFHIRAllergyIntolerance creates a FHIRAllergyIntolerance instance
func (fh *StoreImpl) CreateFHIRAllergyIntolerance(_ context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", allergyIntoleranceResourceType, err)
	}

	resource := &domain.FHIRAllergyIntolerance{}

	err = fh.Dataset.CreateFHIRResource(allergyIntoleranceResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", allergyIntoleranceResourceType, err)
	}

	output := &domain.FHIRAllergyIntoleranceRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// UpdateFHIRAllergyIntolerance updates a FHIRAllergyIntolerance instance
// The resource must have it's ID set.
func (fh *StoreImpl) UpdateFHIRAllergyIntolerance(_ context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", allergyIntoleranceResourceType, err)
	}

	resource := &domain.FHIRAllergyIntolerance{}

	err = fh.Dataset.UpdateFHIRResource(allergyIntoleranceResourceType, *input.ID, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", allergyIntoleranceResourceType, err)
	}

	output := &domain.FHIRAllergyIntoleranceRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// SearchFHIRComposition provides a search API for FHIRComposition
func (fh *StoreImpl) SearchFHIRComposition(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIRCompositionRelayConnection, error) {
	output := domain.FHIRCompositionRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(compositionResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRComposition

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", compositionResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRCompositionRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// CreateFHIRComposition creates a FHIRComposition instance
func (fh *StoreImpl) CreateFHIRComposition(_ context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", compositionResourceType, err)
	}

	resource := &domain.FHIRComposition{}

	err = fh.Dataset.CreateFHIRResource(compositionResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", compositionResourceType, err)
	}

	output := &domain.FHIRCompositionRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// UpdateFHIRComposition updates a FHIRComposition instance
// The resource must have it's ID set.
func (fh *StoreImpl) UpdateFHIRComposition(_ context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", compositionResourceType, err)
	}

	resource := &domain.FHIRComposition{}

	err = fh.Dataset.UpdateFHIRResource(compositionResourceType, *input.ID, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", compositionResourceType, err)
	}

	output := &domain.FHIRCompositionRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// DeleteFHIRComposition deletes the FHIRComposition identified by the supplied ID
func (fh *StoreImpl) DeleteFHIRComposition(_ context.Context, id string) (bool, error) {
	err := fh.Dataset.DeleteFHIRResource(compositionResourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, error: %w",
			compositionResourceType, err,
		)
	}

	return true, nil
}

// SearchFHIRCondition provides a search API for FHIRCondition
func (fh *StoreImpl) SearchFHIRCondition(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIRConditionRelayConnection, error) {
	output := domain.FHIRConditionRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(conditionResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRCondition

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", conditionResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRConditionRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// UpdateFHIRCondition updates a FHIRCondition instance
// The resource must have it's ID set.
func (fh *StoreImpl) UpdateFHIRCondition(_ context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", conditionResourceType, err)
	}

	resource := &domain.FHIRCondition{}

	err = fh.Dataset.UpdateFHIRResource(conditionResourceType, *input.ID, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", conditionResourceType, err)
	}

	output := &domain.FHIRConditionRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// GetFHIREncounter retrieves instances of FHIREncounter by ID
func (fh *StoreImpl) GetFHIREncounter(_ context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	resource := &domain.FHIREncounter{}

	err := fh.Dataset.GetFHIRResource(encounterResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", encounterResourceType, id, err)
	}

	payload := &domain.FHIREncounterRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// SearchFHIREncounter provides a search API for FHIREncounter
func (fh *StoreImpl) SearchFHIREncounter(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIREncounterRelayConnection, error) {
	output := domain.FHIREncounterRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(encounterResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIREncounter

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", encounterResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIREncounterRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// SearchFHIRMedicationRequest provides a search API for FHIRMedicationRequest
func (fh *StoreImpl) SearchFHIRMedicationRequest(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIRMedicationRequestRelayConnection, error) {
	output := domain.FHIRMedicationRequestRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(medicationRequestResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRMedicationRequest

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", medicationRequestResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRMedicationRequestRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// CreateFHIRMedicationRequest creates a FHIRMedicationRequest instance
func (fh *StoreImpl) CreateFHIRMedicationRequest(_ context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationRequestResourceType, err)
	}

	resource := &domain.FHIRMedicationRequest{}

	err = fh.Dataset.CreateFHIRResource(medicationRequestResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", medicationRequestResourceType, err)
	}

	output := &domain.FHIRMedicationRequestRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// UpdateFHIRMedicationRequest updates a FHIRMedicationRequest instance
// The resource must have it's ID set.
func (fh *StoreImpl) UpdateFHIRMedicationRequest(_ context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationRequestResourceType, err)
	}

	resource := &domain.FHIRMedicationRequest{}

	err = fh.Dataset.UpdateFHIRResource(medicationRequestResourceType, *input.ID, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", medicationRequestResourceType, err)
	}

	output := &domain.FHIRMedicationRequestRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// DeleteFHIRMedicationRequest deletes the FHIRMedicationRequest identified by the supplied ID
func (fh *StoreImpl) DeleteFHIRMedicationRequest(_ context.Context, id string) (bool, error) {
	err := fh.Dataset.DeleteFHIRResource(medicationRequestResourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, error: %w",
			medicationRequestResourceType, err,
		)
	}

	return true, nil
}

// SearchFHIRObservation provides a search API for FHIRObservation
func (fh *StoreImpl) SearchFHIRObservation(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIRObservationRelayConnection, error) {
	output := domain.FHIRObservationRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(observationResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRObservation

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", observationResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRObservationRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// CreateFHIRObservation creates a FHIRObservation instance
func (fh *StoreImpl) CreateFHIRObservation(_ context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", observationResourceType, err)
	}

	resource := &domain.FHIRObservation{}

	err = fh.Dataset.CreateFHIRResource(observationResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", observationResourceType, err)
	}

	output := &domain.FHIRObservationRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// DeleteFHIRObservation deletes the FHIRObservation identified by the passed ID
func (fh *StoreImpl) DeleteFHIRObservation(_ context.Context, id string) (bool, error) {
	err := fh.Dataset.DeleteFHIRResource(observationResourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, error: %w",
			observationResourceType, err,
		)
	}

	return true, nil
}

// GetFHIRPatient retrieves instances of FHIRPatient by ID
func (fh *StoreImpl) GetFHIRPatient(_ context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
	resource := &domain.FHIRPatient{}

	err := fh.Dataset.GetFHIRResource(patientResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", patientResourceType, id, err)
	}

	payload := &domain.FHIRPatientRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// DeleteFHIRPatient deletes the FHIRPatient identified by the supplied ID
func (fh *StoreImpl) DeleteFHIRPatient(_ context.Context, id string) (bool, error) {
	patientEverythingBs, err := fh.Dataset.GetFHIRPatientAllData(id)
	if err != nil {
		return false, fmt.Errorf("unable to get patient's compartment: %w", err)
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
		case encounterResourceType:
			encounters = append(
				encounters,
				resourceTypeIDMap,
			)

			continue

		case episodeOfCareResourceType:
			episodesOfCare = append(
				episodesOfCare,
				resourceTypeIDMap,
			)

			continue

		case patientResourceType:
			patient = append(
				patient,
				resourceTypeIDMap,
			)

			continue

		case observationResourceType:
			observations = append(
				observations,
				resourceTypeIDMap,
			)

			continue

		case medicationRequestResourceType:
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

		err := fh.Dataset.DeleteFHIRResource(
			resourceType,
			resourceID,
		)
		if err != nil {
			return fmt.Errorf(
				"unable to delete %s:%s, error: %w",
				resourceType, resourceID, err,
			)
		}
	}

	return nil
}

// DeleteFHIRServiceRequest deletes the FHIRServiceRequest identified by the supplied ID
func (fh *StoreImpl) DeleteFHIRServiceRequest(_ context.Context, id string) (bool, error) {
	err := fh.Dataset.DeleteFHIRResource(serviceRequestResourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, error: %w",
			serviceRequestResourceType, err,
		)
	}

	return true, nil
}

// CreateFHIRMedicationStatement creates a new FHIR Medication statement instance
func (fh *StoreImpl) CreateFHIRMedicationStatement(_ context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationStatementResourceType, err)
	}

	resource := &domain.FHIRMedicationStatement{}

	err = fh.Dataset.CreateFHIRResource(medicationStatementResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", medicationStatementResourceType, err)
	}

	output := &domain.FHIRMedicationStatementRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// CreateFHIRMedication creates a new FHIR Medication instance
func (fh *StoreImpl) CreateFHIRMedication(_ context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", medicationResourceType, err)
	}

	resource := &domain.FHIRMedication{}

	err = fh.Dataset.CreateFHIRResource(medicationResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", medicationResourceType, err)
	}

	output := &domain.FHIRMedicationRelayPayload{
		Resource: resource,
	}

	return output, nil
}

// SearchFHIRMedicationStatement used to search for a fhir medication statement
func (fh *StoreImpl) SearchFHIRMedicationStatement(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers) (*domain.FHIRMedicationStatementRelayConnection, error) {
	output := domain.FHIRMedicationStatementRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(medicationStatementResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource domain.FHIRMedicationStatement

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", medicationStatementResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRMedicationStatementRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// CreateFHIRPatient creates a patient on FHIR
func (fh *StoreImpl) CreateFHIRPatient(_ context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", patientResourceType, err)
	}

	resource := &domain.FHIRPatient{}

	err = fh.Dataset.CreateFHIRResource(patientResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", patientResourceType, err)
	}

	output := &domain.PatientPayload{
		PatientRecord: resource,
	}

	return output, nil
}

// PatchFHIRPatient is used to patch a patient resource
func (fh *StoreImpl) PatchFHIRPatient(_ context.Context, id string, params []map[string]interface{}) (*domain.FHIRPatient, error) {
	resource := &domain.FHIRPatient{}

	err := fh.Dataset.PatchFHIRResource(patientResourceType, id, params, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to patch %s resource: %w", patientResourceType, err)
	}

	return resource, nil
}

// UpdateFHIREpisodeOfCare updates a fhir episode of care
func (fh *StoreImpl) UpdateFHIREpisodeOfCare(_ context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error) {
	if fhirResourceID == "" {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	resource := &domain.FHIREpisodeOfCare{}

	err := fh.Dataset.UpdateFHIRResource(episodeOfCareResourceType, fhirResourceID, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to update %s resource: %w", episodeOfCareResourceType, err)
	}

	return resource, nil
}

// SearchFHIRPatient searches for a FHIR patient
func (fh *StoreImpl) SearchFHIRPatient(_ context.Context, searchParams string, tenant dto.TenantIdentifiers) (*domain.PatientConnection, error) {
	params := map[string]interface{}{
		"_content": searchParams,
	}

	resources, err := fh.Dataset.SearchFHIRResource(patientResourceType, params, tenant)
	if err != nil {
		return nil, err
	}

	output := domain.PatientConnection{}

	for _, resource := range resources {
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
			return nil, fmt.Errorf("%s, error:%w", internalError, err)
		}

		output.Edges = append(output.Edges, &domain.PatientEdge{
			Node: &patient,
		})
	}

	return &output, nil
}
