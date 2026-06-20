package accounts

import (
	"context"

	user "encore.app/services/accounts/handlers/user"
	"encore.dev/beta/auth"
)

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

// encore:api private tag:agent
func (s *Service) GetUserCredit(ctx context.Context) (*user.GetUserCreditResponse, error) {
	authData := auth.Data().(*AuthData)
	h := user.NewUserService(s.query)
	return h.GetUserCredit(ctx, authData.UserID)
}

// encore:api private
func (s *Service) GetUserMarkupGross(ctx context.Context) (*user.GetUserMarkupGrossResponse, error) {
	authData := auth.Data().(*AuthData)
	h := user.NewUserService(s.query)
	return h.GetUserMarkupGross(ctx, authData.UserID)
}
