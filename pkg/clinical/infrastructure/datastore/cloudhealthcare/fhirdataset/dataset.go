package fhirdataset

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/serverutils"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/healthcare/v1"
)

// constants used to configure the Google Cloud Healthcare API
const (
	DatasetLocation       = "europe-west4"
	baseFHIRURL           = "https://healthcare.googleapis.com/v1"
	defaultTimeoutSeconds = 10
)

// FHIRRepository ...
type FHIRRepository interface {
	CreateFHIRResource(resourceType string, payload map[string]interface{}) ([]byte, error)
	DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error)
	PatchFHIRResource(resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error)
	UpdateFHIRResource(resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error)
	GetFHIRPatientAllData(fhirResourceID string) ([]byte, error)
	FHIRRestURL() string
	GetFHIRResource(resourceType, fhirResourceID string) ([]byte, error)
	GetFHIRPatientEverything(fhirResourceID string) ([]byte, error)
	POSTRequest(resourceName string, path string, params url.Values, body io.Reader) ([]byte, error)
	FHIRHeaders() (http.Header, error)
	CreateDataset() (*healthcare.Operation, error)
	GetDataset() (*healthcare.Dataset, error)
	GetFHIRStore() (*healthcare.FhirStore, error)
	CreateFHIRStore() (*healthcare.FhirStore, error)
}

// Repository accesses and updates patient data that is stored on Healthcare
// FHIR repository
type Repository struct {
	healthcareService                           *healthcare.Service
	projectID, location, datasetID, fhirStoreID string
}

// NewFHIRRepository initializes a FHIR repository
func NewFHIRRepository() FHIRRepository {
	project := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	_ = serverutils.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
	dataset := serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
	fhirStore := serverutils.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")
	ctx := context.Background()
	hsv, err := healthcare.NewService(ctx)
	if err != nil {
		utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
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
		utils.ReportErrorToSentry(err)
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

// getErrorMessage unmarshalls the error response that is returned from the FHIR server.
// This function should be called when a status code > 299 has been returned
func getErrorMessage(respBytes []byte) (errorText string, diagnostics string, err error) {
	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(respBytes, &errorResponse)
	if err != nil {
		return "", "", fmt.Errorf("could not unmarshal error response: %v", err)
	}

	errorText = errorResponse.Issue[0].Details.Text
	diagnostics = errorResponse.Issue[0].Diagnostics

	return errorText, diagnostics, nil
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
		errorText, diagnostics, err := getErrorMessage(respBytes)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s: %s", resp.Status, errorText, diagnostics)
	}
	return respBytes, nil
}

// DeleteFHIRResource deletes an FHIR resource.
func (fr Repository) DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	fr.checkPreconditions()

	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s",
		fr.projectID, fr.location, fr.datasetID, fr.fhirStoreID,
		resourceType, fhirResourceID)
	resp, err := fhirService.Delete(name).Do()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("delete: %v", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(
			"delete: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// PatchFHIRResource patches a FHIR resource.
// The payload is a JSON patch document that follows guidance on Patch from the
// FHIR standard.
// See:
//
//	payload := []map[string]interface{}{
//		{
//			"op":    "replace",
//			"path":  "/active",
//			"value": active,
//		},
//	}
//
// See: https://www.hl7.org/fhir/http.html#patch
func (fr Repository) PatchFHIRResource(
	resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error) {
	fr.checkPreconditions()
	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	if serverutils.IsDebug() {
		log.Printf("FHIR Payload: %#v", payload)
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("json.Encode: %v", err)
	}
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s",
		fr.projectID, fr.location, fr.datasetID, fr.fhirStoreID,
		resourceType, fhirResourceID)

	call := fhirService.Patch(name, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/json-patch+json")
	resp, err := call.Do()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("patch: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if serverutils.IsDebug() {
		log.Printf("Patch FHIR Resource %d Response: %s", resp.StatusCode, string(respBytes))
	}
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(
			"patch: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// UpdateFHIRResource updates the entire contents of a resource.
func (fr Repository) UpdateFHIRResource(
	resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
	fr.checkPreconditions()
	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	payload["resourceType"] = resourceType
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("json.Encode: %v", err)
	}
	if serverutils.IsDebug() {
		log.Printf("FHIR Update payload: %s", string(jsonPayload))
	}
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s",
		fr.projectID, fr.location, fr.datasetID, fr.fhirStoreID,
		resourceType, fhirResourceID)
	call := fhirService.Update(name, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	resp, err := call.Do()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("update: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(
			"update: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// GetFHIRPatientAllData gets all resources associated with a particular
// patient compartment.
func (fr Repository) GetFHIRPatientAllData(fhirResourceID string) ([]byte, error) {
	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/Patient/%s",
		fr.projectID,
		fr.location,
		fr.datasetID,
		fr.fhirStoreID,
		fhirResourceID,
	)

	resp, err := fhirService.PatientEverything(name).Do()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("PatientAllData: %v", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("could not read response: %v", err)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("PatientAllData: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}

	return respBytes, nil
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
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("read: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("read: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// GetFHIRPatientEverything gets all resources associated with a particular
// patient compartment.
func (fr Repository) GetFHIRPatientEverything(fhirResourceID string) ([]byte, error) {
	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/Patient/%s",
		fr.projectID,
		fr.location,
		fr.datasetID,
		fr.fhirStoreID,
		fhirResourceID,
	)

	resp, err := fhirService.PatientEverything(name).Do()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("PatientEverything: %v", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("could not read response: %v", err)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("PatientEverything: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}

	return respBytes, nil
}

// POSTRequest is used to manually compose POST requests to the FHIR service
//
// - `resourceName` is a FHIR resource name e.g "Patient"
// - `path` is a sub-path e.g `_search` under a resource
// - `params` should be query params, sent as `url.Values`
func (fr Repository) POSTRequest(
	resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
	fhirHeaders, err := fr.FHIRHeaders()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to get FHIR headers: %v", err)
	}
	url := fmt.Sprintf(
		"%s/%s/%s?%s", fr.FHIRRestURL(), resourceName, path, params.Encode())
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to compose FHIR POST request: %v", err)
	}
	for k, v := range fhirHeaders {
		for _, h := range v {
			req.Header.Add(k, h)
		}
	}
	httpClient := &http.Client{Timeout: time.Second * defaultTimeoutSeconds}

	resp, err := httpClient.Do(req)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("HTTP response error: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		errorText, diagnostics, err := getErrorMessage(respBytes)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s: %s", resp.Status, errorText, diagnostics)
	}
	return respBytes, nil
}

// GetBearerToken logs in and gets a Google bearer auth token.
// The user referred to by `cloudhealthEmail` needs to have IAM permissions
// that allow them to read and write from the project's Cloud Healthcare base.
func GetBearerToken() (string, error) {
	ctx := context.Background()
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
	creds, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return "", fmt.Errorf("default creds error: %v", err)
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return "", fmt.Errorf("oauth token error: %v", err)
	}
	return fmt.Sprintf("Bearer %s", token.AccessToken), nil
}

// FHIRHeaders composes suitable FHIR headers, with authentication and content
// type already set
func (fr Repository) FHIRHeaders() (http.Header, error) {
	headers := make(map[string][]string)
	bearerHeader, err := GetBearerToken()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, fmt.Errorf("can't get bearer token: %v", err)
	}
	headers["Content-Type"] = []string{"application/fhir+json; charset=utf-8"}
	headers["Accept"] = []string{"application/fhir+json; charset=utf-8"}
	headers["Authorization"] = []string{bearerHeader}
	headers["Prefer"] = []string{"handling=strict"}
	return headers, nil
}
