package mock

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// FHIRUsecaseMock struct implements mocks of FHIR methods.
type FHIRUsecaseMock struct {
	MockCreateEpisodeOfCareFn    func(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error)
	MockPOSTRequestFn            func(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error)
	MockSearchFHIRConditionFn    func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error)
	MockCreateFHIRConditionFn    func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	MockCreateFHIROrganizationFn func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	MockSearchFHIROrganizationFn func(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error)
	MockFindOrganizationByIDFn   func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error)
	MockSearchEpisodesByParamFn  func(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error)
	MockHasOpenEpisodeFn         func(
		ctx context.Context,
		patient domain.FHIRPatient,
	) (bool, error)
	MockOpenEpisodesFn func(
		ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error)
	MockCreateFHIREncounterFn           func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error)
	MockGetFHIREpisodeOfCareFn          func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error)
	MockEncountersFn                    func(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum) ([]*domain.FHIREncounter, error)
	MockSearchFHIREpisodeOfCareFn       func(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error)
	MockStartEncounterFn                func(ctx context.Context, episodeID string) (string, error)
	MockUpgradeEpisodeFn                func(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error)
	MockSearchEpisodeEncounterFn        func(ctx context.Context, episodeReference string) (*domain.FHIREncounterRelayConnection, error)
	MockEndEncounterFn                  func(ctx context.Context, encounterID string) (bool, error)
	MockEndEpisodeFn                    func(ctx context.Context, episodeID string) (bool, error)
	MockGetActiveEpisodeFn              func(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error)
	MockSearchFHIRServiceRequestFn      func(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error)
	MockCreateFHIRServiceRequestFn      func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error)
	MockSearchFHIRAllergyIntoleranceFn  func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error)
	MockCreateFHIRAllergyIntoleranceFn  func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	MockUpdateFHIRAllergyIntoleranceFn  func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	MockSearchFHIRCompositionFn         func(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error)
	MockCreateFHIRCompositionFn         func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	MockUpdateFHIRCompositionFn         func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	MockDeleteFHIRCompositionFn         func(ctx context.Context, id string) (bool, error)
	MockUpdateFHIRConditionFn           func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	MockGetFHIREncounterFn              func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error)
	MockSearchFHIREncounterFn           func(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error)
	MockSearchFHIRMedicationRequestFn   func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error)
	MockCreateFHIRMedicationRequestFn   func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	MockUpdateFHIRMedicationRequestFn   func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	MockDeleteFHIRMedicationRequestFn   func(ctx context.Context, id string) (bool, error)
	MockSearchFHIRObservationFn         func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error)
	MockCreateFHIRObservationFn         func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error)
	MockDeleteFHIRObservationFn         func(ctx context.Context, id string) (bool, error)
	MockGetFHIRPatientFn                func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error)
	MockDeleteFHIRPatientFn             func(ctx context.Context, id string) (bool, error)
	MockDeleteFHIRServiceRequestFn      func(ctx context.Context, id string) (bool, error)
	MockDeleteFHIRResourceTypeFn        func(results []map[string]string) error
	MockCreateFHIRMedicationStatementFn func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)
	MockCreateFHIRMedicationFn          func(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error)
	MockSearchFHIRMedicationStatementFn func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error)
	MockFHIRHeadersFn                   func() (http.Header, error)
	MockGetBearerTokenFn                func() (string, error)
}

