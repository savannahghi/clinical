package domain

import "github.com/savannahghi/firebasetools"

// FHIRRiskAssessment resource captures predicted outcomes for a patient or population on the basis of source information
// http://hl7.org/fhir/StructureDefinition/RiskAssessment
type FHIRRiskAssessment struct {
	ID                 *string                        `json:"id,omitempty"`
	Meta               *FHIRMeta                      `json:"meta,omitempty"`
	ImplicitRules      *string                        `json:"implicitRules,omitempty"`
	Language           *string                        `json:"language,omitempty"`
	Text               *FHIRNarrative                 `json:"text,omitempty"`
	Extension          []Extension                    `json:"extension,omitempty"`
	ModifierExtension  []Extension                    `json:"modifierExtension,omitempty"`
	Identifier         []FHIRIdentifier               `json:"identifier,omitempty"`
	BasedOn            *Reference                     `json:"basedOn,omitempty"`
	Parent             *Reference                     `json:"parent,omitempty"`
	Status             ObservationStatusEnum          `json:"status,omitempty"`
	Method             *FHIRCodeableConcept           `json:"method,omitempty"`
	Code               *FHIRCodeableConcept           `json:"code,omitempty"`
	Subject            Reference                      `json:"subject,omitempty"`
	Encounter          *Reference                     `json:"encounter,omitempty"`
	OccurrenceDateTime *string                        `json:"occurrenceDateTime,omitempty"`
	OccurrencePeriod   *FHIRPeriod                    `json:"occurrencePeriod,omitempty"`
	Condition          *Reference                     `json:"condition,omitempty"`
	Performer          *Reference                     `json:"performer,omitempty"`
	ReasonCode         []FHIRCodeableConcept          `json:"reasonCode,omitempty"`
	ReasonReference    []Reference                    `json:"reasonReference,omitempty"`
	Basis              []Reference                    `json:"basis,omitempty"`
	Prediction         []FHIRRiskAssessmentPrediction `json:"prediction,omitempty"`
	Mitigation         *string                        `json:"mitigation,omitempty"`
	Note               []FHIRAnnotation               `json:"note,omitempty"`
}

// FHIRRiskAssessmentPrediction describes the predicted outcome
type FHIRRiskAssessmentPrediction struct {
	ID                 *string              `json:"id,omitempty"`
	Extension          []Extension          `json:"extension,omitempty"`
	ModifierExtension  []Extension          `json:"modifierExtension,omitempty"`
	Outcome            *FHIRCodeableConcept `json:"outcome,omitempty"`
	ProbabilityDecimal *float64             `json:"probabilityDecimal,omitempty"`
	ProbabilityRange   *FHIRRange           `json:"probabilityRange,omitempty"`
	QualitativeRisk    *FHIRCodeableConcept `json:"qualitativeRisk,omitempty"`
	RelativeRisk       *float64             `json:"relativeRisk,omitempty"`
	WhenPeriod         *FHIRPeriod          `json:"whenPeriod,omitempty"`
	WhenRange          *FHIRRange           `json:"whenRange,omitempty"`
	Rationale          *string              `json:"rationale,omitempty"`
}

// FHIRRiskAssessmentRelayPayload is used to return single instance of RiskAssessment
type FHIRRiskAssessmentRelayPayload struct {
	Resource *FHIRRiskAssessment `json:"resource,omitempty"`
}

type FHIRRiskAssessmentInput struct {
	ID                 *string                        `json:"id,omitempty"`
	Meta               *FHIRMetaInput                 `json:"meta,omitempty"`
	ImplicitRules      *string                        `json:"implicitRules,omitempty"`
	Language           *string                        `json:"language,omitempty"`
	Text               *FHIRNarrativeInput            `json:"text,omitempty"`
	Extension          []Extension                    `json:"extension,omitempty"`
	ModifierExtension  []Extension                    `json:"modifierExtension,omitempty"`
	Identifier         []FHIRIdentifierInput          `json:"identifier,omitempty"`
	BasedOn            *FHIRReferenceInput            `json:"basedOn,omitempty"`
	Parent             *FHIRReferenceInput            `json:"parent,omitempty"`
	Status             ObservationStatusEnum          `json:"status,omitempty"`
	Method             *FHIRCodeableConceptInput      `json:"method,omitempty"`
	Code               *FHIRCodeableConceptInput      `json:"code,omitempty"`
	Subject            FHIRReferenceInput             `json:"subject,omitempty"`
	Encounter          *FHIRReferenceInput            `json:"encounter,omitempty"`
	OccurrenceDateTime *string                        `json:"occurrenceDateTime,omitempty"`
	OccurrencePeriod   *FHIRPeriodInput               `json:"occurrencePeriod,omitempty"`
	Condition          *FHIRReferenceInput            `json:"condition,omitempty"`
	Performer          *FHIRReferenceInput            `json:"performer,omitempty"`
	ReasonCode         []FHIRCodeableConceptInput     `json:"reasonCode,omitempty"`
	ReasonReference    []FHIRReferenceInput           `json:"reasonReference,omitempty"`
	Basis              []FHIRReferenceInput           `json:"basis,omitempty"`
	Prediction         []FHIRRiskAssessmentPrediction `json:"prediction,omitempty"`
	Mitigation         *string                        `json:"mitigation,omitempty"`
	Note               []FHIRAnnotationInput          `json:"note,omitempty"`
}

// FHIRRiskAssessmentRelayConnection is a Relay connection for RiskAssessment
type FHIRRiskAssessmentRelayConnection struct {
	Edges []*FHIRRiskAssessmentRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIRRiskAssessmentRelayEdge is a Relay edge for RiskAssessment
type FHIRRiskAssessmentRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIRRiskAssessment `json:"node,omitempty"`
}
