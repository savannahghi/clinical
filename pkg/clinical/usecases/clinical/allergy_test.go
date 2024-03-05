package clinical_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
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
			name: "Sad case: fail to get encounter",
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
			name: "Sad case - fail on finished encounter",
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
		{
			name: "Sad Case - fail to get ciel concept",
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
			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)
			codingCode := "20"
			manifestationCodingCode := scalarutils.Code(gofakeit.BS())

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
									Code:   (*scalarutils.Code)(&codingCode),
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
													Code:   &manifestationCodingCode,
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
				codingCode := "20"
				fakeFHIR.MockCreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return &domain.FHIRAllergyIntoleranceRelayPayload{
						Resource: &domain.FHIRAllergyIntolerance{
							ID:          &UUID,
							Criticality: "let",
							Code: &domain.FHIRCodeableConcept{
								Coding: []*domain.FHIRCoding{{
									System: (*scalarutils.URI)(&system),
									Code:   (*scalarutils.Code)(&codingCode),
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

			if tt.name == "Sad case: fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if tt.name == "Sad case: failed to create allergy intolerance" {
				fakeFHIR.MockCreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get fhir encounter" {
				UUID := uuid.New().String()
				PatientRef := "Patient/" + uuid.NewString()
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							ID:            &UUID,
							Text:          &domain.FHIRNarrative{},
							Identifier:    []*domain.FHIRIdentifier{},
							Status:        domain.EncounterStatusEnum(domain.EncounterStatusEnumOnleave),
							StatusHistory: []*domain.FHIREncounterStatushistory{},
							Class:         domain.FHIRCoding{},
							ClassHistory:  []*domain.FHIREncounterClasshistory{},
							Type:          []*domain.FHIRCodeableConcept{},
							ServiceType:   &domain.FHIRCodeableConcept{},
							Priority:      &domain.FHIRCodeableConcept{},
							Subject: &domain.FHIRReference{
								ID:        &UUID,
								Reference: &PatientRef,
							},
							EpisodeOfCare:   []*domain.FHIRReference{},
							BasedOn:         []*domain.FHIRReference{},
							Participant:     []*domain.FHIREncounterParticipant{},
							Appointment:     []*domain.FHIRReference{},
							Period:          &domain.FHIRPeriod{},
							Length:          &domain.FHIRDuration{},
							ReasonReference: []*domain.FHIRReference{},
							Diagnosis:       []*domain.FHIREncounterDiagnosis{},
							Account:         []*domain.FHIRReference{},
							Hospitalization: &domain.FHIREncounterHospitalization{},
							Location:        []*domain.FHIREncounterLocation{},
							ServiceProvider: &domain.FHIRReference{},
							PartOf:          &domain.FHIRReference{},
						},
					}, nil
				}
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org string, source string, concept string, includeMappings bool, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if tt.name == "Sad case - fail on finished encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					UUID := uuid.New().String()
					PatientRef := "Patient/" + uuid.NewString()
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							ID:            &UUID,
							Text:          &domain.FHIRNarrative{},
							Identifier:    []*domain.FHIRIdentifier{},
							Status:        domain.EncounterStatusEnum(domain.EncounterStatusEnumFinished),
							StatusHistory: []*domain.FHIREncounterStatushistory{},
							Class:         domain.FHIRCoding{},
							ClassHistory:  []*domain.FHIREncounterClasshistory{},
							Type:          []*domain.FHIRCodeableConcept{},
							ServiceType:   &domain.FHIRCodeableConcept{},
							Priority:      &domain.FHIRCodeableConcept{},
							Subject: &domain.FHIRReference{
								ID:        &UUID,
								Reference: &PatientRef,
							},
							EpisodeOfCare:   []*domain.FHIRReference{},
							BasedOn:         []*domain.FHIRReference{},
							Participant:     []*domain.FHIREncounterParticipant{},
							Appointment:     []*domain.FHIRReference{},
							Period:          &domain.FHIRPeriod{},
							Length:          &domain.FHIRDuration{},
							ReasonReference: []*domain.FHIRReference{},
							Diagnosis:       []*domain.FHIREncounterDiagnosis{},
							Account:         []*domain.FHIRReference{},
							Hospitalization: &domain.FHIREncounterHospitalization{},
							Location:        []*domain.FHIREncounterLocation{},
							ServiceProvider: &domain.FHIRReference{},
							PartOf:          &domain.FHIRReference{},
						},
					}, nil
				}
			}

			if tt.name == "Sad case: failed to get fhir patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - fail to get ciel concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
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
			name: "Sad case: invalid pagination",
			args: args{
				ctx:  context.Background(),
				name: "Peanuts",
				pagination: dto.Pagination{
					First: &first,
					Last:  &first,
				},
			},
			wantErr: true,
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
			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad case: unable to search for an allergy" {
				fakeOCL.MockListConceptsFn = func(ctx context.Context, org []string, source []string, verbose bool, q, sortAsc, sortDesc, conceptClass, dataType, locale *string, includeRetired, includeMappings, includeInverseMappings *bool, paginationInput *dto.Pagination) (*domain.ConceptPage, error) {
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
			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
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
	first := 3
	type args struct {
		ctx        context.Context
		patientID  string
		pagination dto.Pagination
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
				pagination: dto.Pagination{
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
				pagination: dto.Pagination{
					First: &testInt,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid pagination",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: dto.Pagination{
					First: &first,
					Last:  &first,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get patient allergy intolerances",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
				pagination: dto.Pagination{
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
				pagination: dto.Pagination{
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
			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
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

			_, err := c.ListPatientAllergies(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetPatientAllergies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
