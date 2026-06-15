package location

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type ToggleLocationIsAirportParams struct {
	IsAirport bool `json:"is_airport"`
}

func (s *LocationService) ToggleLocationIsAirport(ctx context.Context, id int64, p ToggleLocationIsAirportParams) error {
	if err := s.query.ToggleIsAirport(ctx, db.ToggleIsAirportParams{
		IsAirport: p.IsAirport,
		ID:        id,
	}); err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to toggle location airport status", "error", err, "id", id, "is_airport", p.IsAirport)
		return api_errors.ErrInternalError
	}

	return nil
}
