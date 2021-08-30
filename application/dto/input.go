package dto

import (
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/scalarutils"
)

//RegisterUserPayload is the payload used to pass user details to onboarding service
type RegisterUserPayload struct {
	UID         string           `json:"uid,omitempty"`
	FirstName   string           `json:"firstName,omitempty"`
	LastName    string           `json:"lastName,omitempty"`
	Gender      enumutils.Gender `json:"gender,omitempty"`
	PhoneNumber string           `json:"phoneNumber,omitempty"`
	Email       string           `json:"email,omitempty"`
	DateOfBirth scalarutils.Date `json:"dateOfBirth,omitempty"`
}
