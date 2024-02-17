package rest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/savannahghi/authutils"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeAdvantageMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/advantage/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	"github.com/savannahghi/clinical/pkg/clinical/presentation"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

var (
	authServerEndpoint = serverutils.MustGetEnvVar("AUTHSERVER_ENDPOINT")
	clientID           = serverutils.MustGetEnvVar("CLIENT_ID")
	clientSecret       = serverutils.MustGetEnvVar("CLIENT_SECRET")
	username           = serverutils.MustGetEnvVar("AUTH_USERNAME")
	password           = serverutils.MustGetEnvVar("AUTH_PASSWORD")
	grantType          = serverutils.MustGetEnvVar("GRANT_TYPE")
)

var (
	authclient, _ = authutils.NewClient(
		authutils.Config{
			AuthServerEndpoint: authServerEndpoint,
			ClientID:           clientID,
			ClientSecret:       clientSecret,
			GrantType:          grantType,
			Username:           username,
			Password:           password,
		},
	)
)

var testMemoryStore = persist.NewMemoryStore(60 * time.Minute)

func TestPresentationHandlersImpl_ReceivePubSubPushMessage(t *testing.T) {
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
		headers    map[string]string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "happy case: publish create organization message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "sad case: publish create organization message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "happy case: publish create medication message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "sad case: publish create medication message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "happy case: publish create patient message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "sad case: publish create patient message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "happy case: publish create vitals message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "sad case: publish create vitals message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "happy case: publish create allergy message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "sad case: publish create allergy message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "happy case: publish create results message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "happy case: publish tenant message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "sad case: publish create results message",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "happy case: publish patient segmentation",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "sad case: unable to publish patient segmentation",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "sad case: verify pubsub request fails",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "sad case: topic name error",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "sad case: unknown topic",
			args: args{
				url:        "/pubsub",
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			fakeFHIR := fakeFHIRMock.NewFHIRMock()
			fakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()
			infra := infrastructure.NewInfrastructureInteractor(fakeExt, fakeFHIR, fakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			usecases := usecases.NewUsecasesInteractor(infra)

			if tt.name == "happy case: publish create patient message" {
				msg := dto.PatientPubSubMessage{
					UserID:         gofakeit.UUID(),
					ClientID:       gofakeit.UUID(),
					Name:           gofakeit.Name(),
					DateOfBirth:    time.Now(),
					Gender:         "male",
					Active:         true,
					PhoneNumber:    gofakeit.Phone(),
					OrganizationID: gofakeit.UUID(),
					FacilityID:     gofakeit.UUID(),
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName), nil
				}
			}

			if tt.name == "sad case: publish create patient message" {
				msg := dto.PatientPubSubMessage{
					UserID:         gofakeit.UUID(),
					ClientID:       gofakeit.UUID(),
					Name:           gofakeit.Name(),
					DateOfBirth:    time.Now(),
					Gender:         "male",
					Active:         true,
					PhoneNumber:    gofakeit.Phone(),
					OrganizationID: gofakeit.UUID(),
					FacilityID:     gofakeit.UUID(),
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName), nil
				}

				fakeFHIR.MockCreateFHIRPatientFn = func(ctx context.Context, input domain.FHIRPatientInput) (*domain.PatientPayload, error) {
					return nil, fmt.Errorf("failed to create patient")
				}
			}

			if tt.name == "happy case: publish create organization message" {
				id := gofakeit.UUID()
				msg := dto.FacilityPubSubMessage{
					ID:          &id,
					Name:        gofakeit.Company(),
					Code:        100,
					Phone:       gofakeit.Phone(),
					Active:      true,
					County:      gofakeit.Country(),
					Description: gofakeit.Company(),
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName), nil
				}
			}

			if tt.name == "sad case: publish create organization message" {
				id := gofakeit.UUID()
				msg := dto.FacilityPubSubMessage{
					ID:          &id,
					Name:        gofakeit.Company(),
					Code:        100,
					Phone:       gofakeit.Phone(),
					Active:      true,
					County:      gofakeit.Country(),
					Description: gofakeit.Company(),
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName), nil
				}

				fakeFHIR.MockCreateFHIROrganizationFn = func(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to create organization")
				}
			}

			if tt.name == "happy case: publish create vitals message" {
				concept := "1234"
				msg := dto.VitalSignPubSubMessage{
					PatientID:      gofakeit.UUID(),
					OrganizationID: gofakeit.UUID(),
					Name:           gofakeit.Name(),
					ConceptID:      &concept,
					Value:          "",
					Date:           time.Now(),
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName), nil
				}
			}

			if tt.name == "sad case: publish create vitals message" {
				concept := "1234"
				msg := dto.VitalSignPubSubMessage{
					PatientID:      gofakeit.UUID(),
					OrganizationID: gofakeit.UUID(),
					Name:           gofakeit.Name(),
					ConceptID:      &concept,
					Value:          "",
					Date:           time.Now(),
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName), nil
				}

				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			if tt.name == "happy case: publish create allergy message" {
				concept := "1234"
				msg := dto.PatientAllergyPubSubMessage{
					PatientID:      gofakeit.UUID(),
					OrganizationID: gofakeit.UUID(),
					Name:           gofakeit.Name(),
					ConceptID:      &concept,
					Date:           time.Now(),
					Reaction:       dto.AllergyReaction{},
					Severity:       dto.AllergySeverity{},
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName), nil
				}
			}

			if tt.name == "sad case: publish create allergy message" {
				concept := "1234"
				msg := dto.PatientAllergyPubSubMessage{
					PatientID:      gofakeit.UUID(),
					OrganizationID: gofakeit.UUID(),
					Name:           gofakeit.Name(),
					ConceptID:      &concept,
					Date:           time.Now(),
					Reaction:       dto.AllergyReaction{},
					Severity:       dto.AllergySeverity{},
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName), nil
				}

				fakeFHIR.MockCreateFHIRAllergyIntoleranceFn = func(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
					return nil, fmt.Errorf("failed to create allergy")
				}
			}

			if tt.name == "happy case: publish create medication message" {
				concept := "1234"
				msg := dto.MedicationPubSubMessage{
					PatientID:      gofakeit.UUID(),
					OrganizationID: gofakeit.UUID(),
					Name:           gofakeit.Name(),
					ConceptID:      &concept,
					Date:           time.Now(),
					Value:          "",
					Drug: &dto.MedicationDrug{
						ConceptID: &concept,
					},
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName), nil
				}
			}

			if tt.name == "sad case: publish create medication message" {
				concept := "1234"
				msg := dto.MedicationPubSubMessage{
					PatientID:      gofakeit.UUID(),
					OrganizationID: gofakeit.UUID(),
					Name:           gofakeit.Name(),
					ConceptID:      &concept,
					Date:           time.Now(),
					Value:          "",
					Drug: &dto.MedicationDrug{
						ConceptID: &concept,
					},
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName), nil
				}

				fakeFHIR.MockCreateFHIRMedicationStatementFn = func(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error) {
					return nil, fmt.Errorf("failed to create observation")
				}
			}

			if tt.name == "happy case: publish create results message" {
				concept := "1234"
				msg := dto.PatientTestResultPubSubMessage{
					PatientID:      gofakeit.UUID(),
					OrganizationID: gofakeit.UUID(),
					Name:           gofakeit.Name(),
					ConceptID:      &concept,
					Date:           time.Now(),
					Result:         dto.TestResult{},
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName), nil
				}
			}

			if tt.name == "happy case: publish tenant message" {
				message := dto.OrganizationInput{
					Name:        gofakeit.BeerName(),
					PhoneNumber: interserviceclient.TestUserPhoneNumber,
					Identifiers: []dto.OrganizationIdentifier{
						{
							Type:  dto.OrganizationIdentifierType("MCHProgram"),
							Value: gofakeit.UUID(),
						},
					},
				}
				data, _ := json.Marshal(message)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.TenantTopicName, common.ClinicalServiceName), nil
				}
			}

			if tt.name == "sad case: publish create results message" {
				concept := "1234"
				msg := dto.PatientTestResultPubSubMessage{
					PatientID:      gofakeit.UUID(),
					OrganizationID: gofakeit.UUID(),
					Name:           gofakeit.Name(),
					ConceptID:      &concept,
					Date:           time.Now(),
					Result:         dto.TestResult{},
				}
				data, _ := json.Marshal(msg)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName), nil
				}

				fakeFHIR.MockCreateFHIRObservationFn = func(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservation, error) {
					return nil, fmt.Errorf("failed to crate observation")
				}
			}

			if tt.name == "sad case: verify pubsub request fails" {
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return nil, fmt.Errorf("failed to verify")
				}
			}

			if tt.name == "sad case: topic name error" {
				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return "", fmt.Errorf("cant get pubsub topic")
				}
			}

			if tt.name == "happy case: publish patient segmentation" {
				message := dto.SegmentationPayload{
					ClinicalID:   gofakeit.UUID(),
					SegmentLabel: "LOW RISK",
				}

				data, _ := json.Marshal(message)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.SegmentationTopicName, common.ClinicalServiceName), nil
				}
			}
			if tt.name == "sad case: unable to publish patient segmentation" {
				message := dto.SegmentationPayload{
					ClinicalID:   gofakeit.UUID(),
					SegmentLabel: "LOW RISK",
				}

				data, _ := json.Marshal(message)
				fakeExt.MockVerifyPubSubJWTAndDecodePayloadFn = func(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error) {
					return &pubsubtools.PubSubPayload{
						Message: pubsubtools.PubSubMessage{
							Data: data,
						},
					}, nil
				}

				fakeExt.MockGetPubSubTopicFn = func(m *pubsubtools.PubSubPayload) (string, error) {
					return utils.AddPubSubNamespace(common.SegmentationTopicName, common.ClinicalServiceName), nil
				}

				fakeAdvantage.MockSegmentPatient = func(ctx context.Context, payload dto.SegmentationPayload) error {
					return fmt.Errorf("failed to segment patient")
				}
			}

			w := httptest.NewRecorder()
			ctx, engine := gin.CreateTestContext(w)

			req, err := http.NewRequestWithContext(
				ctx,
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			for k, v := range tt.args.headers {
				req.Header.Add(k, v)
			}

			presentation.SetupRoutes(engine, testMemoryStore, authclient, usecases, infra)
			engine.ServeHTTP(w, req)

			resp := w.Result()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}

			dataResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

		})
	}
}
