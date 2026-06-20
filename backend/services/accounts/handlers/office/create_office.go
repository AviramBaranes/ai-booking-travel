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

type CreateOfficeParams struct {
	Name           string   `json:"name" validate:"required,notblank"`
	OrganizationID int64    `json:"organizationId" validate:"required,gte=1"`
	IcountClientID *int32   `json:"icountClientId" validate:"omitempty,gte=0" encore:"optional"`
	Phone          *string  `json:"phone" encore:"optional"`
	Address        *string  `json:"address" encore:"optional"`
	Obligo         *int32   `json:"obligo" validate:"omitempty,gte=0" encore:"optional"`
	GrossMarkup    *float64 `json:"grossMarkup" validate:"omitempty,gte=0" encore:"optional"`
}

func (p CreateOfficeParams) Validate() error {
	return validation.ValidateStruct(p)
}

// validateCreateOfficeIcountClientIDConstraint always runs on create.
// It fetches the parent organization's billing state to determine the constraint.
func (s *OfficeService) validateCreateOfficeIcountClientIDConstraint(ctx context.Context, orgID int64, icountClientID *int32) error {
	org, err := s.query.GetOrganizationBillingState(ctx, orgID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to fetch organization billing state", "error", err)
		return api_errors.ErrInternalError
	}
	return validateOfficeIcountClientIDConstraint(org.IsOrganic, icountClientID)
}

// CreateOffice creates a new office.
func (s *OfficeService) CreateOffice(ctx context.Context, params CreateOfficeParams) (*OfficeResponse, error) {
	if err := s.validateCreateOfficeIcountClientIDConstraint(ctx, params.OrganizationID, params.IcountClientID); err != nil {
		return nil, err
	}

	row, err := s.query.CreateOffice(ctx, db.CreateOfficeParams{
		Name:           params.Name,
		OrganizationID: params.OrganizationID,
		IcountClientID: params.IcountClientID,
		Phone:          params.Phone,
		Address:        params.Address,
		Obligo:         params.Obligo,
		GrossMarkup:    dbadapters.NumericFromFloat64Ptr(params.GrossMarkup),
	})
	if err != nil {
		if db.IsUniqueViolation(err) {
			return nil, ErrNameAlreadyExists
		}
		rlog.Error("failed to create office", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := OfficeResponse{
		ID:             row.ID,
		Name:           row.Name,
		OrganizationID: row.OrganizationID,
		IcountClientID: row.IcountClientID,
		Phone:          row.Phone,
		Address:        row.Address,
		GrossMarkup:    dbadapters.NumericToFloat64Ptr(row.GrossMarkup),
	}
	return &resp, nil
}
