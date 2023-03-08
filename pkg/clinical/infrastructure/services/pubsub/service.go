package pubsubmessaging

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

const (
	// ClinicalServiceName defines the service where the topic is created
	ClinicalServiceName = "clinical"

	// MyCareHubServiceName defines the service where some of the topics have been created
	MyCareHubServiceName = "mycarehub"

	// TopicVersion defines the topic version. That standard one is `v1`
	TopicVersion = "v1"

	// HostNameEnvVarName defines the host name
	HostNameEnvVarName = "SERVICE_HOST"

	// TestTopicName is a topic that is used for testing purposes
	TestTopicName = "pubsub.testing.topic"
)

// ServicePubsub represent all the logic required to interact with pubsub
type ServicePubsub interface {
	TopicIDs() []string
	AddPubSubNamespace(topicName string) string
	PublishToPubsub(
		ctx context.Context,
		topicID string,
		payload []byte,
	) error
	EnsureTopicsExist(
		ctx context.Context,
		topicIDs []string,
	) error
	EnsureSubscriptionsExist(
		ctx context.Context,
	) error
	SubscriptionIDs() map[string]string
	ReceivePubSubPushMessages(
		w http.ResponseWriter,
		r *http.Request,
	)
}

// Clinical represents all the patient business logic
type Clinical interface {
	FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error)
	RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error)
	CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error)
	CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)
	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
}

// ServicePubSubMessaging is used to send and receive pubsub notifications
type ServicePubSubMessaging struct {
	client   *pubsub.Client
	baseExt  extensions.BaseExtension
	infra    infrastructure.Infrastructure
	clinical Clinical
}

// NewServicePubSubMessaging returns a new instance of pubsub
func NewServicePubSubMessaging(
	ctx context.Context,
	client *pubsub.Client,
	baseExt extensions.BaseExtension,
	infra infrastructure.Infrastructure,
	clinical Clinical,
) (*ServicePubSubMessaging, error) {
	s := &ServicePubSubMessaging{
		client:   client,
		baseExt:  baseExt,
		infra:    infra,
		clinical: clinical,
	}

	if err := s.EnsureTopicsExist(
		ctx,
		s.TopicIDs(),
	); err != nil {
		return nil, err
	}

	if err := s.EnsureSubscriptionsExist(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

// AddPubSubNamespace creates unique topics. The topics will be in the form
// <service name>-<topicName>-<environment>-v1
func (ps ServicePubSubMessaging) AddPubSubNamespace(topicName, serviceName string) string {
	environment := serverutils.GetRunningEnvironment()
	return pubsubtools.NamespacePubsubIdentifier(
		serviceName,
		topicName,
		environment,
		TopicVersion,
	)
}

// TopicIDs returns the known (registered) topic IDs
func (ps ServicePubSubMessaging) TopicIDs() []string {
	return []string{
		ps.AddPubSubNamespace(TestTopicName, ClinicalServiceName),
		ps.AddPubSubNamespace(common.CreatePatientTopic, ClinicalServiceName),
		ps.AddPubSubNamespace(common.VitalsTopicName, ClinicalServiceName),
		ps.AddPubSubNamespace(common.MedicationTopicName, ClinicalServiceName),
		ps.AddPubSubNamespace(common.AllergyTopicName, ClinicalServiceName),
		ps.AddPubSubNamespace(common.TestResultTopicName, ClinicalServiceName),
		ps.AddPubSubNamespace(common.TestOrderTopicName, ClinicalServiceName),
		ps.AddPubSubNamespace(common.OrganizationTopicName, ClinicalServiceName),
	}
}

// PublishToPubsub publishes a message to a specified topic
func (ps ServicePubSubMessaging) PublishToPubsub(
	ctx context.Context,
	topicID string,
	payload []byte,
) error {
	environment, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return err
	}
	return pubsubtools.PublishToPubsub(
		ctx,
		ps.client,
		topicID,
		environment,
		ClinicalServiceName,
		TopicVersion,
		payload,
	)
}

// EnsureTopicsExist creates the topic(s) in the supplied list if they do not
// exist
func (ps ServicePubSubMessaging) EnsureTopicsExist(
	ctx context.Context,
	topicIDs []string,
) error {
	return pubsubtools.EnsureTopicsExist(
		ctx,
		ps.client,
		topicIDs,
	)
}

// EnsureSubscriptionsExist ensures that the subscriptions named in the supplied
// topic:subscription map exist. If any does not exist, it is created.
func (ps ServicePubSubMessaging) EnsureSubscriptionsExist(
	ctx context.Context,
) error {
	hostName, err := ps.baseExt.GetEnvVar(HostNameEnvVarName)
	if err != nil {
		return err
	}

	callbackURL := fmt.Sprintf(
		"%s%s",
		hostName,
		pubsubtools.PubSubHandlerPath,
	)

	return pubsubtools.EnsureSubscriptionsExist(
		ctx,
		ps.client,
		ps.SubscriptionIDs(),
		callbackURL,
	)
}

// SubscriptionIDs returns a map of topic IDs to subscription IDs
func (ps ServicePubSubMessaging) SubscriptionIDs() map[string]string {
	return pubsubtools.SubscriptionIDs(ps.TopicIDs())
}
