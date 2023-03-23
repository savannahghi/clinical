package clinical

import (
	"context"
	"fmt"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
	"time"
)

// CreateCondition creates a new conditions
func (c *UseCasesClinicalImpl) CreateCondition(ctx context.Context, input dto.ConditionInput) (*dto.Condition, error) {
	today := time.Now()

	date, err := scalarutils.NewDate(today.Day(), int(today.Month()), today.Year())
	if err != nil {
		return nil, err
	}

	conditionConcept, err := c.getICD10Concept(ctx, input.Code)
	if err != nil {
		return nil, err
	}

	statusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
	categorySystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-category")
	verificationSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-ver-status")
	userSelectedFalse := false
	conditionInput := domain.FHIRConditionInput{
		ClinicalStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  &statusSystem,
					Code:    scalarutils.Code(string(input.Status)),
					Display: string(input.Status),
				},
			},
			Text: string(input.Status),
		},
		VerificationStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &verificationSystem,
					Code:         scalarutils.Code("confirmed"),
					Display:      "confirmed",
					UserSelected: &userSelectedFalse,
				},
			},
			Text: "confirmed",
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       &categorySystem,
						Code:         scalarutils.Code("encounter-diagnosis"),
						Display:      "encounter-diagnosis",
						UserSelected: &userSelectedFalse,
					},
				},
				Text: "encounter-diagnosis",
			},
		},
		Code: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&conditionConcept.URL),
					Code:    scalarutils.Code(conditionConcept.ID),
					Display: conditionConcept.DisplayName,
				},
			},
			Text: conditionConcept.DisplayName,
		},
		RecordedDate: date,
	}

	if input.OnsetDate != nil {
		conditionInput.OnsetDateTime = input.OnsetDate
	}

	if input.Note != "" {
		note := scalarutils.Markdown(input.Note)
		noteTime := scalarutils.DateTime(time.Now().Format(scalarutils.DateTimeFormatLayout))

		conditionInput.Note = []*domain.FHIRAnnotationInput{
			{
				Time: &noteTime,
				Text: &note,
			},
		}
	}

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	patientType := scalarutils.URI("Patient")

	conditionInput.Subject = &domain.FHIRReferenceInput{
		ID:        patient.Resource.ID,
		Reference: &patientRef,
		Display:   patient.Resource.Name[0].Text,
		Type:      &patientType,
	}

	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot record a condition in a finished encounter")
	}

	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	encounterType := scalarutils.URI("Encounter")

	conditionInput.Encounter = &domain.FHIRReferenceInput{
		ID:        encounter.Resource.ID,
		Reference: &encounterRef,
		Display:   *encounter.Resource.ID,
		Type:      &encounterType,
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	conditionInput.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	condition, err := c.infrastructure.FHIR.CreateFHIRCondition(ctx, conditionInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRConditionToConditionDTO(*condition.Resource), nil
}

func mapFHIRConditionToConditionDTO(condition domain.FHIRCondition) *dto.Condition {
	output := dto.Condition{
		ID:           *condition.ID,
		Status:       dto.ConditionStatus(condition.ClinicalStatus.Text),
		Name:         condition.Code.Text,
		Code:         string(condition.Code.Coding[0].Code),
		System:       string(*condition.Code.Coding[0].System),
		RecordedDate: *condition.RecordedDate,
		PatientID:    *condition.Subject.ID,
		EncounterID:  *condition.Encounter.ID,
	}

	if condition.Note != nil && len(condition.Note) > 0 {
		output.Note = string(*condition.Note[0].Text)
	}

	if condition.OnsetDateTime != nil {
		output.OnsetDate = *condition.OnsetDateTime
	}

	return &output
}
