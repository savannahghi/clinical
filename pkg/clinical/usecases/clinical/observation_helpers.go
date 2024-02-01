package clinical

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// NOTE: This file contains simplified helper functions to map observation categories to their FHIR observation categories

var (
	observationCategorySystem = "http://terminology.hl7.org/CodeSystem/observation-category"
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
