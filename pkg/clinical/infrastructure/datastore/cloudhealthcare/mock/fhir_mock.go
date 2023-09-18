package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIRMock struct implements mocks of FHIR methods.
type FHIRMock struct {
	MockCreateEpisodeOfCareFn    func(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error)
	MockSearchFHIRConditionFn    func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error)
	MockCreateFHIRConditionFn    func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	MockCreateFHIROrganizationFn func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	MockSearchFHIROrganizationFn func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIROrganizationRelayConnection, error)
	MockGetFHIROrganizationFn    func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error)
	MockSearchEpisodesByParamFn  func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error)
	MockHasOpenEpisodeFn         func(
		ctx context.Context,
		patient domain.FHIRPatient,
		tenant dto.TenantIdentifiers,
		pagination dto.Pagination,
	) (bool, error)
	MockOpenEpisodesFn func(
		ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error)
	MockCreateFHIREncounterFn             func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error)
	MockGetFHIREpisodeOfCareFn            func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error)
	MockSearchPatientEncountersFn         func(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error)
	MockSearchFHIREpisodeOfCareFn         func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error)
	MockStartEncounterFn                  func(ctx context.Context, episodeID string) (string, error)
	MockUpgradeEpisodeFn                  func(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error)
	MockSearchEpisodeEncounterFn          func(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error)
	MockEndEncounterFn                    func(ctx context.Context, encounterID string) (bool, error)
	MockEndEpisodeFn                      func(ctx context.Context, episodeID string) (bool, error)
	MockGetActiveEpisodeFn                func(ctx context.Context, episodeID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCare, error)
	MockSearchFHIRServiceRequestFn        func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRServiceRequestRelayConnection, error)
	MockCreateFHIRServiceRequestFn        func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error)
	MockSearchFHIRAllergyIntoleranceFn    func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error)
	MockCreateFHIRAllergyIntoleranceFn    func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	MockUpdateFHIRAllergyIntoleranceFn    func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	MockSearchFHIRCompositionFn           func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRCompositionRelayConnection, error)
	MockCreateFHIRCompositionFn           func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	MockUpdateFHIRCompositionFn           func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	MockDeleteFHIRCompositionFn           func(ctx context.Context, id string) (bool, error)
	MockUpdateFHIRConditionFn             func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	MockGetFHIREncounterFn                func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error)
	MockPatchFHIREncounterFn              func(ctx context.Context, encounterID string, input domain.FHIREncounterInput) (*domain.FHIREncounter, error)
	MockSearchFHIREncounterFn             func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error)
	MockSearchFHIRMedicationRequestFn     func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationRequestRelayConnection, error)
	MockCreateFHIRMedicationRequestFn     func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	MockUpdateFHIRMedicationRequestFn     func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	MockDeleteFHIRMedicationRequestFn     func(ctx context.Context, id string) (bool, error)
	MockSearchFHIRObservationFn           func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error)
	MockCreateFHIRObservationFn           func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error)
	MockDeleteFHIRObservationFn           func(ctx context.Context, id string) (bool, error)
	MockGetFHIRPatientFn                  func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error)
	MockDeleteFHIRPatientFn               func(ctx context.Context, id string) (bool, error)
	MockDeleteFHIRServiceRequestFn        func(ctx context.Context, id string) (bool, error)
	MockDeleteFHIRResourceTypeFn          func(results []map[string]string) error
	MockCreateFHIRMedicationStatementFn   func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)
	MockCreateFHIRMedicationFn            func(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error)
	MockSearchFHIRMedicationStatementFn   func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error)
	MockCreateFHIRPatientFn               func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error)
	MockPatchFHIRPatientFn                func(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error)
	MockPatchFHIREpisodeOfCareFn          func(ctx context.Context, id string, input domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error)
	MockUpdateFHIREpisodeOfCareFn         func(ctx context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error)
	MockSearchFHIRPatientFn               func(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error)
	MockSearchPatientObservationsFn       func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error)
	MockGetFHIRAllergyIntoleranceFn       func(ctx context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	MockSearchPatientAllergyIntoleranceFn func(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error)
	MockCreateFHIRMediaFn                 func(ctx context.Context, input domain.FHIRMedia) (*domain.FHIRMedia, error)
	MockSearchPatientMediaFn              func(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedia, error)
}

