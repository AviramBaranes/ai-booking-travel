package location

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type ToggleLocationParams struct {
	Enabled bool `json:"enabled"`
}

func (s *LocationService) ToggleLocation(ctx context.Context, id int64, p ToggleLocationParams) error {
	if err := s.query.ToggleLocationBrokerCode(ctx, db.ToggleLocationBrokerCodeParams{
		Enabled: p.Enabled,
		ID:      id,
	}); err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to toggle location broker code", "error", err, "id", id, "enabled", p.Enabled)
		return api_errors.ErrInternalError
	}

	return nil
}
