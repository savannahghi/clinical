package domain

import (
	"github.com/savannahghi/scalarutils"
)

// FHIRDocumentReference represents a reference to a document of any kind for any purpose.
// It provides metadata about the document so that the document can be discovered and managed.
// The scope of a document is any seralized object with a mime-type, so includes formal patient centric documents (CDA), cliical notes, scanned paper, and non-patient centric documents like policy text.
type FHIRDocumentReference struct {
	ID                string                           `json:"id,omitempty"`
	Meta              *FHIRMeta                        `json:"meta,omitempty"`
	ImplicitRules     *string                          `json:"implicitRules,omitempty"`
	Language          *string                          `json:"language,omitempty"`
	Text              *FHIRNarrative                   `json:"text,omitempty"`
	Extension         []FHIRExtension                  `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension                  `json:"modifierExtension,omitempty"`
	MasterIdentifier  *FHIRIdentifier                  `json:"masterIdentifier,omitempty"`
	Identifier        []FHIRIdentifier                 `json:"identifier,omitempty"`
	Status            DocumentReferenceStatusEnum      `json:"status,omitempty"`
	DocStatus         *CompositionStatusEnum           `json:"docStatus,omitempty"`
	Type              *FHIRCodeableConcept             `json:"type,omitempty"`
	Category          []FHIRCodeableConcept            `json:"category,omitempty"`
	Subject           *FHIRReference                   `json:"subject,omitempty"`
	Date              *string                          `json:"date,omitempty"`
	Author            []FHIRReference                  `json:"author,omitempty"`
	Authenticator     *FHIRReference                   `json:"authenticator,omitempty"`
	Custodian         *FHIRReference                   `json:"custodian,omitempty"`
	RelatesTo         []FHIRDocumentReferenceRelatesTo `json:"relatesTo,omitempty"`
	Description       string                           `json:"description,omitempty"`
	SecurityLabel     []FHIRCodeableConcept            `json:"securityLabel,omitempty"`
	Content           []FHIRDocumentReferenceContent   `json:"content,omitempty"`
	Context           *FHIRDocumentReferenceContext    `json:"context,omitempty"`
}

// FHIRDocumentReferenceRelatesTo specifies how this document reference is related to other resources,
// such as being a replacement for, transformation of, or addition to another document reference.
type FHIRDocumentReferenceRelatesTo struct {
	ID                string                       `json:"id,omitempty"`
	Extension         []Extension                  `json:"extension,omitempty"`
	ModifierExtension []Extension                  `json:"modifierExtension,omitempty"`
	Code              DocumentRelationshipTypeEnum `json:"code"`
	Target            Reference                    `json:"target"`
}

// FHIRDocumentReferenceContent describes the content of the document, including the document itself as an attachment, and potentially its format.
type FHIRDocumentReferenceContent struct {
	ID                string         `json:"id,omitempty"`
	Extension         []Extension    `json:"extension,omitempty"`
	ModifierExtension []Extension    `json:"modifierExtension,omitempty"`
	Attachment        FHIRAttachment `json:"attachment"`
	Format            *FHIRCoding    `json:"format,omitempty"`
}

// FHIRDocumentReferenceContext provides the clinical context in which the document was created, such as encounter, period, practice setting.
type FHIRDocumentReferenceContext struct {
	ID                string                `json:"id,omitempty"`
	Extension         []Extension           `json:"extension,omitempty"`
	ModifierExtension []Extension           `json:"modifierExtension,omitempty"`
	Encounter         []Reference           `json:"encounter,omitempty"`
	Event             []FHIRCodeableConcept `json:"event,omitempty"`
	Period            *FHIRPeriod           `json:"period,omitempty"`
	FacilityType      *FHIRCodeableConcept  `json:"facilityType,omitempty"`
	PracticeSetting   *FHIRCodeableConcept  `json:"practiceSetting,omitempty"`
	SourcePatientInfo *FHIRReference        `json:"sourcePatientInfo,omitempty"`
	Related           []*FHIRReference      `json:"related,omitempty"`
}

// FHIRDocumentReferenceInput is the input type for FHIRDocumentReference
type FHIRDocumentReferenceInput struct {
	ID                string                           `json:"id,omitempty"`
	Meta              *FHIRMetaInput                   `json:"meta,omitempty"`
	ImplicitRules     *string                          `json:"implicitRules,omitempty"`
	Language          *string                          `json:"language,omitempty"`
	Text              *FHIRNarrativeInput              `json:"text,omitempty"`
	Extension         []FHIRExtension                  `json:"extension,omitempty"`
	ModifierExtension []FHIRExtension                  `json:"modifierExtension,omitempty"`
	MasterIdentifier  *FHIRIdentifierInput             `json:"masterIdentifier,omitempty"`
	Identifier        []FHIRIdentifierInput            `json:"identifier,omitempty"`
	Status            DocumentReferenceStatusEnum      `json:"status,omitempty"`
	DocStatus         *CompositionStatusEnum           `json:"docStatus,omitempty"`
	Type              *FHIRCodeableConceptInput        `json:"type,omitempty"`
	Category          []FHIRCodeableConceptInput       `json:"category,omitempty"`
	Subject           *FHIRReferenceInput              `json:"subject,omitempty"`
	Date              *scalarutils.Instant             `json:"date,omitempty"`
	Author            []FHIRReferenceInput             `json:"author,omitempty"`
	Authenticator     *FHIRReferenceInput              `json:"authenticator,omitempty"`
	Custodian         *FHIRReferenceInput              `json:"custodian,omitempty"`
	RelatesTo         []FHIRDocumentReferenceRelatesTo `json:"relatesTo,omitempty"`
	Description       string                           `json:"description,omitempty"`
	SecurityLabel     []FHIRCodeableConceptInput       `json:"securityLabel,omitempty"`
	Content           []FHIRDocumentReferenceContent   `json:"content,omitempty"`
	Context           *FHIRDocumentReferenceContext    `json:"context,omitempty"`
}

// FHIRQuestionnaireResponseRelayPayload is used to return a single instance of document response
type FHIRDocumentReferenceRelayPayload struct {
	Resource *FHIRDocumentReference `json:"resource,omitempty"`
}

// PagedFHIRDocumentReference is an FHIR document's paginated model data class
type PagedFHIRDocumentReference struct {
	DocumentReferences []FHIRDocumentReference
	HasNextPage        bool
	NextCursor         string
	HasPreviousPage    bool
	PreviousCursor     string
	TotalCount         int
}
