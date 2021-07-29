package graph

import (
	"context"
	"testing"

	"github.com/savannahghi/clinical/graph/clinical"
	"github.com/segmentio/ksuid"
)

func Test_dummyResolver_ID(t *testing.T) {
	resolver := NewResolver()
	type args struct {
		ctx context.Context
		obj *clinical.Dummy
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil case",
			args: args{
				ctx: context.Background(),
				obj: nil,
			},
			wantErr: false,
		},
		{
			name: "non nil case",
			args: args{
				ctx: context.Background(),
				obj: &clinical.Dummy{
					ID: ksuid.New().String(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := dummyResolver{Resolver: resolver}
			got, err := r.ID(tt.args.ctx, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("dummyResolver.ID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("got nil ID from dummy resolver")
				return
			}
		})
	}
}
