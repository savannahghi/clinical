package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// FakeAdvantage mocks the implementation of advantage API methods
type FakeAdvantage struct {
	MockSegmentPatient func(ctx context.Context, payload dto.SegmentationPayload) error
}

// NewFakeAdvantageMock is the advantage's mock methods constructor
func NewFakeAdvantageMock() *FakeAdvantage {
	return &FakeAdvantage{
		MockSegmentPatient: func(ctx context.Context, payload dto.SegmentationPayload) error {
			return nil
		},
	}
}

// SegmentPatient mocks the implementation of patient segmentation usecase
func (f *FakeAdvantage) SegmentPatient(ctx context.Context, payload dto.SegmentationPayload) error {
	return f.MockSegmentPatient(ctx, payload)
}
