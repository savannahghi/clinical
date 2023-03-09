package mock

import (
	"encoding/json"
	"fmt"
)

// FakeFHIRRepository is a mock FHIR repository
type FakeFHIRRepository struct {
	MockCreateFHIRResourceFn    func(resourceType string, payload map[string]interface{}, resource interface{}) error
	MockDeleteFHIRResourceFn    func(resourceType, fhirResourceID string) ([]byte, error)
	MockPatchFHIRResourceFn     func(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error)
	MockUpdateFHIRResourceFn    func(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error)
	MockGetFHIRPatientAllDataFn func(fhirResourceID string) ([]byte, error)
	MockGetFHIRResourceFn       func(resourceType, fhirResourceID string) ([]byte, error)
	MockSearchFHIRResourceFn    func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error)
}

// NewFakeFHIRRepositoryMock initializes a new FakeFHIRRepositoryMock
func NewFakeFHIRRepositoryMock() *FakeFHIRRepository {
	return &FakeFHIRRepository{
		MockCreateFHIRResourceFn: func(resourceType string, payload map[string]interface{}, resource interface{}) error {
			return nil
		},
		MockDeleteFHIRResourceFn: func(resourceType, fhirResourceID string) ([]byte, error) {
			bs, err := json.Marshal("m")
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
		},
		MockPatchFHIRResourceFn: func(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error) {
			m := map[string]interface{}{
				"key": "value",
			}
			bs, err := json.Marshal(m)
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
		},
		MockUpdateFHIRResourceFn: func(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
			m := map[string]interface{}{
				"key": "value",
			}
			bs, err := json.Marshal(m)
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
		},
		MockGetFHIRPatientAllDataFn: func(fhirResourceID string) ([]byte, error) {
			bs, err := json.Marshal("m")
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
		},
		MockGetFHIRResourceFn: func(resourceType, fhirResourceID string) ([]byte, error) {
			n := map[string]interface{}{"given": []string{"John"}, "family": []string{"Doe"}}
			p := map[string]interface{}{
				"resourceType": "Patient/",
				"id":           "test-UUID",
				"name":         []map[string]interface{}{n},
			}
			m := map[string]interface{}{
				"resourceType":  "Patient/",
				"status":        "active",
				"id":            "test-UUID",
				"patientRecord": p,
			}
			bs, err := json.Marshal(m)
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
		},
		MockSearchFHIRResourceFn: func(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
			n := map[string]interface{}{"given": []string{"John"}, "family": []string{"Doe"}}
			p := map[string]interface{}{
				"resourceType": "Patient/",
				"id":           "test-UUID",
				"name":         []map[string]interface{}{n},
			}

			m := []map[string]interface{}{
				{
					"resourceType":  "Patient/",
					"status":        "active",
					"id":            "test-UUID",
					"patientRecord": p,
				},
			}

			return m, nil
		},
	}
}

// CreateFHIRResource ...
func (f *FakeFHIRRepository) CreateFHIRResource(resourceType string, payload map[string]interface{}, resource interface{}) error {
	return f.MockCreateFHIRResourceFn(resourceType, payload, resource)
}

// DeleteFHIRResource ...
func (f *FakeFHIRRepository) DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return f.MockDeleteFHIRResourceFn(resourceType, fhirResourceID)
}

// PatchFHIRResource ...
func (f *FakeFHIRRepository) PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error) {
	return f.MockPatchFHIRResourceFn(resourceType, fhirResourceID, payload)
}

// UpdateFHIRResource ...
func (f *FakeFHIRRepository) UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
	return f.MockUpdateFHIRResourceFn(resourceType, fhirResourceID, payload)
}

// GetFHIRPatientAllData ...
func (f *FakeFHIRRepository) GetFHIRPatientAllData(fhirResourceID string) ([]byte, error) {
	return f.MockGetFHIRPatientAllDataFn(fhirResourceID)
}

// GetFHIRResource ...
func (f *FakeFHIRRepository) GetFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return f.MockGetFHIRResourceFn(resourceType, fhirResourceID)
}

// SearchFHIRResource ...
func (f *FakeFHIRRepository) SearchFHIRResource(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
	return f.MockSearchFHIRResourceFn(resourceType, params)
}
