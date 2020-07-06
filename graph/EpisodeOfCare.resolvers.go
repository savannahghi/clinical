package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIREpisodeOfCare(ctx context.Context, input clinical.FHIREpisodeOfCareInput) (*clinical.FHIREpisodeOfCareRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIREpisodeOfCare(ctx, input)
}

func (r *mutationResolver) UpdateFHIREpisodeOfCare(ctx context.Context, input clinical.FHIREpisodeOfCareInput) (*clinical.FHIREpisodeOfCareRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIREpisodeOfCare(ctx, input)
}

func (r *mutationResolver) DeleteFHIREpisodeOfCare(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIREpisodeOfCare(ctx, id)
}

func (r *queryResolver) GetFHIREpisodeOfCare(ctx context.Context, id string) (*clinical.FHIREpisodeOfCareRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIREpisodeOfCare(ctx, id)
}

func (r *queryResolver) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*clinical.FHIREpisodeOfCareRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIREpisodeOfCare(ctx, params)
}
