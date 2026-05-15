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
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	IsOrganic      bool    `json:"isOrganic"`
	IcountClientID *int32  `json:"icountClientId"`
	Phone          *string `json:"phone"`
	Address        *string `json:"address"`
	Obligo         *int32  `json:"obligo"`
}

type ListOrganizationsRow struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	IsOrganic      bool    `json:"isOrganic"`
	IcountClientID *int32  `json:"icountClientId"`
	Phone          *string `json:"phone"`
	Address        *string `json:"address"`
	Obligo         *int32  `json:"obligo"`
	OfficeCount    int64   `json:"officeCount"`
	ContactCount   int64   `json:"contactCount"`
	AgentCount     int64   `json:"agentCount"`
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
	Name           string  `json:"name" validate:"required,notblank"`
	IsOrganic      bool    `json:"isOrganic"`
	IcountClientID *int32  `json:"icountClientId" validate:"omitempty,gt=0" encore:"optional"`
	Phone          *string `json:"phone" validate:"omitempty,notblank" encore:"optional"`
	Address        *string `json:"address" validate:"omitempty,notblank" encore:"optional"`
	Obligo         *int32  `json:"obligo" validate:"omitempty,gt=0" encore:"optional"`
}

func (p CreateOrganizationRequest) Validate() error {
	if err := validateIcountClientIDConstraint(p.IsOrganic, p.IcountClientID); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

type UpdateOrganizationRequest struct {
	Name           *string `json:"name" validate:"omitempty,notblank" encore:"optional"`
	IsOrganic      *bool   `json:"isOrganic" encore:"optional"`
	IcountClientID *int32  `json:"icountClientId" validate:"omitempty,gt=0" encore:"optional"`
	Phone          *string `json:"phone" encore:"optional"`
	Address        *string `json:"address" encore:"optional"`
	Obligo         *int32  `json:"obligo" validate:"omitempty,gt=0" encore:"optional"`
}

func (p UpdateOrganizationRequest) Validate() error {
	// If both fields are provided we can validate the constraint immediately.
	if p.IsOrganic != nil && p.IcountClientID != nil {
		if err := validateIcountClientIDConstraint(*p.IsOrganic, p.IcountClientID); err != nil {
			return err
		}
	}
	return validation.ValidateStruct(p)
}

// --- Helpers ---

// validateIcountClientIDConstraint enforces billing rules:
// organic orgs must have an icount_client_id; non-organic orgs must not.
func validateIcountClientIDConstraint(isOrganic bool, icountClientID *int32) error {
	if isOrganic && icountClientID == nil {
		return ErrOrganizationOrganicRequiresIcountClientID
	}
	if !isOrganic && icountClientID != nil {
		return ErrOrganizationNonOrganicForbidsIcountClientID
	}
	return nil
}

// validateUpdateIcountClientIDConstraint handles the case where only one of
// isOrganic / icountClientId is being changed. It fetches the current values
// from the DB, merges the incoming partial update, and re-validates.
func (s *Service) validateUpdateIcountClientIDConstraint(ctx context.Context, id int64, params UpdateOrganizationRequest) error {
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

func toOrganizationResponse(o db.Organization) OrganizationResponse {
	return OrganizationResponse{
		ID:             o.ID,
		Name:           o.Name,
		IsOrganic:      o.IsOrganic,
		IcountClientID: o.IcountClientID,
		Phone:          o.Phone,
		Address:        o.Address,
		Obligo:         o.Obligo,
	}
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
		Name:           params.Name,
		IsOrganic:      params.IsOrganic,
		IcountClientID: params.IcountClientID,
		Phone:          params.Phone,
		Address:        params.Address,
		Obligo:         params.Obligo,
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
func (s *Service) UpdateOrganization(ctx context.Context, id int64, params UpdateOrganizationRequest) (*OrganizationResponse, error) {
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

type OrganicOrganization struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ListOrganicOrganizationResponse struct {
	Organizations []OrganicOrganization `json:"organizations"`
}

// ListOrganicOrganizations lists all organic organizations for accountant use.
//
//encore:api auth method=GET path=/organic-organizations tag:accountant
func (s *Service) ListOrganicOrganizations(ctx context.Context) (*ListOrganicOrganizationResponse, error) {
	rows, err := s.query.ListOrganicOrganizations(ctx)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return &ListOrganicOrganizationResponse{Organizations: []OrganicOrganization{}}, nil
		}
		rlog.Error("failed to list organic organizations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	orgs := make([]OrganicOrganization, 0, len(rows))
	for _, r := range rows {
		orgs = append(orgs, OrganicOrganization{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return &ListOrganicOrganizationResponse{Organizations: orgs}, nil
}
