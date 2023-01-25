package infrastructure

import (
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	fhir "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhir"
	dataset "github.com/savannahghi/clinical/pkg/clinical/infrastructure/datastore/fhirdataset"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
)

// Infrastructure ...
type Infrastructure struct {
	FHIRRepo       dataset.FHIRRepository
	FHIR           fhir.FHIR
	OpenConceptLab openconceptlab.ServiceOCL
	BaseExtension  extensions.BaseExtension
	MyCareHub      mycarehub.IServiceMyCareHub
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor(
	ext extensions.BaseExtension,
	fhirRepository dataset.FHIRRepository,
	fhir fhir.FHIR,
	openconceptlab openconceptlab.ServiceOCL,
) Infrastructure {
	myCareHubClient := common.NewInterServiceClient("mycarehub", ext)
	mycarehub := mycarehub.NewServiceMyCareHub(myCareHubClient, ext)

	return Infrastructure{
		fhirRepository,
		fhir,
		openconceptlab,
		ext,
		mycarehub,
	}
}
