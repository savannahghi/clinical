package clinical

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// CreateQuestionnaireResponse creates a questionnaire response
func (u *UseCasesClinicalImpl) CreateQuestionnaireResponse(ctx context.Context, questionnaireID string, encounterID string, input dto.QuestionnaireResponse) (*dto.QuestionnaireResponse, error) {
	questionnaireResponse := &domain.FHIRQuestionnaireResponse{}

	err := mapstructure.Decode(input, questionnaireResponse)
	if err != nil {
		return nil, err
	}

	encounter, err := u.infrastructure.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot record an observation in a finished encounter")
	}

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)

	questionnaireResponse.Source = &domain.FHIRReference{
		ID:        patientID,
		Reference: &patientReference,
		Display:   encounter.Resource.Subject.Display,
	}

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot patch an observation in a finished encounter")
	}

	questionnaireResponse.Encounter = &domain.FHIRReference{
		ID:        encounter.Resource.ID,
		Reference: &encounterReference,
	}

	questionnaireResponse.Questionnaire = &questionnaireID

	tags, err := u.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	questionnaireResponse.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	questionnaireResponse.Authored = &input.Authored
	questionnaireResponse.Status = input.Status

	resp, err := u.infrastructure.FHIR.CreateFHIRQuestionnaireResponse(ctx, questionnaireResponse)

	if err != nil {
		return nil, err
	}

	output := &dto.QuestionnaireResponse{}
	err = mapstructure.Decode(resp, output)

	if err != nil {
		return nil, err
	}

	return output, nil
}
