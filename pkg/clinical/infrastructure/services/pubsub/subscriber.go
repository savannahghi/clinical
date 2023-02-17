package pubsubmessaging

import (
	"context"
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
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure/services/openconceptlab"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
)

var (
	fhirAllergyIntoleranceClinicalStatusURL     = "http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical"
	fhirAllergyIntoleranceVerificationStatusURL = "http://terminology.hl7.org/CodeSystem/allergyintolerance-verification"
	unknownConceptID                            = "1067"
)

// ReceivePubSubPushMessages receives and processes a pubsub message
func (ps ServicePubSubMessaging) ReceivePubSubPushMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	message, err := pubsubtools.VerifyPubSubJWTAndDecodePayload(w, r)
	if err != nil {
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	topicID, err := pubsubtools.GetPubSubTopic(message)
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

		patient, err := ps.clinical.RegisterPatient(ctx, payload)
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

		response, err := ps.clinical.CreateFHIROrganization(ctx, input)
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

		input, err := ps.ComposeVitalsInput(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		_, err = ps.clinical.CreateFHIRObservation(ctx, *input)
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

		input, err := ps.ComposeAllergyIntoleranceInput(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		_, err = ps.clinical.CreateFHIRAllergyIntolerance(ctx, *input)
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

		input, err := ps.ComposeMedicationStatementInput(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		_, err = ps.clinical.CreateFHIRMedicationStatement(ctx, *input)
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

		input, err := ps.ComposeTestResultInput(ctx, data)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		_, err = ps.clinical.CreateFHIRObservation(ctx, *input)
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

// ComposeTestResultInput composes a test result input from data received
func (ps ServicePubSubMessaging) ComposeTestResultInput(ctx context.Context, input dto.CreatePatientTestResultPubSubMessage) (*domain.FHIRObservationInput, error) {
	var patientName string
	patient, err := ps.clinical.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}
	patientName = *patient.PatientRecord.Name[0].Given[0]

	observationConcept, err := getCIELConcept(ctx, ps.infra.OpenConceptLab, *input.ConceptID)
	if err != nil {
		return nil, err
	}

	system := "http://terminology.hl7.org/CodeSystem/observation-category"
	subjectReference := fmt.Sprintf("Patient/%v", input.PatientID)
	status := domain.ObservationStatusEnumPreliminary
	instant := scalarutils.Instant(input.Date.Format(time.RFC3339))

	observation := domain.FHIRObservationInput{
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
					System:  (*scalarutils.URI)(&observationConcept.URL),
					Code:    scalarutils.Code(observationConcept.ID),
					Display: observationConcept.DisplayName,
				},
			},
			Text: observationConcept.DisplayName,
		},
		ValueString:      &input.Result.Name,
		EffectiveInstant: &instant,
		Subject: &domain.FHIRReferenceInput{
			Reference: &subjectReference,
			Display:   patientName,
		},
	}

	if input.OrganizationID != "" {
		organization, err := ps.clinical.FindOrganizationByID(ctx, input.OrganizationID) // rename organization response
		if err != nil {
			//Should not fail if the organization is not found
			log.Printf("the error is: %v", err)
		}
		if organization != nil {
			performer := fmt.Sprintf("Organization/%v", input.OrganizationID)

			referenceInput := &domain.FHIRReferenceInput{
				Reference: &performer,
				Display:   *organization.Resource.Name,
			}

			observation.Performer = append(observation.Performer, referenceInput)
		}
	}

	return &observation, nil
}

