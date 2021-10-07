package mock

import (
	"context"
	"net/http"

	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
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
