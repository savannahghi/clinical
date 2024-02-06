package domain

import (
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// FHIRQuestionnaireResponse models questionnaire response resource
type FHIRQuestionnaireResponse struct {
	ID                *string                             `json:"id,omitempty"`
	Meta              *FHIRMetaInput                      `json:"meta,omitempty"`
	ImplicitRules     *string                             `json:"implicitRules,omitempty"`
	Language          *string                             `json:"language,omitempty"`
	Text              *FHIRNarrative                      `json:"text,omitempty"`
	Extension         []FHIRExtension                     `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension                     `json:"modifierExtension,omitempty"`
	Identifier        *FHIRIdentifier                     `json:"identifier,omitempty"`
	BasedOn           []FHIRReference                     `json:"basedOn,omitempty"`
	PartOf            []FHIRReference                     `json:"partOf,omitempty"`
	Questionnaire     *string                             `json:"questionnaire,omitempty"`
	Status            dto.QuestionnaireResponseStatusEnum `json:"status"`
	Subject           *FHIRReference                      `json:"subject,omitempty"`
	Encounter         *FHIRReference                      `json:"encounter,omitempty"`
	Authored          *string                             `json:"authored,omitempty"`
	Author            *FHIRReference                      `json:"author,omitempty"`
	Source            *FHIRReference                      `json:"source,omitempty"`
	Item              []FHIRQuestionnaireResponseItem     `json:"item,omitempty"`
}

// FHIRQuestionnaireResponseItem models item object of questionnaire response resource
type FHIRQuestionnaireResponseItem struct {
	ID                *string                               `json:"id,omitempty"`
	Extension         []FHIRExtension                       `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension                       `json:"modifierExtension,omitempty"`
	LinkID            string                                `json:"linkId"`
	Definition        *string                               `json:"definition,omitempty"`
	Text              *string                               `json:"text,omitempty"`
	Answer            []FHIRQuestionnaireResponseItemAnswer `json:"answer,omitempty"`
	Item              []FHIRQuestionnaireResponseItem       `json:"item,omitempty"`
}

// FHIRQuestionnaireResponseItemAnswer models item answer object of questionnaire response resource
type FHIRQuestionnaireResponseItemAnswer struct {
	ID                *string                         `json:"id,omitempty"`
	Extension         []FHIRExtension                 `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension                 `json:"modifierExtension,omitempty"`
	ValueBoolean      *bool                           `json:"valueBoolean,omitempty"`
	ValueDecimal      *float64                        `json:"valueDecimal,omitempty"`
	ValueInteger      *int                            `json:"valueInteger,omitempty"`
	ValueDate         *string                         `json:"valueDate,omitempty"`
	ValueDateTime     *string                         `json:"valueDateTime,omitempty"`
	ValueTime         *string                         `json:"valueTime,omitempty"`
	ValueString       *string                         `json:"valueString,omitempty"`
	ValueURI          *string                         `json:"valueUri,omitempty"`
	ValueAttachment   *FHIRAttachment                 `json:"valueAttachment,omitempty"`
	ValueCoding       *FHIRCoding                     `json:"valueCoding,omitempty"`
	ValueQuantity     *FHIRQuantity                   `json:"valueQuantity,omitempty"`
	ValueReference    *FHIRReference                  `json:"valueReference,omitempty"`
	Item              []FHIRQuestionnaireResponseItem `json:"item,omitempty"`
}

// FHIRQuestionnaireResponseRelayPayload is used to return a single instance of Questionnaire response
type FHIRQuestionnaireResponseRelayPayload struct {
	Resource *FHIRQuestionnaireResponse `json:"resource,omitempty"`
}
