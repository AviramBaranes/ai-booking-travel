package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// DeleteBrokerTranslation deletes a broker translation by ID.
//
//encore:api auth method=DELETE path=/broker-translations/:id tag:admin
func (s *Service) DeleteBrokerTranslation(ctx context.Context, id int32) error {
	err := s.query.DeleteBrokerTranslation(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to delete broker translation", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
