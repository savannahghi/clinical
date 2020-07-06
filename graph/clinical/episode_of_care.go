package clinical

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
)

// GetFHIREpisodeOfCare retrieves instances of FHIREpisodeOfCare by ID
func (s Service) GetFHIREpisodeOfCare(ctx context.Context, id string) (*FHIREpisodeOfCareRelayPayload, error) {
	s.checkPreconditions()

	resourceType := "EpisodeOfCare"
	var resource FHIREpisodeOfCare

	data, err := s.clinicalRepository.GetFHIRResource(resourceType, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s with ID %s, err: %s", resourceType, id, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s data from JSON, err: %v", resourceType, err)
	}

	payload := &FHIREpisodeOfCareRelayPayload{
		Resource: &resource,
	}
	return payload, nil
}

// SearchFHIREpisodeOfCare provides a search API for FHIREpisodeOfCare
func (s Service) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*FHIREpisodeOfCareRelayConnection, error) {
	s.checkPreconditions()

	if params == nil {
		return nil, fmt.Errorf("can't search with nil params")
	}
	urlParams, err := s.validateSearchParams(params)
	if err != nil {
		return nil, err
	}

	resourceName := "EpisodeOfCare"
	path := "_search"
	output := FHIREpisodeOfCareRelayConnection{}

	resources, err := s.searchFilterHelper(ctx, resourceName, path, urlParams)
	if err != nil {
		return nil, err
	}

	for _, result := range resources {
		var resource FHIREpisodeOfCare

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
		output.Edges = append(output.Edges, &FHIREpisodeOfCareRelayEdge{
			Node: &resource,
		})
	}
	return &output, nil
}

// CreateFHIREpisodeOfCare creates a FHIREpisodeOfCare instance
func (s Service) CreateFHIREpisodeOfCare(ctx context.Context, input FHIREpisodeOfCareInput) (*FHIREpisodeOfCareRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "EpisodeOfCare"
	resource := FHIREpisodeOfCare{}

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

	output := &FHIREpisodeOfCareRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// UpdateFHIREpisodeOfCare updates a FHIREpisodeOfCare instance
// The resource must have it's ID set.
func (s Service) UpdateFHIREpisodeOfCare(ctx context.Context, input FHIREpisodeOfCareInput) (*FHIREpisodeOfCareRelayPayload, error) {
	s.checkPreconditions()
	resourceType := "EpisodeOfCare"
	resource := FHIREpisodeOfCare{}

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

	output := &FHIREpisodeOfCareRelayPayload{
		Resource: &resource,
	}
	return output, nil
}

// DeleteFHIREpisodeOfCare deletes the FHIREpisodeOfCare identified by the supplied ID
func (s Service) DeleteFHIREpisodeOfCare(ctx context.Context, id string) (bool, error) {
	resourceType := "EpisodeOfCare"
	resp, err := s.clinicalRepository.DeleteFHIRResource(resourceType, id)
	if err != nil {
		return false, fmt.Errorf(
			"unable to delete %s, response %s, error: %v",
			resourceType, string(resp), err,
		)
	}
	return true, nil
}

// FHIREpisodeOfCare definition: an association between a patient and an organization / healthcare provider(s) during which time encounters may occur. the managing organization assumes a level of responsibility for the patient during this time.
type FHIREpisodeOfCare struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// The EpisodeOfCare may be known by different identifiers for different contexts of use, such as when an external agency is tracking the Episode for funding purposes.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// planned | waitlist | active | onhold | finished | cancelled.
	Status *EpisodeOfCareStatusEnum `json:"status,omitempty"`

	// The history of statuses that the EpisodeOfCare has been through (without requiring processing the history of the resource).
	StatusHistory []*FHIREpisodeofcareStatushistory `json:"statusHistory,omitempty"`

	// A classification of the type of episode of care; e.g. specialist referral, disease management, type of funded care.
	Type []*FHIRCodeableConcept `json:"type,omitempty"`

	// The list of diagnosis relevant to this episode of care.
	Diagnosis []*FHIREpisodeofcareDiagnosis `json:"diagnosis,omitempty"`

	// The patient who is the focus of this episode of care.
	Patient *FHIRReference `json:"patient,omitempty"`

	// The organization that has assumed the specific responsibilities for the specified duration.
	ManagingOrganization *FHIRReference `json:"managingOrganization,omitempty"`

	// The interval during which the managing organization assumes the defined responsibility.
	Period *FHIRPeriod `json:"period,omitempty"`

	// Referral Request(s) that are fulfilled by this EpisodeOfCare, incoming referrals.
	ReferralRequest []*FHIRReference `json:"referralRequest,omitempty"`

	// The practitioner that is the care manager/care coordinator for this patient.
	CareManager *FHIRReference `json:"careManager,omitempty"`

	// The list of practitioners that may be facilitating this episode of care for specific purposes.
	Team []*FHIRReference `json:"team,omitempty"`

	// The set of accounts that may be used for billing for this EpisodeOfCare.
	Account []*FHIRReference `json:"account,omitempty"`
}

