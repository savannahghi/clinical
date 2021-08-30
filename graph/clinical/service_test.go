package clinical_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/graph/clinical"
	"github.com/savannahghi/firebasetools"
)

func TestService_DeleteFHIRObservation(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	service := clinical.NewService()
	id := uuid.New().String()
	codeID := uuid.New().String()
	coding := []*clinical.FHIRCodingInput{
		{
			ID:      &codeID,
			Display: "This is a sample display",
			Code:    "123",
		},
	}
	var status clinical.ObservationStatusEnum = "registered"
	testObservationItemCode := clinical.FHIRCodeableConceptInput{
		ID:     &codeID,
		Coding: coding,
		Text:   "This is a test Text",
	}

	testObservationItem := clinical.FHIRObservationInput{
		ID:     &id,
		Status: &status,
		Code:   testObservationItemCode,
	}

	// Here we create the FHIRObservation so that we attempt to delete it in the test
	newFHIRObservation, err := service.CreateFHIRObservation(ctx, testObservationItem)
	if err != nil {
		t.Errorf("failed to create a FHIR observation :%v", err)
		return
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: delete a created FHIR observation",
			args: args{
				ctx: ctx,
				id:  *newFHIRObservation.Resource.ID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: delete a non-existent FHIR observation",
			args: args{
				ctx: ctx,
				id:  "",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.DeleteFHIRObservation(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.DeleteFHIRObservation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.DeleteFHIRObservation() = %v, want %v", got, tt.want)
			}
		})
	}
}
