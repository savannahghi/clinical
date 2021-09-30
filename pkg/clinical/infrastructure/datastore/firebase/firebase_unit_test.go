package firebase_test

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	fb "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase"
	extMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase/mock"
	"github.com/savannahghi/firebasetools"
)

var fakeFireBaseClientExt extMock.FirebaseClientExtension
var fireBaseClientExt fb.FBClientExtension = &fakeFireBaseClientExt
var fakeFireStoreClientExt extMock.FirestoreClientExtension

func TestRepository_SaveEmailOTP_Unittest(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)
	type args struct {
		ctx   context.Context
		email string
		optIn bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid save email OTP",
			args: args{
				ctx:   ctx,
				email: firebasetools.TestUserEmail,
				optIn: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid save email OTP:  nil email address",
			args: args{
				ctx:   ctx,
				optIn: true,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid save email OTP" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, nil
				}
			}
			if tt.name == "invalid save email OTP:  nil email address" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, fmt.Errorf("cannot create firestore document")
				}
			}
			err := repo.SaveEmailOTP(tt.args.ctx, tt.args.email, tt.args.optIn)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}

		})
	}
}
