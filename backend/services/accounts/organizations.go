package accounts

import (
	"context"
	"errors"
	"strconv"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// --- Request / Response types ---

const orgPageSize = 15

type OrganizationResponse struct {
	ID        int32    `json:"id"`
	Name      string   `json:"name"`
	IsOrganic bool     `json:"isOrganic"`
	Phone     *string  `json:"phone"`
	Address   *string  `json:"address"`
	Obligo    *float64 `json:"obligo"`
}

type ListOrganizationsRow struct {
	ID           int32    `json:"id"`
	Name         string   `json:"name"`
	IsOrganic    bool     `json:"isOrganic"`
	Phone        *string  `json:"phone"`
	Address      *string  `json:"address"`
	Obligo       *float64 `json:"obligo"`
	OfficeCount  int64    `json:"officeCount"`
	ContactCount int64    `json:"contactCount"`
	AgentCount   int64    `json:"agentCount"`
}

type ListOrganizationsResponse struct {
	Organizations []ListOrganizationsRow `json:"organizations"`
	Total         int64                  `json:"total"`
}

type ListOrganizationsRequest struct {
	Search    string `query:"search" encore:"optional"`
	IsOrganic string `query:"isOrganic" encore:"optional"`
	Page      int64  `query:"page" validate:"required,gte=1"`
}

func (p ListOrganizationsRequest) Validate() error {
	_, err := strconv.ParseBool(p.IsOrganic)
	if err != nil && p.IsOrganic != "" {
		return api_errors.NewValidationError("isOrganic is invalid")
	}
	return validation.ValidateStruct(p)
}

type CreateOrganizationRequest struct {
	Name      string   `json:"name" validate:"required,notblank"`
	IsOrganic bool     `json:"isOrganic"`
	Phone     *string  `json:"phone" validate:"omitempty,notblank" encore:"optional"`
	Address   *string  `json:"address" validate:"omitempty,notblank" encore:"optional"`
	Obligo    *float64 `json:"obligo" validate:"omitempty,gte=0" encore:"optional"`
}

func (p CreateOrganizationRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type UpdateOrganizationRequest struct {
	Name      *string  `json:"name" validate:"omitempty,notblank" encore:"optional"`
	IsOrganic *bool    `json:"isOrganic" encore:"optional"`
	Phone     *string  `json:"phone" encore:"optional"`
	Address   *string  `json:"address" encore:"optional"`
	Obligo    *float64 `json:"obligo" validate:"omitempty,gte=0" encore:"optional"`
}

func (p UpdateOrganizationRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// --- Helpers ---

func toOrganizationResponse(o db.Organization) OrganizationResponse {
	return OrganizationResponse{
		ID:        o.ID,
		Name:      o.Name,
		IsOrganic: o.IsOrganic,
		Phone:     o.Phone,
		Address:   o.Address,
		Obligo:    db.FloatFromNumeric(o.Obligo),
	}
}

func toListOrganizationsRow(o db.ListOrganizationsRow) ListOrganizationsRow {
	return ListOrganizationsRow{
		ID:           o.ID,
		Name:         o.Name,
		IsOrganic:    o.IsOrganic,
		Phone:        o.Phone,
		Address:      o.Address,
		Obligo:       db.FloatFromNumeric(o.Obligo),
		OfficeCount:  o.OfficeCount,
		ContactCount: o.ContactCount,
		AgentCount:   o.AgentCount,
	}
}

// --- Endpoints ---

// ListOrganizations lists organizations with optional search and pagination.
//
//encore:api auth method=GET path=/organizations tag:admin
func (s *Service) ListOrganizations(ctx context.Context, params *ListOrganizationsRequest) (*ListOrganizationsResponse, error) {
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

// CreateOrganization creates a new organization.
//
//encore:api auth method=POST path=/organizations tag:admin
func (s *Service) CreateOrganization(ctx context.Context, params CreateOrganizationRequest) (*OrganizationResponse, error) {
	row, err := s.query.CreateOrganization(ctx, db.CreateOrganizationParams{
		Name:      params.Name,
		IsOrganic: params.IsOrganic,
		Phone:     params.Phone,
		Address:   params.Address,
		Obligo:    db.NumericParam(params.Obligo),
	})
	if err != nil {
		if db.IsUniqueViolation(err) {
			return nil, ErrNameAlreadyExists
		}
		rlog.Error("failed to create organization", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toOrganizationResponse(row)
	return &resp, nil
}

// UpdateOrganization updates an existing organization.
//
//encore:api auth method=PUT path=/organizations/:id tag:admin
func (s *Service) UpdateOrganization(ctx context.Context, id int32, params UpdateOrganizationRequest) (*OrganizationResponse, error) {
	row, err := s.query.UpdateOrganization(ctx, db.UpdateOrganizationParams{
		ID:        id,
		Name:      params.Name,
		IsOrganic: params.IsOrganic,
		Phone:     params.Phone,
		Address:   params.Address,
		Obligo:    db.NumericParam(params.Obligo),
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
