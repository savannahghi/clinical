package graph_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/clinical/graph"
)

const (
	testHTTPClientTimeout = 180
	testProviderCode      = "123"
	dateFormat            = "2006-01-02"
	testProviderPhone     = "+254721000111"
	testProviderUID       = "0b1fcd62-46df-4cbc-9096-7849859dcd76"
)

// these are set up once in TestMain and used by all the acceptance tests in
// this package
var srv *http.Server
var baseURL string
var serverErr error

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

func TestMain(m *testing.M) {
	// setup
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("CLOUD_HEALTH_PUBSUB_TOPIC", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_DATASET_ID", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_FHIRSTORE_ID", "healthcloud-bewell-fhir-staging")

	ctx := context.Background()
	srv, baseURL, serverErr = base.StartTestServer(
		ctx,
		graph.PrepareServer,
		graph.ClinicalAllowedOrigins,
	) // set the globals
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

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

func TestRouter(t *testing.T) {
	router, err := graph.Router()
	if err != nil {
		t.Errorf("can't initialize router: %v", err)
		return
	}

	if router == nil {
		t.Errorf("nil router")
		return
	}
}

func TestHealthStatusCheck(t *testing.T) {
	client := http.DefaultClient

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful health check",
			args: args{
				url: fmt.Sprintf(
					"%s/health",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

func TestGraphQLRegisterPatient(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	simplePatientRegInput, _, err := getTestSimplePatientRegistration()
	if err != nil {
		t.Errorf("can't genereate simple patient reg inpit: %v", err)
		return
	}

	patientRegInput, err := base.StructToMap(simplePatientRegInput)
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
					"query": `
mutation SimplePatientRegistration($input: SimplePatientRegistrationInput!) {
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
				if ok {
					t.Errorf("error not expected got error: %s", errors)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQFindPatientsByMSISDN(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	_, msisdn, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("error in getting test patient: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLFindPatients(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQGetPatient(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("can't get test patient: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLStartEpisodeByOTP(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	phoneNumber, otp, err := getTestVerifiedPhoneandOTP()
	if err != nil {
		t.Errorf("unable to get verified phone number and OTP")
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
							"msisdn":       phoneNumber,
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
					t.Errorf("error not expected got: %w", errMsg)
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLStartEpisodeByBreakGlass(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	validPhone, otp, err := getTestVerifiedPhoneandOTP()
	if err != nil {
		t.Errorf("unable to get verified phone number and OTP")
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
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
							"msisdn":          validPhone,
							"fullAccess":      false,
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
					t.Errorf("error not expected got: %w", errMsg)
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLUpgradeEpisode(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	phoneNumber, otp, err := getTestVerifiedPhoneandOTP()
	if err != nil {
		t.Errorf("unable to get verified phone number and OTP")
		return
	}

	episode, _, err := getTestEpisodeOfCare(
		ctx,
		base.TestUserPhoneNumber,
		false, testProviderCode,
	)
	if err != nil {
		t.Errorf("can't create test episode: %w", err)
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
							"msisdn":    phoneNumber,
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
					t.Errorf("error not expected got: %w", errMsg)
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
				t.Errorf("Bad status reponse returned")
				return
			}
		})
	}
}

func TestGraphQLEndEpisode(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	episode, _, err := getTestEpisodeOfCare(
		ctx,
		base.TestUserPhoneNumber,
		false, testProviderCode,
	)
	if err != nil {
		t.Errorf("can't create test episode: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
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
				t.Errorf("Bad status reponse returned")
				return
			}
		})
	}
}

func TestGraphQLStartEncounter(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}
	episode, _, _, err := getTestEncounterID(
		ctx, base.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLEndEncounter(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	_, _, encounterID, err := getTestEncounterID(
		ctx, base.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQOpenEpisodes(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	_, patient, err := getTestEpisodeOfCare(
		ctx,
		base.TestUserPhoneNumber,
		false, testProviderCode,
	)
	if err != nil {
		t.Errorf("can't create test episode: %w", err)
		return
	}
	if patient.ID == nil {
		t.Errorf("nil patient ID")
		return
	}

	patientID := *patient.ID
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIREncounter(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	episode, _, _, err := getTestEncounterID(
		ctx, base.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %w", err)
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
					t.Errorf("error not expected, got error: %s", errorMessage)
					return
				}

				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}

					log.Printf("response: \n%s\n", nestedMap)

					_, ok = nestedMap["searchFHIREncounter"].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast nested map key value of %v to type map[string]interface{}", key)
						return
					}
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphqlOpenOrganizationEpisodes(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	_, _, err = getTestEpisodeOfCare(
		ctx,
		base.TestUserPhoneNumber,
		false, testProviderCode,
	)
	if err != nil {
		t.Errorf("can't create test episode: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLAddNextOfKin(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	phoneNumber, otp, err := getTestVerifiedPhoneandOTP()
	if err != nil {
		t.Errorf("unable to get verified phone number and OTP")
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
								"msisdn":             phoneNumber,
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
					t.Errorf("error not expected got: %w", errMsg)
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLUpdatePatient(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("can't get test patient: %w", err)
		return
	}

	newPatientInputData, _, err := getTestSimplePatientRegistration()
	if err != nil {
		t.Errorf("can't genereate simple patient reg inpit: %v", err)
		return
	}

	patientInputWithUpdatedData, err := base.StructToMap(newPatientInputData)
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}
		})
	}
}

func TestGraphQLAddNHIF(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	patient, _, err := getTestPatient(ctx)
	if err != nil {
		t.Errorf("error in getting test patient: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLCreateUpdatePatientExtraInformation(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
								"email":              base.GenerateRandomEmail(),
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
					t.Errorf("error not expected got: %w", errMsg)
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQVisitSummary(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @timeline @mathenge Test visit summary
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQPatientTimelineWithCount(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @timeline @mashaa Test patient timeline with count
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQProblemSummary(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @timeline @sala Test problem summary
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLCreateFHIRMedicationRequest(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, base.TestUserPhoneNumber, false, testProviderCode)
	if err != nil {
		t.Errorf("error creating test encounter ID: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
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

							id, prs := resource["ID"]
							if !prs {
								t.Errorf("ID not present in medication request resource")
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLUpdateFHIRMedicationRequest(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @medication @maluki Test update medication request
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLDeleteFHIRMedicationRequest(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @medication @maluki Test delete medication request
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRMedicationRequest(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @medication @maluki Test search FHIR medication request
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLCreateFHIRAllergyIntolerance(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @allergy @mashaa Test create allergy intolerance
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLUpdateFHIRAllergyIntolerance(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @allergy @mashaa Test update FHIR allergy intolerance
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRAllergyIntolerance(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @allergy @mashaa Test search FHIR allergy intolerance
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLCreateFHIRCondition(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, base.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}
		})
	}
}

func TestGraphQUpdateFHIRCondition(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @condition @sala Test update FHIR condition
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRCondition(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @condition @sala Test search FHIR condition
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQCreateFHIRServiceRequest(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @servicerequest @mathenge Test create FHIR service request
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQDeleteFHIRServiceRequest(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @servicerequest @mathenge Test delete FHIR service request
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRServiceRequest(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @servicerequest @mathenge Test search FHIR service request
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQCreateFHIRObservation(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @observation @ngure Test create FHIR observation
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQSearchFHIRObservation(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @observation @ngure Test search FHIR observation
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}
		})
	}
}

func TestGraphQCreateFHIRComposition(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
		return
	}

	_, patient, encounterID, err := getTestEncounterID(
		ctx, base.TestUserPhoneNumber, true, testProviderCode)
	if err != nil {
		t.Errorf("unable to generate test encounter ID: %w", err)
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQUpdateFHIRComposition(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @composition @ngure Test update FHIR composition
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLDeleteFHIRComposition(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @composition @ngure Delete a composition that exists (additional test case)
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQlSearchFHIRComposition(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
		// TODO @composition @ngure Test search FHIR composition
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestGraphQLListConcepts(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	if ctx == nil {
		t.Errorf("nil context")
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting GraphQL headers: %w", err)
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
					t.Errorf("error not expected got: %w", errMsg)
					return
				}
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}
