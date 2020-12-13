package clinical

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"gitlab.slade360emr.com/go/base"
)

// simple patient registration defaults
// should be reviewed later (ticket created)
const (
	DefaultCountry    = CountryKe
	DefaultAddressUse = AddressUseHome
	DefaultContactUse = ContactPointUseHome

	IDIdentifierSystem     = "healthcloud.iddocument"
	MSISDNIdentifierSystem = "healthcloud.msisdn"
	DefaultVersion         = "0.0.1"
	DefaultPhotoTitle      = "Patient Photo"
	DefaultPhotoFilename   = "photo.jpg"
)

// NameToHumanName translates the simple name input of simple
// patient registration to FHIR human names
func NameToHumanName(names []*NameInput) []*FHIRHumanNameInput {
	if names == nil {
		return nil
	}
	output := []*FHIRHumanNameInput{}
	for _, name := range names {
		otherNames := ""
		if name.OtherNames != nil {
			otherNames = *name.OtherNames
		}
		fullName := fmt.Sprintf(
			"%s, %s %s", name.LastName, name.FirstName, otherNames)
		use := HumanNameUseEnumOfficial
		humanName := &FHIRHumanNameInput{
			Given:  &name.FirstName,
			Family: &name.LastName,
			Use:    use,
			Period: DefaultPeriodInput(),
			Text:   fullName,
		}
		output = append(output, humanName)
	}
	return output
}

// IDToIdentifier translates simple identification
// document details to FHIR identifiers
func IDToIdentifier(
	ids []*IdentificationDocument, phones []*PhoneNumberInput) ([]*FHIRIdentifierInput, error) {
	if ids == nil || phones == nil {
		return nil, nil
	}
	output := []*FHIRIdentifierInput{}
	identificationDocumentIdentifierSystem := base.URI(IDIdentifierSystem)
	msisdnIdentifierSystem := base.URI(MSISDNIdentifierSystem)
	userSelected := true
	idSystem := base.URI(identificationDocumentIdentifierSystem)
	version := DefaultVersion

	for _, id := range ids {
		identifier := &FHIRIdentifierInput{
			Use: IdentifierUseEnumOfficial,
			Type: FHIRCodeableConceptInput{
				Coding: []*FHIRCodingInput{
					{
						System:       &identificationDocumentIdentifierSystem,
						Version:      &version,
						Code:         base.Code(id.DocumentNumber),
						Display:      id.DocumentNumber,
						UserSelected: &userSelected,
					},
				},
				Text: id.DocumentNumber,
			},
			System: &idSystem,
			Value:  id.DocumentNumber,
			Period: DefaultPeriodInput(),
		}
		output = append(output, identifier)
	}

	for _, phone := range phones {
		// assume already verified by the contact input transform step
		normalized, err := base.NormalizeMSISDN(phone.Msisdn)
		if err != nil {
			return nil, fmt.Errorf("invalid phone number: %v", err)
		}
		identifier := &FHIRIdentifierInput{
			Use: IdentifierUseEnumOfficial,
			Type: FHIRCodeableConceptInput{
				// Coding
				Coding: []*FHIRCodingInput{
					{
						System:       &identificationDocumentIdentifierSystem,
						Version:      &version,
						Code:         base.Code(normalized),
						Display:      normalized,
						UserSelected: &userSelected,
					},
				},
				Text: normalized,
			},
			System: &msisdnIdentifierSystem,
			Value:  normalized,
			Period: DefaultPeriodInput(),
		}
		output = append(output, identifier)
	}

	return output, nil
}

// ContactsToContactPointInput translates phone and email contacts to
// FHIR contact points
func ContactsToContactPointInput(
	phones []*PhoneNumberInput, emails []*EmailInput,
	firestoreClient *firestore.Client,
) ([]*FHIRContactPointInput, error) {
	if phones == nil && emails == nil {
		return nil, nil
	}
	output := []*FHIRContactPointInput{}
	rank := int64(1)
	phoneSystem := ContactPointSystemEnumPhone
	use := ContactPointUseEnumHome

	for _, phone := range phones {
		validPhone, err := base.ValidateMSISDN(
			phone.Msisdn, phone.VerificationCode, false, firestoreClient)
		if err != nil {
			return nil, fmt.Errorf("invalid phone number: %s", err)
		}
		phoneContact := &FHIRContactPointInput{
			System: &phoneSystem,
			Use:    &use,
			Rank:   &rank,
			Period: DefaultPeriodInput(),
			Value:  &validPhone,
		}
		output = append(output, phoneContact)
		rank++
	}

	emailSystem := ContactPointSystemEnumEmail
	for _, email := range emails {
		err := base.ValidateEmail(
			email.Email, email.CommunicationOptIn, firestoreClient)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %v", err)
		}
		emailContact := &FHIRContactPointInput{
			System: &emailSystem,
			Use:    &use,
			Rank:   &rank,
			Period: DefaultPeriodInput(),
			Value:  &email.Email,
		}
		output = append(output, emailContact)
		rank++
	}

	return output, nil
}

