package engagement

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/savannahghi/clinical/pkg/clinical/application/common"
	"github.com/savannahghi/clinical/pkg/clinical/application/extensions"
	"github.com/savannahghi/clinical/pkg/clinical/domain"
	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
)

const (
	timeFormatStr = "2006-01-02T15:04:05+03:00"
	// DefaultLanguage ...
	DefaultLanguage = "English"
	// EmailWelcomeSubject ...
	EmailWelcomeSubject = "Welcome to Be.Well"
)

// engagement ISC paths
const (
	// verifyOTPEndpoint ISC endpoint to verify OTP
	VerifyOTPEndpoint = "internal/verify_otp/"
	SendOtp           = "internal/send_otp/"

	// sendEmailEndpoint ISC endpoint to send email
	sendEmailEndpoint = "internal/send_email"

	// engagement ISC uploads paths
	uploadEndpoint    = "internal/upload/"
	getUploadEndpoint = "internal/upload/%s/"
)

// ServiceEngagement represents engagement usecases
type ServiceEngagement interface {
	VerifyOTP(
		ctx context.Context,
		msisdn string,
		otp string,
	) (bool, string, error)
	RequestOTP(
		ctx context.Context,
		msisdn string,
	) (string, error)
	PhotosToAttachments(
		ctx context.Context,
		photos []*domain.PhotoInput,
	) ([]*domain.FHIRAttachmentInput, error)
	SendPatientWelcomeEmail(ctx context.Context, emailaddress string) error
}

// ServiceEngagementImpl represents engagement usecases
type ServiceEngagementImpl struct {
	Engage extensions.ISCClientExtension
	Basext extensions.BaseExtension
}

// NewServiceEngagementImpl returns new instance of ServiceEngagementImpl
func NewServiceEngagementImpl(
	eng extensions.ISCClientExtension,
	ext extensions.BaseExtension,
) *ServiceEngagementImpl {
	return &ServiceEngagementImpl{
		Engage: eng,
		Basext: ext,
	}
}

// VerifyOTP sends an inter-service API call to the OTP service and checks if
// the supplied verification code is valid.
func (en *ServiceEngagementImpl) VerifyOTP(
	ctx context.Context,
	msisdn string,
	otp string,
) (bool, string, error) {
	en.checkPreconditions()

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

	resp, err := en.Engage.MakeRequest(
		ctx, http.MethodPost, VerifyOTPEndpoint, verifyPayload)
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
func (en *ServiceEngagementImpl) RequestOTP(
	ctx context.Context,
	msisdn string,
) (string, error) {
	en.checkPreconditions()

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
	resp, err := en.Engage.MakeRequest(
		ctx, http.MethodPost, SendOtp, requestPayload)
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

// PhotosToAttachments translates patient photos to FHIR attachments
func (en *ServiceEngagementImpl) PhotosToAttachments(
	ctx context.Context,
	photos []*domain.PhotoInput,
) ([]*domain.FHIRAttachmentInput, error) {
	if photos == nil {
		return []*domain.FHIRAttachmentInput{}, nil
	}

	output := []*domain.FHIRAttachmentInput{}
	for _, photo := range photos {
		uploadInput := profileutils.UploadInput{
			Title:       "Patient Photo",
			ContentType: photo.PhotoContentType.String(),
			Language:    enumutils.LanguageEn.String(),
			Base64data:  photo.PhotoBase64data,
			Filename:    photo.PhotoFilename,
		}

		resp, err := en.Engage.MakeRequest(
			ctx,
			http.MethodPost,
			uploadEndpoint,
			uploadInput,
		)
		if err != nil {
			return nil, fmt.Errorf("error sending upload: %w", err)
		}

		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading upload response: %w", err)
		}

		upload := profileutils.Upload{}
		err = json.Unmarshal(respData, &upload)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal upload response: %w", err)
		}

		data, err := base64.StdEncoding.DecodeString(upload.Base64data)
		if err != nil {
			return nil, errors.Wrap(err, "upload base64 decode error")
		}

		hash := scalarutils.Base64Binary(upload.Hash)
		size := len(data)
		url := scalarutils.URL(upload.URL)
		now := scalarutils.DateTime(time.Now().Format(timeFormatStr))
		contentType := scalarutils.Code(photo.PhotoContentType.String())
		language := scalarutils.Code(DefaultLanguage)
		photoData := scalarutils.Base64Binary(photo.PhotoBase64data)
		attachment := &domain.FHIRAttachmentInput{
			ContentType: &contentType,
			Language:    &language,
			Data:        &photoData,
			URL:         &url,
			Size:        &size,
			Hash:        &hash,
			Creation:    &now,
		}
		output = append(output, attachment)
	}
	return output, nil
}

// SendPatientWelcomeEmail will send a welcome email to the practitioner
func (en ServiceEngagementImpl) SendPatientWelcomeEmail(ctx context.Context, emailaddress string) error {
	en.checkPreconditions()

	text := common.GeneratePatientWelcomeEmailTemplate()
	if !govalidator.IsEmail(emailaddress) {
		return nil
	}
	body := map[string]interface{}{
		"to":      []string{emailaddress},
		"text":    text,
		"subject": EmailWelcomeSubject,
	}

	resp, err := en.Engage.MakeRequest(ctx, http.MethodPost, sendEmailEndpoint, body)
	if err != nil {
		return fmt.Errorf("unable to send welcome email: %w", err)
	}
	if serverutils.IsDebug() {
		b, _ := httputil.DumpResponse(resp, true)
		log.Println(string(b))
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got error status %s from email service", resp.Status)
	}

	return nil
}

func (en ServiceEngagementImpl) checkPreconditions() {
	if en.Engage == nil {
		log.Panicf("engagement ISC call is nil")
	}
}
