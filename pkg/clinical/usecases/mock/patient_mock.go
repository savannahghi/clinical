package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// ClinicalMock contains all mock methods
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

// ProblemSummary is a mock implementation of ProblemSummary method
func (p *ClinicalMock) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	return p.ProblemSummaryFn(ctx, patientID)
}

// VisitSummary is a mock implementation of VisitSummary method
func (p *ClinicalMock) VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
	return p.VisitSummaryFn(ctx, encounterID, count)
}

// PatientTimelineWithCount is a mock implementation of PatientTimelineWithCount method
func (p *ClinicalMock) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	return p.PatientTimelineWithCountFn(ctx, episodeID, count)
}

// ContactsToContactPointInput is a mock implementation of ContactsToContactPointInput method
func (p *ClinicalMock) ContactsToContactPointInput(ctx context.Context, phones []*domain.PhoneNumberInput, emails []*domain.EmailInput) ([]*domain.FHIRContactPointInput, error) {
	return p.ContactsToContactPointInputFn(ctx, phones, emails)
}

// RegisterPatient is a mock implementation of RegisterPatient method
func (p *ClinicalMock) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return p.RegisterPatientFn(ctx, input)
}

// CreatePatient is a mock implementation of CreatePatient method
func (p *ClinicalMock) CreatePatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
	return p.CreatePatientFn(ctx, input)
}

// FindPatientByID is a mock implementation of FindPatientByID method
func (p *ClinicalMock) FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error) {
	return p.FindPatientByIDFn(ctx, id)
}

// PatientSearch is a mock implementation of PatientSearch method
func (p *ClinicalMock) PatientSearch(ctx context.Context, search string) (*domain.PatientConnection, error) {
	return p.PatientSearchFn(ctx, search)
}

// UpdatePatient is a mock implementation of UpdatePatient method
func (p *ClinicalMock) UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return p.UpdatePatientFn(ctx, input)
}

// AddNextOfKin is a mock implementation of AddNextOfKin method
func (p *ClinicalMock) AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
	return p.AddNextOfKinFn(ctx, input)
}

// AddNHIF is a mock implementation of AddNHIF method
func (p *ClinicalMock) AddNHIF(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
	return p.AddNHIFFn(ctx, input)
}

// RegisterUser is a mock implementation of RegisterUser method
func (p *ClinicalMock) RegisterUser(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	return p.RegisterUserFn(ctx, input)
}

// CreateUpdatePatientExtraInformation is a mock implementation of CreateUpdatePatientExtraInformation method
func (p *ClinicalMock) CreateUpdatePatientExtraInformation(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
	return p.CreateUpdatePatientExtraInformationFn(ctx, input)
}

// AllergySummary is a mock implementation of AllergySummary method
func (p *ClinicalMock) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	return p.AllergySummaryFn(ctx, patientID)
}

// DeleteFHIRPatientByPhone is a mock implementation of DeleteFHIRPatientByPhone method
func (p *ClinicalMock) DeleteFHIRPatientByPhone(ctx context.Context, phoneNumber string) (bool, error) {
	return p.DeleteFHIRPatientByPhoneFn(ctx, phoneNumber)
}

// StartEpisodeByBreakGlass is a mock implementation of StartEpisodeByBreakGlass method
func (p *ClinicalMock) StartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	return p.StartEpisodeByBreakGlassFn(ctx, input)
}
