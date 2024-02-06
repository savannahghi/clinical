package clinical

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
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
		return nil, fmt.Errorf("cannot create a questionnaire response in a finished encounter")
	}

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)

	questionnaireResponse.Source = &domain.FHIRReference{
		ID:        patientID,
		Reference: &patientReference,
		Display:   encounter.Resource.Subject.Display,
	}

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

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

	// TODO: This will affect the API performance. Optimize it
	err = u.generateQuestionnaireReviewSummary(
		ctx,
		questionnaireID,
		*resp.ID,
		encounterID,
		output,
	)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// generateQuestionnaireReviewSummary takes a questionnaire response and
// analyzes it to determine the risk stratification based on three distinct groups:
// symptoms, family history, and risk factors. The assumption is that the
// questionnaire has groups with linkIds: symptoms, family_history, and risk-factors.
// The function looks into the responses saved under the tags <group_name>-score,
// calculates the total scores for each group, and returns a summary indicating
// whether the individual is high risk, low risk, or average risk.
func (u *UseCasesClinicalImpl) generateQuestionnaireReviewSummary(
	ctx context.Context,
	questionnaireID,
	questionnaireResponseID,
	encounterID string,
	questionnaireResponse *dto.QuestionnaireResponse,
) error {
	questionnaire, err := u.infrastructure.FHIR.GetFHIRQuestionnaire(ctx, questionnaireID)
	if err != nil {
		return err
	}

	switch *questionnaire.Resource.Name {
	// TODO: Make this a controlled enum?
	case "Cervical Cancer Screening":
		var symptomsScore, riskFactorsScore, totalScore int

		for _, item := range questionnaireResponse.Item {
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
		}

		totalScore = symptomsScore + riskFactorsScore

		switch {
		case totalScore >= 2:
			err := u.recordRiskAssessment(
				ctx,
				encounterID,
				questionnaireResponseID,
				common.HighRiskCIELCode,
				"High Risk",
			)
			if err != nil {
				return err
			}

		case totalScore < 2:
			err := u.recordRiskAssessment(
				ctx,
				encounterID,
				questionnaireResponseID,
				common.LowRiskCIELCode,
				"Low Risk",
			)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("questionnaire does not exist")
	}

	return nil
}

func (u *UseCasesClinicalImpl) recordRiskAssessment(
	ctx context.Context,
	encounterID,
	questionnaireResponseID string,
	outcomeCode, outcomeDisplay string,
) error {
	CIELTerminologySystem := scalarutils.URI(common.CIELTerminologySystem)
	codingCode := scalarutils.Code(outcomeCode)

	outcome := domain.FHIRCodeableConcept{
		Coding: []*domain.FHIRCoding{
			{
				System:  &CIELTerminologySystem,
				Code:    &codingCode,
				Display: outcomeDisplay,
			},
		},
		Text: outcomeDisplay,
	}

	_, err := u.RecordRiskAssessment(ctx, encounterID, questionnaireResponseID, outcome)
	if err != nil {
		return err
	}

	return nil
}

// GetQuestionnaireResponseRiskLevel fetches the risk level associated with a questionnaire response. This is based off the scoring
// of the questionnaire response. Outcome: High Risk / Low Risk
func (u *UseCasesClinicalImpl) GetQuestionnaireResponseRiskLevel(ctx context.Context, questionnaireResponseID string) (string, error) {
	questionnaireResponse, err := u.infrastructure.FHIR.GetFHIRQuestionnaireResponse(ctx, questionnaireResponseID)
	if err != nil {
		return "", err
	}

	encounterReference := questionnaireResponse.Resource.Encounter.Reference
	patientReference := questionnaireResponse.Resource.Source.Reference

	riskAssessmentSearchParams := map[string]interface{}{
		"patient":   *patientReference,
		"encounter": *encounterReference,
		"_sort":     "date",
		"_count":    "1",
	}

	identifiers, err := u.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	riskAssessment, err := u.infrastructure.FHIR.SearchFHIRRiskAssessment(ctx, riskAssessmentSearchParams, *identifiers, dto.Pagination{})
	if err != nil {
		return "", err
	}

	riskLevel := ""

	for _, assessment := range riskAssessment.Edges {
		riskLevel = assessment.Node.Prediction[0].Outcome.Text
	}

	return riskLevel, nil
}
