package graph_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

const (
	testPhotoBase64 = "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"
)

func getTestSimplePatientRegistration() (*clinical.SimplePatientRegistrationInput, string, error) {
	otherNames := gofakeit.Name()
	dob := gofakeit.Date()
	msisdn := gofakeit.Phone()
	birthDate, err := scalarutils.NewDate(dob.Day(), int(dob.Month()), dob.Year())
	if err != nil {
		return nil, "", fmt.Errorf("can't create valid birth date: %w", err)
	}
	return &clinical.SimplePatientRegistrationInput{
		ID: ksuid.New().String(),
		Names: []*clinical.NameInput{
			{
				FirstName:  gofakeit.FirstName(),
				LastName:   gofakeit.LastName(),
				OtherNames: &otherNames,
			},
		},
		IdentificationDocuments: []*clinical.IdentificationDocument{
			{
				DocumentType:   clinical.IDDocumentTypeNationalID,
				DocumentNumber: strconv.Itoa(gofakeit.Number(11111111, 111111111)),
			},
		},
		BirthDate: *birthDate,
		PhoneNumbers: []*clinical.PhoneNumberInput{
			{
				Msisdn:             msisdn,
				VerificationCode:   ksuid.New().String(),
				IsUssd:             true, // this will turn off OTP verification
				CommunicationOptIn: true,
			},
		},
		Photos: []*clinical.PhotoInput{
			{
				PhotoContentType: enumutils.ContentTypeJpg,
				PhotoBase64data:  testPhotoBase64,
				PhotoFilename:    fmt.Sprintf("%s.jpg", gofakeit.Name()),
			},
		},
		Emails: []*clinical.EmailInput{
			{
				Email:              "test@bewell.co.ke",
				CommunicationOptIn: true,
			},
		},
		PhysicalAddresses: []*clinical.PhysicalAddress{
			{
				PhysicalAddress: gofakeit.Address().Address,
				MapsCode:        gofakeit.Address().Zip,
			},
		},
		PostalAddresses: []*clinical.PostalAddress{
			{
				PostalAddress: gofakeit.Address().Address,
				PostalCode:    gofakeit.Address().City,
			},
		},
		Gender:        enumutils.GenderFemale.String(),
		Active:        true,
		MaritalStatus: clinical.MaritalStatusS,
		Languages: []enumutils.Language{
			enumutils.LanguageEn,
			enumutils.LanguageSw,
		},
	}, msisdn, nil
}

func getTestVerifiedPhoneandOTP() (string, string, error) {
	config, err := base.LoadDepsFromYAML()
	if err != nil {
		return "", "", fmt.Errorf("unable to load dependencies from YAML: %s", err)
	}

	engagementClient, err := base.SetupISCclient(*config, clinical.EngagementService)
	if err != nil {
		return "", "", fmt.Errorf("unable to set up engagement ISC client: %v", err)
	}

	validPhone := base.TestUserPhoneNumber
	validOTP, err := clinical.RequestOTP(validPhone, engagementClient)
	if err != nil {
		return "", "", fmt.Errorf("unable to generate OTP: %v", err)
	}

	return validPhone, validOTP, err
}

func getTestPatient(ctx context.Context) (*clinical.FHIRPatient, string, error) {
	srv := clinical.NewService()
	simplePatientRegInput, msisdn, err := getTestSimplePatientRegistration()
	if err != nil {
		return nil, "", fmt.Errorf("can't genereate simple patient reg input: %v", err)
	}

	patientPayload, err := srv.RegisterPatient(ctx, *simplePatientRegInput)
	if err != nil {
		return nil, "", fmt.Errorf("can't register patient: %v", err)
	}
	return patientPayload.PatientRecord, msisdn, nil
}

