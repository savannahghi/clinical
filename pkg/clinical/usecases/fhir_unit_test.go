package usecases_test

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/firebasetools"
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
				fakeFhir.CreateEpisodeOfCareFn = func(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
					return &domain.EpisodeOfCarePayload{
						TotalVisits: 1,
					}, nil
				}
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
				fakeFhir.CreateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					id := ksuid.New().String()
					return &domain.FHIRConditionRelayPayload{
						Resource: &domain.FHIRCondition{
							ID: &id,
						},
					}, nil
				}
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
				fakeFhir.OpenOrganizationEpisodesFn = func(ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error) {
					id := ksuid.New().String()
					episodeofcare := &domain.FHIREpisodeOfCare{
						ID: &id,
					}
					return []*domain.FHIREpisodeOfCare{
						episodeofcare,
					}, nil
				}
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
				fakeFhir.GetORCreateOrganizationFn = func(ctx context.Context, providerSladeCode string) (*string, error) {
					org := "test-organization"
					return &org, nil
				}
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
				fakeFhir.CreateFHIROrganizationFn = func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
					ID := ksuid.New().String()
					return &domain.FHIROrganizationRelayPayload{
						Resource: &domain.FHIROrganization{
							ID: &ID,
						},
					}, nil
				}
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
				fakeFhir.CreateOrganizationFn = func(ctx context.Context, providerSladeCode string) (*string, error) {
					org := "test-organization"
					return &org, nil
				}
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
				fakeFhir.SearchFHIROrganizationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {
					return &domain.FHIROrganizationRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}
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
				fakeFhir.GetOrganizationFn = func(ctx context.Context, providerSladeCode string) (*string, error) {
					org := "test-organization"
					return &org, nil
				}
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
				fakeFhir.SearchEpisodesByParamFn = func(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
					id := ksuid.New().String()
					fhirEpisode := &domain.FHIREpisodeOfCare{
						ID: &id,
					}
					return []*domain.FHIREpisodeOfCare{
						fhirEpisode,
					}, nil
				}
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
				fakeFhir.POSTRequestFn = func(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
					return []byte(""), nil
				}
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