// NewFHIRMock initializes a new instance of FHIR mock
func NewFHIRMock() *FHIRMock {
	return &FHIRMock{
		MockCreateEpisodeOfCareFn: func(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error) {
			UUID := uuid.New().String()
			PatientRef := "Patient/1"
			OrgRef := "Organization/1"
			status := domain.EpisodeOfCareStatusEnumActive
			return &domain.EpisodeOfCarePayload{
				EpisodeOfCare: &domain.FHIREpisodeOfCare{
					ID:            &UUID,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
					Status:        &status,
					StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
					Type:          []*domain.FHIRCodeableConcept{},
					Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
					Patient: &domain.FHIRReference{
						ID:        &UUID,
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
		MockCreateFHIRConditionFn: func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
			UUID := uuid.New().String()
			statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
			status := "active"
			note := scalarutils.Markdown("Fever Fever")
			noteTime := time.Now()
			uri := scalarutils.URI("1234567")

			return &domain.FHIRConditionRelayPayload{
				Resource: &domain.FHIRCondition{
					ID:         &UUID,
					Text:       &domain.FHIRNarrative{},
					Identifier: []*domain.FHIRIdentifier{},
					ClinicalStatus: &domain.FHIRCodeableConcept{
						Coding: []*domain.FHIRCoding{
							{
								System:  &statusSystem,
								Code:    scalarutils.Code(string(status)),
								Display: string(status),
							},
						},
						Text: string(status),
					},
					Code: &domain.FHIRCodeableConcept{
						Coding: []*domain.FHIRCoding{
							{
								System:  &uri,
								Code:    scalarutils.Code("1234"),
								Display: "1234",
							},
						},
						Text: "1234",
					},
					OnsetDateTime: &scalarutils.Date{},
					RecordedDate:  &scalarutils.Date{},
					Note: []*domain.FHIRAnnotation{
						{
							Time: &noteTime,
							Text: &note,
						},
					},
					Subject: &domain.FHIRReference{
						ID: &UUID,
					},
					Encounter: &domain.FHIRReference{
						ID: &UUID,
					},
					Category: []*domain.FHIRCodeableConcept{
						{
							ID: &UUID,
							Coding: []*domain.FHIRCoding{
								{
									ID:           &UUID,
									System:       (*scalarutils.URI)(&UUID),
									Version:      &UUID,
									Code:         "PROBLEM_LIST_ITEM",
									Display:      gofakeit.BeerAlcohol(),
									UserSelected: new(bool),
								},
							},
							Text: "PROBLEM_LIST_ITEM",
						},
					},
				},
			}, nil
		},
		MockCreateFHIROrganizationFn: func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
			UUID := uuid.New().String()
			active := true
			name := gofakeit.Name()
			uri := ""

			use := domain.ContactPointUseEnumWork
			rank := int64(1)
			phoneSystem := domain.ContactPointSystemEnumPhone
			phoneNumber := gofakeit.Phone()

			return &domain.FHIROrganizationRelayPayload{
				Resource: &domain.FHIROrganization{
					ID:     &UUID,
					Name:   &name,
					Active: &active,
					Identifier: []*domain.FHIRIdentifier{
						{
							Use: "official",
							Type: domain.FHIRCodeableConcept{
								Text: "type",
							},
							System:   (*scalarutils.URI)(&uri),
							Value:    "",
							Period:   &domain.FHIRPeriod{},
							Assigner: &domain.FHIRReference{},
						},
					},
					Telecom: []*domain.FHIRContactPoint{
						{
							System: &phoneSystem,
							Value:  &phoneNumber,
							Use:    &use,
							Rank:   &rank,
							Period: &domain.FHIRPeriod{},
						},
					},
				},
			}, nil
		},
		MockSearchFHIROrganizationFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIROrganizationRelayConnection, error) {
			return &domain.FHIROrganizationRelayConnection{}, nil
		},
		MockSearchEpisodesByParamFn: func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error) {
			return []*domain.FHIREpisodeOfCare{}, nil
		},
		MockHasOpenEpisodeFn: func(ctx context.Context, patient domain.FHIRPatient, tenant dto.TenantIdentifiers, pagination dto.Pagination) (bool, error) {
			return true, nil
		},
		MockSearchPatientMediaFn: func(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedia, error) {
			UUID := uuid.New().String()
			PatientRef := "Patient/1"
			OrgRef := "Organization/1"
			contentType := "video/mp4"
			url := gofakeit.URL()
			title := "test"
			return &domain.PagedFHIRMedia{
				Media: []domain.FHIRMedia{
					{
						ID: &UUID,
						Subject: &domain.FHIRReferenceInput{
							ID:        &UUID,
							Reference: &PatientRef,
						},
						Operator: &domain.FHIRReferenceInput{
							ID:        &UUID,
							Reference: &OrgRef,
						},
						Content: &domain.FHIRAttachmentInput{
							ID:          &UUID,
							ContentType: (*scalarutils.Code)(&contentType),
							URL:         (*scalarutils.URL)(&url),
							Title:       &title,
						},
					},
				},
				HasNextPage:     true,
				NextCursor:      "",
				HasPreviousPage: true,
				PreviousCursor:  "",
				TotalCount:      10,
			}, nil
		},
		MockOpenEpisodesFn: func(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error) {
			UUID := uuid.New().String()
			PatientRef := "Patient/1"
			OrgRef := "Organization/1"
			return []*domain.FHIREpisodeOfCare{
				{
					ID:            &UUID,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
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
			resourceID := uuid.New().String()
			return &domain.FHIREncounterRelayPayload{
				Resource: &domain.FHIREncounter{
					ID: &resourceID,
				},
			}, nil
		},
		MockGetFHIREpisodeOfCareFn: func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
			UUID := uuid.New().String()
			PatientRef := "Patient/1"
			OrgRef := "Organization/1"
			st := domain.EpisodeOfCareStatusEnumActive
			return &domain.FHIREpisodeOfCareRelayPayload{
				Resource: &domain.FHIREpisodeOfCare{
					ID:            &UUID,
					Status:        &st,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
					StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
					Type:          []*domain.FHIRCodeableConcept{},
					Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
					Patient: &domain.FHIRReference{
						ID:        &UUID,
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
		MockSearchPatientEncountersFn: func(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
			encounterID := uuid.New().String()
			patientID := uuid.New().String()
			episodeID := uuid.New().String()
			return &domain.PagedFHIREncounter{
				Encounters: []domain.FHIREncounter{
					{
						ID:     &encounterID,
						Status: "finished",
						Class: domain.FHIRCoding{
							Display: "ambulatory",
						},
						Subject: &domain.FHIRReference{
							ID: &patientID,
						},
						EpisodeOfCare: []*domain.FHIRReference{{
							ID: &episodeID,
						}},
					},
				},
				HasNextPage:     false,
				NextCursor:      "",
				HasPreviousPage: false,
				PreviousCursor:  "",
				TotalCount:      0,
			}, nil
		},
		MockSearchFHIREpisodeOfCareFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
			PatientRef := "Patient/1"
			OrgRef := "Organization/1"
			return &domain.FHIREpisodeOfCareRelayConnection{
				Edges: []*domain.FHIREpisodeOfCareRelayEdge{
					{
						Cursor: new(string),
						Node: &domain.FHIREpisodeOfCare{
							ID:            new(string),
							Text:          &domain.FHIRNarrative{},
							Identifier:    []*domain.FHIRIdentifier{},
							StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
							Type:          []*domain.FHIRCodeableConcept{},
							Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
							Patient:       &domain.FHIRReference{Reference: &PatientRef},
							ManagingOrganization: &domain.FHIRReference{
								Reference: &OrgRef,
							},
							Period:          &domain.FHIRPeriod{},
							ReferralRequest: []*domain.FHIRReference{},
							CareManager:     &domain.FHIRReference{},
							Team:            []*domain.FHIRReference{},
							Account:         []*domain.FHIRReference{},
							Meta:            &domain.FHIRMeta{},
							Extension:       []*domain.FHIRExtension{},
						},
					},
				},
				PageInfo: &firebasetools.PageInfo{},
			}, nil
		},
		MockStartEncounterFn: func(ctx context.Context, episodeID string) (string, error) {
			return "test-encounter", nil
		},
		MockUpgradeEpisodeFn: func(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
			return &domain.EpisodeOfCarePayload{}, nil
		},
		MockSearchEpisodeEncounterFn: func(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
			id := gofakeit.UUID()
			return &domain.PagedFHIREncounter{
				Encounters: []domain.FHIREncounter{
					{
						ID: &id,
						Text: &domain.FHIRNarrative{
							ID: &id,
						},
						Identifier: []*domain.FHIRIdentifier{
							{
								ID: &id,
							},
						},
						Status:        "",
						StatusHistory: []*domain.FHIREncounterStatushistory{},
						Class:         domain.FHIRCoding{},
						ClassHistory:  []*domain.FHIREncounterClasshistory{},
						Type:          []*domain.FHIRCodeableConcept{},
						ServiceType:   &domain.FHIRCodeableConcept{},
						Priority:      &domain.FHIRCodeableConcept{},
						Subject: &domain.FHIRReference{
							ID: &id,
						},
						EpisodeOfCare: []*domain.FHIRReference{
							{
								ID: &id,
							},
						},
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
						Meta:            &domain.FHIRMeta{},
						Extension:       []*domain.FHIRExtension{},
					},
					{
						ID: &id,
					},
				},
				HasNextPage:     true,
				NextCursor:      "",
				HasPreviousPage: false,
				PreviousCursor:  "",
				TotalCount:      0,
			}, nil
		},
		MockEndEncounterFn: func(ctx context.Context, encounterID string) (bool, error) {
			return true, nil
		},
		MockEndEpisodeFn: func(ctx context.Context, episodeID string) (bool, error) {
			return true, nil
		},
		MockGetActiveEpisodeFn: func(ctx context.Context, episodeID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCare, error) {
			return &domain.FHIREpisodeOfCare{}, nil
		},
		MockSearchFHIRServiceRequestFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRServiceRequestRelayConnection, error) {
			return &domain.FHIRServiceRequestRelayConnection{}, nil
		},
		MockCreateFHIRServiceRequestFn: func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
			return &domain.FHIRServiceRequestRelayPayload{}, nil
		},
		MockSearchFHIRAllergyIntoleranceFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
			UID := gofakeit.UUID()
			system := scalarutils.URI("/orgs/CIEL/sources/CIEL/concepts/148888/")
			severityStatus := domain.AllergyIntoleranceReactionSeverityEnumSevere
			return &domain.PagedFHIRAllergy{
				Allergies: []domain.FHIRAllergyIntolerance{
					{
						ID: &UID,
						Encounter: &domain.FHIRReference{
							ID: &UID,
						},
						Code: &domain.FHIRCodeableConcept{
							ID: &UID,
							Coding: []*domain.FHIRCoding{
								{
									ID:     &UID,
									Code:   scalarutils.Code("124"),
									System: &system,
								},
							},
						},
						Patient: &domain.FHIRReference{
							ID: &UID,
						},
						Reaction: []*domain.FHIRAllergyintoleranceReaction{
							{
								ID:        &UID,
								Substance: &domain.FHIRCodeableConcept{},
								Manifestation: []*domain.FHIRCodeableConcept{
									{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												ID:     new(string),
												System: &system,
												Code:   scalarutils.Code("1234"),
											},
										},
										Text: gofakeit.Name(),
									},
								},
								Severity: &severityStatus,
							},
						},
						Meta:      &domain.FHIRMeta{},
						Extension: []*domain.FHIRExtension{},
					},
				},
				HasNextPage:     false,
				NextCursor:      "",
				HasPreviousPage: false,
				PreviousCursor:  "",
				TotalCount:      0,
			}, nil
		},
		MockSearchPatientAllergyIntoleranceFn: func(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
			UID := gofakeit.UUID()
			system := scalarutils.URI("/orgs/CIEL/sources/CIEL/concepts/148888/")
			severityStatus := domain.AllergyIntoleranceReactionSeverityEnumSevere
			return &domain.PagedFHIRAllergy{
				Allergies: []domain.FHIRAllergyIntolerance{
					{
						ID: &UID,
						Encounter: &domain.FHIRReference{
							ID: &UID,
						},
						Code: &domain.FHIRCodeableConcept{
							ID: &UID,
							Coding: []*domain.FHIRCoding{
								{
									ID:     &UID,
									Code:   scalarutils.Code("124"),
									System: &system,
								},
							},
						},
						Patient: &domain.FHIRReference{
							ID: &UID,
						},
						Reaction: []*domain.FHIRAllergyintoleranceReaction{
							{
								ID:        &UID,
								Substance: &domain.FHIRCodeableConcept{},
								Manifestation: []*domain.FHIRCodeableConcept{
									{
										ID: new(string),
										Coding: []*domain.FHIRCoding{
											{
												ID:     new(string),
												System: &system,
												Code:   scalarutils.Code("1234"),
											},
										},
										Text: gofakeit.Name(),
									},
								},
								Severity: &severityStatus,
							},
						},
						Meta:      &domain.FHIRMeta{},
						Extension: []*domain.FHIRExtension{},
					},
				},
				HasNextPage:     false,
				NextCursor:      "",
				HasPreviousPage: false,
				PreviousCursor:  "",
				TotalCount:      0,
			}, nil
		},
		MockGetFHIRAllergyIntoleranceFn: func(ctx context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
			UID := uuid.NewString()
			system := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
			return &domain.FHIRAllergyIntoleranceRelayPayload{
				Resource: &domain.FHIRAllergyIntolerance{
					ID:         &UID,
					Text:       &domain.FHIRNarrative{},
					Identifier: []*domain.FHIRIdentifier{},
					ClinicalStatus: domain.FHIRCodeableConcept{
						ID:     &UID,
						Coding: []*domain.FHIRCoding{},
						Text:   "test",
					},
					VerificationStatus: domain.FHIRCodeableConcept{},
					Category:           []*domain.AllergyIntoleranceCategoryEnum{},
					Criticality:        "",
					Code: &domain.FHIRCodeableConcept{
						ID: &UID,
						Coding: []*domain.FHIRCoding{
							{
								ID:     &UID,
								System: &system,
								Code:   scalarutils.Code("1234"),
							},
						},
						Text: "",
					},
					Patient: &domain.FHIRReference{
						ID: &UID,
					},
					Encounter: &domain.FHIRReference{
						ID: &UID,
					},
					OnsetDateTime: &scalarutils.Date{},
					OnsetAge:      &domain.FHIRAge{},
					OnsetPeriod:   &domain.FHIRPeriod{},
					OnsetRange:    &domain.FHIRRange{},
					OnsetString:   new(string),
					RecordedDate:  &scalarutils.Date{},
					Recorder:      &domain.FHIRReference{},
					Asserter:      &domain.FHIRReference{},
					Note:          []*domain.FHIRAnnotation{},
					Reaction:      []*domain.FHIRAllergyintoleranceReaction{},
					Meta:          &domain.FHIRMeta{},
					Extension:     []*domain.FHIRExtension{},
				},
			}, nil
		},
		MockCreateFHIRAllergyIntoleranceFn: func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
			system := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
			return &domain.FHIRAllergyIntoleranceRelayPayload{
				Resource: &domain.FHIRAllergyIntolerance{
					ID:   new(string),
					Text: &domain.FHIRNarrative{},
					Reaction: []*domain.FHIRAllergyintoleranceReaction{
						{
							ID:        new(string),
							Substance: &domain.FHIRCodeableConcept{},
							Manifestation: []*domain.FHIRCodeableConcept{
								{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											ID:     new(string),
											System: &system,
											Code:   scalarutils.Code("1234"),
										},
									},
									Text: gofakeit.Name(),
								},
							},
						},
					},
					// RecordedDate: &scalarutils.Date{},
					Meta:      &domain.FHIRMeta{},
					Extension: []*domain.FHIRExtension{},
				},
			}, nil
		},
		MockUpdateFHIRAllergyIntoleranceFn: func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
			return &domain.FHIRAllergyIntoleranceRelayPayload{}, nil
		},
		MockSearchFHIRCompositionFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRCompositionRelayConnection, error) {
			return &domain.FHIRCompositionRelayConnection{}, nil
		},
		MockCreateFHIRCompositionFn: func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
			UUID := uuid.New().String()
			compositionTitle := gofakeit.Name() + "assessment note"
			typeSystem := scalarutils.URI("http://hl7.org/fhir/ValueSet/doc-typecodes")
			categorySystem := scalarutils.URI("http://hl7.org/fhir/ValueSet/referenced-item-category")
			category := "Assessment + plan"
			compositionType := "Progress note"
			treatmentPlan := "Treatment Plan"
			compositionStatus := "active"
			note := scalarutils.Markdown("Fever Fever")
			PatientRef := "Patient/" + uuid.NewString()
			patientType := "Patient"
			organizationRef := "Organization/" + uuid.NewString()
			compositionSectionTextStatus := "generated"

			return &domain.FHIRCompositionRelayPayload{
				Resource: &domain.FHIRComposition{
					ID:         &UUID,
					Text:       &domain.FHIRNarrative{},
					Identifier: &domain.FHIRIdentifier{},
					Status:     (*domain.CompositionStatusEnum)(&compositionStatus),
					Type: &domain.FHIRCodeableConcept{
						ID: new(string),
						Coding: []*domain.FHIRCoding{
							{
								ID:      &UUID,
								System:  &typeSystem,
								Code:    scalarutils.Code(string(common.LOINCProgressNoteCode)),
								Display: compositionType,
							},
						},
						Text: "Progress note",
					},
					Category: []*domain.FHIRCodeableConcept{
						{
							// ID: new(string),
							Coding: []*domain.FHIRCoding{
								{
									ID:      &UUID,
									System:  &categorySystem,
									Version: new(string),
									Code:    scalarutils.Code(string(common.LOINCAssessmentPlanCode)),
									Display: category,
								},
							},
							Text: "Assessment + plan",
						},
					},
					Subject: &domain.FHIRReference{
						ID:        &UUID,
						Reference: &PatientRef,
						Type:      (*scalarutils.URI)(&patientType),
					},
					Encounter: &domain.FHIRReference{
						ID: &UUID,
					},
					Date: &scalarutils.Date{
						Year:  2023,
						Month: 9,
						Day:   25,
					},
					Author: []*domain.FHIRReference{
						{
							Reference: &organizationRef,
						},
					},
					Title: &compositionTitle,
					Section: []*domain.FHIRCompositionSection{
						{
							ID:    &UUID,
							Title: &treatmentPlan,
							Code: &domain.FHIRCodeableConceptInput{
								ID: new(string),
								Coding: []*domain.FHIRCodingInput{
									{
										ID:      &UUID,
										System:  &categorySystem,
										Version: new(string),
										Code:    scalarutils.Code(string(common.LOINCAssessmentPlanCode)),
										Display: category,
									},
								},
								Text: "Assessment + plan",
							},
							Author: []*domain.FHIRReference{
								{
									Reference: new(string),
								},
							},
							Text: &domain.FHIRNarrative{
								ID:     &UUID,
								Status: (*domain.NarrativeStatusEnum)(&compositionSectionTextStatus),
								Div:    scalarutils.XHTML(note),
							},
						},
					},
				},
			}, nil
		},
		MockUpdateFHIRCompositionFn: func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
			return &domain.FHIRCompositionRelayPayload{}, nil
		},
		MockDeleteFHIRCompositionFn: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},
		MockSearchFHIRConditionFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
			id := gofakeit.UUID()
			statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
			status := "inactive"
			note := scalarutils.Markdown("Fever Fever")
			noteTime := time.Now()
			uri := scalarutils.URI("1234567345")

			condition := domain.FHIRCondition{
				ID:         &id,
				Text:       &domain.FHIRNarrative{},
				Identifier: []*domain.FHIRIdentifier{},
				ClinicalStatus: &domain.FHIRCodeableConcept{
					Coding: []*domain.FHIRCoding{
						{
							System:  &statusSystem,
							Code:    scalarutils.Code(string(status)),
							Display: string(status),
						},
					},
					Text: string(status),
				},
				Code: &domain.FHIRCodeableConcept{
					Coding: []*domain.FHIRCoding{
						{
							System:  &uri,
							Code:    scalarutils.Code("1234"),
							Display: "1234567",
						},
					},
					Text: "1234",
				},
				OnsetDateTime: &scalarutils.Date{},
				RecordedDate:  &scalarutils.Date{},
				Subject: &domain.FHIRReference{
					ID: &id,
				},
				Note: []*domain.FHIRAnnotation{
					{
						Time: &noteTime,
						Text: &note,
					},
				},
				Encounter: &domain.FHIRReference{
					ID: &id,
				},
				Category: []*domain.FHIRCodeableConcept{
					{
						ID: &id,
						Coding: []*domain.FHIRCoding{
							{
								ID:           &id,
								System:       (*scalarutils.URI)(&id),
								Version:      &id,
								Code:         "PROBLEM_LIST_ITEM",
								Display:      gofakeit.BeerAlcohol(),
								UserSelected: new(bool),
							},
						},
						Text: "PROBLEM_LIST_ITEM",
					},
				},
			}

			return &domain.PagedFHIRCondition{
				Conditions:      []domain.FHIRCondition{condition},
				HasNextPage:     false,
				HasPreviousPage: false,
			}, nil
		},
		MockUpdateFHIRConditionFn: func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
			return &domain.FHIRConditionRelayPayload{}, nil
		},
		MockGetFHIREncounterFn: func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
			UUID := uuid.New().String()
			PatientRef := "Patient/" + uuid.NewString()
			return &domain.FHIREncounterRelayPayload{
				Resource: &domain.FHIREncounter{
					ID:            &UUID,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
					Status:        domain.EncounterStatusEnum(domain.EncounterStatusEnumOnleave),
					StatusHistory: []*domain.FHIREncounterStatushistory{},
					Class:         domain.FHIRCoding{},
					ClassHistory:  []*domain.FHIREncounterClasshistory{},
					Type:          []*domain.FHIRCodeableConcept{},
					ServiceType:   &domain.FHIRCodeableConcept{},
					Priority:      &domain.FHIRCodeableConcept{},
					Subject: &domain.FHIRReference{
						ID:        &UUID,
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
		MockPatchFHIREncounterFn: func(ctx context.Context, encounterID string, input domain.FHIREncounterInput) (*domain.FHIREncounter, error) {
			UUID := uuid.New().String()
			PatientRef := "Patient/" + uuid.NewString()
			return &domain.FHIREncounter{
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
					ID:        &UUID,
					Reference: &PatientRef,
				},
				EpisodeOfCare: []*domain.FHIRReference{
					{
						ID: &UUID,
					},
				},
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
			}, nil
		},
		MockSearchFHIREncounterFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
			PatientRef := "Patient/" + uuid.NewString()
			UID := gofakeit.UUID()
			return &domain.PagedFHIREncounter{
				Encounters: []domain.FHIREncounter{
					{
						ID: &UID,
						Text: &domain.FHIRNarrative{
							ID: &UID,
						},
						Identifier: []*domain.FHIRIdentifier{
							{
								ID: &UID,
							},
						},
						Status:        "",
						StatusHistory: []*domain.FHIREncounterStatushistory{},
						Class:         domain.FHIRCoding{},
						ClassHistory:  []*domain.FHIREncounterClasshistory{},
						Type:          []*domain.FHIRCodeableConcept{},
						ServiceType:   &domain.FHIRCodeableConcept{},
						Priority:      &domain.FHIRCodeableConcept{},
						Subject: &domain.FHIRReference{
							ID:        &UID,
							Reference: &PatientRef,
						},
						EpisodeOfCare: []*domain.FHIRReference{
							{
								ID: &UID,
							},
						},
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
						Meta:            &domain.FHIRMeta{},
						Extension:       []*domain.FHIRExtension{},
					},
				},
				HasNextPage:     false,
				NextCursor:      "",
				HasPreviousPage: false,
				PreviousCursor:  "",
				TotalCount:      0,
			}, nil
		},
		MockSearchFHIRMedicationRequestFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationRequestRelayConnection, error) {
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
		MockSearchFHIRObservationFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
			uuid := uuid.New().String()
			finalStatus := domain.ObservationStatusEnumFinal
			return &domain.PagedFHIRObservations{
				Observations: []domain.FHIRObservation{
					{
						ID:     &uuid,
						Status: &finalStatus,
						Subject: &domain.FHIRReference{
							ID: &uuid,
						},
						Encounter: &domain.FHIRReference{
							ID: &uuid,
						},
					},
				},
				HasNextPage:     false,
				NextCursor:      "",
				HasPreviousPage: false,
				PreviousCursor:  "",
				TotalCount:      0,
			}, nil
		},
		MockCreateFHIRObservationFn: func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
			uuid := uuid.New().String()
			instant := gofakeit.TimeZone()
			finalStatus := domain.ObservationStatusEnumFinal
			return &domain.FHIRObservation{
				ID:         new(string),
				Text:       &domain.FHIRNarrative{},
				Identifier: []*domain.FHIRIdentifier{},
				BasedOn:    []*domain.FHIRReference{},
				PartOf:     []*domain.FHIRReference{},
				Category:   []*domain.FHIRCodeableConcept{},
				Code: domain.FHIRCodeableConcept{
					ID: new(string),
					Coding: []*domain.FHIRCoding{
						{
							ID:           new(string),
							Version:      new(string),
							Code:         "",
							Display:      "Vital",
							UserSelected: new(bool),
						},
					},
					Text: "",
				},
				Subject: &domain.FHIRReference{
					ID: &uuid,
				},
				Status: &finalStatus,
				Focus:  []*domain.FHIRReference{},
				Encounter: &domain.FHIRReference{
					ID: &uuid,
				},
				EffectiveDateTime:    &scalarutils.Date{},
				EffectivePeriod:      &domain.FHIRPeriod{},
				EffectiveTiming:      &domain.FHIRTiming{},
				EffectiveInstant:     (*scalarutils.Instant)(&instant),
				Performer:            []*domain.FHIRReference{},
				ValueQuantity:        &domain.FHIRQuantity{},
				ValueCodeableConcept: (*scalarutils.Code)(&uuid),
				ValueString:          new(string),
				ValueBoolean:         new(bool),
				ValueInteger:         new(string),
				ValueRange:           &domain.FHIRRange{},
				ValueRatio:           &domain.FHIRRatio{},
				ValueSampledData: &domain.FHIRSampledData{
					ID: &uuid,
				},
				ValueTime:        &time.Time{},
				ValueDateTime:    &scalarutils.Date{},
				ValuePeriod:      &domain.FHIRPeriod{},
				DataAbsentReason: &domain.FHIRCodeableConcept{},
				Interpretation:   []*domain.FHIRCodeableConcept{},
				Note:             []*domain.FHIRAnnotation{},
				BodySite:         &domain.FHIRCodeableConcept{},
				Method:           &domain.FHIRCodeableConcept{},
				Specimen:         &domain.FHIRReference{},
				Device:           &domain.FHIRReference{},
				ReferenceRange:   []*domain.FHIRObservationReferencerange{},
				HasMember:        []*domain.FHIRReference{},
				DerivedFrom:      []*domain.FHIRReference{},
				Component:        []*domain.FHIRObservationComponent{},
				Meta:             &domain.FHIRMeta{},
				Extension:        []*domain.FHIRExtension{},
			}, nil
		},
		MockDeleteFHIRObservationFn: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},
		MockGetFHIRPatientFn: func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
			patientID := uuid.New().String()
			patientName := gofakeit.Name()
			return &domain.FHIRPatientRelayPayload{
				Resource: &domain.FHIRPatient{
					ID: &patientID,
					Name: []*domain.FHIRHumanName{
						{
							Given: []*string{&patientName},
						},
					},
				},
			}, nil
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
		MockSearchFHIRMedicationStatementFn: func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
			return &domain.FHIRMedicationStatementRelayConnection{}, nil
		},
		MockGetFHIROrganizationFn: func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
			id := uuid.New().String()
			name := "Test Organisation"
			return &domain.FHIROrganizationRelayPayload{
				Resource: &domain.FHIROrganization{
					ID:   &id,
					Name: &name,
				},
			}, nil
		},
		MockCreateFHIRMedicationStatementFn: func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
			return &domain.FHIRMedicationStatementRelayPayload{}, nil
		},
		MockCreateFHIRMedicationFn: func(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
			return &domain.FHIRMedicationRelayPayload{}, nil
		},
		MockCreateFHIRMediaFn: func(ctx context.Context, input domain.FHIRMedia) (*domain.FHIRMedia, error) {
			id := uuid.New().String()
			url := gofakeit.URL()
			title := gofakeit.BeerName()
			return &domain.FHIRMedia{
				Status: "",
				Subject: &domain.FHIRReferenceInput{
					ID: &id,
				},
				Content: &domain.FHIRAttachmentInput{
					URL:   (*scalarutils.URL)(&url),
					Title: &title,
				},
			}, nil
		},
		MockCreateFHIRPatientFn: func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
			male := domain.PatientGenderEnumMale
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID:         new(string),
					Text:       &domain.FHIRNarrative{},
					Gender:     &male,
					Identifier: []*domain.FHIRIdentifier{},
					Active:     new(bool),
					Name: []*domain.FHIRHumanName{
						{
							Text: gofakeit.Name(),
						},
					},
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
			}, nil
		},
		MockPatchFHIRPatientFn: func(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error) {
			male := domain.PatientGenderEnumMale
			return &domain.FHIRPatient{
				ID:         new(string),
				Text:       &domain.FHIRNarrative{},
				Gender:     &male,
				Identifier: []*domain.FHIRIdentifier{},
				Active:     new(bool),
				Name: []*domain.FHIRHumanName{
					{
						Text: gofakeit.Name(),
					},
				},
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
		MockPatchFHIREpisodeOfCareFn: func(ctx context.Context, id string, input domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error) {
			return &domain.FHIREpisodeOfCare{
				ID:            new(string),
				Text:          &domain.FHIRNarrative{},
				Identifier:    []*domain.FHIRIdentifier{},
				Status:        new(domain.EpisodeOfCareStatusEnum),
				StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
				Type:          []*domain.FHIRCodeableConcept{},
				Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
				Patient: &domain.FHIRReference{
					ID: new(string),
				},
				ManagingOrganization: &domain.FHIRReference{},
				Period:               &domain.FHIRPeriod{},
				ReferralRequest:      []*domain.FHIRReference{},
				CareManager:          &domain.FHIRReference{},
				Team:                 []*domain.FHIRReference{},
				Account:              []*domain.FHIRReference{},
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
		MockSearchFHIRPatientFn: func(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error) {
			return &domain.PatientConnection{
				Edges:    []*domain.PatientEdge{},
				PageInfo: &firebasetools.PageInfo{},
			}, nil
		},
		MockSearchPatientObservationsFn: func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
			uuid := uuid.New().String()
			instant := gofakeit.TimeZone()
			finalStatus := domain.ObservationStatusEnumFinal
			return &domain.PagedFHIRObservations{
				Observations: []domain.FHIRObservation{
					{
						ID:     &uuid,
						Status: &finalStatus,
						Subject: &domain.FHIRReference{
							ID: &uuid,
						},
						Encounter: &domain.FHIRReference{
							ID: &uuid,
						},
						Code: domain.FHIRCodeableConcept{
							ID: new(string),
							Coding: []*domain.FHIRCoding{
								{
									ID:           new(string),
									Version:      new(string),
									Code:         "",
									Display:      "Vital",
									UserSelected: new(bool),
								},
							},
							Text: "",
						},
						EffectiveInstant: (*scalarutils.Instant)(&instant),
					},
					{
						ID:     &uuid,
						Status: &finalStatus,
						Encounter: &domain.FHIRReference{
							ID: &uuid,
						},
						Code: domain.FHIRCodeableConcept{
							ID: new(string),
							Coding: []*domain.FHIRCoding{
								{
									ID:           new(string),
									Version:      new(string),
									Code:         "",
									Display:      "Vital",
									UserSelected: new(bool),
								},
							},
							Text: "",
						},
						EffectiveInstant: (*scalarutils.Instant)(&instant),
					},
					{
						ID:     &uuid,
						Status: &finalStatus,
						Encounter: &domain.FHIRReference{
							ID: &uuid,
						},
						Subject: &domain.FHIRReference{},
						Code: domain.FHIRCodeableConcept{
							ID: new(string),
							Coding: []*domain.FHIRCoding{
								{
									ID:           new(string),
									Version:      new(string),
									Code:         "",
									Display:      "Vital",
									UserSelected: new(bool),
								},
							},
							Text: "",
						},
						EffectiveInstant: (*scalarutils.Instant)(&instant),
					},
				},
				HasNextPage:     true,
				NextCursor:      "",
				HasPreviousPage: true,
				PreviousCursor:  "",
				TotalCount:      0,
			}, nil
		},
	}
}

// CreateEpisodeOfCare is a mock implementation of CreateEpisodeOfCare method
func (fh *FHIRMock) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error) {
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
func (fh *FHIRMock) SearchFHIROrganization(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIROrganizationRelayConnection, error) {
	return fh.MockSearchFHIROrganizationFn(ctx, params, tenant, pagination)
}

// SearchEpisodesByParam is a mock implementation of SearchEpisodesByParam method
func (fh *FHIRMock) SearchEpisodesByParam(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.MockSearchEpisodesByParamFn(ctx, searchParams, tenant, pagination)
}

// HasOpenEpisode is a mock implementation of HasOpenEpisode method
func (fh *FHIRMock) HasOpenEpisode(ctx context.Context, patient domain.FHIRPatient, tenant dto.TenantIdentifiers, pagination dto.Pagination) (bool, error) {
	return fh.MockHasOpenEpisodeFn(ctx, patient, tenant, pagination)
}

// OpenEpisodes is a mock implementation of OpenEpisodes method
func (fh *FHIRMock) OpenEpisodes(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error) {
	return fh.MockOpenEpisodesFn(ctx, patientReference, tenant, pagination)
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
func (fh *FHIRMock) SearchPatientEncounters(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
	return fh.MockSearchPatientEncountersFn(ctx, patientReference, status, tenant, pagination)
}

// SearchFHIREpisodeOfCare is a mock implementation of SearchFHIREpisodeOfCare method
func (fh *FHIRMock) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
	return fh.MockSearchFHIREpisodeOfCareFn(ctx, params, tenant, pagination)
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
func (fh *FHIRMock) SearchEpisodeEncounter(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
	return fh.MockSearchEpisodeEncounterFn(ctx, episodeReference, tenant, pagination)
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
func (fh *FHIRMock) GetActiveEpisode(ctx context.Context, episodeID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCare, error) {
	return fh.MockGetActiveEpisodeFn(ctx, episodeID, tenant, pagination)
}

// SearchFHIRServiceRequest is a mock implementation of SearchFHIRServiceRequest method
func (fh *FHIRMock) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRServiceRequestRelayConnection, error) {
	return fh.MockSearchFHIRServiceRequestFn(ctx, params, tenant, pagination)
}

// CreateFHIRServiceRequest is a mock implementation of CreateFHIRServiceRequest method
func (fh *FHIRMock) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	return fh.MockCreateFHIRServiceRequestFn(ctx, input)
}

// SearchFHIRAllergyIntolerance is a mock implementation of SearchFHIRAllergyIntolerance method
func (fh *FHIRMock) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
	return fh.MockSearchFHIRAllergyIntoleranceFn(ctx, params, tenant, pagination)
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
func (fh *FHIRMock) SearchFHIRComposition(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRCompositionRelayConnection, error) {
	return fh.MockSearchFHIRCompositionFn(ctx, params, tenant, pagination)
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
func (fh *FHIRMock) SearchFHIRCondition(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
	return fh.MockSearchFHIRConditionFn(ctx, params, tenant, pagination)
}

// UpdateFHIRCondition is a mock implementation of UpdateFHIRCondition method
func (fh *FHIRMock) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	return fh.MockUpdateFHIRConditionFn(ctx, input)
}

// GetFHIREncounter is a mock implementation of GetFHIREncounter method
func (fh *FHIRMock) GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
	return fh.MockGetFHIREncounterFn(ctx, id)
}

