package repository

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FHIR represents all the FHIR logic
type FHIR interface {
	FHIROrganization
	FHIRPatient
	FHIREpisodeOfCare
	FHIRObservation
	FHIRAllergyIntolerance
	FHIRServiceRequest
	FHIRMedicationRequest
	FHIRCondition
	FHIREncounter
	FHIRComposition
	FHIRMedicationStatement
	FHIRMedication
}

type FHIROrganization interface {
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	SearchFHIROrganization(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIROrganizationRelayConnection, error)
	GetFHIROrganization(ctx context.Context, id string) (*domain.FHIROrganizationRelayPayload, error)
}

type FHIRPatient interface {
	GetFHIRPatient(ctx context.Context, id string) (*domain.FHIRPatientRelayPayload, error)
	DeleteFHIRPatient(ctx context.Context, id string) (bool, error)
	CreateFHIRPatient(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error)
	PatchFHIRPatient(ctx context.Context, id string, params []map[string]interface{}) (*domain.FHIRPatient, error)
	SearchFHIRPatient(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error)
}
type FHIREpisodeOfCare interface {
	SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error)
	SearchEpisodesByParam(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error)
	GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error)
	CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error)
	UpdateFHIREpisodeOfCare(ctx context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error)
	HasOpenEpisode(ctx context.Context, patient domain.FHIRPatient, tenant dto.TenantIdentifiers, pagination dto.Pagination) (bool, error)
	OpenEpisodes(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error)
	EndEpisode(ctx context.Context, episodeID string) (bool, error)
	GetActiveEpisode(ctx context.Context, episodeID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCare, error)
}
type FHIRObservation interface {
	SearchFHIRObservation(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error)
	CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error)
	DeleteFHIRObservation(ctx context.Context, id string) (bool, error)
	SearchPatientObservations(ctx context.Context, patientReference, observationCode string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error)
}
type FHIRAllergyIntolerance interface {
	SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error)
	CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	GetFHIRAllergyIntolerance(ctx context.Context, id string) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	SearchPatientAllergyIntolerance(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRAllergy, error)
}
type FHIRServiceRequest interface {
	SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRServiceRequestRelayConnection, error)
	CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error)
	DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error)
}
type FHIRMedicationRequest interface {
	SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationRequestRelayConnection, error)
	CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error)
	DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error)
}
type FHIRCondition interface {
	SearchFHIRCondition(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRCondition, error)
	CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
	UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
}
type FHIREncounter interface {
	CreateFHIREncounter(ctx context.Context, input domain.FHIREncounterInput) (*domain.FHIREncounterRelayPayload, error)
	SearchPatientEncounters(ctx context.Context, patientReference string, status *domain.EncounterStatusEnum, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error)
	StartEncounter(ctx context.Context, episodeID string) (string, error)
	SearchEpisodeEncounter(ctx context.Context, episodeReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error)
	EndEncounter(ctx context.Context, encounterID string) (bool, error)
	GetFHIREncounter(ctx context.Context, id string) (*domain.FHIREncounterRelayPayload, error)
	SearchFHIREncounter(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIREncounter, error)
}
type FHIRComposition interface {
	SearchFHIRComposition(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRCompositionRelayConnection, error)
	CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	DeleteFHIRComposition(ctx context.Context, id string) (bool, error)
}
type FHIRMedicationStatement interface {
	CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)
	SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error)
}

type FHIRMedication interface {
	CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error)
}
