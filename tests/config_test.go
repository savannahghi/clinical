package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/imroc/req"
	"github.com/savannahghi/authutils"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/testutils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/presentation"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/interactor"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
)

const ( // Repo the env to identify which repo to use
	Repo = "REPOSITORY"
	//FirebaseRepository is the value of the env when using firebase
	FirebaseRepository = "firebase"
	instantFormat      = "2006-01-02T15:04:05.999-07:00"
)

/// these are set up once in TestMain and used by all the acceptance tests in
// this package
var (
	srv              *http.Server
	baseURL          string
	serverErr        error
	testInteractor   interactor.Usecases
	testProviderCode = "1234"
)

var (
	authServerEndpoint = serverutils.MustGetEnvVar("AUTHSERVER_ENDPOINT")
	clientID           = serverutils.MustGetEnvVar("CLIENT_ID")
	clientSecret       = serverutils.MustGetEnvVar("CLIENT_SECRET")
	username           = serverutils.MustGetEnvVar("AUTH_USERNAME")
	password           = serverutils.MustGetEnvVar("AUTH_PASSWORD")
	grantType          = serverutils.MustGetEnvVar("GRANT_TYPE")
)

var oauthPayload *authutils.OAUTHResponse
var headers map[string]string

func TestMain(m *testing.M) {
	ctx := context.Background()

	// setup
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("CLOUD_HEALTH_PUBSUB_TOPIC", "healthcloud-sghi-staging")
	os.Setenv("CLOUD_HEALTH_DATASET_ID", "healthcloud-sghi-staging")
	os.Setenv("CLOUD_HEALTH_FHIRSTORE_ID", "healthcloud-sghi-fhir-staging")
	os.Setenv("REPOSITORY", "firebase")

	authServerConfig := authutils.Config{
		AuthServerEndpoint: authServerEndpoint,
		ClientID:           clientID,
		ClientSecret:       clientSecret,
		GrantType:          grantType,
		Username:           username,
		Password:           password,
	}
	authClient, err := authutils.NewClient(authServerConfig)
	if err != nil {
		log.Printf("an error occurred: %v", err)
	}

	oauthPayload, err = authClient.Authenticate()
	if err != nil {
		log.Printf("unable to authenticate with slade 360 auth server")
		return
	}

	headers, err = GetGraphQLHeaders(ctx)
	if err != nil {
		log.Printf("error adding the graphql headers")
		return
	}

	srv, baseURL, serverErr = serverutils.StartTestServer(
		ctx,
		presentation.PrepareServer,
		presentation.ClinicalAllowedOrigins,
	) // set the globals
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
		return
	}

	i, err := testutils.InitializeTestService(ctx)
	if err != nil {
		log.Printf("unable to initialize test service: %v", err)
		return
	}

	testInteractor = i

	// run the tests
	log.Printf("about to run tests")
	code := m.Run()
	log.Printf("finished running tests")

	// cleanup here
	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()
	os.Exit(code)
}

// GetGraphQLHeaders gets relevant GraphQLHeaders
func GetGraphQLHeaders(ctx context.Context) (map[string]string, error) {
	accessToken := fmt.Sprintf("Bearer %s", oauthPayload.AccessToken)

	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": accessToken,
	}, nil
}

func generateTestOTP(t *testing.T, msisdn string) (string, error) {
	// TODO:(engagement) Replace engagement
	return "", nil
}

