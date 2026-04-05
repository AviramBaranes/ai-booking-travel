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

// --- Request / Response types ---

type UpdateUserRequest struct {
	Email       *string `json:"email" validate:"omitempty,email" encore:"optional"`
	PhoneNumber *string `json:"phoneNumber" encore:"optional"`
	OfficeID    *int32  `json:"officeId" validate:"omitempty,gte=1" encore:"optional"`
	Password    *string `json:"password" validate:"omitempty,min=8" encore:"sensitive,optional"`
}

func (p UpdateUserRequest) Validate() error {
	if p.Password != nil {
		if err := validatePasswordForAPI(*p.Password); err != nil {
			return err
		}
	}
	return validation.ValidateStruct(p)
}

type UpdateUserResponse struct {
	ID          int32   `json:"id"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	OfficeID    *int32  `json:"officeId"`
}

// --- Endpoints ---

// UpdateUser updates an existing user.
//
//encore:api auth method=PUT path=/users/:id tag:admin
func (s *Service) UpdateUser(ctx context.Context, id int32, params UpdateUserRequest) (*UpdateUserResponse, error) {
	// Check email uniqueness
	if params.Email != nil {
		existingID, err := s.query.CheckUserExists(ctx, *params.Email)
		if err != nil && !errors.Is(err, db.ErrNoRows) {
			rlog.Error("failed to check email uniqueness", "error", err)
			return nil, api_errors.ErrInternalError
		}
		if existingID != 0 && existingID != id {
			return nil, ErrEmailAlreadyExists
		}
	}

	// Check phone uniqueness
	if params.PhoneNumber != nil {
		existingID, err := s.query.GetUserByPhone(ctx, params.PhoneNumber)
		if err != nil && !errors.Is(err, db.ErrNoRows) {
			rlog.Error("failed to check phone uniqueness", "error", err)
			return nil, api_errors.ErrInternalError
		}
		if existingID != 0 && existingID != id {
			return nil, ErrPhoneAlreadyExists
		}
	}

	// Hash password if provided
	var hashedPtr *string
	if params.Password != nil {
		hashed, err := password.HashPassword(*params.Password)
		if err != nil {
			rlog.Error("failed to hash password", "error", err)
			return nil, api_errors.ErrInternalError
		}
		hashedPtr = &hashed
	}

	row, err := s.query.UpdateUser(ctx, db.UpdateUserParams{
		ID:           id,
		Email:        params.Email,
		PhoneNumber:  params.PhoneNumber,
		OfficeID:     params.OfficeID,
		PasswordHash: hashedPtr,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		rlog.Error("failed to update user", "id", id, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &UpdateUserResponse{
		ID:          row.ID,
		Email:       row.Email,
		PhoneNumber: row.PhoneNumber,
		OfficeID:    row.OfficeID,
	}, nil
}
