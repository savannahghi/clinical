package mock

import (
	"context"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
)

// FakeServiceOnboarding is a mock of the Onboarding service.
type FakeServiceOnboarding struct {
	CreateUserProfileFn func(ctx context.Context, payload *dto.RegisterUserInput) error
}

// CreateUserProfile is a fake implementation of the CreateUserProfile method
func (f *FakeServiceOnboarding) CreateUserProfile(ctx context.Context, payload *dto.RegisterUserInput) error {
	return f.CreateUserProfileFn(ctx, payload)
}
