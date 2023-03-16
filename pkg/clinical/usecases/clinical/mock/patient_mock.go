package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/scalarutils"
)

// FakeClinical ....
type FakeClinical struct {
	MockGetMedicalDataFn      func(ctx context.Context, patientID string) (*domain.MedicalData, error)
	MockCreatePubsubPatientFn func(ctx context.Context, payload dto.CreatePatientPubSubMessage) error
	MockCreatePatientFn       func(ctx context.Context, input dto.PatientInput) (*dto.Patient, error)
}

// NewFakeClinicalMock ...
func NewFakeClinicalMock() *FakeClinical {
	return &FakeClinical{
		MockGetMedicalDataFn: func(ctx context.Context, patientID string) (*domain.MedicalData, error) {
			return &domain.MedicalData{
				Regimen:   []*domain.FHIRMedicationStatement{},
				Allergies: []*domain.FHIRAllergyIntolerance{},
				Weight:    []*domain.FHIRObservation{},
				BMI:       []*domain.FHIRObservation{},
				ViralLoad: []*domain.FHIRObservation{},
				CD4Count:  []*domain.FHIRObservation{},
			}, nil
		},
		MockCreatePubsubPatientFn: func(ctx context.Context, payload dto.CreatePatientPubSubMessage) error {
			return nil
		},
		MockCreatePatientFn: func(ctx context.Context, input dto.PatientInput) (*dto.Patient, error) {
			return &dto.Patient{
				ID:          uuid.NewString(),
				Active:      true,
				Name:        gofakeit.BeerName(),
				PhoneNumber: []string{interserviceclient.TestUserPhoneNumber},
				Gender:      string(enumutils.GenderMale),
				BirthDate: scalarutils.Date{
					Year:  2000,
					Month: 3,
					Day:   20,
				},
			}, nil
		},
	}
}

// GetMedicalData ...
func (f *FakeClinical) GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error) {
	return f.MockGetMedicalDataFn(ctx, patientID)
}

// CreatePubsubPatient mocks the implementation os creating a user using pubsub
func (f *FakeClinical) CreatePubsubPatient(ctx context.Context, payload dto.CreatePatientPubSubMessage) error {
	return f.MockCreatePubsubPatientFn(ctx, payload)
}

// CreatePatient mocks the implementation of creating a patient
func (f *FakeClinical) CreatePatient(ctx context.Context, input dto.PatientInput) (*dto.Patient, error) {
	return f.MockCreatePatientFn(ctx, input)
}
