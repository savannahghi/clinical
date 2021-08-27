package onboarding

import (
	"context"
	"fmt"
	"net/http"

	"github.com/savannahghi/clinical/application/dto"
	"github.com/savannahghi/interserviceclient"
)

// internal apis definitions
const (
	registerUserURL = "internal/register_user"
)

//Service represents the implemented methods in this ISC
type Service interface {
	CreateUserProfile(ctx context.Context, payload dto.RegisterUserPayload) error
	RemoveTestUser(ctx context.Context, phoneNumber string) error
}

//ServiceImpl represents the implemented methods in this ISC
type ServiceImpl struct {
	isc *interserviceclient.InterServiceClient
}

//NewService initializes a new instance of ServiceImpl
func NewService(isc *interserviceclient.InterServiceClient) *ServiceImpl {
	return &ServiceImpl{
		isc: isc,
	}
}

//CreateUserProfile makes the request to register a user
func (o *ServiceImpl) CreateUserProfile(ctx context.Context, payload dto.RegisterUserPayload) error {

	res, err := o.isc.MakeRequest(ctx, http.MethodPost, registerUserURL, payload)
	if err != nil {
		return fmt.Errorf("unable to send request, error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("register user failed with status code: %v", res.StatusCode)
	}

	return nil
}
