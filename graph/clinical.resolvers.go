package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
	"gitlab.slade360emr.com/go/clinical/graph/generated"
)

func (r *mutationResolver) DeleteFHIRAllergyIntolerance(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIRAllergyIntolerance(ctx, id)
}

func (r *mutationResolver) DeleteFHIRCondition(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIRCondition(ctx, id)
}

func (r *mutationResolver) CreateFHIREncounter(ctx context.Context, input clinical.FHIREncounterInput) (*clinical.FHIREncounterRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIREncounter(ctx, input)
}

func (r *mutationResolver) UpdateFHIREncounter(ctx context.Context, input clinical.FHIREncounterInput) (*clinical.FHIREncounterRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIREncounter(ctx, input)
}

func (r *mutationResolver) DeleteFHIREncounter(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIREncounter(ctx, id)
}

func (r *mutationResolver) CreateFHIREpisodeOfCare(ctx context.Context, input clinical.FHIREpisodeOfCareInput) (*clinical.FHIREpisodeOfCareRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIREpisodeOfCare(ctx, input)
}

func (r *mutationResolver) UpdateFHIREpisodeOfCare(ctx context.Context, input clinical.FHIREpisodeOfCareInput) (*clinical.FHIREpisodeOfCareRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIREpisodeOfCare(ctx, input)
}

func (r *mutationResolver) DeleteFHIREpisodeOfCare(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIREpisodeOfCare(ctx, id)
}

func (r *mutationResolver) UpdateFHIRObservation(ctx context.Context, input clinical.FHIRObservationInput) (*clinical.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIRObservation(ctx, input)
}

func (r *mutationResolver) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIRObservation(ctx, id)
}

func (r *mutationResolver) CreateFHIROrganization(ctx context.Context, input clinical.FHIROrganizationInput) (*clinical.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIROrganization(ctx, input)
}

func (r *mutationResolver) UpdateFHIROrganization(ctx context.Context, input clinical.FHIROrganizationInput) (*clinical.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIROrganization(ctx, input)
}

func (r *mutationResolver) DeleteFHIROrganization(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIROrganization(ctx, id)
}

func (r *mutationResolver) CreateFHIRPatient(ctx context.Context, input clinical.FHIRPatientInput) (*clinical.FHIRPatientRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.CreateFHIRPatient(ctx, input)
}

func (r *mutationResolver) UpdateFHIRPatient(ctx context.Context, input clinical.FHIRPatientInput) (*clinical.FHIRPatientRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIRPatient(ctx, input)
}

func (r *mutationResolver) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.DeleteFHIRPatient(ctx, id)
}

func (r *mutationResolver) UpdateFHIRServiceRequest(ctx context.Context, input clinical.FHIRServiceRequestInput) (*clinical.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.UpdateFHIRServiceRequest(ctx, input)
}

func (r *queryResolver) GetFHIRAllergyIntolerance(ctx context.Context, id string) (*clinical.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIRAllergyIntolerance(ctx, id)
}

func (r *queryResolver) GetFHIRComposition(ctx context.Context, id string) (*clinical.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIRComposition(ctx, id)
}

func (r *queryResolver) GetFHIRCondition(ctx context.Context, id string) (*clinical.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIRCondition(ctx, id)
}

func (r *queryResolver) GetFHIREncounter(ctx context.Context, id string) (*clinical.FHIREncounterRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIREncounter(ctx, id)
}

func (r *queryResolver) GetFHIREpisodeOfCare(ctx context.Context, id string) (*clinical.FHIREpisodeOfCareRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIREpisodeOfCare(ctx, id)
}

func (r *queryResolver) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*clinical.FHIREpisodeOfCareRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIREpisodeOfCare(ctx, params)
}

func (r *queryResolver) GetFHIRMedicationRequest(ctx context.Context, id string) (*clinical.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIRMedicationRequest(ctx, id)
}

func (r *queryResolver) GetFHIRObservation(ctx context.Context, id string) (*clinical.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIRObservation(ctx, id)
}

func (r *queryResolver) GetFHIROrganization(ctx context.Context, id string) (*clinical.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIROrganization(ctx, id)
}

func (r *queryResolver) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*clinical.FHIROrganizationRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIROrganization(ctx, params)
}

func (r *queryResolver) GetFHIRPatient(ctx context.Context, id string) (*clinical.FHIRPatientRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIRPatient(ctx, id)
}

func (r *queryResolver) SearchFHIRPatient(ctx context.Context, params map[string]interface{}) (*clinical.FHIRPatientRelayConnection, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.SearchFHIRPatient(ctx, params)
}

func (r *queryResolver) GetFHIRServiceRequest(ctx context.Context, id string) (*clinical.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.clinicalService.GetFHIRServiceRequest(ctx, id)
}

func (r *queryResolver) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.AllergySummary(ctx, patientID)
}

func (r *queryResolver) PatientTimeline(ctx context.Context, episodeID string) ([]map[string]interface{}, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.PatientTimeline(ctx, episodeID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
