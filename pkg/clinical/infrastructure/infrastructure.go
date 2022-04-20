package infrastructure

import (
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/engagement"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/onboarding"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
)

// Infrastructure ...
type Infrastructure struct {
	FHIRRepo       FHIRRepository
	FirestoreRepo  Repository
	Engagement     engagement.ServiceEngagement
	Onboarding     onboarding.ServiceOnboarding
	OpenConceptLab openconceptlab.ServiceOCL
	BaseExtension  extensions.BaseExtension
	MyCareHub      mycarehub.IServiceMyCareHub
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor(
	baseExtension extensions.BaseExtension,
) Infrastructure {
	fhirRepository := NewFHIRService()

	engagementClient := common.NewInterServiceClient("engagement", baseExtension)
	engagement := engagement.NewServiceEngagementImpl(engagementClient, baseExtension)

	firestoreDB := NewDBService()

	onboardingClient := common.NewInterServiceClient("onboarding", baseExtension)
	onboarding := onboarding.NewServiceOnboardingImpl(onboardingClient, baseExtension)

	myCareHubClient := common.NewInterServiceClient("mycarehub", baseExtension)
	mycarehub := mycarehub.NewServiceMyCareHub(myCareHubClient, baseExtension)

	openconceptlab := openconceptlab.NewServiceOCL()

	return Infrastructure{
		fhirRepository,
		firestoreDB,
		engagement,
		onboarding,
		openconceptlab,
		baseExtension,
		mycarehub,
	}
}
