package clinical

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func Test_trimString(t *testing.T) {
	type args struct {
		inp       string
		maxLength int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "short string",
			args: args{
				inp:       "a short string",
				maxLength: 20,
			},
			want: "a short string",
		},
		{
			name: "a long string",
			args: args{
				inp:       "this string is longer than the indicated max length",
				maxLength: 20,
			},
			want: "this string is lo...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimString(tt.args.inp, tt.args.maxLength); got != tt.want {
				t.Errorf("trimString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRPatient_RenderOfficialName(t *testing.T) {
	type fields struct {
		ID                   *string
		Text                 *FHIRNarrative
		Identifier           []*FHIRIdentifier
		Active               *bool
		Name                 []*FHIRHumanName
		Telecom              []*FHIRContactPoint
		Gender               *PatientGenderEnum
		BirthDate            *base.Date
		DeceasedBoolean      *bool
		DeceasedDateTime     *base.Date
		Address              []*FHIRAddress
		MaritalStatus        *FHIRCodeableConcept
		MultipleBirthBoolean *bool
		MultipleBirthInteger *string
		Photo                []*FHIRAttachment
		Contact              []*FHIRPatientContact
		Communication        []*FHIRPatientCommunication
		GeneralPractitioner  []*FHIRReference
		ManagingOrganization *FHIRReference
		Link                 []*FHIRPatientLink
	}
	name := FHIRHumanName{
		Text: "John Doe",
		Use:  "official",
	}
	name2 := []*FHIRHumanName{&name}
	tests := []struct {
		name   string
		fields fields
		want   base.Markdown
	}{
		{
			name: "good case: valid name",
			fields: fields{
				Name: name2,
			},
			want: base.Markdown("John Doe"),
		},

		{
			name: "bad case: unknown name",
			fields: fields{
				Name: []*FHIRHumanName{},
			},
			want: base.Markdown("UNKNOWN NAME"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FHIRPatient{
				ID:                   tt.fields.ID,
				Text:                 tt.fields.Text,
				Identifier:           tt.fields.Identifier,
				Active:               tt.fields.Active,
				Name:                 tt.fields.Name,
				Telecom:              tt.fields.Telecom,
				Gender:               tt.fields.Gender,
				BirthDate:            tt.fields.BirthDate,
				DeceasedBoolean:      tt.fields.DeceasedBoolean,
				DeceasedDateTime:     tt.fields.DeceasedDateTime,
				Address:              tt.fields.Address,
				MaritalStatus:        tt.fields.MaritalStatus,
				MultipleBirthBoolean: tt.fields.MultipleBirthBoolean,
				MultipleBirthInteger: tt.fields.MultipleBirthInteger,
				Photo:                tt.fields.Photo,
				Contact:              tt.fields.Contact,
				Communication:        tt.fields.Communication,
				GeneralPractitioner:  tt.fields.GeneralPractitioner,
				ManagingOrganization: tt.fields.ManagingOrganization,
				Link:                 tt.fields.Link,
			}
			if got := p.RenderOfficialName(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FHIRPatient.RenderOfficialName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRPatient_RenderIDDocuments(t *testing.T) {
	type fields struct {
		ID                   *string
		Text                 *FHIRNarrative
		Identifier           []*FHIRIdentifier
		Active               *bool
		Name                 []*FHIRHumanName
		Telecom              []*FHIRContactPoint
		Gender               *PatientGenderEnum
		BirthDate            *base.Date
		DeceasedBoolean      *bool
		DeceasedDateTime     *base.Date
		Address              []*FHIRAddress
		MaritalStatus        *FHIRCodeableConcept
		MultipleBirthBoolean *bool
		MultipleBirthInteger *string
		Photo                []*FHIRAttachment
		Contact              []*FHIRPatientContact
		Communication        []*FHIRPatientCommunication
		GeneralPractitioner  []*FHIRReference
		ManagingOrganization *FHIRReference
		Link                 []*FHIRPatientLink
	}

	identifier := FHIRIdentifier{
		Value: "test",
		Use:   "official",
	}
	identifier2 := []*FHIRIdentifier{&identifier}
	tests := []struct {
		name   string
		fields fields
		want   base.Markdown
	}{
		{
			name: "good case: valid ids",
			fields: fields{
				Identifier: identifier2,
			},
			want: base.Markdown("test (official)"),
		},
		{
			name: "bad case: no ids",
			fields: fields{
				Identifier: []*FHIRIdentifier{},
			},
			want: base.Markdown("No Identification documents found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FHIRPatient{
				ID:                   tt.fields.ID,
				Text:                 tt.fields.Text,
				Identifier:           tt.fields.Identifier,
				Active:               tt.fields.Active,
				Name:                 tt.fields.Name,
				Telecom:              tt.fields.Telecom,
				Gender:               tt.fields.Gender,
				BirthDate:            tt.fields.BirthDate,
				DeceasedBoolean:      tt.fields.DeceasedBoolean,
				DeceasedDateTime:     tt.fields.DeceasedDateTime,
				Address:              tt.fields.Address,
				MaritalStatus:        tt.fields.MaritalStatus,
				MultipleBirthBoolean: tt.fields.MultipleBirthBoolean,
				MultipleBirthInteger: tt.fields.MultipleBirthInteger,
				Photo:                tt.fields.Photo,
				Contact:              tt.fields.Contact,
				Communication:        tt.fields.Communication,
				GeneralPractitioner:  tt.fields.GeneralPractitioner,
				ManagingOrganization: tt.fields.ManagingOrganization,
				Link:                 tt.fields.Link,
			}
			if got := p.RenderIDDocuments(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FHIRPatient.RenderIDDocuments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRPatient_RenderAge(t *testing.T) {
	type fields struct {
		ID                   *string
		Text                 *FHIRNarrative
		Identifier           []*FHIRIdentifier
		Active               *bool
		Name                 []*FHIRHumanName
		Telecom              []*FHIRContactPoint
		Gender               *PatientGenderEnum
		BirthDate            *base.Date
		DeceasedBoolean      *bool
		DeceasedDateTime     *base.Date
		Address              []*FHIRAddress
		MaritalStatus        *FHIRCodeableConcept
		MultipleBirthBoolean *bool
		MultipleBirthInteger *string
		Photo                []*FHIRAttachment
		Contact              []*FHIRPatientContact
		Communication        []*FHIRPatientCommunication
		GeneralPractitioner  []*FHIRReference
		ManagingOrganization *FHIRReference
		Link                 []*FHIRPatientLink
	}
	date := base.Date{Year: 1990, Month: 4, Day: 12}
	tests := []struct {
		name   string
		fields fields
		want   base.Markdown
	}{

		{
			name: "good case: age is present",
			fields: fields{
				BirthDate: &date,
			},
			want: base.Markdown("Age: 30 yrs"),
		},
		{
			name: "bad case: age isnt present",
			fields: fields{
				BirthDate: nil,
			},
			want: base.Markdown("Age: UNKNOWN AGE"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FHIRPatient{
				ID:                   tt.fields.ID,
				Text:                 tt.fields.Text,
				Identifier:           tt.fields.Identifier,
				Active:               tt.fields.Active,
				Name:                 tt.fields.Name,
				Telecom:              tt.fields.Telecom,
				Gender:               tt.fields.Gender,
				BirthDate:            tt.fields.BirthDate,
				DeceasedBoolean:      tt.fields.DeceasedBoolean,
				DeceasedDateTime:     tt.fields.DeceasedDateTime,
				Address:              tt.fields.Address,
				MaritalStatus:        tt.fields.MaritalStatus,
				MultipleBirthBoolean: tt.fields.MultipleBirthBoolean,
				MultipleBirthInteger: tt.fields.MultipleBirthInteger,
				Photo:                tt.fields.Photo,
				Contact:              tt.fields.Contact,
				Communication:        tt.fields.Communication,
				GeneralPractitioner:  tt.fields.GeneralPractitioner,
				ManagingOrganization: tt.fields.ManagingOrganization,
				Link:                 tt.fields.Link,
			}
			if got := p.RenderAge(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FHIRPatient.RenderAge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRPatient_RenderGender(t *testing.T) {
	type fields struct {
		ID                   *string
		Text                 *FHIRNarrative
		Identifier           []*FHIRIdentifier
		Active               *bool
		Name                 []*FHIRHumanName
		Telecom              []*FHIRContactPoint
		Gender               *PatientGenderEnum
		BirthDate            *base.Date
		DeceasedBoolean      *bool
		DeceasedDateTime     *base.Date
		Address              []*FHIRAddress
		MaritalStatus        *FHIRCodeableConcept
		MultipleBirthBoolean *bool
		MultipleBirthInteger *string
		Photo                []*FHIRAttachment
		Contact              []*FHIRPatientContact
		Communication        []*FHIRPatientCommunication
		GeneralPractitioner  []*FHIRReference
		ManagingOrganization *FHIRReference
		Link                 []*FHIRPatientLink
	}
	gender := "male"
	genderEnum := PatientGenderEnum(gender)
	tests := []struct {
		name   string
		fields fields
		want   base.Markdown
	}{
		{
			name: "good case: gender specified",
			fields: fields{
				Gender: &genderEnum,
			},
			want: base.Markdown("Gender: male"),
		},
		{
			name: "bad case: no gender specified",
			fields: fields{
				Gender: nil,
			},
			want: base.Markdown("Gender: UNKNOWN GENDER"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FHIRPatient{
				ID:                   tt.fields.ID,
				Text:                 tt.fields.Text,
				Identifier:           tt.fields.Identifier,
				Active:               tt.fields.Active,
				Name:                 tt.fields.Name,
				Telecom:              tt.fields.Telecom,
				Gender:               tt.fields.Gender,
				BirthDate:            tt.fields.BirthDate,
				DeceasedBoolean:      tt.fields.DeceasedBoolean,
				DeceasedDateTime:     tt.fields.DeceasedDateTime,
				Address:              tt.fields.Address,
				MaritalStatus:        tt.fields.MaritalStatus,
				MultipleBirthBoolean: tt.fields.MultipleBirthBoolean,
				MultipleBirthInteger: tt.fields.MultipleBirthInteger,
				Photo:                tt.fields.Photo,
				Contact:              tt.fields.Contact,
				Communication:        tt.fields.Communication,
				GeneralPractitioner:  tt.fields.GeneralPractitioner,
				ManagingOrganization: tt.fields.ManagingOrganization,
				Link:                 tt.fields.Link,
			}
			if got := p.RenderGender(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FHIRPatient.RenderGender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRPatient_RenderProblems(t *testing.T) {
	type fields struct {
		ID                   *string
		Text                 *FHIRNarrative
		Identifier           []*FHIRIdentifier
		Active               *bool
		Name                 []*FHIRHumanName
		Telecom              []*FHIRContactPoint
		Gender               *PatientGenderEnum
		BirthDate            *base.Date
		DeceasedBoolean      *bool
		DeceasedDateTime     *base.Date
		Address              []*FHIRAddress
		MaritalStatus        *FHIRCodeableConcept
		MultipleBirthBoolean *bool
		MultipleBirthInteger *string
		Photo                []*FHIRAttachment
		Contact              []*FHIRPatientContact
		Communication        []*FHIRPatientCommunication
		GeneralPractitioner  []*FHIRReference
		ManagingOrganization *FHIRReference
		Link                 []*FHIRPatientLink
	}
	id := "a7942fb4-61b4-4cf2-ab39-a2904d3090c3"
	tests := []struct {
		name   string
		fields fields
		want   base.Markdown
	}{
		{
			name: "Case 1: invalid id",
			fields: fields{
				ID: nil,
			},
			want: base.Markdown("Problems: No known problems"),
		},
		{
			name: "Case 2: valid id",
			fields: fields{
				ID: &id,
			},
			want: base.Markdown("Problems: No known problems"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FHIRPatient{
				ID:                   tt.fields.ID,
				Text:                 tt.fields.Text,
				Identifier:           tt.fields.Identifier,
				Active:               tt.fields.Active,
				Name:                 tt.fields.Name,
				Telecom:              tt.fields.Telecom,
				Gender:               tt.fields.Gender,
				BirthDate:            tt.fields.BirthDate,
				DeceasedBoolean:      tt.fields.DeceasedBoolean,
				DeceasedDateTime:     tt.fields.DeceasedDateTime,
				Address:              tt.fields.Address,
				MaritalStatus:        tt.fields.MaritalStatus,
				MultipleBirthBoolean: tt.fields.MultipleBirthBoolean,
				MultipleBirthInteger: tt.fields.MultipleBirthInteger,
				Photo:                tt.fields.Photo,
				Contact:              tt.fields.Contact,
				Communication:        tt.fields.Communication,
				GeneralPractitioner:  tt.fields.GeneralPractitioner,
				ManagingOrganization: tt.fields.ManagingOrganization,
				Link:                 tt.fields.Link,
			}
			if got := p.RenderProblems(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FHIRPatient.RenderProblems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFHIRPatient_RenderAllergies(t *testing.T) {
	type fields struct {
		ID                   *string
		Text                 *FHIRNarrative
		Identifier           []*FHIRIdentifier
		Active               *bool
		Name                 []*FHIRHumanName
		Telecom              []*FHIRContactPoint
		Gender               *PatientGenderEnum
		BirthDate            *base.Date
		DeceasedBoolean      *bool
		DeceasedDateTime     *base.Date
		Address              []*FHIRAddress
		MaritalStatus        *FHIRCodeableConcept
		MultipleBirthBoolean *bool
		MultipleBirthInteger *string
		Photo                []*FHIRAttachment
		Contact              []*FHIRPatientContact
		Communication        []*FHIRPatientCommunication
		GeneralPractitioner  []*FHIRReference
		ManagingOrganization *FHIRReference
		Link                 []*FHIRPatientLink
	}
	id := "a7942fb4-61b4-4cf2-ab39-a2904d3090c3"
	tests := []struct {
		name   string
		fields fields
		want   base.Markdown
	}{
		{
			name: "Case 1: invalid id",
			fields: fields{
				ID: nil,
			},
			want: base.Markdown("Allergies: No known allergies"),
		},
		{
			name: "Case 2: valid id",
			fields: fields{
				ID: &id,
			},
			want: base.Markdown("Allergies: No known allergies"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := FHIRPatient{
				ID:                   tt.fields.ID,
				Text:                 tt.fields.Text,
				Identifier:           tt.fields.Identifier,
				Active:               tt.fields.Active,
				Name:                 tt.fields.Name,
				Telecom:              tt.fields.Telecom,
				Gender:               tt.fields.Gender,
				BirthDate:            tt.fields.BirthDate,
				DeceasedBoolean:      tt.fields.DeceasedBoolean,
				DeceasedDateTime:     tt.fields.DeceasedDateTime,
				Address:              tt.fields.Address,
				MaritalStatus:        tt.fields.MaritalStatus,
				MultipleBirthBoolean: tt.fields.MultipleBirthBoolean,
				MultipleBirthInteger: tt.fields.MultipleBirthInteger,
				Photo:                tt.fields.Photo,
				Contact:              tt.fields.Contact,
				Communication:        tt.fields.Communication,
				GeneralPractitioner:  tt.fields.GeneralPractitioner,
				ManagingOrganization: tt.fields.ManagingOrganization,
				Link:                 tt.fields.Link,
			}
			if got := p.RenderAllergies(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FHIRPatient.RenderAllergies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_RequestUSSDFullHistory(t *testing.T) {
	type args struct {
		ctx   context.Context
		input USSDClinicalRequest
	}

	service := NewService()
	ctx := context.Background()

	patientPayload := PatientResourceFHIRPayload(t)
	patient, err := service.CreateFHIRPatient(ctx, patientPayload)
	if err != nil {
		t.Fatalf("unable to retrieve patient resource %s: ", err)
	}
	patientID := *patient.Resource.ID

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid full history USSD response",
			args: args{
				ctx: ctx,
				input: USSDClinicalRequest{
					PatientID: patientID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.RequestUSSDFullHistory(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RequestUSSDFullHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}
