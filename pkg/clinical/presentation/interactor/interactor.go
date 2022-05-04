package interactor

import (
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	fhir "github.com/savannahghi/clinical/pkg/clinical/usecases/fhir"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/ocl"
)

// Usecases is an interface that combines of all usescases
type Usecases interface {
	clinical.UseCasesClinical
	fhir.UseCasesFHIR
	ocl.UseCasesOCL
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	clinical.UseCasesClinical
	fhir.UseCasesFHIR
	infrastructure.Infrastructure
	ocl.UseCasesOCL
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Usecases {
	clinical := clinical.NewUseCasesClinicalImpl(infrastructure)
	fhir := fhir.NewUseCasesFHIRImpl(infrastructure)
	ocl := ocl.NewUseCasesOCLImpl(infrastructure)

	impl := &Interactor{
		clinical,
		fhir,
		infrastructure,
		ocl,
	}

	return impl
}
