package pubsubmessaging

import (
	"encoding/json"
	"net/http"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
)

// ReceivePubSubPushMessages receives and processes a pubsub message
func (ps ServicePubSubMessaging) ReceivePubSubPushMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	message, err := ps.baseExt.VerifyPubSubJWTAndDecodePayload(w, r)
	if err != nil {
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	topicID, err := ps.baseExt.GetPubSubTopic(message)
	if err != nil {
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	switch topicID {
	case ps.AddPubSubNamespace(common.CreatePatientTopic, ClinicalServiceName):
		var data dto.CreatePatientPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		profile, err := ps.infra.MyCareHub.UserProfile(ctx, data.UserID)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		year, month, day := profile.DateOfBirth.Date()
		payload := domain.SimplePatientRegistrationInput{
			ID:                      *profile.ID,
			Names:                   []*domain.NameInput{{FirstName: profile.Name, OtherNames: &profile.Username}},
			IdentificationDocuments: []*domain.IdentificationDocument{},
			BirthDate: scalarutils.Date{
				Year:  year,
				Month: int(month),
				Day:   day,
			},
			PhoneNumbers: []*domain.PhoneNumberInput{{Msisdn: profile.Contacts.ContactValue}},
			Gender:       string(profile.Gender),
			Active:       profile.Active,
		}

		patient, err := ps.usecases.RegisterPatient(ctx, payload)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		err = ps.infra.MyCareHub.AddFHIRIDToPatientProfile(ctx, *patient.PatientRecord.ID, *profile.ID)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
	}

	resp := map[string]string{"Status": "Success"}
	returnedResponse, err := json.Marshal(resp)
	if err != nil {
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	_, _ = w.Write(returnedResponse)
}
