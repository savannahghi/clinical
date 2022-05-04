package clinical_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

func TestClinicalUseCaseImpl_ProblemSummary(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, patient, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		MFLCode,
	)
	if err != nil {
		t.Errorf("cant get test episode of care: %v\n", err)
	}

	episodePayload, err := u.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
	}

	encounterInput, err := getTestEncounterInput(t, episodePayload)
	if err != nil {
		t.Errorf("unable to get episode: %v", err)
	}

	encounter, err := u.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to create FHIREncounter: %v", err)
	}

	input, err := createTestConditionInput(*encounter.Resource.ID, *patient.ID)
	if err != nil {
		t.Errorf("cant create condition: %v\n", err)
		return
	}

	condition, err := u.CreateFHIRCondition(ctx, *input)
	if err != nil {
		t.Errorf("failed to create fhir condition: %v", err)
	}

	input.ID = condition.Resource.ID

	type args struct {
		ctx       context.Context
		patientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				patientID: *patient.ID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.ProblemSummary(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.ProblemSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_VisitSummary(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	episode, _, err := createTestEpisodeOfCare(
		ctx,
		msisdn,
		false,
		MFLCode,
	)
	if err != nil {
		t.Errorf("cant create test episode of care: %v\n", err)
		return
	}

	episodePayload, err := u.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
		return
	}

	activeEpisodeStatus := domain.EpisodeOfCareStatusEnumActive
	activeEncounterStatus := domain.EncounterStatusEnumInProgress
	if episodePayload.Resource.Status.String() != activeEpisodeStatus.String() {
		t.Errorf("an encounter can only be started for an active episode")
		return
	}
	episodeRef := fmt.Sprintf("EpisodeOfCare/%s", *episodePayload.Resource.ID)

	now := time.Now()
	startTime := scalarutils.DateTime(now.Format("2006-01-02T15:04:05+03:00"))

	encounterClassCode := scalarutils.Code("AMB")
	encounterClassSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/v3-ActCode")
	encounterClassVersion := "2018-08-12"
	encounterClassDisplay := "ambulatory"
	encounterClassUserSelected := false

	encounterInput := domain.FHIREncounterInput{
		Status: activeEncounterStatus,
		Class: domain.FHIRCodingInput{
			System:       &encounterClassSystem,
			Version:      &encounterClassVersion,
			Code:         encounterClassCode,
			Display:      encounterClassDisplay,
			UserSelected: &encounterClassUserSelected,
		},
		Subject: &domain.FHIRReferenceInput{
			Reference: episodePayload.Resource.Patient.Reference,
			Display:   episodePayload.Resource.Patient.Display,
			Type:      episodePayload.Resource.Patient.Type,
		},
		EpisodeOfCare: []*domain.FHIRReferenceInput{
			{
				Reference: &episodeRef,
			},
		},
		ServiceProvider: &domain.FHIRReferenceInput{
			Display: episodePayload.Resource.ManagingOrganization.Display,
			Type:    episodePayload.Resource.ManagingOrganization.Type,
		},
		Period: &domain.FHIRPeriodInput{
			Start: startTime,
		},
	}

	enc, err := u.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to CreateFHIREncounter: %v", err)
	}

	type args struct {
		ctx         context.Context
		encounterID string
		count       int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				encounterID: *enc.Resource.ID,
				count:       0,
			},
			wantErr: false,
		},

		{
			name: "Sad case: nil encounter ID",
			args: args{
				ctx:   ctx,
				count: 0,
			},
			wantErr: true,
		},

		{
			name: "Sad case: no count",
			args: args{
				ctx:         ctx,
				encounterID: ksuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.VisitSummary(tt.args.ctx, tt.args.encounterID, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.VisitSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_PatientTimelineWithCount(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	normalized, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		t.Errorf("can't normalize phone number: %v \n", err)
		return
	}
	episode, _, err := createTestEpisodeOfCare(
		ctx,
		*normalized,
		false,
		MFLCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}

	episodePayload, err := u.GetFHIREpisodeOfCare(ctx, *episode.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *episode.ID, err)
		return
	}

	activeEpisodeStatus := domain.EpisodeOfCareStatusEnumActive
	activeEncounterStatus := domain.EncounterStatusEnumInProgress
	if episodePayload.Resource.Status.String() != activeEpisodeStatus.String() {
		t.Errorf("an encounter can only be started for an active episode")
		return
	}
	episodeRef := fmt.Sprintf("EpisodeOfCare/%s", *episodePayload.Resource.ID)

	now := time.Now()
	startTime := scalarutils.DateTime(now.Format("2006-01-02T15:04:05+03:00"))

	encounterClassCode := scalarutils.Code("AMB")
	encounterClassSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/v3-ActCode")
	encounterClassVersion := "2018-08-12"
	encounterClassDisplay := "ambulatory"
	encounterClassUserSelected := false

	encounterInput := domain.FHIREncounterInput{
		Status: activeEncounterStatus,
		Class: domain.FHIRCodingInput{
			System:       &encounterClassSystem,
			Version:      &encounterClassVersion,
			Code:         encounterClassCode,
			Display:      encounterClassDisplay,
			UserSelected: &encounterClassUserSelected,
		},
		Subject: &domain.FHIRReferenceInput{
			Reference: episodePayload.Resource.Patient.Reference,
			Display:   episodePayload.Resource.Patient.Display,
			Type:      episodePayload.Resource.Patient.Type,
		},
		EpisodeOfCare: []*domain.FHIRReferenceInput{
			{
				Reference: &episodeRef,
			},
		},
		ServiceProvider: &domain.FHIRReferenceInput{
			Display: episodePayload.Resource.ManagingOrganization.Display,
			Type:    episodePayload.Resource.ManagingOrganization.Type,
		},
		Period: &domain.FHIRPeriodInput{
			Start: startTime,
		},
	}

	encounter, err := u.CreateFHIREncounter(ctx, encounterInput)
	if err != nil {
		t.Errorf("unable to CreateFHIREncounter: %v", err)
	}

	_, err = u.GetFHIREncounter(ctx, *encounter.Resource.ID)
	if err != nil {
		t.Errorf("unable to GetFHIREncounter: %v", err)
	}

	type args struct {
		ctx       context.Context
		episodeID string
		count     int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				episodeID: *episode.ID,
				count:     0,
			},
			wantErr: false,
		},

		{
			name: "Sad case: nil episode ID",
			args: args{
				ctx:   ctx,
				count: 0,
			},
			wantErr: true,
		},

		{
			name: "Sad case: no count",
			args: args{
				ctx:       ctx,
				episodeID: ksuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.PatientTimelineWithCount(tt.args.ctx, tt.args.episodeID, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientTimelineWithCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_PatientSearch(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	searchPatient := "Test user"

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	_, err = u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		t.Errorf("unable to create patient: %v", err)
		return
	}

	type args struct {
		ctx    context.Context
		search string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				search: searchPatient,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				search: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.PatientSearch(tt.args.ctx, tt.args.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_ContactsToContactPointInput(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	_, err = u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		t.Errorf("unable to create patient: %v", err)
		return
	}

	phone := &domain.PhoneNumberInput{
		Msisdn:             interserviceclient.TestUserPhoneNumber,
		VerificationCode:   "300300",
		IsUssd:             true,
		CommunicationOptIn: false,
	}

	type args struct {
		ctx    context.Context
		phones []*domain.PhoneNumberInput
		emails []*domain.EmailInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				phones: []*domain.PhoneNumberInput{
					phone,
				},
				emails: []*domain.EmailInput{},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				phones: nil,
			},
			//should not return an error but a nil
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.ContactsToContactPointInput(tt.args.ctx, tt.args.phones, tt.args.emails)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.ContactsToContactPointInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClinicalUseCaseImpl_CreatePatient(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	ID := ksuid.New().String()

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	_, err = u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		t.Errorf("unable to create patient: %v", err)
		return
	}

	gender := domain.PatientGenderEnumMale
	identifier := &domain.FHIRIdentifierInput{
		ID: &ID,
	}

	active := true

	humanName := &domain.FHIRHumanNameInput{
		ID:  &ID,
		Use: domain.HumanNameUseEnumNickname,
	}

	type args struct {
		ctx   context.Context
		input domain.FHIRPatientInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: domain.FHIRPatientInput{
					ID: &ID,
					Identifier: []*domain.FHIRIdentifierInput{
						identifier,
					},
					Active: &active,
					Name: []*domain.FHIRHumanNameInput{
						humanName,
					},
					Gender: &gender,
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: domain.FHIRPatientInput{
					ID: &ID,
					Identifier: []*domain.FHIRIdentifierInput{
						identifier,
					},
					Active: &active,
					Name: []*domain.FHIRHumanNameInput{
						humanName,
					},
					Gender: nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.CreatePatient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.CreatePatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_FindPatientByID(t *testing.T) {
	ctx := context.Background()

	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	var patient domain.FHIRPatient
	p, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			p, err := u.FindPatientsByMSISDN(ctx, interserviceclient.TestUserPhoneNumber)
			if err != nil {
				t.Errorf("can't find existing patient by MSISDN: %v", err)
			}
			patient = *p.Edges[0].Node
		} else {
			t.Errorf("unable to create patient: %v", err)
			return
		}
	} else {
		patient = *p.PatientRecord
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  *patient.ID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.FindPatientByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.FindPatientByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_UpdatePatient(t *testing.T) {
	TestClinicalUseCaseImpl_DeleteFHIRPatientByPhone(t)

	ctx := context.Background()
	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		t.Errorf("unable to create patient: %v", err)
		return
	}

	updatePhone := interserviceclient.TestUserPhoneNumber
	otp, err := generateTestOTP(t, updatePhone)
	if err != nil {
		log.Errorf("unable to get verified phone number and OTP")
		return
	}

	phone := &domain.PhoneNumberInput{
		Msisdn:           updatePhone,
		VerificationCode: otp,
	}
	patientInput.PhoneNumbers = []*domain.PhoneNumberInput{phone}
	date := scalarutils.Date{
		Year:  1900,
		Month: 12,
		Day:   20,
	}

	type args struct {
		ctx   context.Context
		input domain.SimplePatientRegistrationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: domain.SimplePatientRegistrationInput{
					ID:           *patient.PatientRecord.ID,
					BirthDate:    date,
					PhoneNumbers: patientInput.PhoneNumbers,
					Gender:       "male",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: domain.SimplePatientRegistrationInput{
					ID:           "",
					PhoneNumbers: patientInput.PhoneNumbers,
					Gender:       "male",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.UpdatePatient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.UpdatePatient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_AddNextOfKin(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		t.Errorf("unable to create patient: %v", err)
		return
	}

	otherTestName := "other test names"
	name := &domain.NameInput{
		FirstName:  gofakeit.FirstName(),
		LastName:   gofakeit.LastName(),
		OtherNames: &otherTestName,
	}
	date := scalarutils.Date{
		Year:  1900,
		Month: 12,
		Day:   20,
	}

	type args struct {
		ctx   context.Context
		input domain.SimpleNextOfKinInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: domain.SimpleNextOfKinInput{
					PatientID: *patient.PatientRecord.ID,
					Names: []*domain.NameInput{
						name,
					},
					PhoneNumbers:      []*domain.PhoneNumberInput{},
					Emails:            []*domain.EmailInput{},
					PhysicalAddresses: []*domain.PhysicalAddress{},
					PostalAddresses:   []*domain.PostalAddress{},
					Gender:            "male",
					Relationship:      "",
					Active:            false,
					BirthDate:         date,
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.AddNextOfKin(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.AddNextOfKin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_AddNHIF(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	testInput := ksuid.New().String()

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		t.Errorf("unable to create patient: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		input *domain.SimpleNHIFInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: &domain.SimpleNHIFInput{
					PatientID:        *patient.PatientRecord.ID,
					MembershipNumber: ksuid.New().String(),
					FrontImageBase64: &testInput,
					RearImageBase64:  &testInput,
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.AddNHIF(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.AddNHIF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_CreateUpdatePatientExtraInformation(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	maritalStatus := ksuid.New().String()

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		t.Errorf("unable to create patient: %v", err)
		return
	}

	pid, err := u.FindPatientByID(ctx, *patient.PatientRecord.ID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		input domain.PatientExtraInformationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: domain.PatientExtraInformationInput{
					PatientID:     *pid.PatientRecord.ID,
					MaritalStatus: (*domain.MaritalStatus)(&maritalStatus),
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: domain.PatientExtraInformationInput{
					PatientID:     "",
					MaritalStatus: (*domain.MaritalStatus)(&maritalStatus),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.CreateUpdatePatientExtraInformation(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.CreateUpdatePatientExtraInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_AllergySummary(t *testing.T) {

	ctx := context.Background()
	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		t.Errorf("unable to create patient: %v", err)
		return
	}

	type args struct {
		ctx       context.Context
		patientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				patientID: *patient.PatientRecord.ID,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				patientID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.AllergySummary(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.AllergySummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("the error was not expected")
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("an error was expected: %v", err)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_DeleteFHIRPatientByPhone(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	_, err = u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		t.Errorf("unable to create patient: %v", err)
		return
	}

	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// {
		// 	name: "Happy case",
		// 	args: args{
		// 		ctx:         ctx,
		// 		phoneNumber: interserviceclient.TestUserPhoneNumber,
		// 	},
		// 	wantErr: false,
		// }, TODO: Investigate
		{
			name: "Sad case: empty phone",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid phone",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.DeleteFHIRPatientByPhone(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.DeleteFHIRPatientByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClinicalUseCaseImpl_PatientTimeline(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		// TODO: investigate register patient
		// t.Errorf("unable to create patient: %v", err)
		return
	}

	episode, _, err := createTestEpisodeOfCare(
		context.Background(),
		interserviceclient.TestUserPhoneNumber,
		false,
		MFLCode,
	)
	if err != nil {
		t.Errorf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := u.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	oInput, err := getFhirObservationInput(*patient.PatientRecord, encounterID)
	if err != nil {
		t.Errorf("failed to get fhir observation input: %v", err)
	}

	_, err = u.CreateFHIRObservation(ctx, *oInput)
	if err != nil {
		t.Errorf("failed to create fhir observation: %v", err)
	}

	msInput, err := getFhirMedicationStatementInput(*patient.PatientRecord)
	if err != nil {
		t.Errorf("failed to get fhir medication statement input: %v", err)
	}

	_, err = u.CreateFHIRMedicationStatement(ctx, *msInput)
	if err != nil {
		t.Errorf("failed to create fhir medication statement: %v", err)
	}

	type args struct {
		ctx       context.Context
		patientID string
		count     int
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "Happy case: patient timeline",
			args: args{
				ctx:       context.Background(),
				patientID: *patient.PatientRecord.ID,
				count:     4,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.PatientTimeline(tt.args.ctx, tt.args.patientID, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient timeline to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient timeline not to be nil for %v", tt.name)
				return
			}
		})
	}

}

func TestClinicalUseCaseImpl_PatientHealthTimeline(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		// TODO: investigate register patient
		// t.Errorf("unable to create patient: %v", err)
		return
	}

	episode, _, err := createTestEpisodeOfCare(
		context.Background(),
		interserviceclient.TestUserPhoneNumber,
		false,
		MFLCode,
	)
	if err != nil {
		t.Errorf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := u.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	oInput, err := getFhirObservationInput(*patient.PatientRecord, encounterID)
	if err != nil {
		t.Errorf("failed to get fhir observation input: %v", err)
	}

	_, err = u.CreateFHIRObservation(ctx, *oInput)
	if err != nil {
		t.Errorf("failed to create fhir observation: %v", err)
	}

	msInput, err := getFhirMedicationStatementInput(*patient.PatientRecord)
	if err != nil {
		t.Errorf("failed to get fhir medication statement input: %v", err)
	}

	_, err = u.CreateFHIRMedicationStatement(ctx, *msInput)
	if err != nil {
		t.Errorf("failed to create fhir medication statement: %v", err)
	}

	type args struct {
		ctx   context.Context
		input domain.HealthTimelineInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.HealthTimeline
		wantErr bool
	}{
		{
			name: "Happy case: patient timeline",
			args: args{
				ctx: context.Background(),
				input: domain.HealthTimelineInput{
					PatientID: *patient.PatientRecord.ID,
					Offset:    0,
					Limit:     20,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.PatientHealthTimeline(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.PatientHealthTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient timeline to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient timeline not to be nil for %v", tt.name)
				return
			}
		})
	}

}

func TestClinicalUseCaseImpl_GetMedicalData(t *testing.T) {
	ctx := context.Background()
	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		// TODO: investigate register patient
		// t.Errorf("unable to create patient: %v", err)
		return
	}

	episode, _, err := createTestEpisodeOfCare(
		context.Background(),
		interserviceclient.TestUserPhoneNumber,
		false,
		MFLCode,
	)
	if err != nil {
		t.Errorf("cant get test encounter id: %v\n", err)
		return
	}

	encounterID, err := u.StartEncounter(ctx, *episode.ID)
	if err != nil {
		t.Errorf("failed to start encounter: %v\n", err)
	}

	oInput, err := getFhirObservationInput(*patient.PatientRecord, encounterID)
	if err != nil {
		t.Errorf("failed to get fhir observation input: %v", err)
	}

	_, err = u.CreateFHIRObservation(ctx, *oInput)
	if err != nil {
		t.Errorf("failed to create fhir observation: %v", err)
	}

	msInput, err := getFhirMedicationStatementInput(*patient.PatientRecord)
	if err != nil {
		t.Errorf("failed to get fhir medication statement input: %v", err)
	}

	_, err = u.CreateFHIRMedicationStatement(ctx, *msInput)
	if err != nil {
		t.Errorf("failed to create fhir medication statement: %v", err)
	}

	type args struct {
		ctx       context.Context
		patientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.MedicalData
		wantErr bool
	}{
		{
			name: "Happy case: patient timeline",
			args: args{
				ctx:       context.Background(),
				patientID: *patient.PatientRecord.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.GetMedicalData(tt.args.ctx, tt.args.patientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.GetMedicalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected patient medical data to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected patient medical data not to be nil for %v", tt.name)
				return
			}
		})
	}

}
