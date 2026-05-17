package auth

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// LoginParams defines the parameters required for user login.
type LoginParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required" encore:"sensitive"`
}

func (p LoginParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *AuthService) Login(ctx context.Context, p LoginParams) (*LoginResponse, error) {
	row, err := s.query.GetUserByEmail(ctx, p.Email)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get user by email", "email", p.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	if !password.ComparePassword(row.PasswordHash, p.Password) {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, row, nil)
	if err != nil {
		rlog.Error("failed to generate tokens", "user_id", row.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &LoginResponse{
		ID:           row.ID,
		Role:         row.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        row.Email,
		PhoneNumber:  ptrToStr(row.PhoneNumber),
		OfficeID:     row.OfficeID,
	}, nil
}
