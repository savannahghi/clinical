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

func TestUseCasesClinicalImpl_RecordConsent(t *testing.T) {
	ID := gofakeit.UUID()
	status := dto.ConsentStatusActive
	provisionType := dto.ConsentProvisionTypePermit

	type args struct {
		ctx   context.Context
		input dto.ConsentInput
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create a fhir consent",
			args: args{
				ctx: context.Background(),
				input: dto.ConsentInput{
					EncounterID: ID,
					Provision:   dto.ConsentProvisionTypeEnum(provisionType),
					Status:      dto.ConsentStatusEnum(status),
					DenyReason:  "",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to create consent",
			args: args{
				ctx: context.Background(),
				input: dto.ConsentInput{
					EncounterID: ID,
					Status:      dto.ConsentStatusEnum(status),
					DenyReason:  "",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid encounter id",
			args: args{
				ctx: context.Background(),
				input: dto.ConsentInput{
					EncounterID: "",
					Status:      dto.ConsentStatusEnum(status),
					DenyReason:  "",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid context",
			args: args{
				ctx: nil,
				input: dto.ConsentInput{
					EncounterID: "",
					Status:      dto.ConsentStatusEnum(status),
					DenyReason:  "",
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

			if tt.name == "Sad Case: failed to create consent" {
				fakeFHIR.MockCreateFHIRConsentFn = func(ctx context.Context, input domain.FHIRConsent) (*domain.FHIRConsent, error) {
					return nil, fmt.Errorf("an error occurred")
				}

			}
			if tt.name == "Sad Case: invalid encounter id" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}

			}
			if tt.name == "Sad Case: invalid context" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}

			}

			_, err := c.RecordConsent(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordConsent() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}

}
