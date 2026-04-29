package reservation

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/services/reservation/db"
	"encore.dev/et"
	"go.uber.org/mock/gomock"
)

// futurePickup returns a pickup date string sufficiently far in the future to be
// outside the cancellation window (i.e. cancellable).
func futurePickup() string {
	return time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02")
}

// outboxRowForReservation returns the most recent outbox row whose payload
// references the given reservation id, or nil if none exists.
func outboxRowForReservation(t *testing.T, ctx context.Context, reservationID int64) *BookingCancellationEvent {
	t.Helper()
	rows, err := query.GetOutboxByTopic(ctx, "booking-cancellation-events")
	if err != nil {
		t.Fatalf("query outbox: %v", err)
	}
	for _, row := range rows {
		var ev BookingCancellationEvent
		if err := json.Unmarshal(row.Data, &ev); err != nil {
			t.Fatalf("unmarshal outbox event: %v", err)
		}
		if ev.ReservationID == reservationID {
			return &ev
		}
	}
	return nil
}

func TestCancelReservation(t *testing.T) {
	ctx := context.Background()

	t.Run("returns 404 for non-existent reservation", func(t *testing.T) {
		err := CancelReservation(authContext(1), 99999999)
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns 404 when reservation belongs to another user", func(t *testing.T) {
		const ownerID int32 = 1001
		params := validCreateReservationParams()
		params.UserID = ownerID
		params.PickupDate = futurePickup()
		res, err := CreateReservation(ctx, *params)
		if err != nil {
			t.Fatalf("failed to create reservation: %v", err)
		}

		err = CancelReservation(authContext(2002), res.ID)
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns FailedPrecondition when past cancellation window", func(t *testing.T) {
		const userID int32 = 1003
		// Pickup tomorrow at noon -> within the 48h cancellation window.
		pickup := time.Now().Add(24 * time.Hour)
		params := validCreateReservationParams()
		params.UserID = userID
		params.PickupDate = pickup.Format("2006-01-02")
		params.PickupTime = "12:00"
		params.ReturnDate = pickup.Add(48 * time.Hour).Format("2006-01-02")
		res, err := CreateReservation(ctx, *params)
		if err != nil {
			t.Fatalf("failed to create reservation: %v", err)
		}

		err = CancelReservation(authContext(userID), res.ID)
		api_errors.AssertApiError(t, ErrCancellationWindowExceeded, err)
	})

	t.Run("does not enqueue a cancellation event when DB lookup fails", func(t *testing.T) {
		const userID int32 = 1004
		q, s := mockService(t)
		q.EXPECT().GetReservationByID(gomock.Any(), int64(42)).Return(db.GetReservationByIDRow{}, errors.New("db error"))
		et.MockService[Interface]("reservation", s)

		err := CancelReservation(authContext(userID), 42)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)

		if row := outboxRowForReservation(t, ctx, 42); row != nil {
			t.Fatalf("expected no outbox row for reservation 42, got %+v", row)
		}
	})

	t.Run("atomically cancels reservation and writes event to the outbox", func(t *testing.T) {
		const userID int32 = 1005
		params := validCreateReservationParams()
		params.UserID = userID
		params.PickupDate = futurePickup()
		params.ReturnDate = time.Now().Add(35 * 24 * time.Hour).Format("2006-01-02")
		res, err := CreateReservation(ctx, *params)
		if err != nil {
			t.Fatalf("failed to create reservation: %v", err)
		}

		authCtx := authContext(userID)
		if err := CancelReservation(authCtx, res.ID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Reservation status must be updated (transaction committed).
		got, err := GetReservation(authCtx, res.ID)
		if err != nil {
			t.Fatalf("failed to fetch reservation: %v", err)
		}
		if got.ReservationStatus != string(db.ReservationStatusCanceled) {
			t.Fatalf("expected reservation_status %q, got %q", db.ReservationStatusCanceled, got.ReservationStatus)
		}
		if got.PaymentStatus != string(db.PaymentStatusRefundPending) {
			t.Fatalf("expected payment_status %q, got %q", db.PaymentStatusRefundPending, got.PaymentStatus)
		}

		// And an outbox row must have been written in the same transaction.
		ev := outboxRowForReservation(t, ctx, res.ID)
		if ev == nil {
			t.Fatal("expected an outbox row for the cancellation, got none")
		}
		if ev.BrokerReservationID != params.BrokerReservationID {
			t.Fatalf("expected event BrokerReservationID %q, got %q", params.BrokerReservationID, ev.BrokerReservationID)
		}
		if string(ev.Broker) != params.Broker {
			t.Fatalf("expected event Broker %q, got %q", params.Broker, ev.Broker)
		}
		if ev.SupplierCode != params.SupplierCode {
			t.Fatalf("expected event SupplierCode %q, got %q", params.SupplierCode, ev.SupplierCode)
		}
		if ev.LastName != params.DriverLastName {
			t.Fatalf("expected event LastName %q, got %q", params.DriverLastName, ev.LastName)
		}
	})
}
