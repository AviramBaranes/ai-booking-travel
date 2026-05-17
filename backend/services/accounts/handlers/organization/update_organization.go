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
	Name           *string `json:"name" validate:"omitempty,notblank" encore:"optional"`
	IsOrganic      *bool   `json:"isOrganic" encore:"optional"`
	IcountClientID *int32  `json:"icountClientId" validate:"omitempty,gt=0" encore:"optional"`
	Phone          *string `json:"phone" encore:"optional"`
	Address        *string `json:"address" encore:"optional"`
	Obligo         *int32  `json:"obligo" validate:"omitempty,gt=0" encore:"optional"`
}

func (p UpdateOrganizationParams) Validate() error {
	// If both fields are provided we can validate the constraint immediately.
	if p.IsOrganic != nil && p.IcountClientID != nil {
		if err := validateIcountClientIDConstraint(*p.IsOrganic, p.IcountClientID); err != nil {
			return err
		}
	}
	return validation.ValidateStruct(p)
}

// validateUpdateIcountClientIDConstraint handles the case where only one of
// isOrganic / icountClientId is being changed. It fetches the current values
// from the DB, merges the incoming partial update, and re-validates.
func (s *OrganizationService) validateUpdateIcountClientIDConstraint(ctx context.Context, id int64, params UpdateOrganizationParams) error {
	current, err := s.query.GetOrganizationBillingState(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to fetch organization billing state", "error", err)
		return api_errors.ErrInternalError
	}

	finalIsOrganic := current.IsOrganic
	if params.IsOrganic != nil {
		finalIsOrganic = *params.IsOrganic
	}

	finalIcount := current.IcountClientID
	if params.IcountClientID != nil {
		finalIcount = params.IcountClientID
	}

	return validateIcountClientIDConstraint(finalIsOrganic, finalIcount)
}

// UpdateOrganization updates an existing organization.
func (s *OrganizationService) UpdateOrganization(ctx context.Context, id int64, params UpdateOrganizationParams) (*OrganizationResponse, error) {
	// When only one of IsOrganic / IcountClientID is provided we need to
	// resolve the effective pair from the DB before validating.
	if (params.IsOrganic != nil) != (params.IcountClientID != nil) {
		if err := s.validateUpdateIcountClientIDConstraint(ctx, id, params); err != nil {
			return nil, err
		}
	}

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
