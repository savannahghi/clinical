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

	return c.RecordDiagnosticReport(ctx, common.MammogramTerminologyCode, input, observationOutput, nil)
}

// RecordBiopsy is used to record biopsy test results as a diagnostic report
// FHIR recommends use of diagnostic resource to record the findings and interpretation of biopsy test results
// performed on patients, groups of patients, devices, and locations, and/or specimens.
func (c *UseCasesClinicalImpl) RecordBiopsy(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	observationInput := &dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: input.EncounterID,
		Value:       input.Findings,
	}

	observationOutput, err := c.RecordObservation(ctx, *observationInput, common.BiopsyTerminologySystem, []ObservationInputMutatorFunc{addProcedureCategory})
	if err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, common.BiopsyTerminologySystem, input, observationOutput, []DiagnosticReportMutatorFunc{addCytopathologyCategory})
}

// RecordMRI is used to record MRI scan results as a diagnostic report
func (c *UseCasesClinicalImpl) RecordMRI(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	observationInput := &dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: input.EncounterID,
		Value:       input.Findings,
	}

	observationOutput, err := c.RecordObservation(ctx, *observationInput, common.MRITerminologySystem, []ObservationInputMutatorFunc{addProcedureCategory})
	if err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, common.MRITerminologySystem, input, observationOutput, []DiagnosticReportMutatorFunc{addNuclearMagneticResonanceCategory})
}

// RecordUltrasound is used to record the breast ultrasound diagnostic reports
func (c *UseCasesClinicalImpl) RecordUltrasound(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	observationInput := &dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: input.EncounterID,
		Value:       input.Findings,
	}

	observationOutput, err := c.RecordObservation(ctx, *observationInput, common.BilateralConceptTerminologySystem, []ObservationInputMutatorFunc{addImagingCategory})
	if err != nil {
		return nil, err
	}

	// TODO: `BilateralConceptTerminologySystem` is a `PLACE HOLDER`. It should be adjusted accordingly when the breast cancer designs are revamped
	return c.RecordDiagnosticReport(ctx, common.BilateralConceptTerminologySystem, input, observationOutput, []DiagnosticReportMutatorFunc{addRadiologyUltrasoundCategory})
}

// RecordCBE is used to record clinical based examination test results for a patient
func (c *UseCasesClinicalImpl) RecordCBE(ctx context.Context, input *dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	observationInput := &dto.ObservationInput{
		Status:      dto.ObservationStatusFinal,
		EncounterID: input.EncounterID,
		Value:       input.Findings,
	}

	observationOutput, err := c.RecordObservation(ctx, *observationInput, common.BreastExaminationCIELTerminologySystem, []ObservationInputMutatorFunc{addExamCategory})
	if err != nil {
		return nil, err
	}

	return c.RecordDiagnosticReport(ctx, common.BreastExaminationCIELTerminologySystem, *input, observationOutput, []DiagnosticReportMutatorFunc{addOtherCategory})
}

// RecordDiagnosticReport is a re-usable method to help with diagnostic report recording
func (c *UseCasesClinicalImpl) RecordDiagnosticReport(ctx context.Context, conceptID string, input dto.DiagnosticReportInput, observation *dto.Observation, mutators []DiagnosticReportMutatorFunc) (*dto.DiagnosticReport, error) {
	observationsReference := fmt.Sprintf("Observation/%s", observation.ID)
	observationType := scalarutils.URI("Observation")
	encounterReference := fmt.Sprintf("Encounter/%s", observation.EncounterID)
	encounterType := scalarutils.URI("Encounter")
	patientReference := fmt.Sprintf("Patient/%s", observation.PatientID)
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
	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	concept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, conceptID)
	if err != nil {
		return nil, err
	}

	diagnosticReport := &domain.FHIRDiagnosticReportInput{
		ID:       &id,
		Status:   domain.DiagnosticReportStatusFinal,
		Category: []*domain.FHIRCodeableConceptInput{},
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
			ID:        &observation.PatientID,
			Reference: &patientReference,
			Type:      &patientType,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        &observation.EncounterID,
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
		Conclusion: &input.Note,
		Result: []*domain.FHIRReferenceInput{
			{
				ID:        &observation.ID,
				Reference: &observationsReference,
				Type:      &observationType,
			},
		},
	}

	diagnosticReport.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	if input.Media != nil {
		mediaReference := fmt.Sprintf("Media/%s", input.Media.ID)
		mediaType := scalarutils.URI("Media")

		diagnosticReport.Media = []*domain.FHIRDiagnosticReportMediaInput{
			{
				ID: &input.Media.ID,
				Link: &domain.FHIRReferenceInput{
					Reference: &mediaReference,
					Type:      &mediaType,
				},
			},
		}
	}

	if len(mutators) > 0 {
		for _, mutator := range mutators {
			err = mutator(ctx, diagnosticReport)
			if err != nil {
				return nil, err
			}
		}
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
