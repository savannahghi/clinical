package clinical

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	allergyIntoleranceInput.Meta = domain.FHIRMetaInput{
		Tag: tags,
	}

	allergyIntolerance, err := c.infrastructure.FHIR.CreateFHIRAllergyIntolerance(ctx, allergyIntoleranceInput)
	if err != nil {
		return nil, err
	}

	allergyIntoleranceObj := mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(*allergyIntolerance.Resource)
	allergyIntoleranceObj.TerminologySource = input.TerminologySource

	return allergyIntoleranceObj, nil
}

// SearchAllergy is used to retrieve allergy from OCL
func (c *UseCasesClinicalImpl) SearchAllergy(ctx context.Context, name string, pagination dto.Pagination) (*dto.TerminologyConnection, error) {
	err := pagination.Validate()
	if err != nil {
		return nil, err
	}

	conceptPage, err := c.infrastructure.OpenConceptLab.
		ListConcepts(ctx, string(dto.TerminologySourceCIEL), string(dto.TerminologySourceCIEL), true, &name, nil, nil, nil, nil, nil, nil, nil, nil, &pagination)
	if err != nil {
		return nil, err
	}

	terminologyPage := &dto.TerminologyConnection{
		TotalCount: conceptPage.Count,
		Edges:      []dto.TerminologyEdge{},
		PageInfo:   dto.PageInfo{},
	}

	if conceptPage.Next != nil {
		terminologyPage.PageInfo.HasNextPage = true
		terminologyPage.PageInfo.StartCursor = conceptPage.Next
	}

	if conceptPage.Previous != nil {
		terminologyPage.PageInfo.HasPreviousPage = true
		terminologyPage.PageInfo.EndCursor = conceptPage.Previous
	}

	for _, concept := range conceptPage.Results {
		terminologyPage.Edges = append(terminologyPage.Edges, dto.TerminologyEdge{
			Node: dto.Terminology{
				Code:   concept.ID,
				System: dto.TerminologySource(concept.Source),
				Name:   concept.DisplayName,
			},
		})
	}

	return terminologyPage, nil
}

// GetAllergyIntolerance fetches all the allergy intolerance from FHIR by allergy intolerance ID
func (c *UseCasesClinicalImpl) GetAllergyIntolerance(ctx context.Context, id string) (*dto.Allergy, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid allergy intolerance id: %s", id)
	}

	allergyIntolerance, err := c.infrastructure.FHIR.GetFHIRAllergyIntolerance(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to search for allergy intolerance: %w", err)
	}

	intolerance := mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(*allergyIntolerance.Resource)

	return intolerance, nil
}

// ListPatientAllergies is used to list all allergies associated with a specific patient
func (c *UseCasesClinicalImpl) ListPatientAllergies(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.AllergyConnection, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
	}

	err = pagination.Validate()
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", patientID)

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	allergyResponses, err := c.infrastructure.FHIR.SearchPatientAllergyIntolerance(ctx, patientReference, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	patientAllergyIntolerances := []*dto.Allergy{}

	for _, allergyResponse := range allergyResponses.Allergies {
		patientAllergyIntolerances = append(patientAllergyIntolerances, mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(allergyResponse))
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     allergyResponses.HasNextPage,
		EndCursor:       &allergyResponses.NextCursor,
		HasPreviousPage: allergyResponses.HasPreviousPage,
		StartCursor:     &allergyResponses.PreviousCursor,
	}

	connection := dto.CreateAllergyConnection(patientAllergyIntolerances, pageInfo, allergyResponses.TotalCount)

	return &connection, nil
}
