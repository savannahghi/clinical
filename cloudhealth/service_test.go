package cloudhealth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("CLOUD_HEALTH_PUBSUB_TOPIC", "pubtopic")
	os.Setenv("CLOUD_HEALTH_DATASET_ID", "datasetid")
	os.Setenv("CLOUD_HEALTH_FHIRSTORE_ID", "fhirid")
	m.Run()
}

func TestNewService(t *testing.T) {
	s := NewService()
	assert.NotNil(t, s)
	if s != nil {
		s.checkPreconditions() // should not panic
	}
}
