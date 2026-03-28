package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// VerifyBrokerTranslation marks a broker translation as verified by ID.
//
//encore:api auth method=PATCH path=/broker-translations/:id/verify tag:admin
func (s *Service) VerifyBrokerTranslation(ctx context.Context, id int32) error {
	err := s.query.VerifyBrokerTranslation(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to verify broker translation", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
