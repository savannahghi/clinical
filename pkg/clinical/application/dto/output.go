package dto

// TenantIdentifiers models the json object used to store some of the tenant identifiers
type TenantIdentifiers struct {
	OrganizationID string `json:"organizationID,omitempty"`
}

type Organization struct {
	ID           string                   `json:"id"`
	Active       bool                     `json:"active"`
	Name         string                   `json:"name"`
	Identifiers  []OrganizationIdentifier `json:"identifiers"`
	PhoneNumbers []string                 `json:"phoneNumbers"`
}
