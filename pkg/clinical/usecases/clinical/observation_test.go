package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeMyCarehubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

func TestUseCasesClinicalImpl_RecordObservation(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx                context.Context
		input              dto.ObservationInput
		vitalSignConceptID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record observation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail validation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail validation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status: dto.ObservationStatusFinal,
					Value:  "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail validation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get CIEL concept",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
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

			if tt.name == "Sad Case - Fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}

			if tt.name == "Sad Case - return a finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get CIEL concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (map[string]interface{}, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordObservation(tt.args.ctx, tt.args.input, tt.args.vitalSignConceptID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordTemperature(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record temperature",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "12",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record temperature",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					EncounterID: uuid.New().String(),
					Value:       "12",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get CIEL concept",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
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

			if tt.name == "Sad Case - Fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}

			if tt.name == "Sad Case - return a finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get CIEL concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (map[string]interface{}, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordTemperature(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordTemperature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordHeight(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record height",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "185.21",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record height",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					EncounterID: uuid.New().String(),
					Value:       "12",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get CIEL concept",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
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

			if tt.name == "Sad Case - Fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}

			if tt.name == "Sad Case - return a finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get CIEL concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (map[string]interface{}, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordHeight(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordWeight(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record weight",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "185.21",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record weight",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					EncounterID: uuid.New().String(),
					Value:       "12",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get CIEL concept",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
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

			if tt.name == "Sad Case - Fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}

			if tt.name == "Sad Case - return a finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get CIEL concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (map[string]interface{}, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordWeight(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordWeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordRespiratoryRate(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record respiratory rate",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "185.21",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record respiratory rate",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					EncounterID: uuid.New().String(),
					Value:       "12",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get CIEL concept",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
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

			if tt.name == "Sad Case - Fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}

			if tt.name == "Sad Case - return a finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get CIEL concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (map[string]interface{}, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordRespiratoryRate(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordRespiratoryRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordPulseRate(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record pulse rate",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "185.21",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record pulse rate",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					EncounterID: uuid.New().String(),
					Value:       "12",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get CIEL concept",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
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

			if tt.name == "Sad Case - Fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}

			if tt.name == "Sad Case - return a finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get CIEL concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (map[string]interface{}, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordPulseRate(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordPulseRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordBloodPressure(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record blood pressure",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "185.21",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record blood pressure",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					EncounterID: uuid.New().String(),
					Value:       "12",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get CIEL concept",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
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

			if tt.name == "Sad Case - Fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}

			if tt.name == "Sad Case - return a finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get CIEL concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (map[string]interface{}, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordBloodPressure(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordBloodPressure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordBMI(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.ObservationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully record BMI",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "185.21",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record BMI",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					EncounterID: uuid.New().String(),
					Value:       "12",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - return a finished encounter",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get CIEL concept",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get tenant meta tags",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create observation",
			args: args{
				ctx: ctx,
				input: dto.ObservationInput{
					Status:      dto.ObservationStatusFinal,
					EncounterID: uuid.New().String(),
					Value:       "1234",
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

			if tt.name == "Sad Case - Fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}

			if tt.name == "Sad Case - return a finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get CIEL concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (map[string]interface{}, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordBMI(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordBMI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientObservations(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx             context.Context
		patientID       string
		observationCode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - successfully get patient observations",
			args: args{
				ctx:             ctx,
				patientID:       uuid.New().String(),
				observationCode: "1234",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:             ctx,
				observationCode: "1234",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient",
			args: args{
				ctx:             ctx,
				patientID:       uuid.New().String(),
				observationCode: "1234",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:             ctx,
				patientID:       uuid.New().String(),
				observationCode: "1234",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to search patient observations",
			args: args{
				ctx:             ctx,
				patientID:       uuid.New().String(),
				observationCode: "1234",
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

			if tt.name == "Sad Case - fail to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get patient")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to search patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers) ([]*domain.FHIRObservation, error) {
					return nil, fmt.Errorf("failed to search patient observations")
				}
			}

			got, err := u.GetPatientObservations(tt.args.ctx, tt.args.patientID, tt.args.observationCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientObservations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientTemperatureEntries(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		patientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient temperature entries",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx: ctx,
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

			got, err := u.GetPatientTemperatureEntries(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientTemperatureEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got %v", got)
					return
				}
			}
		})
	}
}