func getTestEpisodeOfCare(
	ctx context.Context,
	msisdn string,
	fullAccess bool,
	providerCode string,
) (*clinical.FHIREpisodeOfCare, *clinical.FHIRPatient, error) {
	srv := clinical.NewService()
	normalized, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		return nil, nil, fmt.Errorf("can't normalize phone number: %w", err)
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create test patient: %w", err)
	}

	orgID, err := srv.GetORCreateOrganization(ctx, providerCode)
	if err != nil {
		return nil, nil, fmt.Errorf("can't get or create test organization : %v", err)
	}

	ep := clinical.ComposeOneHealthEpisodeOfCare(
		*normalized,
		fullAccess,
		*orgID,
		providerCode,
		*patient.ID,
	)
	epPayload, err := srv.CreateEpisodeOfCare(ctx, ep)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create episode of care: %w", err)
	}
	return epPayload.EpisodeOfCare, patient, nil
}

func getTestEncounterID(
	ctx context.Context,
	msisdn string,
	fullAccess bool,
	providerCode string,
) (*clinical.FHIREpisodeOfCare, *clinical.FHIRPatient, string, error) {
	srv := clinical.NewService()
	episode, patient, err := getTestEpisodeOfCare(
		ctx, msisdn, fullAccess, providerCode)
	if err != nil {
		return nil, nil, "", fmt.Errorf("can't create episode of care: %w", err)
	}

	encounterID, err := srv.StartEncounter(ctx, *episode.ID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("unable to start encounter: %w", err)
	}

	return episode, patient, encounterID, nil
}

func getTestFHIRMedicationRequestID(ctx context.Context, encounterID string) (string, error) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		return "", fmt.Errorf("error in getting GraphQL headers: %w", err)
	}
	patient, _, err := getTestPatient(ctx)
	if err != nil {
		return "", fmt.Errorf("error creating test patient: %w", err)
	}

	patientName := patient.Names()
	requester := gofakeit.Name()
	dateRecorded := time.Now().Format(dateFormat)

	query := map[string]interface{}{
		"query": `mutation CreateMedicationRequest($input: FHIRMedicationRequestInput!) {
			createFHIRMedicationRequest(input: $input) {
			  resource {
				ID
			  }
			}
		  }`,
		"variables": map[string]interface{}{
			"input": map[string]interface{}{
				"Status":   "active",
				"Intent":   "proposal",
				"Priority": "routine",
				"Subject": map[string]interface{}{
					"Reference": fmt.Sprintf("Patient/%s", *patient.ID),
					"Type":      "Patient",
					"Display":   patientName,
				},
				"Encounter": map[string]interface{}{
					"Reference": fmt.Sprintf("Encounter/%s", encounterID),
					"Type":      "Encounter",
					"Display":   fmt.Sprintf("Encounter/%s", encounterID),
				},
				"SupportingInformation": []map[string]interface{}{
					{
						"ID":        "113488",
						"Reference": fmt.Sprintf("Encounter/%s", encounterID),
						"Display":   "Pulmonary Tuberculosis",
					},
				},
				"Requester": map[string]interface{}{
					"Display": requester,
				},
				"Note": []map[string]interface{}{
					{
						"AuthorString": requester,
						"Text":         gofakeit.HipsterSentence(10),
					},
				},
				"MedicationCodeableConcept": map[string]interface{}{
					"Text": "Panadol Extra",
					"Coding": []map[string]interface{}{
						{
							"System":       "OCL:/orgs/CIEL/sources/CIEL/",
							"Code":         "999999",
							"Display":      "Panadol Extra",
							"UserSelected": true,
						},
					},
				},
				"DosageInstruction": []map[string]interface{}{
					{
						"Text":               "500 mg 5/7 B.D.",
						"PatientInstruction": "Take two tablets after meals, three times a day",
					},
				},
				"AuthoredOn": dateRecorded,
			},
		},
	}

	body, err := mapToJSONReader(query)
	if err != nil {
		return "", fmt.Errorf("unable to get GQL JSON io Reader: %s", err)
	}

	r, err := http.NewRequest(
		http.MethodPost,
		graphQLURL,
		body,
	)
	if err != nil {
		return "", fmt.Errorf("unable to compose request: %s", err)
	}

	if r == nil {
		return "", fmt.Errorf("nil request")
	}

	for k, v := range headers {
		r.Header.Add(k, v)
	}
	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}
	resp, err := client.Do(r)
	if err != nil {
		return "", fmt.Errorf("request error: %s", err)
	}

	dataResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't read request body: %s", err)
	}

	if dataResponse == nil {
		return "", fmt.Errorf("nil response data")
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(dataResponse, &data)
	if err != nil {
		return "", fmt.Errorf("bad data returned")
	}

	for key := range data {
		nestedMap, ok := data[key].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
		}

		for nestedKey := range nestedMap {
			if nestedKey == "createFHIRMedicationRequest" {
				output, ok := nestedMap[nestedKey].(map[string]interface{})
				if !ok {
					return "", fmt.Errorf("can't cast output to map[string]interface{}")
				}

				resource, ok := output["resource"].(map[string]interface{})
				if !ok {
					return "", fmt.Errorf("can't cast resource to map[string]interface{}")
				}

				id, prs := resource["ID"]
				if !prs {
					return "", fmt.Errorf("ID not present in medication request resource")
				}
				if id == "" {
					return "", fmt.Errorf("blank medication request ID")
				}

				idString := id.(string)
				return idString, nil
			}
		}
	}
	return "", err
}