func createFHIRTestObservation(ctx context.Context, encounterID string) (
	*domain.FHIRObservation,
	*domain.FHIRPatient,
	*domain.ObservationStatusEnum,
	error,
) {
	instantRecorded := scalarutils.Instant(time.Now().Format(instantFormat))
	patient, _, err := getTestPatient(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't create test patient: %w", err)
	}
	srv := testInteractor
	status := domain.ObservationStatusEnumPreliminary
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

	inp := domain.FHIRObservationInput{
		Status: &status,
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Text: "Vital Signs",
				Coding: []*domain.FHIRCodingInput{
					{
						Code:         "vital-signs",
						System:       &categorySystem,
						Display:      "Vital Signs",
						UserSelected: &notSelected,
					},
				},
			},
		},
		Code: domain.FHIRCodeableConceptInput{
			Text: "Body weight",
			Coding: []*domain.FHIRCodingInput{
				{
					Code:         "29463-7",
					System:       &loincSystem,
					Display:      "Body Weight",
					UserSelected: &selected,
				},
			},
		},
		ValueQuantity: &domain.FHIRQuantityInput{
			Value:  72,
			Unit:   "kg",
			System: scalarutils.URI("http://unitsofmeasure.org"),
			Code:   scalarutils.Code("kg"),
		},
		ReferenceRange: []*domain.FHIRObservationReferencerangeInput{
			{
				Low: &domain.FHIRQuantityInput{
					Value:  0,
					Unit:   "kg",
					System: scalarutils.URI("http://unitsofmeasure.org"),
					Code:   "kg",
				},
				High: &domain.FHIRQuantityInput{
					Value:  300,
					Unit:   "kg",
					System: scalarutils.URI("http://unitsofmeasure.org"),
					Code:   "kg",
				},
				Text: &refrangeText,
				Type: &domain.FHIRCodeableConceptInput{
					Text: "Normal Range",
					Coding: []*domain.FHIRCodingInput{
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
		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterRef,
			Type:      &encounterType,
			Display:   encounterRef,
		},
		Subject: &domain.FHIRReferenceInput{
			Reference: &patientRef,
			Type:      &patientType,
			Display:   patientRef,
		},
		Interpretation: []*domain.FHIRCodeableConceptInput{
			{
				Text: "Normal",
				Coding: []*domain.FHIRCodingInput{
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
	*domain.FHIRServiceRequestInput,
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
	active := scalarutils.Code(domain.EpisodeOfCareStatusEnumActive)
	system := scalarutils.URI("OCL:/orgs/CIEL/sources/CIEL/")
	userSelected := true
	intent := scalarutils.Code("proposal")
	priority := scalarutils.Code("routine")

	return &domain.FHIRServiceRequestInput{
		Status:   &active,
		Intent:   &intent,
		Priority: &priority,
		Subject: &domain.FHIRReferenceInput{
			Reference: &patientRef,
			Type:      &patientType,
			Display:   patientName,
		},
		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterRef,
			Type:      &encounterType,
			Display:   encounterRef,
		},
		SupportingInfo: []*domain.FHIRReferenceInput{
			{
				Reference: &encounterRef,
				Display:   "Pulmonary Tuberculosis",
			},
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Text: "Laboratory procedure",
				Coding: []*domain.FHIRCodingInput{
					{
						System:       &system,
						Code:         "108252007",
						Display:      "Laboratory procedure",
						UserSelected: &userSelected,
					},
				},
			},
		},
		Requester: &domain.FHIRReferenceInput{
			Display: requester,
		},
		Code: &domain.FHIRCodeableConceptInput{
			Text: "Hospital re-admission",
			Coding: []*domain.FHIRCodingInput{
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
	*domain.FHIRServiceRequest,
	string,
	error,
) {
	srv := testInteractor
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
) (*domain.FHIRComposition,
	*domain.FHIRPatient,
	error,
) {
	srv := testInteractor
	status := domain.CompositionStatusEnumPreliminary
	now := time.Now()
	title := gofakeit.HipsterSentence(10)
	author := gofakeit.Name()
	system := scalarutils.URI("http://loinc.org")
	historyTitle := "Patient History"
	notSelected := false
	generatedStatus := domain.NarrativeStatusEnumGenerated

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
	inp := domain.FHIRCompositionInput{
		Status: &status,
		Date:   recorded,
		Title:  &title,
		Author: []*domain.FHIRReferenceInput{
			{
				Display: author,
			},
		},
		Type: &domain.FHIRCodeableConceptInput{
			Text: "Consult Note",
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &system,
					Code:         scalarutils.Code("11488-4"),
					Display:      "Consult Note",
					UserSelected: &notSelected,
				},
			},
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Text: "Consult Note",
				Coding: []*domain.FHIRCodingInput{
					{
						System:       &system,
						Code:         scalarutils.Code("11488-4"),
						Display:      "Consult Note",
						UserSelected: &notSelected,
					},
				},
			},
		},
		Section: []*domain.FHIRCompositionSectionInput{
			{
				Title: &historyTitle,
				Text: &domain.FHIRNarrativeInput{
					Status: &generatedStatus,
					Div:    scalarutils.XHTML(gofakeit.HipsterSentence(10)),
				},
			},
		},
		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterRef,
			Type:      &encounterType,
			Display:   encounterRef,
		},
		Subject: &domain.FHIRReferenceInput{
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

func getTestEpisodeOfCare(
	ctx context.Context,
	msisdn string,
	fullAccess bool,
	providerCode string,
) (*domain.FHIREpisodeOfCare, *domain.FHIRPatient, error) {
	normalized, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		return nil, nil, fmt.Errorf("can't normalize phone number: %w", err)
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create test patient: %w", err)
	}

	orgID, err := testInteractor.GetORCreateOrganization(ctx, providerCode)
	if err != nil {
		return nil, nil, fmt.Errorf("can't get or create test organization : %v", err)
	}

	ep := helpers.ComposeOneHealthEpisodeOfCare(
		*normalized,
		fullAccess,
		*orgID,
		providerCode,
		*patient.ID,
	)
	epPayload, err := testInteractor.CreateEpisodeOfCare(ctx, ep)
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
) (*domain.FHIREpisodeOfCare, *domain.FHIRPatient, string, error) {
	episode, patient, err := getTestEpisodeOfCare(
		ctx, msisdn, fullAccess, providerCode)
	if err != nil {
		return nil, nil, "", fmt.Errorf("can't create episode of care: %w", err)
	}

	encounterID, err := testInteractor.StartEncounter(ctx, *episode.ID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("unable to start encounter: %w", err)
	}

	return episode, patient, encounterID, nil
}

func getTestFHIRMedicationRequestID(ctx context.Context, encounterID string) (string, error) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
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

func createTestConditionInput(
	encounterID string,
	patientID string,
) (*domain.FHIRConditionInput, error) {
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

	return &domain.FHIRConditionInput{
		Code: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &system,
					Code:         scalarutils.Code("113488"),
					Display:      "Pulmonary Tuberculosis",
					UserSelected: &userSelected,
				},
			},
			Text: "Pulmonary Tuberculosis",
		},
		ClinicalStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &clinicalSystem,
					Code:         scalarutils.Code("active"),
					Display:      "Active",
					UserSelected: &falseUserSelected,
				},
			},
			Text: "Active",
		},
		VerificationStatus: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
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
		Category: []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
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
		Subject: &domain.FHIRReferenceInput{
			Reference: &patRef,
			Type:      &subjectType,
			Display:   patRef,
		},
		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterRef,
			Type:      &encounterType,
			Display:   "Encounter",
		},
		Note: []*domain.FHIRAnnotationInput{
			{
				AuthorString: &name,
				Text:         &text,
			},
		},
		Recorder: &domain.FHIRReferenceInput{
			Display: gofakeit.Name(),
		},
		Asserter: &domain.FHIRReferenceInput{
			Display: gofakeit.Name(),
		},
	}, nil
}

