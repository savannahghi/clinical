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
	Allergies []*Allergy
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

// Terminology models the OCL terminology output
type Terminology struct {
	Code   string            `json:"code"`
	System TerminologySource `json:"system"`
	Name   string            `json:"name"`
}

// TerminologyEdge is a terminology edge
type TerminologyEdge struct {
	Node   Terminology
	Cursor string
}

// TerminologyConnection is the terminology connection
type TerminologyConnection struct {
	TotalCount int               `json:"totalCount,omitempty"`
	Edges      []TerminologyEdge `json:"edges,omitempty"`
	PageInfo   PageInfo          `json:"pageInfo,omitempty"`
}
