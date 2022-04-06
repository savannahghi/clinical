package infrastructure

import (
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	"github.com/savannahghi/firebasetools"
)

// Infrastructure ...
type Infrastructure struct {
	FHIRRepo       FHIRRepository
	FirestoreRepo  Repository
	OpenConceptLab openconceptlab.ServiceOCL
	BaseExtension  extensions.BaseExtension
	MyCareHub      mycarehub.IServiceMyCareHub
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor() Infrastructure {
	baseExtension := extensions.NewBaseExtensionImpl(&firebasetools.FirebaseClient{})
	fhirRepository := NewFHIRService()

	firestoreDB := NewDBService()

	myCareHubClient := common.NewInterServiceClient("mycarehub", baseExtension)
	mycarehub := mycarehub.NewServiceMyCareHub(myCareHubClient, baseExtension)

	openconceptlab := openconceptlab.NewServiceOCL()

	return Infrastructure{
		fhirRepository,
		firestoreDB,
		openconceptlab,
		baseExtension,
		mycarehub,
	}
}
