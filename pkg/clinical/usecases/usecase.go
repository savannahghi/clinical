package usecase

import (
	"github.com/savannahghi/clinical/pkg/clinical/usecases/patient"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	patient.UseCasePatient
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor() *Interactor {

	patient := patient.NewUseCasePatient()

	impl := &Interactor{
		patient,
	}

	return impl
}
