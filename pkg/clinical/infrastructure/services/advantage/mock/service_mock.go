package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// FakeAdvantage mocks the implementation of advantage API methods
type FakeAdvantage struct {
	MockSegmentPatient func(ctx context.Context, payload dto.SegmentationPayload) error
	MockSendSMSFn      func(ctx context.Context, workstationID string, payload dto.SMSPayload) error
}

// NewFakeAdvantageMock is the advantage's mock methods constructor
func NewFakeAdvantageMock() *FakeAdvantage {
	return &FakeAdvantage{
		MockSegmentPatient: func(ctx context.Context, payload dto.SegmentationPayload) error {
			return nil
		},
		MockSendSMSFn: func(ctx context.Context, workstationID string, payload dto.SMSPayload) error {
			return nil
		},
	}
}

// SegmentPatient mocks the implementation of patient segmentation usecase
func (f *FakeAdvantage) SegmentPatient(ctx context.Context, payload dto.SegmentationPayload) error {
	return f.MockSegmentPatient(ctx, payload)
}

// SendSMS mocks the implementation of SMS notification to patient
func (f *FakeAdvantage) SendSMS(ctx context.Context, workstationID string, payload dto.SMSPayload) error {
	return f.MockSendSMSFn(ctx, workstationID, payload)
}
