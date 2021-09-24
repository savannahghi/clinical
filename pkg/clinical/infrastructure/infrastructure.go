package infrastructure

// Infrastructure ...
type Infrastructure struct {
	FHIRRepo FHIRRepository
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor() Infrastructure {

	fhirRepository := NewFHIRService()
	return Infrastructure{
		fhirRepository,
	}
}
