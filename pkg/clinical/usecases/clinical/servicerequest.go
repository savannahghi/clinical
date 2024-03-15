package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// ReferPatient creates a new patient referral based on the specified inputs.
// This method handles the logic for generating a referral for different purposes,
// such as diagnostics and testing, specialist consultation, or treatment.
// The method constructs a ServiceRequest object following the FHIR R4 standards
func (c *UseCasesClinicalImpl) ReferPatient(
	ctx context.Context,
	input *dto.ReferralInput,
) (*dto.ServiceRequest, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot record a referral in a finished encounter")
	}

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	startTime := scalarutils.DateTime(time.Now().Format("2006-01-02T15:04:05+03:00"))

	serviceRequest := domain.FHIRServiceRequestInput{
		Status:     domain.ServiceRequestStatusActive,
		Intent:     domain.ServiceRequestIntentOrder,
		Priority:   domain.ServiceRequestPriorityUrgent,
		AuthoredOn: &startTime,
		Subject: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.Subject.ID,
			Reference: &patientReference,
			Display:   encounter.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.ID,
			Reference: &encounterReference,
		},
		Note: []*domain.FHIRAnnotationInput{
			{
				Time: &startTime,
				Text: (*scalarutils.Markdown)(&input.ReferralNote),
			},
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	serviceRequest.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	referralRequest, err := c.infrastructure.FHIR.CreateFHIRServiceRequest(ctx, serviceRequest)
	if err != nil {
		return nil, err
	}

	return &dto.ServiceRequest{
		ID:       *referralRequest.Resource.ID,
		Status:   string(referralRequest.Resource.Status),
		Intent:   string(referralRequest.Resource.Intent),
		Priority: string(referralRequest.Resource.Priority),
		Note: []dto.Annotation{
			{
				Text: *referralRequest.Resource.Note[0].Text,
			},
		},
		Subject: dto.Reference{
			ID:        *referralRequest.Resource.Subject.ID,
			Reference: *referralRequest.Resource.Subject.Reference,
			Display:   referralRequest.Resource.Subject.Display,
		},
		Encounter: &dto.Reference{
			ID:        *referralRequest.Resource.Encounter.ID,
			Reference: *referralRequest.Resource.Encounter.Reference,
			Display:   referralRequest.Resource.Encounter.Display,
		},
	}, nil
}
