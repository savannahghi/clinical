package dto

// ServiceRequest is a record of a request for a procedure or diagnostic or other service to be planned
type ServiceRequest struct {
	ID        string       `json:"id,omitempty"`
	Status    string       `json:"status,omitempty"`
	Intent    string       `json:"intent,omitempty"`
	Priority  string       `json:"priority,omitempty"`
	Note      []Annotation `json:"note,omitempty"`
	Subject   Reference    `json:"subject,omitempty"`
	Encounter *Reference   `json:"encounter,omitempty"`
}
