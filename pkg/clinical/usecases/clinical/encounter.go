package clinical

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
)

// StartEncounter starts an encounter within an episode of care
func (c *UseCasesClinicalImpl) StartEncounter(ctx context.Context, episodeID string) (string, error) {
	if episodeID == "" {
		return "", fmt.Errorf("an episode of care ID is required")
	}

	episodeOfCare, err := c.infrastructure.FHIR.GetFHIREpisodeOfCare(ctx, episodeID)
	if err != nil {
		return "", err
	}

	encounterClassCode := scalarutils.Code("AMB")
	encounterClassSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/v3-ActCode")
	encounterClassVersion := "2018-08-12"
	encounterClassDisplay := string(dto.EncounterClassEnumAmbulatory)
	encounterClassUserSelected := false
	now := time.Now()
	startTime := scalarutils.DateTime(now.Format("2006-01-02T15:04:05+03:00"))

	episodeReference := fmt.Sprintf("EpisodeOfCare/%s", *episodeOfCare.Resource.ID)
	encounterPayload := domain.FHIREncounterInput{
		Status: domain.EncounterStatusEnum(dto.EncounterStatusEnumInProgress),
		Class: domain.FHIRCodingInput{
			System:       &encounterClassSystem,
			Version:      &encounterClassVersion,
			Code:         encounterClassCode,
			Display:      encounterClassDisplay,
			UserSelected: &encounterClassUserSelected,
		},
		Subject: &domain.FHIRReferenceInput{
			ID:        episodeOfCare.Resource.Patient.ID,
			Reference: episodeOfCare.Resource.Patient.Reference,
			Type:      episodeOfCare.Resource.Patient.Type,
			Display:   episodeOfCare.Resource.Patient.Display,
		},
		EpisodeOfCare: []*domain.FHIRReferenceInput{
			{
				ID:        episodeOfCare.Resource.ID,
				Reference: &episodeReference,
			},
		},
		ServiceProvider: &domain.FHIRReferenceInput{
			Reference: episodeOfCare.Resource.ManagingOrganization.Reference,
			Type:      episodeOfCare.Resource.ManagingOrganization.Type,
			Display:   episodeOfCare.Resource.ManagingOrganization.Display,
		},
		Period: &domain.FHIRPeriodInput{
			Start: startTime,
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return "", err
	}

	encounterPayload.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	encounter, err := c.infrastructure.FHIR.CreateFHIREncounter(ctx, encounterPayload)
	if err != nil {
		return "", err
	}

	// Create a blank composition
	_, err = c.CreateInitialComposition(ctx, encounter, tags, episodeOfCare)
	if err != nil {
		return "", err
	}

	// Create Subscription
	_, err = c.CreateSubscription(ctx, *encounter.Resource.ID, tags)
	if err != nil {
		return "", err
	}

	return *encounter.Resource.ID, nil
}

// CreateInitialComposition this method is to be specifically used when a new encounter is created.
// We create the initial(blank) composition so that we can use the composition to create a patients 'clinical document'
// that can be used for various purposes such as referral.
// This resource is to be updated on need basis such as (but not limited to) when a new diagnostic resource, observation etc. is created
// We want to consolidate the patients information into a single source for ease of retrieval and usage across
// depending on the current business case.
func (c *UseCasesClinicalImpl) CreateInitialComposition(ctx context.Context,
	encounter *domain.FHIREncounterRelayPayload,
	tags []domain.FHIRCodingInput, episodeOfCare *domain.FHIREpisodeOfCareRelayPayload) (*domain.FHIRCompositionRelayPayload, error) {
	encounterRef := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)
	encounterType := scalarutils.URI("Encounter")

	today := time.Now()

	date, err := scalarutils.NewDate(today.Day(), int(today.Month()), today.Year())
	if err != nil {
		return nil, err
	}

	preliminaryStatus := domain.CompositionStatusEnumPreliminary

	organizationRef := fmt.Sprintf("Organization/%s", *encounter.Resource.Meta.Tag[0].Code)

	compositionCategoryCode, err := c.mapCategoryEnumToCode(dto.ProviderUnspecifiedProgressNote)
	if err != nil {
		return nil, err
	}

	compositionConcept, err := c.mapCompositionConcepts(ctx, compositionCategoryCode, common.LOINCProviderUnspecifiedProgressNote)
	if err != nil {
		return nil, err
	}

	compositionTitle := fmt.Sprintf("%s's %s", episodeOfCare.Resource.Patient.Display, compositionConcept.CompositionCategoryConcept.DisplayName)

	compositionInput := domain.FHIRCompositionInput{
		Status: &preliminaryStatus,
		Meta: &domain.FHIRMetaInput{
			Tag: tags,
		},
		Type: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&compositionConcept.CompositionTypeConcept.URL),
					Code:    scalarutils.Code(compositionConcept.CompositionTypeConcept.ID),
					Display: compositionConcept.CompositionTypeConcept.DisplayName,
				},
			},
			Text: compositionConcept.CompositionTypeConcept.DisplayName,
		},
		Subject: &domain.FHIRReferenceInput{
			ID:        episodeOfCare.Resource.Patient.ID,
			Reference: episodeOfCare.Resource.Patient.Reference,
			Type:      episodeOfCare.Resource.Patient.Type,
			Display:   episodeOfCare.Resource.Patient.Display,
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        encounter.Resource.ID,
			Reference: &encounterRef,
			Display:   *encounter.Resource.ID,
			Type:      &encounterType,
		},
		Date: date,
		Author: []*domain.FHIRReferenceInput{
			{
				Reference: &organizationRef,
			},
		},
		Title: &compositionTitle,
	}

	output, err := c.infrastructure.FHIR.CreateFHIRComposition(ctx, compositionInput)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// CreateSubscription is used to define a push-based subscription from a server to another system.
