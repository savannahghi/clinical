package mycarehub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

const (
	getUserProfile      = "internal/user-profile/%s"
	addFHIRIDToProfile  = "internal/add-fhir-id"
	addFHIRIDToFacility = "internal/facilities"
)

// IServiceMyCareHub represents mycarehub usecases
type IServiceMyCareHub interface {
	UserProfile(
		ctx context.Context,
		userID string,
	) (*domain.User, error)

	AddFHIRIDToPatientProfile(
		ctx context.Context,
		fhirID string,
		clientID string,
	) error
	AddFHIRIDToFacility(
		ctx context.Context,
		fhirID string,
		facilityID string,
	) error
}

// ServiceMyCareHubImpl represents mycarehub usecases
type ServiceMyCareHubImpl struct {
	MyCareHubClient extensions.ISCClientExtension
	baseExt         extensions.BaseExtension
}

// NewServiceMyCareHub returns new instance of ServiceMyCareHubImpl
func NewServiceMyCareHub(
	pr extensions.ISCClientExtension,
	baseExt extensions.BaseExtension,
) *ServiceMyCareHubImpl {
	return &ServiceMyCareHubImpl{
		MyCareHubClient: pr,
		baseExt:         baseExt,
	}
}

// UserProfile gets the profile of the user with the indicated ID
func (s ServiceMyCareHubImpl) UserProfile(
	ctx context.Context,
	userID string,
) (*domain.User, error) {
	getUserProfileURL := fmt.Sprintf(getUserProfile, userID)
	resp, err := s.MyCareHubClient.MakeRequest(
		ctx,
		http.MethodGet,
		getUserProfileURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch client profile from mycarehub service: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user profile: %w, with status code %v",
			err,
			resp.StatusCode,
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read data returned by mycarehub service: %w",
			err,
		)
	}

	var profile domain.User
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return nil, fmt.Errorf(
			"can't unmarshal user profile data: %w",
			err,
		)
	}

	// just a sanity check
	if profile.ID == nil {
		return nil, fmt.Errorf("failed to get the user profile")
	}

	return &profile, nil
}

// AddFHIRIDToPatientProfile makes an interservice call to mycarehub service so that the created patient
// FHIR ID can be added to their profile on mycarehub
func (s ServiceMyCareHubImpl) AddFHIRIDToPatientProfile(
	ctx context.Context,
	fhirID string,
	clientID string,
) error {
	type requestPayload struct {
		FhirID   string `json:"fhirID"`
		ClientID string `json:"clientID"`
	}

	resp, err := s.MyCareHubClient.MakeRequest(
		ctx,
		http.MethodPatch,
		addFHIRIDToProfile,
		&requestPayload{FhirID: fhirID, ClientID: clientID},
	)
	if err != nil {
		return fmt.Errorf("failed to make a request to mycarehub service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update patient fhir ID : %w, with status code %v",
			err,
			resp.StatusCode,
		)
	}

	return nil
}

// AddFHIRIDToFacility makes an interservice call to mycarehub service and updated the FHIR Organization ID of a given facility
func (s ServiceMyCareHubImpl) AddFHIRIDToFacility(
	ctx context.Context,
	fhirID string,
	facilityID string,
) error {
	type requestPayload struct {
		FacilityID string `json:"facilityID"`
		FhirID     string `json:"fhirOrganisationID"`
	}

	resp, err := s.MyCareHubClient.MakeRequest(
		ctx,
		http.MethodPost,
		addFHIRIDToFacility,
		&requestPayload{FhirID: fhirID, FacilityID: facilityID},
	)
	if err != nil {
		return fmt.Errorf("failed to make a request to mycarehub service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update facility fhir ID : %w, with status code %v",
			err,
			resp.StatusCode,
		)
	}

	return nil
}
