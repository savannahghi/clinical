package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIRMedicationRequest(ctx context.Context, input clinical.FHIRMedicationRequestInput) (*clinical.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRMedicationRequest(ctx, input)
}

func (r *mutationResolver) UpdateFHIRMedicationRequest(ctx context.Context, input clinical.FHIRMedicationRequestInput) (*clinical.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRMedicationRequest(ctx, input)
}

func (r *mutationResolver) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRMedicationRequest(ctx, id)
}

func (r *queryResolver) GetFHIRMedicationRequest(ctx context.Context, id string) (*clinical.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRMedicationRequest(ctx, id)
}

func (r *queryResolver) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*clinical.FHIRMedicationRequestRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRMedicationRequest(ctx, params)
}
