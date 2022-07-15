package fhir

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/converterandformatter"
)

const (
	fullAccessLevel = "FULL_ACCESS"
)

// UseCasesFHIR represents all the FHIR business logic
type UseCasesFHIR interface {
	CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error)
	CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	OpenOrganizationEpisodes(
		ctx context.Context, MFLCode string) ([]*domain.FHIREpisodeOfCare, error)
	GetORCreateOrganization(ctx context.Context, org domain.FHIROrganizationInput) (*string, error)
	GetOrganization(ctx context.Context, MFLCode string) (*string, error)
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	CreateOrganization(ctx context.Context, org domain.FHIROrganizationInput) (*string, error)
	SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error)
	FindOrganizationByID(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error)
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
	CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)
	CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error)
	SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error)
}

// UseCasesFHIRImpl represents the FHIR usecase implementation
type UseCasesFHIRImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewUseCasesFHIRImpl initializes the new FHIR implementation
func NewUseCasesFHIRImpl(infra infrastructure.Infrastructure) UseCasesFHIR {
	return &UseCasesFHIRImpl{
		infrastructure: infra,
	}
}

// Encounters returns encounters that belong to the indicated patient.
//
// The patientReference should be a [string] in the format "Patient/<patient resource ID>".
func (fh UseCasesFHIRImpl) Encounters(
	ctx context.Context,
	patientReference string,
	status *domain.EncounterStatusEnum,
) ([]*domain.FHIREncounter, error) {

	return fh.infrastructure.FHIR.Encounters(ctx, patientReference, status)
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (fh UseCasesFHIRImpl) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return fh.infrastructure.FHIR.SearchFHIREpisodeOfCare(ctx, params)
}

//CreateEpisodeOfCare is the final common pathway for creation of episodes of care.
func (fh UseCasesFHIRImpl) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
	return fh.infrastructure.FHIR.CreateEpisodeOfCare(ctx, episode)
}

// CreateFHIRCondition creates a FHIRCondition instance
func (fh UseCasesFHIRImpl) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIRCondition(ctx, input)
}

// CreateFHIROrganization creates a FHIROrganization instance
func (fh UseCasesFHIRImpl) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIROrganization(ctx, input)
}

// OpenOrganizationEpisodes return all organization specific open episodes
func (fh UseCasesFHIRImpl) OpenOrganizationEpisodes(
	ctx context.Context, MFLCode string) ([]*domain.FHIREpisodeOfCare, error) {
	org := domain.FHIROrganizationInput{
		Identifier: []*domain.FHIRIdentifierInput{
			{
				Value: MFLCode,
			},
		},
	}

	organizationID, err := fh.GetORCreateOrganization(ctx, org)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"internal server error in retrieving service provider : %v", err)
	}
	organizationReference := fmt.Sprintf("Organization/%s", *organizationID)
	searchParams := url.Values{}
	searchParams.Add("status", domain.EpisodeOfCareStatusEnumActive.String())
	searchParams.Add("organization", organizationReference)
	return fh.SearchEpisodesByParam(ctx, searchParams)
}

// CreateOrganization creates an organization given ist provider code
func (fh UseCasesFHIRImpl) CreateOrganization(ctx context.Context, org domain.FHIROrganizationInput) (*string, error) {
	createdOrganization, err := fh.CreateFHIROrganization(ctx, org)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}
	organisationID := createdOrganization.Resource.ID
	return organisationID, nil
}

// SearchFHIROrganization provides a search API for FHIROrganization
func (fh UseCasesFHIRImpl) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return fh.infrastructure.FHIR.SearchFHIROrganization(ctx, params)
}

// FindOrganizationByID finds and retrieves organization details using the specified organization ID
func (fh UseCasesFHIRImpl) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("organization ID is required")
	}
	return fh.infrastructure.FHIR.FindOrganizationByID(ctx, organizationID)
}

// GetORCreateOrganization retrieve an organisation via its code if not found create a new one.
func (fh UseCasesFHIRImpl) GetORCreateOrganization(ctx context.Context, org domain.FHIROrganizationInput) (*string, error) {
	retrievedOrg, err := fh.GetOrganization(ctx, org.Identifier[0].Value)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"internal server error in getting organisation : %v", err)
	}
	if retrievedOrg == nil {
		createdOrg, err := fh.CreateOrganization(ctx, org)
		if err != nil {
			utils.ReportErrorToSentry(err)
			return nil, fmt.Errorf(
				"internal server error in creating organisation : %v", err)
		}
		return createdOrg, nil
	}
	return retrievedOrg, nil
}

// GetOrganization retrieves an organization given its code
func (fh UseCasesFHIRImpl) GetOrganization(ctx context.Context, MFLCode string) (*string, error) {

	searchParam := map[string]interface{}{
		"identifier": MFLCode,
	}
	organization, err := fh.SearchFHIROrganization(ctx, searchParam)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}
	if organization.Edges == nil {
		return nil, nil
	}
	ORGID := organization.Edges[0].Node.ID
	return ORGID, nil
}

