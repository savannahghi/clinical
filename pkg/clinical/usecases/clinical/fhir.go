package clinical

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
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
		ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error)
	GetORCreateOrganization(ctx context.Context, providerSladeCode string) (*string, error)
	GetOrganization(ctx context.Context, providerSladeCode string) (*string, error)
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	CreateOrganization(ctx context.Context, providerSladeCode string) (*string, error)
	SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error)
	FindOrganizationByID(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error)
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

// Encounters returns encounters that belong to the indicated patient.
//
// The patientReference should be a [string] in the format "Patient/<patient resource ID>".
func (c *UseCasesClinicalImpl) Encounters(
	ctx context.Context,
	patientReference string,
	status *domain.EncounterStatusEnum,
) ([]*domain.FHIREncounter, error) {

	return c.infrastructure.FHIR.Encounters(ctx, patientReference, status)
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (c *UseCasesClinicalImpl) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return c.infrastructure.FHIR.SearchFHIREpisodeOfCare(ctx, params)
}

// CreateEpisodeOfCare is the final common pathway for creation of episodes of care.
func (c *UseCasesClinicalImpl) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
	return c.infrastructure.FHIR.CreateEpisodeOfCare(ctx, episode)
}

// CreateFHIRCondition creates a FHIRCondition instance
func (c *UseCasesClinicalImpl) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	return c.infrastructure.FHIR.CreateFHIRCondition(ctx, input)
}

// CreateFHIROrganization creates a FHIROrganization instance
func (c *UseCasesClinicalImpl) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	organizationRelayPayload, err := c.infrastructure.FHIR.CreateFHIROrganization(ctx, input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}
	return organizationRelayPayload, nil
}

// OpenOrganizationEpisodes return all organization specific open episodes
func (c *UseCasesClinicalImpl) OpenOrganizationEpisodes(
	ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error) {

	organizationID, err := c.GetORCreateOrganization(ctx, providerSladeCode)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"internal server error in retrieving service provider : %v", err)
	}
	organizationReference := fmt.Sprintf("Organization/%s", *organizationID)
	searchParams := url.Values{}
	searchParams.Add("status", domain.EpisodeOfCareStatusEnumActive.String())
	searchParams.Add("organization", organizationReference)
	return c.SearchEpisodesByParam(ctx, searchParams)
}

// CreateOrganization creates an organization given ist provider code
func (c *UseCasesClinicalImpl) CreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {

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
	createdOrganization, err := c.CreateFHIROrganization(ctx, organizationInput)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}
	organisationID := createdOrganization.Resource.ID
	return organisationID, nil
}

// SearchFHIROrganization provides a search API for FHIROrganization
func (c *UseCasesClinicalImpl) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return c.infrastructure.FHIR.SearchFHIROrganization(ctx, params)
}

// FindOrganizationByID finds and retrieves organization details using the specified organization ID
func (c *UseCasesClinicalImpl) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("organization ID is required")
	}
	return c.infrastructure.FHIR.FindOrganizationByID(ctx, organizationID)
}

// GetORCreateOrganization retrieve an organisation via its code if not found create a new one.
func (c *UseCasesClinicalImpl) GetORCreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	retrievedOrg, err := c.GetOrganization(ctx, providerSladeCode)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf(
			"internal server error in getting organisation : %v", err)
	}
	if retrievedOrg == nil {
		createdOrg, err := c.CreateOrganization(ctx, providerSladeCode)
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
func (c *UseCasesClinicalImpl) GetOrganization(ctx context.Context, providerSladeCode string) (*string, error) {

	searchParam := map[string]interface{}{
		"identifier": providerSladeCode,
	}
	organization, err := c.SearchFHIROrganization(ctx, searchParam)
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
func (c *UseCasesClinicalImpl) SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
	return c.infrastructure.FHIR.SearchEpisodesByParam(ctx, searchParams)
}

// OpenEpisodes returns the IDs of a patient's open episodes
func (c *UseCasesClinicalImpl) OpenEpisodes(
	ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	return c.infrastructure.FHIR.OpenEpisodes(ctx, patientReference)
}

// HasOpenEpisode determines if a patient has an open episode
func (c *UseCasesClinicalImpl) HasOpenEpisode(
	ctx context.Context,
	patient domain.FHIRPatient,
) (bool, error) {
	return c.infrastructure.FHIR.HasOpenEpisode(ctx, patient)
}

// FHIRHeaders composes suitable FHIR headers, with authentication and content
// type already set
func (c *UseCasesClinicalImpl) FHIRHeaders() (http.Header, error) {
	return c.infrastructure.FHIR.FHIRHeaders()
}

// CreateFHIREncounter creates a FHIREncounter instance
func (c *UseCasesClinicalImpl) CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {

	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	return c.infrastructure.FHIR.CreateFHIREncounter(ctx, input)
}

// GetFHIREpisodeOfCare retrieves instances of FHIREpisodeOfCare by ID
func (c *UseCasesClinicalImpl) GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
	return c.infrastructure.FHIR.GetFHIREpisodeOfCare(ctx, id)
}

