package testutils

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	fhir "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare"
	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/cloudhealthcare/fhirdataset"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
)

// InitializeTestService sets up the structure that will be used by the usecase layer for
// integration tests
func InitializeTestService(ctx context.Context) (usecases.Interactor, error) {
	baseExtension := extensions.NewBaseExtensionImpl()
	repo := dataset.NewFHIRRepository()
	fhir := fhir.NewFHIRStoreImpl(repo)
	ocl := openconceptlab.NewServiceOCL()

	infrastructure := infrastructure.NewInfrastructureInteractor(baseExtension, fhir, ocl)

	usecases := usecases.NewUsecasesInteractor(infrastructure)

	return usecases, nil
}
