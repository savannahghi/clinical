package fhir

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
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
	organizationResource              = "Organization"
	patientResourceType               = "Patient"
	episodeOfCareResourceType         = "EpisodeOfCare"
	observationResourceType           = "Observation"
	allergyIntoleranceResourceType    = "AllergyIntolerance"
	serviceRequestResourceType        = "ServiceRequest"
	medicationRequestResourceType     = "MedicationRequest"
	conditionResourceType             = "Condition"
	encounterResourceType             = "Encounter"
	compositionResourceType           = "Composition"
	medicationStatementResourceType   = "MedicationStatement"
	medicationResourceType            = "Medication"
	mediaResourceType                 = "Media"
	questionnaireResourceType         = "Questionnaire"
	consentResourceType               = "Consent"
	questionnaireResponseResourceType = "QuestionnaireResponse"
	riskAssessmentResourceType        = "RiskAssessment"
	diagnosticReportResourceType      = "DiagnosticReport"
	subscriptionResourceType          = "Subscription"
)

// Dataset ...
type Dataset interface {
	GetFHIRResource(resourceType, fhirResourceID string, resource interface{}) error
	CreateFHIRResource(resourceType string, payload map[string]interface{}, resource interface{}) error
	DeleteFHIRResource(resourceType, fhirResourceID string) error
	PatchFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error
	UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error
	SearchFHIRResource(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error)

	GetFHIRPatientAllData(fhirResourceID string, params map[string]interface{}) ([]byte, error)
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

// SearchPatientObservations fetches all observations that belong to a specific patient
func (fh StoreImpl) SearchPatientObservations(
	_ context.Context,
	searchParameters map[string]interface{},
	tenant dto.TenantIdentifiers,
	pagination dto.Pagination,
) (*domain.PagedFHIRObservations, error) {
	observations, err := fh.Dataset.SearchFHIRResource(observationResourceType, searchParameters, tenant, pagination)
	if err != nil {
		return nil, err
	}

	observationOutput := domain.PagedFHIRObservations{
		Observations:    []domain.FHIRObservation{},
		HasNextPage:     observations.HasNextPage,
		NextCursor:      observations.NextCursor,
		HasPreviousPage: observations.HasPreviousPage,
		PreviousCursor:  observations.PreviousCursor,
		TotalCount:      observations.TotalCount,
	}

	for _, obs := range observations.Resources {
		var observation domain.FHIRObservation

		resourceBs, err := json.Marshal(obs)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &observation)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		observationOutput.Observations = append(observationOutput.Observations, observation)
	}

	return &observationOutput, nil
}

// Encounters returns encounters that belong to the indicated patient.
//
// The patientReference should be a [string] in the format "Patient/<patient resource ID>".
func (fh StoreImpl) SearchPatientEncounters(
	_ context.Context,
	patientReference string,
	status *domain.EncounterStatusEnum,
	tenant dto.TenantIdentifiers,
	pagination dto.Pagination,
) (*domain.PagedFHIREncounter, error) {
	params := map[string]interface{}{
		"patient": patientReference,
	}
	if status != nil {
		params["status:exact"] = status.String()
	}

	resources, err := fh.Dataset.SearchFHIRResource(encounterResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	encounterOutput := domain.PagedFHIREncounter{
		Encounters:      []domain.FHIREncounter{},
		HasNextPage:     resources.HasNextPage,
		NextCursor:      resources.NextCursor,
		HasPreviousPage: resources.HasPreviousPage,
		PreviousCursor:  resources.PreviousCursor,
		TotalCount:      resources.TotalCount,
	}

	for _, resource := range resources.Resources {
		var encounter domain.FHIREncounter

		resourceBs, err := json.Marshal(resource)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &encounter)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		encounterOutput.Encounters = append(encounterOutput.Encounters, encounter)
	}

	return &encounterOutput, nil
}

