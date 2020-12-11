package clinical

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"gitlab.slade360emr.com/go/base"
)

// Dummy ..
type Dummy struct {
	id string
}

//IsEntity ...
func (d Dummy) IsEntity() {}

// IsNode ..
func (d *Dummy) IsNode() {}

// SetID sets the trace's ID
func (d *Dummy) SetID(id string) {
	d.id = id
}

// CodeableConcept description from FHIR: A concept that may be defined by a
// formal reference to a terminology or ontology or may be provided by text.
type CodeableConcept struct {
	ID     string    `json:"id,omitempty"`
	Coding []*Coding `json:"coding,omitempty"`
	Text   string    `json:"text,omitempty"`
}

//IsEntity ...
func (c CodeableConcept) IsEntity() {}

// CodeableConceptInput is used to create codeable concepts.
type CodeableConceptInput struct {
	ID     *string        `json:"id"`
	Coding []*CodingInput `json:"coding"`
	Text   string         `json:"text"`
}

//IsEntity ...
func (c CodeableConceptInput) IsEntity() {}

// Coding description from FHIR: A reference to a code defined by a terminology system.
type Coding struct {
	System       string  `json:"system"`
	Version      *string `json:"version"`
	Code         string  `json:"code"`
	Display      *string `json:"display"`
	UserSelected *bool   `json:"userSelected"`
}

//IsEntity ...
func (c Coding) IsEntity() {}

// CodingInput is used to set coding.
type CodingInput struct {
	System       string `json:"system"`
	Version      string `json:"version"`
	Code         string `json:"code"`
	Display      string `json:"display"`
	UserSelected bool   `json:"userSelected"`
}

// Communication description from FHIR: Information about a person
// that is involved in the care for a patient, but who is not the target of
// healthcare, nor has a formal responsibility in the care process.
type Communication struct {
	Language  *CodeableConcept `json:"language"`
	Preferred bool             `json:"preferred"`
}

//IsEntity ...
func (c Communication) IsEntity() {}

// CommunicationInput is used to create send communication data in GraphQL.
type CommunicationInput struct {
	Language  *CodeableConceptInput `json:"language"`
	Preferred bool                  `json:"preferred"`
}

// ContactPoint description from FHIR: Details for all kinds of technology
// mediated contact points for a person or organization, including telephone,
// email, etc.
type ContactPoint struct {
	ID     string             `json:"id"`
	System ContactPointSystem `json:"system"`
	Value  string             `json:"value"`
	Use    ContactPointUse    `json:"use"`
	Rank   int64              `json:"rank"`
	Period *Period            `json:"period"`
}

//IsEntity ...
func (c ContactPoint) IsEntity() {}

// ContactPointInput is used to input contact details e.g phone, email etc.
type ContactPointInput struct {
	ID     *string            `json:"id"`
	System ContactPointSystem `json:"system"`
	Value  string             `json:"value"`
	Use    ContactPointUse    `json:"use"`
	Rank   int64              `json:"rank"`
	Period *FHIRPeriodInput   `json:"period"`
}

// HumanName description from FHIR: A human's name with the ability to identify
// parts and usage.
type HumanName struct {
	ID     string    `json:"id"`
	Use    NameUse   `json:"use"`
	Text   *string   `json:"text"`
	Family string    `json:"family"`
	Given  []string  `json:"given"`
	Prefix []*string `json:"prefix"`
	Suffix []*string `json:"suffix"`
	Period *Period   `json:"period"`
}

//IsEntity ...
func (h HumanName) IsEntity() {}

// HumanNameInput is used to input patient names.
type HumanNameInput struct {
	ID     *string          `json:"id"`
	Use    NameUse          `json:"use"`
	Text   *string          `json:"text"`
	Family string           `json:"family"`
	Given  []string         `json:"given"`
	Prefix []*string        `json:"prefix"`
	Suffix []*string        `json:"suffix"`
	Period *FHIRPeriodInput `json:"period"`
}

// Address description from FHIR: An address expressed using postal conventions
// (as opposed to GPS or other location definition formats).
//
// This data type may be used to convey addresses for use in delivering mail as
// well as for visiting locations which might not be valid for mail delivery.
//
// There are a variety of postal address formats defined around the world.
type Address struct {
	ID         string      `json:"id"`
	Use        AddressUse  `json:"use"`
	Type       AddressType `json:"type"`
	Text       string      `json:"text"`
	Line       []*string   `json:"line"`
	City       *string     `json:"city"`
	District   *string     `json:"district"`
	State      *string     `json:"state"`
	PostalCode *string     `json:"postalCode"`
	Country    Country     `json:"country"`
	Period     *Period     `json:"period"`
}

//IsEntity ...
func (a Address) IsEntity() {}

// AddressInput is used to create postal and physical addresses.
//
// IMPORTANT:
//
// For physical addresses, use Google Maps co-ordinates or plus codes.
// See: https://support.google.com/maps/answer/18539?co=GENIE.Platform%3DDesktop&hl=en
type AddressInput struct {
	ID         *string          `json:"id"`
	Use        AddressUse       `json:"use"`
	Type       AddressType      `json:"type"`
	Text       string           `json:"text"`
	Line       []*string        `json:"line"`
	City       *string          `json:"city"`
	District   *string          `json:"district"`
	State      *string          `json:"state"`
	PostalCode *string          `json:"postalCode"`
	Country    Country          `json:"country"`
	Period     *FHIRPeriodInput `json:"period"`
}

// Identifier is used to represent a numeric or alphanumeric string that is
// associated with a single object or entity within a given system. Typically,
// identifiers are used to connect content in resources to external content
// available in other frameworks or protocols. Identifiers are associated with
// objects and may be changed or retired due to human or system process and
// errors.
type Identifier struct {
	ID     string           `json:"id,omitempty"`
	Use    IdentifierUse    `json:"use,omitempty"`
	Type   *CodeableConcept `json:"type,omitempty"`
	System *string          `json:"system,omitempty"`
	Value  *string          `json:"value,omitempty"`
	Period *Period          `json:"period,omitempty"`
}

//IsEntity ...
func (i Identifier) IsEntity() {}

// IdentifierInput is used to create and update identifiers.
type IdentifierInput struct {
	ID     string                `json:"id"`
	Use    IdentifierUse         `json:"use"`
	Type   *CodeableConceptInput `json:"type"`
	System *string               `json:"system"`
	Value  *string               `json:"value"`
	Period *FHIRPeriodInput      `json:"period"`
}

