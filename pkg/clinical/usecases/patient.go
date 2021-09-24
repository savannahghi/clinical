package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/clinical/pkg/clinical/infrastructure"
	"github.com/savannahghi/converterandformatter"
)

// ClinicalUseCase ...
type ClinicalUseCase interface {
	ProblemSummary(ctx context.Context, patientID string) ([]string, error)
	VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error)
	PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error)
	CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error)
	CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error)
}

// ClinicalUseCaseImpl ...
type ClinicalUseCaseImpl struct {
	infrastructure infrastructure.Infrastructure
}

// NewClinicalUseCaseImpl ...
func NewClinicalUseCaseImpl(infra infrastructure.Infrastructure) ClinicalUseCase {
	return &ClinicalUseCaseImpl{
		infrastructure: infra,
	}
}

// ProblemSummary ...
func (c ClinicalUseCaseImpl) ProblemSummary(ctx context.Context, patientID string) ([]string, error) {
	return nil, nil
}

// VisitSummary ...
func (c ClinicalUseCaseImpl) VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error) {
	return nil, nil
}

// PatientTimelineWithCount ...
func (c ClinicalUseCaseImpl) PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error) {
	return nil, nil
}

// CreateEpisodeOfCare ...
func (c ClinicalUseCaseImpl) CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error) {
	return nil, nil
}

// CreateFHIRCondition creates a FHIRCondition instance
func (c ClinicalUseCaseImpl) CreateFHIRCondition(ctx context.Context, input domain.FHIRConditionInput) (*domain.FHIRConditionRelayPayload, error) {
	// TODO: return casbin and check precondition
	resourceType := "Condition"
	resource := domain.FHIRCondition{}

	payload, err := converterandformatter.StructToMap(input)
	if err != nil {
		return nil, fmt.Errorf("unable to turn %s input into a map: %v", resourceType, err)
	}

	data, err := c.infrastructure.FHIRRepo.CreateFHIRResource(resourceType, payload)
	if err != nil {
		return nil, fmt.Errorf("unable to create/update %s resource: %v", resourceType, err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to unmarshal %s response JSON: data: %v\n, error: %v",
			resourceType, string(data), err)
	}

	output := &domain.FHIRConditionRelayPayload{
		Resource: &resource,
	}
	return output, nil
}
