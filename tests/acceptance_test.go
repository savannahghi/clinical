package acceptance_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/brianvoe/gofakeit"
	"github.com/imroc/req"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"
)

const (
	testHTTPClientTimeout = 180
	testProviderCode      = "123"
	dateFormat            = "2006-01-02"
	instantFormat         = "2006-01-02T15:04:05.999-07:00"
	testProviderPhone     = "+254721000111"
	testProviderUID       = "0b1fcd62-46df-4cbc-9096-7849859dcd76"
)

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %v", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

// GetGraphQLHeaders gets relevant GraphQLHeaders
func GetGraphQLHeaders(ctx context.Context) (map[string]string, error) {
	authorization, err := GetBearerTokenHeader(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't Generate Bearer Token: %s", err)
	}
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": authorization,
	}, nil
}

// GetBearerTokenHeader gets bearer Token Header
func GetBearerTokenHeader(ctx context.Context) (string, error) {
	TestUserEmail := "test@bewell.co.ke"
	user, err := firebasetools.GetOrCreateFirebaseUser(ctx, TestUserEmail)
	if err != nil {
		return "", fmt.Errorf("can't get or create firebase user: %s", err)
	}

	if user == nil {
		return "", fmt.Errorf("nil firebase user")
	}

	customToken, err := firebasetools.CreateFirebaseCustomToken(ctx, user.UID)
	if err != nil {
		return "", fmt.Errorf("can't create custom token: %s", err)
	}

	if customToken == "" {
		return "", fmt.Errorf("blank custom token: %s", err)
	}

	idTokens, err := firebasetools.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return "", fmt.Errorf("can't authenticate custom token: %s", err)
	}
	if idTokens == nil {
		return "", fmt.Errorf("nil idTokens")
	}

	return fmt.Sprintf("Bearer %s", idTokens.IDToken), nil
}

const (
	testPhotoBase64 = "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"
)

// func getTestSimpleUserRegistration() (*domain.SimplePatientRegistrationInput, string, error) {
// 	otherNames := gofakeit.Name()
// 	dob := gofakeit.Date()
// 	msisdn := interserviceclient.TestUserPhoneNumberWithPin
// 	birthDate, err := scalarutils.NewDate(dob.Day(), int(dob.Month()), dob.Year())
// 	if err != nil {
// 		return nil, "", fmt.Errorf("can't create valid birth date: %v", err)
// 	}
// 	return &domain.SimplePatientRegistrationInput{
// 		ID: ksuid.New().String(),
// 		Names: []*domain.NameInput{
// 			{
// 				FirstName:  gofakeit.FirstName(),
// 				LastName:   gofakeit.LastName(),
// 				OtherNames: &otherNames,
// 			},
// 		},
// 		IdentificationDocuments: []*domain.IdentificationDocument{
// 			{
// 				DocumentType:   domain.IDDocumentTypeNationalID,
// 				DocumentNumber: strconv.Itoa(gofakeit.Number(11111111, 111111111)),
// 			},
// 		},
// 		BirthDate: *birthDate,
// 		PhoneNumbers: []*domain.PhoneNumberInput{
// 			{
// 				Msisdn:             msisdn,
// 				VerificationCode:   ksuid.New().String(),
// 				IsUssd:             true, // this will turn off OTP verification
// 				CommunicationOptIn: true,
// 			},
// 		},
// 		Photos: []*domain.PhotoInput{
// 			{
// 				PhotoContentType: enumutils.ContentTypeJpg,
// 				PhotoBase64data:  testPhotoBase64,
// 				PhotoFilename:    fmt.Sprintf("%s.jpg", gofakeit.Name()),
// 			},
// 		},
// 		Emails: []*domain.EmailInput{
// 			{
// 				Email:              "test@bewell.co.ke",
// 				CommunicationOptIn: true,
// 			},
// 		},
// 		PhysicalAddresses: []*domain.PhysicalAddress{
// 			{
// 				PhysicalAddress: gofakeit.Address().Address,
// 				MapsCode:        gofakeit.Address().Zip,
// 			},
// 		},
// 		PostalAddresses: []*domain.PostalAddress{
// 			{
// 				PostalAddress: gofakeit.Address().Address,
// 				PostalCode:    gofakeit.Address().City,
// 			},
// 		},
// 		Gender:        enumutils.GenderFemale.String(),
// 		Active:        true,
// 		MaritalStatus: domain.MaritalStatusS,
// 		Languages: []enumutils.Language{
// 			enumutils.LanguageEn,
// 			enumutils.LanguageSw,
// 		},
// 	}, msisdn, nil
// }

func getTestSimplePatientRegistration() (*domain.SimplePatientRegistrationInput, string, error) {
	otherNames := gofakeit.Name()
	dob := gofakeit.Date()
	msisdn := gofakeit.Phone()
	birthDate, err := scalarutils.NewDate(dob.Day(), int(dob.Month()), dob.Year())
	if err != nil {
		return nil, "", fmt.Errorf("can't create valid birth date: %v", err)
	}
	return &domain.SimplePatientRegistrationInput{
		ID: ksuid.New().String(),
		Names: []*domain.NameInput{
			{
				FirstName:  gofakeit.FirstName(),
				LastName:   gofakeit.LastName(),
				OtherNames: &otherNames,
			},
		},
		IdentificationDocuments: []*domain.IdentificationDocument{
			{
				DocumentType:   domain.IDDocumentTypeNationalID,
				DocumentNumber: strconv.Itoa(gofakeit.Number(11111111, 111111111)),
			},
		},
		BirthDate: *birthDate,
		PhoneNumbers: []*domain.PhoneNumberInput{
			{
				Msisdn:             msisdn,
				VerificationCode:   ksuid.New().String(),
				IsUssd:             true, // this will turn off OTP verification
				CommunicationOptIn: true,
			},
		},
		Photos: []*domain.PhotoInput{
			{
				PhotoContentType: enumutils.ContentTypeJpg,
				PhotoBase64data:  testPhotoBase64,
				PhotoFilename:    fmt.Sprintf("%s.jpg", gofakeit.Name()),
			},
		},
		Emails: []*domain.EmailInput{
			{
				Email:              "test@bewell.co.ke",
				CommunicationOptIn: true,
			},
		},
		PhysicalAddresses: []*domain.PhysicalAddress{
			{
				PhysicalAddress: gofakeit.Address().Address,
				MapsCode:        gofakeit.Address().Zip,
			},
		},
		PostalAddresses: []*domain.PostalAddress{
			{
				PostalAddress: gofakeit.Address().Address,
				PostalCode:    gofakeit.Address().City,
			},
		},
		Gender:        enumutils.GenderFemale.String(),
		Active:        true,
		MaritalStatus: domain.MaritalStatusS,
		Languages: []enumutils.Language{
			enumutils.LanguageEn,
			enumutils.LanguageSw,
		},
	}, msisdn, nil
}

func getTestPatient(ctx context.Context) (*domain.FHIRPatient, string, error) {
	simplePatientRegInput, msisdn, err := getTestSimplePatientRegistration()
	if err != nil {
		return nil, "", fmt.Errorf("can't genereate simple patient reg input: %v", err)
	}
	patientPayload, err := testInteractor.RegisterPatient(ctx, *simplePatientRegInput)
	if err != nil {
		return nil, "", fmt.Errorf("can't register patient: %v", err)
	}
	return patientPayload.PatientRecord, msisdn, nil
}

