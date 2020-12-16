// Package cloudhealth manages our interaction with the Google Cloud Healthcare
// base.
package cloudhealth

import (
	"context"
	"log"

	"gitlab.slade360emr.com/go/base"
	healthcare "google.golang.org/api/healthcare/v1"
)

// constants used to configure the Google Cloud Healthcare API
const (
	DatasetLocation = "europe-west4"
)

// NewService initializes a Google Cloud Healthcare API service
func NewService() *Service {
	projectID := base.MustGetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	_ = base.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
	datasetID := base.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
	fhirStoreID := base.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")

	ctx := context.Background()
	hsv, err := healthcare.NewService(ctx)
	if err != nil {
		log.Panicf("unable to initialize new Google Cloud Healthcare Service: %s", err)
	}
	service := &Service{
		healthcareService: hsv,
		projectID:         projectID,
		location:          DatasetLocation,
		datasetID:         datasetID,
		fhirStoreID:       fhirStoreID,
	}

	// ensure that the dataset exists
	_, err = service.GetDataset()
	if err != nil {
		_, err = service.CreateDataset()
		if err != nil {
			log.Printf(
				"Unable to get or create dataset with projectID %s, "+
					"location %s, datasetID %s; got error %s",
				service.projectID, service.location, service.datasetID, err,
			)
		}
	}

	// ensure that the FHIR store exists
	_, err = service.GetFHIRStore()
	if err != nil {
		_, err = service.CreateFHIRStore()
		if err != nil {
			log.Printf(
				"Unable to get or create FHIR store with projectID %s, "+
					"location %s, datasetID %s, fhirStoreID %s; got error %s",
				service.projectID, service.location, service.datasetID, service.fhirStoreID, err,
			)
		}
	}

	return service
}

// Service is a gateway to the Google Cloud Healthcare base.
// It holds the configuration needed to talk to the Cloud Healthcare base.
type Service struct {
	healthcareService                           *healthcare.Service
	projectID, location, datasetID, fhirStoreID string
}

func (s Service) checkPreconditions() {
	if s.healthcareService == nil {
		log.Panicf("cloudhealth.Service *healthcare.Service is nil")
	}
}
