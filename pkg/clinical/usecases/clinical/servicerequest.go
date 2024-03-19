package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
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
	count int,
) (*dto.ServiceRequest, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, err
	}

	includedResources := []string{
		"Observation:encounter",
	}

	searchParams := map[string]interface{}{
		"_id":         input.EncounterID,
		"_sort":       "_lastUpdated",
		"_revinclude": includedResources,
	}

	encounterAllData, err := c.infrastructure.FHIR.SearchFHIREncounterAllData(ctx, searchParams, *identifiers, dto.Pagination{
		First: &count,
	})
	if err != nil {
		return nil, err
	}

	output, err := c.mapFHIREncounterDataToEncounterAssociatedDTO(encounterAllData)
	if err != nil {
		return nil, err
	}

	payload := &CompositionPayload{
		ConceptID: common.LOINCReferralNote,
		CompositionInput: &dto.CompositionInput{
			EncounterID: input.EncounterID,
			Category:    dto.ReferralNote,
			Status:      dto.CompositionStatusEnumFinal,
		},
		SectionData: c.mapEncounterAssociatedDToFHIRSectionInput(output),
	}

	compositionOutput, err := c.RecordComposition(ctx, *payload)
	if err != nil {
		return nil, err
	}

	// Create service request
	patientID := compositionOutput.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)

	encounterReference := fmt.Sprintf("Encounter/%s", *compositionOutput.Resource.Encounter.ID)
	startTime := scalarutils.DateTime(time.Now().Format("2006-01-02T15:04:05+03:00"))
	compositionReference := fmt.Sprintf("Composition/%s", *compositionOutput.Resource.ID)

	serviceRequest := domain.FHIRServiceRequestInput{
		Status:     domain.ServiceRequestStatusActive,
		Intent:     domain.ServiceRequestIntentOrder,
		Priority:   domain.ServiceRequestPriorityUrgent,
		AuthoredOn: &startTime,
		Subject: &domain.FHIRReferenceInput{
			ID:        compositionOutput.Resource.Subject.ID,
			Reference: &patientReference,
			Display:   compositionOutput.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        compositionOutput.Resource.Encounter.ID,
			Reference: &encounterReference,
		},
		Note: []*domain.FHIRAnnotationInput{
			{
				Time: &startTime,
				Text: (*scalarutils.Markdown)(&input.ReferralNote),
			},
		},
		SupportingInfo: []*domain.FHIRReferenceInput{
			{
				ID:        compositionOutput.Resource.ID,
				Reference: &compositionReference,
			},
		},
	}

	if input.Facility != "" {
		facilityExtension := &domain.FHIRExtension{
			URL: "http://savannahghi.org/fhir/StructureDefinition/referred-facility",
			Extension: []domain.Extension{
				{
					URL:         "facilityName",
					ValueString: input.Facility,
				},
			},
		}
		serviceRequest.Extension = append(serviceRequest.Extension, facilityExtension)
	}

	if input.Specialist != "" {
		specialistExtension := &domain.FHIRExtension{
			URL: "http://savannahghi.org/fhir/StructureDefinition/referred-specialist",
			Extension: []domain.Extension{
				{
					URL:         "specialistName",
					ValueString: input.Specialist,
				},
			},
		}
		serviceRequest.Extension = append(serviceRequest.Extension, specialistExtension)
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

func (c *UseCasesClinicalImpl) mapEncounterAssociatedDToFHIRSectionInput(associatedEncounterResources *dto.EncounterAssociatedResources) []*domain.FHIRCompositionSectionInput {
	compositionSectionInput := []*domain.FHIRCompositionSectionInput{}
	compositionSectionTextStatus := "generated"

	for _, observation := range associatedEncounterResources.Observation {
		resourceName := "Observation"
		compositionSectionInput = append(compositionSectionInput, &domain.FHIRCompositionSectionInput{
			Title: &resourceName,
			Section: []*domain.FHIRCompositionSectionInput{
				{
					Title: &observation.Name,
					Text: &domain.FHIRNarrativeInput{
						Div:    scalarutils.XHTML(observation.Value),
						Status: (*domain.NarrativeStatusEnum)(&compositionSectionTextStatus),
					},
				}},
		})
	}

	return compositionSectionInput
}