// ComposeVitalsInput composes a vitals observation from data received
func (ps ServicePubSubMessaging) ComposeVitalsInput(ctx context.Context, input dto.CreateVitalSignPubSubMessage) (*domain.FHIRObservationInput, error) {
	vitalsConcept, err := getCIELConcept(ctx, ps.infra.OpenConceptLab, *input.ConceptID)
	if err != nil {
		return nil, err
	}

	system := "http://terminology.hl7.org/CodeSystem/observation-category"
	status := domain.ObservationStatusEnumPreliminary
	instant := scalarutils.Instant(input.Date.Format(time.RFC3339))
	observation := domain.FHIRObservationInput{
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
					System:  (*scalarutils.URI)(&vitalsConcept.URL),
					Code:    scalarutils.Code(vitalsConcept.ID),
					Display: vitalsConcept.DisplayName,
				},
			},
			Text: vitalsConcept.DisplayName,
		},
		ValueString: &input.Value,
	}

	patient, err := ps.clinical.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}
	patientReference := fmt.Sprintf("Patient/%v", patient.PatientRecord.ID)
	patientName := *patient.PatientRecord.Name[0].Given[0]
	observation.Subject = &domain.FHIRReferenceInput{
		Reference: &patientReference,
		Display:   patientName,
	}

	if input.OrganizationID != "" {
		organization, err := ps.clinical.FindOrganizationByID(ctx, input.OrganizationID)
		if err != nil {
			//Should not fail if organization is not found
			log.Printf("the error is: %v", err)
		}

		if organization != nil {
			performerReference := fmt.Sprintf("Organization/%v", input.OrganizationID)
			referenceInput := &domain.FHIRReferenceInput{
				Reference: &performerReference,
				Display:   *organization.Resource.Name,
			}

			observation.Performer = append(observation.Performer, referenceInput)
		}
	}

	return &observation, nil
}

// ComposeMedicationStatementInput composes a medication statement input from received data
func (ps ServicePubSubMessaging) ComposeMedicationStatementInput(ctx context.Context, input dto.CreateMedicationPubSubMessage) (*domain.FHIRMedicationStatementInput, error) {
	medicationConcept, err := getCIELConcept(ctx, ps.infra.OpenConceptLab, *input.ConceptID)
	if err != nil {
		return nil, err
	}

	drugConcept, err := getCIELConcept(ctx, ps.infra.OpenConceptLab, *input.Drug.ConceptID)
	if err != nil {
		return nil, err
	}

	year, month, day := input.Date.Date()
	status := domain.MedicationStatementStatusEnumUnknown
	medicationStatement := domain.FHIRMedicationStatementInput{
		Status: &status,
		Category: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&medicationConcept.URL),
					Code:    scalarutils.Code(medicationConcept.ID),
					Display: medicationConcept.DisplayName,
				},
			},
			Text: medicationConcept.DisplayName,
		},
		MedicationCodeableConcept: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&drugConcept.URL),
					Code:    scalarutils.Code(drugConcept.ID),
					Display: drugConcept.DisplayName,
				},
			},
			Text: drugConcept.DisplayName,
		},
		EffectiveDateTime: &scalarutils.Date{
			Year:  year,
			Month: int(month),
			Day:   day,
		},
	}

	patient, err := ps.clinical.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}
	patientReference := fmt.Sprintf("Patient/%v", patient.PatientRecord.ID)
	patientName := *patient.PatientRecord.Name[0].Given[0]
	medicationStatement.Subject = &domain.FHIRReferenceInput{
		Reference: &patientReference,
		Display:   patientName,
	}

	if input.OrganizationID != "" {
		organization, err := ps.clinical.FindOrganizationByID(ctx, input.OrganizationID) // rename organization response
		if err != nil {
			log.Printf("the error is: %v", err)
		}
		if organization != nil {
			informationSourceReference := fmt.Sprintf("Organization/%v", input.OrganizationID)

			reference := &domain.FHIRReferenceInput{
				Reference: &informationSourceReference,
				Display:   *organization.Resource.Name,
			}

			medicationStatement.InformationSource = reference
		}
	}

	return &medicationStatement, nil
}

