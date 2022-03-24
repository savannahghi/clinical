package pubsubmessaging

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
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
		patientName := strings.Split(profile.Name, " ")
		payload := domain.SimplePatientRegistrationInput{
			ID:    *profile.ID,
			Names: []*domain.NameInput{{FirstName: patientName[0], LastName: patientName[1]}},
			BirthDate: scalarutils.Date{
				Year:  year,
				Month: int(month),
				Day:   day,
			},
			PhoneNumbers: []*domain.PhoneNumberInput{{Msisdn: profile.Contacts.ContactValue, CommunicationOptIn: true}},
			Gender:       string(profile.Gender),
			Active:       profile.Active,
		}

		patient, err := ps.patient.RegisterPatient(ctx, payload)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		err = ps.infra.MyCareHub.AddFHIRIDToPatientProfile(ctx, *patient.PatientRecord.ID, data.ID)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.OrganizationTopicName, ClinicalServiceName):
		var data dto.CreateFacilityPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		use := domain.ContactPointUseEnumWork
		rank := int64(1)
		phoneSystem := domain.ContactPointSystemEnumPhone
		input := domain.FHIROrganizationInput{
			ID:     data.ID,
			Active: &data.Active,
			Name:   &data.Name,
			Telecom: []*domain.FHIRContactPointInput{
				{
					System: &phoneSystem,
					Value:  &data.Phone,
					Use:    &use,
					Rank:   &rank,
					Period: common.DefaultPeriodInput(),
				},
			},
		}

		_, err = ps.fhir.CreateFHIROrganization(ctx, input)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.VitalsTopicName, ClinicalServiceName):
		var data dto.CreateVitalSignPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := ps.ocl.GetConcept(
			ctx,
			"CIEL",
			"CIEL",
			*data.ConceptID,
			false,
			false,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		var ConceptPayload domain.Concept
		err = mapstructure.Decode(response, &ConceptPayload)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		value, _ := strconv.ParseFloat(data.Value, 64)
		system := "http://terminology.hl7.org/CodeSystem/observation-category"
		subjectReference := fmt.Sprintf("patient/%v", data.PatientID)
		status := domain.ObservationStatusEnumPreliminary
		input := domain.FHIRObservationInput{
			Status: &status,
			Category: []*domain.FHIRCodeableConceptInput{
				{
					Coding: []*domain.FHIRCodingInput{
						{
							System:  (*scalarutils.URI)(&system),
							Code:    "vital-signs",
							Display: "Vital Signs",
						},
					},
					Text: "Vital Signs",
				},
			},

			Code: domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&ConceptPayload.URL),
						Code:    scalarutils.Code(ConceptPayload.ID),
						Display: ConceptPayload.DisplayName,
					},
				},
				Text: ConceptPayload.DisplayName,
			},

			ValueQuantity: &domain.FHIRQuantityInput{
				Value:  value,
				Unit:   "°C",
				System: "http://unitsofmeasure.org",
				Code:   "(C)",
			},

			Subject: &domain.FHIRReferenceInput{
				Reference: &subjectReference,
				Display:   data.PatientID,
			},
		}

		_, err = ps.fhir.CreateFHIRObservation(ctx, input)
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
