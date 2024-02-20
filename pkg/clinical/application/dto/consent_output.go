package dto

// Consent models a fhir consent resource.
type Consent struct {
	ID        *string            `json:"id,omitempty"`
	Status    *ConsentStatusEnum `json:"status"`
	Provision *ConsentProvision  `json:"provision,omitempty"`
	Patient   *Reference         `json:"patient,omitempty"`
}

// ConsentProvision models a  consent provision
type ConsentProvision struct {
	ID   *string                   `json:"id,omitempty"`
	Type *ConsentProvisionTypeEnum `json:"type,omitempty"`
}
