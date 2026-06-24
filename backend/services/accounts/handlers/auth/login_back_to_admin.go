package auth

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// LoginBackToAdmin logs back into the admin account from an impersonated agent session.
// adminRefID is extracted from the caller's auth context by the thin wrapper.
func (s *AuthService) LoginBackToAdmin(ctx context.Context, adminRefID *int64) (*LoginResponse, error) {
	if adminRefID == nil {
		return nil, ErrInvalidCredentials
	}

	admin, err := s.query.GetUserById(ctx, *adminRefID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		rlog.Error("failed to get admin by ID in login back to admin", "admin_id", *adminRefID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	tokens, err := s.generateTokens(ctx, jwt.AccessTokenData{
		UserID: admin.ID,
		Role:   admin.Role,
	})

	if err != nil {
		rlog.Error("failed to generate tokens in login back to admin", "user_id", admin.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &LoginResponse{
		ID:                   admin.ID,
		FirstName:            admin.FirstName,
		LastName:             admin.LastName,
		Role:                 admin.Role,
		AccessToken:          tokens.AccessToken,
		RefreshToken:         tokens.RefreshToken,
		AccessTokenExpiresAt: tokens.AccessTokenExpiresAt,
		Email:                admin.Email,
	}, nil
}