// ContactsToContactPoint translates phone and email contacts to
// FHIR contact points
func ContactsToContactPoint(
	phones []*PhoneNumberInput, emails []*EmailInput,
	firestoreClient *firestore.Client,
) ([]*FHIRContactPoint, error) {
	if phones == nil && emails == nil {
		return nil, nil
	}
	output := []*FHIRContactPoint{}
	rank := int64(1)
	contactUse := ContactPointUseEnumHome
	emailSystem := ContactPointSystemEnumEmail
	phoneSystem := ContactPointSystemEnumPhone

	for _, phone := range phones {
		validPhone, err := base.ValidateMSISDN(
			phone.Msisdn, phone.VerificationCode, false, firestoreClient)
		if err != nil {
			return nil, fmt.Errorf("invalid phone number: %v", err)
		}
		phoneContact := &FHIRContactPoint{
			System: &phoneSystem,
			Use:    &contactUse,
			Rank:   &rank,
			Period: DefaultPeriod(),
			Value:  &validPhone,
		}
		output = append(output, phoneContact)
		rank++
	}

	for _, email := range emails {
		err := base.ValidateEmail(
			email.Email, email.CommunicationOptIn, firestoreClient)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %v", err)
		}
		emailContact := &FHIRContactPoint{
			System: &emailSystem,
			Use:    &contactUse,
			Rank:   &rank,
			Period: DefaultPeriod(),
			Value:  &email.Email,
		}
		output = append(output, emailContact)
		rank++
	}

	return output, nil
}

// PhotosToAttachments translates patient photos to FHIR attachments
func PhotosToAttachments(
	ctx context.Context,
	photos []*PhotoInput,
	engagement *base.InterServiceClient,
) ([]*FHIRAttachmentInput, error) {
	if photos == nil {
		return []*FHIRAttachmentInput{}, nil
	}

	output := []*FHIRAttachmentInput{}
	for _, photo := range photos {
		uploadInput := base.UploadInput{
			Title:       "Patient Photo",
			ContentType: photo.PhotoContentType.String(),
			Language:    base.LanguageEn.String(),
			Base64data:  photo.PhotoBase64data,
			Filename:    photo.PhotoFilename,
		}

		resp, err := engagement.MakeRequest(
			http.MethodPost,
			uploadEndpoint,
			uploadInput,
		)
		if err != nil {
			return nil, fmt.Errorf("error sending upload: %w", err)
		}

		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading upload response: %w", err)
		}

		upload := base.Upload{}
		err = json.Unmarshal(respData, &upload)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal upload response: %w", err)
		}

		data, err := base64.StdEncoding.DecodeString(upload.Base64data)
		if err != nil {
			return nil, errors.Wrap(err, "upload base64 decode error")
		}

		hash := base.Base64Binary(upload.Hash)
		size := len(data)
		url := base.URL(upload.URL)
		now := base.DateTime(time.Now().Format(timeFormatStr))
		contentType := base.Code(photo.PhotoContentType.String())
		language := base.Code(DefaultLanguage)
		photoData := base.Base64Binary(photo.PhotoBase64data)
		attachment := &FHIRAttachmentInput{
			ContentType: &contentType,
			Language:    &language,
			Data:        &photoData,
			URL:         &url,
			Size:        &size,
			Hash:        &hash,
			Creation:    &now,
		}
		output = append(output, attachment)
	}
	return output, nil
}

// PhysicalPostalAddressesToFHIRAddresses translates address inputs to FHIR addresses
func PhysicalPostalAddressesToFHIRAddresses(
	physical []*PhysicalAddress, postal []*PostalAddress) []*FHIRAddressInput {
	if physical == nil && postal == nil {
		return nil
	}
	output := []*FHIRAddressInput{}
	addrUse := AddressUseEnumHome
	country := DefaultCountry.String()
	physicalAddrType := AddressTypeEnumPhysical
	postalAddrType := AddressTypeEnumPostal

	for _, postal := range postal {
		text := fmt.Sprintf("%s\n%s", postal.PostalAddress, postal.PostalCode)
		postalCode := base.Code(postal.PostalCode)
		postalAddr := &FHIRAddressInput{
			Use:        &addrUse,
			Type:       &postalAddrType,
			Country:    &country,
			Period:     DefaultPeriodInput(),
			PostalCode: &postalCode,
			Line:       &postal.PostalAddress,
			Text:       text,
		}
		output = append(output, postalAddr)
	}

	for _, physical := range physical {
		text := fmt.Sprintf(
			"%s\n%s", physical.MapsCode, physical.PhysicalAddress)
		mapsCode := base.Code(physical.MapsCode)
		physicalAddr := &FHIRAddressInput{
			Use:        &addrUse,
			Type:       &physicalAddrType,
			Country:    &country,
			Period:     DefaultPeriodInput(),
			PostalCode: &mapsCode,
			Line:       &physical.PhysicalAddress,
			Text:       text,
		}
		output = append(output, physicalAddr)
	}

	return output
}

