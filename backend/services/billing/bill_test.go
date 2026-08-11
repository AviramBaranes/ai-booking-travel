package billing

import (
	"context"
	"math"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/currency"
	"encore.app/services/accounts"
	contact "encore.app/services/accounts/handlers/contact"
	"encore.app/services/accounts/handlers/office"
	"encore.app/services/reservation"
	"encore.dev/et"
)

// --- Helpers ---

func newTestService() *Service {
	return &Service{ratesCache: currency.NewCurrenciesCache(currenciesRates)}
}

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
	items := buildInvoiceItems([]int64{1}, []int64{10}, reservations, penalties, 0)
	if len(items) != 3 {
		t.Fatalf("expected 3 invoice items, got %d", len(items))
	}

	if items[2].SKU != penaltySKU(10) {
		t.Errorf("expected the penalty to be the last item, got SKU %s", items[2].SKU)
	}
}

// The discount comes out of our profit, so it has to be priced like the profit line — VAT-inclusive
// and taxable. iCount's document-level discount field apportions across the exempt lines instead,
// which both under-reduces the VAT and leaves the invoice above what the receipt collects.
func TestDiscountInvoiceItemIsTaxableAndNegative(t *testing.T) {
	reservations := reservationSet{1: billingReservation(1, 500, "USD")}

	items := buildInvoiceItems([]int64{1}, nil, reservations, penaltySet{}, 20.92)
	if len(items) != 3 {
		t.Fatalf("expected the discount to add a third line, got %d", len(items))
	}

	discount := items[2]
	if discount.IsTaxExempt {
		t.Error("expected the discount line to carry VAT, so the whole of it comes off the profit")
	}
	if discount.UnitPriceIncvat == nil || *discount.UnitPriceIncvat != -20.92 {
		t.Errorf("expected a VAT-inclusive line of -20.92, got %v", discount.UnitPriceIncvat)
	}
	// A net unit price alongside it would double-count the discount.
	if discount.UnitPrice != nil {
		t.Errorf("expected no net unit price, got %v", *discount.UnitPrice)
	}
}

// The invoice states what is owed and the receipt states how it was settled, so the two have to
// land on the same number: total price less the discount equals what was paid plus what was
// withheld. Figures are those of invoice 2051.
func TestInvoiceAndReceiptAgreeOnWhatIsOwed(t *testing.T) {
	const (
		vat       = 0.18
		discount  = 26.17
		deduction = 26.62
		totalPaid = 480.0
	)
	reservations := reservationSet{
		27: {ID: 27, CarPurchasePrice: 172.98, TotalProfit: 63.89, TotalPrice: 236.87, ERPSellingPrice: 12, CurrencyCode: "GBP"},
		29: {ID: 29, CarPurchasePrice: 218.40, TotalProfit: 77.52, TotalPrice: 295.92, ERPSellingPrice: 12, CurrencyCode: "GBP"},
	}
	ids := []int64{27, 29}

	// The invoice prices exempt lines net and taxable lines VAT-inclusive, so the amount owed is
	// the exempt lines plus the taxable lines grossed back up.
	var exempt, taxableIncVat float64
	for _, item := range buildInvoiceItems(ids, nil, reservations, penaltySet{}, discount) {
		if item.IsTaxExempt {
			exempt += *item.UnitPrice
			continue
		}
		taxableIncVat += *item.UnitPriceIncvat
	}
	owed := exempt + taxableIncVat

	settled := totalPaid + deduction
	if math.Abs(owed-settled) > 0.001 {
		t.Errorf("expected the invoice to state %v owed, got %v", settled, owed)
	}
	if math.Abs(owed-506.62) > 0.001 {
		t.Errorf("expected 506.62 owed, got %v", owed)
	}

	// And the whole discount lands on the taxable side rather than being spread over the exempt
	// purchase price: VAT falls by exactly the discount's own VAT.
	vatCharged := taxableIncVat / (1 + vat) * vat
	vatWithoutDiscount := (taxableIncVat + discount) / (1 + vat) * vat
	if math.Abs((vatWithoutDiscount-vatCharged)-(discount/(1+vat)*vat)) > 0.001 {
		t.Errorf("expected the discount to reduce VAT by its own VAT, got %v", vatWithoutDiscount-vatCharged)
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

// iCount rejects a receipt whose payments do not cover its items, so what was collected plus what
// was withheld has to come back to the discounted total of the lines.
func TestDiscountReceiptItemBalancesTheReceipt(t *testing.T) {
	reservations := reservationSet{1: billingReservation(1, 500, "USD")}
	penalties := penaltySet{10: billingPenalty(10, 150, "USD", reservation.PenaltyTypeNoShow)}

	const (
		discount  = 20.92
		deduction = 27.44
		totalPaid = 650 - discount - deduction
	)

	items := buildReceiptItems([]int64{1}, []int64{10}, reservations, penalties)
	items = append(items, discountReceiptItem(discount, 2))

	var total float64
	for _, item := range items {
		if item.UnitPrice == nil {
			t.Fatalf("expected every receipt line to carry a unit price, got %+v", item)
		}
		total += *item.UnitPrice
	}

	if math.Abs(total-(totalPaid+deduction)) > 0.001 {
		t.Errorf("expected the lines to total %v, got %v", totalPaid+deduction, total)
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

		if _, err := newTestService().Bill(ctx, BillParams{
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

		if _, err := newTestService().Bill(ctx, BillParams{
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

		_, err := newTestService().Bill(ctx, BillParams{
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
