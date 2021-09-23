package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/domain/model"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph/generated"
)

func (r *entityResolver) FindDummyByID(ctx context.Context, id string) (*model.Dummy, error) {
	panic(fmt.Errorf("not implemented"))
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
