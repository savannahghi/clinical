package rest

import (
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
)

// PresentationHandlers represents all the REST API logic
type PresentationHandlers interface {
}

// PresentationHandlersImpl represents the usecase implementation object
type PresentationHandlersImpl struct {
	usecases usecases.Interactor
}

// NewPresentationHandlers initializes a new rest handlers usecase
func NewPresentationHandlers(usecases usecases.Interactor) PresentationHandlers {
	return &PresentationHandlersImpl{usecases: usecases}
}
