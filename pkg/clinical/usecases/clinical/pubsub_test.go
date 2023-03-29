package clinical_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeMyCarehubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/firebasetools"
)

func TestUseCasesClinicalImpl_CreatePubsubPatient(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx     context.Context
		payload dto.CreatePatientPubSubMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create pubsub patient",
			args: args{
				ctx: ctx,
				payload: dto.CreatePatientPubSubMessage{
					ID:         uuid.New().String(),
					Active:     true,
					Counselled: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user profile",
			args: args{
				ctx: ctx,
				payload: dto.CreatePatientPubSubMessage{
					ID:         uuid.New().String(),
					Active:     true,
					Counselled: false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to check patient existence",
			args: args{
				ctx: ctx,
				payload: dto.CreatePatientPubSubMessage{
					ID:         uuid.New().String(),
					Active:     true,
					Counselled: false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Patient already exists",
			args: args{
				ctx: ctx,
				payload: dto.CreatePatientPubSubMessage{
					ID:         uuid.New().String(),
					Active:     true,
					Counselled: false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create patient",
			args: args{
				ctx: ctx,
				payload: dto.CreatePatientPubSubMessage{
					ID:         uuid.New().String(),
					Active:     true,
					Counselled: false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to add FHIR ID to profile",
			args: args{
				ctx: ctx,
				payload: dto.CreatePatientPubSubMessage{
					ID:         uuid.New().String(),
					Active:     true,
					Counselled: false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid patient input",
			args: args{
				ctx: ctx,
				payload: dto.CreatePatientPubSubMessage{
					ID:         uuid.New().String(),
					Active:     true,
					Counselled: false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to get user profile" {
				fakeMCH.MockUserProfileFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "Sad Case - Fail to check patient existence" {
				fakeFHIR.MockSearchFHIRPatientFn = func(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers) (*domain.PatientConnection, error) {
					return nil, fmt.Errorf("fail to get a user")
				}
			}

			if tt.name == "Sad Case - Patient already exists" {
				fakeFHIR.MockSearchFHIRPatientFn = func(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers) (*domain.PatientConnection, error) {
					return &domain.PatientConnection{
						Edges: []*domain.PatientEdge{
							{
								Cursor: "",
								Node:   &domain.FHIRPatient{},
							},
							{
								Cursor: "",
								Node:   &domain.FHIRPatient{},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to create patient" {
				fakeFHIR.MockCreateFHIRPatientFn = func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
					return nil, fmt.Errorf("failed to create patient")
				}
			}

			if tt.name == "Sad Case - Fail to add FHIR ID to profile" {
				fakeMCH.MockAddFHIRIDToPatientProfileFn = func(ctx context.Context, fhirID, clientID string) error {
					return fmt.Errorf("failed to add fhir ID to profile")
				}
			}

			if tt.name == "Sad Case - invalid patient input" {
				fakeMCH.MockUserProfileFn = func(ctx context.Context, userID string) (*domain.User, error) {
					dob := time.Now()
					return &domain.User{
						ID:       new(string),
						Username: gofakeit.Username(),
						UserType: "STAFF",
						Name:     gofakeit.Name(),
						Gender:   "MALE",
						Active:   false,
						Flavour:  "PRO",
						Avatar:   "",
						Contacts: &domain.Contact{
							ID:           new(string),
							ContactType:  "PHONE",
							ContactValue: "gofakeit.PhoneFormatted()",
							Active:       true,
							OptedIn:      true,
						},
						DateOfBirth: &dob,
					}, nil
				}
			}

			if err := u.CreatePubsubPatient(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubPatient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubOrganization(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx  context.Context
		data dto.CreateFacilityPubSubMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create pubsub organization",
			args: args{
				ctx: ctx,
				data: dto.CreateFacilityPubSubMessage{
					ID:          new(string),
					Name:        "Test Facility",
					Code:        0,
					Phone:       "",
					Active:      false,
					County:      "",
					Description: "",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create pubsub organization",
			args: args{
				ctx: ctx,
				data: dto.CreateFacilityPubSubMessage{
					ID:          new(string),
					Name:        "Test Facility",
					Code:        0,
					Phone:       "",
					Active:      false,
					County:      "",
					Description: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to add fhir id to facility",
			args: args{
				ctx: ctx,
				data: dto.CreateFacilityPubSubMessage{
					ID:          new(string),
					Name:        "Test Facility",
					Code:        0,
					Phone:       "",
					Active:      false,
					County:      "",
					Description: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to create pubsub organization" {
				fakeFHIR.MockCreateFHIROrganizationFn = func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create organization")
				}
			}

			if tt.name == "Sad Case - Fail to add fhir id to facility" {
				fakeMCH.MockAddFHIRIDToFacilityFn = func(ctx context.Context, fhirID, facilityID string) error {
					return fmt.Errorf("failed to add fhir ID to facility")
				}
			}

			if err := u.CreatePubsubOrganization(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubVitals(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx  context.Context
		data dto.CreateVitalSignPubSubMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create pubsub vitals",
			args: args{
				ctx: ctx,
				data: dto.CreateVitalSignPubSubMessage{
					PatientID:      uuid.NewString(),
					OrganizationID: "",
					Name:           "",
					ConceptID:      new(string),
					Value:          "",
					Date:           time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create pubsub vitals - available organizationID",
			args: args{
				ctx: ctx,
				data: dto.CreateVitalSignPubSubMessage{
					PatientID:      uuid.NewString(),
					OrganizationID: uuid.NewString(),
					Name:           "",
					ConceptID:      new(string),
					Value:          "",
					Date:           time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to find patient",
			args: args{
				ctx: ctx,
				data: dto.CreateVitalSignPubSubMessage{
					PatientID:      uuid.NewString(),
					OrganizationID: uuid.NewString(),
					Name:           "",
					ConceptID:      new(string),
					Value:          "",
					Date:           time.Time{},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to find organisation using org ID",
			args: args{
				ctx: ctx,
				data: dto.CreateVitalSignPubSubMessage{
					PatientID:      uuid.NewString(),
					OrganizationID: uuid.NewString(),
					Name:           "",
					ConceptID:      new(string),
					Value:          "",
					Date:           time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create observation",
			args: args{
				ctx: ctx,
				data: dto.CreateVitalSignPubSubMessage{
					PatientID:      uuid.NewString(),
					OrganizationID: uuid.NewString(),
					Name:           "",
					ConceptID:      new(string),
					Value:          "",
					Date:           time.Time{},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get ciel concept",
			args: args{
				ctx: ctx,
				data: dto.CreateVitalSignPubSubMessage{
					PatientID:      uuid.NewString(),
					OrganizationID: uuid.NewString(),
					Name:           "",
					ConceptID:      new(string),
					Value:          "",
					Date:           time.Time{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to find patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get fhir patient")
				}
			}

			if tt.name == "Sad Case - Fail to find organisation using org ID" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to find org by ID")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			if tt.name == "Sad Case - fail to get ciel concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if err := u.CreatePubsubVitals(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubVitals() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubTestResult(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx  context.Context
		data dto.CreatePatientTestResultPubSubMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create test result",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientTestResultPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: "",
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Result:         dto.TestResult{},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create test result - with organisation ID",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientTestResultPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Result:         dto.TestResult{},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get fhir patient",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientTestResultPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: "",
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Result:         dto.TestResult{},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get organisation",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientTestResultPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Result:         dto.TestResult{},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to create observation",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientTestResultPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Result:         dto.TestResult{},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get ciel concept",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientTestResultPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Result:         dto.TestResult{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - fail to get fhir patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get fhir patient")
				}
			}

			if tt.name == "Sad Case - fail to get organisation" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to find org by ID")
				}
			}

			if tt.name == "Sad Case - fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			if tt.name == "Sad Case - fail to get ciel concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if err := u.CreatePubsubTestResult(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubTestResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubMedicationStatement(t *testing.T) {
	ctx := context.Background()
	conceptID := "12345"
	type args struct {
		ctx  context.Context
		data dto.CreateMedicationPubSubMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create medication statement",
			args: args{
				ctx: ctx,
				data: dto.CreateMedicationPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: "",
					Name:           "",
					ConceptID:      &conceptID,
					Date:           time.Time{},
					Value:          "",
					Drug: &dto.MedicationDrug{
						ConceptID: &conceptID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully create medication statement - with organisation ID",
			args: args{
				ctx: ctx,
				data: dto.CreateMedicationPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "",
					ConceptID:      &conceptID,
					Date:           time.Time{},
					Value:          "",
					Drug: &dto.MedicationDrug{
						ConceptID: &conceptID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get patient",
			args: args{
				ctx: ctx,
				data: dto.CreateMedicationPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "",
					ConceptID:      &conceptID,
					Date:           time.Time{},
					Value:          "",
					Drug: &dto.MedicationDrug{
						ConceptID: &conceptID,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get organisation",
			args: args{
				ctx: ctx,
				data: dto.CreateMedicationPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "",
					ConceptID:      &conceptID,
					Date:           time.Time{},
					Value:          "",
					Drug: &dto.MedicationDrug{
						ConceptID: &conceptID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to get medication statement",
			args: args{
				ctx: ctx,
				data: dto.CreateMedicationPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "",
					ConceptID:      &conceptID,
					Date:           time.Time{},
					Value:          "",
					Drug: &dto.MedicationDrug{
						ConceptID: &conceptID,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get ciel concept",
			args: args{
				ctx: ctx,
				data: dto.CreateMedicationPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "",
					ConceptID:      &conceptID,
					Date:           time.Time{},
					Value:          "",
					Drug: &dto.MedicationDrug{
						ConceptID: &conceptID,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get fhir patient")
				}
			}

			if tt.name == "Sad Case - fail to get organisation" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to find org by ID")
				}
			}

			if tt.name == "Sad Case - fail to get medication statement" {
				fakeFHIR.MockCreateFHIRMedicationStatementFn = func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
					return nil, fmt.Errorf("failed to create medication statement")
				}
			}

			if tt.name == "Sad Case - fail to get ciel concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if err := u.CreatePubsubMedicationStatement(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubMedicationStatement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreatePubsubAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx  context.Context
		data dto.CreatePatientAllergyPubSubMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create allergy intolerance",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientAllergyPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: "",
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Reaction:       dto.AllergyReaction{},
					Severity:       dto.AllergySeverity{},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user profile",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientAllergyPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: "",
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Reaction:       dto.AllergyReaction{},
					Severity:       dto.AllergySeverity{},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create allergy intolerance",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientAllergyPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: "",
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Reaction:       dto.AllergyReaction{},
					Severity:       dto.AllergySeverity{},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get ciel concept",
			args: args{
				ctx: ctx,
				data: dto.CreatePatientAllergyPubSubMessage{
					PatientID:      uuid.New().String(),
					OrganizationID: "",
					Name:           "",
					ConceptID:      new(string),
					Date:           time.Time{},
					Reaction:       dto.AllergyReaction{},
					Severity:       dto.AllergySeverity{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to get user profile" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get patient")
				}
			}

			if tt.name == "Sad Case - Fail to create allergy intolerance" {
				fakeFHIR.MockCreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("failed to create allergy intolerance")
				}
			}

			if tt.name == "Sad Case - fail to get ciel concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if err := u.CreatePubsubAllergyIntolerance(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreatePubsubAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_getConcept(t *testing.T) {
	type args struct {
		ctx       context.Context
		source    dto.TerminologySource
		conceptID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sad case: failed to get icd10 concept",
			args: args{
				ctx:       context.Background(),
				source:    dto.TerminologySourceICD10,
				conceptID: gofakeit.BS(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get ciel concept",
			args: args{
				ctx:       context.Background(),
				source:    dto.TerminologySourceCIEL,
				conceptID: gofakeit.BS(),
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed to get snomed concept",
			args: args{
				ctx:       context.Background(),
				source:    dto.TerminologySourceSNOMEDCT,
				conceptID: gofakeit.BS(),
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed to get loinc concept",
			args: args{
				ctx:       context.Background(),
				source:    dto.TerminologySourceLOINC,
				conceptID: gofakeit.BS(),
			},
			wantErr: true,
		},

		{
			name: "Sad case: invalid concept source",
			args: args{
				ctx:       context.Background(),
				source:    dto.TerminologySource("invalid"),
				conceptID: gofakeit.BS(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad case: failed to get icd10 concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get ciel concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get snomed concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get loinc concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := c.GetConcept(tt.args.ctx, tt.args.source, tt.args.conceptID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.getConcept() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
