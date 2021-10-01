package mock

import "context"

// FakeFHIRRepository is a mock FHIR repository
type FakeFHIRRepository struct {
	CreateFHIRResourceFn       func(resourceType string, payload map[string]interface{}) ([]byte, error)
	DeleteFHIRResourceFn       func(resourceType, fhirResourceID string) ([]byte, error)
	PatchFHIRResourceFn        func(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error)
	UpdateFHIRResourceFn       func(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error)
	GetFHIRPatientAllDataFn    func(fhirResourceID string) ([]byte, error)
	FHIRRestURLFn              func() string
	GetFHIRResourceFn          func(resourceType, fhirResourceID string) ([]byte, error)
	GetFHIRPatientEverythingFn func(fhirResourceID string) ([]byte, error)
}

// CreateFHIRResource ...
func (f *FakeFHIRRepository) CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error) {
	return f.CreateFHIRResourceFn(resourceType, payload)
}

// DeleteFHIRResource ...
func (f *FakeFHIRRepository) DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return f.DeleteFHIRResourceFn(resourceType, fhirResourceID)
}

// PatchFHIRResource ...
func (f *FakeFHIRRepository) PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error) {
	return f.PatchFHIRResourceFn(resourceType, fhirResourceID, payload)
}

// UpdateFHIRResource ...
func (f *FakeFHIRRepository) UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
	return f.UpdateFHIRResourceFn(resourceType, fhirResourceID, payload)
}

// GetFHIRPatientAllData ...
func (f *FakeFHIRRepository) GetFHIRPatientAllData(fhirResourceID string) ([]byte, error) {
	return f.GetFHIRPatientAllDataFn(fhirResourceID)
}

// FHIRRestURL ...
func (f *FakeFHIRRepository) FHIRRestURL() string {
	return f.FHIRRestURLFn()
}

// GetFHIRResource ...
func (f *FakeFHIRRepository) GetFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return f.GetFHIRResourceFn(resourceType, fhirResourceID)
}

//GetFHIRPatientEverything ....
func (f *FakeFHIRRepository) GetFHIRPatientEverything(fhirResourceID string) ([]byte, error) {
	return f.GetFHIRPatientEverythingFn(fhirResourceID)
}

// FakeRepository is a mock firebase repository
type FakeRepository struct {
	SaveEmailOTPFn func(ctx context.Context, email string, optIn bool) error
}

// SaveEmailOTP ...
func (fb *FakeRepository) SaveEmailOTP(ctx context.Context, email string, optIn bool) error {
	return fb.SaveEmailOTPFn(ctx, email, optIn)
}