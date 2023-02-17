package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/serverutils"
)

// TenantIdentifier is a type representing a header name and a corresponding context key
// The header name is what will be used to extract the specified header and the context key
// Will be the key value used when adding the header in the request context
type TenantIdentifier struct {
	HeaderKey  string
	ContextKey utils.ContextKey
}

type errResponse struct {
	Err string `json:"error"`
}

func handleError(w http.ResponseWriter, err error) {
	serverutils.WriteJSONResponse(
		w,
		errResponse{
			Err: err.Error(),
		},
		http.StatusBadRequest,
	)
}

// TenantIdentifierExtractionMiddleware is a middleware function that extracts the `organizationID`,
// `programID`, and `facilityID` values from the request and adds them to the request
// context. These IDs can then be used by downstream handlers or middleware to perform
// tasks such as filtering, or database queries
// Note that this middleware assumes that the IDs are included in the request as headers
// and it does not perform any validation or sanitization of the ID values.
func TenantIdentifierExtractionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := []TenantIdentifier{
			{
				HeaderKey:  "OrganizationID",
				ContextKey: utils.OrganizationIDContextKey,
			},
			// TODO: Add more headers here as needed e.g FacilityID, ProgramID
		}

		for _, header := range headers {
			headerValue := r.Header.Get(header.HeaderKey)
			if headerValue == "" {
				err := fmt.Errorf("expected `%s` header to be included in the request", header.HeaderKey)
				handleError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), header.ContextKey, headerValue)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
