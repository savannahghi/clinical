package clinical

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/dto"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

type Patient struct {
	Name        string
	EmpowerID   string
	NationalID  string
	PhoneNumber string
	DateOfBirth string
	Age         int
	Sex         string
}

type NextOfKin struct {
	Name         string
	Relationship string
	PhoneNumber  string
}

type Facility struct {
	Name     string
	Location string
	Contact  string
}

type Referral struct {
	Reason string
}

type Test struct {
	Name    string
	Results string
	Date    string
}

type MedicalHistory struct {
	Procedure     string
	Medication    string
	ReferralNotes string
	Tests         []Test
}

type Footer struct {
	Phone   string
	Email   string
	Address string
}

type TemplateData struct {
	Date              string
	Time              string
	Reason            string
	ReferringFacility Facility
	Patient           Patient
	NextOfKin         NextOfKin
	Facility          Facility
	Referral          Referral
	MedicalHistory    MedicalHistory
	Footer            Footer
}

// DocumentReferencePayload models data used to create caller specific reference document payload
type DocumentReferencePayload struct {
	Subject           *domain.FHIRReferenceInput
	Attachment        *domain.FHIRAttachment
	Related           *domain.FHIRReference
	TerminologySystem string
}

// GenerateReferralReportPDF generates a PDF report for a given referral.
//
// The serviceRequestID is unique to each ServiceRequest resource, which,
// according to FHIR standards, is how referrals are represented. In FHIR, a referral
// is a specific type of ServiceRequest, which typically contains details such
// as the requester, the patient, the requested service, and other clinical information.
//
// By leveraging the serviceRequestID, this function retrieves the associated ServiceRequest
// from the FHIR server. It then extracts relevant data, including patient and encounter
// information, to construct a comprehensive referral report. The report is formatted
// as a PDF, making it suitable for clinical review, record-keeping, or sharing with
// other healthcare professionals.
func (c *UseCasesClinicalImpl) GenerateReferralReportPDF(ctx context.Context, serviceRequestID string) ([]byte, error) {
	if serviceRequestID == "" {
		return nil, fmt.Errorf("service request ID cannot be empty")
	}

	serviceRequest, err := c.infrastructure.FHIR.GetFHIRServiceRequest(ctx, serviceRequestID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	patient, err := c.infrastructure.FHIR.GetFHIRPatient(ctx, *serviceRequest.Resource.Subject.ID)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	var nationalID string

	for _, identifier := range patient.Resource.Identifier {
		system := scalarutils.URI(helpers.IDIdentifierSystem)
		if identifier.System == &system && identifier.Value != "" {
			nationalID = identifier.Value
		}
	}

	age := time.Since(patient.Resource.BirthDate.AsTime()).Hours() / 24 / 365
	patientData := Patient{
		Name:        patient.Resource.Name[0].Text,
		EmpowerID:   "",
		NationalID:  nationalID,
		PhoneNumber: *patient.Resource.Telecom[0].Value,
		DateOfBirth: patient.Resource.BirthDate.String(),
		Age:         int(age),
		Sex:         patient.Resource.Gender.String(),
	}

	var referredFacilityName, referredFacilityCounty, referredFacilityContact string

	for _, extension := range serviceRequest.Resource.Extension {
		if extension.URL == "http://savannahghi.org/fhir/StructureDefinition/referred-facility" {
			for _, ext := range extension.Extension {
				if ext.URL == "facilityName" {
					referredFacilityName = ext.ValueString
				}

				if ext.URL == "facilityCounty" {
					referredFacilityCounty = ext.ValueString
				}

				if ext.URL == "facilityContact" {
					referredFacilityContact = ext.ValueString
				}
			}
		}
	}

	var referralReason string
	if len(serviceRequest.Resource.Note) > 0 && serviceRequest.Resource.Note[0].Text != nil {
		referralReason = string(*serviceRequest.Resource.Note[0].Text)
	}

	facilityID, err := extensions.GetFacilityIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	facility, err := c.infrastructure.FHIR.GetFHIROrganization(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	var facilityContact string

	if facility.Resource.Telecom != nil {
		for _, tel := range facility.Resource.Telecom {
			if tel.Value != nil {
				facilityContact = *tel.Value
				break
			}
		}
	}

	observations, err := c.GetPatientObservations(
		ctx,
		*patient.Resource.ID,
		serviceRequest.Resource.Encounter.ID,
		nil,
		"",
		&dto.Pagination{},
	)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	var tests []Test

	for _, observation := range observations.Edges {
		obs := Test{
			Name:    observation.Node.Name,
			Results: observation.Node.Value,
			Date:    observation.Node.TimeRecorded,
		}

		tests = append(tests, obs)
	}

	data := TemplateData{
		Date: time.Now().Format("Jan 2, 2006"),
		Time: time.Now().Format("15:04"),
		ReferringFacility: Facility{
			Name: *facility.Resource.Name,
		},
		Patient:   patientData,
		NextOfKin: NextOfKin{},
		Facility: Facility{
			Name:     referredFacilityName,
			Contact:  referredFacilityContact,
			Location: referredFacilityCounty,
		},
		Referral: Referral{
			Reason: referralReason,
		},
		MedicalHistory: MedicalHistory{Procedure: "Screening", Medication: "None", ReferralNotes: referralReason, Tests: tests},
		Footer: Footer{
			Phone: facilityContact,
		},
	}

	var htmlBuffer bytes.Buffer

	tmpl := template.Must(template.New("ReferralFormTemplate").Parse(utils.ReferralFormTemplate))

	err = tmpl.Execute(&htmlBuffer, data)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	htmlContent := htmlBuffer.String()

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlContent)))

	err = pdfg.Create()
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	pdfBytes := pdfg.Bytes()

	currentTime := time.Now().Format("20060102T150405")
	filename := fmt.Sprintf("%s_%s.pdf", patientData.Name, currentTime)

	result, err := c.infrastructure.Upload.UploadMedia(ctx, filename, bytes.NewReader(pdfBytes), "")
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	title := fmt.Sprintf("%s's Referral report", serviceRequest.Resource.Subject.Display)
	serviceRequestReference := fmt.Sprintf("ServiceRequest/%s", *serviceRequest.Resource.ID)

	payload := &DocumentReferencePayload{
		Subject: &domain.FHIRReferenceInput{
			ID:        serviceRequest.Resource.Subject.ID,
			Reference: serviceRequest.Resource.Subject.Reference,
		},
		Attachment: &domain.FHIRAttachment{
			ContentType: (*scalarutils.Code)(&result.ContentType),
			URL:         (*scalarutils.URL)(&result.SignedURL),
			Title:       &title,
		},
		Related: &domain.FHIRReference{
			Reference: &serviceRequestReference,
		},
		TerminologySystem: common.ReferralNoteLOINCTerminologySystem,
	}

	err = c.CreateDocumentReference(ctx, payload)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return pdfBytes, nil
}

