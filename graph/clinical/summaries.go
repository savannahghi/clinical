package clinical

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.slade360emr.com/go/base"
	"jaytaylor.com/html2text"
)

const emptyMarkdown = base.Markdown("")

// VisitSummary is used to serialize the visit summary to the client
type VisitSummary struct {
	Encounter          []base.Markdown `json:"encounter,omitempty"`
	AllergyIntolerance []base.Markdown `json:"allergyIntolerance,omitempty"`
	Condition          []base.Markdown `json:"condition,omitempty"`
	Observation        []base.Markdown `json:"observation,omitempty"`
	Composition        []base.Markdown `json:"composition,omitempty"`
	MedicationRequest  []base.Markdown `json:"medicationRequest,omitempty"`
	ServiceRequest     []base.Markdown `json:"serviceRequest,omitempty"`
}

// HumanDateTime renders a FHIR DateTime in a human friendly way
func HumanDateTime(dt base.DateTime) base.Markdown {
	// reference time: Mon Jan 2 15:04:05 MST 2006
	formatted := string(dt)
	ts := string(dt)
	if len(dt) == 4 {
		t, err := time.Parse("2006", ts)
		if err != nil {
			formatted = string(dt)
			return base.Markdown(formatted)
		}
		formatted = t.Format("2006")
		return base.Markdown(formatted)
	}
	if len(dt) == 7 {
		// YYYY-MM
		t, err := time.Parse("2006-01", ts)
		if err != nil {
			formatted = string(dt)
			return base.Markdown(formatted)
		}
		formatted = t.Format("Jan 2006")
		return base.Markdown(formatted)
	}
	if len(dt) == 10 {
		// YYYY-MM-DD
		t, err := time.Parse("2006-01-02", ts)
		if err != nil {
			formatted = string(dt)
			return base.Markdown(formatted)
		}
		formatted = t.Format("Jan 01 2006")
		return base.Markdown(formatted)
	}
	if len(dt) == 25 {
		// YYYY-MM-DDThh:mm:ss+zz:zz e.g "2006-01-02T15:04:05+07:00"
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			formatted = string(dt)
			return base.Markdown(formatted)
		}
		formatted = t.Format(time.RFC1123)
		return base.Markdown(formatted)
	}
	return base.Markdown(formatted)
}

// HumanInstant renders a FHIR instant in a human friendly way
func HumanInstant(i base.Instant) base.Markdown {
	t, err := time.Parse(time.RFC3339Nano, string(i))
	if err != nil {
		return base.Markdown(string(i))
	}
	formatted := t.Format(time.RFC1123)
	return base.Markdown(formatted)
}

// HumanNarrative renders a FHIR narrative in a human friendly way.
//
// The `div` (xhtml) is converted to text or markdown. The `status` is ignored.
func HumanNarrative(n FHIRNarrative) base.Markdown {
	rawHTML := string(n.Div)
	text, err := html2text.FromString(rawHTML, html2text.Options{PrettyTables: true})
	if err != nil {
		return emptyMarkdown
	}
	return base.Markdown(text)
}

// HumanQuantity renders a FHIR quantity in a human friendly way.
//
// All fields apart from `system` (a `uri`) and `code` (a `code`) are used.
//
// - `value`, a `decimal`
// - `comparator`, a `code` (if defined)
// - `unit`, a `string`
func HumanQuantity(q FHIRQuantity) base.Markdown {
	comparator := ""
	if q.Comparator != nil {
		comparator = q.Comparator.String() + " "
	}
	qty := fmt.Sprintf("%s%.4f %s", comparator, q.Value, q.Unit)
	return base.Markdown(qty)
}

// HumanRange renders a FHIR range in a human friendly way.
//
// Both fields - `low` and `high` are rendered as quantities
func HumanRange(q FHIRRange) base.Markdown {
	formatted := fmt.Sprintf("%s - %s", HumanQuantity(q.Low), HumanQuantity(q.High))
	return base.Markdown(formatted)
}

