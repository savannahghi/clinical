package usecases_test

import (
	"context"
	"os"
	"testing"

	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	extensionMock "github.com/savannahghi/clinical/pkg/clinical/application/extensions/mock"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fhir "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset"
	fhirRepoMock "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset/mock"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/interactor"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	oclMock "github.com/savannahghi/clinical/pkg/clinical/usecases/ocl/mock"
	log "github.com/sirupsen/logrus"
)

var (
	testUsecaseInteractor interactor.Usecases
	testInfrastructure    infrastructure.Infrastructure

	testFakeInfra usecases.Interactor

	FHIRRepoMock fhirRepoMock.FakeFHIRRepository
	fakeOCL      oclMock.OCLMock

	fakeBaseExtension extensionMock.FakeBaseExtension
)

func TestMain(m *testing.M) {
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("CLOUD_HEALTH_PUBSUB_TOPIC", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_DATASET_ID", "sghi-healthcare-staging")
	os.Setenv("CLOUD_HEALTH_FHIRSTORE_ID", "sghi-healthcare-fhir-staging")
	os.Setenv("REPOSITORY", "firebase")

	ctx := context.Background()

	infra, err := InitializeTestInfrastructure(ctx)
	if err != nil {
		log.Printf("failed to initialize infrastructure: %v", err)
	}
	testInfrastructure = infra

	fakeInfras, err := InitializeFakeTestlInteractor(ctx)
	if err != nil {
		log.Printf("failed to initialize fake usecase interractor: %v", err)
	}

	testFakeInfra = fakeInfras

	svc, err := InitializeTestService(ctx, testInfrastructure)
	if err != nil {
		log.Printf("failed to initialize test service: %v", err)
	}

	testUsecaseInteractor = svc

	// run the tests
	log.Printf("about to run tests\n")
	code := m.Run()
	log.Printf("finished running tests\n")

	// cleanup here
	os.Exit(code)
}

func InitializeTestService(ctx context.Context, infra infrastructure.Infrastructure) (interactor.Usecases, error) {
	i := interactor.NewUsecasesInteractor(infra)
	return i, nil

}

func InitializeTestInfrastructure(ctx context.Context) (infrastructure.Infrastructure, error) {
	baseExtension := extensions.NewBaseExtensionImpl()
	repo := dataset.NewFHIRRepository(ctx)
	fhir := fhir.NewFHIRStoreImpl(repo)
	ocl := openconceptlab.NewServiceOCL()

	return infrastructure.NewInfrastructureInteractor(baseExtension, fhir, ocl), nil
}

func InitializeFakeTestlInteractor(ctx context.Context) (usecases.Interactor, error) {
	var ocl infrastructure.ServiceOCL = &fakeOCL
	var baseExt extensions.BaseExtension = &fakeBaseExtension
	infra := func() infrastructure.Infrastructure {
		return infrastructure.Infrastructure{
			OpenConceptLab: ocl,
			BaseExtension:  baseExt,
		}
	}()

	i := usecases.NewUsecasesInteractor(infra)

	return i, nil
}
