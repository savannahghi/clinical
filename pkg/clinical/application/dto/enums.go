package dto

type OrganizationIdentifierType string

const (
	SladeCode OrganizationIdentifierType = "SladeCode"
	MFLCode   OrganizationIdentifierType = "MFLCode"
	ProgramID OrganizationIdentifierType = "MCHProgram"
	Other     OrganizationIdentifierType = "Other"
)

type EpisodeOfCareStatusEnum string

const (
	EpisodeOfCareStatusEnumPlanned   EpisodeOfCareStatusEnum = "planned"
	EpisodeOfCareStatusEnumActive    EpisodeOfCareStatusEnum = "active"
	EpisodeOfCareStatusEnumFinished  EpisodeOfCareStatusEnum = "finished"
	EpisodeOfCareStatusEnumCancelled EpisodeOfCareStatusEnum = "cancelled"
)

// EncounterStatusEnum is a FHIR enum
type EncounterStatusEnum string

const (
	// EncounterStatusEnumPlanned ...
	EncounterStatusEnumPlanned EncounterStatusEnum = "planned"
	// EncounterStatusEnumArrived ...
	EncounterStatusEnumArrived EncounterStatusEnum = "arrived"
	// EncounterStatusEnumTriaged ...
	EncounterStatusEnumTriaged EncounterStatusEnum = "triaged"
	// EncounterStatusEnumInProgress ...
	EncounterStatusEnumInProgress EncounterStatusEnum = "in_progress"
	// EncounterStatusEnumOnleave ...
	EncounterStatusEnumOnleave EncounterStatusEnum = "onleave"
	// EncounterStatusEnumFinished ...
	EncounterStatusEnumFinished EncounterStatusEnum = "finished"
	// EncounterStatusEnumCancelled ...
	EncounterStatusEnumCancelled EncounterStatusEnum = "cancelled"
	// EncounterStatusEnumEnteredInError ...
	EncounterStatusEnumEnteredInError EncounterStatusEnum = "entered_in_error"
	// EncounterStatusEnumUnknown ...
	EncounterStatusEnumUnknown EncounterStatusEnum = "unknown"
)

type EncounterClass string

const (
	// Also referred to as outpatient - For now we'll start with outpatient only
	EncounterClassAmbulatory EncounterClass = "ambulatory"
)

type ResourceType string

const (
	ResourceTypeAllergyIntolerance  ResourceType = "AllergyIntolerance"
	ResourceTypeObservation         ResourceType = "Observation"
	ResourceTypeCondition           ResourceType = "Condition"
	ResourceTypeMedicationStatement ResourceType = "MedicationStatement"
)

type AllergyIntoleranceReactionSeverityEnum string

const (
	AllergyIntoleranceReactionSeverityEnumMild     AllergyIntoleranceReactionSeverityEnum = "mild"
	AllergyIntoleranceReactionSeverityEnumModerate AllergyIntoleranceReactionSeverityEnum = "moderate"
	AllergyIntoleranceReactionSeverityEnumSevere   AllergyIntoleranceReactionSeverityEnum = "severe"
)

type ObservationStatus string

const (
	ObservationStatusFinal     ObservationStatus = "final"
	ObservationStatusCancelled ObservationStatus = "cancelled"
)

type MedicationStatementStatusEnum string

const (
	MedicationStatementStatusEnumActive   MedicationStatementStatusEnum = "active"
	MedicationStatementStatusEnumInActive MedicationStatementStatusEnum = "inactive"
	MedicationStatementStatusEnumUnknown  MedicationStatementStatusEnum = "unknown"
)

type IdentifierType string

const (
	IdentifierTypeNationalID IdentifierType = "NATIONAL_ID"
	IdentifierTypePassport   IdentifierType = "PASSPORT"
	IdentifierTypeAlienID    IdentifierType = "ALIEN_ID"
	IdentifierTypeCCCNumber  IdentifierType = "CCC_NUMBER"
)

type ContactType string

const (
	ContactTypePhoneNumber ContactType = "PHONE_NUMBER"
)

// Gender is a FHIR enum
type Gender string

const (
	// GenderMale ...
	GenderMale Gender = "male"
	// GenderFemale ...
	GenderFemale Gender = "female"
	// GenderOther ...
	GenderOther Gender = "other"
)

// ConditionStatus represents status of a FHIR condition
type ConditionStatus string

const (
	ConditionStatusActive   ConditionStatus = "active"
	ConditionStatusInactive ConditionStatus = "inactive"
	ConditionStatusResolved ConditionStatus = "resolved"
)

// TerminologySource represents various concept sources
type TerminologySource string

const (
	TerminologySourceICD10    TerminologySource = "ICD10"
	TerminologySourceCIEL     TerminologySource = "CIEL"
	TerminologySourceSNOMEDCT TerminologySource = "SNOMED_CT"
	TerminologySourceLOINC    TerminologySource = "LOINC"
)
