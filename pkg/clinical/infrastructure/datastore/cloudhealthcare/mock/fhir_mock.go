package mock

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIRMock struct implements mocks of FHIR methods.
type FHIRMock struct {
	MockCreateEpisodeOfCareFn    func(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error)
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
	MockCreateFHIRPatientFn             func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error)
	MockPatchFHIRPatientFn              func(ctx context.Context, id string, params []map[string]interface{}) (*domain.FHIRPatient, error)
	MockUpdateFHIREpisodeOfCareFn       func(ctx context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error)
	MockSearchFHIRPatientFn             func(ctx context.Context, searchParams string) (*domain.PatientConnection, error)
}

// NewFHIRMock initializes a new instance of FHIR mock
func NewFHIRMock() *FHIRMock {
	UUID := uuid.New().String()
	PatientRef := "Patient/1"
	OrgRef := "Organization/1"
	status := domain.EpisodeOfCareStatusEnumFinished
	return &FHIRMock{
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
					ID:                 &UUID,
					Text:               &domain.FHIRNarrative{},
					Identifier:         []*domain.FHIRIdentifier{},
					ClinicalStatus:     &domain.FHIRCodeableConcept{},
					VerificationStatus: &domain.FHIRCodeableConcept{},
					Category:           []*domain.FHIRCodeableConcept{},
					Severity:           &domain.FHIRCodeableConcept{},
					Code:               &domain.FHIRCodeableConcept{},
					BodySite:           []*domain.FHIRCodeableConcept{},
					Subject:            &domain.FHIRReference{},
					Encounter:          &domain.FHIRReference{},
					OnsetDateTime:      &scalarutils.Date{},
					OnsetAge:           &domain.FHIRAge{},
					OnsetPeriod:        &domain.FHIRPeriod{},
					OnsetRange:         &domain.FHIRRange{},
					OnsetString:        new(string),
					AbatementDateTime:  &scalarutils.Date{},
					AbatementAge:       &domain.FHIRAge{},
					AbatementPeriod:    &domain.FHIRPeriod{},
					AbatementRange:     &domain.FHIRRange{},
					AbatementString:    new(string),
					RecordedDate:       &scalarutils.Date{},
					Recorder:           &domain.FHIRReference{},
					Asserter:           &domain.FHIRReference{},
					Stage:              []*domain.FHIRConditionStage{},
					Evidence:           []*domain.FHIRConditionEvidence{},
					Note:               []*domain.FHIRAnnotation{},
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
			PatientRef := "Patient/" + uuid.NewString()
			return &domain.FHIREncounterRelayPayload{
				Resource: &domain.FHIREncounter{
					ID:            &UUID,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
					Status:        "",
					StatusHistory: []*domain.FHIREncounterStatushistory{},
					Class:         domain.FHIRCoding{},
					ClassHistory:  []*domain.FHIREncounterClasshistory{},
					Type:          []*domain.FHIRCodeableConcept{},
					ServiceType:   &domain.FHIRCodeableConcept{},
					Priority:      &domain.FHIRCodeableConcept{},
					Subject: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					EpisodeOfCare:   []*domain.FHIRReference{},
					BasedOn:         []*domain.FHIRReference{},
					Participant:     []*domain.FHIREncounterParticipant{},
					Appointment:     []*domain.FHIRReference{},
					Period:          &domain.FHIRPeriod{},
					Length:          &domain.FHIRDuration{},
					ReasonReference: []*domain.FHIRReference{},
					Diagnosis:       []*domain.FHIREncounterDiagnosis{},
					Account:         []*domain.FHIRReference{},
					Hospitalization: &domain.FHIREncounterHospitalization{},
					Location:        []*domain.FHIREncounterLocation{},
					ServiceProvider: &domain.FHIRReference{},
					PartOf:          &domain.FHIRReference{},
				},
			}, nil
		},
		MockSearchFHIREncounterFn: func(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
			PatientRef := "Patient/" + uuid.NewString()
			return &domain.FHIREncounterRelayConnection{
				Edges: []*domain.FHIREncounterRelayEdge{
					{
						Cursor: new(string),
						Node: &domain.FHIREncounter{
							ID:            new(string),
							Text:          &domain.FHIRNarrative{},
							Identifier:    []*domain.FHIRIdentifier{},
							Status:        "",
							StatusHistory: []*domain.FHIREncounterStatushistory{},
							Class:         domain.FHIRCoding{},
							ClassHistory:  []*domain.FHIREncounterClasshistory{},
							Type:          []*domain.FHIRCodeableConcept{},
							ServiceType:   &domain.FHIRCodeableConcept{},
							Priority:      &domain.FHIRCodeableConcept{},
							Subject: &domain.FHIRReference{
								Reference: &PatientRef,
							},
							EpisodeOfCare:   []*domain.FHIRReference{},
							BasedOn:         []*domain.FHIRReference{},
							Participant:     []*domain.FHIREncounterParticipant{},
							Appointment:     []*domain.FHIRReference{},
							Period:          &domain.FHIRPeriod{},
							Length:          &domain.FHIRDuration{},
							ReasonReference: []*domain.FHIRReference{},
							Diagnosis:       []*domain.FHIREncounterDiagnosis{},
							Account:         []*domain.FHIRReference{},
							Hospitalization: &domain.FHIREncounterHospitalization{},
							Location:        []*domain.FHIREncounterLocation{},
							ServiceProvider: &domain.FHIRReference{},
							PartOf:          &domain.FHIRReference{},
						},
					},
				},
				PageInfo: &firebasetools.PageInfo{},
			}, nil
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
		MockCreateFHIRPatientFn: func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID:                   new(string),
					Text:                 &domain.FHIRNarrative{},
					Identifier:           []*domain.FHIRIdentifier{},
					Active:               new(bool),
					Name:                 []*domain.FHIRHumanName{},
					Telecom:              []*domain.FHIRContactPoint{},
					BirthDate:            &scalarutils.Date{},
					DeceasedBoolean:      new(bool),
					DeceasedDateTime:     &scalarutils.Date{},
					Address:              []*domain.FHIRAddress{},
					MaritalStatus:        &domain.FHIRCodeableConcept{},
					MultipleBirthBoolean: new(bool),
					MultipleBirthInteger: new(string),
					Photo:                []*domain.FHIRAttachment{},
					Contact:              []*domain.FHIRPatientContact{},
					Communication:        []*domain.FHIRPatientCommunication{},
					GeneralPractitioner:  []*domain.FHIRReference{},
					ManagingOrganization: &domain.FHIRReference{},
					Link:                 []*domain.FHIRPatientLink{},
				},
				HasOpenEpisodes: false,
				OpenEpisodes:    []*domain.FHIREpisodeOfCare{},
			}, nil
		},
		MockPatchFHIRPatientFn: func(ctx context.Context, id string, params []map[string]interface{}) (*domain.FHIRPatient, error) {
			return &domain.FHIRPatient{
				ID:                   new(string),
				Text:                 &domain.FHIRNarrative{},
				Identifier:           []*domain.FHIRIdentifier{},
				Active:               new(bool),
				Name:                 []*domain.FHIRHumanName{},
				Telecom:              []*domain.FHIRContactPoint{},
				BirthDate:            &scalarutils.Date{},
				DeceasedBoolean:      new(bool),
				DeceasedDateTime:     &scalarutils.Date{},
				Address:              []*domain.FHIRAddress{},
				MaritalStatus:        &domain.FHIRCodeableConcept{},
				MultipleBirthBoolean: new(bool),
				MultipleBirthInteger: new(string),
				Photo:                []*domain.FHIRAttachment{},
				Contact:              []*domain.FHIRPatientContact{},
				Communication:        []*domain.FHIRPatientCommunication{},
				GeneralPractitioner:  []*domain.FHIRReference{},
				ManagingOrganization: &domain.FHIRReference{},
				Link:                 []*domain.FHIRPatientLink{},
			}, nil
		},
		MockUpdateFHIREpisodeOfCareFn: func(ctx context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error) {
			return &domain.FHIREpisodeOfCare{
				ID:                   new(string),
				Text:                 &domain.FHIRNarrative{},
				Identifier:           []*domain.FHIRIdentifier{},
				StatusHistory:        []*domain.FHIREpisodeofcareStatushistory{},
				Type:                 []*domain.FHIRCodeableConcept{},
				Diagnosis:            []*domain.FHIREpisodeofcareDiagnosis{},
				Patient:              &domain.FHIRReference{},
				ManagingOrganization: &domain.FHIRReference{},
				Period:               &domain.FHIRPeriod{},
				ReferralRequest:      []*domain.FHIRReference{},
				CareManager:          &domain.FHIRReference{},
				Team:                 []*domain.FHIRReference{},
				Account:              []*domain.FHIRReference{},
			}, nil
		},
		MockSearchFHIRPatientFn: func(ctx context.Context, searchParams string) (*domain.PatientConnection, error) {
			return &domain.PatientConnection{
				Edges:    []*domain.PatientEdge{},
				PageInfo: &firebasetools.PageInfo{},
			}, nil
		},
	}
}

