package mock

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FakeClinical ....
type FakeClinical struct {
	MockProblemSummaryFn                      func(ctx context.Context, patientID string) ([]string, error)
	MockVisitSummaryFn                        func(ctx context.Context, encounterID string, count int) (map[string]interface{}, error)
	MockPatientTimelineWithCountFn            func(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error)
	MockContactsToContactPointInputFn         func(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error)
	MockRegisterPatientFn                     func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	MockCreatePatientFn                       func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error)
	MockFindPatientByIDFn                     func(ctx context.Context, id string) (*domain.PatientPayload, error)
	MockPatientSearchFn                       func(ctx context.Context, search string) (*domain.PatientConnection, error)
	MockUpdatePatientFn                       func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	MockAddNextOfKinFn                        func(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error)
	MockCreateUpdatePatientExtraInformationFn func(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error)
	MockAllergySummaryFn                      func(ctx context.Context, patientID string) ([]string, error)
	MockDeleteFHIRPatientByPhoneFn            func(ctx context.Context, phoneNumber string) (bool, error)
	MockStartEpisodeByBreakGlassFn            func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error)
	MockFindPatientsByMSISDNFn                func(ctx context.Context, msisdn string) (*domain.PatientConnection, error)
	MockPatientTimelineFn                     func(ctx context.Context, patientID string, count int) ([]map[string]interface{}, error)
	MockGetMedicalDataFn                      func(ctx context.Context, patientID string) (*domain.MedicalData, error)
}

// NewFakeClinicalMock ...
func NewFakeClinicalMock() *FakeClinical {
	UUID := uuid.New().String()
	return &FakeClinical{
		MockProblemSummaryFn: func(ctx context.Context, patientID string) ([]string, error) { return nil, nil },
		MockVisitSummaryFn: func(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {

			encounterRef := fmt.Sprintf("Encounter/%s", uuid.New().String())
			m := map[string]interface{}{
				"key":       "value",
				"encounter": encounterRef,
			}
			return m, nil
		},
		MockPatientTimelineWithCountFn: func(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
			m := map[string]interface{}{
				"key":    "value",
				"active": true,
			}
			return []map[string]interface{}{m}, nil
		},
		MockContactsToContactPointInputFn: func(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
			return []*domain.FHIRContactPointInput{
				{
					ID:     &UUID,
					Value:  new(string),
					Rank:   new(int64),
					Period: &domain.FHIRPeriodInput{},
				},
			}, nil
		},
		MockRegisterPatientFn: func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
			return &domain.PatientPayload{
				PatientRecord:   &domain.FHIRPatient{},
				HasOpenEpisodes: false,
				OpenEpisodes:    []*domain.FHIREpisodeOfCare{},
			}, nil
		},
		MockCreatePatientFn: func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
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
		MockPatientSearchFn: func(ctx context.Context, search string) (*domain.PatientConnection, error) {
			return &domain.PatientConnection{
				Edges:    []*domain.PatientEdge{},
				PageInfo: &firebasetools.PageInfo{},
			}, nil
		},
		MockUpdatePatientFn: func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
			UUID := uuid.New().String()
			g := domain.PatientContactGenderEnumFemale
			PatientRef := fmt.Sprintf("Patient/%s", UUID)
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID:                   &UUID,
					Text:                 &domain.FHIRNarrative{},
					Identifier:           []*domain.FHIRIdentifier{},
					Active:               new(bool),
					Name:                 []*domain.FHIRHumanName{},
					Telecom:              []*domain.FHIRContactPoint{},
					Gender:               (*domain.PatientGenderEnum)(&g),
					BirthDate:            &scalarutils.Date{},
					DeceasedBoolean:      new(bool),
					DeceasedDateTime:     &scalarutils.Date{},
					Address:              []*domain.FHIRAddress{},
					MaritalStatus:        &domain.FHIRCodeableConcept{},
					MultipleBirthBoolean: new(bool),
					MultipleBirthInteger: &UUID,
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
						ID:            new(string),
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
		MockAddNextOfKinFn: func(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
			return &domain.PatientPayload{
				PatientRecord:   &domain.FHIRPatient{},
				HasOpenEpisodes: false,
				OpenEpisodes:    []*domain.FHIREpisodeOfCare{},
			}, nil
		},
		MockCreateUpdatePatientExtraInformationFn: func(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) { return true, nil },
		MockAllergySummaryFn: func(ctx context.Context, patientID string) ([]string, error) {
			return []string{"test"}, nil
		},
		MockDeleteFHIRPatientByPhoneFn: func(ctx context.Context, phoneNumber string) (bool, error) { return true, nil },
		MockStartEpisodeByBreakGlassFn: func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
			return &domain.EpisodeOfCarePayload{
				EpisodeOfCare: &domain.FHIREpisodeOfCare{},
				TotalVisits:   0,
			}, nil
		},
		MockFindPatientsByMSISDNFn: func(ctx context.Context, msisdn string) (*domain.PatientConnection, error) {
			return &domain.PatientConnection{
				Edges:    []*domain.PatientEdge{},
				PageInfo: &firebasetools.PageInfo{},
			}, nil
		},
		MockPatientTimelineFn: func(ctx context.Context, patientID string, count int) ([]map[string]interface{}, error) {
			m := map[string]interface{}{
				"key": "value",
			}
			return []map[string]interface{}{m}, nil
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

// ProblemSummary ...
func (f *FakeClinical) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	return f.MockProblemSummaryFn(ctx, patientID)
}

// VisitSummary ...
func (f *FakeClinical) VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
	return f.MockVisitSummaryFn(ctx, encounterID, count)
}

