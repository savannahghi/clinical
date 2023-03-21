package clinical

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
)

// UseCasesClinicalImpl represents the patient usecase implementation
type UseCasesClinicalImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewUseCasesClinicalImpl initializes new Clinical/Patient implementation
func NewUseCasesClinicalImpl(infra infrastructure.Infrastructure) *UseCasesClinicalImpl {
	return &UseCasesClinicalImpl{
		infrastructure: infra,
	}
}

// GetMedicalData returns a limited subset of specific medical data that for a specific patient
// These include: Allergies, Viral Load, Body Mass Index, Weight, CD4 Count using their respective OCL CIEL Terminology
// For each category the latest three records are fetched
func (c *UseCasesClinicalImpl) GetMedicalData(ctx context.Context, patientID string) (*dto.MedicalData, error) {
	data := &dto.MedicalData{}

	filterParams := map[string]interface{}{
		"patient": fmt.Sprintf("Patient/%v", patientID),
		"_count":  common.MedicalDataCount,
		"_sort":   "-date",
	}

	fields := []string{
		"Regimen",
		"AllergyIntolerance",
		"Weight",
		"BMI",
		"ViralLoad",
		"CD4Count",
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	for _, field := range fields {
		switch field {
		case "Regimen":
			conn, err := c.infrastructure.FHIR.SearchFHIRMedicationStatement(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				if edge.Node.ID == nil {
					continue
				}

				if edge.Node.Status == nil {
					continue
				}

				if edge.Node.MedicationCodeableConcept.Coding == nil {
					continue
				}

				if len(edge.Node.MedicationCodeableConcept.Coding) < 1 {
					continue
				}

				if edge.Node.Subject.ID == nil {
					continue
				}

				data.Regimen = append(data.Regimen, mapFHIRMedicationStatementToMedicationStatementDTO(edge.Node))
			}
		case "AllergyIntolerance":
			conn, err := c.infrastructure.FHIR.SearchFHIRAllergyIntolerance(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				if edge.Node.ID == nil {
					continue
				}

				if edge.Node.Patient == nil {
					continue
				}

				if edge.Node.Patient.ID == nil {
					continue
				}

				if edge.Node.Encounter == nil {
					continue
				}

				if edge.Node.Encounter.ID == nil {
					continue
				}

				if edge.Node.Code == nil {
					continue
				}

				if edge.Node.Code.Coding == nil {
					continue
				}

				if len(edge.Node.Code.Coding) < 1 {
					continue
				}

				data.Allergies = append(data.Allergies, mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(edge.Node))
			}

		case "Weight":
			filterParams["code"] = common.WeightCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				if edge.Node.Code.Coding == nil {
					continue
				}

				if len(edge.Node.Code.Coding) < 1 {
					continue
				}

				if edge.Node.Status == nil {
					continue
				}

				data.Weight = append(data.Weight, mapFHIRObservationToObservationDTO(edge.Node))
			}

		case "BMI":
			filterParams["code"] = common.BMICIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				if edge.Node.Code.Coding == nil {
					continue
				}

				if len(edge.Node.Code.Coding) < 1 {
					continue
				}

				if edge.Node.Status == nil {
					continue
				}

				data.BMI = append(data.BMI, mapFHIRObservationToObservationDTO(edge.Node))
			}

		case "ViralLoad":
			filterParams["code"] = common.ViralLoadCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				if edge.Node.Code.Coding == nil {
					continue
				}

				if len(edge.Node.Code.Coding) < 1 {
					continue
				}

				if edge.Node.Status == nil {
					continue
				}

				data.ViralLoad = append(data.ViralLoad, mapFHIRObservationToObservationDTO(edge.Node))
			}

		case "CD4Count":
			filterParams["code"] = common.CD4CountCIELTerminologyCode

			conn, err := c.infrastructure.FHIR.SearchFHIRObservation(ctx, filterParams, *identifiers)
			if err != nil {
				utils.ReportErrorToSentry(err)
				return nil, fmt.Errorf("%s search error: %w", field, err)
			}

			for _, edge := range conn.Edges {
				if edge.Node == nil {
					continue
				}

				if edge.Node.Code.Coding == nil {
					continue
				}

				if len(edge.Node.Code.Coding) < 1 {
					continue
				}

				if edge.Node.Status == nil {
					continue
				}

				data.CD4Count = append(data.CD4Count, mapFHIRObservationToObservationDTO(edge.Node))
			}
		}
	}

	return data, nil
}