// StartEncounter starts an encounter within an episode of care
func (c *UseCasesClinicalImpl) StartEncounter(
	ctx context.Context, episodeID string) (string, error) {
	return c.infrastructure.FHIR.StartEncounter(ctx, episodeID)
}

// StartEpisodeByOtp starts a patient OTP verified episode
func (c *UseCasesClinicalImpl) StartEpisodeByOtp(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {

	normalized, err := converterandformatter.NormalizeMSISDN(input.Msisdn)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to normalize phone number: %w", err)
	}

	organizationID, err := c.GetORCreateOrganization(ctx, input.ProviderCode)
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
	return c.CreateEpisodeOfCare(ctx, ep)
}

// UpgradeEpisode starts a patient OTP verified episode
func (c *UseCasesClinicalImpl) UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
	// retrieve and validate the episode
	episode, err := c.GetActiveEpisode(ctx, input.EpisodeID)
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
	encounters, err := c.Encounters(ctx, *episode.Patient.Reference, nil)
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

	episodeOfCare, err := c.infrastructure.FHIR.UpdateFHIREpisodeOfCare(ctx, *episode.ID, payload)
	if err != nil {
		return nil, err
	}

	return &domain.EpisodeOfCarePayload{
		EpisodeOfCare: episodeOfCare,
		TotalVisits:   len(encounters),
	}, nil
}

// SearchEpisodeEncounter returns all encounters in a visit
func (c *UseCasesClinicalImpl) SearchEpisodeEncounter(
	ctx context.Context,
	episodeReference string,
) (*domain.FHIREncounterRelayConnection, error) {
	return c.infrastructure.FHIR.SearchEpisodeEncounter(ctx, episodeReference)
}

// EndEncounter ends an encounter
func (c *UseCasesClinicalImpl) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	return c.infrastructure.FHIR.EndEncounter(ctx, encounterID)
}

// EndEpisode ends an episode of care by patching it's status to "finished"
func (c *UseCasesClinicalImpl) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	return c.infrastructure.FHIR.EndEpisode(ctx, episodeID)
}

// GetActiveEpisode returns any ACTIVE episode that has to the indicated ID
func (c *UseCasesClinicalImpl) GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
	return c.infrastructure.FHIR.GetActiveEpisode(ctx, episodeID)
}

// SearchFHIRServiceRequest provides a search API for FHIRServiceRequest
func (c *UseCasesClinicalImpl) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return c.infrastructure.FHIR.SearchFHIRServiceRequest(ctx, params)
}

// CreateFHIRServiceRequest creates a FHIRServiceRequest instance
func (c *UseCasesClinicalImpl) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	return c.infrastructure.FHIR.CreateFHIRServiceRequest(ctx, input)
}

// SearchFHIRAllergyIntolerance provides a search API for FHIRAllergyIntolerance
func (c *UseCasesClinicalImpl) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	return c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, params)
}

// CreateFHIRAllergyIntolerance creates a FHIRAllergyIntolerance instance
func (c *UseCasesClinicalImpl) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {

	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	return c.infrastructure.FHIR.CreateFHIRAllergyIntolerance(ctx, input)
}

// UpdateFHIRAllergyIntolerance updates a FHIRAllergyIntolerance instance
// The resource must have it's ID set.
func (c *UseCasesClinicalImpl) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return c.infrastructure.FHIR.UpdateFHIRAllergyIntolerance(ctx, input)
}

// SearchFHIRComposition provides a search API for FHIRComposition
func (c *UseCasesClinicalImpl) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return c.infrastructure.FHIR.SearchFHIRComposition(ctx, params)
}

// CreateFHIRComposition creates a FHIRComposition instance
func (c *UseCasesClinicalImpl) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {

	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	return c.infrastructure.FHIR.CreateFHIRComposition(ctx, input)
}

