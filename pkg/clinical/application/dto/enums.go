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
	EpisodeOfCareStatusEnumPlanned        EpisodeOfCareStatusEnum = "PLANNED"
	EpisodeOfCareStatusEnumActive         EpisodeOfCareStatusEnum = "ACTIVE"
	EpisodeOfCareStatusEnumFinished       EpisodeOfCareStatusEnum = "FINISHED"
	EpisodeOfCareStatusEnumCancelled      EpisodeOfCareStatusEnum = "CANCELLED"
	EpisodeOfCareStatusEnumEnteredInError EpisodeOfCareStatusEnum = "ENTERED_IN_ERROR"
)

// EncounterStatusEnum is a FHIR enum
type EncounterStatusEnum string

const (
	// EncounterStatusEnumPlanned ...
	EncounterStatusEnumPlanned EncounterStatusEnum = "PLANNED"
	// EncounterStatusEnumArrived ...
	EncounterStatusEnumArrived EncounterStatusEnum = "ARRIVED"
	// EncounterStatusEnumTriaged ...
	EncounterStatusEnumTriaged EncounterStatusEnum = "TRIAGED"
	// EncounterStatusEnumInProgress ...
	EncounterStatusEnumInProgress EncounterStatusEnum = "IN_PROGRESS"
	// EncounterStatusEnumFinished ...
	EncounterStatusEnumFinished EncounterStatusEnum = "FINISHED"
	// EncounterStatusEnumCancelled ...
	EncounterStatusEnumCancelled EncounterStatusEnum = "CANCELLED"
)

type EncounterClass string

const (
	// Also referred to as outpatient - For now we'll start with outpatient only
	EncounterClassAmbulatory EncounterClass = "AMBULATORY"
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
	AllergyIntoleranceReactionSeverityEnumMild     AllergyIntoleranceReactionSeverityEnum = "MILD"
	AllergyIntoleranceReactionSeverityEnumModerate AllergyIntoleranceReactionSeverityEnum = "MODERATE"
	AllergyIntoleranceReactionSeverityEnumSevere   AllergyIntoleranceReactionSeverityEnum = "SEVERE"
)

type ObservationStatus string

const (
	ObservationStatusFinal     ObservationStatus = "FINAL"
	ObservationStatusCancelled ObservationStatus = "CANCELLED"
)

type MedicationStatementStatusEnum string

const (
	MedicationStatementStatusEnumActive     MedicationStatementStatusEnum = "ACTIVE"
	MedicationStatementStatusEnumInActive   MedicationStatementStatusEnum = "INACTIVE"
	MedicationStatementStatusEnumUnknown    MedicationStatementStatusEnum = "UNKNOWN"
	MedicationStatementStatusEnumRecurrence MedicationStatementStatusEnum = "RECURRENCE"
	MedicationStatementStatusEnumRelapse    MedicationStatementStatusEnum = "RELAPSE"
	MedicationStatementStatusEnumRemission  MedicationStatementStatusEnum = "REMISSSION"
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
	ConditionStatusActive   ConditionStatus = "ACTIVE"
	ConditionStatusInactive ConditionStatus = "INACTIVE"
	ConditionStatusResolved ConditionStatus = "RESOLVED"
	ConditionStatusUnknown  ConditionStatus = "UNKNOWN"
)

// ConditionStatus represents status of a FHIR condition
type ConditionCategory string

const (
	ConditionCategoryProblemList ConditionCategory = "PROBLEM_LIST_ITEM"
	ConditionCategoryDiagnosis   ConditionCategory = "ENCOUNTER_DIAGNOSIS"
)

// TerminologySource represents various concept sources
type TerminologySource string

const (
	TerminologySourceICD10    TerminologySource = "ICD10"
	TerminologySourceCIEL     TerminologySource = "CIEL"
	TerminologySourceSNOMEDCT TerminologySource = "SNOMED_CT"
	TerminologySourceLOINC    TerminologySource = "LOINC"
)
