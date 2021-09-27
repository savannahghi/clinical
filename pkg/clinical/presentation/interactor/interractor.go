package interactor

import (
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
)

// Usecases is an interface that combines of all usescases
type Usecases interface {
	usecases.ClinicalUseCase
	usecases.FHIRUseCase
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	usecases.ClinicalUseCase
	usecases.FHIRUseCase
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Usecases {
	fhir := usecases.NewFHIRUseCaseImpl(infrastructure)
	clinical := usecases.NewClinicalUseCaseImpl(infrastructure, fhir)

	impl := &Interactor{
		clinical,
		fhir,
	}

	return impl
}
