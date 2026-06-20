package office

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type UpdateOfficeParams struct {
	Name           string   `json:"name" validate:"required,notblank"`
	OrganizationID int64    `json:"organizationId" validate:"required,gte=1"`
	IcountClientID *int32   `json:"icountClientId" validate:"omitempty,gte=0" encore:"optional"`
	Phone          *string  `json:"phone" encore:"optional"`
	Address        *string  `json:"address" encore:"optional"`
	Obligo         *int32   `json:"obligo" validate:"omitempty,gte=0" encore:"optional"`
	GrossMarkup    *float64 `json:"grossMarkup" validate:"omitempty,gte=0" encore:"optional"`
}

func (p UpdateOfficeParams) Validate() error {
	return validation.ValidateStruct(p)
}

// validateUpdateOfficeIcountClientIDConstraint fetches the office's current billing state and
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
	if params.OrganizationID != billing.OrganizationID {
		newOrg, err := s.query.GetOrganizationBillingState(ctx, params.OrganizationID)
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
	if err := s.validateUpdateOfficeIcountClientIDConstraint(ctx, id, params); err != nil {
		return nil, err
	}

	row, err := s.query.UpdateOffice(ctx, db.UpdateOfficeParams{
		ID:             id,
		Name:           params.Name,
		OrganizationID: params.OrganizationID,
		IcountClientID: params.IcountClientID,
		Phone:          params.Phone,
		Address:        params.Address,
		Obligo:         params.Obligo,
		GrossMarkup:    dbadapters.NumericFromFloat64Ptr(params.GrossMarkup),
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
		Obligo:         row.Obligo,
		GrossMarkup:    dbadapters.NumericToFloat64Ptr(row.GrossMarkup),
	}
	return &resp, nil
}
