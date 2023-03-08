package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FHIRUsecaseMock struct implements mocks of FHIR methods.
type FHIRUsecaseMock struct {
	MockFindOrganizationByIDFn          func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error)
	MockCreateFHIRObservationFn         func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error)
	MockCreateFHIRAllergyIntoleranceFn  func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	MockCreateFHIRMedicationStatementFn func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)
	MockCreateFHIROrganizationFn        func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
}

// NewFHIRUsecaseMock initializes a new instance of FHIR mock
func NewFHIRUsecaseMock() *FHIRUsecaseMock {
	UUID := uuid.New().String()

	return &FHIRUsecaseMock{

		MockCreateFHIROrganizationFn: func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
			return &domain.FHIROrganizationRelayPayload{
				Resource: &domain.FHIROrganization{
					ID: &UUID,
				},
			}, nil
		},
		MockCreateFHIRAllergyIntoleranceFn: func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
			return &domain.FHIRAllergyIntoleranceRelayPayload{}, nil
		},

		MockCreateFHIRObservationFn: func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
			return &domain.FHIRObservationRelayPayload{}, nil
		},

		MockFindOrganizationByIDFn: func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
			return &domain.FHIROrganizationRelayPayload{}, nil
		},
		MockCreateFHIRMedicationStatementFn: func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
			return &domain.FHIRMedicationStatementRelayPayload{}, nil
		},
	}
}

// CreateFHIROrganization is a mock implementation of CreateFHIROrganization method
func (fh *FHIRUsecaseMock) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.MockCreateFHIROrganizationFn(ctx, input)
}

// CreateFHIRAllergyIntolerance is a mock implementation of CreateFHIRAllergyIntolerance method
func (fh *FHIRUsecaseMock) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	return fh.MockCreateFHIRAllergyIntoleranceFn(ctx, input)
}

// FindOrganizationByID is a mock implementation of FindOrganizationByID method
func (fh *FHIRUsecaseMock) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	return fh.MockFindOrganizationByIDFn(ctx, organizationID)
}

// CreateFHIRMedicationStatement is a mock implementation of CreateFHIRMedicationStatement method
func (fh *FHIRUsecaseMock) CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
	return fh.MockCreateFHIRMedicationStatementFn(ctx, input)
}
