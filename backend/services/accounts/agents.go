package accounts

import (
	"context"
	"errors"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// --- Request / Response types ---

type AgentResponse struct {
	ID               int32      `json:"id"`
	Email            string     `json:"email"`
	PhoneNumber      *string    `json:"phoneNumber"`
	OfficeID         *int32     `json:"officeId"`
	OfficeName       *string    `json:"officeName"`
	OrganizationName *string    `json:"organizationName"`
	LastLogin        *time.Time `json:"lastLogin"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type ListAgentsRequest struct {
	Search   string `query:"search"`
	OfficeID int32  `query:"officeId"`
	OrgID    int32  `query:"orgId"`
	Page     int32  `query:"page" validate:"required,gte=1"`
}

func (p ListAgentsRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type ListAgentsResponse struct {
	Agents []AgentResponse `json:"agents"`
	Total  int64           `json:"total"`
}

type CreateAgentRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8" encore:"sensitive"`
	PhoneNumber string `json:"phoneNumber" validate:"required,israeli_phone"`
	OfficeID    int32  `json:"officeId" validate:"required,gte=1"`
}

func (p CreateAgentRequest) Validate() error {
	if err := validatePasswordForAPI(p.Password); err != nil {
		return err
	}
	return validation.ValidateStruct(p)
}

type CreateAgentResponse struct {
	ID int32 `json:"id"`
}

// --- Helpers ---

const agentsPageSize int32 = 15

func toAgentResponse(r db.ListAgentsRow) AgentResponse {
	return AgentResponse{
		ID:               r.ID,
		Email:            r.Email,
		PhoneNumber:      r.PhoneNumber,
		OfficeID:         r.OfficeID,
		OfficeName:       r.OfficeName,
		OrganizationName: r.OrganizationName,
		LastLogin:        db.TimePtrFromDB(r.LastLogin),
		CreatedAt:        db.TimeFromDB(r.CreatedAt),
		UpdatedAt:        db.TimeFromDB(r.UpdatedAt),
	}
}

// --- Endpoints ---

// ListAgents lists agents with optional filtering and pagination.
//
//encore:api auth method=GET path=/agents tag:admin
func (s *Service) ListAgents(ctx context.Context, params *ListAgentsRequest) (*ListAgentsResponse, error) {
	offset := (params.Page - 1) * agentsPageSize

	var searchPtr *string
	if params.Search != "" {
		searchPtr = &params.Search
	}

	var officeIDPtr *int32
	if params.OfficeID != 0 {
		officeIDPtr = &params.OfficeID
	}

	var orgIDPtr *int32
	if params.OrgID != 0 {
		orgIDPtr = &params.OrgID
	}

	rows, err := s.query.ListAgents(ctx, db.ListAgentsParams{
		Search:         searchPtr,
		OfficeID:       officeIDPtr,
		OrganizationID: orgIDPtr,
		PageOffset:     offset,
		PageSize:       agentsPageSize,
	})
	if err != nil {
		rlog.Error("failed to list agents", "error", err)
		return nil, api_errors.ErrInternalError
	}

	total, err := s.query.CountAgents(ctx, db.CountAgentsParams{
		Search:         searchPtr,
		OfficeID:       officeIDPtr,
		OrganizationID: orgIDPtr,
	})
	if err != nil {
		rlog.Error("failed to count agents", "error", err)
		return nil, api_errors.ErrInternalError
	}

	agents := make([]AgentResponse, 0, len(rows))
	for _, r := range rows {
		agents = append(agents, toAgentResponse(r))
	}

	return &ListAgentsResponse{Agents: agents, Total: total}, nil
}

// CreateAgent creates a new agent user.
//
//encore:api auth method=POST path=/agents tag:admin
func (s *Service) CreateAgent(ctx context.Context, params CreateAgentRequest) (*CreateAgentResponse, error) {
	userID, err := s.query.CheckUserExists(ctx, params.Email)
	if err != nil && !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to check if user exists", "email", params.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}
	if userID != 0 {
		return nil, ErrEmailAlreadyExists
	}

	hashed, err := password.HashPassword(params.Password)
	if err != nil {
		rlog.Error("failed to hash password", "email", params.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	row, err := s.query.CreateAgent(ctx, db.CreateAgentParams{
		Email:        params.Email,
		PhoneNumber:  &params.PhoneNumber,
		PasswordHash: hashed,
		OfficeID:     &params.OfficeID,
	})
	if err != nil {
		rlog.Error("failed to create agent user", "email", params.Email, "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &CreateAgentResponse{
		ID: row.ID,
	}, nil
}
