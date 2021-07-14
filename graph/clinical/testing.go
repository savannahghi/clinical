// Package clinical -  common helper functions and test fixtures
package clinical

import (
	"github.com/savannahghi/scalarutils"
)

// SingleIdentifierInput - single indentifier input
type SingleIdentifierInput struct {
	IdentifierUse IdentifierUseEnum
	Value         string
	System        scalarutils.URI
	UserSelected  bool
	Version       string
	Code          scalarutils.Code
}

// ReferenceInput - FHIR reference input
type ReferenceInput struct {
	Reference  string
	URL        scalarutils.URI
	Display    string
	Identifier *SingleIdentifierInput
}

// SingleFHIRCodingPayload - compose an FHIRCodingInput
func SingleFHIRCodingPayload(code scalarutils.Code, display string) *FHIRCodingInput {
	userSelected := true
	var system scalarutils.URI = "http://terminology.hl7.org/CodeSystem/v2-0131"
	version := "2.0"

	return &FHIRCodingInput{
		System:       &system,
		Code:         code,
		Version:      &version,
		Display:      display,
		UserSelected: &userSelected,
	}
}
