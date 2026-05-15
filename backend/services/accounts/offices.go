package accounts

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// validateOfficeIcountClientIDConstraint enforces billing rules for offices.
// Offices under organic orgs must NOT have an icount_client_id;
// offices under non-organic orgs MUST have one.
func validateOfficeIcountClientIDConstraint(isOrganic bool, icountClientID *int32) error {
	if isOrganic && icountClientID != nil {
		return ErrOfficeOrganicForbidsIcountClientID
	}
	if !isOrganic && icountClientID == nil {
		return ErrOfficeNonOrganicRequiresIcountClientID
	}
	return nil
}

// validateCreateOfficeIcountClientIDConstraint always runs on create.
// It fetches the parent organization's billing state to determine the constraint.
func (s *Service) validateCreateOfficeIcountClientIDConstraint(ctx context.Context, orgID int64, icountClientID *int32) error {
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

// validateUpdateOfficeIcountClientIDConstraint runs only when IcountClientID is
// present in the update payload. Fetches the office's current billing state and
// resolves the final organization if it is being changed.
func (s *Service) validateUpdateOfficeIcountClientIDConstraint(ctx context.Context, id int64, params UpdateOfficeRequest) error {
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

// --- Request / Response types ---

type OfficeResponse struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	OrganizationID   int64   `json:"organizationId"`
	OrganizationName string  `json:"organizationName"`
	IcountClientID   *int32  `json:"icountClientId"`
	Phone            *string `json:"phone"`
	Address          *string `json:"address"`
	ContactCount     int64   `json:"contactCount"`
	AgentCount       int64   `json:"agentCount"`
}

type ListOfficesRequest struct {
	Search string `query:"search"`
	OrgID  int64  `query:"orgId"`
	Page   int32  `query:"page" validate:"required,gte=1"`
}

func (p ListOfficesRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type ListOfficesResponse struct {
	Offices []OfficeResponse `json:"offices"`
	Total   int64            `json:"total"`
}

type CreateOfficeRequest struct {
	Name           string  `json:"name" validate:"required,notblank"`
	OrganizationID int64   `json:"organizationId" validate:"required,gte=1"`
	IcountClientID *int32  `json:"icountClientId" validate:"omitempty,gte=0" encore:"optional"`
	Phone          *string `json:"phone" encore:"optional"`
	Address        *string `json:"address" encore:"optional"`
}

func (p CreateOfficeRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type UpdateOfficeRequest struct {
	Name           *string `json:"name" validate:"omitempty,notblank" encore:"optional"`
	OrganizationID *int64  `json:"organizationId" validate:"omitempty,gte=1" encore:"optional"`
	IcountClientID *int32  `json:"icountClientId" validate:"omitempty,gte=0" encore:"optional"`
	Phone          *string `json:"phone" encore:"optional"`
	Address        *string `json:"address" encore:"optional"`
}

func (p UpdateOfficeRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// --- Helpers ---

const officesPageSize int64 = 15

func toOfficeResponse(o db.ListOfficesRow) OfficeResponse {
	return OfficeResponse{
		ID:               o.ID,
		Name:             o.Name,
		OrganizationID:   o.OrganizationID,
		OrganizationName: o.OrganizationName,
		IcountClientID:   o.IcountClientID,
		Phone:            o.Phone,
		Address:          o.Address,
		ContactCount:     o.ContactCount,
		AgentCount:       o.AgentCount,
	}
}

// --- Endpoints ---

// ListOffices lists offices with optional filtering and pagination.
//
//encore:api auth method=GET path=/offices tag:admin
func (s *Service) ListOffices(ctx context.Context, params *ListOfficesRequest) (*ListOfficesResponse, error) {
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

// CreateOffice creates a new office.
//
//encore:api auth method=POST path=/offices tag:admin
func (s *Service) CreateOffice(ctx context.Context, params CreateOfficeRequest) (*OfficeResponse, error) {
	if err := s.validateCreateOfficeIcountClientIDConstraint(ctx, params.OrganizationID, params.IcountClientID); err != nil {
		return nil, err
	}

	row, err := s.query.CreateOffice(ctx, db.CreateOfficeParams{
		Name:           params.Name,
		OrganizationID: params.OrganizationID,
		IcountClientID: params.IcountClientID,
		Phone:          params.Phone,
		Address:        params.Address,
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
	}
	return &resp, nil
}

// UpdateOffice updates an existing office.
//
//encore:api auth method=PUT path=/offices/:id tag:admin
func (s *Service) UpdateOffice(ctx context.Context, id int64, params UpdateOfficeRequest) (*OfficeResponse, error) {
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

type InorganicOffice struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ListInorganicOfficeResponse struct {
	Offices []InorganicOffice `json:"offices"`
}

// ListInorganicOffices lists all inorganic offices for accountant use.
//
//encore:api auth method=GET path=/inorganic-offices tag:accountant
func (s *Service) ListInorganicOffices(ctx context.Context) (*ListInorganicOfficeResponse, error) {
	rows, err := s.query.ListInorganicOffices(ctx)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return &ListInorganicOfficeResponse{Offices: []InorganicOffice{}}, nil
		}
		rlog.Error("failed to list inorganic offices", "error", err)
		return nil, api_errors.ErrInternalError
	}

	offices := make([]InorganicOffice, 0, len(rows))
	for _, r := range rows {
		offices = append(offices, InorganicOffice{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return &ListInorganicOfficeResponse{Offices: offices}, nil
}
