package fhir

// Repository accesses and updates patient data that is stored on Healthcare
// FHIR repository
type Repository struct {
	Dataset DatasetExtension
}

// NewFHIRRepository initializes a FHIR repository
func NewFHIRRepository(
	dataset DatasetExtension,
) *Repository {
	return &Repository{
		Dataset: dataset,
	}
}

// CreateFHIRResource creates an FHIR resource.
//
// The payload should be the result of marshalling a resource to JSON
func (fr Repository) CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error) {
	return nil, nil
}

// DeleteFHIRResource deletes an FHIR resource.
func (fr Repository) DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	return nil, nil
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
func (fr Repository) PatchFHIRResource(
	resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error) {
	return nil, nil
}

// UpdateFHIRResource updates the entire contents of a resource.
func (fr Repository) UpdateFHIRResource(
	resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
	return nil, nil
}

// GetFHIRPatientEverything gets all resources associated with a particular
// patient compartment.
func (fr Repository) GetFHIRPatientEverything(fhirResourceID string) ([]byte, error) {
	return nil, nil
}
