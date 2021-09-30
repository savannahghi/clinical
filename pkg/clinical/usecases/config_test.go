package usecases_test

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fb "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/interactor"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	usecaseMock "github.com/savannahghi/clinical/pkg/clinical/usecases/mock"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
)

const ( // Repo the env to identify which repo to use
	Repo = "REPOSITORY"
	//FirebaseRepository is the value of the env when using firebase
	FirebaseRepository = "firebase"
)

var (
	fakeFhir        usecaseMock.FHIRMock
	fakePatient     usecaseMock.ClinicalMock
	fakeUsecaseIntr usecases.Interactor

	testUsecaseInteractor interactor.Usecases
	testInfrastructure    infrastructure.Infrastructure
)

func TestMain(m *testing.M) {
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("CLOUD_HEALTH_PUBSUB_TOPIC", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_DATASET_ID", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_FHIRSTORE_ID", "healthcloud-bewell-fhir-staging")
	os.Setenv("REPOSITORY", "firebase")

	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to Initialize Test Firestore Client")
	}

	if fbc == nil {
		log.Panicf("failed to Initialize Test Firebase Client")
	}

	infra, err := InitializeTestInfrastructure(ctx)
	if err != nil {
		log.Printf("failed to initialize infrastructure: %v", err)
	}

	testInfrastructure = infra

	svc, err := InitializeTestService(ctx, testInfrastructure)
	if err != nil {
		log.Printf("failed to initialize test service: %v", err)
	}

	testUsecaseInteractor = svc

	fakeUsecaseIntr, err = InitializeFakeTestService(&fakePatient, &fakeFhir)
	if err != nil {
		log.Printf("failed to initialize fake test service: %v", err)
	}

	purgeRecords := func() {
		if serverutils.MustGetEnvVar(Repo) == FirebaseRepository {
			r := fb.Repository{}
			collections := []string{
				r.GetEmailOptInCollectionName(),
			}
			for _, collection := range collections {
				ref := fsc.Collection(collection)
				firebasetools.DeleteCollection(ctx, fsc, ref, 10)
			}
		}

	}

	purgeRecords()

	// run the tests
	log.Printf("about to run tests")
	code := m.Run()
	log.Printf("finished running tests")

	// cleanup here
	os.Exit(code)
}

func InitializeTestService(
	ctx context.Context,
	infra infrastructure.Infrastructure,
) (interactor.Usecases, error) {
	i := interactor.NewUsecasesInteractor(infra)
	return i, nil

}

func InitializeTestInfrastructure(
	ctx context.Context,
) (infrastructure.Infrastructure, error) {
	return infrastructure.NewInfrastructureInteractor(), nil
}

func InitializeTestFirebaseClient(
	ctx context.Context,
) (*firestore.Client, *auth.Client) {
	fc := firebasetools.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore: %s", err)
	}

	fbc, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up tests: %s", err)
	}
	return fsc, fbc
}

func InitializeFakeTestService(
	patient usecases.ClinicalUseCase,
	fhir usecases.FHIRUseCase,
) (
	usecases.Interactor,
	error,
) {
	itr := usecases.Interactor{
		patient,
		fhir,
	}
	return itr, nil
}
