package fhir_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	FHIR "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"

	fakeDataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset/mock"
)

func TestStoreImpl_CreateFHIRCondition(t *testing.T) {

	ID := uuid.New().String()

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			_, err := fh.CreateFHIRCondition(tt.args.ctx, tt.args.input)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			_, err := fh.CreateFHIROrganization(tt.args.ctx, tt.args.input)
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
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case" {
				dataset.MockGetFHIRResourceFn = func(resourceType string, id string) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			_, err := fh.CreateFHIREncounter(tt.args.ctx, tt.args.input)
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
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string) ([]byte, error) {
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
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
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
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := fh.GetFHIREncounter(tt.args.ctx, tt.args.id)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			_, err := fh.CreateFHIRServiceRequest(tt.args.ctx, tt.args.input)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			_, err := fh.CreateFHIRAllergyIntolerance(tt.args.ctx, tt.args.input)
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
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case" {
				dataset.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			_, err := fh.CreateFHIRComposition(tt.args.ctx, tt.args.input)
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
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case" {
				dataset.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
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
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case" {
				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) ([]byte, error) {
					return nil, fmt.Errorf("an error occurred")
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
