package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph/generated"
)

func (r *mutationResolver) StartEpisodeByOtp(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.StartEpisodeByOtp(ctx, input)
}

func (r *mutationResolver) StartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.StartEpisodeByBreakGlass(ctx, input)
}

func (r *mutationResolver) UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.UpgradeEpisode(ctx, input)
}

func (r *mutationResolver) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.EndEpisode(ctx, episodeID)
}

func (r *mutationResolver) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.usecases.StartEncounter(ctx, episodeID)
}

func (r *mutationResolver) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.EndEncounter(ctx, encounterID)
}

func (r *mutationResolver) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.usecases.RegisterPatient(ctx, input)
}

func (r *mutationResolver) UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.UpdatePatient(ctx, input)
}

func (r *mutationResolver) AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.AddNextOfKin(ctx, input)
}

func (r *mutationResolver) AddNhif(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.AddNHIF(ctx, input)
}

func (r *mutationResolver) CreateUpdatePatientExtraInformation(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.CreateUpdatePatientExtraInformation(ctx, input)
}

func (r *mutationResolver) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.CreateFHIRMedicationRequest(ctx, input)
}

func (r *mutationResolver) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.UpdateFHIRMedicationRequest(ctx, input)
}

func (r *mutationResolver) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.DeleteFHIRMedicationRequest(ctx, id)
}

func (r *mutationResolver) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.CreateFHIRAllergyIntolerance(ctx, input)
}

func (r *mutationResolver) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.UpdateFHIRAllergyIntolerance(ctx, input)
}

func (r *mutationResolver) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.CreateFHIRCondition(ctx, input)
}

func (r *mutationResolver) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.CreateFHIRCondition(ctx, input)
}

func (r *mutationResolver) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.CreateFHIRServiceRequest(ctx, input)
}

func (r *mutationResolver) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.DeleteFHIRServiceRequest(ctx, id)
}

func (r *mutationResolver) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.CreateFHIRObservation(ctx, input)
}

func (r *mutationResolver) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.CreateFHIRComposition(ctx, input)
}

func (r *mutationResolver) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.UpdateFHIRComposition(ctx, input)
}

func (r *mutationResolver) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.DeleteFHIRComposition(ctx, id)
}

func (r *mutationResolver) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.DeleteFHIRPatient(ctx, id)
}

func (r *mutationResolver) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.DeleteFHIRObservation(ctx, id)
}

func (r *queryResolver) PatientHealthTimeline(ctx context.Context, input domain.HealthTimelineInput) (*domain.HealthTimeline, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.PatientHealthTimeline(ctx, input)
}

func (r *queryResolver) FindPatientsByMsisdn(ctx context.Context, msisdn string) (*domain.PatientConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.FindPatientsByMSISDN(ctx, msisdn)
}

func (r *queryResolver) FindPatients(ctx context.Context, search string) (*domain.PatientConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.PatientSearch(ctx, search)
}

func (r *queryResolver) GetPatient(ctx context.Context, id string) (*domain.PatientPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.FindPatientByID(ctx, id)
}

func (r *queryResolver) OpenEpisodes(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.OpenEpisodes(ctx, patientReference)
}

func (r *queryResolver) OpenOrganizationEpisodes(ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.OpenOrganizationEpisodes(ctx, providerSladeCode)
}

func (r *queryResolver) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.ProblemSummary(ctx, patientID)
}

func (r *queryResolver) VisitSummary(ctx context.Context, encounterID string) (map[string]interface{}, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.VisitSummary(ctx, encounterID, common.MaxClinicalRecordPageSize)
}

func (r *queryResolver) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.PatientTimelineWithCount(ctx, episodeID, count)
}

func (r *queryResolver) PatientTimeline(ctx context.Context, patientID string, count int) ([]map[string]interface{}, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.PatientTimeline(ctx, patientID, count)
}

func (r *queryResolver) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.SearchFHIREncounter(ctx, params)
}

func (r *queryResolver) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.SearchFHIRCondition(ctx, params)
}

func (r *queryResolver) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.SearchFHIRAllergyIntolerance(ctx, params)
}

func (r *queryResolver) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.SearchFHIRObservation(ctx, params)
}

func (r *queryResolver) SearchFHIRMedicationStatement(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationStatementRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.SearchFHIRMedicationStatement(ctx, params)
}

func (r *queryResolver) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.SearchFHIRMedicationRequest(ctx, params)
}

func (r *queryResolver) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.SearchFHIRServiceRequest(ctx, params)
}

func (r *queryResolver) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.SearchFHIRComposition(ctx, params)
}

func (r *queryResolver) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.AllergySummary(ctx, patientID)
}

func (r *queryResolver) GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.GetMedicalData(ctx, patientID)
}

func (r *queryResolver) SearchOrganization(ctx context.Context, params map[string]interface{}) (*domain.FHIROrganizationRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.SearchFHIROrganization(ctx, params)
}

func (r *queryResolver) FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.FindOrganizationByID(ctx, organizationID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