// ComposeAllergyIntoleranceInput composes an allergy intolerance input from the data received
func (ps ServicePubSubMessaging) ComposeAllergyIntoleranceInput(ctx context.Context, input dto.CreatePatientAllergyPubSubMessage) (*domain.FHIRAllergyIntoleranceInput, error) {
	allergyType := domain.AllergyIntoleranceTypeEnumAllergy
	allergyCategory := domain.AllergyIntoleranceCategoryEnumMedication
	allergy := &domain.FHIRAllergyIntoleranceInput{
		Type:     &allergyType,
		Category: []*domain.AllergyIntoleranceCategoryEnum{&allergyCategory},
		ClinicalStatus: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&fhirAllergyIntoleranceClinicalStatusURL),
					Code:    scalarutils.Code("active"),
					Display: "Active",
				},
			},
			Text: "Active",
		},
		VerificationStatus: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&fhirAllergyIntoleranceVerificationStatusURL),
					Code:    scalarutils.Code("confirmed"),
					Display: "Confirmed",
				},
			},
			Text: "Confirmed",
		},
		Reaction: []*domain.FHIRAllergyintoleranceReactionInput{},
	}

	year, month, day := input.Date.Date()
	allergy.RecordedDate = &scalarutils.Date{
		Year:  year,
		Month: int(month),
		Day:   day,
	}

	patient, err := ps.clinical.FindPatientByID(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}
	subjectReference := fmt.Sprintf("Patient/%v", input.PatientID)
	patientName := *patient.PatientRecord.Name[0].Given[0]

	allergy.Patient = &domain.FHIRReferenceInput{
		Reference: &subjectReference,
		Display:   patientName,
	}

	allergenConcept, err := getCIELConcept(ctx, ps.infra.OpenConceptLab, *input.ConceptID)
	if err != nil {
		return nil, err
	}

	allergy.Code = domain.FHIRCodeableConceptInput{
		Coding: []*domain.FHIRCodingInput{
			{
				System:  (*scalarutils.URI)(&allergenConcept.URL),
				Code:    scalarutils.Code(allergenConcept.ID),
				Display: allergenConcept.DisplayName,
			},
		},
		Text: allergenConcept.DisplayName,
	}

	// create the allergy reaction
	var reaction *domain.FHIRAllergyintoleranceReactionInput

	// reaction manifestation is required
	//
	// check if there is a reaction manifestation,
	// if no reaction use unknown
	var manifestationConcept *domain.Concept
	if input.Reaction.ConceptID != nil {
		manifestationConcept, err = getCIELConcept(ctx, ps.infra.OpenConceptLab, *input.Reaction.ConceptID)
		if err != nil {
			return nil, err
		}

	} else {
		manifestationConcept, err = getCIELConcept(ctx, ps.infra.OpenConceptLab, unknownConceptID)
		if err != nil {
			return nil, err
		}

	}

	manifestation := &domain.FHIRCodeableConceptInput{
		Coding: []*domain.FHIRCodingInput{
			{
				System:  (*scalarutils.URI)(&manifestationConcept.URL),
				Code:    scalarutils.Code(manifestationConcept.ID),
				Display: manifestationConcept.DisplayName,
			},
		},
		Text: manifestationConcept.DisplayName,
	}

	// add reaction manifestation
	reaction.Manifestation = append(reaction.Manifestation, manifestation)

	if input.Severity.ConceptID != nil {
		severityConcept, err := getCIELConcept(ctx, ps.infra.OpenConceptLab, *input.Severity.ConceptID)
		if err != nil {
			return nil, err
		}

		reaction.Description = &severityConcept.DisplayName
	}

	// add allergy reaction
	allergy.Reaction = append(allergy.Reaction, reaction)

	return allergy, nil
}

func getCIELConcept(ctx context.Context, ocl openconceptlab.ServiceOCL, conceptID string) (*domain.Concept, error) {
	response, err := ocl.GetConcept(
		ctx,
		"CIEL",
		"CIEL",
		conceptID,
		false,
		false,
	)
	if err != nil {
		return nil, err
	}

	var concept *domain.Concept
	err = mapstructure.Decode(response, &concept)
	if err != nil {
		return nil, err
	}

	return concept, nil
}
