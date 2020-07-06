package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
	"gitlab.slade360emr.com/go/clinical/graph/generated"
)

func (r *mutationResolver) CreateFHIRAllergyIntolerance(ctx context.Context, input clinical.FHIRAllergyIntoleranceInput) (*clinical.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRAllergyIntolerance(ctx, input)
}

func (r *mutationResolver) UpdateFHIRAllergyIntolerance(ctx context.Context, input clinical.FHIRAllergyIntoleranceInput) (*clinical.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRAllergyIntolerance(ctx, input)
}

func (r *mutationResolver) DeleteFHIRAllergyIntolerance(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRAllergyIntolerance(ctx, id)
}

func (r *queryResolver) GetFHIRAllergyIntolerance(ctx context.Context, id string) (*clinical.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRAllergyIntolerance(ctx, id)
}

func (r *queryResolver) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*clinical.FHIRAllergyIntoleranceRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRAllergyIntolerance(ctx, params)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
