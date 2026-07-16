package auth

import (
	"context"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/jwt"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
)

// OTPRateLimiter abstracts OTP rate limiting to avoid an import cycle with the accounts package.
type OTPRateLimiter interface {
	CheckAndIncrementSend(ctx context.Context, phone string) error
	RecordFailedAttempt(ctx context.Context, phone string) (limitReached bool, err error)
	ClearAttempts(ctx context.Context, phone string)
	ValidateRateLimitErr() error
}

// AuthService holds the handler logic for all authentication endpoints.
type AuthService struct {
	query       db.Querier
	rateLimiter OTPRateLimiter
}

// NewAuthService creates a new AuthService.
func NewAuthService(query db.Querier) *AuthService {
	return &AuthService{query: query}
}

// WithOTPRateLimiter attaches a rate limiter for OTP endpoints and returns the same AuthService for chaining.
func (s *AuthService) WithOTPRateLimiter(rl OTPRateLimiter) *AuthService {
	s.rateLimiter = rl
	return s
}

// LoginResponse is the shared response type for all login/refresh endpoints.
// The refresh token is never serialized in the JSON body; it is delivered as an
// httpOnly cookie via the SetCookies field (Set-Cookie header).
type LoginResponse struct {
	ID                   int64       `json:"id"`
	Email                string      `json:"email,omitempty"`
	FirstName            string      `json:"firstName"`
	LastName             string      `json:"lastName"`
	Role                 db.UserRole `json:"role,omitempty"`
	AccessToken          string      `json:"accessToken"`
	AccessTokenExpiresAt int64       `json:"accessTokenExpiresAt"` //UNIX timestamp
	PhoneNumber          string      `json:"phoneNumber,omitempty"`
	OfficeID             *int64      `json:"officeId,omitempty"`

	// SetCookies carries the httpOnly refresh token and session hint cookies.
	SetCookies     []string `header:"Set-Cookie"`
	IsAdminAsAgent bool     `json:"isAdminAsAgent,omitempty" encore:"optional"`
}

// LogoutResponse clears the auth cookies via the Set-Cookie header.
type LogoutResponse struct {
	SetCookies []string `header:"Set-Cookie"`
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

	ErrCustomerPasswordResetNotAllowed = api_errors.NewErrorWithDetail(
		errs.PermissionDenied, "Password reset is not allowed for customers",
		api_errors.ErrorDetails{Code: api_errors.CodeCustomerPasswordResetNotAllowed},
	)
)

type generateTokenResponse struct {
	AccessToken          string
	AccessTokenExpiresAt int64
	RefreshToken         string
}

func (s *AuthService) generateTokens(ctx context.Context, data jwt.AccessTokenData) (generateTokenResponse, error) {
	accessToken, accessTokenExp, err := jwt.SignAccessToken(data)
	if err != nil {
		return generateTokenResponse{}, errs.Wrap(err, "failed to sign access token")
	}

	refreshToken, jti, exp, err := jwt.SignRefreshToken(data.UserID)
	if err != nil {
		return generateTokenResponse{}, errs.Wrap(err, "failed to sign refresh token")
	}

	err = s.query.SaveRefreshToken(ctx, db.SaveRefreshTokenParams{
		Jti:        jti,
		UserID:     data.UserID,
		AdminRefID: data.AdminRefID,
		ExpiresAt:  dbadapters.DBTime(exp),
	})
	if err != nil {
		return generateTokenResponse{}, errs.Wrap(err, "failed to save refresh token")
	}

	return generateTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExp.Unix(),
		RefreshToken:         refreshToken,
	}, nil
}

func ptrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
