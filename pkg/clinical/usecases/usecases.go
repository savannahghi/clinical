package usecases

import (
	"context"
	"io"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/scalarutils"
)

// Clinical represents all the patient business logic
type Clinical interface {
	RegisterTenant(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error)
	RegisterFacility(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error)

	CreateEpisodeOfCare(ctx context.Context, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error)
	GetEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error)
	PatchEpisodeOfCare(ctx context.Context, id string, input dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error)
	EndEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error)

	CreateCondition(ctx context.Context, input dto.ConditionInput) (*dto.Condition, error)
	ListPatientConditions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ConditionConnection, error)

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
	PatchPatient(ctx context.Context, id string, input dto.PatientInput) (*dto.Patient, error)
	DeletePatient(ctx context.Context, id string) (bool, error)

	StartEncounter(ctx context.Context, episodeID string) (string, error)
	PatchEncounter(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	ListPatientEncounters(ctx context.Context, patientID string, pagination *dto.Pagination) (*dto.EncounterConnection, error)

	RecordTemperature(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordHeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordRespiratoryRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordPulseRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordBMI(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordViralLoad(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordMuac(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordOxygenSaturation(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordBloodSugar(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordLastMenstrualPeriod(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordDiastolicBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordColposcopy(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error)
	RecordObservation(ctx context.Context, input dto.ObservationInput, vitalSignConceptID string) (*dto.Observation, error)

	GetPatientObservations(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, observationCode string, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientTemperatureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientHeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientRespiratoryRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientPulseRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientBMIEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientWeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientMuacEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientOxygenSaturationEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientViralLoad(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientBloodSugarEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientLastMenstrualPeriodEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)
	GetPatientDiastolicBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination *dto.Pagination) (*dto.ObservationConnection, error)

	CreateAllergyIntolerance(ctx context.Context, input dto.AllergyInput) (*dto.Allergy, error)
	SearchAllergy(ctx context.Context, name string, pagination dto.Pagination) (*dto.TerminologyConnection, error)
	GetAllergyIntolerance(ctx context.Context, id string) (*dto.Allergy, error)
	ListPatientAllergies(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.AllergyConnection, error)

	UploadMedia(ctx context.Context, encounterID string, file io.Reader, contentType string) (*dto.Media, error)
	ListPatientMedia(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.MediaConnection, error)

	CreateComposition(ctx context.Context, input dto.CompositionInput) (*dto.Composition, error)
	AppendNoteToComposition(ctx context.Context, id string, input dto.PatchCompositionInput) (*dto.Composition, error)
	ListPatientCompositions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.CompositionConnection, error)

	PatchPatientObservations(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientHeight(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientWeight(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientBMI(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientTemperature(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientDiastolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientSystolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientRespiratoryRate(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientOxygenSaturation(ctx context.Context, id string, value string) (*dto.Observation, error)
	PatchPatientPulseRate(ctx context.Context, id string, value string) (*dto.Observation, error)

	// Questionnaire
	CreateQuestionnaire(ctx context.Context, questionnaireInput *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error)
	RecordConsent(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error)
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
