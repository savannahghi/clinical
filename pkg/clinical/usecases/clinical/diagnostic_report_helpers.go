package clinical

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// DiagnosticReportMutatorFunc is a helper function that is used by the "caller" of "RecordDiagnosticReport" to modify the diagnostic report data class model
// with the aapropriate data to suit its use case.
type DiagnosticReportMutatorFunc func(context.Context, *domain.FHIRDiagnosticReportInput) error

// addCytopathologyCategory is used to add a cytopathology category for various diagnostic reports such as biopsy etc.
var addCytopathologyCategory = func(ctx context.Context, diagnosticReport *domain.FHIRDiagnosticReportInput) error {
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

// addNuclearMagneticResonanceCategory is used to add a nuclear magnetic resonance category for diagnostic report such as MRI reports etc.
var addNuclearMagneticResonanceCategory = func(ctx context.Context, diagnosticReport *domain.FHIRDiagnosticReportInput) error {
	category := []*domain.FHIRCodeableConceptInput{
		{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&diagnosticReportCategorySystem),
					Code:    scalarutils.Code("NMR"),
					Display: "Nuclear Magnetic Resonance",
				},
			},
			Text: "Nuclear Magnetic Resonance",
		},
	}

	diagnosticReport.Category = append(diagnosticReport.Category, category...)

	return nil
}