// PatchFHIREncounter is a mock implementation of PatchFHIREncounter method
func (fh *FHIRMock) PatchFHIREncounter(ctx context.Context, encounterID string, input domain.FHIREncounterInput) (*domain.FHIREncounter, error) {
	return fh.MockPatchFHIREncounterFn(ctx, encounterID, input)
}

// SearchFHIREncounter is a mock implementation of SearchFHIREncounter method
func (fh *FHIRMock) SearchFHIREncounter(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
	return fh.MockSearchFHIREncounterFn(ctx, params, tenant, pagination)
}

// SearchFHIRMedicationRequest is a mock implementation of SearchFHIRMedicationRequest method
func (fh *FHIRMock) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationRequestRelayConnection, error) {
	return fh.MockSearchFHIRMedicationRequestFn(ctx, params, tenant, pagination)
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
func (fh *FHIRMock) SearchFHIRObservation(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
	return fh.MockSearchFHIRObservationFn(ctx, params, tenant, pagination)
}

// CreateFHIRObservation is a mock implementation of CreateFHIRObservation method
func (fh *FHIRMock) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
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
func (fh *FHIRMock) SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
	return fh.MockSearchFHIRMedicationStatementFn(ctx, params, tenant, pagination)
}

// GetFHIROrganization is a mock implementation of GetFHIROrganization method
func (fh *FHIRMock) GetFHIROrganization(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.MockGetFHIROrganizationFn(ctx, organizationID)
}

