package accounts

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// RegisterAdminParams defines the parameters required to register an admin user.
type RegisterAdminParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8" encore:"sensitive"`
}

// RegisterAdminResponse represents the response returned after registering an admin user.
type RegisterAdminResponse struct {
	ID int32 `json:"id"`
}

// Validate performs basic validation of the registration params.
func (p RegisterAdminParams) Validate() error {
	if err := validatePasswordForAPI(p.Password); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

var (
	ErrOfficeAgentCodeMismatch = api_errors.NewValidationError("office code and agent code must be provided together")
)

// RegisterAdmin registers a new admin user.
// encore:api auth path=/register-admin method=POST tag:admin
func (s *Service) RegisterAdmin(ctx context.Context, params RegisterAdminParams) (*RegisterAdminResponse, error) {
	userID, err := s.query.CheckUserExists(ctx, params.Email)
	if err != nil && !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to check if user exists", "email", params.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}
	if userID != 0 {
		return nil, ErrEmailAlreadyExists
	}

	hashed, err := password.HashPassword(params.Password)
	if err != nil {
		rlog.Error("failed to hash password", "email", params.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	row, err := s.query.RegisterAdmin(ctx, db.RegisterAdminParams{
		Email:        params.Email,
		PasswordHash: hashed,
	})

	if err != nil {
		rlog.Error("failed to register admin user", "email", params.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &RegisterAdminResponse{
		ID: row.ID,
	}, nil
}
