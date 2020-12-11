package clinical

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
)

// ParseDate parses a formatted string and returns the time value it represents
// Try different layout due to Inconsistent date formats
// They include "2018-01-01", "2020-09-24T18:02:38.661033Z", "2018-Jan-01"
func ParseDate(date string) time.Time {
	// parses "2020-09-24T18:02:38.661033Z"
	timeValue, err := time.Parse(time.RFC3339, date)
	if err != nil {
		// parses when month is a number e.g "2018-01-01"
		timeValue, err := time.Parse(StringTimeParseMonthNameLayout, date)

		if err != nil {
			// parses when month is a name e.g "2018-Jan-01"
			timeValue, err := time.Parse(StringTimeParseMonthNumberLayout, date)

			if err != nil {
				log.Errorf("cannot parse date %v with error %v", date, err)
			}
			return timeValue
		}
		return timeValue
	}

	return timeValue
}

// ComposeOneHealthEpisodeOfCare is used to create an episode of care
func ComposeOneHealthEpisodeOfCare(
	validPhone string, fullAccess bool, organizationID, providerCode, patientID string,
) FHIREpisodeOfCare {
	accessLevel := ""
	if fullAccess {
		accessLevel = fullAccessLevel
	} else {
		accessLevel = partialAccessLevel
	}
	now := time.Now()
	farFuture := time.Now().Add(time.Hour * CenturyHours)
	orgIdentifier := &FHIRIdentifier{
		Use:   "official",
		Value: providerCode,
	}
	active := EpisodeOfCareStatusEnumActive
	orgRef := fmt.Sprintf("Organization/%s", organizationID)
	patientRef := fmt.Sprintf("Patient/%s", patientID)
	orgType := base.URI("Organization")
	patientType := base.URI("Patient")
	return FHIREpisodeOfCare{
		Status: &active,
		Period: &FHIRPeriod{
			Start: base.DateTime(now.Format(timeFormatStr)),
			End:   base.DateTime(farFuture.Format(timeFormatStr)),
		},
		ManagingOrganization: &FHIRReference{
			Reference:  &orgRef,
			Display:    providerCode,
			Type:       &orgType,
			Identifier: orgIdentifier,
		},
		Patient: &FHIRReference{
			Reference: &patientRef,
			Display:   validPhone,
			Type:      &patientType,
		},
		Type: []*FHIRCodeableConcept{
			{
				Text: accessLevel, // FULL_ACCESS or PROFILE_AND_RECENT_VISITS_ACCESS
			},
		},
	}
}
