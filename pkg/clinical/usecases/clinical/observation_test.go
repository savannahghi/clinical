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
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/scalarutils"
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
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

func TestUseCasesClinicalImpl_RecordMuac(t *testing.T) {
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
			name: "Happy Case - Successfully record muac",
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
			name: "Sad Case - Fail to record muac",
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infa := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infa)

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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordMuac(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordMuac() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_RecordOxygenSaturation(t *testing.T) {
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
			name: "Happy Case - Successfully record oxygen saturation",
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
			name: "Sad Case - Fail to record oxygen saturation",
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infa := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infa)

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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			got, err := u.RecordOxygenSaturation(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordOxygenSaturation() error = %v, wantErr %v", err, tt.wantErr)
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
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
	first := 10
	ctx := context.Background()
	type args struct {
		ctx             context.Context
		patientID       string
		observationCode string
		pagination      *dto.Pagination
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
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:             ctx,
				observationCode: "1234",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient",
			args: args{
				ctx:             ctx,
				patientID:       uuid.New().String(),
				observationCode: "1234",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:             ctx,
				patientID:       uuid.New().String(),
				observationCode: "1234",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to search patient observations",
			args: args{
				ctx:             ctx,
				patientID:       uuid.New().String(),
				observationCode: "1234",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search observation - nil subject",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil subject id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil encounter",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
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
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("failed to search patient observations")
				}
			}

			if tt.name == "Happy case: get patient height" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Subject: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil subject" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     &UUID,
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Encounter: &domain.FHIRReference{
									ID: &UUID,
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil subject id" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Subject: &domain.FHIRReference{},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil encounter" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					instant := gofakeit.TimeZone()
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Subject: &domain.FHIRReference{
									ID: &UUID,
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								EffectiveInstant:     (*scalarutils.Instant)(&instant),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			got, err := u.GetPatientObservations(tt.args.ctx, tt.args.patientID, tt.args.observationCode, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientObservations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
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
	first := 10
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
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
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			got, err := u.GetPatientTemperatureEntries(tt.args.ctx, tt.args.patientID, tt.args.pagination)
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

func TestUseCasesClinicalImpl_GetPatientBloodPressureEntries(t *testing.T) {
	ctx := context.Background()
	first := 10
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient blood pressure entries",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			got, err := u.GetPatientBloodPressureEntries(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientBloodPressureEntries() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetHeight(t *testing.T) {
	first := 10
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get patient height",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			_, err := c.GetPatientHeightEntries(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientHeightEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientPulseRateEntries(t *testing.T) {
	first := 10
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get patient pulse rate",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			_, err := c.GetPatientPulseRateEntries(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientPulseRateEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientRespiratoryRateEntries(t *testing.T) {
	ctx := context.Background()
	first := 10
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient respiratory rate entries",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			got, err := u.GetPatientRespiratoryRateEntries(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientRespiratoryRateEntries() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetPatientBMIEntries(t *testing.T) {
	first := 10
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get patient bmi entries",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			got, err := u.GetPatientBMIEntries(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientBMIEntries() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetPatientWeightEntries(t *testing.T) {
	first := 10
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get patient pulse rate",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			_, err := c.GetPatientWeightEntries(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientWeightEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientMuacEntries(t *testing.T) {
	first := 10
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get patient muac",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			_, err := c.GetPatientMuacEntries(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientMuacEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientOxygenSaturationEntries(t *testing.T) {
	first := 10
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get patient oxygen saturation",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			_, err := c.GetPatientOxygenSaturationEntries(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientOxygenSaturationEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientViralLoad(t *testing.T) {
	first := 10
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get patient viral load",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					Skip: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid patient ID",
			args: args{
				ctx:       ctx,
				patientID: "invalid",
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Incorrect patient ID",
			args: args{
				ctx:       ctx,
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get tenant identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get patient observations",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get patient viral load",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.BeerMalt(),
				pagination: &dto.Pagination{
					First: &first,
					Skip:  true,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Invalid patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Incorrect patient ID" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("incorrect patient ID")
				}
			}

			if tt.name == "Sad Case - fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - fail to get patient observations" {
				fakeFHIR.MockSearchPatientObservationsFn = func(ctx context.Context, patientReference, conceptID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					return nil, fmt.Errorf("an error occured")
				}
			}

			_, err := c.GetPatientViralLoad(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientViralLoad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordViralLoad(t *testing.T) {
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
			name: "Happy Case - Successfully record viral load",
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
			name: "Sad Case - Fail to record viral load",
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

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
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("fail to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to create observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			_, err := c.RecordViralLoad(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordViralLoad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
