package accounts

import (
	"context"
	"errors"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// CreateStaffRequest is the shared request type for creating admin and accountant users.
type CreateStaffRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8" encore:"sensitive"`
}

func (p CreateStaffRequest) Validate() error {
	if err := validatePasswordForAPI(p.Password); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

// CreateStaffResponse is the shared response type for creating admin and accountant users.
type CreateStaffResponse struct {
	ID int32 `json:"id"`
}

// StaffResponse is the shared response type for listing admin and accountant users.
type StaffResponse struct {
	ID        int32      `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	LastLogin *time.Time `json:"lastLogin"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// checkEmailAvailable returns ErrEmailAlreadyExists if the email is taken,
// or ErrInternalError on a database failure.
func (s *Service) checkEmailAvailable(ctx context.Context, email string) error {
	userID, err := s.query.CheckUserExists(ctx, email)
	if err != nil && !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to check if user exists", "email", email, "error", err)
		return api_errors.ErrInternalError
	}
	if userID != 0 {
		return ErrEmailAlreadyExists
	}
	return nil
}

// hashPasswordForInsert hashes rawPassword, returning ErrInternalError on failure.
func hashPasswordForInsert(email, rawPassword string) (string, error) {
	hashed, err := password.HashPassword(rawPassword)
	if err != nil {
		rlog.Error("failed to hash password", "email", email, "error", err)
		return "", api_errors.ErrInternalError
	}
	return hashed, nil
}

// toStaffResponse maps a ListStaffByRoleRow to StaffResponse.
func toStaffResponse(r db.ListStaffByRoleRow) StaffResponse {
	return StaffResponse{
		ID:        r.ID,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		LastLogin: db.TimePtrFromDB(r.LastLogin),
		CreatedAt: db.TimeFromDB(r.CreatedAt),
		UpdatedAt: db.TimeFromDB(r.UpdatedAt),
	}
}

// createStaffUser checks email availability, hashes the password, and inserts the user.
func (s *Service) createStaffUser(ctx context.Context, role db.UserRole, params CreateStaffRequest) (*CreateStaffResponse, error) {
	if err := s.checkEmailAvailable(ctx, params.Email); err != nil {
		return nil, err
	}

	hashed, err := hashPasswordForInsert(params.Email, params.Password)
	if err != nil {
		return nil, err
	}

	row, err := s.query.CreateStaffUser(ctx, db.CreateStaffUserParams{
		Role:         role,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Email:        params.Email,
		PasswordHash: hashed,
	})
	if err != nil {
		rlog.Error("failed to create staff user", "role", role, "email", params.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &CreateStaffResponse{ID: row.ID}, nil
}

// listStaffByRole fetches all users with the given role.
func (s *Service) listStaffByRole(ctx context.Context, role db.UserRole) ([]StaffResponse, error) {
	rows, err := s.query.ListStaffByRole(ctx, role)
	if err != nil {
		rlog.Error("failed to list staff by role", "role", role, "error", err)
		return nil, api_errors.ErrInternalError
	}

	result := make([]StaffResponse, 0, len(rows))
	for _, r := range rows {
		result = append(result, toStaffResponse(r))
	}
	return result, nil
}
