package usecases_test

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	profileUtilsMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/profileutils/mock"
	usecaseMock "github.com/savannahghi/clinical/pkg/clinical/usecases/mock"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"
)

const (
	baseFHIRURL = "https://healthcare.googleapis.com/v1"
)

var fakeProfileUtils profileUtilsMock.FakeUserProfileRepository

func TestClinicalUseCaseImpl_ProblemSummary_Unittest(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeClinicalInteractor(ctx)
	if err != nil {
		t.Errorf("failed to fake initialize fake clinical interactor: %v", err)
		return
	}

	patientId := "test1234"

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
				patientID: patientId,
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
			if tt.name == "Happy case" {

				fakeProfileUtils.GetLoggedInUserFn = func(ctx context.Context) (*profileutils.UserInfo, error) {
					return &profileutils.UserInfo{
						DisplayName: "test user",
						Email:       firebasetools.TestUserEmail,
						PhoneNumber: interserviceclient.TestUserPhoneNumber,
						PhotoURL:    gofakeit.ImageURL(50, 50),
					}, nil
				}

				fakeFHIR.CreateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					return &domain.FHIRConditionRelayPayload{
						Resource: &domain.FHIRCondition{
							ID: &patientId,
						},
					}, nil
				}

				fakeFHIR.SearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
					cursor := "Composition"
					return &domain.FHIRConditionRelayConnection{
						Edges: []*domain.FHIRConditionRelayEdge{
							{
								Cursor: &cursor,
								Node:   &domain.FHIRCondition{},
							},
						},
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}

				fakeFHIR.POSTRequestFn = func(resourceName, path string, params url.Values, body io.Reader) ([]byte, error) {
					return []byte("some-byte"), nil
				}

				FHIRRepoMock.FHIRRestURLFn = func() string {
					return "https://healthcare.googleapis.com/v1/projects/bewell-app-ci/locations/europe-west4/datasets/healthcloud-bewell-staging/fhirStores/healthcloud-bewell-fhir-staging/fhir/EpisodeOfCare/_search?patient=Patient%2F1e216562-3f8a-4ec9-977b-2e12b9fdeb39"
				}
			}

			if tt.name == "Sad case" {
				fakePatient.ProblemSummaryFn = func(ctx context.Context, patientID string) ([]string, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.ClinicalUseCase.ProblemSummary(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.ProblemSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_VisitSummary_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	Id := "test1234ID"

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
				encounterID: Id,
				count:       0,
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
				encounterID: Id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFHIR.CreateFHIREncounterFn = func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							ID: &Id,
						},
					}, nil
				}

				fakeFHIR.GetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							ID: &Id,
						},
					}, nil
				}

				fakeFHIR.CreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return &domain.FHIRAllergyIntoleranceRelayPayload{
						Resource: &domain.FHIRAllergyIntolerance{
							ID: &Id,
						},
					}, nil
				}

				fakeFHIR.SearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return &domain.FHIRAllergyIntoleranceRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}

				fakeFHIR.SearchFHIREncounterFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
					return &domain.FHIREncounterRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}
			}

			if tt.name == "Sad case: nil encounter ID" {
				fakePatient.VisitSummaryFn = func(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: no count" {
				fakePatient.VisitSummaryFn = func(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
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

func TestClinicalUseCaseImpl_PatientTimelineWithCount_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	type args struct {
		ctx       context.Context
		episodeID string
		count     int
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
				episodeID: ksuid.New().String(),
				count:     0,
			},
			wantErr: false,
		},

		{
			name: "Sad case: nil episode ID",
			args: args{
				ctx:   ctx,
				count: 0,
			},
			wantErr: true,
		},

		{
			name: "Sad case: no count",
			args: args{
				ctx:       ctx,
				episodeID: ksuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakePatient.PatientTimelineWithCountFn = usecaseMock.NewClinicalMock().PatientTimelineWithCount
			}

			if tt.name == "Sad case: nil episode ID" {
				fakePatient.PatientTimelineWithCountFn = func(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: no count" {
				fakePatient.PatientTimelineWithCountFn = func(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.PatientTimelineWithCount(tt.args.ctx, tt.args.episodeID, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientTimelineWithCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_PatientSearch_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

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
			if tt.name == "Happy case" {
				fakePatient.PatientSearchFn = usecaseMock.NewClinicalMock().PatientSearch
			}

			if tt.name == "Sad case" {
				fakePatient.PatientSearchFn = func(ctx context.Context, search string) (*domain.PatientConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.PatientSearch(tt.args.ctx, tt.args.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_ContactsToContactPointInput_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

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
				ctx:    ctx,
				phones: []*domain.PhoneNumberInput{},
				emails: []*domain.EmailInput{},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				phones: []*domain.PhoneNumberInput{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakePatient.ContactsToContactPointInputFn = usecaseMock.NewClinicalMock().ContactsToContactPointInput
			}

			if tt.name == "Sad case" {
				fakePatient.ContactsToContactPointInputFn = func(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.ContactsToContactPointInput(tt.args.ctx, tt.args.phones, tt.args.emails)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.ContactsToContactPointInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_CreatePatient_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	ID := ksuid.New().String()

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
					ID: &ID,
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
			if tt.name == "Happy case" {
				fakePatient.CreatePatientFn = usecaseMock.NewClinicalMock().CreatePatient
			}

			if tt.name == "Sad case" {
				fakePatient.CreatePatientFn = func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.CreatePatient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.CreatePatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_FindPatientByID_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	patientID := ksuid.New().String()

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
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  patientID,
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
			if tt.name == "Happy case" {
				fakePatient.FindPatientByIDFn = usecaseMock.NewClinicalMock().FindPatientByID
			}

			if tt.name == "Sad case" {
				fakePatient.FindPatientByIDFn = func(ctx context.Context, id string) (*domain.PatientPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.FindPatientByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.FindPatientByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_UpdatePatient_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	patientID := ksuid.New().String()

	type args struct {
		ctx   context.Context
		input domain.SimplePatientRegistrationInput
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
				input: domain.SimplePatientRegistrationInput{
					ID:                      patientID,
					Names:                   []*domain.NameInput{},
					IdentificationDocuments: []*domain.IdentificationDocument{},
					BirthDate:               scalarutils.Date{},
					PhoneNumbers:            []*domain.PhoneNumberInput{},
					Photos:                  []*domain.PhotoInput{},
					Emails:                  []*domain.EmailInput{},
					PhysicalAddresses:       []*domain.PhysicalAddress{},
					PostalAddresses:         []*domain.PostalAddress{},
					Gender:                  "TEST GENDER",
					Active:                  true,
					MaritalStatus:           "TEST MARITAL STATUS",
					Languages:               []enumutils.Language{},
					ReplicateUSSD:           true,
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
			if tt.name == "Happy case" {
				fakePatient.UpdatePatientFn = usecaseMock.NewClinicalMock().UpdatePatient
			}

			if tt.name == "Sad case" {
				fakePatient.UpdatePatientFn = func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.UpdatePatient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.UpdatePatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_AddNextOfKin_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	type args struct {
		ctx   context.Context
		input domain.SimpleNextOfKinInput
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
				input: domain.SimpleNextOfKinInput{
					PatientID: ksuid.New().String(),
					Gender:    "TEST GENDER",
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
			if tt.name == "Happy case" {
				fakePatient.AddNextOfKinFn = usecaseMock.NewClinicalMock().AddNextOfKin
			}

			if tt.name == "Sad case" {
				fakePatient.AddNextOfKinFn = func(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.AddNextOfKin(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.AddNextOfKin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_AddNHIF_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	testInput := ksuid.New().String()

	type args struct {
		ctx   context.Context
		input *domain.SimpleNHIFInput
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
				input: &domain.SimpleNHIFInput{
					PatientID:        ksuid.New().String(),
					MembershipNumber: ksuid.New().String(),
					FrontImageBase64: &testInput,
					RearImageBase64:  &testInput,
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
			if tt.name == "Happy case" {
				fakePatient.AddNHIFFn = usecaseMock.NewClinicalMock().AddNHIF
			}

			if tt.name == "Sad case" {
				fakePatient.AddNHIFFn = func(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.AddNHIF(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.AddNHIF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_RegisterUser_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	type args struct {
		ctx   context.Context
		input domain.SimplePatientRegistrationInput
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
				input: domain.SimplePatientRegistrationInput{
					ID:                      ksuid.New().String(),
					Names:                   []*domain.NameInput{},
					IdentificationDocuments: []*domain.IdentificationDocument{},
					BirthDate:               scalarutils.Date{},
					PhoneNumbers:            []*domain.PhoneNumberInput{},
					Photos:                  []*domain.PhotoInput{},
					Emails:                  []*domain.EmailInput{},
					PhysicalAddresses:       []*domain.PhysicalAddress{},
					PostalAddresses:         []*domain.PostalAddress{},
					Gender:                  "TEST GENDER",
					Active:                  true,
					MaritalStatus:           "TEST MARITAL STATUS",
					Languages:               []enumutils.Language{},
					ReplicateUSSD:           true,
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
			if tt.name == "Happy case" {
				fakePatient.RegisterUserFn = usecaseMock.NewClinicalMock().RegisterUser
			}

			if tt.name == "Sad case" {
				fakePatient.RegisterUserFn = func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := i.RegisterUser(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_CreateUpdatePatientExtraInformation_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	maritalStatus := ksuid.New().String()

	type args struct {
		ctx   context.Context
		input domain.PatientExtraInformationInput
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
				input: domain.PatientExtraInformationInput{
					MaritalStatus: (*domain.MaritalStatus)(&maritalStatus),
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
			if tt.name == "Happy case" {
				fakePatient.CreateUpdatePatientExtraInformationFn = usecaseMock.NewClinicalMock().CreateUpdatePatientExtraInformation
			}

			if tt.name == "Sad case" {
				fakePatient.CreateUpdatePatientExtraInformationFn = func(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			_, err := i.CreateUpdatePatientExtraInformation(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.CreateUpdatePatientExtraInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_AllergySummary_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	patientID := ksuid.New().String()

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
				patientID: patientID,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				patientID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakePatient.AllergySummaryFn = usecaseMock.NewClinicalMock().AllergySummary
			}

			if tt.name == "Sad case" {
				fakePatient.AllergySummaryFn = func(ctx context.Context, patientID string) ([]string, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := i.AllergySummary(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.AllergySummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_DeleteFHIRPatientByPhone_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	testPhone := "+254720000000"

	type args struct {
		ctx         context.Context
		phoneNumber string
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
				phoneNumber: testPhone,
			},
			wantErr: false,
		},
		{
			name: "Sad case: empty phone",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
			},
			wantErr: true,
		},

		{
			name: "Sad case: invalid phone",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakePatient.DeleteFHIRPatientByPhoneFn = usecaseMock.NewClinicalMock().DeleteFHIRPatientByPhone
			}

			if tt.name == "Sad case: empty phone" {
				fakePatient.DeleteFHIRPatientByPhoneFn = func(ctx context.Context, phoneNumber string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: invalid phone" {
				fakePatient.DeleteFHIRPatientByPhoneFn = func(ctx context.Context, phoneNumber string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			_, err := i.DeleteFHIRPatientByPhone(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.DeleteFHIRPatientByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_StartEpisodeByBreakGlass_Unittest(t *testing.T) {
	ctx := context.Background()
	i := fakeUsecaseIntr

	type args struct {
		ctx   context.Context
		input domain.BreakGlassEpisodeCreationInput
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
				input: domain.BreakGlassEpisodeCreationInput{
					PatientID:       ksuid.New().String(),
					ProviderCode:    ksuid.New().String(),
					PractitionerUID: ksuid.New().String(),
					ProviderPhone:   ksuid.New().String(),
					Otp:             "000000",
					FullAccess:      true,
					PatientPhone:    ksuid.New().String(),
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case: empty patient ID",
			args: args{
				ctx: ctx,
				input: domain.BreakGlassEpisodeCreationInput{
					PatientID:       "",
					ProviderCode:    ksuid.New().String(),
					PractitionerUID: ksuid.New().String(),
					ProviderPhone:   ksuid.New().String(),
					Otp:             "000000",
					FullAccess:      true,
					PatientPhone:    ksuid.New().String(),
				},
			},
			wantErr: true,
		},

		{
			name: "Sad case: no input",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakePatient.StartEpisodeByBreakGlassFn = usecaseMock.NewClinicalMock().StartEpisodeByBreakGlass
			}

			if tt.name == "Sad case: empty patient ID" {
				fakePatient.StartEpisodeByBreakGlassFn = func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: no input" {
				fakePatient.StartEpisodeByBreakGlassFn = func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := i.StartEpisodeByBreakGlass(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.StartEpisodeByBreakGlass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
