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

			if tt.name == "happy case: create an episode of care" {
				fakeFHIR.MockSearchFHIREpisodeOfCareFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
					return &domain.FHIREpisodeOfCareRelayConnection{
						Edges: []*domain.FHIREpisodeOfCareRelayEdge{},
					}, nil
				}
			}

			if tt.name == "sad case: error fetching facility" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
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

func TestUseCasesClinicalImpl_EndEpisodeOfCare(t *testing.T) {

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: end episode of care",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid episode of care id",
			args: args{
				ctx: context.Background(),
				id:  "",
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving episode of care",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get tenant identifiers",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to search encounters",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to end encounter",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to end episode of care",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
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

			if tt.name == "sad case: error retrieving episode of care" {
				fakeFHIR.MockGetFHIREpisodeOfCareFn = func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
					return nil, fmt.Errorf("failed to get episode of care")
				}
			}

			if tt.name == "sad case: fail to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "sad case: fail to search encounters" {
				fakeFHIR.MockSearchEpisodeEncounterFn = func(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error) {
					return nil, fmt.Errorf("error searching for encounters")
				}
			}

			if tt.name == "sad case: fail to end encounter" {
				fakeFHIR.MockEndEncounterFn = func(ctx context.Context, encounterID string) (bool, error) {
					return false, fmt.Errorf("failed to end encounter")
				}
			}

			if tt.name == "sad case: fail to end episode of care" {
				fakeFHIR.MockEndEpisodeFn = func(ctx context.Context, episodeID string) (bool, error) {
					return false, fmt.Errorf("error ending episode")
				}
			}

			got, err := c.EndEpisodeOfCare(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetEpisodeOfCare(t *testing.T) {

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get episode of care",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid episode of care id",
			args: args{
				ctx: context.Background(),
				id:  "bonoko",
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get episode of care",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
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

			if tt.name == "sad case: fail to get episode of care" {
				fakeFHIR.MockGetFHIREpisodeOfCareFn = func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
					return nil, fmt.Errorf("failed to get episode of care")
				}
			}

			got, err := c.GetEpisodeOfCare(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}
