package clinical_test

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"testing"

	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhir/mock"
	fakeDatasetMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhirdataset/mock"
	fakeFirebaseMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/interserviceclient"
	"github.com/segmentio/ksuid"
)

func TestClinicalUseCaseImpl_ProblemSummary_Unittest(t *testing.T) {

	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	FakefirestoreDB := fakeFirebaseMock.NewFakeFirebaseMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, FakefirestoreDB, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

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
				ctx:       ctx,
				patientID: ksuid.New().String(),
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				Fakefhir.MockSearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
					return nil, fmt.Errorf("error")
				}
			}
			_, err := m.ProblemSummary(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.ProblemSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_VisitSummary_Unittest(t *testing.T) {
	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	FakefirestoreDB := fakeFirebaseMock.NewFakeFirebaseMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, FakefirestoreDB, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

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
				encounterID: ksuid.New().String(),
				count:       10,
			},
			wantErr: false,
		},

		{
			name: "Sad case: nil encounter ID",
			args: args{
				ctx:   ctx,
				count: 0,
			},
			wantErr: true,
		},

		{
			name: "Sad case: no count",
			args: args{
				ctx:         ctx,
				encounterID: ksuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: nil encounter ID" {
				Fakefhir.MockGetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("error")
				}
			}

			if tt.name == "Sad case: no count" {
				Fakefhir.MockSearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return nil, fmt.Errorf("error")
				}
			}
			_, err := m.VisitSummary(tt.args.ctx, tt.args.encounterID, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.VisitSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_PatientSearch_Unittest(t *testing.T) {
	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	FakefirestoreDB := fakeFirebaseMock.NewFakeFirebaseMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, FakefirestoreDB, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	searchPatient := "Test user"

	type args struct {
		ctx    context.Context
		search string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				search: searchPatient,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				FakefhirRepository.MockPOSTRequestFn = func(resourceName, path string, params url.Values, body io.Reader) ([]byte, error) {
					return nil, fmt.Errorf("error")
				}
			}
			_, err := m.PatientSearch(tt.args.ctx, tt.args.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_ContactsToContactPointInput_Unittest(t *testing.T) {
	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	FakefirestoreDB := fakeFirebaseMock.NewFakeFirebaseMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, FakefirestoreDB, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	type args struct {
		ctx    context.Context
		phones []*domain.PhoneNumberInput
		emails []*domain.EmailInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				phones: []*domain.PhoneNumberInput{
					{
						Msisdn:             interserviceclient.TestUserPhoneNumber,
						VerificationCode:   "20002",
						IsUssd:             false,
						CommunicationOptIn: true,
					},
				},
				emails: []*domain.EmailInput{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := m.ContactsToContactPointInput(tt.args.ctx, tt.args.phones, tt.args.emails)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.ContactsToContactPointInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_CreatePatient_Unittest(t *testing.T) {
	ctx := context.Background()

	FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
	FakefhirRepository := fakeDatasetMock.NewFakeFHIRRepositoryMock()
	Fakefhir := fakeFHIRMock.NewFHIRMock()
	FakefirestoreDB := fakeFirebaseMock.NewFakeFirebaseMock()
	Fakeopenconceptlab := fakeOCLMock.NewFakeOCLMock()

	n := infrastructure.NewInfrastructureInteractor(FakeExt, FakefhirRepository, Fakefhir, FakefirestoreDB, Fakeopenconceptlab)
	m := clinicalUsecase.NewUseCasesClinicalImpl(n)

	ID := ksuid.New().String()

	gender := domain.PatientGenderEnumMale

	type args struct {
		ctx   context.Context
		input domain.FHIRPatientInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: domain.FHIRPatientInput{
					ID:     &ID,
					Gender: &gender,
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				FakefhirRepository.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}) ([]byte, error) {
					return nil, fmt.Errorf("error")
				}
			}
			_, err := m.CreatePatient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.CreatePatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// func TestClinicalUseCaseImpl_AddNextOfKin_Unittest(t *testing.T) {
// 	ctx := context.Background()
// 	i := fakeUsecaseIntr

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
// 			if tt.name == "Happy case" {
// 				fakePatient.AddNextOfKinFn = usecaseMock.NewClinicalMock().AddNextOfKin
// 			}

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
// 	i := fakeUsecaseIntr

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
// 			if tt.name == "Happy case" {
// 				fakePatient.AddNHIFFn = usecaseMock.NewClinicalMock().AddNHIF
// 			}

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
// 	i := fakeUsecaseIntr

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
// 			if tt.name == "Happy case" {
// 				fakePatient.RegisterUserFn = usecaseMock.NewClinicalMock().RegisterUser
// 			}

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
// 	i := fakeUsecaseIntr

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
// 			if tt.name == "Happy case" {
// 				fakePatient.CreateUpdatePatientExtraInformationFn = usecaseMock.NewClinicalMock().CreateUpdatePatientExtraInformation
// 			}

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
// 	i := fakeUsecaseIntr

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
// 			if tt.name == "Happy case" {
// 				fakePatient.AllergySummaryFn = usecaseMock.NewClinicalMock().AllergySummary
// 			}

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
// 	i := fakeUsecaseIntr

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
// 			if tt.name == "Happy case" {
// 				fakePatient.DeleteFHIRPatientByPhoneFn = usecaseMock.NewClinicalMock().DeleteFHIRPatientByPhone
// 			}

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
// 	i := fakeUsecaseIntr

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
// 					MFLCode:    ksuid.New().String(),
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
// 					MFLCode:    ksuid.New().String(),
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
// 			if tt.name == "Happy case" {
// 				fakePatient.StartEpisodeByBreakGlassFn = usecaseMock.NewClinicalMock().StartEpisodeByBreakGlass
// 			}

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
