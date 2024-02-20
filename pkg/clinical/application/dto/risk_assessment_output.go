package dto

// RiskAssessment ...
type RiskAssessment struct {
	ID         *string                    `json:"id,omitempty" mapstructure:"id"`
	Subject    Reference                  `json:"subject,omitempty"`
	Encounter  *Reference                 `json:"encounter,omitempty"`
	Prediction []RiskAssessmentPrediction `json:"prediction,omitempty"`
	Note       []Annotation               `json:"note,omitempty"`
}

// RiskAssessmentPrediction describes the predicted outcome
type RiskAssessmentPrediction struct {
	ID                 *string          `json:"id,omitempty"`
	Outcome            *CodeableConcept `json:"outcome,omitempty"`
	ProbabilityDecimal *float64         `json:"probabilityDecimal,omitempty"`
}
