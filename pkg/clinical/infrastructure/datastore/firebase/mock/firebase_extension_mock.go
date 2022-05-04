package mock

import (
	"context"
	"time"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	fb "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase"
	"github.com/savannahghi/profileutils"

	"cloud.google.com/go/firestore"
)

// FakeFirebase represents a `firestore.Client` fake
type FakeFirebase struct {
	CollectionFn func(path string) *firestore.CollectionRef
	GetAllFn     func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error)
	CreateFn     func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error)
	UpdateFn     func(ctx context.Context, command *fb.UpdateCommand) error
	DeleteFn     func(ctx context.Context, command *fb.DeleteCommand) error
	GetFn        func(ctx context.Context, query *fb.GetSingleQuery) (*firestore.DocumentSnapshot, error)

	//Firebase
	GetUserByPhoneNumberFn func(ctx context.Context, phone string) (*auth.UserRecord, error)
	CreateUserFn           func(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error)
	DeleteUserFn           func(ctx context.Context, uid string) error
	GetUserProfileByIDFn   func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error)

	MockSaveEmailOTPFn                  func(ctx context.Context, email string, optIn bool) error
	MockStageStartEpisodeByBreakGlassFn func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) error
}

// NewFakeFirebaseMock initializes a new NewFakeFirebaseMock
func NewFakeFirebaseMock() *FakeFirebase {
	return &FakeFirebase{
		CollectionFn: func(path string) *firestore.CollectionRef {
			return &firestore.CollectionRef{}
		},
		GetAllFn: func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
			return []*firestore.DocumentSnapshot{
				{
					Ref:        &firestore.DocumentRef{},
					CreateTime: time.Time{},
					UpdateTime: time.Time{},
					ReadTime:   time.Time{},
				},
			}, nil
		},
		CreateFn: func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
			return &firestore.DocumentRef{
				Parent: &firestore.CollectionRef{},
				Path:   "",
				ID:     "",
			}, nil
		},
		UpdateFn: func(ctx context.Context, command *fb.UpdateCommand) error {
			return nil
		},
		DeleteFn: func(ctx context.Context, command *fb.DeleteCommand) error {
			return nil
		},
		GetFn: func(ctx context.Context, query *fb.GetSingleQuery) (*firestore.DocumentSnapshot, error) {
			return &firestore.DocumentSnapshot{
				Ref:        &firestore.DocumentRef{},
				CreateTime: time.Time{},
				UpdateTime: time.Time{},
				ReadTime:   time.Time{},
			}, nil
		},

		//Firebase
		GetUserByPhoneNumberFn: func(ctx context.Context, phone string) (*auth.UserRecord, error) {
			return &auth.UserRecord{
				UserInfo: &auth.UserInfo{
					DisplayName: "",
					Email:       "",
					PhoneNumber: "",
					PhotoURL:    "",
					UID:         "",
				},
				CustomClaims: map[string]interface{}{},
				Disabled:     false,
			}, nil
		},
		CreateUserFn: func(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
			return &auth.UserRecord{
				UserInfo: &auth.UserInfo{
					DisplayName: "",
					Email:       "",
					PhoneNumber: "",
					PhotoURL:    "",
					UID:         "",
				},
				CustomClaims: map[string]interface{}{},
				Disabled:     false,
			}, nil
		},
		DeleteUserFn: func(ctx context.Context, uid string) error {
			return nil
		},
		GetUserProfileByIDFn: func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
			return &profileutils.UserProfile{
				ID: "",
			}, nil
		},
		MockSaveEmailOTPFn: func(ctx context.Context, email string, optIn bool) error {
			return nil
		},
		MockStageStartEpisodeByBreakGlassFn: func(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) error {
			return nil
		},
	}
}

// Collection ...
func (f *FakeFirebase) Collection(path string) *firestore.CollectionRef {
	return f.CollectionFn(path)
}

// GetAll retrieve a value from the store
func (f *FakeFirebase) GetAll(ctx context.Context, getQuery *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
	return f.GetAllFn(ctx, getQuery)
}

// Create persists data to a firestore collection
func (f *FakeFirebase) Create(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
	return f.CreateFn(ctx, command)
}

// Update updates data to a firestore collection
func (f *FakeFirebase) Update(ctx context.Context, command *fb.UpdateCommand) error {
	return f.UpdateFn(ctx, command)
}

// Delete deletes data to a firestore collection
func (f *FakeFirebase) Delete(ctx context.Context, command *fb.DeleteCommand) error {
	return f.DeleteFn(ctx, command)
}

// Get retrieves data to a firestore collection
func (f *FakeFirebase) Get(ctx context.Context, query *fb.GetSingleQuery) (*firestore.DocumentSnapshot, error) {
	return f.GetFn(ctx, query)
}

// GetUserByPhoneNumber ...
func (f *FakeFirebase) GetUserByPhoneNumber(ctx context.Context, phone string) (*auth.UserRecord, error) {
	return f.GetUserByPhoneNumberFn(ctx, phone)
}

// CreateUser ...
func (f *FakeFirebase) CreateUser(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
	return f.CreateUserFn(ctx, user)
}

// DeleteUser ...
func (f *FakeFirebase) DeleteUser(ctx context.Context, uid string) error {
	return f.DeleteUserFn(ctx, uid)
}

// GetUserProfileByID ...
func (f *FakeFirebase) GetUserProfileByID(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
	return f.GetUserProfileByIDFn(ctx, id, suspended)
}

// SaveEmailOTP ...
func (f *FakeFirebase) SaveEmailOTP(ctx context.Context, email string, optIn bool) error {
	return f.MockSaveEmailOTPFn(ctx, email, optIn)
}

// StageStartEpisodeByBreakGlass ...
func (f *FakeFirebase) StageStartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) error {
	return f.MockStageStartEpisodeByBreakGlassFn(ctx, input)
}
