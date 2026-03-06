package auth

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/auth/db"
	"encore.dev/rlog"
)

type UpdateUserResponse struct {
	ID          int32       `json:"id"`
	Role        db.UserRole `json:"role"`
	Username    string      `json:"username"`
	AgentCode   string      `json:"agent_code"`
	OfficeCode  string      `json:"office_code"`
	PhoneNumber string      `json:"phone_number"`
	LastLogin   string      `json:"last_login"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type UpdateUserParams struct {
	ID          int32   `json:"id" validate:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	OfficeCode  *string `json:"office_code,omitempty"`
	AgentCode   *string `json:"agent_code,omitempty"`
}

func (p UpdateUserParams) Validate() error {
	return validation.ValidateStruct(p)
}

// UpdateUser updates user details such as phone number, office code, and agent code.
// encore:api auth path=/update-user method=PUT tag:admin
func (s *Service) UpdateUser(ctx context.Context, params UpdateUserParams) (*UpdateUserResponse, error) {
	row, err := s.query.UpdateUser(ctx, db.UpdateUserParams{
		ID:          params.ID,
		AgentCode:   db.TextParam(params.AgentCode),
		PhoneNumber: db.TextParam(params.PhoneNumber),
		OfficeCode:  db.TextParam(params.OfficeCode),
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		rlog.Error("failed to update user", "userID", params.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &UpdateUserResponse{
		ID:          row.ID,
		Role:        row.Role,
		Username:    row.Username,
		AgentCode:   db.StringFromTextParam(row.AgentCode),
		OfficeCode:  db.StringFromTextParam(row.OfficeCode),
		PhoneNumber: db.StringFromTextParam(row.PhoneNumber),
		LastLogin:   db.StringFromTimeParam(row.LastLogin),
		CreatedAt:   db.StringFromTimeParam(row.CreatedAt),
		UpdatedAt:   db.StringFromTimeParam(row.UpdatedAt),
	}, nil
}
