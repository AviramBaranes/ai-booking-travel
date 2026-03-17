package auth

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/auth/db"
	"encore.dev/rlog"
)

// ChangePasswordParams defines the parameters for changing a user's password.
type ChangePasswordParams struct {
	ID          int32  `json:"id" validate:"required,gt=0"`
	NewPassword string `json:"new_password" validate:"required,min=8" encore:"sensitive"`
}

// Validate performs validation on ChangePasswordParams.
func (p ChangePasswordParams) Validate() error {
	if err := validatePasswordForAPI(p.NewPassword); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

// ChangePassword changes the password for a given user.
// encore:api auth path=/change-password method=PUT tag:admin
func (s *Service) ChangePassword(ctx context.Context, params ChangePasswordParams) error {
	hashed, err := password.HashPassword(params.NewPassword)
	if err != nil {
		rlog.Error("failed to hash password", "id", params.ID, "error", err)
		return api_errors.ErrInternalError
	}

	err = s.query.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           params.ID,
		PasswordHash: hashed,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return ErrUserNotFound
		}

		rlog.Error("failed to update user password", "id", params.ID, "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
