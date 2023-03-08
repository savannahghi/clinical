package fhir_test

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	FHIR "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"

	fakeDataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset/mock"
)

func TestStoreImpl_CreateEpisodeOfCare_Unittest(t *testing.T) {

	UUID := uuid.New().String()
	PatientRef := "Patient/1"
	OrgRef := "Organization/1"
	status := domain.EpisodeOfCareStatusEnumFinished

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
				ctx: context.Background(),
				episode: domain.FHIREpisodeOfCare{
					ID:            &UUID,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
					Status:        &(status),
					StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
					Type:          []*domain.FHIRCodeableConcept{},
					Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
					Patient: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReference{
						Reference: &OrgRef,
					},
					Period:          &domain.FHIRPeriod{},
					ReferralRequest: []*domain.FHIRReference{},
					CareManager:     &domain.FHIRReference{},
					Team:            []*domain.FHIRReference{},
					Account:         []*domain.FHIRReference{},
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
				episode: domain.FHIREpisodeOfCare{
					ID:            nil,
					Text:          &domain.FHIRNarrative{},
					Identifier:    []*domain.FHIRIdentifier{},
					StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
					Type:          []*domain.FHIRCodeableConcept{},
					Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
					Patient: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReference{
						Reference: &OrgRef,
					},
					Period:          &domain.FHIRPeriod{},
					ReferralRequest: []*domain.FHIRReference{},
					CareManager:     &domain.FHIRReference{},
					Team:            []*domain.FHIRReference{},
					Account:         []*domain.FHIRReference{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}) ([]byte, error) {
					return nil, fmt.Errorf("error")
				}
			}
			_, err := h.CreateEpisodeOfCare(tt.args.ctx, tt.args.episode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRStoreImpl.CreateEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRCondition(t *testing.T) {

	ID := uuid.New().String()

	invalidfhirconditionInput := domain.FHIRConditionInput{
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
				ctx: context.Background(),
				input: domain.FHIRConditionInput{
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
					OnsetDateTime: &scalarutils.Date{
						Year:  2000,
						Month: 3,
						Day:   30,
					},
					OnsetAge:    &domain.FHIRAgeInput{},
					OnsetPeriod: &domain.FHIRPeriodInput{},
					OnsetRange:  &domain.FHIRRangeInput{},
					OnsetString: new(string),
					AbatementDateTime: &scalarutils.Date{
						Year:  2000,
						Month: 3,
						Day:   30,
					},
					AbatementAge:    &domain.FHIRAgeInput{},
					AbatementPeriod: &domain.FHIRPeriodInput{},
					AbatementRange:  &domain.FHIRRangeInput{},
					AbatementString: new(string),
					RecordedDate: &scalarutils.Date{
						Year:  2000,
						Month: 3,
						Day:   30,
					},
					Recorder: &domain.FHIRReferenceInput{},
					Asserter: &domain.FHIRReferenceInput{},
					Stage:    []*domain.FHIRConditionStageInput{},
					Evidence: []*domain.FHIRConditionEvidenceInput{},
					Note:     []*domain.FHIRAnnotationInput{},
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   context.Background(),
				input: invalidfhirconditionInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()
			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}) ([]byte, error) {
					return nil, fmt.Errorf("error")
				}
			}
			_, err := h.CreateFHIRCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIROrganization_Unittest(t *testing.T) {

	ID := ksuid.New().String()
	active := true
	testname := gofakeit.FirstName()

	orgInput := domain.FHIROrganizationInput{
		ID:         &ID,
		Active:     &active,
		Identifier: []*domain.FHIRIdentifierInput{},
		Type:       []*domain.FHIRCodeableConceptInput{},
		Name:       &testname,
		Alias:      []string{"alias test"},
		Telecom:    []*domain.FHIRContactPointInput{},
		Address:    []*domain.FHIRAddressInput{},
	}

	invalidOrgInput := domain.FHIROrganizationInput{
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
				ctx:   context.Background(),
				input: orgInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   context.Background(),
				input: invalidOrgInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				if tt.name == "Sad case" {
					d.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}) ([]byte, error) {
						return nil, fmt.Errorf("error")
					}
				}
			}
			_, err := h.CreateFHIROrganization(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_FindOrganizationByID_Unittest(t *testing.T) {

	type args struct {
		ctx            context.Context
		organizationID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FHIROrganizationRelayPayload
		wantErr bool
	}{
		{
			name: "Sad case",
			args: args{
				ctx:            context.Background(),
				organizationID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockGetFHIRResourceFn = func(resourceType string, id string) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := h.FindOrganizationByID(tt.args.ctx, tt.args.organizationID)
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
}

func TestStoreImpl_SearchEpisodesByParam(t *testing.T) {

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
				ctx: context.Background(),
				searchParams: map[string][]string{
					"patient": {"12233"},
				},
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
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockPOSTRequestFn = func(resourceName, path string, params url.Values, body io.Reader) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.SearchEpisodesByParam(tt.args.ctx, tt.args.searchParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchEpisodesByParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIREpisodeOfCare(t *testing.T) {

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
				ctx: context.Background(),
				params: map[string]interface{}{
					"patient": "12233",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"patient": "12233",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockPOSTRequestFn = func(resourceName, path string, params url.Values, body io.Reader) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.SearchFHIREpisodeOfCare(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_OpenEpisodes(t *testing.T) {

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
				ctx:              context.Background(),
				patientReference: patientReference,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:              context.Background(),
				patientReference: patientReference,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockPOSTRequestFn = func(resourceName, path string, params url.Values, body io.Reader) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.OpenEpisodes(tt.args.ctx, tt.args.patientReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.OpenEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_HasOpenEpisode(t *testing.T) {

	UUID := uuid.New().String()

	patient := domain.FHIRPatient{
		ID: &UUID,
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
	}{

		{
			name: "Happy case",
			args: args{
				ctx:     context.Background(),
				patient: patient,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     context.Background(),
				patient: patient,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockPOSTRequestFn = func(resourceName, path string, params url.Values, body io.Reader) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.HasOpenEpisode(tt.args.ctx, tt.args.patient)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.HasOpenEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIREncounter(t *testing.T) {

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
				ctx:   context.Background(),
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.CreateFHIREncounter(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestStoreImpl_GetFHIREpisodeOfCare(t *testing.T) {

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
				ctx: context.Background(),
				id:  id,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.GetFHIREpisodeOfCare(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_StartEncounter(t *testing.T) {

	episodeID := uuid.New().String()

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
			name: "Sad case",
			args: args{
				ctx:       context.Background(),
				episodeID: episodeID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.StartEncounter(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIREncounter(t *testing.T) {

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
				ctx: context.Background(),
				id:  uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.GetFHIREncounter(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRServiceRequest(t *testing.T) {

	UUID := uuid.New().String()

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
				ctx: context.Background(),
				input: domain.FHIRServiceRequestInput{
					ID: &UUID,
				},
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
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.CreateFHIRServiceRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRAllergyIntolerance(t *testing.T) {

	UUID := uuid.New().String()

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
				ctx: context.Background(),
				input: domain.FHIRAllergyIntoleranceInput{
					ID: &UUID,
				},
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
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.CreateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIRAllergyIntolerance(t *testing.T) {

	UUID := uuid.New().String()

	input := domain.FHIRAllergyIntoleranceInput{
		ID: &UUID,
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
			name: "Happy case",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.UpdateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRComposition(t *testing.T) {

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
				ctx:    context.Background(),
				params: params,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    context.Background(),
				params: params,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockPOSTRequestFn = func(resourceName, path string, params url.Values, body io.Reader) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := h.SearchFHIRComposition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.SearchFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRComposition(t *testing.T) {

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
				ctx:   context.Background(),
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := h.CreateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.CreateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIRComposition(t *testing.T) {

	UUID := uuid.New().String()

	input := domain.FHIRCompositionInput{
		ID: &UUID,
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
			name: "Happy case",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := h.UpdateFHIRComposition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.UpdateFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRComposition(t *testing.T) {

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
				ctx: context.Background(),
				id:  id,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d = fakeDataset.NewFakeFHIRRepositoryMock()

			h := FHIR.NewFHIRStoreImpl(d)

			if tt.name == "Sad case" {
				d.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := h.DeleteFHIRComposition(tt.args.ctx, tt.args.id)
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
