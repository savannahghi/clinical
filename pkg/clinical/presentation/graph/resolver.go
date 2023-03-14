package graph

import (
	"context"
	"log"

	"github.com/savannahghi/authutils"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
)

//go:generate go run github.com/99designs/gqlgen

// Resolver wires up the resolvers needed for the clinical services
type Resolver struct {
	usecases usecases.Interactor
}

// NewResolver initializes a working top leve Resolver that has been initialized
// with all necessary dependencies
func NewResolver(usecases usecases.Interactor) (*Resolver, error) {
	return &Resolver{
		usecases: usecases,
	}, nil
}

// CheckUserTokenInContext ensures that the context has a valid slade360 auth token
func (r *Resolver) CheckUserTokenInContext(ctx context.Context) *authutils.TokenIntrospectionResponse {
	token, err := authutils.GetUserTokenFromContext(ctx)
	if err != nil {
		log.Panicf("graph.Resolver: context user token is nil")
	}

	return token
}

// CheckDependencies ensures that the resolver has what it needs in order to work
func (r *Resolver) CheckDependencies() {

}
