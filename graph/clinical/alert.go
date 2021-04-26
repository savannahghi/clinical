package clinical

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

// createAlertMessage Create a nice message to be sent.
func createAlertMessage(names string) string {
	text := fmt.Sprintf(
		"Dear %s. Your health record has been accessed for an emergency. "+
			"If you are not aware of the circumstances of this, please call %s",
		names,
		CallCenterNumber,
	)
	return text
}

// createNextOfKinAlertMessage creates a message to be sent to the next of kin
func createNextOfKinAlertMessage(names, patientName string) string {
	text := fmt.Sprintf(
		"Dear %s. The health record for %s has been accessed for an emergency. "+
			"If you are not aware of the circumstances of this, please call %s",
		names,
		patientName,
		CallCenterNumber,
	)
	return text
}

// generatePatientWelcomeEmailTemplate generates a welcome email
func generatePatientWelcomeEmailTemplate() string {
	t := template.Must(template.New("welcomeEmail").Parse(patientWelcomeEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, "")
	if err != nil {
		log.Fatalf("Error while generating patient welcome email template: %s", err)
	}
	return buf.String()
}
