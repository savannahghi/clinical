package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/usecases"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
)

// BaseExtension is an interface that represents some methods in base helper libs
type BaseExtension interface {
	VerifyPubSubJWTAndDecodePayload(w http.ResponseWriter, r *http.Request) (*pubsubtools.PubSubPayload, error)
	GetPubSubTopic(m *pubsubtools.PubSubPayload) (string, error)
}

// PresentationHandlersImpl represents the usecase implementation object
type PresentationHandlersImpl struct {
	usecases usecases.Interactor
	baseExt  BaseExtension
}

// NewPresentationHandlers initializes a new rest handlers usecase
func NewPresentationHandlers(usecases usecases.Interactor, extension BaseExtension) *PresentationHandlersImpl {
	return &PresentationHandlersImpl{usecases: usecases}
}

// ReceivePubSubPushMessage receives and processes a pubsub message
func (p PresentationHandlersImpl) ReceivePubSubPushMessage(c *gin.Context) {
	ctx := c.Request.Context()

	message, err := p.baseExt.VerifyPubSubJWTAndDecodePayload(c.Writer, c.Request)
	if err != nil {
		serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	topicID, err := p.baseExt.GetPubSubTopic(message)
	if err != nil {
		serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	switch topicID {
	case utils.AddPubSubNamespace(common.CreatePatientTopic, common.ClinicalServiceName):
		var data dto.CreatePatientPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.usecases.CreatePubsubPatient(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.OrganizationTopicName, common.ClinicalServiceName):
		var data dto.CreateFacilityPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.usecases.CreatePubsubOrganization(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.VitalsTopicName, common.ClinicalServiceName):
		var data dto.CreateVitalSignPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.usecases.CreatePubsubVitals(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.AllergyTopicName, common.ClinicalServiceName):
		var data dto.CreatePatientAllergyPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.usecases.CreatePubsubAllergyIntolerance(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.MedicationTopicName, common.ClinicalServiceName):
		var data dto.CreateMedicationPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.usecases.CreatePubsubMedicationStatement(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName):
		var data dto.CreatePatientTestResultPubSubMessage

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.usecases.CreatePubsubTestResult(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}
	}

	resp := map[string]string{"Status": "Success"}

	returnedResponse, err := json.Marshal(resp)
	if err != nil {
		serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	_, _ = c.Writer.Write(returnedResponse)
}
