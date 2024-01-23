package domain

import (
	"time"

	"github.com/savannahghi/scalarutils"
)

// FHIRQuestionnaire models the FHIR questionnaire model as described in https://www.hl7.org/fhir/questionnaire.html
type FHIRQuestionnaire struct {
	ID                *string                  `json:"id,omitempty"`
	Meta              *FHIRMetaInput           `json:"meta,omitempty"`
	ImplicitRules     *string                  `json:"implicitRules,omitempty"`
	Language          *string                  `json:"language,omitempty"`
	Text              *FHIRNarrative           `json:"text,omitempty"`
	FHIRExtension     []*Extension             `json:"extension,omitempty"`
	ModifierExtension []*Extension             `json:"modifierExtension,omitempty"`
	URL               *scalarutils.URI         `json:"url,omitempty"`
	Identifier        []*FHIRIdentifier        `json:"identifier,omitempty"`
	Version           *string                  `json:"version,omitempty"`
	Name              *string                  `json:"name,omitempty"`
	Title             *string                  `json:"title,omitempty"`
	DerivedFrom       []*string                `json:"derivedFrom,omitempty"`
	Status            *scalarutils.Code        `json:"status,omitempty"`
	Experimental      *bool                    `json:"experimental,omitempty"`
	Date              *scalarutils.DateTime    `json:"date,omitempty"`
	Publisher         *string                  `json:"publisher,omitempty"`
	Description       *string                  `json:"description,omitempty"`
	UseContext        *FHIRUsageContext        `json:"useContext,omitempty"`
	Jurisdiction      []*FHIRCodeableConcept   `json:"jurisdiction,omitempty"`
	Purpose           *string                  `json:"purpose,omitempty"`
	EffectivePeriod   *FHIRPeriod              `json:"effectivePeriod,omitempty"`
	Code              []*FHIRCoding            `json:"code,omitempty"`
	Item              []*FHIRQuestionnaireItem `json:"item,omitempty"`
}

// FHIRQuestionnaireItem represents the questions and sections within a FHIR questionnaire
type FHIRQuestionnaireItem struct {
	ID                *string                              `json:"id,omitempty"`
	Meta              *FHIRMeta                            `json:"meta,omitempty"`
	FHIRExtension     []*Extension                         `json:"extension,omitempty"`
	ModifierExtension []*Extension                         `json:"modifierExtension,omitempty"`
	LinkID            *string                              `json:"linkId,omitempty"`
	Definition        *scalarutils.URI                     `json:"definition,omitempty"`
	Code              []*FHIRCoding                        `json:"code,omitempty"`
	Prefix            *string                              `json:"prefix,omitempty"`
	Text              *string                              `json:"text,omitempty"`
	Type              *scalarutils.Code                    `json:"type,omitempty"`
	EnableWhen        []*FHIRQuestionnaireItemEnableWhen   `json:"enableWhen,omitempty"`
	EnableBehavior    *scalarutils.Code                    `json:"enableBehavior,omitempty"`
	DisabledDisplay   *scalarutils.Code                    `json:"disabledDisplay,omitempty"`
	Required          *bool                                `json:"required,omitempty"`
	Repeats           *bool                                `json:"repeats,omitempty"`
	ReadOnly          *bool                                `json:"readOnly,omitempty"`
	MaxLength         *int                                 `json:"maxLength,omitempty"`
	AnswerValueSet    *string                              `json:"answerValueSet,omitempty"`
	AnswerOption      []*FHIRQuestionnaireItemAnswerOption `json:"answerOption,omitempty"`
	Initial           []*FHIRQuestionnaireItemInitial      `json:"initial,omitempty"`
	Item              []*FHIRQuestionnaireItem             `json:"item,omitempty"`
}

