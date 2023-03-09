package fhir_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
			name: "Happy case: create FHIR condition",
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
			name: "Sad case: failed to create FHIR condition",
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to create FHIR condition" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

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
			name: "Happy case: create FHIR organization",
			args: args{
				ctx:   context.Background(),
				input: orgInput,
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to create FHIR organization",
			args: args{
				ctx:   context.Background(),
				input: orgInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to create FHIR organization" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("an error occurred")
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
			name: "Happy case: find organization by ID",
			args: args{
				ctx:            context.Background(),
				organizationID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: missing organization ID",
			args: args{
				ctx:            context.Background(),
				organizationID: "",
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to find organization by ID",
			args: args{
				ctx:            context.Background(),
				organizationID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to find organization by ID" {
				dataset.MockGetFHIRResourceFn = func(resourceType string, id string, resource interface{}) error {
					return fmt.Errorf("an error occurred")
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
			name: "Happy case: create FHIR encounter",
			args: args{
				ctx:   context.Background(),
				input: input,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create FHIR resource",
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

			if tt.name == "Sad Case - Fail to create FHIR resource" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to create fhir service request")
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
			name: "Happy case: get FHIR episode of care",
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get FHIR resource",
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

			if tt.name == "Sad case: failed to get FHIR resource" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					return fmt.Errorf("an error occurred")
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
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					return fmt.Errorf("an error occurred")
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
		{
			name: "Sad Case - Fail to create FHIR service request",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - Fail to create FHIR service request" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to create fhir service request")
				}
			}

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
		{
			name: "Sad Case - Fail to create FHIR allergy intolerance",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - Fail to create FHIR allergy intolerance" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to create fhir service request")
				}
			}

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
		{
			name: "Sad Case - missing input",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case" {
				dataset.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("an error occurred")
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
				input: domain.FHIRCompositionInput{},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create FHIR composition",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - Fail to create FHIR composition" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to create FHIR composition")
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
		{
			name: "Sad Case - Missing user ID",
			args: args{
				ctx:   context.Background(),
				input: domain.FHIRCompositionInput{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()

			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case" {
				dataset.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("an error occurred")
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
				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					return fmt.Errorf("an error occurred")
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

func TestStoreImpl_CreateFHIRMedicationRequest(t *testing.T) {
	id := uuid.New().String()
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
			name: "Happy Case - Successfully create medication request",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRMedicationRequestInput{
					ID: &id,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to create medication request",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRMedicationRequestInput{
					ID: &id,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to create medication request" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to create fhir service request")
				}
			}

			got, err := fh.CreateFHIRMedicationRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRServiceRequest(t *testing.T) {
	ctx := context.Background()

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
			name: "Happy Case - successfully search fhir service request",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"id": "1234",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a service request",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to search a service request" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRServiceRequest(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRAllergyIntolerance(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - successfully search fhir allergy intolerance",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"id": "1234",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search an allergy intolerance",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to search an allergy intolerance" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRAllergyIntolerance(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRComposition(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - successfully search fhir composition",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"id": "1234",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a composition",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to search a composition" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRComposition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRComposition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRCondition(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - successfully search fhir condition",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"id": "1234",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a condition",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to search a condition" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRCondition(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIREncounter(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - successfully search fhir encounter",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"id": "1234",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search an encounter",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to search an encounter" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIREncounter(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIREncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - successfully search fhir medication request",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"id": "1234",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a medication request",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to search a medication request" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRMedicationRequest(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIRCondition(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
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
			name: "Happy Case - successfully update fhir condition",
			args: args{ctx: ctx, input: domain.FHIRConditionInput{
				ID: &id,
			}},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to update fhir condition",
			args: args{ctx: ctx, input: domain.FHIRConditionInput{
				ID: &id,
			}},
			wantErr: true,
		},
		{
			name:    "Sad Case - missing ID",
			args:    args{ctx: ctx, input: domain.FHIRConditionInput{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to update fhir condition" {
				dataset.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to update condition")
				}
			}

			got, err := fh.UpdateFHIRCondition(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.UpdateFHIRCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
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
			name: "Happy Case - successfully update fhir medication request",
			args: args{ctx: ctx, input: domain.FHIRMedicationRequestInput{
				ID: &id,
			}},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to update fhir medication request",
			args: args{ctx: ctx, input: domain.FHIRMedicationRequestInput{
				ID: &id,
			}},
			wantErr: true,
		},
		{
			name:    "Sad Case - missing ID",
			args:    args{ctx: ctx, input: domain.FHIRMedicationRequestInput{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to update fhir medication request" {
				dataset.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to update medication request")
				}
			}

			got, err := fh.UpdateFHIRMedicationRequest(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.UpdateFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRMedicationRequest(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - successfully delete a medication request",
			args: args{
				ctx: ctx,
				id:  "1234567890",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - fail to delete a medication request",
			args: args{
				ctx: ctx,
				id:  "12445",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to delete a medication request" {
				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					return fmt.Errorf("failed to update resource")
				}
			}

			got, err := fh.DeleteFHIRMedicationRequest(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRMedicationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.DeleteFHIRMedicationRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_Encounters(t *testing.T) {
	status := domain.EncounterStatusEnumPlanned
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
			name: "Happy case: get encounters",
			args: args{
				ctx:              context.Background(),
				patientReference: gofakeit.BS(),
				status:           &status,
			},
			wantErr: false,
		},
		{
			name: "Happy case: nil status",
			args: args{
				ctx:              context.Background(),
				patientReference: gofakeit.BS(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			args: args{
				ctx:              context.Background(),
				patientReference: gofakeit.BS(),
				status:           &status,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to search FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.Encounters(tt.args.ctx, tt.args.patientReference, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.Encounters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
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
			name: "Happy Case: search FHIR episode of care",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"test": "search",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"test": "search",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to search FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.SearchFHIREpisodeOfCare(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_CreateEpisodeOfCare(t *testing.T) {
	status := domain.EpisodeOfCareStatusEnumPlanned
	UUID := gofakeit.UUID()
	PatientRef := "Patient/1"
	OrgRef := "Organization/"
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
			name: "Happy case: create episode of care",
			args: args{
				ctx: context.Background(),
				episode: domain.FHIREpisodeOfCare{
					ID:     &UUID,
					Status: &status,
					Patient: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReference{
						Reference: &OrgRef,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: create episode of care, episode does not exist",
			args: args{
				ctx: context.Background(),
				episode: domain.FHIREpisodeOfCare{
					ID:     &UUID,
					Status: &status,
					Patient: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReference{
						Reference: &OrgRef,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR episode of care",
			args: args{
				ctx: context.Background(),
				episode: domain.FHIREpisodeOfCare{
					ID:     &UUID,
					Status: &status,
					Patient: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReference{
						Reference: &OrgRef,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to create FHIR resource",
			args: args{
				ctx: context.Background(),
				episode: domain.FHIREpisodeOfCare{
					ID:     &UUID,
					Status: &status,
					Patient: &domain.FHIRReference{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReference{
						Reference: &OrgRef,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Happy case: create episode of care, episode does not exist" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, nil
				}
			}

			if tt.name == "Sad case: failed to search FHIR episode of care" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to create FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, nil
				}
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.CreateEpisodeOfCare(tt.args.ctx, tt.args.episode)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateEpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIROrganization(t *testing.T) {
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
			name: "Happy case: search FHIR organisation",
			args: args{
				ctx: nil,
				params: map[string]interface{}{
					"test": "params",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR organisation",
			args: args{
				ctx: nil,
				params: map[string]interface{}{
					"test": "params",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to search FHIR organisation" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := fh.SearchFHIROrganization(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchEpisodesByParam(t *testing.T) {
	type args struct {
		ctx          context.Context
		searchParams map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: search episode by param",
			args: args{
				ctx: context.Background(),
				searchParams: map[string]interface{}{
					"period": map[string]interface{}{
						"start": time.February.String(),
						"end":   time.February.String(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			args: args{
				ctx: context.Background(),
				searchParams: map[string]interface{}{
					"period": map[string]interface{}{
						"start": time.February.String(),
						"end":   time.February.String(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to search FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.SearchEpisodesByParam(tt.args.ctx, tt.args.searchParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchEpisodesByParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_OpenEpisodes(t *testing.T) {
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
			name: "Happy case: open episodes",
			args: args{
				ctx:              context.Background(),
				patientReference: gofakeit.BS(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			args: args{
				ctx:              context.Background(),
				patientReference: gofakeit.BS(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to search FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.OpenEpisodes(tt.args.ctx, tt.args.patientReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.OpenEpisodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_HasOpenEpisode(t *testing.T) {
	UUID := gofakeit.UUID()
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
			name: "Happy case: has open episodes",
			args: args{
				ctx: context.Background(),
				patient: domain.FHIRPatient{
					ID: &UUID,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			args: args{
				ctx: context.Background(),
				patient: domain.FHIRPatient{
					ID: &UUID,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to search FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.HasOpenEpisode(tt.args.ctx, tt.args.patient)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.HasOpenEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.HasOpenEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_SearchEpisodeEncounter(t *testing.T) {
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
			name: "Happy case: search episode encounter",
			args: args{
				ctx:              context.Background(),
				episodeReference: gofakeit.BS(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			args: args{
				ctx:              context.Background(),
				episodeReference: gofakeit.BS(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to search FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := fh.SearchEpisodeEncounter(tt.args.ctx, tt.args.episodeReference)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchEpisodeEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_GetActiveEpisode(t *testing.T) {
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
			name: "Happy case: get active episodes",
			args: args{
				ctx:       context.Background(),
				episodeID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search FHIR resource",
			args: args{
				ctx:       context.Background(),
				episodeID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: empty FHIR resource",
			args: args{
				ctx:       context.Background(),
				episodeID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad case: failed to search FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: empty FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
					return []map[string]interface{}{}, nil
				}
			}
			got, err := fh.GetActiveEpisode(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetActiveEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_StartEncounter(t *testing.T) {
	status := domain.EpisodeOfCareStatusEnumActive
	finishedStatus := domain.EpisodeOfCareStatusEnumFinished
	UUID := gofakeit.UUID()
	dummyString := gofakeit.BS()
	uri := "foo://example.com:8042"
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
			name: "Happy case: start encounter",
			args: args{
				ctx:       context.Background(),
				episodeID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get FHIR episode of care",
			args: args{
				ctx:       context.Background(),
				episodeID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: episode  not active",
			args: args{
				ctx:       context.Background(),
				episodeID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to create encounter",
			args: args{
				ctx:       context.Background(),
				episodeID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Happy case: start encounter" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Patient: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
						ManagingOrganization: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}

				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID: &UUID,
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}

			}

			if tt.name == "Sad case: failed to get FHIR episode of care" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Patient: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
						ManagingOrganization: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return fmt.Errorf("an error occurred")
				}

			}

			if tt.name == "Sad case: episode  not active" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &finishedStatus,
						Patient: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
						ManagingOrganization: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}

				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID: &UUID,
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}

			}

			if tt.name == "Sad case: failed to create encounter" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Patient: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
						ManagingOrganization: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}

				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID: &UUID,
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return fmt.Errorf("an error occurred")
				}

			}

			if tt.name == "Happy case: start encounter" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Patient: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
						ManagingOrganization: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}

				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID: &UUID,
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}

			}

			if tt.name == "Happy case: start encounter" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Patient: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
						ManagingOrganization: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &dummyString,
							Type:      (*scalarutils.URI)(&uri),
							Display:   dummyString,
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}

				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID: &UUID,
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}

			}

			got, err := fh.StartEncounter(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.StartEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("expected a value to be returned, got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_EndEncounter(t *testing.T) {
	status := domain.EpisodeOfCareStatusEnumActive
	UUID := gofakeit.UUID()
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
			name: "Happy case: end encounter",
			args: args{
				ctx:         context.Background(),
				encounterID: gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: failed to get encounter",
			args: args{
				ctx:         context.Background(),
				encounterID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to get update resource",
			args: args{
				ctx:         context.Background(),
				encounterID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Happy case: end encounter" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Period: &domain.FHIRPeriod{
							ID:    &UUID,
							Start: scalarutils.DateTime(time.February.String()),
							End:   scalarutils.DateTime(time.March.String()),
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}
			}

			if tt.name == "Sad case: failed to get encounter" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Period: &domain.FHIRPeriod{
							ID:    &UUID,
							Start: scalarutils.DateTime(time.February.String()),
							End:   scalarutils.DateTime(time.March.String()),
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to get update resource" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Period: &domain.FHIRPeriod{
							ID:    &UUID,
							Start: scalarutils.DateTime(time.February.String()),
							End:   scalarutils.DateTime(time.March.String()),
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}
				dataset.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.EndEncounter(tt.args.ctx, tt.args.encounterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.EndEncounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.EndEncounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_EndEpisode(t *testing.T) {
	status := domain.EpisodeOfCareStatusEnumActive
	UUID := gofakeit.UUID()
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
			name: "Happy case: end episode",
			args: args{
				ctx:       context.Background(),
				episodeID: UUID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: failed to get encounter",
			args: args{
				ctx:       context.Background(),
				episodeID: UUID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to get update resource",
			args: args{
				ctx:       context.Background(),
				episodeID: UUID,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Happy case: end episode" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Period: &domain.FHIRPeriod{
							ID:    &UUID,
							Start: scalarutils.DateTime(time.February.String()),
							End:   scalarutils.DateTime(time.March.String()),
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}
			}

			if tt.name == "Sad case: failed to get encounter" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Period: &domain.FHIRPeriod{
							ID:    &UUID,
							Start: scalarutils.DateTime(time.February.String()),
							End:   scalarutils.DateTime(time.March.String()),
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get update resource" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					episode := domain.FHIREpisodeOfCare{
						ID:     &UUID,
						Status: &status,
						Period: &domain.FHIRPeriod{
							ID:    &UUID,
							Start: scalarutils.DateTime(time.February.String()),
							End:   scalarutils.DateTime(time.March.String()),
						},
					}
					bs, err := json.Marshal(episode)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}
					return nil
				}
				dataset.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.EndEpisode(tt.args.ctx, tt.args.episodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.EndEpisode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.EndEpisode() = %v, want %v", got, tt.want)
			}
		})
	}
}
