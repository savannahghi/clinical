package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *mutationResolver) StartEpisodeByOtp(ctx context.Context, input clinical.OTPEpisodeCreationInput) (*clinical.EpisodeOfCarePayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.StartEpisodeByOtp(ctx, input)
}

func (r *mutationResolver) UpgradeEpisode(ctx context.Context, input clinical.OTPEpisodeUpgradeInput) (*clinical.EpisodeOfCarePayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.UpgradeEpisode(ctx, input)
}

func (r *mutationResolver) StartEpisodeByBreakGlass(ctx context.Context, input clinical.BreakGlassEpisodeCreationInput) (*clinical.EpisodeOfCarePayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.StartEpisodeByBreakGlass(ctx, input)
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

func (r *mutationResolver) RegisterPatient(ctx context.Context, input clinical.SimplePatientRegistrationInput) (*clinical.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.RegisterPatient(ctx, input)
}

func (r *mutationResolver) UpdatePatient(ctx context.Context, input clinical.SimplePatientRegistrationInput) (*clinical.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.UpdatePatient(ctx, input)
}

func (r *mutationResolver) DeletePatient(ctx context.Context, input clinical.RetirePatientInput) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.DeletePatient(ctx, input)
}

func (r *mutationResolver) AddNextOfKin(ctx context.Context, input clinical.SimpleNextOfKinInput) (*clinical.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.AddNextOfKin(ctx, input)
}

func (r *mutationResolver) AddNhif(ctx context.Context, input *clinical.SimpleNHIFInput) (*clinical.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.AddNhif(ctx, input)
}

func (r *queryResolver) OpenOrganizationEpisodes(ctx context.Context, providerSladeCode string) ([]*clinical.FHIREpisodeOfCare, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.OpenOrganizationEpisodes(ctx, providerSladeCode)
}

func (r *queryResolver) FindPatientsByMsisdn(ctx context.Context, msisdn string) (*clinical.PatientConnection, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.FindPatientsByMSISDN(ctx, msisdn)
}

func (r *queryResolver) PatientSearch(ctx context.Context, search string) (*clinical.PatientConnection, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.PatientSearch(ctx, search)
}

func (r *queryResolver) OpenEpisodes(ctx context.Context, patientReference string) ([]*clinical.FHIREpisodeOfCare, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.OpenEpisodes(ctx, patientReference)
}
