package extensions

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
)

// ISCClientExtension represents the base ISC client
type ISCClientExtension interface {
	MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)
}

// ISCExtensionImpl ...
type ISCExtensionImpl struct{}

// NewISCExtension initializes an ISC extension
func NewISCExtension() *ISCExtensionImpl {
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
	GetTenantIdentifiers(ctx context.Context) (*dto.TenantIdentifiers, error)
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
}

// NewBaseExtensionImpl ...
func NewBaseExtensionImpl() *BaseExtensionImpl {
	return &BaseExtensionImpl{}
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

// GetTenantIdentifiers retrieves the tenant identifiers e.g OrganizationID, FacilityID from the context
func (b *BaseExtensionImpl) GetTenantIdentifiers(ctx context.Context) (*dto.TenantIdentifiers, error) {
	organizationID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve organization ID from context: %w", err)
	}

	return &dto.TenantIdentifiers{
		OrganizationID: organizationID,
	}, nil
}

// GetOrganizationIDFromContext is a function that retrieves the organization ID from the context.
func GetOrganizationIDFromContext(ctx context.Context) (string, error) {
	organizationID, ok := ctx.Value(utils.OrganizationIDContextKey).(string)
	if !ok {
		return "", errors.New("unable to get organization ID from context")
	}

	return organizationID, nil
}
