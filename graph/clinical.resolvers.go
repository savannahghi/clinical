package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
	"gitlab.slade360emr.com/go/clinical/graph/generated"
)

func (r *mutationResolver) CreateFHIRAllergyIntolerance(ctx context.Context, input clinical.FHIRAllergyIntoleranceInput) (*clinical.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRAllergyIntolerance(ctx, input)
}

func (r *mutationResolver) UpdateFHIRAllergyIntolerance(ctx context.Context, input clinical.FHIRAllergyIntoleranceInput) (*clinical.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRAllergyIntolerance(ctx, input)
}

func (r *mutationResolver) DeleteFHIRAllergyIntolerance(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRAllergyIntolerance(ctx, id)
}

func (r *mutationResolver) CreateFHIRAppointment(ctx context.Context, input clinical.FHIRAppointmentInput) (*clinical.FHIRAppointmentRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRAppointment(ctx, input)
}

func (r *mutationResolver) UpdateFHIRAppointment(ctx context.Context, input clinical.FHIRAppointmentInput) (*clinical.FHIRAppointmentRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRAppointment(ctx, input)
}

func (r *mutationResolver) DeleteFHIRAppointment(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRAppointment(ctx, id)
}

func (r *mutationResolver) CreateFHIRComposition(ctx context.Context, input clinical.FHIRCompositionInput) (*clinical.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRComposition(ctx, input)
}

func (r *mutationResolver) UpdateFHIRComposition(ctx context.Context, input clinical.FHIRCompositionInput) (*clinical.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRComposition(ctx, input)
}

func (r *mutationResolver) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRComposition(ctx, id)
}

func (r *mutationResolver) CreateFHIRCondition(ctx context.Context, input clinical.FHIRConditionInput) (*clinical.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRCondition(ctx, input)
}

func (r *mutationResolver) UpdateFHIRCondition(ctx context.Context, input clinical.FHIRConditionInput) (*clinical.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRCondition(ctx, input)
}

func (r *mutationResolver) DeleteFHIRCondition(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRCondition(ctx, id)
}

func (r *mutationResolver) CreateFHIREncounter(ctx context.Context, input clinical.FHIREncounterInput) (*clinical.FHIREncounterRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIREncounter(ctx, input)
}

func (r *mutationResolver) UpdateFHIREncounter(ctx context.Context, input clinical.FHIREncounterInput) (*clinical.FHIREncounterRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIREncounter(ctx, input)
}

func (r *mutationResolver) DeleteFHIREncounter(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIREncounter(ctx, id)
}

func (r *mutationResolver) CreateFHIREpisodeOfCare(ctx context.Context, input clinical.FHIREpisodeOfCareInput) (*clinical.FHIREpisodeOfCareRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIREpisodeOfCare(ctx, input)
}

func (r *mutationResolver) UpdateFHIREpisodeOfCare(ctx context.Context, input clinical.FHIREpisodeOfCareInput) (*clinical.FHIREpisodeOfCareRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIREpisodeOfCare(ctx, input)
}

func (r *mutationResolver) DeleteFHIREpisodeOfCare(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIREpisodeOfCare(ctx, id)
}

func (r *mutationResolver) CreateFHIRMedicationRequest(ctx context.Context, input clinical.FHIRMedicationRequestInput) (*clinical.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRMedicationRequest(ctx, input)
}

func (r *mutationResolver) UpdateFHIRMedicationRequest(ctx context.Context, input clinical.FHIRMedicationRequestInput) (*clinical.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRMedicationRequest(ctx, input)
}

func (r *mutationResolver) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRMedicationRequest(ctx, id)
}

func (r *mutationResolver) CreateFHIRObservation(ctx context.Context, input clinical.FHIRObservationInput) (*clinical.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRObservation(ctx, input)
}

func (r *mutationResolver) UpdateFHIRObservation(ctx context.Context, input clinical.FHIRObservationInput) (*clinical.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRObservation(ctx, input)
}

func (r *mutationResolver) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRObservation(ctx, id)
}

func (r *mutationResolver) CreateFHIROrganization(ctx context.Context, input clinical.FHIROrganizationInput) (*clinical.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIROrganization(ctx, input)
}

func (r *mutationResolver) UpdateFHIROrganization(ctx context.Context, input clinical.FHIROrganizationInput) (*clinical.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIROrganization(ctx, input)
}

func (r *mutationResolver) DeleteFHIROrganization(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIROrganization(ctx, id)
}

func (r *mutationResolver) CreateFHIRPatient(ctx context.Context, input clinical.FHIRPatientInput) (*clinical.FHIRPatientRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRPatient(ctx, input)
}

func (r *mutationResolver) UpdateFHIRPatient(ctx context.Context, input clinical.FHIRPatientInput) (*clinical.FHIRPatientRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRPatient(ctx, input)
}

func (r *mutationResolver) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRPatient(ctx, id)
}

func (r *mutationResolver) CreateFHIRServiceRequest(ctx context.Context, input clinical.FHIRServiceRequestInput) (*clinical.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.CreateFHIRServiceRequest(ctx, input)
}

func (r *mutationResolver) UpdateFHIRServiceRequest(ctx context.Context, input clinical.FHIRServiceRequestInput) (*clinical.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.UpdateFHIRServiceRequest(ctx, input)
}

func (r *mutationResolver) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	r.CheckDependencies()
	return r.clinicalService.DeleteFHIRServiceRequest(ctx, id)
}

func (r *mutationResolver) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.StartEncounter(ctx, episodeID)
}

