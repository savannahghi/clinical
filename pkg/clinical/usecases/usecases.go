package usecases

import (
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	ClinicalUseCase
	FHIRUseCase
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Interactor {
	fhir := NewFHIRUseCaseImpl(infrastructure)
	clinical := NewClinicalUseCaseImpl(infrastructure, fhir)

	impl := Interactor{
		clinical,
		fhir,
	}

	return impl
}
