package dto

type OrganizationIdentifier struct {
	Type  OrganizationIdentifierType `json:"type"`
	Value string                     `json:"value"`
}

type OrganizationInput struct {
	Name        string                   `json:"name"`
	PhoneNumber string                   `json:"phoneNumber"`
	Identifiers []OrganizationIdentifier `json:"identifiers"`
}
