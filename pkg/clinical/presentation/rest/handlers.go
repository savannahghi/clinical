package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
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
	return &PresentationHandlersImpl{usecases: usecases, baseExt: extension}
}

// ReceivePubSubPushMessage receives and processes a pubsub message
func (p PresentationHandlersImpl) ReceivePubSubPushMessage(c *gin.Context) {
	ctx := context.Background()

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
		var data dto.PatientPubSubMessage

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
		var data dto.FacilityPubSubMessage

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
		var data dto.VitalSignPubSubMessage

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
		var data dto.PatientAllergyPubSubMessage

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
		var data dto.MedicationPubSubMessage

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

	case utils.AddPubSubNamespace(common.TenantTopicName, common.ClinicalServiceName):
		var data dto.OrganizationInput

		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

		err = p.usecases.CreatePubsubTenant(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)

			return
		}

	case utils.AddPubSubNamespace(common.TestResultTopicName, common.ClinicalServiceName):
		var data dto.PatientTestResultPubSubMessage

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

	default:
		err := fmt.Errorf("unknown topic ID: %v", topicID)

		serverutils.WriteJSONResponse(c.Writer, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)

		return
	}

	resp := map[string]string{"Status": "Success"}
	c.JSON(http.StatusOK, resp)
}

func (p PresentationHandlersImpl) RegisterTenant(c *gin.Context) {
	input := dto.OrganizationInput{}

	err := c.BindJSON(&input)
	if err != nil {
		resp := map[string]string{"error": err.Error()}
		c.JSON(http.StatusBadRequest, resp)
	}

	organization, err := p.usecases.RegisterTenant(c.Request.Context(), input)
	if err != nil {
		resp := map[string]string{"error": err.Error()}
		c.JSON(http.StatusBadRequest, resp)
	}

	c.JSON(http.StatusOK, organization)
}

func jsonErrorResponse(c *gin.Context, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, gin.H{"error": err.Error()})
}

// RegisterFacility creates a facility in fhir.
func (p PresentationHandlersImpl) RegisterFacility(c *gin.Context) {
	input := dto.OrganizationInput{}

	err := c.BindJSON(&input)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	organization, err := p.usecases.RegisterFacility(c.Request.Context(), input)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, organization)
}

// UploadMedia uploads media to GCS and stores the URL in FHIR attachment
func (p PresentationHandlersImpl) UploadMedia(c *gin.Context) {
	input := &dto.MediaInput{
		EncounterID: c.Request.FormValue("encounterID"),
		File:        c.Request.MultipartForm.File,
	}

	if err := c.Request.ParseMultipartForm(500 << 20); err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	var response []*dto.Media

	for _, fileHeaders := range input.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				jsonErrorResponse(c, http.StatusBadRequest, err)
				return
			}
			defer file.Close()

			contentType := fileHeader.Header.Get("Content-Type")

			output, err := p.usecases.UploadMedia(c.Request.Context(), input.EncounterID, file, contentType)
			if err != nil {
				jsonErrorResponse(c, http.StatusInternalServerError, err)
				return
			}

			response = append(response, output)
		}
	}

	c.JSON(http.StatusOK, response)
}

// LoadQuestionnaire is used to upload a user defined questionnaire for the purpose of soliciting client data.
func (p PresentationHandlersImpl) LoadQuestionnaire(c *gin.Context) {
	input := domain.FHIRQuestionnaire{}

	err := c.BindJSON(&input)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	questionnaire, err := p.usecases.CreateQuestionnaire(c.Request.Context(), &input)
	if err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, questionnaire)
}
