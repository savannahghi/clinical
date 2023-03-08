package mock

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// FakeClinical ....
type FakeClinical struct {
	MockRegisterPatientFn func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	MockFindPatientByIDFn func(ctx context.Context, id string) (*domain.PatientPayload, error)
	MockGetMedicalDataFn  func(ctx context.Context, patientID string) (*domain.MedicalData, error)
}

// NewFakeClinicalMock ...
func NewFakeClinicalMock() *FakeClinical {
	return &FakeClinical{

		MockRegisterPatientFn: func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
			return &domain.PatientPayload{
				PatientRecord:   &domain.FHIRPatient{},
				HasOpenEpisodes: false,
				OpenEpisodes:    []*domain.FHIREpisodeOfCare{},
			}, nil
		},

		MockFindPatientByIDFn: func(ctx context.Context, id string) (*domain.PatientPayload, error) {
			UUID := "FDGFSDG33222"
			PatientRef := fmt.Sprintf("Patient/%s", UUID)
			gender := domain.PatientContactGenderEnumMale
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID:                   &UUID,
					Text:                 &domain.FHIRNarrative{},
					Identifier:           []*domain.FHIRIdentifier{},
					Active:               new(bool),
					Name:                 []*domain.FHIRHumanName{},
					Telecom:              []*domain.FHIRContactPoint{},
					Gender:               (*domain.PatientGenderEnum)(&gender),
					BirthDate:            &scalarutils.Date{},
					DeceasedBoolean:      new(bool),
					DeceasedDateTime:     &scalarutils.Date{},
					Address:              []*domain.FHIRAddress{},
					MaritalStatus:        &domain.FHIRCodeableConcept{},
					MultipleBirthBoolean: new(bool),
					MultipleBirthInteger: new(string),
					Photo:                []*domain.FHIRAttachment{},
					Contact:              []*domain.FHIRPatientContact{},
					Communication:        []*domain.FHIRPatientCommunication{},
					GeneralPractitioner:  []*domain.FHIRReference{},
					ManagingOrganization: &domain.FHIRReference{},
					Link:                 []*domain.FHIRPatientLink{},
				},
				HasOpenEpisodes: false,
				OpenEpisodes: []*domain.FHIREpisodeOfCare{
					{
						ID:            &UUID,
						Text:          &domain.FHIRNarrative{},
						Identifier:    []*domain.FHIRIdentifier{},
						StatusHistory: []*domain.FHIREpisodeofcareStatushistory{},
						Type:          []*domain.FHIRCodeableConcept{},
						Diagnosis:     []*domain.FHIREpisodeofcareDiagnosis{},
						Patient: &domain.FHIRReference{
							ID:        &UUID,
							Reference: &PatientRef,
						},
						ManagingOrganization: &domain.FHIRReference{},
						Period:               &domain.FHIRPeriod{},
						ReferralRequest:      []*domain.FHIRReference{},
						CareManager:          &domain.FHIRReference{},
						Team:                 []*domain.FHIRReference{},
						Account:              []*domain.FHIRReference{},
					},
				},
			}, nil
		},
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
	}
}

// RegisterPatient ...
func (f *FakeClinical) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return f.MockRegisterPatientFn(ctx, input)
}

// FindPatientByID ...
func (f *FakeClinical) FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error) {
	return f.MockFindPatientByIDFn(ctx, id)
}

// GetMedicalData ...
func (f *FakeClinical) GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error) {
	return f.MockGetMedicalDataFn(ctx, patientID)
}
