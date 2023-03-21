package interactor

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

// UseCasesClinical represents all the patient business logic
type Clinical interface {
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	PatientHealthTimeline(ctx context.Context, input dto.HealthTimelineInput) (*dto.HealthTimeline, error)
	GetMedicalData(ctx context.Context, patientID string) (*dto.MedicalData, error)
}

// Usecases is an interface that combines of all usescases
type Usecases interface {
	Clinical
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	Clinical
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) *Interactor {
	clinical := clinical.NewUseCasesClinicalImpl(infrastructure)

	impl := &Interactor{
		clinical,
	}

	return impl
}