func TestGraphQLRegisterPatient(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	simplePatientRegInput, _, err := getTestSimplePatientRegistration()
	if err != nil {
		t.Errorf("can't genereate simple patient reg input: %v", err)
		return
	}

	patientRegInput, err := converterandformatter.StructToMap(simplePatientRegInput)
	if err != nil {
		t.Errorf("can't convert simple patient reg input to map: %v", err)
		return
	}
	validInput := map[string]interface{}{
		"input": patientRegInput,
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation SimplePatientRegistration($input: SimplePatientRegistrationInput!) {
	registerPatient(input: $input) {
		patientRecord {
			ID
			Identifier {
				ID
				Use
				Type {
					ID
					Text
					Coding {
						System
						Version
						Display
						Code
						UserSelected
					}
				}
				System
				Value
				Period {
					ID
					Start
					End
				}
			}
			Active
			Name {
				ID
				Use
				Text
				Family
				Given
				Prefix
				Suffix
				Period {
					ID
					Start
					End
				}
			}
			Telecom {
				ID
				System
				Value
				Use
				Rank
				Period {
					ID
					Start
					End
				}
			}
			Gender
			BirthDate
			Address {
				ID
				Use
				Type
				Text
				Line
				City
				District
				State
				PostalCode
				Country
				Period {
					ID
					Start
					End
				}
			}     
			Photo {
				Data
			}
			Contact {
				ID
				Relationship {
					ID
					Text
					Coding {
						System
						Version
						Display
						Code
						UserSelected
					}
				}
				Name {
					ID
					Use
					Text
					Family
					Given
					Prefix
					Suffix
					Period {
						ID
						Start
						End
					}
				}
				Gender
				Period {
					ID
					Start
					End
				}
				Address {
					ID
					Use
					Type
					Text
					Line
					City
					District
					State
					PostalCode
					Country
					Period {
						ID
						Start
						End
					}
				}
				Telecom {
					ID
					System
					Value
					Use
					Rank
					Period {
						ID
						Start
						End
					}
				}
			}
		}
	}
	}
					`,
					"variables": validInput,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errors, ok := data["errors"]
				log.Printf("The data returned...%v\n\n\n\n", data)
				if ok {
					t.Errorf("error not expected got error: %s", errors)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQFindPatientsByMSISDN(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, msisdn, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("error in getting test patient: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query FindPatientByMSISDN($msisdn: String!) {
						findPatientsByMSISDN(msisdn: $msisdn) {
						  edges {
							hasOpenEpisodes
							node {
							  ID
							  Active
							  Gender
							  BirthDate
							  Telecom{
								System
								Value
							  }
							  Name {
								Given
								Family
								Use
								Text
								Prefix
								Suffix
								Period {
								  Start
								  End
								}
							  }
							  Photo{
								ID
								ContentType
								Language
								Data
								URL
								Size
								Hash
								Title
							  }
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"msisdn": msisdn,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	usecase := testInteractor
	payload := &domain.PhoneNumberPayload{}
	usecase.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLFindPatients(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("could not get patient: %v", err)
		return
	}
	names := patient.Names()

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query FindPatient($search:String!) {
						findPatients(search:$search){
						  edges{     
							node{
							  Active
							  Name{
								Text
								Family
								Given
							  }
							  Telecom{
								Value
							  }
							}
						  }   
						}
					  }`,
					"variables": map[string]interface{}{
						"search": names,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					if key == "data" {
						_, present := nestedMap["findPatients"]
						if !present {
							t.Errorf("can't find patient data")
							return
						}

						patientMap, ok := nestedMap["findPatients"].(map[string]interface{})
						if !ok {
							t.Errorf("cannot cast key value of %v to type map[string]interface{}", patientMap)
							return
						}

						_, found := patientMap["edges"]
						if !found {
							t.Errorf("can't find patient edges data")
							return
						}
						edges, ok := patientMap["edges"].([]interface{})
						if !ok {
							t.Errorf("cannot cast key value of %v to type []interface{}", edges)
							return
						}

						if len(edges) == 0 {
							t.Error("can't find the patient")
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQGetPatient(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("can't get test patient: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `                    
					query GetPatientInfo($id: ID!) {
					  getPatient(id: $id) {
						hasOpenEpisodes
						openEpisodes{
						  ID
						  Status
						  Patient{
							Reference
							Type
							Display
						  }
						}
						patientRecord {
						  ID
						  Name {
							Text
						  }
						  Telecom {
							Value
						  }
						}
					  }
					}`,
					"variables": map[string]interface{}{
						"id": *patient.ID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLStartEpisodeByOTP(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}
	msisdn := interserviceclient.TestUserPhoneNumber
	otp, err := generateTestOTP(t, msisdn)
	if err != nil {
		log.Errorf("unable to get verified phone number and OTP")
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("could not get patient: %v", err)
		return
	}

	if patient.ID == nil {
		t.Errorf("nil patient ID")
		return
	}

	patientID := *patient.ID

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation StartOTPEpisode($input: OTPEpisodeCreationInput!) {
						startEpisodeByOTP(input: $input) {
						  episodeOfCare {
							ID
							Status
							Period {
							  Start
							}
							ManagingOrganization {
							  Display
							}
							Patient {
							  Identifier {
								Value
							  }
							  Display
							}
							Type {
							  Text
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"patientID":    patientID,
							"providerCode": testProviderCode,
							"otp":          otp,
							"msisdn":       msisdn,
							"fullAccess":   false,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					if key == "episodeOfCare" {
						if nestedMap["ID"] == "" {
							t.Errorf("got blank ID")
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLStartEpisodeByBreakGlass(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	msisdn := interserviceclient.TestUserPhoneNumber
	otp, err := generateTestOTP(t, msisdn)
	if err != nil {
		log.Errorf("unable to get verified phone number and OTP")
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("could not get patient: %v", err)
		return
	}

	if patient.ID == nil {
		t.Errorf("nil patient ID")
		return
	}

	patientID := *patient.ID

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation StartBreakGlassEpisode($input: BreakGlassEpisodeCreationInput!) {
						startEpisodeByBreakGlass(input: $input) {
							episodeOfCare {
							  ID
							  Status
							  Period {
								Start
							  }
							  ManagingOrganization {
								Display
							  }
							  Patient {
								Identifier {
								  Value
								}
								Display
							  }
							  Type {
								Text
							  }
							}
						  }
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"practitionerUID": testProviderUID,
							"patientID":       patientID,
							"providerCode":    testProviderCode,
							"otp":             otp,
							"providerPhone":   msisdn,
							"fullAccess":      false,
							"patientPhone":    interserviceclient.TestUserPhoneNumber,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					if key == "episodeOfCare" {
						if nestedMap["ID"] == "" {
							t.Errorf("got blank ID")
							return
						}
					}
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLUpgradeEpisode(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}
	msisdn := interserviceclient.TestUserPhoneNumber

	otp, err := generateTestOTP(t, msisdn)
	if err != nil {
		log.Errorf("unable to get verified phone number and OTP")
		return
	}

	episode, _, err := getTestEpisodeOfCare(
		ctx,
		interserviceclient.TestUserPhoneNumber,
		false, testProviderCode,
	)
	if err != nil {
		t.Errorf("can't create test episode: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation UpgradeEpisode($input: OTPEpisodeUpgradeInput!){
						upgradeEpisode(input: $input){
						  episodeOfCare{
							ID
						  }
						  totalVisits
						}
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"episodeID": episode.ID,
							"otp":       otp,
							"msisdn":    msisdn,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					if key == "episodeOfCare" {
						if nestedMap["ID"] == "" {
							t.Errorf("got blank ID")
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
}

func TestGraphQLEndEpisode(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	episode, _, err := getTestEpisodeOfCare(
		ctx,
		interserviceclient.TestUserPhoneNumber,
		false, testProviderCode,
	)
	if err != nil {
		t.Errorf("can't create test episode: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation EndEpisode($episodeID: String!) {
						endEpisode(episodeID: $episodeID)
					  }`,
					"variables": map[string]interface{}{
						"episodeID": episode.ID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					if nestedMap["endEpisode"] != true {
						t.Errorf("endEpisode result is not `true`")
						return
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
}

func TestGraphQLStartEncounter(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}
	episode, _, _, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `
					mutation StartEncounter($episodeID: String!) {
						startEncounter(episodeID: $episodeID)  
					  }
					`,
					"variables": map[string]interface{}{
						"episodeID": episode.ID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					if key == "data" {
						if nestedMap["startEncounter"] == "" {
							t.Errorf("got blank encounter ID")
							return
						}
					}

				}

			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLEndEncounter(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation EndExam($encounterID: String!) {
						endEncounter(encounterID: $encounterID)
					  }`,
					"variables": map[string]interface{}{
						"encounterID": encounterID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLOpenEpisodes(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	_, patient, err := getTestEpisodeOfCare(
		ctx,
		interserviceclient.TestUserPhoneNumber,
		false, testProviderCode,
	)
	if err != nil {
		t.Errorf("can't create test episode: %v", err)
		return
	}
	if patient.ID == nil {
		t.Errorf("nil patient ID")
		return
	}

	patientID := *patient.ID
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}
	patientRef := fmt.Sprintf("Patient/%s", patientID)

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `
					query searchOpenEpisodes($patientReference: String!) {
						openEpisodes(patientReference: $patientReference) {
						  ID
						  Status
						  Patient {
							Identifier {
							  Value
							}
							Display
						  }
						}
					  }
					`,
					"variables": map[string]interface{}{
						"patientReference": patientRef,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					if key == "data" {
						if nestedMap["openEpisodes"] == nil {
							t.Errorf("empty open episodes found")
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLSearchFHIREncounter(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	episode, _, _, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}

	// we have intermittent CI failures that could be related to replication
	// lag or latency issues on the backing data store (unproven).
	// If that is the case, then this sleep would reduce the failure rate.
	time.Sleep(time.Second * 5)

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query SearchFHIREncounter($params: Map!) {
						searchFHIREncounter(params: $params) {
						  edges {
							node {
							  ID
							  Status
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"params": map[string]interface{}{
							"episode-of-care": fmt.Sprintf(
								"EpisodeOfCare/%s", *episode.ID),
							"status": "in-progress",
							"_count": "1",
							"_sort":  "-_last_updated",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			if tt.wantErr {
				data := map[string]interface{}{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("unexpected data format: %s", string(dataResponse))
					return
				}
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				// example:
				// {"data":{"searchFHIREncounter":{"edges":[{"node":{"ID":"49e2f11e-7e1f-4107-8a2e-9c25c309122a","Status":"in-progress"}}]}}}
				data := map[string]map[string]map[string][]map[string]map[string]string{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("unexpected data format: %s", string(dataResponse))
					return
				}

				searchResult, present := data["data"]["searchFHIREncounter"]
				if !present {
					t.Errorf("key searchFHIREncounter not found in %v", data)
					return
				}

				edges, present := searchResult["edges"]
				if !present {
					t.Errorf("key edges not found in %v", searchResult)
					return
				}

				if len(edges) != 1 {
					t.Errorf("expected exactly one result, got %d", len(edges))
					return
				}

				result := edges[0]["node"]
				if result["Status"] != "in-progress" {
					t.Errorf("unexpected result status: %s", result["Status"])
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
}

func TestGraphqlOpenOrganizationEpisodes(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, err = getTestEpisodeOfCare(
		ctx,
		interserviceclient.TestUserPhoneNumber,
		false, testProviderCode,
	)
	if err != nil {
		t.Errorf("can't create test episode: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},

		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `
					query openOrganizationEpisodes($providerSladeCode: String!) {
						openOrganizationEpisodes(providerSladeCode: $providerSladeCode) {
						  ID
						  Status
						  Patient {
							Identifier {
							  Value
							}
							Display
						  }
						}
					  }
					`,
					"variables": map[string]interface{}{
						"providerSladeCode": testProviderCode,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLAddNextOfKin(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}
	msisdn := interserviceclient.TestUserPhoneNumber
	otp, err := generateTestOTP(t, msisdn)
	if err != nil {
		log.Errorf("unable to get verified phone number and OTP")
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("could not get patient: %v", err)
		return
	}

	if patient.ID == nil {
		t.Errorf("nil patient ID")
		return
	}

	patientID := *patient.ID

	var names []map[string]interface{}
	var phoneNumbers []map[string]interface{}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation AddKin($input: SimpleNextOfKinInput! ) {
						addNextOfKin(input:$input){
						  patientRecord{   
						  ID
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"patientID": patientID,
							"names": append(names, map[string]interface{}{
								"firstName": "Dennis",
								"lastName":  "Menace",
							}),
							"gender":       "male",
							"birthDate":    "1900-01-01",
							"relationship": "C",
							"active":       true,
							"phoneNumbers": append(phoneNumbers, map[string]interface{}{
								"msisdn":             msisdn,
								"verificationCode":   otp,
								"isUSSD":             false,
								"communicationOptIn": true,
							}),
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					if key == "data" {
						_, present := nestedMap["addNextOfKin"]
						if !present {
							t.Errorf("can't find next of kin payload")
							return
						}
						addNextOfKinMap := nestedMap["addNextOfKin"].(map[string]interface{})

						_, found := addNextOfKinMap["patientRecord"]
						if !found {
							t.Errorf("can't find patient record")
							return
						}
						patientRecordMap := addNextOfKinMap["patientRecord"].(map[string]interface{})
						if patientRecordMap["ID"] == "" {
							t.Errorf("got blank ID")
							return
						}
					}
				}

			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLUpdatePatient(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("can't get test patient: %v", err)
		return
	}

	newPatientInputData, _, err := getTestSimplePatientRegistration()
	if err != nil {
		t.Errorf("can't genereate simple patient reg inpit: %v", err)
		return
	}

	patientInputWithUpdatedData, err := converterandformatter.StructToMap(newPatientInputData)
	if err != nil {
		t.Errorf("can't convert simple patient reg input to map: %v", err)
		return
	}
	patientInputWithUpdatedData["id"] = *patient.ID

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation SimplePatientUpdate($input: SimplePatientRegistrationInput!) {
						updatePatient(input: $input) {
						  patientRecord {
							ID
							Identifier {
							  ID
							  Use
							  Type {
								ID
								Text
								Coding {
								  System
								  Version
								  Display
								  Code
								  UserSelected
								}
							  }
							  System
							  Value
							  Period {
								ID
								Start
								End
							  }
							}
							Active
							Name {
							  ID
							  Use
							  Text
							  Family
							  Given
							  Prefix
							  Suffix
							  Period {
								ID
								Start
								End
							  }
							}
							Telecom {
							  ID
							  System
							  Value
							  Use
							  Rank
							  Period {
								ID
								Start
								End
							  }
							}
							Gender
							BirthDate
							Address {
							  ID
							  Use
							  Type
							  Text
							  Line
							  City
							  District
							  State
							  PostalCode
							  Country
							  Period {
								ID
								Start
								End
							  }
							}     
							Photo {
							  ID
							  ContentType
							  Language
							  Data
							  URL
							  Size
							  Hash
							  Title
							  Creation
							}
							Contact {
							  ID
							  Relationship {
								ID
								Text
								Coding {
								  System
								  Version
								  Display
								  Code
								  UserSelected
								}
							  }
							  Name {
								ID
								Use
								Text
								Family
								Given
								Prefix
								Suffix
								Period {
								  ID
								  Start
								  End
								}
							  }
							  Gender
							  Period {
								ID
								Start
								End
							  }
							  Address {
								ID
								Use
								Type
								Text
								Line
								City
								District
								State
								PostalCode
								Country
								Period {
								  ID
								  Start
								  End
								}
							  }
							  Telecom {
								ID
								System
								Value
								Use
								Rank
								Period {
								  ID
								  Start
								  End
								}
							  }
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"input": patientInputWithUpdatedData,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLAddNHIF(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("error in getting test patient: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation AddNHIF($input: SimpleNHIFInput) {
						addNHIF(input: $input) {
							patientRecord {
							  ID
							}
						} 
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"patientID":        *patient.ID,
							"membershipNumber": gofakeit.BuzzWord(),
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLCreateUpdatePatientExtraInformation(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("could not get patient: %v", err)
		return
	}

	if patient.ID == nil {
		t.Errorf("nil patient ID")
		return
	}

	patientID := *patient.ID

	var emails []map[string]interface{}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation updatePatientExtraInformation($input: PatientExtraInformationInput!){
						createUpdatePatientExtraInformation(input: $input)
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"patientID":     patientID,
							"languages":     []string{"en"},
							"maritalStatus": "S",
							"emails": append(emails, map[string]interface{}{
								"email":              converterandformatter.GenerateRandomEmail(),
								"communicationOptIn": true,
							}),
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					if key == "data" {
						if nestedMap["createUpdatePatientExtraInformation"] == false {
							t.Errorf("expected true but got false instead")
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLVisitSummary(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	encounterID, _, err := patientVisitSummary(ctx)
	if err != nil {
		t.Errorf("can't create a visit summary")
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query VisitSummary($encounterID: String!) {
						visitSummary(encounterID: $encounterID)
					  }`,
					"variables": map[string]interface{}{
						"encounterID": encounterID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					for nestedKey := range nestedMap {
						if nestedKey == "visitSummary" {
							output, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast output to map[string]interface{}")
								return
							}

							for summary := range output {
								switch summary {
								case "Condition":
									condition, ok := output["Condition"].([]interface{})
									if !ok {
										t.Errorf("can't cast output to []interface{}")
										return
									}

									if len(condition) == 0 {
										t.Errorf("can't find a condition count")
										return
									}

									for _, c := range condition {
										conditionMap, ok := c.(map[string]interface{})
										if !ok {
											t.Errorf("can't cast to map[string]interface{}")
											return
										}
										if conditionMap["subject"] == nil {
											t.Errorf("can't find a condition")
											return
										}
									}

								case "Encounter":
									encounter, ok := output["Encounter"].([]interface{})
									if !ok {
										t.Errorf("can't cast encounter to []interface{}")
										return
									}

									if len(encounter) == 0 {
										t.Errorf("can't find a encounter count")
										return
									}

									for _, e := range encounter {
										encounterMap, ok := e.(map[string]interface{})
										if !ok {
											t.Errorf("can't cast to map[string]interface{}")
											return
										}
										if encounterMap["subject"] == nil {
											t.Errorf("can't find a encounter")
											return
										}
									}
								case "Observation":
									observation, ok := output["Observation"].([]interface{})
									if !ok {
										t.Errorf("can't cast observation to []interface{}")
										return
									}

									if len(observation) == 0 {
										t.Errorf("can't find a observation count")
										return
									}

									for _, o := range observation {
										observationMap, ok := o.(map[string]interface{})
										if !ok {
											t.Errorf("can't cast to map[string]interface{}")
											return
										}
										if observationMap["subject"] == nil {
											t.Errorf("can't find a observation")
											return
										}
									}
								case "Composition":
									composition, ok := output["Composition"].([]interface{})
									if !ok {
										t.Errorf("can't cast composition to []interface{}")
										return
									}

									if len(composition) == 0 {
										t.Errorf("can't find a composition count")
										return
									}

									for _, co := range composition {
										compositionMap, ok := co.(map[string]interface{})
										if !ok {
											t.Errorf("can't cast to map[string]interface{}")
											return
										}
										log.Printf("%v", compositionMap)
										if compositionMap["subject"] == nil {
											t.Errorf("can't find a composition")
											return
										}
									}
								case "ServiceRequest":
									serviceRequest, ok := output["ServiceRequest"].([]interface{})
									if !ok {
										t.Errorf("can't cast serviceRequest to []interface{}")
										return
									}

									if len(serviceRequest) == 0 {
										t.Errorf("can't find a serviceRequest count")
										return
									}

									for _, sr := range serviceRequest {
										serviceRequestMap, ok := sr.(map[string]interface{})
										if !ok {
											t.Errorf("can't cast to map[string]interface{}")
											return
										}
										if serviceRequestMap["subject"] == nil {
											t.Errorf("can't find a serviceRequest")
											return
										}
									}
								case "MedicationRequest":
									medicationRequest, ok := output["MedicationRequest"].([]interface{})
									if !ok {
										t.Errorf("can't cast medicationRequest to []interface{}")
										return
									}

									if len(medicationRequest) == 0 {
										t.Errorf("can't find a medicationRequest count")
										return
									}

									for _, mr := range medicationRequest {
										medicationRequestMap, ok := mr.(map[string]interface{})
										if !ok {
											t.Errorf("can't cast to map[string]interface{}")
											return
										}
										if medicationRequestMap["subject"] == nil {
											t.Errorf("can't find a medicationRequest")
											return
										}
									}
								case "AllergyIntolerance":
									allergyIntolerance, ok := output["AllergyIntolerance"].([]interface{})
									if !ok {
										t.Errorf("can't cast allergyIntolerance to []interface{}")
										return
									}

									if len(allergyIntolerance) == 0 {
										t.Errorf("can't find a allergyIntolerance count")
										return
									}

									for _, a := range allergyIntolerance {
										allergyIntoleranceMap, ok := a.(map[string]interface{})
										if !ok {
											t.Errorf("can't cast to map[string]interface{}")
											return
										}
										if allergyIntoleranceMap["encounter"] == nil {
											t.Errorf("can't find a allergyIntolerance")
											return
										}
									}
								}
							}

						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLPatientTimelineWithCount(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, episodeID, err := patientVisitSummary(ctx)
	if err != nil {
		t.Errorf("can't create a visit summary")
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query PatientTimelineWithCount($episodeID: String!, $count: Int!){
						patientTimelineWithCount(episodeID: $episodeID, count: $count)
					  }`,
					"variables": map[string]interface{}{
						"episodeID": episodeID,
						"count":     1,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					for nestedKey := range nestedMap {
						if nestedKey == "patientTimelineWithCount" {
							output, ok := nestedMap[nestedKey].([]interface{})
							if !ok {
								t.Errorf("can't cast output to []interface {}")
								return
							}

							for _, summary := range output {
								summaryMap, ok := summary.(map[string]interface{})
								if !ok {
									t.Errorf("can't cast output to map[string]interface{}")
									return
								}
								for summaryKey := range summaryMap {
									switch summaryKey {
									case "Condition":
										condition, ok := summaryMap["Condition"].([]interface{})
										if !ok {
											t.Errorf("can't cast summary to []interface{}")
											return
										}

										if len(condition) == 0 {
											t.Errorf("can't find a condition count")
											return
										}

										for _, c := range condition {
											conditionMap, ok := c.(map[string]interface{})
											if !ok {
												t.Errorf("can't cast to map[string]interface{}")
												return
											}
											if conditionMap["subject"] == nil {
												t.Errorf("can't find a condition")
												return
											}
										}

									case "Encounter":
										encounter, ok := summaryMap["Encounter"].([]interface{})
										if !ok {
											t.Errorf("can't cast encounter to []interface{}")
											return
										}

										if len(encounter) == 0 {
											t.Errorf("can't find a encounter count")
											return
										}

										for _, e := range encounter {
											encounterMap, ok := e.(map[string]interface{})
											if !ok {
												t.Errorf("can't cast to map[string]interface{}")
												return
											}
											if encounterMap["subject"] == nil {
												t.Errorf("can't find a encounter")
												return
											}
										}
									case "Observation":
										observation, ok := summaryMap["Observation"].([]interface{})
										if !ok {
											t.Errorf("can't cast observation to []interface{}")
											return
										}

										if len(observation) == 0 {
											t.Errorf("can't find a observation count")
											return
										}

										for _, o := range observation {
											observationMap, ok := o.(map[string]interface{})
											if !ok {
												t.Errorf("can't cast to map[string]interface{}")
												return
											}
											if observationMap["subject"] == nil {
												t.Errorf("can't find a observation")
												return
											}
										}
									case "Composition":
										composition, ok := summaryMap["Composition"].([]interface{})
										if !ok {
											t.Errorf("can't cast composition to []interface{}")
											return
										}

										if len(composition) == 0 {
											t.Errorf("can't find a composition count")
											return
										}

										for _, co := range composition {
											compositionMap, ok := co.(map[string]interface{})
											if !ok {
												t.Errorf("can't cast to map[string]interface{}")
												return
											}
											if compositionMap["subject"] == nil {
												t.Errorf("can't find a composition")
												return
											}
										}
									case "ServiceRequest":
										serviceRequest, ok := summaryMap["ServiceRequest"].([]interface{})
										if !ok {
											t.Errorf("can't cast serviceRequest to []interface{}")
											return
										}

										if len(serviceRequest) == 0 {
											t.Errorf("can't find a serviceRequest count")
											return
										}

										for _, sr := range serviceRequest {
											serviceRequestMap, ok := sr.(map[string]interface{})
											if !ok {
												t.Errorf("can't cast to map[string]interface{}")
												return
											}
											if serviceRequestMap["subject"] == nil {
												t.Errorf("can't find a serviceRequest")
												return
											}
										}
									case "MedicationRequest":
										medicationRequest, ok := summaryMap["MedicationRequest"].([]interface{})
										if !ok {
											t.Errorf("can't cast medicationRequest to []interface{}")
											return
										}

										if len(medicationRequest) == 0 {
											t.Errorf("can't find a medicationRequest count")
											return
										}

										for _, mr := range medicationRequest {
											medicationRequestMap, ok := mr.(map[string]interface{})
											if !ok {
												t.Errorf("can't cast to map[string]interface{}")
												return
											}
											if medicationRequestMap["subject"] == nil {
												t.Errorf("can't find a medicationRequest")
												return
											}
										}
									case "AllergyIntolerance":
										allergyIntolerance, ok := summaryMap["AllergyIntolerance"].([]interface{})
										if !ok {
											t.Errorf("can't cast allergyIntolerance to []interface{}")
											return
										}
										if len(allergyIntolerance) == 0 {
											t.Errorf("can't find a allergyIntolerance count")
											return
										}

										for _, a := range allergyIntolerance {
											allergyIntoleranceMap, ok := a.(map[string]interface{})
											if !ok {
												t.Errorf("can't cast to map[string]interface{}")
												return
											}
											if allergyIntoleranceMap["encounter"] == nil {
												t.Errorf("can't find a allergyIntolerance")
												return
											}
										}
									}
								}

							}

						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLProblemSummary(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("could not get patient: %v", err)
		return
	}

	if patient.ID == nil {
		t.Errorf("nil patient ID")
		return
	}

	patientID := *patient.ID

	_, err = createTestCondition(ctx, encounterID, patientID)
	if err != nil {
		t.Errorf("error creating a test condition: %v", err)
		return
	}

	// we have intermittent CI failures that could be related to replication
	// lag or latency issues on the backing data store (unproven).
	// If that is the case, then this sleep would reduce the failure rate.
	time.Sleep(time.Second * 5)

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query ProblemSummary($patientID: String!) {
						problemSummary(patientID: $patientID)
					}`,
					"variables": map[string]interface{}{
						"patientID": patientID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			if tt.wantErr {
				data := map[string]interface{}{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned: %s", string(dataResponse))
					return
				}
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}
			if !tt.wantErr {
				data := map[string]map[string][]string{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned: %s", string(dataResponse))
					return
				}
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %s", errMsg)
					return
				}

				for key := range data {
					nestedMap, present := data[key]
					if !present {
						t.Errorf("key %s not found in %v", key, data)
						return
					}
					expected := map[string][]string{
						"problemSummary": {
							"Pulmonary Tuberculosis",
						},
					}
					if !reflect.DeepEqual(expected, nestedMap) {
						t.Errorf("expected %v got %v", expected, nestedMap)
						return
					}
				}

			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLCreateFHIRMedicationRequest(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}

	patientName := patient.Names()
	requester := gofakeit.Name()
	dateRecorded := time.Now().Format(dateFormat)

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
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
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "createFHIRMedicationRequest" {
							output, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast output to map[string]interface{}")
								return
							}

							resource, ok := output["resource"].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast resource to map[string]interface{}")
								return
							}

							log.Printf("resource: %v", resource)

							idI, prs := resource["ID"]
							if !prs {
								t.Errorf("ID not present in medication request resource")
								return
							}
							id, ok := idI.(string)
							if !ok {
								t.Errorf("nil id")
								return
							}

							if id == "" {
								t.Errorf("blank medication request ID")
								return
							}
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLUpdateFHIRMedicationRequest(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}

	// create FHIRMedicationRequest
	medicationRequestID, err := getTestFHIRMedicationRequestID(ctx, encounterID)
	if err != nil {
		t.Errorf("error creating test medication request")
		return
	}
	if medicationRequestID == "" {
		t.Errorf("failed to create medication request: %v", err)
		return
	}

	patientName := patient.Names()
	requester := gofakeit.Name()
	dateRecorded := time.Now().Format(dateFormat)

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `
					mutation UpdateMedicationRequest($input: FHIRMedicationRequestInput!) {
						updateFHIRMedicationRequest(input: $input) {
						  resource {
							ID
							Status
							Intent
							Priority
							Subject {
							  Reference
							  Type
							  Display
							}
							Encounter {
							  Reference
							  Type
							  Display
							}
							MedicationCodeableConcept {
							  Text
							  Coding {
								System
								Code
								Display
								UserSelected
							  }
							}
							DosageInstruction {
							  Text
							  PatientInstruction
							}
							AuthoredOn
							Note{
							  AuthorString
							  Text
							}
							Requester{
							  Display
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"ID":       medicationRequestID,
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
							"Note": []map[string]interface{}{
								{
									"AuthorString": requester,
									"Text":         gofakeit.HipsterSentence(10),
								},
							},
							"Requester": map[string]interface{}{
								"Display": requester,
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "updateFHIRMedicationRequest" {
							output, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast output to map[string]interface{}")
								return
							}

							resource, ok := output["resource"].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast resource to map[string]interface{}")
								return
							}

							log.Printf("resource: %v", resource)

							id, prs := resource["ID"]
							if !prs {
								t.Errorf("ID not present in medication request resource")
								return
							}
							if id == "" {
								t.Errorf("blank medication request ID")
								return
							}

							if id != medicationRequestID {
								t.Errorf("wrong medicationRequestUpdate")
							}
						}
					}
				}

			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

// TODO restore this
func TestGraphQLDeleteFHIRMedicationRequest(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}

	medicationRequestID, err := getTestFHIRMedicationRequestID(ctx, encounterID)
	if err != nil {
		t.Errorf("error creating test medication request")
		return
	}
	if medicationRequestID == "" {
		t.Errorf("empty medication request ID")
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `
					mutation DeleteFHIRMedicationRequest($id: ID!) {
						deleteFHIRMedicationRequest(id: $id)
					  }`,
					"variables": map[string]interface{}{
						"id": medicationRequestID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: pass unknown variables",
			args: args{
				query: map[string]interface{}{
					"query": `
					mutation DeleteFHIRMedicationRequest($id: ID!) {
						deleteFHIRMedicationRequest(id: $id)
					  }`,
					"variables": map[string]interface{}{
						"ID": "some_unknown_id",
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid: wrong ID provided",
			args: args{
				query: map[string]interface{}{
					"query": `
					mutation DeleteFHIRMedicationRequest($id: ID!) {
						deleteFHIRMedicationRequest(id: $id)
					  }`,
					"variables": map[string]interface{}{
						"id": "some_unknown_id",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "deleteFHIRMedicationRequest" {
							result, resultFound := nestedMap[nestedKey]
							if !resultFound {
								t.Errorf("response[deleteFHIRMedicationRequest] = ' '")
								return
							}
							resultBool, castOk := result.(bool)
							if !castOk {
								t.Errorf("failed to delete medicationRequest")
								return
							}
							if !resultBool {
								t.Errorf("failed to delete medicationRequest")
								return
							}
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRMedicationRequest(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	// create a medication request
	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}

	medicationRequestID, err := getTestFHIRMedicationRequestID(ctx, encounterID)
	if err != nil {
		t.Errorf("error creating test medication request")
		return
	}
	if medicationRequestID == "" {
		t.Errorf("empty medicationrequest ID")
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `
					query SearchMedicationRequests($params: Map!) {
						searchFHIRMedicationRequest(params: $params) {
						  edges {
							node {
							  ID
							  Status
							  Intent
							  Priority
							  Subject {
								Reference
								Type
								Display
							  }
							  MedicationCodeableConcept{
								Text
								Coding{
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							  DosageInstruction{
								Text
								PatientInstruction
							  }
							  Requester{
								Display
							  }
							  Encounter {
								Reference
								Type
								Display
							  }
							  SupportingInformation {
								ID
								Reference
								Display
							  }
							  AuthoredOn
							  Note{
								AuthorString
							  }
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"params": map[string]interface{}{
							"status": "active",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: wrong parameter",
			args: args{
				query: map[string]interface{}{
					"query": `
					query SearchMedicationRequests($params: Map!) {
						searchFHIRMedicationRequest(params: $params) {
						  edges {
							node {
							  ID
							  Status
							  Intent
							  Priority
							  Subject {
								Reference
								Type
								Display
							  }
							  MedicationCodeableConcept{
								Text
								Coding{
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							  DosageInstruction{
								Text
								PatientInstruction
							  }
							  Requester{
								Display
							  }
							  Encounter {
								Reference
								Type
								Display
							  }
							  SupportingInformation {
								ID
								Reference
								Display
							  }
							  AuthoredOn
							  Note{
								AuthorString
							  }
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"params": map[string]interface{}{
							"status": "111111111",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "searchFHIRMedicationRequest" {
							resultI, resultFound := nestedMap[nestedKey]
							if !resultFound {
								t.Errorf("response[searchFHIRMedicationRequest] = ' '")
								return
							}
							result, resultConvert := resultI.(map[string]interface{})
							if !resultConvert {
								t.Errorf("cannot convert result to a map[string]interface{}")
								return
							}

							edgesI, edgesFound := result["edges"]
							if !edgesFound {
								t.Errorf("no medication request was returned")
								return
							}
							if edgesI == nil {
								t.Errorf("no medication request was found")
								return
							}

						}
					}
				}

			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLCreateFHIRAllergyIntolerance(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}

	patientName := patient.Names()
	requester := gofakeit.Name()
	dateRecorded := time.Now().Format(dateFormat)
	recordingDoctor := gofakeit.Name()
	substanceID := "1234"
	substanceDisplayName := gofakeit.Name()

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `
					mutation CreateAllergy($input: FHIRAllergyIntoleranceInput!) {
						createFHIRAllergyIntolerance(input:$input) {
						resource {
						  ID
						}
					  }
					}  
					`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"Type":        "allergy",
							"Criticality": "high",
							"ClinicalStatus": map[string]interface{}{
								"Text": "Panadol Extra",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical",
										"Code":         "active",
										"Display":      "Active",
										"UserSelected": false,
									},
								},
							},
							"VerificationStatus": map[string]interface{}{
								"Text": "Panadol Extra",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/allergyintolerance-verification",
										"Code":         "confirmed",
										"Display":      "confirmed",
										"UserSelected": false,
									},
								},
							},
							"Encounter": map[string]interface{}{
								"Reference": fmt.Sprintf("Encounter/%s", encounterID),
								"Type":      "Encounter",
								"Display":   fmt.Sprintf("Encounter/%s", encounterID),
							},
							"Patient": map[string]interface{}{
								"Reference": fmt.Sprintf("Patient/%s", *patient.ID),
								"Type":      "Patient",
								"Display":   patientName,
							},
							"Recorder": map[string]interface{}{
								"Display": recordingDoctor,
							},
							"Asserter": map[string]interface{}{
								"Display": recordingDoctor,
							},
							"Note": []map[string]interface{}{
								{
									"AuthorString": requester,
									"Text":         gofakeit.HipsterSentence(10),
								},
							},
							"Reaction": []map[string]interface{}{
								{
									"Description": requester,
									"Severity":    "severe",
									"Substance": map[string]interface{}{
										"Text": "Panadol Extra",
										"Coding": []map[string]interface{}{
											{
												"System":       "OCL:/orgs/CIEL/sources/CIEL/'",
												"Code":         substanceID,
												"Display":      substanceDisplayName,
												"UserSelected": true,
											},
										},
									},
									"Manifestation": []map[string]interface{}{
										{
											"Text": "Panadol Extra",
											"Coding": []map[string]interface{}{
												{
													"System":       "OCL:/orgs/CIEL/sources/CIEL/'",
													"Code":         substanceID,
													"Display":      substanceDisplayName,
													"UserSelected": true,
												},
											},
										},
									},
								},
							},
							"Code": map[string]interface{}{
								"Text": "Panadol Extra",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/allergyintolerance-verification",
										"Code":         substanceID,
										"Display":      substanceDisplayName,
										"UserSelected": false,
									},
								},
							},
							"RecordedDate": dateRecorded,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "createFHIRAllergyIntolerance" {
							output, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast output to map[string]interface{}")
								return
							}

							resource, ok := output["resource"].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast resource to map[string]interface{}")
								return
							}

							log.Printf("resource: %v", resource)

							id, prs := resource["ID"]
							if !prs {
								t.Errorf("ID not present in allergy intolerance resource")
								return
							}
							if id == "" {
								t.Errorf("blank allergy intolerance request ID")
								return
							}
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLUpdateFHIRAllergyIntolerance(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}

	allergyID, err := createTestAllergy(ctx, patient, encounterID)
	if err != nil {
		t.Errorf("can't get test allergy intolerance")
		return
	}
	if allergyID == "" {
		t.Errorf("emtyp allergy intolerance ID")
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query": `
					mutation UpdateAllergy($input: FHIRAllergyIntoleranceInput!) {
						updateFHIRAllergyIntolerance(input: $input) {
						  resource {
							ID
						  }
						}
					  }
					  `,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"ID":          allergyID,
							"Type":        "allergy",
							"Category":    []string{},
							"Criticality": domain.AllergyIntoleranceCriticalityEnumLow,
							"ClinicalStatus": map[string]interface{}{
								"Text": "Resolved",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/condition-clinical",
										"Code":         "resolved",
										"Display":      "Resolved",
										"UserSelected": true,
									},
								},
							},
							"VerificationStatus": map[string]interface{}{
								"Text": "Panadol Extra",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/allergyintolerance-verification",
										"Code":         "confirmed",
										"Display":      "confirmed",
										"UserSelected": false,
									},
								},
							},
							"Code": map[string]interface{}{
								"Text": "Panadol Extra",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/allergyintolerance-verification",
										"Code":         "1234",
										"Display":      "substanceDisplayName",
										"UserSelected": false,
									},
								},
							},
							"Patient": map[string]interface{}{
								"Reference": fmt.Sprintf("Patient/%s", *patient.ID),
								"Type":      "Patient",
								"Display":   "patientName",
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "updateFHIRAllergyIntolerance" {
							output, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast output to map[string]interface{}")
								return
							}

							resource, ok := output["resource"].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast resource to map[string]interface{}")
								return
							}

							idI, prs := resource["ID"]
							if !prs {
								t.Errorf("ID not present in allergy intolerance resource")
								return
							}
							id, idConvert := idI.(string)
							if !idConvert {
								t.Errorf("mulformed ID returned")
								return
							}
							if id == "" {
								t.Errorf("blank allergy intolerance request ID")
								return
							}

							if id != allergyID {
								t.Errorf("wrong allergy intolerance request ID returned")
								return
							}
						}
					}
				}

			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRAllergyIntolerance(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}
	allergyID, err := createTestAllergy(ctx, patient, encounterID)
	if err != nil {
		t.Errorf("can't get test allergy intolerance")
		return
	}
	if allergyID == "" {
		t.Errorf("empty allergy intolerance ID")
		return
	}

	encounterRef := fmt.Sprintf("Encounter/%s", encounterID)
	patientRef := fmt.Sprintf("Patient/%s", *patient.ID)

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `
					query AllergySearch($params: Map!) {
						searchFHIRAllergyIntolerance(params: $params) {
						  edges {
							node {
							  ID
							  Type
							  RecordedDate
							  Category
							  Criticality
							  ClinicalStatus{
								Text
								Coding{
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							  VerificationStatus{
								Text
								Coding{
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							  Patient{
								Reference
								Type
								Display
							  }
							  Code {
								Text
								Coding {
								  Code
								  System
								  Display
								}
							  }
							  Encounter{
								Reference
								Type
							  }
							  Asserter{
								Display
							  }
							  Note{
								AuthorString
								Text
							  }
							  Reaction{
								Description
								Severity
								Substance{
								  Text
								  Coding {
									System
									Code
									Display
									UserSelected
								  }
								}         
							  }
							  Recorder{
								Display
							  }
							}
						  }
						}
					  }
					  `,
					"variables": map[string]interface{}{
						"params": map[string]interface{}{
							"patient":   patientRef,
							"encounter": encounterRef,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "searchFHIRAllergyIntolerance" {
							resultI, resultFound := nestedMap[nestedKey]
							if !resultFound {
								t.Errorf("response[searchFHIRAllergyIntolerance] = ' '")
								return
							}
							result, resultConvert := resultI.(map[string]interface{})
							if !resultConvert {
								t.Errorf("cannot convert result to a map[string]interface{}")
								return
							}

							edgesI, edgesFound := result["edges"]
							if !edgesFound {
								t.Errorf("no allergy intollerance request was returned")
								return
							}
							if edgesI == nil {
								log.Errorf("no allergy intollerance request was found")
								return
							}

						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLCreateFHIRCondition(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}

	recordedDate := time.Now().Format(dateFormat)

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation CreateFHIRCondition($input: FHIRConditionInput!) {
						createFHIRCondition(input: $input) {
						  resource {
							ID
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"Code": map[string]interface{}{
								"Coding": []map[string]interface{}{
									{
										"System":       "OCL:/orgs/CIEL/sources/CIEL/",
										"Code":         "113488",
										"Display":      "Pulmonary Tuberculosis",
										"UserSelected": true,
									},
								},
								"Text": "Pulmonary Tuberculosis",
							},
							"ClinicalStatus": map[string]interface{}{
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/condition-clinical",
										"Code":         "active",
										"Display":      "Active",
										"UserSelected": false,
									},
								},
								"Text": "Active",
							},
							"VerificationStatus": map[string]interface{}{
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/condition-ver-status",
										"Code":         "confirmed",
										"Display":      "Confirmed",
										"UserSelected": false,
									},
								},
								"Text": "Active",
							},
							"RecordedDate": recordedDate,
							"Category": []map[string]interface{}{
								{
									"Coding": []map[string]interface{}{
										{
											"System":       "http://terminology.hl7.org/CodeSystem/condition-category",
											"Code":         "encounter-diagnosis",
											"Display":      "encounter-diagnosis",
											"UserSelected": false,
										},
									},
									"Text": "encounter-diagnosis",
								},
							},
							"Subject": map[string]interface{}{
								"Reference": fmt.Sprintf("Patient/%s", *patient.ID),
								"Type":      "Patient",
								"Display":   fmt.Sprintf("Patient/%s", *patient.ID),
							},
							"Encounter": map[string]interface{}{
								"Reference": fmt.Sprintf("Encounter/%s", encounterID),
								"Type":      "Encounter",
								"Display":   "Encounter",
							},
							"Note": []map[string]interface{}{
								{
									"AuthorString": gofakeit.Name(),
									"Text":         gofakeit.HipsterSentence(20),
								},
							},
							"Recorder": map[string]interface{}{
								"Display": gofakeit.Name(),
							},
							"Asserter": map[string]interface{}{
								"Display": gofakeit.Name(),
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
}

func TestGraphQUpdateFHIRCondition(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}

	patientName := patient.Names()
	recorder := gofakeit.Name()
	asserter := gofakeit.Name()
	dateRecorded := time.Now().Format(dateFormat)

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation UpdateFHIRCondition($input: FHIRConditionInput!) {
						updateFHIRCondition(input: $input) {
						  resource {
							ID
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"ID":           "113488",
							"RecordedDate": dateRecorded,
							"ClinicalStatus": map[string]interface{}{
								"Text": "Resolved",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/condition-clinical",
										"Code":         "resolved",
										"Display":      "Resolved",
										"UserSelected": true,
									},
								},
							},
							"VerificationStatus": map[string]interface{}{
								"Text": "Active",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://terminology.hl7.org/CodeSystem/condition-ver-status",
										"Code":         "confirmed",
										"Display":      "confirmed",
										"UserSelected": true,
									},
								},
							},
							"Category": []map[string]interface{}{
								{
									"Text": "encounter-diagnosis",
									"Coding": []map[string]interface{}{
										{
											"System":       "http://terminology.hl7.org/CodeSystem/condition-category",
											"Code":         "encounter-diagnosis",
											"Display":      "encounter-diagnosis",
											"UserSelected": true,
										},
									},
								},
							},
							"Code": map[string]interface{}{
								"Coding": []map[string]interface{}{
									{
										"System":       "OCL:/orgs/CIEL/sources/CIEL/",
										"Code":         "113488",
										"Display":      "Pulmonary Tuberculosis",
										"UserSelected": true,
									},
								},
								"Text": "Pulmonary Tuberculosis",
							},
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
							"Recorder": map[string]interface{}{
								"Display": recorder,
							},
							"Asserter": map[string]interface{}{
								"Display": asserter,
							},
							"Note": []map[string]interface{}{
								{
									"AuthorString": recorder,
									"Text":         "A good reason",
								},
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "updateFHIRCondition" {
							output, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast output to map[string]interface{}")
								return
							}

							resource, ok := output["resource"].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast resource to map[string]interface{}")
								return
							}

							id, prs := resource["ID"]
							if !prs {
								t.Errorf("ID not present in service request resource")
								return
							}
							if id == "" {
								t.Errorf("blank service request ID")
								return
							}
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRCondition(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("could not get test patient: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query ConditionSearch($params: Map!) {
						searchFHIRCondition(params: $params) {
						  edges {
							node {
							  ID
							  RecordedDate
							  ClinicalStatus{
								Text
								Coding{
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							  VerificationStatus{
								Text
								Coding{
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							  Category{
								Text
								Coding{
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							  Subject{
								Reference
								Type
								Display
							  }
							  Encounter{
								Reference
								Type
							  }
							  Evidence{
								Detail{
								  Display
								}
							  }
							  Asserter{
								Display
							  }
							  Note{
								AuthorString
								Text
							  }
							  Severity{
								ID
								Coding{
								  Display
								}
								Text
							  }
							  OnsetString
							  Recorder{
								Display
							  }
							  Code{
								Coding{
								  System
								  Code
								  Display
								  UserSelected
								}
								Text
							  }
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"params": map[string]interface{}{
							"patient": fmt.Sprintf("Patient/%s", *patient.ID),
							"_count":  "1",
							"_sort":   "-_last_updated",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					if key == "data" {
						_, present := nestedMap["searchFHIRCondition"]
						if !present {
							t.Errorf("can't find medication request data")
							return
						}

						conditionMap, ok := nestedMap["searchFHIRCondition"].(map[string]interface{})
						if !ok {
							t.Errorf("cannot cast key value of %v to type map[string]interface{}", conditionMap)
							return
						}

						_, found := conditionMap["edges"]
						if !found {
							t.Errorf("can't find FHIR condition edges data")
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLCreateFHIRServiceRequest(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}

	patientName := patient.Names()
	requester := gofakeit.Name()

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation CreateServiceRequest($input: FHIRServiceRequestInput!) {
						createFHIRServiceRequest(input: $input) {
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
							"SupportingInfo": []map[string]interface{}{
								{
									"ID":        "113488",
									"Reference": fmt.Sprintf("Encounter/%s", encounterID),
									"Display":   "Pulmonary Tuberculosis",
								},
							},
							"Category": []map[string]interface{}{
								{
									"Text": "Laboratory procedure",
									"Coding": []map[string]interface{}{
										{
											"System":       "OCL:/orgs/CIEL/sources/CIEL/",
											"Code":         "108252007",
											"Display":      "Laboratory procedure",
											"UserSelected": true,
										},
									},
								},
							},
							"Requester": map[string]interface{}{
								"Display": requester,
							},
							"Code": map[string]interface{}{
								"Text": "Hospital re-admission",
								"Coding": []map[string]interface{}{
									{
										"System":       "OCL:/orgs/CIEL/sources/CIEL/",
										"Code":         "417005",
										"Display":      "Hospital re-admission",
										"UserSelected": true,
									},
								},
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}
			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "createFHIRServiceRequest" {
							output, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast output to map[string]interface{}")
								return
							}

							resource, ok := output["resource"].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast resource to map[string]interface{}")
								return
							}

							id, present := resource["ID"]
							if !present {
								t.Errorf("ID not present in service request resource")
								return
							}
							if id == "" {
								t.Errorf("blank service request ID")
								return
							}
						}
					}
				}

			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLDeleteFHIRServiceRequest(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}
	serviceRequest, _, err := getTestServiceRequest(ctx, encounterID)
	if err != nil {
		t.Errorf("error creating test service request: %v", err)
		return
	}

	if serviceRequest.ID == nil {
		t.Errorf("can't find service request ID")
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation DeleteFHIRServiceRequest($id: ID!) {
						deleteFHIRServiceRequest(id: $id)
					  }`,
					"variables": map[string]interface{}{
						"id": *serviceRequest.ID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					if nestedMap["deleteFHIRServiceRequest"] == false {
						t.Errorf("service request was not deleted successfully")
						return
					}

				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLSearchFHIRServiceRequest(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %v", err)
		return
	}
	serviceRequest, patientID, err := getTestServiceRequest(ctx, encounterID)
	if err != nil {
		t.Errorf("error creating test service request: %v", err)
		return
	}

	if serviceRequest.ID == nil {
		t.Errorf("can't find service request ID")
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query SearchFHIRServiceRequests($params: Map!) {
						searchFHIRServiceRequest(params: $params) {
						  edges {
							node {
							  ID
							  Status
							  Intent
							  Priority
							  Subject {
								Reference
								Type
								Display
							  }
							  Encounter {
								Reference
								Type
								Display
							  }
							  SupportingInfo{
								ID
								Reference
								Display
							  }
							  Requester{
								Display
							  }
							  Code{
								Coding{
								  Display
								  Code
								}
							  }
							  Category {
								Text
								Coding {
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"params": map[string]interface{}{
							"patient":   patientID,
							"encounter": fmt.Sprintf("Encounter/%s", encounterID),
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					_, ok = nestedMap["searchFHIRServiceRequest"].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast nested map key value of %v to type map[string]interface{}", key)
						return
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQCreateFHIRObservation(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	instantRecorded := time.Now().Format(instantFormat)
	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("can't create test encounter: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation AddObservation($input: FHIRObservationInput!) {
						createFHIRObservation(input: $input) {
						  resource {
							ID
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"Status": "preliminary",
							"Category": []map[string]interface{}{
								{
									"Text": "Vital Signs",
									"Coding": []map[string]interface{}{
										{
											"Code":         "vital-signs",
											"System":       "http://terminology.hl7.org/CodeSystem/observation-category",
											"Display":      "Vital Signs",
											"UserSelected": false,
										},
									},
								},
							},
							"Code": map[string]interface{}{
								"Text": "Body Weight",
								"Coding": []map[string]interface{}{
									{
										"Display":      "Body Weight",
										"Code":         "29463-7",
										"System":       "http://loinc.org",
										"UserSelected": true,
									},
								},
							},
							"ValueQuantity": map[string]interface{}{
								"Value":  72,
								"Unit":   "kg",
								"System": "http://unitsofmeasure.org",
								"Code":   "kg",
							},
							"ReferenceRange": []map[string]interface{}{
								{
									"Low": map[string]interface{}{
										"Value":  "0",
										"Unit":   "kg",
										"System": "http://unitsofmeasure.org",
										"Code":   "kg",
									},
									"High": map[string]interface{}{
										"Value":  "300",
										"Unit":   "kg",
										"System": "http://unitsofmeasure.org",
										"Code":   "kg",
									},
									"Text": "0kg to 300kg",
									"Type": map[string]interface{}{
										"Text": "Normal Range",
										"Coding": []map[string]interface{}{
											{
												"Code":         "normal",
												"UserSelected": false,
												"System":       "http://terminology.hl7.org/CodeSystem/referencerange-meaning",
												"Display":      "Normal Range",
											},
										},
									},
								},
							},
							"Issued":           instantRecorded,
							"EffectiveInstant": instantRecorded,
							"Encounter": map[string]interface{}{
								"Reference": fmt.Sprintf("Encounter/%s", encounterID),
								"Type":      "Encounter",
								"Display":   fmt.Sprintf("Encounter/%s", encounterID),
							},
							"Subject": map[string]interface{}{
								"Reference": fmt.Sprintf("Patient/%s", *patient.ID),
								"Type":      "Patient",
								"Display":   patient.Names(),
							},
							"Interpretation": []map[string]interface{}{
								{
									"Text": "Normal",
									"Coding": []map[string]interface{}{
										{
											"Display":      "Normal",
											"Code":         "N",
											"System":       "http://terminology.hl7.org/CodeSystem/v3-ObservationInterpretation",
											"UserSelected": false,
										},
									},
								},
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRObservation(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("can't create test encounter: %v", err)
		return
	}
	_, patient, status, err := createFHIRTestObservation(ctx, encounterID)
	if err != nil {
		t.Errorf("can't create FHIR test observation: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query SearchObservations($params: Map!) {
						searchFHIRObservation(params: $params) {
						  edges {
							node {
							  ID
							Status
							Code {
							  Text
							  Coding {
								Display
								Code
								System
								UserSelected
							  }
							}
							ValueQuantity {
							  Value
							  Unit
							  System
							  Code
							}
							ReferenceRange {
							  Text
							  Low {
								Value
								Unit
								System
								Code
							  }
							  High {
								Value
								Unit
								System
								Code
							  }
							  Type {
								Text
								Coding {
								  Code
								  UserSelected
								  System
								  Display
								}
							  }
							}
							Interpretation {
							  Text
							  Coding {
								Display
								Code
								System
								UserSelected
							  }
							}
							Category {
							  Text
							  Coding {
								Code
								System
								Display
								UserSelected
							  }
							}
							Issued
							EffectiveInstant
							Subject {
							  Reference
							  Type
							  Display
							}
							Encounter {
							  Reference
							  Type
							  Display
							  }
							}
						  }
						  }
					  }`,
					"variables": map[string]interface{}{
						"params": map[string]interface{}{
							"patient":   fmt.Sprintf("Patient/%s", *patient.ID),
							"status":    status.String(),
							"encounter": fmt.Sprintf("Encounter/%s", encounterID),
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
}

func TestGraphQCreateFHIRComposition(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}

	recorded := time.Now().Format(dateFormat)

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation CreateComposition($input: FHIRCompositionInput!) {
						createFHIRComposition(input: $input) {
						  resource {
							ID
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"Status": "preliminary",
							"Date":   recorded,
							"Title":  gofakeit.HipsterSentence(10),
							"Type": map[string]interface{}{
								"Text": "Consult Note",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://loinc.org",
										"Code":         "11488-4",
										"Display":      "Consult Note",
										"UserSelected": false,
									},
								},
							},
							"Category": []map[string]interface{}{
								{
									"Text": "Consult Note",
									"Coding": []map[string]interface{}{
										{
											"System":       "http://loinc.org",
											"Code":         "11488-4",
											"Display":      "Consult Note",
											"UserSelected": false,
										},
									},
								},
							},
							"Subject": map[string]interface{}{
								"Reference": fmt.Sprintf("Patient/%s", *patient.ID),
								"Type":      "Patient",
								"Display":   patient.Names(),
							},
							"Encounter": map[string]interface{}{
								"Reference": fmt.Sprintf("Encounter/%s", encounterID),
								"Type":      "Encounter",
								"Display":   fmt.Sprintf("Encounter/%s", encounterID),
							},
							"Section": []map[string]interface{}{
								{
									"Title": "patientHistory",
									"Text": map[string]interface{}{
										"Status": "generated",
										"Div":    gofakeit.HipsterSentence(10),
									},
								},
							},
							"Author": []map[string]interface{}{
								{
									"Display": gofakeit.Name(),
								},
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errorMessage, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got error: %s", errorMessage)
					return
				}

				log.Printf("response: \n%s\n", data)
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					for nestedKey := range nestedMap {
						if nestedKey == "createFHIRComposition" {
							compositionData, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast nested composition data to map")
								return
							}
							if compositionData["id"] == "" {
								t.Errorf("got back blank ID for new composition")
								return
							}
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQUpdateFHIRComposition(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("can't create test encounter: %v", err)
		return
	}

	composition, _, err := createTestFHIRComposition(ctx, encounterID)
	if err != nil {
		t.Errorf("can't create test composition: %v", err)
		return
	}
	recorded := time.Now().Format(dateFormat)

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `mutation UpdateComposition($input: FHIRCompositionInput!) {
						updateFHIRComposition(input: $input) {
						  resource {
							ID  
							Status
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"ID":     *composition.ID,
							"Status": "final", // this is the edit...make the composition final
							"Date":   recorded,
							"Title":  gofakeit.HipsterSentence(10),
							"Type": map[string]interface{}{
								"Text": "Consult Note",
								"Coding": []map[string]interface{}{
									{
										"System":       "http://loinc.org",
										"Code":         "11488-4",
										"Display":      "Consult Note",
										"UserSelected": false,
									},
								},
							},
							"Category": []map[string]interface{}{
								{
									"Text": "Consult Note",
									"Coding": []map[string]interface{}{
										{
											"System":       "http://loinc.org",
											"Code":         "11488-4",
											"Display":      "Consult Note",
											"UserSelected": false,
										},
									},
								},
							},
							"Subject": map[string]interface{}{
								"Reference": fmt.Sprintf("Patient/%s", *patient.ID),
								"Type":      "Patient",
								"Display":   patient.Names(),
							},
							"Encounter": map[string]interface{}{
								"Reference": fmt.Sprintf("Encounter/%s", encounterID),
								"Type":      "Encounter",
								"Display":   fmt.Sprintf("Encounter/%s", encounterID),
							},
							"Section": []map[string]interface{}{
								{
									"Title": "patientHistory",
									"Text": map[string]interface{}{
										"Status": "generated",
										"Div":    gofakeit.HipsterSentence(10),
									},
								},
							},
							"Author": []map[string]interface{}{
								{
									"Display": gofakeit.Name(),
								},
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					for nestedKey := range nestedMap {
						if nestedKey == "createFHIRComposition" {
							compositionData, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast nested composition data to map")
								return
							}
							if compositionData["ID"] == "" {
								t.Errorf("got back blank ID for new composition")
								return
							}
							if compositionData["Status"] != "final" {
								t.Errorf("got back non final status after updating composition status to final")
								return
							}
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLDeleteFHIRComposition(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}
	composition, _, err := createTestFHIRComposition(ctx, encounterID)
	if err != nil {
		t.Errorf("can't create test composition: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query - case that exists",
			args: args{
				query: map[string]interface{}{
					"query": `mutation DeleteComposition($id: ID!) {
						deleteFHIRComposition(id: $id)
					  }`,
					"variables": map[string]interface{}{
						"id": *composition.ID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "valid composition delete",
			args: args{
				query: map[string]interface{}{
					"query": `mutation deleteFHIRComposition($id: ID!) {
						deleteFHIRComposition(
							id: $id
						)
					  }`,
					"variables": map[string]interface{}{
						"id": ksuid.New().String(),
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid format",
			args: args{
				query: map[string]interface{}{
					"query": `mutation deleteFHIRComposition($id: ID!) {
						deleteFHIRComposition
						// bad format
					  }`,
					"variables": map[string]interface{}{
						"id": ksuid.New().String(),
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					if key == "data" {
						respMap, present := nestedMap["deleteFHIRComposition"]
						if !present {
							t.Errorf("can't find delete response")
							return
						}

						deleted, ok := nestedMap["deleteFHIRComposition"].(bool)
						if !ok {
							t.Errorf("cannot cast key value of %v to type map[string]interface{}", respMap)
							return
						}

						if !deleted {
							t.Errorf("expected the composition to have been successfully deleted, it wasn't")
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQlSearchFHIRComposition(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}
	_, patient, err := createTestFHIRComposition(ctx, encounterID)
	if err != nil {
		t.Errorf("can't create test composition: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query SearchCompositions($params: Map!) {
						searchFHIRComposition(params: $params) {
						  edges {
							node {
							  ID
							  Status
							  Type {
								Text
								Coding {
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							  Category {
								Text
								Coding {
								  System
								  Code
								  Display
								  UserSelected
								}
							  }
							  Author {
								Reference
								Display
							  }
							  Title
							  Section {
								Title
								Text {
								  Status
								  Div
								}
							  }
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"params": map[string]interface{}{
							"status":    domain.CompositionStatusEnumPreliminary.String(),
							"patient":   fmt.Sprintf("Patient/%s", *patient.ID),
							"encounter": fmt.Sprintf("Encounter/%s", encounterID),
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", data[key])
						return
					}

					if key == "data" {
						respMap, present := nestedMap["searchFHIRComposition"]
						if !present {
							t.Errorf("can't find delete response")
							return
						}

						_, ok := respMap.(map[string]interface{})
						if !ok {
							t.Errorf("can't cast respMap %v to map[string]interface{}", respMap)
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQLAllergySummary(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("could not get patient: %v", err)
		return
	}

	if patient.ID == nil {
		t.Errorf("nil patient ID")
		return
	}

	patientID := *patient.ID

	_, err = createTestAllergy(ctx, patient, encounterID)
	if err != nil {
		t.Errorf("error creating a test condition: %v", err)
		return
	}

	// we have intermittent CI failures that could be related to replication
	// lag or latency issues on the backing data store (unproven).
	// If that is the case, then this sleep would reduce the failure rate.
	time.Sleep(time.Second * 5)

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query AllergySummary($patientID: String!) {
						allergySummary(patientID: $patientID)
					}`,
					"variables": map[string]interface{}{
						"patientID": patientID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			if tt.wantErr {
				data := map[string]interface{}{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned: %s", string(dataResponse))
					return
				}
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}
			if !tt.wantErr {
				data := map[string]map[string][]string{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned: %s", string(dataResponse))
					return
				}
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %s", errMsg)
					return
				}

				for key := range data {
					nestedMap, present := data[key]
					if !present {
						t.Errorf("key %s not found in %v", key, data)
						return
					}
					expected := map[string][]string{
						"allergySummary": {
							"Panadol Extra",
						},
					}
					if !reflect.DeepEqual(expected, nestedMap) {
						t.Errorf("expected %v got %v", expected, nestedMap)
						return
					}
				}

			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLDeleteFHIRPatient(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query - case that exists",
			args: args{
				query: map[string]interface{}{
					"query": `mutation DeletePatient($id: ID!) {
						deleteFHIRPatient(id: $id)
					  }`,
					"variables": map[string]interface{}{
						"id": *patient.ID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid format",
			args: args{
				query: map[string]interface{}{
					"query": `mutation deleteFHIRPatient($id: ID!) {
						deleteFHIRPatient
						// bad format
					  }`,
					"variables": map[string]interface{}{
						"id": ksuid.New().String(),
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					if key == "data" {
						respMap, present := nestedMap["deleteFHIRPatient"]
						if !present {
							t.Errorf("can't find delete response")
							return
						}

						deleted, ok := nestedMap["deleteFHIRPatient"].(bool)
						if !ok {
							t.Errorf("cannot cast key value of %v to type map[string]interface{}", respMap)
							return
						}

						if !deleted {
							t.Errorf("expected the composition to have been successfully deleted, it wasn't")
							return
						}
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
	s := testInteractor
	payload := &domain.PhoneNumberPayload{}
	s.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
}

func TestGraphQLListConcepts(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid concept list query",
			args: args{
				query: map[string]interface{}{
					"query": `
						query ListConcepts(
							$org: String!, $source: String!, $verbose: Boolean!, 
							$q: String, $conceptClass: String, $includeRetired: Boolean,
							$includeMappings: Boolean, $includeInverseMappings: Boolean,
							$sortAsc: String, $locale: String
						) {
							listConcepts(
							org: $org, 
							source:$source,
							verbose: $verbose,
							q: $q,
							conceptClass: $conceptClass,
							includeRetired: $includeRetired,
							includeMappings: $includeMappings,
							includeInverseMappings: $includeInverseMappings,
							sortAsc: $sortAsc,
							locale: $locale
							)
						}
						`,
					"variables": map[string]interface{}{
						"org":                    "CIEL",
						"source":                 "CIEL",
						"q":                      "cold",
						"conceptClass":           "Diagnosis",
						"verbose":                false,
						"includeRetired":         false,
						"includeMappings":        false,
						"includeInverseMappings": false,
						"sortAsc":                "bestMatch",
						"locale":                 "en",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid concept list query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRMedicationStatement(t *testing.T) {
	ctx := firebasetools.GetAuthenticatedContext(t)
	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %v", err)
		return
	}

	_, patient, _, err := getTestEncounterID(
		ctx, interserviceclient.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("can't create test encounter: %v", err)
		return
	}

	type args struct {
		query map[string]interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query",
			args: args{
				query: map[string]interface{}{
					"query": `query SearchMedicationStatement($params: Map!) {
						searchFHIRMedicationStatement(params: $params) {
						  edges {
							cursor
							node {
							  ID
							  Identifier {
								ID
								Use
								Type {
								  ID
								  Coding {
									ID
									System
									Version
								  }
								  Text
								}
							  }
							  Status
							  Category {
								ID
								Coding {
								  ID
								  System
								  Version
								  Display
								  UserSelected
								  Code
								}
								Text
							  }
							  Subject {
								ID
								Reference
								Type
								Display
							  }
							  MedicationCodeableConcept {
								ID
								Text
							  }
							  MedicationReference {
								ID
								Code {
								  ID
								  Text
								}
							  }
							  EffectiveDateTime
							  DateAsserted
							}
						  }
						}
					  }`,
					"variables": map[string]interface{}{
						"params": map[string]interface{}{
							"patient": fmt.Sprintf("Patient/%s", *patient.ID),
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid query",
			args: args{
				query: map[string]interface{}{
					"query":     `bad format query`,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}
			if tt.wantErr {
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				errMsg, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got: %v", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
}