// HumanRatio renders a FHIR ratio in a human friendly way.
//
// Both fields - `numerator` and `denominator` are rendered as quantities.
func HumanRatio(r FHIRRatio) base.Markdown {
	num := strings.TrimSpace(string(HumanQuantity(r.Numerator)))
	denom := strings.TrimSpace(string(HumanQuantity(r.Denominator)))
	formatted := fmt.Sprintf("%s / %s", num, denom)
	return base.Markdown(formatted)
}

// HumanTiming renders a FHIR timing in a human friendly way.
//
// The only field that is rendered in the summary is `code` - a `code`
// concept.
func HumanTiming(t FHIRTiming) base.Markdown {
	return base.Markdown(string(t.Code))
}

// HumanPeriod is used to render a FHIR period in a human friendly way.
//
// Both `start` and `end` are rendered as `dateTime`s.
func HumanPeriod(p FHIRPeriod) base.Markdown {
	formatted := fmt.Sprintf("%s - %s", HumanDateTime(p.Start), HumanDateTime(p.End))
	return base.Markdown(formatted)
}

// FHIRCodingSummary summarizes a FHIR coding in a human friendly way.
//
// For the human friendly display, the following fields are rendered:
//
//  - `code`, a `code`
//  - `display`, a `string`
//
// The following fields are ignored, despite being marked as summary fields by
// FHIR R4:
//
//  - `system` - a `uri`
//  - `version` - a `string`
//  - `userSelected` - a `boolean`
//
// The rendering will follow the format below:
//
//     {{display}} _({{code}})
//
// e.g Acute nasopharyngitis [common cold] _(J00)_
func FHIRCodingSummary(coding FHIRCoding) base.Markdown {
	formatted := fmt.Sprintf("%s _%s_", coding.Display, coding.Code)
	return base.Markdown(formatted)
}

// FHIRCodeableConceptSummary summarizes a FHIR codeable concept in a human friendly way
//
// For the human friendly summary, ONE of the following fields will be rendered:
//
//  - `text`, a `string`
//  - `coding`, a `Coding`
//
// If the `coding` is set, the coding is rendered. Otherwise, the text is rendered.
func FHIRCodeableConceptSummary(concept FHIRCodeableConcept) base.Markdown {
	if concept.Text != "" {
		return base.Markdown(concept.Text)
	}
	if len(concept.Coding) > 0 && concept.Coding[0] != nil {
		return FHIRCodingSummary(*concept.Coding[0])
	}
	return emptyMarkdown
}

// FHIRIdentifierSummary summarizes a FHIR identifier in a human friendly way.
//
// For the human friendly summary, the following fields will be rendered:
//
//  - `use` - a `code`
//  - `type` - a `FHIRCodeableConcept`
//  - `value` - a `string`
//
// The following fields will be omitted, intentionally:
//
//  - `system` - a `uri`
//  - `period` - a `Period`
//  - `assigner` - a `Reference(Organization)`
//
// The rendering follows the pattern below:
//
//     {{value}} _({{type.text}} - {{use}})_
// e.g
//     24291415 _(National ID - official)_
func FHIRIdentifierSummary(identifier FHIRIdentifier) base.Markdown {
	formatted := fmt.Sprintf(
		"%s _(%s) - %s_", identifier.Value, identifier.Type.Text, identifier.Use)
	return base.Markdown(formatted)
}

// FHIRReferenceSummary summarizes a FHIR reference in a human friendly way.
//
// The human display summary only renders the `display` field.
//
// The following fields are intentionally omitted from the display, despite being
// specified as part of the summary by FHIR R4:
//
// - `reference` - a `string`
// - `type` - a `uri`
// - `identifier` - an `Identifier`
//
// For that reason, this summary method does not need to know the referenced
// resource type.
func FHIRReferenceSummary(reference FHIRReference) base.Markdown {
	return base.Markdown(reference.Display)
}

