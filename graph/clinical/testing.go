// Package clinical -  common helper functions and test fixtures
package clinical

import "gitlab.slade360emr.com/go/base"

// SingleIdentifierInput - single indentifier input
type SingleIdentifierInput struct {
	IdentifierUse IdentifierUseEnum
	Value         string
	System        base.URI
	UserSelected  bool
	Version       string
	Code          base.Code
}

// ReferenceInput - FHIR reference input
type ReferenceInput struct {
	Reference  string
	URL        base.URI
	Display    string
	Identifier *SingleIdentifierInput
}

// SingleFHIRCodingPayload - compose an FHIRCodingInput
func SingleFHIRCodingPayload(code base.Code, display string) *FHIRCodingInput {
	userSelected := true
	var system base.URI = "http://terminology.hl7.org/CodeSystem/v2-0131"
	version := "2.0"

	return &FHIRCodingInput{
		System:       &system,
		Code:         code,
		Version:      &version,
		Display:      display,
		UserSelected: &userSelected,
	}
}
