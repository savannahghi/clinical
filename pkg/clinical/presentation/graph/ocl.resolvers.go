package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *queryResolver) ListConcepts(ctx context.Context, org string, source string, verbose bool, q *string, sortAsc *string, sortDesc *string, conceptClass *string, dataType *string, locale *string, includeRetired *bool, includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.usecases.ListConcepts(
		ctx, org, source, verbose, q, sortAsc, sortDesc, conceptClass, dataType,
		locale, includeRetired, includeMappings, includeInverseMappings)
}
