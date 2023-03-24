package clinical_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeMyCarehubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/scalarutils"
	"testing"
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
					PatientID:   gofakeit.UUID(),
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
					PatientID:   gofakeit.UUID(),
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
					PatientID:   gofakeit.UUID(),
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
					PatientID:   gofakeit.UUID(),
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
					PatientID:   gofakeit.UUID(),
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
					PatientID:   gofakeit.UUID(),
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
					PatientID:   gofakeit.UUID(),
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

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "sad case: error fetching concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org string, source string, concept string, includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error) {
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