// GetBearerToken is a mock implementation of get bearer token method
func (fh *FHIRMock) GetBearerToken() (string, error) {
	return fh.MockGetBearerTokenFn()
}

// CreateEpisodeOfCare is a mock implementation of CreateEpisodeOfCare method
func (fh *FHIRMock) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
	return fh.MockCreateEpisodeOfCareFn(ctx, episode)
}

// CreateFHIRCondition is a mock implementation of CreateFHIRCondition method
func (fh *FHIRMock) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return fh.MockCreateFHIRConditionFn(ctx, input)
}

// CreateFHIROrganization is a mock implementation of CreateFHIROrganization method
func (fh *FHIRMock) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.MockCreateFHIROrganizationFn(ctx, input)
}

// SearchFHIROrganization is a mock implementation of SearchFHIROrganization method
func (fh *FHIRMock) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {
	return fh.MockSearchFHIROrganizationFn(ctx, params)
}

// SearchEpisodesByParam is a mock implementation of SearchEpisodesByParam method
func (fh *FHIRMock) SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.MockSearchEpisodesByParamFn(ctx, searchParams)
}

// HasOpenEpisode is a mock implementation of HasOpenEpisode method
func (fh *FHIRMock) HasOpenEpisode(ctx context.Context, patient domain.FHIRPatient) (bool, error) {
	return fh.MockHasOpenEpisodeFn(ctx, patient)
}

