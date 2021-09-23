package domain

import (
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/scalarutils"
)

// FHIRServiceRequest definition: a record of a request for service such as diagnostic investigations, treatments, or operations to be performed.
type FHIRServiceRequest struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// A human-readable narrative that contains a summary of the resource and can be used to represent the content of the resource to a human. The narrative need not encode all the structured data, but is required to contain sufficient detail to make it "clinically safe" for a human to just read the narrative. Resource definitions may define what content should be represented in the narrative to ensure clinical safety.
	Text *FHIRNarrative `json:"text,omitempty"`

	// Identifiers assigned to this order instance by the orderer and/or the receiver and/or order fulfiller.
	Identifier []*FHIRIdentifier `json:"identifier,omitempty"`

	// The URL pointing to a FHIR-defined protocol, guideline, orderset or other definition that is adhered to in whole or in part by this ServiceRequest.
	InstantiatesCanonical *scalarutils.Canonical `json:"instantiatesCanonical,omitempty"`

	// The URL pointing to an externally maintained protocol, guideline, orderset or other definition that is adhered to in whole or in part by this ServiceRequest.
	InstantiatesURI *scalarutils.Instant `json:"instantiatesURI,omitempty"`

	// Plan/proposal/order fulfilled by this request.
	BasedOn []*FHIRReference `json:"basedOn,omitempty"`

	// The request takes the place of the referenced completed or terminated request(s).
	Replaces []*FHIRReference `json:"replaces,omitempty"`

	// A shared identifier common to all service requests that were authorized more or less simultaneously by a single author, representing the composite or group identifier.
	Requisition *FHIRIdentifier `json:"requisition,omitempty"`

	// The status of the order.
	Status *scalarutils.Code `json:"status,omitempty"`

	// Whether the request is a proposal, plan, an original order or a reflex order.
	Intent *scalarutils.Code `json:"intent,omitempty"`

	// A code that classifies the service for searching, sorting and display purposes (e.g. "Surgical Procedure").
	Category []*FHIRCodeableConcept `json:"category,omitempty"`

	// Indicates how quickly the ServiceRequest should be addressed with respect to other requests.
	Priority *scalarutils.Code `json:"priority,omitempty"`

	// Set this to true if the record is saying that the service/procedure should NOT be performed.
	DoNotPerform *bool `json:"doNotPerform,omitempty"`

	// A code that identifies a particular service (i.e., procedure, diagnostic investigation, or panel of investigations) that have been requested.
	Code *FHIRCodeableConcept `json:"code,omitempty"`

	// Additional details and instructions about the how the services are to be delivered.   For example, and order for a urinary catheter may have an order detail for an external or indwelling catheter, or an order for a bandage may require additional instructions specifying how the bandage should be applied.
	OrderDetail []*FHIRCodeableConcept `json:"orderDetail,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityQuantity *FHIRQuantity `json:"quantityQuantity,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityRatio *FHIRRatio `json:"quantityRatio,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityRange *FHIRRange `json:"quantityRange,omitempty"`

	// On whom or what the service is to be performed. This is usually a human patient, but can also be requested on animals, groups of humans or animals, devices such as dialysis machines, or even locations (typically for environmental scans).
	Subject *FHIRReference `json:"subject,omitempty"`

	// An encounter that provides additional information about the healthcare context in which this request is made.
	Encounter *FHIRReference `json:"encounter,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrenceDateTime *scalarutils.Date `json:"occurrenceDateTime,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrencePeriod *FHIRPeriod `json:"occurrencePeriod,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrenceTiming *FHIRTiming `json:"occurrenceTiming,omitempty"`

	// If a CodeableConcept is present, it indicates the pre-condition for performing the service.  For example "pain", "on flare-up", etc.
	AsNeededBoolean *bool `json:"asNeededBoolean,omitempty"`

	// If a CodeableConcept is present, it indicates the pre-condition for performing the service.  For example "pain", "on flare-up", etc.
	AsNeededCodeableConcept *scalarutils.Code `json:"asNeededCodeableConcept,omitempty"`

	// When the request transitioned to being actionable.
	AuthoredOn *scalarutils.DateTime `json:"authoredOn,omitempty"`

	// The individual who initiated the request and has responsibility for its activation.
	Requester *FHIRReference `json:"requester,omitempty"`

	// Desired type of performer for doing the requested service.
	PerformerType *FHIRCodeableConcept `json:"performerType,omitempty"`

	// The desired performer for doing the requested service.  For example, the surgeon, dermatopathologist, endoscopist, etc.
	Performer []*FHIRReference `json:"performer,omitempty"`

	// The preferred location(s) where the procedure should actually happen in coded or free text form. E.g. at home or nursing day care center.
	LocationCode *scalarutils.Code `json:"locationCode,omitempty"`

	// A reference to the the preferred location(s) where the procedure should actually happen. E.g. at home or nursing day care center.
	LocationReference []*FHIRReference `json:"locationReference,omitempty"`

	// An explanation or justification for why this service is being requested in coded or textual form.   This is often for billing purposes.  May relate to the resources referred to in `supportingInfo`.
	ReasonCode *scalarutils.Code `json:"reasonCode,omitempty"`

	// Indicates another resource that provides a justification for why this service is being requested.   May relate to the resources referred to in `supportingInfo`.
	ReasonReference []*FHIRReference `json:"reasonReference,omitempty"`

	// Insurance plans, coverage extensions, pre-authorizations and/or pre-determinations that may be needed for delivering the requested service.
	Insurance []*FHIRReference `json:"insurance,omitempty"`

	// Additional clinical information about the patient or specimen that may influence the services or their interpretations.     This information includes diagnosis, clinical findings and other observations.  In laboratory ordering these are typically referred to as "ask at order entry questions (AOEs)".  This includes observations explicitly requested by the producer (filler) to provide context or supporting information needed to complete the order. For example,  reporting the amount of inspired oxygen for blood gas measurements.
	SupportingInfo []*FHIRReference `json:"supportingInfo,omitempty"`

	// One or more specimens that the laboratory procedure will use.
	Specimen []*FHIRReference `json:"specimen,omitempty"`

	// Anatomic location where the procedure should be performed. This is the target site.
	BodySite []*FHIRCodeableConcept `json:"bodySite,omitempty"`

	// Any other notes and comments made about the service request. For example, internal billing notes.
	Note []*FHIRAnnotation `json:"note,omitempty"`

	// Instructions in terms that are understood by the patient or consumer.
	PatientInstruction *string `json:"patientInstruction,omitempty"`

	// Key events in the history of the request.
	RelevantHistory []*FHIRReference `json:"relevantHistory,omitempty"`
}

