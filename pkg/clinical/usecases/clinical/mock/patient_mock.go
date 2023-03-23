package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// FakeClinical ....
type FakeClinical struct {
	MockGetMedicalDataFn      func(ctx context.Context, patientID string) (*dto.MedicalData, error)
	MockCreatePubsubPatientFn func(ctx context.Context, payload dto.CreatePatientPubSubMessage) error
}

// NewFakeClinicalMock ...
func NewFakeClinicalMock() *FakeClinical {
	return &FakeClinical{
		MockGetMedicalDataFn: func(ctx context.Context, patientID string) (*dto.MedicalData, error) {
			return &dto.MedicalData{
				Regimen:   []*dto.MedicationStatement{},
				Allergies: []*dto.Allergy{},
				Weight:    []*dto.Observation{},
				BMI:       []*dto.Observation{},
				ViralLoad: []*dto.Observation{},
				CD4Count:  []*dto.Observation{},
			}, nil
		},
		MockCreatePubsubPatientFn: func(ctx context.Context, payload dto.CreatePatientPubSubMessage) error {
			return nil
		},
	}
}

// GetMedicalData ...
func (f *FakeClinical) GetMedicalData(ctx context.Context, patientID string) (*dto.MedicalData, error) {
	return f.MockGetMedicalDataFn(ctx, patientID)
}

// CreatePubsubPatient mocks the implementation os creating a user using pubsub
func (f *FakeClinical) CreatePubsubPatient(ctx context.Context, payload dto.CreatePatientPubSubMessage) error {
	return f.MockCreatePubsubPatientFn(ctx, payload)
}
