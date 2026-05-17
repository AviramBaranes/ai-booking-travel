package user

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// UpdateUserParams are the params for updating a user.
type UpdateUserParams struct {
	FirstName   *string `json:"firstName" encore:"optional"`
	LastName    *string `json:"lastName" encore:"optional"`
	Email       *string `json:"email" validate:"omitempty,email" encore:"optional"`
	PhoneNumber *string `json:"phoneNumber" encore:"optional"`
	OfficeID    *int64  `json:"officeId" validate:"omitempty,gte=1" encore:"optional"`
	Password    *string `json:"password" validate:"omitempty,min=8" encore:"sensitive,optional"`
}

func (p UpdateUserParams) Validate() error {
	if p.Password != nil {
		if err := validatePasswordForAPI(*p.Password); err != nil {
			return err
		}
	}
	return validation.ValidateStruct(p)
}

// UpdateUserResponse is the response for updating a user.
type UpdateUserResponse struct {
	ID          int64   `json:"id"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	OfficeID    *int64  `json:"officeId"`
}

// UpdateUser updates an existing user.
func (s *UserService) UpdateUser(ctx context.Context, id int64, params UpdateUserParams) (*UpdateUserResponse, error) {
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
		user, err := s.query.GetUserByPhone(ctx, params.PhoneNumber)
		if err != nil && !errors.Is(err, db.ErrNoRows) {
			rlog.Error("failed to check phone uniqueness", "error", err)
			return nil, api_errors.ErrInternalError
		}
		if user.ID != 0 && user.ID != id {
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
		FirstName:    params.FirstName,
		LastName:     params.LastName,
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
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		Email:       row.Email,
		PhoneNumber: row.PhoneNumber,
		OfficeID:    row.OfficeID,
	}, nil
}
