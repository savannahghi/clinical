package clinical

import (
	"context"
	"testing"

	"github.com/savannahghi/scalarutils"
	"github.com/stretchr/testify/assert"
)

// organizationCodeableConcept - compose an organization codeable concept
func organizationCodeableConcept() []*FHIRCodeableConceptInput {
	display := "Healthcare Provider"
	var code scalarutils.Code = "prov"
	text := "An organization that provides healthcare services"

	codeableConceptInput := SingleCodeableConceptPayload(code, display, text)
	return []*FHIRCodeableConceptInput{codeableConceptInput}
}

//ContactPointPayload - compose a test FHIR contact point input
func organizationContactPointPayload() []*FHIRContactPointInput {
	msisdn := "+254723002959"
	var contactPointSystem ContactPointSystemEnum = "phone"
	var contactPointUse ContactPointUseEnum = "work"
	var rank int64 = 1
	return []*FHIRContactPointInput{
		{
			System: &contactPointSystem,
			Use:    &contactPointUse,
			Value:  &msisdn,
			Rank:   &rank,
		},
	}
}

func validOrganizationPayload() FHIROrganizationInput {
	name := "Ubora Test Hospital"
	alias := []string{"Ubora Clinic", "Clinic Bora"}
	active := true
	return FHIROrganizationInput{
		Active:     &active,
		Identifier: IdentifierPayload(),
		Type:       organizationCodeableConcept(),
		Name:       &name,
		Alias:      alias,
		Telecom:    organizationContactPointPayload(),
		Address:    AddressPayload(),
	}
}

// The telecom for this organization is 'home'
func invalidOrganizationPayload() FHIROrganizationInput {
	name := "Ubora Test Hospital"
	alias := []string{"Ubora Clinic", "Clinic Bora"}
	active := true
	return FHIROrganizationInput{
		Active:     &active,
		Identifier: IdentifierPayload(),
		Type:       organizationCodeableConcept(),
		Name:       &name,
		Alias:      alias,
		Telecom:    ContactPointPayload(),
		Address:    AddressPayload(),
	}
}

func TestService_CreateFHIROrganization(t *testing.T) {
	ctx := context.Background()
	service := NewService()

	type args struct {
		ctx   context.Context
		input FHIROrganizationInput
	}
	tests := []struct {
		name    string
		args    args
		want    *FHIROrganizationRelayPayload
		wantErr bool
	}{
		{
			name:    "Test successful creation of an organization",
			args:    args{input: validOrganizationPayload(), ctx: ctx},
			wantErr: false,
		},
		{
			name:    "Test failure in creation of an organization",
			args:    args{input: invalidOrganizationPayload(), ctx: ctx},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.CreateFHIROrganization(tt.args.ctx, tt.args.input)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, res)
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestService_SearchFHIROrganization(t *testing.T) {
	ctx := context.Background()
	service := NewService()
	// create an organization
	orgPayload := validOrganizationPayload()
	res, err := service.CreateFHIROrganization(ctx, orgPayload)
	if err != nil {
		t.Fatalf("unable to search organization resource %s: ", err)
	}
	validSearchParams := map[string]interface{}{
		"name": *res.Resource.Name,
	}
	inValidSearchParams := map[string]interface{}{
		"name": "Test Hospital Aiko",
	}
	type args struct {
		ctx    context.Context
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Successfully search an existing organization",
			args:    args{ctx: ctx, params: validSearchParams},
			wantErr: false,
		},
		{
			name:    "Search a non existent organization",
			args:    args{ctx: ctx, params: inValidSearchParams},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := service.SearchFHIROrganization(tt.args.ctx, tt.args.params)
			if tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, len(res.Edges), 0)
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}
