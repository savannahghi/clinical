package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
)

// AllergyUsecaseMock defines the Allergy mock methods
type AllergyUsecaseMock struct {
	MockSearchAllergyFn func(ctx context.Context, name *string) ([]*dto.Terminology, error)
}

// NewAllergyUsecaseMock constructor initializes a new instance of allergy mock methods
func NewAllergyUsecaseMock() *AllergyUsecaseMock {
	return &AllergyUsecaseMock{
		MockSearchAllergyFn: func(ctx context.Context, name *string) ([]*dto.Terminology, error) {
			return []*dto.Terminology{
				{
					Code:   "test",
					System: dto.TerminologySourceCIEL,
					Name:   "test",
				},
			}, nil
		},
	}
}

// SearchAllergy is the mock implementation of searching for an allergy
func (a *AllergyUsecaseMock) SearchAllergy(ctx context.Context, name *string) ([]*dto.Terminology, error) {
	return a.MockSearchAllergyFn(ctx, name)
}
