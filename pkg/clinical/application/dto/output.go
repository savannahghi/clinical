package dto

import (
	"time"

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

// TimelineResource represents a FHIR timeline resources
type TimelineResource struct {
	ID           string           `json:"id"`
	ResourceType ResourceType     `json:"resourceType"`
	Name         string           `json:"name"`
	Value        string           `json:"value"`
	Status       string           `json:"status"`
	Date         scalarutils.Date `json:"date"`
	TimeRecorded time.Time        `json:"timeRecorded"`
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

// Media is the output to show the results of the created media resource item
type Media struct {
	ID          string `json:"id,omitempty"`
	PatientID   string `json:"patientID,omitempty"`
	PatientName string `json:"patientName,omitempty"`
	URL         string `json:"url,omitempty"`
	Name        string `json:"name,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

// MediaEdge is an media connection edge
type MediaEdge struct {
	Node   Media  `json:"node,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// MediaConnection is a media connection
type MediaConnection struct {
	TotalCount int         `json:"totalCount,omitempty"`
	Edges      []MediaEdge `json:"edges,omitempty"`
	PageInfo   PageInfo    `json:"pageInfo,omitempty"`
}

// CreateMediaConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateMediaConnection(mediaList []*Media, pageInfo PageInfo, total int) MediaConnection {
	connection := MediaConnection{
		TotalCount: total,
		Edges:      []MediaEdge{},
		PageInfo:   pageInfo,
	}

	for _, media := range mediaList {
		edge := MediaEdge{
			Node:   *media,
			Cursor: media.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}

// Section defines a composition section output model
type Section struct {
	ID      *string    `json:"id,omitempty"`
	Title   *string    `json:"title,omitempty"`
	Code    *string    `json:"code,omitempty"`
	Author  *string    `json:"author,omitempty"`
	Text    string     `json:"text,omitempty"`
	Section []*Section `json:"section,omitempty"`
}

// Consent models a fhir consent resource
type ConsentOutput struct {
	Status *ConsentStatusEnum `json:"status"`
}