// NewFHIRUsecaseMock initializes a new instance of FHIR mock
func NewFHIRUsecaseMock() *FHIRUsecaseMock {
	UUID := uuid.New().String()
	PatientRef := "Patient/1"
	OrgRef := "Organization/1"
	status := domain.EpisodeOfCareStatusEnumFinished
	return &FHIRUsecaseMock{
		MockCreateEpisodeOfCareFn: func(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
			return &domain.EpisodeOfCarePayload{
				EpisodeOfCare: &domain.FHIREpisodeOfCare{
					ID:            &UUID,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
					Status:        &(status),
					StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
					Type:          []*domain.FHIRCodeableConcept{},
					Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
					Patient: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReference{
						Reference: &OrgRef,
					},
					Period:          &domain.FHIRPeriod{},
					ReferralRequest: []*domain.FHIRReference{},
					CareManager:     &domain.FHIRReference{},
					Team:            []*domain.FHIRReference{},
					Account:         []*domain.FHIRReference{},
				},
				TotalVisits: 0,
			}, nil
		},
		MockCreateFHIRConditionFn: func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
			return &domain.FHIRConditionRelayPayload{
				Resource: &domain.FHIRCondition{
					ID: &UUID,
					Text: &domain.FHIRNarrative{
						ID:  new(string),
						Div: "",
					},
					Identifier: []*domain.FHIRIdentifier{},
					ClinicalStatus: &domain.FHIRCodeableConcept{
						ID:   new(string),
						Text: "",
					},
					VerificationStatus: &domain.FHIRCodeableConcept{},
					Category:           []*domain.FHIRCodeableConcept{},
					Severity:           &domain.FHIRCodeableConcept{},
					Code:               &domain.FHIRCodeableConcept{},
					BodySite:           []*domain.FHIRCodeableConcept{},
					Subject:            &domain.FHIRReference{},
					Encounter:          &domain.FHIRReference{},
					OnsetDateTime: &scalarutils.Date{
						Year:  2000,
						Month: 10,
						Day:   10,
					},
					OnsetAge:    &domain.FHIRAge{},
					OnsetPeriod: &domain.FHIRPeriod{},
					OnsetRange:  &domain.FHIRRange{},
					OnsetString: &UUID,
					AbatementDateTime: &scalarutils.Date{
						Year:  2000,
						Month: 10,
						Day:   10,
					},
					AbatementAge:    &domain.FHIRAge{},
					AbatementPeriod: &domain.FHIRPeriod{},
					AbatementRange:  &domain.FHIRRange{},
					AbatementString: &UUID,
					RecordedDate: &scalarutils.Date{
						Year:  2000,
						Month: 10,
						Day:   10,
					},
					Recorder: &domain.FHIRReference{},
					Asserter: &domain.FHIRReference{},
					Stage: []*domain.FHIRConditionStage{
						{
							ID:         new(string),
							Summary:    &domain.FHIRCodeableConcept{},
							Assessment: []*domain.FHIRReference{},
							Type:       &domain.FHIRCodeableConcept{},
						},
					},
					Evidence: []*domain.FHIRConditionEvidence{
						{
							ID:     new(string),
							Code:   &domain.FHIRCodeableConcept{},
							Detail: []*domain.FHIRReference{},
						},
					},
					Note: []*domain.FHIRAnnotation{
						{
							ID:              new(string),
							AuthorReference: &domain.FHIRReference{},
							AuthorString:    new(string),
						},
					},
				},
			}, nil
		},
		MockCreateFHIROrganizationFn: func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
			return &domain.FHIROrganizationRelayPayload{
				Resource: &domain.FHIROrganization{
					ID: &UUID,
				},
			}, nil
		},
		MockSearchFHIROrganizationFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {
			return &domain.FHIROrganizationRelayConnection{}, nil
		},
		MockPOSTRequestFn: func(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
			return []byte("m"), nil
		},
		MockSearchEpisodesByParamFn: func(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
			return []*domain.FHIREpisodeOfCare{}, nil
		},
		MockHasOpenEpisodeFn: func(ctx context.Context, patient domain.FHIRPatient) (bool, error) {
			return true, nil
		},
		MockOpenEpisodesFn: func(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
			return []*domain.FHIREpisodeOfCare{
				{
					ID:            &UUID,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
					Status:        &(status),
					StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
					Type:          []*domain.FHIRCodeableConcept{},
					Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
					Patient: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReference{
						Reference: &OrgRef,
					},
					Period:          &domain.FHIRPeriod{},
					ReferralRequest: []*domain.FHIRReference{},
					CareManager:     &domain.FHIRReference{},
					Team:            []*domain.FHIRReference{},
					Account:         []*domain.FHIRReference{},
				},
			}, nil
		},
		MockCreateFHIREncounterFn: func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
			return &domain.FHIREncounterRelayPayload{}, nil
		},
		MockGetFHIREpisodeOfCareFn: func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
			return &domain.FHIREpisodeOfCareRelayPayload{
				Resource: &domain.FHIREpisodeOfCare{
					ID:            &UUID,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
					Status:        &(status),
					StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
					Type:          []*domain.FHIRCodeableConcept{},
					Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
					Patient: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReference{
						Reference: &OrgRef,
					},
					Period:          &domain.FHIRPeriod{},
					ReferralRequest: []*domain.FHIRReference{},
					CareManager:     &domain.FHIRReference{},
					Team:            []*domain.FHIRReference{},
					Account:         []*domain.FHIRReference{},
				},
			}, nil
		},
		MockEncountersFn: func(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum) ([]*domain.FHIREncounter, error) {
			return []*domain.FHIREncounter{}, nil
		},
		MockSearchFHIREpisodeOfCareFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
			return &domain.FHIREpisodeOfCareRelayConnection{}, nil
		},
		MockStartEncounterFn: func(ctx context.Context, episodeID string) (string, error) {
			return "test-encounter", nil
		},
		MockUpgradeEpisodeFn: func(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
			return &domain.EpisodeOfCarePayload{}, nil
		},
		MockSearchEpisodeEncounterFn: func(ctx context.Context, episodeReference string) (*domain.FHIREncounterRelayConnection, error) {
			return &domain.FHIREncounterRelayConnection{}, nil
		},
		MockEndEncounterFn: func(ctx context.Context, encounterID string) (bool, error) {
			return true, nil
		},
		MockEndEpisodeFn: func(ctx context.Context, episodeID string) (bool, error) {
			return true, nil
		},
		MockGetActiveEpisodeFn: func(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
			return &domain.FHIREpisodeOfCare{}, nil
		},
		MockSearchFHIRServiceRequestFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
			return &domain.FHIRServiceRequestRelayConnection{}, nil
		},
		MockCreateFHIRServiceRequestFn: func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
			return &domain.FHIRServiceRequestRelayPayload{}, nil
		},
		MockSearchFHIRAllergyIntoleranceFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
			return &domain.FHIRAllergyIntoleranceRelayConnection{}, nil
		},
		MockCreateFHIRAllergyIntoleranceFn: func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
			return &domain.FHIRAllergyIntoleranceRelayPayload{}, nil
		},
		MockUpdateFHIRAllergyIntoleranceFn: func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
			return &domain.FHIRAllergyIntoleranceRelayPayload{}, nil
		},
		MockSearchFHIRCompositionFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
			return &domain.FHIRCompositionRelayConnection{}, nil
		},
		MockCreateFHIRCompositionFn: func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
			return &domain.FHIRCompositionRelayPayload{}, nil
		},
		MockUpdateFHIRCompositionFn: func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
			return &domain.FHIRCompositionRelayPayload{}, nil
		},
		MockDeleteFHIRCompositionFn: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},
		MockSearchFHIRConditionFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
			return &domain.FHIRConditionRelayConnection{}, nil
		},
		MockUpdateFHIRConditionFn: func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
			return &domain.FHIRConditionRelayPayload{}, nil
		},
		MockGetFHIREncounterFn: func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
			return &domain.FHIREncounterRelayPayload{}, nil
		},
		MockSearchFHIREncounterFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
			return &domain.FHIREncounterRelayConnection{}, nil
		},
		MockSearchFHIRMedicationRequestFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
			return &domain.FHIRMedicationRequestRelayConnection{}, nil
		},
		MockCreateFHIRMedicationRequestFn: func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
			return &domain.FHIRMedicationRequestRelayPayload{}, nil
		},
		MockUpdateFHIRMedicationRequestFn: func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
			return &domain.FHIRMedicationRequestRelayPayload{}, nil
		},
		MockDeleteFHIRMedicationRequestFn: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},
		MockSearchFHIRObservationFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
			return &domain.FHIRObservationRelayConnection{}, nil
		},
		MockCreateFHIRObservationFn: func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
			return &domain.FHIRObservationRelayPayload{}, nil
		},
		MockDeleteFHIRObservationFn: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},
		MockGetFHIRPatientFn: func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
			return &domain.FHIRPatientRelayPayload{}, nil
		},
		MockDeleteFHIRPatientFn: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},
		MockDeleteFHIRResourceTypeFn: func(results []map[string]string) error {
			return nil
		},
		MockDeleteFHIRServiceRequestFn: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},
		MockSearchFHIRMedicationStatementFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
			return &domain.FHIRMedicationStatementRelayConnection{}, nil
		},
		MockFindOrganizationByIDFn: func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
			return &domain.FHIROrganizationRelayPayload{}, nil
		},
		MockCreateFHIRMedicationStatementFn: func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
			return &domain.FHIRMedicationStatementRelayPayload{}, nil
		},
		MockCreateFHIRMedicationFn: func(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
			return &domain.FHIRMedicationRelayPayload{}, nil
		},
		MockFHIRHeadersFn: func() (http.Header, error) {
			return http.Header{
				"Authorization": []string{"Bearer " + uuid.NewString()},
			}, nil
		},
		MockGetBearerTokenFn: func() (string, error) {
			return fmt.Sprintf("Bearer %s", uuid.NewString()), nil
		},
	}
}

