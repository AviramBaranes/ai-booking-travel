package price_offer

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// ApprovePriceOffer changes the status of a price offer to "approved" and notify the agent.
func (s *PriceOfferService) ApprovePriceOffer(ctx context.Context, id int64) error {
	_, err := s.query.ApprovePriceOffer(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to approve price offer", "error", err, "price_offer_id", id)
		return api_errors.ErrInternalError
	}

	return nil
}
