package reservation

import (
	"context"
	"errors"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	"encore.dev/beta/errs"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

const (
	cancellationWindowHours = 48
)

var (
	ErrCancellationWindowExceeded = api_errors.NewErrorWithDetail(errs.FailedPrecondition, "cancellation window exceeded", api_errors.ErrorDetails{
		Code: api_errors.CancellationWindowExceeded,
	})
)

// encore:api auth method=POST path=/api/reservation/:id/cancel tag:agent
func (s *Service) CancelReservation(ctx context.Context, id int64) error {
	reservation, err := s.query.GetReservationByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api_errors.ErrNotFound
		}
		rlog.Error("failed to get reservation by id", "error", err, "reservationId", id)
		return api_errors.ErrInternalError
	}

	authData := accounts.GetAuthData()
	if reservation.UserID != authData.UserID {
		rlog.Warn("user attempted to cancel a reservation that does not belong to them", "userId", authData.UserID, "reservationId", id)
		return api_errors.ErrNotFound
	}

	if !canCancel(reservation) {
		return ErrCancellationWindowExceeded
	}

	if err := db.WithTx(ctx, s.pool, func(q db.Querier) error {
		if err := q.CancelReservation(ctx, id); err != nil {
			rlog.Error("failed to cancel reservation", "error", err, "reservationId", id)
			return err
		}

		event := &BookingCancellationEvent{
			ReservationID:       reservation.ID,
			Broker:              reservation.Broker,
			BrokerReservationID: reservation.BrokerReservationID,
		}

		if _, err := BookingCancellationEvents.Publish(ctx, event); err != nil {
			rlog.Error("failed to publish booking cancellation event", "error", err, "reservationId", reservation.ID)
			return err
		}

		return nil
	}); err != nil {
		return api_errors.ErrInternalError
	}

	return nil
}

// canCancel checks if the reservation can be canceled based on the current time and the pickup time.
func canCancel(reservation db.GetReservationByIDRow) bool {
	now := time.Now()
	pickupDateTime, err := db.CombineDateTime(reservation.PickupDate, reservation.PickupTime)
	if err != nil {
		rlog.Error("failed to combine pickup date and time", "error", err, "reservationId", reservation.ID)
		return false
	}
	cancellationDeadline := pickupDateTime.Add(-cancellationWindowHours * time.Hour)
	return now.Before(cancellationDeadline)
}

// BookingCancellationEvent represents the details of a reservation cancellation event.
type BookingCancellationEvent struct {
	ReservationID       int64
	Broker              db.Broker
	BrokerReservationID string
}

// BookingCancellationEvents is a pub/sub topic that publishes events whenever a reservation is canceled.
var BookingCancellationEvents = pubsub.NewTopic[*BookingCancellationEvent]("booking-cancellation-events", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})
