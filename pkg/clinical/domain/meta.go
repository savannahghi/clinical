package domain

// Meta field in FHIR is used to capture metadata about a resource
type Meta struct {
	VersionId   string          `json:"versionId,omitempty"`
	LastUpdated string          `json:"lastUpdated,omitempty"`
	Source      string          `json:"source,omitempty"`
	Profile     []string        `json:"profile,omitempty"`
	Security    []SecurityLabel `json:"security,omitempty"`
	Tag         []FHIRCoding    `json:"tag,omitempty"`
}

// SecurityLabel is used to manage access control to the resource
type SecurityLabel struct {
	System  string `json:"system,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}
