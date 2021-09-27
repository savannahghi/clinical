package infrastructure

import (
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/engagement"
	"github.com/savannahghi/firebasetools"
)

// Infrastructure ...
type Infrastructure struct {
	FHIRRepo      FHIRRepository
	FirestoreRepo Repository
	Engagement    engagement.ServiceEngagement
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor() Infrastructure {
	baseExtension := extensions.NewBaseExtensionImpl(&firebasetools.FirebaseClient{})
	fhirRepository := NewFHIRService()
	engagementClient := common.NewInterServiceClient("engagement", baseExtension)
	engagement := engagement.NewServiceEngagementImpl(engagementClient, baseExtension)
	firestore := NewDBService()
	return Infrastructure{
		fhirRepository,
		firestore,
		engagement,
	}
}
