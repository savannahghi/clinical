package clinical_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestUseCasesClinicalImpl_CreateQuestionnaire(t *testing.T) {
	ID := gofakeit.UUID()
	type args struct {
		ctx                context.Context
		questionnaireInput *domain.FHIRQuestionnaire
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create questionnaire",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				questionnaireInput: &domain.FHIRQuestionnaire{
					ID: &ID,
					Meta: &domain.FHIRMetaInput{
						VersionID:   ID,
						LastUpdated: time.Now(),
						Source:      "",
						Tag: []domain.FHIRCodingInput{
							{
								ID:           &ID,
								Version:      &ID,
								Code:         "",
								Display:      "",
								UserSelected: new(bool),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get tenant tags",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				questionnaireInput: &domain.FHIRQuestionnaire{
					ID: &ID,
					Meta: &domain.FHIRMetaInput{
						VersionID:   ID,
						LastUpdated: time.Now(),
						Source:      "",
						Tag: []domain.FHIRCodingInput{
							{
								ID:           &ID,
								Version:      &ID,
								Code:         "",
								Display:      "",
								UserSelected: new(bool),
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create questionnaire",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				questionnaireInput: &domain.FHIRQuestionnaire{
					ID: &ID,
					Meta: &domain.FHIRMetaInput{
						VersionID:   ID,
						LastUpdated: time.Now(),
						Source:      "",
						Tag: []domain.FHIRCodingInput{
							{
								ID:           &ID,
								Version:      &ID,
								Code:         "",
								Display:      "",
								UserSelected: new(bool),
							},
						},
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			q := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad case: unable to get tenant tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to create questionnaire" {
				fakeFHIR.MockCreateFHIRQuestionnaireFn = func(ctx context.Context, input *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := q.CreateQuestionnaire(tt.args.ctx, tt.args.questionnaireInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateQuestionnaire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}