// FHIREpisodeOfCareInput is the input type for EpisodeOfCare
type FHIREpisodeOfCareInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// The EpisodeOfCare may be known by different identifiers for different contexts of use, such as when an external agency is tracking the Episode for funding purposes.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// planned | waitlist | active | onhold | finished | cancelled.
	Status *EpisodeOfCareStatusEnum `json:"status,omitempty"`

	// The history of statuses that the EpisodeOfCare has been through (without requiring processing the history of the resource).
	StatusHistory []*FHIREpisodeofcareStatushistoryInput `json:"statusHistory,omitempty"`

	// A classification of the type of episode of care; e.g. specialist referral, disease management, type of funded care.
	Type []*FHIRCodeableConceptInput `json:"type,omitempty"`

	// The list of diagnosis relevant to this episode of care.
	Diagnosis []*FHIREpisodeofcareDiagnosisInput `json:"diagnosis,omitempty"`

	// The patient who is the focus of this episode of care.
	Patient *FHIRReferenceInput `json:"patient,omitempty"`

	// The organization that has assumed the specific responsibilities for the specified duration.
	ManagingOrganization *FHIRReferenceInput `json:"managingOrganization,omitempty"`

	// The interval during which the managing organization assumes the defined responsibility.
	Period *FHIRPeriodInput `json:"period,omitempty"`

	// Referral Request(s) that are fulfilled by this EpisodeOfCare, incoming referrals.
	ReferralRequest []*FHIRReferenceInput `json:"referralRequest,omitempty"`

	// The practitioner that is the care manager/care coordinator for this patient.
	CareManager *FHIRReferenceInput `json:"careManager,omitempty"`

	// The list of practitioners that may be facilitating this episode of care for specific purposes.
	Team []*FHIRReferenceInput `json:"team,omitempty"`

	// The set of accounts that may be used for billing for this EpisodeOfCare.
	Account []*FHIRReferenceInput `json:"account,omitempty"`
}

// FHIREpisodeOfCareRelayConnection is a Relay connection for EpisodeOfCare
type FHIREpisodeOfCareRelayConnection struct {
	Edges []*FHIREpisodeOfCareRelayEdge `json:"edges,omitempty"`

	PageInfo *base.PageInfo `json:"pageInfo,omitempty"`
}

// FHIREpisodeOfCareRelayEdge is a Relay edge for EpisodeOfCare
type FHIREpisodeOfCareRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIREpisodeOfCare `json:"node,omitempty"`
}

// FHIREpisodeOfCareRelayPayload is used to return single instances of EpisodeOfCare
type FHIREpisodeOfCareRelayPayload struct {
	Resource *FHIREpisodeOfCare `json:"resource,omitempty"`
}

// FHIREpisodeofcareDiagnosis definition: an association between a patient and an organization / healthcare provider(s) during which time encounters may occur. the managing organization assumes a level of responsibility for the patient during this time.
type FHIREpisodeofcareDiagnosis struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A list of conditions/problems/diagnoses that this episode of care is intended to be providing care for.
	Condition *FHIRReference `json:"condition,omitempty"`

	// Role that this diagnosis has within the episode of care (e.g. admission, billing, discharge …).
	Role *FHIRCodeableConcept `json:"role,omitempty"`

	// Ranking of the diagnosis (for each role type).
	Rank *string `json:"rank,omitempty"`
}

