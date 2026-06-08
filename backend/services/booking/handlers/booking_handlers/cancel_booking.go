package booking_handlers

import (
	"context"
	"fmt"

	"encore.app/services/booking/db"
	emailevents "encore.app/services/notifications/events"
	"encore.app/services/reservation"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

var emailRequestedTopic = pubsub.TopicRef[pubsub.Publisher[*emailevents.EmailEvent]](emailevents.EmailRequestedTopic)
var emailPublisher = emailevents.NewPublisher(emailRequestedTopic)

// CancelBooking handles the cancellation of a booking by processing the BookingCancellationEvent received from the reservation service.
func (s *BookingService) CancelBooking(ctx context.Context, e *reservation.BookingCancellationEvent) error {
	b, err := getCanceler(db.Broker(e.Broker))
	if err != nil {
		rlog.Error("unsupported broker for cancellation", "broker", b, "reservationId", e.ReservationID)
		if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
			Subject: "Unsupported broker for booking cancellation",
			Message: fmt.Sprintf("unsupported broker for cancellation: %v, reservationId: %v", b, e.ReservationID),
		}); publishErr != nil {
			rlog.Error("failed to publish critical error email event", "reservationId", e.ReservationID, "error", publishErr)
		}
		return err
	}

	err = b.Cancel(e.BrokerReservationID, e.LastName, e.SupplierCode)
	if err != nil {
		rlog.Error("failed to cancel booking", "broker", b, "reservationId", e.ReservationID, "error", err)
		if _, publishErr := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
			Subject: "Failed to cancel booking",
			Message: fmt.Sprintf("failed to cancel booking: %s, bookingID: %v, reservationId: %v, error: %v", b.Name(), e.BrokerReservationID, e.ReservationID, err),
		}); publishErr != nil {
			rlog.Error("failed to publish critical error email event", "reservationId", e.ReservationID, "error", publishErr)
		}
		return nil // return nil to avoid retrying the cancellation, as the failure is likely unrecoverable (e.g. invalid reservation ID or last name)
	}

	return nil
}
