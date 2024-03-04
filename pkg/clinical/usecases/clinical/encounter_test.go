package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeAdvantageMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/advantage/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

func TestUseCasesClinicalImpl_StartEncounter(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		episodeID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully start an encounter",
			args: args{
				ctx:       ctx,
				episodeID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Empty episode of care ID",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get episode of care",
			args: args{
				ctx:       ctx,
				episodeID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create encounter",
			args: args{
				ctx:       ctx,
				episodeID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to get tenant identifiers",
			args: args{
				ctx:       ctx,
				episodeID: uuid.New().String(),
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
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to get episode of care" {
				fakeFHIR.MockGetFHIREpisodeOfCareFn = func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
					return nil, fmt.Errorf("failed to get episode of care")
				}
			}

			if tt.name == "Sad Case - Fail to create encounter" {
				fakeFHIR.MockCreateFHIREncounterFn = func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to create encounter")
				}
			}

			if tt.name == "Sad Case - failed to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			got, err := u.StartEncounter(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == "" {
					t.Errorf("expected an episode of care ID but got %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_PatchEncounter(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		encounterID string
		input       dto.EncounterInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully patch encounter",
			args: args{
				ctx:         ctx,
				encounterID: gofakeit.UUID(),
				input: dto.EncounterInput{
					Status: dto.EncounterStatusEnumInProgress,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid encounterID",
			args: args{
				ctx: ctx,
				input: dto.EncounterInput{
					Status: dto.EncounterStatusEnumInProgress,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get encounter",
			args: args{
				ctx:         ctx,
				encounterID: gofakeit.UUID(),
				input: dto.EncounterInput{
					Status: dto.EncounterStatusEnumCancelled,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to patch encounter",
			args: args{
				ctx:         ctx,
				encounterID: gofakeit.UUID(),
				input: dto.EncounterInput{
					Status: dto.EncounterStatusEnumFinished,
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
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Unable to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Unable to patch encounter" {
				fakeFHIR.MockPatchFHIREncounterFn = func(ctx context.Context, encounterID string, input domain.FHIREncounterInput) (*domain.FHIREncounter, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.PatchEncounter(ctx, tt.args.encounterID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_EndEncounter(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		encounterID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully end encounter",
			args: args{
				ctx:         ctx,
				encounterID: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Missing encounter ID",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to end encounter",
			args: args{
				ctx:         ctx,
				encounterID: uuid.New().String(),
			},
			want:    false,
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
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to end encounter" {
				fakeFHIR.MockEndEncounterFn = func(ctx context.Context, encounterID string) (bool, error) {
					return false, fmt.Errorf("failed to update encounter")
				}
			}

			got, err := u.EndEncounter(tt.args.ctx, tt.args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.EndEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesClinicalImpl.EndEncounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesClinicalImpl_ListPatientEncounters(t *testing.T) {
	ctx := context.Background()
	first := 3
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
			name: "Happy Case - Successfully list patient encounter",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
					Skip:  false,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing patient ID",
			args: args{
				ctx: ctx,
				pagination: &dto.Pagination{
					First: &first,
					Skip:  false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get fhir patient",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
					Skip:  false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get identifiers",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
					Skip:  false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get patient encounters",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
					Skip:  false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid pagination",
			args: args{
				ctx:       ctx,
				patientID: uuid.New().String(),
				pagination: &dto.Pagination{
					First: &first,
					Last:  &first,
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
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to get fhir patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get fhir patient")
				}
			}

			if tt.name == "Sad Case: Fail to get identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case - Fail to get patient encounters" {
				fakeFHIR.MockSearchPatientEncountersFn = func(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
					return nil, fmt.Errorf("error searching for encounters")
				}
			}

			got, err := u.ListPatientEncounters(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ListPatientEncounters() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesClinicalImpl_GetEncounterAssociatedResources(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		encounterID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully list encounter's all data",
			args: args{
				ctx:         ctx,
				encounterID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing encounter ID",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to search all fhir encounter data",
			args: args{
				ctx:         ctx,
				encounterID: uuid.New().String(),
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
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Missing encounter ID" {
				fakeFHIR.MockSearchFHIREncounterAllDataFn = func(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}
			if tt.name == "Sad Case - unable to search all fhir encounter data" {
				fakeFHIR.MockSearchFHIREncounterAllDataFn = func(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := u.GetEncounterAssociatedResources(tt.args.ctx, tt.args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetEncounterAssociatedResources() error = %v, wantErr %v", err, tt.wantErr)
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
