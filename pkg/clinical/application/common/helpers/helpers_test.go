package helpers

import (
	"context"
	"encoding/base64"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/enums"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/enumutils"
	"github.com/stretchr/testify/assert"
)

func TestComposeOneHealthEpisodeOfCare(t *testing.T) {
	type args struct {
		validPhone     string
		fullAccess     bool
		organizationID string
		providerCode   string
		patientID      string
	}
	tests := []struct {
		name      string
		args      args
		wantEmpty bool
	}{
		{
			name: "Happy case: compose episode, full access",
			args: args{
				validPhone:     "0888888888",
				fullAccess:     true,
				organizationID: gofakeit.UUID(),
				providerCode:   "1",
				patientID:      gofakeit.UUID(),
			},
			wantEmpty: false,
		},
		{
			name: "Happy case: compose episode",
			args: args{
				validPhone:     "0888888888",
				fullAccess:     false,
				organizationID: gofakeit.UUID(),
				providerCode:   "1",
				patientID:      gofakeit.UUID(),
			},
			wantEmpty: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComposeOneHealthEpisodeOfCare(tt.args.validPhone, tt.args.fullAccess, tt.args.organizationID, tt.args.providerCode, tt.args.patientID)
			if !tt.wantEmpty {
				assert.NotEmpty(t, got, "ComposeOneHealthEpisodeOfCare() = empty struct, want non-empty struct")
				assert.NotEmpty(t, got.Patient.Reference, "ComposeOneHealthEpisodeOfCare() = %v, PatientReference field is empty", got)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Happy case: valid date",
			args: args{
				date: "2006-01-02T00:00:00",
			},
			want: time.Time{},
		},
		{
			name: "Sad case: invalid date",
			args: args{
				date: "invalid",
			},
			want: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseDate(tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIDToIdentifier(t *testing.T) {
	dummyString := gofakeit.BS()
	base64String := base64.StdEncoding.EncodeToString([]byte(dummyString))
	contentType := enumutils.ContentTypePng
	type args struct {
		ids    []*dto.IdentificationDocument
		phones []*domain.PhoneNumberInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: convert ID to  identifier",
			args: args{
				ids: []*dto.IdentificationDocument{{
					Type:             enums.IDDocumentTypeNationalID,
					Number:           "2121212121",
					Title:            &dummyString,
					ImageContentType: &contentType,
					ImageBase64:      &base64String,
				}},
				phones: []*domain.PhoneNumberInput{{
					Msisdn:             "0888888888",
					VerificationCode:   "0000",
					IsUssd:             false,
					CommunicationOptIn: false,
				}},
			},
			wantErr: false,
		},
		{
			name: "Sad case: missing IDs",
			args: args{
				phones: []*domain.PhoneNumberInput{{
					Msisdn:             "0888888888",
					VerificationCode:   "0000",
					IsUssd:             false,
					CommunicationOptIn: false,
				}},
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid phone",
			args: args{
				ids: []*dto.IdentificationDocument{{
					Type:             enums.IDDocumentTypeNationalID,
					Number:           "2121212121",
					Title:            &dummyString,
					ImageContentType: &contentType,
					ImageBase64:      &base64String,
				}},
				phones: []*domain.PhoneNumberInput{{
					Msisdn:             "invalid phone number",
					VerificationCode:   "0000",
					IsUssd:             false,
					CommunicationOptIn: false,
				}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := IDToIdentifier(tt.args.ids, tt.args.phones)
			if (err != nil) != tt.wantErr {
				t.Errorf("IDToIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNameToHumanName(t *testing.T) {
	dummyString := "test"
	type args struct {
		names []*domain.NameInput
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "Happy case: convert name to human name",
			args: args{
				names: []*domain.NameInput{{
					FirstName:  "test",
					LastName:   "test",
					OtherNames: &dummyString,
				}},
			},
			wantNil: false,
		},
		{
			name:    "Sad case: nil name",
			args:    args{},
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NameToHumanName(tt.args.names)
			if got != nil && tt.wantNil {
				t.Errorf("NameToHumanName() = '%v', want '%v'", got == nil, tt.wantNil)
				return
			}
		})
	}
}

func TestPhysicalPostalAddressesToFHIRAddresses(t *testing.T) {
	type args struct {
		physical []*domain.PhysicalAddress
		postal   []*domain.PostalAddress
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "Happy case: convert postal address to FHIR address",
			args: args{
				physical: []*domain.PhysicalAddress{{
					MapsCode:        gofakeit.BS(),
					PhysicalAddress: gofakeit.BS(),
				}},
				postal: []*domain.PostalAddress{{
					PostalAddress: gofakeit.BS(),
					PostalCode:    gofakeit.BS(),
				}},
			},
			wantNil: false,
		},
		{
			name:    "Sad case: nil address",
			args:    args{},
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PhysicalPostalAddressesToFHIRAddresses(tt.args.physical, tt.args.postal); got != nil && tt.wantNil {
				t.Errorf("PhysicalPostalAddressesToFHIRAddresses() = %v, want %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestMaritalStatusEnumToCodeableConcept(t *testing.T) {
	type args struct {
		val domain.MaritalStatus
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "Happy Case: convert marital status",
			args: args{
				val: domain.MaritalStatusA,
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaritalStatusEnumToCodeableConcept(tt.args.val); got != nil && tt.wantNil {
				t.Errorf("MaritalStatusEnumToCodeableConcept() = %v, want %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestLanguagesToCommunicationInputs(t *testing.T) {
	type args struct {
		languages []enumutils.Language
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "Happy case: convert languages",
			args: args{
				languages: []enumutils.Language{enumutils.LanguageEn},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LanguagesToCommunicationInputs(tt.args.languages); got != nil && tt.wantNil {
				t.Errorf("LanguagesToCommunicationInputs() = %v, want %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestPhysicalPostalAddressesToCombinedFHIRAddress(t *testing.T) {
	type args struct {
		physical []*domain.PhysicalAddress
		postal   []*domain.PostalAddress
	}
	tests := []struct {
		name      string
		args      args
		wantEmpty bool
	}{
		{
			name: "Happy case: convert physical postal address",
			args: args{
				physical: []*domain.PhysicalAddress{{
					MapsCode:        gofakeit.BS(),
					PhysicalAddress: gofakeit.BS(),
				}},
				postal: []*domain.PostalAddress{{
					PostalAddress: gofakeit.BS(),
					PostalCode:    gofakeit.BS(),
				}},
			},
			wantEmpty: false,
		},
		{
			name:      "Sad case: empty input",
			args:      args{},
			wantEmpty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PhysicalPostalAddressesToCombinedFHIRAddress(tt.args.physical, tt.args.postal)
			if !tt.wantEmpty {
				assert.NotEmpty(t, got, "PhysicalPostalAddressesToCombinedFHIRAddress() = empty struct, want non-empty struct")
				assert.NotEmpty(t, got.Period, "PhysicalPostalAddressesToCombinedFHIRAddress() = %v, Period field is empty", got)
			}
			if got != nil && tt.wantEmpty {
				t.Errorf("ContactsToContactPoint() = %v, want %v", got == nil, tt.wantEmpty)
			}
		})
	}
}

func TestContactsToContactPoint(t *testing.T) {
	type args struct {
		ctx             context.Context
		phones          []*domain.PhoneNumberInput
		emails          []*domain.EmailInput
		firestoreClient *firestore.Client
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "Happy case: convert contact",
			args: args{
				ctx: context.Background(),
				phones: []*domain.PhoneNumberInput{{
					Msisdn:             "0888888888",
					VerificationCode:   "0000",
					IsUssd:             false,
					CommunicationOptIn: false,
				}},
				emails: []*domain.EmailInput{{
					Email:              "user@domain.com",
					CommunicationOptIn: false,
				}},
				firestoreClient: &firestore.Client{},
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "Sad case: nil input",
			args: args{
				ctx: context.Background(),
			},
			wantNil: true,
			wantErr: false,
		},
		{
			name: "Sad case: invalid phone",
			args: args{
				ctx: context.Background(),
				phones: []*domain.PhoneNumberInput{{
					Msisdn:             "invalid phone",
					VerificationCode:   "0000",
					IsUssd:             false,
					CommunicationOptIn: false,
				}},
				emails: []*domain.EmailInput{{
					Email:              "user@domain.com",
					CommunicationOptIn: false,
				}},
				firestoreClient: &firestore.Client{},
			},
			wantNil: true,
			wantErr: true,
		},
		{
			name: "Sad case: invalid email",
			args: args{
				ctx: context.Background(),
				phones: []*domain.PhoneNumberInput{{
					Msisdn:             "0888888888",
					VerificationCode:   "0000",
					IsUssd:             false,
					CommunicationOptIn: false,
				}},
				emails: []*domain.EmailInput{{
					Email:              "invalid email address",
					CommunicationOptIn: false,
				}},
				firestoreClient: &firestore.Client{},
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ContactsToContactPoint(tt.args.ctx, tt.args.phones, tt.args.emails, tt.args.firestoreClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactsToContactPoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.wantNil {
				t.Errorf("ContactsToContactPoint() = %v, want %v", got == nil, tt.wantNil)
			}
		})
	}
}

func TestMaritalStatusEnumToCodeableConceptInput(t *testing.T) {
	type args struct {
		val domain.MaritalStatus
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "Happy case: convert marital status",
			args: args{
				val: domain.MaritalStatusA,
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaritalStatusEnumToCodeableConceptInput(tt.args.val); got != nil && tt.wantNil {
				t.Errorf("ContactsToContactPoint() = %v, want %v", got == nil, tt.wantNil)
			}
		})
	}
}
