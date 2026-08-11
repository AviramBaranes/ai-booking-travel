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

// seededReservations indexes the reservations of a response by reservation ID. The endpoint is not
// user-scoped and the test database is shared with parallel tests, so assertions are made per
// seeded ID rather than on the whole response.
type seededReservations map[int64]UnpaidSupplierReservation

func indexByID(t *testing.T, resp *ListUnpaidSupplierReservationsResponse) seededReservations {
	t.Helper()
	indexed := make(seededReservations)
	for _, r := range resp.Reservations {
		if _, exists := indexed[r.ID]; exists {
			t.Errorf("reservation %d appears more than once in the response", r.ID)
		}
		indexed[r.ID] = r
	}
	return indexed
}

// seededPenalties indexes the penalties of a response by penalty ID, for the same reason.
type seededPenalties map[int64]UnpaidSupplierPenalty

func indexPenaltiesByID(t *testing.T, resp *ListUnpaidSupplierReservationsResponse) seededPenalties {
	t.Helper()
	indexed := make(seededPenalties)
	for _, p := range resp.Penalties {
		if _, exists := indexed[p.ID]; exists {
			t.Errorf("penalty %d appears more than once in the response", p.ID)
		}
		indexed[p.ID] = p
	}
	return indexed
}

func (s seededReservations) assertInCurrency(t *testing.T, id int64, currencyCode string) UnpaidSupplierReservation {
	t.Helper()
	r, ok := s[id]
	if !ok {
		t.Fatalf("expected reservation %d in the response, but it was missing", id)
	}
	if r.CurrencyCode != currencyCode {
		t.Errorf("reservation %d has currency %q, want %q", id, r.CurrencyCode, currencyCode)
	}
	return r
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

	t.Run("orders reservations of the same currency by pickup date", func(t *testing.T) {
		var earlyIdx, lateIdx = -1, -1
		for i, r := range resp.Reservations {
			switch r.ID {
			case usdEarlyID:
				earlyIdx = i
			case usdLateID:
				lateIdx = i
			}
		}
		if earlyIdx == -1 || lateIdx == -1 {
			t.Fatalf("expected both USD reservations in the response, got indexes %d and %d", earlyIdx, lateIdx)
		}
		if earlyIdx > lateIdx {
			t.Errorf("reservation picked up on 2026-09-01 is listed after the one picked up on 2026-09-20")
		}
	})

	t.Run("maps the reservation fields", func(t *testing.T) {
		got := indexed.assertInCurrency(t, usdEarlyID, "USD")

		// What we owe the supplier, from validCreateReservationParams: 100 purchase price plus
		// the 15 broker ERP day charge, before our markup and discount.
		const wantAmountOwed = 115.0

		want := UnpaidSupplierReservation{
			ID:                  usdEarlyID,
			BrokerReservationID: "SUPPAY-USD-EARLY",
			DriverName:          "Mr John Doe",
			PickupDate:          "2026-09-01",
			PickupLocationName:  "Airport Terminal 1",
			RentalDays:          4,
			AmountOwed:          wantAmountOwed,
			CurrencyCode:        "USD",
			ReservationStatus:   ReservationStatusBooked,
			PaymentStatus:       PaymentStatusUnpaid,
		}
		if got != want {
			t.Errorf("unexpected reservation:\ngot  %+v\nwant %+v", got, want)
		}
	})

	t.Run("lists unpaid penalties of the broker", func(t *testing.T) {
		penalty := seedPenalty(t, ctx, s, canceledID, db.PenaltyTypeNoShow, 90.00, "USD")

		// A fee already settled with the supplier, and one on a reservation of another broker.
		paidPenalty := seedPenalty(t, ctx, s, usdLateID, db.PenaltyTypeCancellation, 60.00, "USD")
		if err := s.query.MarkPenaltiesPaidToSupplier(ctx, db.MarkPenaltiesPaidToSupplierParams{
			Ids:               []int64{paidPenalty.ID},
			SupplierExpenseID: &expenseID,
			SupplierPaidAt:    dbadapters.DBTime(time.Now()),
		}); err != nil {
			t.Fatalf("failed to mark penalty as paid to supplier: %v", err)
		}
		flexPenalty := seedPenalty(t, ctx, s, flexID, db.PenaltyTypeCancellation, 45.00, "USD")

		resp, err := ListUnpaidSupplierReservations(ctx, &ListUnpaidSupplierReservationsParams{Broker: "hertz"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		indexedPenalties := indexPenaltiesByID(t, resp)

		got, ok := indexedPenalties[penalty.ID]
		if !ok {
			t.Fatalf("expected penalty %d in the response, but it was missing", penalty.ID)
		}

		want := UnpaidSupplierPenalty{
			ID:                  penalty.ID,
			ReservationID:       canceledID,
			BrokerReservationID: "SUPPAY-CANCELED",
			Type:                string(db.PenaltyTypeNoShow),
			DriverName:          "Mr John Doe",
			PickupDate:          "2026-04-01",
			PickupLocationName:  "Airport Terminal 1",
			AmountOwed:          90.00,
			CurrencyCode:        "USD",
		}
		if got != want {
			t.Errorf("unexpected penalty:\ngot  %+v\nwant %+v", got, want)
		}

		if _, ok := indexedPenalties[paidPenalty.ID]; ok {
			t.Errorf("penalty %d should have been excluded (supplier was already paid), but it was returned", paidPenalty.ID)
		}
		if _, ok := indexedPenalties[flexPenalty.ID]; ok {
			t.Errorf("penalty %d should have been excluded (belongs to another broker), but it was returned", flexPenalty.ID)
		}
	})
}

func TestPaySupplier_Validation(t *testing.T) {
	validParams := func() PaySupplierParams {
		return PaySupplierParams{
			ReservationIDs: []int64{1},
			PenaltyIDs:     []int64{},
			PaymentDate:    "2026-09-01",
		}
	}

	tests := []struct {
		name    string
		modify  func(p *PaySupplierParams)
		wantErr error
	}{
		{
			name:   "accepts reservations only",
			modify: func(p *PaySupplierParams) {},
		},
		{
			name: "accepts penalties only",
			modify: func(p *PaySupplierParams) {
				p.ReservationIDs = []int64{}
				p.PenaltyIDs = []int64{5}
			},
		},
		{
			name:   "accepts reservations and penalties together",
			modify: func(p *PaySupplierParams) { p.PenaltyIDs = []int64{5} },
		},
		{
			name:    "rejects neither reservations nor penalties",
			modify:  func(p *PaySupplierParams) { p.ReservationIDs = []int64{} },
			wantErr: ErrNothingToPay,
		},
		{
			name:    "rejects missing payment date",
			modify:  func(p *PaySupplierParams) { p.PaymentDate = "" },
			wantErr: invalidValueErr("paymentDate"),
		},
		{
			name:    "rejects a payment date that is not ISO formatted",
			modify:  func(p *PaySupplierParams) { p.PaymentDate = "01/09/2026" },
			wantErr: invalidValueErr("paymentDate"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := validParams()
			tt.modify(&p)

			err := p.Validate()
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			api_errors.AssertApiError(t, tt.wantErr, err)
		})
	}
}

// TestPaySupplier_SkipsSettledItems covers the paths that never reach iCount: an id that is no
// longer outstanding is skipped before any expense is created.
func TestPaySupplier_SkipsSettledItems(t *testing.T) {
	const userID int64 = 90003
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	expenseID := "EXPENSE-PAYSUPPLIER"

	// Already settled with the supplier.
	paidID := seedReservation(t, ctx, s, userID, func(p *CreateReservationParams) {
		p.BrokerReservationID = "PAYSUP-ALREADY-PAID"
		p.Broker = "flex"
	})
	if err := s.query.MarkReservationsPaidToSupplier(ctx, db.MarkReservationsPaidToSupplierParams{
		Ids:               []int64{paidID},
		SupplierExpenseID: &expenseID,
		SupplierPaidAt:    dbadapters.DBTime(time.Now()),
	}); err != nil {
		t.Fatalf("failed to mark reservation as paid to supplier: %v", err)
	}

	// Canceled, so we owe the supplier for the fee rather than the rental.
	canceledID := seedReservation(t, ctx, s, userID, func(p *CreateReservationParams) {
		p.BrokerReservationID = "PAYSUP-CANCELED"
		p.Broker = "flex"
	})
	if err := s.query.CancelReservation(ctx, canceledID); err != nil {
		t.Fatalf("failed to cancel reservation: %v", err)
	}

	// A fee whose expense was already recorded.
	paidPenalty := seedPenalty(t, ctx, s, canceledID, db.PenaltyTypeNoShow, 80.00, "EUR")
	if err := s.query.MarkPenaltiesPaidToSupplier(ctx, db.MarkPenaltiesPaidToSupplierParams{
		Ids:               []int64{paidPenalty.ID},
		SupplierExpenseID: &expenseID,
		SupplierPaidAt:    dbadapters.DBTime(time.Now()),
	}); err != nil {
		t.Fatalf("failed to mark penalty as paid to supplier: %v", err)
	}

	const unknownID int64 = 99999999

	resp, err := PaySupplier(ctx, PaySupplierParams{
		ReservationIDs: []int64{paidID, canceledID, unknownID},
		PenaltyIDs:     []int64{paidPenalty.ID, unknownID},
		PaymentDate:    "2026-09-01",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertIDs(t, "skipped reservations", resp.Skipped, []int64{paidID, canceledID, unknownID})
	assertIDs(t, "skipped penalties", resp.SkippedPenalties, []int64{paidPenalty.ID, unknownID})

	if len(resp.Paid) != 0 || len(resp.PaidPenalties) != 0 {
		t.Errorf("nothing was outstanding, but the response reports payments: %+v %+v", resp.Paid, resp.PaidPenalties)
	}
	if len(resp.Failed) != 0 || len(resp.FailedPenalties) != 0 {
		t.Errorf("nothing was outstanding, so nothing should have failed: %+v %+v", resp.Failed, resp.FailedPenalties)
	}
}

// assertIDs compares an id list against the expected set, order-independently.
func assertIDs(t *testing.T, label string, got, want []int64) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: got %v, want %v", label, got, want)
		return
	}

	gotSet := make(map[int64]struct{}, len(got))
	for _, id := range got {
		gotSet[id] = struct{}{}
	}
	for _, id := range want {
		if _, ok := gotSet[id]; !ok {
			t.Errorf("%s: %d missing from %v", label, id, got)
		}
	}
}
