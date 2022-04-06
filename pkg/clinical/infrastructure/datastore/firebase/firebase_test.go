package firebase_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fb "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase"
	"github.com/savannahghi/firebasetools"
)

var testRepo *fb.Repository

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")
	envOriginalValue := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "staging")
	debugEnvValue := os.Getenv("DEBUG")
	os.Setenv("DEBUG", "true")
	os.Setenv("REPOSITORY", "firebase")
	collectionEnvValue := os.Getenv("ROOT_COLLECTION_SUFFIX")
	os.Setenv("CLOUD_HEALTH_PUBSUB_TOPIC", "pubtopic")
	os.Setenv("CLOUD_HEALTH_DATASET_ID", "datasetid")
	os.Setenv("CLOUD_HEALTH_FHIRSTORE_ID", "fhirid")

	// !NOTE!
	// Under no circumstances should you remove this env var when testing
	// You risk purging important collections, like our prod collections
	os.Setenv("ROOT_COLLECTION_SUFFIX", fmt.Sprintf("clinical_ci_%v", time.Now().Unix()))
	ctx := context.Background()
	r := fb.Repository{} // They are nil
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Printf("failed to initialize test FireStore client")
		return
	}
	if fbc == nil {
		log.Printf("failed to initialize test FireBase client")
		return
	}
	infra, err := InitializeTestFirebaseRepository()
	if err != nil {
		log.Printf("unable to initialize test repository: %v", err)
		return
	}

	testRepo = infra

	purgeRecords := func() {
		collections := []string{
			r.GetEmailOptInCollectionName(),
		}
		for _, collection := range collections {
			ref := fsc.Collection(collection)
			_ = firebasetools.DeleteCollection(ctx, fsc, ref, 10)
		}
	}

	// try clean up first
	purgeRecords()

	// do clean up
	log.Printf("Running tests ...")
	code := m.Run()

	log.Printf("Tearing tests down ...")
	purgeRecords()

	// restore environment varibles to original values
	os.Setenv(envOriginalValue, "ENVIRONMENT")
	os.Setenv("DEBUG", debugEnvValue)
	os.Setenv("ROOT_COLLECTION_SUFFIX", collectionEnvValue)

	os.Exit(code)
}

func InitializeTestInfrastructure(ctx context.Context) (infrastructure.Infrastructure, error) {

	return infrastructure.NewInfrastructureInteractor(), nil
}

func InitializeTestFirebaseRepository() (*fb.Repository, error) {
	ctx := context.Background()
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Printf("failed to initialize test Firestore client")
		return nil, fmt.Errorf("failed to initialize test Firestore client")
	}
	if fbc == nil {
		log.Printf("failed to initialize test Firebase client")
		return nil, fmt.Errorf("failed to initialize test Firebase client")
	}
	firestoreExtension := fb.NewFirestoreClientExtension(fsc)
	fr := fb.NewFirebaseRepository(firestoreExtension, fbc)
	return fr, nil
}
func InitializeTestFirebaseClient(ctx context.Context) (*firestore.Client, *auth.Client) {
	fc := firebasetools.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore: %s", err)
	}

	fbc, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up tests: %s", err)
	}
	return fsc, fbc
}

func TestRepository_SaveEmailOTP(t *testing.T) {
	fr := testRepo
	ctx := context.Background()
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
			err := fr.SaveEmailOTP(tt.args.ctx, tt.args.email, tt.args.optIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveEmailOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})

	}
}
