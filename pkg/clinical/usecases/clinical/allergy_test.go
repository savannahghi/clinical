package clinical_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeMyCarehubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/scalarutils"
)

func TestUseCasesClinicalImpl_CreateAllergyIntolerance(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.AllergyInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create allergy intolerance",
			args: args{
				ctx: context.Background(),
				input: dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "C12345",
					TerminologySource: dto.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: create allergy intolerance, no reaction",
			args: args{
				ctx: context.Background(),
				input: dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: dto.TerminologySourceLOINC,
					EncounterID:       gofakeit.UUID(),
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case: unsupported concept source",
			args: args{
				ctx: context.Background(),
				input: dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: dto.TerminologySource("invalid"),
					EncounterID:       gofakeit.UUID(),
				},
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed to create allergy intolerance",
			args: args{
				ctx: context.Background(),
				input: dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: dto.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
					},
				},
			},
			wantErr: true,
		},

		{
			name: "Sad case: no encounter id passed",
			args: args{
				ctx: context.Background(),
				input: dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: dto.TerminologySourceCIEL,
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
					},
				},
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed to get fhir encounter",
			args: args{
				ctx: context.Background(),
				input: dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: dto.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get fhir patient",
			args: args{
				ctx: context.Background(),
				input: dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: dto.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail to get tags",
			args: args{
				ctx: nil,
				input: dto.AllergyInput{
					PatientID:         gofakeit.UUID(),
					Code:              "100",
					TerminologySource: dto.TerminologySourceCIEL,
					EncounterID:       gofakeit.UUID(),
					Reaction: &dto.ReactionInput{
						Code:     "2000",
						System:   gofakeit.BS(),
						Severity: "fatal",
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

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Happy case: create allergy intolerance" {
				system := gofakeit.URL()
				UUID := gofakeit.UUID()
				mildSeverity := domain.AllergyIntoleranceReactionSeverityEnumMild

				fakeFHIR.MockCreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return &domain.FHIRAllergyIntoleranceRelayPayload{
						Resource: &domain.FHIRAllergyIntolerance{
							ID:          &UUID,
							Criticality: "fatal",
							Code: &domain.FHIRCodeableConcept{
								Coding: []*domain.FHIRCoding{{
									System: (*scalarutils.URI)(&system),
									Code:   "20",
								}},
							},
							OnsetPeriod: &domain.FHIRPeriod{
								Start: scalarutils.DateTime("2000-01-01T00:00:00"),
							},
							Patient: &domain.FHIRReference{
								ID: &UUID,
							},
							Encounter: &domain.FHIRReference{
								ID: &UUID,
							},
							Reaction: []*domain.FHIRAllergyintoleranceReaction{
								{
									Severity: &mildSeverity,
									Manifestation: []*domain.FHIRCodeableConcept{
										{
											Coding: []*domain.FHIRCoding{
												{
													System: (*scalarutils.URI)(&system),
													Code:   scalarutils.Code(gofakeit.BS()),
												},
											},
										},
									},
								},
							},
						},
					}, nil
				}
			}

			if tt.name == "Happy case: create allergy intolerance, no reaction" {
				system := gofakeit.URL()
				UUID := gofakeit.UUID()
				fakeFHIR.MockCreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return &domain.FHIRAllergyIntoleranceRelayPayload{
						Resource: &domain.FHIRAllergyIntolerance{
							ID:          &UUID,
							Criticality: "let",
							Code: &domain.FHIRCodeableConcept{
								Coding: []*domain.FHIRCoding{{
									System: (*scalarutils.URI)(&system),
									Code:   "20",
								}},
							},
							OnsetPeriod: &domain.FHIRPeriod{
								Start: scalarutils.DateTime("2000-01-01T00:00:00"),
							},
							Patient: &domain.FHIRReference{
								ID: &UUID,
							},
							Encounter: &domain.FHIRReference{
								ID: &UUID,
							},
						},
					}, nil
				}
			}

			if tt.name == "Sad case: failed to create allergy intolerance" {
				fakeFHIR.MockCreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get fhir encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to get fhir patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: fail to get tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tags")
				}
			}
			_, err := c.CreateAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_SearchAllergy(t *testing.T) {
	first := 5
	type args struct {
		ctx        context.Context
		name       string
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: search for an allergy",
			args: args{
				ctx:  context.Background(),
				name: "Peanuts",
				pagination: dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search for an allergy",
			args: args{
				ctx:  context.Background(),
				name: "test",
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

			if tt.name == "Sad case: unable to search for an allergy" {
				fakeOCL.MockListConceptsFn = func(ctx context.Context, org, source string, verbose bool, q, sortAsc, sortDesc, conceptClass, dataType, locale *string, includeRetired, includeMappings, includeInverseMappings *bool, paginationInput *dto.Pagination) (*domain.ConceptPage, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := c.SearchAllergy(tt.args.ctx, tt.args.name, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.SearchAllergy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetAllergyIntolerance(t *testing.T) {
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
			name: "Happy case: get allergy intolerance",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid uuid",
			args: args{
				ctx: context.Background(),
				id:  "1",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get allergy intolerance",
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

			if tt.name == "Sad case: unable to get allergy intolerance" {
				fakeFHIR.MockGetFHIRAllergyIntoleranceFn = func(ctx context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("error")
				}
			}

			_, err := c.GetAllergyIntolerance(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_GetPatientAllergies(t *testing.T) {
	testInt := 5
	type args struct {
		ctx        context.Context
		patientID  string
		pagination *dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get allergy patient intolerances",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &testInt,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid uuid",
			args: args{
				ctx:       context.Background(),
				patientID: "1",
				pagination: &dto.Pagination{
					First: &testInt,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get patient allergy intolerances",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &testInt,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get meta tags",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: &dto.Pagination{
					First: &testInt,
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

			if tt.name == "Sad case: unable to get patient allergy intolerances" {
				fakeFHIR.MockSearchPatientAllergyIntoleranceFn = func(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					return nil, errors.New("unable to get patient allergy intolerance")
				}
			}
			if tt.name == "Sad case: unable to get meta tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, errors.New("unable to get meta tags")
				}
			}

			_, err := c.ListPatientAllergies(tt.args.ctx, tt.args.patientID, *tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientAllergies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
