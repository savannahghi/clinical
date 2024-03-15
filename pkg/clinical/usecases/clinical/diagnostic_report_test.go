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

func TestUseCasesClinicalImpl_RecordMammographyResult(t *testing.T) {
	url := gofakeit.URL()
	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.DiagnosticReport
		wantErr bool
	}{
		{
			name: "Happy case: record mammography report",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
					Findings:    "BI-RADs 0",
					Media: &dto.Media{
						ID:   gofakeit.UUID(),
						URL:  url,
						Name: gofakeit.BeerName(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get encounter",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
					Findings:    "BI-RADs 0",
					Media: &dto.Media{
						ID:   gofakeit.UUID(),
						URL:  url,
						Name: gofakeit.BeerName(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: finished encounter",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
					Findings:    "BI-RADs 0",
					Media: &dto.Media{
						ID:   gofakeit.UUID(),
						URL:  url,
						Name: gofakeit.BeerName(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get meta tags",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
					Findings:    "BI-RADs 0",
					Media: &dto.Media{
						ID:   gofakeit.UUID(),
						URL:  url,
						Name: gofakeit.BeerName(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get facility ID from context",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
					Findings:    "BI-RADs 0",
					Media: &dto.Media{
						ID:   gofakeit.UUID(),
						URL:  url,
						Name: gofakeit.BeerName(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get FHIR organisation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
					Findings:    "BI-RADs 0",
					Media: &dto.Media{
						ID:   gofakeit.UUID(),
						URL:  url,
						Name: gofakeit.BeerName(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get concept",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
					Findings:    "BI-RADs 0",
					Media: &dto.Media{
						ID:   gofakeit.UUID(),
						URL:  url,
						Name: gofakeit.BeerName(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create FHIR observation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
					Findings:    "BI-RADs 0",
					Media: &dto.Media{
						ID:   gofakeit.UUID(),
						URL:  url,
						Name: gofakeit.BeerName(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create FHIR diagnostic report",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
					Findings:    "BI-RADs 0",
					Media: &dto.Media{
						ID:   gofakeit.UUID(),
						URL:  url,
						Name: gofakeit.BeerName(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: required field omitted",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Test",
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

			if tt.name == "Sad case: unable to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							Status: domain.EncounterStatusEnumFinished,
						},
					}, nil
				}
			}
			if tt.name == "Sad case: unable to get meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get facility ID from context" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get FHIR organisation" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to create FHIR observation" {
				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to create FHIR diagnostic report" {
				fakeFHIR.MockCreateFHIRDiagnosticReportFn = func(_ context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := u.RecordMammographyResult(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordMammographyResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordBiopsy(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: successfully record biopsy test",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Go for biopsy test",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to successfully record biopsy test",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Go for biopsy test",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail validation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record observation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Go for biopsy test",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Go for biopsy test",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get organization",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "Go for biopsy test",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
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

			if tt.name == "Sad case: unable to successfully record biopsy test" {
				fakeFHIR.MockCreateFHIRDiagnosticReportFn = func(_ context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to record observation" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get organization" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := u.RecordBiopsy(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordBiopsy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordMRI(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: successfully record mri results",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No Tumours observed",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to successfully record mri results",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No Tumours observed",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail validation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record observation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No Tumours observed",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No Tumours observed",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get organization",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No Tumours observed",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
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

			if tt.name == "Sad case: unable to successfully record mri results" {
				fakeFHIR.MockCreateFHIRDiagnosticReportFn = func(_ context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to record observation" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get organization" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := u.RecordMRI(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordMRI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordUltrasound(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.DiagnosticReport
		wantErr bool
	}{
		{
			name: "Happy case: record ultrasound",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to record ultrasound results",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail validation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record observation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get organization",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
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

			if tt.name == "Sad case: unable to record ultrasound results" {
				fakeFHIR.MockCreateFHIRDiagnosticReportFn = func(_ context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to record observation" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get organization" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := u.RecordUltrasound(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordUltrasound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_RecordCBE(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *dto.DiagnosticReportInput
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.DiagnosticReport
		wantErr bool
	}{
		{
			name: "Happy case: record CBE test",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: &dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to record CBE test",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: &dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail validation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: &dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to record observation",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: &dto.DiagnosticReportInput{
					Note: "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: &dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get organization",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: &dto.DiagnosticReportInput{
					EncounterID: "12345678905432345",
					Note:        "No lumps felt",
					Media: &dto.Media{
						URL:  gofakeit.URL(),
						Name: gofakeit.Name(),
					},
					Findings: gofakeit.HipsterSentence(20),
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

			if tt.name == "Sad case: unable to record CBE test" {
				fakeFHIR.MockCreateFHIRDiagnosticReportFn = func(_ context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to record observation" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get organization" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := u.RecordCBE(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RecordCBE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
