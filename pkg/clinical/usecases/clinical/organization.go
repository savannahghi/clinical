package clinical

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// This file holds all the business logic for creating a FHIR organization. We have a notion of tenants and facilities
// The tenant ID will be used as a logical partitioning key since we want to show that this data resource belongs to this patient who is part of a certain organization(tenant).

// RegisterTenant is used to create an organisation/tenant in the FHIR stores. The tenant ID will be used for logical
// partitioning of data
func (c *UseCasesClinicalImpl) RegisterTenant(ctx context.Context, input dto.OrganizationInput) (*dto.Organization, error) {
	if len(input.Identifiers) == 0 {
		err := fmt.Errorf("expected at least one tenant identifier")
		message := "please provide at least one identifier"

		return nil, utils.NewCustomError(err, message)
	}

	if input.Name == "" {
		err := fmt.Errorf("expected name to be defined")
		message := "please provide the tenant name"

		return nil, utils.NewCustomError(err, message)
	}

	payload := mapOrganizationInputToFHIROrganizationInput(input)

	organisationPayload, err := c.infrastructure.FHIR.CreateFHIROrganization(ctx, *payload)
	if err != nil {
		return nil, err
	}

	// fhir organization to dto organization
	return mapFHIROrganizationToDTOOrganization(organisationPayload.Resource), nil
}

func mapIdentifierToFHIRIdentifierInput(idType, value string) *domain.FHIRIdentifierInput {
	identificationDocumentIdentifierSystem := scalarutils.URI(idType)
	userSelected := true
	idSystem := scalarutils.URI(identificationDocumentIdentifierSystem)
	version := helpers.DefaultVersion

	identifier := domain.FHIRIdentifierInput{
		Use: domain.IdentifierUseEnumOfficial,
		Type: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &identificationDocumentIdentifierSystem,
					Version:      &version,
					Code:         scalarutils.Code(value),
					Display:      value,
					UserSelected: &userSelected,
				},
			},
			Text: value,
		},
		System: &idSystem,
		Value:  value,
		Period: common.DefaultPeriodInput(),
	}

	return &identifier
}

func mapPhoneNumberToFHIRContactPointInput(phoneNumber string) *domain.FHIRContactPointInput {
	use := domain.ContactPointUseEnumWork
	rank := int64(1)
	phoneSystem := domain.ContactPointSystemEnumPhone

	return &domain.FHIRContactPointInput{
		System: &phoneSystem,
		Value:  &phoneNumber,
		Use:    &use,
		Rank:   &rank,
		Period: common.DefaultPeriodInput(),
	}
}

func mapOrganizationInputToFHIROrganizationInput(organization dto.OrganizationInput) *domain.FHIROrganizationInput {
	active := true
	org := domain.FHIROrganizationInput{
		Name:       &organization.Name,
		Active:     &active,
		Telecom:    []*domain.FHIRContactPointInput{},
		Identifier: []*domain.FHIRIdentifierInput{},
	}

	contact := mapPhoneNumberToFHIRContactPointInput(organization.PhoneNumber)
	org.Telecom = append(org.Telecom, contact)

	for _, id := range organization.Identifiers {
		identifier := mapIdentifierToFHIRIdentifierInput(string(id.Type), id.Value)
		org.Identifier = append(org.Identifier, identifier)
	}

	return &org
}

func mapFHIROrganizationToDTOOrganization(organisation *domain.FHIROrganization) *dto.Organization {
	org := &dto.Organization{
		ID:           *organisation.ID,
		Active:       *organisation.Active,
		Name:         *organisation.Name,
		Identifiers:  make([]dto.OrganizationIdentifier, 0),
		PhoneNumbers: make([]string, 0),
	}

	for _, identifier := range organisation.Identifier {
		org.Identifiers = append(org.Identifiers, dto.OrganizationIdentifier{
			Type:  dto.OrganizationIdentifierType(*identifier.System),
			Value: identifier.Value,
		})
	}

	for _, telecom := range organisation.Telecom {
		org.PhoneNumbers = append(org.PhoneNumbers, *telecom.Value)
	}

	return org
}
