package testutils

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhir"
	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhirdataset"
	fb "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/firebase"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/firebasetools"
)

// InitializeTestService sets up the structure that will be used by the usecase layer for
// integration tests
func InitializeTestService(ctx context.Context) (usecases.Interactor, error) {
	fc := &firebasetools.FirebaseClient{}
	baseExtension := extensions.NewBaseExtensionImpl(fc)
	repo := dataset.NewFHIRRepository()
	fhir := fhir.NewFHIRStoreImpl(repo)
	firestoreExt := fb.NewFirestoreClientExtension(&firestore.Client{})
	fbExt := fb.NewFBClientExtensionImpl()
	f := fb.NewFirebaseRepository(firestoreExt, fbExt)
	ocl := openconceptlab.NewServiceOCL()

	infrastructure := infrastructure.NewInfrastructureInteractor(baseExtension, repo, fhir, f, ocl)

	usecases := usecases.NewUsecasesInteractor(infrastructure)

	return usecases, nil
}
