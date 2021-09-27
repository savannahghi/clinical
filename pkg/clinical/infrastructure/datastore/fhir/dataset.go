package fhir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/savannahghi/serverutils"
	"google.golang.org/api/healthcare/v1"
)

// Repository accesses and updates patient data that is stored on Healthcare
// FHIR repository
type Repository struct {
	healthcareService                           *healthcare.Service
	projectID, location, datasetID, fhirStoreID string
}

// NewFHIRRepository initializes a FHIR repository
func NewFHIRRepository() *Repository {
	project := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	_ = serverutils.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
	dataset := serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
	fhirStore := serverutils.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")
	ctx := context.Background()
	hsv, err := healthcare.NewService(ctx)
	if err != nil {
		log.Panicf("unable to initialize new Google Cloud Healthcare Service: %s", err)
	}
	return &Repository{
		healthcareService: hsv,
		projectID:         project,
		location:          DatasetLocation,
		datasetID:         dataset,
		fhirStoreID:       fhirStore,
	}
}

// CreateDataset creates a dataset and returns it's name
func (fr Repository) CreateDataset() (*healthcare.Operation, error) {
	fr.checkPreconditions()
	datasetsService := fr.healthcareService.Projects.Locations.Datasets
	parent := fmt.Sprintf("projects/%s/locations/%s", fr.projectID, fr.location)
	resp, err := datasetsService.Create(parent, &healthcare.Dataset{}).DatasetId(fr.datasetID).Do()
	if err != nil {
		return nil, fmt.Errorf("create Data Set: %v", err)
	}
	return resp, nil
}

// GetDataset gets a dataset.
func (fr Repository) GetDataset() (*healthcare.Dataset, error) {
	fr.checkPreconditions()
	datasetsService := fr.healthcareService.Projects.Locations.Datasets
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", fr.projectID, fr.location, fr.datasetID)
	resp, err := datasetsService.Get(name).Do()
	if err != nil {
		return nil, fmt.Errorf("get Data Set: %v", err)
	}
	return resp, nil
}

// CreateFHIRStore creates an FHIR store.
func (fr Repository) CreateFHIRStore() (*healthcare.FhirStore, error) {
	fr.checkPreconditions()
	storesService := fr.healthcareService.Projects.Locations.Datasets.FhirStores
	store := &healthcare.FhirStore{
		DisableReferentialIntegrity: false,
		DisableResourceVersioning:   false,
		EnableUpdateCreate:          true,
		Version:                     "R4",
	}
	parent := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", fr.projectID, fr.location, fr.datasetID)
	resp, err := storesService.Create(parent, store).FhirStoreId(fr.fhirStoreID).Do()
	if err != nil {
		return nil, fmt.Errorf("create FHIR Store: %v", err)
	}
	return resp, nil
}

// GetFHIRStore gets an FHIR store.
func (fr Repository) GetFHIRStore() (*healthcare.FhirStore, error) {
	fr.checkPreconditions()
	storesService := fr.healthcareService.Projects.Locations.Datasets.FhirStores
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s",
		fr.projectID, fr.location, fr.datasetID, fr.fhirStoreID)
	store, err := storesService.Get(name).Do()
	if err != nil {
		return nil, fmt.Errorf("get FHIR Store: %v", err)
	}
	return store, nil
}

func (fr Repository) checkPreconditions() {
	if fr.healthcareService == nil {
		log.Panicf("cloudhealth.Service *healthcare.Service is nil")
	}
}

// FHIRRestURL composes a FHIR REST URL for manual calls to the FHIR REST API
func (fr Repository) FHIRRestURL() string {
	return fmt.Sprintf(
		"%s/projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir",
		baseFHIRURL, fr.projectID, fr.location, fr.datasetID, fr.fhirStoreID)
}

// CreateFHIRResource creates an FHIR resource.
//
// The payload should be the result of marshalling a resource to JSON
func (fr Repository) CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error) {
	fr.checkPreconditions()
	payload["resourceType"] = resourceType
	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	payload["resourceType"] = resourceType
	payload["language"] = "EN"
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json.Encode: %v", err)
	}

	parent := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s",
		fr.projectID, fr.location, fr.datasetID, fr.fhirStoreID)
	call := fhirService.Create(
		parent, resourceType, bytes.NewReader(jsonPayload))

	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("create: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(
			"create: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
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

// GetFHIRResource gets an FHIR resource.
func (fr Repository) GetFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	fr.checkPreconditions()
	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s",
		fr.projectID, fr.location, fr.datasetID, fr.fhirStoreID,
		resourceType, fhirResourceID)
	call := fhirService.Read(name)
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("read: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("read: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}
