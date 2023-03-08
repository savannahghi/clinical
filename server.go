package main

import (
	"context"
	"strconv"

	"github.com/savannahghi/clinical/pkg/clinical/presentation"
	"github.com/savannahghi/serverutils"
)

func init() {
	// check if must have env variables exist
	// expects the server to die if this not explicitly set
	serverutils.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
	serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
	serverutils.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")
}

func main() {
	ctx := context.Background()

	port, err := strconv.Atoi(serverutils.MustGetEnvVar(serverutils.PortEnvVarName))
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	presentation.StartGin(ctx, port, presentation.ClinicalAllowedOrigins)
}
