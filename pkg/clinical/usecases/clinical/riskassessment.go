package clinical

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// RecordRiskAssessment records a risk assessment based on the provided parameters.
func (c *UseCasesClinicalImpl) RecordRiskAssessment(
	ctx context.Context,
	riskAssessment *domain.FHIRRiskAssessmentInput,
) (*domain.FHIRRiskAssessment, error) {
	assessment, err := c.infrastructure.FHIR.CreateFHIRRiskAssessment(ctx, riskAssessment)
	if err != nil {
		return nil, err
	}

	return assessment.Resource, nil
}
