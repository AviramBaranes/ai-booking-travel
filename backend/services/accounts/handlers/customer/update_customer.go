package customer

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type UpdateCustomerParams struct {
	FirstName   string `json:"firstName" validate:"required,notblank"`
	LastName    string `json:"lastName" validate:"required,notblank"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"omitempty,israeli_phone"`
}

func (p UpdateCustomerParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, p UpdateCustomerParams, userID int64) error {
	if _, err := s.query.UpdateUser(ctx, db.UpdateUserParams{
		ID:          userID,
		FirstName:   &p.FirstName,
		LastName:    &p.LastName,
		Email:       &p.Email,
		PhoneNumber: &p.PhoneNumber,
	}); err != nil {
		rlog.Error("failed to update customer", "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
