package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
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
	return r.usecases.Clinical.StartEncounter(ctx, episodeID)
}

// PatchEncounter is the resolver for the patchEncounter field.
func (r *mutationResolver) PatchEncounter(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.PatchEncounter(ctx, encounterID, input)
}

// EndEncounter is the resolver for the endEncounter field.
func (r *mutationResolver) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.EndEncounter(ctx, encounterID)
}

// RecordTemperature is the resolver for the recordTemperature field.
func (r *mutationResolver) RecordTemperature(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.RecordTemperature(ctx, input)
}

// RecordHeight is the resolver for the recordHeight field.
func (r *mutationResolver) RecordHeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.RecordHeight(ctx, input)
}

// RecordWeight is the resolver for the recordWeight field.
func (r *mutationResolver) RecordWeight(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.RecordWeight(ctx, input)
}

// RecordRespiratoryRate is the resolver for the recordRespiratoryRate field.
func (r *mutationResolver) RecordRespiratoryRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.RecordRespiratoryRate(ctx, input)
}

// RecordPulseRate is the resolver for the recordPulseRate field.
func (r *mutationResolver) RecordPulseRate(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.RecordPulseRate(ctx, input)
}

// RecordBloodPressure is the resolver for the recordBloodPressure field.
func (r *mutationResolver) RecordBloodPressure(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.RecordBloodPressure(ctx, input)
}

// RecordBmi is the resolver for the recordBMI field.
func (r *mutationResolver) RecordBmi(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.RecordBMI(ctx, input)
}

// RecordViralLoad is the resolver for the recordViralLoad field.
func (r *mutationResolver) RecordViralLoad(ctx context.Context, input dto.ObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()

	return r.usecases.Clinical.RecordViralLoad(ctx, input)
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

// PatchHeight is the resolver for the patchHeight field.
func (r *mutationResolver) PatchHeight(ctx context.Context, id string, input dto.PatchObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientHeightEntries(ctx, id, input)
}

// PatchWeight is the resolver for the patchWeight field.
func (r *mutationResolver) PatchWeight(ctx context.Context, id string, input dto.PatchObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientWeightEntries(ctx, id, input)
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

	return r.usecases.Clinical.ListPatientEncounters(ctx, patientID, &pagination)
}

// GetPatientTemperatureEntries is the resolver for the getPatientTemperatureEntries field.
func (r *queryResolver) GetPatientTemperatureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.GetPatientTemperatureEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientBloodPressureEntries is the resolver for the getPatientBloodPressureEntries field.
func (r *queryResolver) GetPatientBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.GetPatientBloodPressureEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientHeightEntries is the resolver for the getPatientHeightEntries field.
func (r *queryResolver) GetPatientHeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	return r.usecases.Clinical.GetPatientHeightEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientRespiratoryRateEntries is the resolver for the getPatientRespiratoryRateEntries field.
func (r *queryResolver) GetPatientRespiratoryRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.GetPatientRespiratoryRateEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientPulseRateEntries is the resolver for the GetPatientPulseRateEntries field.
func (r *queryResolver) GetPatientPulseRateEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	return r.usecases.Clinical.GetPatientPulseRateEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientBMIEntries is the resolver for the getPatientBMIEntries field.
func (r *queryResolver) GetPatientBMIEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.GetPatientBMIEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientWeightEntries is the resolver for the getPatientWeightEntries field.
func (r *queryResolver) GetPatientWeightEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	return r.usecases.Clinical.GetPatientWeightEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientMuacEntries is the resolver for the getPatientMuacEntries field.
func (r *queryResolver) GetPatientMuacEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.GetPatientMuacEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientOxygenSaturationEntries is the resolver for the getPatientOxygenSaturationEntries field.
func (r *queryResolver) GetPatientOxygenSaturationEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.GetPatientOxygenSaturationEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientViralLoad is the resolver for the getPatientViralLoad field.
func (r *queryResolver) GetPatientViralLoad(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()

	return r.usecases.Clinical.GetPatientViralLoad(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientBloodSugarEntries is the resolver for the getPatientBloodSugarEntries field.
func (r *queryResolver) GetPatientBloodSugarEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.GetPatientBloodSugarEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientLastMenstrualPeriodEntries is the resolver for the GetPatientLastMenstrualPeriodEntries field.
func (r *queryResolver) GetPatientLastMenstrualPeriodEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.GetPatientLastMenstrualPeriodEntries(ctx, patientID, encounterID, date, &pagination)
}

// GetPatientDiastolicBloodPressureEntries is the resolver for the getPatientDiastolicBloodPressureEntries field.
func (r *queryResolver) GetPatientDiastolicBloodPressureEntries(ctx context.Context, patientID string, encounterID *string, date *scalarutils.Date, pagination dto.Pagination) (*dto.ObservationConnection, error) {
	r.CheckDependencies()
	return r.usecases.Clinical.GetPatientDiastolicBloodPressureEntries(ctx, patientID, encounterID, date, &pagination)
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) Patchweight(ctx context.Context, id string, input dto.PatchObservationInput) (*dto.Observation, error) {
	r.CheckDependencies()
	return r.usecases.PatchPatientWeightEntries(ctx, id, input)
}
