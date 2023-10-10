package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chenyahui/gin-cache/persist"
	"github.com/savannahghi/authutils"
)

func TestHasValidCachedToken(t *testing.T) {
	type args struct {
		cacheStore persist.CacheStore
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy case: generate cached token",
			args: args{
				cacheStore: persist.NewMemoryStore(1 * time.Second),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasValidCachedToken(tt.args.cacheStore); got == nil {
				t.Errorf("HasValidCachedToken() is nil")
			}
		})
	}
}

func TestCachedToken(t *testing.T) {
	cacheStore := persist.NewMemoryStore(1 * time.Second)

	testRequestWithStoredToken := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	testRequestWithStoredToken.Header.Set("Authorization", "Bearer xyz")
	cacheStore.Set("xyz", authutils.TokenIntrospectionResponse{}, time.Until(time.Now().Add(1*time.Hour)))

	testRequestWithToken := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	testRequestWithToken.Header.Set("Authorization", "Bearer zxy")

	testRequestNoToken := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)

	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "happy case: cached token",
			args: args{
				ctx: nil,
				r:   testRequestWithStoredToken,
			},
			want: true,
		},
		{
			name: "sad case: missing bearer token in request",
			args: args{
				ctx: nil,
				r:   testRequestNoToken,
			},
			want: false,
		},
		{
			name: "sad case: missing token in cache",
			args: args{
				ctx: nil,
				r:   testRequestWithToken,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := HasValidCachedToken(cacheStore)

			got, _, _ := middleware(tt.args.ctx, tt.args.r)
			if got != tt.want {
				t.Errorf("CachedToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
