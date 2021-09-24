package usecases

import "github.com/savannahghi/clinical/pkg/clinical/infrastructure"

// Usecases is an interface that combines of all usescases
type Usecases interface {
	ClinicalUseCase
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	ClinicalUseCase
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) *Interactor {

	clinical := NewClinicalUseCaseImpl(infrastructure)

	impl := &Interactor{
		clinical,
	}

	return impl
}
