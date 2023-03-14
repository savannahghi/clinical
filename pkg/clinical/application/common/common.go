package common

import (
	"github.com/savannahghi/interserviceclient"
	"github.com/sirupsen/logrus"
)

const (
	// MaxClinicalRecordPageSize is the maximum number of encounters we can show on a timeline
	MaxClinicalRecordPageSize = 50
)

// BaseExtension is an interface that represents some methods in base
// The `onboarding` service has a dependency on `base` library.
// Our first step to making some functions are testable is to remove the base dependency.
// This can be achieved with the below interface.
type BaseExtension interface {
	LoadDepsFromYAML() (*interserviceclient.DepsConfig, error)
	SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error)
}

// NewInterServiceClient initializes an external service in the correct environment given its name
func NewInterServiceClient(serviceName string, baseExt BaseExtension) *interserviceclient.InterServiceClient {
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
