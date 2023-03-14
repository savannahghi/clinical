package helpers

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
)

func TestInitializeInterServiceClient(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name   string
		args   args
		panics bool
	}{

		{
			name: "Happy case",
			args: args{
				serviceName: gofakeit.Name(),
			},
			panics: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panics {
				assert.Panics(t, func() { _ = InitializeInterServiceClient(tt.args.serviceName) })
			} else {
				assert.NotPanics(t, func() { _ = InitializeInterServiceClient(tt.args.serviceName) })
			}
		})
	}
}
