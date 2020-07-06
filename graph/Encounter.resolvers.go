package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIREncounter(ctx context.Context, input clinical.FHIREncounterInput) (*clinical.FHIREncounterRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIREncounter(ctx, input)
}

func (r *mutationResolver) UpdateFHIREncounter(ctx context.Context, input clinical.FHIREncounterInput) (*clinical.FHIREncounterRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIREncounter(ctx, input)
}

func (r *mutationResolver) DeleteFHIREncounter(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIREncounter(ctx, id)
}

func (r *queryResolver) GetFHIREncounter(ctx context.Context, id string) (*clinical.FHIREncounterRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIREncounter(ctx, id)
}

func (r *queryResolver) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*clinical.FHIREncounterRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIREncounter(ctx, params)
}
