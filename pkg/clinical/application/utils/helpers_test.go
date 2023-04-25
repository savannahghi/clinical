package utils

import (
	"fmt"
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

func TestReportErrorToSentry(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy case",
			args: args{
				err: fmt.Errorf("test error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReportErrorToSentry(tt.args.err)
		})
	}
}

func TestAddPubSubNamespace(t *testing.T) {
	type args struct {
		topicName   string
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy case: add pubsub namespace",
			args: args{
				topicName:   "test",
				serviceName: "service",
			},
			want: "service-test-staging-v2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddPubSubNamespace(tt.args.topicName, tt.args.serviceName); got != tt.want {
				t.Errorf("AddPubSubNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}
