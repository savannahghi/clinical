package usecases

import (
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/ocl"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	infra infrastructure.Infrastructure
	clinicalUsecase.UseCasesClinical
	ocl.UseCasesOCL
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Interactor {
	clinical := clinicalUsecase.NewUseCasesClinicalImpl(infrastructure)
	ocl := ocl.NewUseCasesOCLImpl(infrastructure)

	impl := Interactor{
		infrastructure,
		clinical,
		ocl,
	}

	return impl
}
