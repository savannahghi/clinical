package rest

import (
	"fmt"
	"net/http"

	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/serverutils"
)

// PresentationHandlers represents all the REST API logic
type PresentationHandlers interface {
	DeleteFHIRPatientByPhone() http.HandlerFunc
}

// PresentationHandlersImpl represents the usecase implementation object
type PresentationHandlersImpl struct {
	usecases usecases.Interactor
}

// NewPresentationHandlers initializes a new rest handlers usecase
func NewPresentationHandlers(usecases usecases.Interactor) PresentationHandlers {
	return &PresentationHandlersImpl{usecases: usecases}
}

// DeleteFHIRPatientByPhone handler exposes an endpoint that takes a
// patient's phone number and deletes the patient's FHIR compartment
func (p PresentationHandlersImpl) DeleteFHIRPatientByPhone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &domain.PhoneNumberPayload{}
		type errResponse struct {
			Err string `json:"error"`
		}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.PhoneNumber == "" {
			serverutils.WriteJSONResponse(
				w,
				errResponse{
					Err: "expected a phone number to be defined",
				},
				http.StatusBadRequest,
			)
			return
		}
		deleted, err := p.usecases.DeleteFHIRPatientByPhone(ctx, payload.PhoneNumber)
		if err != nil {
			utils.ReportErrorToSentry(err)
			err := fmt.Sprintf("unable to delete patient: %v", err.Error())
			serverutils.WriteJSONResponse(
				w,
				errResponse{
					Err: err,
				},
				http.StatusInternalServerError,
			)
			return
		}

		type response struct {
			Deleted bool `json:"deleted"`
		}
		serverutils.WriteJSONResponse(
			w,
			response{Deleted: deleted},
			http.StatusOK,
		)
	}
}
