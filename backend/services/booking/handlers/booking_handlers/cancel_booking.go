package booking_handlers

import (
	"context"
	"fmt"

	emailevents "encore.app/internal/email_events"
	"encore.app/services/booking/db"
	"encore.app/services/reservation"
	"encore.dev/rlog"
)

// CancelBooking handles the cancellation of a booking by processing the BookingCancellationEvent received from the reservation service.
func (s *BookingService) CancelBooking(ctx context.Context, e *reservation.BookingCancellationEvent) error {
	b, err := getCanceler(db.Broker(e.Broker))
	if err != nil {
		rlog.Error("unsupported broker for cancellation", "broker", b, "reservationId", e.ReservationID)
		emailevents.PublishEmailEvent(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
			Subject: "Unsupported broker for booking cancellation",
			Message: fmt.Sprintf("unsupported broker for cancellation: %v, reservationId: %v", b, e.ReservationID),
		})
		return err
	}

	err = b.Cancel(e.BrokerReservationID, e.LastName, e.SupplierCode)
	if err != nil {
		rlog.Error("failed to cancel booking", "broker", b, "reservationId", e.ReservationID, "error", err)
		emailevents.PublishEmailEvent(ctx, emailevents.EmailEventTypeCriticalError, emailevents.CriticalErrorEmailPayload{
			Subject: "Failed to cancel booking",
			Message: fmt.Sprintf("failed to cancel booking: %s, bookingID: %v, reservationId: %v, error: %v", b.Name(), e.BrokerReservationID, e.ReservationID, err),
		})
		return nil // return nil to avoid retrying the cancellation, as the failure is likely unrecoverable (e.g. invalid reservation ID or last name)
	}

	return nil
}
