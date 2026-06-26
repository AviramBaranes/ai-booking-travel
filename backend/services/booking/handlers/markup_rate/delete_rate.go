package markup_rate

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

func (s *MarkupRateService) DeleteMarkupRate(ctx context.Context, id int64) error {
	_, err := s.query.DeleteMarkupRate(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to delete markup rate", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
