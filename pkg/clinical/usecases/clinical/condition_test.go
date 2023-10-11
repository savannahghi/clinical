package clinical_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/scalarutils"
)

func TestUseCasesClinicalImpl_CreateCondition(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.ConditionInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create condition - problem list",
			args: args{
				ctx: context.Background(),
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      dto.TerminologySourceCIEL,
					Status:      dto.ConditionStatusActive,
					Category:    dto.ConditionCategoryProblemList,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: create condition - encounter diagnosis",
			args: args{
				ctx: context.Background(),
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      dto.TerminologySourceCIEL,
					Status:      dto.ConditionStatusActive,
					Category:    dto.ConditionCategoryDiagnosis,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: create condition -  invalid category",
			args: args{
				ctx: context.Background(),
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      dto.TerminologySourceCIEL,
					Status:      dto.ConditionStatusActive,
					Category:    dto.ConditionCategoryProblemList,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error fetching concept",
			args: args{
				ctx: context.Background(),
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					Category:    dto.ConditionCategoryProblemList,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail on finished encounter",
			args: args{
				ctx: context.Background(),
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      dto.TerminologySourceCIEL,
					Status:      dto.ConditionStatusActive,
					Category:    dto.ConditionCategoryDiagnosis,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get encounter",
			args: args{
				ctx: context.Background(),
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					Category:    dto.ConditionCategoryProblemList,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get patient",
			args: args{
				ctx: context.Background(),
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					Category:    dto.ConditionCategoryProblemList,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to create condition",
			args: args{
				ctx: context.Background(),
				input: dto.ConditionInput{
					Code:        "386661006",
					System:      "SNOMED",
					Status:      dto.ConditionStatusActive,
					Category:    dto.ConditionCategoryProblemList,
					EncounterID: gofakeit.UUID(),
					Note:        "Fever Fever",
					OnsetDate: &scalarutils.Date{
						Year:  2022,
						Month: 12,
						Day:   12,
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "happy case: create condition - encounter diagnosis" {
				fakeFHIR.MockCreateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					UUID := uuid.New().String()
					statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
					status := "active"
					note := scalarutils.Markdown("Fever Fever")
					noteTime := time.Now()
					uri := scalarutils.URI("1234567")
					return &domain.FHIRConditionRelayPayload{
						Resource: &domain.FHIRCondition{
							ID:         &UUID,
							Text:       &domain.FHIRNarrative{},
							Identifier: []*domain.FHIRIdentifier{},
							ClinicalStatus: &domain.FHIRCodeableConcept{
								Coding: []*domain.FHIRCoding{
									{
										System:  &statusSystem,
										Code:    scalarutils.Code(string(status)),
										Display: string(status),
									},
								},
								Text: string(status),
							},
							Code: &domain.FHIRCodeableConcept{
								Coding: []*domain.FHIRCoding{
									{
										System:  &uri,
										Code:    scalarutils.Code("1234"),
										Display: "1234",
									},
								},
								Text: "1234",
							},
							OnsetDateTime: &scalarutils.Date{},
							RecordedDate:  &scalarutils.Date{},
							Note: []*domain.FHIRAnnotation{
								{
									Time: &noteTime,
									Text: &note,
								},
							},
							Subject: &domain.FHIRReference{
								ID: &UUID,
							},
							Encounter: &domain.FHIRReference{
								ID: &UUID,
							},
							Category: []*domain.FHIRCodeableConcept{
								{
									ID: &UUID,
									Coding: []*domain.FHIRCoding{
										{
											ID:           &UUID,
											System:       (*scalarutils.URI)(&UUID),
											Version:      &UUID,
											Code:         "ENCOUNTER_DIAGNOSIS",
											Display:      gofakeit.BeerAlcohol(),
											UserSelected: new(bool),
										},
									},
									Text: "ENCOUNTER_DIAGNOSIS",
								},
							},
						},
					}, nil
				}
			}

			if tt.name == "happy case: create condition -  invalid category" {
				fakeFHIR.MockCreateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					UUID := uuid.New().String()
					statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
					status := "active"
					note := scalarutils.Markdown("Fever Fever")
					noteTime := time.Now()
					uri := scalarutils.URI("1234567")
					return &domain.FHIRConditionRelayPayload{
						Resource: &domain.FHIRCondition{
							ID:         &UUID,
							Text:       &domain.FHIRNarrative{},
							Identifier: []*domain.FHIRIdentifier{},
							ClinicalStatus: &domain.FHIRCodeableConcept{
								Coding: []*domain.FHIRCoding{
									{
										System:  &statusSystem,
										Code:    scalarutils.Code(string(status)),
										Display: string(status),
									},
								},
								Text: string(status),
							},
							Code: &domain.FHIRCodeableConcept{
								Coding: []*domain.FHIRCoding{
									{
										System:  &uri,
										Code:    scalarutils.Code("1234"),
										Display: "1234",
									},
								},
								Text: "1234",
							},
							OnsetDateTime: &scalarutils.Date{},
							RecordedDate:  &scalarutils.Date{},
							Note: []*domain.FHIRAnnotation{
								{
									Time: &noteTime,
									Text: &note,
								},
							},
							Subject: &domain.FHIRReference{
								ID: &UUID,
							},
							Encounter: &domain.FHIRReference{
								ID: &UUID,
							},
							Category: []*domain.FHIRCodeableConcept{
								{
									ID: &UUID,
									Coding: []*domain.FHIRCoding{
										{
											ID:           &UUID,
											System:       (*scalarutils.URI)(&UUID),
											Version:      &UUID,
											Code:         "INVALID",
											Display:      gofakeit.BeerAlcohol(),
											UserSelected: new(bool),
										},
									},
									Text: "INVALID",
								},
							},
						},
					}, nil
				}
			}

			if tt.name == "sad case: error fetching concept" {
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

			if tt.name == "Sad Case - fail on finished encounter" {
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

			if tt.name == "sad case: fail to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("fail to get patient")
				}
			}

			if tt.name == "sad case: fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to ge encounter")
				}
			}

			if tt.name == "sad case: fail in completed encounter" {
				finished := domain.EncounterStatusEnumFinished
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					resourceID := uuid.New().String()
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							ID:     &resourceID,
							Status: finished,
						},
					}, nil
				}
			}

			if tt.name == "sad case: failed to create condition" {
				fakeFHIR.MockCreateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := c.CreateCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateCondition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesClinicalImpl_ListPatientConditions(t *testing.T) {
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
			name: "happy case: list conditions",
			args: args{
				ctx:        context.Background(),
				patientID:  gofakeit.UUID(),
				pagination: dto.Pagination{},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid patient id",
			args: args{
				ctx:        context.Background(),
				patientID:  "invalid",
				pagination: dto.Pagination{},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid pagination",
			args: args{
				ctx:       context.Background(),
				patientID: uuid.New().String(),
				pagination: dto.Pagination{
					First: &first,
					Last:  &first,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get identifiers",
			args: args{
				ctx:        context.Background(),
				patientID:  gofakeit.UUID(),
				pagination: dto.Pagination{},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get patient",
			args: args{
				ctx:        context.Background(),
				patientID:  gofakeit.UUID(),
				pagination: dto.Pagination{},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to search condition",
			args: args{
				ctx:        context.Background(),
				patientID:  gofakeit.UUID(),
				pagination: dto.Pagination{},
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

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "sad case: fail to get identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get identifiers")
				}
			}

			if tt.name == "sad case: fail to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get patient")
				}
			}

			if tt.name == "sad case: fail to search condition" {
				fakeFHIR.MockSearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error) {
					return nil, fmt.Errorf("failed to find condition")
				}
			}

			got, err := c.ListPatientConditions(tt.args.ctx, tt.args.patientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPatientConditions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}
