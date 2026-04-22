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

var secrets struct {
	FirstAdminEmail    string
	FirstAdminPassword string
}

// --- Request / Response types ---

type AdminResponse struct {
	ID        int32      `json:"id"`
	Email     string     `json:"email"`
	LastLogin *time.Time `json:"lastLogin"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type ListAdminsResponse struct {
	Admins []AdminResponse `json:"admins"`
}

type CreateAdminRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8" encore:"sensitive"`
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

func toAdminResponse(r db.ListAdminsRow) AdminResponse {
	return AdminResponse{
		ID:        r.ID,
		Email:     r.Email,
		LastLogin: db.TimePtrFromDB(r.LastLogin),
		CreatedAt: db.TimeFromDB(r.CreatedAt),
		UpdatedAt: db.TimeFromDB(r.UpdatedAt),
	}
}

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

	_, err = query.CreateAdmin(ctx, db.CreateAdminParams{
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
	rows, err := s.query.ListAdmins(ctx)
	if err != nil {
		rlog.Error("failed to list admins", "error", err)
		return nil, api_errors.ErrInternalError
	}

	admins := make([]AdminResponse, 0, len(rows))
	for _, r := range rows {
		admins = append(admins, toAdminResponse(r))
	}

	return &ListAdminsResponse{Admins: admins}, nil
}

// CreateAdmin creates a new admin user.
//
//encore:api auth method=POST path=/admins tag:admin
func (s *Service) CreateAdmin(ctx context.Context, params CreateAdminRequest) (*CreateAdminResponse, error) {
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

	row, err := s.query.CreateAdmin(ctx, db.CreateAdminParams{
		Email:        params.Email,
		PasswordHash: hashed,
	})
	if err != nil {
		rlog.Error("failed to create admin user", "email", params.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &CreateAdminResponse{
		ID: row.ID,
	}, nil
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
