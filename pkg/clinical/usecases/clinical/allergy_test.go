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
					PatientID:   gofakeit.UUID(),
					Code:        "100",
					System:      gofakeit.BS(),
					EncounterID: gofakeit.UUID(),
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
					PatientID:   gofakeit.UUID(),
					Code:        "100",
					System:      gofakeit.BS(),
					EncounterID: gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to create allergy intolerance",
			args: args{
				ctx: context.Background(),
				input: dto.AllergyInput{
					PatientID:   gofakeit.UUID(),
					Code:        "100",
					System:      gofakeit.BS(),
					EncounterID: gofakeit.UUID(),
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
					PatientID: gofakeit.UUID(),
					Code:      "100",
					System:    gofakeit.BS(),
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
					PatientID:   gofakeit.UUID(),
					Code:        "100",
					System:      gofakeit.BS(),
					EncounterID: gofakeit.UUID(),
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
					PatientID:   gofakeit.UUID(),
					Code:        "100",
					System:      gofakeit.BS(),
					EncounterID: gofakeit.UUID(),
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
			_, err := c.CreateAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