func (r *mutationResolver) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.EndEncounter(ctx, encounterID)
}

func (r *mutationResolver) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.EndEpisode(ctx, episodeID)
}

func (r *queryResolver) GetFHIRAllergyIntolerance(ctx context.Context, id string) (*clinical.FHIRAllergyIntoleranceRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRAllergyIntolerance(ctx, id)
}

func (r *queryResolver) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*clinical.FHIRAllergyIntoleranceRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRAllergyIntolerance(ctx, params)
}

func (r *queryResolver) GetFHIRAppointment(ctx context.Context, id string) (*clinical.FHIRAppointmentRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRAppointment(ctx, id)
}

func (r *queryResolver) SearchFHIRAppointment(ctx context.Context, params map[string]interface{}) (*clinical.FHIRAppointmentRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRAppointment(ctx, params)
}

func (r *queryResolver) GetFHIRComposition(ctx context.Context, id string) (*clinical.FHIRCompositionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRComposition(ctx, id)
}

func (r *queryResolver) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*clinical.FHIRCompositionRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRComposition(ctx, params)
}

func (r *queryResolver) GetFHIRCondition(ctx context.Context, id string) (*clinical.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRCondition(ctx, id)
}

func (r *queryResolver) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*clinical.FHIRConditionRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRCondition(ctx, params)
}

func (r *queryResolver) GetFHIREncounter(ctx context.Context, id string) (*clinical.FHIREncounterRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIREncounter(ctx, id)
}

func (r *queryResolver) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*clinical.FHIREncounterRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIREncounter(ctx, params)
}

func (r *queryResolver) GetFHIREpisodeOfCare(ctx context.Context, id string) (*clinical.FHIREpisodeOfCareRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIREpisodeOfCare(ctx, id)
}

func (r *queryResolver) SearchFHIREpisodeOfCare(ctx context.Context, params map[string]interface{}) (*clinical.FHIREpisodeOfCareRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIREpisodeOfCare(ctx, params)
}

func (r *queryResolver) GetFHIRMedicationRequest(ctx context.Context, id string) (*clinical.FHIRMedicationRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRMedicationRequest(ctx, id)
}

func (r *queryResolver) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*clinical.FHIRMedicationRequestRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRMedicationRequest(ctx, params)
}

func (r *queryResolver) GetFHIRObservation(ctx context.Context, id string) (*clinical.FHIRObservationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRObservation(ctx, id)
}

func (r *queryResolver) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*clinical.FHIRObservationRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRObservation(ctx, params)
}

func (r *queryResolver) GetFHIROrganization(ctx context.Context, id string) (*clinical.FHIROrganizationRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIROrganization(ctx, id)
}

func (r *queryResolver) SearchFHIROrganization(ctx context.Context, params map[string]interface{}) (*clinical.FHIROrganizationRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIROrganization(ctx, params)
}

func (r *queryResolver) GetFHIRPatient(ctx context.Context, id string) (*clinical.FHIRPatientRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRPatient(ctx, id)
}

func (r *queryResolver) SearchFHIRPatient(ctx context.Context, params map[string]interface{}) (*clinical.FHIRPatientRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRPatient(ctx, params)
}

func (r *queryResolver) GetFHIRServiceRequest(ctx context.Context, id string) (*clinical.FHIRServiceRequestRelayPayload, error) {
	r.CheckDependencies()
	return r.clinicalService.GetFHIRServiceRequest(ctx, id)
}

func (r *queryResolver) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*clinical.FHIRServiceRequestRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.SearchFHIRServiceRequest(ctx, params)
}

func (r *queryResolver) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.AllergySummary(ctx, patientID)
}

func (r *queryResolver) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.ProblemSummary(ctx, patientID)
}

func (r *queryResolver) RequestUSSDFullHistory(ctx context.Context, input clinical.USSDClinicalRequest) (*clinical.USSDMedicalHistoryClinicalResponse, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.RequestUSSDFullHistory(ctx, input)
}

func (r *queryResolver) RequestUSSDLastVisit(ctx context.Context, input clinical.USSDClinicalRequest) (*clinical.USSDLastVisitClinicalResponse, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.RequestUSSDLastVisit(ctx, input)
}

func (r *queryResolver) RequestUSSDPatientProfile(ctx context.Context, input clinical.USSDClinicalRequest) (*clinical.USSDPatientProfileClinicalResponse, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.RequestUSSDPatientProfile(ctx, input)
}

func (r *queryResolver) PatientTimeline(ctx context.Context, episodeID string) ([]map[string]interface{}, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.PatientTimeline(ctx, episodeID)
}

func (r *queryResolver) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.PatientTimelineWithCount(ctx, episodeID, count)
}

func (r *queryResolver) VisitSummary(ctx context.Context, encounterID string) (map[string]interface{}, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.VisitSummary(ctx, encounterID, clinical.MaxClinicalRecordPageSize)
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
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) ListAppointments(ctx context.Context, providerSladeCode int) (*clinical.FHIRAppointmentRelayConnection, error) {
	r.CheckDependencies()
	return r.clinicalService.ListAppointments(ctx, providerSladeCode)
}
func (r *queryResolver) GetROrganization(ctx context.Context, providerSladeCode int) (*string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.GetOrganization(ctx, providerSladeCode)
}
func (r *queryResolver) CreateOrganization(ctx context.Context, providerSladeCode int) (*string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.CreateOrganization(ctx, providerSladeCode)
}
func (r *queryResolver) GetORCreateOrganization(ctx context.Context, providerSladeCode int) (*string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.GetORCreateOrganization(ctx, providerSladeCode)
}
