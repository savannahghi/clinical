package repository

import (
	"context"
	"net/http"
	"net/url"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FHIR represents all the FHIR logic
type FHIR interface {
	CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error)
	UpdateFHIREpisodeOfCare(ctx context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error)
	SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error)
	CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error)
	FindOrganizationByID(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error)
	SearchEpisodesByParam(ctx context.Context, searchParams url.Values) ([]*domain.FHIREpisodeOfCare, error)
	HasOpenEpisode(
		ctx context.Context,
		patient domain.FHIRPatient,
	) (bool, error)
	OpenEpisodes(
		ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error)
	CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error)
	GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error)
	Encounters(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum) ([]*domain.FHIREncounter, error)
	SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*domain.FHIREpisodeOfCareRelayConnection, error)
	StartEncounter(ctx context.Context, episodeID string) (string, error)
	SearchEpisodeEncounter(
		ctx context.Context,
		episodeReference string,
	) (*domain.FHIREncounterRelayConnection, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	EndEpisode(ctx context.Context, episodeID string) (bool, error)
	GetActiveEpisode(ctx context.Context, episodeID string) (*domain.FHIREpisodeOfCare, error)
	SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error)
	CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error)
	SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error)
	CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error)
	CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	DeleteFHIRComposition(ctx context.Context, id string) (bool, error)
	UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error)
	SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error)
	SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error)
	CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error)
	SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error)
	CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error)
	DeleteFHIRObservation(ctx context.Context, id string) (bool, error)
	GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error)
	DeleteFHIRPatient(ctx context.Context, id string) (bool, error)
	DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error)
	DeleteFHIRResourceType(results []map[string]string) error
	CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)
	CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error)
	SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error)
	CreateFHIRPatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error)
	PatchFHIRPatient(ctx context.Context, id string, params []map[string]interface{}) (*domain.FHIRPatient, error)
	SearchFHIRPatient(ctx context.Context, searchParams string) (*domain.PatientConnection, error)

	FHIRHeaders() (http.Header, error)
}