// Period is a FHIR https://www.hl7.org/fhir/datatypes.html#Period.
//
// A period should have a start, end or both. It's an error to have a period that
// has null start and end times.
type Period struct {
	ID    string     `json:"id,omitempty"`
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

//IsEntity ...
func (p Period) IsEntity() {}

// PeriodInput is used to set time ranges e.g validity.
//
// A period should have a start, end or both. It's an error to have a period that
// has null start and end times.
type PeriodInput struct {
	ID    *string    `json:"id"`
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

//IsEntity ...
func (p PeriodInput) IsEntity() {}

// NameInput is used to input patient names.
type NameInput struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	OtherNames *string `json:"otherNames"`
}

// IdentificationDocument is used to input e.g National ID or passport document
// numbers at patient registration.
type IdentificationDocument struct {
	DocumentType     IDDocumentType    `json:"documentType"`
	DocumentNumber   string            `json:"documentNumber"`
	Title            *string           `json:"title,omitempty"`
	ImageContentType *base.ContentType `json:"imageContentType,omitempty"`
	ImageBase64      *string           `json:"imageBase64,omitempty"`
}

// PhoneNumberInput is used to input phone numbers.
type PhoneNumberInput struct {
	Msisdn             string `json:"msisdn"`
	VerificationCode   string `json:"verificationCode"`
	IsUssd             bool   `json:"isUSSD"`
	CommunicationOptIn bool   `json:"communicationOptIn"`
}

// PhotoInput is used to upload patient photos.
type PhotoInput struct {
	PhotoContentType base.ContentType `json:"photoContentType"`
	PhotoBase64data  string           `json:"photoBase64data"`
	PhotoFilename    string           `json:"photoFilename"`
}

// EmailInput is used to register patient emails.
type EmailInput struct {
	Email              string `json:"email"`
	CommunicationOptIn bool   `json:"communicationOptIn"`
}

// PhysicalAddress is used to record a precise physical address.
type PhysicalAddress struct {
	MapsCode        string `json:"mapsCode"`
	PhysicalAddress string `json:"physicalAddress"`
}

// PostalAddress is used to record patient's postal addresses
type PostalAddress struct {
	PostalAddress string `json:"postalAddress"`
	PostalCode    string `json:"postalCode"`
}

// SimplePatientRegistrationInput provides a simplified API to support registration
// of patients.
type SimplePatientRegistrationInput struct {
	ID                      string                    `json:"id"`
	Names                   []*NameInput              `json:"names"`
	IdentificationDocuments []*IdentificationDocument `json:"identificationDocuments"`
	BirthDate               base.Date                 `json:"birthDate"`
	PhoneNumbers            []*PhoneNumberInput       `json:"phoneNumbers"`
	Photos                  []*PhotoInput             `json:"photos"`
	Emails                  []*EmailInput             `json:"emails"`
	PhysicalAddresses       []*PhysicalAddress        `json:"physicalAddresses"`
	PostalAddresses         []*PostalAddress          `json:"postalAddresses"`
	Gender                  string                    `json:"gender"`
	Active                  bool                      `json:"active"`
	MaritalStatus           MaritalStatus             `json:"maritalStatus"`
	Languages               []base.Language           `json:"languages"`
	ReplicateUSSD           bool                      `json:"replicate_ussd,omitempty"`
}

// BreakGlassEpisodeCreationInput is used to start emergency episodes via a
// break glass protocol
type BreakGlassEpisodeCreationInput struct {
	PatientID       string `json:"patientID" firestore:"patientID"`
	ProviderCode    string `json:"providerCode" firestore:"providerCode"`
	PractitionerUID string `json:"practitionerUID" firestore:"practitionerUID"`
	Msisdn          string `json:"msisdn" firestore:"msisdn"`
	PatientPhone    string `json:"patientPhone" firestore:"patientPhone"`
	Otp             string `json:"otp" firestore:"otp"`
	FullAccess      bool   `json:"fullAccess" firestore:"fullAccess"`
}

// OTPEpisodeCreationInput is used to start patient visits via OTP
type OTPEpisodeCreationInput struct {
	PatientID    string `json:"patientID"`
	ProviderCode string `json:"providerCode"`
	Msisdn       string `json:"msisdn"`
	Otp          string `json:"otp"`
	FullAccess   bool   `json:"fullAccess"`
}

// OTPEpisodeUpgradeInput is used to upgrade existing open episodes
type OTPEpisodeUpgradeInput struct {
	EpisodeID string `json:"episodeID"`
	Msisdn    string `json:"msisdn"`
	Otp       string `json:"otp"`
}

// SimpleNHIFInput adds NHIF membership details as an extra identifier.
type SimpleNHIFInput struct {
	PatientID             string            `json:"patientID"`
	MembershipNumber      string            `json:"membershipNumber"`
	FrontImageBase64      *string           `json:"frontImageBase64"`
	FrontImageContentType *base.ContentType `json:"frontImageContentType"`
	RearImageBase64       *string           `json:"rearImageBase64"`
	RearImageContentType  *base.ContentType `json:"rearImageContentType"`
}

// SimpleNextOfKinInput is used to add next of kin to a patient.
type SimpleNextOfKinInput struct {
	PatientID         string              `json:"patientID"`
	Names             []*NameInput        `json:"names"`
	PhoneNumbers      []*PhoneNumberInput `json:"phoneNumbers"`
	Emails            []*EmailInput       `json:"emails"`
	PhysicalAddresses []*PhysicalAddress  `json:"physicalAddresses"`
	PostalAddresses   []*PostalAddress    `json:"postalAddresses"`
	Gender            string              `json:"gender"`
	Relationship      RelationshipType    `json:"relationship"`
	Active            bool                `json:"active"`
	BirthDate         base.Date           `json:"birthDate"`
}

// USSDEpisodeCreationInput is used to start episodes via USSD
type USSDEpisodeCreationInput struct {
	PatientID    string `json:"patientID"`
	ProviderCode string `json:"providerCode"`
	SessionID    string `json:"sessionID"`
	Msisdn       string `json:"msisdn"`
	FullAccess   bool   `json:"fullAccess"`
}

// Reference defines references to other FHIR resources.
type Reference struct {
	Reference  string         `json:"reference,omitempty"`
	Type       string         `json:"type,omitempty"`
	Identifier FHIRIdentifier `json:"identifier,omitempty"`
	Display    *string        `json:"display,omitempty"`
}

// Names renders the patient's names as a string
func (p FHIRPatient) Names() string {
	name := ""
	if p.Name == nil {
		return name
	}

	names := []string{}
	for _, hn := range p.Name {
		if hn == nil {
			continue
		}
		if hn.Text == "" {
			continue
		}
		names = append(names, hn.Text)
	}
	name = strings.Join(names, " | ")
	return name
}

// IsEntity ...
func (p FHIRPatient) IsEntity() {}

// PatientContact is a contact party (e.g. guardian, partner, friend) for the
// patient.
type PatientContact struct {
	ID           string                 `json:"id"`
	Relationship []*FHIRCodeableConcept `json:"relationship"`
	Name         *FHIRHumanName         `json:"name"`
	Telecom      []*FHIRContactPoint    `json:"telecom"`
	Address      *FHIRAddress           `json:"address"`
	Gender       *string                `json:"gender"`
	Period       *FHIRPeriod            `json:"period"`
}

// IsEntity ...
func (p PatientContact) IsEntity() {}

// PatientContactInput is used to create and update patient contacts
type PatientContactInput struct {
	ID           *string                     `json:"id"`
	Relationship []*FHIRCodeableConceptInput `json:"relationship"`
	Name         *FHIRHumanNameInput         `json:"name"`
	Telecom      []*FHIRContactPointInput    `json:"telecom"`
	Address      *FHIRAddressInput           `json:"address"`
	Gender       *string                     `json:"gender"`
	Period       *FHIRPeriodInput            `json:"period"`
}

// PatientInput is used to create patient records.
type PatientInput struct {
	ID            *string                   `json:"id"`
	Identifier    []*FHIRIdentifierInput    `json:"identifier"`
	Active        bool                      `json:"active"`
	Name          []*FHIRHumanNameInput     `json:"name"`
	Telecom       []*FHIRContactPointInput  `json:"telecom"`
	Gender        string                    `json:"gender"`
	BirthDate     base.Date                 `json:"birthDate"`
	Address       []*FHIRAddressInput       `json:"address"`
	MaritalStatus *FHIRCodeableConceptInput `json:"maritalStatus"`
	Photo         []*FHIRAttachmentInput    `json:"photo"`
	Contact       []*PatientContactInput    `json:"contact"`
	Communication []*CommunicationInput     `json:"communication"`
}

// PatientPayload is used to return patient records and ancillary data after
// mutations.
type PatientPayload struct {
	PatientRecord   *FHIRPatient         `json:"patientRecord,omitempty"`
	HasOpenEpisodes bool                 `json:"hasOpenEpisodes,omitempty"`
	OpenEpisodes    []*FHIREpisodeOfCare `json:"openEpisodes,omitempty"`
}

// RetirePatientInput is used to retire patient records.
type RetirePatientInput struct {
	ID string `json:"id"`
}

// PatientExtraInformationInput is used to update patient records metadata.
type PatientExtraInformationInput struct {
	PatientID     string           `json:"patientID"`
	MaritalStatus *MaritalStatus   `json:"maritalStatus"`
	Languages     []*base.Language `json:"languages"`
	Emails        []*EmailInput    `json:"emails"`
}

// USSDNextOfKinCreationInput is used to register next of kin via USSD.
type USSDNextOfKinCreationInput struct {
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	OtherNames string    `json:"otherNames"`
	BirthDate  time.Time `json:"birthDate"`
	Gender     string    `json:"gender"`
	Active     bool      `json:"active"`
	ParentID   string    `json:"parentID"`
}

// EpisodeOfCareStatus is used to record the status of an episode of care.
type EpisodeOfCareStatus string

// Episode of care valueset status values
const (
	// EpisodeOfCareStatusPlanned ...
	EpisodeOfCareStatusPlanned EpisodeOfCareStatus = "planned"
	// EpisodeOfCareStatusWaitlist ...
	EpisodeOfCareStatusWaitlist EpisodeOfCareStatus = "waitlist"
	// EpisodeOfCareStatusActive ...
	EpisodeOfCareStatusActive EpisodeOfCareStatus = "active"
	// EpisodeOfCareStatusOnhold ...
	EpisodeOfCareStatusOnhold EpisodeOfCareStatus = "onhold"
	// EpisodeOfCareStatusFinished ...
	EpisodeOfCareStatusFinished EpisodeOfCareStatus = "finished"
	// EpisodeOfCareStatusCancelled ...
	EpisodeOfCareStatusCancelled EpisodeOfCareStatus = "cancelled"
	// EpisodeOfCareStatusEnteredInError ...
	EpisodeOfCareStatusEnteredInError EpisodeOfCareStatus = "entered_in_error"
	// EpisodeOfCareStatusEnteredInErrorCanonical ...
	EpisodeOfCareStatusEnteredInErrorCanonical EpisodeOfCareStatus = "entered-in-error"
)

// AllEpisodeOfCareStatus is a list of episode of care statuses
var AllEpisodeOfCareStatus = []EpisodeOfCareStatus{
	EpisodeOfCareStatusPlanned,
	EpisodeOfCareStatusWaitlist,
	EpisodeOfCareStatusActive,
	EpisodeOfCareStatusOnhold,
	EpisodeOfCareStatusFinished,
	EpisodeOfCareStatusCancelled,
	EpisodeOfCareStatusEnteredInError,
	EpisodeOfCareStatusEnteredInErrorCanonical,
}

// IsValid validates episode of care status values
func (e EpisodeOfCareStatus) IsValid() bool {
	switch e {
	case EpisodeOfCareStatusPlanned, EpisodeOfCareStatusWaitlist, EpisodeOfCareStatusActive, EpisodeOfCareStatusOnhold, EpisodeOfCareStatusFinished, EpisodeOfCareStatusCancelled, EpisodeOfCareStatusEnteredInError, EpisodeOfCareStatusEnteredInErrorCanonical:
		return true
	}
	return false
}

// String ...
func (e EpisodeOfCareStatus) String() string {
	return string(e)
}

// UnmarshalGQL converts the input value into an episode of care
func (e *EpisodeOfCareStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EpisodeOfCareStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EpisodeOfCareStatus", str)
	}
	return nil
}