// RelationshipTypeDisplay computes friendly string for relationship types
func RelationshipTypeDisplay(val RelationshipType) string {
	switch val {
	case RelationshipTypeC:
		return "Emergency Contact"
	case RelationshipTypeE:
		return "Employer"
	case RelationshipTypeF:
		return "Federal Agency"
	case RelationshipTypeI:
		return "Insurance Company"
	case RelationshipTypeN:
		return "Next-of-Kin"
	case RelationshipTypeO:
		return "Other"
	case RelationshipTypeS:
		return "State Agency"
	case RelationshipTypeU:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// MaritalStatusDisplay calculates the text display for a marital status
// See: https://www.hl7.org/fhir/valueset-marital-status.html
func MaritalStatusDisplay(val MaritalStatus) string {
	switch val {
	case MaritalStatusA:
		return "Annulled"
	case MaritalStatusD:
		return "Divorced"
	case MaritalStatusI:
		return "Interlocutory"
	case MaritalStatusL:
		return "Legally Separated"
	case MaritalStatusM:
		return "Married"
	case MaritalStatusP:
		return "Polygamous"
	case MaritalStatusS:
		return "Never Married"
	case MaritalStatusT:
		return "Domestic Partner"
	case MaritalStatusU:
		return "unmarried"
	case MaritalStatusW:
		return "Widowed"
	case MaritalStatusUnk:
		return "unknown"
	default:
		return "unknown"
	}
}

// MaritalStatusEnumToCodeableConceptInput turns the simple enum selected in the
// user interface to a FHIR codeable concept
func MaritalStatusEnumToCodeableConceptInput(val MaritalStatus) *FHIRCodeableConceptInput {
	userSelected := true
	output := &FHIRCodeableConceptInput{
		Coding: []*FHIRCodingInput{
			{
				Code:         base.Code(val.String()),
				Display:      MaritalStatusDisplay(val),
				UserSelected: &userSelected,
			},
		},
		Text: MaritalStatusDisplay(val),
	}
	return output
}

// LanguagesToCommunicationInputs translates the supplied languages to FHIR
// communication preferences
func LanguagesToCommunicationInputs(languages []base.Language) []*CommunicationInput {
	output := []*CommunicationInput{}
	for _, language := range languages {
		comm := &CommunicationInput{
			Language: &CodeableConceptInput{
				Coding: []*CodingInput{
					{
						Code:         language.String(),
						Display:      base.LanguageNames[language],
						UserSelected: true,
						System:       base.LanguageCodingSystem,
						Version:      base.LanguageCodingVersion,
					},
				},
				Text: base.LanguageNames[language],
			},
			Preferred: false,
		}
		output = append(output, comm)
	}
	return output
}

// GetPatientIDFromEpisode gets the patientID, given an episode of care
func GetPatientIDFromEpisode(patientRef string) (string, error) {
	patientRefParts := strings.Split(patientRef, "/")
	if len(patientRefParts) != 2 {
		return "", fmt.Errorf(
			"invalid patient ref '%s'; expected to have two parts separated by a /", patientRef)
	}
	patientID := patientRefParts[1]
	return patientID, nil
}

// MaritalStatusEnumToCodeableConcept turns the simple enum selected in the
// user interface to a FHIR codeable concept
func MaritalStatusEnumToCodeableConcept(val MaritalStatus) *FHIRCodeableConcept {
	sel := true
	disp := MaritalStatusDisplay(val)
	output := &FHIRCodeableConcept{
		Coding: []*FHIRCoding{
			{
				Code:         base.Code(val.String()),
				Display:      disp,
				UserSelected: &sel,
			},
		},
		Text: MaritalStatusDisplay(val),
	}
	return output
}

// PhysicalPostalAddressesToCombinedFHIRAddress translates address inputs to
// a single FHIR address.
//
// It is used for patient contacts (e.g next of kin) where the spec has only
// one address per next of kin.
func PhysicalPostalAddressesToCombinedFHIRAddress(
	physical []*PhysicalAddress,
	postal []*PostalAddress,
) *FHIRAddressInput {
	if physical == nil && postal == nil {
		return nil
	}
	addressUse := AddressUseEnumHome
	postalAddrType := AddressTypeEnumPostal
	ke := CountryKe.String()

	addr := &FHIRAddressInput{
		Use:     &addressUse,
		Type:    &postalAddrType,
		Country: &ke,
		Period:  DefaultPeriodInput(),
		Line:    nil, // will be replaced below
		Text:    "",  // will be replaced below
	}

	postalAddressLines := []string{}
	for _, postal := range postal {
		postalAddressLines = append(postalAddressLines, postal.PostalAddress)
		postalAddressLines = append(postalAddressLines, postal.PostalCode)
		if addr.PostalCode == nil {
			postalCode := base.Code(postal.PostalCode)
			addr.PostalCode = &postalCode
		}
	}
	combinedPostalAddress := strings.Join(postalAddressLines, "\n")
	addr.Line = &combinedPostalAddress

	physicalAddressLines := []string{}
	for _, physical := range physical {
		physicalAddressLines = append(physicalAddressLines, physical.PhysicalAddress)
		physicalAddressLines = append(physicalAddressLines, physical.MapsCode)
	}
	combinedPhysicalAddress := strings.Join(physicalAddressLines, "\n")
	addr.Text = combinedPhysicalAddress

	return addr
}
