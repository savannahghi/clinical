package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
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

func TestUseCasesClinicalImpl_CreateComposition(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.CompositionInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create composition",
			args: args{
				ctx: context.Background(),
				input: dto.CompositionInput{
					EncounterID: gofakeit.UUID(),
					Type:        dto.ProgressNote,
					Category:    dto.AssessmentAndPlan,
					Status:      "final",
					Note:        "Patient is deteriorating",
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error fetching concept",
			args: args{
				ctx: context.Background(),
				input: dto.CompositionInput{
					EncounterID: gofakeit.UUID(),
					Type:        dto.ProgressNote,
					Category:    dto.AssessmentAndPlan,
					Status:      "final",
					Note:        "Patient is deteriorating",
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get encounter",
			args: args{
				ctx: context.Background(),
				input: dto.CompositionInput{
					EncounterID: gofakeit.UUID(),
					Type:        dto.ProgressNote,
					Category:    dto.AssessmentAndPlan,
					Status:      "final",
					Note:        "Patient is deteriorating",
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get patient",
			args: args{
				ctx: context.Background(),
				input: dto.CompositionInput{
					EncounterID: gofakeit.UUID(),
					Type:        dto.ProgressNote,
					Category:    dto.AssessmentAndPlan,
					Status:      "final",
					Note:        "Patient is deteriorating",
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get identifiers",
			args: args{
				ctx: context.Background(),
				input: dto.CompositionInput{
					EncounterID: gofakeit.UUID(),
					Type:        dto.ProgressNote,
					Category:    dto.AssessmentAndPlan,
					Status:      "final",
					Note:        "Patient is deteriorating",
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail on finished encounter",
			args: args{
				ctx: context.Background(),
				input: dto.CompositionInput{
					EncounterID: gofakeit.UUID(),
					Type:        dto.ProgressNote,
					Category:    dto.AssessmentAndPlan,
					Status:      "final",
					Note:        "Patient is deteriorating",
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

			if tt.name == "happy case: get encounter" {
				UUID := uuid.New().String()
				compositionTitle := gofakeit.Name() + "assessment note"
				typeSystem := scalarutils.URI("http://hl7.org/fhir/ValueSet/doc-typecodes")
				categorySystem := scalarutils.URI("http://hl7.org/fhir/ValueSet/referenced-item-category")
				category := "Assessment + plan"
				compositionType := "Progress note"
				treatmentPlan := "Treatment Plan"
				compositionStatus := "active"
				note := scalarutils.Markdown("Fever Fever")
				PatientRef := "Patient/" + uuid.NewString()
				patientType := "Patient"
				organizationRef := "Organization/" + uuid.NewString()
				compositionSectionTextStatus := "generated"

				fakeFHIR.MockCreateFHIRCompositionFn = func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
					return &domain.FHIRCompositionRelayPayload{
						Resource: &domain.FHIRComposition{
							ID:         &UUID,
							Text:       &domain.FHIRNarrative{},
							Identifier: &domain.FHIRIdentifier{},
							Status:     (*domain.CompositionStatusEnum)(&compositionStatus),
							Type: &domain.FHIRCodeableConcept{
								ID: new(string),
								Coding: []*domain.FHIRCoding{
									{
										ID:      &UUID,
										System:  &typeSystem,
										Code:    scalarutils.Code(string(common.LOINCProgressNoteCode)),
										Display: compositionType,
									},
								},
								Text: "Progress note",
							},
							Category: []*domain.FHIRCodeableConcept{
								{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											ID:      &UUID,
											System:  &categorySystem,
											Version: new(string),
											Code:    scalarutils.Code(string(common.LOINCAssessmentPlanCode)),
											Display: category,
										},
									},
									Text: "Assessment + plan",
								},
							},
							Subject: &domain.FHIRReference{
								ID:        &UUID,
								Reference: &PatientRef,
								Type:      (*scalarutils.URI)(&patientType),
							},
							Encounter: &domain.FHIRReference{
								ID: &UUID,
							},
							Date: &scalarutils.Date{
								Year:  2023,
								Month: 9,
								Day:   25,
							},
							Author: []*domain.FHIRReference{
								{
									Reference: &organizationRef,
								},
							},
							Title: &compositionTitle,
							Section: []*domain.FHIRCompositionSection{
								{
									ID:    &UUID,
									Title: &treatmentPlan,
									Code: &domain.FHIRCodeableConceptInput{
										ID: new(string),
										Coding: []*domain.FHIRCodingInput{
											{
												ID:      &UUID,
												System:  &categorySystem,
												Version: new(string),
												Code:    scalarutils.Code(string(common.LOINCAssessmentPlanCode)),
												Display: category,
											},
										},
										Text: "Assessment + plan",
									},
									Author: []*domain.FHIRReference{
										{
											Reference: new(string),
										},
									},
									Text: &domain.FHIRNarrative{
										ID:     &UUID,
										Status: (*domain.NarrativeStatusEnum)(&compositionSectionTextStatus),
										Div:    scalarutils.XHTML(note),
									},
								},
							},
						},
					}, nil
				}
			}

			if tt.name == "sad case: error fetching concept" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org string, source string, concept string, includeMappings bool, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if tt.name == "sad case: fail on finished encounter" {
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

			if tt.name == "sad case: fail to get encounter" {
				fakeFHIR.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get concept")
				}
			}

			if tt.name == "sad case: fail to get patient" {
				fakeFHIR.MockGetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
					return nil, fmt.Errorf("failed to get patient")
				}
			}

			if tt.name == "sad case: fail to get identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get identifiers")
				}
			}

			_, err := c.CreateComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CreateComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}