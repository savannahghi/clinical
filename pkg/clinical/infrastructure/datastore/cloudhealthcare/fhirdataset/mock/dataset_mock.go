package mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"google.golang.org/api/healthcare/v1"
)

// FakeFHIRRepository is a mock FHIR repository
type FakeFHIRRepository struct {
	MockCreateFHIRResourceFn       func(resourceType string, payload map[string]interface{}) ([]byte, error)
	MockDeleteFHIRResourceFn       func(resourceType, fhirResourceID string) ([]byte, error)
	MockPatchFHIRResourceFn        func(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error)
	MockUpdateFHIRResourceFn       func(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error)
	MockGetFHIRPatientAllDataFn    func(fhirResourceID string) ([]byte, error)
	MockFHIRRestURLFn              func() string
	MockFHIRHeadersFn              func() (http.Header, error)
	MockGetFHIRResourceFn          func(resourceType, fhirResourceID string) ([]byte, error)
	MockGetFHIRPatientEverythingFn func(fhirResourceID string) ([]byte, error)
	MockPOSTRequestFn              func(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error)
	MockCreateDatasetFn            func() (*healthcare.Operation, error)
	MockGetDatasetFn               func() (*healthcare.Dataset, error)
	MockGetFHIRStoreFn             func() (*healthcare.FhirStore, error)
	MockCreateFHIRStoreFn          func() (*healthcare.FhirStore, error)
}

// NewFakeFHIRRepositoryMock initializes a new FakeFHIRRepositoryMock
func NewFakeFHIRRepositoryMock() *FakeFHIRRepository {
	return &FakeFHIRRepository{
		MockCreateFHIRResourceFn: func(resourceType string, payload map[string]interface{}) ([]byte, error) {
			m := map[string]interface{}{
				"key": "value",
			}
			bs, err := json.Marshal(m)
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
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
		MockFHIRRestURLFn: func() string {
			return "https://healthcare.googleapis.com/v1"
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
		MockGetFHIRPatientEverythingFn: func(fhirResourceID string) ([]byte, error) {
			bs, err := json.Marshal("m")
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
		},
		MockPOSTRequestFn: func(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
			m := map[string]string{
				"resourceType": "Bundle",
				"type":         "searchset",
				"total":        "10",
				"link":         "test",
			}
			bs, err := json.Marshal(m)
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
		},
		MockFHIRHeadersFn: func() (http.Header, error) {
			return http.Header{
				"Authorization": []string{"Bearer " + uuid.NewString()},
			}, nil
		},
		MockCreateDatasetFn: func() (*healthcare.Operation, error) {
			return &healthcare.Operation{
				Done:     false,
				Error:    &healthcare.Status{},
				Metadata: []byte{},
				Name:     "",
				Response: []byte{},
			}, nil
		},
		MockGetDatasetFn: func() (*healthcare.Dataset, error) {
			return &healthcare.Dataset{
				Name: "test",
			}, nil
		},
		MockGetFHIRStoreFn: func() (*healthcare.FhirStore, error) {
			return &healthcare.FhirStore{
				DefaultSearchHandlingStrict: true,
				DisableReferentialIntegrity: true,
			}, nil
		},
		MockCreateFHIRStoreFn: func() (*healthcare.FhirStore, error) {
			return &healthcare.FhirStore{
				DefaultSearchHandlingStrict: true,
				DisableReferentialIntegrity: true,
			}, nil
		},
	}
}

// CreateFHIRResource ...
func (f *FakeFHIRRepository) CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error) {
	return f.MockCreateFHIRResourceFn(resourceType, payload)
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

// FHIRRestURL ...
func (f *FakeFHIRRepository) FHIRRestURL() string {
	return f.MockFHIRRestURLFn()
}

// GetFHIRResource ...
func (f *FakeFHIRRepository) GetFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return f.MockGetFHIRResourceFn(resourceType, fhirResourceID)
}

//GetFHIRPatientEverything ....
func (f *FakeFHIRRepository) GetFHIRPatientEverything(fhirResourceID string) ([]byte, error) {
	return f.MockGetFHIRPatientEverythingFn(fhirResourceID)
}

// POSTRequest is a mock implementation of POSTRequest method
func (f *FakeFHIRRepository) POSTRequest(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
	return f.MockPOSTRequestFn(resourceName, path, params, body)
}

// FHIRHeaders is a mock implementation of CreateFHIRMedication method
func (f *FakeFHIRRepository) FHIRHeaders() (http.Header, error) {
	return f.MockFHIRHeadersFn()
}

// CreateDataset is a mock implementation of CreateDataset method
func (f *FakeFHIRRepository) CreateDataset() (*healthcare.Operation, error) {
	return f.MockCreateDatasetFn()
}

// GetDataset is a mock implementation of GetDataset method
func (f *FakeFHIRRepository) GetDataset() (*healthcare.Dataset, error) {
	return f.MockGetDatasetFn()
}

// GetFHIRStore is a mock implementation of GetFHIRStore method
func (f *FakeFHIRRepository) GetFHIRStore() (*healthcare.FhirStore, error) {
	return f.MockGetFHIRStoreFn()
}

// CreateFHIRStore is a mock implementation of CreateFHIRStore method
func (f *FakeFHIRRepository) CreateFHIRStore() (*healthcare.FhirStore, error) {
	return f.MockCreateFHIRStoreFn()
}