func createTestCondition(
	ctx context.Context,
	encounterID,
	patientID string,
) (*domain.FHIRCondition, error) {
	srv := testInteractor
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
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
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
	patient *domain.FHIRPatient,
	encounterID string,
) (string, error) {
	srv := testInteractor
	patientName := patient.Names()
	now := time.Now()
	dateRecorded, err := scalarutils.NewDate(now.Day(), int(now.Month()), now.Year())
	if err != nil {
		return "", fmt.Errorf("can't initialize date recorded")
	}
	recordingDoctor := gofakeit.Name()
	substanceID := "1234"
	substanceDisplayName := gofakeit.Name()
	allergyType := domain.AllergyIntoleranceTypeEnumAllergy
	allergySystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/allergyintolerance-verification")
	clinicalStatusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/allergyintolerance-domain")
	verificationSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/allergyintolerance-verification")
	notSelected := false
	selected := true
	encounterReference := fmt.Sprintf("Encounter/%s", encounterID)
	encounterType := scalarutils.URI("Encounter")
	patientReference := fmt.Sprintf("Patient/%s", *patient.ID)
	patientType := scalarutils.URI("Patient")
	annotationText := scalarutils.Markdown(gofakeit.HipsterSentence(10))
	reactionDescription := "some reaction"
	reactionSeverity := domain.AllergyIntoleranceReactionSeverityEnumMild
	oclSystem := scalarutils.URI("OCL:/orgs/CIEL/sources/CIEL/")

	inp := domain.FHIRAllergyIntoleranceInput{
		Type:         &allergyType,
		Criticality:  domain.AllergyIntoleranceCriticalityEnumHigh,
		RecordedDate: dateRecorded,
		Code: domain.FHIRCodeableConceptInput{
			Text: "Panadol Extra",
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &allergySystem,
					Code:         scalarutils.Code(substanceID),
					Display:      substanceDisplayName,
					UserSelected: &notSelected,
				},
			},
		},
		ClinicalStatus: domain.FHIRCodeableConceptInput{
			Text: "Panadol Extra",
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &clinicalStatusSystem,
					Code:         scalarutils.Code("active"),
					Display:      "Active",
					UserSelected: &notSelected,
				},
			},
		},
		VerificationStatus: domain.FHIRCodeableConceptInput{
			Text: "Panadol Extra",
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &verificationSystem,
					Code:         "confirmed",
					Display:      "Confirmed",
					UserSelected: &notSelected,
				},
			},
		},
		Recorder: &domain.FHIRReferenceInput{
			Display: recordingDoctor,
		},
		Asserter: &domain.FHIRReferenceInput{
			Display: recordingDoctor,
		},
		Encounter: &domain.FHIRReferenceInput{
			Reference: &encounterReference,
			Type:      &encounterType,
			Display:   fmt.Sprintf("Encounter/%s", encounterID),
		},
		Patient: &domain.FHIRReferenceInput{
			Reference: &patientReference,
			Type:      &patientType,
			Display:   patientName,
		},
		Note: []*domain.FHIRAnnotationInput{
			{
				AuthorString: &recordingDoctor,
				Text:         &annotationText,
			},
		},
		Reaction: []*domain.FHIRAllergyintoleranceReactionInput{
			{
				Description: &reactionDescription,
				Severity:    &reactionSeverity,
				Substance: &domain.FHIRCodeableConceptInput{
					Text: "Panadol Extra",
					Coding: []*domain.FHIRCodingInput{
						{
							System:       &oclSystem,
							Code:         scalarutils.Code(substanceID),
							Display:      substanceDisplayName,
							UserSelected: &selected,
						},
					},
				},
				Manifestation: []*domain.FHIRCodeableConceptInput{
					{
						Text: "Rashes",
						Coding: []*domain.FHIRCodingInput{
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
