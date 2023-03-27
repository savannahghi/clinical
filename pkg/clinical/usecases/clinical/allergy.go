package clinical

import (
	"context"
	"time"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// CreateAllergyIntolerance creates a new allergy intolerance
func (c *UseCasesClinicalImpl) CreateAllergyIntolerance(ctx context.Context, input dto.AllergyInput) (*dto.Allergy, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, *encounter.Resource.Subject.ID)
	if err != nil {
		return nil, err
	}

	allergyConcept, err := c.GetConcept(ctx, input.TerminologySource, input.Code)
	if err != nil {
		return nil, err
	}

	clinicalStatusSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical")
	verificationSystem := "http://terminology.hl7.org/CodeSystem/allergyintolerance-verification"

	allergyIntoleranceTypeAllergy := domain.AllergyIntoleranceTypeEnumAllergy

	clinicalStatusCodeActive := "active"
	verificationDisplay := "confirmed"

	allergyIntoleranceInput := domain.FHIRAllergyIntoleranceInput{
		ClinicalStatus: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{{
				System:  &clinicalStatusSystem,
				Code:    scalarutils.Code(clinicalStatusCodeActive),
				Display: clinicalStatusCodeActive,
			}},
		},
		Code: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&allergyConcept.URL),
					Code:    scalarutils.Code(allergyConcept.ID),
					Display: allergyConcept.DisplayName,
				},
			},
			Text: allergyConcept.DisplayName,
		},
		Patient: &domain.FHIRReferenceInput{
			ID:        patient.Resource.ID,
			Reference: encounter.Resource.Subject.Reference,
			Type:      encounter.Resource.Subject.Type,
			Display:   patient.Resource.Names(),
		},

		Encounter: &domain.FHIRReferenceInput{
			ID: encounter.Resource.ID,
		},
		RecordedDate: &scalarutils.Date{
			Year:  time.Now().Year(),
			Month: int(time.Now().Month()),
			Day:   time.Now().Day(),
		},
		Type: &allergyIntoleranceTypeAllergy,
		VerificationStatus: domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&verificationSystem),
					Code:    scalarutils.Code(verificationDisplay),
					Display: verificationDisplay,
				},
			},
			Text: verificationDisplay,
		},
	}

	if input.Reaction != nil {
		manifestationConcept, err := c.GetConcept(ctx, dto.TerminologySourceCIEL, input.Reaction.Code)
		if err != nil {
			return nil, err
		}

		allergyIntoleranceInput.Reaction = []*domain.FHIRAllergyintoleranceReactionInput{{
			Description: (*string)(&input.Reaction.Severity),
			Manifestation: []*domain.FHIRCodeableConceptInput{{
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&manifestationConcept.URL),
						Code:    scalarutils.Code(manifestationConcept.ID),
						Display: manifestationConcept.DisplayName,
					},
				},
				Text: manifestationConcept.DisplayName,
			}},
			Severity: (*domain.AllergyIntoleranceReactionSeverityEnum)(&input.Reaction.Severity),
		}}
	}

	allergyIntolerance, err := c.infrastructure.FHIR.CreateFHIRAllergyIntolerance(ctx, allergyIntoleranceInput)
	if err != nil {
		return nil, err
	}

	allergyIntoleranceObj := mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(allergyIntolerance.Resource)
	allergyIntoleranceObj.TerminologySource = input.TerminologySource

	return allergyIntoleranceObj, nil
}
