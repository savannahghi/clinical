package ocl

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
)

// UseCases represents all the Open Concept Lab business logic
type UseCases interface {
	MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error)
	ListConcepts(
		ctx context.Context, org string, source string, verbose bool, q *string,
		sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
		locale *string, includeRetired *bool,
		includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error)
	GetConcept(
		ctx context.Context, org string, source string, concept string,
		includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error)
}

// UseCasesImpl represents the OCL usecase implementation
type UseCasesImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewUseCasesImpl initializes Open Concept Lab implementation
func NewUseCasesImpl(infrastructure infrastructure.Infrastructure) *UseCasesImpl {
	return &UseCasesImpl{infrastructure: infrastructure}
}

// MakeRequest composes an authenticated OCL request that has the correct content type
func (ocl *UseCasesImpl) MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error) {
	return ocl.infrastructure.OpenConceptLab.MakeRequest(method, path, params, body)
}

// ListConcepts retrieves a single concept from OpenConceptLab.
func (ocl *UseCasesImpl) ListConcepts(
	ctx context.Context, org string, source string, verbose bool, q *string,
	sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
	locale *string, includeRetired *bool,
	includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error) {
	return ocl.infrastructure.OpenConceptLab.ListConcepts(
		ctx,
		org,
		source,
		verbose,
		q,
		sortAsc,
		sortDesc,
		conceptClass,
		dataType,
		locale,
		includeRetired,
		includeMappings,
		includeInverseMappings,
	)
}

// GetConcept searches for matching concepts on OpenConceptLab
func (ocl *UseCasesImpl) GetConcept(
	ctx context.Context, org string, source string, concept string,
	includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error) {
	return ocl.infrastructure.OpenConceptLab.GetConcept(
		ctx,
		org,
		source,
		concept,
		includeMappings,
		includeInverseMappings,
	)
}
