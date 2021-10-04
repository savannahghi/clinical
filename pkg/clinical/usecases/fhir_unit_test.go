package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/segmentio/ksuid"
)

func TestFHIRUseCaseImpl_CreateEpisodeOfCare_Unittest(t *testing.T) {
	_, token := firebasetools.GetAuthenticatedContextAndToken(t)
	fmt.Println("token", token.UID)
	ctx := context.Background()
	ctx = context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		token,
	)

	fh, err := InitializeFakeClinicalInteractor(context.Background())
	if err != nil {
		t.Errorf("failed to initialize fake test service: %v", err)
	}

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
				fakeProfileUtils.GetLoggedInUserFn = func(ctx context.Context) (*profileutils.UserInfo, error) {
					return &profileutils.UserInfo{
						DisplayName: "test user",
						Email:       firebasetools.TestUserEmail,
						PhoneNumber: interserviceclient.TestUserPhoneNumber,
						PhotoURL:    gofakeit.ImageURL(50, 50),
						ProviderID:  token.Firebase.SignInProvider,
						UID:         token.UID,
					}, nil
				}
				// fakeFhir.CreateEpisodeOfCareFn = func(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
				// 	return &domain.EpisodeOfCarePayload{}, nil
				// }
				fakeFhir.SearchFHIREpisodeOfCareFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error) {
					return &domain.FHIREpisodeOfCareRelayConnection{}, nil
				}

				fakeFhirRepo.CreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}) ([]byte, error) {
					return []byte("test"), nil
				}
				fakeFhir.EncountersFn = func(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum) ([]*domain.FHIREncounter, error) {
					return []*domain.FHIREncounter{}, nil
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

// func TestFHIRUseCaseImpl_CreateFHIRCondition_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	ID := ksuid.New().String()

// 	fhirconditionInput := &domain.FHIRConditionInput{
// 		ID:                 &ID,
// 		Identifier:         []*domain.FHIRIdentifierInput{},
// 		ClinicalStatus:     &domain.FHIRCodeableConceptInput{},
// 		VerificationStatus: &domain.FHIRCodeableConceptInput{},
// 		Category:           []*domain.FHIRCodeableConceptInput{},
// 		Severity:           &domain.FHIRCodeableConceptInput{},
// 		Code:               &domain.FHIRCodeableConceptInput{},
// 		BodySite:           []*domain.FHIRCodeableConceptInput{},
// 		Subject:            &domain.FHIRReferenceInput{},
// 		Encounter:          &domain.FHIRReferenceInput{},
// 		OnsetDateTime:      &scalarutils.Date{},
// 		OnsetAge:           &domain.FHIRAgeInput{},
// 		OnsetPeriod:        &domain.FHIRPeriodInput{},
// 		OnsetRange:         &domain.FHIRRangeInput{},
// 		OnsetString:        new(string),
// 		AbatementDateTime:  &scalarutils.Date{},
// 		AbatementAge:       &domain.FHIRAgeInput{},
// 		AbatementPeriod:    &domain.FHIRPeriodInput{},
// 		AbatementRange:     &domain.FHIRRangeInput{},
// 		AbatementString:    new(string),
// 		RecordedDate:       &scalarutils.Date{},
// 		Recorder:           &domain.FHIRReferenceInput{},
// 		Asserter:           &domain.FHIRReferenceInput{},
// 		Stage:              []*domain.FHIRConditionStageInput{},
// 		Evidence:           []*domain.FHIRConditionEvidenceInput{},
// 		Note:               []*domain.FHIRAnnotationInput{},
// 	}

// 	invalidfhirconditionInput := &domain.FHIRConditionInput{
// 		ID:                 nil,
// 		Identifier:         []*domain.FHIRIdentifierInput{},
// 		ClinicalStatus:     &domain.FHIRCodeableConceptInput{},
// 		VerificationStatus: &domain.FHIRCodeableConceptInput{},
// 		Category:           []*domain.FHIRCodeableConceptInput{},
// 		Severity:           &domain.FHIRCodeableConceptInput{},
// 		Code:               &domain.FHIRCodeableConceptInput{},
// 	}

// 	type args struct {
// 		ctx   context.Context
// 		input domain.FHIRConditionInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:   ctx,
// 				input: *fhirconditionInput,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:   ctx,
// 				input: *invalidfhirconditionInput,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.CreateFHIRConditionFn = usecaseMock.NewFHIRMock().CreateFHIRCondition
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.CreateFHIRConditionFn = func(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.CreateFHIRCondition(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.CreateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestFHIRUseCaseImpl_OpenOrganizationEpisodes_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	type args struct {
// 		ctx               context.Context
// 		providerSladeCode string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:               ctx,
// 				providerSladeCode: "1234",
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:               ctx,
// 				providerSladeCode: "",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.OpenOrganizationEpisodesFn = usecaseMock.NewFHIRMock().OpenOrganizationEpisodes
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.OpenOrganizationEpisodesFn = func(ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.OpenOrganizationEpisodes(tt.args.ctx, tt.args.providerSladeCode)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.OpenOrganizationEpisodes() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestFHIRUseCaseImpl_GetORCreateOrganization_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	type args struct {
// 		ctx               context.Context
// 		providerSladeCode string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:               ctx,
// 				providerSladeCode: "1234",
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:               ctx,
// 				providerSladeCode: "",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.GetORCreateOrganizationFn = usecaseMock.NewFHIRMock().CreateOrganization
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.GetORCreateOrganizationFn = func(ctx context.Context, providerSladeCode string) (*string, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			_, err := fh.GetORCreateOrganization(tt.args.ctx, tt.args.providerSladeCode)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.GetORCreateOrganization() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestFHIRUseCaseImpl_CreateFHIROrganization_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	ID := ksuid.New().String()
// 	active := true
// 	testname := gofakeit.FirstName()

// 	orgInput := &domain.FHIROrganizationInput{
// 		ID:         &ID,
// 		Active:     &active,
// 		Identifier: []*domain.FHIRIdentifierInput{},
// 		Type:       []*domain.FHIRCodeableConceptInput{},
// 		Name:       &testname,
// 		Alias:      []string{"alias test"},
// 		Telecom:    []*domain.FHIRContactPointInput{},
// 		Address:    []*domain.FHIRAddressInput{},
// 	}

// 	invalidOrgInput := &domain.FHIROrganizationInput{
// 		ID:         &ID,
// 		Active:     new(bool),
// 		Identifier: []*domain.FHIRIdentifierInput{},
// 		Type:       []*domain.FHIRCodeableConceptInput{},
// 		Alias:      []string{"alias test"},
// 		Telecom:    []*domain.FHIRContactPointInput{},
// 		Address:    []*domain.FHIRAddressInput{},
// 	}

// 	type args struct {
// 		ctx   context.Context
// 		input domain.FHIROrganizationInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:   ctx,
// 				input: *orgInput,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:   ctx,
// 				input: *invalidOrgInput,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.CreateFHIROrganizationFn = usecaseMock.NewFHIRMock().CreateFHIROrganization
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.CreateFHIROrganizationFn = func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.CreateFHIROrganization(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.CreateFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestFHIRUseCaseImpl_CreateOrganization_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	type args struct {
// 		ctx               context.Context
// 		providerSladeCode string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:               ctx,
// 				providerSladeCode: "1234",
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:               ctx,
// 				providerSladeCode: "",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.CreateOrganizationFn = usecaseMock.NewFHIRMock().CreateOrganization
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.CreateOrganizationFn = func(ctx context.Context, providerSladeCode string) (*string, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.CreateOrganization(tt.args.ctx, tt.args.providerSladeCode)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.CreateOrganization() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestFHIRUseCaseImpl_SearchFHIROrganization_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	type args struct {
// 		ctx    context.Context
// 		params map[string]interface{}
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
// 				params: map[string]interface{}{
// 					"test": "test",
// 				},
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:    ctx,
// 				params: map[string]interface{}{},
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.SearchFHIROrganizationFn = usecaseMock.NewFHIRMock().SearchFHIROrganization
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.SearchFHIROrganizationFn = func(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.SearchFHIROrganization(tt.args.ctx, tt.args.params)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.SearchFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestFHIRUseCaseImpl_GetOrganization_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	type args struct {
// 		ctx               context.Context
// 		providerSladeCode string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:               ctx,
// 				providerSladeCode: "1234",
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:               ctx,
// 				providerSladeCode: "",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.GetOrganizationFn = usecaseMock.NewFHIRMock().GetOrganization
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.GetOrganizationFn = func(ctx context.Context, providerSladeCode string) (*string, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.GetOrganization(tt.args.ctx, tt.args.providerSladeCode)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.GetOrganization() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestFHIRUseCaseImpl_SearchEpisodesByParam_Unittest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	params := url.Values{}

// 	type args struct {
// 		ctx          context.Context
// 		searchParams url.Values
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:          ctx,
// 				searchParams: params,
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
// 				fakeFhir.SearchEpisodesByParamFn = usecaseMock.NewFHIRMock().SearchEpisodesByParam
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.SearchEpisodesByParamFn = func(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.SearchEpisodesByParam(tt.args.ctx, tt.args.searchParams)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.SearchEpisodesByParam() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestFHIRUseCaseImpl_POSTRequest_Unittest(t *testing.T) {
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	searchParams := url.Values{}

// 	type args struct {
// 		resourceName string
// 		path         string
// 		params       url.Values
// 		body         io.Reader
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				resourceName: "Encounter",
// 				path:         "_search",
// 				params:       searchParams,
// 				body:         nil,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				resourceName: "Encounter",
// 				path:         "",
// 				params:       searchParams,
// 				body:         nil,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.POSTRequestFn = usecaseMock.NewFHIRMock().POSTRequest
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.POSTRequestFn = func(resourceName, path string, params url.Values, body io.Reader) ([]byte, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.POSTRequest(tt.args.resourceName, tt.args.path, tt.args.params, tt.args.body)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.POSTRequest() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_OpenEpisodes(t *testing.T) {
// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}
// 	patientReference := fmt.Sprintf("Patient/%s", ksuid.New().String())

// 	type args struct {
// 		ctx              context.Context
// 		patientReference string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:              ctx,
// 				patientReference: patientReference,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:              ctx,
// 				patientReference: patientReference,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.OpenEpisodesFn = usecaseMock.NewFHIRMock().OpenEpisodesFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.OpenEpisodesFn = func(
// 					ctx context.Context,
// 					patientReference string,
// 				) ([]*domain.FHIREpisodeOfCare, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.OpenEpisodes(tt.args.ctx, tt.args.patientReference)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.OpenEpisodes() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_HasOpenEpisode(t *testing.T) {
// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	patient := domain.FHIRPatient{}

// 	type args struct {
// 		ctx     context.Context
// 		patient domain.FHIRPatient
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{

// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:     ctx,
// 				patient: patient,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:     ctx,
// 				patient: patient,
// 			},
// 			wantErr: true,
// 			want:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.HasOpenEpisodeFn = usecaseMock.NewFHIRMock().HasOpenEpisodeFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.HasOpenEpisodeFn = func(ctx context.Context, patient domain.FHIRPatient) (bool, error) {
// 					return false, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			got, err := fh.HasOpenEpisode(tt.args.ctx, tt.args.patient)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.HasOpenEpisode() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("FHIRUseCaseImpl.HasOpenEpisode() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestUnit_CreateFHIREncounter(t *testing.T) {

// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	input := domain.FHIREncounterInput{}
// 	type args struct {
// 		ctx   context.Context
// 		input domain.FHIREncounterInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *domain.FHIREncounterRelayPayload
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:   ctx,
// 				input: input,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:   ctx,
// 				input: input,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.CreateFHIREncounterFn = usecaseMock.NewFHIRMock().CreateFHIREncounterFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.CreateFHIREncounterFn = func(
// 					ctx context.Context,
// 					input domain.FHIREncounterInput,
// 				) (*domain.FHIREncounterRelayPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.CreateFHIREncounter(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.CreateFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 		})
// 	}
// }

// func TestUnit_GetFHIREpisodeOfCare(t *testing.T) {
// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	id := ksuid.New().String()

// 	type args struct {
// 		ctx context.Context
// 		id  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *domain.FHIREpisodeOfCareRelayPayload
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				id:  id,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 				id:  id,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.GetFHIREpisodeOfCareFn = usecaseMock.NewFHIRMock().GetFHIREpisodeOfCareFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.GetFHIREpisodeOfCareFn = func(
// 					ctx context.Context,
// 					id string,
// 				) (*domain.FHIREpisodeOfCareRelayPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.GetFHIREpisodeOfCare(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.GetFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 		})
// 	}
// }

// func TestClinicalUseCaseImpl_StartEncounter(t *testing.T) {
// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	episodeID := ksuid.New().String()

// 	type args struct {
// 		ctx       context.Context
// 		episodeID string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:       ctx,
// 				episodeID: episodeID,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:       ctx,
// 				episodeID: episodeID,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Happy case" {
// 				fakeFhir.StartEncounterFn = usecaseMock.NewFHIRMock().StartEncounterFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.StartEncounterFn = func(ctx context.Context, episodeID string) (string, error) {
// 					return "", fmt.Errorf("an error occurred")
// 				}
// 			}

// 			_, err := fh.StartEncounter(tt.args.ctx, tt.args.episodeID)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ClinicalUseCaseImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_UpdateFHIRAllergyIntolerance(t *testing.T) {
// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}
// 	input := domain.FHIRAllergyIntoleranceInput{}

// 	type args struct {
// 		ctx   context.Context
// 		input domain.FHIRAllergyIntoleranceInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *domain.FHIRAllergyIntoleranceRelayPayload
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:   ctx,
// 				input: input,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:   ctx,
// 				input: input,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Happy case" {
// 				fakeFhir.UpdateFHIRAllergyIntoleranceFn = usecaseMock.NewFHIRMock().UpdateFHIRAllergyIntoleranceFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.UpdateFHIRAllergyIntoleranceFn = func(
// 					ctx context.Context,
// 					input domain.FHIRAllergyIntoleranceInput,
// 				) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			_, err := fh.UpdateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.UpdateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_SearchFHIRComposition(t *testing.T) {
// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	params := map[string]interface{}{"test": "123"}

// 	type args struct {
// 		ctx    context.Context
// 		params map[string]interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *domain.FHIRCompositionRelayConnection
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:    ctx,
// 				params: params,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:    ctx,
// 				params: params,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.SearchFHIRCompositionFn = usecaseMock.NewFHIRMock().SearchFHIRCompositionFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.SearchFHIRCompositionFn = func(
// 					ctx context.Context,
// 					params map[string]interface{},
// 				) (*domain.FHIRCompositionRelayConnection, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			_, err := fh.SearchFHIRComposition(tt.args.ctx, tt.args.params)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.SearchFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_CreateFHIRComposition(t *testing.T) {

// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	input := domain.FHIRCompositionInput{}

// 	type args struct {
// 		ctx   context.Context
// 		input domain.FHIRCompositionInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *domain.FHIRCompositionRelayPayload
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:   ctx,
// 				input: input,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:   ctx,
// 				input: input,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Happy case" {
// 				fakeFhir.CreateFHIRCompositionFn = usecaseMock.NewFHIRMock().CreateFHIRCompositionFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.CreateFHIRCompositionFn = func(
// 					ctx context.Context,
// 					input domain.FHIRCompositionInput,
// 				) (*domain.FHIRCompositionRelayPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			_, err := fh.CreateFHIRComposition(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.CreateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_UpdateFHIRComposition(t *testing.T) {

// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	input := domain.FHIRCompositionInput{}

// 	type args struct {
// 		ctx   context.Context
// 		input domain.FHIRCompositionInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *domain.FHIRCompositionRelayPayload
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:   ctx,
// 				input: input,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx:   ctx,
// 				input: input,
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.UpdateFHIRCompositionFn = usecaseMock.NewFHIRMock().UpdateFHIRCompositionFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.UpdateFHIRCompositionFn = func(
// 					ctx context.Context,
// 					input domain.FHIRCompositionInput,
// 				) (*domain.FHIRCompositionRelayPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.UpdateFHIRComposition(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.UpdateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_DeleteFHIRComposition(t *testing.T) {

// 	ctx := context.Background()
// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	id := ksuid.New().String()

// 	type args struct {
// 		ctx context.Context
// 		id  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				id:  id,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},
// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 				id:  id,
// 			},
// 			wantErr: true,
// 			want:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.DeleteFHIRCompositionFn = usecaseMock.NewFHIRMock().DeleteFHIRCompositionFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.DeleteFHIRCompositionFn = func(ctx context.Context, id string) (bool, error) {
// 					return false, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			got, err := fh.DeleteFHIRComposition(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRComposition() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestUnit_DeleteFHIRMedicationRequest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	id := ksuid.New().String()

// 	type args struct {
// 		ctx context.Context
// 		id  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				id:  id,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 			want:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				// fakePatient.CreatePatientFn = func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
// 				// 	return &domain.PatientPayload{}, nil
// 				// }
// 				// fakeFhir.CreateEpisodeOfCareFn = func(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
// 				// 	return &domain.EpisodeOfCarePayload{}, nil
// 				// }

// 				// fakeFhir.CreateFHIRMedicationRequestFn = func(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
// 				// 	return &domain.FHIRMedicationRequestRelayPayload{}, nil
// 				// }
// 				fakeFhir.DeleteFHIRMedicationRequestFn = usecaseMock.NewFHIRMock().DeleteFHIRMedicationRequestFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.DeleteFHIRMedicationRequestFn = func(ctx context.Context, id string) (bool, error) {
// 					return false, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			got, err := fh.DeleteFHIRMedicationRequest(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRMedicationRequest() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestUnit_SearchFHIRObservation(t *testing.T) {

// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	params := map[string]interface{}{"test": "123"}

// 	type args struct {
// 		ctx    context.Context
// 		params map[string]interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *domain.FHIRObservationRelayConnection
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:    ctx,
// 				params: params,
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
// 				fakeFhir.SearchFHIRObservationFn = usecaseMock.NewFHIRMock().SearchFHIRObservationFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.SearchFHIRObservationFn = func(
// 					ctx context.Context,
// 					params map[string]interface{},
// 				) (*domain.FHIRObservationRelayConnection, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.SearchFHIRObservation(tt.args.ctx, tt.args.params)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.SearchFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_CreateFHIRObservation(t *testing.T) {

// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	input := domain.FHIRObservationInput{}

// 	type args struct {
// 		ctx   context.Context
// 		input domain.FHIRObservationInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *domain.FHIRObservationRelayPayload
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx:   ctx,
// 				input: input,
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
// 				fakeFhir.CreateFHIRObservationFn = usecaseMock.NewFHIRMock().CreateFHIRObservationFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.CreateFHIRObservationFn = func(
// 					ctx context.Context,
// 					input domain.FHIRObservationInput,
// 				) (*domain.FHIRObservationRelayPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.CreateFHIRObservation(tt.args.ctx, tt.args.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.CreateFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_DeleteFHIRObservation(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	id := ksuid.New().String()

// 	type args struct {
// 		ctx context.Context
// 		id  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				id:  id,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 			want:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.DeleteFHIRObservationFn = usecaseMock.NewFHIRMock().DeleteFHIRObservationFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.DeleteFHIRObservationFn = func(ctx context.Context, id string) (bool, error) {
// 					return false, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			got, err := fh.DeleteFHIRObservation(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRObservation() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestUnit_GetFHIRPatient(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	id := ksuid.New().String()

// 	type args struct {
// 		ctx context.Context
// 		id  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *domain.FHIRPatientRelayPayload
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				id:  id,
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
// 				fakeFhir.GetFHIRPatientFn = usecaseMock.NewFHIRMock().GetFHIRPatientFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.GetFHIRPatientFn = func(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error) {
// 					return nil, fmt.Errorf("an error occurred")
// 				}
// 			}
// 			_, err := fh.GetFHIRPatient(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.GetFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestUnit_DeleteFHIRPatient(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	id := ksuid.New().String()

// 	type args struct {
// 		ctx context.Context
// 		id  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				id:  id,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 			want:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.DeleteFHIRPatientFn = usecaseMock.NewFHIRMock().DeleteFHIRPatientFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.DeleteFHIRPatientFn = func(ctx context.Context, id string) (bool, error) {
// 					return false, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			got, err := fh.DeleteFHIRPatient(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRPatient() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestUnit_DeleteFHIRResourceType(t *testing.T) {

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	results := []map[string]string{{"service-request": ksuid.New().String()}}

// 	type args struct {
// 		results []map[string]string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				results: results,
// 			},
// 			wantErr: false,
// 		},

// 		{
// 			name:    "Sad case",
// 			args:    args{},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.DeleteFHIRResourceTypeFn = usecaseMock.NewFHIRMock().DeleteFHIRResourceTypeFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.DeleteFHIRResourceTypeFn = func(results []map[string]string) error {
// 					return fmt.Errorf("an error occurred")
// 				}
// 			}
// 			if err := fh.DeleteFHIRResourceType(tt.args.results); (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRResourceType() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestUnit_DeleteFHIRServiceRequest(t *testing.T) {
// 	ctx := context.Background()

// 	fh, err := InitializeFakeClinicalInteractor(context.Background())
// 	if err != nil {
// 		t.Errorf("failed to initialize fake test service: %v", err)
// 	}

// 	id := ksuid.New().String()

// 	type args struct {
// 		ctx context.Context
// 		id  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy case",
// 			args: args{
// 				ctx: ctx,
// 				id:  id,
// 			},
// 			wantErr: false,
// 			want:    true,
// 		},

// 		{
// 			name: "Sad case",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			wantErr: true,
// 			want:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.name == "Happy case" {
// 				fakeFhir.DeleteFHIRServiceRequestFn = usecaseMock.NewFHIRMock().DeleteFHIRServiceRequestFn
// 			}

// 			if tt.name == "Sad case" {
// 				fakeFhir.DeleteFHIRServiceRequestFn = func(ctx context.Context, id string) (bool, error) {
// 					return false, fmt.Errorf("an error occurred")
// 				}
// 			}

// 			got, err := fh.DeleteFHIRServiceRequest(tt.args.ctx, tt.args.id)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("FHIRUseCaseImpl.DeleteFHIRServiceRequest() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