// FHIREncounterSummary summarizes the key details of an encounter in a human friendly way.
//
// The summary includes the following fields:
//
//  - `identifier` - a `FHIRIdentifier`
//  - `status` - a `code`
//  - `class` - a `FHIRCoding`
//
// The following fields are left out from this summary, despite being eligible
// for the summary as per FHIR R4:
//  - identifier -> a FHIRIdentifier
//  - type -> CodeableConcept
//  - serviceType -> CodeableConcept
//  - subject -> Reference(Patient|Group)
//  - episodeOfCare -> Reference(EpisodeOfCare)
//  - participant -> embedded resource with:
//      type -> CodeableConcept
//      individual -> Reference(Practitioner | PractitionerRole | RelatedPerson)
//  - appointment -> 	Reference(Appointment)
//  - reasonCode -> CodeableConcept
//  - reasonReference -> Reference(Condition | Procedure | Observation | ImmunizationRecommendation)
//  - diagnosis -> Embedded resource with:
//    - condition -> Reference(Condition | Procedure)
func FHIREncounterSummary(resource *FHIREncounter) base.Markdown {
	if resource == nil {
		return emptyMarkdown
	}
	if resource.Identifier == nil || len(resource.Identifier) == 0 || resource.Identifier[0] == nil {
		return emptyMarkdown
	}
	encSumm := fmt.Sprintf(
		"Encounter: %s _(%s)_ |  _(%s)_",
		FHIRIdentifierSummary(*resource.Identifier[0]),
		resource.Status.String(),
		FHIRCodingSummary(resource.Class),
	)
	return base.Markdown(encSumm)
}

// FHIRAllergyIntoleranceSummary returns a human friendly summary of a FHIR allergy intolerance resource.
//
// The summary includes the following fields:
//
//  - `identifier` - a `FHIRIdentifier`
//  - `type` - a `code`
//  - `category` - a `code`
//  - `code` - a `FHIRCodeableConcept`
//  - `asserter` - a `FHIRReference`
//
// The following fields are excluded, despite being eligible for a summary as per
// FHIR R4:
//
// - `patient` - a reference to a `FHIRPatient`, omitted because the patient is
// "known" in any context where this summary should be displayed
// - `clinicalStatus`, a `FHIRCodeableConcept`; we only show allergies and
// intolerances that are in the `active` status
// - `verificationStatus`, a `FHIRCodeableConcept`; we only show allergies and
// intolerances with `verificationStatus` `confirmed`
// - `criticality`, a `code`; we only show allergies and intolerances with
// `criticality` `high`
func FHIRAllergyIntoleranceSummary(resource *FHIRAllergyIntolerance) base.Markdown {
	if resource == nil {
		return emptyMarkdown
	}
	if resource.Criticality != AllergyIntoleranceCriticalityEnumHigh {
		return emptyMarkdown
	}
	if resource.ClinicalStatus.Coding == nil || len(resource.ClinicalStatus.Coding) == 0 {
		return emptyMarkdown
	}
	if resource.ClinicalStatus.Coding[0].Code != base.Code("active") {
		return emptyMarkdown
	}
	if resource.VerificationStatus.Coding == nil || len(resource.VerificationStatus.Coding) == 0 {
		return emptyMarkdown
	}
	if resource.VerificationStatus.Coding[0].Code != base.Code("confirmed ") {
		return emptyMarkdown
	}
	if resource.Code == nil {
		return emptyMarkdown
	}
	if resource.Type == nil {
		return emptyMarkdown
	}
	if resource.Category == nil || len(resource.Category) == 0 {
		return emptyMarkdown
	}
	if resource.Recorder == nil {
		return emptyMarkdown
	}
	if resource.RecordedDate == nil {
		return emptyMarkdown
	}
	summary := fmt.Sprintf(`
### %s 
Status: (%s)

_%s allergy, reported by %s at %s_
	`,
		FHIRCodeableConceptSummary(*resource.Code),
		resource.Type,
		resource.Category,
		FHIRReferenceSummary(*resource.Recorder),
		resource.RecordedDate,
	)
	return base.Markdown(summary)
}

