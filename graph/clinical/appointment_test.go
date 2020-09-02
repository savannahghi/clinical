package clinical

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func appointmentCancelationReason() *FHIRCodeableConceptInput {
	display := "Patient: Canceled via automated reminder system"
	var code base.Code = "pat-crs"
	text := "Patient: Canceled via automated reminder system"
	return SingleCodeableConceptPayload(code, display, text)
}

func appointmentServiceCategory() *FHIRCodeableConceptInput {
	display := "General Practice"
	var code base.Code = "17"
	text := "General Practice/GP (doctor)"
	return SingleCodeableConceptPayload(code, display, text)
}

func appointmentServiceType() *FHIRCodeableConceptInput {
	display := "Case Management for Older Persons"
	var code base.Code = "5"
	text := "Case Management for Older Persons"
	return SingleCodeableConceptPayload(code, display, text)
}

func appointmentSpecialty() *FHIRCodeableConceptInput {
	display := "Clinical microbiology"
	var code base.Code = "408454008"
	text := "Clinical microbiology"
	return SingleCodeableConceptPayload(code, display, text)
}

func appointmentAppointmentType() *FHIRCodeableConceptInput {
	display := "en: A routine check-up, such as an annual physicalnl:"
	var code base.Code = "CHECKUP"
	text := "en: A routine check-up, such as an annual physicalnl:"
	return SingleCodeableConceptPayload(code, display, text)
}

func appointmentReasonReference() []*FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "reason",
		System:        "reason",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "reason",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "reason",
		Identifier: &identifier,
	}
	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func appointmentSupportingInformation() []*FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "reason",
		System:        "reason",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "reason",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "reason",
		Identifier: &identifier,
	}
	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func appointmentSlot() []*FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "reason",
		System:        "reason",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "reason",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "reason",
		Identifier: &identifier,
	}
	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func appointmentBasedOn() []*FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "reason",
		System:        "reason",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "reason",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "reason",
		Identifier: &identifier,
	}
	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func appointmentPeriod() *FHIRPeriodInput {
	var start base.DateTime = "2020-02-06"
	var end base.DateTime = "2020-02-08"
	return &FHIRPeriodInput{
		Start: start,
		End:   end,
	}
}

func appointmentPatientActor(patientID string) *FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "reason",
		System:        "reason",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "reason",
	}
	patientReference := "Patient/" + patientID
	ref := ReferenceInput{
		Reference:  patientReference,
		URL:        "https://healthcloud.co.ke",
		Display:    "Patient",
		Identifier: &identifier,
	}
	return SingleFHIRReferencePayload(ref)
}

func appointmentParticipantInput(actor *FHIRReferenceInput) *FHIRAppointmentParticipantInput {
	var required AppointmentParticipantRequiredEnum = "required"
	var status AppointmentParticipantStatusEnum = "accepted"
	return &FHIRAppointmentParticipantInput{
		Type:     []*FHIRCodeableConceptInput{appointmentParticipantType()},
		Actor:    actor,
		Required: &required,
		Status:   &status,
		Period:   appointmentPeriod(),
	}
}

func appointmentParticipantType() *FHIRCodeableConceptInput {
	display := "callback contact"
	var code base.Code = "CALLBCK"
	text := "A person or organization who should be contacted for follow-up questions about the act in place of the author."
	return SingleCodeableConceptPayload(code, display, text)
}

func appointmentIdentifierPayload(organizationID *string) []*FHIRIdentifierInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         *organizationID,
		System:        "organization",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "org",
	}
	return []*FHIRIdentifierInput{SingleIdentifierPayload(&identifier)}
}