// CreateDocumentReference is a helper method to abstract the creation of a document reference
func (c *UseCasesClinicalImpl) CreateDocumentReference(ctx context.Context, payload *DocumentReferencePayload) error {
	concept, err := c.GetConcept(ctx, dto.TerminologySourceLOINC, payload.TerminologySystem)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return err
	}

	finalDocStatus := domain.CompositionStatusEnumFinal
	status := domain.DocumentReferenceStatusEnumCurrent
	instant := scalarutils.Instant(time.Now().Format(time.RFC3339))

	documentReference := &domain.FHIRDocumentReferenceInput{
		Status:    status,
		DocStatus: &finalDocStatus,
		Type: &domain.FHIRCodeableConceptInput{
			Coding: []*domain.FHIRCodingInput{
				{
					System:  (*scalarutils.URI)(&concept.URL),
					Code:    scalarutils.Code(concept.ID),
					Display: concept.DisplayName,
				},
			},
			Text: concept.DisplayName,
		},
		Subject: payload.Subject,
		Date:    &instant,
		Content: []domain.FHIRDocumentReferenceContent{
			{
				Attachment: *payload.Attachment,
			},
		},
		Context: &domain.FHIRDocumentReferenceContext{
			Related: []*domain.FHIRReference{
				payload.Related,
			},
		},
	}

	tags, err := c.GetTenantMetaTags(ctx)
	if err != nil {
		return err
	}

	documentReference.Meta = &domain.FHIRMetaInput{
		Tag: tags,
	}

	_, err = c.infrastructure.FHIR.CreateFHIRDocumentReference(ctx, documentReference)
	if err != nil {
		utils.ReportErrorToSentry(err)
		return err
	}

	return nil
}

// ShareReferralForm is searched for a document reference associated with a service request, retrieves the document URL and sends it to the
// patient via SMS
func (c *UseCasesClinicalImpl) ShareReferralForm(ctx context.Context, serviceRequestID string) (bool, error) {
	identifiers, err := c.infrastructure.BaseExtension.GetTenantIdentifiers(ctx)
	if err != nil {
		return false, err
	}

	params := map[string]interface{}{
		"related": fmt.Sprintf("ServiceRequest/%s", serviceRequestID),
		"_sort":   "_lastUpdated",
		"_count":  "1",
	}

	output, err := c.infrastructure.FHIR.SearchFHIRDocumentReference(ctx, params, *identifiers, dto.Pagination{})
	if err != nil {
		utils.ReportErrorToSentry(err)
		return false, err
	}

	if len(output.DocumentReferences) == 0 {
		return false, errors.New("no document reference found")
	}

	// TODO: Send SMS here

	return true, nil
}
