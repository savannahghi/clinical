package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/rest"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/clinical/mock"
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
				utils.OrganizationIDContextKey: "123",
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
		req, _ := http.NewRequest("GET", "http://localhost:8000", nil)
		for key, value := range test.requestHeaders {
			req.Header.Set(key, value)
		}

		res := httptest.NewRecorder()

		engine := gin.New()
		engine.Use(rest.TenantIdentifierExtractionMiddleware(mock.NewFHIRUsecaseMock()))

		engine.GET("", func(c *gin.Context) {
			ctx := c.Request.Context()
			for expectedKey, expectedValue := range test.expectedContext {
				ctxValue := ctx.Value(expectedKey)
				if ctxValue != expectedValue {
					t.Errorf("expected context key %v to have value %v, but got %v", expectedKey, expectedValue, ctxValue)
					return
				}
			}
			c.String(http.StatusOK, "OK")
		})

		if !test.wantErr {
			if res.Code != http.StatusOK {
				t.Errorf("expected error %v, but got %v", test.wantErr, res.Code)
				return
			}
		}
	}
}
