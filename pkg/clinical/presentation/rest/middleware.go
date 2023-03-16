package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/serverutils"
)

// Validators defines the methods used to validate the various identifiers that the api expects
type Validators interface {
	FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error)
}

// OrganisationValidator verifies that the provided organisation exists in clinical
// to ensure the request comes from a known/registered organisation
func OrganisationValidator(v Validators, identifier string) error {
	_, err := v.FindOrganizationByID(context.Background(), identifier)
	if err != nil {
		return fmt.Errorf("failed to find provided organisation")
	}

	return nil
}

// TenantIdentifier is a type representing a header name and a corresponding context key
// The header name is what will be used to extract the specified header and the context key
// Will be the key value used when adding the header in the request context
type TenantIdentifier struct {
	HeaderKey     string
	ContextKey    utils.ContextKey
	ValidatorFunc func(v Validators, identifier string) error
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
func TenantIdentifierExtractionMiddleware(validator Validators) gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := []TenantIdentifier{
			{
				HeaderKey:     "Clinical-Organization-ID",
				ContextKey:    utils.OrganizationIDContextKey,
				ValidatorFunc: OrganisationValidator,
			},
			{
				HeaderKey:     "Clinical-Facility-ID",
				ContextKey:    utils.FacilityIDContextKey,
				ValidatorFunc: OrganisationValidator,
			},
		}

		for _, header := range headers {
			headerValue := c.GetHeader(header.HeaderKey)
			if headerValue == "" {
				err := fmt.Errorf("expected `%s` header to be included in the request", header.HeaderKey)
				handleError(c.Writer, err)
				c.Abort()

				return
			}

			err := header.ValidatorFunc(validator, headerValue)
			if err != nil {
				err := fmt.Errorf("invalid `%s` header value: %s", header.HeaderKey, headerValue)
				handleError(c.Writer, err)
				c.Abort()

				return
			}

			c.Set(string(header.ContextKey), headerValue)
		}

		c.Next()
	}
}
