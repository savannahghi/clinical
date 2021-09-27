package acceptance_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fb "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase"
	"github.com/savannahghi/clinical/pkg/clinical/presentation"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/interactor"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
)

const ( // Repo the env to identify which repo to use
	Repo = "REPOSITORY"
	//FirebaseRepository is the value of the env when using firebase
	FirebaseRepository = "firebase"
)

/// these are set up once in TestMain and used by all the acceptance tests in
// this package
var (
	srv            *http.Server
	baseURL        string
	serverErr      error
	testInteractor interactor.Usecases
	testInfra      infrastructure.Infrastructure
)

func initializeAcceptanceTestFirebaseClient(ctx context.Context) (*firestore.Client, *auth.Client) {
	fc := firebasetools.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firestore for the Feed: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore: %s", err)
	}

	fbc, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
	}
	return fsc, fbc
}
func TestMain(m *testing.M) {
	// setup
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	os.Setenv("CLOUD_HEALTH_PUBSUB_TOPIC", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_DATASET_ID", "healthcloud-bewell-staging")
	os.Setenv("CLOUD_HEALTH_FHIRSTORE_ID", "healthcloud-bewell-fhir-staging")
	os.Setenv("REPOSITORY", "firebase")

	ctx := context.Background()
	srv, baseURL, serverErr = serverutils.StartTestServer(
		ctx,
		presentation.PrepareServer,
		presentation.ClinicalAllowedOrigins,
	) // set the globals
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
		return
	}

	fsc, _ := initializeAcceptanceTestFirebaseClient(ctx)

	infra, err := InitializeTestInfrastructure(ctx)
	if err != nil {
		log.Printf("unable to initialize test infrastructure: %v", err)
		return
	}

	testInfra = infra

	i, err := InitializeTestService(ctx, infra)
	if err != nil {
		log.Printf("unable to initialize test service: %v", err)
		return
	}

	testInteractor = i

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
	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()
	os.Exit(code)
}

func InitializeTestService(ctx context.Context, infra infrastructure.Infrastructure) (interactor.Usecases, error) {
	usecases := interactor.NewUsecasesInteractor(
		infra,
	)
	return usecases, nil
}

func InitializeTestInfrastructure(ctx context.Context) (infrastructure.Infrastructure, error) {

	return infrastructure.NewInfrastructureInteractor(), nil
}
