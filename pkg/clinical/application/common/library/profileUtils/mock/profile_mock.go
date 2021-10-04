package mock

import (
	"context"

	"github.com/savannahghi/profileutils"
)

// FakeUserProfileRepository is a fake implementation of UserProfile repository in profileutils
type FakeUserProfileRepository struct {
	GetLoggedInUserFn func(ctx context.Context) (*profileutils.UserInfo, error)
}

// GetLoggedInUser is a mock implementation of GetLoggedInUser method
func (f *FakeUserProfileRepository) GetLoggedInUser(ctx context.Context) (*profileutils.UserInfo, error) {
	return f.GetLoggedInUserFn(ctx)
}