// GetBearerToken is a mock implementation of get bearer token method
func (fh *FHIRUsecaseMock) GetBearerToken() (string, error) {
	return fh.MockGetBearerTokenFn()
}

// CreateEpisodeOfCare is a mock implementation of CreateEpisodeOfCare method
func (fh *FHIRUsecaseMock) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
	return fh.MockCreateEpisodeOfCareFn(ctx, episode)
}

// CreateFHIRCondition is a mock implementation of CreateFHIRCondition method
func (fh *FHIRUsecaseMock) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return fh.MockCreateFHIRConditionFn(ctx, input)
}

// CreateFHIROrganization is a mock implementation of CreateFHIROrganization method
func (fh *FHIRUsecaseMock) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.MockCreateFHIROrganizationFn(ctx, input)
}

// SearchFHIROrganization is a mock implementation of SearchFHIROrganization method
func (fh *FHIRUsecaseMock) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {
	return fh.MockSearchFHIROrganizationFn(ctx, params)
}

// POSTRequest is a mock implementation of POSTRequest method
func (fh *FHIRUsecaseMock) POSTRequest(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
	return fh.MockPOSTRequestFn(resourceName, path, params, body)
}

// SearchEpisodesByParam is a mock implementation of SearchEpisodesByParam method
func (fh *FHIRUsecaseMock) SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.MockSearchEpisodesByParamFn(ctx, searchParams)
}

