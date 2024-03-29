package dto

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

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

// EncounterClassEnum ...
type EncounterClassEnum string

const (
	// Also referred to as outpatient - For now we'll start with outpatient only
	EncounterClassEnumAmbulatory EncounterClassEnum = "AMBULATORY"
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

// ConditionCategory represents status of a FHIR condition
type ConditionCategory string

const (
	ConditionCategoryProblemList ConditionCategory = "PROBLEM_LIST_ITEM"
	ConditionCategoryDiagnosis   ConditionCategory = "ENCOUNTER_DIAGNOSIS"
)

// TerminologySource represents various concept sources
type TerminologySource string

const (
	TerminologySourceICD10    TerminologySource = "ICD10"
	TerminologySourceICD10WHO TerminologySource = "ICD-10-WHO"
	TerminologySourceCIEL     TerminologySource = "CIEL"
	TerminologySourceSNOMEDCT TerminologySource = "SNOMED_CT"
	TerminologySourceLOINC    TerminologySource = "LOINC"
)

type OrganisationSource string

const (
	OrganisationSourceWHO  OrganisationSource = "WHO"
	OrganisationSourceCIEL OrganisationSource = "CIEL"
)

// LOINCCodes represents LOINC assessment codes
type LOINCCodes string

const (
	LOINCPlanOfCareCode     LOINCCodes = "18776-5"
	LOINCAssessmentPlanCode LOINCCodes = "51847-2"
)

// CompositionCategory enum represents category composition attribute
type CompositionCategory string

const (
	AssessmentAndPlan               CompositionCategory = "ASSESSMENT_PLAN"
	HistoryOfPresentingIllness      CompositionCategory = "HISTORY_OF_PRESENTING_ILLNESS"
	SocialHistory                   CompositionCategory = "SOCIAL_HISTORY"
	FamilyHistory                   CompositionCategory = "FAMILY_HISTORY"
	Examination                     CompositionCategory = "EXAMINATION"
	PlanOfCare                      CompositionCategory = "PLAN_OF_CARE"
	ProviderUnspecifiedProgressNote CompositionCategory = "PROGRESS_NOTE"
)

// Type enum represents type composition attribute
type CompositionType string

const (
	ProgressNote CompositionType = "PROGRESS_NOTE"
)

// CompositionStatus enum represents status composition attribute
type CompositionStatusEnum string

const (
	CompositionStatuEnumPreliminary               CompositionStatusEnum = "PRELIMINARY"
	CompositionStatuEnumFinal                     CompositionStatusEnum = "FINAL"
	CompositionStatuEnumAmended                   CompositionStatusEnum = "AMENDED"
	CompositionStatuEnumEnteredInErrorPreliminary CompositionStatusEnum = "ENTERED_IN_ERROR"
)

// ConsentStatusEnum a type enum tha represents a Consent Status field of consent resource
type ConsentStatusEnum string

const (
	ConsentStatusActive   = "active"
	ConsentStatusInactive = "inactive"
)

// IsValid ...
func (c ConsentStatusEnum) IsValid() bool {
	switch c {
	case ConsentStatusActive, ConsentStatusInactive:
		return true
	}

	return false
}

// String converts status to string
func (c ConsentStatusEnum) String() string {
	return string(c)
}

// MarshalGQL writes the consent status as a quoted string
func (c ConsentStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a consent status enum
func (c *ConsentStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*c = ConsentStatusEnum(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid ConsentStatus Enum", str)
	}

	return nil
}

// ConsentProvisionTypeEnum a type enum tha represents a Consent Provision field of consent resource
type ConsentProvisionTypeEnum string

const (
	ConsentProvisionTypeDeny   = "deny"
	ConsentProvisionTypePermit = "permit"
)

// IsValid ...
func (c ConsentProvisionTypeEnum) IsValid() bool {
	switch c {
	case ConsentProvisionTypeDeny, ConsentProvisionTypePermit:
		return true
	}

	return false
}

// String converts consent provision type to string
func (c ConsentProvisionTypeEnum) String() string {
	return string(c)
}

// MarshalGQL writes the consent provision type as a quoted string
func (c ConsentProvisionTypeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a consent provision type enum
func (c *ConsentProvisionTypeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*c = ConsentProvisionTypeEnum(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid ConsentProvisionTypeEnum Enum", str)
	}

	return nil
}

// ConsentDataMeaningEnum represents the meaning of consent data
type ConsentDataMeaningEnum string

const (
	ConsentDataMeaningInstance   = "instance"
	ConsentDataMeaningRelated    = "related"
	ConsentDataMeaningDependents = "dependents"
	ConsentDataMeaningAuthoredBy = "authoredby"
)

// IsValid checks if the consent data meaning is valid
func (c ConsentDataMeaningEnum) IsValid() bool {
	switch c {
	case ConsentDataMeaningInstance, ConsentDataMeaningRelated, ConsentDataMeaningDependents, ConsentDataMeaningAuthoredBy:
		return true
	}

	return false
}

// String converts the consent data meaning to string
func (c ConsentDataMeaningEnum) String() string {
	return string(c)
}

// MarshalGQL writes the consent data meaning as a quoted string
func (c ConsentDataMeaningEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a JSON and converts it to a consent data meaning enum
func (c *ConsentDataMeaningEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*c = ConsentDataMeaningEnum(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid ConsentDataMeaningEnum", str)
	}

	return nil
}

// QuestionnaireResponseStatusEnum a type enum tha represents a questionnaire response status field of questionnaire response
type QuestionnaireResponseStatusEnum string

const (
	QuestionnaireResponseStatusEnumCompleted = "completed"
)

// IsValid ...
func (c QuestionnaireResponseStatusEnum) IsValid() bool {
	return c == QuestionnaireResponseStatusEnumCompleted
}

// String converts questionnaire response status type to string
func (c QuestionnaireResponseStatusEnum) String() string {
	return string(c)
}

// MarshalGQL writes the questionnaire response status type as a quoted string
func (c QuestionnaireResponseStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a questionnaire response status type enum
func (c *QuestionnaireResponseStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*c = QuestionnaireResponseStatusEnum(strings.ReplaceAll(str, "_", "-"))
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid QuestionnaireResponseStatus Enum", str)
	}

	return nil
}

// QuantityComparatorEnum is a FHIR enum
type QuantityComparatorEnum string

const (
	QuantityComparatorEnumLessThan             QuantityComparatorEnum = "less_than"
	QuantityComparatorEnumLessThanOrEqualTo    QuantityComparatorEnum = "less_than_or_equal_to"
	QuantityComparatorEnumGreaterThanOrEqualTo QuantityComparatorEnum = "greater_than_or_equal_to"
	QuantityComparatorEnumGreaterThan          QuantityComparatorEnum = "greater_than"
)

// VIAOutcomeEnum a type enum that represents the results of a VIA test
// VIA (Visual Inspection with Acetic Acid)
type VIAOutcomeEnum string

const (
	VIAOutcomeNegative               = "negative"
	VIAOutcomePositive               = "positive"
	VIAOutcomePositiveInvasiveCancer = "suspicious_for_cancer"
)

// IsValid checks validity of a VIAOutcomeEnum enum
func (c VIAOutcomeEnum) IsValid() bool {
	switch c {
	case VIAOutcomeNegative, VIAOutcomePositive, VIAOutcomePositiveInvasiveCancer:
		return true
	}

	return false
}

// String converts VIA to string
func (c VIAOutcomeEnum) String() string {
	return string(c)
}

// MarshalGQL writes the VIA as a quoted string
func (c VIAOutcomeEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a VIA enum
func (c *VIAOutcomeEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be a string")
	}

	*c = VIAOutcomeEnum(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid VIAOutcomeEnum Enum", str)
	}

	return nil
}

// ObservationStatusEnum is a FHIR enum
type ObservationStatusEnum string

const (
	// ObservationStatusEnumRegistered ...
	ObservationStatusEnumRegistered ObservationStatusEnum = "registered"
	// ObservationStatusEnumPreliminary ...
	ObservationStatusEnumPreliminary ObservationStatusEnum = "preliminary"
	// ObservationStatusEnumFinal ...
	ObservationStatusEnumFinal ObservationStatusEnum = "final"
	// ObservationStatusEnumAmended ...
	ObservationStatusEnumAmended ObservationStatusEnum = "amended"
	// ObservationStatusEnumCorrected ...
	ObservationStatusEnumCorrected ObservationStatusEnum = "corrected"
	// ObservationStatusEnumCancelled ...
	ObservationStatusEnumCancelled ObservationStatusEnum = "cancelled"
	// ObservationStatusEnumEnteredInError ...
	ObservationStatusEnumEnteredInError ObservationStatusEnum = "entered_in_error"
	// ObservationStatusEnumUnknown ...
	ObservationStatusEnumUnknown ObservationStatusEnum = "unknown"
)

// AllObservationStatusEnum ...
var AllObservationStatusEnum = []ObservationStatusEnum{
	ObservationStatusEnumRegistered,
	ObservationStatusEnumPreliminary,
	ObservationStatusEnumFinal,
	ObservationStatusEnumAmended,
	ObservationStatusEnumCorrected,
	ObservationStatusEnumCancelled,
	ObservationStatusEnumEnteredInError,
	ObservationStatusEnumUnknown,
}

// IsValid ...
func (e ObservationStatusEnum) IsValid() bool {
	switch e {
	case ObservationStatusEnumRegistered, ObservationStatusEnumPreliminary, ObservationStatusEnumFinal, ObservationStatusEnumAmended, ObservationStatusEnumCorrected, ObservationStatusEnumCancelled, ObservationStatusEnumEnteredInError, ObservationStatusEnumUnknown:
		return true
	}

	return false
}

// String ...
func (e ObservationStatusEnum) String() string {
	if e == ObservationStatusEnumEnteredInError {
		return "entered-in-error"
	}

	return string(e)
}

// UnmarshalGQL ...
func (e *ObservationStatusEnum) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ObservationStatusEnum(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ObservationStatusEnum", str)
	}

	return nil
}

// MarshalGQL writes the observation status to the supplied writer as a quoted string
func (e ObservationStatusEnum) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// SegmentationCategory models advantage segmentation categories for clients
type SegmentationCategory string

const (
	SegmentationCategoryNoRisk           SegmentationCategory = "CERVICAL_CANCER_TIPS"
	SegmentationCategoryLowRisk          SegmentationCategory = "CERVICAL_CANCER_LOW_RISK"
	SegmentationCategoryHighRiskPositive SegmentationCategory = "CERVICAL_CANCER_POSITIVE"
	SegmentationCategoryHighRiskNegative SegmentationCategory = "CERVICAL_CANCER_HIGH_RISK"
)

// IsValid checks validity of a SegmentationCategory enum
func (c SegmentationCategory) IsValid() bool {
	switch c {
	case SegmentationCategoryNoRisk, SegmentationCategoryLowRisk, SegmentationCategoryHighRiskPositive, SegmentationCategoryHighRiskNegative:
		return true
	}

	return false
}

// String converts segmentation to string
func (c SegmentationCategory) String() string {
	return string(c)
}

// MarshalGQL writes the segmentation as a quoted string
func (c SegmentationCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}

// UnmarshalGQL reads a json and converts it to a segmentation enum
func (c *SegmentationCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be a string")
	}

	*c = SegmentationCategory(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid SegmentationCategory Enum", str)
	}

	return nil
}

// ReferralTypeEnum defines the type of referral
type ReferralTypeEnum string

const (
	// DIAGNOSTICS refers to diagnostics and testing
	DIAGNOSTICS ReferralTypeEnum = "DIAGNOSTICS"
	// SPECIALIST refers to a specialist consultation
	SPECIALIST ReferralTypeEnum = "SPECIALIST"
	// TREATMENT refers to treatment
	TREATMENT ReferralTypeEnum = "TREATMENT"
)
