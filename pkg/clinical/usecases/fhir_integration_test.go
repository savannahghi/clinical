package usecases_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/interserviceclient"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	testName = "test"
)

func TestFHIRUseCaseImpl_Encounters(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	_, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v", err)
		return
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.ID)

	status := domain.EncounterStatusEnumArrived

	type args struct {
		ctx              context.Context
		patientReference string
		status           *domain.EncounterStatusEnum
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.FHIREncounter
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:              ctx,
				patientReference: patientRef,
				status:           &status,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:              context.Background(),
				patientReference: patientRef,
				status:           &status,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.Encounters(tt.args.ctx, tt.args.patientReference, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.Encounters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_SearchFHIREpisodeOfCare(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	params := map[string]interface{}{}

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
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIREpisodeOfCare(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateEpisodeOfCare(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := getTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	type args struct {
		ctx     context.Context
		episode domain.FHIREpisodeOfCare
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.EpisodeOfCarePayload
		wantErr bool
		panics  bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:     ctx,
				episode: *episode,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:     context.Background(),
				episode: *episode,
			},
			wantErr: true,
		},
		{
			name: "invalid: empty episode",
			args: args{
				ctx:     ctx,
				episode: domain.FHIREpisodeOfCare{},
			},
			panics: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		if !tt.panics {
			t.Run(tt.name, func(t *testing.T) {
				_, err := fh.CreateEpisodeOfCare(tt.args.ctx, tt.args.episode)
				if (err != nil) != tt.wantErr {
					t.Errorf("FHIRUseCaseImpl.CreateEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		}
		if tt.panics {
			fc := func() { _, _ = fh.CreateEpisodeOfCare(tt.args.ctx, tt.args.episode) }
			assert.Panics(t, fc)
		}
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_CreateFHIRCondition(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
	}

	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
	}

	encounter, err := fh.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to create FHIREncounter: %v", err)
	}

	input, err := createTestConditionInput(*encounter.Resource.ID, *patient.ID)
	if err != nil {
		fmt.Printf("cant create condition: %v\n", err)
		return
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRConditionInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRConditionRelayPayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:   context.Background(),
				input: *input,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.CreateFHIRCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_CreateFHIROrganization(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	identifier := []*domain.FHIRIdentifierInput{
		{
			Use:   "official",
			Value: testProviderCode,
		},
	}

	input := domain.FHIROrganizationInput{
		Identifier: identifier,
		Name:       &testName,
	}

	type args struct {
		ctx   context.Context
		input domain.FHIROrganizationInput
	}
	tests := []struct {
		name string

		args    args
		want    *domain.FHIROrganizationRelayPayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.CreateFHIROrganization(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_OpenOrganizationEpisodes(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
	}
	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
	}

	_, err = fh.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to create FHIREncounter: %v", err)
	}

	type args struct {
		ctx               context.Context
		providerSladeCode string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.FHIREpisodeOfCare
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:               ctx,
				providerSladeCode: testProviderCode,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:               context.Background(),
				providerSladeCode: testProviderCode,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.OpenOrganizationEpisodes(tt.args.ctx, tt.args.providerSladeCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.OpenOrganizationEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_GetORCreateOrganization(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	type args struct {
		ctx               context.Context
		providerSladeCode string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:               ctx,
				providerSladeCode: testProviderCode,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:               context.Background(),
				providerSladeCode: testProviderCode,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.GetORCreateOrganization(tt.args.ctx, tt.args.providerSladeCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetORCreateOrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateOrganization(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	type args struct {
		ctx               context.Context
		providerSladeCode string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:               ctx,
				providerSladeCode: testProviderCode,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:               context.Background(),
				providerSladeCode: testProviderCode,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.CreateOrganization(tt.args.ctx, tt.args.providerSladeCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateOrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_SearchFHIROrganization(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	// Create FHIR organization
	identifier := []*domain.FHIRIdentifierInput{
		{
			Use:   "official",
			Value: testProviderCode,
		},
	}

	input := domain.FHIROrganizationInput{
		Identifier: identifier,
		Name:       &testName,
	}

	_, err = fh.CreateFHIROrganization(ctx, input)
	if err != nil {
		t.Errorf("failed to create fhir organization: %v", err)
	}

	params := map[string]interface{}{"provider": "123"}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIROrganizationRelayConnection
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing params",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIROrganization(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

func TestFHIRUseCaseImpl_FindOrganizationByID(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	organizationID := "a1bf993c-c2b6-44dd-8991-3a47b54a6789"

	// Create FHIR organization
	identifier := []*domain.FHIRIdentifierInput{
		{
			Use:   "official",
			Value: testProviderCode,
		},
	}

	input := domain.FHIROrganizationInput{
		ID:         &organizationID,
		Identifier: identifier,
		Name:       &testName,
	}

	organization, err := fh.CreateFHIROrganization(ctx, input)
	if err != nil {
		t.Errorf("failed to create fhir organization: %v", err)
	}

	type args struct {
		ctx            context.Context
		organizationID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:            ctx,
				organizationID: *organization.Resource.ID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:            ctx,
				organizationID: "",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid ID",
			args: args{
				ctx:            ctx,
				organizationID: "testID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fh.FindOrganizationByID(tt.args.ctx, tt.args.organizationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.FindOrganizationByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected response to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected response not to be nil for %v", tt.name)
				return
			}
		})
	}

	//Clean up
	resource := []map[string]string{
		{"resourceType": "Organization", "resourceID": *organization.Resource.ID},
	}
	err = fh.DeleteFHIRResourceType(resource)
	if err != nil {
		t.Errorf("failed to delete fhir organization: %v", err)
	}
}

func TestFHIRUseCaseImpl_GetOrganization(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	type args struct {
		ctx               context.Context
		providerSladeCode string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:               ctx,
				providerSladeCode: testProviderCode,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:               context.Background(),
				providerSladeCode: testProviderCode,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.GetOrganization(tt.args.ctx, tt.args.providerSladeCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetOrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_SearchEpisodesByParam(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	_, _, err = createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	searchParams := url.Values{"": []string{""}}
	type args struct {
		ctx          context.Context
		searchParams url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.FHIREpisodeOfCare
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:          ctx,
				searchParams: searchParams,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchEpisodesByParam(tt.args.ctx, tt.args.searchParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchEpisodesByParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

// func TestFHIRUseCaseImpl_POSTRequest(t *testing.T) {
// ctx, err := getTestAuthenticatedContext(t)
// if err != nil {
// 	t.Errorf("cant get phone number authenticated context token: %v", err)
// 	return
// }
// 	fh := testUsecaseInteractor

// 	params := url.Values{"name": []string{"test"}}

// 	rawBody := map[string]interface{}{"test": "body"}

// 	body, err := mapToJSONReader(rawBody)
// 	if err != nil {
// 		t.Errorf("failed to marshal body: %v", err)
// 	}

// 	type args struct {
// 		resourceName string
// 		path         string
// 		params       url.Values
// 		body         io.Reader
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []byte
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid: correct params passed",
// 			args: args{
// 				resourceName: "patient",
// 				path:         "/",
// 				params:       params,
// 				body:         body,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := fh.POSTRequest(tt.args.resourceName, tt.args.path, tt.args.params, tt.args.body)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FHIRUseCaseImpl.POSTRequest() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			// if !reflect.DeepEqual(got, tt.want) {
// 			// 	t.Errorf("FHIRUseCaseImpl.POSTRequest() = %v, want %v", got, tt.want)
// 			// }
// 		})
// 	}
// }

func TestFHIRUseCaseImpl_OpenEpisodes(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	_, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.ID)

	type args struct {
		ctx              context.Context
		patientReference string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.FHIREpisodeOfCare
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:              ctx,
				patientReference: patientRef,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.OpenEpisodes(tt.args.ctx, tt.args.patientReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.OpenEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_HasOpenEpisode(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	_, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}
	type args struct {
		ctx     context.Context
		patient domain.FHIRPatient
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
		panics  bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:     ctx,
				patient: *patient,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				_, err := fh.HasOpenEpisode(tt.args.ctx, tt.args.patient)
				if (err != nil) != tt.wantErr {
					t.Errorf("FHIRUseCaseImpl.HasOpenEpisode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.panics {
				fc := func() { _, _ = fh.HasOpenEpisode(tt.args.ctx, tt.args.patient) }
				assert.Panics(t, fc)
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_CreateFHIREncounter(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
	}
	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
	}

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
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: encounterInput,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:   context.Background(),
				input: encounterInput,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.CreateFHIREncounter(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_GetFHIREpisodeOfCare(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	id := episode.ID

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
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				id:  *id,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.GetFHIREpisodeOfCare(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_StartEncounter(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodeID := episode.ID

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
			name: "valid: correct params passed",
			args: args{
				ctx:       ctx,
				episodeID: *episodeID,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:       context.Background(),
				episodeID: *episodeID,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.StartEncounter(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_StartEpisodeByOtp(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	_, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	testOtp, err := generateTestOTP(t, msisdn)
	if err != nil {
		log.Errorf("unable to get verified phone number and OTP")
		return
	}

	input := domain.OTPEpisodeCreationInput{
		PatientID:    *patient.ID,
		ProviderCode: testProviderCode,
		Msisdn:       msisdn,
		Otp:          testOtp,
		FullAccess:   false,
	}

	type args struct {
		ctx   context.Context
		input domain.OTPEpisodeCreationInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.EpisodeOfCarePayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.StartEpisodeByOtp(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.StartEpisodeByOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_UpgradeEpisode(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodeID := episode.ID

	testOtp, err := generateTestOTP(t, msisdn)
	if err != nil {
		log.Errorf("unable to get verified phone number and OTP")
		return
	}

	input := domain.OTPEpisodeUpgradeInput{
		EpisodeID: *episodeID,
		Msisdn:    msisdn,
		Otp:       testOtp,
	}

	type args struct {
		ctx   context.Context
		input domain.OTPEpisodeUpgradeInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.EpisodeOfCarePayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.UpgradeEpisode(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpgradeEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_SearchEpisodeEncounter(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodeReference := episode.Patient.Reference

	type args struct {
		ctx              context.Context
		episodeReference string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIREncounterRelayConnection
		wantErr bool
		panics  bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:              ctx,
				episodeReference: *episodeReference,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:              context.Background(),
				episodeReference: *episodeReference,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchEpisodeEncounter(tt.args.ctx, tt.args.episodeReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchEpisodeEncounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_EndEncounter(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}
	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	type args struct {
		ctx         context.Context
		encounterID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:         ctx,
				encounterID: encounterID,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:         context.Background(),
				encounterID: encounterID,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.EndEncounter(tt.args.ctx, tt.args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.EndEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_EndEpisode(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodeID := episode.ID

	type args struct {
		ctx       context.Context
		episodeID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:       ctx,
				episodeID: *episodeID,
			},
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:       context.Background(),
				episodeID: *episodeID,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.EndEpisode(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.EndEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_GetActiveEpisode(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodeID := episode.ID

	_, err = fh.StartEncounter(ctx, *episodeID)
	if err != nil {
		t.Errorf("failed to start encounter: %v", err)
	}

	type args struct {
		ctx       context.Context
		episodeID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIREpisodeOfCare
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:       ctx,
				episodeID: *episodeID,
			},
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.GetActiveEpisode(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetActiveEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_SearchFHIRServiceRequest(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	params := map[string]interface{}{"name": "123"}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRServiceRequestRelayConnection
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIRServiceRequest(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateFHIRServiceRequest(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}
	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	input, _, err := getTestSimpleServiceRequest(ctx, encounterID, patient)
	if err != nil {
		t.Errorf("cant get simpleservice request: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRServiceRequestInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRServiceRequestRelayPayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.CreateFHIRServiceRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_SearchFHIRAllergyIntolerance(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}
	fh := testUsecaseInteractor

	params := map[string]interface{}{}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRAllergyIntoleranceRelayConnection
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIRAllergyIntolerance(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateFHIRAllergyIntolerance(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	input, err := getTestAlergyIntorelaceInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get allergy intolerance input: %v", err)
	}

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
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.CreateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_UpdateFHIRAllergyIntolerance(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}
	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	input, err := getTestAlergyIntorelaceInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get allergy intolerance input: %v", err)
	}

	intolerance, err := fh.CreateFHIRAllergyIntolerance(ctx, *input)
	if err != nil {
		t.Errorf("failed to create allergy tolerance input: %v", err)
	}

	input.ID = intolerance.Resource.ID

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
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.UpdateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_SearchFHIRComposition(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	params := map[string]interface{}{"name": "123"}

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
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIRComposition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateFHIRComposition(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	input, err := getFhirComposition(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
	}

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
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:   context.Background(),
				input: *input,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.CreateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_UpdateFHIRComposition(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	input, err := getFhirComposition(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
	}

	composition, err := fh.CreateFHIRComposition(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
	}

	input.ID = composition.Resource.ID

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
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:   context.Background(),
				input: *input,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.UpdateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_DeleteFHIRComposition(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	input, err := getFhirComposition(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
	}

	composition, err := fh.CreateFHIRComposition(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
	}

	id := composition.Resource.ID

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
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				id:  *id,
			},
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx: context.Background(),
				id:  *id,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.DeleteFHIRComposition(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_SearchFHIRCondition(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	params := map[string]interface{}{"name": "123"}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRConditionRelayConnection
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIRCondition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_UpdateFHIRCondition(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
	}

	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
	}

	encounter, err := fh.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to create FHIREncounter: %v", err)
	}

	input, err := createTestConditionInput(*encounter.Resource.ID, *patient.ID)
	if err != nil {
		fmt.Printf("cant create condition: %v\n", err)
		return
	}

	condition, err := fh.CreateFHIRCondition(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir condition: %v", err)
	}

	input.ID = condition.Resource.ID

	type args struct {
		ctx   context.Context
		input domain.FHIRConditionInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRConditionRelayPayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:   context.Background(),
				input: *input,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.UpdateFHIRCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_GetFHIREncounter(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodePayload, err := fh.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
	}

	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
	}

	encounter, err := fh.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to create FHIREncounter: %v", err)
	}

	id := encounter.Resource.ID

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIREncounterRelayPayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				id:  *id,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx: context.Background(),
				id:  *id,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.GetFHIREncounter(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_SearchFHIREncounter(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	params := map[string]interface{}{"name": "123"}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIREncounterRelayConnection
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIREncounter(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_SearchFHIRMedicationRequest(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	params := map[string]interface{}{"name": "123"}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRMedicationRequestRelayConnection
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIRMedicationRequest(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_CreateFHIRMedicationRequest(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	conditionInput, err := createTestConditionInput(encounterID, *patient.ID)
	if err != nil {
		fmt.Printf("cant create condition: %v\n", err)
		return
	}

	condition, err := fh.CreateFHIRCondition(ctx, *conditionInput)
	if err != nil {
		t.Errorf("failed to create fhir condition: %v", err)
	}

	input, err := getFHIRMedicationRequestInput(*patient, encounterID, *condition.Resource.ID)
	if err != nil {
		t.Errorf("failed to get fhir medication request: %v", err)
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationRequestInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRMedicationRequestRelayPayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:   context.Background(),
				input: *input,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.CreateFHIRMedicationRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_UpdateFHIRMedicationRequest(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	conditionInput, err := createTestConditionInput(encounterID, *patient.ID)
	if err != nil {
		fmt.Printf("cant create condition: %v\n", err)
		return
	}

	condition, err := fh.CreateFHIRCondition(ctx, *conditionInput)
	if err != nil {
		t.Errorf("failed to create fhir condition: %v", err)
	}

	input, err := getFHIRMedicationRequestInput(*patient, encounterID, *condition.Resource.ID)
	if err != nil {
		t.Errorf("failed to get fhir medication request: %v", err)
	}

	medication, err := fh.CreateFHIRMedicationRequest(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir medications request: %v", err)
	}

	input.ID = medication.Resource.ID

	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationRequestInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRMedicationRequestRelayPayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx:   context.Background(),
				input: *input,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.UpdateFHIRMedicationRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_DeleteFHIRMedicationRequest(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	conditionInput, err := createTestConditionInput(encounterID, *patient.ID)
	if err != nil {
		fmt.Printf("cant create condition: %v\n", err)
		return
	}

	condition, err := fh.CreateFHIRCondition(ctx, *conditionInput)
	if err != nil {
		t.Errorf("failed to create fhir condition: %v", err)
	}

	input, err := getFHIRMedicationRequestInput(*patient, encounterID, *condition.Resource.ID)
	if err != nil {
		t.Errorf("failed to get fhir medication request: %v", err)
	}

	medication, err := fh.CreateFHIRMedicationRequest(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir medications request: %v", err)
	}

	id := medication.Resource.ID

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
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				id:  *id,
			},
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx: context.Background(),
				id:  *id,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.DeleteFHIRMedicationRequest(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_SearchFHIRObservation(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	input, err := getFhirObservationInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get fhir observation input: %v", err)
	}

	observation, err := fh.CreateFHIRObservation(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir observation: %v", err)
	}

	id := observation.Resource.ID

	params := map[string]interface{}{"id": *id}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRObservationRelayConnection
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIRObservation(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_CreateFHIRObservation(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	input, err := getFhirObservationInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get fhir observation input: %v", err)
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRObservationInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRObservationRelayPayload
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:   ctx,
				input: *input,
			},
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.CreateFHIRObservation(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_DeleteFHIRObservation(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	input, err := getFhirObservationInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get fhir observation input: %v", err)
	}

	observation, err := fh.CreateFHIRObservation(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir observation: %v", err)
	}

	id := observation.Resource.ID

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
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				id:  *id,
			},
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx: context.Background(),
				id:  *id,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.DeleteFHIRObservation(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_GetFHIRPatient(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	_, _, err = createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}
	_, err = fh.GetORCreateOrganization(ctx, testProviderCode)
	if err != nil {
		log.Printf("can't get or create test organization : %v\n", err)
	}
	patientFhirInput := getTestFHIRPatientInput()

	fhirPatient, err := fh.CreatePatient(ctx, patientFhirInput)
	if err != nil {
		t.Fatalf("Failed to create patient %v: %v", patientFhirInput, err)
	}

	id := fhirPatient.PatientRecord.ID

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRPatientRelayPayload
		wantErr bool
		panics  bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				id:  *id,
			},
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		if !tt.panics {
			t.Run(tt.name, func(t *testing.T) {
				_, err := fh.GetFHIRPatient(tt.args.ctx, tt.args.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("FHIRUseCaseImpl.GetFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}
		if tt.panics {
			fc := func() { _, _ = fh.GetFHIRPatient(tt.args.ctx, tt.args.id) }
			assert.Panics(t, fc)
		}
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_DeleteFHIRPatient(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor
	patientFhirInput := getTestFHIRPatientInput()

	fhirPatient, err := fh.CreatePatient(ctx, patientFhirInput)
	if err != nil {
		t.Fatalf("Failed to create patient %v: %v", patientFhirInput, err)
	}

	id := fhirPatient.PatientRecord.ID

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
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				id:  *id,
			},
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx: context.Background(),
				id:  *id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.DeleteFHIRPatient(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRUseCaseImpl_DeleteFHIRServiceRequest(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := fh.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	serviceRequest, err := getFhirServiceRequest(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get service request: %v", err)
	}
	request, err := fh.CreateFHIRServiceRequest(ctx, *serviceRequest)
	if err != nil {
		t.Errorf("failed to create service request: %v", err)
	}

	id := request.Resource.ID
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
			name: "valid: correct params passed",
			args: args{
				ctx: ctx,
				id:  *id,
			},
		},
		{
			name: "invalid: missing parameters",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.DeleteFHIRServiceRequest(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.DeleteFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}

func TestFHIRUseCaseImpl_SearchFHIRMedicationStatement(t *testing.T) {
	ctx, err := getTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	fh := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	patient, _, err := createTestPatient(ctx)
	if err != nil {
		log.Printf("cant create test patient: %v\n", err)
		return
	}

	input, err := getFhirMedicationStatementInput(*patient)
	if err != nil {
		t.Errorf("failed to get fhir medication statement input: %v", err)
	}

	statement, err := fh.CreateFHIRMedicationStatement(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir medication statement: %v", err)
	}

	id := statement.Resource.ID

	params := map[string]interface{}{"id": *id}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIRMedicationStatementRelayConnection
		wantErr bool
	}{
		{
			name: "valid: correct params passed",
			args: args{
				ctx:    ctx,
				params: params,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.SearchFHIRMedicationStatement(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRMedicationStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
	// teardown
	deleteTestPatient(ctx, msisdn)
}
