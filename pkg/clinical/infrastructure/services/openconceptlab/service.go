// Package openconceptlab provides APIs to interact with an OpenConceptLab API
// server
package openconceptlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/savannahghi/serverutils"
)

// constants used to configure the OCL service
const (
	OCLAPIURLEnvVarName  = "OPENCONCEPTLAB_API_URL"
	OCLTokenEnvVarName   = "OPENCONCEPTLAB_TOKEN"
	OCLAPITimeoutSeconds = 30
)

// NewServiceOCL creates a new open conceptlab Service
func NewServiceOCL() *Service {
	baseURL := serverutils.MustGetEnvVar(OCLAPIURLEnvVarName)
	token := serverutils.MustGetEnvVar(OCLTokenEnvVarName)
	header := fmt.Sprintf("Authorization: Token %s", token)

	srv := &Service{
		baseURL: baseURL,
		header:  header,
	}
	srv.enforcePreconditions()

	return srv
}

// Service is an OpenConceptLab service
type Service struct {
	baseURL string
	header  string
}

func (s Service) enforcePreconditions() {
	if s.baseURL == "" {
		log.Panicf("Open Concept Lab API Base URL not set in service")
	}

	if s.header == "" {
		log.Panicf("Open Concept Lab API Token header not set in service")
	}
}

// MakeRequest composes an authenticated OCL request that has the correct content type
func (s Service) MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/?%s", s.baseURL, path, params.Encode())

	req, reqErr := http.NewRequest(method, url, body)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.header)

	httpClient := &http.Client{Timeout: time.Second * OCLAPITimeoutSeconds}
	resp, respErr := httpClient.Do(req)

	if respErr != nil {
		return nil, respErr
	}

	return resp, nil
}

// GetConcept retrieves a single concept from OpenConceptLab.
// The URL that is composed follows this pattern: GET /orgs/:org/sources/:source/[:sourceVersion/]concepts/:concept/
// e.g GET /orgs/WHO/sources/ICD-10-2010/concepts/A15.1/?includeInverseMappings=true
func (s Service) GetConcept(
	_ context.Context, org string, source string, concept string,
	includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error) {
	s.enforcePreconditions()

	path := fmt.Sprintf("orgs/%s/sources/%s/concepts/%s", org, source, concept)

	params := url.Values{}
	params.Add("includeMappings", strconv.FormatBool(includeMappings))
	params.Add("includeMappings", strconv.FormatBool(includeInverseMappings))

	resp, err := s.MakeRequest("GET", path, params, nil)

	if err != nil {
		return nil, fmt.Errorf("OCL API request error: %w", err)
	}
	defer resp.Body.Close()

	output := make(map[string]interface{})

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read OCL API response body: %w", err)
	}

	err = json.Unmarshal(data, &output)

	if err != nil {
		return nil, fmt.Errorf(
			"unable to marshal OCL get concept response %s to JSON: %w", string(data), err)
	}

	return output, nil
}

// ListConcepts searches for matching concepts on OpenConceptLab
// The URL that is composed follows this pattern: GET /orgs/:org/sources/:source/[:sourceVersion/]concepts/
// e.g GET /orgs/PEPFAR-Test7/sources/MER/concepts/?conceptClass="Symptom"+OR+"Diagnosis"
func (s Service) ListConcepts(
	_ context.Context, org string, source string, verbose bool, q *string,
	sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
	locale *string, includeRetired *bool,
	includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error) {
	s.enforcePreconditions()

	path := fmt.Sprintf("orgs/%s/sources/%s/concepts", org, source)

	params := url.Values{}
	params.Add("verbose", strconv.FormatBool(verbose))

	if q != nil {
		params.Add("q", *q)
	}

	if sortAsc != nil {
		params.Add("sortAsc", *sortAsc)
	}

	if sortDesc != nil {
		params.Add("sortDesc", *sortDesc)
	}

	if conceptClass != nil {
		params.Add("conceptClass", *conceptClass)
	}

	if dataType != nil {
		params.Add("dataType", *dataType)
	}

	if locale != nil {
		params.Add("locale", *locale)
	}

	if includeRetired != nil {
		if *includeRetired {
			params.Add("includeRetired", "1")
		} else {
			params.Add("includeRetired", "0")
		}
	}

	if includeMappings != nil && *includeMappings {
		params.Add("includeMappings", "true")
	}

	if includeInverseMappings != nil && *includeInverseMappings {
		params.Add("includeReverseMappings", "true")
	}

	resp, err := s.MakeRequest("GET", path, params, nil)
	if err != nil {
		return nil, fmt.Errorf("OCL API request error: %w", err)
	}

	defer resp.Body.Close()

	output := []map[string]interface{}{}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read OCL API response body: %w", err)
	}

	err = json.Unmarshal(data, &output)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to marshal OCL get concept response %s to JSON: %w", string(data), err)
	}

	return output, nil
}
