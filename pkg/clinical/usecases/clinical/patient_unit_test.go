package clinical_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

func TestClinicalUseCaseImpl_GetMedicalData(t *testing.T) {

	type args struct {
		ctx       context.Context
		patientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.MedicalData
		wantErr bool
	}{
		{
			name: "Happy case: patient timeline",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get tenant identifiers from context",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},

		{
			name: "Happy Case - Successfully search medication statement",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil node",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil node id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil status",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil coding",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - empty coding",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil subject id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully search allergy intolerance",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},

		{
			name: "Sad Case - Fail to search allergy intolerance",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil node",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil node id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil patient",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil patient id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil encounter",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil encounter id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil code",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - nil coding",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance - empty coding",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},

		{
			name: "Happy Case - Successfully search observation",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil coding",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - empty coding",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil status",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil subject",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil subject id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil encounter",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search observation - nil encounter id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search weight",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search BMI",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search viralLoad",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search cd4Count",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
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
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Happy Case - Successfully search medication statement" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					status := dto.MedicationStatementStatusEnumActive
					code := "123"
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRMedicationStatement{
									ID:     new(string),
									Status: (*domain.MedicationStatementStatusEnum)(&status),
									MedicationCodeableConcept: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
										}},
									},
									Subject: &domain.FHIRReference{
										ID: new(string),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}
			if tt.name == "Sad case: failed to get tenant identifiers from context" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement - nil node" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement - nil node id" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					status := dto.MedicationStatementStatusEnumActive
					code := "123"
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRMedicationStatement{
									Status: (*domain.MedicationStatementStatusEnum)(&status),
									MedicationCodeableConcept: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
										}},
									},
									Subject: &domain.FHIRReference{
										ID: new(string),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}
			if tt.name == "Sad Case - Fail to search medication statement - nil status" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					code := "123"
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRMedicationStatement{
									ID: new(string),
									MedicationCodeableConcept: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
										}},
									},
									Subject: &domain.FHIRReference{
										ID: new(string),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}
			if tt.name == "Sad Case - Fail to search medication statement - nil coding" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					status := dto.MedicationStatementStatusEnumActive
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRMedicationStatement{
									ID:     new(string),
									Status: (*domain.MedicationStatementStatusEnum)(&status),
									MedicationCodeableConcept: &domain.FHIRCodeableConcept{
										ID: new(string),
									},
									Subject: &domain.FHIRReference{
										ID: new(string),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}
			if tt.name == "Sad Case - Fail to search medication statement - empty coding" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					status := dto.MedicationStatementStatusEnumActive
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRMedicationStatement{
									ID:     new(string),
									Status: (*domain.MedicationStatementStatusEnum)(&status),
									MedicationCodeableConcept: &domain.FHIRCodeableConcept{
										ID:     new(string),
										Coding: []*domain.FHIRCoding{},
									},
									Subject: &domain.FHIRReference{
										ID: new(string),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}
			if tt.name == "Sad Case - Fail to search medication statement - nil subject id" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					status := dto.MedicationStatementStatusEnumActive
					code := "123"
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRMedicationStatement{
									ID:     new(string),
									Status: (*domain.MedicationStatementStatusEnum)(&status),
									MedicationCodeableConcept: &domain.FHIRCodeableConcept{
										ID: new(string),
										Coding: []*domain.FHIRCoding{{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
										}},
									},
									Subject: &domain.FHIRReference{},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					return nil, fmt.Errorf("failed to search medication statement")
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					return nil, fmt.Errorf("failed to search allergy intolerance")
				}
			}

			if tt.name == "Happy Case - Successfully search allergy intolerance" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					code := "123"
					system := gofakeit.URL()
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								ID: new(string),
								Code: &domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
											System:  (*scalarutils.URI)(&system),
										},
									},
								},
								Patient: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: "2000-01-01",
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance - nil node" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								ID: new(string),
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance - nil node id" {
				code := "123"
				system := gofakeit.URL()
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								Code: &domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
											System:  (*scalarutils.URI)(&system),
										},
									},
								},
								Patient: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: "2000-01-01",
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance - nil patient" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					code := "123"
					system := gofakeit.URL()
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								Code: &domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
											System:  (*scalarutils.URI)(&system),
										},
									},
								},
								Patient: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: "2000-01-01",
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance - nil patient id" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					code := "123"
					system := gofakeit.URL()
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								Code: &domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
											System:  (*scalarutils.URI)(&system),
										},
									},
								},
								Patient: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: "2000-01-01",
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance - nil encounter" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					code := "123"
					system := gofakeit.URL()
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								Code: &domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
											System:  (*scalarutils.URI)(&system),
										},
									},
								},
								Patient: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: "2000-01-01",
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance - nil encounter id" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					code := "123"
					system := gofakeit.URL()
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								Code: &domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
											System:  (*scalarutils.URI)(&system),
										},
									},
								},
								Patient: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: "2000-01-01",
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance - nil code" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					code := "123"
					system := gofakeit.URL()
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								Code: &domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
											System:  (*scalarutils.URI)(&system),
										},
									},
								},
								Patient: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: "2000-01-01",
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance - nil coding" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					code := "123"
					system := gofakeit.URL()
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								Code: &domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
											System:  (*scalarutils.URI)(&system),
										},
									},
								},
								Patient: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: "2000-01-01",
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}
			if tt.name == "Sad Case - Fail to search allergy intolerance - empty coding" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					code := "123"
					system := gofakeit.URL()
					return &domain.PagedFHIRAllergy{
						Allergies: []domain.FHIRAllergyIntolerance{
							{
								Code: &domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{
										{
											Code:    scalarutils.Code(code),
											Display: gofakeit.BS(),
											System:  (*scalarutils.URI)(&system),
										},
									},
								},
								Patient: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								OnsetPeriod: &domain.FHIRPeriod{
									Start: "2000-01-01",
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Happy Case - Successfully search observation" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Subject: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil coding" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
								},
								Subject: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - empty coding" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID:     new(string),
									Coding: []*domain.FHIRCoding{},
								},
								Subject: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil status" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID: new(string),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil subject" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil subject id" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Subject: &domain.FHIRReference{},
								Encounter: &domain.FHIRReference{
									ID: new(string),
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil encounter" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Subject: &domain.FHIRReference{
									ID: new(string),
								},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search observation - nil status" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					status := dto.ObservationStatusFinal
					valueConcept := "222"
					UUID := gofakeit.UUID()
					return &domain.PagedFHIRObservations{
						Observations: []domain.FHIRObservation{
							{
								ID:     new(string),
								Status: (*domain.ObservationStatusEnum)(&status),
								Code: domain.FHIRCodeableConcept{
									ID: new(string),
									Coding: []*domain.FHIRCoding{{
										Display: gofakeit.BS(),
									}},
								},
								Subject: &domain.FHIRReference{
									ID: new(string),
								},
								Encounter: &domain.FHIRReference{},
								ValueQuantity: &domain.FHIRQuantity{
									Value: 100,
									Unit:  "cm",
								},
								ValueCodeableConcept: (*scalarutils.Code)(&valueConcept),
								ValueString:          new(string),
								ValueBoolean:         new(bool),
								ValueInteger:         new(string),
								ValueRange: &domain.FHIRRange{
									Low: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									High: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueRatio: &domain.FHIRRatio{
									Numerator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
									Denominator: domain.FHIRQuantity{
										Value: 100,
										Unit:  "cm",
									},
								},
								ValueSampledData: &domain.FHIRSampledData{
									ID: &UUID,
								},
								ValueTime: &time.Time{},
								ValueDateTime: &scalarutils.Date{
									Year:  2000,
									Month: 1,
									Day:   1,
								},
								ValuePeriod: &domain.FHIRPeriod{
									Start: scalarutils.DateTime(time.Wednesday.String()),
									End:   scalarutils.DateTime(time.Thursday.String()),
								},
							},
						},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}
			if tt.name == "Sad Case - Fail to search weight" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					if params["code"] == common.WeightCIELTerminologyCode {
						return nil, fmt.Errorf("failed to search observation")
					}

					return &domain.PagedFHIRObservations{
						Observations:    []domain.FHIRObservation{},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search BMI" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					if params["code"] == common.BMICIELTerminologyCode {
						return nil, fmt.Errorf("failed to search observation")
					}

					return &domain.PagedFHIRObservations{
						Observations:    []domain.FHIRObservation{},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search viralLoad" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					if params["code"] == common.ViralLoadCIELTerminologyCode {
						return nil, fmt.Errorf("failed to search observation")
					}

					return &domain.PagedFHIRObservations{
						Observations:    []domain.FHIRObservation{},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search cd4Count" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error) {
					if params["code"] == common.CD4CountCIELTerminologyCode {
						return nil, fmt.Errorf("failed to search observation")
					}

					return &domain.PagedFHIRObservations{
						Observations:    []domain.FHIRObservation{},
						HasNextPage:     false,
						NextCursor:      "",
						HasPreviousPage: false,
						PreviousCursor:  "",
						TotalCount:      0,
					}, nil
				}
			}

			got, err := u.GetMedicalData(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.GetMedicalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient medical data to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient medical data not to be nil for %v", tt.name)
				return
			}
		})
	}

}

