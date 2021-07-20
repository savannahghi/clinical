package clinical

import (
	"github.com/savannahghi/firebasetools"
)

// FHIROrganization definition: The organization (facility) responsible for this organization
type FHIROrganization struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`
	// Whether the organization's record is still in active use
	Active *bool `json:"active,omitempty"`
	// Identifier(s) by which this organization is known.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`
	// Specific type of organization (e.g. Healthcare Provider, Hospital Department, Insurance Company).
	Type []*FHIRCodeableConcept `json:"type,omitempty"`
	// Name used for the organization
	Name *string `json:"name,omitempty"`
	// A list of alternate names that the organization is known as, or was known as in the past
	Alias []string `json:"alias,omitempty"`
	//A contact detail for the organization
	// ! Rule: The telecom of an organization can never be of use 'home'
	Telecom []*FHIRContactPoint `json:"telecom,omitempty"`

	// An address for the organization.
	Address []*FHIRAddress `json:"address,omitempty"`
}

// FHIROrganizationInput definition: The organization (facility) responsible for this organization
type FHIROrganizationInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`
	// Whether the organization's record is still in active use
	Active *bool `json:"active,omitempty"`
	// Identifier(s) by which this organization is known.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`
	// Specific type of organization (e.g. Healthcare Provider, Hospital Department, Insurance Company).
	Type []*FHIRCodeableConceptInput `json:"type,omitempty"`
	// Name used for the organization
	Name *string `json:"name,omitempty"`
	// A list of alternate names that the organization is known as, or was known as in the past
	Alias []string `json:"alias,omitempty"`
	//A contact detail for the organization
	// ! Rule: The telecom of an organization can never be of use 'home'
	Telecom []*FHIRContactPointInput `json:"telecom,omitempty"`

	// An address for the organization.
	Address []*FHIRAddressInput `json:"address,omitempty"`
}

// FHIROrganizationRelayPayload is used to return single instances of Organization
type FHIROrganizationRelayPayload struct {
	Resource *FHIROrganization `json:"resource,omitempty"`
}

// FHIROrganizationRelayConnection is a Relay connection for Organization
type FHIROrganizationRelayConnection struct {
	Edges    []*FHIROrganizationRelayEdge `json:"edges,omitempty"`
	PageInfo *firebasetools.PageInfo      `json:"pageInfo,omitempty"`
}

// FHIROrganizationRelayEdge is a Relay edge for Organization
type FHIROrganizationRelayEdge struct {
	Cursor *string           `json:"cursor,omitempty"`
	Node   *FHIROrganization `json:"node,omitempty"`
}
