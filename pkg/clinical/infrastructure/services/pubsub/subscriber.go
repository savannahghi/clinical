package pubsubmessaging

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

		response, err := ps.fhir.CreateFHIROrganization(ctx, input)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		err = ps.infra.MyCareHub.AddFHIRIDToFacility(ctx, *response.Resource.ID, *data.ID)
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

		var patientName string
		patient, err := ps.patient.FindPatientByID(ctx, data.PatientID)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		patientName = *patient.PatientRecord.Name[0].Given[0]

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

		system := "http://terminology.hl7.org/CodeSystem/observation-category"
		subjectReference := fmt.Sprintf("Patient/%v", data.PatientID)
		status := domain.ObservationStatusEnumPreliminary
		instant := scalarutils.Instant(data.Date.Format(time.RFC3339))
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
			EffectiveInstant: &instant,
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
			ValueString: &data.Value,
			Subject: &domain.FHIRReferenceInput{
				Reference: &subjectReference,
				Display:   patientName,
			},
		}

		if data.OrganizationID != "" {
			organization, err := ps.fhir.FindOrganizationByID(ctx, data.OrganizationID)
			if err != nil {
				//Should not fail incase the organization is not found
				log.Printf("the error is: %v", err)
			}

			if organization != nil {
				performerReference := fmt.Sprintf("Organization/%v", data.OrganizationID)
				referenceInput := &domain.FHIRReferenceInput{
					Reference: &performerReference,
					Display:   *organization.Resource.Name,
				}

				input.Performer = append(input.Performer, referenceInput)
			}
		}

		_, err = ps.fhir.CreateFHIRObservation(ctx, input)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.AllergyTopicName, ClinicalServiceName):
		var data dto.CreatePatientAllergyPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		var patientName string
		patient, err := ps.patient.FindPatientByID(ctx, data.PatientID)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		patientName = *patient.PatientRecord.Name[0].Given[0]

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

		allergyType := domain.AllergyIntoleranceTypeEnumAllergy
		allergyCategory := domain.AllergyIntoleranceCategoryEnumMedication
		year, month, day := data.Date.Date()
		subjectReference := fmt.Sprintf("Patient/%v", data.PatientID)
		severity := data.Severity.Name
		input := domain.FHIRAllergyIntoleranceInput{
			Type: &allergyType,
			RecordedDate: &scalarutils.Date{
				Year:  year,
				Month: int(month),
				Day:   day,
			},
			Category: []*domain.AllergyIntoleranceCategoryEnum{&allergyCategory},
			ClinicalStatus: domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&ConceptPayload.URL),
						Code:    scalarutils.Code(ConceptPayload.ID),
						Display: ConceptPayload.DisplayName,
					},
				},
				Text: ConceptPayload.DisplayName,
			},
			VerificationStatus: domain.FHIRCodeableConceptInput{},
			Patient: &domain.FHIRReferenceInput{
				Reference: &subjectReference,
				Display:   patientName,
			},
			Reaction: []*domain.FHIRAllergyintoleranceReactionInput{
				{
					Substance: &domain.FHIRCodeableConceptInput{
						Coding: []*domain.FHIRCodingInput{
							{
								System:  (*scalarutils.URI)(&ConceptPayload.URL),
								Code:    scalarutils.Code(ConceptPayload.ID),
								Display: ConceptPayload.DisplayName,
							},
						},
						Text: ConceptPayload.DisplayName,
					},
					Description: &data.Name,
					Severity:    (*domain.AllergyIntoleranceReactionSeverityEnum)(&severity),
				},
			},
		}

		_, err = ps.fhir.CreateFHIRAllergyIntolerance(ctx, input)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.MedicationTopicName, ClinicalServiceName):
		var data dto.CreateMedicationPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		var patientName string
		patient, err := ps.patient.FindPatientByID(ctx, data.PatientID)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		patientName = *patient.PatientRecord.Name[0].Given[0]

		resp, err := ps.ocl.GetConcept(
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

		response, err := ps.ocl.GetConcept(
			ctx,
			"CIEL",
			"CIEL",
			*data.Drug.ConceptID,
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

		var StatementConceptPayload domain.Concept
		err = mapstructure.Decode(resp, &StatementConceptPayload)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		var DrugConceptPayload domain.Concept
		err = mapstructure.Decode(response, &DrugConceptPayload)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		year, month, day := data.Date.Date()
		subjectReference := fmt.Sprintf("Patient/%v", data.PatientID)
		status := domain.MedicationStatementStatusEnumUnknown
		msInput := domain.FHIRMedicationStatementInput{
			Status: &status,
			Category: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&StatementConceptPayload.URL),
						Code:    scalarutils.Code(StatementConceptPayload.ID),
						Display: StatementConceptPayload.DisplayName,
					},
				},
				Text: StatementConceptPayload.DisplayName,
			},
			MedicationCodeableConcept: &domain.FHIRCodeableConceptInput{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&DrugConceptPayload.URL),
						Code:    scalarutils.Code(DrugConceptPayload.ID),
						Display: DrugConceptPayload.DisplayName,
					},
				},
				Text: DrugConceptPayload.DisplayName,
			},
			Subject: &domain.FHIRReferenceInput{
				Reference: &subjectReference,
				Display:   patientName,
			},
			EffectiveDateTime: &scalarutils.Date{
				Year:  year,
				Month: int(month),
				Day:   day,
			},
		}

		if data.OrganizationID != "" {
			organization, err := ps.fhir.FindOrganizationByID(ctx, data.OrganizationID) // rename organization resposne
			if err != nil {
				//Should not fail incase the organization is not found
				log.Printf("the error is: %v", err)
			}
			if organization != nil {
				informationSourceReference := fmt.Sprintf("Organization/%v", data.OrganizationID)

				referenceInput := &domain.FHIRReferenceInput{
					Reference: &informationSourceReference,
					Display:   *organization.Resource.Name,
				}

				msInput.InformationSource = referenceInput
			}
		}

		_, err = ps.fhir.CreateFHIRMedicationStatement(ctx, msInput)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

	case ps.AddPubSubNamespace(common.TestResultTopicName, ClinicalServiceName):
		var data dto.CreatePatientTestResultPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		var patientName string
		patient, err := ps.patient.FindPatientByID(ctx, data.PatientID)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		patientName = *patient.PatientRecord.Name[0].Given[0]

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

		system := "http://terminology.hl7.org/CodeSystem/observation-category"
		subjectReference := fmt.Sprintf("Patient/%v", data.PatientID)
		status := domain.ObservationStatusEnumPreliminary
		instant := scalarutils.Instant(data.Date.Format(time.RFC3339))
		input := domain.FHIRObservationInput{
			Status: &status,
			Category: []*domain.FHIRCodeableConceptInput{
				{
					Coding: []*domain.FHIRCodingInput{
						{
							System:  (*scalarutils.URI)(&system),
							Code:    "laboratory",
							Display: "Laboratory",
						},
					},
					Text: "Laboratory",
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
			ValueString:      &data.Result.Name,
			EffectiveInstant: &instant,
			Subject: &domain.FHIRReferenceInput{
				Reference: &subjectReference,
				Display:   patientName,
			},
		}

		if data.OrganizationID != "" {
			organization, err := ps.fhir.FindOrganizationByID(ctx, data.OrganizationID) // rename organization resposne
			if err != nil {
				//Should not fail incase the organization is not found
				log.Printf("the error is: %v", err)
			}
			if organization != nil {
				performerReference := fmt.Sprintf("Organization/%v", data.OrganizationID)

				referenceInput := &domain.FHIRReferenceInput{
					Reference: &performerReference,
					Display:   *organization.Resource.Name,
				}

				input.Performer = append(input.Performer, referenceInput)
			}
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
