package patient

import (
	"context"
	"net/url"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// UseCasePatient defines the patient usecase
type UseCasePatient interface {
	ProblemSummary(ctx context.Context, patientID string) ([]string, error)
	VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error)
	PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error)
	CreateEpisodeOfCare(ctx context.Context, ep domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error)
	StartEpisodeByOtp(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error)
	UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error)
	StartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error)
	GetOrganization(ctx context.Context, providerSladeCode string) (*string, error)
	CreateOrganization(ctx context.Context, providerSladeCode string) (*string, error)
	GetORCreateOrganization(ctx context.Context, providerSladeCode string) (*string, error)
	OpenOrganizationEpisodes(ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error)
	SearchEpisodeEncounter(
		ctx context.Context,
		episodeReference string,
	) (*domain.FHIREncounterRelayConnection, error)
	StartEncounter(ctx context.Context, episodeID string) (string, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	EndEpisode(ctx context.Context, episodeID string) (bool, error)
	CreatePatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error)
	FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error)
	PatientSearch(ctx context.Context, search string) (*domain.PatientConnection, error)
	UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error)
	AddNhif(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error)
	GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error)
	SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error)
	OpenEpisodes(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error)
	HasOpenEpisode(
		ctx context.Context,
		patient domain.FHIRPatient,
	) (bool, error)
	FindPatientsByMSISDN(ctx context.Context, msisdn string) (*domain.PatientConnection, error)
	RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	CheckPatientExistenceUsingPhoneNumber(ctx context.Context, patientInput domain.SimplePatientRegistrationInput) (bool, error)
	SimplePatientRegistrationInputToPatientInput(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.FHIRPatientInput, error)
	SendPatientWelcomeEmail(ctx context.Context, emailaddress string) error
	SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error)
	CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	DeleteFHIRComposition(ctx context.Context, id string) (bool, error)
	SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error)
	CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error)
	SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error)
	CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error)
	GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error)
	SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error)
	SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error)
	CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error)
	SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error)
	CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error)
	DeleteFHIRObservation(ctx context.Context, id string) (bool, error)
	SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error)
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error)
	CreateUpdatePatientExtraInformation(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error)
	DeleteFHIRPatient(ctx context.Context, id string) (bool, error)
	DeleteFHIRResourceType(results []map[string]string) error
	AllergySummary(ctx context.Context, patientID string) ([]string, error)
	DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error)
}

// UseCasePatientImpl represents the usecase implementation object
type UseCasePatientImpl struct{}

// NewUseCasePatient initializes a new patient usecase
func NewUseCasePatient() UseCasePatient {
	return &UseCasePatientImpl{}
}

//TODO: implement the methods after adding the infra interactor

// ProblemSummary returns a short list of the patient's active and confirmed
// problems (by name).
func (u *UseCasePatientImpl) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	return nil, nil
}

// VisitSummary returns a narrative friendly display of the data that has
// been associated with a single visit
func (u *UseCasePatientImpl) VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
	return nil, nil
}

// PatientTimelineWithCount returns the patient's visit note timeline (a list of
// narratives that are sorted with the most recent one first), while
// respecting the approval level AND limiting the number
func (u *UseCasePatientImpl) PatientTimelineWithCount(
	ctx context.Context,
	episodeID string,
	count int) ([]map[string]interface{}, error) {
	return nil, nil
}

// CreateEpisodeOfCare is the final common pathway for creation of episodes of
// care.
func (u *UseCasePatientImpl) CreateEpisodeOfCare(
	ctx context.Context,
	ep domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
	return nil, nil
}

// Encounters returns encounters that belong to the indicated patient.
//
// The patientReference should be a [string] in the format "Patient/<patient resource ID>".
func (u *UseCasePatientImpl) Encounters(
	ctx context.Context,
	patientReference string,
	status *domain.EncounterStatusEnum,
) ([]*domain.FHIREncounter, error) {
	return nil, nil
}

// StartEpisodeByOtp starts a patient OTP verified episode
func (u *UseCasePatientImpl) StartEpisodeByOtp(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	return nil, nil
}

// UpgradeEpisode starts a patient OTP verified episode
func (u *UseCasePatientImpl) UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
	return nil, nil
}

// StartEpisodeByBreakGlass starts an emergency episode
func (u *UseCasePatientImpl) StartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	return nil, nil
}

