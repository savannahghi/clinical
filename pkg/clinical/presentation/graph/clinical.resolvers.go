package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph/generated"
	"github.com/savannahghi/scalarutils"
)

// CreateEpisodeOfCare is the resolver for the createEpisodeOfCare field.
func (r *mutationResolver) CreateEpisodeOfCare(ctx context.Context, episodeOfCare dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error) {
	r.CheckDependencies()
	return r.usecases.CreateEpisodeOfCare(ctx, episodeOfCare)
}

// PatchEpisodeOfCare is the resolver for the patchEpisodeOfCare field.
func (r *mutationResolver) PatchEpisodeOfCare(ctx context.Context, id string, episodeOfCare dto.EpisodeOfCareInput) (*dto.EpisodeOfCare, error) {
	r.CheckDependencies()
	return r.usecases.PatchEpisodeOfCare(ctx, id, episodeOfCare)
}

// EndEpisodeOfCare is the resolver for the endEpisodeOfCare field.
func (r *mutationResolver) EndEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error) {
	r.CheckDependencies()
	return r.usecases.EndEpisodeOfCare(ctx, id)
}

// StartEncounter is the resolver for the startEncounter field.
func (r *mutationResolver) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	r.CheckDependencies()
	return r.usecases.StartEncounter(ctx, episodeID)
}

// PatchEncounter is the resolver for the patchEncounter field.
func (r *mutationResolver) PatchEncounter(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error) {
	r.CheckDependencies()
	return r.usecases.PatchEncounter(ctx, encounterID, input)
}

// EndEncounter is the resolver for the endEncounter field.
func (r *mutationResolver) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	r.CheckDependencies()
	return r.usecases.EndEncounter(ctx, encounterID)
}

// RecordTemperature is the resolver for the recordTemperature field.
func (r *mutationResolver) RecordTemperature(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordTemperature(ctx, input)
}

// RecordHeight is the resolver for the recordHeight field.
func (r *mutationResolver) RecordHeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordHeight(ctx, input)
}

// RecordWeight is the resolver for the recordWeight field.
func (r *mutationResolver) RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordWeight(ctx, input)
}

// RecordRespiratoryRate is the resolver for the recordRespiratoryRate field.
func (r *mutationResolver) RecordRespiratoryRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordRespiratoryRate(ctx, input)
}

// RecordPulseRate is the resolver for the recordPulseRate field.
func (r *mutationResolver) RecordPulseRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordPulseRate(ctx, input)
}

// RecordBloodPressure is the resolver for the recordBloodPressure field.
func (r *mutationResolver) RecordBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordBloodPressure(ctx, input)
}

// RecordBmi is the resolver for the recordBMI field.
func (r *mutationResolver) RecordBmi(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordBMI(ctx, input)
}

// RecordViralLoad is the resolver for the recordViralLoad field.
func (r *mutationResolver) RecordViralLoad(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()

	return r.usecases.RecordViralLoad(ctx, input)
}

// RecordMuac is the resolver for the recordMUAC field.
func (r *mutationResolver) RecordMuac(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordMuac(ctx, input)
}

// RecordOxygenSaturation is the resolver for the recordOxygenSaturation field.
func (r *mutationResolver) RecordOxygenSaturation(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordOxygenSaturation(ctx, input)
}

// RecordBloodSugar is the resolver for the recordBloodSugar field.
func (r *mutationResolver) RecordBloodSugar(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordBloodSugar(ctx, input)
}

// RecordLastMenstrualPeriod is the resolver for the recordLastMenstrualPeriod field.
func (r *mutationResolver) RecordLastMenstrualPeriod(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordLastMenstrualPeriod(ctx, input)
}

// RecordDiastolicBloodPressure is the resolver for the recordDiastolicBloodPressure field.
func (r *mutationResolver) RecordDiastolicBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordDiastolicBloodPressure(ctx, input)
}

// RecordColposcopy is the resolver for the recordColposcopy field.
func (r *mutationResolver) RecordColposcopy(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	return r.usecases.RecordColposcopy(ctx, input)
}

