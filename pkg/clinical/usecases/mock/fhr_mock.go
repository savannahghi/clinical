package mock

import (
	"context"
	"io"
	"net/url"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FHIRMock contains all mock methods
type FHIRMock struct {
	CreateEpisodeOfCareFn          func(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error)
	CreateFHIRConditionFn          func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	OpenOrganizationEpisodesFn     func(ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error)
	GetORCreateOrganizationFn      func(ctx context.Context, providerSladeCode string) (*string, error)
	GetOrganizationFn              func(ctx context.Context, providerSladeCode string) (*string, error)
	CreateFHIROrganizationFn       func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	CreateOrganizationFn           func(ctx context.Context, providerSladeCode string) (*string, error)
	SearchFHIROrganizationFn       func(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error)
	POSTRequestFn                  func(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error)
	SearchEpisodesByParamFn        func(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error)
	HasOpenEpisodeFn               func(ctx context.Context, patient domain.FHIRPatient) (bool, error)
	OpenEpisodesFn                 func(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error)
	CreateFHIREncounterFn          func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error)
	GetFHIREpisodeOfCareFn         func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error)
	EncountersFn                   func(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum) ([]*domain.FHIREncounter, error)
	SearchFHIREpisodeOfCareFn      func(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error)
	StartEncounterFn               func(ctx context.Context, episodeID string) (string, error)
	StartEpisodeByOtpFn            func(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error)
	UpgradeEpisodeFn               func(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error)
	SearchEpisodeEncounterFn       func(ctx context.Context, episodeReference string) (*domain.FHIREncounterRelayConnection, error)
	EndEncounterFn                 func(ctx context.Context, encounterID string) (bool, error)
	EndEpisodeFn                   func(ctx context.Context, episodeID string) (bool, error)
	GetActiveEpisodeFn             func(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error)
	SearchFHIRServiceRequestFn     func(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error)
	CreateFHIRServiceRequestFn     func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error)
	SearchFHIRAllergyIntoleranceFn func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error)
	CreateFHIRAllergyIntoleranceFn func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	UpdateFHIRAllergyIntoleranceFn func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	SearchFHIRCompositionFn        func(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error)
	CreateFHIRCompositionFn        func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	UpdateFHIRCompositionFn        func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	DeleteFHIRCompositionFn        func(ctx context.Context, id string) (bool, error)
	SearchFHIRConditionFn          func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error)
	UpdateFHIRConditionFn          func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	GetFHIREncounterFn             func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error)
	SearchFHIREncounterFn          func(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error)
	SearchFHIRMedicationRequestFn  func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error)
	CreateFHIRMedicationRequestFn  func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	UpdateFHIRMedicationRequestFn  func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	DeleteFHIRMedicationRequestFn  func(ctx context.Context, id string) (bool, error)
	SearchFHIRObservationFn        func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error)
	CreateFHIRObservationFn        func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error)
	DeleteFHIRObservationFn        func(ctx context.Context, id string) (bool, error)
	GetFHIRPatientFn               func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error)
	DeleteFHIRPatientFn            func(ctx context.Context, id string) (bool, error)
	DeleteFHIRResourceTypeFn       func(results []map[string]string) error
	DeleteFHIRServiceRequestFn     func(ctx context.Context, id string) (bool, error)
}

// CreateEpisodeOfCare is a mock implementation of CreateEpisodeOfCare method
func (fh *FHIRMock) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
	return fh.CreateEpisodeOfCareFn(ctx, episode)
}

// CreateFHIRCondition is a mock implementation of CreateFHIRCondition method
func (fh *FHIRMock) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return fh.CreateFHIRConditionFn(ctx, input)
}

// OpenOrganizationEpisodes is a mock implementation of OpenOrganizationEpisodes method
func (fh *FHIRMock) OpenOrganizationEpisodes(ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.OpenOrganizationEpisodesFn(ctx, providerSladeCode)
}

// GetORCreateOrganization is a mock implementation of GetORCreateOrganization method
func (fh *FHIRMock) GetORCreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	return fh.GetORCreateOrganizationFn(ctx, providerSladeCode)
}

// GetOrganization is a mock implementation of GetOrganization method
func (fh *FHIRMock) GetOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	return fh.GetOrganizationFn(ctx, providerSladeCode)
}

