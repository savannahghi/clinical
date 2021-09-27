package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph/generated"
)

func (r *mutationResolver) StartEpisodeByOtp(ctx context.Context, input domain.OTPEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StartEpisodeByBreakGlass(ctx context.Context, input domain.BreakGlassEpisodeCreationInput) (*domain.EpisodeOfCarePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpgradeEpisode(ctx context.Context, input domain.OTPEpisodeUpgradeInput) (*domain.EpisodeOfCarePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.usecases.StartEncounter(ctx, episodeID)
}

func (r *mutationResolver) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RegisterPatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.usecases.RegisterPatient(ctx, input)
}

func (r *mutationResolver) RegisterUser(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePatient(ctx context.Context, input domain.SimplePatientRegistrationInput) (*domain.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddNextOfKin(ctx context.Context, input domain.SimpleNextOfKinInput) (*domain.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddNhif(ctx context.Context, input *domain.SimpleNHIFInput) (*domain.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUpdatePatientExtraInformation(ctx context.Context, input domain.PatientExtraInformationInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFHIRMedicationRequest(ctx context.Context, input domain.FHIRMedicationRequestInput) (*domain.FHIRMedicationRequestRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFHIRAllergyIntolerance(ctx context.Context, input domain.FHIRAllergyIntoleranceInput) (*domain.FHIRAllergyIntoleranceRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.usecases.CreateFHIRCondition(ctx, input)
}

func (r *mutationResolver) UpdateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRServiceRequest(ctx context.Context, input domain.FHIRServiceRequestInput) (*domain.FHIRServiceRequestRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRObservation(ctx context.Context, input domain.FHIRObservationInput) (*domain.FHIRObservationRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFHIRComposition(ctx context.Context, input domain.FHIRCompositionInput) (*domain.FHIRCompositionRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFHIRComposition(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFHIRPatient(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFHIRObservation(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindPatientsByMsisdn(ctx context.Context, msisdn string) (*domain.PatientConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindPatients(ctx context.Context, search string) (*domain.PatientConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetPatient(ctx context.Context, id string) (*domain.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OpenEpisodes(ctx context.Context, patientReference string) ([]*domain.FHIREpisodeOfCare, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OpenOrganizationEpisodes(ctx context.Context, providerSladeCode string) ([]*domain.FHIREpisodeOfCare, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) VisitSummary(ctx context.Context, encounterID string) (map[string]interface{}, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*domain.FHIREncounterRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*domain.FHIRConditionRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*domain.FHIRAllergyIntoleranceRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*domain.FHIRObservationRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRMedicationRequestRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*domain.FHIRServiceRequestRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*domain.FHIRCompositionRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AllergySummary(ctx context.Context, patientID string) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
