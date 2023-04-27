package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// AllergyUsecaseMock defines the Allergy mock methods
type AllergyUsecaseMock struct {
	MockSearchAllergyFn         func(ctx context.Context, name *string, pagination dto.Pagination) (*dto.TerminologyConnection, error)
	MockGetAllergyIntoleranceFn func(ctx context.Context, id string) (*dto.Allergy, error)
}

// NewAllergyUsecaseMock constructor initializes a new instance of allergy mock methods
func NewAllergyUsecaseMock() *AllergyUsecaseMock {
	return &AllergyUsecaseMock{
		MockSearchAllergyFn: func(ctx context.Context, name *string, pagination dto.Pagination) (*dto.TerminologyConnection, error) {
			next := "https://api.openconceptlab.org/concepts/?page=3&q=peanuts&limit=2&exact_match=off&verbose=false&includeRetired=false&includeInverseMappings=false"
			previous := "https://api.openconceptlab.org/concepts/?page=1&q=peanuts&limit=2&exact_match=off&verbose=false&includeRetired=false&includeInverseMappings=false"
			return &dto.TerminologyConnection{
				TotalCount: 10,
				Edges: []dto.TerminologyEdge{
					{
						Node: dto.Terminology{
							Code:   "test",
							System: dto.TerminologySourceCIEL,
							Name:   "test",
						},
					},
				},
				PageInfo: dto.PageInfo{
					HasNextPage:     true,
					HasPreviousPage: true,
					EndCursor:       &previous,
					StartCursor:     &next,
				},
			}, nil
		},
		MockGetAllergyIntoleranceFn: func(ctx context.Context, id string) (*dto.Allergy, error) {
			return &dto.Allergy{
				ID:          id,
				PatientID:   id,
				EncounterID: id,
			}, nil
		},
	}
}

// SearchAllergy is the mock implementation of searching for an allergy
func (a *AllergyUsecaseMock) SearchAllergy(ctx context.Context, name *string, pagination dto.Pagination) (*dto.TerminologyConnection, error) {
	return a.MockSearchAllergyFn(ctx, name, pagination)
}

// GetAllergyIntolerance mocks the implementation of getting allergy intolerance by id
func (a *AllergyUsecaseMock) GetAllergyIntolerance(ctx context.Context, id string) (*dto.Allergy, error) {
	return a.MockGetAllergyIntoleranceFn(ctx, id)
}
