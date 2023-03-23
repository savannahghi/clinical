package dto

import (
	"github.com/go-playground/validator"
	"github.com/savannahghi/scalarutils"
)

type OrganizationIdentifier struct {
	Type  OrganizationIdentifierType `json:"type,omitempty"`
	Value string                     `json:"value,omitempty"`
}

type OrganizationInput struct {
	Name        string                   `json:"name,omitempty"`
	PhoneNumber string                   `json:"phoneNumber,omitempty"`
	Identifiers []OrganizationIdentifier `json:"identifiers,omitempty"`
}

type EpisodeOfCareInput struct {
	Status    EpisodeOfCareStatusEnum `json:"status"`
	PatientID string                  `json:"patientID"`
}

// ObservationInput models the observation input
type ObservationInput struct {
	Status      ObservationStatus `json:"status,omitempty" validate:"required"`
	EncounterID string            `json:"encounterID,omitempty" validate:"required"`
	Value       string            `json:"value,omitempty" validate:"required"`
}

func (o ObservationInput) Validate() error {
	v := validator.New()
	err := v.Struct(o)

	return err
}

type PatientInput struct {
	FirstName   string            `json:"firstName"`
	LastName    string            `json:"lastName"`
	OtherNames  *string           `json:"otherNames"`
	BirthDate   scalarutils.Date  `json:"birthDate"`
	Gender      Gender            `json:"gender"`
	Identifiers []IdentifierInput `json:"identifiers"`
	Contacts    []ContactInput    `json:"contacts"`
}

type IdentifierInput struct {
	Type  IdentifierType `json:"type"`
	Value string         `json:"value"`
}

type ContactInput struct {
	Type  ContactType `json:"type"`
	Value string      `json:"value"`
}

// ConditionInput represents input for creating a FHIR condition
type ConditionInput struct {
	Code        string            `json:"condition"`
	System      string            `json:"system"`
	Status      ConditionStatus   `json:"status"`
	PatientID   string            `json:"patientID"`
	EncounterID string            `json:"encounterID"`
	Note        string            `json:"note"`
	OnsetDate   *scalarutils.Date `json:"onsetDate"`
}
