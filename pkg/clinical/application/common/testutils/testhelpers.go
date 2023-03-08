package testutils

import (
	"context"
	"log"

	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fhir "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/serverutils"
	"google.golang.org/api/healthcare/v1"
)

// InitializeTestService sets up the structure that will be used by the usecase layer for
// integration tests
func InitializeTestService(ctx context.Context) (usecases.Interactor, error) {
	project := serverutils.MustGetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	_ = serverutils.MustGetEnvVar("CLOUD_HEALTH_PUBSUB_TOPIC")
	datasetID := serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_ID")
	datasetLocation := serverutils.MustGetEnvVar("CLOUD_HEALTH_DATASET_LOCATION")
	fhirStoreID := serverutils.MustGetEnvVar("CLOUD_HEALTH_FHIRSTORE_ID")
	hsv, err := healthcare.NewService(ctx)
	if err != nil {
		log.Panicf("unable to initialize new Google Cloud Healthcare Service: %s", err)
	}

	repo := dataset.NewFHIRRepository(ctx, hsv, project, datasetID, datasetLocation, fhirStoreID)

	baseExtension := extensions.NewBaseExtensionImpl()
	fhir := fhir.NewFHIRStoreImpl(repo)
	ocl := openconceptlab.NewServiceOCL()

	infrastructure := infrastructure.NewInfrastructureInteractor(baseExtension, fhir, ocl)

	usecases := usecases.NewUsecasesInteractor(infrastructure)

	return usecases, nil
}
