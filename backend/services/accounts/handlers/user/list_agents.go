package user

import (
	"context"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// AgentResponse is the response type for a single agent.
type AgentResponse struct {
	ID               int64      `json:"id"`
	FirstName        string     `json:"firstName"`
	LastName         string     `json:"lastName"`
	Email            string     `json:"email"`
	PhoneNumber      *string    `json:"phoneNumber"`
	OfficeID         *int64     `json:"officeId"`
	OfficeName       *string    `json:"officeName"`
	OrganizationName *string    `json:"organizationName"`
	LastLogin        *time.Time `json:"lastLogin"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

// ListAgentsParams are the params for listing agents.
type ListAgentsParams struct {
	Search   string `query:"search"`
	OfficeID int64  `query:"officeId"`
	OrgID    int64  `query:"orgId"`
	Page     int32  `query:"page" validate:"required,gte=1"`
}

func (p ListAgentsParams) Validate() error {
	return validation.ValidateStruct(p)
}

// ListAgentsResponse is the response for listing agents.
type ListAgentsResponse struct {
	Agents []AgentResponse `json:"agents"`
	Total  int64           `json:"total"`
}

const agentsPageSize int32 = 15

func toAgentResponse(r db.ListAgentsRow) AgentResponse {
	lastLogin := dbadapters.TimeFromDB(r.LastLogin)
	return AgentResponse{
		ID:               r.ID,
		FirstName:        r.FirstName,
		LastName:         r.LastName,
		Email:            r.Email,
		PhoneNumber:      r.PhoneNumber,
		OfficeID:         r.OfficeID,
		OfficeName:       r.OfficeName,
		OrganizationName: r.OrganizationName,
		LastLogin:        &lastLogin,
		CreatedAt:        dbadapters.TimeFromDB(r.CreatedAt),
		UpdatedAt:        dbadapters.TimeFromDB(r.UpdatedAt),
	}
}

// ListAgents lists agents with optional filtering and pagination.
func (s *UserService) ListAgents(ctx context.Context, params *ListAgentsParams) (*ListAgentsResponse, error) {
	offset := (params.Page - 1) * agentsPageSize

	var searchPtr *string
	if params.Search != "" {
		searchPtr = &params.Search
	}

	var officeIDPtr *int64
	if params.OfficeID != 0 {
		officeIDPtr = &params.OfficeID
	}

	var orgIDPtr *int64
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
