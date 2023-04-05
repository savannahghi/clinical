package usecases

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

// Clinical represents all the patient business logic
type Clinical interface {
	RegisterTenant(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error)
	RegisterFacility(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error)

	CreateEpisodeOfCare(ctx context.Context, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error)
	GetEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error)
	EndEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error)

	CreateCondition(ctx context.Context, input dto.ConditionInput) (*dto.Condition, error)
	ListPatientConditions(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.ConditionConnection, error)

	PatientHealthTimeline(ctx context.Context, input dto.HealthTimelineInput) (*dto.HealthTimeline, error)
	GetMedicalData(ctx context.Context, patientID string) (*dto.MedicalData, error)

	CreatePubsubPatient(ctx context.Context, payload dto.PatientPubSubMessage) error
	CreatePubsubOrganization(ctx context.Context, payload dto.FacilityPubSubMessage) error
	CreatePubsubVitals(ctx context.Context, data dto.VitalSignPubSubMessage) error
	CreatePubsubAllergyIntolerance(ctx context.Context, data dto.PatientAllergyPubSubMessage) error
	CreatePubsubTestResult(ctx context.Context, data dto.PatientTestResultPubSubMessage) error
	CreatePubsubMedicationStatement(ctx context.Context, data dto.MedicationPubSubMessage) error
	CreatePubsubTenant(ctx context.Context, data dto.OrganizationInput) error

	CreatePatient(ctx context.Context, input dto.PatientInput) (*dto.Patient, error)

	StartEncounter(ctx context.Context, episodeID string) (string, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	ListPatientEncounters(ctx context.Context, patientID string) ([]*dto.Encounter, error)

	RecordTemperature(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordHeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordRespiratoryRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordPulseRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordBMI(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordObservation(ctx context.Context, input dto.ObservationInput, vitalSignConceptID string) (*dto.Observation, error)

	GetPatientObservations(ctx context.Context, patientID string, observationCode string) ([]*dto.Observation, error)
	GetPatientTemperatureEntries(ctx context.Context, patientID string) ([]*dto.Observation, error)
	GetPatientBloodPressureEntries(ctx context.Context, patientID string) ([]*dto.Observation, error)
	GetPatientHeightEntries(ctx context.Context, patientID string) ([]*dto.Observation, error)
	GetPatientRespiratoryRateEntries(ctx context.Context, patientID string) ([]*dto.Observation, error)
	GetPatientPulseRateEntries(ctx context.Context, patientID string) ([]*dto.Observation, error)
	GetPatientBMIEntries(ctx context.Context, patientID string) ([]*dto.Observation, error)
	GetPatientWeightEntries(ctx context.Context, patientID string) ([]*dto.Observation, error)

	CreateAllergyIntolerance(ctx context.Context, input dto.AllergyInput) (*dto.Allergy, error)
	SearchAllergy(ctx context.Context, name string) ([]*dto.Terminology, error)
	GetAllergyIntolerance(ctx context.Context, id string) (*dto.Allergy, error)
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
