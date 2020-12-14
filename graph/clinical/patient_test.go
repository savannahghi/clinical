package clinical

import (
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

//ContactPointPayload - compose a test FHIR contact point input
func ContactPointPayload() []*FHIRContactPointInput {
	msisdn := "+254723002959"
	var contactPointSystem ContactPointSystemEnum = "phone"
	var contactPointUse ContactPointUseEnum = "work"
	var rank int64 = 1
	var ContactPointSystemEmail ContactPointSystemEnum = "email"
	var ContactPointUseHome ContactPointUseEnum = "home"
	email := base.GenerateRandomEmail()
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
