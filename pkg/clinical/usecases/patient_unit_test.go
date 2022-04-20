package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/profileutils"
	"github.com/segmentio/ksuid"
)

func TestClinicalUseCaseImpl_ProblemSummary_Unittest(t *testing.T) {
	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

	type args struct {
		ctx       context.Context
		patientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       context.Background(),
				patientID: ksuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search fhir condition",
			args: args{
				ctx:       context.Background(),
				patientID: ksuid.New().String(),
			},
			wantErr: true,
		},

		{
			name: "Sad case failed to get problem summary",
			args: args{
				ctx:       context.Background(),
				patientID: ksuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case nil condition code",
			args: args{
				ctx:       context.Background(),
				patientID: ksuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case empty condition code text",
			args: args{
				ctx:       context.Background(),
				patientID: ksuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: failed to search fhir condition" {
				fakeFHIR.SearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
					return nil, fmt.Errorf("failed to search fhir condition")
				}
			}

			if tt.name == "Sad case failed to get problem summary" {
				fakePatient.ProblemSummaryFn = func(ctx context.Context, patientID string) ([]string, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case nil condition code" {
				fakeFHIR.SearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
					cursor := "1"

					testID := ksuid.New().String()

					return &domain.FHIRConditionRelayConnection{
						Edges: []*domain.FHIRConditionRelayEdge{
							{
								Cursor: &cursor,
								Node: &domain.FHIRCondition{
									ID: &testID,

									Code: nil,
								},
							},
						},
					}, nil
				}
			}
			if tt.name == "Sad case empty condition code text" {
				fakeFHIR.SearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
					cursor := "1"

					testID := ksuid.New().String()

					return &domain.FHIRConditionRelayConnection{
						Edges: []*domain.FHIRConditionRelayEdge{
							{
								Cursor: &cursor,
								Node: &domain.FHIRCondition{
									ID: &testID,

									Code: &domain.FHIRCodeableConcept{
										Text: "",
									},
								},
							},
						},
					}, nil
				}
			}

			_, err := i.ProblemSummary(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.ProblemSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_VisitSummary_Unittest(t *testing.T) {
	ctx := context.Background()
	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

	type args struct {
		ctx         context.Context
		encounterID string
		count       int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get logged in user",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get fhir encounter",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to search fhir allergy intollerance",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to search fhir encounter",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to search fhir condition",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to search fhir observation",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to search fhir composition",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to search fhir medication request",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to search fhir service request",
			args: args{
				ctx:         ctx,
				encounterID: "1234",
				count:       0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad case: failed to get logged in user" {
				fakeBaseExtension.GetLoggedInUserFn = func(ctx context.Context) (*profileutils.UserInfo, error) {
					return nil, fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "Sad case: failed to get fhir encounter" {
				fakeFHIR.GetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("failed to get fhir encounter")
				}
			}

			if tt.name == "Sad case: failed to search fhir allergy intollerance" {
				fakeFHIR.SearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return nil, fmt.Errorf("failed to search fhir allergy intollerance")
				}
			}

			if tt.name == "Sad case: failed to search fhir encounter" {
				fakeFHIR.SearchFHIREncounterFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
					return nil, fmt.Errorf("failed to search fhir encounter")
				}
			}

			if tt.name == "Sad case: failed to search fhir condition" {
				fakeFHIR.SearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
					return nil, fmt.Errorf("failed to search fhir condition")
				}
			}

			if tt.name == "Sad case: failed to search fhir observation" {
				fakeFHIR.SearchFHIRObservationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
					return nil, fmt.Errorf("failed to search fhir observation")
				}
			}

			if tt.name == "Sad case: failed to search fhir composition" {
				fakeFHIR.SearchFHIRCompositionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
					return nil, fmt.Errorf("failed to search fhir composition")
				}
			}

			if tt.name == "Sad case: failed to search fhir medication request" {
				fakeFHIR.SearchFHIRMedicationRequestFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
					return nil, fmt.Errorf("failed to search fhir medication request")
				}
			}

			if tt.name == "Sad case: failed to search fhir service request" {
				fakeFHIR.SearchFHIRServiceRequestFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
					return nil, fmt.Errorf("failed to search fhir service request")
				}
			}

			_, err := i.VisitSummary(tt.args.ctx, tt.args.encounterID, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.VisitSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// func TestClinicalUseCaseImpl_PatientTimelineWithCount_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	type args struct {
// 		ctx       context.Context
// 		episodeID string
// 		count     int
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:       ctx,
// 				episodeID: ksuid.New().String(),
// 				count:     0,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case: nil episode ID",
// 			args: args{
// 				ctx:   ctx,
// 				count: 0,
// 			},
// 			wantErr: true,
// 		},

// 		{
// 			name: "Sad case: no count",
// 			args: args{
// 				ctx:       ctx,
// 				episodeID: ksuid.New().String(),
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case: nil episode ID" {
// 				fakePatient.PatientTimelineWithCountFn = func(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			if tt.name == "Sad case: no count" {
// 				fakePatient.PatientTimelineWithCountFn = func(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := i.PatientTimelineWithCount(tt.args.ctx, tt.args.episodeID, tt.args.count)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.PatientTimelineWithCount() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_PatientSearch_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	searchPatient := "Test user"

// 	type args struct {
// 		ctx    context.Context
// 		search string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:    ctx,
// 				search: searchPatient,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.PatientSearchFn = func(ctx context.Context, search string) (*domain.PatientConnection, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := i.PatientSearch(tt.args.ctx, tt.args.search)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.PatientSearch() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_ContactsToContactPointInput_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	type args struct {
// 		ctx    context.Context
// 		phones []*domain.PhoneNumberInput
// 		emails []*domain.EmailInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:    ctx,
// 				phones: []*domain.PhoneNumberInput{},
// 				emails: []*domain.EmailInput{},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:    ctx,
// 				phones: []*domain.PhoneNumberInput{},
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.ContactsToContactPointInputFn = func(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := i.ContactsToContactPointInput(tt.args.ctx, tt.args.phones, tt.args.emails)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.ContactsToContactPointInput() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_CreatePatient_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	ID := ksuid.New().String()

// 	type args struct {
// 		ctx   context.Context
// 		input domain.FHIRPatientInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				input: domain.FHIRPatientInput{
// 					ID: &ID,
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.CreatePatientFn = func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := i.CreatePatient(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.CreatePatient() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_FindPatientByID_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	patientID := ksuid.New().String()

// 	type args struct {
// 		ctx context.Context
// 		id  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				id:  patientID,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.FindPatientByIDFn = func(ctx context.Context, id string) (*domain.PatientPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := i.FindPatientByID(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.FindPatientByID() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_UpdatePatient_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	patientID := ksuid.New().String()

// 	type args struct {
// 		ctx   context.Context
// 		input domain.SimplePatientRegistrationInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				input: domain.SimplePatientRegistrationInput{
// 					ID:                      patientID,
// 					Names:                   []*domain.NameInput{},
// 					IdentificationDocuments: []*domain.IdentificationDocument{},
// 					BirthDate:               scalarutils.Date{},
// 					PhoneNumbers:            []*domain.PhoneNumberInput{},
// 					Photos:                  []*domain.PhotoInput{},
// 					Emails:                  []*domain.EmailInput{},
// 					PhysicalAddresses:       []*domain.PhysicalAddress{},
// 					PostalAddresses:         []*domain.PostalAddress{},
// 					Gender:                  "TEST GENDER",
// 					Active:                  true,
// 					MaritalStatus:           "TEST MARITAL STATUS",
// 					Languages:               []enumutils.Language{},
// 					ReplicateUSSD:           true,
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.UpdatePatientFn = func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := i.UpdatePatient(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.UpdatePatient() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_AddNextOfKin_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	type args struct {
// 		ctx   context.Context
// 		input domain.SimpleNextOfKinInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				input: domain.SimpleNextOfKinInput{
// 					PatientID: ksuid.New().String(),
// 					Gender:    "TEST GENDER",
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.AddNextOfKinFn = func(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := i.AddNextOfKin(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.AddNextOfKin() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_AddNHIF_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	testInput := ksuid.New().String()

// 	type args struct {
// 		ctx   context.Context
// 		input *domain.SimpleNHIFInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				input: &domain.SimpleNHIFInput{
// 					PatientID:        ksuid.New().String(),
// 					MembershipNumber: ksuid.New().String(),
// 					FrontImageBase64: &testInput,
// 					RearImageBase64:  &testInput,
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.AddNHIFFn = func(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := i.AddNHIF(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.AddNHIF() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_RegisterUser_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	type args struct {
// 		ctx   context.Context
// 		input domain.SimplePatientRegistrationInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				input: domain.SimplePatientRegistrationInput{
// 					ID:                      ksuid.New().String(),
// 					Names:                   []*domain.NameInput{},
// 					IdentificationDocuments: []*domain.IdentificationDocument{},
// 					BirthDate:               scalarutils.Date{},
// 					PhoneNumbers:            []*domain.PhoneNumberInput{},
// 					Photos:                  []*domain.PhotoInput{},
// 					Emails:                  []*domain.EmailInput{},
// 					PhysicalAddresses:       []*domain.PhysicalAddress{},
// 					PostalAddresses:         []*domain.PostalAddress{},
// 					Gender:                  "TEST GENDER",
// 					Active:                  true,
// 					MaritalStatus:           "TEST MARITAL STATUS",
// 					Languages:               []enumutils.Language{},
// 					ReplicateUSSD:           true,
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.RegisterUserFn = func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			_, err := i.RegisterUser(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_CreateUpdatePatientExtraInformation_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	maritalStatus := ksuid.New().String()

// 	type args struct {
// 		ctx   context.Context
// 		input domain.PatientExtraInformationInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				input: domain.PatientExtraInformationInput{
// 					MaritalStatus: (*domain.MaritalStatus)(&maritalStatus),
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.CreateUpdatePatientExtraInformationFn = func(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
// 					return false, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := i.CreateUpdatePatientExtraInformation(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.CreateUpdatePatientExtraInformation() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_AllergySummary_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	patientID := ksuid.New().String()

// 	type args struct {
// 		ctx       context.Context
// 		patientID string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:       ctx,
// 				patientID: patientID,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:       ctx,
// 				patientID: "",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case" {
// 				fakePatient.AllergySummaryFn = func(ctx context.Context, patientID string) ([]string, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			_, err := i.AllergySummary(tt.args.ctx, tt.args.patientID)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.AllergySummary() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_DeleteFHIRPatientByPhone_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	testPhone := "+254720000000"

// 	type args struct {
// 		ctx         context.Context
// 		phoneNumber string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:         ctx,
// 				phoneNumber: testPhone,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case: empty phone",
// 			args: args{
// 				ctx:         ctx,
// 				phoneNumber: "",
// 			},
// 			wantErr: true,
// 		},

// 		{
// 			name: "Sad case: invalid phone",
// 			args: args{
// 				ctx:         ctx,
// 				phoneNumber: "+254",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case: empty phone" {
// 				fakePatient.DeleteFHIRPatientByPhoneFn = func(ctx context.Context, phoneNumber string) (bool, error) {
// 					return false, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			if tt.name == "Sad case: invalid phone" {
// 				fakePatient.DeleteFHIRPatientByPhoneFn = func(ctx context.Context, phoneNumber string) (bool, error) {
// 					return false, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			_, err := i.DeleteFHIRPatientByPhone(tt.args.ctx, tt.args.phoneNumber)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.DeleteFHIRPatientByPhone() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_StartEpisodeByBreakGlass_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := usecases.NewClinicalUseCaseImpl(fakeInfra, &fakeFHIR)

// 	type args struct {
// 		ctx   context.Context
// 		input domain.BreakGlassEpisodeCreationInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				input: domain.BreakGlassEpisodeCreationInput{
// 					PatientID:       ksuid.New().String(),
// 					ProviderCode:    ksuid.New().String(),
// 					PractitionerUID: ksuid.New().String(),
// 					ProviderPhone:   ksuid.New().String(),
// 					Otp:             "000000",
// 					FullAccess:      true,
// 					PatientPhone:    ksuid.New().String(),
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case: empty patient ID",
// 			args: args{
// 				ctx: ctx,
// 				input: domain.BreakGlassEpisodeCreationInput{
// 					PatientID:       "",
// 					ProviderCode:    ksuid.New().String(),
// 					PractitionerUID: ksuid.New().String(),
// 					ProviderPhone:   ksuid.New().String(),
// 					Otp:             "000000",
// 					FullAccess:      true,
// 					PatientPhone:    ksuid.New().String(),
// 				},
// 			},
// 			wantErr: true,
// 		},

// 		{
// 			name: "Sad case: no input",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Sad case: empty patient ID" {
// 				fakePatient.StartEpisodeByBreakGlassFn = func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			if tt.name == "Sad case: no input" {
// 				fakePatient.StartEpisodeByBreakGlassFn = func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			_, err := i.StartEpisodeByBreakGlass(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.StartEpisodeByBreakGlass() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }
