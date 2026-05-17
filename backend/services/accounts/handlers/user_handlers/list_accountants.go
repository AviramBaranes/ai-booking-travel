package user_handlers

import (
	"context"
	"time"

	"encore.app/services/accounts/db"
)

// AccountantResponse is the response type for a single accountant.
type AccountantResponse struct {
	ID        int64      `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	LastLogin *time.Time `json:"lastLogin"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// ListAccountantsResponse is the response for listing accountants.
type ListAccountantsResponse struct {
	Accountants []AccountantResponse `json:"accountants"`
}

// ListAccountants returns all accountant users.
func (s *UserService) ListAccountants(ctx context.Context) (*ListAccountantsResponse, error) {
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
