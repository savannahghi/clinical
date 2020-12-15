package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/clinical/graph/clinical"
	"gitlab.slade360emr.com/go/clinical/graph/generated"
)

func (r *dummyResolver) ID(ctx context.Context, obj *clinical.Dummy) (*string, error) {
	if obj != nil && obj.ID != "" {
		return &obj.ID, nil
	}
	randomID := ksuid.New().String()
	return &randomID, nil
}

// Dummy returns generated.DummyResolver implementation.
func (r *Resolver) Dummy() generated.DummyResolver { return &dummyResolver{r} }

type dummyResolver struct{ *Resolver }
