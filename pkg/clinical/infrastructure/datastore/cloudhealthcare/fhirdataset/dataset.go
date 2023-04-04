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
	"strconv"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/domain"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/serverutils"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/healthcare/v1"
)

// constants used to configure the Google Cloud Healthcare API
const (
	baseFHIRURL           = "https://healthcare.googleapis.com/v1"
	defaultTimeoutSeconds = 10
)

// Repository accesses and updates patient data that is stored on Healthcare
// FHIR repository
type Repository struct {
	healthcareService                           *healthcare.Service
	projectID, location, datasetID, fhirStoreID string
	parent                                      string
	datasetName                                 string
	fhirStoreName                               string
}

// NewFHIRRepository initializes a FHIR repository
func NewFHIRRepository(_ context.Context, hsv *healthcare.Service, projectID, datasetID, datasetLocation, fhirStoreID string) *Repository {
	return &Repository{
		healthcareService: hsv,
		projectID:         projectID,
		location:          datasetLocation,
		datasetID:         datasetID,
		fhirStoreID:       fhirStoreID,
		parent:            fmt.Sprintf("projects/%s/locations/%s", projectID, datasetLocation),
		datasetName:       fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, datasetLocation, datasetID),
		fhirStoreName:     fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, datasetLocation, datasetID, fhirStoreID),
	}
}

// CreateDataset creates a dataset and returns it's name
func (fr Repository) CreateDataset() (*healthcare.Operation, error) {
	fr.checkPreconditions()
	datasetsService := fr.healthcareService.Projects.Locations.Datasets

	resp, err := datasetsService.Create(fr.parent, &healthcare.Dataset{}).DatasetId(fr.datasetID).Do()
	if err != nil {
		return nil, fmt.Errorf("create Data Set: %w", err)
	}

	return resp, nil
}

// GetDataset gets a dataset.
func (fr Repository) GetDataset() (*healthcare.Dataset, error) {
	fr.checkPreconditions()
	datasetsService := fr.healthcareService.Projects.Locations.Datasets

	resp, err := datasetsService.Get(fr.datasetName).Do()
	if err != nil {
		return nil, fmt.Errorf("get Data Set: %w", err)
	}

	return resp, nil
}

// CreateFHIRStore creates a FHIR store.
func (fr Repository) CreateFHIRStore() (*healthcare.FhirStore, error) {
	fr.checkPreconditions()
	storesService := fr.healthcareService.Projects.Locations.Datasets.FhirStores
	store := &healthcare.FhirStore{
		DisableReferentialIntegrity: false,
		DisableResourceVersioning:   false,
		EnableUpdateCreate:          true,
		Version:                     "R4",
	}

	resp, err := storesService.Create(fr.datasetName, store).FhirStoreId(fr.fhirStoreID).Do()
	if err != nil {
		return nil, fmt.Errorf("create FHIR Store: %w", err)
	}

	return resp, nil
}

// GetFHIRStore gets an FHIR store.
func (fr Repository) GetFHIRStore() (*healthcare.FhirStore, error) {
	fr.checkPreconditions()
	storesService := fr.healthcareService.Projects.Locations.Datasets.FhirStores

	store, err := storesService.Get(fr.fhirStoreName).Do()
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
func (fr Repository) CreateFHIRResource(resourceType string, payload map[string]interface{}, resource interface{}) error {
	fr.checkPreconditions()

	payload["resourceType"] = resourceType

	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	payload["resourceType"] = resourceType
	payload["language"] = "EN"

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Encode: %w", err)
	}

	call := fhirService.Create(fr.fhirStoreName, resourceType, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")

	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode > 299 {
		errorText, diagnostics, err := getErrorMessage(respBytes)
		if err != nil {
			return err
		}

		return fmt.Errorf("%s: %s: %s", resp.Status, errorText, diagnostics)
	}

	err = json.Unmarshal(respBytes, resource)
	if err != nil {
		return fmt.Errorf("unable to unmarshal %s response JSON: data: %v\n, error: %w", resourceType, string(respBytes), err)
	}

	return nil
}

// DeleteFHIRResource deletes an FHIR resource.
func (fr Repository) DeleteFHIRResource(resourceType, fhirResourceID string) error {
	fr.checkPreconditions()

	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	fhirResource := fmt.Sprintf("%s/fhir/%s/%s", fr.fhirStoreName, resourceType, fhirResourceID)

	resp, err := fhirService.Delete(fhirResource).Do()
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf(
			"delete: status %d %s: %s", resp.StatusCode, resp.Status, string(respBytes))
	}

	return nil
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
	resourceType, fhirResourceID string, payload []map[string]interface{}, resource interface{}) error {
	fr.checkPreconditions()

	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir

	if serverutils.IsDebug() {
		log.Printf("FHIR Payload: %#v", payload)
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Encode: %w", err)
	}

	fhirResource := fmt.Sprintf("%s/fhir/%s/%s", fr.fhirStoreName, resourceType, fhirResourceID)
	call := fhirService.Patch(fhirResource, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/json-patch+json")

	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("patch: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if serverutils.IsDebug() {
		log.Printf("Patch FHIR Resource %d Response: %s", resp.StatusCode, string(respBytes))
	}

	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf(
			"patch: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}

	err = json.Unmarshal(respBytes, resource)
	if err != nil {
		return fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %w",
			resourceType, string(respBytes), err)
	}

	return nil
}

// UpdateFHIRResource updates the entire contents of a resource.
func (fr Repository) UpdateFHIRResource(
	resourceType, fhirResourceID string, payload map[string]interface{}, resource interface{}) error {
	fr.checkPreconditions()

	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir

	payload["resourceType"] = resourceType

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Encode: %w", err)
	}

	if serverutils.IsDebug() {
		log.Printf("FHIR Update payload: %s", string(jsonPayload))
	}

	fhirResource := fmt.Sprintf("%s/fhir/%s/%s", fr.fhirStoreName, resourceType, fhirResourceID)
	call := fhirService.Update(fhirResource, bytes.NewReader(jsonPayload))
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")

	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf(
			"update: status %d %s: %s", resp.StatusCode, resp.Status, respBytes)
	}

	err = json.Unmarshal(respBytes, resource)
	if err != nil {
		return fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %w",
			resourceType, string(respBytes), err)
	}

	return nil
}

