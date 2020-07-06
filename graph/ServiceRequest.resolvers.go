package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIRServiceRequest(ctx context.Context, input clinical.FHIRServiceRequestInput) (*clinical.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRServiceRequest(ctx, input)
}

func (r *mutationResolver) UpdateFHIRServiceRequest(ctx context.Context, input clinical.FHIRServiceRequestInput) (*clinical.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRServiceRequest(ctx, input)
}

func (r *mutationResolver) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRServiceRequest(ctx, id)
}

func (r *queryResolver) GetFHIRServiceRequest(ctx context.Context, id string) (*clinical.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRServiceRequest(ctx, id)
}

func (r *queryResolver) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*clinical.FHIRServiceRequestRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRServiceRequest(ctx, params)
}
