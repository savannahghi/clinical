package utils

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/getsentry/sentry-go"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

const (
	// TopicVersion defines the topic version. That standard one is `v1`
	TopicVersion = "v1"
)

// ContextKey is a custom type used as a key value when adding IDs to the context
// This is because using a built-in type as a context key can cause collisions with keys used by other packages, leading to unexpected behavior.
type ContextKey string

const (
	// OrganizationIDContextKey is the key used to add an organizationID to the context
	OrganizationIDContextKey = ContextKey("OrganizationID")

	// FacilityIDContextKey is the key used to add a facility to the context
	FacilityIDContextKey = ContextKey("FacilityID")
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

func AddPubSubNamespace(topicName, serviceName string) string {
	environment := serverutils.GetRunningEnvironment()

	return pubsubtools.NamespacePubsubIdentifier(
		serviceName,
		topicName,
		environment,
		TopicVersion,
	)
}

// CustomError represents a custom error struct
// Reference https://blog.golang.org/error-handling-and-go
type CustomError struct {
	Err     error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// Error implements the error interface
func (e CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
}

// NewCustomError is a helper function to create a new custom error
func NewCustomError(err error, message string) CustomError {
	return CustomError{
		Err:     err,
		Message: message,
	}
}
