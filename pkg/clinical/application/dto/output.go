package dto

import (
	"github.com/savannahghi/scalarutils"
)

// TenantIdentifiers models the json object used to store some of the tenant identifiers
type TenantIdentifiers struct {
	OrganizationID string `json:"organizationID,omitempty"`
	FacilityID     string `json:"facilityID,omitempty"`
}

type Organization struct {
	ID           string                   `json:"id"`
	Active       bool                     `json:"active"`
	Name         string                   `json:"name"`
	Identifiers  []OrganizationIdentifier `json:"identifiers"`
	PhoneNumbers []string                 `json:"phoneNumbers"`
}

type EpisodeOfCare struct {
	ID        string                  `json:"id"`
	Status    EpisodeOfCareStatusEnum `json:"status"`
	PatientID string                  `json:"patientID"`
}

// Encounter definition: an interaction between a patient and healthcare provider(s) for the purpose of providing healthcare service(s) or assessing the health status of a patient.
type Encounter struct {
	ID              string              `json:"id,omitempty"`
	Status          EncounterStatusEnum `json:"status,omitempty"`
	Class           EncounterClass      `json:"class,omitempty"`
	PatientID       string              `json:"patientID,omitempty"`
	EpisodeOfCareID string              `json:"episodeOfCareID,omitempty"`
}

// HealthTimelineInput is the input for fetching a patient's health timeline
type HealthTimelineInput struct {
	PatientID string `json:"patientID"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

// TimelineResource represents a FHIR timeline resource
type TimelineResource struct {
	ID           string           `json:"id"`
	ResourceType ResourceType     `json:"resourceType"`
	Name         string           `json:"name"`
	Value        string           `json:"value"`
	Status       string           `json:"status"`
	Date         scalarutils.Date `json:"date"`
}

// HealthTimeline represents a health timeline containing various FHIR resources
type HealthTimeline struct {
	Timeline   []TimelineResource `json:"timeline"`
	TotalCount int                `json:"totalCount"`
}

// AllergyIntolerance represents an allergy intolerance containing minimal FHIR resources
type AllergyIntolerance struct {
	ID              string                                 `json:"id"`
	PatientID       string                                 `json:"patientID"`
	EncounterID     string                                 `json:"encounterID"`
	OnsetDateTime   scalarutils.DateTime                   `json:"onsetDateTime"`
	Severity        AllergyIntoleranceReactionSeverityEnum `json:"severity"`
	SubstanceCode   string                                 `json:"substanceCode"`
	SubstanceSystem string                                 `json:"substanceSystem"`
}

// Observation is a minimal representation of a fhir Observation
type Observation struct {
	ID          string            `json:"id,omitempty"`
	Status      ObservationStatus `json:"status,omitempty"`
	PatientID   string            `json:"patientID,omitempty"`
	EncounterID string            `json:"encounterID,omitempty"`
	Name        string            `json:"name,omitempty"`
	Value       string            `json:"value,omitempty"`
}

// Medication is a minimal representation of a fhir Medication
type Medication struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// MedicationStatement is a minimal representation of a fhir MedicationStatement
type MedicationStatement struct {
	ID         string                        `json:"id"`
	Status     MedicationStatementStatusEnum `json:"status"`
	Medication Medication                    `json:"medication"`
	PatientID  string                        `json:"patientID"`
}

// MedicalData is a minimal representation of a fhir MedicalData
type MedicalData struct {
	Regimen   []*MedicationStatement
	Allergies []*AllergyIntolerance
	Weight    []*Observation
	BMI       []*Observation
	ViralLoad []*Observation
	CD4Count  []*Observation
}

type Patient struct {
	ID          string           `json:"id"`
	Active      bool             `json:"active"`
	Name        string           `json:"name"`
	PhoneNumber []string         `json:"phoneNumber"`
	Gender      Gender           `json:"gender"`
	BirthDate   scalarutils.Date `json:"birthDate"`
}
