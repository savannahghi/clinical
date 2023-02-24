package domain

import (
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIRMedication definition:
type FHIRMedication struct {
	ID *string `json:"id,omitempty"`

	Text *FHIRNarrative `json:"text,omitempty"`

	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	Code *FHIRCodeableConcept `json:"code,omitempty"`

	Status MedicationStatusEnum `json:"status,omitempty"`

	Manufacturer *FHIROrganization `json:"manufacturer,omitempty"`

	Form *FHIRCodeableConcept `json:"form,omitempty"`

	Amount *FHIRRatio `json:"amount,omitempty"`

	Ingredient []*MedicationIngredient `json:"ingredient,omitempty"`

	Batch *MedicationBatch `json:"batch,omitempty"`

	// Meta stores more information about the resource
	Meta *FHIRMeta `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

// MedicationBatch definition:
type MedicationBatch struct {
	LotNumber string `json:"lotNumber,omitempty"`

	ExpirationDate *scalarutils.Date `json:"expirationDate,omitempty"`
}

// MedicationIngredient definition:
type MedicationIngredient struct {
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

	Status MedicationStatusEnum `json:"status,omitempty"`

	Manufacturer *FHIROrganizationInput `json:"manufacturer,omitempty"`

	Form *FHIRCodeableConceptInput `json:"form,omitempty"`

	Amount *FHIRRatioInput `json:"amount,omitempty"`

	Ingredient []*MedicationIngredientInput `json:"ingredient,omitempty"`

	Batch *MedicationBatchInput `json:"batch,omitempty"`

	// Meta stores more information about the resource
	Meta FHIRMetaInput `json:"meta,omitempty"`

	// Extension is an optional element that provides additional information not captured in the basic resource definition
	Extension []*FHIRExtension `json:"extension,omitempty"`
}

// MedicationBatchInput ...
type MedicationBatchInput struct {
	LotNumber string `json:"lotNumber,omitempty"`

	ExpirationDate *scalarutils.Date `json:"expirationDate,omitempty"`
}

// MedicationIngredientInput definition:
type MedicationIngredientInput struct {
	ItemCodelabConcept *FHIRCodeableConceptInput `json:"itemCodelabConcept,omitempty"`

	ItemReference *FHIRReferenceInput `json:"itemReference,omitempty"`

	IsActive bool `json:"isActive,omitempty"`

	Strength *FHIRRatioInput `json:"strength,omitempty"`
}

// FHIRMedicationRelayConnection is a Relay connection for Medication
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
