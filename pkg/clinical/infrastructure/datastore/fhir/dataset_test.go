package fhir_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhir"
	fb "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/interactor"
	"github.com/savannahghi/firebasetools"
	"github.com/tj/assert"
)

var testUsecase interactor.Usecases

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

	i, err := InitializeTestService(ctx)
	if err != nil {
		log.Printf("unable to initialize test service: %v", err)
		return
	}

	testUsecase = i

	purgeRecords := func() {
		collections := []string{
			r.GetEmailOptInCollectionName(),
		}
		for _, collection := range collections {
			ref := fsc.Collection(collection)
			firebasetools.DeleteCollection(ctx, fsc, ref, 10)
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

func InitializeTestService(ctx context.Context) (interactor.Usecases, error) {
	infrastructure := infrastructure.NewInfrastructureInteractor()

	usecases := interactor.NewUsecasesInteractor(infrastructure)
	return usecases, nil
}

func TestRepository_CreateDataset(t *testing.T) {
	fhirRepo := fhir.NewFHIRRepository()
	tests := []struct {
		name    string
		wantNil bool
		wantErr bool
	}{
		{
			name:    "valid dataset create",
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := fhirRepo
			got, err := fr.CreateDataset()
			if !tt.wantErr && err != nil {
				if !strings.Contains(
					err.Error(),
					"googleapi: Error 409: already exists",
				) {
					t.Errorf("unexpected error: %w", err)
					return
				}
			}
			if !tt.wantNil && err == nil && got == nil {
				t.Errorf("got nil dataset")
				return
			}
		})
	}
}

func TestRepository_GetDataset(t *testing.T) {
	fhirRepo := fhir.NewFHIRRepository()
	tests := []struct {
		name    string
		wantNil bool
		wantErr bool
	}{
		{
			name:    "default case",
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := fhirRepo
			got, err := fr.GetDataset()
			if !tt.wantErr && err != nil {
				if !strings.Contains(
					err.Error(),
					"googleapi: Error 404: does not exist exists",
				) {
					t.Errorf("unexpected error: %w", err)
					return
				}
				assert.NotNil(t, got)
			}
			if !tt.wantNil && err == nil && got == nil {
				t.Errorf("got nil dataset")
				return
			}
		})
	}
}

func TestRepository_GetFHIRStore(t *testing.T) {
	fhirRepo := fhir.NewFHIRRepository()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "default case",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fhirRepo.GetFHIRStore()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFHIRStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func TestRepository_CreateFHIRStore(t *testing.T) {
	fhirRepo := fhir.NewFHIRRepository()
	tests := []struct {
		name    string
		wantNil bool
		wantErr bool
	}{
		{
			name:    "valid fhir store create",
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fhirRepo.CreateFHIRStore()
			if !tt.wantErr && err != nil {
				if !strings.Contains(
					err.Error(),
					"googleapi: Error 409: already exists",
				) {
					t.Errorf("unexpected error: %w", err)
					return
				}
			}
			if !tt.wantNil && err == nil && got == nil {
				t.Errorf("got nil fhir store")
				return
			}
		})
	}
}