// FHIRConditionSummary returns a human friendly summary of a FHIR condition resource.
//
// The following fields are included in the summary:
//
//  - `identifier`, a `FHIRIdentifier`
//  - `clinicalStatus`, a `FHIRCodeableConcept`
//  - `code`, a `FHIRCodeableConcept`
//  - `onsetDateTime` (a `dateTime`) or `onsetString`, whichever is set
//  - `recordedDate`, a `dateTime`
//  - `recorder`, a `FHIRReference`
//  - `asserter`, a `FHIRReference`
//
// The following fields are omitted from the human friendly condition summary,
// despite being eligible under FHIR R4:
//
//  - `identifier`
//  - `onsetAge`, `onsetPeriod`, `onsetRange`, `onsetDatetime` and `onsetString`
//  - `bodySite`, a FHIRCodeableConcept
//  - `subject`; this summary is intended to display within the context of a
//    patient
//  - `encounter`; this summary displays within the context of an encounter
//  - `asserter`, a reference
func FHIRConditionSummary(resource *FHIRCondition) base.Markdown {
	if resource == nil {
		return emptyMarkdown
	}
	if resource.VerificationStatus == nil || resource.VerificationStatus.Coding == nil || len(resource.VerificationStatus.Coding) == 0 {
		return emptyMarkdown
	}
	if resource.VerificationStatus.Coding[0].Code != base.Code("confirmed ") {
		return emptyMarkdown
	}
	if resource.Code == nil {
		return emptyMarkdown
	}
	if resource.ClinicalStatus == nil || resource.ClinicalStatus.Coding == nil || len(resource.ClinicalStatus.Coding) == 0 {
		return emptyMarkdown
	}
	if resource.Recorder == nil {
		return emptyMarkdown
	}
	if resource.RecordedDate == nil {
		return emptyMarkdown
	}
	summary := fmt.Sprintf(`
### %s 
Status: _%s_

_Recorded by %s at %s_
	`,
		FHIRCodeableConceptSummary(*resource.Code),
		FHIRCodeableConceptSummary(*resource.ClinicalStatus),
		FHIRReferenceSummary(*resource.Recorder),
		resource.RecordedDate,
	)
	return base.Markdown(summary)
}

// FHIRServiceRequestSummary returns a human friendly summary of a FHIR service request resource.
//
// The following fields are included in the summary:
//
//  - `identifier`, a `FHIRIdentifier`
//  - `status`, a `code`
//  - `intent`, a `code`
//  - `category`, a `FHIRCodeableConcept`
//  - `priority`, a `code`
//  - `code`, a `FHIRCodeableConcept`
//  - `occurrenceDatetime`, a `dateTime`
//  - `asNeededBoolean`
//  - `asNeededCodeableConcept`
//  - `authoredOn`, a `dateTime`
//  - `requester`, a `FHIRReference`
//  - `reasonCode`, a `FHIRCodeableConcept`
//  - `patientInstruction`, a `string`
//
// The following fields have been omitted from the summary, even though they
// are required for the summary by FHIR R4:
//
//  - `identifier`
//  - `instantiatesCanonical`
//  - `instantiatesUri`
//  - `doNotPerform`, used as a filter
//  - `encounter`; this summary is displayed within the context of an encounter
//  - subject -> Reference(Patient | Group | Location | Device)
// - occurrencePeriod	-> Period
// - occurrenceTiming	-> Timing
// - `performerType`, a `FHIRCodeableConcept`
// - `performer`, a `FHIRReference`
// - `locationCode`, a `FHIRCodeableConcept`
// - `locationReference`, a `FHIRReference`
// - `reasonReference`, a `FHIRReference`
// - `specimen`, a `FHIRReference`
// - `bodySite`, a `FHIRCodeableConcept`
//
// The following fields may be added to the summary in future:
//
// - category, a list of FHIR CodeableConcepts
// - basedOn -> Reference(CarePlan | ServiceRequest | MedicationRequest)
// - replaces -> Reference(ServiceRequest)
// - requisition -> Identifier
// - orderDetail -> CodeableConcept
// - quantityQuantity	-> Quantity
// - quantityRatio -> Ratio
// - quantityRange -> Range
// - intent -> code
// priority -> code
// asNeededBoolean -> boolean
// asNeededCodeableConcept -> CodeableConcept
// reasonCode -> CodeableConcept
// reasonReference -> Reference(Condition | Observation | DiagnosticReport | DocumentReference)
// specimen -> Reference(Specimen)
// bodySite -> CodeableConcept
// patientInstruction -> string
func FHIRServiceRequestSummary(resource *FHIRServiceRequest) base.Markdown {
	if resource == nil {
		return emptyMarkdown
	}
	if resource.DoNotPerform != nil && !*resource.DoNotPerform {
		return emptyMarkdown
	}
	if resource.AuthoredOn == nil {
		return emptyMarkdown
	}
	if resource.Code == nil {
		return emptyMarkdown
	}
	if resource.Status == nil {
		return emptyMarkdown
	}
	if resource.Requester == nil {
		return emptyMarkdown
	}
	if resource.AuthoredOn == nil {
		return emptyMarkdown
	}
	summary := fmt.Sprintf(`
### %s 
Status: _%s_

_Requested by %s at %s_
	`,
		FHIRCodeableConceptSummary(*resource.Code),
		*resource.Status,
		FHIRReferenceSummary(*resource.Requester),
		HumanDateTime(*resource.AuthoredOn),
	)
	return base.Markdown(summary)
}

