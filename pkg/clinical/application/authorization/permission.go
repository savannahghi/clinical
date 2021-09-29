package authentication

import (
	"github.com/savannahghi/profileutils"
)

// ProblemSummaryView describes view permissions on viewing a patient problem summary resource
var ProblemSummaryView = profileutils.PermissionInput{
	Resource: "problem_summary_view",
	Action:   "view",
}

// VisitSummaryView describes view permissions on a user past visits
var VisitSummaryView = profileutils.PermissionInput{
	Resource: "visit_summary_view",
	Action:   "view",
}

// PatientTimelineWithCountView describes view permissions on a user patienttimeline
var PatientTimelineWithCountView = profileutils.PermissionInput{
	Resource: "patient_timeline_with_count_view",
	Action:   "view",
}

// EpisodeOfCareCreate describes write permissions on a patient episodeofcare
var EpisodeOfCareCreate = profileutils.PermissionInput{
	Resource: "episode_of_care_create",
	Action:   "create",
}

// EncountersList describes read permissions on a patient encounters
var EncountersList = profileutils.PermissionInput{
	Resource: "encounters_list",
	Action:   "view",
}

// StartEpisodeByOtpCreate describes write permissions on a patient StartEpisodeByOtp
var StartEpisodeByOtpCreate = profileutils.PermissionInput{
	Resource: "start_episode_by_otp_create",
	Action:   "create",
}

// UpgradeEpisodeUpdate describes update permissions on a patient episode resource
var UpgradeEpisodeUpdate = profileutils.PermissionInput{
	Resource: "upgrade_episode_update",
	Action:   "edit",
}

// StartEpisodeByBreakGlassCreate describes write permissions on a patient start episode by break glass
var StartEpisodeByBreakGlassCreate = profileutils.PermissionInput{
	Resource: "start_episode_by_break_glass_create",
	Action:   "create",
}

// OrganizationView describe view permissions on getting an organisation
var OrganizationView = profileutils.PermissionInput{
	Resource: "organisation_view",
	Action:   "view",
}

// OrganizationCreate describe write permissions on creating an organisation
var OrganizationCreate = profileutils.PermissionInput{
	Resource: "organisation_create",
	Action:   "create",
}

// OpenEpisodesView describe view permissions on an patient open episodes
var OpenEpisodesView = profileutils.PermissionInput{
	Resource: "open_episodes_view",
	Action:   "view",
}

// SearchEpisodesView describes permissions on an patient ability to search open episodes
var SearchEpisodesView = profileutils.PermissionInput{
	Resource: "search_episodes_view",
	Action:   "view",
}

// EncounterCreate describes write permissions on an encounter resource
var EncounterCreate = profileutils.PermissionInput{
	Resource: "encounter_create",
	Action:   "create",
}

// EncounterUpdate describes edit permissions on an encounter resource
var EncounterUpdate = profileutils.PermissionInput{
	Resource: "encounter_edit",
	Action:   "edit",
}

// PatientCreate describes write permissions on a patient resource
var PatientCreate = profileutils.PermissionInput{
	Resource: "patient_create",
	Action:   "create",
}

// PatientGet describes read permissions on patient resource
var PatientGet = profileutils.PermissionInput{
	Resource: "patient_get",
	Action:   "view",
}

// FHIRCompositionCreate describes write permissions on patient FHIR composition resource
var FHIRCompositionCreate = profileutils.PermissionInput{
	Resource: "fhir_composition_create",
	Action:   "create",
}

// FHIRCompositionEdit describes edit permissions on patient FHIR composition resource
var FHIRCompositionEdit = profileutils.PermissionInput{
	Resource: "fhir_composition_edit",
	Action:   "edit",
}

// FHIRCompositionDelete describes delete permissions on patient FHIR composition resource
var FHIRCompositionDelete = profileutils.PermissionInput{
	Resource: "fhir_composition_delete",
	Action:   "delete",
}

// FHIRConditionView describes view permissions on patient FHIR condition resource
var FHIRConditionView = profileutils.PermissionInput{
	Resource: "fhir_condition_view",
	Action:   "view",
}

// FHIRConditionCreate describes write permissions on patient FHIR condition resource
var FHIRConditionCreate = profileutils.PermissionInput{
	Resource: "fhir_condition_create",
	Action:   "create",
}

// FHIRConditionEdit describes edit permissions on patient FHIR condition resource
var FHIRConditionEdit = profileutils.PermissionInput{
	Resource: "fhir_condition_edit",
	Action:   "edit",
}

// FHIREncounterView describes view permissions on patient FHIR encounter resource
var FHIREncounterView = profileutils.PermissionInput{
	Resource: "fhir_encounter_view",
	Action:   "view",
}

// FHIREncounterCreate describes write permissions on patient FHIR condition resource
var FHIREncounterCreate = profileutils.PermissionInput{
	Resource: "fhir_enconter_create",
	Action:   "create",
}

// FHIRMedicationRequestView describes view permissions on a FHIR medication request resource
var FHIRMedicationRequestView = profileutils.PermissionInput{
	Resource: "fhir_medication_view",
	Action:   "view",
}

// FHIRMedicationRequestCreate describes write permissions on a FHIR medication request resource
var FHIRMedicationRequestCreate = profileutils.PermissionInput{
	Resource: "fhir_medication_create",
	Action:   "create",
}

// FHIRMedicationRequestEdit describes edit permissions on a FHIR medication request resource
var FHIRMedicationRequestEdit = profileutils.PermissionInput{
	Resource: "fhir_medication_edit",
	Action:   "edit",
}

// FHIRMedicationRequestDelete describes delete permissions on a FHIR medication request resource
var FHIRMedicationRequestDelete = profileutils.PermissionInput{
	Resource: "fhir_medication_delete",
	Action:   "delete",
}

// FHIRObservationView describes view permissions on a FHIR observation resource
var FHIRObservationView = profileutils.PermissionInput{
	Resource: "fhir_observation_view",
	Action:   "view",
}

// FHIRObservationCreate describes create permissions on a FHIR observation resource
var FHIRObservationCreate = profileutils.PermissionInput{
	Resource: "fhir_observation_create",
	Action:   "create",
}

// FHIRObservationDelete describes delete permissions on a FHIR observation resource
var FHIRObservationDelete = profileutils.PermissionInput{
	Resource: "fhir_observation_delete",
	Action:   "delete",
}

// PatientExtraInformationEdit describes edit permissions on a patient resource
var PatientExtraInformationEdit = profileutils.PermissionInput{
	Resource: "patient_update_extra_info_edit",
	Action:   "edit",
}

// FHIRPatientDelete describes delete permissions on a FHIR patient resource
var FHIRPatientDelete = profileutils.PermissionInput{
	Resource: "fhir_patient_delete",
	Action:   "delete",
}

// AllergyIntoleranceView describes view permissions on a FHIR allergy resource
var AllergyIntoleranceView = profileutils.PermissionInput{
	Resource: "allergy_intolerance_view",
	Action:   "view",
}