// Once a subscription is registered with the server, the server checks every resource that is created
// or updated, and if the resource matches the given criteria, it sends a message on the defined "channel" so that another system can take an appropriate action.
//
// An example could be, update a composition resource when a new observation is created.
func (c *UseCasesClinicalImpl) CreateSubscription(ctx context.Context, encounterID string, tags []domain.FHIRCodingInput) (*domain.FHIRSubscription, error) {
	payload := "application/fhir+json"
	baseURL := serverutils.MustGetEnvVar("SERVICE_HOST")
	endpoint := fmt.Sprintf("%s/api/v1/subscriptions", baseURL)
	criteria := fmt.Sprintf("Observation?encounter=%s", encounterID)

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	clinicalOrgHeader := fmt.Sprintf("Clinical-Organization-ID: %s", identifiers.OrganizationID)
	facilityHeader := fmt.Sprintf("Clinical-Facility-ID: %s", identifiers.FacilityID)

	headers := []string{clinicalOrgHeader, facilityHeader}

	subscriptionInput := &domain.FHIRSubscriptionInput{
		Meta: &domain.FHIRMetaInput{
			Tag: tags,
		},
		Status:   domain.SubscriptionStatusRequested,
		Reason:   "Subscription is used to assemble together a single logical 'document' that provides a single and coherent statement of meaning",
		Criteria: criteria,
		Channel: &domain.FHIRSubscriptionChannel{
			Type:     domain.SubscriptionTypeRestHook,
			Endpoint: &endpoint,
			Payload:  &payload,
			Header:   headers,
		},
	}

	subscription, err := c.infrastructure.FHIR.CreateFHIRSubscription(ctx, subscriptionInput)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// PatchEncounter is used to update the details of an encounter resource
func (c *UseCasesClinicalImpl) PatchEncounter(ctx context.Context, encounterID string, input dto.EncounterInput) (*dto.Encounter, error) {
	if encounterID == "" {
		return nil, fmt.Errorf("an encounterID is required")
	}

	status := domain.EncounterStatusEnum(strings.ToLower(string(input.Status)))
	encounterInput := domain.FHIREncounterInput{
		Status: status,
	}

	if status.IsFinal() {
		encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, encounterID)
		if err != nil {
			return nil, fmt.Errorf("unable to get encounter with ID %s: %w", encounterID, err)
		}

		var startTime scalarutils.DateTime
		if encounter.Resource.Period == nil {
			startTime = scalarutils.DateTime(time.Now().Format(timeFormatStr))
		} else {
			startTime = encounter.Resource.Period.Start
		}

		// workaround for odd date comparison behavior on the Google Cloud Healthcare API
		end := startTime.Time().Add(time.Hour * 24)
		endTime := scalarutils.DateTime(end.Format(timeFormatStr))

		encounterInput.Period = &domain.FHIRPeriodInput{Start: startTime, End: endTime}
	}

	fhirEncounter, err := c.infrastructure.FHIR.PatchFHIREncounter(ctx, encounterID, encounterInput)

	if err != nil {
		return nil, err
	}

	encounters := []*dto.Encounter{}

	err = mapstructure.Decode([]domain.FHIREncounter{*fhirEncounter}, &encounters)

	if err != nil {
		return nil, err
	}

	return encounters[0], nil
}

