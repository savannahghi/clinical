package fhir

import (
	"context"
	"fmt"
	"log"

	"github.com/savannahghi/serverutils"
	"google.golang.org/api/healthcare/v1"
)

// constants used to configure the Google Cloud Healthcare API
const (
	DatasetLocation = "europe-west4"
)

// CreateDataset creates a  cloud healthcare dataset and returns it's name
func (s Service) CreateDataset() (*healthcare.Operation, error) {
	s.checkPreconditions()
	datasetsService := s.healthcareService.Projects.Locations.Datasets
	parent := fmt.Sprintf("projects/%s/locations/%s", s.projectID, s.location)
	resp, err := datasetsService.Create(parent, &healthcare.Dataset{}).DatasetId(s.datasetID).Do()
	if err != nil {
		return nil, fmt.Errorf("create Data Set: %v", err)
	}
	return resp, nil
}

// GetDataset gets a cloud healthcare dataset.
func (s Service) GetDataset() (*healthcare.Dataset, error) {
	s.checkPreconditions()
	datasetsService := s.healthcareService.Projects.Locations.Datasets
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", s.projectID, s.location, s.datasetID)
	resp, err := datasetsService.Get(name).Do()
	if err != nil {
		return nil, fmt.Errorf("get Data Set: %v", err)
	}
	return resp, nil
}

// Service is an implementation of the healthcare fhir repository
// It is implementation agnostic i.e logic should be handled using
// the preferred datastore/database
type Service struct {
	healthcareService                           *healthcare.Service
	projectID, location, datasetID, fhirStoreID string
}

// NewService initializes a Google Cloud Healthcare API service
func NewService() *Service {
	projectID := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	_ = serverutils.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
	datasetID := serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
	fhirStoreID := serverutils.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")

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

func (s Service) checkPreconditions() {
	if s.healthcareService == nil {
		log.Panicf("cloudhealth.Service *healthcare.Service is nil")
	}
}