// MarshalGQL writes the episode of care status value to the supplied writer
func (e EpisodeOfCareStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))

}

// IDDocumentType is an internal code system for identification document types.
type IDDocumentType string

// ID type constants
const (
	// IDDocumentTypeNationalID ...
	IDDocumentTypeNationalID IDDocumentType = "national_id"
	// IDDocumentTypePassport ...
	IDDocumentTypePassport IDDocumentType = "passport"
	// IDDocumentTypeAlienID ...
	IDDocumentTypeAlienID IDDocumentType = "alien_id"
)

// AllIDDocumentType is a list of known ID types
var AllIDDocumentType = []IDDocumentType{
	IDDocumentTypeNationalID,
	IDDocumentTypePassport,
	IDDocumentTypeAlienID,
}

// IsValid checks that the ID type is valid
func (e IDDocumentType) IsValid() bool {
	switch e {
	case IDDocumentTypeNationalID, IDDocumentTypePassport, IDDocumentTypeAlienID:
		return true
	}
	return false
}

// String ...
func (e IDDocumentType) String() string {
	return string(e)
}

// UnmarshalGQL translates the input value to an ID type
func (e *IDDocumentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IDDocumentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IDDocumentType", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e IDDocumentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))

}

// MaritalStatus is used to code individuals' marital statuses.
//
// See: https://www.hl7.org/fhir/valueset-marital-status.html
type MaritalStatus string

// known marital statuses
const (
	// MaritalStatusA ...
	MaritalStatusA MaritalStatus = "A"
	// MaritalStatusD ...
	MaritalStatusD MaritalStatus = "D"
	// MaritalStatusI ...
	MaritalStatusI MaritalStatus = "I"
	// MaritalStatusL ...
	MaritalStatusL MaritalStatus = "L"
	// MaritalStatusM ...
	MaritalStatusM MaritalStatus = "M"
	// MaritalStatusP ...
	MaritalStatusP MaritalStatus = "P"
	// MaritalStatusS ...
	MaritalStatusS MaritalStatus = "S"
	// MaritalStatusT ...
	MaritalStatusT MaritalStatus = "T"
	// MaritalStatusU ...
	MaritalStatusU MaritalStatus = "U"
	// MaritalStatusW ...
	MaritalStatusW MaritalStatus = "W"
	// MaritalStatusUnk ...
	MaritalStatusUnk MaritalStatus = "UNK"
)

