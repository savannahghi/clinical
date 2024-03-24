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
	FHIRMedia
	FHIRQuestionnaire
	FHIRConsent
	FHIRQuestionnaireResponse
	FHIRRiskAssessment
	FHIRDiagnosticReport
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
	PatchFHIRPatient(ctx context.Context, id string, input domain.FHIRPatientInput) (*domain.FHIRPatient, error)
	SearchFHIRPatient(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error)
	GetFHIRPatientEverything(ctx context.Context, id string, params map[string]interface{}) (*domain.PagedFHIRResource, error)
}
type FHIREpisodeOfCare interface {
	SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCareRelayConnection, error)
	SearchEpisodesByParam(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error)
	GetFHIREpisodeOfCare(ctx context.Context, id string) (*domain.FHIREpisodeOfCareRelayPayload, error)
	CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCareInput) (*domain.EpisodeOfCarePayload, error)
	PatchFHIREpisodeOfCare(ctx context.Context, id string, episode domain.FHIREpisodeOfCareInput) (*domain.FHIREpisodeOfCare, error)
	UpdateFHIREpisodeOfCare(ctx context.Context, fhirResourceID string, payload map[string]interface{}) (*domain.FHIREpisodeOfCare, error)
	HasOpenEpisode(ctx context.Context, patient domain.FHIRPatient, tenant dto.TenantIdentifiers, pagination dto.Pagination) (bool, error)
	OpenEpisodes(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) ([]*domain.FHIREpisodeOfCare, error)
	EndEpisode(ctx context.Context, episodeID string) (bool, error)
	GetActiveEpisode(ctx context.Context, episodeID string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIREpisodeOfCare, error)
}
type FHIRObservation interface {
	GetFHIRObservation(ctx context.Context, id string) (*domain.FHIRObservationRelayPayload, error)
	SearchFHIRObservation(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error)
	CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error)
	DeleteFHIRObservation(ctx context.Context, id string) (bool, error)
	SearchPatientObservations(ctx context.Context, searchParameters map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRObservations, error)
	PatchFHIRObservation(ctx context.Context, id string, input domain.FHIRObservationInput) (*domain.FHIRObservation, error)
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
	GetFHIRServiceRequest(ctx context.Context, id string) (*domain.FHIRServiceRequestRelayPayload, error)
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
	PatchFHIREncounter(ctx context.Context, encounterID string, input domain.FHIREncounterInput) (*domain.FHIREncounter, error)
	SearchFHIREncounterAllData(_ context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRResource, error)
}
type FHIRComposition interface {
	GetFHIRComposition(ctx context.Context, id string) (*domain.FHIRCompositionRelayPayload, error)
	SearchFHIRComposition(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRComposition, error)
	CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error)
	UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error)
	PatchFHIRComposition(ctx context.Context, id string, input domain.FHIRCompositionInput) (*domain.FHIRComposition, error)
	DeleteFHIRComposition(ctx context.Context, id string) (bool, error)
}
type FHIRMedicationStatement interface {
	CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)
	SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRMedicationStatementRelayConnection, error)
}

type FHIRMedication interface {
	CreateFHIRMedication(ctx context.Context, input domain.FHIRMedicationInput) (*domain.FHIRMedicationRelayPayload, error)
}

type FHIRMedia interface {
	CreateFHIRMedia(ctx context.Context, input domain.FHIRMedia) (*domain.FHIRMedia, error)
	SearchPatientMedia(ctx context.Context, patientReference string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRMedia, error)
}

type FHIRQuestionnaire interface {
	GetFHIRQuestionnaire(ctx context.Context, id string) (*domain.FHIRQuestionnaireRelayPayload, error)
	CreateFHIRQuestionnaire(ctx context.Context, input *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error)
	ListFHIRQuestionnaire(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PagedFHIRQuestionnaires, error)
}
type FHIRConsent interface {
	CreateFHIRConsent(ctx context.Context, input domain.FHIRConsent) (*domain.FHIRConsent, error)
}

type FHIRQuestionnaireResponse interface {
	CreateFHIRQuestionnaireResponse(_ context.Context, input *domain.FHIRQuestionnaireResponse) (*domain.FHIRQuestionnaireResponse, error)
	GetFHIRQuestionnaireResponse(ctx context.Context, id string) (*domain.FHIRQuestionnaireResponseRelayPayload, error)
}

type FHIRRiskAssessment interface {
	CreateFHIRRiskAssessment(_ context.Context, input *domain.FHIRRiskAssessmentInput) (*domain.FHIRRiskAssessmentRelayPayload, error)
	SearchFHIRRiskAssessment(ctx context.Context, params map[string]interface{}, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.FHIRRiskAssessmentRelayConnection, error)
}

type FHIRDiagnosticReport interface {
	CreateFHIRDiagnosticReport(_ context.Context, input *domain.FHIRDiagnosticReportInput) (*domain.FHIRDiagnosticReport, error)
}
