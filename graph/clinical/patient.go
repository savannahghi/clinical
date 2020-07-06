package clinical

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
)

func (s Service) lookupUSSDSessionPatient(ctx context.Context, input USSDClinicalRequest) (*FHIRPatient, error) {
	patient, err := s.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}
	return patient.Resource, nil
}

// GetFHIRPatient retrieves instances of FHIRPatient by ID
func (s Service) GetFHIRPatient(ctx context.Context, id string) (*FHIRPatientRelayPayload, error) {
	s.checkPreconditions()

	resourceType := "Patient"
	var resource FHIRPatient

	data, err := s.clinicalRepository.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &FHIRPatientRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// SearchFHIRPatient provides a search API for FHIRPatient
func (s Service) SearchFHIRPatient(ctx context.Context, params map[string]interface{}) (*FHIRPatientRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Patient"
	path := "_search"
	output := FHIRPatientRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIRPatient

		resourceBs, err := json.Marshal(result)
		if err != nil {
			logrus.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			logrus.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &FHIRPatientRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRPatient creates a FHIRPatient instance
func (s Service) CreateFHIRPatient(ctx context.Context, input FHIRPatientInput) (*FHIRPatientRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Patient"
	resource := FHIRPatient{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRPatientRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIRPatient updates a FHIRPatient instance
// The resource must have it's ID set.
func (s Service) UpdateFHIRPatient(ctx context.Context, input FHIRPatientInput) (*FHIRPatientRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Patient"
	resource := FHIRPatient{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIRPatientRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIRPatient deletes the FHIRPatient identified by the supplied ID
func (s Service) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	resourceType := "Patient"
	resp, err := s.clinicalRepository.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

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
	Edges []*FHIRPatientRelayEdge `json:"edges,omitempty"`

	PageInfo *base.PageInfo `json:"pageInfo,omitempty"`
}

// FHIRPatientRelayEdge is a Relay edge for Patient
type FHIRPatientRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIRPatient `json:"node,omitempty"`
}

// FHIRPatientRelayPayload is used to return single instances of Patient
type FHIRPatientRelayPayload struct {
	Resource *FHIRPatient `json:"resource,omitempty"`
}

// USSDClinicalRequest is used to request the patient profile, medical
// history etc
type USSDClinicalRequest struct {
	PatientID     string `json:"patientID" firestore:"patientID"`
	Msisdn        string `json:"msisdn" firestore:"msisdn"`
	UssdSessionID string `json:"ussdSessionID" firestore:"ussdSessionID"`
}

// USSDClinicalResponse is used to return the patient profile, medical history
// or visit information etc
type USSDClinicalResponse struct {
	ShortLink string `json:"shortLink"`
	Summary   string `json:"summary"`
	Text      string `json:"text"`
}

// PatientGenderEnum is a FHIR enum
type PatientGenderEnum string

const (
	// PatientGenderEnumMale ...
	PatientGenderEnumMale PatientGenderEnum = "male"
	// PatientGenderEnumFemale ...
	PatientGenderEnumFemale PatientGenderEnum = "female"
	// PatientGenderEnumOther ...
	PatientGenderEnumOther PatientGenderEnum = "other"
	// PatientGenderEnumUnknown ...
	PatientGenderEnumUnknown PatientGenderEnum = "unknown"
)

// AllPatientGenderEnum ...
var AllPatientGenderEnum = []PatientGenderEnum{
	PatientGenderEnumMale,
	PatientGenderEnumFemale,
	PatientGenderEnumOther,
	PatientGenderEnumUnknown,
}

// IsValid ...
func (e PatientGenderEnum) IsValid() bool {
	switch e {
	case PatientGenderEnumMale, PatientGenderEnumFemale, PatientGenderEnumOther, PatientGenderEnumUnknown:
		return true
	}
	return false
}

// String ...
func (e PatientGenderEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *PatientGenderEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PatientGenderEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PatientGenderEnum", str)
	}
	return nil
}

// MarshalGQL writes the patient gender to the supplied writer as a quoted string
func (e PatientGenderEnum) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// PatientContactGenderEnum is a FHIR enum
type PatientContactGenderEnum string

const (
	// PatientContactGenderEnumMale ...
	PatientContactGenderEnumMale PatientContactGenderEnum = "male"
	// PatientContactGenderEnumFemale ...
	PatientContactGenderEnumFemale PatientContactGenderEnum = "female"
	// PatientContactGenderEnumOther ...
	PatientContactGenderEnumOther PatientContactGenderEnum = "other"
	// PatientContactGenderEnumUnknown ...
	PatientContactGenderEnumUnknown PatientContactGenderEnum = "unknown"
)

// AllPatientContactGenderEnum ...
var AllPatientContactGenderEnum = []PatientContactGenderEnum{
	PatientContactGenderEnumMale,
	PatientContactGenderEnumFemale,
	PatientContactGenderEnumOther,
	PatientContactGenderEnumUnknown,
}

// IsValid ...
func (e PatientContactGenderEnum) IsValid() bool {
	switch e {
	case PatientContactGenderEnumMale, PatientContactGenderEnumFemale, PatientContactGenderEnumOther, PatientContactGenderEnumUnknown:
		return true
	}
	return false
}

// String ...
func (e PatientContactGenderEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *PatientContactGenderEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PatientContactGenderEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Patient_ContactGenderEnum", str)
	}
	return nil
}

// MarshalGQL writes the patient contact gender to the supplied writer as a quoted string
func (e PatientContactGenderEnum) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// PatientLinkTypeEnum is a FHIR enum
type PatientLinkTypeEnum string

const (
	// PatientLinkTypeEnumReplacedBy ...
	PatientLinkTypeEnumReplacedBy PatientLinkTypeEnum = "replaced_by"
	// PatientLinkTypeEnumReplaces ...
	PatientLinkTypeEnumReplaces PatientLinkTypeEnum = "replaces"
	// PatientLinkTypeEnumRefer ...
	PatientLinkTypeEnumRefer PatientLinkTypeEnum = "refer"
	// PatientLinkTypeEnumSeealso ...
	PatientLinkTypeEnumSeealso PatientLinkTypeEnum = "seealso"
)

// AllPatientLinkTypeEnum ...
var AllPatientLinkTypeEnum = []PatientLinkTypeEnum{
	PatientLinkTypeEnumReplacedBy,
	PatientLinkTypeEnumReplaces,
	PatientLinkTypeEnumRefer,
	PatientLinkTypeEnumSeealso,
}

// IsValid ...
func (e PatientLinkTypeEnum) IsValid() bool {
	switch e {
	case PatientLinkTypeEnumReplacedBy, PatientLinkTypeEnumReplaces, PatientLinkTypeEnumRefer, PatientLinkTypeEnumSeealso:
		return true
	}
	return false
}

// String ...
func (e PatientLinkTypeEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *PatientLinkTypeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PatientLinkTypeEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Patient_LinkTypeEnum", str)
	}
	return nil
}

// MarshalGQL writes the patient link type to the supplied writer as a quoted string
func (e PatientLinkTypeEnum) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}
