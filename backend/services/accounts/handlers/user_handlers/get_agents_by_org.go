package user_handlers

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// GetAgentsByOrganizationIDParams are the params for retrieving agents by organization ID.
type GetAgentsByOrganizationIDParams struct {
	OrgID int64
}

// GetAgentsByOrganizationID retrieves agent IDs for a given organization ID.
func (s *UserService) GetAgentsByOrganizationID(ctx context.Context, params GetAgentsByOrganizationIDParams) (*GetAgentsResponse, error) {
	rows, err := s.query.GetAgentsByOrganizationID(ctx, params.OrgID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get agents by organization ID", "orgID", params.OrgID, "error", err)
		return nil, api_errors.ErrInternalError
	}

	if len(rows) == 0 {
		return nil, api_errors.ErrNotFound
	}

	ids := make([]int64, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}

	return &GetAgentsResponse{
		IDs:       ids,
		IsOrganic: rows[0].IsOrganic,
	}, nil
}
