package cloudhealth

import (
	"strings"
	"testing"
)

func TestService_CreateDataset(t *testing.T) {
	s := NewService()
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
			got, err := s.CreateDataset()
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
