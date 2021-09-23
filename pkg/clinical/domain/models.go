package domain

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
)

// Dummy ..
type Dummy struct {
	ID string `json:"id"`
}

//IsEntity ...
func (d Dummy) IsEntity() {}

// IsNode ..
func (d *Dummy) IsNode() {}

// SetID sets the trace's ID
func (d *Dummy) SetID(id string) {
	d.ID = id
}

// NameInput is used to input patient names.
type NameInput struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	OtherNames *string `json:"otherNames"`
}

// IdentificationDocument is used to input e.g National ID or passport document
// numbers at patient registration.
type IdentificationDocument struct {
	DocumentType     IDDocumentType         `json:"documentType"`
	DocumentNumber   string                 `json:"documentNumber"`
	Title            *string                `json:"title,omitempty"`
	ImageContentType *enumutils.ContentType `json:"imageContentType,omitempty"`
	ImageBase64      *string                `json:"imageBase64,omitempty"`
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
	PhotoContentType enumutils.ContentType `json:"photoContentType"`
	PhotoBase64data  string                `json:"photoBase64data"`
	PhotoFilename    string                `json:"photoFilename"`
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
	BirthDate               scalarutils.Date          `json:"birthDate"`
	PhoneNumbers            []*PhoneNumberInput       `json:"phoneNumbers"`
	Photos                  []*PhotoInput             `json:"photos"`
	Emails                  []*EmailInput             `json:"emails"`
	PhysicalAddresses       []*PhysicalAddress        `json:"physicalAddresses"`
	PostalAddresses         []*PostalAddress          `json:"postalAddresses"`
	Gender                  string                    `json:"gender"`
	Active                  bool                      `json:"active"`
	MaritalStatus           MaritalStatus             `json:"maritalStatus"`
	Languages               []enumutils.Language      `json:"languages"`
	ReplicateUSSD           bool                      `json:"replicate_ussd,omitempty"`
}

// BreakGlassEpisodeCreationInput is used to start emergency episodes via a
// break glass protocol
type BreakGlassEpisodeCreationInput struct {
	PatientID       string `json:"patientID" firestore:"patientID"`
	ProviderCode    string `json:"providerCode" firestore:"providerCode"`
	PractitionerUID string `json:"practitionerUID" firestore:"practitionerUID"`
	// ProviderPhone is the provider phone number
	ProviderPhone string `json:"providerPhone" firestore:"providerPhone"`
	Otp           string `json:"otp" firestore:"otp"`
	FullAccess    bool   `json:"fullAccess" firestore:"fullAccess"`
	// PatientPhone is the patient phone number used to send alert to patient
	PatientPhone string `json:"patient_phone" firestore:"patient_phone"`
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
	PatientID             string                 `json:"patientID"`
	MembershipNumber      string                 `json:"membershipNumber"`
	FrontImageBase64      *string                `json:"frontImageBase64"`
	FrontImageContentType *enumutils.ContentType `json:"frontImageContentType"`
	RearImageBase64       *string                `json:"rearImageBase64"`
	RearImageContentType  *enumutils.ContentType `json:"rearImageContentType"`
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
	BirthDate         scalarutils.Date    `json:"birthDate"`
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
	PatientID     string                `json:"patientID"`
	MaritalStatus *MaritalStatus        `json:"maritalStatus"`
	Languages     []*enumutils.Language `json:"languages"`
	Emails        []*EmailInput         `json:"emails"`
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

// PhoneNumberPayload is a D.T.O that accepts a phone number
type PhoneNumberPayload struct {
	PhoneNumber string `json:"phoneNumber"`
}

// EmailOptIn is used to persist and manage email communication whitelists
type EmailOptIn struct {
	Email   string `json:"email" firestore:"optedIn"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
}

//IsEntity ...
func (e EmailOptIn) IsEntity() {}
