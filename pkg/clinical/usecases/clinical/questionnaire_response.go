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

	u.generateSummary(ctx, questionnaireID, *resp)

	output := &dto.QuestionnaireResponse{}

	err = mapstructure.Decode(resp, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// GetQuestionnaireReviewSummary fetches a summary of a completed questionnaire response
func (u *UseCasesClinicalImpl) GetQuestionnaireReviewSummary(_ context.Context, _ string) (*dto.QuestionnaireReviewSummary, error) {
	return nil, nil
}

func (u UseCasesClinicalImpl) generateSummary(ctx context.Context, questionnaireID string, questionnaireResponse domain.FHIRQuestionnaireResponse) error {
	/*
		1. Get questionnaire
		2. Switch on the questionnaire (cervical, breast, mental etc)
		3. calculate and save in Obs resource
	*/
	questionnaire, err := u.infrastructure.FHIR.GetFHIRQuestionnaire(ctx, questionnaireID)
	if err != nil {
		return err
	}

	switch *questionnaire.Resource.Name {
	case "CERVICAL":
		/*
			Total symptoms
			Total family history
			Total risk factors
		*/
		var familyHistoryScore, symptomsScore, riskFactorsScore int

		for _, item := range questionnaireResponse.Item {
			if item.LinkID == "family-history" {
				for _, answer := range item.Item {
					if answer.LinkID == "family-history-score" {
						familyHistoryScore = *answer.Answer[0].ValueInteger
					}
				}
			}

			if item.LinkID == "symptoms" {
				for _, answer := range item.Item {
					if answer.LinkID == "symptoms-score" {
						symptomsScore = *answer.Answer[0].ValueInteger
					}
				}
			}

			if item.LinkID == "risk-factors" {
				for _, answer := range item.Item {
					if answer.LinkID == "risk-factors-score" {
						riskFactorsScore = *answer.Answer[0].ValueInteger
					}
				}
			}

			if familyHistoryScore > 1 {

			}
			if symptomsScore > 1 {

			}
			if riskFactorsScore > 1 {

			}
		}
	case "BREAST":

	}

	return nil
}
