package dto

import "github.com/savannahghi/scalarutils"

// Condition represents a FHIR condition
type Condition struct {
	ID     string          `json:"id"`
	Status ConditionStatus `json:"status"`
	Name   string          `json:"condition"`
	Code   string          `json:"code"`
	System string          `json:"system"`

	Category ConditionCategory `json:"category"`

	OnsetDate    *scalarutils.Date `json:"onsetDate"`
	RecordedDate scalarutils.Date  `json:"recordedDate"`

	Note string `json:"note"`

	PatientID   string `json:"patientID"`
	EncounterID string `json:"encounterID"`
}

// ConditionEdge is a condition edge
type ConditionEdge struct {
	Node   Condition
	Cursor string
}

// ConditionConnection  is a Condition Connection Type
type ConditionConnection struct {
	TotalCount int
	Edges      []ConditionEdge
	PageInfo   PageInfo
}

// CreateConditionConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateConditionConnection(conditions []Condition, pageInfo PageInfo, total int) ConditionConnection {
	connection := ConditionConnection{
		TotalCount: total,
		Edges:      []ConditionEdge{},
		PageInfo:   pageInfo,
	}

	for _, condition := range conditions {
		edge := ConditionEdge{
			Node:   condition,
			Cursor: condition.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
