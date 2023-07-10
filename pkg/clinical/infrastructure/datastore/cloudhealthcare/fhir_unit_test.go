package fhir_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	FHIR "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"

	fakeDataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset/mock"
)

func TestStoreImpl_SearchFHIRObservation(t *testing.T) {

	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: search observation",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"_sort": "-date",
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: search resource error",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"_sort": "-date",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "sad case: search resource error" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to search fhir resource")
				}
			}

			got, err := fh.SearchFHIRObservation(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRObservation(t *testing.T) {

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
			name: "happy case: delete observation",
			args: args{
				ctx: context.Background(),
				id:  uuid.NewString(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: error deleting resource",
			args: args{
				ctx: context.Background(),
				id:  uuid.NewString(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "sad case: error deleting resource" {
				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					return fmt.Errorf("failed to delete resource")
				}
			}

			got, err := fh.DeleteFHIRObservation(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.DeleteFHIRObservation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_CreateFHIRObservation(t *testing.T) {

	type args struct {
		ctx   context.Context
		input domain.FHIRObservationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create fhir observation",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRObservationInput{
					Code: domain.FHIRCodeableConceptInput{
						Text: "Obs",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating resource",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRObservationInput{
					Code: domain.FHIRCodeableConceptInput{
						Text: "Obs",
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

			if tt.name == "sad case: error creating resource" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to create resource")
				}
			}

			got, err := fh.CreateFHIRObservation(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRPatient(t *testing.T) {

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
			name: "happy case: get patient",
			args: args{
				ctx: context.Background(),
				id:  uuid.NewString(),
			},
			wantErr: false,
		},
		{
			name: "sad case: error retrieving fhir resource",
			args: args{
				ctx: context.Background(),
				id:  uuid.NewString(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "happy case: get patient" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					id := gofakeit.UUID()
					patient := &domain.FHIRPatient{
						ID: &id,
					}

					bs, err := json.Marshal(patient)
					if err != nil {
						return err
					}

					err = json.Unmarshal(bs, resource)
					if err != nil {
						return err
					}

					return nil
				}

				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					episode := domain.FHIREpisodeOfCare{
						Period: &domain.FHIRPeriod{
							Start: "2020-09-24T18:02:38.661033Z",
							End:   "2020-09-24T18:02:38.661033Z",
						},
					}

					payload, err := converterandformatter.StructToMap(episode)
					if err != nil {
						return nil, err
					}

					return &domain.PagedFHIRResource{
						Resources: []map[string]interface{}{
							payload,
						},
					}, nil
				}
			}

			if tt.name == "sad case: error retrieving fhir resource" {
				dataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					return fmt.Errorf("failed to get resource")
				}
			}

			got, err := fh.GetFHIRPatient(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRResourceType(t *testing.T) {

	type args struct {
		results []map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: delete resource",
			args: args{
				results: []map[string]string{
					{"Patient": gofakeit.UUID()},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: delete resource error",
			args: args{
				results: []map[string]string{
					{"Patient": gofakeit.UUID()},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "sad case: delete resource error" {
				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					return fmt.Errorf("failed to delete resource")
				}
			}

			if err := fh.DeleteFHIRResourceType(tt.args.results); (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRResourceType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRServiceRequest(t *testing.T) {

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
			name: "happy case: delete service request",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: delete resource error",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "sad case: delete resource error" {
				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					return fmt.Errorf("failed to delete resource")
				}
			}

			got, err := fh.DeleteFHIRServiceRequest(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.DeleteFHIRServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreImpl_CreateFHIRMedicationStatement(t *testing.T) {

	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationStatementInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create medication statement",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRMedicationStatementInput{
					Category: &domain.FHIRCodeableConceptInput{
						Text: "dawa",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating resource",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRMedicationStatementInput{
					Category: &domain.FHIRCodeableConceptInput{
						Text: "dawa",
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

			if tt.name == "sad case: error creating resource" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to create resource")
				}
			}

			got, err := fh.CreateFHIRMedicationStatement(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRMedicationStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRMedication(t *testing.T) {

	type args struct {
		ctx   context.Context
		input domain.FHIRMedicationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create medication",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRMedicationInput{
					Code: &domain.FHIRCodeableConceptInput{
						Text: "ARV",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating resource",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRMedicationInput{
					Code: &domain.FHIRCodeableConceptInput{
						Text: "ARV",
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

			if tt.name == "sad case: error creating resource" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to create resource")
				}
			}

			got, err := fh.CreateFHIRMedication(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRMedication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRPatient(t *testing.T) {

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
			name: "happy case: create patient",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRPatientInput{
					MaritalStatus: &domain.FHIRCodeableConceptInput{
						Text: "single",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating resource",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRPatientInput{
					MaritalStatus: &domain.FHIRCodeableConceptInput{
						Text: "single",
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

			if tt.name == "sad case: error creating resource" {
				dataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to create resource")
				}
			}

			got, err := fh.CreateFHIRPatient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_PatchFHIRPatient(t *testing.T) {

	type args struct {
		ctx   context.Context
		id    string
		input domain.FHIRPatientInput
	}
	is_active := false
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: patch patient",
			args: args{
				ctx:   context.Background(),
				id:    gofakeit.UUID(),
				input: domain.FHIRPatientInput{Active: &is_active},
			},
			wantErr: false,
		},
		{
			name: "sad case: error patching resource",
			args: args{
				ctx:   context.Background(),
				id:    gofakeit.UUID(),
				input: domain.FHIRPatientInput{Active: &is_active},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "sad case: error patching resource" {
				dataset.MockPatchFHIRResourceFn = func(resourceType string, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed to patch resource")
				}
			}

			got, err := fh.PatchFHIRPatient(tt.args.ctx, tt.args.id, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.PatchFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_UpdateFHIREpisodeOfCare(t *testing.T) {

	type args struct {
		ctx            context.Context
		fhirResourceID string
		payload        map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: update episode of care",
			args: args{
				ctx:            context.Background(),
				fhirResourceID: gofakeit.UUID(),
				payload: map[string]interface{}{
					"episode": "one",
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error updating resource",
			args: args{
				ctx:            context.Background(),
				fhirResourceID: gofakeit.UUID(),
				payload: map[string]interface{}{
					"episode": "one",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "sad case: error updating resource" {
				dataset.MockUpdateFHIRResourceFn = func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("failed ro update resource")
				}
			}

			got, err := fh.UpdateFHIREpisodeOfCare(tt.args.ctx, tt.args.fhirResourceID, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.UpdateFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_SearchFHIRMedicationStatement(t *testing.T) {
	type args struct {
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: search medication statement",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"name": "ARVs",
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: search resource error",
			args: args{
				ctx: context.Background(),
				params: map[string]interface{}{
					"name": "ARVs",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "sad case: search resource error" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to search")
				}
			}

			got, err := fh.SearchFHIRMedicationStatement(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRMedicationStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func fakePatient() (*domain.FHIRPatient, error) {
	id := gofakeit.UUID()
	system := scalarutils.URI("test")
	version := "0.0.1"
	userSelected := false
	active := true
	phoneSystem := domain.ContactPointSystemEnumPhone
	use := domain.ContactPointUseEnumHome
	rank := int64(1)
	phone := gofakeit.Phone()
	date, err := scalarutils.NewDate(12, 12, 2000)
	if err != nil {
		return nil, err
	}
	male := domain.PatientGenderEnumMale
	maleContact := domain.PatientContactGenderEnumMale
	address := gofakeit.Address()
	addrUse := domain.AddressUseEnumHome
	country := helpers.DefaultCountry
	postalAddrType := domain.AddressTypeEnumPostal
	name := gofakeit.Name()
	nameUse := domain.HumanNameUseEnumOfficial

	creation := scalarutils.DateTime("2020-09-24T18:02:38.661033Z")

	patient := domain.FHIRPatient{
		ID: &id,
		Identifier: []*domain.FHIRIdentifier{
			{
				ID:  &id,
				Use: domain.IdentifierUseEnumOfficial,
				Type: domain.FHIRCodeableConcept{
					Text: "MR",
					Coding: []*domain.FHIRCoding{
						{
							System:       &system,
							Version:      &version,
							Code:         scalarutils.Code(id),
							Display:      id,
							UserSelected: &userSelected,
						},
					},
				},
				System:   &system,
				Value:    id,
				Period:   common.DefaultPeriod(),
				Assigner: &domain.FHIRReference{},
			},
		},
		Active: &active,
		Name: []*domain.FHIRHumanName{
			{
				Given:  []*string{&name},
				Family: &name,
				Use:    nameUse,
				Period: common.DefaultPeriod(),
				Text:   name,
			},
		},
		Telecom: []*domain.FHIRContactPoint{
			{
				System: &phoneSystem,
				Use:    &use,
				Rank:   &rank,
				Period: common.DefaultPeriod(),
				Value:  &phone,
			},
		},
		Gender:    &male,
		BirthDate: date,
		Address: []*domain.FHIRAddress{
			{
				Use:     &addrUse,
				Type:    &postalAddrType,
				Country: &country,
				Period:  common.DefaultPeriod(),
				Line:    []*string{&address.Address},
				Text:    address.Address,
			},
		},
		Photo: []*domain.FHIRAttachment{
			{
				ID:       &id,
				Creation: &creation,
			},
		},
		Contact: []*domain.FHIRPatientContact{
			{
				ID:           new(string),
				Relationship: []*domain.FHIRCodeableConcept{},
				Name: &domain.FHIRHumanName{
					Given:  []*string{&name},
					Family: &name,
					Use:    nameUse,
					Period: common.DefaultPeriod(),
					Text:   name,
				},
				Telecom: []*domain.FHIRContactPoint{
					{
						System: &phoneSystem,
						Use:    &use,
						Rank:   &rank,
						Period: common.DefaultPeriod(),
						Value:  &phone,
					},
				},
				Address: &domain.FHIRAddress{

					Use:     &addrUse,
					Type:    &postalAddrType,
					Country: &country,
					Period:  common.DefaultPeriod(),
					Line:    []*string{&address.Address},
					Text:    address.Address,
				},
				Gender: &maleContact,
				Period: common.DefaultPeriod(),
			},
		},
		MaritalStatus: &domain.FHIRCodeableConcept{
			Coding: []*domain.FHIRCoding{
				{
					Code:         scalarutils.Code(domain.MaritalStatusA.String()),
					Display:      domain.MaritalStatusDisplay(domain.MaritalStatusA),
					UserSelected: &userSelected,
				},
			},
			Text: domain.MaritalStatusDisplay(domain.MaritalStatusA),
		},
	}

	return &patient, nil
}

func TestStoreImpl_SearchFHIRPatient(t *testing.T) {

	type args struct {
		ctx          context.Context
		searchParams string
		tenant       dto.TenantIdentifiers
		pagination   dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: search patient",
			args: args{
				ctx:          context.Background(),
				searchParams: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: search patient error",
			args: args{
				ctx:          context.Background(),
				searchParams: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "happy case: search patient" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					var payload map[string]interface{}

					switch resourceType {
					case "Patient":
						patient, err := fakePatient()
						if err != nil {
							return nil, err
						}

						payload, err = converterandformatter.StructToMap(patient)
						if err != nil {
							return nil, err
						}
					case "EpisodeOfCare":

						episode := domain.FHIREpisodeOfCare{
							Period: &domain.FHIRPeriod{
								Start: "2020-09-24T18:02:38.661033Z",
								End:   "2020-09-24T18:02:38.661033Z",
							},
						}

						p, err := converterandformatter.StructToMap(episode)
						if err != nil {
							return nil, err
						}

						payload = p
					}

					return &domain.PagedFHIRResource{
						Resources: []map[string]interface{}{
							payload,
						},
					}, nil
				}
			}

			if tt.name == "sad case: search patient error" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					var payload map[string]interface{}

					switch resourceType {
					case "Patient":
						return nil, fmt.Errorf("failed to find patient")
					case "EpisodeOfCare":

						episode := domain.FHIREpisodeOfCare{
							Period: &domain.FHIRPeriod{
								Start: "2020-09-24T18:02:38.661033Z",
								End:   "2020-09-24T18:02:38.661033Z",
							},
						}

						p, err := converterandformatter.StructToMap(episode)
						if err != nil {
							return nil, err
						}

						payload = p
					}

					return &domain.PagedFHIRResource{
						Resources: []map[string]interface{}{
							payload,
						},
					}, nil
				}
			}

			got, err := fh.SearchFHIRPatient(tt.args.ctx, tt.args.searchParams, tt.args.tenant, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestStoreImpl_DeleteFHIRPatient(t *testing.T) {

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
			name: "happy case: delete all patient data",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: all patient data error",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: all patient data invalid entry",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: all patient data invalid entry type",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: all patient data entry invalid resource type",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: error deleting medication request",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting encounters",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting episode of care",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting observation",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting patient",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: error deleting other types",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "happy case: delete all patient data" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": []map[string]interface{}{
							{
								"resource": map[string]interface{}{
									"resourceType": "EpisodeOfCare",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "Observation",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "AllergyIntolerance",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "ServiceRequest",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "MedicationRequest",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "Condition",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "Encounter",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "Composition",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "MedicationStatement",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "Medication",
									"id":           gofakeit.UUID(),
								},
							},
							{
								"resource": map[string]interface{}{
									"resourceType": "Patient",
									"id":           gofakeit.UUID(),
								},
							},
						},
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}
			}

			if tt.name == "sad case: all patient data error" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					return nil, fmt.Errorf("failed to get data")
				}
			}

			if tt.name == "sad case: all patient data invalid entry" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": "invalid",
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}
			}

			if tt.name == "sad case: all patient data invalid entry type" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": []map[int]string{
							{
								1: "bad",
							},
						},
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}
			}

			if tt.name == "sad case: all patient data entry invalid resource type" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": []map[string]interface{}{

							{
								"resource": "invalid",
							},
						},
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}
			}

			if tt.name == "sad case: error deleting medication request" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": []map[string]interface{}{
							{
								"resource": map[string]interface{}{
									"resourceType": "MedicationRequest",
									"id":           gofakeit.UUID(),
								},
							},
						},
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}

				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					if resourceType == "MedicationRequest" {
						return fmt.Errorf("failed")
					}
					return nil
				}
			}

			if tt.name == "sad case: error deleting other types" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": []map[string]interface{}{
							{
								"resource": map[string]interface{}{
									"resourceType": "Composition",
									"id":           gofakeit.UUID(),
								},
							},
						},
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}

				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					if resourceType == "Composition" {
						return fmt.Errorf("failed")
					}
					return nil
				}
			}

			if tt.name == "sad case: error deleting patient" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": []map[string]interface{}{
							{
								"resource": map[string]interface{}{
									"resourceType": "Patient",
									"id":           gofakeit.UUID(),
								},
							},
						},
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}

				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					if resourceType == "Patient" {
						return fmt.Errorf("failed")
					}
					return nil
				}
			}

			if tt.name == "sad case: error deleting observation" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": []map[string]interface{}{
							{
								"resource": map[string]interface{}{
									"resourceType": "Observation",
									"id":           gofakeit.UUID(),
								},
							},
						},
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}

				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					if resourceType == "Observation" {
						return fmt.Errorf("failed")
					}
					return nil
				}
			}

			if tt.name == "sad case: error deleting encounters" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": []map[string]interface{}{

							{
								"resource": map[string]interface{}{
									"resourceType": "Encounter",
									"id":           gofakeit.UUID(),
								},
							},
						},
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}

				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					if resourceType == "Encounter" {
						return fmt.Errorf("failed")
					}
					return nil
				}
			}

			if tt.name == "sad case: error deleting episode of care" {
				dataset.MockGetFHIRPatientAllDataFn = func(fhirResourceID string) ([]byte, error) {
					data := map[string]interface{}{
						"entry": []map[string]interface{}{
							{
								"resource": map[string]interface{}{
									"resourceType": "EpisodeOfCare",
									"id":           gofakeit.UUID(),
								},
							},
						},
					}

					bs, err := json.Marshal(data)
					if err != nil {
						return nil, err
					}

					return bs, err
				}

				dataset.MockDeleteFHIRResourceFn = func(resourceType, fhirResourceID string) error {
					if resourceType == "EpisodeOfCare" {
						return fmt.Errorf("failed")
					}
					return nil
				}
			}

			got, err := fh.DeleteFHIRPatient(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.DeleteFHIRPatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StoreImpl.DeleteFHIRPatient() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestStoreImpl_GetFHIROrganization_Unittest(t *testing.T) {

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

			got, err := fh.GetFHIROrganization(tt.args.ctx, tt.args.organizationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FHIRUseCaseImpl.GetFHIROrganization() error = %v, wantErr %v", err, tt.wantErr)
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
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRServiceRequest(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
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
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.PagedFHIRAllergy
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRAllergyIntolerance(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
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
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRComposition(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
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
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.PagedFHIRCondition
		wantErr bool
	}{
		{
			name: "Happy Case - successfully search fhir condition",
			args: args{
				ctx: ctx,
				params: map[string]interface{}{
					"id": "1234",
				},
				pagination: dto.Pagination{},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search a condition",
			args: args{
				ctx:        ctx,
				pagination: dto.Pagination{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to search a condition" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRCondition(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
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
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIREncounter(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
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
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to search resource")
				}
			}

			got, err := fh.SearchFHIRMedicationRequest(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
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

func TestStoreImpl_SearchPatientEncounters(t *testing.T) {
	status := domain.EncounterStatusEnumPlanned
	type args struct {
		ctx              context.Context
		patientReference string
		status           *domain.EncounterStatusEnum
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.SearchPatientEncounters(tt.args.ctx, tt.args.patientReference, tt.args.status, tt.args.tenant, tt.args.pagination)
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
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.SearchFHIREpisodeOfCare(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
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
		episode domain.FHIREpisodeOfCareInput
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
				episode: domain.FHIREpisodeOfCareInput{
					ID:     &UUID,
					Status: &status,
					Patient: &domain.FHIRReferenceInput{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReferenceInput{
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
				episode: domain.FHIREpisodeOfCareInput{
					ID:     &UUID,
					Status: &status,
					Patient: &domain.FHIRReferenceInput{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReferenceInput{
						Reference: &OrgRef,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to create FHIR resource",
			args: args{
				ctx: context.Background(),
				episode: domain.FHIREpisodeOfCareInput{
					ID:     &UUID,
					Status: &status,
					Patient: &domain.FHIRReferenceInput{
						Reference: &PatientRef,
					},
					ManagingOrganization: &domain.FHIRReferenceInput{
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, nil
				}
			}

			if tt.name == "Sad case: failed to create FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
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
		ctx        context.Context
		params     map[string]interface{}
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := fh.SearchFHIROrganization(tt.args.ctx, tt.args.params, tt.args.tenant, tt.args.pagination)
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
		tenant       dto.TenantIdentifiers
		pagination   dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.SearchEpisodesByParam(tt.args.ctx, tt.args.searchParams, tt.args.tenant, tt.args.pagination)
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
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.OpenEpisodes(tt.args.ctx, tt.args.patientReference, tt.args.tenant, tt.args.pagination)
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
		ctx        context.Context
		patient    domain.FHIRPatient
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := fh.HasOpenEpisode(tt.args.ctx, tt.args.patient, tt.args.tenant, tt.args.pagination)
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
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := fh.SearchEpisodeEncounter(tt.args.ctx, tt.args.episodeReference, tt.args.tenant, tt.args.pagination)
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
		ctx        context.Context
		episodeID  string
		tenant     dto.TenantIdentifiers
		pagination dto.Pagination
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
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: empty FHIR resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return &domain.PagedFHIRResource{
						Resources: []map[string]interface{}{},
					}, nil
				}
			}
			got, err := fh.GetActiveEpisode(tt.args.ctx, tt.args.episodeID, tt.args.tenant, tt.args.pagination)
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

func TestStoreImpl_SearchPatientObservations(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx              context.Context
		patientReference string
		observationCode  string
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully search patient observation",
			args: args{
				ctx:              ctx,
				patientReference: fmt.Sprintf("Patient/%s", gofakeit.UUID()),
				observationCode:  "5088",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - fail to search fhir resource",
			args: args{
				ctx:              ctx,
				patientReference: fmt.Sprintf("Patient/%s", gofakeit.UUID()),
				observationCode:  "5088",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(dataset)

			if tt.name == "Sad Case - fail to search fhir resource" {
				dataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, fmt.Errorf("failed to search observation resource")
				}
			}

			got, err := fh.SearchPatientObservations(tt.args.ctx, tt.args.patientReference, tt.args.observationCode, tt.args.tenant, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchPatientObservations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestStoreImpl_GetFHIRAllergyIntolerance(t *testing.T) {
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
			name: "Happy case: get allergy intolerance by ID",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get allergy intolerance by ID",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(fakeDataset)

			if tt.name == "Sad case: unable to get allergy intolerance by ID" {
				fakeDataset.MockGetFHIRResourceFn = func(resourceType, fhirResourceID string, resource interface{}) error {
					return fmt.Errorf("error")
				}
			}

			_, err := fh.GetFHIRAllergyIntolerance(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.GetFHIRAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchPatientAllergyIntolerance(t *testing.T) {
	patientReference := fmt.Sprintf("Patient/%s", gofakeit.UUID())
	type args struct {
		ctx              context.Context
		patientReference string
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: search allergy intolerance",
			args: args{
				ctx:              context.Background(),
				patientReference: patientReference,
				tenant: dto.TenantIdentifiers{
					OrganizationID: gofakeit.UUID(),
					FacilityID:     gofakeit.UUID(),
				},
				pagination: dto.Pagination{},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search allergy intolerance",
			args: args{
				ctx:              context.Background(),
				patientReference: patientReference,
				tenant: dto.TenantIdentifiers{
					OrganizationID: gofakeit.UUID(),
					FacilityID:     gofakeit.UUID(),
				},
				pagination: dto.Pagination{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(fakeDataset)

			if tt.name == "Sad case: unable to search allergy intolerance" {
				fakeDataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, errors.New("some error")
				}
			}

			_, err := fh.SearchPatientAllergyIntolerance(tt.args.ctx, tt.args.patientReference, tt.args.tenant, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchPatientAllergyIntolerance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_CreateFHIRMedia(t *testing.T) {
	id := uuid.New().String()
	type args struct {
		ctx   context.Context
		input domain.FHIRMedia
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create FHIR media",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRMedia{
					Subject: &domain.FHIRReferenceInput{
						ID: &id,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create FHIR media",
			args: args{
				ctx: context.Background(),
				input: domain.FHIRMedia{
					Subject: &domain.FHIRReferenceInput{
						ID: &id,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(fakeDataset)

			if tt.name == "Sad case: unable to create FHIR media" {
				fakeDataset.MockCreateFHIRResourceFn = func(resourceType string, payload map[string]interface{}, resource interface{}) error {
					return fmt.Errorf("an error ocurred")
				}
			}

			_, err := fh.CreateFHIRMedia(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.CreateFHIRMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStoreImpl_SearchPatientMedia(t *testing.T) {
	first := 1
	type args struct {
		ctx              context.Context
		patientReference string
		tenant           dto.TenantIdentifiers
		pagination       dto.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: search patient media",
			args: args{
				ctx:              context.Background(),
				patientReference: gofakeit.BeerIbu(),
				tenant: dto.TenantIdentifiers{
					OrganizationID: gofakeit.UUID(),
					FacilityID:     gofakeit.UUID(),
				},
				pagination: dto.Pagination{
					First: &first,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search patient media",
			args: args{
				ctx:              context.Background(),
				patientReference: gofakeit.BeerIbu(),
				tenant: dto.TenantIdentifiers{
					OrganizationID: gofakeit.UUID(),
					FacilityID:     gofakeit.UUID(),
				},
				pagination: dto.Pagination{
					First: &first,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDataset := fakeDataset.NewFakeFHIRRepositoryMock()
			fh := FHIR.NewFHIRStoreImpl(fakeDataset)

			if tt.name == "Sad case: unable to search patient media" {
				fakeDataset.MockSearchFHIRResourceFn = func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
					return nil, errors.New("unable to search patient media")
				}
			}

			_, err := fh.SearchPatientMedia(tt.args.ctx, tt.args.patientReference, tt.args.tenant, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreImpl.SearchPatentMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
