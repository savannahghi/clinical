package dto

// UserInfo is a collection of standard profile information for a user.
type UserInfo struct {
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	PhotoURL    string `json:"photoUrl,omitempty"`
	// In the ProviderUserInfo[] ProviderID can be a short domain name (e.g. google.com),
	// or the identity of an OpenID identity provider.
	// In UserRecord.UserInfo it will return the constant string "firebase".
	ProviderID string `json:"providerId,omitempty"`
	UID        string `json:"rawId,omitempty"`
}

// ErrorDetails contains more details about the error that occurred while making a REST API call to FHIR servers
type ErrorDetails struct {
	Text string `json:"text"`
}

// ErrorIssue models the error issue returned from FHIR
type ErrorIssue struct {
	Details     ErrorDetails `json:"details"`
	Diagnostics string       `json:"diagnostics"`
}

// ErrorResponse models the json object returned when an error occurs while calling the
// FHIR server REST APIs
type ErrorResponse struct {
	Issue []ErrorIssue `json:"issue"`
}

// TenantIdentifiers models the json object used to store some of the tenant identifiers
type TenantIdentifiers struct {
	OrganizationID string `json:"organizationID,omitempty"`
}
