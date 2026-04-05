package accounts

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// --- Request / Response types ---

type OfficeResponse struct {
	ID             int32   `json:"id"`
	Name           string  `json:"name"`
	OrganizationID int32   `json:"organizationId"`
	Phone          *string `json:"phone"`
	Address        *string `json:"address"`
	ContactCount   int64   `json:"contactCount"`
	AgentCount     int64   `json:"agentCount"`
}

type ListOfficesRequest struct {
	Search string `query:"search"`
	OrgID  int32  `query:"orgId"`
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
	OrganizationID int32   `json:"organizationId" validate:"required,gte=1"`
	Phone          *string `json:"phone"`
	Address        *string `json:"address"`
}

func (p CreateOfficeRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type UpdateOfficeRequest struct {
	Name           *string `json:"name" validate:"omitempty,notblank"`
	OrganizationID *int32  `json:"organizationId" validate:"omitempty,gte=1"`
	Phone          *string `json:"phone"`
	Address        *string `json:"address"`
}

func (p UpdateOfficeRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// --- Helpers ---

const officesPageSize int64 = 15

func toOfficeResponse(o db.ListOfficesRow) OfficeResponse {
	return OfficeResponse{
		ID:             o.ID,
		Name:           o.Name,
		OrganizationID: o.OrganizationID,
		Phone:          o.Phone,
		Address:        o.Address,
		ContactCount:   o.ContactCount,
		AgentCount:     o.AgentCount,
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

	var orgIDPtr *int32
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
	row, err := s.query.CreateOffice(ctx, db.CreateOfficeParams{
		Name:           params.Name,
		OrganizationID: params.OrganizationID,
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
		Phone:          row.Phone,
		Address:        row.Address,
	}
	return &resp, nil
}

// UpdateOffice updates an existing office.
//
//encore:api auth method=PUT path=/offices/:id tag:admin
func (s *Service) UpdateOffice(ctx context.Context, id int32, params UpdateOfficeRequest) (*OfficeResponse, error) {
	row, err := s.query.UpdateOffice(ctx, db.UpdateOfficeParams{
		ID:             id,
		Name:           params.Name,
		OrganizationID: params.OrganizationID,
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
		Phone:          row.Phone,
		Address:        row.Address,
	}
	return &resp, nil
}
