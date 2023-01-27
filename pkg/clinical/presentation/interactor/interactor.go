package interactor

import (
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/ocl"
)

// Usecases is an interface that combines of all usescases
type Usecases interface {
	clinical.UseCasesClinical
	ocl.UseCasesOCL
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	clinical.UseCasesClinical
	infrastructure.Infrastructure
	ocl.UseCasesOCL
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Usecases {
	clinical := clinical.NewUseCasesClinicalImpl(infrastructure)
	ocl := ocl.NewUseCasesOCLImpl(infrastructure)

	impl := &Interactor{
		clinical,
		infrastructure,
		ocl,
	}

	return impl
}
