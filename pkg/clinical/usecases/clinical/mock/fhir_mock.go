package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FHIRUsecaseMock struct implements mocks of FHIR methods.
type FHIRUsecaseMock struct {
	MockGetFHIROrganizationFn func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error)
}

// NewFHIRUsecaseMock initializes a new instance of FHIR mock
func NewFHIRUsecaseMock() *FHIRUsecaseMock {
	return &FHIRUsecaseMock{
		MockGetFHIROrganizationFn: func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
			return &domain.FHIROrganizationRelayPayload{}, nil
		},
	}
}

// GetFHIROrganization is a mock implementation of GetFHIROrganization method
func (fh *FHIRUsecaseMock) GetFHIROrganization(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.MockGetFHIROrganizationFn(ctx, organizationID)
}
