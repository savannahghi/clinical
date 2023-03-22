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
	RegisterTenant(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error)
	RegisterFacility(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error)

	CreateEpisodeOfCare(ctx context.Context, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error)

	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	PatientHealthTimeline(ctx context.Context, input dto.HealthTimelineInput) (*dto.HealthTimeline, error)
	GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error)

	CreatePubsubPatient(ctx context.Context, payload dto.CreatePatientPubSubMessage) error
	CreatePubsubOrganization(ctx context.Context, payload dto.CreateFacilityPubSubMessage) error
	CreatePubsubVitals(ctx context.Context, data dto.CreateVitalSignPubSubMessage) error
	CreatePubsubAllergyIntolerance(ctx context.Context, data dto.CreatePatientAllergyPubSubMessage) error
	CreatePubsubTestResult(ctx context.Context, data dto.CreatePatientTestResultPubSubMessage) error
	CreatePubsubMedicationStatement(ctx context.Context, data dto.CreateMedicationPubSubMessage) error

	StartEncounter(ctx context.Context, episodeID string) (string, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	ListPatientEncounters(ctx context.Context, patientID string) ([]*dto.Encounter, error)

	RecordTemperature(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordHeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordRespiratoryRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordPulseRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordObservation(ctx context.Context, input dto.ObservationInput, vitalSignConceptID string) (*dto.Observation, error)
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
