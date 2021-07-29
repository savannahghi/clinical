package graph

import (
	"context"
	"log"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/firebasetools"

	"github.com/savannahghi/clinical/graph/clinical"
	"github.com/savannahghi/clinical/graph/openconceptlab"
)

//go:generate go run github.com/99designs/gqlgen

// Resolver wires up the resolvers needed for the clinical services
type Resolver struct {
	clinicalService *clinical.Service
	oclService      *openconceptlab.Service
}

// NewResolver initializes a working top leve Resolver that has been initialized
// with all necessary dependencies
func NewResolver() *Resolver {
	clinicalService := clinical.NewService()
	oclService := openconceptlab.NewService()
	return &Resolver{
		clinicalService: clinicalService,
		oclService:      oclService,
	}
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
	if r.clinicalService == nil {
		log.Panicf("nil clinicalService in resolver")
	}

	if r.oclService == nil {
		log.Panicf("nil oclService in resolver")
	}
}
