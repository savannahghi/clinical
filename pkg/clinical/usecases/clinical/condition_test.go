package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeMyCarehubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/scalarutils"
)

func TestUseCasesClinicalImpl_CreateCondition(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.ConditionInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create condition",
			args: args{
				ctx: nil,
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error fetching concept",
			args: args{
				ctx: nil,
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get patient",
			args: args{
				ctx: nil,
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get encounter",
			args: args{
				ctx: nil,
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail in completed encounter",
			args: args{
				ctx: nil,
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get tags",
			args: args{
				ctx: nil,
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to  create condition",
			args: args{
				ctx: nil,
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "sad case: error fetching concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org string, source string, concept string, includeMappings bool, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if tt.name == "sad case: fail to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("fail to get patient")
				}
			}

			if tt.name == "sad case: fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to ge encounter")
				}
			}

			if tt.name == "sad case: fail to  create condition" {
				fakeFHIR.MockCreateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					return nil, fmt.Errorf("failed to create condition")
				}
			}

			if tt.name == "sad case: fail in completed encounter" {
				finished := domain.EncounterStatusEnumFinished
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: finished,
						},
					}, nil
				}
			}

			if tt.name == "sad case: fail to get tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tags")
				}
			}

			got, err := c.CreateCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ListPatientConditions(t *testing.T) {

	type args struct {
		ctx        context.Context
		patientID  string
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: list conditions",
			args: args{
				ctx:        context.Background(),
				patientID:  gofakeit.UUID(),
				pagination: dto.Pagination{},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid patient id",
			args: args{
				ctx:        context.Background(),
				patientID:  "invalid",
				pagination: dto.Pagination{},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get identifiers",
			args: args{
				ctx:        context.Background(),
				patientID:  gofakeit.UUID(),
				pagination: dto.Pagination{},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get patient",
			args: args{
				ctx:        context.Background(),
				patientID:  gofakeit.UUID(),
				pagination: dto.Pagination{},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to search condition",
			args: args{
				ctx:        context.Background(),
				patientID:  gofakeit.UUID(),
				pagination: dto.Pagination{},
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
			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "sad case: fail to get identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get identifiers")
				}
			}

			if tt.name == "sad case: fail to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get patient")
				}
			}

			if tt.name == "sad case: fail to search condition" {
				fakeFHIR.MockSearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
					return nil, fmt.Errorf("failed to find condition")
				}
			}

			got, err := c.ListPatientConditions(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPatientConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}
