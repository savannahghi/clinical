package infrastructure

import "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhir"

// Service is a clinical service
type Infrastructure struct {
	fhirDatastore *fhir.Service
}

// NewService initializes a new clinical service
func NewInfrastructureInteractor() *Infrastructure {

	clinicalRepository := fhir.NewService()
	return &Infrastructure{
		fhirDatastore: clinicalRepository,
	}
}
