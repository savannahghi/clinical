package dto

import "github.com/go-playground/validator"

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
