package openconceptlab

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing")
	m.Run()
}

func TestNewService(t *testing.T) {
	service := NewService()

	tests := []struct {
		name string
		want *Service
	}{
		{
			name: "good_case",
			want: service,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			} else {
				service.enforcePreconditions()
			}
		})
	}
}

func TestService_GetConcept(t *testing.T) {
	service := NewService()
	type args struct {
		ctx                    context.Context
		org                    string
		source                 string
		concept                string
		includeMappings        bool
		includeInverseMappings bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid_who_icd_concept",
			args: args{
				ctx:                    context.Background(),
				org:                    "CIEL",
				source:                 "CIEL",
				concept:                "106",
				includeMappings:        true,
				includeInverseMappings: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.GetConcept(
				tt.args.ctx, tt.args.org, tt.args.source, tt.args.concept,
				tt.args.includeMappings, tt.args.includeInverseMappings)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetConcept() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Contains(t, got, "display_name")
			}
		})
	}
}

func TestService_ListConcepts(t *testing.T) {
	type args struct {
		ctx                    context.Context
		org                    string
		source                 string
		verbose                bool
		q                      *string
		sortAsc                *string
		sortDesc               *string
		conceptClass           *string
		dataType               *string
		locale                 *string
		includeRetired         *bool
		includeMappings        *bool
		includeInverseMappings *bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx:     context.Background(),
				org:     "CIEL",
				source:  "CIEL",
				verbose: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.ListConcepts(tt.args.ctx, tt.args.org, tt.args.source, tt.args.verbose, tt.args.q, tt.args.sortAsc, tt.args.sortDesc, tt.args.conceptClass, tt.args.dataType, tt.args.locale, tt.args.includeRetired, tt.args.includeMappings, tt.args.includeInverseMappings)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListConcepts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}