// AllMaritalStatus is a list of known marital statuses
var AllMaritalStatus = []MaritalStatus{
	MaritalStatusA,
	MaritalStatusD,
	MaritalStatusI,
	MaritalStatusL,
	MaritalStatusM,
	MaritalStatusP,
	MaritalStatusS,
	MaritalStatusT,
	MaritalStatusU,
	MaritalStatusW,
	MaritalStatusUnk,
}

// IsValid checks that the marital status is valid
func (e MaritalStatus) IsValid() bool {
	switch e {
	case MaritalStatusA, MaritalStatusD, MaritalStatusI, MaritalStatusL, MaritalStatusM, MaritalStatusP, MaritalStatusS, MaritalStatusT, MaritalStatusU, MaritalStatusW, MaritalStatusUnk:
		return true
	}
	return false
}

// String ...
func (e MaritalStatus) String() string {
	return string(e)
}

// UnmarshalGQL turns the supplied input into a marital status enum value
func (e *MaritalStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MaritalStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MaritalStatus", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e MaritalStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// RelationshipType defines relationship types for patients.
//
// See: https://www.hl7.org/fhir/valueset-relatedperson-relationshiptype.html
type RelationshipType string

// list of known relationship types
const (
	// RelationshipTypeC ...
	RelationshipTypeC RelationshipType = "C"
	// RelationshipTypeE ...
	RelationshipTypeE RelationshipType = "E"
	// RelationshipTypeF ...
	RelationshipTypeF RelationshipType = "F"
	// RelationshipTypeI ...
	RelationshipTypeI RelationshipType = "I"
	// RelationshipTypeN ...
	RelationshipTypeN RelationshipType = "N"
	// RelationshipTypeO ...
	RelationshipTypeO RelationshipType = "O"
	// RelationshipTypeS ...
	RelationshipTypeS RelationshipType = "S"
	// RelationshipTypeU ...
	RelationshipTypeU RelationshipType = "U"
)

// AllRelationshipType is a list of all known relationship types
var AllRelationshipType = []RelationshipType{
	RelationshipTypeC,
	RelationshipTypeE,
	RelationshipTypeF,
	RelationshipTypeI,
	RelationshipTypeN,
	RelationshipTypeO,
	RelationshipTypeS,
	RelationshipTypeU,
}

// IsValid ensures that the relationship type is valid
func (e RelationshipType) IsValid() bool {
	switch e {
	case RelationshipTypeC, RelationshipTypeE, RelationshipTypeF, RelationshipTypeI, RelationshipTypeN, RelationshipTypeO, RelationshipTypeS, RelationshipTypeU:
		return true
	}
	return false
}

// String ...
func (e RelationshipType) String() string {
	return string(e)
}

// UnmarshalGQL converts its input (if valid) into a relationship type
func (e *RelationshipType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RelationshipType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RelationshipType", str)
	}
	return nil
}

