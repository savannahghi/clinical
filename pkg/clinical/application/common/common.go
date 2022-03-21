package common

import (
	"bytes"
	"html/template"
	"log"

	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/interserviceclient"
	"github.com/sirupsen/logrus"
)

// NewInterServiceClient initializes an external service in the correct environment given its name
func NewInterServiceClient(serviceName string, baseExt extensions.BaseExtension) *interserviceclient.InterServiceClient {
	config, err := baseExt.LoadDepsFromYAML()
	if err != nil {
		logrus.Panicf("occurred while opening deps file %v", err)
		return nil
	}

	client, err := baseExt.SetupISCclient(*config, serviceName)
	if err != nil {
		logrus.Panicf("unable to initialize inter service client for %v service: %s", err, serviceName)
		return nil
	}
	return client
}

// GeneratePatientWelcomeEmailTemplate generates a welcome email
func GeneratePatientWelcomeEmailTemplate() string {
	t := template.Must(template.New("welcomeEmail").Parse(utils.PatientWelcomeEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, "")
	if err != nil {
		log.Fatalf("Error while generating patient welcome email template: %s", err)
	}
	return buf.String()
}