// FHIREpisodeofcareDiagnosisInput is the input type for EpisodeofcareDiagnosis
type FHIREpisodeofcareDiagnosisInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// A list of conditions/problems/diagnoses that this episode of care is intended to be providing care for.
	Condition *FHIRReferenceInput `json:"condition,omitempty"`

	// Role that this diagnosis has within the episode of care (e.g. admission, billing, discharge …).
	Role *FHIRCodeableConceptInput `json:"role,omitempty"`

	// Ranking of the diagnosis (for each role type).
	Rank *string `json:"rank,omitempty"`
}

// FHIREpisodeofcareStatushistory definition: an association between a patient and an organization / healthcare provider(s) during which time encounters may occur. the managing organization assumes a level of responsibility for the patient during this time.
type FHIREpisodeofcareStatushistory struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// planned | waitlist | active | onhold | finished | cancelled.
	Status *EpisodeOfCareStatusHistoryStatusEnum `json:"status,omitempty"`

	// The period during this EpisodeOfCare that the specific status applied.
	Period *FHIRPeriod `json:"period,omitempty"`
}

// FHIREpisodeofcareStatushistoryInput is the input type for EpisodeofcareStatushistory
type FHIREpisodeofcareStatushistoryInput struct {
	// Unique id for the element within a resource (for internal references). This may be any string value that does not contain spaces.
	ID *string `json:"id,omitempty"`

	// planned | waitlist | active | onhold | finished | cancelled.
	Status *EpisodeOfCareStatusHistoryStatusEnum `json:"status,omitempty"`

	// The period during this EpisodeOfCare that the specific status applied.
	Period *FHIRPeriodInput `json:"period,omitempty"`
}

// EpisodeOfCareStatusEnum is a FHIR enum
type EpisodeOfCareStatusEnum string

const (
	// EpisodeOfCareStatusEnumPlanned ...
	EpisodeOfCareStatusEnumPlanned EpisodeOfCareStatusEnum = "planned"
	// EpisodeOfCareStatusEnumWaitlist ...
	EpisodeOfCareStatusEnumWaitlist EpisodeOfCareStatusEnum = "waitlist"
	// EpisodeOfCareStatusEnumActive ...
	EpisodeOfCareStatusEnumActive EpisodeOfCareStatusEnum = "active"
	// EpisodeOfCareStatusEnumOnhold ...
	EpisodeOfCareStatusEnumOnhold EpisodeOfCareStatusEnum = "onhold"
	// EpisodeOfCareStatusEnumFinished ...
	EpisodeOfCareStatusEnumFinished EpisodeOfCareStatusEnum = "finished"
	// EpisodeOfCareStatusEnumCancelled ...
	EpisodeOfCareStatusEnumCancelled EpisodeOfCareStatusEnum = "cancelled"
	// EpisodeOfCareStatusEnumEnteredInError ...
	EpisodeOfCareStatusEnumEnteredInError EpisodeOfCareStatusEnum = "entered_in_error"
)

// AllEpisodeOfCareStatusEnum ...
var AllEpisodeOfCareStatusEnum = []EpisodeOfCareStatusEnum{
	EpisodeOfCareStatusEnumPlanned,
	EpisodeOfCareStatusEnumWaitlist,
	EpisodeOfCareStatusEnumActive,
	EpisodeOfCareStatusEnumOnhold,
	EpisodeOfCareStatusEnumFinished,
	EpisodeOfCareStatusEnumCancelled,
	EpisodeOfCareStatusEnumEnteredInError,
}

// IsValid ...
func (e EpisodeOfCareStatusEnum) IsValid() bool {
	switch e {
	case EpisodeOfCareStatusEnumPlanned, EpisodeOfCareStatusEnumWaitlist, EpisodeOfCareStatusEnumActive, EpisodeOfCareStatusEnumOnhold, EpisodeOfCareStatusEnumFinished, EpisodeOfCareStatusEnumCancelled, EpisodeOfCareStatusEnumEnteredInError:
		return true
	}
	return false
}