// GetOrganization retrieves an organization given its code
func (u *UseCasePatientImpl) GetOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	return nil, nil
}

// CreateOrganization creates an organization given ist provider code
func (u *UseCasePatientImpl) CreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	return nil, nil
}

// GetORCreateOrganization retrieve an organisation via its code if not found create a new one.
func (u *UseCasePatientImpl) GetORCreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	return nil, nil
}

// OpenOrganizationEpisodes return all organization specific open episodes
func (u *UseCasePatientImpl) OpenOrganizationEpisodes(ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error) {
	return nil, nil
}

// SearchEpisodeEncounter returns all encounters in a visit
func (u *UseCasePatientImpl) SearchEpisodeEncounter(
	ctx context.Context,
	episodeReference string,
) (*domain.FHIREncounterRelayConnection, error) {
	return nil, nil
}

// StartEncounter starts an encounter within an episode of care
func (u *UseCasePatientImpl) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	return "", nil
}

// EndEncounter ends an encounter
func (u *UseCasePatientImpl) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	return false, nil
}

// EndEpisode ends an episode of care by patching it's status to "finished"
func (u *UseCasePatientImpl) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	return false, nil
}

// CreatePatient creates or updates a patient record on FHIR
func (u *UseCasePatientImpl) CreatePatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	return nil, nil
}

// FindPatientByID retrieves a single patient by their ID
func (u *UseCasePatientImpl) FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error) {
	return nil, nil
}

// PatientSearch searches for a patient by identifiers and names
func (u *UseCasePatientImpl) PatientSearch(ctx context.Context, search string) (*domain.PatientConnection, error) {
	return nil, nil
}

// UpdatePatient patches a patient record with fresh data.
// It updates elements that are set and ignores the ones that are nil.
func (u *UseCasePatientImpl) UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return nil, nil
}

// AddNextOfKin patches a patient with next of kin
func (u *UseCasePatientImpl) AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
	return nil, nil
}

// AddNhif patches a patient with NHIF details
func (u *UseCasePatientImpl) AddNhif(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
	return nil, nil
}

// GetActiveEpisode returns any ACTIVE episode that has to the indicated ID
func (u *UseCasePatientImpl) GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
	return nil, nil
}

// SearchEpisodesByParam search episodes by params
func (u *UseCasePatientImpl) SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
	return nil, nil
}

// OpenEpisodes returns the IDs of a patient's open episodes
func (u *UseCasePatientImpl) OpenEpisodes(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	return nil, nil
}

// HasOpenEpisode determines if a patient has an open episode
func (u *UseCasePatientImpl) HasOpenEpisode(
	ctx context.Context,
	patient domain.FHIRPatient,
) (bool, error) {
	return false, nil
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
func (u *UseCasePatientImpl) FindPatientsByMSISDN(ctx context.Context, msisdn string) (*domain.PatientConnection, error) {
	return nil, nil
}

// RegisterPatient implements simple patient registration
func (u *UseCasePatientImpl) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return nil, nil
}

// CheckPatientExistenceUsingPhoneNumber checks whether a patient with the phone number they're trying to register with exists
func (u *UseCasePatientImpl) CheckPatientExistenceUsingPhoneNumber(ctx context.Context, patientInput domain.SimplePatientRegistrationInput) (bool, error) {
	return false, nil
}

// SimplePatientRegistrationInputToPatientInput transforms a patient input into
// a
func (u *UseCasePatientImpl) SimplePatientRegistrationInputToPatientInput(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.FHIRPatientInput, error) {
	return nil, nil
}

// SendPatientWelcomeEmail will send a welcome email to the practitioner
func (u *UseCasePatientImpl) SendPatientWelcomeEmail(ctx context.Context, emailaddress string) error {
	return nil
}

// SearchFHIRServiceRequest provides a search API for FHIRServiceRequest
func (u *UseCasePatientImpl) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
	return nil, nil
}

// CreateFHIRServiceRequest creates a FHIRServiceRequest instance
func (u *UseCasePatientImpl) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	return nil, nil
}

// DeleteFHIRServiceRequest deletes the FHIRServiceRequest identified by the supplied ID
func (u *UseCasePatientImpl) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	return false, nil
}

