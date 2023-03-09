package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/firebasetools"
	"github.com/segmentio/ksuid"
)

func TestUsecaseImpl_CreateFHIROrganization_Unittest(t *testing.T) {
	ctx := context.Background()

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
			FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			Fakefhir := fakeFHIRMock.NewFHIRMock()
			FakeOCL := fakeOCLMock.NewFakeOCLMock()

			infra := infrastructure.NewInfrastructureInteractor(FakeExt, Fakefhir, FakeOCL)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad case" {
				Fakefhir.MockCreateFHIROrganizationFn = func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("Error creating fhir organization")
				}
			}
			_, err := u.CreateFHIROrganization(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_PatientTimeline(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	type args struct {
		ctx       context.Context
		patientID string
		count     int
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "Happy case: patient timeline",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get allergy intolerance - nil node",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get allergy intolerance",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get observation",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get observation - nil node",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully get medication statement",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil node",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				count:     4,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to search allergy intolerance" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return &domain.FHIRAllergyIntoleranceRelayConnection{}, fmt.Errorf("failed to search allergy")
				}
			}

			if tt.name == "Happy Case - Successfully get allergy intolerance" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return &domain.FHIRAllergyIntoleranceRelayConnection{
						Edges: []*domain.FHIRAllergyIntoleranceRelayEdge{
							{
								Cursor: new(string),
								Node:   &domain.FHIRAllergyIntolerance{},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get allergy intolerance - nil node" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return &domain.FHIRAllergyIntoleranceRelayConnection{
						Edges: []*domain.FHIRAllergyIntoleranceRelayEdge{
							{
								Cursor: new(string),
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
					return &domain.FHIRObservationRelayConnection{}, fmt.Errorf("failed to get observation")
				}
			}

			if tt.name == "Happy Case - Successfully get observation" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node:   &domain.FHIRObservation{},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get observation - nil node" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
					return &domain.FHIRMedicationStatementRelayConnection{}, fmt.Errorf("failed to get medication statement")
				}
			}

			if tt.name == "Happy Case - Successfully get medication statement" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node:   &domain.FHIRMedicationStatement{},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement - nil node" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			got, err := u.PatientTimeline(tt.args.ctx, tt.args.patientID, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient timeline to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient timeline not to be nil for %v", tt.name)
				return
			}
		})
	}

}

func TestClinicalUseCaseImpl_PatientHealthTimeline(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	type args struct {
		ctx   context.Context
		input domain.HealthTimelineInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.HealthTimeline
		wantErr bool
	}{
		{
			name: "Happy case: patient timeline",
			args: args{
				ctx: ctx,
				input: domain.HealthTimelineInput{
					PatientID: gofakeit.UUID(),
					Offset:    0,
					Limit:     20,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			Fakefhir := fakeFHIRMock.NewFHIRMock()
			FakeOCL := fakeOCLMock.NewFakeOCLMock()

			infra := infrastructure.NewInfrastructureInteractor(FakeExt, Fakefhir, FakeOCL)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			got, err := u.PatientHealthTimeline(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientHealthTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient timeline to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient timeline not to be nil for %v", tt.name)
				return
			}
		})
	}

}

func TestClinicalUseCaseImpl_GetMedicalData(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	type args struct {
		ctx       context.Context
		patientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.MedicalData
		wantErr bool
	}{
		{
			name: "Happy case: patient timeline",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil node",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully search medication statement",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil node",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully search allergy intolerance",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search observation - nil node",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully search observation",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to search medication statement" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
					return nil, fmt.Errorf("failed to search medication statement")
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement - nil node" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Happy Case - Successfully search medication statement" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node:   &domain.FHIRMedicationStatement{},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return nil, fmt.Errorf("failed to search allergy intolerance")
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance - nil node" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return &domain.FHIRAllergyIntoleranceRelayConnection{
						Edges: []*domain.FHIRAllergyIntoleranceRelayEdge{
							{
								Cursor: new(string),
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Happy Case - Successfully search allergy intolerance" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return &domain.FHIRAllergyIntoleranceRelayConnection{
						Edges: []*domain.FHIRAllergyIntoleranceRelayEdge{
							{
								Cursor: new(string),
								Node:   &domain.FHIRAllergyIntolerance{},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
					return nil, fmt.Errorf("failed to search observation")
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil node" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Happy Case - Successfully search observation" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node:   &domain.FHIRObservation{},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			got, err := u.GetMedicalData(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.GetMedicalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient medical data to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient medical data not to be nil for %v", tt.name)
				return
			}
		})
	}

}
