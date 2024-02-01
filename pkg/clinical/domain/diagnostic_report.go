package domain

// DiagnosticReport is documented here http://hl7.org/fhir/StructureDefinition/DiagnosticReport
type DiagnosticReport struct {
	ID                 *string                    `json:"id,omitempty"`
	Meta               *FHIRMeta                  `json:"meta,omitempty"`
	ImplicitRules      *string                    `json:"implicitRules,omitempty"`
	Language           *string                    `json:"language,omitempty"`
	Text               *FHIRNarrative             `json:"text,omitempty"`
	Extension          []*Extension               `json:"extension,omitempty"`
	ModifierExtension  []*Extension               `json:"modifierExtension,omitempty"`
	Identifier         []*FHIRIdentifier          `json:"identifier,omitempty"`
	BasedOn            []*FHIRReference           `json:"basedOn,omitempty"`
	Status             DiagnosticReportStatusEnum `json:"status"`
	Category           []*FHIRCodeableConcept     `json:"category,omitempty"`
	Code               FHIRCodeableConcept        `json:"code"`
	Subject            *FHIRReference             `json:"subject,omitempty"`
	Encounter          *FHIRReference             `json:"encounter,omitempty"`
	EffectiveDateTime  *string                    `json:"effectiveDateTime,omitempty"`
	EffectivePeriod    *FHIRPeriod                `json:"effectivePeriod,omitempty"`
	Issued             *string                    `json:"issued,omitempty"`
	Performer          []*FHIRReference           `json:"performer,omitempty"`
	ResultsInterpreter []*FHIRReference           `json:"resultsInterpreter,omitempty"`
	Specimen           []*FHIRReference           `json:"specimen,omitempty"`
	Result             []*FHIRReference           `json:"result,omitempty"`
	ImagingStudy       []*FHIRReference           `json:"imagingStudy,omitempty"`
	Media              []*DiagnosticReportMedia   `json:"media,omitempty"`
	Conclusion         *string                    `json:"conclusion,omitempty"`
	ConclusionCode     []*FHIRCodeableConcept     `json:"conclusionCode,omitempty"`
	PresentedForm      []*FHIRAttachment          `json:"presentedForm,omitempty"`
}

// DiagnosticReportMedia represents the key images associated with this report
type DiagnosticReportMedia struct {
	ID                *string        `json:"id,omitempty"`
	Extension         []*Extension   `json:"extension,omitempty"`
	ModifierExtension []*Extension   `json:"modifierExtension,omitempty"`
	Comment           *string        `json:"comment,omitempty"`
	Link              *FHIRReference `json:"link"`
}
