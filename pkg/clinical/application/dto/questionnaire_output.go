package dto

import (
	"github.com/savannahghi/scalarutils"
)

// QuestionnaireEdge is an instance of QuestionnaireEdge
type QuestionnaireEdge struct {
	Node   Questionnaire
	Cursor string
}

// QuestionnaireConnection models the questionnaire connection data class
type QuestionnaireConnection struct {
	TotalCount int
	Edges      []QuestionnaireEdge
	PageInfo   PageInfo
}

// QuestionnaireOutput models the data class to display questionnaire.
type QuestionnaireOutput struct {
	TotalCount int
	Edges      []QuestionnaireEdge
	PageInfo   PageInfo
}

// CreateQuestionnaireConnection creates a connection to map out questionnaire results following the GraphQl Cursor Connection Specification
func CreateQuestionnaireConnection(questionnaires []*Questionnaire, pageInfo PageInfo, total int) QuestionnaireConnection {
	connection := QuestionnaireConnection{
		TotalCount: total,
		Edges:      []QuestionnaireEdge{},
		PageInfo:   pageInfo,
	}

	for _, questionnaire := range questionnaires {
		edge := QuestionnaireEdge{
			Node:   *questionnaire,
			Cursor: questionnaire.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}

// Questionnaire models the dataclass to display questionnaires
type Questionnaire struct {
	ID                string               `json:"id,omitempty"`
	Meta              Meta                 `json:"meta,omitempty"`
	ImplicitRules     string               `json:"implicitRules,omitempty"`
	Language          string               `json:"language,omitempty"`
	Text              Narrative            `json:"text,omitempty"`
	Extension         []Extension          `json:"extension,omitempty"`
	ModifierExtension []Extension          `json:"modifierExtension,omitempty"`
	URL               scalarutils.URI      `json:"url,omitempty"`
	Identifier        []Identifier         `json:"identifier,omitempty"`
	Version           string               `json:"version,omitempty"`
	Name              string               `json:"name,omitempty"`
	Title             string               `json:"title,omitempty"`
	DerivedFrom       []string             `json:"derivedFrom,omitempty"`
	Status            scalarutils.Code     `json:"status,omitempty"`
	Experimental      bool                 `json:"experimental,omitempty"`
	Date              scalarutils.DateTime `json:"date,omitempty"`
	Publisher         string               `json:"publisher,omitempty"`
	Description       string               `json:"description,omitempty"`
	UseContext        UsageContext         `json:"useContext,omitempty"`
	Jurisdiction      []CodeableConcept    `json:"jurisdiction,omitempty"`
	Purpose           string               `json:"purpose,omitempty"`
	EffectivePeriod   Period               `json:"effectivePeriod,omitempty"`
	Code              []Coding             `json:"code,omitempty"`
	Item              []QuestionnaireItem  `json:"item,omitempty"`
}

type Narrative struct {
	ID     string            `json:"id,omitempty"`
	Status string            `json:"status,omitempty"`
	Div    scalarutils.XHTML `json:"div,omitempty"`
}

type MetaInput struct {
	VersionID   string               `json:"versionId,omitempty"`
	LastUpdated scalarutils.DateTime `json:"lastUpdated,omitempty"`
	Source      string               `json:"source,omitempty"`
	Tag         []CodingInput        `json:"tag,omitempty"`
	Security    []CodingInput        `json:"security,omitempty"`
}

type CodingInput struct {
	ID           string           `json:"id,omitempty"`
	System       scalarutils.URI  `json:"system,omitempty"`
	Version      string           `json:"version,omitempty"`
	Code         scalarutils.Code `json:"code,omitempty"`
	Display      string           `json:"display,omitempty"`
	UserSelected bool             `json:"userSelected,omitempty"`
}

type Extension struct {
	URL                  string           `json:"url,omitempty"`
	ValueBoolean         bool             `json:"valueBoolean,omitempty"`
	ValueInteger         int              `json:"valueInteger,omitempty"`
	ValueDecimal         float64          `json:"valueDecimal,omitempty"`
	ValueBase64Binary    string           `json:"valueBase64Binary,omitempty"`
	ValueInstant         string           `json:"valueInstant,omitempty"`
	ValueString          string           `json:"valueString,omitempty"`
	ValueURI             string           `json:"valueURI,omitempty"`
	ValueDate            string           `json:"valueDate,omitempty"`
	ValueDateTime        string           `json:"valueDateTime,omitempty"`
	ValueTime            string           `json:"valueTime,omitempty"`
	ValueCode            string           `json:"valueCode,omitempty"`
	ValueOid             string           `json:"valueOid,omitempty"`
	ValueUUID            string           `json:"valueUUID,omitempty"`
	ValueID              string           `json:"valueID,omitempty"`
	ValueUnsignedInt     int              `json:"valueUnsignedInt,omitempty"`
	ValuePositiveInt     int              `json:"valuePositiveInt,omitempty"`
	ValueMarkdown        string           `json:"valueMarkdown,omitempty"`
	ValueAnnotation      *Annotation      `json:"valueAnnotation,omitempty"`
	ValueAttachment      *Attachment      `json:"valueAttachment,omitempty"`
	ValueIdentifier      *Identifier      `json:"valueIdentifier,omitempty"`
	ValueCodeableConcept *CodeableConcept `json:"valueCodeableConcept,omitempty"`
	ValueCoding          *Coding          `json:"valueCoding,omitempty"`
	ValueQuantity        *Quantity        `json:"valueQuantity,omitempty"`
	ValueRange           *Range           `json:"valueRange,omitempty"`
	ValuePeriod          *Period          `json:"valuePeriod,omitempty"`
	ValueRatio           *Ratio           `json:"valueRatio,omitempty"`
	ValueReference       *Reference       `json:"valueReference,omitempty"`
}

type Annotation struct {
	ID              *string               `json:"id,omitempty"`
	AuthorReference *Reference            `json:"authorReference,omitempty"`
	AuthorString    *string               `json:"authorString,omitempty"`
	Time            *scalarutils.DateTime `json:"time,omitempty"`
	Text            *scalarutils.Markdown `json:"text,omitempty"`
}

type Range struct {
	ID   string   `json:"id,omitempty"`
	Low  Quantity `json:"low,omitempty"`
	High Quantity `json:"high,omitempty"`
}

type Ratio struct {
	ID          string   `json:"id,omitempty"`
	Numerator   Quantity `json:"numerator,omitempty"`
	Denominator Quantity `json:"denominator,omitempty"`
}

type QuestionnaireItem struct {
	ID                string                          `json:"id,omitempty"`
	Meta              Meta                            `json:"meta,omitempty"`
	Extension         []Extension                     `json:"extension,omitempty"`
	ModifierExtension []Extension                     `json:"modifierExtension,omitempty"`
	LinkID            string                          `json:"linkId,omitempty"`
	Definition        scalarutils.URI                 `json:"definition,omitempty"`
	Code              []Coding                        `json:"code,omitempty"`
	Prefix            string                          `json:"prefix,omitempty"`
	Text              string                          `json:"text,omitempty"`
	Type              scalarutils.Code                `json:"type,omitempty"`
	EnableWhen        []QuestionnaireItemEnableWhen   `json:"enableWhen,omitempty"`
	EnableBehavior    scalarutils.Code                `json:"enableBehavior,omitempty"`
	DisabledDisplay   scalarutils.Code                `json:"disabledDisplay,omitempty"`
	Required          bool                            `json:"required,omitempty"`
	Repeats           bool                            `json:"repeats,omitempty"`
	ReadOnly          bool                            `json:"readOnly,omitempty"`
	MaxLength         int                             `json:"maxLength,omitempty"`
	AnswerValueSet    string                          `json:"answerValueSet,omitempty"`
	AnswerOption      []QuestionnaireItemAnswerOption `json:"answerOption,omitempty"`
	Initial           []QuestionnaireItemInitial      `json:"initial,omitempty"`
	Item              []QuestionnaireItem             `json:"item,omitempty"`
}

type QuestionnaireItemEnableWhen struct {
	ID                string               `json:"id,omitempty"`
	Extension         []Extension          `json:"extension,omitempty"`
	ModifierExtension []Extension          `json:"modifierExtension,omitempty"`
	Question          string               `json:"question,omitempty"`
	Operator          scalarutils.Code     `json:"operator,omitempty"`
	AnswerBoolean     bool                 `json:"answerBoolean,omitempty"`
	AnswerDecimal     float64              `json:"answerDecimal,omitempty"`
	AnswerInteger     int                  `json:"answerInteger,omitempty"`
	AnswerDate        scalarutils.Date     `json:"answerDate,omitempty"`
	AnswerDateTime    scalarutils.DateTime `json:"answerDateTime,omitempty"`
	AnswerTime        scalarutils.DateTime `json:"answerTime,omitempty"`
	AnswerString      string               `json:"answerString,omitempty"`
	AnswerCoding      Coding               `json:"answerCoding,omitempty"`
	AnswerQuantity    Quantity             `json:"answerQuantity,omitempty"`
	AnswerReference   Reference            `json:"answerReference,omitempty"`
}

type QuestionnaireItemAnswerOption struct {
	ID                string           `json:"id,omitempty"`
	Extension         []Extension      `json:"extension,omitempty"`
	ModifierExtension []Extension      `json:"modifierExtension,omitempty"`
	ValueInteger      int              `json:"valueInteger,omitempty"`
	ValueDate         scalarutils.Date `json:"valueDate,omitempty"`
	ValueString       string           `json:"valueString,omitempty"`
	ValueCoding       Coding           `json:"valueCoding,omitempty"`
	ValueReference    Reference        `json:"valueReference,omitempty"`
	InitialSelected   bool             `json:"initialSelected,omitempty"`
}

type QuestionnaireItemInitial struct {
	ID                string               `json:"id,omitempty"`
	Extension         []Extension          `json:"extension,omitempty"`
	ModifierExtension []Extension          `json:"modifierExtension,omitempty"`
	ValueBoolean      bool                 `json:"valueBoolean,omitempty"`
	ValueDecimal      float64              `json:"valueDecimal,omitempty"`
	ValueInteger      int                  `json:"valueInteger,omitempty"`
	ValueDate         scalarutils.Date     `json:"valueDate,omitempty"`
	ValueDateTime     scalarutils.DateTime `json:"valuescalarutils.DateTime,omitempty"`
	ValueString       string               `json:"valueString,omitempty"`
	ValueURI          scalarutils.URI      `json:"valueUri,omitempty"`
	ValueAttachment   Attachment           `json:"valueAttachment,omitempty"`
	ValueCoding       Coding               `json:"valueCoding,omitempty"`
	ValueQuantity     Quantity             `json:"valueQuantity,omitempty"`
	ValueReference    Reference            `json:"valueReference,omitempty"`
}

type Quantity struct {
	ID         string           `json:"id,omitempty"`
	Value      float64          `json:"value,omitempty"`
	Comparator string           `json:"comparator,omitempty"`
	Unit       string           `json:"unit,omitempty"`
	System     scalarutils.URI  `json:"system,omitempty"`
	Code       scalarutils.Code `json:"code,omitempty"`
}

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

type UsageContext struct {
	ID                   string          `json:"id,omitempty"`
	Extension            []Extension     `json:"extension,omitempty"`
	Code                 Coding          `json:"code,omitempty"`
	ValueCodeableConcept CodeableConcept `json:"valueCodeableConcept,omitempty"`
	ValueQuantity        Quantity        `json:"valueQuantity,omitempty"`
	ValueRange           Range           `json:"valueRange,omitempty"`
	ValueReference       Reference       `json:"valueReference,omitempty"`
}

type Identifier struct {
	ID       string          `json:"id,omitempty"`
	Use      string          `json:"use,omitempty"`
	Type     CodeableConcept `json:"type,omitempty"`
	System   scalarutils.URI `json:"system,omitempty"`
	Value    string          `json:"value,omitempty"`
	Period   Period          `json:"period,omitempty"`
	Assigner *Reference      `json:"assigner,omitempty"`
}

type Reference struct {
	ID         string          `json:"id,omitempty"`
	Reference  string          `json:"reference,omitempty"`
	Type       scalarutils.URI `json:"type,omitempty"`
	Identifier Identifier      `json:"identifier,omitempty"`
	Display    string          `json:"display,omitempty"`
}

type CodeableConcept struct {
	ID     string   `json:"id,omitempty"`
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}

type Coding struct {
	ID           string           `json:"id,omitempty"`
	System       scalarutils.URI  `json:"system,omitempty"`
	Version      string           `json:"version,omitempty"`
	Code         scalarutils.Code `json:"code,omitempty"`
	Display      string           `json:"display,omitempty"`
	UserSelected bool             `json:"userSelected,omitempty"`
}

type Period struct {
	ID    string               `json:"id,omitempty"`
	Start scalarutils.DateTime `json:"start,omitempty"`
	End   scalarutils.DateTime `json:"end,omitempty"`
}

type Meta struct {
	VersionID string `json:"versionId,omitempty"`
	// LastUpdated time.Time `json:"lastUpdated,omitempty"`
	Source   string   `json:"source,omitempty"`
	Tag      []Coding `json:"tag,omitempty"`
	Security []Coding `json:"security,omitempty"`
}
