package common

import (
	"time"

	"github.com/rs/xid"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/scalarutils"
)

// constants and defaults
const (

	// CenturyHours is the number of hours in a (fictional) century of leap years
	CenturyHours                  = 878400
	timeFormatStr                 = "2006-01-02T15:04:05+03:00"
	healthCloudIdentifiers        = "healthcloud.identifiers"
	healthCloudIdentifiersVersion = "0.0.1"
)

// DefaultIdentifier assigns a patient a code to function as their
// medical record number.
func DefaultIdentifier() *domain.FHIRIdentifierInput {
	xid := xid.New().String()
	system := scalarutils.URI(healthCloudIdentifiers)
	version := healthCloudIdentifiersVersion
	userSelected := false
	return &domain.FHIRIdentifierInput{
		Use:   domain.IdentifierUseEnumOfficial,
		Value: xid,
		Type: domain.FHIRCodeableConceptInput{
			Text: "MR",
			Coding: []*domain.FHIRCodingInput{
				{
					System:       &system,
					Version:      &version,
					Code:         scalarutils.Code(xid),
					Display:      xid,
					UserSelected: &userSelected,
				},
			},
		},
		System: &system,
		Period: DefaultPeriodInput(),
	}
}

// DefaultPeriodInput sets up a period input covering roughly a century from when it's run
func DefaultPeriodInput() *domain.FHIRPeriodInput {
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * CenturyHours)
	return &domain.FHIRPeriodInput{
		Start: scalarutils.DateTime(now.Format(timeFormatStr)),
		End:   scalarutils.DateTime(farFuture.Format(timeFormatStr)),
	}
}

// DefaultPeriod sets up a period input covering roughly a century from when it's run
func DefaultPeriod() *domain.FHIRPeriod {
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * CenturyHours)
	return &domain.FHIRPeriod{
		Start: scalarutils.DateTime(now.Format(timeFormatStr)),
		End:   scalarutils.DateTime(farFuture.Format(timeFormatStr)),
	}
}