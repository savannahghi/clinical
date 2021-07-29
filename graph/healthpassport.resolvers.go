package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/clinical/graph/clinical"
	"github.com/savannahghi/clinical/graph/generated"
)

func (r *mutationResolver) StartEpisodeByOtp(ctx context.Context, input clinical.OTPEpisodeCreationInput) (*clinical.EpisodeOfCarePayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.StartEpisodeByOtp(ctx, input)
}

func (r *mutationResolver) StartEpisodeByBreakGlass(ctx context.Context, input clinical.BreakGlassEpisodeCreationInput) (*clinical.EpisodeOfCarePayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.StartEpisodeByBreakGlass(ctx, input)
}

func (r *mutationResolver) UpgradeEpisode(ctx context.Context, input clinical.OTPEpisodeUpgradeInput) (*clinical.EpisodeOfCarePayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.UpgradeEpisode(ctx, input)
}

func (r *mutationResolver) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.EndEpisode(ctx, episodeID)
}

func (r *mutationResolver) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.StartEncounter(ctx, episodeID)
}

func (r *mutationResolver) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.EndEncounter(ctx, encounterID)
}

func (r *mutationResolver) RegisterPatient(ctx context.Context, input clinical.SimplePatientRegistrationInput) (*clinical.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.RegisterPatient(ctx, input)
}

func (r *mutationResolver) UpdatePatient(ctx context.Context, input clinical.SimplePatientRegistrationInput) (*clinical.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.UpdatePatient(ctx, input)
}

func (r *mutationResolver) AddNextOfKin(ctx context.Context, input clinical.SimpleNextOfKinInput) (*clinical.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.AddNextOfKin(ctx, input)
}

func (r *mutationResolver) AddNhif(ctx context.Context, input *clinical.SimpleNHIFInput) (*clinical.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.AddNhif(ctx, input)
}

func (r *mutationResolver) CreateUpdatePatientExtraInformation(ctx context.Context, input clinical.PatientExtraInformationInput) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.CreateUpdatePatientExtraInformation(ctx, input)
}

func (r *mutationResolver) CreateFHIRMedicationRequest(ctx context.Context, input clinical.FHIRMedicationRequestInput) (*clinical.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIRMedicationRequest(ctx, input)
}

func (r *mutationResolver) UpdateFHIRMedicationRequest(ctx context.Context, input clinical.FHIRMedicationRequestInput) (*clinical.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIRMedicationRequest(ctx, input)
}

func (r *mutationResolver) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIRMedicationRequest(ctx, id)
}

func (r *mutationResolver) CreateFHIRAllergyIntolerance(ctx context.Context, input clinical.FHIRAllergyIntoleranceInput) (*clinical.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIRAllergyIntolerance(ctx, input)
}

func (r *mutationResolver) UpdateFHIRAllergyIntolerance(ctx context.Context, input clinical.FHIRAllergyIntoleranceInput) (*clinical.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIRAllergyIntolerance(ctx, input)
}

func (r *mutationResolver) CreateFHIRCondition(ctx context.Context, input clinical.FHIRConditionInput) (*clinical.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIRCondition(ctx, input)
}

func (r *mutationResolver) UpdateFHIRCondition(ctx context.Context, input clinical.FHIRConditionInput) (*clinical.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIRCondition(ctx, input)
}

func (r *mutationResolver) CreateFHIRServiceRequest(ctx context.Context, input clinical.FHIRServiceRequestInput) (*clinical.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIRServiceRequest(ctx, input)
}

func (r *mutationResolver) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIRServiceRequest(ctx, id)
}

func (r *mutationResolver) CreateFHIRObservation(ctx context.Context, input clinical.FHIRObservationInput) (*clinical.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIRObservation(ctx, input)
}

func (r *mutationResolver) CreateFHIRComposition(ctx context.Context, input clinical.FHIRCompositionInput) (*clinical.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIRComposition(ctx, input)
}

func (r *mutationResolver) UpdateFHIRComposition(ctx context.Context, input clinical.FHIRCompositionInput) (*clinical.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIRComposition(ctx, input)
}

func (r *mutationResolver) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIRComposition(ctx, id)
}

func (r *mutationResolver) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIRPatient(ctx, id)
}

func (r *mutationResolver) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIRObservation(ctx, id)
}

func (r *queryResolver) FindPatientsByMsisdn(ctx context.Context, msisdn string) (*clinical.PatientConnection, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.FindPatientsByMSISDN(ctx, msisdn)
}

func (r *queryResolver) FindPatients(ctx context.Context, search string) (*clinical.PatientConnection, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.PatientSearch(ctx, search)
}

func (r *queryResolver) GetPatient(ctx context.Context, id string) (*clinical.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.FindPatientByID(ctx, id)
}

func (r *queryResolver) OpenEpisodes(ctx context.Context, patientReference string) ([]*clinical.FHIREpisodeOfCare, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.OpenEpisodes(ctx, patientReference)
}

func (r *queryResolver) OpenOrganizationEpisodes(ctx context.Context, providerSladeCode string) ([]*clinical.FHIREpisodeOfCare, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.OpenOrganizationEpisodes(ctx, providerSladeCode)
}

func (r *queryResolver) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.ProblemSummary(ctx, patientID)
}

func (r *queryResolver) VisitSummary(ctx context.Context, encounterID string) (map[string]interface{}, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.VisitSummary(ctx, encounterID, clinical.MaxClinicalRecordPageSize)
}

func (r *queryResolver) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.PatientTimelineWithCount(ctx, episodeID, count)
}

func (r *queryResolver) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*clinical.FHIREncounterRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIREncounter(ctx, params)
}

func (r *queryResolver) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*clinical.FHIRConditionRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIRCondition(ctx, params)
}

func (r *queryResolver) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*clinical.FHIRAllergyIntoleranceRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIRAllergyIntolerance(ctx, params)
}

func (r *queryResolver) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*clinical.FHIRObservationRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIRObservation(ctx, params)
}

func (r *queryResolver) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*clinical.FHIRMedicationRequestRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIRMedicationRequest(ctx, params)
}

func (r *queryResolver) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*clinical.FHIRServiceRequestRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIRServiceRequest(ctx, params)
}

func (r *queryResolver) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*clinical.FHIRCompositionRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIRComposition(ctx, params)
}

func (r *queryResolver) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.AllergySummary(ctx, patientID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
