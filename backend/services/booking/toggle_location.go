package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type ToggleLocationRequest struct {
	Enabled bool `json:"enabled"`
}

//encore:api auth method=PATCH path=/locations/:id tag:admin
func (s *Service) ToggleLocation(ctx context.Context, id int64, p *ToggleLocationRequest) error {
	var err error
	if p.Enabled {
		err = s.query.EnableLocationBrokerCode(ctx, id)
	} else {
		err = s.query.DisableLocationBrokerCode(ctx, id)
	}

	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to toggle location broker code", "error", err, "id", id, "enabled", p.Enabled)
		return api_errors.ErrInternalError
	}

	return nil
}
