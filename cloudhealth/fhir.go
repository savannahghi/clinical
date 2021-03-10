package cloudhealth

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

	"gitlab.slade360emr.com/go/base"
	"golang.org/x/oauth2/google"
	healthcare "google.golang.org/api/healthcare/v1"
)

const (
	baseFHIRURL           = "https://healthcare.googleapis.com/v1"
	cloudhealthEmail      = "cloudhealth@healthcloud.co.ke"
	defaultTimeoutSeconds = 10
)

// FHIRRestURL composes a FHIR REST URL for manual calls to the FHIR REST API
func (s Service) FHIRRestURL() string {
	return fmt.Sprintf(
		"%s/projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir",
		baseFHIRURL, s.projectID, s.location, s.datasetID, s.fhirStoreID)
}

// GetBearerToken logs in and gets a Google bearer auth token.
// The user referred to by `cloudhealthEmail` needs to have IAM permissions
// that allow them to read and write from teh project's Cloud Healthcare base.
func (s Service) GetBearerToken() (string, error) {
	ctx := context.Background()
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
	creds, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		return "", fmt.Errorf("default creds error: %v", err)
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("oauth token error: %v", err)
	}
	return fmt.Sprintf("Bearer %s", token.AccessToken), nil
}

// FHIRHeaders composes suitable FHIR headers, with authentication and content
// type already set
func (s Service) FHIRHeaders() (http.Header, error) {
	headers := make(map[string][]string)
	bearerHeader, err := s.GetBearerToken()
	if err != nil {
		return nil, fmt.Errorf("can't get bearer token: %v", err)
	}
	headers["Content-Type"] = []string{"application/fhir+json; charset=utf-8"}
	headers["Accept"] = []string{"application/fhir+json; charset=utf-8"}
	headers["Authorization"] = []string{bearerHeader}
	return headers, nil
}

// POSTRequest is used to manually compose POST requests to the FHIR service
//
// - `resourceName` is a FHIR resource name e.g "Patient"
// - `path` is a sub-path e.g `_search` under a resource
// - `params` should be query params, sent as `url.Values`
func (s Service) POSTRequest(
	resourceName string, path string, params url.Values, body io.Reader) ([]byte, error) {
	fhirHeaders, err := s.FHIRHeaders()
	if err != nil {
		return nil, fmt.Errorf("unable to get FHIR headers: %v", err)
	}

	url := fmt.Sprintf(
		"%s/%s/%s?%s", s.FHIRRestURL(), resourceName, path, params.Encode())
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
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
		return nil, fmt.Errorf("HTTP response error: %v", err)
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
			"search: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// CreateFHIRResource creates an FHIR resource.
//
// The payload should be the result of marshalling a resource to JSON
func (s Service) CreateFHIRResource(
	resourceType string, payload map[string]interface{}) ([]byte, error) {
	s.checkPreconditions()
	payload["resourceType"] = resourceType
	fhirService := s.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	payload["resourceType"] = resourceType
	payload["language"] = "EN"
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json.Encode: %v", err)
	}

	parent := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s",
		s.projectID, s.location, s.datasetID, s.fhirStoreID)
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
func (s Service) DeleteFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	s.checkPreconditions()

	fhirService := s.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s",
		s.projectID, s.location, s.datasetID, s.fhirStoreID,
		resourceType, fhirResourceID)
	resp, err := fhirService.Delete(name).Do()
	if err != nil {
		return nil, fmt.Errorf("delete: %v", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(
			"delete: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// GetFHIRResource gets an FHIR resource.
func (s Service) GetFHIRResource(resourceType, fhirResourceID string) ([]byte, error) {
	s.checkPreconditions()
	fhirService := s.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s",
		s.projectID, s.location, s.datasetID, s.fhirStoreID,
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
func (s Service) PatchFHIRResource(
	resourceType, fhirResourceID string, payload []map[string]interface{}) ([]byte, error) {
	s.checkPreconditions()
	fhirService := s.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	if base.IsDebug() {
		log.Printf("FHIR Payload: %#v", payload)
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json.Encode: %v", err)
	}
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s",
		s.projectID, s.location, s.datasetID, s.fhirStoreID,
		resourceType, fhirResourceID)

	call := fhirService.Patch(name, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/json-patch+json")
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("patch: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if base.IsDebug() {
		log.Printf("Patch FHIR Resource %d Response: %s", resp.StatusCode, string(respBytes))
	}
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(
			"patch: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// UpdateFHIRResource updates the entire contents of a resource.
func (s Service) UpdateFHIRResource(
	resourceType, fhirResourceID string, payload map[string]interface{}) ([]byte, error) {
	s.checkPreconditions()
	fhirService := s.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	payload["resourceType"] = resourceType
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json.Encode: %v", err)
	}
	if base.IsDebug() {
		log.Printf("FHIR Update payload: %s", string(jsonPayload))
	}
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/%s/%s",
		s.projectID, s.location, s.datasetID, s.fhirStoreID,
		resourceType, fhirResourceID)
	call := fhirService.Update(name, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("update: %v", err)
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
			"update: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}
	return respBytes, nil
}

// CreateFHIRStore creates an FHIR store.
func (s Service) CreateFHIRStore() (*healthcare.FhirStore, error) {
	s.checkPreconditions()
	storesService := s.healthcareService.Projects.Locations.Datasets.FhirStores
	store := &healthcare.FhirStore{
		DisableReferentialIntegrity: false,
		DisableResourceVersioning:   false,
		EnableUpdateCreate:          true,
		Version:                     "R4",
	}
	parent := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s", s.projectID, s.location, s.datasetID)
	resp, err := storesService.Create(parent, store).FhirStoreId(s.fhirStoreID).Do()
	if err != nil {
		return nil, fmt.Errorf("create FHIR Store: %v", err)
	}
	return resp, nil
}

// GetFHIRStore gets an FHIR store.
func (s Service) GetFHIRStore() (*healthcare.FhirStore, error) {
	s.checkPreconditions()
	storesService := s.healthcareService.Projects.Locations.Datasets.FhirStores
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s",
		s.projectID, s.location, s.datasetID, s.fhirStoreID)
	store, err := storesService.Get(name).Do()
	if err != nil {
		return nil, fmt.Errorf("get FHIR Store: %v", err)
	}
	return store, nil
}

// GetFHIRPatientEverything gets all resources associated with a particular
// patient compartment.
func (s Service) GetFHIRPatientEverything(fhirResourceID string) ([]byte, error) {
	fhirService := s.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	name := fmt.Sprintf(
		"projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir/Patient/%s",
		s.projectID,
		s.location,
		s.datasetID,
		s.fhirStoreID,
		fhirResourceID,
	)

	resp, err := fhirService.PatientEverything(name).Do()
	if err != nil {
		return nil, fmt.Errorf("PatientEverything: %v", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("PatientEverything: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}

	return respBytes, nil
}