// RecordHpv is the resolver for the recordHPV field.
func (r *mutationResolver) RecordHpv(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordHPV(ctx, input)
}

// RecordVia is the resolver for the recordVIA field.
func (r *mutationResolver) RecordVia(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.RecordVIA(ctx, input)
}

// RecordPapSmear is the resolver for the recordPapSmear field.
func (r *mutationResolver) RecordPapSmear(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()

	return r.usecases.RecordPapSmear(ctx, input)
}

// CreatePatient is the resolver for the createPatient field.
func (r *mutationResolver) CreatePatient(ctx context.Context, input dto.PatientInput) (*dto.Patient, error) {
	r.CheckDependencies()
	return r.usecases.CreatePatient(ctx, input)
}

// PatchPatient is the resolver for the patchPatient field.
func (r *mutationResolver) PatchPatient(ctx context.Context, id string, input dto.PatientInput) (*dto.Patient, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatient(ctx, id, input)
}

// DeletePatient is the resolver for the deletePatient field.
func (r *mutationResolver) DeletePatient(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.usecases.DeletePatient(ctx, id)
}

// CreateCondition is the resolver for the createCondition field.
func (r *mutationResolver) CreateCondition(ctx context.Context, input dto.ConditionInput) (*dto.Condition, error) {
	r.CheckDependencies()
	return r.usecases.CreateCondition(ctx, input)
}

// CreateAllergyIntolerance is the resolver for the createAllergyIntolerance field.
func (r *mutationResolver) CreateAllergyIntolerance(ctx context.Context, input dto.AllergyInput) (*dto.Allergy, error) {
	return r.usecases.CreateAllergyIntolerance(ctx, input)
}

// CreateComposition is the resolver for the createComposition field.
func (r *mutationResolver) CreateComposition(ctx context.Context, input dto.CompositionInput) (*dto.Composition, error) {
	r.CheckDependencies()
	return r.usecases.CreateComposition(ctx, input)
}

// AppendNoteToComposition is the resolver for the AppendNoteToComposition field.
func (r *mutationResolver) AppendNoteToComposition(ctx context.Context, id string, input dto.PatchCompositionInput) (*dto.Composition, error) {
	r.CheckDependencies()
	return r.usecases.AppendNoteToComposition(ctx, id, input)
}

// PatchPatientHeight is the resolver for the patchHeight field.
func (r *mutationResolver) PatchPatientHeight(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientHeight(ctx, id, value)
}

// PatchPatientWeight is the resolver for the patchWeight field.
func (r *mutationResolver) PatchPatientWeight(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientWeight(ctx, id, value)
}

// PatchPatientBmi is the resolver for the patchPatientBMI field.
func (r *mutationResolver) PatchPatientBmi(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientBMI(ctx, id, value)
}

// PatchPatientTemperature is the resolver for the patchPatientTemperature field.
func (r *mutationResolver) PatchPatientTemperature(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientTemperature(ctx, id, value)
}

// PatchPatientDiastolicBloodPressure is the resolver for the patchPatientDiastolicBloodPressure field.
func (r *mutationResolver) PatchPatientDiastolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientDiastolicBloodPressure(ctx, id, value)
}

// PatchPatientSystolicBloodPressure is the resolver for the patchPatientSystolicBloodPressure field.
func (r *mutationResolver) PatchPatientSystolicBloodPressure(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientSystolicBloodPressure(ctx, id, value)
}

// PatchPatientRespiratoryRate is the resolver for the patchPatientRespiratoryRate field.
func (r *mutationResolver) PatchPatientRespiratoryRate(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientRespiratoryRate(ctx, id, value)
}

// PatchPatientOxygenSaturation is the resolver for the patchPatientOxygenSaturation field.
func (r *mutationResolver) PatchPatientOxygenSaturation(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientOxygenSaturation(ctx, id, value)
}

// PatchPatientPulseRate is the resolver for the PatchPatientPulseRate field.
func (r *mutationResolver) PatchPatientPulseRate(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientPulseRate(ctx, id, value)
}

