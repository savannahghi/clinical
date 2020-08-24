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

// SingleFHIRPeriodPayload - compose a period payload
func SingleFHIRPeriodPayload() *FHIRPeriodInput {
	var start base.DateTime = "2019-01-03"
	var end base.DateTime = "2019-05-03"
	return &FHIRPeriodInput{
		Start: start,
		End:   end,
	}
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

// SingleFHIRDurationPayload - compose a duration payload
func SingleFHIRDurationPayload() *FHIRDurationInput {
	var value base.Decimal
	var comparator DurationComparatorEnum = "less_than"
	var system base.URI = "http://terminology.hl7.org/CodeSystem/v2-0131"
	unit := "Hours"
	var code base.Code = "HU"

	return &FHIRDurationInput{
		Value:      &value,
		Comparator: &comparator,
		Unit:       &unit,
		System:     &system,
		Code:       &code,
	}
}
