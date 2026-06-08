package actions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/accounts"
	emailevents "encore.app/services/notifications/events"
	"encore.app/services/reservation/db"
	"encore.dev/beta/errs"
	"encore.dev/pubsub"
	"encore.dev/rlog"
	"github.com/jackc/pgx/v5"
	"x.encore.dev/infra/pubsub/outbox"
)

const cancellationWindowHours = 48

var ErrCancellationWindowExceeded = api_errors.NewErrorWithDetail(errs.FailedPrecondition, "cancellation window exceeded", api_errors.ErrorDetails{
	Code: api_errors.CancellationWindowExceeded,
})

var emailRequestedTopic = pubsub.TopicRef[pubsub.Publisher[*emailevents.EmailEvent]](emailevents.EmailRequestedTopic)
var emailPublisher = emailevents.NewPublisher(emailRequestedTopic)

// BookingCancellationEvent represents the details of a reservation cancellation event.
type BookingCancellationEvent struct {
	ReservationID       int64
	Broker              db.Broker
	BrokerReservationID string
	LastName            string
	SupplierCode        string
}

func (s *ActionService) CancelReservation(ctx context.Context, id int64) error {
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

	isLateCancellation := !canCancel(reservation)

	if err := db.WithTx(ctx, s.pool, func(q db.Querier, tx pgx.Tx) error {
		if err := q.CancelReservation(ctx, id); err != nil {
			rlog.Error("failed to cancel reservation", "error", err, "reservationId", id)
			return err
		}

		event := &BookingCancellationEvent{
			ReservationID:       reservation.ID,
			Broker:              reservation.Broker,
			BrokerReservationID: reservation.BrokerReservationID,
			LastName:            reservation.DriverLastName,
			SupplierCode:        reservation.SupplierCode,
		}

		topic := outbox.Bind(s.cancellationTopic, outbox.PgxTxPersister(tx))
		if _, err := topic.Publish(ctx, event); err != nil {
			rlog.Error("failed to enqueue booking cancellation event", "error", err, "reservationId", reservation.ID)
			return err
		}

		return nil
	}); err != nil {
		return api_errors.ErrInternalError
	}

	if _, err := emailPublisher.Publish(ctx, emailevents.EmailEventTypeCancellation, emailevents.CancellationEmailPayload{
		UserID:             reservation.UserID,
		BookingReferenceID: reservation.BrokerReservationID,
		DriverFullName:     fmt.Sprintf("%s %s %s", reservation.DriverTitle, reservation.DriverFirstName, reservation.DriverLastName),
	}); err != nil {
		rlog.Error("failed to publish cancellation email event", "error", err, "reservationId", reservation.ID)
	}

	if isLateCancellation {
		if _, err := emailPublisher.Publish(ctx, emailevents.EmailEventTypeLateCancellationAlert, emailevents.LateCancellationAlertEmailPayload{
			ReservationID:       reservation.ID,
			BrokerReservationID: reservation.BrokerReservationID,
			AgentID:             reservation.UserID,
			OfficeID:            reservation.OfficeID,
			OrganizationID:      reservation.OrganizationID,
		}); err != nil {
			rlog.Error("failed to publish late cancellation alert email event", "error", err, "reservationId", reservation.ID)
		}
	}

	return nil
}

// canCancel checks if the reservation can be canceled based on the current time and the pickup time.
func canCancel(reservation db.Reservation) bool {
	pickupDateTime, err := dbadapters.CombineDateTime(reservation.PickupDate, reservation.PickupTime)
	if err != nil {
		rlog.Error("failed to parse pickup date and time", "error", err, "reservationId", reservation.ID, "pickupTime", reservation.PickupTime)
		return false
	}

	cancellationDeadline := pickupDateTime.Add(-cancellationWindowHours * time.Hour)
	return time.Now().Before(cancellationDeadline)
}
