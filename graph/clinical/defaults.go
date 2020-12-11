package clinical

import (
	"time"

	"github.com/rs/xid"
	"gitlab.slade360emr.com/go/base"
)

const (
	healthCloudIdentifiers        = "healthcloud.identifiers"
	healthCloudIdentifiersVersion = "0.0.1"
)

// DefaultCommunication returns default values for patient / person communication
// preferences
func DefaultCommunication() []*CommunicationInput {
	return []*CommunicationInput{
		{
			Language: &CodeableConceptInput{
				Coding: []*CodingInput{
					{
						System:       "urn:ietf:bcp:47",
						Code:         "en",
						Display:      "English",
						UserSelected: false,
					},
					{
						System:       "urn:ietf:bcp:47",
						Code:         "en-US",
						Display:      "English (United States)",
						UserSelected: false,
					},
				},
				Text: DefaultLanguage,
			},
			Preferred: true,
		},
	}
}

// DefaultIdentifier assigns a patient a code to function as their
// medical record number.
func DefaultIdentifier() *FHIRIdentifierInput {
	xid := xid.New().String()
	system := base.URI(healthCloudIdentifiers)
	version := healthCloudIdentifiersVersion
	userSelected := false
	return &FHIRIdentifierInput{
		Use:   IdentifierUseEnumOfficial,
		Value: xid,
		Type: FHIRCodeableConceptInput{
			Text: "MR",
			Coding: []*FHIRCodingInput{
				{
					System:       &system,
					Version:      &version,
					Code:         base.Code(xid),
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
func DefaultPeriodInput() *FHIRPeriodInput {
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * CenturyHours)
	return &FHIRPeriodInput{
		Start: base.DateTime(now.Format(timeFormatStr)),
		End:   base.DateTime(farFuture.Format(timeFormatStr)),
	}
}

// DefaultPeriod sets up a period input covering roughly a century from when it's run
func DefaultPeriod() *Period {
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * CenturyHours)
	return &Period{
		Start: &now,
		End:   &farFuture,
	}
}
