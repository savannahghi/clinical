package domain

import "time"

// Concept models a concept type from OpenConceptLab
type Concept struct {
	ConceptClass     string     `json:"concept_class"`
	DataType         string     `json:"datatype"`
	DisplayLocale    string     `json:"display_locale"`
	DisplayName      string     `json:"display_name"`
	ExternalID       string     `json:"external_id"`
	ID               string     `json:"id"`
	IsLatestVersion  bool       `json:"is_latest_version"`
	Locale           *string    `json:"locale"`
	Owner            string     `json:"owner"`
	OwnerType        string     `json:"owner_type"`
	OwnerURL         string     `json:"owner_url"`
	Retired          bool       `json:"retired"`
	Source           string     `json:"source"`
	Type             string     `json:"type"`
	UpdateComment    string     `json:"update_comment"`
	URL              string     `json:"url"`
	UUID             string     `json:"uuid"`
	Version          string     `json:"version"`
	VersionCreatedBy string     `json:"version_created_by"`
	VersionCreatedOn *time.Time `json:"version_created_on"`
	VersionURL       string     `json:"version_url"`
	VersionsURL      string     `json:"versions_url"`
}
