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

// CreateComposition creates a new composition
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

	id := uuid.New().String()

	var compositionCategoryCode string

	switch input.Category {
	case "ASSESSMENT_PLAN":
		compositionCategoryCode = common.LOINCAssessmentPlanCode
	case "HISTORY_OF_PRESENTING_ILLNESS":
		compositionCategoryCode = common.LOINCHistoryOfPresentingIllness
	case "SOCIAL_HISTORY":
		compositionCategoryCode = common.LOINCSocialHistory
	case "FAMILY_HISTORY":
		compositionCategoryCode = common.LOINCFamilyHistory
	case "EXAMINATION":
		compositionCategoryCode = common.LOINCExamination
	case "PLAN_OF_CARE":
		compositionCategoryCode = common.LOINCPLANOFCARE
	default:
		return nil, fmt.Errorf("category is needed")
	}

	compositionCategoryConcept, err := c.GetConcept(ctx, dto.TerminologySourceLOINC, compositionCategoryCode)
	if err != nil {
		return nil, err
	}

	compositionTypeConcept, err := c.GetConcept(ctx, dto.TerminologySourceLOINC, common.LOINCProgressNoteCode)
	if err != nil {
		return nil, err
	}

	status := strings.ToLower(string(input.Status))

	compositionInput := domain.FHIRCompositionInput{
		ID:     &id,
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
				ID: &id,
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
				ID:    &id,
				Title: &compositionCategoryConcept.DisplayName,
				Code: &domain.FHIRCodeableConceptInput{
					ID: &id,
					Coding: []*domain.FHIRCodingInput{
						{
							ID:      &id,
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
					ID:     &id,
					Status: (*domain.NarrativeStatusEnum)(&compositionSectionTextStatus),
					Div:    scalarutils.XHTML(input.Note),
				},
			},
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	compositionInput.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	composition, err := c.infrastructure.FHIR.CreateFHIRComposition(ctx, compositionInput)
	if err != nil {
		return nil, err
	}

	compositionInput.Category = []*domain.FHIRCodeableConceptInput{
		{
			ID: &id,
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&compositionCategoryConcept.URL),
					Code:    scalarutils.Code(compositionCategoryConcept.ID),
					Display: compositionCategoryConcept.DisplayName,
				},
			},
			Text: compositionCategoryConcept.DisplayName,
		},
	}

	result := mapFHIRCompositionToCompositionDTO(*composition.Resource)

	return &result.Edges[0].Node, nil
}

func mapFHIRCompositionToCompositionDTO(composition domain.FHIRComposition) *dto.CompositionConnection {
	var compositionSection []*dto.Section

	for _, item := range composition.Section {
		var itemSubSections []*dto.Section

		if len(item.Section) > 0 {
			for _, section := range item.Section {
				itemSubSections = append(itemSubSections, &dto.Section{
					ID:     section.ID,
					Title:  section.Title,
					Code:   section.Code.ID,
					Author: section.Author[0].Reference,
					Text:   string(section.Text.Div),
				})
			}
		}

		compositionSection = append(compositionSection, &dto.Section{
			ID:      item.ID,
			Title:   item.Title,
			Code:    &item.Code.Coding[0].Display,
			Author:  item.Author[0].Reference,
			Text:    string(item.Text.Div),
			Section: itemSubSections,
		})
	}

	output := dto.Composition{
		ID:          *composition.ID,
		Text:        string(composition.Section[0].Text.Div),
		Type:        dto.CompositionType(composition.Type.Text),
		Category:    dto.CompositionCategory(composition.Category[0].Text),
		Status:      dto.CompositionStatusEnum(*composition.Status),
		PatientID:   *composition.Subject.ID,
		EncounterID: *composition.Encounter.ID,
		Date:        composition.Date,
		Section:     compositionSection,
	}

	return &dto.CompositionConnection{
		TotalCount: 0,
		Edges: []dto.CompositionEdge{
			{
				Node:   output,
				Cursor: "",
			},
		},
		PageInfo: dto.PageInfo{},
	}
}

