package infrastructure

import "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhir"

// Infrastructure ...
type Infrastructure struct {
	fhirDatastore *fhir.Service
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor() *Infrastructure {

	clinicalRepository := fhir.NewService()
	return &Infrastructure{
		fhirDatastore: clinicalRepository,
	}
}
