package price_offer

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	emailevents "encore.app/services/notifications/events"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

var emailRequestedTopic = pubsub.TopicRef[pubsub.Publisher[*emailevents.EmailEvent]](emailevents.EmailRequestedTopic)
var emailPublisher = emailevents.NewPublisher(emailRequestedTopic)

// ApprovePriceOffer changes the status of a price offer to "approved" and notify the agent.
func (s *PriceOfferService) ApprovePriceOffer(ctx context.Context, id int64) error {
	priceOffer, err := s.query.ApprovePriceOffer(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to approve price offer", "error", err, "price_offer_id", id)
		return api_errors.ErrInternalError
	}

	if _, err := emailPublisher.Publish(ctx, emailevents.EmailEventPriceOfferApproved, emailevents.PriceOfferApprovedEmailPayload{
		PriceOfferID:   id,
		PriceOfferName: priceOffer.Name,
		AgentID:        priceOffer.AgentID,
		Price:          float64(priceOffer.OfferedPrice),
		Currency:       priceOffer.CurrencyCode,
	}); err != nil {
		rlog.Error("failed to publish price offer approved email event", "price_offer_id", id, "error", err)
	}

	return nil
}
