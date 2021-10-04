package usecases_test

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	usecaseMock "github.com/savannahghi/clinical/pkg/clinical/usecases/mock"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"
)

func TestFHIRUseCaseImpl_CreateEpisodeOfCare_Unittest(t *testing.T) {
	ctx := context.Background()

	fh := fakeUsecaseIntr

	id := ksuid.New().String()

	episode := &domain.FHIREpisodeOfCare{
		ID:                   &id,
		Text:                 &domain.FHIRNarrative{},
		Identifier:           []*domain.FHIRIdentifier{},
		StatusHistory:        []*domain.FHIREpisodeofcareStatushistory{},
		Type:                 []*domain.FHIRCodeableConcept{},
		Diagnosis:            []*domain.FHIREpisodeofcareDiagnosis{},
		Patient:              &domain.FHIRReference{},
		ManagingOrganization: &domain.FHIRReference{},
		Period:               &domain.FHIRPeriod{},
		ReferralRequest:      []*domain.FHIRReference{},
		CareManager:          &domain.FHIRReference{},
		Team:                 []*domain.FHIRReference{},
		Account:              []*domain.FHIRReference{},
	}

	type args struct {
		ctx     context.Context
		episode domain.FHIREpisodeOfCare
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				episode: *episode,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				episode: domain.FHIREpisodeOfCare{
					ID:                   nil,
					Text:                 &domain.FHIRNarrative{},
					Identifier:           []*domain.FHIRIdentifier{},
					StatusHistory:        []*domain.FHIREpisodeofcareStatushistory{},
					Type:                 []*domain.FHIRCodeableConcept{},
					Diagnosis:            []*domain.FHIREpisodeofcareDiagnosis{},
					Patient:              &domain.FHIRReference{},
					ManagingOrganization: &domain.FHIRReference{},
					Period:               &domain.FHIRPeriod{},
					ReferralRequest:      []*domain.FHIRReference{},
					CareManager:          &domain.FHIRReference{},
					Team:                 []*domain.FHIRReference{},
					Account:              []*domain.FHIRReference{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateEpisodeOfCareFn = usecaseMock.NewFHIRMock().CreateEpisodeOfCareFn
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateEpisodeOfCareFn = func(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.CreateEpisodeOfCare(tt.args.ctx, tt.args.episode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateFHIRCondition_Unittest(t *testing.T) {
	ctx := context.Background()

	fh := fakeUsecaseIntr

	ID := ksuid.New().String()

	fhirconditionInput := &domain.FHIRConditionInput{
		ID:                 &ID,
		Identifier:         []*domain.FHIRIdentifierInput{},
		ClinicalStatus:     &domain.FHIRCodeableConceptInput{},
		VerificationStatus: &domain.FHIRCodeableConceptInput{},
		Category:           []*domain.FHIRCodeableConceptInput{},
		Severity:           &domain.FHIRCodeableConceptInput{},
		Code:               &domain.FHIRCodeableConceptInput{},
		BodySite:           []*domain.FHIRCodeableConceptInput{},
		Subject:            &domain.FHIRReferenceInput{},
		Encounter:          &domain.FHIRReferenceInput{},
		OnsetDateTime:      &scalarutils.Date{},
		OnsetAge:           &domain.FHIRAgeInput{},
		OnsetPeriod:        &domain.FHIRPeriodInput{},
		OnsetRange:         &domain.FHIRRangeInput{},
		OnsetString:        new(string),
		AbatementDateTime:  &scalarutils.Date{},
		AbatementAge:       &domain.FHIRAgeInput{},
		AbatementPeriod:    &domain.FHIRPeriodInput{},
		AbatementRange:     &domain.FHIRRangeInput{},
		AbatementString:    new(string),
		RecordedDate:       &scalarutils.Date{},
		Recorder:           &domain.FHIRReferenceInput{},
		Asserter:           &domain.FHIRReferenceInput{},
		Stage:              []*domain.FHIRConditionStageInput{},
		Evidence:           []*domain.FHIRConditionEvidenceInput{},
		Note:               []*domain.FHIRAnnotationInput{},
	}

	invalidfhirconditionInput := &domain.FHIRConditionInput{
		ID:                 nil,
		Identifier:         []*domain.FHIRIdentifierInput{},
		ClinicalStatus:     &domain.FHIRCodeableConceptInput{},
		VerificationStatus: &domain.FHIRCodeableConceptInput{},
		Category:           []*domain.FHIRCodeableConceptInput{},
		Severity:           &domain.FHIRCodeableConceptInput{},
		Code:               &domain.FHIRCodeableConceptInput{},
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRConditionInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: *fhirconditionInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: *invalidfhirconditionInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateFHIRConditionFn = usecaseMock.NewFHIRMock().CreateFHIRConditionFn
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.CreateFHIRCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_OpenOrganizationEpisodes_Unittest(t *testing.T) {
	ctx := context.Background()

	fh := fakeUsecaseIntr

	type args struct {
		ctx               context.Context
		providerSladeCode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:               ctx,
				providerSladeCode: "1234",
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:               ctx,
				providerSladeCode: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.OpenOrganizationEpisodesFn = usecaseMock.NewFHIRMock().OpenOrganizationEpisodesFn
			}

			if tt.name == "Sad case" {
				fakeFhir.OpenOrganizationEpisodesFn = func(ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.OpenOrganizationEpisodes(tt.args.ctx, tt.args.providerSladeCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.OpenOrganizationEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_GetORCreateOrganization_Unittest(t *testing.T) {
	ctx := context.Background()

	fh := fakeUsecaseIntr

	type args struct {
		ctx               context.Context
		providerSladeCode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:               ctx,
				providerSladeCode: "1234",
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:               ctx,
				providerSladeCode: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.GetORCreateOrganizationFn = usecaseMock.NewFHIRMock().CreateOrganizationFn
			}

			if tt.name == "Sad case" {
				fakeFhir.GetORCreateOrganizationFn = func(ctx context.Context, providerSladeCode string) (*string, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.GetORCreateOrganization(tt.args.ctx, tt.args.providerSladeCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetORCreateOrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateFHIROrganization_Unittest(t *testing.T) {
	ctx := context.Background()

	fh := fakeUsecaseIntr

	ID := ksuid.New().String()
	active := true
	testname := gofakeit.FirstName()

	orgInput := &domain.FHIROrganizationInput{
		ID:         &ID,
		Active:     &active,
		Identifier: []*domain.FHIRIdentifierInput{},
		Type:       []*domain.FHIRCodeableConceptInput{},
		Name:       &testname,
		Alias:      []string{"alias test"},
		Telecom:    []*domain.FHIRContactPointInput{},
		Address:    []*domain.FHIRAddressInput{},
	}

	invalidOrgInput := &domain.FHIROrganizationInput{
		ID:         &ID,
		Active:     new(bool),
		Identifier: []*domain.FHIRIdentifierInput{},
		Type:       []*domain.FHIRCodeableConceptInput{},
		Alias:      []string{"alias test"},
		Telecom:    []*domain.FHIRContactPointInput{},
		Address:    []*domain.FHIRAddressInput{},
	}

	type args struct {
		ctx   context.Context
		input domain.FHIROrganizationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: *orgInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: *invalidOrgInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateFHIROrganizationFn = usecaseMock.NewFHIRMock().CreateFHIROrganizationFn
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateFHIROrganizationFn = func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.CreateFHIROrganization(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateOrganization_Unittest(t *testing.T) {
	ctx := context.Background()

	fh := fakeUsecaseIntr

	type args struct {
		ctx               context.Context
		providerSladeCode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:               ctx,
				providerSladeCode: "1234",
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:               ctx,
				providerSladeCode: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateOrganizationFn = usecaseMock.NewFHIRMock().CreateOrganizationFn
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateOrganizationFn = func(ctx context.Context, providerSladeCode string) (*string, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.CreateOrganization(tt.args.ctx, tt.args.providerSladeCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateOrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_SearchFHIROrganization_Unittest(t *testing.T) {
	ctx := context.Background()

	fh := fakeUsecaseIntr

	type args struct {
		ctx    context.Context
		params map[string]interface{}
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
				params: map[string]interface{}{
					"test": "test",
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				params: map[string]interface{}{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.SearchFHIROrganizationFn = usecaseMock.NewFHIRMock().SearchFHIROrganizationFn
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchFHIROrganizationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.SearchFHIROrganization(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_GetOrganization_Unittest(t *testing.T) {
	ctx := context.Background()

	fh := fakeUsecaseIntr

	type args struct {
		ctx               context.Context
		providerSladeCode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:               ctx,
				providerSladeCode: "1234",
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:               ctx,
				providerSladeCode: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.GetOrganizationFn = usecaseMock.NewFHIRMock().GetOrganizationFn
			}

			if tt.name == "Sad case" {
				fakeFhir.GetOrganizationFn = func(ctx context.Context, providerSladeCode string) (*string, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.GetOrganization(tt.args.ctx, tt.args.providerSladeCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetOrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_SearchEpisodesByParam_Unittest(t *testing.T) {
	ctx := context.Background()

	fh := fakeUsecaseIntr

	params := url.Values{}

	type args struct {
		ctx          context.Context
		searchParams url.Values
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:          ctx,
				searchParams: params,
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
				fakeFhir.SearchEpisodesByParamFn = usecaseMock.NewFHIRMock().SearchEpisodesByParamFn
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchEpisodesByParamFn = func(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.SearchEpisodesByParam(tt.args.ctx, tt.args.searchParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchEpisodesByParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_POSTRequest_Unittest(t *testing.T) {
	fh := fakeUsecaseIntr

	searchParams := url.Values{}

	type args struct {
		resourceName string
		path         string
		params       url.Values
		body         io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				resourceName: "Encounter",
				path:         "_search",
				params:       searchParams,
				body:         nil,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				resourceName: "Encounter",
				path:         "",
				params:       searchParams,
				body:         nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.POSTRequestFn = usecaseMock.NewFHIRMock().POSTRequestFn
			}

			if tt.name == "Sad case" {
				fakeFhir.POSTRequestFn = func(resourceName, path string, params url.Values, body io.Reader) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.POSTRequest(tt.args.resourceName, tt.args.path, tt.args.params, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.POSTRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_StartEpisodeByOtp_Unittest(t *testing.T) {
	fh := fakeUsecaseIntr

	validEpisodeInput := &domain.OTPEpisodeCreationInput{
		PatientID:    "test",
		ProviderCode: "1234",
		Msisdn:       "+254711223344",
		Otp:          "1234",
		FullAccess:   false,
	}

	invalidEpisodeInput := &domain.OTPEpisodeCreationInput{
		PatientID:    "",
		ProviderCode: "1234",
		Msisdn:       "+254711223344",
		Otp:          "1234",
		FullAccess:   false,
	}

	type args struct {
		ctx   context.Context
		input domain.OTPEpisodeCreationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   nil,
				input: *validEpisodeInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   nil,
				input: *invalidEpisodeInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.StartEpisodeByOtpFn = usecaseMock.NewFHIRMock().StartEpisodeByOtpFn
			}

			if tt.name == "Sad case" {
				fakeFhir.StartEpisodeByOtpFn = func(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.StartEpisodeByOtp(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.StartEpisodeByOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_UpgradeEpisode_Unittest(t *testing.T) {
	fh := fakeUsecaseIntr

	type args struct {
		ctx   context.Context
		input domain.OTPEpisodeUpgradeInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: nil,
				input: domain.OTPEpisodeUpgradeInput{
					EpisodeID: ksuid.New().String(),
					Msisdn:    "+254711223344",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: nil,
				input: domain.OTPEpisodeUpgradeInput{
					EpisodeID: "",
					Msisdn:    "+254711223344",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.UpgradeEpisodeFn = usecaseMock.NewFHIRMock().UpgradeEpisodeFn
			}

			if tt.name == "Sad case" {
				fakeFhir.UpgradeEpisodeFn = func(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.UpgradeEpisode(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpgradeEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_SearchEpisodeEncounter_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	type args struct {
		ctx              context.Context
		episodeReference string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:              ctx,
				episodeReference: "test_episode",
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:              ctx,
				episodeReference: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.SearchEpisodeEncounterFn = usecaseMock.NewFHIRMock().SearchEpisodeEncounterFn
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchEpisodeEncounterFn = func(ctx context.Context, episodeReference string) (*domain.FHIREncounterRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.SearchEpisodeEncounter(tt.args.ctx, tt.args.episodeReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchEpisodeEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_EndEncounter_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	type args struct {
		ctx         context.Context
		encounterID string
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
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				encounterID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.EndEncounterFn = usecaseMock.NewFHIRMock().EndEncounterFn
			}

			if tt.name == "Sad case" {
				fakeFhir.EndEncounterFn = func(ctx context.Context, encounterID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.EndEncounter(tt.args.ctx, tt.args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.EndEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_EndEpisode_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	type args struct {
		ctx       context.Context
		episodeID string
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
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				episodeID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.EndEpisodeFn = usecaseMock.NewFHIRMock().EndEpisodeFn
			}

			if tt.name == "Sad case" {
				fakeFhir.EndEpisodeFn = func(ctx context.Context, episodeID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.EndEpisode(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.EndEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_GetActiveEpisode_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	type args struct {
		ctx       context.Context
		episodeID string
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
			},
			wantErr: false,
		},

		{
			name: "Sad case",
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
				fakeFhir.GetActiveEpisodeFn = usecaseMock.NewFHIRMock().GetActiveEpisodeFn
			}

			if tt.name == "Sad case" {
				fakeFhir.GetActiveEpisodeFn = func(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.GetActiveEpisode(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetActiveEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_SearchFHIRServiceRequest_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	type args struct {
		ctx    context.Context
		params map[string]interface{}
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
				params: map[string]interface{}{
					"test": "test",
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"test": "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.SearchFHIRServiceRequestFn = usecaseMock.NewFHIRMock().SearchFHIRServiceRequestFn
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchFHIRServiceRequestFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.SearchFHIRServiceRequest(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateFHIRServiceRequest_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	ID := ksuid.New().String()

	type args struct {
		ctx   context.Context
		input domain.FHIRServiceRequestInput
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
				input: domain.FHIRServiceRequestInput{
					ID:          &ID,
					Identifier:  []*domain.FHIRIdentifierInput{},
					BasedOn:     []*domain.FHIRReferenceInput{},
					Replaces:    []*domain.FHIRReferenceInput{},
					Requisition: &domain.FHIRIdentifierInput{},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: domain.FHIRServiceRequestInput{
					Identifier:  []*domain.FHIRIdentifierInput{},
					BasedOn:     []*domain.FHIRReferenceInput{},
					Replaces:    []*domain.FHIRReferenceInput{},
					Requisition: &domain.FHIRIdentifierInput{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateFHIRServiceRequestFn = usecaseMock.NewFHIRMock().CreateFHIRServiceRequestFn
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateFHIRServiceRequestFn = func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.CreateFHIRServiceRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_SearchFHIRAllergyIntolerance_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	type args struct {
		ctx    context.Context
		params map[string]interface{}
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
				params: map[string]interface{}{
					"test": "test",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"test": "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.SearchFHIRAllergyIntoleranceFn = usecaseMock.NewFHIRMock().SearchFHIRAllergyIntoleranceFn
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.SearchFHIRAllergyIntolerance(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateFHIRAllergyIntolerance_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	ID := ksuid.New().String()

	type args struct {
		ctx   context.Context
		input domain.FHIRAllergyIntoleranceInput
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
				input: domain.FHIRAllergyIntoleranceInput{
					ID:                 &ID,
					Identifier:         []*domain.FHIRIdentifierInput{},
					ClinicalStatus:     domain.FHIRCodeableConceptInput{},
					VerificationStatus: domain.FHIRCodeableConceptInput{},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: domain.FHIRAllergyIntoleranceInput{
					Identifier:         []*domain.FHIRIdentifierInput{},
					ClinicalStatus:     domain.FHIRCodeableConceptInput{},
					VerificationStatus: domain.FHIRCodeableConceptInput{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateFHIRAllergyIntoleranceFn = usecaseMock.NewFHIRMock().CreateFHIRAllergyIntoleranceFn
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.CreateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
