package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeMyCarehubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

func TestUseCasesClinicalImpl_RegisterTenant(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.OrganizationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully register tenant",
			args: args{
				ctx: ctx,
				input: dto.OrganizationInput{
					Name:        "Test facility",
					PhoneNumber: "0700000000",
					Identifiers: []dto.OrganizationIdentifier{
						{
							Type:  "Other",
							Value: "001",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing identifiers",
			args: args{
				ctx: ctx,
				input: dto.OrganizationInput{
					Name:        "Test facility",
					PhoneNumber: "0700000000",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing name",
			args: args{
				ctx: ctx,
				input: dto.OrganizationInput{
					PhoneNumber: "0700000000",
					Identifiers: []dto.OrganizationIdentifier{
						{
							Type:  "Other",
							Value: "001",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create organisation",
			args: args{
				ctx: ctx,
				input: dto.OrganizationInput{
					Name:        "Test facility",
					PhoneNumber: "0700000000",
					Identifiers: []dto.OrganizationIdentifier{
						{
							Type:  "Other",
							Value: "001",
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
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH, fakeUpload, fakePubSub)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to create organisation" {
				fakeFHIR.MockCreateFHIROrganizationFn = func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create organisation")
				}
			}

			got, err := u.RegisterTenant(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RegisterTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got: %v", got)
					return
				}
			}
		})
	}
}

func TestUseCasesClinicalImpl_RegisterFacility(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.OrganizationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - successfully register facility",
			args: args{
				ctx: ctx,
				input: dto.OrganizationInput{
					Name:        "Test",
					PhoneNumber: "Number",
					Identifiers: []dto.OrganizationIdentifier{
						{
							Type:  "MFLCode",
							Value: "1234",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case - fail to register facility",
			args: args{
				ctx: ctx,
				input: dto.OrganizationInput{
					Name:        "Test",
					PhoneNumber: "Number",
					Identifiers: []dto.OrganizationIdentifier{
						{
							Type:  "MFLCode",
							Value: "1234",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing slade code / mfl code identifier",
			args: args{
				ctx: ctx,
				input: dto.OrganizationInput{
					Name:        "Test",
					PhoneNumber: "Number",
					Identifiers: []dto.OrganizationIdentifier{
						{
							Type:  "Other",
							Value: "1234",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Missing name",
			args: args{
				ctx: ctx,
				input: dto.OrganizationInput{
					PhoneNumber: "Number",
					Identifiers: []dto.OrganizationIdentifier{
						{
							Type:  "SladeCode",
							Value: "1234",
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
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH, fakeUpload, fakePubSub)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad case - fail to register facility" {
				fakeFHIR.MockCreateFHIROrganizationFn = func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to register facility")
				}
			}

			got, err := u.RegisterFacility(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.RegisterFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("expected a response but got: %v", got)
					return
				}
			}
		})
	}
}
