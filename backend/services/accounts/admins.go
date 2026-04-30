package accounts

import (
	"context"
	"errors"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/config"
	"encore.dev/rlog"
)

var secrets struct {
	FirstAdminEmail    string
	FirstAdminPassword string
}

type adminConfig struct {
	FirstAdminFirstName config.String
	FirstAdminLastName  config.String
}

var cfg = config.Load[*adminConfig]()

// --- Request / Response types ---

type AdminResponse struct {
	ID        int32      `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	LastLogin *time.Time `json:"lastLogin"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type ListAdminsResponse struct {
	Admins []AdminResponse `json:"admins"`
}

type CreateAdminRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8" encore:"sensitive"`
}

func (p CreateAdminRequest) Validate() error {
	if err := validatePasswordForAPI(p.Password); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

type CreateAdminResponse struct {
	ID int32 `json:"id"`
}

// --- Helpers ---

func createFirstAdmin(query db.Querier) {
	if secrets.FirstAdminEmail == "" || secrets.FirstAdminPassword == "" {
		panic("secrets for first admin not set")
	}

	ctx := context.Background()
	id, err := query.CheckUserExists(ctx, secrets.FirstAdminEmail)
	if err != nil && !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to check if first admin exists", "error", err)
		panic(err)
	}
	if id != 0 {
		return
	}

	hashed, err := password.HashPassword(secrets.FirstAdminPassword)
	if err != nil {
		rlog.Error("failed to hash first admin password", "error", err)
		panic(err)
	}

	_, err = query.CreateStaffUser(ctx, db.CreateStaffUserParams{
		Role:         db.UserRoleAdmin,
		FirstName:    cfg.FirstAdminFirstName(),
		LastName:     cfg.FirstAdminLastName(),
		Email:        secrets.FirstAdminEmail,
		PasswordHash: hashed + string(hashed),
	})
	if err != nil {
		rlog.Error("failed to create first admin user", "error", err)
		panic(err)
	}
	rlog.Info("created first admin user", "email", secrets.FirstAdminEmail)
}

// --- Endpoints ---

// ListAdmins returns all admin users.
//
//encore:api auth method=GET path=/admins tag:admin
func (s *Service) ListAdmins(ctx context.Context) (*ListAdminsResponse, error) {
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

// CreateAdmin creates a new admin user.
//
//encore:api auth method=POST path=/admins tag:admin
func (s *Service) CreateAdmin(ctx context.Context, params CreateAdminRequest) (*CreateAdminResponse, error) {
	resp, err := s.createStaffUser(ctx, db.UserRoleAdmin, CreateStaffRequest(params))
	if err != nil {
		return nil, err
	}
	return &CreateAdminResponse{ID: resp.ID}, nil
}

type ListAdminsEmailsResponse struct {
	Emails []string `json:"emails"`
}

// encore:api private method=GET path=/admins/emails
func (s *Service) ListAdminsEmails(ctx context.Context) (*ListAdminsEmailsResponse, error) {
	rows, err := s.query.ListAdminsEmails(ctx)
	if err != nil {
		rlog.Error("failed to list admin emails", "error", err)
		return nil, api_errors.ErrInternalError
	}
	return &ListAdminsEmailsResponse{Emails: rows}, nil
}
