package clinical

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

var (
	observationCategorySystem      = "http://terminology.hl7.org/CodeSystem/observation-category"
	diagnosticReportCategorySystem = "http://terminology.hl7.org/CodeSystem/v2-0074"
)

// addVitalSignCategory is used to add laboratory categories for various observations records.
var addVitalSignCategory = func(ctx context.Context, observation *domain.FHIRObservationInput) error {
	userSelected := false
	category := []*domain.FHIRCodeableConceptInput{
		{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       (*scalarutils.URI)(&observationCategorySystem),
					Code:         "vital-signs",
					Display:      "Vital Signs",
					UserSelected: &userSelected,
				},
			},
			Text: "Vital Signs",
		},
	}

	observation.Category = append(observation.Category, category...)
	return nil
}

// addLabCategory is used to add laboratory categories for various observations records.
var addLabCategory = func(ctx context.Context, observation *domain.FHIRObservationInput) error {
	userSelected := false
	category := []*domain.FHIRCodeableConceptInput{
		{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       (*scalarutils.URI)(&observationCategorySystem),
					Code:         "laboratory",
					Display:      "Laboratory",
					UserSelected: &userSelected,
				},
			},
			Text: "Laboratory",
		},
	}

	observation.Category = append(observation.Category, category...)
	return nil
}

// addImagingCategory is used to add imaging categories for various observations records.
var addImagingCategory = func(ctx context.Context, observation *domain.FHIRObservationInput) error {
	userSelected := false
	category := []*domain.FHIRCodeableConceptInput{
		{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       (*scalarutils.URI)(&observationCategorySystem),
					Code:         "imaging",
					Display:      "Imaging",
					UserSelected: &userSelected,
				},
			},
			Text: "Imaging",
		},
	}

	observation.Category = append(observation.Category, category...)
	return nil
}

// addProcedureCategory is used to add procedure category for various observations records such as biopsy etc.
var addProcedureCategory = func(ctx context.Context, observation *domain.FHIRObservationInput) error {
	userSelected := false
	category := []*domain.FHIRCodeableConceptInput{
		{
			Coding: []*domain.FHIRCodingInput{
				{
					System:       (*scalarutils.URI)(&observationCategorySystem),
					Code:         "procedure",
					Display:      "Procedure",
					UserSelected: &userSelected,
				},
			},
			Text: "Procedure",
		},
	}

	observation.Category = append(observation.Category, category...)
	return nil
}
