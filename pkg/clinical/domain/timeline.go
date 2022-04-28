package domain

// HealthTimelineInput is the input for fetching a patient's health timeline
type HealthTimelineInput struct {
	PatientID string `json:"patientID"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

// HealthTimeline represents a health timeline containing various FHIR resources
type HealthTimeline struct {
	Timeline   []map[string]interface{} `json:"timeline"`
	TotalCount int                      `json:"totalCount"`
}
