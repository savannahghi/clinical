package extensions

import (
	"context"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

// ISCClientExtension represents the base ISC client
type ISCClientExtension interface {
	MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)
}

// ISCExtensionImpl ...
type ISCExtensionImpl struct{}

// NewISCExtension initializes an ISC extension
func NewISCExtension() ISCClientExtension {
	return &ISCExtensionImpl{}
}

// MakeRequest performs an inter service http request and returns a response
func (i *ISCExtensionImpl) MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	var isc interserviceclient.InterServiceClient
	return isc.MakeRequest(ctx, method, path, body)
}

// BaseExtension is an interface that represents some methods in base
// The `onboarding` service has a dependency on `base` library.
// Our first step to making some functions are testable is to remove the base dependency.
// This can be achieved with the below interface.
type BaseExtension interface {
	GetLoggedInUser(ctx context.Context) (*profileutils.UserInfo, error)
	GetLoggedInUserUID(ctx context.Context) (string, error)
	NormalizeMSISDN(msisdn string) (*string, error)
	LoadDepsFromYAML() (*interserviceclient.DepsConfig, error)
	SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error)
	GetEnvVar(envName string) (string, error)
	ErrorMap(err error) map[string]string
	WriteJSONResponse(
		w http.ResponseWriter,
		source interface{},
		status int,
	)

	// PubSub
	EnsureTopicsExist(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicIDs []string,
	) error
	NamespacePubsubIdentifier(
		serviceName string,
		topicID string,
		environment string,
		version string,
	) string
	PublishToPubsub(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicID string,
		environment string,
		serviceName string,
		version string,
		payload []byte,
	) error
	EnsureSubscriptionsExist(
		ctx context.Context,
		pubsubClient *pubsub.Client,
		topicSubscriptionMap map[string]string,
		callbackURL string,
	) error
	SubscriptionIDs(topicIDs []string) map[string]string
	PubSubHandlerPath() string
	VerifyPubSubJWTAndDecodePayload(
		w http.ResponseWriter,
		r *http.Request,
	) (*pubsubtools.PubSubPayload, error)
	GetPubSubTopic(m *pubsubtools.PubSubPayload) (string, error)
}

// BaseExtensionImpl ...
type BaseExtensionImpl struct {
	fc firebasetools.IFirebaseClient
}

// NewBaseExtensionImpl ...
func NewBaseExtensionImpl(fc firebasetools.IFirebaseClient) BaseExtension {
	return &BaseExtensionImpl{
		fc: fc,
	}
}

// GetLoggedInUser retrieves logged in user information
func (b *BaseExtensionImpl) GetLoggedInUser(ctx context.Context) (*profileutils.UserInfo, error) {
	return profileutils.GetLoggedInUser(ctx)
}

// GetLoggedInUserUID get the logged in user uid
func (b *BaseExtensionImpl) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return firebasetools.GetLoggedInUserUID(ctx)
}

// NormalizeMSISDN validates the input phone number.
func (b *BaseExtensionImpl) NormalizeMSISDN(msisdn string) (*string, error) {
	return converterandformatter.NormalizeMSISDN(msisdn)
}

// LoadDepsFromYAML ...
func (b *BaseExtensionImpl) LoadDepsFromYAML() (*interserviceclient.DepsConfig, error) {
	return interserviceclient.LoadDepsFromYAML()
}

// SetupISCclient ...
func (b *BaseExtensionImpl) SetupISCclient(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
	return interserviceclient.SetupISCclient(config, serviceName)
}

// GetEnvVar ...
func (b *BaseExtensionImpl) GetEnvVar(envName string) (string, error) {
	return serverutils.GetEnvVar(envName)
}

// WriteJSONResponse writes the content supplied via the `source` parameter to
// the supplied http ResponseWriter. The response is returned with the indicated
// status.
func (b *BaseExtensionImpl) WriteJSONResponse(
	w http.ResponseWriter,
	source interface{},
	status int,
) {
	serverutils.WriteJSONResponse(w, source, status)
}

// ErrorMap turns the supplied error into a map with "error" as the key
func (b *BaseExtensionImpl) ErrorMap(err error) map[string]string {
	return serverutils.ErrorMap(err)
}

// EnsureTopicsExist creates the topic(s) in the suppplied list if they do not
// already exist.
func (b *BaseExtensionImpl) EnsureTopicsExist(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicIDs []string,
) error {
	return pubsubtools.EnsureTopicsExist(ctx, pubsubClient, topicIDs)
}

// NamespacePubsubIdentifier uses the service name, environment and version to
// create a "namespaced" pubsub identifier. This could be a topicID or
// subscriptionID.
func (b *BaseExtensionImpl) NamespacePubsubIdentifier(
	serviceName string,
	topicID string,
	environment string,
	version string,
) string {
	return pubsubtools.NamespacePubsubIdentifier(
		serviceName,
		topicID,
		environment,
		version,
	)
}

// PublishToPubsub sends the supplied payload to the indicated topic
func (b *BaseExtensionImpl) PublishToPubsub(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicID string,
	environment string,
	serviceName string,
	version string,
	payload []byte,
) error {
	return pubsubtools.PublishToPubsub(
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
func (b *BaseExtensionImpl) EnsureSubscriptionsExist(
	ctx context.Context,
	pubsubClient *pubsub.Client,
	topicSubscriptionMap map[string]string,
	callbackURL string,
) error {
	return pubsubtools.EnsureSubscriptionsExist(
		ctx,
		pubsubClient,
		topicSubscriptionMap,
		callbackURL,
	)
}

// SubscriptionIDs returns a map of topic IDs to subscription IDs
func (b *BaseExtensionImpl) SubscriptionIDs(topicIDs []string) map[string]string {
	return pubsubtools.SubscriptionIDs(topicIDs)
}

// PubSubHandlerPath returns pubsub hander path `/pubsub`
func (b *BaseExtensionImpl) PubSubHandlerPath() string {
	return pubsubtools.PubSubHandlerPath
}

// VerifyPubSubJWTAndDecodePayload confirms that there is a valid Google signed
// JWT and decodes the pubsub message payload into a struct.
//
// It's use will simplify & shorten the handler funcs that process Cloud Pubsub
// push notifications.
func (b *BaseExtensionImpl) VerifyPubSubJWTAndDecodePayload(
	w http.ResponseWriter,
	r *http.Request,
) (*pubsubtools.PubSubPayload, error) {
	return pubsubtools.VerifyPubSubJWTAndDecodePayload(
		w,
		r,
	)
}

// GetPubSubTopic retrieves a pubsub topic from a pubsub payload.
func (b *BaseExtensionImpl) GetPubSubTopic(m *pubsubtools.PubSubPayload) (string, error) {
	return pubsubtools.GetPubSubTopic(m)
}
