package common

import (
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/interserviceclient"
	"github.com/sirupsen/logrus"
)

const (
	// MaxClinicalRecordPageSize is the maximum number of encounters we can show on a timeline
	MaxClinicalRecordPageSize = 50
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
