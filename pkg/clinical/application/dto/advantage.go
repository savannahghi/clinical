package dto

// SegmentationPayload is used to stratify clients in advantage EMR.
type SegmentationPayload struct {
	// ClinicalID represents the patient's ID in this service
	ClinicalID   string               `json:"clinical_id,omitempty"`
	SegmentLabel SegmentationCategory `json:"segment_label,omitempty"`
}

// SMSPayload is used to model data used to send sms to patients through the advantage service
type SMSPayload struct {
	Intention  string   `json:"intention"`
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
}
