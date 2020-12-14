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
func DefaultCommunication() []*FHIRPatientCommunicationInput {
	userSelected := false
	preferred := true
	system := base.URI("urn:ietf:bcp:47")
	return []*FHIRPatientCommunicationInput{
		{
			Language: &FHIRCodeableConceptInput{
				Coding: []*FHIRCodingInput{
					{
						System:       &system,
						Code:         "en",
						Display:      "English",
						UserSelected: &userSelected,
					},
					{
						System:       &system,
						Code:         "en-US",
						Display:      "English (United States)",
						UserSelected: &userSelected,
					},
				},
				Text: DefaultLanguage,
			},
			Preferred: &preferred,
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
func DefaultPeriod() *FHIRPeriod {
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * CenturyHours)
	return &FHIRPeriod{
		Start: base.DateTime(now.Format(timeFormatStr)),
		End:   base.DateTime(farFuture.Format(timeFormatStr)),
	}
}
