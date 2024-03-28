package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// CreateQuestionnaireResponse creates a questionnaire response
func (u *UseCasesClinicalImpl) CreateQuestionnaireResponse(ctx context.Context, questionnaireID string, encounterID string, input dto.QuestionnaireResponse) (string, error) {
	questionnaireResponse := &domain.FHIRQuestionnaireResponse{}

	err := mapstructure.Decode(input, questionnaireResponse)
	if err != nil {
		return "", err
	}

	encounter, err := u.infrastructure.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return "", err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return "", fmt.Errorf("cannot create a questionnaire response in a finished encounter")
	}

	// TODO: Ensure user cannot submit the same risk assessment twice in the same encounter

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
		return "", err
	}

	questionnaireResponse.Meta.Tag = tags

	questionnaireResponse.Authored = &input.Authored
	questionnaireResponse.Status = input.Status

	resp, err := u.infrastructure.FHIR.CreateFHIRQuestionnaireResponse(ctx, questionnaireResponse)
	if err != nil {
		return "", err
	}

	output := &dto.QuestionnaireResponse{}

	err = mapstructure.Decode(resp, output)
	if err != nil {
		return "", err
	}

	// TODO: This will affect the API performance. Optimize it
	riskLevel, err := u.generateQuestionnaireReviewSummary(
		ctx,
		questionnaireID,
		*resp.ID,
		encounter,
		output,
	)
	if err != nil {
		return "", err
	}

	return riskLevel, nil
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
	questionnaireResponseID string,
	encounter *domain.FHIREncounterRelayPayload,
	questionnaireResponse *dto.QuestionnaireResponse,
) (string, error) {
	riskLevel := ""

	questionnaire, err := u.infrastructure.FHIR.GetFHIRQuestionnaire(ctx, questionnaireID)
	if err != nil {
		return "", err
	}

	patient, err := u.infrastructure.FHIR.GetFHIRPatient(ctx, *encounter.Resource.Subject.ID)
	if err != nil {
		return "", err
	}

	switch *questionnaire.Resource.Title {
	// TODO: Make this a controlled enum?
	case "Cervical Cancer Screening":
		scores := calculateScoresFromResponses(questionnaire.Resource, questionnaireResponse)

		totalScore := scores["symptoms"] + scores["risk-factors"]

		switch {
		case totalScore >= 2:
			riskLevel, err = u.recordRiskAssessment(
				ctx,
				encounter,
				questionnaireResponseID,
				common.HighRiskCIELCode,
				"High Risk",
				domain.CervicalCancerScreeningTypeEnum.Text(), // TODO: This is TEMPORARY. A follow up PR is to follow supplying the value from params
			)
			if err != nil {
				return "", err
			}

			err := u.infrastructure.Pubsub.NotifySegmentation(ctx, dto.SegmentationPayload{
				ClinicalID:   *patient.Resource.ID,
				SegmentLabel: dto.SegmentationCategoryHighRiskNegative,
			})
			if err != nil {
				return "", err
			}

		case totalScore < 2:
			riskLevel, err = u.recordRiskAssessment(
				ctx,
				encounter,
				questionnaireResponseID,
				common.LowRiskCIELCode,
				"Low Risk",
				domain.CervicalCancerScreeningTypeEnum.Text(),
			)
			if err != nil {
				return "", err
			}

			err := u.infrastructure.Pubsub.NotifySegmentation(ctx, dto.SegmentationPayload{
				ClinicalID:   *patient.Resource.ID,
				SegmentLabel: dto.SegmentationCategoryLowRisk,
			})
			if err != nil {
				return "", err
			}
		}

	case "Breast Cancer Screening":
		scores := calculateScoresFromResponses(questionnaire.Resource, questionnaireResponse)

		if scores["risk-assessment"] >= 1 {
			riskLevel, err = u.recordRiskAssessment(
				ctx,
				encounter,
				questionnaireResponseID,
				common.HighRiskCIELCode,
				"High Risk",
				domain.BreastCancerScreeningTypeEnum.String(),
			)
			if err != nil {
				return "", err
			}
		} else {
			riskLevel, err = u.recordRiskAssessment(
				ctx,
				encounter,
				questionnaireResponseID,
				common.HighRiskCIELCode,
				"Average Risk",
				domain.BreastCancerScreeningTypeEnum.String(),
			)
			if err != nil {
				return "", err
			}
		}

	default:
		return "", fmt.Errorf("questionnaire does not exist")
	}

	return riskLevel, nil
}

