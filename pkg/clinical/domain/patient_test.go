package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPatientLink_GetID(t *testing.T) {
	id := uuid.New().String()
	type fields struct {
		ID        string
		PatientID string
		OpaqueID  string
		Expires   time.Time
		Deleted   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy case",
			fields: fields{
				ID: id,
			},
			want: id,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := &PatientLink{
				ID:        tt.fields.ID,
				PatientID: tt.fields.PatientID,
				OpaqueID:  tt.fields.OpaqueID,
				Expires:   tt.fields.Expires,
				Deleted:   tt.fields.Deleted,
			}
			if got := pl.GetID(); got != tt.want {
				t.Errorf("PatientLink.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatientLink_SetID(t *testing.T) {
	id := uuid.New().String()
	type fields struct {
		ID        string
		PatientID string
		OpaqueID  string
		Expires   time.Time
		Deleted   bool
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Happy case",
			fields: fields{
				ID: id,
			},
			args: args{
				id: id,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := &PatientLink{
				ID:        tt.fields.ID,
				PatientID: tt.fields.PatientID,
				OpaqueID:  tt.fields.OpaqueID,
				Expires:   tt.fields.Expires,
				Deleted:   tt.fields.Deleted,
			}
			pl.SetID(tt.args.id)
		})
	}
}
