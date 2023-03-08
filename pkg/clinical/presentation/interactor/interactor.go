package interactor

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/clinical"
	"github.com/savannahghi/clinical/pkg/clinical/usecases/ocl"
)

// UseCasesClinical represents all the patient business logic
type Clinical interface {
	FindOrganizationByID(ctx context.Context, organizationID string) (*domain.FHIROrganizationRelayPayload, error)
	RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error)
	CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error)
	CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error)
	FindPatientByID(ctx context.Context, id string) (*domain.PatientPayload, error)
	CreateFHIRMedicationStatement(ctx context.Context, input domain.FHIRMedicationStatementInput) (*domain.FHIRMedicationStatementRelayPayload, error)

	CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error)
	PatientTimeline(ctx context.Context, patientID string, count int) ([]map[string]interface{}, error)
	PatientHealthTimeline(ctx context.Context, input domain.HealthTimelineInput) (*domain.HealthTimeline, error)
	GetMedicalData(ctx context.Context, patientID string) (*domain.MedicalData, error)
}

// OCLUsecase represents all the Open Concept Lab business logic
type OCLUsecase interface {
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

// Usecases is an interface that combines of all usescases
type Usecases interface {
	Clinical
	OCLUsecase
}

// Interactor is an implementation of the usecases interface
type Interactor struct {
	Clinical
	infrastructure.Infrastructure
	OCLUsecase
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
) *Interactor {
	clinical := clinical.NewUseCasesClinicalImpl(infrastructure)
	ocl := ocl.NewUseCasesOCLImpl(infrastructure)

	impl := &Interactor{
		clinical,
		infrastructure,
		ocl,
	}

	return impl
}
