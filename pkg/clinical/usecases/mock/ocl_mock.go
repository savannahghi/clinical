package mock

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

// OCLMock instantiates all the mock functions
type OCLMock struct {
	MakeRequestFn  func(method string, path string, params url.Values, body io.Reader) (*http.Response, error)
	ListConceptsFn func(
		ctx context.Context, org string, source string, verbose bool, q *string,
		sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
		locale *string, includeRetired *bool,
		includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error)
	GetConceptFn func(
		ctx context.Context, org string, source string, concept string,
		includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error)
}

// NewOCLMock is a new instance of OCLMock
func NewOCLMock() *OCLMock {
	return &OCLMock{
		MakeRequestFn: func(method string, path string, params url.Values, body io.Reader) (*http.Response, error) {
			return &http.Response{}, nil
		},
		ListConceptsFn: func(
			ctx context.Context, org string, source string, verbose bool, q *string,
			sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
			locale *string, includeRetired *bool,
			includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error) {
			return nil, nil
		},
		GetConceptFn: func(
			ctx context.Context, org string, source string, concept string,
			includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error) {
			return nil, nil
		},
	}
}

//MakeRequest is the MakeRequest mock
func (ocl *OCLMock) MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error) {
	return ocl.MakeRequestFn(method, path, params, body)
}

//ListConcepts is the ListConcepts mock
func (ocl *OCLMock) ListConcepts(
	ctx context.Context, org string, source string, verbose bool, q *string,
	sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
	locale *string, includeRetired *bool,
	includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error) {
	return ocl.ListConceptsFn(
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

// GetConcept is the GetConcept mock
func (ocl *OCLMock) GetConcept(
	ctx context.Context, org string, source string, concept string,
	includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error) {
	return ocl.GetConceptFn(
		ctx,
		org,
		source,
		concept,
		includeMappings,
		includeInverseMappings,
	)
}
