package clinical_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	fakeExtMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fakeFHIRMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/mock"
	fakeAdvantageMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/advantage/mock"
	fakeOCLMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab/mock"
	fakePubSubMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/pubsub/mock"
	fakeUploadMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/upload/mock"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
)

func TestUseCasesClinicalImpl_GetTenantMetaTags(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get tenant org from context",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "sad case: missing tenant org in context",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving organisation",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
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
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "sad case: missing tenant org in context" {
				fakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to to get identifiers")
				}
			}

			if tt.name == "sad case: error retrieving organisation" {
				fakeFHIR.MockGetFHIROrganizationFn = func(ctx context.Context, organisationID string) (*domain.FHIROrganizationRelayPayload, error) {
					return nil, fmt.Errorf("failed to find organization")
				}
			}

			got, err := c.GetTenantMetaTags(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.GetTenantMetaTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected result not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_CheckPatientExistenceUsingPhoneNumber(t *testing.T) {
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
			name: "happy case: nil inputs",
			args: args{
				ctx: context.Background(),
				input: domain.SimplePatientRegistrationInput{
					PhoneNumbers: []*domain.PhoneNumberInput{
						{
							Msisdn:             gofakeit.Phone(),
							VerificationCode:   "1234",
							IsUssd:             false,
							CommunicationOptIn: false,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: contacts to contact point",
			args: args{
				ctx: context.Background(),
				input: domain.SimplePatientRegistrationInput{
					PhoneNumbers: []*domain.PhoneNumberInput{
						{
							Msisdn:             gofakeit.Phone(),
							VerificationCode:   "1234",
							IsUssd:             false,
							CommunicationOptIn: false,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: missing FHIR patient",
			args: args{
				ctx: context.Background(),
				input: domain.SimplePatientRegistrationInput{
					PhoneNumbers: []*domain.PhoneNumberInput{
						{
							Msisdn:             gofakeit.Phone(),
							VerificationCode:   "1234",
							IsUssd:             false,
							CommunicationOptIn: false,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: missing tenant org in context",
			args: args{
				ctx: context.Background(),
				input: domain.SimplePatientRegistrationInput{
					PhoneNumbers: []*domain.PhoneNumberInput{
						{
							Msisdn:             gofakeit.Phone(),
							VerificationCode:   "1234",
							IsUssd:             false,
							CommunicationOptIn: false,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid phone",
			args: args{
				ctx: context.Background(),
				input: domain.SimplePatientRegistrationInput{
					PhoneNumbers: []*domain.PhoneNumberInput{
						{
							Msisdn:             "0722",
							VerificationCode:   "1234",
							IsUssd:             false,
							CommunicationOptIn: false,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FakeExt := fakeExtMock.NewFakeBaseExtensionMock()
			Fakefhir := fakeFHIRMock.NewFHIRMock()
			FakeOCL := fakeOCLMock.NewFakeOCLMock()
			fakeUpload := fakeUploadMock.NewFakeUploadMock()
			fakePubSub := fakePubSubMock.NewPubSubServiceMock()
			fakeAdvantage := fakeAdvantageMock.NewFakeAdvantageMock()

			infra := infrastructure.NewInfrastructureInteractor(FakeExt, Fakefhir, FakeOCL, fakeUpload, fakePubSub, fakeAdvantage)
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			if tt.name == "sad case: missing tenant org in context" {
				FakeExt.MockGetTenantIdentifiersFn = func(ctx context.Context) (*dto.TenantIdentifiers, error) {
					return nil, fmt.Errorf("failed to to get identifiers")
				}
			}

			if tt.name == "sad case: missing FHIR patient" {
				Fakefhir.MockSearchFHIRPatientFn = func(ctx context.Context, searchParams string, tenant dto.TenantIdentifiers, pagination dto.Pagination) (*domain.PatientConnection, error) {
					return nil, fmt.Errorf("failed to to get fhir patient")
				}
			}

			got, err := c.CheckPatientExistenceUsingPhoneNumber(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.CheckPatientExistenceUsingPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != false {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_ContactsToContactPointInput(t *testing.T) {
	email := gofakeit.Email()
	invalidEmail := "gofakeit.Email()"
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
			name: "happy case: nil inputs",
			args: args{
				ctx:    context.Background(),
				phones: nil,
				emails: nil,
			},
			wantErr: false,
		},
		{
			name: "happy case: contacts to contact point",
			args: args{
				ctx: context.Background(),
				phones: []*domain.PhoneNumberInput{
					{
						Msisdn:             gofakeit.Phone(),
						VerificationCode:   "1234",
						IsUssd:             false,
						CommunicationOptIn: false,
					},
				},
				emails: []*domain.EmailInput{
					{
						Email:              &email,
						CommunicationOptIn: false,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid phone",
			args: args{
				ctx: context.Background(),
				phones: []*domain.PhoneNumberInput{
					{
						Msisdn:             "0722",
						VerificationCode:   "1234",
						IsUssd:             false,
						CommunicationOptIn: false,
					},
				},
				emails: []*domain.EmailInput{
					{
						Email:              &email,
						CommunicationOptIn: false,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid email",
			args: args{
				ctx: context.Background(),
				phones: []*domain.PhoneNumberInput{
					{
						Msisdn:             gofakeit.Phone(),
						VerificationCode:   "1234",
						IsUssd:             false,
						CommunicationOptIn: false,
					},
				},
				emails: []*domain.EmailInput{
					{
						Email:              &invalidEmail,
						CommunicationOptIn: false,
					},
				},
			},
			wantErr: true,
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
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			got, err := c.ContactsToContactPointInput(tt.args.ctx, tt.args.phones, tt.args.emails)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.ContactsToContactPointInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesClinicalImpl_SimplePatientRegistrationInputToPatientInput(t *testing.T) {
	email := gofakeit.Email()
	invalidEmail := "invalid"
	address := gofakeit.BS()
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
			name: "happy case: fhir patient input",
			args: args{
				ctx: context.Background(),
				input: domain.SimplePatientRegistrationInput{
					ID: gofakeit.UUID(),
					Names: []*domain.NameInput{
						{
							FirstName: gofakeit.Name(),
							LastName:  gofakeit.Name(),
						},
					},
					IdentificationDocuments: []*domain.IdentificationDocument{
						{
							DocumentType:   domain.IDDocumentTypePassport,
							DocumentNumber: gofakeit.SSN(),
						},
					},
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 10,
						Day:   10,
					},
					PhoneNumbers: []*domain.PhoneNumberInput{},
					Photos: []*domain.PhotoInput{
						{
							PhotoContentType: enumutils.ContentTypeJpg,
							PhotoBase64data:  "qweqwdwedwed",
							PhotoFilename:    "test",
						},
					},
					Emails: []*domain.EmailInput{
						{
							Email:              &email,
							CommunicationOptIn: false,
						},
					},
					PhysicalAddresses: []*domain.PhysicalAddress{
						{
							MapsCode:        "123",
							PhysicalAddress: &address,
						},
					},
					PostalAddresses: []*domain.PostalAddress{
						{
							PostalAddress: &address,
							PostalCode:    "1234",
						},
					},
					Gender:        "",
					Active:        true,
					MaritalStatus: "",
					Languages:     []enumutils.Language{"en"},
					ReplicateUSSD: false,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid email",
			args: args{
				ctx: context.Background(),
				input: domain.SimplePatientRegistrationInput{
					ID: gofakeit.UUID(),
					Names: []*domain.NameInput{
						{
							FirstName: gofakeit.Name(),
							LastName:  gofakeit.Name(),
						},
					},
					IdentificationDocuments: []*domain.IdentificationDocument{
						{
							DocumentType:   domain.IDDocumentTypePassport,
							DocumentNumber: gofakeit.SSN(),
						},
					},
					BirthDate: &scalarutils.Date{
						Year:  2000,
						Month: 10,
						Day:   10,
					},
					PhoneNumbers: []*domain.PhoneNumberInput{},
					Photos: []*domain.PhotoInput{
						{
							PhotoContentType: enumutils.ContentTypeJpg,
							PhotoBase64data:  "qweqwdwedwed",
							PhotoFilename:    "test",
						},
					},
					Emails: []*domain.EmailInput{
						{
							Email:              &invalidEmail,
							CommunicationOptIn: false,
						},
					},
					PhysicalAddresses: []*domain.PhysicalAddress{
						{
							MapsCode:        "123",
							PhysicalAddress: &address,
						},
					},
					PostalAddresses: []*domain.PostalAddress{
						{
							PostalAddress: &address,
							PostalCode:    "1234",
						},
					},
					Active:        true,
					Languages:     []enumutils.Language{"en"},
					ReplicateUSSD: false,
				},
			},
			wantErr: true,
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
			c := clinicalUsecase.NewUseCasesClinicalImpl(infra)

			got, err := c.SimplePatientRegistrationInputToPatientInput(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClinicalImpl.SimplePatientRegistrationInputToPatientInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected result to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected result not to be nil for %v", tt.name)
				return
			}
		})
	}
}
