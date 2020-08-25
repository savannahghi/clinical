package clinical

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

// SingleEncounterStatushistoryPayload - compose an encounter status history payload
func SingleEncounterStatushistoryPayload() FHIREncounterStatushistoryInput {
	var status EncounterStatusHistoryStatusEnum = "planned"
	return FHIREncounterStatushistoryInput{
		Status: &status,
		Period: SingleFHIRPeriodPayload(),
	}
}

func encounterStatushistoryPayloadSecond() *FHIREncounterStatushistoryInput {
	var status EncounterStatusHistoryStatusEnum = "planned"
	return &FHIREncounterStatushistoryInput{
		Status: &status,
		Period: SingleFHIRPeriodPayload(),
	}
}

// EncounterStatushistoryPayload - compose a list of FHIRReference inputs
func EncounterStatushistoryPayload() []*FHIREncounterStatushistoryInput {
	return []*FHIREncounterStatushistoryInput{encounterStatushistoryPayloadSecond()}
}

// EncounterFHIRCodingPayload - compose an encounter coding input
func EncounterFHIRCodingPayload() FHIRCodingInput {
	display := "inpatient acute"
	var code base.Code = "ACUTE"

	userSelected := true
	var system base.URI = "http://terminology.hl7.org/CodeSystem/v2-0131"
	version := "2.0"

	return FHIRCodingInput{
		System:       &system,
		Code:         code,
		Version:      &version,
		Display:      display,
		UserSelected: &userSelected,
	}
}

func encounterCodingPayload() *FHIRCodingInput {
	display := "inpatient acute"
	var code base.Code = "ACUTE"

	userSelected := true
	var system base.URI = "http://terminology.hl7.org/CodeSystem/v2-0131"
	version := "2.0"

	return &FHIRCodingInput{
		System:       &system,
		Code:         code,
		Version:      &version,
		Display:      display,
		UserSelected: &userSelected,
	}
}

// SigleEncounterClasshistoryPayload - compose a n encounter class history payload
func SingleEncounterClasshistoryPayload() FHIREncounterClasshistoryInput {
	return FHIREncounterClasshistoryInput{
		Class:  encounterCodingPayload(),
		Period: SingleFHIRPeriodPayload(),
	}
}

func encounterClasshistoryPayload() *FHIREncounterClasshistoryInput {
	return &FHIREncounterClasshistoryInput{
		Class:  encounterCodingPayload(),
		Period: SingleFHIRPeriodPayload(),
	}
}

// EncounterClasshistoryPayload - compose an encounter class history payload
func EncounterClasshistoryPayload() []*FHIREncounterClasshistoryInput {
	return []*FHIREncounterClasshistoryInput{encounterClasshistoryPayload()}
}

//
func EncounterTypePayload() []*FHIRCodeableConceptInput {
	display := "Outpatient Kenacort injection"
	var code base.Code = "KI"
	text := "Encounter"
	codeableConcept := SingleCodeableConceptPayload(code, display, text)
	return []*FHIRCodeableConceptInput{codeableConcept}
}

func EncounterServiceTypePayload() *FHIRCodeableConceptInput {
	display := "Adoption/Permanent Care Info/Support"
	var code base.Code = "1"
	text := "Encounter"
	return SingleCodeableConceptPayload(code, display, text)
}

// EncounterPriorityPayload - compose an encounter prioroty payload
func EncounterPriorityPayload() *FHIRCodeableConceptInput {
	display := "urgent"
	var code base.Code = "UR"
	text := "Calls for prompt action."
	return SingleCodeableConceptPayload(code, display, text)
}

// EncounterParticipantPayload - compose an FHIR participant payload
func EncounterParticipantPayload() *FHIREncounterParticipantInput {
	display := "consultant"
	var code base.Code = "CON"
	text := "An advisor participating in the service."
	codeableConcept := SingleCodeableConceptPayload(code, display, text)

	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "admitter",
		System:        "http://terminology.hl7.org/CodeSystem/v3-ParticipationType",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "ADM",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "admitter",
		Identifier: &identifier,
	}

	return &FHIREncounterParticipantInput{
		Type:       []*FHIRCodeableConceptInput{codeableConcept},
		Period:     SingleFHIRPeriodPayload(),
		Individual: SingleFHIRReferencePayload(ref),
	}

}

