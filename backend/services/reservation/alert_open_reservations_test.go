package reservation

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	dbadapters "encore.app/internal/db_adapters"
	emailevents "encore.app/services/notifications/events"
	"encore.app/services/reservation/db"
	"encore.dev/et"
)

func TestAlertOpenReservations(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	now := time.Now().UTC()

	// Case 1: vouchered + pickup within 7 days.
	// The query filters for reservation_status = 'booked', so this should be ignored.
	p1 := validCreateReservationParams()
	p1.UserID = 8001
	p1.PickupDate = now.Add(2 * 24 * time.Hour).Format("2006-01-02")
	p1.DropoffDate = now.Add(4 * 24 * time.Hour).Format("2006-01-02")
	p1.PickupTime = "12:00"
	res1, err := CreateReservation(ctx, *p1)
	if err != nil {
		t.Fatalf("seed vouchered reservation: %v", err)
	}
	voucher := "TESTVOUCH-001"
	if _, err := q.ApplyVoucher(ctx, db.ApplyVoucherParams{
		ID:            res1.ID,
		UserID:        p1.UserID,
		VoucherNumber: &voucher,
		CurrencyRate:  dbadapters.NumericFromFloat64(1),
	}); err != nil {
		t.Fatalf("apply voucher: %v", err)
	}

	// Case 2: booked + far pickup (>7 days from now).
	// The query filters pickup_date <= CURRENT_DATE + 7 days, so this should be ignored.
	p2 := validCreateReservationParams()
	p2.UserID = 8002
	p2.PickupDate = now.Add(14 * 24 * time.Hour).Format("2006-01-02")
	p2.DropoffDate = now.Add(16 * 24 * time.Hour).Format("2006-01-02")
	p2.PickupTime = "12:00"
	if _, err := CreateReservation(ctx, *p2); err != nil {
		t.Fatalf("seed far-pickup reservation: %v", err)
	}

	// Case 3: booked + within the 48h cancellation window → should be auto-cancelled.
	// Pickup is tomorrow at 12:00 UTC — always within 48 hours.
	p3 := validCreateReservationParams()
	p3.UserID = 8003
	p3.PickupDate = now.Add(24 * time.Hour).Format("2006-01-02")
	p3.DropoffDate = now.Add(3 * 24 * time.Hour).Format("2006-01-02")
	p3.PickupTime = "12:00"
	if _, err := CreateReservation(ctx, *p3); err != nil {
		t.Fatalf("seed within-window reservation: %v", err)
	}

	// Case 4: booked + inside the alert window (>48h, ≤7 days) → should receive an alert email.
	// Pickup is 4 days from now at 12:00 UTC — always >48h and within 7 days.
	p4 := validCreateReservationParams()
	p4.UserID = 8004
	p4.PickupDate = now.Add(4 * 24 * time.Hour).Format("2006-01-02")
	p4.DropoffDate = now.Add(6 * 24 * time.Hour).Format("2006-01-02")
	p4.PickupTime = "12:00"
	res4, err := CreateReservation(ctx, *p4)
	if err != nil {
		t.Fatalf("seed alert-window reservation: %v", err)
	}

	// Mock the CancelReservation endpoint so the auth check is bypassed and we can
	// assert which IDs it was called with.
	var cancelledIDs []int64
	et.MockEndpoint(CancelReservation, func(ctx context.Context, id int64) error {
		cancelledIDs = append(cancelledIDs, id)
		return nil
	})

	if err := AlertOpenReservations(ctx); err != nil {
		t.Fatalf("AlertOpenReservations: %v", err)
	}

	// Auto cancellation is currently disabled.
	// if len(cancelledIDs) != 1 {
	// 	t.Fatalf("expected CancelReservation to be called once, got %d calls: %v", len(cancelledIDs), cancelledIDs)
	// }
	// if cancelledIDs[0] != res3.ID {
	// 	t.Fatalf("expected CancelReservation called with reservation %d, got %d", res3.ID, cancelledIDs[0])
	// }

	// Assert: case 4 received an open-order-alert email.
	msgs := et.Topic(emailevents.EmailRequestedTopic).PublishedMessages()
	var alertMsg *emailevents.EmailEvent
	for _, msg := range msgs {
		if msg.Type == emailevents.EmailEventTypeOpenOrderAlert {
			var payload emailevents.OpenOrderAlertEmailPayload
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				t.Fatalf("unmarshal open order alert payload: %v", err)
			}
			if payload.UserID == p4.UserID {
				alertMsg = msg
				break
			}
		}
	}
	if alertMsg == nil {
		t.Fatalf("expected an open order alert email for user %d (reservation %d), published messages: %v", p4.UserID, res4.ID, msgs)
	}

	// Assert: cases 1 and 2 were untouched (not cancelled, no alert email).
	for _, id := range cancelledIDs {
		if id == res1.ID {
			t.Errorf("vouchered reservation %d should not have been cancelled", res1.ID)
		}
	}
	for _, msg := range msgs {
		if msg.Type == emailevents.EmailEventTypeOpenOrderAlert {
			var payload emailevents.OpenOrderAlertEmailPayload
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				continue
			}
			if payload.UserID == p1.UserID || payload.UserID == p2.UserID {
				t.Errorf("unexpected alert email for user %d (should be skipped)", payload.UserID)
			}
		}
	}
}
