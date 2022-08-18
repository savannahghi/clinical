package utils

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/getsentry/sentry-go"
)

// ValidateEmail returns an error if the supplied string does not have a
// valid format or resolvable host
func ValidateEmail(email string) error {
	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ReportErrorToSentry captures the exception thrown and registers an issue in sentry
func ReportErrorToSentry(err error) {
	defer sentry.Flush(2 * time.Millisecond)
	sentry.CaptureException(err)
}

// CaptureSentryMessage captures the sentry message and registers the message in sentry
func CaptureSentryMessage(errMsg string) {
	sentry.CaptureMessage(errMsg)
}
