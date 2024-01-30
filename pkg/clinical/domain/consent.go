package domain

import (
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// Consent models a fhir consent resource.
type FHIRConsent struct {
	ID         *string                `json:"id,omitempty"`
	Status     *dto.ConsentStatusEnum `json:"status"`
	Scope      *FHIRCodeableConcept   `json:"scope"`
	Category   []*FHIRCodeableConcept `json:"category"`
	PolicyRule *FHIRCodeableConcept   `json:"policyRule,omitempty"`
	Provision  *FHIRConsentProvision  `json:"provision,omitempty"`
	Patient    *FHIRReference         `json:"patient,omitempty"`
	Meta       *FHIRMetaInput         `json:"meta,omitempty"`
	Extension  []Extension            `json:"extension,omitempty"`
}

// ConsentProvision models a fhir consent provision
type FHIRConsentProvision struct {
	ID   *string                       `json:"id,omitempty"`
	Type *dto.ConsentProvisionTypeEnum `json:"type,omitempty"`
}