// SingleEncounterDiagnosisPayload - compose a single encounter diagnosis payload
func SingleEncounterDiagnosisPayload() *FHIREncounterDiagnosisInput {
	// rank := "1" FHIR API expects an int
	display := "Admission diagnosis"
	var code base.Code = "AD"
	text := "admission."
	codeableConcept := SingleCodeableConceptPayload(code, display, text)

	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Amputee-BKA",
		System:        "Amputee-BKA",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "BKA",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "Amputee-BKA",
		Identifier: &identifier,
	}

	return &FHIREncounterDiagnosisInput{
		Condition: SingleFHIRReferencePayload(ref),
		Use:       codeableConcept,
		// Rank:      &rank, FHIR API expects an int
	}
}

func encounterAdmitSource() *FHIRCodeableConceptInput {
	display := "From outpatient department"
	var code base.Code = "outp"
	text := "The patient has been transferred from an outpatient department within the hospital."
	return SingleCodeableConceptPayload(code, display, text)
}

func encounterReAdmission() *FHIRCodeableConceptInput {
	display := "Re-admission"
	var code base.Code = "R"
	text := "Re-admission"
	return SingleCodeableConceptPayload(code, display, text)
}

func encounterDietPreference() *FHIRCodeableConceptInput {
	display := "Vegetarian"
	var code base.Code = "vegetarian"
	text := "Vegetarian"
	return SingleCodeableConceptPayload(code, display, text)
}

func encounterSpecialCourtesy() *FHIRCodeableConceptInput {
	display := "professional courtesy"
	var code base.Code = "PRF"
	text := "professional courtesy"
	return SingleCodeableConceptPayload(code, display, text)
}

func encounterSpecialArrangement() *FHIRCodeableConceptInput {
	display := "Wheelchair"
	var code base.Code = "wheel"
	text := "Wheelchair"
	return SingleCodeableConceptPayload(code, display, text)
}

func encounterDischargeDisposition() *FHIRCodeableConceptInput {
	display := "Home"
	var code base.Code = "home"
	text := "The patient was dicharged and has indicated that they are going to return home afterwards."
	return SingleCodeableConceptPayload(code, display, text)
}

func encounterOrigin() *FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Ambulance",
		System:        "Ambulance",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "Ambulance",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "Ambulance",
		Identifier: &identifier,
	}
	return SingleFHIRReferencePayload(ref)
}

func encounterDestination() *FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Ward",
		System:        "Ward",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "Ward",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "Ward",
		Identifier: &identifier,
	}
	return SingleFHIRReferencePayload(ref)
}

// EncounterHospitalizationPayload -  compose an encounter hospitalization payload
func EncounterHospitalizationPayload() *FHIREncounterHospitalizationInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "+254723002959",
		System:        "phone",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "msisdn",
	}
	return &FHIREncounterHospitalizationInput{
		PreAdmissionIdentifier: SingleIdentifierPayload(&identifier),
		Origin:                 encounterOrigin(),
		AdmitSource:            encounterAdmitSource(),
		ReAdmission:            encounterReAdmission(),
		DietPreference:         []*FHIRCodeableConceptInput{encounterDietPreference()},
		SpecialCourtesy:        []*FHIRCodeableConceptInput{encounterSpecialCourtesy()},
		SpecialArrangement:     []*FHIRCodeableConceptInput{encounterSpecialArrangement()},
		Destination:            encounterDestination(),
		DischargeDisposition:   encounterDischargeDisposition(),
	}
}

// EncounterLocationPayload - compose an encounter location payload
func EncounterLocationPayload() *FHIREncounterLocationInput {
	var status EncounterLocationStatusEnum = "active"
	display := "Building"
	var code base.Code = "bu"
	text := "Any Building or structure"
	codeableConcept := SingleCodeableConceptPayload(code, display, text)

	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Building",
		System:        "Building",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "Building",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "Building",
		Identifier: &identifier,
	}

	return &FHIREncounterLocationInput{
		Location:     SingleFHIRReferencePayload(ref),
		Status:       &status,
		PhysicalType: codeableConcept,
		Period:       SingleFHIRPeriodPayload(),
	}
}