func TestUseCasesClinicalImpl_CreatePatient(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.PatientInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: register a patient",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid phone number",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "070000",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get tenant tags",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to create patient",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: no facility id in context",
			args: args{
				ctx: context.Background(),
				input: dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to find facility",
			args: args{
				ctx: addTenantIdentifierContext(context.Background()),
				input: dto.PatientInput{
					FirstName: gofakeit.Name(),
					LastName:  gofakeit.Name(),
					BirthDate: &scalarutils.Date{
						Year:  1997,
						Month: 12,
						Day:   12,
					},
					Gender: dto.GenderFemale,
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "0700000000",
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
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "sad case: fail to get tenant tags" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to get tenant identifiers")
				}
			}

			if tt.name == "sad case: fail to create patient" {
				fakeFHIR.MockCreateFHIRPatientFn = func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
					return nil, fmt.Errorf("failed to create patient")
				}
			}

			if tt.name == "sad case: fail to find facility" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to find facility")
				}
			}

			got, err := c.CreatePatient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patients not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_PatchPatient(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		id    string
		input dto.PatientInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully patch a patient (single field)",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
				input: dto.PatientInput{
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 6,
						Day:   14,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully patch a patient (multiple fields)",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
				input: dto.PatientInput{
					FirstName: gofakeit.Name(),
					Identifiers: []dto.IdentifierInput{
						{
							Type:  dto.IdentifierTypeNationalID,
							Value: "12345678",
						},
					},
					Gender: dto.GenderFemale,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing patient ID",
			args: args{
				ctx: ctx,
				input: dto.PatientInput{
					Gender: dto.GenderFemale,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid phone number",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
				input: dto.PatientInput{
					Contacts: []dto.ContactInput{
						{
							Type:  dto.ContactTypePhoneNumber,
							Value: "070000",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to patch patient",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
				input: dto.PatientInput{
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 6,
						Day:   14,
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
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to patch patient" {
				fakeFHIR.MockPatchFHIRPatientFn = func(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error) {
					return nil, fmt.Errorf("failed to patch patient")
				}
			}

			got, err := u.PatchPatient(tt.args.ctx, tt.args.id, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("Expected patient to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("Expected patient not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_DeletePatient(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete patient",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Missing patient ID",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to delete patient",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			want:    false,
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
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad Case - Fail to delete patient" {
				fakeFHIR.MockDeleteFHIRPatientFn = func(ctx context.Context, id string) (bool, error) {
					return false, fmt.Errorf("failed to delete patient")
				}
			}

			got, err := u.DeletePatient(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.DeletePatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesClinicalImpl.DeletePatient() = %v, want %v", got, tt.want)
			}
		})
	}
}
