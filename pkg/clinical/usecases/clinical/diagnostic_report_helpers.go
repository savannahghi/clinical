package clinical

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// DiagnosticReportMutatorFunc is a helper function that is used by the "caller" of "RecordDiagnosticReport" to modify the diagnostic report data class model
// with the aapropriate data to suit its use case.
type DiagnosticReportMutatorFunc func(context.Context, *domain.FHIRDiagnosticReportInput) error

// addProcedureCategory is used to add procedure category for various observations records such as biopsy etc.
var addCytologyCategory = func(ctx context.Context, diagnosticReport *domain.FHIRDiagnosticReportInput) error {
	category := []*domain.FHIRCodeableConceptInput{
		{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&diagnosticReportCategorySystem),
					Code:    scalarutils.Code("CP"),
					Display: "Cytopathology",
				},
			},
			Text: "Cytopathology",
		},
	}

	diagnosticReport.Category = append(diagnosticReport.Category, category...)

	return nil
}
