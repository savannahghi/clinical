package clinical

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/scalarutils"
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
	orgType := scalarutils.URI("Organization")
	patientType := scalarutils.URI("Patient")
	return FHIREpisodeOfCare{
		Status: &active,
		Period: &FHIRPeriod{
			Start: scalarutils.DateTime(now.Format(timeFormatStr)),
			End:   scalarutils.DateTime(farFuture.Format(timeFormatStr)),
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

// VerifyOTP sends an inter-service API call to the OTP service and checks if
// the supplied verification code is valid.
func VerifyOTP(
	msisdn string,
	otp string,
	engagementClient *base.InterServiceClient,
) (bool, string, error) {
	if engagementClient == nil {
		return false, "", fmt.Errorf("nil engagement client")
	}

	normalized, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, "", fmt.Errorf("invalid phone format: %w", err)
	}
	type VerifyOTP struct {
		Msisdn           string `json:"msisdn"`
		VerificationCode string `json:"verificationCode"`
	}

	verifyPayload := VerifyOTP{
		Msisdn:           msisdn,
		VerificationCode: otp,
	}

	resp, err := engagementClient.MakeRequest(
		http.MethodPost, verifyOTPEndpoint, verifyPayload)
	if err != nil {
		return false, "", fmt.Errorf(
			"can't complete OTP verification request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf(
			"OTP verification call got non OK status: %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("can't read OTP response data: %w", err)
	}

	type otpResponse struct {
		IsVerified bool `json:"IsVerified"`
	}

	var r otpResponse
	err = json.Unmarshal(data, &r)
	if err != nil {
		return false, "", fmt.Errorf(
			"can't unmarshal OTP response data from JSON: %w", err)
	}

	return r.IsVerified, *normalized, nil
}

// RequestOTP sends an inter-service API call to the OTP service to request
// a new OTP.
func RequestOTP(
	msisdn string,
	engagementClient *base.InterServiceClient,
) (string, error) {
	if engagementClient == nil {
		return "", fmt.Errorf("nil engagement client")
	}

	normalized, err := converterandformatter.NormalizeMSISDN(msisdn)
	if err != nil {
		return "", fmt.Errorf("invalid phone format: %w", err)
	}

	type Msisdn struct {
		Msisdn string `json:"msisdn"`
	}
	requestPayload := Msisdn{
		Msisdn: *normalized,
	}
	resp, err := engagementClient.MakeRequest(
		http.MethodPost, sendOTPEndpoint, requestPayload)
	if err != nil {
		return "", fmt.Errorf("can't complete OTP request: %w", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't read OTP response data: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("OTP API response data: \n%s\n", string(data))
		return "", fmt.Errorf(
			"OTP request got non OK status: %s", resp.Status)
	}

	var r string

	err = json.Unmarshal(data, &r)
	if err != nil {
		return "", fmt.Errorf(
			"can't unmarshal OTP response data from JSON: %w", err)
	}

	return r, nil
}
