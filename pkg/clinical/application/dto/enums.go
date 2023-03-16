package dto

type OrganizationIdentifierType string

const (
	SladeCode OrganizationIdentifierType = "SladeCode"
	MFLCode   OrganizationIdentifierType = "MFLCode"
	Other     OrganizationIdentifierType = "Other"
)

type EpisodeOfCareStatusEnum string

const (
	EpisodeOfCareStatusEnumPlanned   EpisodeOfCareStatusEnum = "planned"
	EpisodeOfCareStatusEnumActive    EpisodeOfCareStatusEnum = "active"
	EpisodeOfCareStatusEnumFinished  EpisodeOfCareStatusEnum = "finished"
	EpisodeOfCareStatusEnumCancelled EpisodeOfCareStatusEnum = "cancelled"
)
