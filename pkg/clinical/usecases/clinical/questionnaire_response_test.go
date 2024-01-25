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
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

func TestUseCasesClinicalImpl_CreateQuestionnaireResponse(t *testing.T) {
	ID := gofakeit.UUID()
	type args struct {
		ctx             context.Context
		input           dto.QuestionnaireResponse
		questionnaireID string
		encounterID     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create questionnaire",
			args: args{
				ctx:             addTenantIdentifierContext(context.Background()),
				input:           dto.QuestionnaireResponse{},
				questionnaireID: ID,
				encounterID:     ID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get tenant tags",
			args: args{
				ctx:             nil,
				input:           dto.QuestionnaireResponse{},
				questionnaireID: ID,
				encounterID:     ID,
			},
			wantErr: true,
		},

		{
			name: "Sad Case: invalid encounter id",
			args: args{
				ctx:             addTenantIdentifierContext(context.Background()),
				input:           dto.QuestionnaireResponse{},
				questionnaireID: ID,
				encounterID:     "",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create questionnaire response",
			args: args{
				ctx:             addTenantIdentifierContext(context.Background()),
				questionnaireID: ID,
				encounterID:     ID,
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
			q := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad case: unable to create questionnaire response" {
				fakeFHIR.MockCreateFHIRQuestionnaireResponseFn = func(ctx context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get tenant tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case: invalid encounter id" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}

			}

			_, err := q.CreateQuestionnaireResponse(tt.args.ctx, tt.args.questionnaireID, tt.args.encounterID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateQuestionnaireResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
