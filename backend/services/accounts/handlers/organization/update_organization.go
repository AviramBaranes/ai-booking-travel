package organization

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type UpdateOrganizationParams struct {
	Name           string  `json:"name" validate:"required,notblank"`
	IsOrganic      bool    `json:"isOrganic"`
	IcountClientID *int32  `json:"icountClientId" validate:"omitempty,gt=0" encore:"optional"`
	Phone          *string `json:"phone" validate:"omitempty,notblank" encore:"optional"`
	Address        *string `json:"address" validate:"omitempty,notblank" encore:"optional"`
	Obligo         *int32  `json:"obligo" validate:"omitempty,gt=0" encore:"optional"`
}

func (p UpdateOrganizationParams) Validate() error {
	if err := validateIcountClientIDConstraint(p.IsOrganic, p.IcountClientID); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

// UpdateOrganization updates an existing organization.
func (s *OrganizationService) UpdateOrganization(ctx context.Context, id int64, params UpdateOrganizationParams) (*OrganizationResponse, error) {
	row, err := s.query.UpdateOrganization(ctx, db.UpdateOrganizationParams{
		ID:             id,
		Name:           params.Name,
		IsOrganic:      params.IsOrganic,
		IcountClientID: params.IcountClientID,
		Phone:          params.Phone,
		Address:        params.Address,
		Obligo:         params.Obligo,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		if db.IsUniqueViolation(err) {
			return nil, ErrNameAlreadyExists
		}
		rlog.Error("failed to update organization", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toOrganizationResponse(row)
	return &resp, nil
}
