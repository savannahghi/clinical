package clinical

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
)

// GetFHIROrganization retrieves instances of FHIROrganization by ID
func (s Service) GetFHIROrganization(ctx context.Context, id string) (*FHIROrganizationRelayPayload, error) {
	s.checkPreconditions()

	resourceType := "Organization"
	var resource FHIROrganization

	data, err := s.clinicalRepository.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &FHIROrganizationRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// SearchFHIROrganization provides a search API for FHIROrganization
func (s Service) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*FHIROrganizationRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "Organization"
	path := "_search"
	output := FHIROrganizationRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIROrganization

		resourceBs, err := json.Marshal(result)
		if err != nil {
			logrus.Errorf("unable to marshal map to JSON: %v", err)
			return nil, fmt.Errorf("server error: Unable to marshal map to JSON: %s", err)
		}

		err = json.Unmarshal(resourceBs, &resource)
		if err != nil {
			logrus.Errorf("unable to unmarshal %s: %v", resourceName, err)
			return nil, fmt.Errorf(
				"server error: Unable to unmarshal %s: %s", resourceName, err)
		}
		output.Edges = append(output.Edges, &FHIROrganizationRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIROrganization creates a FHIROrganization instance
func (s Service) CreateFHIROrganization(ctx context.Context, input FHIROrganizationInput) (*FHIROrganizationRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Organization"
	resource := FHIROrganization{}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIROrganizationRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIROrganization updates a FHIROrganization instance
// ! The resource must have it's ID set.
func (s Service) UpdateFHIROrganization(ctx context.Context, input FHIROrganizationInput) (*FHIROrganizationRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "Organization"
	resource := FHIROrganization{}

	if input.ID == nil {
		return nil, fmt.Errorf("can't update with a nil ID")
	}

	payload, err := base.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := s.clinicalRepository.UpdateFHIRResource(resourceType, *input.ID, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &FHIROrganizationRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIROrganization deletes the FHIROrganization identified by the supplied ID
func (s Service) DeleteFHIROrganization(ctx context.Context, id string) (bool, error) {
	resourceType := "Organization"
	resp, err := s.clinicalRepository.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// GetOrganization retrieves an organization given its code
func (s Service) GetOrganization(ctx context.Context, providerSladeCode int) (*string, error) {
	s.checkPreconditions()
	providerCode := strconv.Itoa(providerSladeCode)
	searchParam := map[string]interface{}{
		"identifier": providerCode,
	}
	organization, err := s.SearchFHIROrganization(ctx, searchParam)
	if err != nil {
		return nil, err
	}
	if organization.Edges == nil {
		return nil, nil
	}
	ORGID := organization.Edges[0].Node.ID
	return ORGID, nil

}

// CreateOrganization creates an organization given ist provider code
func (s Service) CreateOrganization(ctx context.Context, providerSladeCode int) (*string, error) {
	s.checkPreconditions()
	providerCode := strconv.Itoa(providerSladeCode)
	identifier := []*FHIRIdentifierInput{
		{
			Use:   "official",
			Value: providerCode,
		},
	}
	organizationInput := FHIROrganizationInput{
		Identifier: identifier,
		Name:       &providerCode,
	}
	createdOrganization, err := s.CreateFHIROrganization(ctx, organizationInput)
	if err != nil {
		return nil, err
	}
	organisationID := createdOrganization.Resource.ID
	return organisationID, nil
}

// GetORCreateOrganization retrieve an organisation via its code if not found create a new one.
func (s Service) GetORCreateOrganization(ctx context.Context, providerSladeCode int) (*string, error) {
	s.checkPreconditions()
	retrievedOrg, err := s.GetOrganization(ctx, providerSladeCode)
	if err != nil {
		return nil, fmt.Errorf(
			"internal server error in getting organisation : %v", err)
	}
	if retrievedOrg == nil {
		createdOrg, err := s.CreateOrganization(ctx, providerSladeCode)
		if err != nil {
			return nil, fmt.Errorf(
				"internal server error in creating organisation : %v", err)
		}
		return createdOrg, nil
	}
	return retrievedOrg, nil
}

// FHIROrganization definition: The organization (facility) responsible for this organization
type FHIROrganization struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`
	// Whether the organization's record is still in active use
	Active *bool `json:"active,omitempty"`
	// Identifier(s) by which this organization is known.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`
	// Specific type of organization (e.g. Healthcare Provider, Hospital Department, Insurance Company).
	Type []*FHIRCodeableConcept `json:"type,omitempty"`
	// Name used for the organization
	Name *string `json:"name,omitempty"`
	// A list of alternate names that the organization is known as, or was known as in the past
	Alias []string `json:"alias,omitempty"`
	//A contact detail for the organization
	// ! Rule: The telecom of an organization can never be of use 'home'
	Telecom []*FHIRContactPoint `json:"telecom,omitempty"`

	// An address for the organization.
	Address []*FHIRAddress `json:"address,omitempty"`
}

// FHIROrganizationInput definition: The organization (facility) responsible for this organization
type FHIROrganizationInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`
	// Whether the organization's record is still in active use
	Active *bool `json:"active,omitempty"`
	// Identifier(s) by which this organization is known.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`
	// Specific type of organization (e.g. Healthcare Provider, Hospital Department, Insurance Company).
	Type []*FHIRCodeableConceptInput `json:"type,omitempty"`
	// Name used for the organization
	Name *string `json:"name,omitempty"`
	// A list of alternate names that the organization is known as, or was known as in the past
	Alias []string `json:"alias,omitempty"`
	//A contact detail for the organization
	// ! Rule: The telecom of an organization can never be of use 'home'
	Telecom []*FHIRContactPointInput `json:"telecom,omitempty"`

	// An address for the organization.
	Address []*FHIRAddressInput `json:"address,omitempty"`
}

// FHIROrganizationRelayPayload is used to return single instances of Organization
type FHIROrganizationRelayPayload struct {
	Resource *FHIROrganization `json:"resource,omitempty"`
}

// FHIROrganizationRelayConnection is a Relay connection for Organization
type FHIROrganizationRelayConnection struct {
	Edges    []*FHIROrganizationRelayEdge `json:"edges,omitempty"`
	PageInfo *base.PageInfo               `json:"pageInfo,omitempty"`
}

// FHIROrganizationRelayEdge is a Relay edge for Organization
type FHIROrganizationRelayEdge struct {
	Cursor *string           `json:"cursor,omitempty"`
	Node   *FHIROrganization `json:"node,omitempty"`
}
