package clinical

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

// SingleIdentifierPayload - compose a single FHIRIdentifierInput
func SingleIdentifierPayload(input *SingleIdentifierInput) *FHIRIdentifierInput {
	return &FHIRIdentifierInput{
		Use: input.IdentifierUse, // usual | official | temp | secondary | old (If known
		Type: FHIRCodeableConceptInput{ // Description of identifier
			Text: input.Value,
			Coding: []*FHIRCodingInput{
				{
					System:       &input.System,  // Identity of the terminology system
					Version:      &input.Version, // Version of the system - if relevant
					Code:         input.Code,     // Symbol in syntax defined by the system
					Display:      input.Value,    //Representation defined by the system
					UserSelected: &input.UserSelected,
				},
			},
		},
		System: &input.System, // uri: The namespace for the identifier value
		Value:  input.Value,   // string: The value that is unique
	}
}

// IdentifierPayload - compose a list of FHIRIdentifierInput
func IdentifierPayload() []*FHIRIdentifierInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "+254723002959",
		System:        "phone",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "msisdn",
	}
	return []*FHIRIdentifierInput{SingleIdentifierPayload(&identifier)}
}

// SingleNamePayload - compose a single FHIR name input
func SingleNamePayload() *FHIRHumanNameInput {
	fullName := "Prof Pogrey Kiogothe Sr"
	var nameUseOfficial HumanNameUseEnum = "official"
	contactFamilyName := "Paul Kiogothe"
	return &FHIRHumanNameInput{
		Use:    nameUseOfficial,
		Text:   fullName,
		Family: &contactFamilyName,
		// FHIR API expects an array for the below fields
		// Given:  &givenName,
		// Prefix: &prefix,
		// Suffix: &suffix,
	}
}

// NamePayload - compose a list of FHIR input names
func NamePayload() []*FHIRHumanNameInput {
	names := SingleNamePayload()
	return []*FHIRHumanNameInput{names}
}

//ContactPointPayload - compose a test FHIR contact point input
func ContactPointPayload() []*FHIRContactPointInput {
	msisdn := "+254723002959"
	var contactPointSystem ContactPointSystemEnum = "phone"
	var contactPointUse ContactPointUseEnum = "work"
	var rank int64 = 1
	var ContactPointSystemEmail ContactPointSystemEnum = "email"
	var ContactPointUseHome ContactPointUseEnum = "home"
	email := "juha@kalulu.com"
	var rank2 int64 = 2
	return []*FHIRContactPointInput{
		{
			System: &contactPointSystem,
			Use:    &contactPointUse,
			Value:  &msisdn,
			Rank:   &rank,
		},
		{

			System: &ContactPointSystemEmail,
			Use:    &ContactPointUseHome,
			Value:  &email,
			Rank:   &rank2,
		},
	}
}

// SingleAddressPayload - compose a single test FHIRAddress
func SingleAddressPayload() *FHIRAddressInput {
	var address AddressUseEnum = "work"
	var addressType AddressTypeEnum = "postal"
	// addressLine := "Juha Kalulu's Foxhole"
	city := "Nairobi"
	district := "Langata"
	state := "Nairobi"
	postalCode := base.Code("00200")
	countryKe := "KE"

	return &FHIRAddressInput{
		Use:  &address,
		Type: &addressType,
		Text: "GF77+PC Kiangima",
		// Line:       &addressLine,
		City:       &city,
		District:   &district,
		State:      &state,
		PostalCode: &postalCode,
		Country:    &countryKe,
	}
}

// AddressPayload - compose a list of FHIR addresses
func AddressPayload() []*FHIRAddressInput {
	adresses := SingleAddressPayload()
	return []*FHIRAddressInput{adresses}
}

// MaritalStatusPayload - compose an FHIR marital status input
func MaritalStatusPayload() *FHIRCodeableConceptInput {
	var maritalSystemURI base.URI = "http://terminology.hl7.org/CodeSystem/v3-MaritalStatus"
	userSelected := true
	version := "2018-08-12"
	return &FHIRCodeableConceptInput{
		Text: "Married",
		Coding: []*FHIRCodingInput{
			{
				System:       &maritalSystemURI,
				Code:         "M",
				Version:      &version,
				Display:      "Married",
				UserSelected: &userSelected,
			},
		},
	}
}