// ListPatientCompositions lists a patient's compositions
func (c UseCasesClinicalImpl) ListPatientCompositions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.CompositionConnection, error) {
	_, err := uuid.Parse(patientID)
	if err != nil {
		return nil, fmt.Errorf("invalid patient id: %s", patientID)
	}

	err = pagination.Validate()
	if err != nil {
		return nil, err
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	patientRef := fmt.Sprintf("Patient/%s", *patient.Resource.ID)
	params := map[string]interface{}{
		"subject": patientRef,
		"_sort":   "date",
	}

	if encounterID != nil {
		encounterReference := fmt.Sprintf("Encounter/%s", *encounterID)
		params["encounter"] = encounterReference
	}

	if date != nil {
		params["date"] = date.AsTime().Format(dateFormatStr)
	}

	compositionsResponse, err := c.infrastructure.FHIR.SearchFHIRComposition(ctx, params, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	compositions := []dto.Composition{}

	for _, resource := range compositionsResponse.Compositions {
		composition := mapFHIRCompositionToCompositionDTO(resource)
		compositions = append(compositions, composition.Edges[0].Node)
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     compositionsResponse.HasNextPage,
		EndCursor:       &compositionsResponse.NextCursor,
		HasPreviousPage: compositionsResponse.HasPreviousPage,
		StartCursor:     &compositionsResponse.PreviousCursor,
	}

	connection := dto.CreateCompositionConnection(compositions, pageInfo, compositionsResponse.TotalCount)

	return &connection, nil
}

// AppendNoteToComposition appends a note to the patient's composition information such as section
func (c *UseCasesClinicalImpl) AppendNoteToComposition(ctx context.Context, id string, input dto.PatchCompositionInput) (*dto.Composition, error) {
	if id == "" {
		return nil, fmt.Errorf("a composition id is required")
	}

	composition, err := c.infrastructure.FHIR.GetFHIRComposition(ctx, id)
	if err != nil {
		return nil, err
	}

	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, *composition.Resource.Encounter.ID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot record a composition in a finished encounter")
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	organizationRef := fmt.Sprintf("Organization/%s", identifiers.OrganizationID)

	var compositionCategoryCode string

	switch input.Category {
	case "ASSESSMENT_PLAN":
		compositionCategoryCode = common.LOINCAssessmentPlanCode
	case "HISTORY_OF_PRESENTING_ILLNESS":
		compositionCategoryCode = common.LOINCHistoryOfPresentingIllness
	case "SOCIAL_HISTORY":
		compositionCategoryCode = common.LOINCSocialHistory
	case "FAMILY_HISTORY":
		compositionCategoryCode = common.LOINCFamilyHistory
	case "EXAMINATION":
		compositionCategoryCode = common.LOINCExamination
	case "PLAN_OF_CARE":
		compositionCategoryCode = common.LOINCPLANOFCARE
	}

	compositionCategoryConcept, err := c.GetConcept(ctx, dto.TerminologySourceLOINC, compositionCategoryCode)
	if err != nil {
		return nil, err
	}

	idd := uuid.New().String()
	compositionSectionTextStatus := "generated"

	compositionSection := &domain.FHIRCompositionSection{
		ID:    &idd,
		Title: &compositionCategoryConcept.DisplayName,
		Code: &domain.FHIRCodeableConceptInput{
			ID: &idd,
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&compositionCategoryConcept.URL),
					Code:    scalarutils.Code(compositionCategoryConcept.ID),
					Display: compositionCategoryConcept.DisplayName,
				},
			},
			Text: compositionCategoryConcept.DisplayName,
		},
		Author: []*domain.FHIRReference{
			{
				Reference: &organizationRef,
			},
		},
		Focus: &domain.FHIRReference{},
		Text: &domain.FHIRNarrative{
			ID:     &idd,
			Status: (*domain.NarrativeStatusEnum)(&compositionSectionTextStatus),
			Div:    scalarutils.XHTML(input.Note),
		},
	}

	composition.Resource.Section = append(composition.Resource.Section, compositionSection)

	var sectionInput []*domain.FHIRCompositionSectionInput

	for _, s := range composition.Resource.Section {
		sectionInput = append(sectionInput, &domain.FHIRCompositionSectionInput{
			ID:    s.ID,
			Title: s.Title,
			Code: &domain.FHIRCodeableConceptInput{
				ID:     s.Code.ID,
				Coding: s.Code.Coding,
				Text:   compositionSectionTextStatus,
			},
			Author: []*domain.FHIRReferenceInput{
				{
					Reference: &organizationRef,
				},
			},
			Text: &domain.FHIRNarrativeInput{
				ID:     s.Text.ID,
				Status: s.Text.Status,
				Div:    s.Text.Div,
			},
		})

		if len(s.Section) > 0 {
			var nestedsectionInput []*domain.FHIRCompositionSectionInput

			for _, r := range s.Section {
				nestedsectionInput = append(nestedsectionInput, &domain.FHIRCompositionSectionInput{
					ID:    r.ID,
					Title: r.Title,
					Code: &domain.FHIRCodeableConceptInput{
						ID:     r.Code.ID,
						Coding: r.Code.Coding,
						Text:   compositionSectionTextStatus,
					},
					Author: []*domain.FHIRReferenceInput{
						{
							Reference: &organizationRef,
						},
					},
					Text: &domain.FHIRNarrativeInput{
						ID:     r.Text.ID,
						Status: r.Text.Status,
						Div:    r.Text.Div,
					},
				})
			}

			for _, m := range sectionInput {
				m.Section = append(m.Section, nestedsectionInput...)
			}
		}
	}

	compositionInput := &domain.FHIRCompositionInput{
		Section: sectionInput,
	}

	output, err := c.infrastructure.FHIR.PatchFHIRComposition(ctx, id, *compositionInput)
	if err != nil {
		return nil, err
	}

	result := mapFHIRCompositionToCompositionDTO(*output)

	return &result.Edges[0].Node, nil
}