func mapFHIRMedicationStatementToMedicationStatementDTO(fhirAllergyIntolerance *domain.FHIRMedicationStatement) *dto.MedicationStatement {
	return &dto.MedicationStatement{
		ID:     *fhirAllergyIntolerance.ID,
		Status: dto.MedicationStatementStatusEnum(*fhirAllergyIntolerance.Status),
		Medication: dto.Medication{
			Name: fhirAllergyIntolerance.MedicationCodeableConcept.Coding[0].Display,
			Code: string(fhirAllergyIntolerance.MedicationCodeableConcept.Coding[0].Code),
		},
		PatientID: *fhirAllergyIntolerance.Subject.ID,
	}
}

func mapFHIRAllergyIntoleranceToAllergyIntoleranceDTO(fhirAllergyIntolerance *domain.FHIRAllergyIntolerance) *dto.AllergyIntolerance {
	return &dto.AllergyIntolerance{
		ID:              *fhirAllergyIntolerance.ID,
		PatientID:       *fhirAllergyIntolerance.Patient.ID,
		EncounterID:     *fhirAllergyIntolerance.Encounter.ID,
		OnsetDateTime:   fhirAllergyIntolerance.OnsetPeriod.Start,
		Severity:        dto.AllergyIntoleranceReactionSeverityEnum(fhirAllergyIntolerance.Criticality),
		SubstanceCode:   string(fhirAllergyIntolerance.Code.Coding[0].Code),
		SubstanceSystem: string(*fhirAllergyIntolerance.Code.Coding[0].System),
	}
}

func mapFHIRObservationToObservationDTO(fhirObservation *domain.FHIRObservation) *dto.Observation {
	var value string

	if fhirObservation.ValueQuantity != nil {
		value = fmt.Sprintf("%v %v", fhirObservation.ValueQuantity.Value, fhirObservation.ValueQuantity.Unit)
	}

	if fhirObservation.ValueCodeableConcept != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueCodeableConcept)
	}

	if fhirObservation.ValueString != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueString)
	}

	if fhirObservation.ValueBoolean != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueBoolean)
	}

	if fhirObservation.ValueInteger != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueInteger)
	}

	if fhirObservation.ValueRange != nil {
		value = fmt.Sprintf("%v %v - %v %v", fhirObservation.ValueRange.High.Value, fhirObservation.ValueRange.High.Unit, fhirObservation.ValueRange.Low.Value, fhirObservation.ValueRange.Low.Unit)
	}

	if fhirObservation.ValueRatio != nil {
		value = fmt.Sprintf("%v %v : %v %v", fhirObservation.ValueRatio.Numerator.Value, fhirObservation.ValueRatio.Numerator.Unit, fhirObservation.ValueRatio.Denominator.Value, fhirObservation.ValueRatio.Denominator.Unit)
	}

	if fhirObservation.ValueSampledData != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueSampledData.ID)
	}

	if fhirObservation.ValueTime != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueTime)
	}

	if fhirObservation.ValueDateTime != nil {
		value = fmt.Sprintf("%v", *fhirObservation.ValueDateTime)
	}

	if fhirObservation.ValuePeriod != nil {
		value = fmt.Sprintf("%v - %v", fhirObservation.ValuePeriod.Start, fhirObservation.ValuePeriod.End)
	}

	return &dto.Observation{
		ID:          *fhirObservation.ID,
		Status:      dto.ObservationStatus(*fhirObservation.Status),
		Name:        fhirObservation.Code.Coding[0].Display,
		Value:       value,
		PatientID:   *fhirObservation.Subject.ID,
		EncounterID: *fhirObservation.Encounter.ID,
	}
}

// CreateFHIROrganization creates a FHIROrganization instance
func (c *UseCasesClinicalImpl) CreateFHIROrganization(ctx context.Context, input domain.FHIROrganizationInput) (*domain.FHIROrganizationRelayPayload, error) {
	organizationRelayPayload, err := c.infrastructure.FHIR.CreateFHIROrganization(ctx, input)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return organizationRelayPayload, nil
}
