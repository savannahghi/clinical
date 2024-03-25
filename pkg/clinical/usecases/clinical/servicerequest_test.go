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
	fakeAdvantageMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/advantage/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

func TestUseCasesClinicalImpl_ReferPatient(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *dto.ReferralInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully refer patient",
			args: args{
				ctx: context.Background(),
				input: &dto.ReferralInput{
					EncounterID:  gofakeit.UUID(),
					ReferralType: "DIAGNOSTICS",
					Tests:        []string{"VIA"},
					Specialist:   "Oncologist",
					Facility:     "KNH",
					ReferralNote: "",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to get encounter",
			args: args{
				ctx: context.Background(),
				input: &dto.ReferralInput{
					EncounterID:  gofakeit.UUID(),
					ReferralType: "DIAGNOSTICS",
					Tests:        []string{"VIA"},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Record referral in finished encounter",
			args: args{
				ctx: context.Background(),
				input: &dto.ReferralInput{
					EncounterID:  gofakeit.UUID(),
					ReferralType: "DIAGNOSTICS",
					Tests:        []string{"VIA"},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to get tenant meta tags",
			args: args{
				ctx: context.Background(),
				input: &dto.ReferralInput{
					EncounterID:  gofakeit.UUID(),
					ReferralType: "DIAGNOSTICS",
					Tests:        []string{"VIA"},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Fail to create service request",
			args: args{
				ctx: context.Background(),
				input: &dto.ReferralInput{
					EncounterID:  gofakeit.UUID(),
					ReferralType: "DIAGNOSTICS",
					Tests:        []string{"VIA"},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Input validation - missing encounter ID",
			args: args{
				ctx: context.Background(),
				input: &dto.ReferralInput{
					ReferralType: "DIAGNOSTICS",
					Tests:        []string{"VIA"},
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

			if tt.name == "Sad Case: Fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get encounter")
				}
			}

			if tt.name == "Sad Case: Record referral in finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}

			if tt.name == "Sad Case: Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "Sad Case: Fail to create service request" {
				fakeFHIR.MockCreateFHIRServiceRequestFn = func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
					return nil, fmt.Errorf("failed to record service request")
				}
			}

			got, err := c.ReferPatient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ReferPatient() error = %v, wantErr %v", err, tt.wantErr)
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
