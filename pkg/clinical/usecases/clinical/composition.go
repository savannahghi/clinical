package clinical

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

func (c *UseCasesClinicalImpl) CreateComposition(ctx context.Context, input dto.CompositionInput) (*dto.Composition, error) {
	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	// check encounter status
	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot record a composition in a finished encounter")
	}

	// get patient from encounter
	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, *encounter.Resource.Subject.ID)
	if err != nil {
		return nil, err
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	patientType := "Patient"
	compositionTitle := fmt.Sprintf("%s's assessment note", patient.Resource.Name[0].Text)
	compositionSectionTextStatus := "generated"

	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	encounterType := scalarutils.URI("Encounter")

	organizationRef := fmt.Sprintf("Organization/%s", identifiers.OrganizationID)

	today := time.Now()

	date, err := scalarutils.NewDate(today.Day(), int(today.Month()), today.Year())
	if err != nil {
		return nil, err
	}

	ID := uuid.New().String()

	compositionCategoryConcept, err := c.GetConcept(ctx, dto.TerminologySourceLOINC, common.LOINCAssessmentPlanCode)
	if err != nil {
		return nil, err
	}

	compositionTypeConcept, err := c.GetConcept(ctx, dto.TerminologySourceLOINC, common.LOINCProgressNoteCode)
	if err != nil {
		return nil, err
	}

	status := strings.ToLower(string(input.Status))

	compositionInput := domain.FHIRCompositionInput{
		ID:     &ID,
		Status: (*domain.CompositionStatusEnum)(&status),
		Type: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&compositionTypeConcept.URL),
					Code:    scalarutils.Code(compositionTypeConcept.ID),
					Display: compositionTypeConcept.DisplayName,
				},
			},
			Text: compositionTypeConcept.DisplayName,
		},
		Category: []*domain.FHIRCodeableConceptInput{
			{
				ID: &ID,
				Coding: []*domain.FHIRCodingInput{
					{
						System:  (*scalarutils.URI)(&compositionCategoryConcept.URL),
						Code:    scalarutils.Code(compositionCategoryConcept.ID),
						Display: compositionCategoryConcept.DisplayName,
					},
				},
				Text: compositionCategoryConcept.DisplayName,
			},
		},
		Subject: &domain.FHIRReferenceInput{
			ID:        patient.Resource.ID,
			Reference: &patientRef,
			Type:      (*scalarutils.URI)(&patientType),
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.ID,
			Reference: &encounterRef,
			Display:   *encounter.Resource.ID,
			Type:      &encounterType,
		},
		Date: date,
		Author: []*domain.FHIRReferenceInput{
			{
				Reference: &organizationRef,
			},
		},
		Title: &compositionTitle,
		Section: []*domain.FHIRCompositionSectionInput{
			{
				ID:    &ID,
				Title: &compositionCategoryConcept.DisplayName,
				Code: &domain.FHIRCodeableConceptInput{
					ID: new(string),
					Coding: []*domain.FHIRCodingInput{
						{
							ID:      &ID,
							System:  (*scalarutils.URI)(&compositionCategoryConcept.URL),
							Code:    scalarutils.Code(compositionCategoryConcept.ID),
							Display: compositionCategoryConcept.DisplayName,
						},
					},
					Text: compositionTypeConcept.DisplayName,
				},
				Author: []*domain.FHIRReferenceInput{
					{
						Reference: &organizationRef,
					},
				},
				Text: &domain.FHIRNarrativeInput{
					ID:     new(string),
					Status: (*domain.NarrativeStatusEnum)(&compositionSectionTextStatus),
					Div:    scalarutils.XHTML(input.Note),
				},
			},
		},
	}

	composition, err := c.infrastructure.FHIR.CreateFHIRComposition(ctx, compositionInput)
	if err != nil {
		return nil, err
	}

	return mapFHIRCompositionToCompositionDTO(*composition.Resource), nil
}

func mapFHIRCompositionToCompositionDTO(composition domain.FHIRComposition) *dto.Composition {
	output := dto.Composition{
		ID:          *composition.ID,
		Text:        string(composition.Section[0].Text.Div),
		Type:        dto.CompositionType(composition.Type.Text),
		Category:    dto.CompositionCategory(composition.Category[0].Text),
		Status:      dto.CompositionStatusEnum(*composition.Status),
		PatientID:   *composition.Subject.ID,
		EncounterID: *composition.Encounter.ID,
		Date:        composition.Date,
	}

	return &output
}