// SearchEpisodesByParam search episodes by params
func (fh UseCasesFHIRImpl) SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.infrastructure.FHIR.SearchEpisodesByParam(ctx, searchParams)
}

// POSTRequest is used to manually compose POST requests to the FHIR service
//
// - `resourceName` is a FHIR resource name e.g "Patient"
// - `path` is a sub-path e.g `_search` under a resource
// - `params` should be query params, sent as `url.Values`
func (fh UseCasesFHIRImpl) POSTRequest(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
	return fh.infrastructure.FHIRRepo.POSTRequest(resourceName, path, params, body)
}

// OpenEpisodes returns the IDs of a patient's open episodes
func (fh UseCasesFHIRImpl) OpenEpisodes(
	ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.infrastructure.FHIR.OpenEpisodes(ctx, patientReference)
}

// HasOpenEpisode determines if a patient has an open episode
func (fh UseCasesFHIRImpl) HasOpenEpisode(
	ctx context.Context,
	patient domain.FHIRPatient,
) (bool, error) {
	return fh.infrastructure.FHIR.HasOpenEpisode(ctx, patient)
}

// FHIRHeaders composes suitable FHIR headers, with authentication and content
// type already set
func (fh UseCasesFHIRImpl) FHIRHeaders() (http.Header, error) {
	return fh.infrastructure.FHIRRepo.FHIRHeaders()
}

// CreateFHIREncounter creates a FHIREncounter instance
func (fh UseCasesFHIRImpl) CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIREncounter(ctx, input)
}

// GetFHIREpisodeOfCare retrieves instances of FHIREpisodeOfCare by ID
func (fh UseCasesFHIRImpl) GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
	return fh.infrastructure.FHIR.GetFHIREpisodeOfCare(ctx, id)
}

// StartEncounter starts an encounter within an episode of care
func (fh *UseCasesFHIRImpl) StartEncounter(
	ctx context.Context, episodeID string) (string, error) {
	return fh.infrastructure.FHIR.StartEncounter(ctx, episodeID)
}

// StartEpisodeByOtp starts a patient OTP verified episode
func (fh *UseCasesFHIRImpl) StartEpisodeByOtp(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {

	normalized, err := converterandformatter.NormalizeMSISDN(input.Msisdn)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to normalize phone number: %w", err)
	}

	p := domain.FHIROrganizationInput{
		Identifier: []*domain.FHIRIdentifierInput{
			{
				Value: input.MFLCode,
			},
		},
	}
	organizationID, err := fh.GetORCreateOrganization(ctx, p)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"internal server error in retrieving service provider : %v", err)
	}
	ep := helpers.ComposeOneHealthEpisodeOfCare(
		*normalized,
		input.FullAccess,
		*organizationID,
		input.MFLCode,
		input.PatientID,
	)
	return fh.CreateEpisodeOfCare(ctx, ep)
}

// UpgradeEpisode starts a patient OTP verified episode
func (fh *UseCasesFHIRImpl) UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
	// retrieve and validate the episode
	episode, err := fh.GetActiveEpisode(ctx, input.EpisodeID)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
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

	// patch the episode status
	episode.Type = []*domain.FHIRCodeableConcept{
		{Text: fullAccessLevel},
	}
	payload, err := converterandformatter.StructToMap(episode)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to turn episode of care input into a map: %v", err)
	}

	_, err = fh.infrastructure.FHIRRepo.UpdateFHIRResource(
		"EpisodeOfCare", *episode.ID, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"unable to update episode of care resource: %v", err)
	}
	return &domain.EpisodeOfCarePayload{
		EpisodeOfCare: episode,
		TotalVisits:   len(encounters),
	}, nil
}

// SearchEpisodeEncounter returns all encounters in a visit
func (fh *UseCasesFHIRImpl) SearchEpisodeEncounter(
	ctx context.Context,
	episodeReference string,
) (*domain.FHIREncounterRelayConnection, error) {
	return fh.infrastructure.FHIR.SearchEpisodeEncounter(ctx, episodeReference)
}

// EndEncounter ends an encounter
func (fh *UseCasesFHIRImpl) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	return fh.infrastructure.FHIR.EndEncounter(ctx, encounterID)
}

// EndEpisode ends an episode of care by patching it's status to "finished"
func (fh *UseCasesFHIRImpl) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	return fh.infrastructure.FHIR.EndEpisode(ctx, episodeID)
}

// GetActiveEpisode returns any ACTIVE episode that has to the indicated ID
func (fh *UseCasesFHIRImpl) GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
	return fh.infrastructure.FHIR.GetActiveEpisode(ctx, episodeID)
}

// SearchFHIRServiceRequest provides a search API for FHIRServiceRequest
func (fh *UseCasesFHIRImpl) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return fh.infrastructure.FHIR.SearchFHIRServiceRequest(ctx, params)
}

// CreateFHIRServiceRequest creates a FHIRServiceRequest instance
func (fh *UseCasesFHIRImpl) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIRServiceRequest(ctx, input)
}

// SearchFHIRAllergyIntolerance provides a search API for FHIRAllergyIntolerance
func (fh *UseCasesFHIRImpl) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	return fh.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, params)
}

