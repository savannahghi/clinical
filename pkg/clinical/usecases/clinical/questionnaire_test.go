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
						VersionID: ID,
						Source:    "",
						Tag: []domain.FHIRCodingInput{
							{
								ID:           &ID,
								Version:      &ID,
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
						VersionID: ID,
						Source:    "",
						Tag: []domain.FHIRCodingInput{
							{
								ID:           &ID,
								Version:      &ID,
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
						VersionID: ID,
						Source:    "",
						Tag: []domain.FHIRCodingInput{
							{
								ID:           &ID,
								Version:      &ID,
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
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
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

func TestUseCasesClinicalImpl_ListQuestionnaires(t *testing.T) {
	type args struct {
		ctx        context.Context
		title      string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list questionnaire",
			args: args{
				ctx:        addTenantIdentifierContext(context.Background()),
				title:      "Cervical Cancer Screening Form",
				pagination: &dto.Pagination{},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to list questionnaire",
			args: args{
				ctx:        addTenantIdentifierContext(context.Background()),
				title:      "Cervical Cancer Screening Form",
				pagination: &dto.Pagination{},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			args: args{
				ctx:        addTenantIdentifierContext(context.Background()),
				title:      "Cervical Cancer Screening Form",
				pagination: &dto.Pagination{},
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
			q := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad case: unable to list questionnaire" {
				fakeFHIR.MockListFHIRQuestionnaireFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRQuestionnaires, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := q.ListQuestionnaires(tt.args.ctx, tt.args.title, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ListQuestionnaires() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
