package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIRObservation(ctx context.Context, input clinical.FHIRObservationInput) (*clinical.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRObservation(ctx, input)
}

func (r *mutationResolver) UpdateFHIRObservation(ctx context.Context, input clinical.FHIRObservationInput) (*clinical.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRObservation(ctx, input)
}

func (r *mutationResolver) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRObservation(ctx, id)
}

func (r *queryResolver) GetFHIRObservation(ctx context.Context, id string) (*clinical.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRObservation(ctx, id)
}

func (r *queryResolver) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*clinical.FHIRObservationRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRObservation(ctx, params)
}
