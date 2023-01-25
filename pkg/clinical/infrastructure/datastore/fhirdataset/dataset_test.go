package fhirdataset_test

import (
	"log"
	"os"
	"strings"
	"testing"

	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhirdataset"
	"github.com/tj/assert"
)

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

	// do clean up
	log.Printf("Running tests ...")
	code := m.Run()

	log.Printf("Tearing tests down ...")

	// restore environment varibles to original values
	os.Setenv(envOriginalValue, "ENVIRONMENT")
	os.Setenv("DEBUG", debugEnvValue)
	os.Setenv("ROOT_COLLECTION_SUFFIX", collectionEnvValue)

	os.Exit(code)
}

func TestRepository_CreateDataset(t *testing.T) {
	fhirRepo := dataset.NewFHIRRepository()
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
					t.Errorf("unexpected error: %v", err)
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
	fhirRepo := dataset.NewFHIRRepository()
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
					t.Errorf("unexpected error: %v", err)
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
	fhirRepo := dataset.NewFHIRRepository()
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
	fhirRepo := dataset.NewFHIRRepository()
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
					t.Errorf("unexpected error: %v", err)
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
