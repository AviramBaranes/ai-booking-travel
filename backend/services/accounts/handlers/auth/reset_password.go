package auth

import (
	"context"
	"errors"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.app/services/accounts/handlers/user"
	"encore.dev/rlog"
)

// ResetPasswordParams is the request payload for the ResetPassword endpoint.
type ResetPasswordParams struct {
	Token       string `json:"token" validate:"required" encore:"sensitive"`
	NewPassword string `json:"newPassword" validate:"required,min=8" encore:"sensitive"`
}

func (p ResetPasswordParams) Validate() error {
	if err := user.ValidatePasswordForAPI(p.NewPassword); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

func (s *AuthService) ResetPassword(ctx context.Context, p ResetPasswordParams) error {
	hashedToken := hashPasswordResetToken(p.Token)

	token, err := s.query.GetPasswordResetTokenByHash(ctx, hashedToken)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to get password reset token by hash", "error", err)
		return api_errors.ErrInternalError
	}

	if token.ExpiresAt.Time.Before(time.Now()) {
		return api_errors.ErrNotFound
	}

	newHashedPassword, err := password.HashPassword(p.NewPassword)
	if err != nil {
		rlog.Error("failed to hash new password", "error", err)
		return api_errors.ErrInternalError
	}

	if _, err := s.query.UpdateUser(ctx, db.UpdateUserParams{
		ID:           token.UserID,
		PasswordHash: &newHashedPassword,
	}); err != nil {
		rlog.Error("failed to update user password", "user_id", token.UserID, "error", err)
		return api_errors.ErrInternalError
	}

	if err := s.query.DeletePasswordResetTokenByID(ctx, token.ID); err != nil {
		rlog.Error("failed to delete password reset token", "token_id", token.ID, "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
