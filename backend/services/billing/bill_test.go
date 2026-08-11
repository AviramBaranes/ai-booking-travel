package billing

import (
	"context"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/accounts"
	contact "encore.app/services/accounts/handlers/contact"
	"encore.app/services/accounts/handlers/office"
	"encore.app/services/reservation"
	"encore.dev/et"
)

// --- Helpers ---

func billingPenalty(id int64, amount float64, currencyCode string, penaltyType string) reservation.BillingPenalty {
	return reservation.BillingPenalty{
		ID:                  id,
		ReservationID:       id * 10,
		BrokerReservationID: "BRK-PEN",
		Type:                penaltyType,
		Amount:              amount,
		CurrencyCode:        currencyCode,
		CurrencyRate:        3.5,
	}
}

func billingReservation(id int64, totalPrice float64, currencyCode string) reservation.BillingReservation {
	return reservation.BillingReservation{
		ID:                  id,
		BrokerReservationID: "BRK-RES",
		PaymentStatus:       reservation.PaymentStatusUnpaid,
		ReservationStatus:   reservation.ReservationStatusVouchered,
		CarPurchasePrice:    100,
		TotalProfit:         20,
		TotalPrice:          totalPrice,
		CurrencyCode:        currencyCode,
		CurrencyRate:        3.5,
	}
}

// --- Tests ---

