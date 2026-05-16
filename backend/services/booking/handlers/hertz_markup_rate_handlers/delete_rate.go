package hertz_markup_rate_handlers

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

func (s *HertzMarkupRateService) DeleteHertzMarkupRate(ctx context.Context, id int64) error {
	_, err := s.query.DeleteHertzMarkupRate(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to delete hertz markup rate", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
