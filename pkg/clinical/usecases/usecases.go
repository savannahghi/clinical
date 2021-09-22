package usecases

// Interactor is an implementation of the usecases interface
type Interactor struct {
	ClinicalUseCase
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor() *Interactor {

	clinical := NewClinicalUseCaseImpl()

	impl := &Interactor{
		clinical,
	}

	return impl
}
