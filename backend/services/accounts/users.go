package accounts

import (
	"context"
	"errors"

	"encore.app/internal/password"
	"encore.app/services/accounts/db"
	user_handlers "encore.app/services/accounts/handlers/user_handlers"
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
func (s *Service) ListAdmins(ctx context.Context) (*user_handlers.ListAdminsResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.ListAdmins(ctx)
}

// CreateAdmin creates a new admin user.
//
//encore:api auth method=POST path=/admins tag:admin
func (s *Service) CreateAdmin(ctx context.Context, params user_handlers.CreateAdminParams) (*user_handlers.CreateAdminResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.CreateAdmin(ctx, params)
}

// encore:api private method=GET path=/admins/emails
func (s *Service) ListAdminsEmails(ctx context.Context) (*user_handlers.ListAdminsEmailsResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.ListAdminsEmails(ctx)
}

// UpdateUser updates an existing user.
//
//encore:api auth method=PUT path=/users/:id tag:admin
func (s *Service) UpdateUser(ctx context.Context, id int64, params user_handlers.UpdateUserParams) (*user_handlers.UpdateUserResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.UpdateUser(ctx, id, params)
}

// ListAgents lists agents with optional filtering and pagination.
//
//encore:api auth method=GET path=/agents tag:admin
func (s *Service) ListAgents(ctx context.Context, params *user_handlers.ListAgentsParams) (*user_handlers.ListAgentsResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.ListAgents(ctx, params)
}

// CreateAgent creates a new agent user.
//
//encore:api auth method=POST path=/agents tag:admin
func (s *Service) CreateAgent(ctx context.Context, params user_handlers.CreateAgentParams) (*user_handlers.CreateAgentResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.CreateAgent(ctx, params)
}

// GetAgentsByOfficeID retrieves agent IDs for a given office ID.
//
// encore:api private
func (s *Service) GetAgentsByOfficeID(ctx context.Context, params user_handlers.GetAgentsByOfficeIDParams) (*user_handlers.GetAgentsResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.GetAgentsByOfficeID(ctx, params)
}

// GetAgentsByOrganizationID retrieves agent IDs for a given organization ID.
//
// encore:api private
func (s *Service) GetAgentsByOrganizationID(ctx context.Context, params user_handlers.GetAgentsByOrganizationIDParams) (*user_handlers.GetAgentsResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.GetAgentsByOrganizationID(ctx, params)
}

// ListAccountants returns all accountant users.
//
//encore:api auth method=GET path=/accountants tag:admin
func (s *Service) ListAccountants(ctx context.Context) (*user_handlers.ListAccountantsResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.ListAccountants(ctx)
}

// CreateAccountant creates a new accountant user.
//
//encore:api auth method=POST path=/accountants tag:admin
func (s *Service) CreateAccountant(ctx context.Context, params user_handlers.CreateAccountantParams) (*user_handlers.CreateAccountantResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.CreateAccountant(ctx, params)
}

// encore:api private
func (s *Service) GetUserEmail(ctx context.Context, params user_handlers.GetUserEmailParams) (*user_handlers.GetUserEmailResponse, error) {
	h := user_handlers.NewUserService(s.query)
	return h.GetUserEmail(ctx, params)
}
