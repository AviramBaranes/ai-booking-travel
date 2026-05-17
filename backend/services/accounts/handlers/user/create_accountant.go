package user

import (
	"context"

	"encore.app/services/accounts/db"
)

// CreateAccountantParams are the params for creating an accountant.
type CreateAccountantParams struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8" encore:"sensitive"`
}

func (p CreateAccountantParams) Validate() error {
	return CreateStaffParams(p).Validate()
}

// CreateAccountantResponse is the response for creating an accountant.
type CreateAccountantResponse struct {
	ID int64 `json:"id"`
}

// CreateAccountant creates a new accountant user.
func (s *UserService) CreateAccountant(ctx context.Context, params CreateAccountantParams) (*CreateAccountantResponse, error) {
	resp, err := s.createStaffUser(ctx, db.UserRoleAccountant, CreateStaffParams(params))
	if err != nil {
		return nil, err
	}
	return &CreateAccountantResponse{ID: resp.ID}, nil
}
