package fhirdataset

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
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

// Repository accesses and updates patient data that is stored on Healthcare
// FHIR repository
type Repository struct {
	healthcareService                           *healthcare.Service
	projectID, location, datasetID, fhirStoreID string
}

// NewFHIRRepository initializes a FHIR repository
func NewFHIRRepository(ctx context.Context) *Repository {
	project := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	_ = serverutils.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
	dataset := serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
	fhirStore := serverutils.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")
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
		return nil, fmt.Errorf("create Data Set: %w", err)
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
		return nil, fmt.Errorf("get Data Set: %w", err)
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
		return nil, fmt.Errorf("create FHIR Store: %w", err)
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
		return nil, fmt.Errorf("get FHIR Store: %w", err)
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

// getErrorMessage un-marshals the error response that is returned from the FHIR server.
// This function should be called when a status code > 299 has been returned
func getErrorMessage(respBytes []byte) (errorText string, diagnostics string, err error) {
	var errorResponse dto.ErrorResponse
	err = json.Unmarshal(respBytes, &errorResponse)
	if err != nil {
		return "", "", fmt.Errorf("could not unmarshal error response: %w", err)
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
		return nil, fmt.Errorf("json.Encode: %w", err)
	}

	parent := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s",
		fr.projectID, fr.location, fr.datasetID, fr.fhirStoreID)
	call := fhirService.Create(
		parent, resourceType, bytes.NewReader(jsonPayload))

	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
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
		return nil, fmt.Errorf("delete: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
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
		return nil, fmt.Errorf("json.Encode: %w", err)
	}
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s",
		fr.projectID, fr.location, fr.datasetID, fr.fhirStoreID,
		resourceType, fhirResourceID)

	call := fhirService.Patch(name, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/json-patch+json")
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("patch: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBytes, err := io.ReadAll(resp.Body)
	if serverutils.IsDebug() {
		log.Printf("Patch FHIR Resource %d Response: %s", resp.StatusCode, string(respBytes))
	}
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
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
		return nil, fmt.Errorf("json.Encode: %w", err)
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
		return nil, fmt.Errorf("update: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
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
		return nil, fmt.Errorf("PatientAllData: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
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
		return nil, fmt.Errorf("read: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("read: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// SearchFHIRResource ...
func (fr Repository) SearchFHIRResource(resourceType string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}

	urlParams := url.Values{}
	for k, v := range params {
		val, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("the search/filter params should all be sent as strings")
		}
		urlParams.Add(k, val)
	}

	path := "_search"
	bs, err := fr.POSTRequest(resourceType, path, urlParams, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to search: %w", err)
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(bs, &respMap)
	if err != nil {
		return nil, fmt.Errorf(
			"%s could not be found with search params %v: %w", resourceType, params, err)
	}

	mandatoryKeys := []string{"resourceType", "type", "total", "link"}
	for _, k := range mandatoryKeys {
		_, found := respMap[k]
		if !found {
			return nil, fmt.Errorf("server error: mandatory search result key %s not found", k)
		}
	}
	resourceType, ok := respMap["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("server error: the resourceType is not a string")
	}
	if resourceType != "Bundle" {
		return nil, fmt.Errorf(
			"server error: the resourceType value is not 'Bundle' as expected")
	}

	resultType, ok := respMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("server error: the search result type value is not a string")
	}
	if resultType != "searchset" {
		return nil, fmt.Errorf("server error: the type value is not 'searchset' as expected")
	}

	respEntries := respMap["entry"]
	if respEntries == nil {
		return []map[string]interface{}{}, nil
	}
	entries, ok := respEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"server error: entries is not a list of maps, it is: %T", respEntries)
	}

	results := []map[string]interface{}{}
	for _, en := range entries {
		entry, ok := en.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(
				"server error: expected each entry to be map, they are %T instead", en)
		}
		expectedKeys := []string{"fullUrl", "resource", "search"}
		for _, k := range expectedKeys {
			_, found := entry[k]
			if !found {
				return nil, fmt.Errorf("server error: FHIR search entry does not have key '%s'", k)
			}
		}

		resource, ok := entry["resource"].(map[string]interface{})
		if !ok {
			{
				return nil, fmt.Errorf("server error: result entry %#v is not a map", entry["resource"])
			}
		}
		results = append(results, resource)
	}

	return results, nil
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
		return nil, fmt.Errorf("unable to get FHIR headers: %w", err)
	}
	url := fmt.Sprintf(
		"%s/%s/%s?%s", fr.FHIRRestURL(), resourceName, path, params.Encode())
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to compose FHIR POST request: %w", err)
	}
	for k, v := range fhirHeaders {
		for _, h := range v {
			req.Header.Add(k, h)
		}
	}
	httpClient := &http.Client{Timeout: time.Second * defaultTimeoutSeconds}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP response error: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
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
		return "", fmt.Errorf("default creds error: %w", err)
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("oauth token error: %w", err)
	}
	return fmt.Sprintf("Bearer %s", token.AccessToken), nil
}

// FHIRHeaders composes suitable FHIR headers, with authentication and content
// type already set
func (fr Repository) FHIRHeaders() (http.Header, error) {
	headers := make(map[string][]string)
	bearerHeader, err := GetBearerToken()
	if err != nil {
		return nil, fmt.Errorf("can't get bearer token: %w", err)
	}
	headers["Content-Type"] = []string{"application/fhir+json; charset=utf-8"}
	headers["Accept"] = []string{"application/fhir+json; charset=utf-8"}
	headers["Authorization"] = []string{bearerHeader}
	return headers, nil
}
