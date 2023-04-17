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
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

func TestClinicalUseCaseImpl_PatientTimeline(t *testing.T) {
	type args struct {
		ctx       context.Context
		patientID string
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
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
			name: "Sad case: invalid uuid",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.BS(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get tenant identifiers",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to search allergy intolerance",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get allergy intolerance - nil node",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get allergy intolerance - nil node id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get allergy intolerance - nil code",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get allergy intolerance - nil reaction",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get allergy intolerance - empty reaction",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get allergy intolerance - nil manifestation",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get allergy intolerance - empty manifestation",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get allergy intolerance - nil date",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},

		{
			name: "Sad Case - Fail to search observation",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get observation - nil node",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get observation - nil node id",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get observation - nil coding",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get observation - empty coding",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get observation - nil status",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get observation - nil date",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get observation - invalid date",
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
			wantErr: false,
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
			name: "Sad Case - Fail to search medication statement - nil concept",
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
			name: "Sad Case - Fail to search medication statement - nil status",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil subject",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - nil date",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to search medication statement - invalid date",
			args: args{
				ctx:       context.Background(),
				patientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeMCH := fakeMyCarehubMock.NewFakeMyCareHubServiceMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeMCH)
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Happy case: patient timeline" {
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

				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					status := dto.ObservationStatusFinal
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRObservation{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
										Text: gofakeit.BS(),
									},
									EffectiveDateTime: &scalarutils.Date{
										Year:  2000,
										Month: 1,
										Day:   1,
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}

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
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
									},
									EffectiveDateTime: &scalarutils.Date{Year: 2019, Month: 11, Day: 10},
									Subject: &domain.FHIRReference{
										Display: gofakeit.BS(),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad case: failed to get tenant identifiers" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Fail to search allergy intolerance" {
				fakeFHIR.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error) {
					return &domain.PagedFHIRAllergy{}, fmt.Errorf("failed to search allergy")
				}
			}

			if tt.name == "Sad Case - Fail to get allergy intolerance - nil node" {
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

			if tt.name == "Sad Case - Fail to get allergy intolerance - nil node id" {
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

			if tt.name == "Sad Case - Fail to get allergy intolerance - nil code" {
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

			if tt.name == "Sad Case - Fail to get allergy intolerance - nil reaction" {
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

			if tt.name == "Sad Case - Fail to get allergy intolerance - empty reaction" {
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

			if tt.name == "Sad Case - Fail to get allergy intolerance - nil manifestation" {
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

			if tt.name == "Sad Case - Fail to get allergy intolerance - empty manifestation" {
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

			if tt.name == "Sad Case - Fail to get allergy intolerance - nil date" {
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

			if tt.name == "Sad Case - Fail to search observation" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					return &domain.FHIRObservationRelayConnection{}, fmt.Errorf("failed to get observation")
				}
			}

			if tt.name == "Sad Case - Fail to get observation - nil node" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get observation - nil node id" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					status := dto.ObservationStatusFinal
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRObservation{
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
										Text: gofakeit.BS(),
									},
									EffectiveDateTime: &scalarutils.Date{
										Year:  2000,
										Month: 1,
										Day:   1,
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get observation - nil coding" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					status := dto.ObservationStatusFinal
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRObservation{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: domain.FHIRCodeableConcept{
										Text: gofakeit.BS(),
									},
									EffectiveDateTime: &scalarutils.Date{
										Year:  2000,
										Month: 1,
										Day:   1,
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get observation - empty coding" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					status := dto.ObservationStatusFinal
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRObservation{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{},
										Text:   gofakeit.BS(),
									},
									EffectiveDateTime: &scalarutils.Date{
										Year:  2000,
										Month: 1,
										Day:   1,
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get observation - nil status" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRObservation{
									ID: new(string),
									Code: domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
										Text: gofakeit.BS(),
									},
									EffectiveDateTime: &scalarutils.Date{
										Year:  2000,
										Month: 1,
										Day:   1,
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get observation - nil date" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					status := dto.ObservationStatusFinal
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRObservation{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
										Text: gofakeit.BS(),
									},
									EffectiveDateTime: &scalarutils.Date{
										Year:  20000,
										Month: 1,
										Day:   1,
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get observation - invalid date" {
				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					status := dto.ObservationStatusFinal
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRObservation{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
										Text: gofakeit.BS(),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					return &domain.FHIRMedicationStatementRelayConnection{}, fmt.Errorf("failed to get medication statement")
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
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRMedicationStatement{
									Status: (*domain.MedicationStatementStatusEnum)(&status),
									MedicationCodeableConcept: &domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
									},
									EffectiveDateTime: &scalarutils.Date{Year: 2019, Month: 11, Day: 10},
									Subject: &domain.FHIRReference{
										Display: gofakeit.BS(),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement - nil concept" {
				fakeFHIR.MockSearchFHIRMedicationStatementFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error) {
					status := dto.MedicationStatementStatusEnumActive
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRMedicationStatement{
									ID:                new(string),
									Status:            (*domain.MedicationStatementStatusEnum)(&status),
									EffectiveDateTime: &scalarutils.Date{Year: 2019, Month: 11, Day: 10},
									Subject: &domain.FHIRReference{
										Display: gofakeit.BS(),
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
									ID:                        new(string),
									Status:                    (*domain.MedicationStatementStatusEnum)(&status),
									MedicationCodeableConcept: &domain.FHIRCodeableConcept{},
									EffectiveDateTime:         &scalarutils.Date{Year: 2019, Month: 11, Day: 10},
									Subject: &domain.FHIRReference{
										Display: gofakeit.BS(),
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
										Coding: []*domain.FHIRCoding{},
									},
									EffectiveDateTime: &scalarutils.Date{Year: 2019, Month: 11, Day: 10},
									Subject: &domain.FHIRReference{
										Display: gofakeit.BS(),
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
					return &domain.FHIRMedicationStatementRelayConnection{
						Edges: []*domain.FHIRMedicationStatementRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRMedicationStatement{
									ID: new(string),
									MedicationCodeableConcept: &domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
									},
									EffectiveDateTime: &scalarutils.Date{Year: 2019, Month: 11, Day: 10},
									Subject: &domain.FHIRReference{
										Display: gofakeit.BS(),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement - nil date" {
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
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
									},
									Subject: &domain.FHIRReference{
										Display: gofakeit.BS(),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement - nil subject" {
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
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
									},
									EffectiveDateTime: &scalarutils.Date{Year: 2019, Month: 11, Day: 10},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to search medication statement - invalid date" {
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
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
									},
									EffectiveDateTime: &scalarutils.Date{Year: 20190, Month: 11, Day: 10},
									Subject: &domain.FHIRReference{
										Display: gofakeit.BS(),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			got, err := u.PatientTimeline(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient timeline to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient timeline not to be nil for %v", tt.name)
				return
			}
		})
	}

}

func TestClinicalUseCaseImpl_PatientHealthTimeline(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.HealthTimelineInput
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.HealthTimeline
		wantErr bool
	}{
		{
			name: "Happy case: patient timeline",
			args: args{
				ctx: context.Background(),
				input: dto.HealthTimelineInput{
					PatientID: gofakeit.UUID(),
					Offset:    0,
					Limit:     20,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get patient timeline",
			args: args{
				ctx: context.Background(),
				input: dto.HealthTimelineInput{
					PatientID: gofakeit.UUID(),
					Offset:    0,
					Limit:     20,
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
			u := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "Happy case: patient timeline" {
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

				fakeFHIR.MockSearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRObservationRelayConnection, error) {
					status := dto.ObservationStatusFinal
					return &domain.FHIRObservationRelayConnection{
						Edges: []*domain.FHIRObservationRelayEdge{
							{
								Cursor: new(string),
								Node: &domain.FHIRObservation{
									ID:     new(string),
									Status: (*domain.ObservationStatusEnum)(&status),
									Code: domain.FHIRCodeableConcept{
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
										Text: gofakeit.BS(),
									},
									EffectiveDateTime: &scalarutils.Date{
										Year:  2000,
										Month: 1,
										Day:   1,
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}

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
										Coding: []*domain.FHIRCoding{{
											Display: gofakeit.BS(),
										}},
									},
									EffectiveDateTime: &scalarutils.Date{Year: 2019, Month: 11, Day: 10},
									Subject: &domain.FHIRReference{
										Display: gofakeit.BS(),
									},
								},
							},
						},
						PageInfo: &firebasetools.PageInfo{},
					}, nil
				}
			}

			if tt.name == "Sad case: failed to get patient timeline" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := u.PatientHealthTimeline(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientHealthTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient timeline to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient timeline not to be nil for %v", tt.name)
				return
			}
		})
	}

}
