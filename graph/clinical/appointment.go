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

// GetFHIRAppointment retrieves instances of FHIRAppointment by ID
func (s Service) GetFHIRAppointment(ctx context.Context, id string) (*FHIRAppointmentRelayPayload, error) {
	s.checkPreconditions()

	resourceType := "Appointment"
	var resource FHIRAppointment

	data, err := s.clinicalRepository.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &FHIRAppointmentRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// SearchFHIRAppointment provides a search API for FHIRAppointment
func (s Service) SearchFHIRAppointment(ctx context.Context, params map[string]interface{}) (*FHIRAppointmentRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Appointment"
	path := "_search"
	output := FHIRAppointmentRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIRAppointment

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
		output.Edges = append(output.Edges, &FHIRAppointmentRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIRAppointment creates a FHIRAppointment instance
func (s Service) CreateFHIRAppointment(ctx context.Context, input FHIRAppointmentInput) (*FHIRAppointmentRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Appointment"
	resource := FHIRAppointment{}

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

	output := &FHIRAppointmentRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIRAppointment updates a FHIRAppointment instance
// The resource must have it's ID set.
func (s Service) UpdateFHIRAppointment(ctx context.Context, input FHIRAppointmentInput) (*FHIRAppointmentRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Appointment"
	resource := FHIRAppointment{}

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

	output := &FHIRAppointmentRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIRAppointment deletes the FHIRAppointment identified by the supplied ID
func (s Service) DeleteFHIRAppointment(ctx context.Context, id string) (bool, error) {
	resourceType := "Appointment"
	resp, err := s.clinicalRepository.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// ListAppointments filter appointments by their provider code
func (s Service) ListAppointments(ctx context.Context, providerSladeCode int) (*FHIRAppointmentRelayConnection, error) {
	s.checkPreconditions()
	organizaionID, err := s.GetORCreateOrganization(ctx, providerSladeCode)
	if err != nil {
		return nil, fmt.Errorf(
			"internal server error in retrieving organisation : %v", err)
	}
	params := map[string]interface{}{
		"identifier": *organizaionID,
	}
	return s.SearchFHIRAppointment(ctx, params)

}

// FHIRAppointment definition: a booking of a healthcare event among patient(s), practitioner(s), related person(s) and/or device(s) for a specific date/time. this may result in one or more encounter(s).
type FHIRAppointment struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// This records identifiers associated with this appointment concern that are defined by business processes and/or used to refer to it when a direct URL reference to the resource itself is not appropriate (e.g. in CDA documents, or in written / printed documentation).
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// The overall status of the Appointment. Each of the participants has their own participation status which indicates their involvement in the process, however this status indicates the shared status.
	Status *AppointmentStatusEnum `json:"status,omitempty"`

	// The coded reason for the appointment being cancelled. This is often used in reporting/billing/futher processing to determine if further actions are required, or specific fees apply.
	CancelationReason *FHIRCodeableConcept `json:"cancelationReason,omitempty"`

	// A broad categorization of the service that is to be performed during this appointment.
	ServiceCategory []*FHIRCodeableConcept `json:"serviceCategory,omitempty"`

	// The specific service that is to be performed during this appointment.
	ServiceType []*FHIRCodeableConcept `json:"serviceType,omitempty"`

	// The specialty of a practitioner that would be required to perform the service requested in this appointment.
	Specialty []*FHIRCodeableConcept `json:"specialty,omitempty"`

	// The style of appointment or patient that has been booked in the slot (not service type).
	AppointmentType *FHIRCodeableConcept `json:"appointmentType,omitempty"`

	// The coded reason that this appointment is being scheduled. This is more clinical than administrative.
	ReasonCode *base.Code `json:"reasonCode,omitempty"`

	// Reason the appointment has been scheduled to take place, as specified using information from another resource. When the patient arrives and the encounter begins it may be used as the admission diagnosis. The indication will typically be a Condition (with other resources referenced in the evidence.detail), or a Procedure.
	ReasonReference []*FHIRReference `json:"reasonReference,omitempty"`

	// The priority of the appointment. Can be used to make informed decisions if needing to re-prioritize appointments. (The iCal Standard specifies 0 as undefined, 1 as highest, 9 as lowest priority).
	Priority *int `json:"priority,omitempty"`

	// The brief description of the appointment as would be shown on a subject line in a meeting request, or appointment list. Detailed or expanded information should be put in the comment field.
	Description *string `json:"description,omitempty"`

	// Additional information to support the appointment provided when making the appointment.
	SupportingInformation []*FHIRReference `json:"supportingInformation,omitempty"`

	// Date/Time that the appointment is to take place.
	Start *base.Instant `json:"start,omitempty"`

	// Date/Time that the appointment is to conclude.
	End *base.Instant `json:"end,omitempty"`

	// Number of minutes that the appointment is to take. This can be less than the duration between the start and end times.  For example, where the actual time of appointment is only an estimate or if a 30 minute appointment is being requested, but any time would work.  Also, if there is, for example, a planned 15 minute break in the middle of a long appointment, the duration may be 15 minutes less than the difference between the start and end.
	MinutesDuration *int `json:"minutesDuration,omitempty"`

	// The slots from the participants' schedules that will be filled by the appointment.
	Slot []*FHIRReference `json:"slot,omitempty"`

	// The date that this appointment was initially created. This could be different to the meta.lastModified value on the initial entry, as this could have been before the resource was created on the FHIR server, and should remain unchanged over the lifespan of the appointment.
	Created *base.DateTime `json:"created,omitempty"`

	// Additional comments about the appointment.
	Comment *string `json:"comment,omitempty"`

	// While Appointment.comment contains information for internal use, Appointment.patientInstructions is used to capture patient facing information about the Appointment (e.g. please bring your referral or fast from 8pm night before).
	PatientInstruction *string `json:"patientInstruction,omitempty"`

	// The service request this appointment is allocated to assess (e.g. incoming referral or procedure request).
	BasedOn []*FHIRReference `json:"basedOn,omitempty"`

	// List of participants involved in the appointment.
	Participant []*FHIRAppointmentParticipant `json:"participant,omitempty"`

	//     A set of date ranges (potentially including times) that the appointment is preferred to be scheduled within.
	//
	// The duration (usually in minutes) could also be provided to indicate the length of the appointment to fill and populate the start/end times for the actual allocated time. However, in other situations the duration may be calculated by the scheduling system.
	RequestedPeriod []*FHIRPeriod `json:"requestedPeriod,omitempty"`
}

// FHIRAppointmentInput is the input type for Appointment
type FHIRAppointmentInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// This records identifiers associated with this appointment concern that are defined by business processes and/or used to refer to it when a direct URL reference to the resource itself is not appropriate (e.g. in CDA documents, or in written / printed documentation).
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// The overall status of the Appointment. Each of the participants has their own participation status which indicates their involvement in the process, however this status indicates the shared status.
	Status *AppointmentStatusEnum `json:"status,omitempty"`

	// The coded reason for the appointment being cancelled. This is often used in reporting/billing/futher processing to determine if further actions are required, or specific fees apply.
	CancelationReason *FHIRCodeableConceptInput `json:"cancelationReason,omitempty"`

	// A broad categorization of the service that is to be performed during this appointment.
	ServiceCategory []*FHIRCodeableConceptInput `json:"serviceCategory,omitempty"`

	// The specific service that is to be performed during this appointment.
	ServiceType []*FHIRCodeableConceptInput `json:"serviceType,omitempty"`

	// The specialty of a practitioner that would be required to perform the service requested in this appointment.
	Specialty []*FHIRCodeableConceptInput `json:"specialty,omitempty"`

	// The style of appointment or patient that has been booked in the slot (not service type).
	AppointmentType *FHIRCodeableConceptInput `json:"appointmentType,omitempty"`

	// The coded reason that this appointment is being scheduled. This is more clinical than administrative.
	ReasonCode *base.Code `json:"reasonCode,omitempty"`

	// Reason the appointment has been scheduled to take place, as specified using information from another resource. When the patient arrives and the encounter begins it may be used as the admission diagnosis. The indication will typically be a Condition (with other resources referenced in the evidence.detail), or a Procedure.
	ReasonReference []*FHIRReferenceInput `json:"reasonReference,omitempty"`

	// The priority of the appointment. Can be used to make informed decisions if needing to re-prioritize appointments. (The iCal Standard specifies 0 as undefined, 1 as highest, 9 as lowest priority).
	Priority *int `json:"priority,omitempty"`

	// The brief description of the appointment as would be shown on a subject line in a meeting request, or appointment list. Detailed or expanded information should be put in the comment field.
	Description *string `json:"description,omitempty"`

	// Additional information to support the appointment provided when making the appointment.
	SupportingInformation []*FHIRReferenceInput `json:"supportingInformation,omitempty"`

	// Date/Time that the appointment is to take place.
	Start *base.Instant `json:"start,omitempty"`

	// Date/Time that the appointment is to conclude.
	End *base.Instant `json:"end,omitempty"`

	// Number of minutes that the appointment is to take. This can be less than the duration between the start and end times.  For example, where the actual time of appointment is only an estimate or if a 30 minute appointment is being requested, but any time would work.  Also, if there is, for example, a planned 15 minute break in the middle of a long appointment, the duration may be 15 minutes less than the difference between the start and end.
	MinutesDuration *int `json:"minutesDuration,omitempty"`

	// The slots from the participants' schedules that will be filled by the appointment.
	Slot []*FHIRReferenceInput `json:"slot,omitempty"`

	// The date that this appointment was initially created. This could be different to the meta.lastModified value on the initial entry, as this could have been before the resource was created on the FHIR server, and should remain unchanged over the lifespan of the appointment.
	Created *base.DateTime `json:"created,omitempty"`

	// Additional comments about the appointment.
	Comment *string `json:"comment,omitempty"`

	// While Appointment.comment contains information for internal use, Appointment.patientInstructions is used to capture patient facing information about the Appointment (e.g. please bring your referral or fast from 8pm night before).
	PatientInstruction *string `json:"patientInstruction,omitempty"`

	// The service request this appointment is allocated to assess (e.g. incoming referral or procedure request).
	BasedOn []*FHIRReferenceInput `json:"basedOn,omitempty"`

	// List of participants involved in the appointment.
	Participant []*FHIRAppointmentParticipantInput `json:"participant,omitempty"`

	//     A set of date ranges (potentially including times) that the appointment is preferred to be scheduled within.
	//
	// The duration (usually in minutes) could also be provided to indicate the length of the appointment to fill and populate the start/end times for the actual allocated time. However, in other situations the duration may be calculated by the scheduling system.
	RequestedPeriod []*FHIRPeriodInput `json:"requestedPeriod,omitempty"`
}

// FHIRAppointmentParticipant definition: a booking of a healthcare event among patient(s), practitioner(s), related person(s) and/or device(s) for a specific date/time. this may result in one or more encounter(s).
type FHIRAppointmentParticipant struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Role of participant in the appointment.
	Type []*FHIRCodeableConcept `json:"type,omitempty"`

	// A Person, Location/HealthcareService or Device that is participating in the appointment.
	Actor *FHIRReference `json:"actor,omitempty"`

	// Whether this participant is required to be present at the meeting. This covers a use-case where two doctors need to meet to discuss the results for a specific patient, and the patient is not required to be present.
	Required *AppointmentParticipantRequiredEnum `json:"required,omitempty"`

	// Participation status of the actor.
	Status *AppointmentParticipantStatusEnum `json:"status,omitempty"`

	// Participation period of the actor.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIRAppointmentParticipantInput is the input type for AppointmentParticipant
type FHIRAppointmentParticipantInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// Role of participant in the appointment.
	Type []*FHIRCodeableConceptInput `json:"type,omitempty"`

	// A Person, Location/HealthcareService or Device that is participating in the appointment.
	Actor *FHIRReferenceInput `json:"actor,omitempty"`

	// Whether this participant is required to be present at the meeting. This covers a use-case where two doctors need to meet to discuss the results for a specific patient, and the patient is not required to be present.
	Required *AppointmentParticipantRequiredEnum `json:"required,omitempty"`

	// Participation status of the actor.
	Status *AppointmentParticipantStatusEnum `json:"status,omitempty"`

	// Participation period of the actor.
	Period *FHIRPeriodInput `json:"period,omitempty"`
}

// FHIRAppointmentRelayConnection is a Relay connection for Appointment
type FHIRAppointmentRelayConnection struct {
	Edges    []*FHIRAppointmentRelayEdge `json:"edges,omitempty"`
	PageInfo *base.PageInfo              `json:"pageInfo,omitempty"`
}

// FHIRAppointmentRelayEdge is a Relay edge for Appointment
type FHIRAppointmentRelayEdge struct {
	Cursor *string          `json:"cursor,omitempty"`
	Node   *FHIRAppointment `json:"node,omitempty"`
}

// FHIRAppointmentRelayPayload is used to return single instances of Appointment
type FHIRAppointmentRelayPayload struct {
	Resource *FHIRAppointment `json:"resource,omitempty"`
}

// AppointmentStatusEnum is a FHIR enum
type AppointmentStatusEnum string

const (
	// AppointmentStatusEnumProposed ...
	AppointmentStatusEnumProposed AppointmentStatusEnum = "proposed"
	// AppointmentStatusEnumPending ...
	AppointmentStatusEnumPending AppointmentStatusEnum = "pending"
	// AppointmentStatusEnumBooked ...
	AppointmentStatusEnumBooked AppointmentStatusEnum = "booked"
	// AppointmentStatusEnumArrived ...
	AppointmentStatusEnumArrived AppointmentStatusEnum = "arrived"
	// AppointmentStatusEnumFulfilled ...
	AppointmentStatusEnumFulfilled AppointmentStatusEnum = "fulfilled"
	// AppointmentStatusEnumCancelled ...
	AppointmentStatusEnumCancelled AppointmentStatusEnum = "cancelled"
	// AppointmentStatusEnumNoshow ...
	AppointmentStatusEnumNoshow AppointmentStatusEnum = "noshow"
	// AppointmentStatusEnumEnteredInError ...
	AppointmentStatusEnumEnteredInError AppointmentStatusEnum = "entered_in_error"
	// AppointmentStatusEnumCheckedIn ...
	AppointmentStatusEnumCheckedIn AppointmentStatusEnum = "checked_in"
	// AppointmentStatusEnumWaitlist ...
	AppointmentStatusEnumWaitlist AppointmentStatusEnum = "waitlist"
)

// AllAppointmentStatusEnum ...
var AllAppointmentStatusEnum = []AppointmentStatusEnum{
	AppointmentStatusEnumProposed,
	AppointmentStatusEnumPending,
	AppointmentStatusEnumBooked,
	AppointmentStatusEnumArrived,
	AppointmentStatusEnumFulfilled,
	AppointmentStatusEnumCancelled,
	AppointmentStatusEnumNoshow,
	AppointmentStatusEnumEnteredInError,
	AppointmentStatusEnumCheckedIn,
	AppointmentStatusEnumWaitlist,
}

// IsValid ...
func (e AppointmentStatusEnum) IsValid() bool {
	switch e {
	case AppointmentStatusEnumProposed, AppointmentStatusEnumPending, AppointmentStatusEnumBooked, AppointmentStatusEnumArrived, AppointmentStatusEnumFulfilled, AppointmentStatusEnumCancelled, AppointmentStatusEnumNoshow, AppointmentStatusEnumEnteredInError, AppointmentStatusEnumCheckedIn, AppointmentStatusEnumWaitlist:
		return true
	}
	return false
}

// String ...
func (e AppointmentStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *AppointmentStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AppointmentStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AppointmentStatusEnum", str)
	}
	return nil
}

// MarshalGQL writes the appointment status to the supplied writer as a quoted string
func (e AppointmentStatusEnum) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// AppointmentParticipantRequiredEnum is a FHIR enum
type AppointmentParticipantRequiredEnum string

const (
	// AppointmentParticipantRequiredEnumRequired ...
	AppointmentParticipantRequiredEnumRequired AppointmentParticipantRequiredEnum = "required"
	// AppointmentParticipantRequiredEnumOptional ...
	AppointmentParticipantRequiredEnumOptional AppointmentParticipantRequiredEnum = "optional"
	// AppointmentParticipantRequiredEnumInformationOnly ...
	AppointmentParticipantRequiredEnumInformationOnly AppointmentParticipantRequiredEnum = "information_only"
)

// AllAppointmentParticipantRequiredEnum ...
var AllAppointmentParticipantRequiredEnum = []AppointmentParticipantRequiredEnum{
	AppointmentParticipantRequiredEnumRequired,
	AppointmentParticipantRequiredEnumOptional,
	AppointmentParticipantRequiredEnumInformationOnly,
}

// IsValid ...
func (e AppointmentParticipantRequiredEnum) IsValid() bool {
	switch e {
	case AppointmentParticipantRequiredEnumRequired, AppointmentParticipantRequiredEnumOptional, AppointmentParticipantRequiredEnumInformationOnly:
		return true
	}
	return false
}

// String ...
func (e AppointmentParticipantRequiredEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *AppointmentParticipantRequiredEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AppointmentParticipantRequiredEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Appointment_ParticipantRequiredEnum", str)
	}
	return nil
}

// MarshalGQL writes the given enum to the supplied writer as a quoted string
func (e AppointmentParticipantRequiredEnum) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// AppointmentParticipantStatusEnum is a FHIR enum
type AppointmentParticipantStatusEnum string

const (
	// AppointmentParticipantStatusEnumAccepted ...
	AppointmentParticipantStatusEnumAccepted AppointmentParticipantStatusEnum = "accepted"
	// AppointmentParticipantStatusEnumDeclined ...
	AppointmentParticipantStatusEnumDeclined AppointmentParticipantStatusEnum = "declined"
	// AppointmentParticipantStatusEnumTentative ...
	AppointmentParticipantStatusEnumTentative AppointmentParticipantStatusEnum = "tentative"
	// AppointmentParticipantStatusEnumNeedsAction ...
	AppointmentParticipantStatusEnumNeedsAction AppointmentParticipantStatusEnum = "needs_action"
)

// AllAppointmentParticipantStatusEnum ...
var AllAppointmentParticipantStatusEnum = []AppointmentParticipantStatusEnum{
	AppointmentParticipantStatusEnumAccepted,
	AppointmentParticipantStatusEnumDeclined,
	AppointmentParticipantStatusEnumTentative,
	AppointmentParticipantStatusEnumNeedsAction,
}

// IsValid ...
func (e AppointmentParticipantStatusEnum) IsValid() bool {
	switch e {
	case AppointmentParticipantStatusEnumAccepted, AppointmentParticipantStatusEnumDeclined, AppointmentParticipantStatusEnumTentative, AppointmentParticipantStatusEnumNeedsAction:
		return true
	}
	return false
}

// String ...
func (e AppointmentParticipantStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *AppointmentParticipantStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AppointmentParticipantStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Appointment_ParticipantStatusEnum", str)
	}
	return nil
}

// MarshalGQL writes the appointment participant status to the supplied writer as a quoted string
func (e AppointmentParticipantStatusEnum) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}
