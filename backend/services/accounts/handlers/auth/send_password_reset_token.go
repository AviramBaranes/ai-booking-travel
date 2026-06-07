package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

const (
	passwordResetTokenExpiry = time.Hour
)

type SendPasswordResetTokenParams struct {
	Email string `json:"email" validate:"required,email"`
}

func (p SendPasswordResetTokenParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *AuthService) SendPasswordResetToken(ctx context.Context, p SendPasswordResetTokenParams) error {
	user, err := s.query.GetUserByEmail(ctx, p.Email)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil // Don't reveal that the email doesn't exist
		}
		rlog.Error("failed to get user by email", "email", p.Email, "error", err)
		return api_errors.ErrInternalError
	}

	if user.Role == db.UserRoleCustomer {
		return ErrCustomerPasswordResetNotAllowed
	}

	rawToken, err := generatePasswordResetToken()
	if err != nil {
		rlog.Error("failed to generate password reset token", "error", err)
		return api_errors.ErrInternalError
	}

	hashed, err := password.HashPassword(rawToken)
	if err != nil {
		rlog.Error("failed to hash password reset token", "error", err)
		return api_errors.ErrInternalError
	}

	_, err = s.query.InsertPasswordResetToken(ctx, db.InsertPasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: hashed,
		ExpiresAt: dbadapters.DBTime(time.Now().Add(passwordResetTokenExpiry)),
	})

	if err != nil {
		rlog.Error("failed to save password reset token", "user_id", user.ID, "error", err)
		return api_errors.ErrInternalError
	}

	// TODO: FIX import cycle
	// notifications.PublishEmailEvent(ctx, notifications.EmailEventTypePasswordReset, notifications.PasswordResetEmailPayload{
	// 	Email:     user.Email,
	// 	TokenHash: hashed,
	// })

	return nil
}

func (s *AuthService) getUserForPasswordReset() {

}

func generatePasswordResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
