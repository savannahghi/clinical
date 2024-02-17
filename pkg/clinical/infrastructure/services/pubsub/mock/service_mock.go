package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// FakeServicePubsub ...
type FakeServicePubsub struct {
	MockNotifyPatientFHIRIDUpdatefn  func(ctx context.Context, data dto.UpdatePatientFHIRID) error
	MockNotifyFacilityFHIRIDUpdatefn func(ctx context.Context, data dto.UpdateFacilityFHIRID) error
	MockNotifyProgramFHIRIDUpdatefn  func(ctx context.Context, data dto.UpdateProgramFHIRID) error
	MockNotifySegmentationFn         func(ctx context.Context, data dto.SegmentationPayload) error
}

// NewPubSubServiceMock mocks the pubsub service implementation
func NewPubSubServiceMock() *FakeServicePubsub {
	return &FakeServicePubsub{
		MockNotifyPatientFHIRIDUpdatefn: func(ctx context.Context, data dto.UpdatePatientFHIRID) error {
			return nil
		},
		MockNotifyFacilityFHIRIDUpdatefn: func(ctx context.Context, data dto.UpdateFacilityFHIRID) error {
			return nil
		},
		MockNotifyProgramFHIRIDUpdatefn: func(ctx context.Context, data dto.UpdateProgramFHIRID) error {
			return nil
		},
		MockNotifySegmentationFn: func(ctx context.Context, data dto.SegmentationPayload) error {
			return nil
		},
	}
}

// NotifyPatientFHIRIDUpdate publishes to patient fhir id update topic. Mycarehub service will subscribe to this topic
// and update the patient's FHIR ID in it's database
func (f *FakeServicePubsub) NotifyPatientFHIRIDUpdate(ctx context.Context, data dto.UpdatePatientFHIRID) error {
	return f.MockNotifyPatientFHIRIDUpdatefn(ctx, data)
}

// NotifyFacilityFHIRIDUpdate publishes to a topic. The idea is that after a mycarehub facility is created as an organization in FHIR,
// we should send back the ID to mycarehub and store in the database
func (f *FakeServicePubsub) NotifyFacilityFHIRIDUpdate(ctx context.Context, data dto.UpdateFacilityFHIRID) error {
	return f.MockNotifyFacilityFHIRIDUpdatefn(ctx, data)
}

// NotifyProgramFHIRIDUpdate publishes to the program fhir id update topic
func (f *FakeServicePubsub) NotifyProgramFHIRIDUpdate(ctx context.Context, data dto.UpdateProgramFHIRID) error {
	return f.MockNotifyProgramFHIRIDUpdatefn(ctx, data)
}

// NotifySegmentation publishes the mock data used to update the segmentation data in advantage
func (f *FakeServicePubsub) NotifySegmentation(ctx context.Context, data dto.SegmentationPayload) error {
	return f.MockNotifySegmentationFn(ctx, data)
}
