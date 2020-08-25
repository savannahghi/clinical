package clinical

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/base"
)

// SingleFHIRReferencePayload - compose a single FHIRReference input
func SingleFHIRReferencePayload(r ReferenceInput) *FHIRReferenceInput {
	identifier := SingleIdentifierPayload(r.Identifier)
	return &FHIRReferenceInput{
		Reference:  &r.Reference,
		Type:       &r.URL,
		Identifier: identifier,
		Display:    r.Display,
	}
}

// SingleEpisodeofcareStatushistoryPayload - compose an FHIREpisodeofcareStatushistory payload
func SingleEpisodeofcareStatushistoryPayload() *FHIREpisodeofcareStatushistoryInput {
	var status EpisodeOfCareStatusHistoryStatusEnum = "planned"
	return &FHIREpisodeofcareStatushistoryInput{
		Status: &status,
		Period: SingleFHIRPeriodPayload(),
	}
}

// EpisodeofcareStatushistoryPayload - compose an FHIREpisodeofcareStatushistory payload
func EpisodeofcareStatushistoryPayload() []*FHIREpisodeofcareStatushistoryInput {
	episodeOfCare := SingleEpisodeofcareStatushistoryPayload()
	return []*FHIREpisodeofcareStatushistoryInput{episodeOfCare}
}

// SingleEpisodeofcareDiagnosisPayload - compose a SingleEpisodeofcareDiagnosis
func SingleEpisodeofcareDiagnosisPayload() *FHIREpisodeofcareDiagnosisInput {
	// rank := 1
	display := "Admission diagnosis"
	var code base.Code = "AD"
	text := "Diagnosis"

	codeableConceptInput := SingleCodeableConceptPayload(code, display, text)

	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Chlamydia Trachomatis",
		System:        "http://snomed.info/sct",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "260385009",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "Test Organisation",
		Identifier: &identifier,
	}

	return &FHIREpisodeofcareDiagnosisInput{
		Condition: SingleFHIRReferencePayload(ref),
		Role:      codeableConceptInput,
		// Rank:      &rank, // expects an int
	}
}

// EpisodeOfCareType - compose an Episode of care type
func EpisodeOfCareType() []*FHIRCodeableConceptInput {
	display := "Post Acute Care"
	var code base.Code = "pac"
	text := "EpisodeOfCare"
	codeableConcept := SingleCodeableConceptPayload(code, display, text)
	return []*FHIRCodeableConceptInput{codeableConcept}

}

func episodeOfCarePatient(patientID string) *FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "+254723002959",
		System:        "phone",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "msisdn",
	}
	patientReference := "Patient/" + patientID
	ref := ReferenceInput{
		Reference:  patientReference,
		URL:        "https://healthcloud.co.ke",
		Display:    "Test User",
		Identifier: &identifier,
	}
	return SingleFHIRReferencePayload(ref)
}

func episodeOfCareOrganisation() *FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Healthcare Provider",
		System:        "organisation",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "prov",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "2050", // sladecode
		Identifier: &identifier,
	}
	return SingleFHIRReferencePayload(ref)
}

func episodeCareManager() *FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Doctor",
		System:        "http://terminology.hl7.org/CodeSystem/practitioner-role",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "doctor",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "Doctor",
		Identifier: &identifier,
	}
	return SingleFHIRReferencePayload(ref)
}

func referralRequestPayload() []*FHIRReferenceInput {
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
		Display:    "counseling",
		Identifier: &identifier,
	}

	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func teamPayload() []*FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "Primary physician",
		System:        "Primary physician",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "Primary physician",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "Primary physician",
		Identifier: &identifier,
	}

	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func AccountPayload() []*FHIRReferenceInput {
	identifier := SingleIdentifierInput{
		IdentifierUse: "official",
		Value:         "cash",
		System:        "cash",
		UserSelected:  true,
		Version:       "0.1",
		Code:          "ACCTRECEIVABLE",
	}
	ref := ReferenceInput{
		Reference:  "https://healthcloud.co.ke",
		URL:        "https://healthcloud.co.ke",
		Display:    "account receivable",
		Identifier: &identifier,
	}

	return []*FHIRReferenceInput{SingleFHIRReferencePayload(ref)}
}

func getEpisodeOfCarePayload(patientID string) FHIREpisodeOfCareInput {
	var status EpisodeOfCareStatusEnum = "active"

	return FHIREpisodeOfCareInput{
		Identifier:           IdentifierPayload(),
		Status:               &status,
		StatusHistory:        EpisodeofcareStatushistoryPayload(),
		Type:                 EpisodeOfCareType(),
		Diagnosis:            []*FHIREpisodeofcareDiagnosisInput{SingleEpisodeofcareDiagnosisPayload()},
		Patient:              episodeOfCarePatient(patientID),
		ManagingOrganization: episodeOfCareOrganisation(),
		Period:               SingleFHIRPeriodPayload(),
		ReferralRequest:      referralRequestPayload(),
		CareManager:          episodeCareManager(),
		Team:                 teamPayload(),
		Account:              AccountPayload(),
	}
}

// CreateFHIREpisodeOfCarePayload - create an episode of care
func CreateFHIREpisodeOfCarePayload(t *testing.T) FHIREpisodeOfCareRelayPayload {
	service := NewService()
	createdPatient := CreateTestFHIRPatient(t)
	patientID := *createdPatient.Resource.ID
	episodeOfCarePayload := getEpisodeOfCarePayload(patientID)
	ctx := context.Background()
	ep, err := service.CreateFHIREpisodeOfCare(ctx, episodeOfCarePayload)
	if err != nil {
		t.Fatalf("unable to episode of care resource %s: ", err)
	}
	return *ep
}

// GetTestFHIREpisodeOfCare - retrieve a created test patient in FHIR
func GetTestFHIREpisodeOfCare(t *testing.T, id string) FHIREpisodeOfCareRelayPayload {
	service := NewService()
	ctx := context.Background()
	ep, err := service.GetFHIREpisodeOfCare(ctx, id)
	if err != nil {
		t.Fatalf("unable to retrieve the episode of care %s: ", err)
	}
	return *ep
}
