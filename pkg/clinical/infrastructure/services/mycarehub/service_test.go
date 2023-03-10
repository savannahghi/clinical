package mycarehub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	extMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

func TestServiceMyCareHubImpl_UserProfile(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: get user profile",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to make request",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case: non OK status returned",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case: non json body",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case: user id not in body",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := extMock.NewFakeBaseExtensionMock()
			fakeISC := extMock.NewFakeISCExtensionMock()
			s := ServiceMyCareHubImpl{
				MyCareHubClient: fakeISC,
				baseExt:         fakeExt,
			}

			if tt.name == "Happy Case: get user profile" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					UUID := gofakeit.UUID()
					requestBody := domain.User{
						ID: &UUID,
					}
					jsonBody, err := json.Marshal(requestBody)
					if err != nil {
						return nil, err
					}
					return &http.Response{
						Status:     "200",
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(jsonBody)),
					}, nil
				}
			}
			if tt.name == "Sad Case: failed to make request" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case: non OK status returned" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					jsonBody, err := json.Marshal(map[string]interface{}{"error": "error"})
					if err != nil {
						return nil, err
					}
					return &http.Response{
						Status:     "400",
						StatusCode: 400,
						Body:       io.NopCloser(bytes.NewReader(jsonBody)),
					}, nil
				}
			}

			if tt.name == "Sad Case: non json body" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					jsonBody, err := json.Marshal("not json")
					if err != nil {
						return nil, err
					}
					return &http.Response{
						Status:     "200",
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(jsonBody)),
					}, nil
				}
			}

			if tt.name == "Sad Case: user id not in body" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					requestBody := domain.User{}
					jsonBody, err := json.Marshal(requestBody)
					if err != nil {
						return nil, err
					}
					return &http.Response{
						Status:     "200",
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(jsonBody)),
					}, nil
				}
			}
			got, err := s.UserProfile(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceMyCareHubImpl.UserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a value, got %v", got)
				return
			}
		})
	}
}

func TestServiceMyCareHubImpl_AddFHIRIDToPatientProfile(t *testing.T) {
	type args struct {
		ctx      context.Context
		fhirID   string
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: add FHIR patient to profile",
			args: args{
				ctx:      context.Background(),
				fhirID:   gofakeit.UUID(),
				clientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to make request",
			args: args{
				ctx:      context.Background(),
				fhirID:   gofakeit.UUID(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case: non OK status returned",
			args: args{
				ctx:      context.Background(),
				fhirID:   gofakeit.UUID(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := extMock.NewFakeBaseExtensionMock()
			fakeISC := extMock.NewFakeISCExtensionMock()
			s := ServiceMyCareHubImpl{
				MyCareHubClient: fakeISC,
				baseExt:         fakeExt,
			}

			if tt.name == "Happy case: add FHIR patient to profile" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					UUID := gofakeit.UUID()
					requestBody := domain.FHIRPatient{
						ID: &UUID,
					}
					jsonBody, err := json.Marshal(requestBody)
					if err != nil {
						return nil, err
					}
					return &http.Response{
						Status:     "200",
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(jsonBody)),
					}, nil
				}
			}
			if tt.name == "Sad Case: failed to make request" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case: non OK status returned" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					jsonBody, err := json.Marshal(map[string]interface{}{"error": "error"})
					if err != nil {
						return nil, err
					}
					return &http.Response{
						Status:     "400",
						StatusCode: 400,
						Body:       io.NopCloser(bytes.NewReader(jsonBody)),
					}, nil
				}
			}

			if err := s.AddFHIRIDToPatientProfile(tt.args.ctx, tt.args.fhirID, tt.args.clientID); (err != nil) != tt.wantErr {
				t.Errorf("ServiceMyCareHubImpl.AddFHIRIDToPatientProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceMyCareHubImpl_AddFHIRIDToFacility(t *testing.T) {
	type args struct {
		ctx        context.Context
		fhirID     string
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: add FHIR id to facility",
			args: args{
				ctx:        context.Background(),
				fhirID:     gofakeit.UUID(),
				facilityID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to make request",
			args: args{
				ctx:        context.Background(),
				fhirID:     gofakeit.UUID(),
				facilityID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case: non OK status returned",
			args: args{
				ctx:        context.Background(),
				fhirID:     gofakeit.UUID(),
				facilityID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := extMock.NewFakeBaseExtensionMock()
			fakeISC := extMock.NewFakeISCExtensionMock()
			s := ServiceMyCareHubImpl{
				MyCareHubClient: fakeISC,
				baseExt:         fakeExt,
			}

			if tt.name == "Happy Case: add FHIR id to facility" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					UUID := gofakeit.UUID()
					requestBody := domain.User{
						ID: &UUID,
					}
					jsonBody, err := json.Marshal(requestBody)
					if err != nil {
						return nil, err
					}
					return &http.Response{
						Status:     "200",
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(jsonBody)),
					}, nil
				}
			}
			if tt.name == "Sad Case: failed to make request" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case: non OK status returned" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					jsonBody, err := json.Marshal(map[string]interface{}{"error": "error"})
					if err != nil {
						return nil, err
					}
					return &http.Response{
						Status:     "400",
						StatusCode: 400,
						Body:       io.NopCloser(bytes.NewReader(jsonBody)),
					}, nil
				}
			}
			if err := s.AddFHIRIDToFacility(tt.args.ctx, tt.args.fhirID, tt.args.facilityID); (err != nil) != tt.wantErr {
				t.Errorf("ServiceMyCareHubImpl.AddFHIRIDToFacility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