func createFHIRTestObservation(ctx context.Context, encounterID string) (
	*clinical.FHIRObservation,
	*clinical.FHIRPatient,
	*clinical.ObservationStatusEnum,
	error,
) {
	instantRecorded := scalarutils.Instant(time.Now().Format(instantFormat))
	patient, _, err := getTestPatient(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't create test patient: %w", err)
	}
	srv := clinical.NewService()
	status := clinical.ObservationStatusEnumPreliminary
	categorySystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/observation-category")
	loincSystem := scalarutils.URI("http://loinc.org")
	notSelected := false
	selected := true
	refrangeText := "0kg to 300kg"
	refrangeSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/referencerange-meaning")
	interpretationSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/v3-ObservationInterpretation")
	patientType := scalarutils.URI("Patient")
	encounterType := scalarutils.URI("Encounter")
	encounterRef := fmt.Sprintf("Encounter/%s", encounterID)
	patientRef := fmt.Sprintf("Patient/%s", *patient.ID)

	inp := clinical.FHIRObservationInput{
		Status: &status,
		Category: []*clinical.FHIRCodeableConceptInput{
			{
				Text: "Vital Signs",
				Coding: []*clinical.FHIRCodingInput{
					{
						Code:         "vital-signs",
						System:       &categorySystem,
						Display:      "Vital Signs",
						UserSelected: &notSelected,
					},
				},
			},
		},
		Code: clinical.FHIRCodeableConceptInput{
			Text: "Body weight",
			Coding: []*clinical.FHIRCodingInput{
				{
					Code:         "29463-7",
					System:       &loincSystem,
					Display:      "Body Weight",
					UserSelected: &selected,
				},
			},
		},
		ValueQuantity: &clinical.FHIRQuantityInput{
			Value:  72,
			Unit:   "kg",
			System: scalarutils.URI("http://unitsofmeasure.org"),
			Code:   scalarutils.Code("kg"),
		},
		ReferenceRange: []*clinical.FHIRObservationReferencerangeInput{
			{
				Low: &clinical.FHIRQuantityInput{
					Value:  0,
					Unit:   "kg",
					System: scalarutils.URI("http://unitsofmeasure.org"),
					Code:   "kg",
				},
				High: &clinical.FHIRQuantityInput{
					Value:  300,
					Unit:   "kg",
					System: scalarutils.URI("http://unitsofmeasure.org"),
					Code:   "kg",
				},
				Text: &refrangeText,
				Type: &clinical.FHIRCodeableConceptInput{
					Text: "Normal Range",
					Coding: []*clinical.FHIRCodingInput{
						{
							Code:         "normal",
							System:       &refrangeSystem,
							Display:      "Normal Range",
							UserSelected: &notSelected,
						},
					},
				},
			},
		},
		Issued:           &instantRecorded,
		EffectiveInstant: &instantRecorded,
		Encounter: &clinical.FHIRReferenceInput{
			Reference: &encounterRef,
			Type:      &encounterType,
			Display:   encounterRef,
		},
		Subject: &clinical.FHIRReferenceInput{
			Reference: &patientRef,
			Type:      &patientType,
			Display:   patientRef,
		},
		Interpretation: []*clinical.FHIRCodeableConceptInput{
			{
				Text: "Normal",
				Coding: []*clinical.FHIRCodingInput{
					{
						Code:         "N",
						System:       &interpretationSystem,
						Display:      "Normal",
						UserSelected: &notSelected,
					},
				},
			},
		},
	}
	obsPl, err := srv.CreateFHIRObservation(ctx, inp)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't create FHIR observation: %w", err)
	}
	return obsPl.Resource, patient, &status, nil
}

