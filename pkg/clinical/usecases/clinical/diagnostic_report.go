package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// RecordMammographyResult is used to record mammography diagnostic reports as specified in https://hl7.org/fhir/R4/diagnosticreport.html#scope
func (c *UseCasesClinicalImpl) RecordMammographyResult(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	observationInput := &dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: input.EncounterID,
		Value:       input.Findings,
	}

	// !NOTE: The terminology code used here is used TEMPORARILY. Pending discussion about how to represent BI-RADs conclusions/observation
	observationOutput, err := c.RecordObservation(ctx, *observationInput, common.BenignNeoplasmOfBreastOfSkinTerminologyCode, []ObservationInputMutatorFunc{addImagingCategory})
	if err != nil {
		return nil, err
	}

	observationsReference := fmt.Sprintf("Observation/%s", observationOutput.ID)
	observationType := scalarutils.URI("Observation")
	encounterReference := fmt.Sprintf("Encounter/%s", observationOutput.EncounterID)
	encounterType := scalarutils.URI("Encounter")
	patientReference := fmt.Sprintf("Patient/%s", observationOutput.PatientID)
	patientType := scalarutils.URI("Patient")

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	facility, err := c.infrastructure.FHIR.GetFHIROrganization(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	orgRef := fmt.Sprintf("Organization/%s", *facility.Resource.ID)
	orgType := scalarutils.URI("Organization")
	id := uuid.New().String()
	mediaReference := fmt.Sprintf("Media/%s", input.Media.ID)
	mediaType := scalarutils.URI("Media")
	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	concept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, common.MammogramTerminologyCode)
	if err != nil {
		return nil, err
	}

	diagnosticReport := &domain.FHIRDiagnosticReportInput{
		ID:     &id,
		Status: domain.DiagnosticReportStatusPreliminary,
		Code: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&observationCategorySystem),
					Code:    scalarutils.Code(concept.ID),
					Display: concept.DisplayName,
				},
			},
			Text: concept.DisplayName,
		},
		Subject: &domain.FHIRReferenceInput{
			ID:        &observationOutput.PatientID,
			Reference: &patientReference,
			Type:      &patientType,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        &observationOutput.EncounterID,
			Reference: &encounterReference,
			Type:      &encounterType,
		},
		Issued: (*string)(&instant),
		Performer: []*domain.FHIRReferenceInput{
			{
				ID:        facility.Resource.ID,
				Reference: &orgRef,
				Type:      &orgType,
				Display:   *facility.Resource.Name,
			},
		},
		ResultsInterpreter: []*domain.FHIRReferenceInput{
			{
				ID:        facility.Resource.ID,
				Reference: &orgRef,
				Type:      &orgType,
				Display:   *facility.Resource.Name,
			},
		},
		Media: []*domain.FHIRDiagnosticReportMediaInput{
			{
				ID: &input.Media.ID,
				Link: &domain.FHIRReferenceInput{
					Reference: &mediaReference,
					Type:      &mediaType,
				},
			},
		},
		Conclusion: &input.Note,
		Result: []*domain.FHIRReferenceInput{
			{
				ID:        &observationOutput.ID,
				Reference: &observationsReference,
				Type:      &observationType,
			},
		},
	}

	diagnosticReport.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	result, err := c.infrastructure.FHIR.CreateFHIRDiagnosticReport(ctx, diagnosticReport)
	if err != nil {
		return nil, err
	}

	return &dto.DiagnosticReport{
		ID:          *result.ID,
		Status:      dto.ObservationStatus(result.Status),
		PatientID:   *result.Subject.ID,
		EncounterID: *result.Encounter.ID,
		Issued:      *result.Issued,
		Conclusion:  *result.Conclusion,
	}, nil
}
