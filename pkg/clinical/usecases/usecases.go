package usecases

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

// Clinical represents all the patient business logic
type Clinical interface {
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	PatientHealthTimeline(ctx context.Context, input domain.HealthTimelineInput) (*domain.HealthTimeline, error)
	GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error)

	CreatePubsubPatient(ctx context.Context, payload dto.CreatePatientPubSubMessage) error
	CreatePubsubOrganization(ctx context.Context, payload dto.CreateFacilityPubSubMessage) error
	CreatePubsubVitals(ctx context.Context, data dto.CreateVitalSignPubSubMessage) error
	CreatePubsubAllergyIntolerance(ctx context.Context, data dto.CreatePatientAllergyPubSubMessage) error
	CreatePubsubTestResult(ctx context.Context, data dto.CreatePatientTestResultPubSubMessage) error
	CreatePubsubMedicationStatement(ctx context.Context, data dto.CreateMedicationPubSubMessage) error
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	Clinical
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Interactor {
	clinical := clinicalUsecase.NewUseCasesClinicalImpl(infrastructure)

	impl := Interactor{
		clinical,
	}

	return impl
}
