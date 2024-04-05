package clinical

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/savannahghi/clinical/pkg/clinical/application/common/helpers"
	"github.com/savannahghi/clinical/pkg/clinical/application/utils"
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
	Date           string
	Time           string
	Reason         string
	Patient        Patient
	NextOfKin      NextOfKin
	Facility       Facility
	Referral       Referral
	MedicalHistory MedicalHistory
	Footer         Footer
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

	data := TemplateData{
		Date:      time.Now().Format("Monday Jan 2 2023"),
		Time:      time.Now().Format("15:04"),
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
		MedicalHistory: MedicalHistory{Procedure: "Screening", Medication: "None", ReferralNotes: referralReason, Tests: []Test{{Name: "VIA", Results: "Positive", Date: "13th May 2024"}}},
		Footer:         Footer{},
	}

	tmpl, err := template.ParseFiles("templates/referral_report_template.html")
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	var htmlBuffer bytes.Buffer

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

	_, err = c.infrastructure.Upload.UploadMedia(ctx, "referral_report_1", bytes.NewReader(pdfBytes), "")
	if err != nil {
		utils.ReportErrorToSentry(err)
		return nil, err
	}

	return pdfBytes, nil
}
