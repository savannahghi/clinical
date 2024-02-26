package graph

import (
	"github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

//go:generate go run github.com/99designs/gqlgen

// Resolver wires up the resolvers needed for the clinical services
type Resolver struct {
	usecases clinical.UseCasesClinicalImpl
}

// NewResolver initializes a working top leve Resolver that has been initialized
// with all necessary dependencies
func NewResolver(usecases clinical.UseCasesClinicalImpl) (*Resolver, error) {
	return &Resolver{
		usecases: usecases,
	}, nil
}

// CheckDependencies ensures that the resolver has what it needs in order to work
func (r *Resolver) CheckDependencies() {

}
