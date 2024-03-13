package mock

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/domain"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// FakeFHIRRepository is a mock FHIR repository
type FakeFHIRRepository struct {
	MockCreateFHIRResourceFn    func(resourceType string, payload map[string]interface{}, resource interface{}) error
	MockDeleteFHIRResourceFn    func(resourceType, fhirResourceID string) error
	MockPatchFHIRResourceFn     func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error
	MockUpdateFHIRResourceFn    func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error
	MockGetFHIRPatientAllDataFn func(fhirResourceID string, params map[string]interface{}) ([]byte, error)
	MockGetFHIRResourceFn       func(resourceType, fhirResourceID string, resource interface{}) error
	MockSearchFHIRResourceFn    func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error)
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
		MockPatchFHIRResourceFn: func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
			return nil
		},
		MockUpdateFHIRResourceFn: func(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
			return nil
		},
		MockGetFHIRPatientAllDataFn: func(fhirResourceID string, params map[string]interface{}) ([]byte, error) {
			bs, err := json.Marshal(`
			"getPatientEverything": {
				"entry": [
				  {
					"fullUrl": "https://healthcare.googleapis.com/v1/projects/",
					"resource": {
					  "active": true,
					  "birthDate": "2024-03-12",
					  "gender": "male",
					  "id": "2051cdba-6f88-405b-aa51-56e28b75b941",
					  "identifier": [
						{
						  "period": {
							"end": "2124-05-27T15:16:04+03:00",
							"start": "2024-03-12T15:16:04+03:00"
						  },
						  "system": "healthcloud.msisdn",
						  "type": {
							"coding": [
							  {
								"code": "+2547011222222",
								"display": "+2547011222222",
								"system": "healthcloud.iddocument",
								"userSelected": true,
								"version": "0.0.1"
							  }
							],
							"text": "+2547011222222"
						  },
						  "use": "official",
						  "value": "+2547011222222"
						}
					  ],
					  "language": "EN",
					  "managingOrganization": {
						"display": "Chuka",
						"id": "8f5c7e78-5d3e-401f-9148-95b4634bfbde",
						"reference": "Organization/8f5c7e78-5d3e-401f-9148-95b4698bfgde",
						"type": "Organization"
					  },
					  "maritalStatus": {
						"coding": [
						  {
							"display": "unknown",
							"userSelected": true
						  }
						],
						"text": "unknown"
					  },
					  "meta": {
						"lastUpdated": "2024-03-12T15:16:05.080533+00:00",
						"tag": [
						  {
							"code": "85e4b0d3-1d69-47ba-b265-579d125f18e5",
							"display": "Napoleon Health Services",
							"system": "http://mycarehub/tenant-identification/organisation",
							"userSelected": false,
							"version": "1.0"
						  },
						  {
							"code": "8f5c7e78-5d3e-401f-9148-95b4698bfbde",
							"display": "Nairobi",
							"system": "http://mycarehub/tenant-identification/facility",
							"userSelected": false,
							"version": "1.0"
						  }
						],
						"versionId": "MTcxMDI1NjU2NTA4MDUzMzAwMA"
					  },
					  "name": [
						{
						  "family": "Jane",
						  "given": [
							"Brian"
						  ],
						  "period": {
							"end": "2124-05-27T15:16:04+03:00",
							"start": "2024-03-12T15:16:04+03:00"
						  },
						  "text": "Jane, Green ",
						  "use": "official"
						}
					  ],
					  "resourceType": "Patient",
					  "telecom": [
						{
						  "period": {
							"end": "2124-05-27T15:16:04+03:00",
							"start": "2024-03-12T15:16:04+03:00"
						  },
						  "rank": 2,
						  "system": "phone",
						  "use": "home",
						  "value": "+2547011222222"
						}
					  ]
					}
				  }
				],
				"link": [
				  {
					"relation": "next",
					"url": "https://healthcare.googleapis.com/v1/projects/"
				  }
				],
				"resourceType": "Bundle",
				"total": 4,
				"type": "searchset"
			  }
			}`)
			if err != nil {
				return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
			}
			return bs, nil
		},
		MockGetFHIRResourceFn: func(resourceType, fhirResourceID string, resource interface{}) error {
			return nil
		},
		MockSearchFHIRResourceFn: func(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
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

			return &domain.PagedFHIRResource{
				Resources: m,
			}, nil
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
func (f *FakeFHIRRepository) PatchFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
	return f.MockPatchFHIRResourceFn(resourceType, fhirResourceID, payload, resource)
}

// UpdateFHIRResource ...
func (f *FakeFHIRRepository) UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
	return f.MockUpdateFHIRResourceFn(resourceType, fhirResourceID, payload, resource)
}

// GetFHIRPatientAllData ...
func (f *FakeFHIRRepository) GetFHIRPatientAllData(fhirResourceID string, params map[string]interface{}) ([]byte, error) {
	return f.MockGetFHIRPatientAllDataFn(fhirResourceID, params)
}

// GetFHIRResource ...
func (f *FakeFHIRRepository) GetFHIRResource(resourceType, fhirResourceID string, resource interface{}) error {
	return f.MockGetFHIRResourceFn(resourceType, fhirResourceID, resource)
}

// SearchFHIRResource ...
func (f *FakeFHIRRepository) SearchFHIRResource(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
	return f.MockSearchFHIRResourceFn(resourceType, params, tenant, pagination)
}
