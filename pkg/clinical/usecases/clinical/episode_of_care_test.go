package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeMyCarehubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

func addTenantIdentifierContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, utils.OrganizationIDContextKey, gofakeit.UUID())
	ctx = context.WithValue(ctx, utils.FacilityIDContextKey, gofakeit.UUID())
	return ctx
}

func TestUseCasesClinicalImpl_CreateEpisodeOfCare(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.EpisodeOfCareInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create an episode of care",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumActive,
					PatientID: gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: missing facility identifier in context",
			args: args{
				ctx: context.Background(),
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumActive,
					PatientID: gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error fetching facility",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumActive,
					PatientID: gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error fetching patient",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumActive,
					PatientID: gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get tenant identifiers",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumActive,
					PatientID: gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to create episode of care",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumActive,
					PatientID: gofakeit.UUID(),
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

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "sad case: error fetching facility" {
				fakeFHIR.MockFindOrganizationByIDFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to find organization")
				}
			}

			if tt.name == "sad case: error fetching patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to find patient")
				}
			}

			if tt.name == "sad case: failed to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "sad case: failed to create episode of care" {
				fakeFHIR.MockCreateEpisodeOfCareFn = func(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error) {
					return nil, fmt.Errorf("failed to create episode of care")
				}
			}

			got, err := c.CreateEpisodeOfCare(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}
