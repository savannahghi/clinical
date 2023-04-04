package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
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

	MockUpdateProgramFHIRTenantIDFn func(ctx context.Context, programID string, tenantID string) error
}

// NewFakeMyCareHubServiceMock initializes a new instance of mycarehub mock
func NewFakeMyCareHubServiceMock() *FakeMyCareHubService {
	return &FakeMyCareHubService{
		MockUserProfileFn: func(ctx context.Context, userID string) (*domain.User, error) {
			dob := time.Now()
			return &domain.User{
				ID:       new(string),
				Username: gofakeit.Username(),
				UserType: "STAFF",
				Name:     gofakeit.Name(),
				Gender:   "MALE",
				Active:   false,
				Flavour:  "PRO",
				Avatar:   "",
				Contacts: &domain.Contact{
					ID:           new(string),
					ContactType:  "PHONE",
					ContactValue: gofakeit.PhoneFormatted(),
					Active:       true,
					OptedIn:      true,
				},
				DateOfBirth: &dob,
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
		MockUpdateProgramFHIRTenantIDFn: func(ctx context.Context, programID, tenantID string) error {
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

// UpdateProgramFHIRTenantID amocks the update of program mycarehub's program fhir tenant id
func (s *FakeMyCareHubService) UpdateProgramFHIRTenantID(ctx context.Context, programID string, tenantID string) error {
	return s.MockUpdateProgramFHIRTenantIDFn(ctx, programID, tenantID)
}
