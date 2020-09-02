package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIRAppointment(ctx context.Context, input clinical.FHIRAppointmentInput) (*clinical.FHIRAppointmentRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRAppointment(ctx, input)
}

func (r *mutationResolver) UpdateFHIRAppointment(ctx context.Context, input clinical.FHIRAppointmentInput) (*clinical.FHIRAppointmentRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRAppointment(ctx, input)
}

func (r *mutationResolver) DeleteFHIRAppointment(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRAppointment(ctx, id)
}

func (r *queryResolver) GetFHIRAppointment(ctx context.Context, id string) (*clinical.FHIRAppointmentRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRAppointment(ctx, id)
}

func (r *queryResolver) SearchFHIRAppointment(ctx context.Context, params map[string]interface{}) (*clinical.FHIRAppointmentRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRAppointment(ctx, params)
}

func (r *queryResolver) ListAppointments(ctx context.Context, providerSladeCode int) (*clinical.FHIRAppointmentRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.ListAppointments(ctx, providerSladeCode)
}
