package dto

import (
	"mime/multipart"

	"github.com/go-playground/validator"
	"github.com/savannahghi/scalarutils"
)

type OrganizationIdentifier struct {
	Type  OrganizationIdentifierType `json:"type,omitempty"`
	Value string                     `json:"value,omitempty"`
}

type OrganizationInput struct {
	Name        string                   `json:"name,omitempty"`
	PhoneNumber string                   `json:"phoneNumber,omitempty"`
	Identifiers []OrganizationIdentifier `json:"identifiers,omitempty"`
}

type EpisodeOfCareInput struct {
	Status    EpisodeOfCareStatusEnum `json:"status"`
	PatientID string                  `json:"patientID"`
}

type EncounterInput struct {
	Status EncounterStatusEnum `json:"status"`
}

// ObservationInput models the observation input
type ObservationInput struct {
	Status      ObservationStatus `json:"status,omitempty" validate:"required"`
	EncounterID string            `json:"encounterID,omitempty" validate:"required"`
	Value       string            `json:"value,omitempty" validate:"required"`
}

func (o ObservationInput) Validate() error {
	v := validator.New()
	err := v.Struct(o)

	return err
}

type PatientInput struct {
	FirstName   string            `json:"firstName"`
	LastName    string            `json:"lastName"`
	OtherNames  *string           `json:"otherNames"`
	BirthDate   *scalarutils.Date `json:"birthDate,omitempty"`
	Gender      Gender            `json:"gender"`
	Identifiers []IdentifierInput `json:"identifiers"`
	Contacts    []ContactInput    `json:"contacts"`
}

type IdentifierInput struct {
	Type  IdentifierType `json:"type"`
	Value string         `json:"value"`
}

type ContactInput struct {
	Type  ContactType `json:"type"`
	Value string      `json:"value"`
}

// ConditionInput represents input for creating a FHIR condition
type ConditionInput struct {
	Code        string            `json:"condition"`
	System      TerminologySource `json:"system"`
	Status      ConditionStatus   `json:"status"`
	Category    ConditionCategory `json:"category"`
	EncounterID string            `json:"encounterID"`
	Note        string            `json:"note"`
	OnsetDate   *scalarutils.Date `json:"onsetDate"`
}

// AllergyInput models the allergy input
type AllergyInput struct {
	PatientID         string            `json:"patientID"`
	Code              string            `json:"code" validate:"required"`
	TerminologySource TerminologySource `json:"terminologySource" validate:"required"`
	EncounterID       string            `json:"encounterID" validate:"required,uuid4"`
	Reaction          *ReactionInput    `json:"reaction"`
}

// Validate ensures the input is valid
func (o AllergyInput) Validate() error {
	v := validator.New()
	err := v.Struct(o)

	return err
}

// ReactionInput models the reaction input
type ReactionInput struct {
	Code     string                                 `json:"code"`
	System   string                                 `json:"system"`
	Severity AllergyIntoleranceReactionSeverityEnum `json:"severity"`
}

// MediaInput models the dataclass to upload media to FHIR
type MediaInput struct {
	EncounterID string                             `json:"encounterID"`
	File        map[string][]*multipart.FileHeader `form:"file" json:"file"`
}

// CompositionInput models the composition input
type CompositionInput struct {
	EncounterID string                `json:"encounterID"`
	Type        CompositionType       `json:"type"`
	Category    CompositionCategory   `json:"category"`
	Status      CompositionStatusEnum `json:"status"`
	Note        string                `json:"note"`
}

// PatchCompositionInput models the patch composition input
type PatchCompositionInput struct {
	Type     CompositionType       `json:"type"`
	Category CompositionCategory   `json:"category"`
	Status   CompositionStatusEnum `json:"status"`
	Note     string                `json:"note"`
	Section  []*SectionInput       `json:"section"`
}

// SectionInput models the composition section input
type SectionInput struct {
	ID      string          `json:"id,omitempty"`
	Title   string          `json:"title"`
	Code    string          `json:"code"`
	Author  string          `json:"author"`
	Text    string          `json:"text"`
	Section []*SectionInput `json:"section"`
}

