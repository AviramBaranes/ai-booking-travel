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

// GetCustomerTokenParams defines the parameters required to validate a customer login OTP.
type GetCustomerTokenParams struct {
	UserID int64
}

func (p GetCustomerTokenParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *AuthService) GetCustomerToken(ctx context.Context, p GetCustomerTokenParams) (*LoginResponse, error) {
	user, err := s.query.GetCustomerByID(ctx, p.UserID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get user by ID", "user_id", p.UserID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	tokens, err := s.generateTokens(ctx, jwt.AccessTokenData{
		UserID: user.ID,
		Role:   user.Role,
	})
	if err != nil {
		rlog.Error("failed to generate tokens in GetCustomerToken", "user_id", user.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	// respond without refresh token for security reasons, as this endpoint is used to generate a new access token for a customer without requiring their password or OTP.
	return &LoginResponse{
		ID:                   user.ID,
		Role:                 user.Role,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		AccessToken:          tokens.AccessToken,
		AccessTokenExpiresAt: tokens.AccessTokenExpiresAt,
		Email:                user.Email,
		PhoneNumber:          ptrToStr(user.PhoneNumber),
	}, nil
}