// HasOpenEpisode is a mock implementation of HasOpenEpisode method
func (fh *FHIRUsecaseMock) HasOpenEpisode(ctx context.Context, patient domain.FHIRPatient) (bool, error) {
	return fh.MockHasOpenEpisodeFn(ctx, patient)
}

// OpenEpisodes is a mock implementation of OpenEpisodes method
func (fh *FHIRUsecaseMock) OpenEpisodes(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.MockOpenEpisodesFn(ctx, patientReference)
}

// CreateFHIREncounter is a mock implementation of CreateFHIREncounter method
func (fh *FHIRUsecaseMock) CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
	return fh.MockCreateFHIREncounterFn(ctx, input)
}

// GetFHIREpisodeOfCare is a mock implementation of GetFHIREpisodeOfCare method
func (fh *FHIRUsecaseMock) GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
	return fh.MockGetFHIREpisodeOfCareFn(ctx, id)
}

// Encounters is a mock implementation of Encounters method
func (fh *FHIRUsecaseMock) Encounters(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum) ([]*domain.FHIREncounter, error) {
	return fh.MockEncountersFn(ctx, patientReference, status)
}

// SearchFHIREpisodeOfCare is a mock implementation of SearchFHIREpisodeOfCare method
func (fh *FHIRUsecaseMock) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	return fh.MockSearchFHIREpisodeOfCareFn(ctx, params)
}

