package pricing

import (
	"math"
	"testing"

	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
)

func TestDetails(t *testing.T) {
	t.Run("car only", func(t *testing.T) {
		// An agent reservation: no coupon, no ERP on either side.
		const rate = 4.0
		d := New(Params{
			PurchasePrice:    100,
			MarkupPercentage: 20,
			CurrencyCode:     "EUR",
			CurrencyRate:     rate,
		}).Details()

		if d.CurrencyCode != "EUR" || d.CurrencyRate != rate || d.MarkupPercentage != 20 || d.DiscountPercentage != 0 {
			t.Errorf("params not carried through: %+v", d)
		}

		assertAmount(t, "CarCost", d.CarCost, 100, rate)
		assertAmount(t, "BrokerErpCost", d.BrokerErpCost, 0, rate)
		assertAmount(t, "TotalCost", d.TotalCost, 100, rate)
		assertAmount(t, "BtErpPrice", d.BtErpPrice, 0, rate)
		assertAmount(t, "CarWithMarkup", d.CarWithMarkup, 120, rate)
		assertAmount(t, "BrokerErpWithMarkup", d.BrokerErpWithMarkup, 0, rate)
		assertAmount(t, "ErpFullPrice", d.ErpFullPrice, 0, rate)
		assertAmount(t, "CarWithBrokerErpWithMarkup", d.CarWithBrokerErpWithMarkup, 120, rate)
		assertAmount(t, "PriceBeforeDiscount", d.PriceBeforeDiscount, 120, rate)
		assertAmount(t, "TotalDiscount", d.TotalDiscount, 0, rate)
		assertAmount(t, "TotalPrice", d.TotalPrice, 120, rate)
		assertAmount(t, "CarProfit", d.CarProfit, 20, rate)
		assertAmount(t, "ErpProfit", d.ErpProfit, 0, rate)
		assertAmount(t, "TotalProfit", d.TotalProfit, 20, rate)
	})

	t.Run("both erps and a discount", func(t *testing.T) {
		// A customer reservation: both ERP charges and a coupon. The discount comes off the
		// marked-up car and broker ERP; the BT ERP charge is added after it.
		const rate = 4.0
		d := New(Params{
			PurchasePrice:      200,
			BrokerErpPrice:     20,
			BtErpPrice:         10,
			MarkupPercentage:   50,
			DiscountPercentage: 20,
			CurrencyCode:       "USD",
			CurrencyRate:       rate,
		}).Details()

		assertAmount(t, "CarCost", d.CarCost, 200, rate)
		assertAmount(t, "BrokerErpCost", d.BrokerErpCost, 20, rate)
		assertAmount(t, "TotalCost", d.TotalCost, 220, rate)
		assertAmount(t, "BtErpPrice", d.BtErpPrice, 10, rate)
		assertAmount(t, "CarWithMarkup", d.CarWithMarkup, 300, rate)
		assertAmount(t, "BrokerErpWithMarkup", d.BrokerErpWithMarkup, 30, rate)
		assertAmount(t, "ErpFullPrice", d.ErpFullPrice, 40, rate)
		assertAmount(t, "CarWithBrokerErpWithMarkup", d.CarWithBrokerErpWithMarkup, 330, rate)
		assertAmount(t, "PriceBeforeDiscount", d.PriceBeforeDiscount, 340, rate)
		assertAmount(t, "TotalDiscount", d.TotalDiscount, 66, rate)
		assertAmount(t, "TotalPrice", d.TotalPrice, 274, rate)
		assertAmount(t, "CarProfit", d.CarProfit, 40, rate)
		assertAmount(t, "ErpProfit", d.ErpProfit, 14, rate)
		assertAmount(t, "TotalProfit", d.TotalProfit, 54, rate)
	})

	t.Run("no currency rate", func(t *testing.T) {
		// A reservation booked before an exchange rate was captured still prices correctly in its
		// own currency; only the shekel side is empty.
		d := New(Params{
			PurchasePrice:    100,
			MarkupPercentage: 10,
			CurrencyCode:     "GBP",
		}).Details()

		assertAmount(t, "CarCost", d.CarCost, 100, 0)
		assertAmount(t, "CarWithMarkup", d.CarWithMarkup, 110, 0)
		assertAmount(t, "TotalPrice", d.TotalPrice, 110, 0)
		assertAmount(t, "TotalProfit", d.TotalProfit, 10, 0)
	})

	t.Run("rounds to the agora alongside the exact value", func(t *testing.T) {
		// 100 * 1.1 is not 110 in float64, which is exactly why the rounded fields exist.
		d := New(Params{PurchasePrice: 100, MarkupPercentage: 10, CurrencyRate: 3}).Details()

		if d.CarWithMarkup.ValueRounded != 110 {
			t.Errorf("CarWithMarkup.ValueRounded = %v, want 110", d.CarWithMarkup.ValueRounded)
		}
		if d.CarWithMarkup.ILSRounded != 330 {
			t.Errorf("CarWithMarkup.ILSRounded = %v, want 330", d.CarWithMarkup.ILSRounded)
		}

		// A half-agora case: 100 with a 12.345% markup is 112.345.
		d = New(Params{PurchasePrice: 100, MarkupPercentage: 12.345, CurrencyRate: 1}).Details()
		if d.CarWithMarkup.ValueRounded != 112.35 {
			t.Errorf("CarWithMarkup.ValueRounded = %v, want 112.35", d.CarWithMarkup.ValueRounded)
		}
		if !closeEnough(d.CarWithMarkup.Value, 112.345) {
			t.Errorf("CarWithMarkup.Value = %v, want 112.345", d.CarWithMarkup.Value)
		}
	})
}