func getTestSimpleServiceRequest(ctx context.Context, encounterID string) (
	*clinical.FHIRServiceRequestInput,
	string,
	error,
) {
	patient, _, err := getTestPatient(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("can't create a test patient: %w", err)
	}
	patientName := patient.Names()
	requester := gofakeit.Name()
	patientRef := fmt.Sprintf("Patient/%s", *patient.ID)
	patientType := scalarutils.URI("Patient")
	encounterRef := fmt.Sprintf("Encounter/%s", encounterID)
	encounterType := scalarutils.URI("Encounter")
	active := scalarutils.Code(clinical.EpisodeOfCareStatusEnumActive)
	system := scalarutils.URI("OCL:/orgs/CIEL/sources/CIEL/")
	userSelected := true
	intent := scalarutils.Code("proposal")
	priority := scalarutils.Code("routine")

	return &clinical.FHIRServiceRequestInput{
		Status:   &active,
		Intent:   &intent,
		Priority: &priority,
		Subject: &clinical.FHIRReferenceInput{
			Reference: &patientRef,
			Type:      &patientType,
			Display:   patientName,
		},
		Encounter: &clinical.FHIRReferenceInput{
			Reference: &encounterRef,
			Type:      &encounterType,
			Display:   encounterRef,
		},
		SupportingInfo: []*clinical.FHIRReferenceInput{
			{
				Reference: &encounterRef,
				Display:   "Pulmonary Tuberculosis",
			},
		},
		Category: []*clinical.FHIRCodeableConceptInput{
			{
				Text: "Laboratory procedure",
				Coding: []*clinical.FHIRCodingInput{
					{
						System:       &system,
						Code:         "108252007",
						Display:      "Laboratory procedure",
						UserSelected: &userSelected,
					},
				},
			},
		},
		Requester: &clinical.FHIRReferenceInput{
			Display: requester,
		},
		Code: &clinical.FHIRCodeableConceptInput{
			Text: "Hospital re-admission",
			Coding: []*clinical.FHIRCodingInput{
				{
					System:       &system,
					Code:         "417005",
					Display:      "Hospital re-admission",
					UserSelected: &userSelected,
				},
			},
		},
	}, *patient.ID, nil
}

func getTestServiceRequest(ctx context.Context, encounterID string) (
	*clinical.FHIRServiceRequest,
	string,
	error,
) {
	srv := clinical.NewService()
	simpleServiceRequestInput, patientID, err := getTestSimpleServiceRequest(ctx, encounterID)
	if err != nil {
		return nil, "", fmt.Errorf("can't genereate simple service request input: %v", err)
	}

	serviceRequestPayload, err := srv.CreateFHIRServiceRequest(ctx, *simpleServiceRequestInput)
	if err != nil {
		return nil, "", fmt.Errorf("can't create service request: %v", err)
	}

	return serviceRequestPayload.Resource, patientID, nil
}

