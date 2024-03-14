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

// addObservationCategory changes category of an observation
func addObservationCategory(code string) ObservationInputMutatorFunc {
	var display string

	switch code {
	case "exam":
		display = "Exam"
	case "procedure":
		display = "Procedure"
	case "imaging":
		display = "Imaging"
	case "laboratory":
		display = "Laboratory"
	case "vital-signs":
		display = "Vital Signs"
	}

	return func(ctx context.Context, observation *domain.FHIRObservationInput) error {
		userSelected := false
		category := []*domain.FHIRCodeableConceptInput{
			{
				Coding: []*domain.FHIRCodingInput{
					{
						System:       (*scalarutils.URI)(&observationCategorySystem),
						Code:         scalarutils.Code(code),
						Display:      display,
						UserSelected: &userSelected,
					},
				},
				Text: display,
			},
		}

		observation.Category = append(observation.Category, category...)

		return nil
	}
}
