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
	// SendEmailEndpoint is the endpoint used to send an email
	SendEmailEndpoint = "internal/send_email"
	// SendSMSEndpoint is the endpoint used to send sms
	SendSMSEndpoint = "internal/send_sms"
	// EmailWelcomeSubject is the subject of the welcome email
	EmailWelcomeSubject = "Welcome to Be.Well"
	// DefaultLanguage ...
	DefaultLanguage = "English"
	// CallCenterNumber is Savannah's call center number
	CallCenterNumber = "0790 360 360"

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

	// MedicalDataCount is the count of medical records
	MedicalDataCount = "3"

	// WeightCIELTerminologyCode is the terminology code for weight
	WeightCIELTerminologyCode = "5089"

	// HeightCIELTerminologyCode is the terminology code height
	HeightCIELTerminologyCode = "5090"

	// TemperatureCIELTerminologyCode is the terminology code for temperature
	TemperatureCIELTerminologyCode = "5088"

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
