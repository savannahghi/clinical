package clinical

import (
	"gitlab.slade360emr.com/go/base"
)

// FHIRPatient definition: demographics and other administrative information about an individual or animal receiving care or other health-related services.
type FHIRPatient struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// An identifier for this patient.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// Whether this patient record is in active use.
	// Many systems use this property to mark as non-current patients, such as those that have not been seen for a period of time based on an organization's business rules.
	//
	// It is often used to filter patient lists to exclude inactive patients
	//
	// Deceased patients may also be marked as inactive for the same reasons, but may be active for some time after death.
	Active *bool `json:"active,omitempty"`

	// A name associated with the individual.
	Name []*FHIRHumanName `json:"name,omitempty"`

	// A contact detail (e.g. a telephone number or an email address) by which the individual may be contacted.
	Telecom []*FHIRContactPoint `json:"telecom,omitempty"`

	// Administrative Gender - the gender that the patient is considered to have for administration and record keeping purposes.
	Gender *PatientGenderEnum `json:"gender,omitempty"`

	// The date of birth for the individual.
	BirthDate *base.Date `json:"birthDate,omitempty"`

	// Indicates if the individual is deceased or not.
	DeceasedBoolean *bool `json:"deceasedBoolean,omitempty"`

	// Indicates if the individual is deceased or not.
	DeceasedDateTime *base.Date `json:"deceasedDateTime,omitempty"`

	// An address for the individual.
	Address []*FHIRAddress `json:"address,omitempty"`

	// This field contains a patient's most recent marital (civil) status.
	MaritalStatus *FHIRCodeableConcept `json:"maritalStatus,omitempty"`

	// Indicates whether the patient is part of a multiple (boolean) or indicates the actual birth order (integer).
	MultipleBirthBoolean *bool `json:"multipleBirthBoolean,omitempty"`

	// Indicates whether the patient is part of a multiple (boolean) or indicates the actual birth order (integer).
	MultipleBirthInteger *string `json:"multipleBirthInteger,omitempty"`

	// Image of the patient.
	Photo []*FHIRAttachment `json:"photo,omitempty"`

	// A contact party (e.g. guardian, partner, friend) for the patient.
	Contact []*FHIRPatientContact `json:"contact,omitempty"`

	// A language which may be used to communicate with the patient about his or her health.
	Communication []*FHIRPatientCommunication `json:"communication,omitempty"`

	// Patient's nominated care provider.
	GeneralPractitioner []*FHIRReference `json:"generalPractitioner,omitempty"`

	// Organization that is the custodian of the patient record.
	ManagingOrganization *FHIRReference `json:"managingOrganization,omitempty"`

	// Link to another patient resource that concerns the same actual patient.
	Link []*FHIRPatientLink `json:"link,omitempty"`
}

// FHIRPatientCommunication definition: demographics and other administrative information about an individual or animal receiving care or other health-related services.
type FHIRPatientCommunication struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The ISO-639-1 alpha 2 code in lower case for the language, optionally followed by a hyphen and the ISO-3166-1 alpha 2 code for the region in upper case; e.g. "en" for English, or "en-US" for American English versus "en-EN" for England English.
	Language *FHIRCodeableConcept `json:"language,omitempty"`

	// Indicates whether or not the patient prefers this language (over other languages he masters up a certain level).
	Preferred *bool `json:"preferred,omitempty"`
}

// FHIRPatientCommunicationInput is the input type for PatientCommunication
type FHIRPatientCommunicationInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The ISO-639-1 alpha 2 code in lower case for the language, optionally followed by a hyphen and the ISO-3166-1 alpha 2 code for the region in upper case; e.g. "en" for English, or "en-US" for American English versus "en-EN" for England English.
	Language *FHIRCodeableConceptInput `json:"language,omitempty"`

	// Indicates whether or not the patient prefers this language (over other languages he masters up a certain level).
	Preferred *bool `json:"preferred,omitempty"`
}

// FHIRPatientContact definition: demographics and other administrative information about an individual or animal receiving care or other health-related services.
type FHIRPatientContact struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The nature of the relationship between the patient and the contact person.
	Relationship []*FHIRCodeableConcept `json:"relationship,omitempty"`

	// A name associated with the contact person.
	Name *FHIRHumanName `json:"name,omitempty"`

	// A contact detail for the person, e.g. a telephone number or an email address.
	Telecom []*FHIRContactPoint `json:"telecom,omitempty"`

	// Address for the contact person.
	Address *FHIRAddress `json:"address,omitempty"`

	// Administrative Gender - the gender that the contact person is considered to have for administration and record keeping purposes.
	Gender *PatientContactGenderEnum `json:"gender,omitempty"`

	// Organization on behalf of which the contact is acting or for which the contact is working.
	Organization *FHIRReference `json:"organization,omitempty"`

	// The period during which this contact person or organization is valid to be contacted relating to this patient.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIRPatientContactInput is the input type for PatientContact
type FHIRPatientContactInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The nature of the relationship between the patient and the contact person.
	Relationship []*FHIRCodeableConceptInput `json:"relationship,omitempty"`

	// A name associated with the contact person.
	Name *FHIRHumanNameInput `json:"name,omitempty"`

	// A contact detail for the person, e.g. a telephone number or an email address.
	Telecom []*FHIRContactPointInput `json:"telecom,omitempty"`

	// Address for the contact person.
	Address *FHIRAddressInput `json:"address,omitempty"`

	// Administrative Gender - the gender that the contact person is considered to have for administration and record keeping purposes.
	Gender *PatientContactGenderEnum `json:"gender,omitempty"`

	// Organization on behalf of which the contact is acting or for which the contact is working.
	Organization *FHIRReferenceInput `json:"organization,omitempty"`

	// The period during which this contact person or organization is valid to be contacted relating to this patient.
	Period *FHIRPeriodInput `json:"period,omitempty"`
}

// FHIRPatientInput is the input type for Patient
type FHIRPatientInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// An identifier for this patient.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	//     Whether this patient record is in active use.
	// Many systems use this property to mark as non-current patients, such as those that have not been seen for a period of time based on an organization's business rules.
	//
	// It is often used to filter patient lists to exclude inactive patients
	//
	// Deceased patients may also be marked as inactive for the same reasons, but may be active for some time after death.
	Active *bool `json:"active,omitempty"`

	// A name associated with the individual.
	Name []*FHIRHumanNameInput `json:"name,omitempty"`

	// A contact detail (e.g. a telephone number or an email address) by which the individual may be contacted.
	Telecom []*FHIRContactPointInput `json:"telecom,omitempty"`

	// Administrative Gender - the gender that the patient is considered to have for administration and record keeping purposes.
	Gender *PatientGenderEnum `json:"gender,omitempty"`

	// The date of birth for the individual.
	BirthDate *base.Date `json:"birthDate,omitempty"`

	// Indicates if the individual is deceased or not.
	DeceasedBoolean *bool `json:"deceasedBoolean,omitempty"`

	// Indicates if the individual is deceased or not.
	DeceasedDateTime *base.Date `json:"deceasedDateTime,omitempty"`

	// An address for the individual.
	Address []*FHIRAddressInput `json:"address,omitempty"`

	// This field contains a patient's most recent marital (civil) status.
	MaritalStatus *FHIRCodeableConceptInput `json:"maritalStatus,omitempty"`

	// Indicates whether the patient is part of a multiple (boolean) or indicates the actual birth order (integer).
	MultipleBirthBoolean *bool `json:"multipleBirthBoolean,omitempty"`

	// Indicates whether the patient is part of a multiple (boolean) or indicates the actual birth order (integer).
	MultipleBirthInteger *string `json:"multipleBirthInteger,omitempty"`

	// Image of the patient.
	Photo []*FHIRAttachmentInput `json:"photo,omitempty"`

	// A contact party (e.g. guardian, partner, friend) for the patient.
	Contact []*FHIRPatientContactInput `json:"contact,omitempty"`

	// A language which may be used to communicate with the patient about his or her health.
	Communication []*FHIRPatientCommunicationInput `json:"communication,omitempty"`

	// Patient's nominated care provider.
	GeneralPractitioner []*FHIRReferenceInput `json:"generalPractitioner,omitempty"`

	// Organization that is the custodian of the patient record.
	ManagingOrganization *FHIRReferenceInput `json:"managingOrganization,omitempty"`

	// Link to another patient resource that concerns the same actual patient.
	Link []*FHIRPatientLinkInput `json:"link,omitempty"`
}

// FHIRPatientLink definition: demographics and other administrative information about an individual or animal receiving care or other health-related services.
type FHIRPatientLink struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The other patient resource that the link refers to.
	Other *FHIRReference `json:"other,omitempty"`

	// The type of link between this patient resource and another patient resource.
	Type *PatientLinkTypeEnum `json:"type,omitempty"`
}

// FHIRPatientLinkInput is the input type for PatientLink
type FHIRPatientLinkInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// The other patient resource that the link refers to.
	Other *FHIRReferenceInput `json:"other,omitempty"`

	// The type of link between this patient resource and another patient resource.
	Type *PatientLinkTypeEnum `json:"type,omitempty"`
}

// FHIRPatientRelayConnection is a Relay connection for Patient
type FHIRPatientRelayConnection struct {
	Edges           []*FHIRPatientRelayEdge `json:"edges,omitempty"`
	HasOpenEpisodes bool                    `json:"hasOpenEpisodes,omitempty"`
	PageInfo        *base.PageInfo          `json:"pageInfo,omitempty"`
}

// FHIRPatientRelayEdge is a Relay edge for Patient
type FHIRPatientRelayEdge struct {
	Cursor          *string      `json:"cursor,omitempty"`
	HasOpenEpisodes bool         `json:"hasOpenEpisodes,omitempty"`
	Node            *FHIRPatient `json:"node,omitempty"`
}

// FHIRPatientRelayPayload is used to return single instances of Patient
type FHIRPatientRelayPayload struct {
	Resource        *FHIRPatient `json:"resource,omitempty"`
	HasOpenEpisodes bool         `json:"hasOpenEpisodes,omitempty"`
}

// PatientEdge is a Relay style edge for listings of FHIR patient records.
type PatientEdge struct {
	Cursor          string       `json:"cursor"`
	Node            *FHIRPatient `json:"node"`
	HasOpenEpisodes bool         `json:"hasOpenEpisodes"`
}

// PatientConnection is a Relay style connection for use in listings of FHIR
// patient records.
type PatientConnection struct {
	Edges    []*PatientEdge `json:"edges"`
	PageInfo *base.PageInfo `json:"pageInfo"`
}
