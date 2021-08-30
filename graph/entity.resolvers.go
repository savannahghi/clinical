package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/clinical/graph/clinical"
	"github.com/savannahghi/clinical/graph/generated"
)

func (r *entityResolver) FindDummyByID(ctx context.Context, id string) (*clinical.Dummy, error) {
	return nil, nil
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
