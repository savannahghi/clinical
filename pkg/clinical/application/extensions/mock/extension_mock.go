package mock

import (
	"context"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
)

// FakeBaseExtension is an mock of the BaseExtension
type FakeBaseExtension struct {
	GetLoggedInUserFn    func(ctx context.Context) (*profileutils.UserInfo, error)
	GetLoggedInUserUIDFn func(ctx context.Context) (string, error)
	NormalizeMSISDNFn    func(msisdn string) (*string, error)
	LoadDepsFromYAMLFn   func() (*interserviceclient.DepsConfig, error)
	SetupISCclientFn     func(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error)
	GetEnvVarFn          func(envName string) (string, error)
	ErrorMapFn           func(err error) map[string]string
	WriteJSONResponseFn  func(
		w http.ResponseWriter,
		source interface{},
		status int,
	)

	MockEnsureTopicsExistFn func(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicIDs []string,
	) error
	MockNamespacePubsubIdentifierFn func(
		serviceName string,
		topicID string,
		environment string,
		version string,
	) string
	MockPublishToPubsubFn func(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicID string,
		environment string,
		serviceName string,
		version string,
		payload []byte,
	) error
	MockEnsureSubscriptionsExistFn func(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicSubscriptionMap map[string]string,
		callbackURL string,
	) error
	MockSubscriptionIDsFn                 func(topicIDs []string) map[string]string
	MockPubSubHandlerPathFn               func() string
	MockVerifyPubSubJWTAndDecodePayloadFn func(
		w http.ResponseWriter,
		r *http.Request,
	) (*pubsubtools.PubSubPayload, error)
	MockGetPubSubTopicFn func(m *pubsubtools.PubSubPayload) (string, error)
}

// GetLoggedInUser retrieves logged in user information
func (b *FakeBaseExtension) GetLoggedInUser(ctx context.Context) (*profileutils.UserInfo, error) {
	return b.GetLoggedInUserFn(ctx)
}

// GetLoggedInUserUID is a mock implementation of GetLoggedInUserUID method
func (b *FakeBaseExtension) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return b.GetLoggedInUserUIDFn(ctx)
}

// NormalizeMSISDN is a mock implementation of NormalizeMSISDN method
func (b *FakeBaseExtension) NormalizeMSISDN(msisdn string) (*string, error) {
	return b.NormalizeMSISDNFn(msisdn)
}

// LoadDepsFromYAML is a mock implementation of LoadDepsFromYAML method
func (b *FakeBaseExtension) LoadDepsFromYAML() (*interserviceclient.DepsConfig, error) {
	return b.LoadDepsFromYAMLFn()
}

// SetupISCclient is a mock implementation of SetupISCclient method
func (b *FakeBaseExtension) SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
	return b.SetupISCclientFn(config, serviceName)
}

// GetEnvVar is a mock implementation of GetEnvVar method
func (b *FakeBaseExtension) GetEnvVar(envName string) (string, error) {
	return b.GetEnvVarFn(envName)
}

// ErrorMap is a mock implementation of ErrorMapG method
func (b *FakeBaseExtension) ErrorMap(err error) map[string]string {
	return b.ErrorMapFn(err)
}

// WriteJSONResponse is a mock implementation of WriteJSONResponse method
func (b *FakeBaseExtension) WriteJSONResponse(
	w http.ResponseWriter,
	source interface{},
	status int,
) {
}

// EnsureTopicsExist creates the topic(s) in the suppplied list if they do not
// already exist.
func (b *FakeBaseExtension) EnsureTopicsExist(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicIDs []string,
) error {
	return b.MockEnsureTopicsExistFn(ctx, pubsubClient, topicIDs)
}

// NamespacePubsubIdentifier uses the service name, environment and version to
// create a "namespaced" pubsub identifier. This could be a topicID or
// subscriptionID.
func (b *FakeBaseExtension) NamespacePubsubIdentifier(
	serviceName string,
	topicID string,
	environment string,
	version string,
) string {
	return b.MockNamespacePubsubIdentifierFn(serviceName, topicID, environment, version)
}

// PublishToPubsub sends the supplied payload to the indicated topic
func (b *FakeBaseExtension) PublishToPubsub(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicID string,
	environment string,
	serviceName string,
	version string,
	payload []byte,
) error {
	return b.MockPublishToPubsubFn(
		ctx,
		pubsubClient,
		topicID,
		environment,
		serviceName,
		version,
		payload,
	)
}

// EnsureSubscriptionsExist ensures that the subscriptions named in the supplied
// topic:subscription map exist. If any does not exist, it is created.
func (b *FakeBaseExtension) EnsureSubscriptionsExist(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicSubscriptionMap map[string]string,
	callbackURL string,
) error {
	return b.MockEnsureSubscriptionsExistFn(
		ctx,
		pubsubClient,
		topicSubscriptionMap,
		callbackURL,
	)
}

// SubscriptionIDs returns a map of topic IDs to subscription IDs
func (b *FakeBaseExtension) SubscriptionIDs(topicIDs []string) map[string]string {
	return b.MockSubscriptionIDsFn(topicIDs)
}

// PubSubHandlerPath returns pubsub hander path `/pubsub`
func (b *FakeBaseExtension) PubSubHandlerPath() string {
	return b.MockPubSubHandlerPathFn()
}

// VerifyPubSubJWTAndDecodePayload confirms that there is a valid Google signed
// JWT and decodes the pubsub message payload into a struct.
//
// It's use will simplify & shorten the handler funcs that process Cloud Pubsub
// push notifications.
func (b *FakeBaseExtension) VerifyPubSubJWTAndDecodePayload(
	w http.ResponseWriter,
	r *http.Request,
) (*pubsubtools.PubSubPayload, error) {
	return b.MockVerifyPubSubJWTAndDecodePayloadFn(w, r)
}

// GetPubSubTopic retrieves a pubsub topic from a pubsub payload.
func (b *FakeBaseExtension) GetPubSubTopic(m *pubsubtools.PubSubPayload) (string, error) {
	return b.MockGetPubSubTopicFn(m)
}