// StartEncounter is a mock implementation of StartEncounter method
func (fh *FHIRUsecaseMock) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	return fh.MockStartEncounterFn(ctx, episodeID)
}

// UpgradeEpisode is a mock implementation of UpgradeEpisode method
func (fh *FHIRUsecaseMock) UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
	return fh.MockUpgradeEpisodeFn(ctx, input)
}

// SearchEpisodeEncounter is a mock implementation of SearchEpisodeEncounter method
func (fh *FHIRUsecaseMock) SearchEpisodeEncounter(ctx context.Context, episodeReference string) (*domain.FHIREncounterRelayConnection, error) {
	return fh.MockSearchEpisodeEncounterFn(ctx, episodeReference)
}

// EndEncounter is a mock implementation of EndEncounter method
func (fh *FHIRUsecaseMock) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	return fh.MockEndEncounterFn(ctx, encounterID)
}

// EndEpisode is a mock implementation of EndEpisode method
func (fh *FHIRUsecaseMock) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	return fh.MockEndEpisodeFn(ctx, episodeID)
}

// GetActiveEpisode is a mock implementation of GetActiveEpisode method
func (fh *FHIRUsecaseMock) GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
	return fh.MockGetActiveEpisodeFn(ctx, episodeID)
}

// SearchFHIRServiceRequest is a mock implementation of SearchFHIRServiceRequest method
func (fh *FHIRUsecaseMock) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
	return fh.MockSearchFHIRServiceRequestFn(ctx, params)
}

// CreateFHIRServiceRequest is a mock implementation of CreateFHIRServiceRequest method
func (fh *FHIRUsecaseMock) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	return fh.MockCreateFHIRServiceRequestFn(ctx, input)
}

// SearchFHIRAllergyIntolerance is a mock implementation of SearchFHIRAllergyIntolerance method
func (fh *FHIRUsecaseMock) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	return fh.MockSearchFHIRAllergyIntoleranceFn(ctx, params)
}

// CreateFHIRAllergyIntolerance is a mock implementation of CreateFHIRAllergyIntolerance method
func (fh *FHIRUsecaseMock) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.MockCreateFHIRAllergyIntoleranceFn(ctx, input)
}

// UpdateFHIRAllergyIntolerance is a mock implementation of UpdateFHIRAllergyIntolerance method
func (fh *FHIRUsecaseMock) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.MockUpdateFHIRAllergyIntoleranceFn(ctx, input)
}

// SearchFHIRComposition is a mock implementation of SearchFHIRComposition method
func (fh *FHIRUsecaseMock) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
	return fh.MockSearchFHIRCompositionFn(ctx, params)
}

// CreateFHIRComposition is a mock implementation of CreateFHIRComposition method
func (fh *FHIRUsecaseMock) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return fh.MockCreateFHIRCompositionFn(ctx, input)
}

// UpdateFHIRComposition is a mock implementation of UpdateFHIRComposition method
func (fh *FHIRUsecaseMock) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return fh.MockUpdateFHIRCompositionFn(ctx, input)
}

// DeleteFHIRComposition is a mock implementation of DeleteFHIRComposition method
func (fh *FHIRUsecaseMock) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRCompositionFn(ctx, id)
}

// SearchFHIRCondition is a mock implementation of SearchFHIRCondition method
func (fh *FHIRUsecaseMock) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	return fh.MockSearchFHIRConditionFn(ctx, params)
}

