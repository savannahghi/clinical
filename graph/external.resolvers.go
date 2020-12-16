package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *dummyResolver) ID(ctx context.Context, obj *clinical.Dummy) (*string, error) {
	if obj != nil && obj.ID != "" {
		return &obj.ID, nil
	}
	randomID := ksuid.New().String()
	return &randomID, nil
}

type dummyResolver struct{ *Resolver }