// SearchPatentMedia searches all the patients media resources
func (fh StoreImpl) SearchPatientMedia(_ context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedia, error) {
	params := map[string]interface{}{
		"patient": patientReference,
	}

	resources, err := fh.Dataset.SearchFHIRResource(mediaResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	mediaOutput := domain.PagedFHIRMedia{
		Media:           []domain.FHIRMedia{},
		HasNextPage:     resources.HasNextPage,
		NextCursor:      resources.NextCursor,
		HasPreviousPage: resources.HasPreviousPage,
		PreviousCursor:  resources.PreviousCursor,
		TotalCount:      resources.TotalCount,
	}

	for _, resource := range resources.Resources {
		var media domain.FHIRMedia

		resourceBs, err := json.Marshal(resource)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &media)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		mediaOutput.Media = append(mediaOutput.Media, media)
	}

	return &mediaOutput, nil
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (fh StoreImpl) SearchFHIREpisodeOfCare(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	output := domain.FHIREpisodeOfCareRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(episodeOfCareResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	for _, result := range resources.Resources {
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
		return nil, fmt.Errorf("unable to create %s resource: %w", conditionResourceType, err)
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
func (fh StoreImpl) SearchFHIROrganization(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIROrganizationRelayConnection, error) {
	output := domain.FHIROrganizationRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(organizationResource, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	for _, result := range resources.Resources {
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

// GetFHIRAllergyIntolerance fetches the allergy from FHIR repository using its id
func (fh StoreImpl) GetFHIRAllergyIntolerance(_ context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	allergyIntoleranace := &domain.FHIRAllergyIntolerance{}

	err := fh.Dataset.GetFHIRResource(allergyIntoleranceResourceType, id, allergyIntoleranace)
	if err != nil {
		return nil, err
	}

	return &domain.FHIRAllergyIntoleranceRelayPayload{
		Resource: allergyIntoleranace,
	}, nil
}

// SearchEpisodesByParam search episodes by params
func (fh StoreImpl) SearchEpisodesByParam(_ context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error) {
	resources, err := fh.Dataset.SearchFHIRResource(episodeOfCareResourceType, searchParams, tenant, pagination)
	if err != nil {
		return nil, err
	}

	output := []*domain.FHIREpisodeOfCare{}

	for _, resource := range resources.Resources {
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
func (fh StoreImpl) OpenEpisodes(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error) {
	params := map[string]interface{}{
		"status:exact": domain.EpisodeOfCareStatusEnumActive.String(),
		"patient":      patientReference,
	}

	return fh.SearchEpisodesByParam(ctx, params, tenant, pagination)
}

// HasOpenEpisode determines if a patient has an open episode
func (fh StoreImpl) HasOpenEpisode(
	ctx context.Context,
	patient domain.FHIRPatient,
	tenant dto.TenantIdentifiers,
	pagination dto.Pagination,
) (bool, error) {
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)

	episodes, err := fh.OpenEpisodes(ctx, patientReference, tenant, pagination)
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
func (fh StoreImpl) StartEncounter(
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

// PatchFHIREncounter is used to patch an encounter resource
func (fh StoreImpl) PatchFHIREncounter(
	_ context.Context,
	encounterID string,
	input domain.FHIREncounterInput,
) (*domain.FHIREncounter, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", encounterResourceType, err)
	}

	resource := &domain.FHIREncounter{}

	err = fh.Dataset.PatchFHIRResource(encounterResourceType, encounterID, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to patch %s resource: %w", encounterResourceType, err)
	}

	return resource, nil
}

// SearchEpisodeEncounter returns all encounters in a visit
func (fh StoreImpl) SearchEpisodeEncounter(
	ctx context.Context,
	episodeReference string,
	tenant dto.TenantIdentifiers,
	pagination dto.Pagination,
) (*domain.PagedFHIREncounter, error) {
	episodeRef := fmt.Sprintf("Episode/%s", episodeReference)
	encounterFilterParams := map[string]interface{}{
		"episodeOfCare": episodeRef,
		"status":        "in_progress",
	}
	encounterConn, err := fh.SearchFHIREncounter(ctx, encounterFilterParams, tenant, pagination)

	if err != nil {
		return nil, fmt.Errorf("unable to search encounter: %w", err)
	}

	return encounterConn, nil
}

// EndEncounter ends an encounter
func (fh StoreImpl) EndEncounter(
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

// EndEpisode ends an episode of care by patching its status to "finished"
func (fh StoreImpl) EndEpisode(
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
func (fh StoreImpl) GetActiveEpisode(_ context.Context, episodeID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCare, error) {
	params := map[string]interface{}{
		"status:exact": domain.EpisodeOfCareStatusEnumActive.String(),
		"_id":          episodeID,
	}

	resources, err := fh.Dataset.SearchFHIRResource(episodeOfCareResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	if len(resources.Resources) != 1 {
		return nil, fmt.Errorf(
			"expected exactly one ACTIVE episode for episode ID %s, got %d", episodeID, len(resources.Resources))
	}

	var episode domain.FHIREpisodeOfCare

	resourceBs, err := json.Marshal(resources.Resources[0])
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
func (fh StoreImpl) SearchFHIRServiceRequest(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRServiceRequestRelayConnection, error) {
	output := domain.FHIRServiceRequestRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(serviceRequestResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	for _, result := range resources.Resources {
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
func (fh StoreImpl) CreateFHIRServiceRequest(_ context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
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
func (fh StoreImpl) SearchFHIRAllergyIntolerance(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
	resources, err := fh.Dataset.SearchFHIRResource(allergyIntoleranceResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	output := domain.PagedFHIRAllergy{
		Allergies:       []domain.FHIRAllergyIntolerance{},
		HasNextPage:     resources.HasNextPage,
		NextCursor:      resources.NextCursor,
		HasPreviousPage: resources.HasPreviousPage,
		PreviousCursor:  resources.PreviousCursor,
		TotalCount:      resources.TotalCount,
	}

	for _, result := range resources.Resources {
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

		output.Allergies = append(output.Allergies, resource)
	}

	return &output, nil
}

// CreateFHIRAllergyIntolerance creates a FHIRAllergyIntolerance instance
func (fh StoreImpl) CreateFHIRAllergyIntolerance(_ context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
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
// The resource must have its ID set.
func (fh StoreImpl) UpdateFHIRAllergyIntolerance(_ context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
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
func (fh StoreImpl) SearchFHIRComposition(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRComposition, error) {
	resources, err := fh.Dataset.SearchFHIRResource(compositionResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	output := domain.PagedFHIRComposition{
		Compositions:    []domain.FHIRComposition{},
		HasNextPage:     resources.HasNextPage,
		NextCursor:      resources.NextCursor,
		HasPreviousPage: resources.HasPreviousPage,
		PreviousCursor:  resources.PreviousCursor,
		TotalCount:      resources.TotalCount,
	}

	for _, result := range resources.Resources {
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

		output.Compositions = append(output.Compositions, resource)
	}

	return &output, nil
}

// CreateFHIRComposition creates a FHIRComposition instance
func (fh StoreImpl) CreateFHIRComposition(_ context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
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
// The resource must have its ID set.
func (fh StoreImpl) UpdateFHIRComposition(_ context.Context, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error) {
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

	return resource, nil
}

// DeleteFHIRComposition deletes the FHIRComposition identified by the supplied ID
func (fh StoreImpl) DeleteFHIRComposition(_ context.Context, id string) (bool, error) {
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
func (fh StoreImpl) SearchFHIRCondition(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
	resources, err := fh.Dataset.SearchFHIRResource(conditionResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	output := domain.PagedFHIRCondition{
		Conditions:      []domain.FHIRCondition{},
		HasNextPage:     resources.HasNextPage,
		NextCursor:      resources.NextCursor,
		HasPreviousPage: resources.HasPreviousPage,
		PreviousCursor:  resources.PreviousCursor,
		TotalCount:      resources.TotalCount,
	}

	for _, result := range resources.Resources {
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

		output.Conditions = append(output.Conditions, resource)
	}

	return &output, nil
}

// SearchPatientAllergyIntolerance searches for a patient's FHIR allergy intolerance using patient ID
func (fh StoreImpl) SearchPatientAllergyIntolerance(_ context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
	params := map[string]interface{}{
		"patient": patientReference,
	}

	resources, err := fh.Dataset.SearchFHIRResource(allergyIntoleranceResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	allergyOutput := domain.PagedFHIRAllergy{
		Allergies:       []domain.FHIRAllergyIntolerance{},
		HasNextPage:     resources.HasNextPage,
		NextCursor:      resources.NextCursor,
		HasPreviousPage: resources.HasPreviousPage,
		PreviousCursor:  resources.PreviousCursor,
		TotalCount:      resources.TotalCount,
	}

	for _, resource := range resources.Resources {
		var allergyIntolerance domain.FHIRAllergyIntolerance

		resourceBs, err := json.Marshal(resource)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &allergyIntolerance)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", allergyIntoleranceResourceType, err)
		}

		allergyOutput.Allergies = append(allergyOutput.Allergies, allergyIntolerance)
	}

	return &allergyOutput, nil
}

// UpdateFHIRCondition updates a FHIRCondition instance
// The resource must have its ID set.
func (fh StoreImpl) UpdateFHIRCondition(_ context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
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
func (fh StoreImpl) GetFHIREncounter(_ context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
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
func (fh StoreImpl) SearchFHIREncounter(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
	resources, err := fh.Dataset.SearchFHIRResource(encounterResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	encounterOutput := domain.PagedFHIREncounter{
		Encounters:      []domain.FHIREncounter{},
		HasNextPage:     resources.HasNextPage,
		NextCursor:      resources.NextCursor,
		HasPreviousPage: resources.HasPreviousPage,
		PreviousCursor:  resources.PreviousCursor,
		TotalCount:      resources.TotalCount,
	}

	for _, result := range resources.Resources {
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

		encounterOutput.Encounters = append(encounterOutput.Encounters, resource)
	}

	return &encounterOutput, nil
}

// SearchFHIREncounterAllData provides a search API for a FHIREncounter and all other resources that reference the encounter
func (fh StoreImpl) SearchFHIREncounterAllData(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
	resources, err := fh.Dataset.SearchFHIRResource(encounterResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	encounterAllDataOutput := domain.PagedFHIRResource{
		Resources:       resources.Resources,
		HasNextPage:     resources.HasNextPage,
		NextCursor:      resources.NextCursor,
		HasPreviousPage: resources.HasPreviousPage,
		PreviousCursor:  resources.PreviousCursor,
		TotalCount:      resources.TotalCount,
	}

	return &encounterAllDataOutput, nil
}

// SearchFHIRMedicationRequest provides a search API for FHIRMedicationRequest
func (fh StoreImpl) SearchFHIRMedicationRequest(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationRequestRelayConnection, error) {
	output := domain.FHIRMedicationRequestRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(medicationRequestResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	for _, result := range resources.Resources {
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
func (fh StoreImpl) CreateFHIRMedicationRequest(_ context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
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
// The resource must have its ID set.
func (fh StoreImpl) UpdateFHIRMedicationRequest(_ context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
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
func (fh StoreImpl) DeleteFHIRMedicationRequest(_ context.Context, id string) (bool, error) {
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
func (fh StoreImpl) SearchFHIRObservation(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
	resources, err := fh.Dataset.SearchFHIRResource(observationResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	observationOutput := domain.PagedFHIRObservations{
		Observations:    []domain.FHIRObservation{},
		HasNextPage:     resources.HasNextPage,
		NextCursor:      resources.NextCursor,
		HasPreviousPage: resources.HasPreviousPage,
		PreviousCursor:  resources.PreviousCursor,
		TotalCount:      resources.TotalCount,
	}

	for _, result := range resources.Resources {
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

		observationOutput.Observations = append(observationOutput.Observations, resource)
	}

	return &observationOutput, nil
}

// CreateFHIRObservation creates a FHIRObservation instance
func (fh StoreImpl) CreateFHIRObservation(_ context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", observationResourceType, err)
	}

	resource := &domain.FHIRObservation{}

	err = fh.Dataset.CreateFHIRResource(observationResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", observationResourceType, err)
	}

	return resource, nil
}

// DeleteFHIRObservation deletes the FHIRObservation identified by the passed ID
func (fh StoreImpl) DeleteFHIRObservation(_ context.Context, id string) (bool, error) {
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
func (fh StoreImpl) GetFHIRPatient(_ context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
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
func (fh StoreImpl) DeleteFHIRPatient(_ context.Context, id string) (bool, error) {
	patientEverythingBs, err := fh.Dataset.GetFHIRPatientAllData(id, nil)
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
	// i.e. Observations, Medication Request etc
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

		case organizationResource:
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
	// as it ensures ResourceType that refers to the encounter is not found
	if err = fh.DeleteFHIRResourceType(encounters); err != nil {
		return false, err
	}

	// Thirdly, delete the episodes of care. This will bring no conflict
	// as it ensures Encounter that refers to the EpisodeOfCare is not found
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
func (fh StoreImpl) DeleteFHIRResourceType(results []map[string]string) error {
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
func (fh StoreImpl) DeleteFHIRServiceRequest(_ context.Context, id string) (bool, error) {
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
func (fh StoreImpl) CreateFHIRMedicationStatement(_ context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
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
func (fh StoreImpl) CreateFHIRMedication(_ context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
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

// CreateFHIRMedia creates a FHIR media resource
func (fh StoreImpl) CreateFHIRMedia(_ context.Context, input domain.FHIRMedia) (*domain.FHIRMedia, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, err
	}

	resource := &domain.FHIRMedia{}

	err = fh.Dataset.CreateFHIRResource(mediaResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", mediaResourceType, err)
	}

	return resource, nil
}

// SearchFHIRMedicationStatement used to search for a fhir medication statement
func (fh StoreImpl) SearchFHIRMedicationStatement(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
	output := domain.FHIRMedicationStatementRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(medicationStatementResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	for _, result := range resources.Resources {
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
func (fh StoreImpl) CreateFHIRPatient(_ context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
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
func (fh StoreImpl) PatchFHIRPatient(_ context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", patientResourceType, err)
	}

	resource := &domain.FHIRPatient{}

	err = fh.Dataset.PatchFHIRResource(patientResourceType, id, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to patch %s resource: %w", patientResourceType, err)
	}

	return resource, nil
}

// PatchFHIREpisodeOfCare patches a FHIR episode of care
func (fh StoreImpl) PatchFHIREpisodeOfCare(_ context.Context, id string, input domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", episodeOfCareResourceType, err)
	}

	resource := &domain.FHIREpisodeOfCare{}

	err = fh.Dataset.PatchFHIRResource(episodeOfCareResourceType, id, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to patch %s resource: %w", episodeOfCareResourceType, err)
	}

	return resource, nil
}

// UpdateFHIREpisodeOfCare updates a fhir episode of care
func (fh StoreImpl) UpdateFHIREpisodeOfCare(_ context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error) {
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
func (fh StoreImpl) SearchFHIRPatient(_ context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error) {
	params := map[string]interface{}{
		"_content": searchParams,
	}

	resources, err := fh.Dataset.SearchFHIRResource(patientResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	output := domain.PatientConnection{}

	for _, resource := range resources.Resources {
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

// GetFHIRComposition retrieves instances of FHIRComposition by ID
func (fh StoreImpl) GetFHIRComposition(_ context.Context, id string) (*domain.FHIRCompositionRelayPayload, error) {
	resource := &domain.FHIRComposition{}

	err := fh.Dataset.GetFHIRResource(compositionResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", compositionResourceType, id, err)
	}

	payload := &domain.FHIRCompositionRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// PatchFHIRComposition is used to patch a composition resource
func (fh StoreImpl) PatchFHIRComposition(_ context.Context, id string, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", compositionResourceType, err)
	}

	resource := &domain.FHIRComposition{}

	err = fh.Dataset.PatchFHIRResource(compositionResourceType, id, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to patch %s resource: %w", compositionResourceType, err)
	}

	return resource, nil
}

// GetFHIRObservation retrieves instances of FHIRObservation by ID
func (fh StoreImpl) GetFHIRObservation(_ context.Context, id string) (*domain.FHIRObservationRelayPayload, error) {
	resource := &domain.FHIRObservation{}

	err := fh.Dataset.GetFHIRResource(observationResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", observationResourceType, id, err)
	}

	payload := &domain.FHIRObservationRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// PatchFHIRObservation is used to patch an observation resource
func (fh StoreImpl) PatchFHIRObservation(_ context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", observationResourceType, err)
	}

	resource := &domain.FHIRObservation{}

	err = fh.Dataset.PatchFHIRResource(observationResourceType, id, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to patch %s resource: %w", observationResourceType, err)
	}

	return resource, nil
}

// ListFHIRQuestionnaire is used to list questionnaire resource using the name or the title of the resource.
func (fh StoreImpl) ListFHIRQuestionnaire(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRQuestionnaires, error) {
	results, err := fh.Dataset.SearchFHIRResource(questionnaireResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	questionnaireOutput := domain.PagedFHIRQuestionnaires{
		Questionnaires:  []domain.FHIRQuestionnaire{},
		HasNextPage:     results.HasNextPage,
		NextCursor:      results.NextCursor,
		HasPreviousPage: results.HasPreviousPage,
		PreviousCursor:  results.PreviousCursor,
		TotalCount:      results.TotalCount,
	}

	for _, result := range results.Resources {
		var questionnaire domain.FHIRQuestionnaire

		resourceBytes, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal resource to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBytes, &questionnaire)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal resource: %w", err)
		}

		questionnaireOutput.Questionnaires = append(questionnaireOutput.Questionnaires, questionnaire)
	}

	return &questionnaireOutput, nil
}

// CreateFHIRQuestionnaire is used to create a FHIR Questionnaire resource
func (fh StoreImpl) CreateFHIRQuestionnaire(_ context.Context, input *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", questionnaireResourceType, err)
	}

	resource := &domain.FHIRQuestionnaire{}

	err = fh.Dataset.CreateFHIRResource(questionnaireResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", questionnaireResourceType, err)
	}

	return resource, nil
}

// CreateFHIRConsent creates a FHIRConsent instance
func (fh StoreImpl) CreateFHIRConsent(_ context.Context, input domain.FHIRConsent) (*domain.FHIRConsent, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", consentResourceType, err)
	}

	resource := &domain.FHIRConsent{}

	err = fh.Dataset.CreateFHIRResource(consentResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %w", consentResourceType, err)
	}

	return resource, nil
}

// CreateFHIRQuestionnaireResponse is used to create a FHIR Questionnaire response resource
func (fh StoreImpl) CreateFHIRQuestionnaireResponse(_ context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", questionnaireResponseResourceType, err)
	}

	resource := &domain.FHIRQuestionnaireResponse{}

	err = fh.Dataset.CreateFHIRResource(questionnaireResponseResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", questionnaireResponseResourceType, err)
	}

	return resource, nil
}

// CreateFHIRRiskAssessment creates a RiskAssessment on FHIR
// The RiskAssessment resource represents an assessment of the likely outcome(s) for a patient's health over
// a period of time, considering various factors.
func (fh StoreImpl) CreateFHIRRiskAssessment(_ context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", riskAssessmentResourceType, err)
	}

	resource := &domain.FHIRRiskAssessment{}

	err = fh.Dataset.CreateFHIRResource(riskAssessmentResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", riskAssessmentResourceType, err)
	}

	return &domain.FHIRRiskAssessmentRelayPayload{
		Resource: resource,
	}, nil
}

// SearchFHIRRiskAssessment searches for a fhir risk assessment
func (fh StoreImpl) SearchFHIRRiskAssessment(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRRiskAssessmentRelayConnection, error) {
	output := domain.FHIRRiskAssessmentRelayConnection{}

	resources, err := fh.Dataset.SearchFHIRResource(riskAssessmentResourceType, params, tenant, pagination)
	if err != nil {
		return nil, err
	}

	for _, result := range resources.Resources {
		var resource domain.FHIRRiskAssessment

		resourceBs, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %w", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %w", riskAssessmentResourceType, err)
		}

		output.Edges = append(output.Edges, &domain.FHIRRiskAssessmentRelayEdge{
			Node: &resource,
		})
	}

	return &output, nil
}

// GetFHIRQuestionnaire retrieves instances of FHIRQuestionnaire by ID
func (fh StoreImpl) GetFHIRQuestionnaire(_ context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
	resource := &domain.FHIRQuestionnaire{}

	err := fh.Dataset.GetFHIRResource(questionnaireResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", questionnaireResourceType, id, err)
	}

	payload := &domain.FHIRQuestionnaireRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// GetFHIRQuestionnaireResponse retrieves an instance of FHIRQuestionnaireResponse by ID
func (fh StoreImpl) GetFHIRQuestionnaireResponse(_ context.Context, id string) (*domain.FHIRQuestionnaireResponseRelayPayload, error) {
	resource := &domain.FHIRQuestionnaireResponse{}

	err := fh.Dataset.GetFHIRResource(questionnaireResponseResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", questionnaireResponseResourceType, id, err)
	}

	payload := &domain.FHIRQuestionnaireResponseRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// CreateFHIRDiagnosticReport is used to create a diagnostic report resource for a patient
func (fh StoreImpl) CreateFHIRDiagnosticReport(_ context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %w", diagnosticReportResourceType, err)
	}

	resource := &domain.FHIRDiagnosticReport{}

	err = fh.Dataset.CreateFHIRResource(diagnosticReportResourceType, payload, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to create %s resource: %w", diagnosticReportResourceType, err)
	}

	return resource, nil
}

// GetFHIRPatientEverything is used to retrieve all patient related information
func (fh StoreImpl) GetFHIRPatientEverything(ctx context.Context, id string, params map[string]interface{}) (*domain.PagedFHIRResource, error) {
	patientEverythingBs, err := fh.Dataset.GetFHIRPatientAllData(id, params)
	if err != nil {
		return nil, fmt.Errorf("unable to get patient's compartment: %w", err)
	}

	respMap := make(map[string]interface{})

	err = json.Unmarshal(patientEverythingBs, &respMap)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal patient everything")
	}

	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
	for _, k := range mandatoryKeys {
		_, found := respMap[k]
		if !found {
			return nil, fmt.Errorf("server error: mandatory search result key %s not found", k)
		}
	}

	resourceType, ok := respMap["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("server error: the resourceType is not a string")
	}

	if resourceType != "Bundle" {
		return nil, fmt.Errorf(
			"server error: the resourceType value is not 'Bundle' as expected")
	}

	resultType, ok := respMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("server error: the search result type value is not a string")
	}

	if resultType != "searchset" {
		return nil, fmt.Errorf("server error: the type value is not 'searchset' as expected")
	}

	response := domain.PagedFHIRResource{
		Resources:       []map[string]interface{}{},
		HasNextPage:     false,
		NextCursor:      "",
		HasPreviousPage: false,
		PreviousCursor:  "",
		TotalCount:      0,
	}

	response.TotalCount = int(respMap["total"].(float64))

	respEntries := respMap["entry"]
	if respEntries == nil {
		return &response, nil
	}

	entries, ok := respEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"server error: entries is not a list of maps, it is: %T", respEntries)
	}

	for _, en := range entries {
		entry, ok := en.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(
				"server error: expected each entry to be map, they are %T instead", en)
		}

		expectedKeys := []string{"fullUrl", "resource"}
		for _, k := range expectedKeys {
			_, found := entry[k]
			if !found {
				return nil, fmt.Errorf("server error: FHIR search entry does not have key '%s'", k)
			}
		}

		resource, ok := entry["resource"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("server error: result entry %#v is not a map", entry["resource"])
		}

		response.Resources = append(response.Resources, resource)
	}

	linksEntries := respMap["link"]
	if linksEntries == nil {
		return &response, nil
	}

	links, ok := linksEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"server error: entries is not a list of maps, it is: %T", linksEntries)
	}

	for _, en := range links {
		link, ok := en.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(
				"server error: expected each link to be map, they are %T instead", en)
		}

		if link["relation"].(string) == "next" {
			u, err := url.Parse(link["url"].(string))
			if err != nil {
				return nil, fmt.Errorf("server error: cannot parse url in link: %w", err)
			}

			params, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				return nil, fmt.Errorf("server error: cannot parse url params in link: %w", err)
			}

			cursor := params["_page_token"][0]

			response.HasNextPage = true
			response.NextCursor = cursor
		} else if link["relation"].(string) == "previous" {
			u, err := url.Parse(link["url"].(string))
			if err != nil {
				return nil, fmt.Errorf("server error: cannot parse url in link: %w", err)
			}

			params, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				return nil, fmt.Errorf("server error: cannot parse url params in link: %w", err)
			}

			cursor := params["_page_token"][0]

			response.HasPreviousPage = true
			response.PreviousCursor = cursor
		}
	}

	return &response, nil
}

// GetFHIRServiceRequest retrieves a FHIR service request using its primary ID
func (fh StoreImpl) GetFHIRServiceRequest(_ context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
	resource := &domain.FHIRServiceRequest{}

	err := fh.Dataset.GetFHIRResource(serviceRequestResourceType, id, resource)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %w", serviceRequestResourceType, id, err)
	}

	payload := &domain.FHIRServiceRequestRelayPayload{
		Resource: resource,
	}

	return payload, nil
}

// CreateFHIRSubscription is responsible for creating a subscription resource in FHIR repository
func (fh StoreImpl) CreateFHIRSubscription(_ context.Context, subscription *domain.FHIRSubscriptionInput) (*domain.FHIRSubscription, error) {
	payload, err := converterandformatter.StructToMap(subscription)
	if err != nil {
		return nil, fmt.Errorf("unable to convert subscription input into a map: %w", err)
	}

	fhirSubscription := &domain.FHIRSubscription{}

	err = fh.Dataset.CreateFHIRResource(subscriptionResourceType, payload, fhirSubscription)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to create episode of care resource: %w", err)
	}

	return fhirSubscription, nil
}