// CreateFHIROrganization is a mock implementation of CreateFHIROrganization method
func (fh *FHIRMock) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.CreateFHIROrganizationFn(ctx, input)
}

// CreateOrganization is a mock implementation of CreateOrganization method
func (fh *FHIRMock) CreateOrganization(ctx context.Context, providerSladeCode string) (*string, error) {
	return fh.CreateOrganizationFn(ctx, providerSladeCode)
}

// SearchFHIROrganization is a mock implementation of SearchFHIROrganization method
func (fh *FHIRMock) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {
	return fh.SearchFHIROrganizationFn(ctx, params)
}

// POSTRequest is a mock implementation of POSTRequest method
func (fh *FHIRMock) POSTRequest(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
	return fh.POSTRequestFn(resourceName, path, params, body)
}

// SearchEpisodesByParam is a mock implementation of SearchEpisodesByParam method
func (fh *FHIRMock) SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.SearchEpisodesByParamFn(ctx, searchParams)
}

// HasOpenEpisode is a mock implementation of HasOpenEpisode method
func (fh *FHIRMock) HasOpenEpisode(ctx context.Context, patient domain.FHIRPatient) (bool, error) {
	return fh.HasOpenEpisodeFn(ctx, patient)
}

// OpenEpisodes is a mock implementation of OpenEpisodes method
func (fh *FHIRMock) OpenEpisodes(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.OpenEpisodesFn(ctx, patientReference)
}

// CreateFHIREncounter is a mock implementation of CreateFHIREncounter method
func (fh *FHIRMock) CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
	return fh.CreateFHIREncounterFn(ctx, input)
}

// GetFHIREpisodeOfCare is a mock implementation of GetFHIREpisodeOfCare method
func (fh *FHIRMock) GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
	return fh.GetFHIREpisodeOfCareFn(ctx, id)
}

// Encounters is a mock implementation of Encounters method
func (fh *FHIRMock) Encounters(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum) ([]*domain.FHIREncounter, error) {
	return fh.EncountersFn(ctx, patientReference, status)
}

// SearchFHIREpisodeOfCare is a mock implementation of SearchFHIREpisodeOfCare method
func (fh *FHIRMock) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	return fh.SearchFHIREpisodeOfCareFn(ctx, params)
}

// StartEncounter is a mock implementation of StartEncounter method
func (fh *FHIRMock) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	return fh.StartEncounterFn(ctx, episodeID)
}

// StartEpisodeByOtp is a mock implementation of StartEpisodeByOtp method
func (fh *FHIRMock) StartEpisodeByOtp(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	return fh.StartEpisodeByOtpFn(ctx, input)
}

// UpgradeEpisode is a mock implementation of UpgradeEpisode method
func (fh *FHIRMock) UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
	return fh.UpgradeEpisodeFn(ctx, input)
}

// SearchEpisodeEncounter is a mock implementation of SearchEpisodeEncounter method
func (fh *FHIRMock) SearchEpisodeEncounter(ctx context.Context, episodeReference string) (*domain.FHIREncounterRelayConnection, error) {
	return fh.SearchEpisodeEncounterFn(ctx, episodeReference)
}

// EndEncounter is a mock implementation of EndEncounter method
func (fh *FHIRMock) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	return fh.EndEncounterFn(ctx, encounterID)
}

// EndEpisode is a mock implementation of EndEpisode method
func (fh *FHIRMock) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	return fh.EndEpisodeFn(ctx, episodeID)
}

// GetActiveEpisode is a mock implementation of GetActiveEpisode method
func (fh *FHIRMock) GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
	return fh.GetActiveEpisodeFn(ctx, episodeID)
}

// SearchFHIRServiceRequest is a mock implementation of SearchFHIRServiceRequest method
func (fh *FHIRMock) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
	return fh.SearchFHIRServiceRequestFn(ctx, params)
}

// CreateFHIRServiceRequest is a mock implementation of CreateFHIRServiceRequest method
func (fh *FHIRMock) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	return fh.CreateFHIRServiceRequestFn(ctx, input)
}

