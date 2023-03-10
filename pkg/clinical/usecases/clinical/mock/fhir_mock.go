package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FHIRUsecaseMock struct implements mocks of FHIR methods.
type FHIRUsecaseMock struct {
	MockFindOrganizationByIDFn func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error)
}

// NewFHIRUsecaseMock initializes a new instance of FHIR mock
func NewFHIRUsecaseMock() *FHIRUsecaseMock {
	return &FHIRUsecaseMock{
		MockFindOrganizationByIDFn: func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
			return &domain.FHIROrganizationRelayPayload{}, nil
		},
	}
}

// FindOrganizationByID is a mock implementation of FindOrganizationByID method
func (fh *FHIRUsecaseMock) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.MockFindOrganizationByIDFn(ctx, organizationID)
}
