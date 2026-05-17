package location

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

func (s *LocationService) DeleteLocation(ctx context.Context, id int64) error {
	locationID, err := s.query.DeleteLocationBrokerCode(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to delete location broker code", "error", err, "id", id)
		return api_errors.ErrInternalError
	}

	count, err := s.query.CountLocationBrokerCodesByLocationID(ctx, locationID)
	if err != nil {
		rlog.Error("failed to count remaining broker codes", "error", err, "location_id", locationID)
		return api_errors.ErrInternalError
	}

	if count == 0 {
		if err := s.query.DeleteLocationByID(ctx, locationID); err != nil {
			rlog.Error("failed to delete orphaned location", "error", err, "location_id", locationID)
			return api_errors.ErrInternalError
		}
	}

	return nil
}
