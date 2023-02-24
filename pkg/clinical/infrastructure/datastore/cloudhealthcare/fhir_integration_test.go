package fhir_test

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/interserviceclient"
	"github.com/tj/assert"
)

func TestFHIRImpl_Encounters(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	msisdn := interserviceclient.TestUserPhoneNumber

	_, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant create test episode of care: %v", err)
		return
	}

	f := testInfrastructure

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := f.FHIR.Encounters(tt.args.ctx, tt.args.patientReference, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.Encounters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_CreateEpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	f := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := getTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	type args struct {
		ctx     context.Context
		episode domain.FHIREpisodeOfCare
	}
	tests := []struct {
		name    string
		args    args
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
			wantErr: false,
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
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				t.Run(tt.name, func(t *testing.T) {
					_, err := f.FHIR.CreateEpisodeOfCare(tt.args.ctx, tt.args.episode)
					if (err != nil) != tt.wantErr {
						t.Errorf("FHIRStoreImpl.CreateEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
					}
				})
			}
			if tt.panics {
				fc := func() { _, _ = f.FHIR.CreateEpisodeOfCare(tt.args.ctx, tt.args.episode) }
				assert.Panics(t, fc)
			}
		})
	}
}

func TestFHIRImpl_SearchFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

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
			_, err := fh.FHIR.SearchFHIREpisodeOfCare(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_CreateFHIRCondition(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	episodePayload, err := fh.FHIR.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
		return
	}

	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
		return
	}

	encounter, err := fh.FHIR.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to create FHIREncounter: %v", err)
		return
	}

	input, err := createTestConditionInput(*encounter.Resource.ID, *patient.ID)
	if err != nil {
		t.Errorf("cant create condition: %v\n", err)
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
			_, err := fh.FHIR.CreateFHIRCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.CreateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_CreateFHIROrganization(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	identifier := []*domain.FHIRIdentifierInput{
		{
			Use:   "official",
			Value: testProviderCode,
		},
	}
	OrgName := gofakeit.BeerName()

	input := domain.FHIROrganizationInput{
		Identifier: identifier,
		Name:       &OrgName,
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
			_, err := fh.FHIR.CreateFHIROrganization(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.CreateFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchFHIROrganization(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	OrgName := gofakeit.BeerName()

	// Create FHIR organization
	identifier := []*domain.FHIRIdentifierInput{
		{
			Use:   "official",
			Value: testProviderCode,
		},
	}

	input := domain.FHIROrganizationInput{
		Identifier: identifier,
		Name:       &OrgName,
	}

	_, err = fh.FHIR.CreateFHIROrganization(ctx, input)
	if err != nil {
		t.Errorf("failed to create fhir organization: %v", err)
		return
	}

	params := map[string]interface{}{"provider": "1234"}

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
			_, err := fh.FHIR.SearchFHIROrganization(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

func TestFHIRImpl_FindOrganizationByID(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	organizationID := "a1bf993c-c2b6-44dd-8991-3a47b54a6789"
	OrgName := gofakeit.BeerName()

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
		Name:       &OrgName,
	}

	organization, err := fh.FHIR.CreateFHIROrganization(ctx, input)
	if err != nil {
		t.Errorf("failed to create fhir organization: %v", err)
		return
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
			got, err := fh.FHIR.FindOrganizationByID(tt.args.ctx, tt.args.organizationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.FindOrganizationByID() error = %v, wantErr %v", err, tt.wantErr)
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
	err = fh.FHIR.DeleteFHIRResourceType(resource)
	if err != nil {
		t.Errorf("failed to delete fhir organization: %v", err)
	}
}

func TestFHIRImpl_SearchEpisodesByParam(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	_, _, err = createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	searchParams := url.Values{"": []string{"resourceType"}}
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
			_, err := fh.FHIR.SearchEpisodesByParam(tt.args.ctx, tt.args.searchParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchEpisodesByParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_OpenEpisodes(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	_, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
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
			_, err := fh.FHIR.OpenEpisodes(tt.args.ctx, tt.args.patientReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.OpenEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_HasOpenEpisode(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	_, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
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
				_, err := fh.FHIR.HasOpenEpisode(tt.args.ctx, tt.args.patient)
				if (err != nil) != tt.wantErr {
					t.Errorf("FHIRStoreImpl.HasOpenEpisode() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.panics {
				fc := func() { _, _ = fh.FHIR.HasOpenEpisode(tt.args.ctx, tt.args.patient) }
				assert.Panics(t, fc)
			}
		})
	}
}

func TestFHIRImpl_CreateFHIREncounter(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	episodePayload, err := fh.FHIR.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
		return
	}
	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
		return
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
			_, err := fh.FHIR.CreateFHIREncounter(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.CreateFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_GetFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
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
			_, err := fh.FHIR.GetFHIREpisodeOfCare(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.GetFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_StartEncounter(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
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
			_, err := fh.FHIR.StartEncounter(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchEpisodeEncounter(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.FHIR.SearchEpisodeEncounter(tt.args.ctx, tt.args.episodeReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchEpisodeEncounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFHIRImpl_EndEncounter(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
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
			_, err := fh.FHIR.EndEncounter(tt.args.ctx, tt.args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.EndEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_EndEpisode(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
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
			_, err := fh.FHIR.EndEpisode(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.EndEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_GetActiveEpisode(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	episodeID := episode.ID

	_, err = fh.FHIR.StartEncounter(ctx, *episodeID)
	if err != nil {
		t.Errorf("failed to start encounter: %v", err)
		return
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
			_, err := fh.FHIR.GetActiveEpisode(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.GetActiveEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchFHIRServiceRequest(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

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
			_, err := fh.FHIR.SearchFHIRServiceRequest(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_CreateFHIRServiceRequest(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
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
			_, err := fh.FHIR.CreateFHIRServiceRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.CreateFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

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
			_, err := fh.FHIR.SearchFHIRAllergyIntolerance(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_CreateFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	input, err := getTestAlergyIntorelaceInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get allergy intolerance input: %v", err)
		return
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
			_, err := fh.FHIR.CreateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.CreateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_UpdateFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	input, err := getTestAlergyIntorelaceInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get allergy intolerance input: %v", err)
		return
	}

	intolerance, err := fh.FHIR.CreateFHIRAllergyIntolerance(ctx, *input)
	if err != nil {
		t.Errorf("failed to create allergy tolerance input: %v", err)
		return
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
			_, err := fh.FHIR.UpdateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.UpdateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchFHIRComposition(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

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
			_, err := fh.FHIR.SearchFHIRComposition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_CreateFHIRComposition(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	input, err := getFhirComposition(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
		return
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
			_, err := fh.FHIR.CreateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.CreateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_UpdateFHIRComposition(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	input, err := getFhirComposition(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
		return
	}

	composition, err := fh.FHIR.CreateFHIRComposition(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
		return
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
			_, err := fh.FHIR.UpdateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.UpdateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_DeleteFHIRComposition(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	input, err := getFhirComposition(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
		return
	}

	composition, err := fh.FHIR.CreateFHIRComposition(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir composition: %v", err)
		return
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
			_, err := fh.FHIR.DeleteFHIRComposition(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.DeleteFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchFHIRCondition(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

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
			_, err := fh.FHIR.SearchFHIRCondition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_UpdateFHIRCondition(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	episodePayload, err := fh.FHIR.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
		return
	}

	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
		return
	}

	encounter, err := fh.FHIR.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to create FHIREncounter: %v", err)
		return
	}

	input, err := createTestConditionInput(*encounter.Resource.ID, *patient.ID)
	if err != nil {
		t.Errorf("cant create condition: %v\n", err)
		return
	}

	condition, err := fh.FHIR.CreateFHIRCondition(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir condition: %v", err)
		return
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
			_, err := fh.FHIR.UpdateFHIRCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.UpdateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_GetFHIREncounter(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	episodePayload, err := fh.FHIR.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
		return
	}

	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
		return
	}

	encounter, err := fh.FHIR.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to create FHIREncounter: %v", err)
		return
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
			_, err := fh.FHIR.GetFHIREncounter(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.GetFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchFHIREncounter(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

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
			_, err := fh.FHIR.SearchFHIREncounter(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

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
			_, err := fh.FHIR.SearchFHIRMedicationRequest(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_CreateFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	conditionInput, err := createTestConditionInput(encounterID, *patient.ID)
	if err != nil {
		t.Errorf("cant create condition: %v\n", err)
		return
	}

	condition, err := fh.FHIR.CreateFHIRCondition(ctx, *conditionInput)
	if err != nil {
		t.Errorf("failed to create fhir condition: %v", err)
		return
	}

	input, err := getFHIRMedicationRequestInput(*patient, encounterID, *condition.Resource.ID)
	if err != nil {
		t.Errorf("failed to get fhir medication request: %v", err)
		return
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
			_, err := fh.FHIR.CreateFHIRMedicationRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.CreateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_UpdateFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	conditionInput, err := createTestConditionInput(encounterID, *patient.ID)
	if err != nil {
		t.Errorf("cant create condition: %v\n", err)
		return
	}

	condition, err := fh.FHIR.CreateFHIRCondition(ctx, *conditionInput)
	if err != nil {
		t.Errorf("failed to create fhir condition: %v", err)
		return
	}

	input, err := getFHIRMedicationRequestInput(*patient, encounterID, *condition.Resource.ID)
	if err != nil {
		t.Errorf("failed to get fhir medication request: %v", err)
		return
	}

	medication, err := fh.FHIR.CreateFHIRMedicationRequest(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir medications request: %v", err)
		return
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
			_, err := fh.FHIR.UpdateFHIRMedicationRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.UpdateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_DeleteFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	conditionInput, err := createTestConditionInput(encounterID, *patient.ID)
	if err != nil {
		t.Errorf("cant create condition: %v\n", err)
		return
	}

	condition, err := fh.FHIR.CreateFHIRCondition(ctx, *conditionInput)
	if err != nil {
		t.Errorf("failed to create fhir condition: %v", err)
		return
	}

	input, err := getFHIRMedicationRequestInput(*patient, encounterID, *condition.Resource.ID)
	if err != nil {
		t.Errorf("failed to get fhir medication request: %v", err)
		return
	}

	medication, err := fh.FHIR.CreateFHIRMedicationRequest(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir medications request: %v", err)
		return
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
			_, err := fh.FHIR.DeleteFHIRMedicationRequest(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.DeleteFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchFHIRObservation(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	input, err := getFhirObservationInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get fhir observation input: %v", err)
		return
	}

	observation, err := fh.FHIR.CreateFHIRObservation(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir observation: %v", err)
		return
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
			_, err := fh.FHIR.SearchFHIRObservation(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_CreateFHIRObservation(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	input, err := getFhirObservationInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get fhir observation input: %v", err)
		return
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
			_, err := fh.FHIR.CreateFHIRObservation(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.CreateFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_DeleteFHIRObservation(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	input, err := getFhirObservationInput(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get fhir observation input: %v", err)
		return
	}

	observation, err := fh.FHIR.CreateFHIRObservation(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir observation: %v", err)
		return
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
			_, err := fh.FHIR.DeleteFHIRObservation(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.DeleteFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_GetFHIRPatient(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure
	u := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	_, _, err = createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}
	_, err = u.GetORCreateOrganization(ctx, testProviderCode)
	if err != nil {
		log.Printf("can't get or create test organization : %v\n", err)
	}
	patientFhirInput := getTestFHIRPatientInput()

	fhirPatient, err := u.CreatePatient(ctx, patientFhirInput)
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
				_, err := fh.FHIR.GetFHIRPatient(tt.args.ctx, tt.args.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("FHIRStoreImpl.GetFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}
		if tt.panics {
			fc := func() { _, _ = fh.FHIR.GetFHIRPatient(tt.args.ctx, tt.args.id) }
			assert.Panics(t, fc)
		}
	}
}

func TestFHIRImpl_DeleteFHIRPatient(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure
	u := testUsecaseInteractor
	patientFhirInput := getTestFHIRPatientInput()

	fhirPatient, err := u.CreatePatient(ctx, patientFhirInput)
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
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context",
			args: args{
				ctx: context.Background(),
				id:  *id,
			},
			wantErr: true, // patient already deleted
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fh.FHIR.DeleteFHIRPatient(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.DeleteFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_DeleteFHIRServiceRequest(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
		return
	}

	encounterID, err := fh.FHIR.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
		return
	}

	serviceRequest, err := getFhirServiceRequest(*patient, encounterID)
	if err != nil {
		t.Errorf("failed to get service request: %v", err)
		return
	}
	request, err := fh.FHIR.CreateFHIRServiceRequest(ctx, *serviceRequest)
	if err != nil {
		t.Errorf("failed to create service request: %v", err)
		return
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
			_, err := fh.FHIR.DeleteFHIRServiceRequest(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.DeleteFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFHIRImpl_SearchFHIRMedicationStatement(t *testing.T) {
	ctx := context.Background()
	ctx, err := addOrganisationContext(ctx, testProviderCode)
	if err != nil {
		t.Errorf("cant add test organisation context: %v\n", err)
		return
	}

	fh := testInfrastructure

	patient, _, err := createTestPatient(ctx)
	if err != nil {
		t.Errorf("cant create test patient: %v\n", err)
		return
	}

	input, err := getFhirMedicationStatementInput(*patient)
	if err != nil {
		t.Errorf("failed to get fhir medication statement input: %v", err)
		return
	}

	statement, err := fh.FHIR.CreateFHIRMedicationStatement(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir medication statement: %v", err)
		return
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
			_, err := fh.FHIR.SearchFHIRMedicationStatement(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.SearchFHIRMedicationStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
