package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
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

	// ProgramIDContextKey is the key used to add a program to the context
	ProgramIDContextKey = ContextKey("ProgramID")

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

// GeneratePatientWelcomeEmailTemplate generates a welcome email
func GeneratePatientWelcomeEmailTemplate() string {
	t := template.Must(template.New("welcomeEmail").Parse(PatientWelcomeEmail))
	buf := new(bytes.Buffer)

	err := t.Execute(buf, "")
	if err != nil {
		log.Fatalf("Error while generating patient welcome email template: %s", err)
	}

	return buf.String()
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
