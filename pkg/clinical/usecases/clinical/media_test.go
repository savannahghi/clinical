package clinical_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeMyCarehubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

func TestUseCasesClinicalImpl_UploadMedia(t *testing.T) {
	type args struct {
		ctx         context.Context
		encounterID string
		file        io.Reader
		contentType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: upload media",
			args: args{
				ctx:         addTenantIdentifierContext(context.Background()),
				encounterID: uuid.NewString(),
				file:        strings.NewReader("test"),
				contentType: "application/json",
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to get encounter",
			args: args{
				ctx:         addTenantIdentifierContext(context.Background()),
				encounterID: uuid.NewString(),
				file:        strings.NewReader("test"),
				contentType: "application/json",
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to upload media",
			args: args{
				ctx:         addTenantIdentifierContext(context.Background()),
				encounterID: uuid.NewString(),
				file:        strings.NewReader("test"),
				contentType: "application/json",
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to create fhir media",
			args: args{
				ctx:         addTenantIdentifierContext(context.Background()),
				encounterID: uuid.NewString(),
				file:        strings.NewReader("test"),
				contentType: "application/json",
			},
			wantErr: true,
		},
		{
			name: "sad case: missing facility identifier in context",
			args: args{
				ctx:         context.Background(),
				encounterID: uuid.NewString(),
				file:        strings.NewReader("test"),
				contentType: "application/json",
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get fhir organisation",
			args: args{
				ctx:         addTenantIdentifierContext(context.Background()),
				encounterID: uuid.NewString(),
				file:        strings.NewReader("test"),
				contentType: "application/json",
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get patient",
			args: args{
				ctx:         addTenantIdentifierContext(context.Background()),
				encounterID: uuid.NewString(),
				file:        strings.NewReader("test"),
				contentType: "application/json",
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

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH, fakeUpload)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "sad case: unable to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to upload media" {
				fakeUpload.MockUploadMediaFn = func(ctx context.Context, name string, file io.Reader, contentType string) (*dto.MediaOutPut, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to create fhir media" {
				fakeFHIR.MockCreateFHIRMediaFn = func(ctx context.Context, input domain.FHIRMedia) (*domain.FHIRMedia, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get fhir organisation" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := c.UploadMedia(tt.args.ctx, tt.args.encounterID, tt.args.file, tt.args.contentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.UploadMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
