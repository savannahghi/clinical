package clinical

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FindOrganizationByID finds and retrieves organization details using the specified organization ID
func (c *UseCasesClinicalImpl) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("organization ID is required")
	}

	return c.infrastructure.FHIR.FindOrganizationByID(ctx, organizationID)
}
