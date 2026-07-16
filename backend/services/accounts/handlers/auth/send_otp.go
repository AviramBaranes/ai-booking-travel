package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"encore.app/internal/api_errors"
	"encore.app/internal/lang"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

// SendCustomerLoginOTPParams defines the parameters required to send a login OTP to a customer.
type SendCustomerLoginOTPParams struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,israeli_phone"`
}

func (p SendCustomerLoginOTPParams) Validate() error {
	return validation.ValidateStruct(p)
}

// CustomerLoginOTPRequestedEvent is the event published when a customer login OTP is requested.
type CustomerLoginOTPRequestedEvent struct {
	PhoneNumber string
	OTP         string
	LangCode    string
}

// CustomerLoginOTPRequestedTopic is the Pub/Sub topic for customer login OTP events.
var CustomerLoginOTPRequestedTopic = pubsub.NewTopic[*CustomerLoginOTPRequestedEvent](
	"customer-login-otp-requested",
	pubsub.TopicConfig{
		DeliveryGuarantee: pubsub.AtLeastOnce,
	},
)

func (s *AuthService) SendCustomerLoginOTP(ctx context.Context, params SendCustomerLoginOTPParams) error {
	if s.rateLimiter != nil {
		if err := s.rateLimiter.CheckAndIncrementSend(ctx, params.PhoneNumber); err != nil {
			return err
		}
	}

	user, err := s.query.GetUserByPhone(ctx, &params.PhoneNumber)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return ErrInvalidCredentials
		}
		rlog.Error("failed to get user by phone number", "phone_number", params.PhoneNumber, "error", err)
		return api_errors.ErrInternalError
	}

	if user.Role != db.UserRoleCustomer {
		return ErrInvalidCredentials
	}

	otp, err := generateOTP(6)
	if err != nil {
		rlog.Error("generating otp failed", "error", err)
		return api_errors.ErrInternalError
	}

	if err = s.query.SaveOTP(ctx, db.SaveOTPParams{
		ID:  user.ID,
		Otp: &otp,
	}); err != nil {
		rlog.Error("failed to save otp", "user_id", user.ID, "error", err)
		return api_errors.ErrInternalError
	}

	langCode := lang.FromContext(ctx, "he")

	if _, err := CustomerLoginOTPRequestedTopic.Publish(ctx, &CustomerLoginOTPRequestedEvent{
		PhoneNumber: params.PhoneNumber,
		OTP:         otp,
		LangCode:    langCode,
	}); err != nil {
		rlog.Error("failed to publish customer login OTP requested event", "phone_number", params.PhoneNumber, "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}

func generateOTP(length int) (string, error) {
	const charset = "0123456789"
	otp := make([]byte, length)
	randMax := int64(len(charset))
	for i := range otp {
		index, err := rand.Int(rand.Reader, big.NewInt(randMax))
		if err != nil {
			return "", err
		}
		otp[i] = charset[index.Int64()]
	}
	return string(otp), nil
}
