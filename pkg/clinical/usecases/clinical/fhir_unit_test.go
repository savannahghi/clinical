package clinical_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhir/mock"
	fakeDatasetMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhirdataset/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"
)

func TestUsecaseImpl_CreateEpisodeOfCare_Unittest(t *testing.T) {
	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	UUID := uuid.New().String()
	PatientRef := "Patient/1"
	OrgRef := "Organization/1"
	status := domain.EpisodeOfCareStatusEnumFinished

	episode := &domain.FHIREpisodeOfCare{
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
	}

	type args struct {
		ctx     context.Context
		episode domain.FHIREpisodeOfCare
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				episode: *episode,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				episode: domain.FHIREpisodeOfCare{
					ID:         nil,
					Text:       &domain.FHIRNarrative{},
					Identifier: []*domain.FHIRIdentifier{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockCreateEpisodeOfCareFn = func(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
					return nil, fmt.Errorf("Error creating episode of care")
				}
			}
			_, err := m.CreateEpisodeOfCare(tt.args.ctx, tt.args.episode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUsecaseStoreImpl.CreateEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_CreateFHIRCondition(t *testing.T) {
	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	ID := uuid.New().String()

	fhirconditionInput := &domain.FHIRConditionInput{
		ID:                 &ID,
		Identifier:         []*domain.FHIRIdentifierInput{},
		ClinicalStatus:     &domain.FHIRCodeableConceptInput{},
		VerificationStatus: &domain.FHIRCodeableConceptInput{},
		Category:           []*domain.FHIRCodeableConceptInput{},
		Severity:           &domain.FHIRCodeableConceptInput{},
		Code:               &domain.FHIRCodeableConceptInput{},
		BodySite:           []*domain.FHIRCodeableConceptInput{},
		Subject:            &domain.FHIRReferenceInput{},
		Encounter:          &domain.FHIRReferenceInput{},
		OnsetDateTime: &scalarutils.Date{
			Year:  2000,
			Month: 3,
			Day:   30,
		},
		OnsetAge:    &domain.FHIRAgeInput{},
		OnsetPeriod: &domain.FHIRPeriodInput{},
		OnsetRange:  &domain.FHIRRangeInput{},
		OnsetString: new(string),
		AbatementDateTime: &scalarutils.Date{
			Year:  2000,
			Month: 3,
			Day:   30,
		},
		AbatementAge:    &domain.FHIRAgeInput{},
		AbatementPeriod: &domain.FHIRPeriodInput{},
		AbatementRange:  &domain.FHIRRangeInput{},
		AbatementString: new(string),
		RecordedDate: &scalarutils.Date{
			Year:  2000,
			Month: 3,
			Day:   30,
		},
		Recorder: &domain.FHIRReferenceInput{},
		Asserter: &domain.FHIRReferenceInput{},
		Stage:    []*domain.FHIRConditionStageInput{},
		Evidence: []*domain.FHIRConditionEvidenceInput{},
		Note:     []*domain.FHIRAnnotationInput{},
	}

	invalidfhirconditionInput := &domain.FHIRConditionInput{
		ID:                 nil,
		Identifier:         []*domain.FHIRIdentifierInput{},
		ClinicalStatus:     &domain.FHIRCodeableConceptInput{},
		VerificationStatus: &domain.FHIRCodeableConceptInput{},
		Category:           []*domain.FHIRCodeableConceptInput{},
		Severity:           &domain.FHIRCodeableConceptInput{},
		Code:               &domain.FHIRCodeableConceptInput{},
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRConditionInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: *fhirconditionInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: *invalidfhirconditionInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockCreateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					return nil, fmt.Errorf("Error creating fhir condition")
				}
			}
			_, err := m.CreateFHIRCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseStoreImpl.CreateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_CreateFHIROrganization_Unittest(t *testing.T) {
	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	ID := ksuid.New().String()
	active := true
	testname := gofakeit.FirstName()

	orgInput := &domain.FHIROrganizationInput{
		ID:         &ID,
		Active:     &active,
		Identifier: []*domain.FHIRIdentifierInput{},
		Type:       []*domain.FHIRCodeableConceptInput{},
		Name:       &testname,
		Alias:      []string{"alias test"},
		Telecom:    []*domain.FHIRContactPointInput{},
		Address:    []*domain.FHIRAddressInput{},
	}

	invalidOrgInput := &domain.FHIROrganizationInput{
		ID:         &ID,
		Active:     new(bool),
		Identifier: []*domain.FHIRIdentifierInput{},
		Type:       []*domain.FHIRCodeableConceptInput{},
		Alias:      []string{"alias test"},
		Telecom:    []*domain.FHIRContactPointInput{},
		Address:    []*domain.FHIRAddressInput{},
	}

	type args struct {
		ctx   context.Context
		input domain.FHIROrganizationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: *orgInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: *invalidOrgInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockCreateFHIROrganizationFn = func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("Error creating fhir organization")
				}
			}
			_, err := m.CreateFHIROrganization(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_FindOrganizationByID_Unittest(t *testing.T) {
	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	type args struct {
		ctx            context.Context
		organizationID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIROrganizationRelayPayload
		wantErr bool
	}{
		{
			name: "Sad case",
			args: args{
				ctx:            ctx,
				organizationID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockFindOrganizationByIDFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("Error finding organization")
				}
			}

			got, err := m.FindOrganizationByID(tt.args.ctx, tt.args.organizationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.FindOrganizationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected response to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_SearchEpisodesByParam(t *testing.T) {
	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	type args struct {
		ctx          context.Context
		searchParams url.Values
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				searchParams: map[string][]string{
					"patient": {"12233"},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockSearchEpisodesByParamFn = func(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
					return nil, fmt.Errorf("Error searching episodes")
				}
			}
			_, err := m.SearchEpisodesByParam(tt.args.ctx, tt.args.searchParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseStoreImpl.SearchEpisodesByParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_SearchFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"patient": "12233",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"patient": "12233",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockSearchFHIREpisodeOfCareFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
					return nil, fmt.Errorf("Error searching episodes")
				}
			}
			_, err := m.SearchFHIREpisodeOfCare(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseStoreImpl.SearchFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_OpenEpisodes(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	patientReference := fmt.Sprintf("Patient/%s", ksuid.New().String())

	type args struct {
		ctx              context.Context
		patientReference string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:              ctx,
				patientReference: patientReference,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:              ctx,
				patientReference: patientReference,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockOpenEpisodesFn = func(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
					return nil, fmt.Errorf("Error searching episodes")
				}
			}
			_, err := m.OpenEpisodes(tt.args.ctx, tt.args.patientReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.OpenEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_HasOpenEpisode(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	UUID := uuid.New().String()

	patient := domain.FHIRPatient{
		ID: &UUID,
	}

	type args struct {
		ctx     context.Context
		patient domain.FHIRPatient
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{

		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				patient: patient,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				patient: patient,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockHasOpenEpisodeFn = func(ctx context.Context, patient domain.FHIRPatient) (bool, error) {
					return false, fmt.Errorf("Error searching episodes")
				}
			}
			_, err := m.HasOpenEpisode(tt.args.ctx, tt.args.patient)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.HasOpenEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_CreateFHIREncounter(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	input := domain.FHIREncounterInput{}
	type args struct {
		ctx   context.Context
		input domain.FHIREncounterInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIREncounterRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockCreateFHIREncounterFn = func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("Error creating encounter")
				}
			}
			_, err := m.CreateFHIREncounter(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestUsecaseStoreImpl_GetFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	id := ksuid.New().String()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockGetFHIREpisodeOfCareFn = func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
					return nil, fmt.Errorf("Error getting episode of care")
				}
			}
			_, err := m.GetFHIREpisodeOfCare(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_StartEncounter(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	episodeID := uuid.New().String()

	type args struct {
		ctx       context.Context
		episodeID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				episodeID: episodeID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockStartEncounterFn = func(ctx context.Context, episodeID string) (string, error) {
					return "", fmt.Errorf("Error starting encounter")
				}
			}
			_, err := m.StartEncounter(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_GetFHIREncounter(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("Error getting encounter")
				}
			}
			_, err := m.GetFHIREncounter(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseStoreImpl.GetFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_CreateFHIRServiceRequest(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	UUID := uuid.New().String()

	type args struct {
		ctx   context.Context
		input domain.FHIRServiceRequestInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: domain.FHIRServiceRequestInput{
					ID: &UUID,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockCreateFHIRServiceRequestFn = func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
					return nil, fmt.Errorf("Error creating service request")
				}
			}
			_, err := m.CreateFHIRServiceRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseStoreImpl.CreateFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_CreateFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	UUID := uuid.New().String()

	type args struct {
		ctx   context.Context
		input domain.FHIRAllergyIntoleranceInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: domain.FHIRAllergyIntoleranceInput{
					ID: &UUID,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockCreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("Error creating allergy intolerance")
				}
			}
			_, err := m.CreateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseStoreImpl.CreateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_UpdateFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	UUID := uuid.New().String()

	input := domain.FHIRAllergyIntoleranceInput{
		ID: &UUID,
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRAllergyIntoleranceInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRAllergyIntoleranceRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockUpdateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("Error updating allergy intolerance")
				}
			}
			_, err := m.UpdateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_SearchFHIRComposition(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	params := map[string]interface{}{"test": "123"}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRCompositionRelayConnection
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				params: params,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				params: params,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockSearchFHIRCompositionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
					return nil, fmt.Errorf("Error searching composition")
				}
			}

			_, err := m.SearchFHIRComposition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_CreateFHIRComposition(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	input := domain.FHIRCompositionInput{}

	type args struct {
		ctx   context.Context
		input domain.FHIRCompositionInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRCompositionRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockCreateFHIRCompositionFn = func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
					return nil, fmt.Errorf("Error creating composition")
				}
			}

			_, err := m.CreateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_UpdateFHIRComposition(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	UUID := uuid.New().String()

	input := domain.FHIRCompositionInput{
		ID: &UUID,
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRCompositionInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRCompositionRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockUpdateFHIRCompositionFn = func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
					return nil, fmt.Errorf("Error updating composition")
				}
			}
			_, err := m.UpdateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseStoreImpl_DeleteFHIRComposition(t *testing.T) {
	ctx := context.Background()
	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	id := ksuid.New().String()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockDeleteFHIRCompositionFn = func(ctx context.Context, id string) (bool, error) {
					return false, fmt.Errorf("Error deleting composition")
				}
			}
			got, err := m.DeleteFHIRComposition(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRComposition() = %v, want %v", got, tt.want)
			}
		})
	}
}
