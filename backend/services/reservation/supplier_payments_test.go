package reservation

import (
	"context"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
)

// --- Helpers ---

// seededReservations indexes the reservations of a currency-grouped response by reservation ID,
// recording which currency group each one landed in. The endpoint is not user-scoped and the test
// database is shared with parallel tests, so assertions are made per seeded ID rather than on
// the whole response.
type seededReservations map[int64]struct {
	currencyCode string
	reservation  UnpaidSupplierReservation
}

func indexByID(t *testing.T, resp *ListUnpaidSupplierReservationsResponse) seededReservations {
	t.Helper()
	indexed := make(seededReservations)
	for _, group := range resp.CurrencyGroups {
		if len(group.Reservations) == 0 {
			t.Errorf("currency group %q is empty; groups must never be emitted without reservations", group.CurrencyCode)
		}
		for _, r := range group.Reservations {
			if _, exists := indexed[r.ID]; exists {
				t.Errorf("reservation %d appears more than once in the response", r.ID)
			}
			indexed[r.ID] = struct {
				currencyCode string
				reservation  UnpaidSupplierReservation
			}{currencyCode: group.CurrencyCode, reservation: r}
		}
	}
	return indexed
}

func (s seededReservations) assertInCurrency(t *testing.T, id int64, currencyCode string) UnpaidSupplierReservation {
	t.Helper()
	entry, ok := s[id]
	if !ok {
		t.Fatalf("expected reservation %d in the response, but it was missing", id)
	}
	if entry.currencyCode != currencyCode {
		t.Errorf("reservation %d grouped under currency %q, want %q", id, entry.currencyCode, currencyCode)
	}
	return entry.reservation
}

func (s seededReservations) assertAbsent(t *testing.T, id int64, reason string) {
	t.Helper()
	if _, ok := s[id]; ok {
		t.Errorf("reservation %d should have been excluded (%s), but it was returned", id, reason)
	}
}

// --- Tests ---

