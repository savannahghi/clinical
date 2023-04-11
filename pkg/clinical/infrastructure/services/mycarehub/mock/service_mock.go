package mock

import (
	"context"
)

// FakeMyCareHubService is an mock of the MyCareHub service
type FakeMyCareHubService struct {
	MockAddFHIRIDToPatientProfileFn func(
		ctx context.Context,
		fhirID string,
		clientID string,
	) error
	MockAddFHIRIDToFacilityFn func(
		ctx context.Context,
		fhirID string,
		facilityID string,
	) error

	MockUpdateProgramFHIRTenantIDFn func(ctx context.Context, programID string, tenantID string) error
}

// NewFakeMyCareHubServiceMock initializes a new instance of mycarehub mock
func NewFakeMyCareHubServiceMock() *FakeMyCareHubService {
	return &FakeMyCareHubService{
		MockAddFHIRIDToPatientProfileFn: func(
			ctx context.Context,
			fhirID string,
			clientID string,
		) error {
			return nil
		},
		MockAddFHIRIDToFacilityFn: func(
			ctx context.Context,
			fhirID string,
			facilityID string,
		) error {
			return nil
		},
		MockUpdateProgramFHIRTenantIDFn: func(ctx context.Context, programID, tenantID string) error {
			return nil
		},
	}
}

// AddFHIRIDToPatientProfile adds a FHIR ID to a patient
func (s *FakeMyCareHubService) AddFHIRIDToPatientProfile(
	ctx context.Context,
	fhirID string,
	clientID string,
) error {
	return s.MockAddFHIRIDToPatientProfileFn(ctx, fhirID, clientID)
}

// AddFHIRIDToFacility adds a FHIR ID to a facility
func (s *FakeMyCareHubService) AddFHIRIDToFacility(
	ctx context.Context,
	fhirID string,
	facilityID string,
) error {
	return s.MockAddFHIRIDToFacilityFn(ctx, fhirID, facilityID)
}

// UpdateProgramFHIRTenantID amocks the update of program mycarehub's program fhir tenant id
func (s *FakeMyCareHubService) UpdateProgramFHIRTenantID(ctx context.Context, programID string, tenantID string) error {
	return s.MockUpdateProgramFHIRTenantIDFn(ctx, programID, tenantID)
}
