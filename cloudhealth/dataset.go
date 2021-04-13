package cloudhealth

import (
	"fmt"

	healthcare "google.golang.org/api/healthcare/v1"
)

// CreateDataset creates a dataset and returns it's name
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

// GetDataset gets a dataset.
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