func TestBillParamsValidation(t *testing.T) {
	officeID := int64(1)

	validParams := func() BillParams {
		return BillParams{
			IDs:          []int64{1},
			PenaltyIDs:   []int64{},
			TotalPaid:    100,
			TransferDate: "2026-08-01",
			OfficeID:     &officeID,
		}
	}

	t.Run("accepts a bill of reservations only", func(t *testing.T) {
		if err := validParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("accepts a bill of penalties only", func(t *testing.T) {
		p := validParams()
		p.IDs = []int64{}
		p.PenaltyIDs = []int64{7}

		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("rejects a bill with neither reservations nor penalties", func(t *testing.T) {
		p := validParams()
		p.IDs = []int64{}

		api_errors.AssertApiError(t, ErrNothingToBill, p.Validate())
	})
}

func TestResolveBillCurrency(t *testing.T) {
	reservations := reservationSet{
		1: billingReservation(1, 500, "USD"),
		2: billingReservation(2, 300, "EUR"),
	}
	penalties := penaltySet{
		10: billingPenalty(10, 150, "USD", reservation.PenaltyTypeNoShow),
		20: billingPenalty(20, 80, "EUR", reservation.PenaltyTypeCancellation),
	}

	t.Run("returns the shared currency of reservations and penalties", func(t *testing.T) {
		p := BillParams{IDs: []int64{1}, PenaltyIDs: []int64{10}}

		got, err := resolveBillCurrency(p, reservations, penalties)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != "USD" {
			t.Fatalf("expected USD, got %s", got)
		}
	})

	t.Run("returns the currency of a penalties-only bill", func(t *testing.T) {
		p := BillParams{IDs: []int64{}, PenaltyIDs: []int64{20}}

		got, err := resolveBillCurrency(p, reservations, penalties)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != "EUR" {
			t.Fatalf("expected EUR, got %s", got)
		}
	})

	t.Run("rejects a penalty in a different currency than the reservations", func(t *testing.T) {
		p := BillParams{IDs: []int64{1}, PenaltyIDs: []int64{20}}

		_, err := resolveBillCurrency(p, reservations, penalties)
		api_errors.AssertApiError(t, ErrMismatchedCurrencies, err)
	})
}

func TestCalculateTotalAmount(t *testing.T) {
	reservations := reservationSet{1: billingReservation(1, 500, "USD")}
	penalties := penaltySet{10: billingPenalty(10, 150, "USD", reservation.PenaltyTypeNoShow)}

	// Both are converted at their own rate of 3.5: 500*3.5 + 150*3.5 = 2275.
	got := calculateTotalAmount(reservations, penalties, []int64{1}, []int64{10})
	if got != 2275 {
		t.Fatalf("expected total 2275, got %v", got)
	}
}

func TestBuildPenaltyInvoiceItem(t *testing.T) {
	penalty := billingPenalty(10, 150, "USD", reservation.PenaltyTypeNoShow)

	item := buildPenaltyInvoiceItem(penalty, 2, 3.5)

	if !item.IsTaxExempt {
		t.Error("expected a penalty item to be tax exempt")
	}
	if item.UnitPrice == nil || *item.UnitPrice != 150 {
		t.Errorf("expected unit price 150, got %v", item.UnitPrice)
	}
	// A penalty is a pass-through charge with no profit, so it never carries a VAT-inclusive price.
	if item.UnitPriceIncvat != nil {
		t.Errorf("expected no VAT-inclusive price, got %v", *item.UnitPriceIncvat)
	}
	if item.Quantity != 1 {
		t.Errorf("expected quantity 1, got %d", item.Quantity)
	}
	// Reservations and penalties are numbered independently, so their SKUs must not collide.
	if item.SKU == "10" {
		t.Errorf("expected the penalty SKU to be namespaced, got %s", item.SKU)
	}
}

func TestBuildInvoiceItems(t *testing.T) {
	reservations := reservationSet{1: billingReservation(1, 500, "USD")}
	penalties := penaltySet{10: billingPenalty(10, 150, "USD", reservation.PenaltyTypeNoShow)}

	// A reservation contributes two lines (purchase + profit), a penalty exactly one.
	items := buildInvoiceItems([]int64{1}, []int64{10}, reservations, penalties)
	if len(items) != 3 {
		t.Fatalf("expected 3 invoice items, got %d", len(items))
	}

	if items[2].SKU != penaltySKU(10) {
		t.Errorf("expected the penalty to be the last item, got SKU %s", items[2].SKU)
	}
}

func TestBuildReceiptItems(t *testing.T) {
	reservations := reservationSet{1: billingReservation(1, 500, "USD")}
	penalties := penaltySet{10: billingPenalty(10, 150, "USD", reservation.PenaltyTypeNoShow)}

	// Both a reservation and a penalty contribute exactly one receipt line.
	items := buildReceiptItems([]int64{1}, []int64{10}, reservations, penalties)
	if len(items) != 2 {
		t.Fatalf("expected 2 receipt items, got %d", len(items))
	}

	if items[1].UnitPrice == nil || *items[1].UnitPrice != 150 {
		t.Errorf("expected penalty receipt line of 150, got %v", items[1].UnitPrice)
	}
}

// TestBillPenalties always sets SkipInvoiceCreation, since iCount has no mock and a document
// would be created for real.
func TestBillPenalties(t *testing.T) {
	ctx := context.Background()
	officeID := int64(1)

	// mockBillingEntity wires up the endpoints Bill depends on, recording what it settles.
	mockBillingEntity := func(t *testing.T, groups []reservation.CurrencyGroup) (resolvedReservations, resolvedPenalties *[]int64, balanceChange *float64) {
		t.Helper()
		var gotReservations, gotPenalties []int64
		var gotBalanceChange float64

		et.MockEndpoint(accounts.GetIcountClientID, func(_ context.Context, _ contact.GetIcountClientIDParams) (*contact.GetIcountClientIDResponse, error) {
			return &contact.GetIcountClientIDResponse{ClientID: 1}, nil
		})
		et.MockEndpoint(reservation.ListOpenReservationsByBillingEntity, func(_ context.Context, _ *reservation.ListOpenReservationsByBillingEntityParams) (*reservation.ListOpenReservationsByBillingEntityResponse, error) {
			return &reservation.ListOpenReservationsByBillingEntityResponse{CurrencyGroups: groups}, nil
		})
		et.MockEndpoint(reservation.ResolveReservations, func(_ context.Context, p reservation.ResolveReservationsParams) error {
			gotReservations = p.IDs
			return nil
		})
		et.MockEndpoint(reservation.ResolvePenalties, func(_ context.Context, p reservation.ResolvePenaltiesParams) error {
			gotPenalties = p.IDs
			return nil
		})
		et.MockEndpoint(accounts.UpdateOfficeBalanceDue, func(_ context.Context, p office.UpdateOfficeBalanceDueParams) error {
			gotBalanceChange = p.BalanceChange
			return nil
		})

		return &gotReservations, &gotPenalties, &gotBalanceChange
	}

	t.Run("resolves penalties and credits their amount to the balance", func(t *testing.T) {
		groups := []reservation.CurrencyGroup{{
			CurrencyCode: "USD",
			Reservations: []reservation.BillingReservation{billingReservation(1, 500, "USD")},
			Penalties:    []reservation.BillingPenalty{billingPenalty(10, 150, "USD", reservation.PenaltyTypeNoShow)},
		}}
		gotReservations, gotPenalties, gotBalanceChange := mockBillingEntity(t, groups)

		if _, err := Bill(ctx, BillParams{
			IDs:                 []int64{1},
			PenaltyIDs:          []int64{10},
			SkipInvoiceCreation: true,
			TotalPaid:           2275,
			TransferDate:        "2026-08-01",
			OfficeID:            &officeID,
		}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(*gotPenalties) != 1 || (*gotPenalties)[0] != 10 {
			t.Errorf("expected penalty 10 to be resolved, got %v", *gotPenalties)
		}
		if len(*gotReservations) != 1 || (*gotReservations)[0] != 1 {
			t.Errorf("expected reservation 1 to be resolved, got %v", *gotReservations)
		}
		// 500*3.5 + 150*3.5, credited against the balance due.
		if *gotBalanceChange != -2275 {
			t.Errorf("expected balance change -2275, got %v", *gotBalanceChange)
		}
	})

	t.Run("bills penalties with no reservations", func(t *testing.T) {
		groups := []reservation.CurrencyGroup{{
			CurrencyCode: "EUR",
			Reservations: []reservation.BillingReservation{},
			Penalties:    []reservation.BillingPenalty{billingPenalty(20, 80, "EUR", reservation.PenaltyTypeCancellation)},
		}}
		gotReservations, gotPenalties, gotBalanceChange := mockBillingEntity(t, groups)

		if _, err := Bill(ctx, BillParams{
			IDs:                 []int64{},
			PenaltyIDs:          []int64{20},
			SkipInvoiceCreation: true,
			TotalPaid:           280,
			TransferDate:        "2026-08-01",
			OfficeID:            &officeID,
		}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(*gotPenalties) != 1 || (*gotPenalties)[0] != 20 {
			t.Errorf("expected penalty 20 to be resolved, got %v", *gotPenalties)
		}
		// With no reservation ids the resolve endpoint must not be called at all: it rejects an
		// empty list, which would abort the bill before the balance is updated.
		if len(*gotReservations) != 0 {
			t.Errorf("expected no reservations to be resolved, got %v", *gotReservations)
		}
		if *gotBalanceChange != -280 {
			t.Errorf("expected balance change -280, got %v", *gotBalanceChange)
		}
	})

	t.Run("rejects a penalty that belongs to another billing entity", func(t *testing.T) {
		groups := []reservation.CurrencyGroup{{
			CurrencyCode: "USD",
			Reservations: []reservation.BillingReservation{},
			Penalties:    []reservation.BillingPenalty{billingPenalty(10, 150, "USD", reservation.PenaltyTypeNoShow)},
		}}
		mockBillingEntity(t, groups)

		_, err := Bill(ctx, BillParams{
			IDs:                 []int64{},
			PenaltyIDs:          []int64{999},
			SkipInvoiceCreation: true,
			TotalPaid:           100,
			TransferDate:        "2026-08-01",
			OfficeID:            &officeID,
		})
		api_errors.AssertApiError(t, ErrInvalidReservationID, err)
	})
}
