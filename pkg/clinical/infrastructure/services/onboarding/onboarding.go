package onboarding

import (
	"context"
	"fmt"
	"net/http"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
)

// internal apis definitions
const (
	registerUserURL = "internal/register_user"
)

// ServiceOnboardingImpl represents onboarding usecases
type ServiceOnboardingImpl struct {
	Onboarding extensions.ISCClientExtension
	Basext     extensions.BaseExtension
}

// NewServiceOnboardingImpl returns new instance of ServiceEngagementImpl
func NewServiceOnboardingImpl(
	onboarding extensions.ISCClientExtension,
	ext extensions.BaseExtension,
) *ServiceOnboardingImpl {
	return &ServiceOnboardingImpl{
		Onboarding: onboarding,
		Basext:     ext,
	}
}

// ServiceOnboarding represents onboarding usecases
type ServiceOnboarding interface {
	CreateUserProfile(ctx context.Context, payload dto.RegisterUserPayload) error
}

//CreateUserProfile makes the request to register a user
func (on *ServiceOnboardingImpl) CreateUserProfile(ctx context.Context, payload dto.RegisterUserPayload) error {

	res, err := on.Onboarding.MakeRequest(ctx, http.MethodPost, registerUserURL, payload)
	if err != nil {
		return fmt.Errorf("unable to send request, error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("register user failed with status code: %v", res.StatusCode)
	}

	return nil
}
