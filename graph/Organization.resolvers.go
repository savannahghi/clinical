package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) CreateFHIROrganization(ctx context.Context, input clinical.FHIROrganizationInput) (*clinical.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIROrganization(ctx, input)
}

func (r *mutationResolver) UpdateFHIROrganization(ctx context.Context, input clinical.FHIROrganizationInput) (*clinical.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIROrganization(ctx, input)
}

func (r *mutationResolver) DeleteFHIROrganization(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIROrganization(ctx, id)
}

func (r *queryResolver) GetFHIROrganization(ctx context.Context, id string) (*clinical.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIROrganization(ctx, id)
}

func (r *queryResolver) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*clinical.FHIROrganizationRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIROrganization(ctx, params)
}

func (r *queryResolver) GetROrganization(ctx context.Context, providerSladeCode int) (*string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.GetOrganization(ctx, providerSladeCode)
}

func (r *queryResolver) CreateOrganization(ctx context.Context, providerSladeCode int) (*string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.CreateOrganization(ctx, providerSladeCode)
}

func (r *queryResolver) GetORCreateOrganization(ctx context.Context, providerSladeCode int) (*string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.GetORCreateOrganization(ctx, providerSladeCode)
}