func encounterEpisodeOfCare(episodeID string) []*FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Visit",
		System:        "Visit",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "Visit",
	}
	epReference := "EpisodeOfCare/" + episodeID
	ref := ReferenceInput{
		Reference:  epReference,
		URL:        "EpisodeOfCare",
		Display:    "Visit",
		Identifier: &identifier,
	}
	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func encounterSubject(patientID string) *FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "+254723002959",
		System:        "phone",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "msisdn",
	}
	ref := ReferenceInput{
		Reference:  patientID,
		URL:        "Patient",
		Display:    "Patient",
		Identifier: &identifier,
	}
	return SingleFHIRReferencePayload(ref)
}

func encounterBasedOn() []*FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "counseling",
		System:        "counseling",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "counseling",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "Counseling",
		Identifier: &identifier,
	}
	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func encounterAppointment() []*FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "appointment",
		System:        "appointment",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "appointment",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "appointment",
		Identifier: &identifier,
	}
	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func encounterReason() []*FHIRReferenceInput {
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

func encounterServiceProvider() *FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "XYZ Clinic",
		System:        "Hospital Facility",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "001",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "XYZ Clinic",
		Identifier: &identifier,
	}
	return SingleFHIRReferencePayload(ref)
}

// Another Encounter this encounter is part of
func encounterPartOf() *FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Another Encounter",
		System:        "Another Encounter",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "Another Encounter",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "Another Encounter",
		Identifier: &identifier,
	}
	return SingleFHIRReferencePayload(ref)
}

// GetEncounterPayload - compose an FHIR Encounter payload
func GetEncounterPayload(ep FHIREpisodeOfCareRelayPayload) FHIREncounterInput {
	var status EncounterStatusEnum = "planned"
	// var reason base.Code = "Spontaneous abortion with laceration of cervix"
	patientID := *ep.Resource.Patient.Reference
	episodeID := *ep.Resource.ID
	return FHIREncounterInput{
		Identifier:    IdentifierPayload(),
		Status:        status,
		StatusHistory: EncounterStatushistoryPayload(),
		Class:         EncounterFHIRCodingPayload(),
		ClassHistory:  EncounterClasshistoryPayload(),
		Type:          EncounterTypePayload(),
		ServiceType:   EncounterServiceTypePayload(),
		Priority:      EncounterPriorityPayload(),
		Subject:       encounterSubject(patientID),       // The patient
		EpisodeOfCare: encounterEpisodeOfCare(episodeID), //Episode(s) of care that this encounter should be recorded against
		BasedOn:       encounterBasedOn(),
		Participant:   []*FHIREncounterParticipantInput{EncounterParticipantPayload()},
		Appointment:   encounterAppointment(),
		Period:        SingleFHIRPeriodPayload(),
		// Length:        encounterDuration(), // if enabled FHIR API  returns length exists error
		// ReasonCode:      &reason, // if enabled FHIR API expects an array
		ReasonReference: encounterReason(),
		Diagnosis:       []*FHIREncounterDiagnosisInput{SingleEncounterDiagnosisPayload()},
		Account:         AccountPayload(),
		Hospitalization: EncounterHospitalizationPayload(),
		Location:        []*FHIREncounterLocationInput{EncounterLocationPayload()},
		ServiceProvider: encounterServiceProvider(),
		PartOf:          encounterPartOf(),
	}
}

/**
 When an EpisodeOfCare has been started/created
 Given one or more encounters have been added to it
 When i search encounters for the specified patient
 Then it should return a number of encounters
**/
func TestService_CreateFHIREncounter(t *testing.T) {
	service := NewService()
	ep := CreateFHIREpisodeOfCarePayload(t)
	encounterPayload := GetEncounterPayload(ep)
	ctx := context.Background()
	encounter, err := service.CreateFHIREncounter(ctx, encounterPayload)
	if err != nil {
		t.Fatalf("unable to create encounter resource %s: ", err)
	}
	assert.NotNil(t, encounter)
	encounterSearchParams := map[string]interface{}{
		"patient": fmt.Sprintf(*ep.Resource.Patient.Reference),
		"_sort":   "date",
		"count":   "1",
	}
	encounterConnection, err := service.SearchFHIREncounter(ctx, encounterSearchParams)
	if err != nil {
		t.Fatalf("unable to search patient encounter resource %s: ", err)
	}
	if len(encounterConnection.Edges) == 0 {
		t.Fatalf("unable to get patient encounter resource %s: ", err)
	}
}