// FHIRMedicationRequestSummary returns a human friendly summary of a FHIR medication request resource
//
// The following fields are included in the summary:
//
//  - `identifier`, a `FHIRIdentifier`
//  - `status`, a `code`
//  - `intent`, a `code`
//  - `priority`, a `code`
//  - `medicationCodeableConcept`, a `FHIRCodeableConcept`
//  - `authoredOn`, a `dateTime`
//  - `requester`, a `Reference(Practitioner|PractitionerRole|Organization|Patient|RelatedPerson|Device)`
//
// The following fields have been omitted from the summary, even though they
// are required for the summary by FHIR R4:
//
//  - `identifier`
//  - `intent`, a code
//  - `instantiatesCanonical`
//  - `instantiatesUri`
//  - `doNotPerform`, which is used as a filter instead
//  - `reportedBoolean` and `reportedReference`
//  - `subject`, a `FHIRReference`
//  - `performerType`, a `FHIRCodeableConcept`
//  - `basedOn`, a `Reference(CarePlan|MedicationRequest|ServiceRequest|ImmunizationRecommendation)`
//  - `groupIdentifier`, a `Reference(CarePlan|MedicationRequest|ServiceRequest|ImmunizationRecommendation)`
//
//  The following fields may be implemented at some point in the future:
//
//  - medicationReference -> Reference(Medication)
func FHIRMedicationRequestSummary(resource *FHIRMedicationRequest) base.Markdown {
	if resource == nil {
		return emptyMarkdown
	}
	if resource.DoNotPerform != nil && !*resource.DoNotPerform {
		return emptyMarkdown
	}
	if resource.MedicationCodeableConcept == nil {
		return emptyMarkdown
	}
	if resource.Status == nil {
		return emptyMarkdown
	}
	if resource.Priority == nil {
		return emptyMarkdown
	}
	if resource.Requester == nil {
		return emptyMarkdown
	}
	if resource.AuthoredOn == nil {
		return emptyMarkdown
	}
	summary := fmt.Sprintf(`
### %s 
Status: _%s_
Priority: _%s_

_Requested by %s at %s_
	`,
		FHIRCodeableConceptSummary(*resource.MedicationCodeableConcept),
		*resource.Status,
		*resource.Priority,
		FHIRReferenceSummary(*resource.Requester),
		HumanDateTime(*resource.AuthoredOn),
	)
	return base.Markdown(summary)
}

