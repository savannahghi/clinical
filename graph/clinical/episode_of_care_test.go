package clinical

import (
	"bytes"
	"context"
	"reflect"
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
	text := "FULL_ACCESS" // This should be the access level type
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

// getEpisodeOfCarePayload creates FHIREpisodeOfCareInput
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

func TestEpisodeOfCareStatusEnum_IsValid(t *testing.T) {
	const unplanned EpisodeOfCareStatusEnum = "unplanned"
	tests := []struct {
		name string
		e    EpisodeOfCareStatusEnum
		want bool
	}{
		{
			name: "pass case",
			e:    EpisodeOfCareStatusEnumActive,
			want: true,
		},
		{
			name: "fail case",
			e:    unplanned,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEpisodeOfCareStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     EpisodeOfCareStatusEnum
		wantW string
	}{
		{
			name:  "pass case for EpisodeOfCareStatusEnumPlanned",
			e:     EpisodeOfCareStatusEnumPlanned,
			wantW: "\"planned\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusEnumWaitlist",
			e:     EpisodeOfCareStatusEnumWaitlist,
			wantW: "\"waitlist\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusEnumActive",
			e:     EpisodeOfCareStatusEnumActive,
			wantW: "\"active\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusEnumOnhold",
			e:     EpisodeOfCareStatusEnumOnhold,
			wantW: "\"onhold\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusEnumFinished",
			e:     EpisodeOfCareStatusEnumFinished,
			wantW: "\"finished\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusEnumCancelled",
			e:     EpisodeOfCareStatusEnumCancelled,
			wantW: "\"cancelled\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusEnumEnteredInError",
			e:     EpisodeOfCareStatusEnumEnteredInError,
			wantW: "\"entered_in_error\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestEpisodeOfCareStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    EpisodeOfCareStatusEnum
		want string
	}{
		{
			name: "pass case for planned",
			e:    EpisodeOfCareStatusEnumPlanned,
			want: "planned",
		},

		{
			name: "pass case for waitllist",
			e:    EpisodeOfCareStatusEnumWaitlist,
			want: "waitlist",
		},

		{
			name: "pass case for active",
			e:    EpisodeOfCareStatusEnumActive,
			want: "active",
		},

		{
			name: "pass case for onhold",
			e:    EpisodeOfCareStatusEnumOnhold,
			want: "onhold",
		},

		{
			name: "pass case for finished",
			e:    EpisodeOfCareStatusEnumFinished,
			want: "finished",
		},

		{
			name: "pass case for cancelled",
			e:    EpisodeOfCareStatusEnumCancelled,
			want: "cancelled",
		},

		{
			name: "pass case for entered_in_error",
			e:    EpisodeOfCareStatusEnumEnteredInError,
			want: "entered_in_error",
		},

		{
			name: "Fail case for unknown enum",
			e:    "unknown",
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEpisodeOfCareStatusEnum_UnmarshalGQL(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       EpisodeOfCareStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "pass case",
			e:    EpisodeOfCareStatusEnumPlanned,
			args: args{
				string("planned"),
			},
			wantErr: false,
		},

		{
			name: "fail case",
			e:    EpisodeOfCareStatusEnumPlanned,
			args: args{
				string("fail"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestEpisodeOfCareStatusHistoryStatusEnum_IsValid
func TestEpisodeOfCareStatusHistoryStatusEnum_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    EpisodeOfCareStatusHistoryStatusEnum
		want bool
	}{
		{
			name: "EpisodeOfCareStatusHistoryStatusEnumPlanned: valid case",
			e:    EpisodeOfCareStatusHistoryStatusEnumPlanned,
			want: true,
		},
		{
			name: "EpisodeOfCareStatusHistoryStatusEnumWaitlist: valid case",
			e:    EpisodeOfCareStatusHistoryStatusEnumWaitlist,
			want: true,
		},

		{
			name: "EpisodeOfCareStatusHistoryStatusEnumActive: valid case",
			e:    EpisodeOfCareStatusHistoryStatusEnumActive,
			want: true,
		},

		{
			name: "EpisodeOfCareStatusHistoryStatusEnumOnhold: valid case",
			e:    EpisodeOfCareStatusHistoryStatusEnumOnhold,
			want: true,
		},

		{
			name: "EpisodeOfCareStatusHistoryStatusEnumFinished: valid case",
			e:    EpisodeOfCareStatusHistoryStatusEnumFinished,
			want: true,
		},

		{
			name: "EpisodeOfCareStatusHistoryStatusEnumCancelled: valid case",
			e:    EpisodeOfCareStatusHistoryStatusEnumCancelled,
			want: true,
		},

		{
			name: "EpisodeOfCareStatusHistoryStatusEnumEnteredInError: valid case",
			e:    EpisodeOfCareStatusHistoryStatusEnumEnteredInError,
			want: true,
		},

		{
			name: "EpisodeOfCareStatusHistoryStatusEnum: invalid case",
			e:    "fail",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEpisodeOfCareStatusHistoryStatusEnum_MarshalGQL
func TestEpisodeOfCareStatusHistoryStatusEnum_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     EpisodeOfCareStatusHistoryStatusEnum
		wantW string
	}{
		{
			name:  "pass case for EpisodeOfCareStatusHistoryStatusEnumPlanned",
			e:     EpisodeOfCareStatusHistoryStatusEnumPlanned,
			wantW: "\"planned\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusHistoryStatusEnumWaitlist",
			e:     EpisodeOfCareStatusHistoryStatusEnumWaitlist,
			wantW: "\"waitlist\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusHistoryStatusEnumActive",
			e:     EpisodeOfCareStatusHistoryStatusEnumActive,
			wantW: "\"active\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusHistoryStatusEnumOnhold",
			e:     EpisodeOfCareStatusHistoryStatusEnumOnhold,
			wantW: "\"onhold\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusHistoryStatusEnumFinished",
			e:     EpisodeOfCareStatusHistoryStatusEnumFinished,
			wantW: "\"finished\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusHistoryStatusEnumCancelled",
			e:     EpisodeOfCareStatusHistoryStatusEnumCancelled,
			wantW: "\"cancelled\"",
		},

		{
			name:  "pass case for EpisodeOfCareStatusHistoryStatusEnumEnteredInError",
			e:     EpisodeOfCareStatusHistoryStatusEnumEnteredInError,
			wantW: "\"entered_in_error\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

// TestEpisodeOfCareStatusHistoryStatusEnum_String
func TestEpisodeOfCareStatusHistoryStatusEnum_String(t *testing.T) {
	tests := []struct {
		name string
		e    EpisodeOfCareStatusHistoryStatusEnum
		want string
	}{
		{
			name: "pass case for planned",
			e:    EpisodeOfCareStatusHistoryStatusEnumPlanned,
			want: "planned",
		},

		{
			name: "pass case for EpisodeOfCareStatusHistoryStatusEnumWaitlist",
			e:    EpisodeOfCareStatusHistoryStatusEnumWaitlist,
			want: "waitlist",
		},

		{
			name: "pass case for EpisodeOfCareStatusHistoryStatusEnumActive",
			e:    EpisodeOfCareStatusHistoryStatusEnumActive,
			want: "active",
		},

		{
			name: "pass case for EpisodeOfCareStatusHistoryStatusEnumOnhold",
			e:    EpisodeOfCareStatusHistoryStatusEnumOnhold,
			want: "onhold",
		},

		{
			name: "pass case for EpisodeOfCareStatusHistoryStatusEnumFinished",
			e:    EpisodeOfCareStatusHistoryStatusEnumFinished,
			want: "finished",
		},

		{
			name: "pass case for EpisodeOfCareStatusHistoryStatusEnumCancelled",
			e:    EpisodeOfCareStatusHistoryStatusEnumCancelled,
			want: "cancelled",
		},

		{
			name: "pass case for EpisodeOfCareStatusHistoryStatusEnumEnteredInError",
			e:    EpisodeOfCareStatusHistoryStatusEnumEnteredInError,
			want: "entered_in_error",
		},

		{
			name: "Fail case for unknown enum",
			e:    "unknown",
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEpisodeOfCareStatusHistoryStatusEnum_UnmarshalGQL
func TestEpisodeOfCareStatusHistoryStatusEnum_UnmarshalGQL(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       EpisodeOfCareStatusHistoryStatusEnum
		args    args
		wantErr bool
	}{
		{
			name: "pass case",
			e:    EpisodeOfCareStatusHistoryStatusEnumPlanned,
			args: args{
				string("planned"),
			},
			wantErr: false,
		},

		{
			name: "fail case",
			e:    EpisodeOfCareStatusHistoryStatusEnumPlanned,
			args: args{
				string("fail"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestService_CreateFHIREpisodeOfCare
//		create an episodeOfCare then fetch the created episodeOfCARE using it's ID
func TestService_CreateFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	patient := CreateTestFHIRPatient(t)

	type args struct {
		input FHIREpisodeOfCareInput
	}
	tests := []struct {
		name    string
		args    args
		want    *FHIREpisodeOfCareRelayPayload
		wantErr bool
	}{
		{
			name: "pass case",
			args: args{
				input: getEpisodeOfCarePayload(*patient.Resource.ID),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fhirEpisodeOfCareRelayPayload, err := service.CreateFHIREpisodeOfCare(ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fhirEpisodeOfCareID := *fhirEpisodeOfCareRelayPayload.Resource.ID

			episodeOfCareFetched, err := service.GetFHIREpisodeOfCare(ctx, fhirEpisodeOfCareID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(fhirEpisodeOfCareRelayPayload, episodeOfCareFetched) {
				t.Errorf("CreateFHIREpisodeOfCare() fhirEpisodeOfCareRelayPayload = %v, want %v", fhirEpisodeOfCareRelayPayload, tt.want)
			}
		})
	}
}

// TestService_DeleteFHIREpisodeOfCare
//		create an episode of care.
//		fetch the episode of care and confirm it exits.
//		Delete the episode of care.
//		try to fetch the deleted episode of care.
func TestService_DeleteFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	service := NewService()

	// create episode of care.
	patient := CreateTestFHIRPatient(t)
	episodeOfCareInput := getEpisodeOfCarePayload(*patient.Resource.ID)
	fhirEpisodeOfCareRelayPayload, err := service.CreateFHIREpisodeOfCare(ctx, episodeOfCareInput)
	if err != nil {
		t.Errorf("CreateFHIREpisodeOfCare() error = %v", err)
		return
	}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Passes",
			args:    args{id: *fhirEpisodeOfCareRelayPayload.Resource.ID},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.DeleteFHIREpisodeOfCare(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteFHIREpisodeOfCare() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestService_GetFHIREpisodeOfCare
//	create FHIR episode of care.
// 		fetch the created episode of care.
// the values of the fetched episode of care should be equal to the created episode of care.
func TestService_GetFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	service := NewService()

	// create episode of care.
	patient := CreateTestFHIRPatient(t)
	episodeOfCareInput := getEpisodeOfCarePayload(*patient.Resource.ID)
	fhirEpisodeOfCareRelayPayload, err := service.CreateFHIREpisodeOfCare(ctx, episodeOfCareInput)
	if err != nil {
		t.Errorf("CreateFHIREpisodeOfCare() error = %v", err)
		return
	}

	tests := []struct {
		name    string
		id      string
		want    *FHIREpisodeOfCareRelayPayload
		wantErr bool
	}{
		{
			name:    "Pass case",
			id:      *fhirEpisodeOfCareRelayPayload.Resource.ID,
			want:    fhirEpisodeOfCareRelayPayload,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetFHIREpisodeOfCare(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFHIREpisodeOfCare() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestService_SearchFHIREpisodeOfCare
// 	To test this function:
//		create an episode of care.
// 		search episodes of care.
// 		the episode we created should be among the episodes returned.
func TestService_SearchFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	service := NewService()

	// create episode of care.
	patient := CreateTestFHIRPatient(t)
	episodeOfCareInput := getEpisodeOfCarePayload(*patient.Resource.ID)
	_, err := service.CreateFHIREpisodeOfCare(ctx, episodeOfCareInput)
	if err != nil {
		t.Errorf("CreateFHIREpisodeOfCare() error = %v", err)
		return
	}
	patientID := *patient.Resource.ID
	episodeOfCareSearchParams := map[string]interface{}{
		"patient": patientID,
		"_sort":   "date",
		"_count":  "1",
	}

	type args struct {
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid_case",
			args:    args{params: episodeOfCareSearchParams},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SearchFHIREpisodeOfCare(ctx, tt.args.params)
			if (err != nil) != tt.wantErr || got == nil {
				t.Errorf("SearchFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.Edges) < 1 {
				t.Errorf("SearchFHIREpisodeOfCare() error: Returned zero episodes of care. ")
			}

		})
	}
}

// TestService_UpdateFHIREpisodeOfCare
func TestService_UpdateFHIREpisodeOfCare(t *testing.T) {
	ctx := context.Background()
	service := NewService()

	// create episode of care.
	patient := CreateTestFHIRPatient(t)
	episodeOfCareInput := getEpisodeOfCarePayload(*patient.Resource.ID)
	fhirEpisodeOfCareRelayPayload, err := service.CreateFHIREpisodeOfCare(ctx, episodeOfCareInput)
	if err != nil {
		t.Errorf("CreateFHIREpisodeOfCare() error = %v", err)
		return
	}

	type args struct {
		input FHIREpisodeOfCareInput
	}
	tests := []struct {
		name    string
		args    args
		want    *FHIREpisodeOfCareRelayPayload
		wantErr bool
	}{
		{
			name:    "Pass case",
			args:    args{input: episodeOfCareInput},
			want:    fhirEpisodeOfCareRelayPayload,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.input.ID = tt.want.Resource.ID
			got, err := service.UpdateFHIREpisodeOfCare(ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateFHIREpisodeOfCare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateFHIREpisodeOfCare() got = %v, want %v", got, tt.want)
			}
		})
	}
}