// OpenEpisodes is a mock implementation of OpenEpisodes method
func (fh *FHIRMock) OpenEpisodes(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.MockOpenEpisodesFn(ctx, patientReference)
}

// CreateFHIREncounter is a mock implementation of CreateFHIREncounter method
func (fh *FHIRMock) CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
	return fh.MockCreateFHIREncounterFn(ctx, input)
}

// GetFHIREpisodeOfCare is a mock implementation of GetFHIREpisodeOfCare method
func (fh *FHIRMock) GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
	return fh.MockGetFHIREpisodeOfCareFn(ctx, id)
}

// Encounters is a mock implementation of Encounters method
func (fh *FHIRMock) Encounters(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum) ([]*domain.FHIREncounter, error) {
	return fh.MockEncountersFn(ctx, patientReference, status)
}

// SearchFHIREpisodeOfCare is a mock implementation of SearchFHIREpisodeOfCare method
func (fh *FHIRMock) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	return fh.MockSearchFHIREpisodeOfCareFn(ctx, params)
}

// StartEncounter is a mock implementation of StartEncounter method
func (fh *FHIRMock) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	return fh.MockStartEncounterFn(ctx, episodeID)
}

// UpgradeEpisode is a mock implementation of UpgradeEpisode method
func (fh *FHIRMock) UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
	return fh.MockUpgradeEpisodeFn(ctx, input)
}

