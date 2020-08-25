// Package clinical implements a simplified GraphQL interface to a FHIR server
// that acts as a clinical data repository.
package clinical

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_StartEncounter(t *testing.T) {
	service := NewService()
	ctx := context.Background()
	ep := CreateFHIREpisodeOfCarePayload(t)

	type args struct {
		ctx       context.Context
		episodeID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test create successfully encounter start",
			args: args{ctx: ctx, episodeID: *ep.Resource.ID},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := service.StartEncounter(tt.args.ctx, tt.args.episodeID)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, actual)
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, actual)
			}
		})
	}
}