// UpdateFHIRComposition updates a FHIRComposition instance
// The resource must have it's ID set.
func (c *UseCasesClinicalImpl) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return c.infrastructure.FHIR.UpdateFHIRComposition(ctx, input)
}

// DeleteFHIRComposition deletes the FHIRComposition identified by the supplied ID
func (c *UseCasesClinicalImpl) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	return c.infrastructure.FHIR.DeleteFHIRComposition(ctx, id)
}

// SearchFHIRCondition provides a search API for FHIRCondition
func (c *UseCasesClinicalImpl) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return c.infrastructure.FHIR.SearchFHIRCondition(ctx, params)
}

// UpdateFHIRCondition updates a FHIRCondition instance
// The resource must have it's ID set.
func (c *UseCasesClinicalImpl) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return c.infrastructure.FHIR.UpdateFHIRCondition(ctx, input)
}

// GetFHIREncounter retrieves instances of FHIREncounter by ID
func (c *UseCasesClinicalImpl) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	return c.infrastructure.FHIR.GetFHIREncounter(ctx, id)
}

// SearchFHIREncounter provides a search API for FHIREncounter
func (c *UseCasesClinicalImpl) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return c.infrastructure.FHIR.SearchFHIREncounter(ctx, params)
}

// SearchFHIRMedicationRequest provides a search API for FHIRMedicationRequest
func (c *UseCasesClinicalImpl) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return c.infrastructure.FHIR.SearchFHIRMedicationRequest(ctx, params)
}

// CreateFHIRMedicationRequest creates a FHIRMedicationRequest instance
func (c *UseCasesClinicalImpl) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {

	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	return c.infrastructure.FHIR.CreateFHIRMedicationRequest(ctx, input)
}

// UpdateFHIRMedicationRequest updates a FHIRMedicationRequest instance
// The resource must have it's ID set.
func (c *UseCasesClinicalImpl) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return c.infrastructure.FHIR.UpdateFHIRMedicationRequest(ctx, input)
}

// DeleteFHIRMedicationRequest deletes the FHIRMedicationRequest identified by the supplied ID
func (c *UseCasesClinicalImpl) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	return c.infrastructure.FHIR.DeleteFHIRMedicationRequest(ctx, id)
}

// SearchFHIRObservation provides a search API for FHIRObservation
func (c *UseCasesClinicalImpl) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return c.infrastructure.FHIR.SearchFHIRObservation(ctx, params)
}

// CreateFHIRObservation creates a FHIRObservation instance
func (c *UseCasesClinicalImpl) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {

	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	return c.infrastructure.FHIR.CreateFHIRObservation(ctx, input)
}

// DeleteFHIRObservation deletes the FHIRObservation identified by the passed ID
func (c *UseCasesClinicalImpl) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	return c.infrastructure.FHIR.DeleteFHIRObservation(ctx, id)
}

// GetFHIRPatient retrieves instances of FHIRPatient by ID
func (c *UseCasesClinicalImpl) GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
	return c.infrastructure.FHIR.GetFHIRPatient(ctx, id)
}

// DeleteFHIRPatient deletes the FHIRPatient identified by the supplied ID
func (c *UseCasesClinicalImpl) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	return c.infrastructure.FHIR.DeleteFHIRPatient(ctx, id)
}

// DeleteFHIRResourceType takes a ResourceType and ID and deletes them from FHIR
func (c *UseCasesClinicalImpl) DeleteFHIRResourceType(results []map[string]string) error {
	return c.infrastructure.FHIR.DeleteFHIRResourceType(results)
}

// DeleteFHIRServiceRequest deletes the FHIRServiceRequest identified by the supplied ID
func (c *UseCasesClinicalImpl) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	return c.infrastructure.FHIR.DeleteFHIRServiceRequest(ctx, id)
}

// CreateFHIRMedicationStatement creates a new FHIR Medication statement instance
func (c *UseCasesClinicalImpl) CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {

	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	return c.infrastructure.FHIR.CreateFHIRMedicationStatement(ctx, input)
}

// CreateFHIRMedication creates a new FHIR Medication instance
func (c *UseCasesClinicalImpl) CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {

	tenantTags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	input.Meta = domain.FHIRMetaInput{
		Tag: tenantTags,
	}

	return c.infrastructure.FHIR.CreateFHIRMedication(ctx, input)
}

// SearchFHIRMedicationStatement used to search for a fhir medication statement
func (c *UseCasesClinicalImpl) SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	return c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, params)
}
