package user

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

// GetAgentsByOfficeIDParams are the params for retrieving agents by office ID.
type GetAgentsByOfficeIDParams struct {
	OfficeID int64
}

// GetAgentsByOfficeID retrieves agent IDs for a given office ID.
func (s *UserService) GetAgentsByOfficeID(ctx context.Context, params GetAgentsByOfficeIDParams) (*GetAgentsResponse, error) {
	rows, err := s.query.GetAgentsByOfficeID(ctx, params.OfficeID)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to get agents by office ID", "officeID", params.OfficeID, "error", err)
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
