package mock

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// FakeFHIRRepository is a mock FHIR repository
type FakeFHIRRepository struct {
	MockCreateFHIRResourceFn    func(resourceType string, payload map[string]interface{}, resource interface{}) error
	MockDeleteFHIRResourceFn    func(resourceType, fhirResourceID string) error
	MockPatchFHIRResourceFn     func(resourceType, fhirResourceID string, payload []map[string]interface{}, resource interface{}) error
	MockUpdateFHIRResourceFn    func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error
	MockGetFHIRPatientAllDataFn func(fhirResourceID string) ([]byte, error)
	MockGetFHIRResourceFn       func(resourceType, fhirResourceID string, resource interface{}) error
	MockSearchFHIRResourceFn    func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers) ([]map[string]interface{}, error)
}

// NewFakeFHIRRepositoryMock initializes a new FakeFHIRRepositoryMock
func NewFakeFHIRRepositoryMock() *FakeFHIRRepository {
	return &FakeFHIRRepository{
		MockCreateFHIRResourceFn: func(resourceType string, payload map[string]interface{}, resource interface{}) error {
			return nil
		},
		MockDeleteFHIRResourceFn: func(resourceType, fhirResourceID string) error {
			return nil
		},
		MockPatchFHIRResourceFn: func(resourceType, fhirResourceID string, payload []map[string]interface{}, resource interface{}) error {
			return nil
		},
		MockUpdateFHIRResourceFn: func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
			return nil
		},
		MockGetFHIRPatientAllDataFn: func(fhirResourceID string) ([]byte, error) {
			bs, err := json.Marshal("m")
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
		},
		MockGetFHIRResourceFn: func(resourceType, fhirResourceID string, resource interface{}) error {
			return nil
		},
		MockSearchFHIRResourceFn: func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers) ([]map[string]interface{}, error) {
			n := map[string]interface{}{"given": []string{"John"}, "family": []string{"Doe"}}
			p := map[string]interface{}{
				"resourceType": "Patient/",
				"id":           "test-UUID",
				"name":         []map[string]interface{}{n},
				"reference":    "Patient/",
			}

			m := []map[string]interface{}{
				{
					"resourceType": "Patient/",
					"status":       "active",
					"id":           "test-UUID",
					"patient":      p,
					"period": map[string]interface{}{
						"start": time.February.String(),
						"end":   time.February.String(),
					},
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
func (f *FakeFHIRRepository) DeleteFHIRResource(resourceType, fhirResourceID string) error {
	return f.MockDeleteFHIRResourceFn(resourceType, fhirResourceID)
}

// PatchFHIRResource ...
func (f *FakeFHIRRepository) PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}, resource interface{}) error {
	return f.MockPatchFHIRResourceFn(resourceType, fhirResourceID, payload, resource)
}

// UpdateFHIRResource ...
func (f *FakeFHIRRepository) UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
	return f.MockUpdateFHIRResourceFn(resourceType, fhirResourceID, payload, resource)
}

// GetFHIRPatientAllData ...
func (f *FakeFHIRRepository) GetFHIRPatientAllData(fhirResourceID string) ([]byte, error) {
	return f.MockGetFHIRPatientAllDataFn(fhirResourceID)
}

// GetFHIRResource ...
func (f *FakeFHIRRepository) GetFHIRResource(resourceType, fhirResourceID string, resource interface{}) error {
	return f.MockGetFHIRResourceFn(resourceType, fhirResourceID, resource)
}

// SearchFHIRResource ...
func (f *FakeFHIRRepository) SearchFHIRResource(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers) ([]map[string]interface{}, error) {
	return f.MockSearchFHIRResourceFn(resourceType, params, tenant)
}
