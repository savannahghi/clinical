package mock

import (
	"context"

	"github.com/savannahghi/clinical/pkg/clinical/domain"
)

// FakeServiceEngagement ...
type FakeServiceEngagement struct {
	VerifyOTPFn               func(ctx context.Context, msisdn string, otp string) (bool, string, error)
	RequestOTPFn              func(ctx context.Context, msisdn string) (string, error)
	PhotosToAttachmentsFn     func(ctx context.Context, photos []*domain.PhotoInput) ([]*domain.FHIRAttachmentInput, error)
	SendPatientWelcomeEmailFn func(ctx context.Context, emailaddress string) error
}

// VerifyOTP ...
func (e *FakeServiceEngagement) VerifyOTP(ctx context.Context, msisdn string, otp string) (bool, string, error) {
	return e.VerifyOTPFn(ctx, msisdn, otp)
}

// RequestOTP ...
func (e *FakeServiceEngagement) RequestOTP(ctx context.Context, msisdn string) (string, error) {
	return e.RequestOTPFn(ctx, msisdn)
}

// PhotosToAttachments ...
func (e *FakeServiceEngagement) PhotosToAttachments(ctx context.Context, photos []*domain.PhotoInput) ([]*domain.FHIRAttachmentInput, error) {
	return e.PhotosToAttachmentsFn(ctx, photos)
}

// PhotosToAttachments ...
func (e *FakeServiceEngagement) SendPatientWelcomeEmail(ctx context.Context, emailaddress string) error {
	return e.SendPatientWelcomeEmailFn(ctx, emailaddress)
}
