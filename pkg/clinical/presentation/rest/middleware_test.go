package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/savannahghi/clinical/pkg/clinical/presentation/rest"
)

func TestIDExtractionMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		requestHeaders  map[string]string
		expectedContext map[interface{}]interface{}
		wantErr         bool
	}{
		{
			name: "Happy Case: Test with valid headers",
			requestHeaders: map[string]string{
				"OrganizationID": "123",
			},
			expectedContext: map[interface{}]interface{}{
				rest.OrganizationIDContextKey: "123",
			},
			wantErr: false,
		},
		{
			name:            "Sad case: missing headers",
			requestHeaders:  map[string]string{},
			expectedContext: map[interface{}]interface{}{},
			wantErr:         true,
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		for key, value := range test.requestHeaders {
			req.Header.Set(key, value)
		}

		res := httptest.NewRecorder()
		middleware := rest.TenantIdentifierExtractionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			for expectedKey, expectedValue := range test.expectedContext {
				ctxValue := ctx.Value(expectedKey)
				if ctxValue != expectedValue {
					t.Errorf("expected context key %v to have value %v, but got %v", expectedKey, expectedValue, ctxValue)
					return
				}
			}
		}))

		middleware.ServeHTTP(res, req)

		if !test.wantErr {
			if res.Code != http.StatusOK {
				t.Errorf("expected error %v, but got %v", test.wantErr, res.Code)
				return
			}
		}
	}
}
