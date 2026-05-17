package auth_handler

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
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
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user, nil)
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

	return &LoginResponse{
		ID:           user.ID,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
		PhoneNumber:  ptrToStr(user.PhoneNumber),
	}, nil
}
