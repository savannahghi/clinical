package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeAdvantageMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/advantage/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/firebasetools"
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
			name: "sad case: create an episode of care already exists",
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
			name: "sad case: failed to search episode of care",
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
			name: "sad Case - Fail to get tenant meta tags",
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "happy case: create an episode of care" {
				fakeFHIR.MockSearchFHIREpisodeOfCareFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
					return &domain.FHIREpisodeOfCareRelayConnection{
						Edges: []*domain.FHIREpisodeOfCareRelayEdge{},
					}, nil
				}
			}

			if tt.name == "sad case: create an episode of care already exists" {
				fakeFHIR.MockSearchFHIREpisodeOfCareFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
					PatientRef := "Patient/1"
					OrgRef := "Organization/1"
					return &domain.FHIREpisodeOfCareRelayConnection{
						Edges: []*domain.FHIREpisodeOfCareRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIREpisodeOfCare{
									ID:            new(string),
									Text:          &domain.FHIRNarrative{},
									Identifier:    []*domain.FHIRIdentifier{},
									StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
									Type:          []*domain.FHIRCodeableConcept{},
									Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
									Patient:       &domain.FHIRReference{Reference: &PatientRef},
									ManagingOrganization: &domain.FHIRReference{
										Reference: &OrgRef,
									},
									Period:          &domain.FHIRPeriod{},
									ReferralRequest: []*domain.FHIRReference{},
									CareManager:     &domain.FHIRReference{},
									Team:            []*domain.FHIRReference{},
									Account:         []*domain.FHIRReference{},
									Meta:            &domain.FHIRMeta{},
									Extension:       []*domain.FHIRExtension{},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
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

			if tt.name == "sad case: failed to search episode of care" {
				fakeFHIR.MockSearchFHIREpisodeOfCareFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error) {
					return nil, fmt.Errorf("failed to search episode of care")
				}
			}

			if tt.name == "sad case: failed to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "sad Case - Fail to get tenant meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "sad case: failed to create episode of care" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
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

func TestUseCasesClinicalImpl_PatchEpisodeOfCare(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		id    string
		input dto.EpisodeOfCareInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Patch an episode of care",
			args: args{
				ctx: ctx,
				id:  gofakeit.UUID(),
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumCancelled,
					PatientID: gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid episode of care ID",
			args: args{
				ctx: ctx,
				id:  "123",
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumEnteredInError,
					PatientID: gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get episode of care",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumCancelled,
					PatientID: gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Failed to patch episode of care",
			args: args{
				ctx: ctx,
				id:  gofakeit.UUID(),
				input: dto.EpisodeOfCareInput{
					Status:    dto.EpisodeOfCareStatusEnumPlanned,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Unable to get episode of care" {
				fakeFHIR.MockGetFHIREpisodeOfCareFn = func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
					return nil, fmt.Errorf("failed to get episode of care")
				}
			}

			if tt.name == "Sad Case - Failed to patch episode of care" {
				fakeFHIR.MockPatchFHIREpisodeOfCareFn = func(ctx context.Context, id string, input domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.PatchEpisodeOfCare(ctx, tt.args.id, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
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