// FHIRCompositionSummary returns a human friendly summary of a FHIR composition resource.
//
// The following fields are included in a composition summary:
//
//  - `identifier`, a `FHIRIdentifier`
//  - `title`, a `string`
//  - `status`, a `code`
//  - `date`, a `dateTime`
//  - `author`, a `Reference(Practitioner|PractitionerRole|Device|Patient|RelatedPerson|Organization)`
//  - `confidentiality`, a `FHIRCodeableConcept`
//  - `section`, with the `text` in each section summarized
//
// The following fields are excluded:
//
//  - `identifier`
//  - `subject`, a `Reference(Any)`
//  - `type`, a `FHIRCodeableConcept`
//  - `encounter`, a `Reference(Encounter)`
//  - `custodian`, a `Reference(Organization)`
//  - `category`, a `FHIRCodeableConcept`
//  - `event`, an embedded structure that includes a `code`, `period` and `detail` ref
func FHIRCompositionSummary(resource *FHIRComposition) base.Markdown {
	if resource == nil {
		return emptyMarkdown
	}
	if resource.Status != nil && *resource.Status != CompositionStatusEnumEnteredInError {
		return emptyMarkdown
	}
	if resource.Title == nil {
		return emptyMarkdown
	}
	if resource.Confidentiality == nil {
		return emptyMarkdown
	}
	if resource.Status == nil {
		return emptyMarkdown
	}
	if resource.Date == nil {
		return emptyMarkdown
	}
	if resource.Section == nil || len(resource.Section) == 0 {
		return emptyMarkdown
	}

	sectionSummaries := []base.Markdown{}
	for _, section := range resource.Section {
		if section.Text == nil {
			continue
		}
		sectionSummaries = append(sectionSummaries, HumanNarrative(*section.Text))
	}

	compositionSummary := fmt.Sprintf(`
### %s (_%s_)
Status: _%s_
Date: %s

Notes:

%s
`,
		*resource.Title,
		*resource.Confidentiality,
		resource.Status,
		resource.Date,
		sectionSummaries,
	)
	return base.Markdown(compositionSummary)
}

// FHIRObservationSummary returns a human friendly summary of a FHIR observation resource.
//
// The following fields are included in the summary:
//
//  - `status`, a `code`
//  - `category`, a `FHIRCodeableConcept`
//  - `code`, a `FHIRCodeableConcept`
//  - `effectiveDatetime` or `effectiveInstant`, whichever is set
//  - `issued`, an `instant`
//  - value, as one of `valueQuantity`/`valueCodeableConcept`/`valueString`/`valueBoolean`
//    `valueInteger`/`valueRange`/`valueRatio`/`valueTime`/`valueDateTime`/`valuePeriod`
//  - `component`, an embedded structure that includes a `code` and a `value`
//    that is rendered as above
//
// The following fields are not included in the summary, despite being eligible
// under FHIR R4:
//
//  - `identifier`, a `FHIRIdentifier`
//  - `effectivePeriod` and `effectiveTiming` (`effectiveDate` and `effectiveInstant`)
//    are supported
//  - `valueSampledData` (both for the observation and it's components)
//  - `basedOn` and `partOf`, both `FHIRReference`s
//  - `subject`, since this summary is viewed within a patient context
//  - `focus`
//  - `encounter`, since this summary is viewed within an encounter context
//  - `performer`, a `FHIRReference`; this detail is deemed to be excessive for
//  a summary that is aimed at doctors
//
// The following advanced fields might be included at a later date:
//
//  - `hasMember`, a Reference(Observation | QuestionnaireResponse | MolecularSequence)
//  - `derivedFrom`, a Reference(DocumentReference | ImagingStudy | Media | QuestionnaireResponse | Observation | MolecularSequence)
func FHIRObservationSummary(resource *FHIRObservation) base.Markdown {
	if resource == nil {
		return emptyMarkdown
	}
	if resource.Status == nil {
		return emptyMarkdown
	}

	summary := fmt.Sprintf(`
### %s 
Status: _%s_
Value(s): %s

As at %s

%s
	`,
		FHIRCodeableConceptSummary(resource.Code),
		*resource.Status,
		summarizeObservationValue(*resource),
		summarizeObservationTime(*resource),
		summarizeObservationComponents(*resource),
	)
	return base.Markdown(summary)
}

