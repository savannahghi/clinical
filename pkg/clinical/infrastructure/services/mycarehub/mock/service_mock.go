package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FakeMyCareHubService is an mock of the MyCareHub service
type FakeMyCareHubService struct {
	MockUserProfileFn func(ctx context.Context, userID string) (*domain.User, error)

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
}

// NewFakeMyCareHubServiceMock initializes a new instance of mycarehub mock
func NewFakeMyCareHubServiceMock() *FakeMyCareHubService {
	return &FakeMyCareHubService{
		MockUserProfileFn: func(ctx context.Context, userID string) (*domain.User, error) {
			return &domain.User{
				ID:       new(string),
				Username: "",
				UserType: "",
				Name:     "",
				Gender:   "",
				Active:   false,
				Flavour:  "",
				Avatar:   "",
			}, nil
		},
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
	}
}

// UserProfile gets the profile of the user with the indicated ID
func (s *FakeMyCareHubService) UserProfile(
	ctx context.Context,
	userID string,
) (*domain.User, error) {
	return s.MockUserProfileFn(ctx, userID)
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
