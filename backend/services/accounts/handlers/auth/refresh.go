package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// RefreshTokensParams defines the parameters required for refreshing tokens.
type RefreshTokensParams struct {
	CookieHeader string `header:"Cookie" encore:"sensitive"`
}

func (s *AuthService) RefreshTokens(ctx context.Context, p RefreshTokensParams) (*LoginResponse, error) {
	claims, err := getRefreshTokenClaims(p.CookieHeader)
	if err != nil {
		return nil, err
	}

	savedToken, err := s.query.GetRefreshToken(ctx, claims.ID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			rlog.Warn("refresh token not found in database", "token_id", claims.ID)
			return nil, ErrInvalidRefreshToken
		}
		rlog.Error("failed to get refresh token from database", "error", err)
		return nil, api_errors.ErrInternalError
	}

	if savedToken.ExpiresAt.Time.Before(time.Now()) {
		rlog.Warn("refresh token has expired", "token_id", claims.ID, "expires_at", savedToken.ExpiresAt)
		return nil, ErrExpiredRefreshToken
	}

	user, err := s.query.GetUserById(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			rlog.Warn("user not found", "user_id", claims.UserID)
			return nil, ErrInvalidRefreshToken
		}
		rlog.Error("failed to get user by ID", "user_id", claims.UserID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	data := accessTokenDataFromUser(user, savedToken.AdminRefID)
	tokens, err := s.generateTokens(ctx, data)
	if err != nil {
		rlog.Error("failed to generate new tokens", "user_id", user.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	// Rotate the refresh token: consume the old jti (single-use)
	// If 0 rows deleted, the token was already used (replay) -> revoke entire session
	deleted, err := s.query.DeleteRefreshTokenChecked(ctx, claims.ID)
	if err != nil {
		rlog.Error("failed to delete old refresh token", "token_id", claims.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}
	if deleted == 0 {
		rlog.Warn("refresh token replay detected - revoking all sessions", "user_id", claims.UserID, "token_id", claims.ID)
		_ = s.query.DeleteRefreshTokensByUserId(ctx, claims.UserID)
		return nil, ErrInvalidRefreshToken
	}

	if err := s.query.UpdateLastLogin(ctx, user.ID); err != nil {
		rlog.Error("failed to update last login", "user_id", user.ID, "error", err)
	}

	return &LoginResponse{
		ID:                   user.ID,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Role:                 user.Role,
		OfficeID:             user.OfficeID,
		AccessToken:          tokens.AccessToken,
		AccessTokenExpiresAt: tokens.AccessTokenExpiresAt,
		Email:                user.Email,
		PhoneNumber:          ptrToStr(user.PhoneNumber),
		SetCookies:           authCookies(tokens.RefreshToken),
		IsAdminAsAgent:       savedToken.AdminRefID != nil,
	}, nil
}

// SessionUser is the read-only identity resolved from a refresh token.
type SessionUser struct {
	ID        int64       `json:"id"`
	Email     string      `json:"email"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Role      db.UserRole `json:"role"`
}

type ResolveSessionParams struct {
	CookieHeader string `header:"Cookie" encore:"sensitive"`
}

func (s *AuthService) ResolveSession(ctx context.Context, p ResolveSessionParams) (*SessionUser, error) {
	claims, err := getRefreshTokenClaims(p.CookieHeader)
	if err != nil {
		return nil, err
	}

	savedToken, err := s.query.GetRefreshToken(ctx, claims.ID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			// Token was revoked/rotated away -> no longer a valid session.
			return nil, ErrInvalidRefreshToken
		}
		rlog.Error("failed to get refresh token from database", "error", err)
		return nil, api_errors.ErrInternalError
	}

	if savedToken.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredRefreshToken
	}

	user, err := s.query.GetUserById(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidRefreshToken
		}
		rlog.Error("failed to get user by ID", "user_id", claims.UserID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &SessionUser{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}, nil
}

func getRefreshTokenClaims(cookieHeader string) (*jwt.RefreshTokenClaims, error) {
	req := &http.Request{Header: http.Header{"Cookie": {cookieHeader}}}
	cookie, err := req.Cookie(refreshCookieName)
	if err != nil || cookie.Value == "" {
		return nil, ErrInvalidRefreshToken
	}

	claims, err := jwt.ValidateRefreshToken(cookie.Value)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	return claims, nil
}