// SearchEpisodeEncounter is a mock implementation of SearchEpisodeEncounter method
func (fh *FHIRMock) SearchEpisodeEncounter(ctx context.Context, episodeReference string) (*domain.FHIREncounterRelayConnection, error) {
	return fh.MockSearchEpisodeEncounterFn(ctx, episodeReference)
}

// EndEncounter is a mock implementation of EndEncounter method
func (fh *FHIRMock) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	return fh.MockEndEncounterFn(ctx, encounterID)
}

// EndEpisode is a mock implementation of EndEpisode method
func (fh *FHIRMock) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	return fh.MockEndEpisodeFn(ctx, episodeID)
}

// GetActiveEpisode is a mock implementation of GetActiveEpisode method
func (fh *FHIRMock) GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
	return fh.MockGetActiveEpisodeFn(ctx, episodeID)
}

// SearchFHIRServiceRequest is a mock implementation of SearchFHIRServiceRequest method
func (fh *FHIRMock) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
	return fh.MockSearchFHIRServiceRequestFn(ctx, params)
}

// CreateFHIRServiceRequest is a mock implementation of CreateFHIRServiceRequest method
func (fh *FHIRMock) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	return fh.MockCreateFHIRServiceRequestFn(ctx, input)
}

// SearchFHIRAllergyIntolerance is a mock implementation of SearchFHIRAllergyIntolerance method
func (fh *FHIRMock) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	return fh.MockSearchFHIRAllergyIntoleranceFn(ctx, params)
}

// CreateFHIRAllergyIntolerance is a mock implementation of CreateFHIRAllergyIntolerance method
func (fh *FHIRMock) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.MockCreateFHIRAllergyIntoleranceFn(ctx, input)
}

// UpdateFHIRAllergyIntolerance is a mock implementation of UpdateFHIRAllergyIntolerance method
func (fh *FHIRMock) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.MockUpdateFHIRAllergyIntoleranceFn(ctx, input)
}

// SearchFHIRComposition is a mock implementation of SearchFHIRComposition method
func (fh *FHIRMock) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
	return fh.MockSearchFHIRCompositionFn(ctx, params)
}

// CreateFHIRComposition is a mock implementation of CreateFHIRComposition method
func (fh *FHIRMock) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return fh.MockCreateFHIRCompositionFn(ctx, input)
}

// UpdateFHIRComposition is a mock implementation of UpdateFHIRComposition method
func (fh *FHIRMock) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	return fh.MockUpdateFHIRCompositionFn(ctx, input)
}

// DeleteFHIRComposition is a mock implementation of DeleteFHIRComposition method
func (fh *FHIRMock) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRCompositionFn(ctx, id)
}

// SearchFHIRCondition is a mock implementation of SearchFHIRCondition method
func (fh *FHIRMock) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	return fh.MockSearchFHIRConditionFn(ctx, params)
}

// UpdateFHIRCondition is a mock implementation of UpdateFHIRCondition method
func (fh *FHIRMock) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return fh.MockUpdateFHIRConditionFn(ctx, input)
}

// GetFHIREncounter is a mock implementation of GetFHIREncounter method
func (fh *FHIRMock) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	return fh.MockGetFHIREncounterFn(ctx, id)
}

// SearchFHIREncounter is a mock implementation of SearchFHIREncounter method
func (fh *FHIRMock) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
	return fh.MockSearchFHIREncounterFn(ctx, params)
}

// SearchFHIRMedicationRequest is a mock implementation of SearchFHIRMedicationRequest method
func (fh *FHIRMock) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
	return fh.MockSearchFHIRMedicationRequestFn(ctx, params)
}

// CreateFHIRMedicationRequest is a mock implementation of CreateFHIRMedicationRequest method
func (fh *FHIRMock) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return fh.MockCreateFHIRMedicationRequestFn(ctx, input)
}

