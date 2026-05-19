package accounts

import (
	"context"
	"errors"

	"encore.app/internal/password"
	"encore.app/services/accounts/db"
	user "encore.app/services/accounts/handlers/user"
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
func (s *Service) ListAdmins(ctx context.Context) (*user.ListAdminsResponse, error) {
	h := user.NewUserService(s.query)
	return h.ListAdmins(ctx)
}

// CreateAdmin creates a new admin user.
//
//encore:api auth method=POST path=/admins tag:admin
func (s *Service) CreateAdmin(ctx context.Context, params user.CreateAdminParams) (*user.CreateAdminResponse, error) {
	h := user.NewUserService(s.query)
	return h.CreateAdmin(ctx, params)
}

// encore:api private method=GET path=/admins/emails
func (s *Service) ListAdminsEmails(ctx context.Context) (*user.ListAdminsEmailsResponse, error) {
	h := user.NewUserService(s.query)
	return h.ListAdminsEmails(ctx)
}

// UpdateUser updates an existing user.
//
//encore:api auth method=PUT path=/users/:id tag:admin
func (s *Service) UpdateUser(ctx context.Context, id int64, params user.UpdateUserParams) (*user.UpdateUserResponse, error) {
	h := user.NewUserService(s.query)
	return h.UpdateUser(ctx, id, params)
}

// ListAgents lists agents with optional filtering and pagination.
//
//encore:api auth method=GET path=/agents tag:admin
func (s *Service) ListAgents(ctx context.Context, params *user.ListAgentsParams) (*user.ListAgentsResponse, error) {
	h := user.NewUserService(s.query)
	return h.ListAgents(ctx, params)
}

// CreateAgent creates a new agent user.
//
//encore:api auth method=POST path=/agents tag:admin
func (s *Service) CreateAgent(ctx context.Context, params user.CreateAgentParams) (*user.CreateAgentResponse, error) {
	h := user.NewUserService(s.query)
	return h.CreateAgent(ctx, params)
}

// ListAccountants returns all accountant users.
//
//encore:api auth method=GET path=/accountants tag:admin
func (s *Service) ListAccountants(ctx context.Context) (*user.ListAccountantsResponse, error) {
	h := user.NewUserService(s.query)
	return h.ListAccountants(ctx)
}

// CreateAccountant creates a new accountant user.
//
//encore:api auth method=POST path=/accountants tag:admin
func (s *Service) CreateAccountant(ctx context.Context, params user.CreateAccountantParams) (*user.CreateAccountantResponse, error) {
	h := user.NewUserService(s.query)
	return h.CreateAccountant(ctx, params)
}

// encore:api private
func (s *Service) GetUserEmail(ctx context.Context, params user.GetUserEmailParams) (*user.GetUserEmailResponse, error) {
	h := user.NewUserService(s.query)
	return h.GetUserEmail(ctx, params)
}
