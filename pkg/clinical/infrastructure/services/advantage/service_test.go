package advantage_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/authutils"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/advantage"
	advantageMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/advantage/mock"
)

func TestServiceAdvantageImpl_PatientSegmentation(t *testing.T) {
	type args struct {
		ctx     context.Context
		payload dto.SegmentationPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: segment patients",
			args: args{
				ctx: context.Background(),
				payload: dto.SegmentationPayload{
					ClinicalID:   gofakeit.UUID(),
					SegmentLabel: dto.SegmentationCategoryHighRiskPositive,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to segment patients",
			args: args{
				ctx: context.Background(),
				payload: dto.SegmentationPayload{
					ClinicalID:   gofakeit.UUID(),
					SegmentLabel: dto.SegmentationCategoryHighRiskPositive,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeAuthUtils := advantageMock.NewAuthUtilsClientMock()
			s := advantage.NewServiceAdvantage(fakeAuthUtils)

			if tt.name == "Sad case: unable to segment patients" {
				fakeAuthUtils.MockAuthenticateFn = func() (*authutils.OAUTHResponse, error) {
					return nil, errors.New("unable to authenticate")
				}
			}

			if err := s.SegmentPatient(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("ServiceAdvantageImpl.SegmentPatient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