// PatchPatientViralLoad is the resolver for the PatchPatientViralLoad field.
func (r *mutationResolver) PatchPatientViralLoad(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientViralLoad(ctx, id, value)
}

// PatchPatientMuac is the resolver for the patchPatientMuac field.
func (r *mutationResolver) PatchPatientMuac(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientMuac(ctx, id, value)
}

// PatchPatientLastMenstrualPeriod is the resolver for the patchPatientLastMenstrualPeriod field.
func (r *mutationResolver) PatchPatientLastMenstrualPeriod(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientLastMenstrualPeriod(ctx, id, value)
}

// PatchPatientBloodSugar is the resolver for the patchPatientBloodSugar field.
func (r *mutationResolver) PatchPatientBloodSugar(ctx context.Context, id string, value string) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientBloodSugar(ctx, id, value)
}

// RecordConsent is the resolver for the recordConsent field.
func (r *mutationResolver) RecordConsent(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error) {
	return r.usecases.RecordConsent(ctx, input)
}

// CreateQuestionnaireResponse is the resolver for the createQuestionnaireResponse field.
func (r *mutationResolver) CreateQuestionnaireResponse(ctx context.Context, questionnaireID string, encounterID string, input dto.QuestionnaireResponse) (string, error) {
	return r.usecases.CreateQuestionnaireResponse(ctx, questionnaireID, encounterID, input)
}

// RecordMammographyResult is the resolver for the recordMammographyResult field.
func (r *mutationResolver) RecordMammographyResult(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	return r.usecases.RecordMammographyResult(ctx, input)
}

// RecordBiopsy is the resolver for the recordBiopsy field.
func (r *mutationResolver) RecordBiopsy(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	return r.usecases.RecordBiopsy(ctx, input)
}

// RecordMri is the resolver for the recordMRI field.
func (r *mutationResolver) RecordMri(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	return r.usecases.RecordMRI(ctx, input)
}

// RecordUltrasound is the resolver for the recordUltrasound field.
func (r *mutationResolver) RecordUltrasound(ctx context.Context, input dto.DiagnosticReportInput) (*dto.DiagnosticReport, error) {
	return r.usecases.RecordUltrasound(ctx, input)
}

// GetEncounterAssociatedResources is the resolver for the getEncounterAssociatedResources field.
func (r *mutationResolver) GetEncounterAssociatedResources(ctx context.Context, encounterID string) (*dto.EncounterAssociatedResources, error) {
	return r.usecases.GetEncounterAssociatedResources(ctx, encounterID)
}

// PatientHealthTimeline is the resolver for the patientHealthTimeline field.
func (r *queryResolver) PatientHealthTimeline(ctx context.Context, input dto.HealthTimelineInput) (*dto.HealthTimeline, error) {
	r.CheckDependencies()

	return r.usecases.PatientHealthTimeline(ctx, input)
}

// GetMedicalData is the resolver for the getMedicalData field.
func (r *queryResolver) GetMedicalData(ctx context.Context, patientID string) (*dto.MedicalData, error) {
	r.CheckDependencies()

	return r.usecases.GetMedicalData(ctx, patientID)
}

// GetEpisodeOfCare is the resolver for the getEpisodeOfCare field.
func (r *queryResolver) GetEpisodeOfCare(ctx context.Context, id string) (*dto.EpisodeOfCare, error) {
	r.CheckDependencies()
	return r.usecases.GetEpisodeOfCare(ctx, id)
}

// ListPatientConditions is the resolver for the listPatientConditions field.
func (r *queryResolver) ListPatientConditions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ConditionConnection, error) {
	r.CheckDependencies()
	return r.usecases.ListPatientConditions(ctx, patientID, encounterID, date, pagination)
}

// ListPatientCompositions is the resolver for the listPatientCompositions field.
func (r *queryResolver) ListPatientCompositions(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.CompositionConnection, error) {
	r.CheckDependencies()
	return r.usecases.ListPatientCompositions(ctx, patientID, encounterID, date, pagination)
}