func getAppointmentPayload(t *testing.T, organizationID *string) FHIRAppointmentInput {

	var status AppointmentStatusEnum = "booked"
	var start base.Instant = "2020-02-07T13:28:17.239+02:00"
	var end base.Instant = "2020-02-07T16:28:17.239+02:00"
	var created base.DateTime = "2020-02-06"
	// var reasonCode base.Code = "109006"
	var priority = 1
	var description = "appointment"
	var minutesDuration = 20
	var comment = "to visit after nect week"
	var patientInstruction = "patient instruction"
	createdPatient := CreateTestFHIRPatient(t)
	patientID := *createdPatient.Resource.ID
	appointmentIdentifier := appointmentIdentifierPayload(organizationID)
	patient := appointmentParticipantInput(appointmentPatientActor(patientID))

	return FHIRAppointmentInput{
		Identifier:        appointmentIdentifier,
		Status:            &status,
		CancelationReason: appointmentCancelationReason(),
		ServiceCategory:   []*FHIRCodeableConceptInput{appointmentServiceCategory()},
		ServiceType:       []*FHIRCodeableConceptInput{appointmentServiceType()},
		Specialty:         []*FHIRCodeableConceptInput{appointmentSpecialty()},
		AppointmentType:   appointmentAppointmentType(),
		// ReasonCode:            &reasonCode, // expects an array
		ReasonReference:       appointmentReasonReference(),
		Priority:              &priority,
		Description:           &description,
		SupportingInformation: appointmentSupportingInformation(),
		Start:                 &start,
		End:                   &end,
		MinutesDuration:       &minutesDuration,
		Slot:                  appointmentSlot(),
		Created:               &created,
		Comment:               &comment,
		PatientInstruction:    &patientInstruction,
		BasedOn:               appointmentBasedOn(),
		Participant:           []*FHIRAppointmentParticipantInput{patient},
		RequestedPeriod:       []*FHIRPeriodInput{appointmentPeriod()},
	}
}

// CreateFHIRAppointment - helper to create a test appointment in FHIR
func CreateFHIRAppointment(t *testing.T) FHIRAppointmentRelayPayload {
	service := NewService()
	ctx := context.Background()
	organizationID, err := service.GetORCreateOrganization(ctx, 123)
	if err != nil {
		t.Fatalf("unable to get or create organization resource %s: ", err)
	}
	appointmentPayload := getAppointmentPayload(t, organizationID)
	appointment, err := service.CreateFHIRAppointment(ctx, appointmentPayload)
	if err != nil {
		t.Fatalf("unable to create appointment resource %s: ", err)
	}
	return *appointment
}

func TestService_SearchFHIRAppointment(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	createdAppointment := CreateFHIRAppointment(t)

	identifiers := createdAppointment.Resource.Identifier
	value := ""
	for _, val := range identifiers {
		value += val.Value
	}
	validSearchParams := map[string]interface{}{
		"identifier": value,
	}

	invalidSearchParams := map[string]interface{}{
		"identifier": "123588",
	}

	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *FHIRAppointmentRelayConnection
		wantErr bool
	}{
		{
			name:    "Sucessfully search for an appointment",
			args:    args{ctx: ctx, params: validSearchParams},
			wantErr: false,
		},
		{
			name:    "search for an appointments returns empty",
			args:    args{ctx: ctx, params: invalidSearchParams},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SearchFHIRAppointment(tt.args.ctx, tt.args.params)
			if tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, len(got.Edges), 0)
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestService_ListAppointments(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	type args struct {
		ctx               context.Context
		providerSladeCode int
	}
	// ensure an  appointment exists for providerCode 123
	CreateFHIRAppointment(t)
	tests := []struct {
		name    string
		args    args
		want    *FHIRAppointmentRelayConnection
		wantErr bool
	}{
		{
			name:    "Sucessfully list appointments",
			args:    args{ctx: ctx, providerSladeCode: 123},
			wantErr: false,
		},
		{
			name:    "appointments missing",
			args:    args{ctx: ctx, providerSladeCode: 55458},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.ListAppointments(tt.args.ctx, tt.args.providerSladeCode)
			if tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, len(res.Edges), 0)
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, res)
			}

		})
	}
}
