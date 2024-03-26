package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
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

func TestUseCasesClinicalImpl_GenerateReferralReportPDF(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx              context.Context
		serviceRequestID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully generate a referral report pdf",
			args: args{
				ctx:              ctx,
				serviceRequestID: uuid.New().String(),
			},
			// TODO: Fix this @salaton
			wantErr: true,
		},
		{
			name: "Sad Case - Missing service request ID",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get service request",
			args: args{
				ctx:              ctx,
				serviceRequestID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get patient",
			args: args{
				ctx:              ctx,
				serviceRequestID: uuid.New().String(),
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

			if tt.name == "Sad Case - Fail to get service request" {
				fakeFHIR.MockGetFHIRServiceRequestFn = func(_ context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error) {
					return nil, fmt.Errorf("failed to get service request ")
				}
			}

			if tt.name == "Sad Case - Fail to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get patient")
				}
			}

			if err := c.GenerateReferralReportPDF(tt.args.ctx, tt.args.serviceRequestID); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GenerateReferralReportPDF() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
