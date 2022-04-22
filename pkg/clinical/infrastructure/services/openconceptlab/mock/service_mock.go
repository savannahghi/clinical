package mock

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

// FakeOCL is an mock of the Open concept lab
type FakeOCL struct {
	MockMakeRequestFn  func(method string, path string, params url.Values, body io.Reader) (*http.Response, error)
	MockListConceptsFn func(
		ctx context.Context, org string, source string, verbose bool, q *string,
		sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
		locale *string, includeRetired *bool,
		includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error)
	MockGetConceptFn func(
		ctx context.Context, org string, source string, concept string,
		includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error)
}

// NewFakeOCLMock initializes a new instance of ocl mock
func NewFakeOCLMock() *FakeOCL {
	return &FakeOCL{
		MockMakeRequestFn: func(method string, path string, params url.Values, body io.Reader) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
			}, nil
		},
		MockListConceptsFn: func(
			ctx context.Context, org string, source string, verbose bool, q *string,
			sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
			locale *string, includeRetired *bool,
			includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error) {
			m := map[string]interface{}{
				"concept": "C00001",
			}
			return []map[string]interface{}{
				m,
			}, nil
		},
		MockGetConceptFn: func(
			ctx context.Context, org string, source string, concept string,
			includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error) {
			m := map[string]interface{}{
				"id": "C12345",
			}
			return m, nil
		},
	}
}

// MakeRequest is a mock implementation of making a http request
func (o *FakeOCL) MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error) {
	return o.MockMakeRequestFn(method, path, params, body)
}

// ListConcepts is a mock implementation of listing concepts
func (o *FakeOCL) ListConcepts(
	ctx context.Context, org string, source string, verbose bool, q *string,
	sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
	locale *string, includeRetired *bool,
	includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error) {
	return o.MockListConceptsFn(ctx, org, source, verbose, q, sortAsc, sortDesc, conceptClass, dataType, locale, includeRetired, includeMappings, includeInverseMappings)
}

// GetConcept is a mock implementation of getting a concept
func (o *FakeOCL) GetConcept(
	ctx context.Context, org string, source string, concept string,
	includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error) {
	return o.MockGetConceptFn(ctx, org, source, concept, includeMappings, includeInverseMappings)
}
