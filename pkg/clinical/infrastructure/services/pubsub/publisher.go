package pubsubmessaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

func (ps ServicePubSubMessaging) newPublish(
	ctx context.Context,
	data interface{},
	topic, serviceName string,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data received: %w", err)
	}

	return ps.PublishToPubsub(
		ctx,
		ps.AddPubSubNamespace(topic, serviceName),
		serviceName,
		payload,
	)
}

// NotifyPatientFHIRIDUpdate publishes to patient fhir id update topic. Mycarehub service will subscribe to this topic
// and update the patient's FHIR ID in it's database
func (ps ServicePubSubMessaging) NotifyPatientFHIRIDUpdate(ctx context.Context, data dto.UpdatePatientFHIRID) error {
	return ps.newPublish(ctx, data, common.AddFHIRIDToPatientProfile, MyCareHubServiceName)
}

// NotifyFacilityFHIRIDUpdate publishes to a topic. The idea is that after a mycarehub facility is created as an organization in FHIR,
// we should send back the ID to mycarehub and store in the database
func (ps ServicePubSubMessaging) NotifyFacilityFHIRIDUpdate(ctx context.Context, data dto.UpdateFacilityFHIRID) error {
	return ps.newPublish(ctx, data, common.AddFHIRIDToFacility, MyCareHubServiceName)
}

// NotifyProgramFHIRIDUpdate publishes to the program fhir id update topic
func (ps ServicePubSubMessaging) NotifyProgramFHIRIDUpdate(ctx context.Context, data dto.UpdateProgramFHIRID) error {
	return ps.newPublish(ctx, data, common.AddFHIRIDToProgram, MyCareHubServiceName)
}

// NotifySegmentation publishes the the data used to update the segmentation data in advantage
func (ps ServicePubSubMessaging) NotifySegmentation(ctx context.Context, data dto.SegmentationPayload) error {
	return ps.newPublish(ctx, data, common.SegmentationTopicName, common.ClinicalServiceName)
}
