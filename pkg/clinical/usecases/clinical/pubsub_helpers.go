package clinical

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// GetConcept is a helper function that returns a concept associated the terminology source passed
func (c *UseCasesClinicalImpl) GetConcept(ctx context.Context, terminologySource dto.TerminologySource, conceptID string) (*domain.Concept, error) {
	var (
		organisation string
		source       string
	)

	switch terminologySource {
	case dto.TerminologySourceICD10:
		organisation = "WHO"
		source = "ICD-10-WHO"

	case dto.TerminologySourceCIEL:
		organisation = "CIEL"
		source = "CIEL"

	case dto.TerminologySourceLOINC:
		organisation = "Regenstrief"
		source = "LOINC"

	case dto.TerminologySourceSNOMEDCT:
		organisation = "Sofya"
		source = "SNOMED-CT"

	default:
		return nil, fmt.Errorf("terminology source %v not supported", source)
	}

	response, err := c.infrastructure.OpenConceptLab.GetConcept(
		ctx,
		organisation,
		source,
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

// ComposeVitalsInput composes a vitals observation from data received
func (c *UseCasesClinicalImpl) ComposeVitalsInput(ctx context.Context, input dto.VitalSignPubSubMessage) (*domain.FHIRObservationInput, error) {
	vitalsConcept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, *input.ConceptID)
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

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%v", patient.Resource.ID)
	patientName := *patient.Resource.Name[0].Given[0]
	observation.Subject = &domain.FHIRReferenceInput{
		ID:        patient.Resource.ID,
		Reference: &patientReference,
		Display:   patientName,
	}

	if input.OrganizationID != "" {
		organization, err := c.infrastructure.FHIR.GetFHIROrganization(ctx, input.OrganizationID)
		if err != nil {
			// Should not fail if organization is not found
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

// ComposeAllergyIntoleranceInput composes an allergy intolerance input from the data received
func (c *UseCasesClinicalImpl) ComposeAllergyIntoleranceInput(ctx context.Context, input dto.PatientAllergyPubSubMessage) (*domain.FHIRAllergyIntoleranceInput, error) {
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

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	subjectReference := fmt.Sprintf("Patient/%v", input.PatientID)
	patientName := *patient.Resource.Name[0].Given[0]

	allergy.Patient = &domain.FHIRReferenceInput{
		ID:        patient.Resource.ID,
		Reference: &subjectReference,
		Display:   patientName,
	}

	allergenConcept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, *input.ConceptID)
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
	var reaction domain.FHIRAllergyintoleranceReactionInput

	// reaction manifestation is required
	//
	// check if there is a reaction manifestation,
	// if no reaction use unknown
	var manifestationConcept *domain.Concept
	if input.Reaction.ConceptID != nil {
		manifestationConcept, err = c.GetConcept(ctx, dto.TerminologySourceCIEL, *input.Reaction.ConceptID)
		if err != nil {
			return nil, err
		}
	} else {
		manifestationConcept, err = c.GetConcept(ctx, dto.TerminologySourceCIEL, unknownConceptID)
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
		severityConcept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, *input.Severity.ConceptID)
		if err != nil {
			return nil, err
		}

		reaction.Description = &severityConcept.DisplayName
	}

	// add allergy reaction
	allergy.Reaction = append(allergy.Reaction, &reaction)

	return allergy, nil
}

// ComposeTestResultInput composes a test result input from data received
func (c *UseCasesClinicalImpl) ComposeTestResultInput(ctx context.Context, input dto.PatientTestResultPubSubMessage) (*domain.FHIRObservationInput, error) {
	var patientName string

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	patientName = *patient.Resource.Name[0].Given[0]

	observationConcept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, *input.ConceptID)
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
			ID:        patient.Resource.ID,
			Reference: &subjectReference,
			Display:   patientName,
		},
	}

	if input.OrganizationID != "" {
		organization, err := c.infrastructure.FHIR.GetFHIROrganization(ctx, input.OrganizationID) // rename organization response
		if err != nil {
			// Should not fail if the organization is not found
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

// ComposeMedicationStatementInput composes a medication statement input from received data
func (c *UseCasesClinicalImpl) ComposeMedicationStatementInput(ctx context.Context, input dto.MedicationPubSubMessage) (*domain.FHIRMedicationStatementInput, error) {
	medicationConcept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, *input.ConceptID)
	if err != nil {
		return nil, err
	}

	drugConcept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, *input.Drug.ConceptID)
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

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%v", patient.Resource.ID)
	patientName := *patient.Resource.Name[0].Given[0]
	medicationStatement.Subject = &domain.FHIRReferenceInput{
		ID:        patient.Resource.ID,
		Reference: &patientReference,
		Display:   patientName,
	}

	if input.OrganizationID != "" {
		organization, err := c.infrastructure.FHIR.GetFHIROrganization(ctx, input.OrganizationID) // rename organization response
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
