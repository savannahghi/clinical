package interactor

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

// UseCasesClinical represents all the patient business logic
type Clinical interface {
	FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error)

	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	PatientTimeline(ctx context.Context, patientID string, count int) ([]map[string]interface{}, error)
	PatientHealthTimeline(ctx context.Context, input domain.HealthTimelineInput) (*domain.HealthTimeline, error)
	GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error)
}

// Usecases is an interface that combines of all usescases
type Usecases interface {
	Clinical
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	Clinical
	infrastructure.Infrastructure
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) *Interactor {
	clinical := clinical.NewUseCasesClinicalImpl(infrastructure)

	impl := &Interactor{
		clinical,
		infrastructure,
	}

	return impl
}
