package accounts

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type UpdateUserResponse struct {
	ID          int32       `json:"id"`
	Role        db.UserRole `json:"role"`
	Email       string      `json:"email,omitempty"`
	OfficeID    *int32      `json:"office_id,omitempty"`
	PhoneNumber string      `json:"phone_number"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type UpdateUserParams struct {
	ID          int32   `json:"id" validate:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	OfficeID    *int32  `json:"office_id,omitempty"`
}

func (p UpdateUserParams) Validate() error {
	return validation.ValidateStruct(p)
}

// UpdateUser updates user details such as phone number, office code, and agent code.
// encore:api auth path=/update-user method=PUT tag:admin
func (s *Service) UpdateUser(ctx context.Context, params UpdateUserParams) (*UpdateUserResponse, error) {
	row, err := s.query.UpdateUser(ctx, db.UpdateUserParams{
		ID:          params.ID,
		PhoneNumber: params.PhoneNumber,
		OfficeID:    params.OfficeID,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		rlog.Error("failed to update user", "userID", params.ID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	var phoneNumber string
	if row.PhoneNumber != nil {
		phoneNumber = *row.PhoneNumber
	}

	return &UpdateUserResponse{
		ID:          row.ID,
		Role:        row.Role,
		Email:       row.Email,
		PhoneNumber: phoneNumber,
		OfficeID:    row.OfficeID,
		CreatedAt:   db.StringFromTimeParam(row.CreatedAt),
		UpdatedAt:   db.StringFromTimeParam(row.UpdatedAt),
	}, nil
}
