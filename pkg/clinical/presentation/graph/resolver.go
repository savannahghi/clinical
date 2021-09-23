package graph

import (
	"context"
	"log"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/firebasetools"
)

//go:generate go run github.com/99designs/gqlgen

// Resolver wires up the resolvers needed for the clinical services
type Resolver struct {
	usecases usecases.Interactor
}

// NewResolver initializes a working top leve Resolver that has been initialized
// with all necessary dependencies
func NewResolver(
	ctx context.Context,
	usecases usecases.Interactor,
) (*Resolver, error) {
	return &Resolver{
		usecases: usecases,
	}, nil
}

// CheckUserTokenInContext ensures that the context has a valid Firebase auth token
func (r *Resolver) CheckUserTokenInContext(ctx context.Context) *auth.Token {
	token, err := firebasetools.GetUserTokenFromContext(ctx)
	if err != nil {
		log.Panicf("graph.Resolver: context user token is nil")
	}
	return token
}

// CheckDependencies ensures that the resolver has what it needs in order to work
func (r *Resolver) CheckDependencies() {

}
