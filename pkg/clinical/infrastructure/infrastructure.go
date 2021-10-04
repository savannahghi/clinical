package infrastructure

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/engagement"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/onboarding"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	"github.com/savannahghi/firebasetools"
)

// Infrastructure ...
type Infrastructure struct {
	FHIRRepo        FHIRRepository
	FirestoreRepo   Repository
	Engagement      engagement.ServiceEngagement
	FirestoreClient *firestore.Client
	Onboarding      onboarding.ServiceOnboarding
	OpenConceptLab  openconceptlab.ServiceOCL
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor() Infrastructure {
	baseExtension := extensions.NewBaseExtensionImpl(&firebasetools.FirebaseClient{})
	fhirRepository := NewFHIRService()
	engagementClient := common.NewInterServiceClient("engagement", baseExtension)
	engagement := engagement.NewServiceEngagementImpl(engagementClient, baseExtension)
	firestoreDB := NewDBService()
	onboardingClient := common.NewInterServiceClient("onboarding", baseExtension)
	onboarding := onboarding.NewServiceOnboardingImpl(onboardingClient, baseExtension)
	openconceptlab := openconceptlab.NewServiceOCL()

	fc := firebasetools.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firestore for the Feed: %s", err)
	}
	ctx := context.Background()
	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore: %s", err)
	}

	return Infrastructure{
		fhirRepository,
		firestoreDB,
		engagement,
		fsc,
		onboarding,
		openconceptlab,
	}
}
