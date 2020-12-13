package cloudhealth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_GetFHIRStore(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "default case",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.GetFHIRStore()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFHIRStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