// String ...
func (e EpisodeOfCareStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *EpisodeOfCareStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EpisodeOfCareStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EpisodeOfCareStatusEnum", str)
	}
	return nil
}

// MarshalGQL writes the episode of care status to the supplied writer as a quoted string
func (e EpisodeOfCareStatusEnum) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// EpisodeOfCareStatusHistoryStatusEnum is a FHIR enum
type EpisodeOfCareStatusHistoryStatusEnum string

const (
	// EpisodeOfCareStatusHistoryStatusEnumPlanned ...
	EpisodeOfCareStatusHistoryStatusEnumPlanned EpisodeOfCareStatusHistoryStatusEnum = "planned"
	// EpisodeOfCareStatusHistoryStatusEnumWaitlist ...
	EpisodeOfCareStatusHistoryStatusEnumWaitlist EpisodeOfCareStatusHistoryStatusEnum = "waitlist"
	// EpisodeOfCareStatusHistoryStatusEnumActive ...
	EpisodeOfCareStatusHistoryStatusEnumActive EpisodeOfCareStatusHistoryStatusEnum = "active"
	// EpisodeOfCareStatusHistoryStatusEnumOnhold ...
	EpisodeOfCareStatusHistoryStatusEnumOnhold EpisodeOfCareStatusHistoryStatusEnum = "onhold"
	// EpisodeOfCareStatusHistoryStatusEnumFinished ...
	EpisodeOfCareStatusHistoryStatusEnumFinished EpisodeOfCareStatusHistoryStatusEnum = "finished"
	// EpisodeOfCareStatusHistoryStatusEnumCancelled ...
	EpisodeOfCareStatusHistoryStatusEnumCancelled EpisodeOfCareStatusHistoryStatusEnum = "cancelled"
	// EpisodeOfCareStatusHistoryStatusEnumEnteredInError ...
	EpisodeOfCareStatusHistoryStatusEnumEnteredInError EpisodeOfCareStatusHistoryStatusEnum = "entered_in_error"
)

// AllEpisodeOfCareStatusHistoryStatusEnum ...
var AllEpisodeOfCareStatusHistoryStatusEnum = []EpisodeOfCareStatusHistoryStatusEnum{
	EpisodeOfCareStatusHistoryStatusEnumPlanned,
	EpisodeOfCareStatusHistoryStatusEnumWaitlist,
	EpisodeOfCareStatusHistoryStatusEnumActive,
	EpisodeOfCareStatusHistoryStatusEnumOnhold,
	EpisodeOfCareStatusHistoryStatusEnumFinished,
	EpisodeOfCareStatusHistoryStatusEnumCancelled,
	EpisodeOfCareStatusHistoryStatusEnumEnteredInError,
}

// IsValid ...
func (e EpisodeOfCareStatusHistoryStatusEnum) IsValid() bool {
	switch e {
	case EpisodeOfCareStatusHistoryStatusEnumPlanned, EpisodeOfCareStatusHistoryStatusEnumWaitlist, EpisodeOfCareStatusHistoryStatusEnumActive, EpisodeOfCareStatusHistoryStatusEnumOnhold, EpisodeOfCareStatusHistoryStatusEnumFinished, EpisodeOfCareStatusHistoryStatusEnumCancelled, EpisodeOfCareStatusHistoryStatusEnumEnteredInError:
		return true
	}
	return false
}

// String ...
func (e EpisodeOfCareStatusHistoryStatusEnum) String() string {
	return string(e)
}

// UnmarshalGQL ...
func (e *EpisodeOfCareStatusHistoryStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EpisodeOfCareStatusHistoryStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EpisodeOfCare_StatusHistoryStatusEnum", str)
	}
	return nil
}

// MarshalGQL writes the status of the episode of care status history to the supplied writer as a quoted string
func (e EpisodeOfCareStatusHistoryStatusEnum) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}
