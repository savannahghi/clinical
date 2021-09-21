package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/domain/model"
	"github.com/savannahghi/clinical/pkg/clinical/presentation/graph/generated"
)

func (r *mutationResolver) StartEpisodeByOtp(ctx context.Context, input model.OTPEpisodeCreationInput) (*model.EpisodeOfCarePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StartEpisodeByBreakGlass(ctx context.Context, input model.BreakGlassEpisodeCreationInput) (*model.EpisodeOfCarePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpgradeEpisode(ctx context.Context, input model.OTPEpisodeUpgradeInput) (*model.EpisodeOfCarePayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) EndEpisode(ctx context.Context, episodeID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RegisterPatient(ctx context.Context, input model.SimplePatientRegistrationInput) (*model.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.SimplePatientRegistrationInput) (*model.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePatient(ctx context.Context, input model.SimplePatientRegistrationInput) (*model.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddNextOfKin(ctx context.Context, input model.SimpleNextOfKinInput) (*model.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddNhif(ctx context.Context, input *model.SimpleNHIFInput) (*model.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUpdatePatientExtraInformation(ctx context.Context, input model.PatientExtraInformationInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRMedicationRequest(ctx context.Context, input model.FHIRMedicationRequestInput) (*model.FHIRMedicationRequestRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFHIRMedicationRequest(ctx context.Context, input model.FHIRMedicationRequestInput) (*model.FHIRMedicationRequestRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFHIRMedicationRequest(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRAllergyIntolerance(ctx context.Context, input model.FHIRAllergyIntoleranceInput) (*model.FHIRAllergyIntoleranceRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFHIRAllergyIntolerance(ctx context.Context, input model.FHIRAllergyIntoleranceInput) (*model.FHIRAllergyIntoleranceRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRCondition(ctx context.Context, input model.FHIRConditionInput) (*model.FHIRConditionRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFHIRCondition(ctx context.Context, input model.FHIRConditionInput) (*model.FHIRConditionRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRServiceRequest(ctx context.Context, input model.FHIRServiceRequestInput) (*model.FHIRServiceRequestRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFHIRServiceRequest(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRObservation(ctx context.Context, input model.FHIRObservationInput) (*model.FHIRObservationRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateFHIRComposition(ctx context.Context, input model.FHIRCompositionInput) (*model.FHIRCompositionRelayPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateFHIRComposition(ctx context.Context, input model.FHIRCompositionInput) (*model.FHIRCompositionRelayPayload, error) {
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

func (r *queryResolver) FindPatientsByMsisdn(ctx context.Context, msisdn string) (*model.PatientConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindPatients(ctx context.Context, search string) (*model.PatientConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetPatient(ctx context.Context, id string) (*model.PatientPayload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OpenEpisodes(ctx context.Context, patientReference string) ([]*model.FHIREpisodeOfCare, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) OpenOrganizationEpisodes(ctx context.Context, providerSladeCode string) ([]*model.FHIREpisodeOfCare, error) {
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

func (r *queryResolver) SearchFHIREncounter(ctx context.Context, params map[string]interface{}) (*model.FHIREncounterRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRCondition(ctx context.Context, params map[string]interface{}) (*model.FHIRConditionRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRAllergyIntolerance(ctx context.Context, params map[string]interface{}) (*model.FHIRAllergyIntoleranceRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRObservation(ctx context.Context, params map[string]interface{}) (*model.FHIRObservationRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRMedicationRequest(ctx context.Context, params map[string]interface{}) (*model.FHIRMedicationRequestRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRServiceRequest(ctx context.Context, params map[string]interface{}) (*model.FHIRServiceRequestRelayConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) SearchFHIRComposition(ctx context.Context, params map[string]interface{}) (*model.FHIRCompositionRelayConnection, error) {
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
