package dto

// Observation is a minimal representation of a fhir Observation
type Observation struct {
	ID             string            `json:"id,omitempty"`
	Name           string            `json:"name,omitempty"`
	Value          string            `json:"value,omitempty"`
	Status         ObservationStatus `json:"status,omitempty"`
	PatientID      string            `json:"patientID,omitempty"`
	EncounterID    string            `json:"encounterID,omitempty"`
	TimeRecorded   string            `json:"timeRecorded,omitempty"`
	Interpretation []string          `json:"interpretation,omitempty"`
	Note           string            `json:"note,omitempty"`
}

// ObservationEdge is an observation edge
type ObservationEdge struct {
	Node   Observation
	Cursor string
}

// ObservationConnection  is an Observation Connection Type
type ObservationConnection struct {
	TotalCount int
	Edges      []ObservationEdge
	PageInfo   PageInfo
}

// CreateObservationConnection creates a connection that follows the GraphQl Cursor Connection Specification
func CreateObservationConnection(observations []*Observation, pageInfo PageInfo, total int) ObservationConnection {
	connection := ObservationConnection{
		TotalCount: total,
		Edges:      []ObservationEdge{},
		PageInfo:   pageInfo,
	}

	for _, observation := range observations {
		edge := ObservationEdge{
			Node:   *observation,
			Cursor: observation.ID,
		}

		connection.Edges = append(connection.Edges, edge)
	}

	return connection
}
