package utils

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

// ValidateEmail returns an error if the supplied string does not have a
// valid format or resolvable host
func ValidateEmail(email string) error {
	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
