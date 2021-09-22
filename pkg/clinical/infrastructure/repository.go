package infrastructure

import "google.golang.org/api/healthcare/v1"

type FHIRRepository interface {
	CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error)
	DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error)
	PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error)
	UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error)
	CreateFHIRStore() (*healthcare.FhirStore, error)
	GetFHIRStore() (*healthcare.FhirStore, error)
	GetFHIRPatientEverything(fhirResourceID string) ([]byte, error)
}
