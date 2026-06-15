package location

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type BulkToggleLocationsParams struct {
	IDs     []int64 `json:"ids" validate:"required,min=1,dive,min=1"`
	Enabled bool    `json:"enabled"`
}

func (p BulkToggleLocationsParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *LocationService) BulkToggleLocations(ctx context.Context, p BulkToggleLocationsParams) error {
	for _, id := range p.IDs {
		if err := s.query.ToggleLocationBrokerCode(ctx, db.ToggleLocationBrokerCodeParams{
			Enabled: p.Enabled,
			ID:      id,
		}); err != nil {
			rlog.Error("failed to toggle location broker code", "error", err, "id", id, "enabled", p.Enabled)
			return api_errors.ErrInternalError
		}
	}
	return nil
}
