package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIRCondition(ctx context.Context, input clinical.FHIRConditionInput) (*clinical.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRCondition(ctx, input)
}

func (r *mutationResolver) UpdateFHIRCondition(ctx context.Context, input clinical.FHIRConditionInput) (*clinical.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRCondition(ctx, input)
}

func (r *mutationResolver) DeleteFHIRCondition(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRCondition(ctx, id)
}

func (r *queryResolver) GetFHIRCondition(ctx context.Context, id string) (*clinical.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRCondition(ctx, id)
}

func (r *queryResolver) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*clinical.FHIRConditionRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRCondition(ctx, params)
}
