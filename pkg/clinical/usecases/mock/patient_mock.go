package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/firebasetools"
	"github.com/segmentio/ksuid"
)

// ClinicalMock instantiates all the mock functions
type ClinicalMock struct {
	ProblemSummaryFn                      func(ctx context.Context, patientID string) ([]string, error)
	VisitSummaryFn                        func(ctx context.Context, encounterID string, count int) (map[string]interface{}, error)
	PatientTimelineWithCountFn            func(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error)
	ContactsToContactPointInputFn         func(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error)
	RegisterPatientFn                     func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	CreatePatientFn                       func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error)
	FindPatientByIDFn                     func(ctx context.Context, id string) (*domain.PatientPayload, error)
	PatientSearchFn                       func(ctx context.Context, search string) (*domain.PatientConnection, error)
	UpdatePatientFn                       func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	AddNextOfKinFn                        func(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error)
	AddNHIFFn                             func(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error)
	RegisterUserFn                        func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	CreateUpdatePatientExtraInformationFn func(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error)
	AllergySummaryFn                      func(ctx context.Context, patientID string) ([]string, error)
	DeleteFHIRPatientByPhoneFn            func(ctx context.Context, phoneNumber string) (bool, error)
	StartEpisodeByBreakGlassFn            func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error)
}

// NewClinicalMock is a new instance of NewClinicalMock
func NewClinicalMock() *ClinicalMock {
	return &ClinicalMock{
		ProblemSummaryFn: func(ctx context.Context, patientID string) ([]string, error) {
			return []string{"Sick"}, nil
		},

		VisitSummaryFn: func(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
			return map[string]interface{}{
				"Condition": "Heart",
			}, nil
		},

		PatientTimelineWithCountFn: func(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
			return []map[string]interface{}{
				1: {
					"Test": "Test",
				},
			}, nil
		},

		ContactsToContactPointInputFn: func(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
			return []*domain.FHIRContactPointInput{}, nil
		},

		RegisterPatientFn: func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
			id := ksuid.New().String()
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID: &id,
				},
			}, nil
		},

		CreatePatientFn: func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
			id := ksuid.New().String()
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID: &id,
				},
			}, nil
		},

		FindPatientByIDFn: func(ctx context.Context, id string) (*domain.PatientPayload, error) {
			patientID := ksuid.New().String()
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID: &patientID,
				},
			}, nil
		},

		PatientSearchFn: func(ctx context.Context, search string) (*domain.PatientConnection, error) {
			return &domain.PatientConnection{
				Edges: []*domain.PatientEdge{},
				PageInfo: &firebasetools.PageInfo{
					HasNextPage:     true,
					HasPreviousPage: true,
				},
			}, nil
		},

		UpdatePatientFn: func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
			ID := ksuid.New().String()
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID: &ID,
				},
			}, nil
		},

		AddNextOfKinFn: func(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
			ID := ksuid.New().String()
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID: &ID,
				},
			}, nil
		},

		AddNHIFFn: func(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
			ID := ksuid.New().String()
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID: &ID,
				},
			}, nil
		},

		RegisterUserFn: func(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
			ID := ksuid.New().String()
			return &domain.PatientPayload{
				PatientRecord: &domain.FHIRPatient{
					ID: &ID,
				},
			}, nil
		},

		CreateUpdatePatientExtraInformationFn: func(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
			return true, nil
		},

		AllergySummaryFn: func(ctx context.Context, patientID string) ([]string, error) {
			return []string{"Test Allergy"}, nil
		},

		DeleteFHIRPatientByPhoneFn: func(ctx context.Context, phoneNumber string) (bool, error) {
			return true, nil
		},

		StartEpisodeByBreakGlassFn: func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
			return &domain.EpisodeOfCarePayload{
				TotalVisits: 5,
			}, nil
		},
	}
}

//ProblemSummary is the ProblemSummary mock
func (p *ClinicalMock) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	return p.ProblemSummaryFn(ctx, patientID)
}

// VisitSummary is the VisitSummary mock
func (p *ClinicalMock) VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
	return p.VisitSummaryFn(ctx, encounterID, count)
}

// PatientTimelineWithCount is the PatientTimelineWithCount mock
func (p *ClinicalMock) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	return p.PatientTimelineWithCountFn(ctx, episodeID, count)
}

// ContactsToContactPointInput is the ContactsToContactPointInput mock
func (p *ClinicalMock) ContactsToContactPointInput(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
	return p.ContactsToContactPointInputFn(ctx, phones, emails)
}

// RegisterPatient is the RegisterPatient mock
func (p *ClinicalMock) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return p.RegisterPatientFn(ctx, input)
}

// CreatePatient is the CreatePatient mock
func (p *ClinicalMock) CreatePatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	return p.CreatePatientFn(ctx, input)
}

// FindPatientByID is the FindPatientByID mock
func (p *ClinicalMock) FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error) {
	return p.FindPatientByIDFn(ctx, id)
}

// PatientSearch initializes PatientSearch mock
func (p *ClinicalMock) PatientSearch(ctx context.Context, search string) (*domain.PatientConnection, error) {
	return p.PatientSearchFn(ctx, search)
}

// UpdatePatient initializes UpdatePatient mock
func (p *ClinicalMock) UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return p.UpdatePatientFn(ctx, input)
}

// AddNextOfKin initializes AddNextOfKin mock
func (p *ClinicalMock) AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
	return p.AddNextOfKinFn(ctx, input)
}

// AddNHIF initializes AddNHIF mock
func (p *ClinicalMock) AddNHIF(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
	return p.AddNHIFFn(ctx, input)
}

// RegisterUser initializes RegisterUser mocks
func (p *ClinicalMock) RegisterUser(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return p.RegisterUserFn(ctx, input)
}

// CreateUpdatePatientExtraInformation initializes CreateUpdatePatientExtraInformation mock
func (p *ClinicalMock) CreateUpdatePatientExtraInformation(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
	return p.CreateUpdatePatientExtraInformationFn(ctx, input)
}

// AllergySummary is the AllergySummary mock initializer
func (p *ClinicalMock) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	return p.AllergySummaryFn(ctx, patientID)
}

// DeleteFHIRPatientByPhone initializes DeleteFHIRPatientByPhone mock
func (p *ClinicalMock) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	return p.DeleteFHIRPatientByPhoneFn(ctx, phoneNumber)
}

// StartEpisodeByBreakGlass initializes StartEpisodeByBreakGlass mock
func (p *ClinicalMock) StartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	return p.StartEpisodeByBreakGlassFn(ctx, input)
}