// FHIRServiceRequestInput is the input type for ServiceRequest
type FHIRServiceRequestInput struct {
	// The logical id of the resource, as used in the URL for the resource. Once assigned, this value never changes.
	ID *string `json:"id,omitempty"`

	// Identifiers assigned to this order instance by the orderer and/or the receiver and/or order fulfiller.
	Identifier []*FHIRIdentifierInput `json:"identifier,omitempty"`

	// The URL pointing to a FHIR-defined protocol, guideline, orderset or other definition that is adhered to in whole or in part by this ServiceRequest.
	InstantiatesCanonical *scalarutils.Canonical `json:"instantiatesCanonical,omitempty"`

	// The URL pointing to an externally maintained protocol, guideline, orderset or other definition that is adhered to in whole or in part by this ServiceRequest.
	InstantiatesURI *scalarutils.Instant `json:"instantiatesURI,omitempty"`

	// Plan/proposal/order fulfilled by this request.
	BasedOn []*FHIRReferenceInput `json:"basedOn,omitempty"`

	// The request takes the place of the referenced completed or terminated request(s).
	Replaces []*FHIRReferenceInput `json:"replaces,omitempty"`

	// A shared identifier common to all service requests that were authorized more or less simultaneously by a single author, representing the composite or group identifier.
	Requisition *FHIRIdentifierInput `json:"requisition,omitempty"`

	// The status of the order.
	Status *scalarutils.Code `json:"status,omitempty"`

	// Whether the request is a proposal, plan, an original order or a reflex order.
	Intent *scalarutils.Code `json:"intent,omitempty"`

	// A code that classifies the service for searching, sorting and display purposes (e.g. "Surgical Procedure").
	Category []*FHIRCodeableConceptInput `json:"category,omitempty"`

	// Indicates how quickly the ServiceRequest should be addressed with respect to other requests.
	Priority *scalarutils.Code `json:"priority,omitempty"`

	// Set this to true if the record is saying that the service/procedure should NOT be performed.
	DoNotPerform *bool `json:"doNotPerform,omitempty"`

	// A code that identifies a particular service (i.e., procedure, diagnostic investigation, or panel of investigations) that have been requested.
	Code *FHIRCodeableConceptInput `json:"code,omitempty"`

	// Additional details and instructions about the how the services are to be delivered.   For example, and order for a urinary catheter may have an order detail for an external or indwelling catheter, or an order for a bandage may require additional instructions specifying how the bandage should be applied.
	OrderDetail []*FHIRCodeableConceptInput `json:"orderDetail,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityQuantity *FHIRQuantityInput `json:"quantityQuantity,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityRatio *FHIRRatioInput `json:"quantityRatio,omitempty"`

	// An amount of service being requested which can be a quantity ( for example $1,500 home modification), a ratio ( for example, 20 half day visits per month), or a range (2.0 to 1.8 Gy per fraction).
	QuantityRange *FHIRRangeInput `json:"quantityRange,omitempty"`

	// On whom or what the service is to be performed. This is usually a human patient, but can also be requested on animals, groups of humans or animals, devices such as dialysis machines, or even locations (typically for environmental scans).
	Subject *FHIRReferenceInput `json:"subject,omitempty"`

	// An encounter that provides additional information about the healthcare context in which this request is made.
	Encounter *FHIRReferenceInput `json:"encounter,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrenceDateTime *scalarutils.Date `json:"occurrenceDateTime,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrencePeriod *FHIRPeriodInput `json:"occurrencePeriod,omitempty"`

	// The date/time at which the requested service should occur.
	OccurrenceTiming *FHIRTimingInput `json:"occurrenceTiming,omitempty"`

	// If a CodeableConcept is present, it indicates the pre-condition for performing the service.  For example "pain", "on flare-up", etc.
	AsNeededBoolean *bool `json:"asNeededBoolean,omitempty"`

	// If a CodeableConcept is present, it indicates the pre-condition for performing the service.  For example "pain", "on flare-up", etc.
	AsNeededCodeableConcept *scalarutils.Code `json:"asNeededCodeableConcept,omitempty"`

	// When the request transitioned to being actionable.
	AuthoredOn *scalarutils.DateTime `json:"authoredOn,omitempty"`

	// The individual who initiated the request and has responsibility for its activation.
	Requester *FHIRReferenceInput `json:"requester,omitempty"`

	// Desired type of performer for doing the requested service.
	PerformerType *FHIRCodeableConceptInput `json:"performerType,omitempty"`

	// The desired performer for doing the requested service.  For example, the surgeon, dermatopathologist, endoscopist, etc.
	Performer []*FHIRReferenceInput `json:"performer,omitempty"`

	// The preferred location(s) where the procedure should actually happen in coded or free text form. E.g. at home or nursing day care center.
	LocationCode *scalarutils.Code `json:"locationCode,omitempty"`

	// A reference to the the preferred location(s) where the procedure should actually happen. E.g. at home or nursing day care center.
	LocationReference []*FHIRReferenceInput `json:"locationReference,omitempty"`

	// An explanation or justification for why this service is being requested in coded or textual form.   This is often for billing purposes.  May relate to the resources referred to in `supportingInfo`.
	ReasonCode *scalarutils.Code `json:"reasonCode,omitempty"`

	// Indicates another resource that provides a justification for why this service is being requested.   May relate to the resources referred to in `supportingInfo`.
	ReasonReference []*FHIRReferenceInput `json:"reasonReference,omitempty"`

	// Insurance plans, coverage extensions, pre-authorizations and/or pre-determinations that may be needed for delivering the requested service.
	Insurance []*FHIRReferenceInput `json:"insurance,omitempty"`

	// Additional clinical information about the patient or specimen that may influence the services or their interpretations.     This information includes diagnosis, clinical findings and other observations.  In laboratory ordering these are typically referred to as "ask at order entry questions (AOEs)".  This includes observations explicitly requested by the producer (filler) to provide context or supporting information needed to complete the order. For example,  reporting the amount of inspired oxygen for blood gas measurements.
	SupportingInfo []*FHIRReferenceInput `json:"supportingInfo,omitempty"`

	// One or more specimens that the laboratory procedure will use.
	Specimen []*FHIRReferenceInput `json:"specimen,omitempty"`

	// Anatomic location where the procedure should be performed. This is the target site.
	BodySite []*FHIRCodeableConceptInput `json:"bodySite,omitempty"`

	// Any other notes and comments made about the service request. For example, internal billing notes.
	Note []*FHIRAnnotationInput `json:"note,omitempty"`

	// Instructions in terms that are understood by the patient or consumer.
	PatientInstruction *string `json:"patientInstruction,omitempty"`

	// Key events in the history of the request.
	RelevantHistory []*FHIRReferenceInput `json:"relevantHistory,omitempty"`
}

// FHIRServiceRequestRelayConnection is a Relay connection for ServiceRequest
type FHIRServiceRequestRelayConnection struct {
	Edges []*FHIRServiceRequestRelayEdge `json:"edges,omitempty"`

	PageInfo *firebasetools.PageInfo `json:"pageInfo,omitempty"`
}

// FHIRServiceRequestRelayEdge is a Relay edge for ServiceRequest
type FHIRServiceRequestRelayEdge struct {
	Cursor *string `json:"cursor,omitempty"`

	Node *FHIRServiceRequest `json:"node,omitempty"`
}

// FHIRServiceRequestRelayPayload is used to return single instances of ServiceRequest
type FHIRServiceRequestRelayPayload struct {
	Resource *FHIRServiceRequest `json:"resource,omitempty"`
}
