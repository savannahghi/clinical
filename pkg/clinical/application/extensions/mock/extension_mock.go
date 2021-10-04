package mock

import (
	"context"
	"net/http"

	"github.com/savannahghi/interserviceclient"
)

// FakeBaseExtension is an mock of the BaseExtension
type FakeBaseExtension struct {
	GetLoggedInUserUIDFn func(ctx context.Context) (string, error)
	NormalizeMSISDNFn    func(msisdn string) (*string, error)
	LoadDepsFromYAMLFn   func() (*interserviceclient.DepsConfig, error)
	SetupISCclientFn     func(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error)
	GetEnvVarFn          func(envName string) (string, error)
	ErrorMapGFn          func(err error) map[string]string
	WriteJSONResponseFn  func(
		w http.ResponseWriter,
		source interface{},
		status int,
	)
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

// ErrorMapG is a mock implementation of ErrorMapG method
func (b *FakeBaseExtension) ErrorMapG(err error) map[string]string {
	return b.ErrorMapGFn(err)
}

// WriteJSONResponse is a mock implementation of WriteJSONResponse method
func (b *FakeBaseExtension) WriteJSONResponse(
	w http.ResponseWriter,
	source interface{},
	status int,
) {
}