// UpdateFHIRCondition is a mock implementation of UpdateFHIRCondition method
func (fh *FHIRUsecaseMock) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return fh.MockUpdateFHIRConditionFn(ctx, input)
}

// GetFHIREncounter is a mock implementation of GetFHIREncounter method
func (fh *FHIRUsecaseMock) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	return fh.MockGetFHIREncounterFn(ctx, id)
}

// SearchFHIREncounter is a mock implementation of SearchFHIREncounter method
func (fh *FHIRUsecaseMock) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
	return fh.MockSearchFHIREncounterFn(ctx, params)
}

// SearchFHIRMedicationRequest is a mock implementation of SearchFHIRMedicationRequest method
func (fh *FHIRUsecaseMock) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
	return fh.MockSearchFHIRMedicationRequestFn(ctx, params)
}

// CreateFHIRMedicationRequest is a mock implementation of CreateFHIRMedicationRequest method
func (fh *FHIRUsecaseMock) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return fh.MockCreateFHIRMedicationRequestFn(ctx, input)
}

// UpdateFHIRMedicationRequest is a mock implementation of UpdateFHIRMedicationRequest method
func (fh *FHIRUsecaseMock) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return fh.MockUpdateFHIRMedicationRequestFn(ctx, input)
}

// DeleteFHIRMedicationRequest is a mock implementation of DeleteFHIRMedicationRequest method
func (fh *FHIRUsecaseMock) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRMedicationRequestFn(ctx, id)
}

// SearchFHIRObservation is a mock implementation of SearchFHIRObservation method
func (fh *FHIRUsecaseMock) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	return fh.MockSearchFHIRObservationFn(ctx, params)
}

// CreateFHIRObservation is a mock implementation of CreateFHIRObservation method
func (fh *FHIRUsecaseMock) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	return fh.MockCreateFHIRObservationFn(ctx, input)
}

// DeleteFHIRObservation is a mock implementation of DeleteFHIRObservation method
func (fh *FHIRUsecaseMock) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRObservationFn(ctx, id)
}

// GetFHIRPatient is a mock implementation of GetFHIRPatient method
func (fh *FHIRUsecaseMock) GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
	return fh.MockGetFHIRPatientFn(ctx, id)
}

// DeleteFHIRPatient is a mock implementation of DeleteFHIRPatient method
func (fh *FHIRUsecaseMock) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRPatientFn(ctx, id)
}

// DeleteFHIRResourceType is a mock implementation of DeleteFHIRResourceType method
func (fh *FHIRUsecaseMock) DeleteFHIRResourceType(results []map[string]string) error {
	return fh.MockDeleteFHIRResourceTypeFn(results)
}

// DeleteFHIRServiceRequest is a mock implementation of DeleteFHIRServiceRequest method
func (fh *FHIRUsecaseMock) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRServiceRequestFn(ctx, id)
}

// SearchFHIRMedicationStatement is a mock implementation of SearchFHIRMedicationStatement method
func (fh *FHIRUsecaseMock) SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
	return fh.MockSearchFHIRMedicationStatementFn(ctx, params)
}

// FindOrganizationByID is a mock implementation of FindOrganizationByID method
func (fh *FHIRUsecaseMock) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.MockFindOrganizationByIDFn(ctx, organizationID)
}

// CreateFHIRMedicationStatement is a mock implementation of CreateFHIRMedicationStatement method
func (fh *FHIRUsecaseMock) CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
	return fh.MockCreateFHIRMedicationStatementFn(ctx, input)
}

// CreateFHIRMedication is a mock implementation of CreateFHIRMedication method
func (fh *FHIRUsecaseMock) CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
	return fh.MockCreateFHIRMedicationFn(ctx, input)

}

// FHIRHeaders is a mock implementation of CreateFHIRMedication method
func (fh *FHIRUsecaseMock) FHIRHeaders() (http.Header, error) {
	return fh.MockFHIRHeadersFn()
}
