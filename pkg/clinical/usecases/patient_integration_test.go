package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

const testProviderCode = "1234"

func TestClinicalUseCaseImpl_ProblemSummary(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

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
				patientID: ksuid.New().String(),
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_VisitSummary(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	_, err = u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
		return
	}

	normalized, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		fmt.Printf("can't normalize phone number: %v \n", err)
	}
	_, patient, _, err := getTestEncounterID(
		ctx,
		msisdn,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}
	orgID, err := u.GetORCreateOrganization(ctx, testProviderCode)
	if err != nil {
		log.Printf("can't get or create test organization : %v\n", err)
	}

	episode := helpers.ComposeOneHealthEpisodeOfCare(
		*normalized,
		false,
		*orgID,
		testProviderCode,
		*patient.ID,
	)

	ep, err := u.CreateEpisodeOfCare(ctx, episode)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *ep.EpisodeOfCare.ID, err)
		return
	}

	episodePayload, err := u.GetFHIREpisodeOfCare(ctx, *ep.EpisodeOfCare.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *ep.EpisodeOfCare.ID, err)
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_PatientTimelineWithCount(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}
	u := testUsecaseInteractor

	msisdn := interserviceclient.TestUserPhoneNumber

	normalized, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		fmt.Printf("can't normalize phone number: %v \n", err)
	}
	_, patient, _, err := getTestEncounterID(
		ctx,
		*normalized,
		false,
		testProviderCode,
	)
	if err != nil {
		log.Printf("cant get test encounter id: %v\n", err)
		return
	}
	orgID, err := u.GetORCreateOrganization(ctx, testProviderCode)
	if err != nil {
		log.Printf("can't get or create test organization : %v\n", err)
	}

	episode := helpers.ComposeOneHealthEpisodeOfCare(
		*normalized,
		false,
		*orgID,
		testProviderCode,
		*patient.ID,
	)

	ep, err := u.CreateEpisodeOfCare(ctx, episode)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *ep.EpisodeOfCare.ID, err)
		return
	}

	episodePayload, err := u.GetFHIREpisodeOfCare(ctx, *ep.EpisodeOfCare.ID)
	if err != nil {
		t.Errorf("unable to get episode with ID %s: %v", *ep.EpisodeOfCare.ID, err)
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
				episodeID: *ep.EpisodeOfCare.ID,
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_PatientSearch(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	searchPatient := "Test user"

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	_, err = u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_ContactsToContactPointInput(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	_, err = u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_CreatePatient(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	ID := ksuid.New().String()

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	_, err = u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_FindPatientByID(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
		return
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
				id:  *patient.PatientRecord.ID,
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_UpdatePatient(t *testing.T) {
	TestClinicalUseCaseImpl_DeleteFHIRPatientByPhone(t)

	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_AddNextOfKin(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_AddNHIF(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	testInput := ksuid.New().String()

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_CreateUpdatePatientExtraInformation(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	maritalStatus := ksuid.New().String()

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_AllergySummary(t *testing.T) {

	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_DeleteFHIRPatientByPhone(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	_, err = u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
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
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
			},
			wantErr: false,
		},
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

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}

func TestClinicalUseCaseImpl_StartEpisodeByBreakGlass(t *testing.T) {
	onboardingClient := onboardingISCClient(t)
	ctx, token, err := interserviceclient.GetPhoneNumberAuthenticatedContextAndToken(t, onboardingClient)
	if err != nil {
		t.Errorf("cant get phone number authenticated context token: %v", err)
		return
	}

	_, err = firebasetools.GetAuthenticatedContextFromUID(ctx, token.UID)
	if err != nil {
		t.Errorf("cant get authenticated context from UID: %v", err)
		return
	}

	u := testUsecaseInteractor

	phone := interserviceclient.TestUserPhoneNumber
	otp, err := generateTestOTP(t, phone)
	if err != nil {
		log.Errorf("unable to get verified phone number and OTP")
		return
	}

	patientInput, err := simplePatientRegistration()
	if err != nil {
		t.Errorf("an error occurred: %v\n", err)
		return
	}

	patient, err := u.RegisterPatient(ctx, *patientInput)
	if err != nil {
		fmt.Printf("unable to create patient: %v", err)
		return
	}

	type args struct {
		ctx   context.Context
		input domain.BreakGlassEpisodeCreationInput
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
				input: domain.BreakGlassEpisodeCreationInput{
					PatientID:       *patient.PatientRecord.ID,
					ProviderCode:    ksuid.New().String(),
					PractitionerUID: ksuid.New().String(),
					ProviderPhone:   phone,
					Otp:             otp,
					FullAccess:      true,
					PatientPhone:    phone,
				},
			},
			wantErr: false,
		},

		{
			name: "Sad case: empty patient ID",
			args: args{
				ctx: ctx,
				input: domain.BreakGlassEpisodeCreationInput{
					PatientID:       "",
					ProviderCode:    ksuid.New().String(),
					PractitionerUID: ksuid.New().String(),
					ProviderPhone:   ksuid.New().String(),
					Otp:             otp,
					FullAccess:      true,
					PatientPhone:    ksuid.New().String(),
				},
			},
			wantErr: true,
		},

		{
			name: "Sad case: no input",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := u.StartEpisodeByBreakGlass(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicalUseCaseImpl.StartEpisodeByBreakGlass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

	_, err = u.DeleteFHIRPatientByPhone(ctx, interserviceclient.TestUserPhoneNumber)
	if err != nil {
		fmt.Printf("unable to delete FHIR patient by phone: %v", err)
		return
	}
}
