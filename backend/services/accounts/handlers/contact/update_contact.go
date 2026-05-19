package contact

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

type UpdateContactParams struct {
	FirstName            string `json:"firstName" validate:"required,notblank"`
	LastName             string `json:"lastName" validate:"required,notblank"`
	Role                 string `json:"role" validate:"required,notblank"`
	Cellphone            string `json:"cellphone" validate:"required,notblank"`
	Email                string `json:"email" validate:"required,email"`
	OfficeID             *int64 `json:"officeId" encore:"optional"`
	OrganizationID       *int64 `json:"organizationId" encore:"optional"`
	IsPaymentResponsible bool   `json:"isPaymentResponsible" encore:"optional"`
}

func (p UpdateContactParams) Validate() error {
	if err := validation.ValidateStruct(p); err != nil {
		return err
	}

	hasOffice := p.OfficeID != nil
	hasOrg := p.OrganizationID != nil
	if hasOffice == hasOrg {
		return api_errors.NewErrorWithDetail(
			errs.InvalidArgument,
			"Exactly one of officeId or organizationId must be provided",
			api_errors.ErrorDetails{Code: api_errors.CodeInvalidValue},
		)
	}

	return nil
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
