package infrastructure

import (
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhir"
)

// FHIRRepository ...
type FHIRRepository interface {
	CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error)
	DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error)
	PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error)
	UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error)
	GetFHIRPatientEverything(fhirResourceID string) ([]byte, error)
}

// FHIRService is an implementation of the database repository
// It is implementation agnostic i.e logic should be handled using
// the preferred database
type FHIRService struct {
	FHIR FHIRRepository
}

// NewFHIRService creates a new database service
func NewFHIRService() FHIRService {
	datasetExtension := fhir.NewDatasetExtension()
	repo := fhir.NewFHIRRepository(datasetExtension)

	return FHIRService{
		repo,
	}
}

// CreateFHIRResource creates an FHIR resource.
//
// The payload should be the result of marshalling a resource to JSON
func (d FHIRService) CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error) {
	return d.FHIR.CreateFHIRResource(resourceType, payload)
}

// DeleteFHIRResource deletes an FHIR resource.
func (d FHIRService) DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return d.FHIR.DeleteFHIRResource(resourceType, fhirResourceID)
}

// PatchFHIRResource patches a FHIR resource.
// The payload is a JSON patch document that follows guidance on Patch from the
// FHIR standard.
// See:
// payload := []map[string]interface{}{
// 	{
// 		"op":    "replace",
// 		"path":  "/active",
// 		"value": active,
// 	},
// }
// See: https://www.hl7.org/fhir/http.html#patch
func (d FHIRService) PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error) {
	return d.FHIR.PatchFHIRResource(resourceType, fhirResourceID, payload)
}

// UpdateFHIRResource updates the entire contents of a resource.
func (d FHIRService) UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
	return d.FHIR.UpdateFHIRResource(resourceType, fhirResourceID, payload)
}

// GetFHIRPatientEverything gets all resources associated with a particular
// patient compartment.
func (d FHIRService) GetFHIRPatientEverything(fhirResourceID string) ([]byte, error) {
	return d.FHIR.GetFHIRPatientEverything(fhir.DatasetLocation)
}
