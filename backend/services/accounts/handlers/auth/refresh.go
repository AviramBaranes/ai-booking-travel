package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// RefreshTokensParams defines the parameters required for refreshing tokens.
type RefreshTokensParams struct {
	RefreshToken string `header:"Authorization"`
}

func (s *AuthService) RefreshTokens(ctx context.Context, p RefreshTokensParams) (*LoginResponse, error) {
	tokenString := strings.TrimPrefix(p.RefreshToken, "Bearer ")
	claims, err := jwt.ValidateRefreshToken(tokenString)
	if err != nil {
		rlog.Error("failed to validate refresh token", "error", err)
		return nil, ErrInvalidRefreshToken
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
	accessToken, refreshToken, err := s.generateTokens(ctx, data)
	if err != nil {
		rlog.Error("failed to generate new tokens", "user_id", user.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &LoginResponse{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		OfficeID:     user.OfficeID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
		PhoneNumber:  ptrToStr(user.PhoneNumber),
	}, nil
}