// PhotoPayload - compose an FHIRAttachmentInput
func PhotoPayload(t *testing.T) []*FHIRAttachmentInput {
	var contentType base.Code = "application/json"
	var language base.Code = "en"
	bs, err := ioutil.ReadFile("testdata/photo.jpg")
	if err != nil {
		t.Fatalf("unable to read test photo %s: ", err)
	}
	var photoBase64 = base.Base64Binary(base64.StdEncoding.EncodeToString(bs))
	var now base.DateTime = "2018-01-01"
	var URL base.URL = "https://healthcloud.co.ke"
	var hash base.Base64Binary = "fake"
	size := 1
	title := "Test Photo Attachment"
	return []*FHIRAttachmentInput{
		{
			ContentType: &contentType,
			Language:    &language,
			Data:        &photoBase64,
			Title:       &title,
			Creation:    &now,
			URL:         &URL,
			Hash:        &hash,
			Size:        &size,
		},
	}
}

// FHIRCodingPayload - compose a list of FHIRCodingInput
func FHIRCodingPayload(code base.Code, display string) []*FHIRCodingInput {
	codingInput := SingleFHIRCodingPayload(code, display)
	return []*FHIRCodingInput{codingInput}
}

// SingleCodeableConceptPayload - compose a single FHIRCodeableConceptInput
func SingleCodeableConceptPayload(code base.Code, display, text string) *FHIRCodeableConceptInput {
	return &FHIRCodeableConceptInput{
		Text:   text,
		Coding: FHIRCodingPayload(code, display),
	}
}

// CodeableConceptPayload - compose many FHIRCodeableConceptInput
func CodeableConceptPayload() []*FHIRCodeableConceptInput {
	display := "Next-of-Kin"
	var code base.Code = "N"
	text := "Married"

	codeableConceptInput := SingleCodeableConceptPayload(code, display, text)
	return []*FHIRCodeableConceptInput{codeableConceptInput}
}

// ContactPayload - compose a test FHIRPatientContactInput
func ContactPayload() []*FHIRPatientContactInput {
	var contactGender PatientContactGenderEnum = "female"
	return []*FHIRPatientContactInput{
		{
			Relationship: CodeableConceptPayload(),
			Name:         SingleNamePayload(),
			Telecom:      ContactPointPayload(),
			Address:      SingleAddressPayload(),
			Gender:       &contactGender,
		},
	}
}

// PatientCommunicationPayload - compose a patient FHIR communication
func PatientCommunicationPayload() []*FHIRPatientCommunicationInput {
	// var languageSystem base.URI = "urn:ietf:bcp:47"
	preferredLanguage := true
	var code base.Code = "sw"

	return []*FHIRPatientCommunicationInput{
		{
			Language:  SingleCodeableConceptPayload(code, "Swahili", "English"),
			Preferred: &preferredLanguage,
		},
	}
}

// PatientResourceFHIRPayload - fixture to create patient in FHIR
func PatientResourceFHIRPayload(t *testing.T) FHIRPatientInput {
	active := true
	var gender PatientGenderEnum = "male"

	return FHIRPatientInput{
		Identifier:    IdentifierPayload(),
		Active:        &active,
		Name:          NamePayload(),
		Telecom:       ContactPointPayload(),
		Gender:        &gender,
		BirthDate:     &base.Date{Year: 1985, Month: 1, Day: 1},
		Address:       AddressPayload(),
		MaritalStatus: MaritalStatusPayload(),
		Photo:         PhotoPayload(t),
		Contact:       ContactPayload(),
		Communication: PatientCommunicationPayload(),
	}
}

// CreateTestFHIRPatient - helper to create a test patient in FHIR
func CreateTestFHIRPatient(t *testing.T) FHIRPatientRelayPayload {
	service := NewService()
	patientPayload := PatientResourceFHIRPayload(t)
	ctx := context.Background()
	patient, err := service.CreateFHIRPatient(ctx, patientPayload)
	if err != nil {
		t.Fatalf("unable to create patient resource %s: ", err)
	}
	return *patient
}

// GetTestFHIRPatient - retrieve a created test patient in FHIR
func GetTestFHIRPatient(t *testing.T, id string) FHIRPatientRelayPayload {
	service := NewService()
	ctx := context.Background()
	patient, err := service.GetFHIRPatient(ctx, id)
	if err != nil {
		t.Fatalf("unable to retrieve patient resource %s: ", err)
	}
	return *patient
}

// TODO: Refactor the test below: use recommended search params

// func TestService_SearchACreatedFHIRPatient(t *testing.T) {
// 	service := NewService()
// 	ctx := context.Background()
// 	params := map[string]interface{}{
// 		"Name": "Kamau",
// 	}
// 	patient, err := service.SearchFHIRPatient(ctx, params)
// 	if err != nil {
// 		t.Fatalf("unable to search patient resource %s: ", err)
// 	}
// 	assert.NotNil(t, patient)
// }

func TestService_GetCreateFHIRPatient(t *testing.T) {
	createdPatient := CreateTestFHIRPatient(t)
	patientID := *createdPatient.Resource.ID
	patient := GetTestFHIRPatient(t, patientID)
	assert.NotNil(t, patient)
}
