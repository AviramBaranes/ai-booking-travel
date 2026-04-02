package accounts

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// RegisterAgentParams defines the parameters required to register an agent user.
type RegisterAgentParams struct {
	Username   string `json:"username" validate:"required,username"`
	Password   string `json:"password" validate:"required,min=8" encore:"sensitive"`
	OfficeCode string `json:"office_code" validate:"required"`
	AgentCode  string `json:"agent_code" validate:"required"`
}

// RegisterAgentResponse represents the response returned after registering an agent user.
type RegisterAgentResponse struct {
	ID int32 `json:"id"`
}

// Validate performs basic validation of the registration params.
func (p RegisterAgentParams) Validate() error {
	if err := validatePasswordForAPI(p.Password); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

// RegisterAgent registers a new agent user.
// encore:api auth path=/register-agent method=POST tag:admin
func (s *Service) RegisterAgent(ctx context.Context, params RegisterAgentParams) (*RegisterAgentResponse, error) {
	userID, err := s.query.CheckUserExists(ctx, params.Username)
	if err != nil && !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to check if user exists", "username", params.Username, "error", err)
		return nil, api_errors.ErrInternalError
	}
	if userID != 0 {
		return nil, ErrUsernameAlreadyExists
	}

	hashed, err := password.HashPassword(params.Password)
	if err != nil {
		rlog.Error("failed to hash password", "username", params.Username, "error", err)
		return nil, api_errors.ErrInternalError
	}

	row, err := s.query.RegisterAgent(ctx, db.RegisterAgentParams{
		Username:     params.Username,
		PasswordHash: hashed,
		AgentCode:    db.TextParam(&params.AgentCode),
		OfficeCode:   db.TextParam(&params.OfficeCode),
	})

	if err != nil {
		rlog.Error("failed to register agent user", "username", params.Username, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &RegisterAgentResponse{
		ID: row.ID,
	}, nil
}
