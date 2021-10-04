package interactor

import (
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/ocl"
)

// Usecases is an interface that combines of all usescases
type Usecases interface {
	usecases.ClinicalUseCase
	usecases.FHIRUseCase
	ocl.UseCases
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	usecases.ClinicalUseCase
	usecases.FHIRUseCase
	ocl.UseCases
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Usecases {
	fhir := usecases.NewFHIRUseCaseImpl(infrastructure)
	clinical := usecases.NewClinicalUseCaseImpl(infrastructure, fhir)
	ocl := ocl.NewUseCasesImpl(infrastructure)

	impl := &Interactor{
		clinical,
		fhir,
		ocl,
	}

	return impl
}
