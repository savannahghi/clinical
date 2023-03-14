package usecases

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	clinicalUsecase "github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
)

// Clinical represents all the patient business logic
type Clinical interface {
	FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error)

	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	PatientHealthTimeline(ctx context.Context, input domain.HealthTimelineInput) (*domain.HealthTimeline, error)
	GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error)

	CreatePubsubPatient(ctx context.Context, payload dto.CreatePatientPubSubMessage) error
	CreatePubsubOrganization(ctx context.Context, payload dto.CreateFacilityPubSubMessage) error
	CreatePubsubVitals(ctx context.Context, data dto.CreateVitalSignPubSubMessage) error
	CreatePubsubAllergyIntolerance(ctx context.Context, data dto.CreatePatientAllergyPubSubMessage) error
	CreatePubsubTestResult(ctx context.Context, data dto.CreatePatientTestResultPubSubMessage) error
	CreatePubsubMedicationStatement(ctx context.Context, data dto.CreateMedicationPubSubMessage) error
}

// OCL represents all the Open Concept Lab business logic
type OCL interface {
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

// Interactor is an implementation of the usecases interface
type Interactor struct {
	Clinical
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) Interactor {
	clinical := clinicalUsecase.NewUseCasesClinicalImpl(infrastructure)

	impl := Interactor{
		clinical,
	}

	return impl
}
