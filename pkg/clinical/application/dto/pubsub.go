package dto

import "time"

// This file will be used to hold the payload for data that has been published to
// different pubsub topics

// CreatePatientPubSubMessage models the payload that is published to the `patient.create`
// topic
type CreatePatientPubSubMessage struct {
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
}
