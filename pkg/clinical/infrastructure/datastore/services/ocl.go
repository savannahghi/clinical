package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/savannahghi/serverutils"
)

// constants used to configure the OCL service
const (
	OCLAPIURLEnvVarName  = "OPENCONCEPTLAB_API_URL"
	OCLTokenEnvVarName   = "OPENCONCEPTLAB_TOKEN"
	OCLAPITimeoutSeconds = 30
)

// ServiceOpenConceptLab represents all the business logic involved in interacting with
// OCL API
type ServiceOpenConceptLab interface {
	MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error)
	GetConcept(
		ctx context.Context, org string, source string, concept string,
		includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error)
	ListConcepts(
		ctx context.Context, org string, source string, verbose bool, q *string,
		sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
		locale *string, includeRetired *bool,
		includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error)
}

// ServiceOpenConceptLabImpl is an OpenConceptLab service
type ServiceOpenConceptLabImpl struct {
	baseURL string
	header  string
}

// NewServiceOpenConceptLabImpl initializes a new OpenConceptLab service
func NewServiceOpenConceptLabImpl() ServiceOpenConceptLab {
	baseURL := serverutils.MustGetEnvVar(OCLAPIURLEnvVarName)
	token := serverutils.MustGetEnvVar(OCLTokenEnvVarName)
	header := fmt.Sprintf("Authorization: Token %s", token)
	return &ServiceOpenConceptLabImpl{
		baseURL: baseURL,
		header:  header,
	}
}

// MakeRequest composes an authenticated OCL request that has the correct content type
func (ocl ServiceOpenConceptLabImpl) MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error) {
	return nil, nil
}

// GetConcept retrieves a single concept from OpenConceptLab.
// The URL that is composed follows this pattern: GET /orgs/:org/sources/:source/[:sourceVersion/]concepts/:concept/
// e.g GET /orgs/WHO/sources/ICD-10-2010/concepts/A15.1/?includeInverseMappings=true
func (ocl ServiceOpenConceptLabImpl) GetConcept(
	ctx context.Context, org string, source string, concept string,
	includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error) {
	return nil, nil
}

// ListConcepts searches for matching concepts on OpenConceptLab
// The URL that is composed follows this pattern: GET /orgs/:org/sources/:source/[:sourceVersion/]concepts/
// e.g GET /orgs/PEPFAR-Test7/sources/MER/concepts/?conceptClass="Symptom"+OR+"Diagnosis"
func (ocl ServiceOpenConceptLabImpl) ListConcepts(
	ctx context.Context, org string, source string, verbose bool, q *string,
	sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
	locale *string, includeRetired *bool,
	includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error) {
	return nil, nil
}
