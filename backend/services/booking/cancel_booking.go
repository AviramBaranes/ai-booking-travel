package booking

import (
	"context"

	"encore.app/services/booking/handlers/booking_handlers"
	"encore.app/services/reservation"
	"encore.dev/pubsub"
)

var _ = pubsub.NewSubscription(
	reservation.BookingCancellationEvents,
	"cancel-booking",
	pubsub.SubscriptionConfig[*reservation.BookingCancellationEvent]{
		Handler: CancelBooking,
	},
)

func CancelBooking(ctx context.Context, e *reservation.BookingCancellationEvent) error {
	bs := booking_handlers.NewBookingService(nil)
	return bs.CancelBooking(ctx, e)
}
