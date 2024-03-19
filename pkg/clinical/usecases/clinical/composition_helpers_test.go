package clinical

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeAdvantageMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/advantage/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
)

func TestUseCasesClinicalImpl_mapCategoryEnumToCode(t *testing.T) {
	type args struct {
		category dto.CompositionCategory
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Happy Case: Assessment plan",
			args:    args{category: dto.AssessmentAndPlan},
			wantErr: false,
		},
		{
			name:    "Happy Case: History of presenting illness",
			args:    args{category: dto.HistoryOfPresentingIllness},
			wantErr: false,
		},
		{
			name:    "Happy Case: social history",
			args:    args{category: dto.SocialHistory},
			wantErr: false,
		},
		{
			name:    "Happy Case: Examination",
			args:    args{category: dto.Examination},
			wantErr: false,
		},
		{
			name:    "Happy Case: Plan of care",
			args:    args{category: dto.PlanOfCare},
			wantErr: false,
		},
		{
			name:    "Happy Case: referral note",
			args:    args{category: dto.ReferralNote},
			wantErr: false,
		},
		{
			name:    "Happy Case: family history",
			args:    args{category: dto.FamilyHistory},
			wantErr: false,
		},
		{
			name:    "Sad Case: unknown category",
			args:    args{category: "test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			c := NewUseCasesClinicalImpl(infra)

			_, err := c.mapCategoryEnumToCode(tt.args.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.mapCategoryEnumToCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_mapCompositionConcepts(t *testing.T) {
	type args struct {
		ctx                     context.Context
		compositionCategoryCode string
		conceptID               string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: map composition concepts",
			args: args{
				ctx:                     context.Background(),
				compositionCategoryCode: string(dto.ReferralNote),
				conceptID:               common.LOINCReferralNote,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to map composition concepts",
			args: args{
				ctx:                     context.Background(),
				compositionCategoryCode: string(dto.ReferralNote),
				conceptID:               common.LOINCReferralNote,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()

			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			c := NewUseCasesClinicalImpl(infra)

			if tt.name == "Sad case: unable to map composition concepts" {
				fakeOCL.MockGetConceptFn = func(ctx context.Context, org, source, concept string, includeMappings, includeInverseMappings bool) (*domain.Concept, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := c.mapCompositionConcepts(tt.args.ctx, tt.args.compositionCategoryCode, tt.args.conceptID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.mapCompositionConcepts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
