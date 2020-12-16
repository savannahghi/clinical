package cloudhealth

import (
	"strings"
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

func TestService_CreateFHIRStore(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
		wantErr bool
	}{
		{
			name:    "valid dataset create",
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.CreateFHIRStore()
			if !tt.wantErr && err != nil {
				if !strings.Contains(
					err.Error(),
					"googleapi: Error 409: already exists",
				) {
					t.Errorf("unexpected error: %w", err)
					return
				}
			}
			if !tt.wantNil && err == nil && got == nil {
				t.Errorf("got nil dataset")
				return
			}
		})
	}
}
