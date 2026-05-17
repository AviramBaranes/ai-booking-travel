package office

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type UpdateOfficeParams struct {
	Name           *string `json:"name" validate:"omitempty,notblank" encore:"optional"`
	OrganizationID *int64  `json:"organizationId" validate:"omitempty,gte=1" encore:"optional"`
	IcountClientID *int32  `json:"icountClientId" validate:"omitempty,gte=0" encore:"optional"`
	Phone          *string `json:"phone" encore:"optional"`
	Address        *string `json:"address" encore:"optional"`
}

func (p UpdateOfficeParams) Validate() error {
	return validation.ValidateStruct(p)
}

// validateUpdateOfficeIcountClientIDConstraint runs only when IcountClientID is
// present in the update payload. Fetches the office's current billing state and
// resolves the final organization if it is being changed.
func (s *OfficeService) validateUpdateOfficeIcountClientIDConstraint(ctx context.Context, id int64, params UpdateOfficeParams) error {
	billing, err := s.query.GetOfficeBillingState(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to fetch office billing state", "error", err)
		return api_errors.ErrInternalError
	}

	finalIsOrganic := billing.IsOrganic
	if params.OrganizationID != nil && *params.OrganizationID != billing.OrganizationID {
		newOrg, err := s.query.GetOrganizationBillingState(ctx, *params.OrganizationID)
		if err != nil {
			if errors.Is(err, db.ErrNoRows) {
				return api_errors.ErrNotFound
			}
			rlog.Error("failed to fetch organization billing state", "error", err)
			return api_errors.ErrInternalError
		}
		finalIsOrganic = newOrg.IsOrganic
	}

	return validateOfficeIcountClientIDConstraint(finalIsOrganic, params.IcountClientID)
}

// UpdateOffice updates an existing office.
func (s *OfficeService) UpdateOffice(ctx context.Context, id int64, params UpdateOfficeParams) (*OfficeResponse, error) {
	if params.IcountClientID != nil {
		if err := s.validateUpdateOfficeIcountClientIDConstraint(ctx, id, params); err != nil {
			return nil, err
		}
	}

	row, err := s.query.UpdateOffice(ctx, db.UpdateOfficeParams{
		ID:             id,
		Name:           params.Name,
		OrganizationID: params.OrganizationID,
		IcountClientID: params.IcountClientID,
		Phone:          params.Phone,
		Address:        params.Address,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		if db.IsUniqueViolation(err) {
			return nil, ErrNameAlreadyExists
		}
		rlog.Error("failed to update office", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := OfficeResponse{
		ID:             row.ID,
		Name:           row.Name,
		OrganizationID: row.OrganizationID,
		IcountClientID: row.IcountClientID,
		Phone:          row.Phone,
		Address:        row.Address,
	}
	return &resp, nil
}
