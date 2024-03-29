package clinical

import (
	"context"
	"fmt"

	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// RecordConsent records a user consent
func (u *UseCasesClinicalImpl) RecordConsent(ctx context.Context, input dto.ConsentInput) (*dto.ConsentOutput, error) {
	encounter, err := u.infrastructure.FHIR.GetFHIREncounter(ctx, input.EncounterID)
	if err != nil {
		return nil, err
	}

	if encounter.Resource.Status == domain.EncounterStatusEnumFinished {
		return nil, fmt.Errorf("cannot create a consent in a finished encounter")
	}

	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	encounterReference := &domain.FHIRReference{
		ID:        encounter.Resource.ID,
		Reference: &encounterRef,
	}

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)
	subjectReference := &domain.FHIRReference{
		ID:        patientID,
		Reference: &patientReference,
	}

	scope := &domain.FHIRCodeableConcept{
		Text: "patient-privacy",
	}

	var system scalarutils.URI = "http://terminology.hl7.org/CodeSystem/consentcategorycodes"

	code := scalarutils.Code("acd")

	coding := &domain.FHIRCoding{
		System:  &system,
		Code:    &code,
		Display: "Advance Directive",
	}
	category := &domain.FHIRCodeableConcept{
		Text:   "Advance Directive",
		Coding: []*domain.FHIRCoding{coding},
	}
	policyRule := &domain.FHIRCodeableConcept{
		Text: "cric",
	}
	consentProvision := &domain.FHIRConsentProvision{
		Type: &input.Provision,
		Data: []domain.FHIRConsentProvisionData{
			{
				Meaning:   dto.ConsentDataMeaningRelated,
				Reference: encounterReference,
			},
		},
	}

	tags, err := u.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	consentMeta := domain.FHIRMetaInput{
		Tag: tags,
	}

	status := dto.ConditionStatusActive
	consent := domain.FHIRConsent{
		Provision:  consentProvision,
		Status:     (*dto.ConsentStatusEnum)(&status),
		Patient:    subjectReference,
		Scope:      scope,
		Category:   []*domain.FHIRCodeableConcept{category},
		PolicyRule: policyRule,
		Meta:       &consentMeta,
	}

	if input.Provision == dto.ConsentProvisionTypeDeny {
		extension := &domain.Extension{
			URL:         "http://hl7.org/fhir/StructureDefinition/event-statusReason",
			ValueString: input.DenyReason,
		}
		consent.Extension = []domain.Extension{*extension}
	}

	resp, err := u.infrastructure.FHIR.CreateFHIRConsent(ctx, consent)
	if err != nil {
		return nil, err
	}

	output := &dto.ConsentOutput{
		Status: resp.Status,
	}

	return output, nil
}
