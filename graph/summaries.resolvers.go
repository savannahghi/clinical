package graph

import (
	"context"

	"gitlab.slade360emr.com/go/clinical/graph/clinical"
)

func (r *queryResolver) AllergySummary(
	ctx context.Context, patientID string) ([]string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.AllergySummary(ctx, patientID)
}

func (r *queryResolver) ProblemSummary(
	ctx context.Context, patientID string) ([]string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.ProblemSummary(ctx, patientID)
}

func (r *queryResolver) RequestUSSDFullHistory(
	ctx context.Context, input clinical.USSDClinicalRequest) (*clinical.USSDClinicalResponse, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.RequestUSSDFullHistory(ctx, input)
}

func (r *queryResolver) RequestUSSDLastVisit(
	ctx context.Context, input clinical.USSDClinicalRequest) (*clinical.USSDClinicalResponse, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.RequestUSSDLastVisit(ctx, input)
}

func (r *queryResolver) RequestUSSDPatientProfile(
	ctx context.Context, input clinical.USSDClinicalRequest) (*clinical.USSDClinicalResponse, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.RequestUSSDPatientProfile(ctx, input)
}

func (r *queryResolver) VisitSummary(
	ctx context.Context, encounterID string) (map[string]interface{}, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.VisitSummary(ctx, encounterID)
}

func (r *queryResolver) PatientTimeline(
	ctx context.Context, episodeID string) ([]map[string]interface{}, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.PatientTimeline(ctx, episodeID)
}

func (r *mutationResolver) EndEncounter(
	ctx context.Context, encounterID string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.EndEncounter(ctx, encounterID)
}

func (r *mutationResolver) EndEpisode(
	ctx context.Context, episodeID string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.EndEpisode(ctx, episodeID)
}

func (r *mutationResolver) StartEncounter(
	ctx context.Context, episodeID string) (string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.clinicalService.StartEncounter(ctx, episodeID)
}
