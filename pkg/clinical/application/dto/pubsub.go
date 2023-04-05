package dto

import "time"

// This file will be used to hold the payload for data that has been published to
// different pubsub topics

// PatientPubSubMessage models the payload that is published to the `patient.create`
// topic
type PatientPubSubMessage struct {
	ID                string     `json:"id"`
	Active            bool       `json:"active"`
	ClientType        string     `json:"client_type"`
	EnrollmentDate    *time.Time `json:"enrollment_date"`
	FHIRPatientID     string     `json:"fhir_patient_id"`
	EMRHealthRecordID string     `json:"emr_health_record_id"`
	TreatmentBuddy    string     `json:"treatment_buddy"`
	Counselled        bool       `json:"counselled"`
	Organisation      string     `json:"organisation"`
	UserID            string     `json:"user"`
	CurrentFacilityID string     `json:"current_facility"`
	CHV               string     `json:"chv"`
	Caregiver         string     `json:"caregiver"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `jsoon:"facilityID"`
}

// FacilityPubSubMessage models the details of healthcare facilities that are on myCareHub platform.
// This will be used to create a FHIR organization
type FacilityPubSubMessage struct {
	ID          *string `json:"id"`
	Name        string  `json:"name"`
	Code        int     `json:"code"`
	Phone       string  `json:"phone"`
	Active      bool    `json:"active"`
	County      string  `json:"county"`
	Description string  `json:"description"`
}

// VitalSignPubSubMessage models the details that will be posted to the vitals pub/sub topic
type VitalSignPubSubMessage struct {
	Name      string    `json:"name"`
	ConceptID *string   `json:"conceptId"`
	Value     string    `json:"value"`
	Date      time.Time `json:"date"`

	PatientID string `json:"patientID"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `jsoon:"facilityID"`
}

// PatientAllergyPubSubMessage contains allergy details for a patient
type PatientAllergyPubSubMessage struct {
	Name      string          `json:"name"`
	ConceptID *string         `json:"conceptID"`
	Date      time.Time       `json:"date"`
	Reaction  AllergyReaction `json:"reaction"`
	Severity  AllergySeverity `json:"severity"`

	PatientID string `json:"patientID"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `jsoon:"facilityID"`
}

// AllergyReaction ...
type AllergyReaction struct {
	Name      string  `json:"name"`
	ConceptID *string `json:"conceptID"`
}

// AllergySeverity ...
type AllergySeverity struct {
	Name      string  `json:"name"`
	ConceptID *string `json:"conceptID"`
}

// MedicationPubSubMessage contains details for medication that a patient/client is prescribed or using
type MedicationPubSubMessage struct {
	Name      string          `json:"medication"`
	ConceptID *string         `json:"conceptId"`
	Date      time.Time       `json:"date"`
	Value     string          `json:"value"`
	Drug      *MedicationDrug `json:"drug"`

	PatientID string `json:"patientID"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `jsoon:"facilityID"`
}

// MedicationDrug ...
type MedicationDrug struct {
	ConceptID *string `json:"conceptId"`
}

// PatientTestResultPubSubMessage models details that are published to the test results topic
type PatientTestResultPubSubMessage struct {
	Name      string     `json:"name"`
	ConceptID *string    `json:"conceptId"`
	Date      time.Time  `json:"date"`
	Result    TestResult `json:"result"`

	PatientID string `json:"patientID"`

	OrganizationID string `json:"organizationID"`
	FacilityID     string `jsoon:"facilityID"`
}

// TestResult ...
type TestResult struct {
	Name      string  `json:"name"`
	ConceptID *string `json:"conceptId"`
}
