package usecases

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// ClinicalUseCase ...
type ClinicalUseCase interface {
	ProblemSummary(ctx context.Context, patientID string) ([]string, error)
	VisitSummary(ctx context.Context, encounterID string, count int) (map[string]interface{}, error)
	PatientTimelineWithCount(ctx context.Context, episodeID string, count int) ([]map[string]interface{}, error)
	CreateEpisodeOfCare(ctx context.Context, episode domain.FHIREpisodeOfCare) (*domain.EpisodeOfCarePayload, error)
}

// ClinicalUseCaseImpl ...
type ClinicalUseCaseImpl struct {
}

// NewClinicalUseCaseImpl ...
func NewClinicalUseCaseImpl() ClinicalUseCase {
	return &ClinicalUseCaseImpl{}
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
