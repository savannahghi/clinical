package pubsubmessaging

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

const (

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

// ServicePubSubMessaging is used to send and receive pubsub notifications
type ServicePubSubMessaging struct {
	client  *pubsub.Client
	baseExt extensions.BaseExtension
	infra   infrastructure.Infrastructure
}

// NewServicePubSubMessaging returns a new instance of pubsub
func NewServicePubSubMessaging(
	ctx context.Context,
	client *pubsub.Client,
	baseExt extensions.BaseExtension,
	infra infrastructure.Infrastructure,
) (*ServicePubSubMessaging, error) {
	s := &ServicePubSubMessaging{
		client:  client,
		baseExt: baseExt,
		infra:   infra,
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
		common.TopicVersion,
	)
}

// TopicIDs returns the known (registered) topic IDs
func (ps ServicePubSubMessaging) TopicIDs() []string {
	return []string{
		ps.AddPubSubNamespace(TestTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.TestOrderTopicName, common.ClinicalServiceName),
		ps.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName),
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
		common.ClinicalServiceName,
		common.TopicVersion,
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