// CreateFHIRAllergyIntolerance creates a FHIRAllergyIntolerance instance
func (fh *UseCasesFHIRImpl) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIRAllergyIntolerance(ctx, input)
}

// UpdateFHIRAllergyIntolerance updates a FHIRAllergyIntolerance instance
// The resource must have it's ID set.
func (fh *UseCasesFHIRImpl) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.infrastructure.FHIR.UpdateFHIRAllergyIntolerance(ctx, input)
}

// SearchFHIRComposition provides a search API for FHIRComposition
func (fh *UseCasesFHIRImpl) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return fh.infrastructure.FHIR.SearchFHIRComposition(ctx, params)
}

// CreateFHIRComposition creates a FHIRComposition instance
func (fh *UseCasesFHIRImpl) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIRComposition(ctx, input)
}

// UpdateFHIRComposition updates a FHIRComposition instance
// The resource must have it's ID set.
func (fh *UseCasesFHIRImpl) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return fh.infrastructure.FHIR.UpdateFHIRComposition(ctx, input)
}

// DeleteFHIRComposition deletes the FHIRComposition identified by the supplied ID
func (fh *UseCasesFHIRImpl) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	return fh.infrastructure.FHIR.DeleteFHIRComposition(ctx, id)
}

// SearchFHIRCondition provides a search API for FHIRCondition
func (fh *UseCasesFHIRImpl) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return fh.infrastructure.FHIR.SearchFHIRCondition(ctx, params)
}

// UpdateFHIRCondition updates a FHIRCondition instance
// The resource must have it's ID set.
func (fh *UseCasesFHIRImpl) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return fh.infrastructure.FHIR.UpdateFHIRCondition(ctx, input)
}

// GetFHIREncounter retrieves instances of FHIREncounter by ID
func (fh *UseCasesFHIRImpl) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	return fh.infrastructure.FHIR.GetFHIREncounter(ctx, id)
}

// SearchFHIREncounter provides a search API for FHIREncounter
func (fh *UseCasesFHIRImpl) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return fh.infrastructure.FHIR.SearchFHIREncounter(ctx, params)
}

// SearchFHIRMedicationRequest provides a search API for FHIRMedicationRequest
func (fh *UseCasesFHIRImpl) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return fh.infrastructure.FHIR.SearchFHIRMedicationRequest(ctx, params)
}

// CreateFHIRMedicationRequest creates a FHIRMedicationRequest instance
func (fh *UseCasesFHIRImpl) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIRMedicationRequest(ctx, input)
}

// UpdateFHIRMedicationRequest updates a FHIRMedicationRequest instance
// The resource must have it's ID set.
func (fh *UseCasesFHIRImpl) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return fh.infrastructure.FHIR.UpdateFHIRMedicationRequest(ctx, input)
}

// DeleteFHIRMedicationRequest deletes the FHIRMedicationRequest identified by the supplied ID
func (fh *UseCasesFHIRImpl) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	return fh.infrastructure.FHIR.DeleteFHIRMedicationRequest(ctx, id)
}

// SearchFHIRObservation provides a search API for FHIRObservation
func (fh *UseCasesFHIRImpl) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return fh.infrastructure.FHIR.SearchFHIRObservation(ctx, params)
}

// CreateFHIRObservation creates a FHIRObservation instance
func (fh *UseCasesFHIRImpl) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIRObservation(ctx, input)
}

// DeleteFHIRObservation deletes the FHIRObservation identified by the passed ID
func (fh *UseCasesFHIRImpl) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	return fh.infrastructure.FHIR.DeleteFHIRObservation(ctx, id)
}

// GetFHIRPatient retrieves instances of FHIRPatient by ID
func (fh *UseCasesFHIRImpl) GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
	return fh.infrastructure.FHIR.GetFHIRPatient(ctx, id)
}

// DeleteFHIRPatient deletes the FHIRPatient identified by the supplied ID
func (fh *UseCasesFHIRImpl) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	return fh.infrastructure.FHIR.DeleteFHIRPatient(ctx, id)
}

// DeleteFHIRResourceType takes a ResourceType and ID and deletes them from FHIR
func (fh *UseCasesFHIRImpl) DeleteFHIRResourceType(results []map[string]string) error {
	return fh.infrastructure.FHIR.DeleteFHIRResourceType(results)
}

// DeleteFHIRServiceRequest deletes the FHIRServiceRequest identified by the supplied ID
func (fh *UseCasesFHIRImpl) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	return fh.infrastructure.FHIR.DeleteFHIRServiceRequest(ctx, id)
}

// CreateFHIRMedicationStatement creates a new FHIR Medication statement instance
func (fh *UseCasesFHIRImpl) CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIRMedicationStatement(ctx, input)
}

// CreateFHIRMedication creates a new FHIR Medication instance
func (fh *UseCasesFHIRImpl) CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
	return fh.infrastructure.FHIR.CreateFHIRMedication(ctx, input)
}

// SearchFHIRMedicationStatement used to search for a fhir medication statement
func (fh *UseCasesFHIRImpl) SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return fh.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, params)
}
