package graph

import (
	"context"
	"testing"

	"github.com/savannahghi/clinical/graph/clinical"
	"github.com/savannahghi/clinical/graph/openconceptlab"
	"github.com/savannahghi/firebasetools"
	"github.com/stretchr/testify/assert"
)

func TestResolver_CheckDependencies(t *testing.T) {
	type fields struct {
		clinicalService *clinical.Service
		oclService      *openconceptlab.Service
	}
	tests := []struct {
		name      string
		fields    fields
		wantPanic bool
	}{
		{
			name: "nil clinical service",
			fields: fields{
				clinicalService: nil,
				oclService:      openconceptlab.NewService(),
			},
			wantPanic: true,
		},
		{
			name: "nil OCL service",
			fields: fields{
				clinicalService: clinical.NewService(),
				oclService:      nil,
			},
			wantPanic: true,
		},
		{
			name: "all dependencies present",
			fields: fields{
				clinicalService: clinical.NewService(),
				oclService:      openconceptlab.NewService(),
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resolver{
				clinicalService: tt.fields.clinicalService,
				oclService:      tt.fields.oclService,
			}
			if !tt.wantPanic {
				r.CheckDependencies()
				return
			}
			assert.Panics(t, func() {
				r.CheckDependencies()
			})
		})
	}
}

func TestResolver_CheckUserTokenInContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "should panic, no auth in context",
			args: args{
				ctx: context.Background(),
			},
			wantPanic: true,
		},
		{
			name: "should not panic, auth in context",
			args: args{
				ctx: firebasetools.GetAuthenticatedContext(t),
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewResolver()
			if !tt.wantPanic {
				token := r.CheckUserTokenInContext(tt.args.ctx)
				if token == nil {
					t.Errorf("unexpected nil token")
				}
				return
			}
			// panic
			assert.Panics(t, func() {
				_ = r.CheckUserTokenInContext(tt.args.ctx)
			})
		})
	}
}
