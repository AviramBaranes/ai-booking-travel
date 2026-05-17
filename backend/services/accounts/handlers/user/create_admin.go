package user

import (
	"context"

	"encore.app/services/accounts/db"
)

// CreateAdminParams are the params for creating an admin.
type CreateAdminParams struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8" encore:"sensitive"`
}

func (p CreateAdminParams) Validate() error {
	return CreateStaffParams(p).Validate()
}

// CreateAdminResponse is the response for creating an admin.
type CreateAdminResponse struct {
	ID int64 `json:"id"`
}

// CreateAdmin creates a new admin user.
func (s *UserService) CreateAdmin(ctx context.Context, params CreateAdminParams) (*CreateAdminResponse, error) {
	resp, err := s.createStaffUser(ctx, db.UserRoleAdmin, CreateStaffParams(params))
	if err != nil {
		return nil, err
	}
	return &CreateAdminResponse{ID: resp.ID}, nil
}
