package organization_handlers

import (
	"context"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

const orgPageSize = 15

type ListOrganizationsRow struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	IsOrganic      bool    `json:"isOrganic"`
	IcountClientID *int32  `json:"icountClientId" encore:"optional"`
	Phone          *string `json:"phone" encore:"optional"`
	Address        *string `json:"address" encore:"optional"`
	Obligo         *int32  `json:"obligo" encore:"optional"`
	OfficeCount    int64   `json:"officeCount"`
	ContactCount   int64   `json:"contactCount"`
	AgentCount     int64   `json:"agentCount"`
}

type ListOrganizationsResponse struct {
	Organizations []ListOrganizationsRow `json:"organizations"`
	Total         int64                  `json:"total"`
}

type ListOrganizationsParams struct {
	Search    string `query:"search" encore:"optional"`
	IsOrganic string `query:"isOrganic" encore:"optional"`
	Page      int64  `query:"page" validate:"required,gte=1"`
}

func (p ListOrganizationsParams) Validate() error {
	_, err := strconv.ParseBool(p.IsOrganic)
	if err != nil && p.IsOrganic != "" {
		return api_errors.NewValidationError("isOrganic is invalid")
	}
	return validation.ValidateStruct(p)
}

func toListOrganizationsRow(o db.ListOrganizationsRow) ListOrganizationsRow {
	return ListOrganizationsRow{
		ID:             o.ID,
		Name:           o.Name,
		IsOrganic:      o.IsOrganic,
		IcountClientID: o.IcountClientID,
		Phone:          o.Phone,
		Address:        o.Address,
		Obligo:         o.Obligo,
		OfficeCount:    o.OfficeCount,
		ContactCount:   o.ContactCount,
		AgentCount:     o.AgentCount,
	}
}

// ListOrganizations lists organizations with optional search and pagination.
func (s *OrganizationService) ListOrganizations(ctx context.Context, params ListOrganizationsParams) (*ListOrganizationsResponse, error) {
	offset := (params.Page - 1) * orgPageSize

	var searchPtr *string
	if params.Search != "" {
		searchPtr = &params.Search
	}

	var isOrganicPtr *bool
	if params.IsOrganic != "" {
		isOrganic, _ := strconv.ParseBool(params.IsOrganic)
		isOrganicPtr = &isOrganic
	}

	rows, err := s.query.ListOrganizations(ctx, db.ListOrganizationsParams{
		Name:       searchPtr,
		IsOrganic:  isOrganicPtr,
		PageSize:   orgPageSize,
		PageOffset: offset,
	})
	if err != nil {
		rlog.Error("failed to list organizations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	total, err := s.query.CountOrganizations(ctx, db.CountOrganizationsParams{
		Name:      searchPtr,
		IsOrganic: isOrganicPtr,
	})
	if err != nil {
		rlog.Error("failed to count organizations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	orgs := make([]ListOrganizationsRow, 0, len(rows))
	for _, r := range rows {
		orgs = append(orgs, toListOrganizationsRow(r))
	}

	return &ListOrganizationsResponse{Organizations: orgs, Total: total}, nil
}
