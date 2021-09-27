package extensions

import (
	"context"
	"net/http"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
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