// CreateFHIRMedicationStatement is a mock implementation of CreateFHIRMedicationStatement method
func (fh *FHIRMock) CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
	return fh.MockCreateFHIRMedicationStatementFn(ctx, input)
}

// CreateFHIRMedication is a mock implementation of CreateFHIRMedication method
func (fh *FHIRMock) CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error) {
	return fh.MockCreateFHIRMedicationFn(ctx, input)
}

// CreateFHIRPatient mocks the implementation of creating a FHIR patient
func (fh *FHIRMock) CreateFHIRPatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	return fh.MockCreateFHIRPatientFn(ctx, input)
}

// PatchFHIRPatient mocks the implementation for patching a fhir patient
func (fh *FHIRMock) PatchFHIRPatient(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error) {
	return fh.MockPatchFHIRPatientFn(ctx, id, input)
}

// PatchFHIREpisodeOfCare mocks the implementation of patching a FHIR episode of care
func (fh *FHIRMock) PatchFHIREpisodeOfCare(ctx context.Context, id string, input domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error) {
	return fh.MockPatchFHIREpisodeOfCareFn(ctx, id, input)
}

// UpdateFHIREpisodeOfCare mocks the implementation of updating a FHIR episode of care
func (fh *FHIRMock) UpdateFHIREpisodeOfCare(ctx context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error) {
	return fh.MockUpdateFHIREpisodeOfCareFn(ctx, fhirResourceID, payload)
}