func createTestFHIRComposition(
	ctx context.Context,
	encounterID string,
) (*clinical.FHIRComposition,
	*clinical.FHIRPatient,
	error,
) {
	srv := clinical.NewService()
	status := clinical.CompositionStatusEnumPreliminary
	now := time.Now()
	title := gofakeit.HipsterSentence(10)
	author := gofakeit.Name()
	system := scalarutils.URI("http://loinc.org")
	historyTitle := "Patient History"
	notSelected := false
	generatedStatus := clinical.NarrativeStatusEnumGenerated

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create patient: %w", err)
	}
	patientRef := fmt.Sprintf("Patient/%s", *patient.ID)
	patientType := scalarutils.URI("Patient")
	encounterRef := fmt.Sprintf("Encounter/%s", encounterID)
	encounterType := scalarutils.URI("Encounter")

	recorded, err := scalarutils.NewDate(now.Day(), int(now.Month()), now.Year())
	if err != nil {
		return nil, nil, fmt.Errorf("can't initialize recorded date: %w", err)
	}
	inp := clinical.FHIRCompositionInput{
		Status: &status,
		Date:   recorded,
		Title:  &title,
		Author: []*clinical.FHIRReferenceInput{
			{
				Display: author,
			},
		},
		Type: &clinical.FHIRCodeableConceptInput{
			Text: "Consult Note",
			Coding: []*clinical.FHIRCodingInput{
				{
					System:       &system,
					Code:         scalarutils.Code("11488-4"),
					Display:      "Consult Note",
					UserSelected: &notSelected,
				},
			},
		},
		Category: []*clinical.FHIRCodeableConceptInput{
			{
				Text: "Consult Note",
				Coding: []*clinical.FHIRCodingInput{
					{
						System:       &system,
						Code:         scalarutils.Code("11488-4"),
						Display:      "Consult Note",
						UserSelected: &notSelected,
					},
				},
			},
		},
		Section: []*clinical.FHIRCompositionSectionInput{
			{
				Title: &historyTitle,
				Text: &clinical.FHIRNarrativeInput{
					Status: &generatedStatus,
					Div:    scalarutils.XHTML(gofakeit.HipsterSentence(10)),
				},
			},
		},
		Encounter: &clinical.FHIRReferenceInput{
			Reference: &encounterRef,
			Type:      &encounterType,
			Display:   encounterRef,
		},
		Subject: &clinical.FHIRReferenceInput{
			Reference: &patientRef,
			Type:      &patientType,
			Display:   patientRef,
		},
	}
	compPl, err := srv.CreateFHIRComposition(ctx, inp)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create composition payload: %w", err)
	}
	return compPl.Resource, patient, nil
}

