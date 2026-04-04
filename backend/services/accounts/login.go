package accounts

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
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

// LoginResponse defines the response structure for user login.
type LoginResponse struct {
	ID           int32       `json:"id"`
	Email        string      `json:"email,omitempty"`
	Role         db.UserRole `json:"role,omitempty"`
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	PhoneNumber  string      `json:"phoneNumber,omitempty"`
	OfficeID     *int32      `json:"officeId,omitempty"`
}

// encore:api public path=/login method=POST
func (s *Service) Login(ctx context.Context, p LoginParams) (*LoginResponse, error) {
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

	accessToken, refreshToken, err := s.generateTokens(ctx, row)
	if err != nil {
		rlog.Error("failed to generate tokens", "user_id", row.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	var phoneNumber string
	if row.PhoneNumber != nil {
		phoneNumber = *row.PhoneNumber
	}

	return &LoginResponse{
		ID:           row.ID,
		Role:         row.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        row.Email,
		PhoneNumber:  phoneNumber,
		OfficeID:     row.OfficeID,
	}, nil
}

func (s *Service) generateTokens(ctx context.Context, user db.User) (string, string, error) {
	accessToken, err := jwt.SignAccessToken(user)
	if err != nil {
		return "", "", errs.Wrap(err, "failed to sign access token")
	}

	refreshToken, jti, exp, err := jwt.SignRefreshToken(user.ID)
	if err != nil {
		return "", "", errs.Wrap(err, "failed to sign refresh token")
	}

	err = s.query.SaveRefreshToken(ctx, db.SaveRefreshTokenParams{
		Jti:       jti,
		UserID:    user.ID,
		ExpiresAt: db.DBTime(exp),
	})
	if err != nil {
		return "", "", errs.Wrap(err, "failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}
