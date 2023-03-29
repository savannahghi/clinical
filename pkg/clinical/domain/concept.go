package domain

// Concept models a concept type from OpenConceptLab
type Concept struct {
	ConceptClass     string  `mapstructure:"concept_class" json:"concept_class"`
	DataType         string  `mapstructure:"datatype" json:"datatype"`
	DisplayLocale    string  `mapstructure:"display_locale" json:"display_locale"`
	DisplayName      string  `mapstructure:"display_name" json:"display_name"`
	ExternalID       string  `mapstructure:"external_id" json:"external_id"`
	ID               string  `mapstructure:"id" json:"id"`
	IsLatestVersion  bool    `mapstructure:"is_latest_version" json:"is_latest_version"`
	Locale           *string `mapstructure:"locale" json:"locale"`
	Owner            string  `mapstructure:"owner" json:"owner"`
	OwnerType        string  `mapstructure:"owner_type" json:"owner_type"`
	OwnerURL         string  `mapstructure:"owner_url" json:"owner_url"`
	Retired          bool    `mapstructure:"retired" json:"retired"`
	Source           string  `mapstructure:"source" json:"source"`
	Type             string  `mapstructure:"type" json:"type"`
	UpdateComment    string  `mapstructure:"update_comment" json:"update_comment"`
	URL              string  `mapstructure:"url" json:"url"`
	UUID             string  `mapstructure:"uuid" json:"uuid"`
	Version          string  `mapstructure:"version" json:"version"`
	VersionCreatedBy string  `mapstructure:"version_created_by" json:"version_created_by"`
	VersionCreatedOn string  `mapstructure:"version_created_on" json:"version_created_on"`
	VersionURL       string  `mapstructure:"version_url" json:"version_url"`
	VersionsURL      string  `mapstructure:"versions_url" json:"versions_url"`
}