// FHIRQuestionnaireItemEnableWhen defines when to enable the FHIR Questionnaire item.
type FHIRQuestionnaireItemEnableWhen struct {
	ID                *string               `json:"id,omitempty"`
	FHIRExtension     []*Extension          `json:"extension,omitempty"`
	ModifierExtension []*Extension          `json:"modifierExtension,omitempty"`
	Question          *string               `json:"question,omitempty"`
	Operator          *scalarutils.Code     `json:"operator,omitempty"`
	AnswerBoolean     *bool                 `json:"answerBoolean,omitempty"`
	AnswerDecimal     *float64              `json:"answerDecimal,omitempty"`
	AnswerInteger     *int                  `json:"answerInteger,omitempty"`
	AnswerDate        *scalarutils.Date     `json:"answerDate,omitempty"`
	AnswerDateTime    *scalarutils.DateTime `json:"answerDateTime,omitempty"`
	AnswerTime        *scalarutils.DateTime `json:"answerTime,omitempty"`
	AnswerString      *string               `json:"answerString,omitempty"`
	AnswerCoding      *FHIRCoding           `json:"answerCoding,omitempty"`
	AnswerQuantity    *FHIRQuantity         `json:"answerQuantity,omitempty"`
	AnswerReference   *FHIRReference        `json:"answerReference,omitempty"`
}

// FHIRQuestionnaireItemAnswerOption represents the permitted answers to a questionnaire.
// ! Rule: A question cannot have both answerOption and answerValueSet
// ! Rule: Only coding, decimal, integer, date, dateTime, time, string or quantity items can have answerOption or answerValueSet
// ! Rule: If one or more answerOption is present, initial cannot be present. Use answerOption.initialSelected instead
type FHIRQuestionnaireItemAnswerOption struct {
	ID                *string           `json:"id,omitempty"`
	FHIRExtension     []*Extension      `json:"extension,omitempty"`
	ModifierExtension []*Extension      `json:"modifierExtension,omitempty"`
	ValueInteger      *int              `json:"valueInteger,omitempty"`
	ValueDate         *scalarutils.Date `json:"valueDate,omitempty"`
	ValueTime         *time.Time        `json:"valueTime,omitempty"`
	ValueString       string            `json:"valueString,omitempty"`
	ValueCoding       *FHIRCoding       `json:"valueCoding,omitempty"`
	ValueReference    *FHIRReference    `json:"valueReference,omitempty"`
	InitialSelected   *bool             `json:"initialSelected,omitempty"`
}

// FHIRQuestionnaireItemInitial defines the initial value(s) when a questionnaire item is first rendered
type FHIRQuestionnaireItemInitial struct {
	ID                *string               `json:"id,omitempty"`
	FHIRExtension     []*Extension          `json:"extension,omitempty"`
	ModifierExtension []*Extension          `json:"modifierExtension,omitempty"`
	ValueBoolean      *bool                 `json:"valueBoolean,omitempty"`
	ValueDecimal      *float64              `json:"valueDecimal,omitempty"`
	ValueInteger      *int                  `json:"valueInteger,omitempty"`
	ValueDate         *scalarutils.Date     `json:"valueDate,omitempty"`
	ValueDateTime     *scalarutils.DateTime `json:"valueDateTime,omitempty"`
	ValueTime         *time.Time            `json:"valueTime,omitempty"`
	ValueString       string                `json:"valueString,omitempty"`
	ValueURI          *scalarutils.URI      `json:"valueUri,omitempty"`
	ValueAttachment   *FHIRAttachment       `json:"valueAttachment,omitempty"`
	ValueCoding       *FHIRCoding           `json:"valueCoding,omitempty"`
	ValueQuantity     *FHIRQuantity         `json:"valueQuantity,omitempty"`
	ValueReference    *FHIRReference        `json:"valueReference,omitempty"`
}

// FHIRUsageContext describes the context that the questionnaire content is intended to support
type FHIRUsageContext struct {
	ID                   *string              `json:"id,omitempty"`
	FHIRExtension        []*Extension         `json:"extension,omitempty"`
	Code                 *FHIRCoding          `json:"code,omitempty"`
	ValueCodeableConcept *FHIRCodeableConcept `json:"valueCodeableConcept,omitempty"`
	ValueQuantity        *FHIRQuantity        `json:"valueQuantity,omitempty"`
	ValueRange           *FHIRRange           `json:"valueRange,omitempty"`
	ValueReference       *FHIRReference       `json:"valueReference,omitempty"`
}
