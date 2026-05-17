package office

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

const officesPageSize int64 = 15

type ListOfficesParams struct {
	Search string `query:"search" encore:"optional"`
	OrgID  int64  `query:"orgId" encore:"optional"`
	Page   int32  `query:"page" validate:"required,gte=1"`
}

func (p ListOfficesParams) Validate() error {
	return validation.ValidateStruct(p)
}

type ListOfficesResponse struct {
	Offices []OfficeResponse `json:"offices"`
	Total   int64            `json:"total"`
}

// ListOffices lists offices with optional filtering and pagination.
func (s *OfficeService) ListOffices(ctx context.Context, params ListOfficesParams) (*ListOfficesResponse, error) {
	offset := int64(params.Page-1) * officesPageSize

	var searchPtr *string
	if params.Search != "" {
		searchPtr = &params.Search
	}

	var orgIDPtr *int64
	if params.OrgID != 0 {
		orgIDPtr = &params.OrgID
	}

	rows, err := s.query.ListOffices(ctx, db.ListOfficesParams{
		Name:           searchPtr,
		OrganizationID: orgIDPtr,
		PageOffset:     offset,
		PageSize:       officesPageSize,
	})
	if err != nil {
		rlog.Error("failed to list offices", "error", err)
		return nil, api_errors.ErrInternalError
	}

	total, err := s.query.CountOffices(ctx, db.CountOfficesParams{
		Name:           searchPtr,
		OrganizationID: orgIDPtr,
	})
	if err != nil {
		rlog.Error("failed to count offices", "error", err)
		return nil, api_errors.ErrInternalError
	}

	offices := make([]OfficeResponse, 0, len(rows))
	for _, r := range rows {
		offices = append(offices, toOfficeResponse(r))
	}

	return &ListOfficesResponse{Offices: offices, Total: total}, nil
}