// ListPatientEncounters is the resolver for the listPatientEncounters field.
func (r *queryResolver) ListPatientEncounters(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.EncounterConnection, error) {
	r.CheckDependencies()

	return r.usecases.ListPatientEncounters(ctx, patientID, &pagination)
}

// GetPatientTemperatureEntries is the resolver for the getPatientTemperatureEntries field.
func (r *queryResolver) GetPatientTemperatureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.GetPatientTemperatureEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientBloodPressureEntries is the resolver for the getPatientBloodPressureEntries field.
func (r *queryResolver) GetPatientBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.GetPatientBloodPressureEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientHeightEntries is the resolver for the getPatientHeightEntries field.
func (r *queryResolver) GetPatientHeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	return r.usecases.GetPatientHeightEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientRespiratoryRateEntries is the resolver for the getPatientRespiratoryRateEntries field.
func (r *queryResolver) GetPatientRespiratoryRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.GetPatientRespiratoryRateEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientPulseRateEntries is the resolver for the GetPatientPulseRateEntries field.
func (r *queryResolver) GetPatientPulseRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	return r.usecases.GetPatientPulseRateEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientBMIEntries is the resolver for the getPatientBMIEntries field.
func (r *queryResolver) GetPatientBMIEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.GetPatientBMIEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientWeightEntries is the resolver for the getPatientWeightEntries field.
func (r *queryResolver) GetPatientWeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	return r.usecases.GetPatientWeightEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientMuacEntries is the resolver for the getPatientMuacEntries field.
func (r *queryResolver) GetPatientMuacEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.GetPatientMuacEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientOxygenSaturationEntries is the resolver for the getPatientOxygenSaturationEntries field.
func (r *queryResolver) GetPatientOxygenSaturationEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.GetPatientOxygenSaturationEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientViralLoad is the resolver for the getPatientViralLoad field.
func (r *queryResolver) GetPatientViralLoad(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()

	return r.usecases.GetPatientViralLoad(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientBloodSugarEntries is the resolver for the getPatientBloodSugarEntries field.
func (r *queryResolver) GetPatientBloodSugarEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.GetPatientBloodSugarEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientLastMenstrualPeriodEntries is the resolver for the GetPatientLastMenstrualPeriodEntries field.
func (r *queryResolver) GetPatientLastMenstrualPeriodEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.GetPatientLastMenstrualPeriodEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientDiastolicBloodPressureEntries is the resolver for the getPatientDiastolicBloodPressureEntries field.
func (r *queryResolver) GetPatientDiastolicBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.GetPatientDiastolicBloodPressureEntries(ctx, patientID, encounterID, date, &pagination)
}

// SearchAllergy is the resolver for the searchAllergy field.
func (r *queryResolver) SearchAllergy(ctx context.Context, name string, pagination dto.Pagination) (*dto.TerminologyConnection, error) {
	r.CheckDependencies()

	return r.usecases.SearchAllergy(ctx, name, pagination)
}

// GetAllergy is the resolver for the getAllergy field.
func (r *queryResolver) GetAllergy(ctx context.Context, id string) (*dto.Allergy, error) {
	r.CheckDependencies()

	return r.usecases.GetAllergyIntolerance(ctx, id)
}

// ListPatientAllergies is the resolver for the listPatientAllergies field.
func (r *queryResolver) ListPatientAllergies(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.AllergyConnection, error) {
	r.CheckDependencies()

	return r.usecases.ListPatientAllergies(ctx, patientID, pagination)
}

// ListPatientMedia is the resolver for the listPatientMedia field.
func (r *queryResolver) ListPatientMedia(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.MediaConnection, error) {
	r.CheckDependencies()

	return r.usecases.ListPatientMedia(ctx, patientID, pagination)
}

// GetQuestionnaireResponseRiskLevel is the resolver for the getQuestionnaireResponseRiskLevel field.
func (r *queryResolver) GetQuestionnaireResponseRiskLevel(ctx context.Context, encounterID string, screeningType domain.ScreeningTypeEnum) (string, error) {
	r.CheckDependencies()

	return r.usecases.GetQuestionnaireResponseRiskLevel(ctx, encounterID, screeningType)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
