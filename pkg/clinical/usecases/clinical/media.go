package clinical

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// UploadMedia uploads media to GCS and creates the resource in FHIR
func (c *UseCasesClinicalImpl) UploadMedia(ctx context.Context, encounterID string, file io.Reader, contentType string) (*dto.Media, error) {
	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	facility, err := c.infrastructure.FHIR.GetFHIROrganization(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	encounter, err := c.infrastructure.FHIR.GetFHIREncounter(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	patientID := encounter.Resource.Subject.ID
	patientReference := fmt.Sprintf("Patient/%s", *patientID)
	mediaObjectName := fmt.Sprintf("%s@%s", patientReference, time.Now())

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, *patientID)
	if err != nil {
		return nil, err
	}

	mediaUploadOutput, err := c.infrastructure.Upload.UploadMedia(ctx, mediaObjectName, file, contentType)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id := uuid.New().String()
	mediaSystem := scalarutils.URI("http://terminology.hl7.org/CodeSystem/media-type")
	encounterReference := fmt.Sprintf("Encounter/%s", *encounter.Resource.ID)

	mediaInput := &domain.FHIRMedia{
		ID: &id,
		Identifier: []*domain.FHIRIdentifier{
			{
				Use:    "official",
				System: &mediaSystem,
			},
		},
		Status: domain.MediaStatusCompleted,
		Subject: &domain.FHIRReferenceInput{
			ID:        patientID,
			Reference: &patientReference,
			Display:   patient.Resource.Names(),
		},
		Encounter: &domain.FHIRReferenceInput{
			ID:        &encounterID,
			Reference: &encounterReference,
		},
		Content: &domain.FHIRAttachmentInput{
			ContentType: (*scalarutils.Code)(&mediaUploadOutput.ContentType),
			URL:         (*scalarutils.URL)(&mediaUploadOutput.URL),
			Title:       &mediaUploadOutput.Name,
		},
		Issued: &now,
		Height: 465,
		Width:  345,
		Frames: 490,
	}

	orgRef := fmt.Sprintf("Organization/%s", *facility.Resource.ID)
	orgType := scalarutils.URI("Organization")

	mediaInput.Operator = &domain.FHIRReferenceInput{
		ID:        facility.Resource.ID,
		Reference: &orgRef,
		Display:   *facility.Resource.Name,
		Type:      &orgType,
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return nil, err
	}

	mediaInput.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	media, err := c.infrastructure.FHIR.CreateFHIRMedia(ctx, *mediaInput)
	if err != nil {
		return nil, err
	}

	output := &dto.Media{
		PatientID:   *patientID,
		PatientName: patient.Resource.Names(),
		URL:         string(*media.Content.URL),
		Name:        *media.Content.Title,
		ContentType: mediaUploadOutput.ContentType,
	}

	return output, nil
}

// ListPatientMedia list the patients media resources
func (c *UseCasesClinicalImpl) ListPatientMedia(ctx context.Context, patientID string, pagination dto.Pagination) (*dto.MediaConnection, error) {
	err := pagination.Validate()
	if err != nil {
		return nil, err
	}

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	patientReference := fmt.Sprintf("Patient/%s", *patient.Resource.ID)

	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant identifiers from context: %w", err)
	}

	mediaResources, err := c.infrastructure.FHIR.SearchPatientMedia(ctx, patientReference, *identifiers, pagination)
	if err != nil {
		return nil, err
	}

	patientMediaList := []*dto.Media{}

	for _, mediaResponse := range mediaResources.Media {
		patientMediaList = append(patientMediaList, mapFHIRMediaToMediaDTO(mediaResponse))
	}

	pageInfo := dto.PageInfo{
		HasNextPage:     mediaResources.HasNextPage,
		EndCursor:       &mediaResources.NextCursor,
		HasPreviousPage: mediaResources.HasPreviousPage,
		StartCursor:     &mediaResources.PreviousCursor,
	}

	connection := dto.CreateMediaConnection(patientMediaList, pageInfo, mediaResources.TotalCount)

	return &connection, nil
}

func mapFHIRMediaToMediaDTO(fhirMedia domain.FHIRMedia) *dto.Media {
	media := &dto.Media{
		ID:          *fhirMedia.ID,
		PatientID:   *fhirMedia.Subject.ID,
		PatientName: fhirMedia.Subject.Display,
		URL:         string(*fhirMedia.Content.URL),
		Name:        *fhirMedia.Content.Title,
		ContentType: string(*fhirMedia.Content.ContentType),
	}

	return media
}
