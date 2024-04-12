package clinical_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
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
	"github.com/savannahghi/scalarutils"
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
				ctx:              addTenantIdentifierContext(ctx),
				serviceRequestID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing service request ID",
			args: args{
				ctx: addTenantIdentifierContext(ctx),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get service request",
			args: args{
				ctx:              addTenantIdentifierContext(ctx),
				serviceRequestID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get patient",
			args: args{
				ctx:              addTenantIdentifierContext(ctx),
				serviceRequestID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to upload media",
			args: args{
				ctx:              addTenantIdentifierContext(ctx),
				serviceRequestID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to create FHIR document reference",
			args: args{
				ctx:              addTenantIdentifierContext(ctx),
				serviceRequestID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get terminology concept",
			args: args{
				ctx:              addTenantIdentifierContext(ctx),
				serviceRequestID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get organization",
			args: args{
				ctx:              addTenantIdentifierContext(ctx),
				serviceRequestID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get tenant meta tags",
			args: args{
				ctx:              addTenantIdentifierContext(ctx),
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
			if tt.name == "Sad Case - unable to upload media" {
				fakeUpload.MockUploadMediaFn = func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.Media, error) {
					return nil, fmt.Errorf("failed to upload media")
				}
			}
			if tt.name == "Sad Case - unable to create FHIR document reference" {
				fakeFHIR.MockCreateFHIRDocumentReferenceFn = func(ctx context.Context, documentReference *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
					return nil, fmt.Errorf("failed to create FHIR document reference")
				}
			}
			if tt.name == "Sad Case - unable to get terminology concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if tt.name == "Sad Case - Fail to get organization" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to get organization")
				}
			}
			if tt.name == "Sad Case - unable to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if _, err := c.GenerateReferralReportPDF(tt.args.ctx, tt.args.serviceRequestID); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GenerateReferralReportPDF() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_CreateDocumentReference(t *testing.T) {
	ref := fmt.Sprintf("ServiceRequest/%s", gofakeit.UUID())
	url := gofakeit.URL()
	mimeType := "application/json"
	title := fmt.Sprintf("%s's Document Reference", gofakeit.Name())
	subjectRef := fmt.Sprintf("Subject/%s", gofakeit.UUID())
	type args struct {
		ctx     context.Context
		payload *clinicalUsecase.DocumentReferencePayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create document reference",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				payload: &clinicalUsecase.DocumentReferencePayload{
					Subject: &domain.FHIRReferenceInput{
						Reference: &subjectRef,
					},
					Attachment: &domain.FHIRAttachment{
						ContentType: (*scalarutils.Code)(&mimeType),
						URL:         (*scalarutils.URL)(&url),
						Title:       &title,
					},
					Related: &domain.FHIRReference{
						Reference: &ref,
					},
					TerminologySystem: common.ReferralNoteLOINCTerminologySystem,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create document reference",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				payload: &clinicalUsecase.DocumentReferencePayload{
					Subject: &domain.FHIRReferenceInput{
						Reference: &subjectRef,
					},
					Attachment: &domain.FHIRAttachment{
						ContentType: (*scalarutils.Code)(&mimeType),
						URL:         (*scalarutils.URL)(&url),
						Title:       &title,
					},
					Related: &domain.FHIRReference{
						Reference: &ref,
					},
					TerminologySystem: common.ReferralNoteLOINCTerminologySystem,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get concept",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				payload: &clinicalUsecase.DocumentReferencePayload{
					Subject: &domain.FHIRReferenceInput{
						Reference: &subjectRef,
					},
					Attachment: &domain.FHIRAttachment{
						ContentType: (*scalarutils.Code)(&mimeType),
						URL:         (*scalarutils.URL)(&url),
						Title:       &title,
					},
					Related: &domain.FHIRReference{
						Reference: &ref,
					},
					TerminologySystem: common.ReferralNoteLOINCTerminologySystem,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant meta tags",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				payload: &clinicalUsecase.DocumentReferencePayload{
					Subject: &domain.FHIRReferenceInput{
						Reference: &subjectRef,
					},
					Attachment: &domain.FHIRAttachment{
						ContentType: (*scalarutils.Code)(&mimeType),
						URL:         (*scalarutils.URL)(&url),
						Title:       &title,
					},
					Related: &domain.FHIRReference{
						Reference: &ref,
					},
					TerminologySystem: common.ReferralNoteLOINCTerminologySystem,
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

			if tt.name == "Sad case: unable to create document reference" {
				fakeFHIR.MockCreateFHIRDocumentReferenceFn = func(ctx context.Context, documentReference *domain.FHIRDocumentReferenceInput) (*domain.FHIRDocumentReference, error) {
					return nil, fmt.Errorf("failed to create FHIR document reference")
				}
			}
			if tt.name == "Sad case: unable to get concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}
			if tt.name == "Sad case: unable to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if err := c.CreateDocumentReference(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateDocumentReference() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_ShareReferralForm(t *testing.T) {
	type args struct {
		ctx              context.Context
		serviceRequestID string
		workstationID    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: share referral form",
			args: args{
				ctx:              addTenantIdentifierContext(context.Background()),
				serviceRequestID: gofakeit.UUID(),
				workstationID:    gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to share referral form",
			args: args{
				ctx:              addTenantIdentifierContext(context.Background()),
				serviceRequestID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get tenant identifiers",
			args: args{
				ctx:              addTenantIdentifierContext(context.Background()),
				serviceRequestID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: no document references found",
			args: args{
				ctx:              addTenantIdentifierContext(context.Background()),
				serviceRequestID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get patient",
			args: args{
				ctx:              addTenantIdentifierContext(context.Background()),
				serviceRequestID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get send SMS",
			args: args{
				ctx:              addTenantIdentifierContext(context.Background()),
				serviceRequestID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to subject associated with the document reference",
			args: args{
				ctx:              addTenantIdentifierContext(context.Background()),
				serviceRequestID: gofakeit.UUID(),
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

			if tt.name == "Sad case: unable to share referral form" {
				fakeFHIR.MockSearchFHIRDocumentReferenceFn = func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
					return nil, fmt.Errorf("failed to search FHIR document reference")
				}
			}
			if tt.name == "Sad case: unable to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}
			if tt.name == "Sad case: no document references found" {
				fakeFHIR.MockSearchFHIRDocumentReferenceFn = func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
					return &domain.PagedFHIRDocumentReference{}, nil
				}
			}
			if tt.name == "Sad case: unable to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get patient")
				}
			}
			if tt.name == "Sad case: unable to get send SMS" {
				fakeAdvantage.MockSendSMSFn = func(ctx context.Context, workstationID string, payload dto.SMSPayload) error {
					return fmt.Errorf("failed to send SMS")
				}
			}
			if tt.name == "Sad case: unable to subject associated with the document reference" {
				fakeFHIR.MockSearchFHIRDocumentReferenceFn = func(ctx context.Context, searchParams map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRDocumentReference, error) {
					resourceID := uuid.NewString()
					return &domain.PagedFHIRDocumentReference{
						DocumentReferences: []domain.FHIRDocumentReference{
							{
								ID:       resourceID,
								Meta:     &domain.FHIRMeta{},
								Type:     &domain.FHIRCodeableConcept{},
								Category: []domain.FHIRCodeableConcept{},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			_, err := c.ShareReferralForm(tt.args.ctx, tt.args.serviceRequestID, tt.args.workstationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ShareReferralForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
