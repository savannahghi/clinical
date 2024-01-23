package clinical

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// CreateQuestionnaire is used to create a new Questionnaire.
// These questionnaire are used to solicit various types of information from patients to server organisation usecases.
func (q *UseCasesClinicalImpl) CreateQuestionnaire(ctx context.Context, questionnaireInput *domain.FHIRQuestionnaire) (*domain.FHIRQuestionnaire, error) {
	tags, err := q.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	questionnaireInput.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	questionnaire, err := q.infrastructure.FHIR.CreateFHIRQuestionnaire(ctx, questionnaireInput)
	if err != nil {
		return nil, err
	}

	return questionnaire, nil
}