// PatientTimelineWithCount ...
func (f *FakeClinical) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	return f.MockPatientTimelineWithCountFn(ctx, episodeID, count)
}

// ContactsToContactPointInput ...
func (f *FakeClinical) ContactsToContactPointInput(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
	return f.MockContactsToContactPointInputFn(ctx, phones, emails)
}

// RegisterPatient ...
func (f *FakeClinical) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return f.MockRegisterPatientFn(ctx, input)
}

// CreatePatient ...
func (f *FakeClinical) CreatePatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	return f.MockCreatePatientFn(ctx, input)
}

// FindPatientByID ...
func (f *FakeClinical) FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error) {
	return f.MockFindPatientByIDFn(ctx, id)
}

// PatientSearch ...
func (f *FakeClinical) PatientSearch(ctx context.Context, search string) (*domain.PatientConnection, error) {
	return f.MockPatientSearchFn(ctx, search)
}

// UpdatePatient ...
func (f *FakeClinical) UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return f.MockUpdatePatientFn(ctx, input)
}

// AddNextOfKin ...
func (f *FakeClinical) AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
	return f.MockAddNextOfKinFn(ctx, input)
}

// CreateUpdatePatientExtraInformation ...
func (f *FakeClinical) CreateUpdatePatientExtraInformation(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
	return f.MockCreateUpdatePatientExtraInformationFn(ctx, input)
}

// AllergySummary ...
func (f *FakeClinical) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	return f.MockAllergySummaryFn(ctx, patientID)
}

// DeleteFHIRPatientByPhone ...
func (f *FakeClinical) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	return f.MockDeleteFHIRPatientByPhoneFn(ctx, phoneNumber)
}

// StartEpisodeByBreakGlass ...
func (f *FakeClinical) StartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	return f.MockStartEpisodeByBreakGlassFn(ctx, input)
}

// FindPatientsByMSISDN ...
func (f *FakeClinical) FindPatientsByMSISDN(ctx context.Context, msisdn string) (*domain.PatientConnection, error) {
	return f.MockFindPatientsByMSISDNFn(ctx, msisdn)
}

// PatientTimeline ...
func (f *FakeClinical) PatientTimeline(ctx context.Context, patientID string, count int) ([]map[string]interface{}, error) {
	return f.MockPatientTimelineFn(ctx, patientID, count)
}

// GetMedicalData ...
func (f *FakeClinical) GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error) {
	return f.MockGetMedicalDataFn(ctx, patientID)
}
