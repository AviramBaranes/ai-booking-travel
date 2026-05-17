package user_handlers

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// CreateAgentParams are the params for creating an agent.
type CreateAgentParams struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8" encore:"sensitive"`
	PhoneNumber string `json:"phoneNumber" validate:"required,israeli_phone"`
	OfficeID    int64  `json:"officeId" validate:"required,gte=1"`
}

func (p CreateAgentParams) Validate() error {
	if err := validatePasswordForAPI(p.Password); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

// CreateAgentResponse is the response for creating an agent.
type CreateAgentResponse struct {
	ID int64 `json:"id"`
}

// CreateAgent creates a new agent user.
func (s *UserService) CreateAgent(ctx context.Context, params CreateAgentParams) (*CreateAgentResponse, error) {
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

	row, err := s.query.CreateAgent(ctx, db.CreateAgentParams{
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Email:        params.Email,
		PhoneNumber:  &params.PhoneNumber,
		PasswordHash: hashed,
		OfficeID:     &params.OfficeID,
	})
	if err != nil {
		rlog.Error("failed to create agent user", "email", params.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &CreateAgentResponse{ID: row.ID}, nil
}
