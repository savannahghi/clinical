package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/savannahghi/authutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
)

var (
	clientID      = serverutils.MustGetEnvVar("MYCAREHUB_CLIENT_ID")
	clientSecret  = serverutils.MustGetEnvVar("MYCAREHUB_CLIENT_SECRET")
	introspectURL = serverutils.MustGetEnvVar("MYCAREHUB_INTROSPECT_URL")
)

type IntrospectResponse struct {
	Active bool   `json:"active"`
	UserID string `json:"user_id"`
}

func HasValidMycarehubBearerToken(_ context.Context, r *http.Request) (bool, map[string]string, *authutils.TokenIntrospectionResponse) {
	token, err := firebasetools.ExtractBearerToken(r)
	if err != nil {
		return false, serverutils.ErrorMap(err), nil
	}

	formData := url.Values{
		"token": []string{token},
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodPost, introspectURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return false, serverutils.ErrorMap(err), nil
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return false, serverutils.ErrorMap(err), nil
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("Introspector() failed to close body:%s", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to introspect token")
		return false, serverutils.ErrorMap(err), nil
	}

	var introspection IntrospectResponse

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, serverutils.ErrorMap(err), nil
	}

	if err := json.Unmarshal(bs, &introspection); err != nil {
		return false, serverutils.ErrorMap(err), nil
	}

	if !introspection.Active {
		err := fmt.Errorf("the supplied access token is invalid")
		return false, serverutils.ErrorMap(err), nil
	}

	return true, nil, &authutils.TokenIntrospectionResponse{Token: token, UserGUID: introspection.UserID, IsValid: introspection.Active}
}

// authCheckFn is a function type for authorization and authentication checks
// there can be several e.g an authentication check runs first then an authorization
// check runs next if the authentication passes etc
type authCheckFn = func(
	ctx context.Context,
	r *http.Request,
) (bool, map[string]string, *authutils.TokenIntrospectionResponse)

// AuthenticationGinMiddleware is an authentication middleware for servers using Gin. It checks the user token and ensures
// that it is valid
func AuthenticationGinMiddleware(cl authutils.Client) gin.HandlerFunc {
	checkFuncs := []authCheckFn{cl.HasValidSlade360BearerToken, HasValidMycarehubBearerToken}

	return func(c *gin.Context) {
		var successful bool

		var tokenResponse *authutils.TokenIntrospectionResponse

		errs := []map[string]string{}

		for _, checkFunc := range checkFuncs {
			valid, errMap, authToken := checkFunc(c.Request.Context(), c.Request)
			if valid {
				successful = true
				tokenResponse = authToken

				break
			}

			errs = append(errs, errMap)
		}

		if !successful {
			serverutils.WriteJSONResponse(c.Writer, errs, http.StatusUnauthorized)
			c.Abort()
		}

		ctx := context.WithValue(c.Request.Context(), authutils.AuthTokenContextKey, tokenResponse)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