// calculateScoresFromResponses Calculates scores for each group in the questionnaire based on responses.
// This function iterates through each item in the QuestionnaireResponse, matching it with its corresponding item in the Questionnaire
// to apply scoring rules based on 'ordinalValue'. This matching is essential to accurately compute scores since it directly ties
// the respondent's answers back to the structured definitions within the Questionnaire, including the consideration of 'ordinalValue'
// for scoring purposes.
func calculateScoresFromResponses(questionnaire *domain.FHIRQuestionnaire, response *dto.QuestionnaireResponse) map[string]int {
	groupScores := make(map[string]int) // Holds the calculated scores for each group.

	// Iterate through the response groups (items)
	for _, responseGroup := range response.Item {
		score := 0
		groupLinkId := responseGroup.LinkID

		// Find matching group in the questionnaire by linkId
		for _, questionnaireGroup := range questionnaire.Item {
			if *questionnaireGroup.LinkID == groupLinkId {
				// Calculate score for this group based on responses
				score = calculateGroupScore(questionnaireGroup.Item, responseGroup.Item)
				break
			}
		}

		groupScores[groupLinkId] = score
	}

	return groupScores
}

func calculateGroupScore(questionnaireItems []*domain.FHIRQuestionnaireItem, responseItems []dto.QuestionnaireResponseItem) int {
	score := 0

	// Extract the ordinal values and their corresponding display text (e.g., "Yes").
	ordinalValues := mapLinkIdsToOrdinalValuesAndDisplay(questionnaireItems)

	for _, responseItem := range responseItems {
		// Each responseItem.LinkId might have an associated ordinal value and display text.
		if val, ok := ordinalValues[responseItem.LinkID]; ok {
			for _, answer := range responseItem.Answer {
				if answer.ValueCoding.Display == val.display {
					score += val.ordinalValue
				}
			}
		}
	}

	return score
}

// mapLinkIdsToOrdinalValuesAndDisplay Maps each linkId in the Questionnaire to its corresponding 'ordinalValue', if available.
// This mapping facilitates a more efficient scoring process by eliminating the need to repeatedly search through the
// questionnaire structure for 'ordinalValue' extensions during scoring.
func mapLinkIdsToOrdinalValuesAndDisplay(items []*domain.FHIRQuestionnaireItem) map[string]struct {
	ordinalValue int
	display      string
} {
	ordinalValues := make(map[string]struct {
		ordinalValue int
		display      string
	})

	for _, item := range items {
		for _, answerOption := range item.AnswerOption {
			// Check if this answerOption has an ordinalValue and is a "Yes".
			hasOrdinalValue := false
			ordinalValue := 0
			displayText := ""

			for _, ext := range answerOption.Extension {
				if ext.URL == "http://hl7.org/fhir/StructureDefinition/ordinalValue" {
					hasOrdinalValue = true
					ordinalValue = int(*ext.ValueDecimal)
				}
			}

			if hasOrdinalValue && answerOption.ValueCoding.Display == "Yes" {
				displayText = answerOption.ValueCoding.Display
				ordinalValues[*item.LinkID] = struct {
					ordinalValue int
					display      string
				}{ordinalValue, displayText}
			}
		}
	}

	return ordinalValues
}

func (u *UseCasesClinicalImpl) recordRiskAssessment(
	ctx context.Context,
	encounter *domain.FHIREncounterRelayPayload,
	questionnaireResponseID, outcomeCode string,
	outcomeDisplay, usageContext string,
) (string, error) {
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

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)

	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	questionnaireResponseReference := fmt.Sprintf("QuestionnaireResponse/%s", questionnaireResponseID)

	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	textStatus := domain.NarrativeStatusEnumAdditional
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
		Basis: []domain.FHIRReferenceInput{
			{
				Reference: &questionnaireResponseReference,
			},
		},
		Text: &domain.FHIRNarrativeInput{
			Status: &textStatus,
			Div:    scalarutils.XHTML(usageContext),
		},
	}

	tags, err := u.GetTenantMetaTags(ctx)
	if err != nil {
		return "", err
	}

	riskAssessment.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	assessment, err := u.RecordRiskAssessment(ctx, &riskAssessment)
	if err != nil {
		return "", err
	}

	riskLevel := assessment.Prediction[0].Outcome.Text

	return riskLevel, nil
}

// GetQuestionnaireResponseRiskLevel fetches the risk level associated with a questionnaire response. This is based off the scoring
// of the questionnaire response. Outcome: High Risk / Low Risk
func (u *UseCasesClinicalImpl) GetQuestionnaireResponseRiskLevel(ctx context.Context, encounterID string, screeningType domain.ScreeningTypeEnum) (string, error) {
	encounter, err := u.infrastructure.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return "", err
	}

	patientID := encounter.Resource.Subject.ID
	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	patientReference := fmt.Sprintf("Patient/%s", *patientID)

	riskAssessmentSearchParams := map[string]interface{}{
		"patient":   patientReference,
		"encounter": encounterReference,
		"_text":     screeningType.Text(),
		"_sort":     "date",
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
