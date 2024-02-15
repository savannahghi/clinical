package domain

import (
	"time"

	"github.com/savannahghi/scalarutils"
)

// FHIRMetaInput is the input to resource meta
type FHIRMetaInput struct {
	VersionID   string            `json:"versionId,omitempty"`
	LastUpdated time.Time         `json:"lastUpdated,omitempty"`
	Source      string            `json:"source,omitempty"`
	Tag         []FHIRCodingInput `json:"tag,omitempty"`
	Security    []FHIRCodingInput `json:"security,omitempty"`
	Profile     []scalarutils.URI `json:"profile,omitempty"`
}

// FHIRMeta is a set of metadata that provides technical and workflow context to a resource.
type FHIRMeta struct {
	VersionID string `json:"versionId,omitempty"`
	// LastUpdated time.Time    `json:"lastUpdated,omitempty"`
	Source   string       `json:"source,omitempty"`
	Tag      []FHIRCoding `json:"tag,omitempty"`
	Security []FHIRCoding `json:"security,omitempty"`
}