// SearchFHIRAllergyIntolerance is a mock implementation of SearchFHIRAllergyIntolerance method
func (fh *FHIRMock) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	return fh.SearchFHIRAllergyIntoleranceFn(ctx, params)
}

// CreateFHIRAllergyIntolerance is a mock implementation of CreateFHIRAllergyIntolerance method
func (fh *FHIRMock) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.CreateFHIRAllergyIntoleranceFn(ctx, input)
}

// UpdateFHIRAllergyIntolerance is a mock implementation of UpdateFHIRAllergyIntolerance method
func (fh *FHIRMock) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.UpdateFHIRAllergyIntoleranceFn(ctx, input)
}

// SearchFHIRComposition is a mock implementation of SearchFHIRComposition method
func (fh *FHIRMock) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
	return fh.SearchFHIRCompositionFn(ctx, params)
}

// CreateFHIRComposition is a mock implementation of CreateFHIRComposition method
func (fh *FHIRMock) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return fh.CreateFHIRCompositionFn(ctx, input)
}

// UpdateFHIRComposition is a mock implementation of UpdateFHIRComposition method
func (fh *FHIRMock) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return fh.UpdateFHIRCompositionFn(ctx, input)
}

// DeleteFHIRComposition is a mock implementation of DeleteFHIRComposition method
func (fh *FHIRMock) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	return fh.DeleteFHIRCompositionFn(ctx, id)
}

// SearchFHIRCondition is a mock implementation of SearchFHIRCondition method
func (fh *FHIRMock) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	return fh.SearchFHIRConditionFn(ctx, params)
}

// UpdateFHIRCondition is a mock implementation of UpdateFHIRCondition method
func (fh *FHIRMock) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return fh.UpdateFHIRConditionFn(ctx, input)
}

// GetFHIREncounter is a mock implementation of GetFHIREncounter method
func (fh *FHIRMock) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	return fh.GetFHIREncounterFn(ctx, id)
}

// SearchFHIREncounter is a mock implementation of SearchFHIREncounter method
func (fh *FHIRMock) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
	return fh.SearchFHIREncounterFn(ctx, params)
}

// SearchFHIRMedicationRequest is a mock implementation of SearchFHIRMedicationRequest method
func (fh *FHIRMock) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
	return fh.SearchFHIRMedicationRequestFn(ctx, params)
}

// CreateFHIRMedicationRequest is a mock implementation of CreateFHIRMedicationRequest method
func (fh *FHIRMock) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return fh.CreateFHIRMedicationRequestFn(ctx, input)
}

// UpdateFHIRMedicationRequest is a mock implementation of UpdateFHIRMedicationRequest method
func (fh *FHIRMock) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return fh.UpdateFHIRMedicationRequestFn(ctx, input)
}

// DeleteFHIRMedicationRequest is a mock implementation of DeleteFHIRMedicationRequest method
func (fh *FHIRMock) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	return fh.DeleteFHIRMedicationRequestFn(ctx, id)
}

// SearchFHIRObservation is a mock implementation of SearchFHIRObservation method
func (fh *FHIRMock) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	return fh.SearchFHIRObservationFn(ctx, params)
}

// CreateFHIRObservation is a mock implementation of CreateFHIRObservation method
func (fh *FHIRMock) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	return fh.CreateFHIRObservationFn(ctx, input)
}

// DeleteFHIRObservation is a mock implementation of DeleteFHIRObservation method
func (fh *FHIRMock) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	return fh.DeleteFHIRObservationFn(ctx, id)
}

// GetFHIRPatient is a mock implementation of GetFHIRPatient method
func (fh *FHIRMock) GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
	return fh.GetFHIRPatientFn(ctx, id)
}

// DeleteFHIRPatient is a mock implementation of DeleteFHIRPatient method
func (fh *FHIRMock) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	return fh.DeleteFHIRPatientFn(ctx, id)
}

// DeleteFHIRResourceType is a mock implementation of DeleteFHIRResourceType method
func (fh *FHIRMock) DeleteFHIRResourceType(results []map[string]string) error {
	return fh.DeleteFHIRResourceTypeFn(results)
}

// DeleteFHIRServiceRequest is a mock implementation of DeleteFHIRServiceRequest method
func (fh *FHIRMock) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	return fh.DeleteFHIRServiceRequestFn(ctx, id)
}