func createTestConditionInput(
	encounterID string,
	patientID string,
) (*clinical.FHIRConditionInput, error) {
	system := scalarutils.URI("OCL:/orgs/CIEL/sources/CIEL/")
	userSelected := true
	falseUserSelected := false
	clinicalSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-clinical")
	verificationStatusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-ver-status")
	categorySystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/condition-category")
	name := gofakeit.Name()
	text := scalarutils.Markdown(gofakeit.HipsterSentence(20))
	encounterType := scalarutils.URI("Encounter")
	encounterRef := fmt.Sprintf("Encounter/%s", encounterID)
	subjectType := scalarutils.URI("Patient")
	patRef := fmt.Sprintf("Patient/%s", patientID)
	dateRecorded := scalarutils.Date{
		Year:  gofakeit.Year(),
		Month: 12,
		Day:   gofakeit.Day(),
	}

	return &clinical.FHIRConditionInput{
		Code: &clinical.FHIRCodeableConceptInput{
			Coding: []*clinical.FHIRCodingInput{
				{
					System:       &system,
					Code:         scalarutils.Code("113488"),
					Display:      "Pulmonary Tuberculosis",
					UserSelected: &userSelected,
				},
			},
			Text: "Pulmonary Tuberculosis",
		},
		ClinicalStatus: &clinical.FHIRCodeableConceptInput{
			Coding: []*clinical.FHIRCodingInput{
				{
					System:       &clinicalSystem,
					Code:         scalarutils.Code("active"),
					Display:      "Active",
					UserSelected: &falseUserSelected,
				},
			},
			Text: "Active",
		},
		VerificationStatus: &clinical.FHIRCodeableConceptInput{
			Coding: []*clinical.FHIRCodingInput{
				{
					System:       &verificationStatusSystem,
					Code:         scalarutils.Code("confirmed"),
					Display:      "Confirmed",
					UserSelected: &falseUserSelected,
				},
			},
			Text: "Active",
		},
		RecordedDate: &dateRecorded,
		Category: []*clinical.FHIRCodeableConceptInput{
			{
				Coding: []*clinical.FHIRCodingInput{
					{
						System:       &categorySystem,
						Code:         scalarutils.Code("problem-list-item"),
						Display:      "problem-list-item",
						UserSelected: &falseUserSelected,
					},
				},
				Text: "problem-list-item",
			},
		},
		Subject: &clinical.FHIRReferenceInput{
			Reference: &patRef,
			Type:      &subjectType,
			Display:   patRef,
		},
		Encounter: &clinical.FHIRReferenceInput{
			Reference: &encounterRef,
			Type:      &encounterType,
			Display:   "Encounter",
		},
		Note: []*clinical.FHIRAnnotationInput{
			{
				AuthorString: &name,
				Text:         &text,
			},
		},
		Recorder: &clinical.FHIRReferenceInput{
			Display: gofakeit.Name(),
		},
		Asserter: &clinical.FHIRReferenceInput{
			Display: gofakeit.Name(),
		},
	}, nil
}

func createTestCondition(
	ctx context.Context,
	encounterID,
	patientID string,
) (*clinical.FHIRCondition, error) {
	srv := clinical.NewService()
	conditionInput, err := createTestConditionInput(encounterID, patientID)
	if err != nil {
		return nil, err
	}
	condition, err := srv.CreateFHIRCondition(ctx, *conditionInput)
	if err != nil {
		return nil, fmt.Errorf("can't create a test condition")
	}
	return condition.Resource, nil
}

func patientVisitSummary(ctx context.Context) (string, string, error) {
	episode, patient, encounterID, err := getTestEncounterID(
		ctx, base.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		return "", "", fmt.Errorf("error creating test encounter ID: %w", err)
	}
	patientID := *patient.ID

	_, _, _, err = createFHIRTestObservation(ctx, encounterID)
	if err != nil {
		return "", "", fmt.Errorf("can't create FHIR test observation: %w", err)
	}

	_, _, err = createTestFHIRComposition(ctx, encounterID)
	if err != nil {
		return "", "", fmt.Errorf("can't create test composition: %w", err)
	}

	_, _, err = getTestServiceRequest(ctx, encounterID)
	if err != nil {
		return "", "", fmt.Errorf("error creating test service request: %w", err)
	}

	_, err = createTestCondition(ctx, encounterID, patientID)
	if err != nil {
		return "", "", fmt.Errorf("error creating test condition: %w", err)
	}

	_, err = getTestFHIRMedicationRequestID(ctx, encounterID)
	if err != nil {
		return "", "", fmt.Errorf("error creating test medication request: %w", err)
	}

	_, err = createTestAllergy(ctx, patient, encounterID)
	if err != nil {
		return "", "", fmt.Errorf("error creating test allergy intolerance: %w", err)
	}

	return encounterID, *episode.ID, nil
}

