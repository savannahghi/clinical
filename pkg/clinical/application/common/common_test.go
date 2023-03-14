package common

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	extMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/interserviceclient"
	"github.com/stretchr/testify/assert"
)

func TestNewInterServiceClient(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name   string
		args   args
		panics bool
	}{
		{
			name: "Happy case: initialize service",
			args: args{
				serviceName: gofakeit.Name(),
			},
			panics: false,
		},
		{
			name: "Sad case: failed to load deps from yaml",
			args: args{
				serviceName: gofakeit.Name(),
			},
			panics: true,
		},
		{
			name: "Sad case: failed load isc client",
			args: args{
				serviceName: gofakeit.Name(),
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := extMock.NewFakeBaseExtensionMock()

			if tt.panics {
				if tt.name == "Sad case: failed to load deps from yaml" {
					fakeExt.LoadDepsFromYAMLFn = func() (*interserviceclient.DepsConfig, error) {
						return nil, fmt.Errorf("an error occurred")
					}
				}
				if tt.name == "Sad case: failed load isc client" {
					fakeExt.SetupISCclientFn = func(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
						return nil, fmt.Errorf("an error occurred")
					}
				}
				assert.Panics(t, func() { NewInterServiceClient(tt.args.serviceName, fakeExt) })
			} else {
				assert.NotPanics(t, func() { NewInterServiceClient(tt.args.serviceName, fakeExt) })
			}
		})
	}
}
