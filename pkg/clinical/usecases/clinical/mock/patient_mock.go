package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FakeClinical ....
type FakeClinical struct {
	MockGetMedicalDataFn      func(ctx context.Context, patientID string) (*domain.MedicalData, error)
	MockCreatePubsubPatientFn func(ctx context.Context, payload dto.CreatePatientPubSubMessage) error
}

// NewFakeClinicalMock ...
func NewFakeClinicalMock() *FakeClinical {
	return &FakeClinical{
		MockGetMedicalDataFn: func(ctx context.Context, patientID string) (*domain.MedicalData, error) {
			return &domain.MedicalData{
				Regimen:   []*domain.FHIRMedicationStatement{},
				Allergies: []*dto.AllergyIntolerance{},
				Weight:    []*domain.FHIRObservation{},
				BMI:       []*domain.FHIRObservation{},
				ViralLoad: []*domain.FHIRObservation{},
				CD4Count:  []*domain.FHIRObservation{},
			}, nil
		},
		MockCreatePubsubPatientFn: func(ctx context.Context, payload dto.CreatePatientPubSubMessage) error {
			return nil
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