func createTestAllergy(
	ctx context.Context,
	patient *clinical.FHIRPatient,
	encounterID string,
) (string, error) {
	srv := clinical.NewService()
	patientName := patient.Names()
	now := time.Now()
	dateRecorded, err := scalarutils.NewDate(now.Day(), int(now.Month()), now.Year())
	if err != nil {
		return "", fmt.Errorf("can't initialize date recorded")
	}
	recordingDoctor := gofakeit.Name()
	substanceID := "1234"
	substanceDisplayName := gofakeit.Name()
	allergyType := clinical.AllergyIntoleranceTypeEnumAllergy
	allergySystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/allergyintolerance-verification")
	clinicalStatusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical")
	verificationSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/allergyintolerance-verification")
	notSelected := false
	selected := true
	encounterReference := fmt.Sprintf("Encounter/%s", encounterID)
	encounterType := scalarutils.URI("Encounter")
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	patientType := scalarutils.URI("Patient")
	annotationText := scalarutils.Markdown(gofakeit.HipsterSentence(10))
	reactionDescription := "some reaction"
	reactionSeverity := clinical.AllergyIntoleranceReactionSeverityEnumMild
	oclSystem := scalarutils.URI("OCL:/orgs/CIEL/sources/CIEL/")

	inp := clinical.FHIRAllergyIntoleranceInput{
		Type:         &allergyType,
		Criticality:  clinical.AllergyIntoleranceCriticalityEnumHigh,
		RecordedDate: dateRecorded,
		Code: clinical.FHIRCodeableConceptInput{
			Text: "Panadol Extra",
			Coding: []*clinical.FHIRCodingInput{
				{
					System:       &allergySystem,
					Code:         scalarutils.Code(substanceID),
					Display:      substanceDisplayName,
					UserSelected: &notSelected,
				},
			},
		},
		ClinicalStatus: clinical.FHIRCodeableConceptInput{
			Text: "Panadol Extra",
			Coding: []*clinical.FHIRCodingInput{
				{
					System:       &clinicalStatusSystem,
					Code:         scalarutils.Code("active"),
					Display:      "Active",
					UserSelected: &notSelected,
				},
			},
		},
		VerificationStatus: clinical.FHIRCodeableConceptInput{
			Text: "Panadol Extra",
			Coding: []*clinical.FHIRCodingInput{
				{
					System:       &verificationSystem,
					Code:         "confirmed",
					Display:      "Confirmed",
					UserSelected: &notSelected,
				},
			},
		},
		Recorder: &clinical.FHIRReferenceInput{
			Display: recordingDoctor,
		},
		Asserter: &clinical.FHIRReferenceInput{
			Display: recordingDoctor,
		},
		Encounter: &clinical.FHIRReferenceInput{
			Reference: &encounterReference,
			Type:      &encounterType,
			Display:   fmt.Sprintf("Encounter/%s", encounterID),
		},
		Patient: &clinical.FHIRReferenceInput{
			Reference: &patientReference,
			Type:      &patientType,
			Display:   patientName,
		},
		Note: []*clinical.FHIRAnnotationInput{
			{
				AuthorString: &recordingDoctor,
				Text:         &annotationText,
			},
		},
		Reaction: []*clinical.FHIRAllergyintoleranceReactionInput{
			{
				Description: &reactionDescription,
				Severity:    &reactionSeverity,
				Substance: &clinical.FHIRCodeableConceptInput{
					Text: "Panadol Extra",
					Coding: []*clinical.FHIRCodingInput{
						{
							System:       &oclSystem,
							Code:         scalarutils.Code(substanceID),
							Display:      substanceDisplayName,
							UserSelected: &selected,
						},
					},
				},
				Manifestation: []*clinical.FHIRCodeableConceptInput{
					{
						Text: "Rashes",
						Coding: []*clinical.FHIRCodingInput{
							{
								System:       &oclSystem,
								Code:         scalarutils.Code("a code"),
								Display:      "Rashes",
								UserSelected: &selected,
							},
						},
					},
				},
			},
		},
	}
	allergyPl, err := srv.CreateFHIRAllergyIntolerance(ctx, inp)
	if err != nil {
		return "", fmt.Errorf("can't create allergy intolerance")
	}
	return *allergyPl.Resource.ID, nil
}
