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

func setupMockFHIRFunctions(fakeFHIR *fakeFHIRMock.FHIRMock, score int) {
	ID := gofakeit.UUID()
	fakeFHIR.MockGetFHIRQuestionnaireFn = func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
		questionnaireName := "Cervical Cancer Screening"
		return &domain.FHIRQuestionnaireRelayPayload{
			Resource: &domain.FHIRQuestionnaire{
				ID:   &ID,
				Name: &questionnaireName,
			},
		}, nil
	}

	fakeFHIR.MockCreateFHIRQuestionnaireResponseFn = func(ctx context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
		return &domain.FHIRQuestionnaireResponse{
			ID: &ID,
			Item: []domain.FHIRQuestionnaireResponseItem{
				{
					LinkID: "symptoms",
					Item: []domain.FHIRQuestionnaireResponseItem{
						{
							LinkID: "symptoms-score",
							Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
								{
									ValueInteger: &score,
								},
							},
						},
					},
				},
				{
					LinkID: "risk-factors",
					Item: []domain.FHIRQuestionnaireResponseItem{
						{
							LinkID: "risk-factors-score",
							Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
								{
									ValueInteger: &score,
								},
							},
						},
					},
				},
			},
		}, nil
	}

	fakeFHIR.MockCreateFHIRRiskAssessmentFn = func(ctx context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error) {
		return nil, fmt.Errorf("failed to record fhir risk assessment")
	}
}

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
		{
			name: "Sad Case - Attempt to record questionnaire response in a finished encounter",
			args: args{
				ctx:             context.Background(),
				encounterID:     gofakeit.UUID(),
				questionnaireID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get fhir questionnaire",
			args: args{
				ctx:             context.Background(),
				encounterID:     gofakeit.UUID(),
				questionnaireID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Happy Case - Create questionnaire response and generate review summary - Cervical Cancer - High Risk",
			args: args{
				ctx:             context.Background(),
				encounterID:     gofakeit.UUID(),
				questionnaireID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Create questionnaire response and generate review summary - Cervical Cancer - Low Risk",
			args: args{
				ctx:             context.Background(),
				encounterID:     gofakeit.UUID(),
				questionnaireID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to record risk assessment - Low Risk",
			args: args{
				ctx:             context.Background(),
				encounterID:     gofakeit.UUID(),
				questionnaireID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to record risk assessment - High Risk",
			args: args{
				ctx:             context.Background(),
				encounterID:     gofakeit.UUID(),
				questionnaireID: gofakeit.UUID(),
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

			if tt.name == "Sad Case - Attempt to record questionnaire response in a finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: "finished",
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get fhir questionnaire" {
				fakeFHIR.MockGetFHIRQuestionnaireFn = func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
					return nil, fmt.Errorf("failed to get questionnaire")
				}
			}

			if tt.name == "Happy Case - Create questionnaire response and generate review summary - Cervical Cancer - High Risk" {
				fakeFHIR.MockGetFHIRQuestionnaireFn = func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
					questionnaireName := "Cervical Cancer Screening"
					return &domain.FHIRQuestionnaireRelayPayload{
						Resource: &domain.FHIRQuestionnaire{
							ID:   &ID,
							Name: &questionnaireName,
						},
					}, nil
				}

				score := 3
				fakeFHIR.MockCreateFHIRQuestionnaireResponseFn = func(ctx context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
					return &domain.FHIRQuestionnaireResponse{
						ID: &ID,
						Item: []domain.FHIRQuestionnaireResponseItem{
							{
								LinkID: "symptoms",
								Item: []domain.FHIRQuestionnaireResponseItem{
									{
										LinkID: "symptoms-score",
										Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
											{
												ValueInteger: &score,
											},
										},
									},
								},
							},
							{
								LinkID: "risk-factors",
								Item: []domain.FHIRQuestionnaireResponseItem{
									{
										LinkID: "risk-factors-score",
										Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
											{
												ValueInteger: &score,
											},
										},
									},
								},
							},
						},
					}, nil
				}

			}

			if tt.name == "Happy Case - Create questionnaire response and generate review summary - Cervical Cancer - Low Risk" {
				fakeFHIR.MockGetFHIRQuestionnaireFn = func(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error) {
					questionnaireName := "Cervical Cancer Screening"
					return &domain.FHIRQuestionnaireRelayPayload{
						Resource: &domain.FHIRQuestionnaire{
							ID:   &ID,
							Name: &questionnaireName,
						},
					}, nil
				}

				score := 0
				fakeFHIR.MockCreateFHIRQuestionnaireResponseFn = func(ctx context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error) {
					return &domain.FHIRQuestionnaireResponse{
						ID: &ID,
						Item: []domain.FHIRQuestionnaireResponseItem{
							{
								LinkID: "symptoms",
								Item: []domain.FHIRQuestionnaireResponseItem{
									{
										LinkID: "symptoms-score",
										Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
											{
												ValueInteger: &score,
											},
										},
									},
								},
							},
							{
								LinkID: "risk-factors",
								Item: []domain.FHIRQuestionnaireResponseItem{
									{
										LinkID: "risk-factors-score",
										Answer: []domain.FHIRQuestionnaireResponseItemAnswer{
											{
												ValueInteger: &score,
											},
										},
									},
								},
							},
						},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to record risk assessment - Low Risk" {
				setupMockFHIRFunctions(fakeFHIR, 0)
			}

			if tt.name == "Sad Case - Fail to record risk assessment - High Risk" {
				setupMockFHIRFunctions(fakeFHIR, 3)
			}

			_, err := q.CreateQuestionnaireResponse(tt.args.ctx, tt.args.questionnaireID, tt.args.encounterID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateQuestionnaireResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