// SearchFHIRPatient mocks the implementation of searching a FHIR patient
func (fh *FHIRMock) SearchFHIRPatient(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error) {
	return fh.MockSearchFHIRPatientFn(ctx, searchParams, tenant, pagination)
}

// SearchPatientObservations mocks the implementation of searching patient observations
func (fh *FHIRMock) SearchPatientObservations(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
	return fh.MockSearchPatientObservationsFn(ctx, patientReference, conceptID, tenant, pagination)
}

// GetFHIRAllergyIntolerance mocks the implementation of getting a resource by its ID
func (fh *FHIRMock) GetFHIRAllergyIntolerance(ctx context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.MockGetFHIRAllergyIntoleranceFn(ctx, id)
}

// SearchPatientAllergyIntolerance mocks the getting of patient allergies
func (fh *FHIRMock) SearchPatientAllergyIntolerance(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
	return fh.MockSearchPatientAllergyIntoleranceFn(ctx, patientReference, tenant, pagination)
}

// SearchPatientAllergyIntolerance mocks the getting of patient allergies
func (fh *FHIRMock) CreateFHIRMedia(ctx context.Context, input domain.FHIRMedia) (*domain.FHIRMedia, error) {
	return fh.MockCreateFHIRMediaFn(ctx, input)
}

// SearchPatentMedia mocks the searching of patient media
func (fh *FHIRMock) SearchPatientMedia(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedia, error) {
	return fh.MockSearchPatientMediaFn(ctx, patientReference, tenant, pagination)
}