// SearchFHIRAllergyIntolerance provides a search API for FHIRAllergyIntolerance
func (u *UseCasePatientImpl) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	return nil, nil
}

// CreateFHIRAllergyIntolerance creates a FHIRAllergyIntolerance instance
func (u *UseCasePatientImpl) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return nil, nil
}

// UpdateFHIRAllergyIntolerance updates a FHIRAllergyIntolerance instance
// The resource must have it's ID set.
func (u *UseCasePatientImpl) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return nil, nil
}

// SearchFHIRComposition provides a search API for FHIRComposition
func (u *UseCasePatientImpl) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
	return nil, nil
}

// CreateFHIRComposition creates a FHIRComposition instance
func (u *UseCasePatientImpl) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return nil, nil
}

// UpdateFHIRComposition updates a FHIRComposition instance
// The resource must have it's ID set.
func (u *UseCasePatientImpl) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return nil, nil
}

// DeleteFHIRComposition deletes the FHIRComposition identified by the supplied ID
func (u *UseCasePatientImpl) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	return false, nil
}

// SearchFHIRCondition provides a search API for FHIRCondition
func (u *UseCasePatientImpl) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	return nil, nil
}

// CreateFHIRCondition creates a FHIRCondition instance
func (u *UseCasePatientImpl) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return nil, nil
}

// UpdateFHIRCondition updates a FHIRCondition instance
// The resource must have it's ID set.
func (u *UseCasePatientImpl) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return nil, nil
}

// GetFHIREncounter retrieves instances of FHIREncounter by ID
func (u *UseCasePatientImpl) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	return nil, nil
}

// SearchFHIREncounter provides a search API for FHIREncounter
func (u *UseCasePatientImpl) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
	return nil, nil
}

// CreateFHIREncounter creates a FHIREncounter instance
func (u *UseCasePatientImpl) CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
	return nil, nil
}

// GetFHIREpisodeOfCare retrieves instances of FHIREpisodeOfCare by ID
func (u *UseCasePatientImpl) GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
	return nil, nil
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (u *UseCasePatientImpl) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	return nil, nil
}

// SearchFHIRMedicationRequest provides a search API for FHIRMedicationRequest
func (u *UseCasePatientImpl) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
	return nil, nil
}

// CreateFHIRMedicationRequest creates a FHIRMedicationRequest instance
func (u *UseCasePatientImpl) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return nil, nil
}

// UpdateFHIRMedicationRequest updates a FHIRMedicationRequest instance
// The resource must have it's ID set.
func (u *UseCasePatientImpl) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return nil, nil
}

// DeleteFHIRMedicationRequest deletes the FHIRMedicationRequest identified by the supplied ID
func (u *UseCasePatientImpl) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	return false, nil
}

// SearchFHIRObservation provides a search API for FHIRObservation
func (u *UseCasePatientImpl) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	return nil, nil
}

// CreateFHIRObservation creates a FHIRObservation instance
func (u *UseCasePatientImpl) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	return nil, nil
}

// DeleteFHIRObservation deletes the FHIRObservation identified by the passed ID
func (u *UseCasePatientImpl) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	return false, nil
}

// SearchFHIROrganization provides a search API for FHIROrganization
func (u *UseCasePatientImpl) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {
	return nil, nil
}

// CreateFHIROrganization creates a FHIROrganization instance
func (u *UseCasePatientImpl) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	return nil, nil
}

// GetFHIRPatient retrieves instances of FHIRPatient by ID
func (u *UseCasePatientImpl) GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
	return nil, nil
}

// CreateUpdatePatientExtraInformation updates a patient's extra info
func (u *UseCasePatientImpl) CreateUpdatePatientExtraInformation(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
	return false, nil
}

// DeleteFHIRPatient deletes the FHIRPatient identified by the supplied ID
func (u *UseCasePatientImpl) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	return false, nil
}

// DeleteFHIRResourceType takes a ResourceType and ID and deletes them from FHIR
func (u *UseCasePatientImpl) DeleteFHIRResourceType(results []map[string]string) error {
	return nil
}

// AllergySummary returns a short list of the patient's active and confirmed
// allergies (by name)
func (u *UseCasePatientImpl) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	return nil, nil
}

// DeleteFHIRPatientByPhone delete's a patient's FHIR compartment
// given their phone number
func (u *UseCasePatientImpl) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	return false, nil
}