// TestDetailsMatchesCalculateTotalPrice guards the reason the reservation's total_price column can
// be recomputed rather than read: the two have to be the same float64, bit for bit.
func TestDetailsMatchesCalculateTotalPrice(t *testing.T) {
	params := []Params{
		{PurchasePrice: 100, MarkupPercentage: 20},
		{PurchasePrice: 173.45, BrokerErpPrice: 15.9, BtErpPrice: 12, MarkupPercentage: 45, DiscountPercentage: 7.5},
		{PurchasePrice: 1234.56, BrokerErpPrice: 78.9, BtErpPrice: 33.33, MarkupPercentage: 12.5, DiscountPercentage: 100},
	}

	for _, p := range params {
		want := CalculateTotalPrice(p.PurchasePrice, p.MarkupPercentage, p.BrokerErpPrice, p.BtErpPrice, p.DiscountPercentage)
		if got := New(p).Details().TotalPrice.Value; got != want {
			t.Errorf("TotalPrice for %+v = %.17g, CalculateTotalPrice = %.17g", p, got, want)
		}
	}
}

// TestStoredTotalPriceWins pins the reason the column outranks the arithmetic. The fixture is the
// one from the billing test: it computes to 170.07499999999999 and Postgres stores it, in a
// NUMERIC(12,2) column, as 170.08 - and the cent that separates them lands on the invoice.
func TestStoredTotalPriceWins(t *testing.T) {
	params := Params{
		PurchasePrice:      100,
		BrokerErpPrice:     15,
		BtErpPrice:         20,
		MarkupPercentage:   45,
		DiscountPercentage: 10,
		CurrencyRate:       1,
	}

	recomputed := New(params).Details()
	if !closeEnough(recomputed.TotalPrice.Value, 170.075) {
		t.Fatalf("recomputed TotalPrice = %.17g, want 170.075", recomputed.TotalPrice.Value)
	}
	if recomputed.TotalProfit.ValueRounded != 55.07 {
		t.Errorf("recomputed profit settles to %v, want 55.07 - the fixture this test exists for", recomputed.TotalProfit.ValueRounded)
	}

	stored := 170.08
	params.StoredTotalPrice = &stored
	d := New(params).Details()

	if d.TotalPrice.Value != stored {
		t.Errorf("TotalPrice = %v, want the stored %v", d.TotalPrice.Value, stored)
	}
	if d.TotalProfit.ValueRounded != 55.08 {
		t.Errorf("TotalProfit = %v, want 55.08 - the number on the invoice", d.TotalProfit.ValueRounded)
	}
	if !closeEnough(d.CarProfit.Value+d.ErpProfit.Value, d.TotalProfit.Value) {
		t.Errorf("profit halves %v + %v do not add up to %v", d.CarProfit.Value, d.ErpProfit.Value, d.TotalProfit.Value)
	}
}

