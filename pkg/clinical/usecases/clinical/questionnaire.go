package clinical

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
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

// ListQuestionnaires is used to list questionnaires from FHIR repository.
// This search is performed using the name or the title of the questionnaire and returns the available questionnaire(s).
func (q *UseCasesClinicalImpl) ListQuestionnaires(ctx context.Context, searchParam string, pagination *dto.Pagination) (*dto.Questionnaire, error) {
	identifiers, err := q.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	params := map[string]interface{}{
		"status": "active",
		"_sort":  "-date",
		"_count": "1",
	}

	if searchParam != "" {
		params["title:exact"] = searchParam
	}

	questionnaire, err := q.infrastructure.FHIR.ListFHIRQuestionnaire(ctx, params, *identifiers, *pagination)
	if err != nil {
		return nil, err
	}

	var dtoQuestionnaire *dto.Questionnaire

	for _, questionnaire := range questionnaire.Questionnaires {
		err := mapstructure.Decode(questionnaire, &dtoQuestionnaire)
		if err != nil {
			return nil, err
		}
	}

	return dtoQuestionnaire, nil
}