func summarizeObservationComponents(obs FHIRObservation) base.Markdown {
	if obs.Component == nil {
		return emptyMarkdown
	}

	summary := ""
	for _, component := range obs.Component {
		if component == nil {
			continue
		}
		summary += fmt.Sprintf(" - %s: %s", FHIRCodeableConceptSummary(component.Code), summarizeComponentValue(*component))
	}
	return base.Markdown(summary)
}

func summarizeObservationTime(obs FHIRObservation) base.Markdown {
	if obs.EffectiveInstant != nil {
		return HumanInstant(*obs.EffectiveInstant)
	}
	if obs.EffectiveDateTime != nil {
		return base.Markdown(obs.EffectiveDateTime.String())
	}
	if obs.Issued != nil {
		return HumanInstant(*obs.Issued)
	}
	return base.Markdown("-")
}

func summarizeObservationValue(obs FHIRObservation) base.Markdown {
	val := "unknown"
	if obs.ValueCodeableConcept != nil {
		return base.Markdown(*obs.ValueCodeableConcept)
	}
	if obs.ValueString != nil {
		return base.Markdown(*obs.ValueString)
	}
	if obs.ValueInteger != nil {
		return base.Markdown(*obs.ValueInteger)
	}
	if obs.ValueBoolean != nil {
		return base.Markdown(strconv.FormatBool(*obs.ValueBoolean))
	}
	if obs.ValueDateTime != nil {
		return base.Markdown(obs.ValueDateTime.String())
	}
	if obs.ValueQuantity != nil {
		return base.Markdown(HumanQuantity(*obs.ValueQuantity))
	}
	if obs.ValueRange != nil {
		return base.Markdown(HumanRange(*obs.ValueRange))
	}
	if obs.ValueRatio != nil {
		return base.Markdown(HumanRatio(*obs.ValueRatio))
	}
	if obs.ValueTime != nil {
		return base.Markdown(obs.ValueTime.Format(time.Kitchen))
	}
	if obs.ValuePeriod != nil {
		return base.Markdown(HumanPeriod(*obs.ValuePeriod))
	}
	return base.Markdown(val)
}

func summarizeComponentValue(cmpn FHIRObservationComponent) base.Markdown {
	val := "unknown"
	if cmpn.ValueCodeableConcept != nil {
		return base.Markdown(*cmpn.ValueCodeableConcept)
	}
	if cmpn.ValueString != nil {
		return base.Markdown(*cmpn.ValueString)
	}
	if cmpn.ValueInteger != nil {
		return base.Markdown(*cmpn.ValueInteger)
	}
	if cmpn.ValueBoolean != nil {
		return base.Markdown(strconv.FormatBool(*cmpn.ValueBoolean))
	}
	if cmpn.ValueDateTime != nil {
		return base.Markdown(cmpn.ValueDateTime.String())
	}
	if cmpn.ValueQuantity != nil {
		return base.Markdown(HumanQuantity(*cmpn.ValueQuantity))
	}
	if cmpn.ValueRange != nil {
		return base.Markdown(HumanRange(*cmpn.ValueRange))
	}
	if cmpn.ValueRatio != nil {
		return base.Markdown(HumanRatio(*cmpn.ValueRatio))
	}
	if cmpn.ValueTime != nil {
		return base.Markdown(cmpn.ValueTime.Format(time.Kitchen))
	}
	if cmpn.ValuePeriod != nil {
		return base.Markdown(HumanPeriod(*cmpn.ValuePeriod))
	}
	return base.Markdown(val)
}
