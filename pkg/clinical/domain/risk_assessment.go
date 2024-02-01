package domain

// FHIRRiskAssessment resource captures predicted outcomes for a patient or population on the basis of source information
// http://hl7.org/fhir/StructureDefinition/RiskAssessment
type FHIRRiskAssessment struct {
	ID                 *string                        `bson:"id,omitempty" json:"id,omitempty"`
	Meta               *FHIRMeta                      `bson:"meta,omitempty" json:"meta,omitempty"`
	ImplicitRules      *string                        `bson:"implicitRules,omitempty" json:"implicitRules,omitempty"`
	Language           *string                        `bson:"language,omitempty" json:"language,omitempty"`
	Text               *FHIRNarrative                 `bson:"text,omitempty" json:"text,omitempty"`
	Extension          []Extension                    `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension  []Extension                    `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Identifier         []FHIRIdentifier               `bson:"identifier,omitempty" json:"identifier,omitempty"`
	BasedOn            *Reference                     `bson:"basedOn,omitempty" json:"basedOn,omitempty"`
	Parent             *Reference                     `bson:"parent,omitempty" json:"parent,omitempty"`
	Status             ObservationStatusEnum          `bson:"status" json:"status"`
	Method             *FHIRCodeableConcept           `bson:"method,omitempty" json:"method,omitempty"`
	Code               *FHIRCodeableConcept           `bson:"code,omitempty" json:"code,omitempty"`
	Subject            Reference                      `bson:"subject" json:"subject"`
	Encounter          *Reference                     `bson:"encounter,omitempty" json:"encounter,omitempty"`
	OccurrenceDateTime *string                        `bson:"occurrenceDateTime,omitempty" json:"occurrenceDateTime,omitempty"`
	OccurrencePeriod   *FHIRPeriod                    `bson:"occurrencePeriod,omitempty" json:"occurrencePeriod,omitempty"`
	Condition          *Reference                     `bson:"condition,omitempty" json:"condition,omitempty"`
	Performer          *Reference                     `bson:"performer,omitempty" json:"performer,omitempty"`
	ReasonCode         []FHIRCodeableConcept          `bson:"reasonCode,omitempty" json:"reasonCode,omitempty"`
	ReasonReference    []Reference                    `bson:"reasonReference,omitempty" json:"reasonReference,omitempty"`
	Basis              []Reference                    `bson:"basis,omitempty" json:"basis,omitempty"`
	Prediction         []FHIRRiskAssessmentPrediction `bson:"prediction,omitempty" json:"prediction,omitempty"`
	Mitigation         *string                        `bson:"mitigation,omitempty" json:"mitigation,omitempty"`
	Note               []FHIRAnnotation               `bson:"note,omitempty" json:"note,omitempty"`
}

// FHIRRiskAssessmentPrediction describes the predicted outcome
type FHIRRiskAssessmentPrediction struct {
	ID                 *string              `bson:"id,omitempty" json:"id,omitempty"`
	Extension          []Extension          `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension  []Extension          `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Outcome            *FHIRCodeableConcept `bson:"outcome,omitempty" json:"outcome,omitempty"`
	ProbabilityDecimal *float64             `bson:"probabilityDecimal,omitempty" json:"probabilityDecimal,omitempty"`
	ProbabilityRange   *FHIRRange           `bson:"probabilityRange,omitempty" json:"probabilityRange,omitempty"`
	QualitativeRisk    *FHIRCodeableConcept `bson:"qualitativeRisk,omitempty" json:"qualitativeRisk,omitempty"`
	RelativeRisk       *float64             `bson:"relativeRisk,omitempty" json:"relativeRisk,omitempty"`
	WhenPeriod         *FHIRPeriod          `bson:"whenPeriod,omitempty" json:"whenPeriod,omitempty"`
	WhenRange          *FHIRRange           `bson:"whenRange,omitempty" json:"whenRange,omitempty"`
	Rationale          *string              `bson:"rationale,omitempty" json:"rationale,omitempty"`
}

// FHIRRiskAssessmentRelayPayload is used to return single instance of RiskAssessment
type FHIRRiskAssessmentRelayPayload struct {
	Resource *FHIRRiskAssessment `json:"resource,omitempty"`
}
