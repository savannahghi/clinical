package mock

import (
	"context"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
)

// FakeOnboarding ...
type FakeOnboarding struct {
	CreateUserProfileFn func(ctx context.Context, payload *dto.RegisterUserInput) error
}

// CreateUserProfile ...
func (e *FakeOnboarding) CreateUserProfile(ctx context.Context, payload *dto.RegisterUserInput) error {
	return e.CreateUserProfileFn(ctx, payload)
}