// EndEncounter marks an encounter as finished and updates the endtime field
func (c *UseCasesClinicalImpl) EndEncounter(ctx context.Context, encounterID string) (bool, error) {
	if encounterID == "" {
		return false, fmt.Errorf("an encounterID is required")
	}

	ok, err := c.infrastructure.FHIR.EndEncounter(ctx, encounterID)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// ListPatientEncounters lists all the encounters that a patient has been part of
func (c *UseCasesClinicalImpl) ListPatientEncounters(ctx context.Context, patientID string, pagination *dto.Pagination) (*dto.EncounterConnection, error) {
	if patientID == "" {
		return nil, fmt.Errorf("a patient ID is required")
	}

	err := pagination.Validate()
	if err != nil {
		return nil, err
	}

	_, err = c.infrastructure.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", patientID)

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	encounterResponses, err := c.infrastructure.FHIR.SearchPatientEncounters(ctx, patientReference, nil, *identifiers, *pagination)
	if err != nil {
		return nil, err
	}

	encounters := []*dto.Encounter{}

	err = mapstructure.Decode(encounterResponses.Encounters, &encounters)
	if err != nil {
		return nil, err
	}

	pagedInfo := dto.PageInfo{
		HasNextPage:     encounterResponses.HasNextPage,
		EndCursor:       &encounterResponses.NextCursor,
		HasPreviousPage: encounterResponses.HasPreviousPage,
		StartCursor:     &encounterResponses.PreviousCursor,
	}

	connection := dto.CreateEncounterConnection(encounters, pagedInfo, encounterResponses.TotalCount)

	return &connection, nil
}

// GetEncounterAssociatedResources get all resources assocuated with an encounter
func (c *UseCasesClinicalImpl) GetEncounterAssociatedResources(ctx context.Context, encounterID string) (*dto.EncounterAssociatedResourceOutput, error) {
	if encounterID == "" {
		return nil, fmt.Errorf("an encounterID is required")
	}

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	includedResources := []string{
		"RiskAssessment:encounter",
		"Consent:data",
		"Observation:encounter",
	}

	encounterSearchParams := map[string]interface{}{
		"_id":         encounterID,
		"_sort":       "_lastUpdated",
		"_count":      "1",
		"_revinclude": includedResources,
	}

	encounterAllData, err := c.infrastructure.FHIR.SearchFHIREncounterAllData(ctx, encounterSearchParams, *identifiers, dto.Pagination{})
	if err != nil {
		return nil, err
	}

	result := dto.EncounterAssociatedResources{}

	for _, encounterData := range encounterAllData.Resources {
		switch encounterData["resourceType"] {
		case "RiskAssessment":
			var riskAssessment dto.RiskAssessment

			riskAssessmentBytes, err := json.Marshal(encounterData)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(riskAssessmentBytes, &riskAssessment); err != nil {
				return nil, err
			}

			result.RiskAssessment = append(result.RiskAssessment, &riskAssessment)
		case "Consent":
			var consent dto.Consent

			consentBytes, err := json.Marshal(encounterData)

			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(consentBytes, &consent); err != nil {
				return nil, err
			}

			result.Consent = append(result.Consent, &consent)

		case "Observation":
			var observation domain.FHIRObservation

			var observationValue, observationNote string

			observationBytes, err := json.Marshal(encounterData)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(observationBytes, &observation); err != nil {
				return nil, err
			}

			if observation.ValueString != nil {
				observationValue = *observation.ValueString
			}

			if observation.Note != nil {
				observationNote = string(*observation.Note[0].Text)
			}

			result.Observation = append(result.Observation, &dto.Observation{
				ID:           *observation.ID,
				Name:         observation.Code.Text,
				Value:        observationValue,
				Status:       dto.ObservationStatus(*observation.Status),
				TimeRecorded: string(*observation.EffectiveInstant),
				Note:         observationNote,
			})
		}
	}

	output := &dto.EncounterAssociatedResourceOutput{}
	if len(result.RiskAssessment) > 0 {
		output.RiskAssessment = result.RiskAssessment[0]
	}

	if len(result.Consent) > 0 {
		output.Consent = result.Consent[0]
	}

	if len(result.Observation) > 0 {
		output.Observation = result.Observation[0]
	}

	return output, nil
}
