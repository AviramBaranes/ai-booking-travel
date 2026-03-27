package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// DeleteLocation deletes a location broker code by its ID.
// If no other broker codes reference the same location, the location is also deleted.
//
//encore:api auth method=DELETE path=/locations/:id tag:admin
func (s *Service) DeleteLocation(ctx context.Context, id int64) error {
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
