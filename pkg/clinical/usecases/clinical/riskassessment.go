package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// RecordRiskAssessment records a risk assessment based on the provided parameters.
func (c *UseCasesClinicalImpl) RecordRiskAssessment(
	ctx context.Context,
	encounterID string,
	questionnaireResponseID string,
	outcome domain.FHIRCodeableConcept,
) (*domain.FHIRRiskAssessment, error) {
	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot record a risk assessment in a finished encounter")
	}

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	riskAssessment := domain.FHIRRiskAssessmentInput{
		Status: domain.ObservationStatusEnumFinal,
		Code:   &domain.FHIRCodeableConceptInput{},
		Subject: domain.FHIRReferenceInput{
			ID:        encounter.Resource.Subject.ID,
			Reference: &patientReference,
			Display:   encounter.Resource.Subject.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.ID,
			Reference: &encounterReference,
		},
		OccurrenceDateTime: (*string)(&instant),
		Prediction: []domain.FHIRRiskAssessmentPrediction{
			{
				Outcome: &outcome,
			},
		},
	}

	if questionnaireResponseID != "" {
		questionnaireResponseReference := fmt.Sprintf("QuestionnaireResponse/%s", questionnaireResponseID)
		riskAssessment.Basis = []domain.FHIRReferenceInput{
			{
				Reference: &questionnaireResponseReference,
			},
		}
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	riskAssessment.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	assessment, err := c.infrastructure.FHIR.CreateFHIRRiskAssessment(ctx, &riskAssessment)
	if err != nil {
		return nil, err
	}

	return assessment.Resource, nil
}