// MarshalGQL writes the relationship type to the supplied writer
func (e RelationshipType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// IdentifierUse is a code system for identifier uses.
//
// See: https://www.hl7.org/fhir/valueset-identifier-use.html
type IdentifierUse string

// known identifier use values
const (
	// IdentifierUseUsual ...
	IdentifierUseUsual IdentifierUse = "usual"
	// IdentifierUseOfficial ...
	IdentifierUseOfficial IdentifierUse = "official"
	// IdentifierUseTemp ...
	IdentifierUseTemp IdentifierUse = "temp"
	// IdentifierUseSecondary ...
	IdentifierUseSecondary IdentifierUse = "secondary"
	// IdentifierUseOld ...
	IdentifierUseOld IdentifierUse = "old"
)

// AllIdentifierUse is a list of all known identifier uses
var AllIdentifierUse = []IdentifierUse{
	IdentifierUseUsual,
	IdentifierUseOfficial,
	IdentifierUseTemp,
	IdentifierUseSecondary,
	IdentifierUseOld,
}

// IsValid returns True if the enum value is valid
func (e IdentifierUse) IsValid() bool {
	switch e {
	case IdentifierUseUsual, IdentifierUseOfficial, IdentifierUseTemp, IdentifierUseSecondary, IdentifierUseOld:
		return true
	}
	return false
}

// String ...
func (e IdentifierUse) String() string {
	return string(e)
}

// UnmarshalGQL translates from the supplied value to a valid enum value
func (e *IdentifierUse) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IdentifierUse(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IdentifierUse", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e IdentifierUse) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// ContactPointSystem defines the type of contact it is.
//
// See: http://hl7.org/fhir/valueset-contact-point-system.html
type ContactPointSystem string

// known contact systems
const (
	// ContactPointSystemPhone ...
	ContactPointSystemPhone ContactPointSystem = "phone"
	// ContactPointSystemFax ...
	ContactPointSystemFax ContactPointSystem = "fax"
	// ContactPointSystemEmail ...
	ContactPointSystemEmail ContactPointSystem = "email"
	// ContactPointSystemPager ...
	ContactPointSystemPager ContactPointSystem = "pager"
	// ContactPointSystemURL ...
	ContactPointSystemURL ContactPointSystem = "url"
	// ContactPointSystemSms ...
	ContactPointSystemSms ContactPointSystem = "sms"
	// ContactPointSystemOther ...
	ContactPointSystemOther ContactPointSystem = "other"
)

// AllContactPointSystem is a list of known contact systems
var AllContactPointSystem = []ContactPointSystem{
	ContactPointSystemPhone,
	ContactPointSystemFax,
	ContactPointSystemEmail,
	ContactPointSystemPager,
	ContactPointSystemURL,
	ContactPointSystemSms,
	ContactPointSystemOther,
}

// IsValid checks that the contact system is valid
func (e ContactPointSystem) IsValid() bool {
	switch e {
	case ContactPointSystemPhone, ContactPointSystemFax, ContactPointSystemEmail, ContactPointSystemPager, ContactPointSystemURL, ContactPointSystemSms, ContactPointSystemOther:
		return true
	}
	return false
}

// String ...
func (e ContactPointSystem) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a contact point system
func (e *ContactPointSystem) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ContactPointSystem(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ContactPointSystem", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e ContactPointSystem) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// ContactPointUse defines the uses of a contact.
//
// See: https://www.hl7.org/fhir/valueset-contact-point-use.html
type ContactPointUse string

// contact point uses
const (
	// ContactPointUseHome ...
	ContactPointUseHome ContactPointUse = "home"
	// ContactPointUseWork ...
	ContactPointUseWork ContactPointUse = "work"
	// ContactPointUseTemp ...
	ContactPointUseTemp ContactPointUse = "temp"
	// ContactPointUseOld ...
	ContactPointUseOld ContactPointUse = "old"
	// ContactPointUseMobile ...
	ContactPointUseMobile ContactPointUse = "mobile"
)

// AllContactPointUse is a list of known contact point uses
var AllContactPointUse = []ContactPointUse{
	ContactPointUseHome,
	ContactPointUseWork,
	ContactPointUseTemp,
	ContactPointUseOld,
	ContactPointUseMobile,
}

// IsValid checks that the enum value is valid
func (e ContactPointUse) IsValid() bool {
	switch e {
	case ContactPointUseHome, ContactPointUseWork, ContactPointUseTemp, ContactPointUseOld, ContactPointUseMobile:
		return true
	}
	return false
}

// String ...
func (e ContactPointUse) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied interface to a contact point use value
func (e *ContactPointUse) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ContactPointUse(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ContactPointUse", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e ContactPointUse) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// NameUse is used to define the uses of a human name.
//
// See: https://www.hl7.org/fhir/valueset-name-use.html
type NameUse string

// known name uses
const (
	// NameUseUsual ...
	NameUseUsual NameUse = "usual"
	// NameUseOfficial ...
	NameUseOfficial NameUse = "official"
	// NameUseTemp ...
	NameUseTemp NameUse = "temp"
	// NameUseNickname ...
	NameUseNickname NameUse = "nickname"
	// NameUseAnonymous ...
	NameUseAnonymous NameUse = "anonymous"
	// NameUseOld ...
	NameUseOld NameUse = "old"
	// NameUseMaiden ...
	NameUseMaiden NameUse = "maiden"
)

// AllNameUse is a list of known name uses
var AllNameUse = []NameUse{
	NameUseUsual,
	NameUseOfficial,
	NameUseTemp,
	NameUseNickname,
	NameUseAnonymous,
	NameUseOld,
	NameUseMaiden,
}

// IsValid checks that the name use is valid
func (e NameUse) IsValid() bool {
	switch e {
	case NameUseUsual, NameUseOfficial, NameUseTemp, NameUseNickname, NameUseAnonymous, NameUseOld, NameUseMaiden:
		return true
	}
	return false
}

// String ...
func (e NameUse) String() string {
	return string(e)
}

// UnmarshalGQL turns the supplied value into a name use enum
func (e *NameUse) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NameUse(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NameUse", str)
	}
	return nil
}

// MarshalGQL writes the name use enum value to the supplied writer
func (e NameUse) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AddressType is used to determine the type of an address.
//
// See: https://www.hl7.org/fhir/valueset-address-type.html
type AddressType string

// known address types
const (
	// AddressTypePostal ...
	AddressTypePostal AddressType = "postal"
	// AddressTypePhysical ...
	AddressTypePhysical AddressType = "physical"
	// AddressTypeBoth ...
	AddressTypeBoth AddressType = "both"
)

// AllAddressType is a list of all known address types
var AllAddressType = []AddressType{
	AddressTypePostal,
	AddressTypePhysical,
	AddressTypeBoth,
}

// IsValid checks that the address type is valid
func (e AddressType) IsValid() bool {
	switch e {
	case AddressTypePostal, AddressTypePhysical, AddressTypeBoth:
		return true
	}
	return false
}

// String renders the address type as a string
func (e AddressType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to an address type
func (e *AddressType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AddressType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AddressType", str)
	}
	return nil
}

// MarshalGQL writes the address type to the supplied writer
func (e AddressType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// AddressUse is used to set address uses.
//
// See: http://hl7.org/fhir/valueset-address-use.html
type AddressUse string

// known address uses
const (
	// AddressUseHome ...
	AddressUseHome AddressUse = "home"
	// AddressUseWork ...
	AddressUseWork AddressUse = "work"
	// AddressUseTemp ...
	AddressUseTemp AddressUse = "temp"
	// AddressUseOld ...
	AddressUseOld AddressUse = "old"
	// AddressUseBilling ...
	AddressUseBilling AddressUse = "billing"
)

// AllAddressUse is a list of all known address uses
var AllAddressUse = []AddressUse{
	AddressUseHome,
	AddressUseWork,
	AddressUseTemp,
	AddressUseOld,
	AddressUseBilling,
}

// IsValid returns true if an address use is valid
func (e AddressUse) IsValid() bool {
	switch e {
	case AddressUseHome, AddressUseWork, AddressUseTemp, AddressUseOld, AddressUseBilling:
		return true
	}
	return false
}

// String ...
func (e AddressUse) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to an address use
func (e *AddressUse) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AddressUse(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AddressUse", str)
	}
	return nil
}

// MarshalGQL writes the address to the supplied writer
func (e AddressUse) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Country codes.
//
// See: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes and
// https://www.iban.com/country-codes .
type Country string

// list of known country codes
const (
	// CountryAf ...
	CountryAf Country = "AF"
	// CountryAx ...
	CountryAx Country = "AX"
	// CountryAl ...
	CountryAl Country = "AL"
	// CountryDz ...
	CountryDz Country = "DZ"
	// CountryAs ...
	CountryAs Country = "AS"
	// CountryAd ...
	CountryAd Country = "AD"
	// CountryAo ...
	CountryAo Country = "AO"
	// CountryAi ...
	CountryAi Country = "AI"
	// CountryAq ...
	CountryAq Country = "AQ"
	// CountryAg ...
	CountryAg Country = "AG"
	// CountryAr ...
	CountryAr Country = "AR"
	// CountryAm ...
	CountryAm Country = "AM"
	// CountryAw ...
	CountryAw Country = "AW"
	// CountryAu ...
	CountryAu Country = "AU"
	// CountryAt ...
	CountryAt Country = "AT"
	// CountryAz ...
	CountryAz Country = "AZ"
	// CountryBs ...
	CountryBs Country = "BS"
	// CountryBh ...
	CountryBh Country = "BH"
	// CountryBd ...
	CountryBd Country = "BD"
	// CountryBb ...
	CountryBb Country = "BB"
	// CountryBy ...
	CountryBy Country = "BY"
	// CountryBe ...
	CountryBe Country = "BE"
	// CountryBz ...
	CountryBz Country = "BZ"
	// CountryBj ...
	CountryBj Country = "BJ"
	// CountryBm ...
	CountryBm Country = "BM"
	// CountryBt ...
	CountryBt Country = "BT"
	// CountryBo ...
	CountryBo Country = "BO"
	// CountryBq ...
	CountryBq Country = "BQ"
	// CountryBa ...
	CountryBa Country = "BA"
	// CountryBw ...
	CountryBw Country = "BW"
	// CountryBv ...
	CountryBv Country = "BV"
	// CountryBr ...
	CountryBr Country = "BR"
	// CountryIo ...
	CountryIo Country = "IO"
	// CountryBn ...
	CountryBn Country = "BN"
	// CountryBg ...
	CountryBg Country = "BG"
	// CountryBf ...
	CountryBf Country = "BF"
	// CountryBi ...
	CountryBi Country = "BI"
	// CountryCv ...
	CountryCv Country = "CV"
	// CountryKh ...
	CountryKh Country = "KH"
	// CountryCm ...
	CountryCm Country = "CM"
	// CountryCa ...
	CountryCa Country = "CA"
	// CountryKy ...
	CountryKy Country = "KY"
	// CountryCf ...
	CountryCf Country = "CF"
	// CountryTd ...
	CountryTd Country = "TD"
	// CountryCl ...
	CountryCl Country = "CL"
	// CountryCn ...
	CountryCn Country = "CN"
	// CountryCx ...
	CountryCx Country = "CX"
	// CountryCc ...
	CountryCc Country = "CC"
	// CountryCo ...
	CountryCo Country = "CO"
	// CountryKm ...
	CountryKm Country = "KM"
	// CountryCg ...
	CountryCg Country = "CG"
	// CountryCd ...
	CountryCd Country = "CD"
	// CountryCk ...
	CountryCk Country = "CK"
	// CountryCr ...
	CountryCr Country = "CR"
	// CountryCi ...
	CountryCi Country = "CI"
	// CountryHr ...
	CountryHr Country = "HR"
	// CountryCu ...
	CountryCu Country = "CU"
	// CountryCw ...
	CountryCw Country = "CW"
	// CountryCy ...
	CountryCy Country = "CY"
	// CountryCz ...
	CountryCz Country = "CZ"
	// CountryDk ...
	CountryDk Country = "DK"
	// CountryDj ...
	CountryDj Country = "DJ"
	// CountryDm ...
	CountryDm Country = "DM"
	// CountryDo ...
	CountryDo Country = "DO"
	// CountryEc ...
	CountryEc Country = "EC"
	// CountryEg ...
	CountryEg Country = "EG"
	// CountrySv ...
	CountrySv Country = "SV"
	// CountryGq ...
	CountryGq Country = "GQ"
	// CountryEr ...
	CountryEr Country = "ER"
	// CountryEe ...
	CountryEe Country = "EE"
	// CountrySz ...
	CountrySz Country = "SZ"
	// CountryEt ...
	CountryEt Country = "ET"
	// CountryFk ...
	CountryFk Country = "FK"
	// CountryFo ...
	CountryFo Country = "FO"
	// CountryFj ...
	CountryFj Country = "FJ"
	// CountryFi ...
	CountryFi Country = "FI"
	// CountryFr ...
	CountryFr Country = "FR"
	// CountryGf ...
	CountryGf Country = "GF"
	// CountryPf ...
	CountryPf Country = "PF"
	// CountryTf ...
	CountryTf Country = "TF"
	// CountryGa ...
	CountryGa Country = "GA"
	// CountryGm ...
	CountryGm Country = "GM"
	// CountryGe ...
	CountryGe Country = "GE"
	// CountryDe ...
	CountryDe Country = "DE"
	// CountryGh ...
	CountryGh Country = "GH"
	// CountryGi ...
	CountryGi Country = "GI"
	// CountryGr ...
	CountryGr Country = "GR"
	// CountryGl ...
	CountryGl Country = "GL"
	// CountryGd ...
	CountryGd Country = "GD"
	// CountryGp ...
	CountryGp Country = "GP"
	// CountryGu ...
	CountryGu Country = "GU"
	// CountryGt ...
	CountryGt Country = "GT"
	// CountryGg ...
	CountryGg Country = "GG"
	// CountryGn ...
	CountryGn Country = "GN"
	// CountryGw ...
	CountryGw Country = "GW"
	// CountryGy ...
	CountryGy Country = "GY"
	// CountryHt ...
	CountryHt Country = "HT"
	// CountryHm ...
	CountryHm Country = "HM"
	// CountryVa ...
	CountryVa Country = "VA"
	// CountryHn ...
	CountryHn Country = "HN"
	// CountryHk ...
	CountryHk Country = "HK"
	// CountryHu ...
	CountryHu Country = "HU"
	// CountryIs ...
	CountryIs Country = "IS"
	// CountryIn ...
	CountryIn Country = "IN"
	// CountryID ...
	CountryID Country = "ID"
	// CountryIr ...
	CountryIr Country = "IR"
	// CountryIq ...
	CountryIq Country = "IQ"
	// CountryIe ...
	CountryIe Country = "IE"
	// CountryIm ...
	CountryIm Country = "IM"
	// CountryIl ...
	CountryIl Country = "IL"
	// CountryIt ...
	CountryIt Country = "IT"
	// CountryJm ...
	CountryJm Country = "JM"
	// CountryJp ...
	CountryJp Country = "JP"
	// CountryJe ...
	CountryJe Country = "JE"
	// CountryJo ...
	CountryJo Country = "JO"
	// CountryKz ...
	CountryKz Country = "KZ"
	// CountryKe ...
	CountryKe Country = "KE"
	// CountryKi ...
	CountryKi Country = "KI"
	// CountryKp ...
	CountryKp Country = "KP"
	// CountryKr ...
	CountryKr Country = "KR"
	// CountryKw ...
	CountryKw Country = "KW"
	// CountryKg ...
	CountryKg Country = "KG"
	// CountryLa ...
	CountryLa Country = "LA"
	// CountryLv ...
	CountryLv Country = "LV"
	// CountryLb ...
	CountryLb Country = "LB"
	// CountryLs ...
	CountryLs Country = "LS"
	// CountryLr ...
	CountryLr Country = "LR"
	// CountryLy ...
	CountryLy Country = "LY"
	// CountryLi ...
	CountryLi Country = "LI"
	// CountryLt ...
	CountryLt Country = "LT"
	// CountryLu ...
	CountryLu Country = "LU"
	// CountryMo ...
	CountryMo Country = "MO"
	// CountryMg ...
	CountryMg Country = "MG"
	// CountryMw ...
	CountryMw Country = "MW"
	// CountryMy ...
	CountryMy Country = "MY"
	// CountryMv ...
	CountryMv Country = "MV"
	// CountryMl ...
	CountryMl Country = "ML"
	// CountryMt ...
	CountryMt Country = "MT"
	// CountryMh ...
	CountryMh Country = "MH"
	// CountryMq ...
	CountryMq Country = "MQ"
	// CountryMr ...
	CountryMr Country = "MR"
	// CountryMu ...
	CountryMu Country = "MU"
	// CountryYt ...
	CountryYt Country = "YT"
	// CountryMx ...
	CountryMx Country = "MX"
	// CountryFm ...
	CountryFm Country = "FM"
	// CountryMd ...
	CountryMd Country = "MD"
	// CountryMc ...
	CountryMc Country = "MC"
	// CountryMn ...
	CountryMn Country = "MN"
	// CountryMe ...
	CountryMe Country = "ME"
	// CountryMs ...
	CountryMs Country = "MS"
	// CountryMa ...
	CountryMa Country = "MA"
	// CountryMz ...
	CountryMz Country = "MZ"
	// CountryMm ...
	CountryMm Country = "MM"
	// CountryNa ...
	CountryNa Country = "NA"
	// CountryNr ...
	CountryNr Country = "NR"
	// CountryNp ...
	CountryNp Country = "NP"
	// CountryNl ...
	CountryNl Country = "NL"
	// CountryNc ...
	CountryNc Country = "NC"
	// CountryNz ...
	CountryNz Country = "NZ"
	// CountryNi ...
	CountryNi Country = "NI"
	// CountryNe ...
	CountryNe Country = "NE"
	// CountryNg ...
	CountryNg Country = "NG"
	// CountryNu ...
	CountryNu Country = "NU"
	// CountryNf ...
	CountryNf Country = "NF"
	// CountryMk ...
	CountryMk Country = "MK"
	// CountryMp ...
	CountryMp Country = "MP"
	// CountryNo ...
	CountryNo Country = "NO"
	// CountryOm ...
	CountryOm Country = "OM"
	// CountryPk ...
	CountryPk Country = "PK"
	// CountryPw ...
	CountryPw Country = "PW"
	// CountryPs ...
	CountryPs Country = "PS"
	// CountryPa ...
	CountryPa Country = "PA"
	// CountryPg ...
	CountryPg Country = "PG"
	// CountryPy ...
	CountryPy Country = "PY"
	// CountryPe ...
	CountryPe Country = "PE"
	// CountryPh ...
	CountryPh Country = "PH"
	// CountryPn ...
	CountryPn Country = "PN"
	// CountryPl ...
	CountryPl Country = "PL"
	// CountryPt ...
	CountryPt Country = "PT"
	// CountryPr ...
	CountryPr Country = "PR"
	// CountryQa ...
	CountryQa Country = "QA"
	// CountryRe ...
	CountryRe Country = "RE"
	// CountryRo ...
	CountryRo Country = "RO"
	// CountryRu ...
	CountryRu Country = "RU"
	// CountryRw ...
	CountryRw Country = "RW"
	// CountryBl ...
	CountryBl Country = "BL"
	// CountrySh ...
	CountrySh Country = "SH"
	// CountryKn ...
	CountryKn Country = "KN"
	// CountryLc ...
	CountryLc Country = "LC"
	// CountryMf ...
	CountryMf Country = "MF"
	// CountryPm ...
	CountryPm Country = "PM"
	// CountryVc ...
	CountryVc Country = "VC"
	// CountryWs ...
	CountryWs Country = "WS"
	// CountrySm ...
	CountrySm Country = "SM"
	// CountrySt ...
	CountrySt Country = "ST"
	// CountrySa ...
	CountrySa Country = "SA"
	// CountrySn ...
	CountrySn Country = "SN"
	// CountryRs ...
	CountryRs Country = "RS"
	// CountrySc ...
	CountrySc Country = "SC"
	// CountrySl ...
	CountrySl Country = "SL"
	// CountrySg ...
	CountrySg Country = "SG"
	// CountrySx ...
	CountrySx Country = "SX"
	// CountrySk ...
	CountrySk Country = "SK"
	// CountrySi ...
	CountrySi Country = "SI"
	// CountrySb ...
	CountrySb Country = "SB"
	// CountrySo ...
	CountrySo Country = "SO"
	// CountryZa ...
	CountryZa Country = "ZA"
	// CountryGs ...
	CountryGs Country = "GS"
	// CountrySs ...
	CountrySs Country = "SS"
	// CountryEs ...
	CountryEs Country = "ES"
	// CountryLk ...
	CountryLk Country = "LK"
	// CountrySd ...
	CountrySd Country = "SD"
	// CountrySr ...
	CountrySr Country = "SR"
	// CountrySj ...
	CountrySj Country = "SJ"
	// CountrySe ...
	CountrySe Country = "SE"
	// CountryCh ...
	CountryCh Country = "CH"
	// CountrySy ...
	CountrySy Country = "SY"
	// CountryTw ...
	CountryTw Country = "TW"
	// CountryTj ...
	CountryTj Country = "TJ"
	// CountryTz ...
	CountryTz Country = "TZ"
	// CountryTh ...
	CountryTh Country = "TH"
	// CountryTl ...
	CountryTl Country = "TL"
	// CountryTg ...
	CountryTg Country = "TG"
	// CountryTk ...
	CountryTk Country = "TK"
	// CountryTo ...
	CountryTo Country = "TO"
	// CountryTt ...
	CountryTt Country = "TT"
	// CountryTn ...
	CountryTn Country = "TN"
	// CountryTr ...
	CountryTr Country = "TR"
	// CountryTm ...
	CountryTm Country = "TM"
	// CountryTc ...
	CountryTc Country = "TC"
	// CountryTv ...
	CountryTv Country = "TV"
	// CountryUg ...
	CountryUg Country = "UG"
	// CountryUa ...
	CountryUa Country = "UA"
	// CountryAe ...
	CountryAe Country = "AE"
	// CountryGb ...
	CountryGb Country = "GB"
	// CountryUs ...
	CountryUs Country = "US"
	// CountryUm ...
	CountryUm Country = "UM"
	// CountryUy ...
	CountryUy Country = "UY"
	// CountryUz ...
	CountryUz Country = "UZ"
	// CountryVu ...
	CountryVu Country = "VU"
	// CountryVe ...
	CountryVe Country = "VE"
	// CountryVn ...
	CountryVn Country = "VN"
	// CountryVg ...
	CountryVg Country = "VG"
	// CountryVi ...
	CountryVi Country = "VI"
	// CountryWf ...
	CountryWf Country = "WF"
	// CountryEh ...
	CountryEh Country = "EH"
	// CountryYe ...
	CountryYe Country = "YE"
	// CountryZm ...
	CountryZm Country = "ZM"
	// CountryZw ...
	CountryZw Country = "ZW"
)

// AllCountry is a list of all known country codes
var AllCountry = []Country{
	CountryAf,
	CountryAx,
	CountryAl,
	CountryDz,
	CountryAs,
	CountryAd,
	CountryAo,
	CountryAi,
	CountryAq,
	CountryAg,
	CountryAr,
	CountryAm,
	CountryAw,
	CountryAu,
	CountryAt,
	CountryAz,
	CountryBs,
	CountryBh,
	CountryBd,
	CountryBb,
	CountryBy,
	CountryBe,
	CountryBz,
	CountryBj,
	CountryBm,
	CountryBt,
	CountryBo,
	CountryBq,
	CountryBa,
	CountryBw,
	CountryBv,
	CountryBr,
	CountryIo,
	CountryBn,
	CountryBg,
	CountryBf,
	CountryBi,
	CountryCv,
	CountryKh,
	CountryCm,
	CountryCa,
	CountryKy,
	CountryCf,
	CountryTd,
	CountryCl,
	CountryCn,
	CountryCx,
	CountryCc,
	CountryCo,
	CountryKm,
	CountryCg,
	CountryCd,
	CountryCk,
	CountryCr,
	CountryCi,
	CountryHr,
	CountryCu,
	CountryCw,
	CountryCy,
	CountryCz,
	CountryDk,
	CountryDj,
	CountryDm,
	CountryDo,
	CountryEc,
	CountryEg,
	CountrySv,
	CountryGq,
	CountryEr,
	CountryEe,
	CountrySz,
	CountryEt,
	CountryFk,
	CountryFo,
	CountryFj,
	CountryFi,
	CountryFr,
	CountryGf,
	CountryPf,
	CountryTf,
	CountryGa,
	CountryGm,
	CountryGe,
	CountryDe,
	CountryGh,
	CountryGi,
	CountryGr,
	CountryGl,
	CountryGd,
	CountryGp,
	CountryGu,
	CountryGt,
	CountryGg,
	CountryGn,
	CountryGw,
	CountryGy,
	CountryHt,
	CountryHm,
	CountryVa,
	CountryHn,
	CountryHk,
	CountryHu,
	CountryIs,
	CountryIn,
	CountryID,
	CountryIr,
	CountryIq,
	CountryIe,
	CountryIm,
	CountryIl,
	CountryIt,
	CountryJm,
	CountryJp,
	CountryJe,
	CountryJo,
	CountryKz,
	CountryKe,
	CountryKi,
	CountryKp,
	CountryKr,
	CountryKw,
	CountryKg,
	CountryLa,
	CountryLv,
	CountryLb,
	CountryLs,
	CountryLr,
	CountryLy,
	CountryLi,
	CountryLt,
	CountryLu,
	CountryMo,
	CountryMg,
	CountryMw,
	CountryMy,
	CountryMv,
	CountryMl,
	CountryMt,
	CountryMh,
	CountryMq,
	CountryMr,
	CountryMu,
	CountryYt,
	CountryMx,
	CountryFm,
	CountryMd,
	CountryMc,
	CountryMn,
	CountryMe,
	CountryMs,
	CountryMa,
	CountryMz,
	CountryMm,
	CountryNa,
	CountryNr,
	CountryNp,
	CountryNl,
	CountryNc,
	CountryNz,
	CountryNi,
	CountryNe,
	CountryNg,
	CountryNu,
	CountryNf,
	CountryMk,
	CountryMp,
	CountryNo,
	CountryOm,
	CountryPk,
	CountryPw,
	CountryPs,
	CountryPa,
	CountryPg,
	CountryPy,
	CountryPe,
	CountryPh,
	CountryPn,
	CountryPl,
	CountryPt,
	CountryPr,
	CountryQa,
	CountryRe,
	CountryRo,
	CountryRu,
	CountryRw,
	CountryBl,
	CountrySh,
	CountryKn,
	CountryLc,
	CountryMf,
	CountryPm,
	CountryVc,
	CountryWs,
	CountrySm,
	CountrySt,
	CountrySa,
	CountrySn,
	CountryRs,
	CountrySc,
	CountrySl,
	CountrySg,
	CountrySx,
	CountrySk,
	CountrySi,
	CountrySb,
	CountrySo,
	CountryZa,
	CountryGs,
	CountrySs,
	CountryEs,
	CountryLk,
	CountrySd,
	CountrySr,
	CountrySj,
	CountrySe,
	CountryCh,
	CountrySy,
	CountryTw,
	CountryTj,
	CountryTz,
	CountryTh,
	CountryTl,
	CountryTg,
	CountryTk,
	CountryTo,
	CountryTt,
	CountryTn,
	CountryTr,
	CountryTm,
	CountryTc,
	CountryTv,
	CountryUg,
	CountryUa,
	CountryAe,
	CountryGb,
	CountryUs,
	CountryUm,
	CountryUy,
	CountryUz,
	CountryVu,
	CountryVe,
	CountryVn,
	CountryVg,
	CountryVi,
	CountryWf,
	CountryEh,
	CountryYe,
	CountryZm,
	CountryZw,
}

// IsValid returns True if a country code is valid
func (e Country) IsValid() bool {
	switch e {
	case CountryAf, CountryAx, CountryAl, CountryDz, CountryAs, CountryAd, CountryAo, CountryAi, CountryAq, CountryAg, CountryAr, CountryAm, CountryAw, CountryAu, CountryAt, CountryAz, CountryBs, CountryBh, CountryBd, CountryBb, CountryBy, CountryBe, CountryBz, CountryBj, CountryBm, CountryBt, CountryBo, CountryBq, CountryBa, CountryBw, CountryBv, CountryBr, CountryIo, CountryBn, CountryBg, CountryBf, CountryBi, CountryCv, CountryKh, CountryCm, CountryCa, CountryKy, CountryCf, CountryTd, CountryCl, CountryCn, CountryCx, CountryCc, CountryCo, CountryKm, CountryCg, CountryCd, CountryCk, CountryCr, CountryCi, CountryHr, CountryCu, CountryCw, CountryCy, CountryCz, CountryDk, CountryDj, CountryDm, CountryDo, CountryEc, CountryEg, CountrySv, CountryGq, CountryEr, CountryEe, CountrySz, CountryEt, CountryFk, CountryFo, CountryFj, CountryFi, CountryFr, CountryGf, CountryPf, CountryTf, CountryGa, CountryGm, CountryGe, CountryDe, CountryGh, CountryGi, CountryGr, CountryGl, CountryGd, CountryGp, CountryGu, CountryGt, CountryGg, CountryGn, CountryGw, CountryGy, CountryHt, CountryHm, CountryVa, CountryHn, CountryHk, CountryHu, CountryIs, CountryIn, CountryID, CountryIr, CountryIq, CountryIe, CountryIm, CountryIl, CountryIt, CountryJm, CountryJp, CountryJe, CountryJo, CountryKz, CountryKe, CountryKi, CountryKp, CountryKr, CountryKw, CountryKg, CountryLa, CountryLv, CountryLb, CountryLs, CountryLr, CountryLy, CountryLi, CountryLt, CountryLu, CountryMo, CountryMg, CountryMw, CountryMy, CountryMv, CountryMl, CountryMt, CountryMh, CountryMq, CountryMr, CountryMu, CountryYt, CountryMx, CountryFm, CountryMd, CountryMc, CountryMn, CountryMe, CountryMs, CountryMa, CountryMz, CountryMm, CountryNa, CountryNr, CountryNp, CountryNl, CountryNc, CountryNz, CountryNi, CountryNe, CountryNg, CountryNu, CountryNf, CountryMk, CountryMp, CountryNo, CountryOm, CountryPk, CountryPw, CountryPs, CountryPa, CountryPg, CountryPy, CountryPe, CountryPh, CountryPn, CountryPl, CountryPt, CountryPr, CountryQa, CountryRe, CountryRo, CountryRu, CountryRw, CountryBl, CountrySh, CountryKn, CountryLc, CountryMf, CountryPm, CountryVc, CountryWs, CountrySm, CountrySt, CountrySa, CountrySn, CountryRs, CountrySc, CountrySl, CountrySg, CountrySx, CountrySk, CountrySi, CountrySb, CountrySo, CountryZa, CountryGs, CountrySs, CountryEs, CountryLk, CountrySd, CountrySr, CountrySj, CountrySe, CountryCh, CountrySy, CountryTw, CountryTj, CountryTz, CountryTh, CountryTl, CountryTg, CountryTk, CountryTo, CountryTt, CountryTn, CountryTr, CountryTm, CountryTc, CountryTv, CountryUg, CountryUa, CountryAe, CountryGb, CountryUs, CountryUm, CountryUy, CountryUz, CountryVu, CountryVe, CountryVn, CountryVg, CountryVi, CountryWf, CountryEh, CountryYe, CountryZm, CountryZw:
		return true
	}
	return false
}

// String ...
func (e Country) String() string {
	return string(e)
}

// UnmarshalGQL turns the value into a country
func (e *Country) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Country(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Country", str)
	}
	return nil
}

// MarshalGQL writes the enum value to the supplied writer
func (e Country) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