func TestListUnpaidSupplierReservations_Validation(t *testing.T) {
	tests := []struct {
		name    string
		params  ListUnpaidSupplierReservationsParams
		wantErr error
	}{
		{
			name:    "rejects missing broker",
			params:  ListUnpaidSupplierReservationsParams{},
			wantErr: invalidValueErr("broker"),
		},
		{
			name:    "rejects unknown broker",
			params:  ListUnpaidSupplierReservationsParams{Broker: "avis"},
			wantErr: invalidValueErr("broker"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			api_errors.AssertApiError(t, tc.wantErr, tc.params.Validate())
		})
	}

	for _, validBroker := range []string{"flex", "hertz"} {
		t.Run("accepts broker "+validBroker, func(t *testing.T) {
			t.Parallel()
			p := ListUnpaidSupplierReservationsParams{Broker: validBroker}
			if err := p.Validate(); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestListUnpaidSupplierReservations(t *testing.T) {
	const userID int64 = 90001
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	// Two hertz reservations in USD, one in EUR — exercises multi-currency grouping.
	usdEarlyID := seedReservation(t, ctx, s, userID, func(p *CreateReservationParams) {
		p.BrokerReservationID = "SUPPAY-USD-EARLY"
		p.Broker = "hertz"
		p.CurrencyCode = "USD"
		p.PickupDate = "2026-09-01"
		p.DropoffDate = "2026-09-04"
	})
	usdLateID := seedReservation(t, ctx, s, userID, func(p *CreateReservationParams) {
		p.BrokerReservationID = "SUPPAY-USD-LATE"
		p.Broker = "hertz"
		p.CurrencyCode = "USD"
		p.PickupDate = "2026-09-20"
		p.DropoffDate = "2026-09-23"
	})
	eurID := seedReservation(t, ctx, s, userID, func(p *CreateReservationParams) {
		p.BrokerReservationID = "SUPPAY-EUR"
		p.Broker = "hertz"
		p.CurrencyCode = "EUR"
		p.PickupDate = "2026-09-10"
		p.DropoffDate = "2026-09-13"
	})

	// A canceled one, one already paid to the supplier, and one belonging to another broker —
	// all three must be filtered out.
	canceledID := seedReservation(t, ctx, s, userID, func(p *CreateReservationParams) {
		p.BrokerReservationID = "SUPPAY-CANCELED"
		p.Broker = "hertz"
		p.CurrencyCode = "USD"
	})
	if err := s.query.CancelReservation(ctx, canceledID); err != nil {
		t.Fatalf("failed to cancel reservation: %v", err)
	}

	paidID := seedReservation(t, ctx, s, userID, func(p *CreateReservationParams) {
		p.BrokerReservationID = "SUPPAY-ALREADY-PAID"
		p.Broker = "hertz"
		p.CurrencyCode = "USD"
	})
	expenseID := "EXPENSE-1"
	if err := s.query.MarkReservationsPaidToSupplier(ctx, db.MarkReservationsPaidToSupplierParams{
		Ids:               []int64{paidID},
		SupplierExpenseID: &expenseID,
		SupplierPaidAt:    dbadapters.DBTime(time.Now()),
	}); err != nil {
		t.Fatalf("failed to mark reservation as paid to supplier: %v", err)
	}

	flexID := seedReservation(t, ctx, s, userID, func(p *CreateReservationParams) {
		p.BrokerReservationID = "SUPPAY-FLEX"
		p.Broker = "flex"
		p.CurrencyCode = "USD"
	})

	resp, err := ListUnpaidSupplierReservations(ctx, &ListUnpaidSupplierReservationsParams{Broker: "hertz"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	indexed := indexByID(t, resp)

	t.Run("groups unpaid reservations by currency", func(t *testing.T) {
		indexed.assertInCurrency(t, usdEarlyID, "USD")
		indexed.assertInCurrency(t, usdLateID, "USD")
		indexed.assertInCurrency(t, eurID, "EUR")
	})

	t.Run("excludes canceled, already-paid and other-broker reservations", func(t *testing.T) {
		indexed.assertAbsent(t, canceledID, "reservation is canceled")
		indexed.assertAbsent(t, paidID, "supplier was already paid")
		indexed.assertAbsent(t, flexID, "reservation belongs to another broker")
	})

	t.Run("orders reservations within a currency group by pickup date", func(t *testing.T) {
		for _, group := range resp.CurrencyGroups {
			if group.CurrencyCode != "USD" {
				continue
			}
			var earlyIdx, lateIdx = -1, -1
			for i, r := range group.Reservations {
				switch r.ID {
				case usdEarlyID:
					earlyIdx = i
				case usdLateID:
					lateIdx = i
				}
			}
			if earlyIdx == -1 || lateIdx == -1 {
				t.Fatalf("expected both USD reservations in the USD group, got indexes %d and %d", earlyIdx, lateIdx)
			}
			if earlyIdx > lateIdx {
				t.Errorf("reservation picked up on 2026-09-01 is listed after the one picked up on 2026-09-20")
			}
		}
	})

	t.Run("maps the reservation fields", func(t *testing.T) {
		got := indexed.assertInCurrency(t, usdEarlyID, "USD")

		// Derived from validCreateReservationParams: (100 purchase + 15 broker ERP) at 45% markup
		// = 166.75, less a 10% discount = 150.075, plus the 20 BT ERP price = 170.075,
		// stored as NUMERIC(12, 2).
		const wantTotalPrice = 170.08

		want := UnpaidSupplierReservation{
			ID:                  usdEarlyID,
			BrokerReservationID: "SUPPAY-USD-EARLY",
			DriverName:          "Mr John Doe",
			PickupDate:          "2026-09-01",
			PickupLocationName:  "Airport Terminal 1",
			RentalDays:          4,
			TotalPrice:          wantTotalPrice,
			ReservationStatus:   ReservationStatusBooked,
			PaymentStatus:       PaymentStatusUnpaid,
		}
		if got != want {
			t.Errorf("unexpected reservation:\ngot  %+v\nwant %+v", got, want)
		}
	})
}
