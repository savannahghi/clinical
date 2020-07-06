package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIRPatient(ctx context.Context, input clinical.FHIRPatientInput) (*clinical.FHIRPatientRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRPatient(ctx, input)
}

func (r *mutationResolver) UpdateFHIRPatient(ctx context.Context, input clinical.FHIRPatientInput) (*clinical.FHIRPatientRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRPatient(ctx, input)
}

func (r *mutationResolver) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRPatient(ctx, id)
}

func (r *queryResolver) GetFHIRPatient(ctx context.Context, id string) (*clinical.FHIRPatientRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRPatient(ctx, id)
}

func (r *queryResolver) SearchFHIRPatient(ctx context.Context, params map[string]interface{}) (*clinical.FHIRPatientRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRPatient(ctx, params)
}
