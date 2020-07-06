package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIRComposition(ctx context.Context, input clinical.FHIRCompositionInput) (*clinical.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRComposition(ctx, input)
}

func (r *mutationResolver) UpdateFHIRComposition(ctx context.Context, input clinical.FHIRCompositionInput) (*clinical.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRComposition(ctx, input)
}

func (r *mutationResolver) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRComposition(ctx, id)
}

func (r *queryResolver) GetFHIRComposition(ctx context.Context, id string) (*clinical.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRComposition(ctx, id)
}

func (r *queryResolver) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*clinical.FHIRCompositionRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRComposition(ctx, params)
}