// TestDetailsIdentities checks the relations the breakdown promises, so a future change to one
// field cannot quietly stop agreeing with the others.
func TestDetailsIdentities(t *testing.T) {
	d := New(Params{
		PurchasePrice:      173.45,
		BrokerErpPrice:     15.9,
		BtErpPrice:         12,
		MarkupPercentage:   45,
		DiscountPercentage: 7.5,
		CurrencyRate:       3.87,
	}).Details()

	checks := []struct {
		name      string
		got, want float64
	}{
		{"car + broker erp with markup", d.CarWithMarkup.Value + d.BrokerErpWithMarkup.Value, d.CarWithBrokerErpWithMarkup.Value},
		{"erp full price", d.BrokerErpWithMarkup.Value + d.BtErpPrice.Value, d.ErpFullPrice.Value},
		{"price before discount", d.CarWithBrokerErpWithMarkup.Value + d.BtErpPrice.Value, d.PriceBeforeDiscount.Value},
		{"total price after discount", d.PriceBeforeDiscount.Value - d.TotalDiscount.Value, d.TotalPrice.Value},
		{"total cost", d.CarCost.Value + d.BrokerErpCost.Value, d.TotalCost.Value},
		{"total profit from parts", d.CarProfit.Value + d.ErpProfit.Value, d.TotalProfit.Value},
		{"total profit from price", d.TotalPrice.Value - d.TotalCost.Value, d.TotalProfit.Value},
		{"ils follows the rate", d.TotalPrice.Value * 3.87, d.TotalPrice.ILS},
	}

	for _, c := range checks {
		if !closeEnough(c.got, c.want) {
			t.Errorf("%s: %.6f != %.6f", c.name, c.got, c.want)
		}
	}
}

func TestNewWithReservation(t *testing.T) {
	reservation := db.Reservation{
		PurchasePrice:      dbadapters.NumericFromFloat64(200),
		BrokerErpPrice:     dbadapters.NumericFromFloat64(20),
		BtErpPrice:         dbadapters.NumericFromFloat64(10),
		MarkupPercentage:   dbadapters.NumericFromFloat64(50),
		DiscountPercentage: dbadapters.NumericFromFloat64(20),
		CurrencyCode:       "USD",
		CurrencyRate:       dbadapters.NumericFromFloat64(4),
	}

	want := New(Params{
		PurchasePrice:      200,
		BrokerErpPrice:     20,
		BtErpPrice:         10,
		MarkupPercentage:   50,
		DiscountPercentage: 20,
		CurrencyCode:       "USD",
		CurrencyRate:       4,
	})

	if got := NewWithReservation(reservation); got != want {
		t.Errorf("NewWithReservation() = %+v, want %+v", got, want)
	}
}

func TestNewWithNumerics(t *testing.T) {
	want := New(Params{
		PurchasePrice:      173.45,
		BrokerErpPrice:     15.9,
		BtErpPrice:         12,
		MarkupPercentage:   45,
		DiscountPercentage: 7.5,
		CurrencyCode:       "EUR",
		CurrencyRate:       3.87,
	})

	got := NewWithNumerics(NumericParams{
		PurchasePrice:      dbadapters.NumericFromFloat64(173.45),
		BrokerErpPrice:     dbadapters.NumericFromFloat64(15.9),
		BtErpPrice:         dbadapters.NumericFromFloat64(12),
		MarkupPercentage:   dbadapters.NumericFromFloat64(45),
		DiscountPercentage: dbadapters.NumericFromFloat64(7.5),
		CurrencyCode:       "EUR",
		CurrencyRate:       dbadapters.NumericFromFloat64(3.87),
	})

	if got != want {
		t.Errorf("NewWithNumerics() = %+v, want %+v", got, want)
	}
}

// assertAmount checks all four faces of one amount against the value it should hold: the exact
// value, its shekel conversion, and both settled to the agora.
func assertAmount(t *testing.T, name string, got Amount, want, currencyRate float64) {
	t.Helper()

	if !closeEnough(got.Value, want) {
		t.Errorf("%s.Value = %v, want %v", name, got.Value, want)
	}
	if !closeEnough(got.ILS, want*currencyRate) {
		t.Errorf("%s.ILS = %v, want %v", name, got.ILS, want*currencyRate)
	}
	if got.ValueRounded != roundMoney(want) {
		t.Errorf("%s.ValueRounded = %v, want %v", name, got.ValueRounded, roundMoney(want))
	}
	if got.ILSRounded != roundMoney(want*currencyRate) {
		t.Errorf("%s.ILSRounded = %v, want %v", name, got.ILSRounded, roundMoney(want*currencyRate))
	}
}

// closeEnough compares prices that float64 cannot hold exactly - 100 marked up 10% is
// 110.00000000000001 - well below the agora the rounded fields settle to.
func closeEnough(got, want float64) bool {
	return math.Abs(got-want) <= 1e-9
}
