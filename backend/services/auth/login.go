package auth

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/auth/db"
	"encore.app/services/auth/jwt"
	"encore.app/services/auth/password"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// LoginParams defines the parameters required for user login.
type LoginParams struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required" encore:"sensitive"`
}

func (p LoginParams) Validate() error {
	return validation.ValidateStruct(p)
}

// LoginResponse defines the response structure for user login.
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Username     string `json:"username,omitempty"`
	PhoneNumber  string `json:"phoneNumber,omitempty"`
	OfficeCode   string `json:"officeCode,omitempty"`
	AgentCode    string `json:"agentCode,omitempty"`
}

// encore:api public path=/login method=POST
func (s *Service) Login(ctx context.Context, p LoginParams) (*LoginResponse, error) {
	row, err := s.query.GetUserByUsername(ctx, p.Username)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get user by username", "username", p.Username, "error", err)
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

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     row.Username,
		PhoneNumber:  row.PhoneNumber.String,
		OfficeCode:   row.OfficeCode.String,
		AgentCode:    row.AgentCode.String,
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
