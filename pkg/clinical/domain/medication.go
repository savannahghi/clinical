package domain

import (
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIRMedicationStatus indicates the medication status
type FHIRMedicationStatus string

const (
	// FHIRMedicationStatusActive is The medication is available for use
	FHIRMedicationStatusActive FHIRMedicationStatus = "active"

	// FHIRMedicationStatusInActive is The medication is not available for use.
	FHIRMedicationStatusInActive FHIRMedicationStatus = "inactive"

	// FHIRMedicationStatusEnteredInError is The medication was entered in error.
	FHIRMedicationStatusEnteredInError FHIRMedicationStatus = "entered-in-error"
)

// FHIRMedication definition:
type FHIRMedication struct {
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrative `json:"text,omitempty"`

	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	Code *FHIRCodeableConcept `json:"code,omitempty"`

	Status FHIRMedicationStatus `json:"status,omitempty"`

	Manufacturer *FHIROrganization `json:"manufacturer,omitempty"`

	Form *FHIRCodeableConcept `json:"form,omitempty"`

	Amount *FHIRRatio `json:"amount,omitempty"`

	Ingredient []*FHIRMedicationIngredient `json:"ingredient,omitempty"`

	Batch *FHIRMedicationBatch `json:"batch,omitempty"`
}

// FHIRMedicationBatch definition:
type FHIRMedicationBatch struct {
	LotNumber string `json:"lotNumber,omitempty"`

	ExpirationDate *scalarutils.Date `json:"expirationDate,omitempty"`
}

// FHIRMedicationIngredient definition:
type FHIRMedicationIngredient struct {
	ItemCodelabConcept *FHIRCodeableConcept `json:"itemCodelabConcept,omitempty"`

	ItemReference *FHIRReference `json:"itemReference,omitempty"`

	IsActive bool `json:"isActive,omitempty"`

	Strength *FHIRRatio `json:"strength,omitempty"`
}

// FHIRMedicationInput ...
type FHIRMedicationInput struct {
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrativeInput `json:"text,omitempty"`

	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	Code *FHIRCodeableConceptInput `json:"code,omitempty"`

	Status FHIRMedicationStatus `json:"status,omitempty"`

	Manufacturer *FHIROrganizationInput `json:"manufacturer,omitempty"`

	Form *FHIRCodeableConceptInput `json:"form,omitempty"`

	Amount *FHIRRatioInput `json:"amount,omitempty"`

	Ingredient []*FHIRMedicationIngredientInput `json:"ingredient,omitempty"`

	Batch *FHIRMedicationBatchInput `json:"batch,omitempty"`
}

// FHIRMedicationBatchInput ...
type FHIRMedicationBatchInput struct {
	LotNumber string `json:"lotNumber,omitempty"`

	ExpirationDate *scalarutils.Date `json:"expirationDate,omitempty"`
}

// FHIRMedicationIngredientInput definition:
type FHIRMedicationIngredientInput struct {
	ItemCodelabConcept *FHIRCodeableConceptInput `json:"itemCodelabConcept,omitempty"`

	ItemReference *FHIRReferenceInput `json:"itemReference,omitempty"`

	IsActive bool `json:"isActive,omitempty"`

	Strength *FHIRRatioInput `json:"strength,omitempty"`
}

// FHIRMedicationRelayConnection is a Relay connection for MedicationStatement
type FHIRMedicationRelayConnection struct {
	Edges []*FHIRMedicationRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIRMedicationRelayEdge is a Relay edge for Medication
type FHIRMedicationRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIRMedication `json:"node,omitempty"`
}

// FHIRMedicationRelayPayload is used to return single instances of Medication
type FHIRMedicationRelayPayload struct {
	Resource *FHIRMedication `json:"resource,omitempty"`
}