func TestUnit_OpenEpisodes(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr
	patientReference := fmt.Sprintf("Patient/%s", ksuid.New().String())

	type args struct {
		ctx              context.Context
		patientReference string
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
				patientReference: patientReference,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:              ctx,
				patientReference: patientReference,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.OpenEpisodesFn = func(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
					id := ksuid.New().String()
					ec := &domain.FHIREpisodeOfCare{
						ID: &id,
					}
					return []*domain.FHIREpisodeOfCare{
						ec,
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.OpenEpisodesFn = func(
					ctx context.Context,
					patientReference string,
				) ([]*domain.FHIREpisodeOfCare, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.OpenEpisodes(tt.args.ctx, tt.args.patientReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.OpenEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUnit_HasOpenEpisode(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	patient := domain.FHIRPatient{}

	type args struct {
		ctx     context.Context
		patient domain.FHIRPatient
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{

		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				patient: patient,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				patient: patient,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.HasOpenEpisodeFn = func(ctx context.Context, patient domain.FHIRPatient) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.HasOpenEpisodeFn = func(ctx context.Context, patient domain.FHIRPatient) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			got, err := fh.HasOpenEpisode(tt.args.ctx, tt.args.patient)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.HasOpenEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FHIRUseCaseImpl.HasOpenEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_CreateFHIREncounter(t *testing.T) {

	ctx := context.Background()
	fh := fakeUsecaseIntr

	input := domain.FHIREncounterInput{}
	type args struct {
		ctx   context.Context
		input domain.FHIREncounterInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIREncounterRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateFHIREncounterFn = func(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error) {
					id := ksuid.New().String()
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							ID: &id,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateFHIREncounterFn = func(
					ctx context.Context,
					input domain.FHIREncounterInput,
				) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.CreateFHIREncounter(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestUnit_GetFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	id := ksuid.New().String()

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIREpisodeOfCareRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.GetFHIREpisodeOfCareFn = func(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error) {
					ID := ksuid.New().String()
					return &domain.FHIREpisodeOfCareRelayPayload{
						Resource: &domain.FHIREpisodeOfCare{
							ID: &ID,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.GetFHIREpisodeOfCareFn = func(
					ctx context.Context,
					id string,
				) (*domain.FHIREpisodeOfCareRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.GetFHIREpisodeOfCare(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_StartEncounter_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	episodeID := ksuid.New().String()

	type args struct {
		ctx       context.Context
		episodeID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				episodeID: episodeID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				episodeID: episodeID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case" {
				fakeFhir.StartEncounterFn = func(ctx context.Context, episodeID string) (string, error) {
					return "test-encounter", nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.StartEncounterFn = func(ctx context.Context, episodeID string) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.StartEncounter(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUnit_UpdateFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr
	input := domain.FHIRAllergyIntoleranceInput{}

	type args struct {
		ctx   context.Context
		input domain.FHIRAllergyIntoleranceInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRAllergyIntoleranceRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case" {
				fakeFhir.UpdateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					id := ksuid.New().String()
					return &domain.FHIRAllergyIntoleranceRelayPayload{
						Resource: &domain.FHIRAllergyIntolerance{
							ID: &id,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.UpdateFHIRAllergyIntoleranceFn = func(
					ctx context.Context,
					input domain.FHIRAllergyIntoleranceInput,
				) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.UpdateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUnit_SearchFHIRComposition(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	params := map[string]interface{}{"test": "123"}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRCompositionRelayConnection
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				params: params,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				params: params,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.SearchFHIRCompositionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
					return &domain.FHIRCompositionRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchFHIRCompositionFn = func(
					ctx context.Context,
					params map[string]interface{},
				) (*domain.FHIRCompositionRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.SearchFHIRComposition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUnit_CreateFHIRComposition(t *testing.T) {

	ctx := context.Background()
	fh := fakeUsecaseIntr

	input := domain.FHIRCompositionInput{}

	type args struct {
		ctx   context.Context
		input domain.FHIRCompositionInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRCompositionRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case" {
				fakeFhir.CreateFHIRCompositionFn = func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
					id := ksuid.New().String()
					return &domain.FHIRCompositionRelayPayload{
						Resource: &domain.FHIRComposition{
							ID: &id,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateFHIRCompositionFn = func(
					ctx context.Context,
					input domain.FHIRCompositionInput,
				) (*domain.FHIRCompositionRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.CreateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUnit_UpdateFHIRComposition(t *testing.T) {

	ctx := context.Background()
	fh := fakeUsecaseIntr

	input := domain.FHIRCompositionInput{}

	type args struct {
		ctx   context.Context
		input domain.FHIRCompositionInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRCompositionRelayPayload
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.UpdateFHIRCompositionFn = func(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
					id := ksuid.New().String()
					return &domain.FHIRCompositionRelayPayload{
						Resource: &domain.FHIRComposition{
							ID: &id,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.UpdateFHIRCompositionFn = func(
					ctx context.Context,
					input domain.FHIRCompositionInput,
				) (*domain.FHIRCompositionRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.UpdateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUnit_DeleteFHIRComposition(t *testing.T) {

	ctx := context.Background()
	fh := fakeUsecaseIntr

	id := ksuid.New().String()

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
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.DeleteFHIRCompositionFn = func(ctx context.Context, id string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.DeleteFHIRCompositionFn = func(ctx context.Context, id string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.DeleteFHIRComposition(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRComposition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRUseCaseImpl_StartEpisodeByOtp_Unittest(t *testing.T) {
	ctx := context.Background()
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
				ctx:   ctx,
				input: *validEpisodeInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: *invalidEpisodeInput,
			},
			wantErr: true,
		},

		{
			name: "Sad case#1",
			args: args{
				ctx:   ctx,
				input: *invalidEpisodeInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.StartEpisodeByOtpFn = func(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
					return &domain.EpisodeOfCarePayload{
						TotalVisits: 1,
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.StartEpisodeByOtpFn = func(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case#1" {
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
	ctx := context.Background()
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
				ctx: ctx,
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
				ctx: ctx,
				input: domain.OTPEpisodeUpgradeInput{
					EpisodeID: "",
					Msisdn:    "+254711223344",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
				input: domain.OTPEpisodeUpgradeInput{
					EpisodeID: ksuid.New().String(),
					Msisdn:    "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.UpgradeEpisodeFn = func(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
					return &domain.EpisodeOfCarePayload{
						TotalVisits: 1,
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.UpgradeEpisodeFn = func(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case#1" {
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
		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.SearchEpisodeEncounterFn = func(ctx context.Context, episodeReference string) (*domain.FHIREncounterRelayConnection, error) {
					return &domain.FHIREncounterRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchEpisodeEncounterFn = func(ctx context.Context, episodeReference string) (*domain.FHIREncounterRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case#1" {
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
		{
			name: "Sad case#1",
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
				fakeFhir.EndEncounterFn = func(ctx context.Context, encounterID string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.EndEncounterFn = func(ctx context.Context, encounterID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case#1" {
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

		{
			name: "Sad case#1",
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
				fakeFhir.EndEpisodeFn = func(ctx context.Context, episodeID string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.EndEpisodeFn = func(ctx context.Context, episodeID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case#1" {
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
				episodeID: "",
			},
			wantErr: true,
		},
		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.GetActiveEpisodeFn = func(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
					id := ksuid.New().String()
					return &domain.FHIREpisodeOfCare{
						ID: &id,
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.GetActiveEpisodeFn = func(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case#1" {
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

		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.SearchFHIRServiceRequestFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
					return &domain.FHIRServiceRequestRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchFHIRServiceRequestFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case#1" {
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
		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateFHIRServiceRequestFn = func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
					id := ksuid.New().String()
					return &domain.FHIRServiceRequestRelayPayload{
						Resource: &domain.FHIRServiceRequest{
							ID: &id,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateFHIRServiceRequestFn = func(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case#1" {
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
		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.SearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return &domain.FHIRAllergyIntoleranceRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchFHIRAllergyIntoleranceFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case#1" {
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
		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					id := ksuid.New().String()
					return &domain.FHIRAllergyIntoleranceRelayPayload{
						Resource: &domain.FHIRAllergyIntolerance{
							ID: &id,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case#1" {
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

func TestFHIRUseCaseImpl_SearchFHIRCondition_Unittest(t *testing.T) {
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
				fakeFhir.SearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
					return &domain.FHIRConditionRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case#1" {
				fakeFhir.SearchFHIRConditionFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.SearchFHIRCondition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_UpdateFHIRCondition_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	id := ksuid.New().String()

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
				ctx: ctx,
				input: domain.FHIRConditionInput{
					ID:                 &id,
					Identifier:         []*domain.FHIRIdentifierInput{},
					ClinicalStatus:     &domain.FHIRCodeableConceptInput{},
					VerificationStatus: &domain.FHIRCodeableConceptInput{},
					Category:           []*domain.FHIRCodeableConceptInput{},
					Severity:           &domain.FHIRCodeableConceptInput{},
					Code:               &domain.FHIRCodeableConceptInput{},
					BodySite:           []*domain.FHIRCodeableConceptInput{},
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: domain.FHIRConditionInput{
					ID:                 &id,
					Identifier:         []*domain.FHIRIdentifierInput{},
					ClinicalStatus:     &domain.FHIRCodeableConceptInput{},
					VerificationStatus: &domain.FHIRCodeableConceptInput{},
					Category:           []*domain.FHIRCodeableConceptInput{},
					Severity:           &domain.FHIRCodeableConceptInput{},
					Code:               &domain.FHIRCodeableConceptInput{},
					BodySite:           []*domain.FHIRCodeableConceptInput{},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.UpdateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					id := ksuid.New().String()
					return &domain.FHIRConditionRelayPayload{
						Resource: &domain.FHIRCondition{
							ID: &id,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.UpdateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case#1" {
				fakeFhir.UpdateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.UpdateFHIRCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_GetFHIREncounter_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	id := ksuid.New().String()

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
				id:  id,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				id:  "",
			},
			wantErr: true,
		},

		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.GetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					Id := ksuid.New().String()
					return &domain.FHIREncounterRelayPayload{
						Resource: &domain.FHIREncounter{
							ID: &Id,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.GetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case#1" {
				fakeFhir.GetFHIREncounterFn = func(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.GetFHIREncounter(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestFHIRUseCaseImpl_SearchFHIREncounter_Unittest(t *testing.T) {
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
		{
			name: "Sad case#1",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.SearchFHIREncounterFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
					return &domain.FHIREncounterRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchFHIREncounterFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case#1" {
				fakeFhir.SearchFHIREncounterFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.SearchFHIREncounter(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_SearchFHIRMedicationRequest_Unittest(t *testing.T) {
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
				fakeFhir.SearchFHIRMedicationRequestFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
					return &domain.FHIRMedicationRequestRelayConnection{
						PageInfo: &firebasetools.PageInfo{
							HasNextPage: true,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.SearchFHIRMedicationRequestFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := fh.SearchFHIRMedicationRequest(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateFHIRMedicationRequest_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	id := ksuid.New().String()

	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationRequestInput
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
				input: domain.FHIRMedicationRequestInput{
					ID:           &id,
					Identifier:   []*domain.FHIRIdentifierInput{},
					StatusReason: &domain.FHIRCodeableConceptInput{},
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: domain.FHIRMedicationRequestInput{
					ID:           &id,
					Identifier:   []*domain.FHIRIdentifierInput{},
					StatusReason: &domain.FHIRCodeableConceptInput{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.CreateFHIRMedicationRequestFn = func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
					ID := ksuid.New().String()
					return &domain.FHIRMedicationRequestRelayPayload{
						Resource: &domain.FHIRMedicationRequest{
							ID: &ID,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.CreateFHIRMedicationRequestFn = func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.CreateFHIRMedicationRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_UpdateFHIRMedicationRequest_Unittest(t *testing.T) {
	ctx := context.Background()
	fh := fakeUsecaseIntr

	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationRequestInput
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
				input: domain.FHIRMedicationRequestInput{
					ID:           new(string),
					Identifier:   []*domain.FHIRIdentifierInput{},
					StatusReason: &domain.FHIRCodeableConceptInput{},
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: domain.FHIRMedicationRequestInput{
					ID:           new(string),
					Identifier:   []*domain.FHIRIdentifierInput{},
					StatusReason: &domain.FHIRCodeableConceptInput{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFhir.UpdateFHIRMedicationRequestFn = func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
					ID := ksuid.New().String()
					return &domain.FHIRMedicationRequestRelayPayload{
						Resource: &domain.FHIRMedicationRequest{
							ID: &ID,
						},
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFhir.UpdateFHIRMedicationRequestFn = func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.UpdateFHIRMedicationRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
