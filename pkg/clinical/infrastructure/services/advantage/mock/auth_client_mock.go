package mock

import "github.com/savannahghi/authutils"

// AuthClientMock mocks the authentication client methods
type AuthClientMock struct {
	MockAuthenticateFn func() (*authutils.OAUTHResponse, error)
}

// NewAuthUtilsClientMock constructor initializes the auth utils client mock
func NewAuthUtilsClientMock() *AuthClientMock {
	return &AuthClientMock{
		MockAuthenticateFn: func() (*authutils.OAUTHResponse, error) {
			return &authutils.OAUTHResponse{
				Scope:        "",
				ExpiresIn:    0,
				AccessToken:  "token",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
			}, nil
		},
	}
}

// Authenticate mocks the implementation of the authentication mechanism provided by auth utils
func (a *AuthClientMock) Authenticate() (*authutils.OAUTHResponse, error) {
	return a.MockAuthenticateFn()
}
