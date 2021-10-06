package mock

import (
	"context"

	"github.com/savannahghi/profileutils"
)

type FakeUserProfileRepository struct {
	GetLoggedInUserFn func(ctx context.Context) (*profileutils.UserInfo, error)
}

func (f *FakeUserProfileRepository) GetLoggedInUser(ctx context.Context) (*profileutils.UserInfo, error) {
	return f.GetLoggedInUserFn(ctx)
}
