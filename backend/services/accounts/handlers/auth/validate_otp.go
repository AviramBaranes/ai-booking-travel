package auth

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// ValidateCustomerLoginOTPParams defines the parameters required to validate a customer login OTP.
type ValidateCustomerLoginOTPParams struct {
	PhoneNumber string `json:"phoneNumber" validate:"required,israeli_phone"`
	OTP         string `json:"otp" validate:"required,len=6"`
}

func (p ValidateCustomerLoginOTPParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *AuthService) ValidateCustomerLoginOTP(ctx context.Context, params ValidateCustomerLoginOTPParams) (*LoginResponse, error) {
	user, err := s.query.GetUserByPhone(ctx, &params.PhoneNumber)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get user by phone number", "phone_number", params.PhoneNumber, "error", err)
		return nil, api_errors.ErrInternalError
	}

	if user.Otp == nil || *user.Otp != params.OTP {
		if rateLimitErr := s.checkValidateRateLimit(ctx, params.PhoneNumber, user.ID); rateLimitErr != nil {
			return nil, rateLimitErr
		}
		return nil, ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(ctx, jwt.AccessTokenData{
		UserID: user.ID,
		Role:   user.Role,
	})
	if err != nil {
		rlog.Error("failed to generate tokens in validate customer login OTP", "user_id", user.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	// Clear OTP after successful login.
	if err = s.query.SaveOTP(ctx, db.SaveOTPParams{
		ID:  user.ID,
		Otp: nil,
	}); err != nil {
		rlog.Error("failed to clear OTP after successful login", "user_id", user.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	if s.rateLimiter != nil {
		s.rateLimiter.ClearAttempts(ctx, params.PhoneNumber)
	}

	return &LoginResponse{
		ID:                   user.ID,
		Role:                 user.Role,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		AccessToken:          tokens.AccessToken,
		AccessTokenExpiresAt: tokens.AccessTokenExpiresAt,
		Email:                user.Email,
		PhoneNumber:          ptrToStr(user.PhoneNumber),
		SetCookies:           authCookies(tokens.RefreshToken),
	}, nil
}

// checkValidateRateLimit records a failed OTP attempt. If the lockout threshold is reached,
// it clears the OTP (forcing a new send) and returns the appropriate rate limit error.
func (s *AuthService) checkValidateRateLimit(ctx context.Context, phone string, userID int64) error {
	if s.rateLimiter == nil {
		return nil
	}
	limitReached, err := s.rateLimiter.RecordFailedAttempt(ctx, phone)
	if err != nil {
		return err
	}
	if limitReached {
		if clearErr := s.query.SaveOTP(ctx, db.SaveOTPParams{ID: userID, Otp: nil}); clearErr != nil {
			rlog.Error("failed to clear OTP on rate limit lockout", "user_id", userID, "error", clearErr)
		}
		return s.rateLimiter.ValidateRateLimitErr()
	}
	return nil
}
