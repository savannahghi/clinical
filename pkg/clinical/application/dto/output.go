package dto

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
