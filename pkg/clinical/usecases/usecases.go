package usecases

import (
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/ocl"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	ClinicalUseCase
	FHIRUseCase
	ocl.UseCases
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Interactor {
	fhir := NewFHIRUseCaseImpl(infrastructure)
	clinical := NewClinicalUseCaseImpl(infrastructure, fhir)
	ocl := ocl.NewUseCasesImpl(infrastructure)

	impl := Interactor{
		clinical,
		fhir,
		ocl,
	}

	return impl
}
