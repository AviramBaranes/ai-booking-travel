package user

import (
	"context"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// AdminResponse is the response type for a single admin.
type AdminResponse struct {
	ID        int64      `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	LastLogin *time.Time `json:"lastLogin"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// ListAdminsResponse is the response for listing admins.
type ListAdminsResponse struct {
	Admins []AdminResponse `json:"admins"`
}

// ListAdminsEmailsResponse is the response for listing admin email addresses.
type ListAdminsEmailsResponse struct {
	Emails []string `json:"emails"`
}

// ListAdmins returns all admin users.
func (s *UserService) ListAdmins(ctx context.Context) (*ListAdminsResponse, error) {
	staff, err := s.listStaffByRole(ctx, db.UserRoleAdmin)
	if err != nil {
		return nil, err
	}

	admins := make([]AdminResponse, len(staff))
	for i, r := range staff {
		admins[i] = AdminResponse(r)
	}
	return &ListAdminsResponse{Admins: admins}, nil
}

// ListAdminsEmails returns all admin email addresses.
func (s *UserService) ListAdminsEmails(ctx context.Context) (*ListAdminsEmailsResponse, error) {
	rows, err := s.query.ListAdminsEmails(ctx)
	if err != nil {
		rlog.Error("failed to list admin emails", "error", err)
		return nil, api_errors.ErrInternalError
	}
	return &ListAdminsEmailsResponse{Emails: rows}, nil
}
