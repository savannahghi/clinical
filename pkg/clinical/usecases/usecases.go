package usecases

import (
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	fhirUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/fhir"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/ocl"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	infra infrastructure.Infrastructure
	clinicalUsecase.UseCasesClinical
	fhirUsecase.UseCasesFHIR
	ocl.UseCasesOCL
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Interactor {
	clinical := clinicalUsecase.NewUseCasesClinicalImpl(infrastructure)
	fhir := fhirUsecase.NewUseCasesFHIRImpl(infrastructure)
	ocl := ocl.NewUseCasesOCLImpl(infrastructure)

	impl := Interactor{
		infrastructure,
		clinical,
		fhir,
		ocl,
	}

	return impl
}
