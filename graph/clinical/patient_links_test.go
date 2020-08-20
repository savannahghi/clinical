package clinical

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("CLOUD_HEALTH_PUBSUB_TOPIC", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_DATASET_ID", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_FHIRSTORE_ID", "healthcloud-bewell-fhir-staging")
	m.Run()
}

func TestGeneratePatientLink(t *testing.T) {
	type args struct {
		patientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "good case",
			args:    args{patientID: "patient1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePatientLink(tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePatientLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.PatientID != tt.args.patientID {
				t.Errorf("GeneratePatientLink() = %v, want %v", got.PatientID, tt.args.patientID)
			}
		})
	}
}

func TestGetPatientID(t *testing.T) {
	testID := uuid.New().String()
	patientLink, err := GeneratePatientLink(testID)
	if err != nil {
		t.Errorf("unable to shorten testID")
	}

	type args struct {
		ctx      context.Context
		opaqueID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				opaqueID: patientLink.OpaqueID,
				ctx:      context.Background(),
			},
			want: patientLink.PatientID,
		},
		{
			name: "bad case invalid opaque id",
			args: args{
				opaqueID: "invalidid",
				ctx:      context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPatientID(tt.args.ctx, tt.args.opaqueID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPatientID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPatientID () = %v, want %v", got, tt.want)
			}
		})
	}
}