// GetFHIRPatientAllData gets all resources associated with a particular
// patient compartment.
func (fr Repository) GetFHIRPatientAllData(fhirResourceID string) ([]byte, error) {
	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	patientResource := fmt.Sprintf("%s/fhir/Patient/%s", fr.fhirStoreName, fhirResourceID)

	resp, err := fhirService.PatientEverything(patientResource).Do()
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
func (fr Repository) GetFHIRResource(resourceType, fhirResourceID string, resource interface{}) error {
	fr.checkPreconditions()
	fhirService := fr.healthcareService.Projects.Locations.Datasets.FhirStores.Fhir
	fhirResource := fmt.Sprintf("%s/fhir/%s/%s", fr.fhirStoreName, resourceType, fhirResourceID)
	call := fhirService.Read(fhirResource)
	call.Header().Set("Content-Type", "application/fhir+json;charset=utf-8")

	resp, err := call.Do()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode > 299 {
		_, diagnostics, err := getErrorMessage(respBytes)
		if err != nil {
			return err
		}

		return fmt.Errorf("%s", diagnostics)
	}

	err = json.Unmarshal(respBytes, resource)
	if err != nil {
		return fmt.Errorf(
			"unable to unmarshal %s , id:%s ,response JSON: data: %v\n, error: %w",
			resourceType, fhirResourceID, string(respBytes), err)
	}

	return nil
}

// SearchFHIRResource ...
func (fr Repository) SearchFHIRResource(resourceType string, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error) {
	err := pagination.Validate()
	if err != nil {
		return nil, err
	}

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}

	if !pagination.Skip {
		params["_count"] = strconv.Itoa(*pagination.First)
		if pagination.After != "" {
			params["_page_token"] = pagination.After
		}
	}

	urlParams := url.Values{}

	for k, v := range params {
		val, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("the search/filter param: %s should all be sent as strings", k)
		}

		urlParams.Add(k, val)
	}

	urlParams.Add("_tag", fmt.Sprintf("http://mycarehub/tenant-identification/organisation|%s", tenant.OrganizationID))
	urlParams.Add("_tag", fmt.Sprintf("http://mycarehub/tenant-identification/facility|%s", tenant.FacilityID))

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

	response := domain.PagedFHIRResource{
		Resources:       []map[string]interface{}{},
		HasNextPage:     false,
		NextCursor:      "",
		HasPreviousPage: false,
		PreviousCursor:  "",
		TotalCount:      0,
	}

	response.TotalCount = int(respMap["total"].(float64))

	respEntries := respMap["entry"]
	if respEntries == nil {
		return &response, nil
	}

	entries, ok := respEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"server error: entries is not a list of maps, it is: %T", respEntries)
	}

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
			return nil, fmt.Errorf("server error: result entry %#v is not a map", entry["resource"])
		}

		response.Resources = append(response.Resources, resource)
	}

	linksEntries := respMap["link"]
	if linksEntries == nil {
		return &response, nil
	}

	links, ok := linksEntries.([]interface{})
	if !ok {
		return nil, fmt.Errorf(
			"server error: entries is not a list of maps, it is: %T", linksEntries)
	}

	for _, en := range links {
		link, ok := en.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(
				"server error: expected each link to be map, they are %T instead", en)
		}

		if link["relation"].(string) == "next" {
			u, err := url.Parse(link["url"].(string))
			if err != nil {
				return nil, fmt.Errorf("server error: cannot parse url in link: %w", err)
			}

			params, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				return nil, fmt.Errorf("server error: cannot parse url params in link: %w", err)
			}

			cursor := params["_page_token"][0]

			response.HasNextPage = true
			response.NextCursor = cursor
		}
	}

	return &response, nil
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
