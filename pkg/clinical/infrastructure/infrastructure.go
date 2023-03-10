package infrastructure

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/mycarehub"
	"github.com/savannahghi/clinical/pkg/clinical/repository"
)

// ServiceOCL ...
type ServiceOCL interface {
	MakeRequest(method string, path string, params url.Values, body io.Reader) (*http.Response, error)
	ListConcepts(
		ctx context.Context, org string, source string, verbose bool, q *string,
		sortAsc *string, sortDesc *string, conceptClass *string, dataType *string,
		locale *string, includeRetired *bool,
		includeMappings *bool, includeInverseMappings *bool) ([]map[string]interface{}, error)
	GetConcept(
		ctx context.Context, org string, source string, concept string,
		includeMappings bool, includeInverseMappings bool) (map[string]interface{}, error)
}

// Infrastructure ...
type Infrastructure struct {
	FHIR           repository.FHIR
	OpenConceptLab ServiceOCL
	BaseExtension  extensions.BaseExtension
	MyCareHub      mycarehub.IServiceMyCareHub
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor(
	ext extensions.BaseExtension,
	fhir repository.FHIR,
	openconceptlab ServiceOCL,
	mycarehub mycarehub.IServiceMyCareHub,
) Infrastructure {
	return Infrastructure{
		fhir,
		openconceptlab,
		ext,
		mycarehub,
	}
}
