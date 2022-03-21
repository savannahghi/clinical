package utils

import (
	"testing"

	"github.com/savannahghi/firebasetools"
)

func TestValidateEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid  email address",
			args: args{
				email: firebasetools.TestUserEmail,
			},
			wantErr: false,
		},
		{
			name: "invalid email address",
			args: args{
				email: "hey@notavalidemail",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateEmail(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
