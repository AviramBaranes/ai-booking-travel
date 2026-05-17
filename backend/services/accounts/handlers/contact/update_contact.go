package contact

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type UpdateContactParams struct {
	FirstName            *string `json:"firstName" validate:"omitempty,notblank" encore:"optional"`
	LastName             *string `json:"lastName" validate:"omitempty,notblank" encore:"optional"`
	Role                 *string `json:"role" validate:"omitempty,notblank" encore:"optional"`
	Cellphone            *string `json:"cellphone" validate:"omitempty,notblank" encore:"optional"`
	Email                *string `json:"email" validate:"omitempty,email" encore:"optional"`
	OfficeID             *int64  `json:"officeId" encore:"optional"`
	OrganizationID       *int64  `json:"organizationId" encore:"optional"`
	IsPaymentResponsible *bool   `json:"isPaymentResponsible" encore:"optional"`
}

func (p UpdateContactParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *ContactService) UpdateContact(ctx context.Context, id int64, p UpdateContactParams) (*ContactResponse, error) {
	row, err := s.query.UpdateContact(ctx, db.UpdateContactParams{
		ID:                   id,
		FirstName:            p.FirstName,
		LastName:             p.LastName,
		Role:                 p.Role,
		Cellphone:            p.Cellphone,
		Email:                p.Email,
		OfficeID:             p.OfficeID,
		OrganizationID:       p.OrganizationID,
		IsPaymentResponsible: p.IsPaymentResponsible,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to update contact", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toContactResponse(row)
	return &resp, nil
}
