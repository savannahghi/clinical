package common

import (
	"time"

	"github.com/rs/xid"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// constants and defaults
const (

	// CenturyHours is the number of hours in a (fictional) century of leap years
	CenturyHours                  = 878400
	timeFormatStr                 = "2006-01-02T15:04:05+03:00"
	healthCloudIdentifiers        = "healthcloud.identifiers"
	healthCloudIdentifiersVersion = "0.0.1"

	// CreatePatientTopic is the topic ID where patient data is published to
	CreatePatientTopic = "patient.create"

	// VitalsTopicName is the topic for publishing a patient's vital signs
	VitalsTopicName = "vitals.create"

	// AllergyTopicName is the topic for publishing a patient's allergy
	AllergyTopicName = "allergy.create"

	// MedicationTopicName is the topic for publishing a patient's medication
	MedicationTopicName = "medication.create"

	// TestResultTopicName is the topic for publishing a patient's test results
	TestResultTopicName = "test.result.create"

	// TestOrderTopicName is the topic for publishing a patient's test order
	TestOrderTopicName = "test.order.create"

	// OrganizationTopicName is the topic where organization(facility) details are published to
	OrganizationTopicName = "organization.create"

	// TenantTopicName is the topic where program is registered in clinical as a tenant
	TenantTopicName = "mycarehub.tenant.create"

	// MedicalDataCount is the count of medical records
	MedicalDataCount = "3"

	// WeightCIELTerminologyCode is the terminology code for weight
	WeightCIELTerminologyCode = "5089"

	// HeightCIELTerminologyCode is the terminology code height
	HeightCIELTerminologyCode = "5090"

	// TemperatureCIELTerminologyCode is the terminology code for temperature
	TemperatureCIELTerminologyCode = "5088"

	// MuacCIELTerminologyCode is the terminology code for mid-upper arm circumference
	MuacCIELTerminologyCode = "1343"

	// BloodSugarCIELTerminologyCode is the terminology code for blood sugar (Serum glucose)
	BloodSugarCIELTerminologyCode = "887"

	// DiastolicBloodPressureTerminologyCode is the terminology code for diastolic blood pressure
	DiastolicBloodPressureCIELTerminologyCode = "5086"

	// LastMenstrualPeriodCIELTerminologyCode is the terminology code for last menstrual period
	LastMenstrualPeriodCIELTerminologyCode = "1427"

	// Spoc2CIELTerminologyCode is the terminology code oxygen saturation
	OxygenSaturationCIELTerminologyCode = "5092"

	// RespiratoryRateCIELTerminologyCode is the terminology code for respiratory rate
	RespiratoryRateCIELTerminologyCode = "5242"

	// PulseCIELTerminologyCode is the terminology code for pulse
	PulseCIELTerminologyCode = "5087"

	// BloodPressureCIELTerminologyCode is the terminology code for blood pressure
	BloodPressureCIELTerminologyCode = "5085"

	// BMICIELTerminologyCode is the terminology code for Body Mass Index
	BMICIELTerminologyCode = "1342"

	// ViralLoadCIELTerminologyCode is the terminology code for Viral Load
	ViralLoadCIELTerminologyCode = "856"

	// CD4CountCIELTerminologyCode is the terminology code for CD$ Count
	CD4CountCIELTerminologyCode = "5497"

	// ClinicalServiceName defines the service where the topic is created
	ClinicalServiceName = "clinical"

	// MyCareHubServiceName defines the service where some of the topics have been created
	MyCareHubServiceName = "mycarehub"

	// TestTopicName is a topic that is used for testing purposes
	TestTopicName = "pubsub.testing.topic"

	// TopicVersion defines the topic version. That standard one is `v1`
	TopicVersion = "v1"

	// AddFHIRIDToPatientProfile is the topic name where the details to update a patient's FHIR ID will be posted to
	AddFHIRIDToPatientProfile = "patient.fhirid.update"

	// AddFHIRIDToFacility is the topic where details to update a facility's fhir ID will be published to
	AddFHIRIDToFacility = "facility.fhirid.update"

	// AddFHIRIDToProgram is the topic where details to update a program's fhir ID will be published to
	AddFHIRIDToProgram = "program.fhirid.update"

	// LOINCProgressNoteCode defines LOINC progress note terminology code
	LOINCProgressNoteCode = "81216-4"

	// LOINCAssessmentPlanCode defines LOINC assessment plan note terminology code
	LOINCAssessmentPlanCode = "51847-2"

	// LOINCHistoryOfPresentingIllness defines LOINC history of presenting illness note terminology code
	LOINCHistoryOfPresentingIllness = "10164-2"

	// LOINCSocialHistory defines LOINC social history note terminology code
	LOINCSocialHistory = "29762-2"

	// LOINCFamilyHistory defines LOINC family history note terminology code
	LOINCFamilyHistory = "10157-6"

	// LOINCExamination defines LOINC Examination note terminology code
	LOINCExamination = "29545-1"

	// LOINCPLANOFCARE defines LOINC Plan of care note terminology code
	LOINCPLANOFCARE = "18776-5"

	// ColposcopyCIELTerminologyCode is the terminology code for colposcopy findings
	ColposcopyCIELTerminologyCode = "162816"

	// VIACIELCode is the terminology code for a VIA test
	VIACIELCode = "164805"

	// VIAPositiveCIELCode is the terminology code for a positive VIA test
	VIAResultPositiveCIELCode = "703"

	// VIANegativeCIELCode is the terminology code for a negative VIA test
	VIAResultNegativeCIELCode = "664"

	// VIASuspiciousOfCancerCIELCode is the terminology code for a suspicious cancer VIA test
	VIAResultSuspiciousOfCancerCIELCode = "159008"

	// HPVCIELTerminologyCode is the terminology code used to represent HPV test.
	HPVCIELTerminologyCode = "1213"

	// PapSmearTerminologyCode is the terminology code used to represent pap smear test.
	PapSmearTerminologyCode = "154451"

	// HighRiskCIELCode represents the CIEL code for a high-risk condition
	HighRiskCIELCode = "166674"

	// LowRiskCIELCode represents the CIEL code for a low-risk condition
	LowRiskCIELCode = "166675"

	// CIELTerminologySystem is the identity of ciel terminology system
	CIELTerminologySystem = "https://CIELterminology.org"

	// MammogramTerminologyCode is the terminology code used to represent mammogram results.
	MammogramTerminologyCode = "163591"

	// BenignNeoplasmOfBreastOfSkinTerminologyCode is the terminology code used to represent benign of skin results.
	BenignNeoplasmOfBreastOfSkinTerminologyCode = "147661"
)

// DefaultIdentifier assigns a patient a code to function as their
// medical record number.
func DefaultIdentifier() *domain.FHIRIdentifierInput {
	xid := xid.New().String()
	system := scalarutils.URI(healthCloudIdentifiers)
	version := healthCloudIdentifiersVersion
	userSelected := false

	return &domain.FHIRIdentifierInput{
		Use:   domain.IdentifierUseEnumOfficial,
		Value: xid,
		Type: domain.FHIRCodeableConceptInput{
			Text: "MR",
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &system,
					Version:      &version,
					Code:         scalarutils.Code(xid),
					Display:      xid,
					UserSelected: &userSelected,
				},
			},
		},
		System: &system,
		Period: DefaultPeriodInput(),
	}
}

// DefaultPeriodInput sets up a period input covering roughly a century from when it's run
func DefaultPeriodInput() *domain.FHIRPeriodInput {
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * CenturyHours)

	return &domain.FHIRPeriodInput{
		Start: scalarutils.DateTime(now.Format(timeFormatStr)),
		End:   scalarutils.DateTime(farFuture.Format(timeFormatStr)),
	}
}

// DefaultPeriod sets up a period input covering roughly a century from when it's run
func DefaultPeriod() *domain.FHIRPeriod {
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * CenturyHours)

	return &domain.FHIRPeriod{
		Start: scalarutils.DateTime(now.Format(timeFormatStr)),
		End:   scalarutils.DateTime(farFuture.Format(timeFormatStr)),
	}
}
