package auth

import (
	"context"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/jwt"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
)

// AuthService holds the handler logic for all authentication endpoints.
type AuthService struct {
	query db.Querier
}

// NewAuthService creates a new AuthService.
func NewAuthService(query db.Querier) *AuthService {
	return &AuthService{query: query}
}

// LoginResponse is the shared response type for all login/refresh endpoints.
type LoginResponse struct {
	ID           int64       `json:"id"`
	Email        string      `json:"email,omitempty"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Role         db.UserRole `json:"role,omitempty"`
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	PhoneNumber  string      `json:"phoneNumber,omitempty"`
	OfficeID     *int64      `json:"officeId,omitempty"`
}

var (
	ErrInvalidCredentials = api_errors.NewErrorWithDetail(
		errs.Unauthenticated, "Invalid credentials",
		api_errors.ErrorDetails{Code: api_errors.CodeInvalidCredentials},
	)

	ErrInvalidRefreshToken = api_errors.NewErrorWithDetail(
		errs.Unauthenticated, "Invalid refresh token",
		api_errors.ErrorDetails{Code: api_errors.CodeInvalidRefreshToken},
	)

	ErrExpiredRefreshToken = api_errors.NewErrorWithDetail(
		errs.Unauthenticated, "Expired refresh token",
		api_errors.ErrorDetails{Code: api_errors.CodeExpiredRefreshToken},
	)
)

func (s *AuthService) generateTokens(ctx context.Context, data jwt.AccessTokenData) (string, string, error) {
	accessToken, err := jwt.SignAccessToken(data)
	if err != nil {
		return "", "", errs.Wrap(err, "failed to sign access token")
	}

	refreshToken, jti, exp, err := jwt.SignRefreshToken(data.UserID)
	if err != nil {
		return "", "", errs.Wrap(err, "failed to sign refresh token")
	}

	err = s.query.SaveRefreshToken(ctx, db.SaveRefreshTokenParams{
		Jti:        jti,
		UserID:     data.UserID,
		AdminRefID: data.AdminRefID,
		ExpiresAt:  dbadapters.DBTime(exp),
	})
	if err != nil {
		return "", "", errs.Wrap(err, "failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}

func ptrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
