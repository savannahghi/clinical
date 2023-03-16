package dto

import (
	"github.com/savannahghi/clinical/pkg/clinical/application/enums"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
	validator "gopkg.in/go-playground/validator.v9"
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

// PatientInput represents the object that hold patient registration input
type PatientInput struct {
	FirstName               string                   `json:"firsname" validate:"required"`
	LastName                string                   `json:"lastName" validate:"required"`
	OtherNames              string                   `json:"otherNames" validate:"required"`
	IdentificationDocuments []IdentificationDocument `json:"identificationDocuments" validate:"required"`
	BirthDate               scalarutils.Date         `json:"birthDate" validate:"required"`
	PhoneNumbers            []string                 `json:"phoneNumbers" validate:"required"`
	Gender                  enumutils.Gender         `json:"gender" validate:"required"`
}

// Validate helps with validation of facility input fields
func (f *PatientInput) Validate() error {
	v := validator.New()

	err := v.Struct(f)

	return err
}

// IdentificationDocument is used to input e.g National ID or passport document
// numbers at patient registration.
type IdentificationDocument struct {
	Type   enums.IDDocumentType `json:"type"`
	Number string               `json:"number"`
}