// UpdateFHIRMedicationRequest is a mock implementation of UpdateFHIRMedicationRequest method
func (fh *FHIRMock) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	return fh.MockUpdateFHIRMedicationRequestFn(ctx, input)
}

// DeleteFHIRMedicationRequest is a mock implementation of DeleteFHIRMedicationRequest method
func (fh *FHIRMock) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRMedicationRequestFn(ctx, id)
}

// SearchFHIRObservation is a mock implementation of SearchFHIRObservation method
func (fh *FHIRMock) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	return fh.MockSearchFHIRObservationFn(ctx, params)
}

// CreateFHIRObservation is a mock implementation of CreateFHIRObservation method
func (fh *FHIRMock) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	return fh.MockCreateFHIRObservationFn(ctx, input)
}

// DeleteFHIRObservation is a mock implementation of DeleteFHIRObservation method
func (fh *FHIRMock) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRObservationFn(ctx, id)
}

// GetFHIRPatient is a mock implementation of GetFHIRPatient method
func (fh *FHIRMock) GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
	return fh.MockGetFHIRPatientFn(ctx, id)
}

// DeleteFHIRPatient is a mock implementation of DeleteFHIRPatient method
func (fh *FHIRMock) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRPatientFn(ctx, id)
}

// DeleteFHIRResourceType is a mock implementation of DeleteFHIRResourceType method
func (fh *FHIRMock) DeleteFHIRResourceType(results []map[string]string) error {
	return fh.MockDeleteFHIRResourceTypeFn(results)
}

// DeleteFHIRServiceRequest is a mock implementation of DeleteFHIRServiceRequest method
func (fh *FHIRMock) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	return fh.MockDeleteFHIRServiceRequestFn(ctx, id)
}

// SearchFHIRMedicationStatement is a mock implementation of SearchFHIRMedicationStatement method
func (fh *FHIRMock) SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
	return fh.MockSearchFHIRMedicationStatementFn(ctx, params)
}

// FindOrganizationByID is a mock implementation of FindOrganizationByID method
func (fh *FHIRMock) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.MockFindOrganizationByIDFn(ctx, organizationID)
}

// CreateFHIRMedicationStatement is a mock implementation of CreateFHIRMedicationStatement method
func (fh *FHIRMock) CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
	return fh.MockCreateFHIRMedicationStatementFn(ctx, input)
}

// CreateFHIRMedication is a mock implementation of CreateFHIRMedication method
func (fh *FHIRMock) CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
	return fh.MockCreateFHIRMedicationFn(ctx, input)

}

// FHIRHeaders is a mock implementation of CreateFHIRMedication method
func (fh *FHIRMock) FHIRHeaders() (http.Header, error) {
	return fh.MockFHIRHeadersFn()
}

// CreateFHIRPatient mocks the implementation of creating a FHIR patient
func (fh *FHIRMock) CreateFHIRPatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	return fh.MockCreateFHIRPatientFn(ctx, input)
}

// PatchFHIRPatient mocks the implementation for patching a fhir patient
func (fh *FHIRMock) PatchFHIRPatient(ctx context.Context, id string, params []map[string]interface{}) (*domain.FHIRPatient, error) {
	return fh.MockPatchFHIRPatientFn(ctx, id, params)
}

// UpdateFHIREpisodeOfCare mocks the implementation of updating a FHIR episode of care
func (fh *FHIRMock) UpdateFHIREpisodeOfCare(ctx context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error) {
	return fh.MockUpdateFHIREpisodeOfCareFn(ctx, fhirResourceID, payload)
}

// SearchFHIRPatient mocks the implementation of searching a FHIR patient
func (fh *FHIRMock) SearchFHIRPatient(ctx context.Context, searchParams string) (*domain.PatientConnection, error) {
	return fh.MockSearchFHIRPatientFn(ctx, searchParams)
}
