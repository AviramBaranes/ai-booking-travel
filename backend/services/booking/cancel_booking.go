package booking

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	"encore.app/services/booking/db"
	"encore.app/services/reservation"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

var _ = pubsub.NewSubscription(
	reservation.BookingCancellationEvents,
	"cancel-booking",
	pubsub.SubscriptionConfig[*reservation.BookingCancellationEvent]{
		Handler: CancelBooking,
	},
)

// CancelBooking handles the cancellation of a booking by processing the BookingCancellationEvent received from the reservation service.
func CancelBooking(ctx context.Context, e *reservation.BookingCancellationEvent) error {
	b, err := getCanceler(db.Broker(e.Broker))
	if err != nil {
		rlog.Error("unsupported broker for cancellation", "broker", b, "reservationId", e.ReservationID)
		return err
	}

	err = b.Cancel(e.BrokerReservationID, e.LastName, e.SupplierCode)
	if err != nil {
		rlog.Error("failed to cancel booking", "broker", b, "reservationId", e.ReservationID, "error", err)
		return err
	}

	return nil
}

// getCanceler returns the appropriate canceler implementation based on the broker.
func getCanceler(b db.Broker) (broker.Canceler, error) {
	switch b {
	case db.BrokerHertz:
		return broker.NewHertz(), nil
	case db.BrokerFlex:
		return broker.NewFlex(), nil
	default:
		return nil, api_errors.ErrInternalError
	}
}