// ConsentInput models the consent input
type ConsentInput struct {
	Status    ConsentStatusEnum        `json:"status"`
	Provision ConsentProvisionTypeEnum `json:"provision,omitempty"`
	PatientID string                   `json:"patientID,omitempty"`
}

// QuestionnaireResponse models input for questionnaire response resource in fhir
type QuestionnaireResponse struct {
	Status   QuestionnaireResponseStatusEnum `json:"status"`
	Authored string                          `json:"authored,omitempty"`
	Item     []QuestionnaireResponseItem     `json:"item,omitempty"`
}

// QuestionnaireResponseItem models input for item object of questionnaire response resource
type QuestionnaireResponseItem struct {
	LinkID string                            `json:"linkId"`
	Text   *string                           `json:"text,omitempty"`
	Answer []QuestionnaireResponseItemAnswer `json:"answer,omitempty"`
	Item   []QuestionnaireResponseItem       `json:"item,omitempty"`
}

// FHIRQuestionnaireResponseItemAnswer models item answer object of questionnaire response resource
type QuestionnaireResponseItemAnswer struct {
	ValueBoolean    *bool                       `json:"valueBoolean,omitempty"`
	ValueDecimal    *float64                    `json:"valueDecimal,omitempty"`
	ValueInteger    *int                        `json:"valueInteger,omitempty"`
	ValueDate       *string                     `json:"valueDate,omitempty"`
	ValueDateTime   *string                     `json:"valueDateTime,omitempty"`
	ValueTime       *string                     `json:"valueTime,omitempty"`
	ValueString     *string                     `json:"valueString,omitempty"`
	ValueURI        *string                     `json:"valueUri,omitempty"`
	ValueAttachment *Attachment                 `json:"valueAttachment,omitempty"`
	ValueCoding     *Coding                     `json:"valueCoding,omitempty"`
	ValueQuantity   *Quantity                   `json:"valueQuantity,omitempty"`
	ValueReference  *Reference                  `json:"valueReference,omitempty"`
	Item            []QuestionnaireResponseItem `json:"item,omitempty"`
}

// Coding : an input for a code defined by a terminology system.
type Coding struct {
	ID           string           `json:"id,omitempty"`
	System       scalarutils.URI  `json:"system,omitempty"`
	Version      string           `json:"version,omitempty"`
	Code         scalarutils.Code `json:"code,omitempty"`
	Display      string           `json:"display,omitempty"`
	UserSelected bool             `json:"userSelected,omitempty"`
}

// Attachment definition: input for referring to data content defined in other formats.
type Attachment struct {
	ID          string                   `json:"id,omitempty"`
	ContentType scalarutils.Code         `json:"contentType,omitempty"`
	Language    scalarutils.Code         `json:"language,omitempty"`
	Data        scalarutils.Base64Binary `json:"data,omitempty"`
	URL         scalarutils.URL          `json:"url,omitempty"`
	Size        int                      `json:"size,omitempty"`
	Hash        scalarutils.Base64Binary `json:"hash,omitempty"`
	Title       string                   `json:"title,omitempty"`
	Creation    scalarutils.DateTime     `json:"creation,omitempty"`
}

// Quantity definition: input for measured amount (or an amount that can potentially be measured). note that measured amounts include amounts that are not precisely quantified, including amounts involving arbitrary units and floating currencies.
type Quantity struct {
	ID         string                  `json:"id,omitempty"`
	Value      float64                 `json:"value"`
	Comparator *QuantityComparatorEnum `json:"comparator,omitempty"`
	Unit       string                  `json:"unit"`
	System     scalarutils.URI         `json:"system"`
	Code       scalarutils.Code        `json:"code"`
}

// Reference definition: input for reference from one resource to another.
type Reference struct {
	ID         string          `json:"id,omitempty"`
	Reference  string          `json:"reference,omitempty"`
	Type       scalarutils.URI `json:"type,omitempty"`
	Identifier Identifier      `json:"identifier,omitempty"`
	Display    string          `json:"display,omitempty"`
}
