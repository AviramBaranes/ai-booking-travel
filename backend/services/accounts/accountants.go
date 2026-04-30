package accounts

import (
	"context"
	"time"

	"encore.app/services/accounts/db"
)

// --- Request / Response types ---

type AccountantResponse struct {
	ID        int32      `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	LastLogin *time.Time `json:"lastLogin"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type ListAccountantsResponse struct {
	Accountants []AccountantResponse `json:"accountants"`
}

type CreateAccountantRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8" encore:"sensitive"`
}

func (p CreateAccountantRequest) Validate() error {
	return CreateStaffRequest(p).Validate()
}

type CreateAccountantResponse struct {
	ID int32 `json:"id"`
}

// --- Endpoints ---

// ListAccountants returns all accountant users.
//
//encore:api auth method=GET path=/accountants tag:admin
func (s *Service) ListAccountants(ctx context.Context) (*ListAccountantsResponse, error) {
	staff, err := s.listStaffByRole(ctx, db.UserRoleAccountant)
	if err != nil {
		return nil, err
	}

	accountants := make([]AccountantResponse, len(staff))
	for i, r := range staff {
		accountants[i] = AccountantResponse(r)
	}
	return &ListAccountantsResponse{Accountants: accountants}, nil
}

// CreateAccountant creates a new accountant user.
//
//encore:api auth method=POST path=/accountants tag:admin
func (s *Service) CreateAccountant(ctx context.Context, params CreateAccountantRequest) (*CreateAccountantResponse, error) {
	resp, err := s.createStaffUser(ctx, db.UserRoleAccountant, CreateStaffRequest(params))
	if err != nil {
		return nil, err
	}
	return &CreateAccountantResponse{ID: resp.ID}, nil
}

