package fhir_test

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("CLOUD_HEALTH_PUBSUB_TOPIC", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_DATASET_ID", "sghi-healthcare-staging")
	os.Setenv("CLOUD_HEALTH_FHIRSTORE_ID", "sghi-healthcare-fhir-staging")
	os.Setenv("REPOSITORY", "firebase")

	// run the tests
	log.Printf("about to run tests\n")
	code := m.Run()
	log.Printf("finished running tests\n")

	// cleanup here
	os.Exit(code)
}
