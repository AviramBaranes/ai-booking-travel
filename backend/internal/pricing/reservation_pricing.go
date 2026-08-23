package pricing

import (
	"math"

	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
)

// Params are the raw money inputs every reservation price is derived from: what the broker
// charges us, the markup we sell at, our own ERP charge, and the coupon discount.
//
// Everything a reservation is worth - to the customer, to the invoice, to a profit report -
// is a pure function of these seven values, which is why Details takes nothing else.
type Params struct {
	// PurchasePrice is what the broker charges us for the car itself.
	PurchasePrice float64
	// BrokerErpPrice is what the broker charges us for their ERP day charge. Like the car,
	// it is a cost we mark up.
	BrokerErpPrice float64
	// BtErpPrice is our own ERP charge, already VAT inclusive. It has no cost behind it and
	// no markup on top of it: it is profit in full, and the discount never touches it.
	BtErpPrice float64
	// MarkupPercentage is the rate applied to the car and the broker ERP, resolved per
	// agent/office/organization at booking time.
	MarkupPercentage float64
	// DiscountPercentage is the coupon discount, applied to the marked-up car and broker ERP
	// only. Customers book with coupons; agent reservations carry 0.
	DiscountPercentage float64
	// CurrencyCode is the currency every Amount.Value is denominated in.
	CurrencyCode string
	// CurrencyRate converts that currency to shekels. A reservation with no rate yields zero
	// ILS amounts rather than an error - the currency values stay correct either way.
	CurrencyRate float64
	// StoredTotalPrice is the total_price column of a reservation that has already been written,
	// and when it is set it *is* the total price.
	//
	// It has to outrank the recomputation, because total_price is NUMERIC(12,2) and Postgres
	// settles it to the agora on write: a reservation of 100 + 15 broker ERP at 45% markup with a
	// 10% coupon and 20 of BT ERP computes to 170.07499999999999 and is stored - and billed, and
	// invoiced - as 170.08.
	//
	// Rounding the recomputed total would get back to 170.08, but what is derived from it does not
	// survive the trip: profit off the computed total is 55.074999999999989, which rounds to
	// 55.07, where profit off the column is 55.08. Go rounds the binary float, Postgres rounded
	// the exact decimal, and a cent of profit is a cent of VAT on an invoice a receipt has to
	// balance. Reading the column also keeps an older row honest, whatever formula wrote it.
	//
	// Only a reservation that has not been written yet has no column to defer to, and only that
	// one recomputes.
	StoredTotalPrice *float64
}

// New returns the pricing params for values that are already plain floats: a plan being priced
// off a snapshot, or a reservation about to be created.
//
// It exists so that every caller reaches Details through the same door as NewWithReservation,
// and so any future normalisation of the inputs has a single place to live.
func New(p Params) Params {
	return p
}

// NewWithReservation returns the pricing params of a stored reservation, priced off its own
// total_price column - see StoredTotalPrice for why the column outranks the arithmetic.
func NewWithReservation(r db.Reservation) Params {
	return New(Params{
		PurchasePrice:      dbadapters.NumericToFloat64(r.PurchasePrice),
		BrokerErpPrice:     dbadapters.NumericToFloat64(r.BrokerErpPrice),
		BtErpPrice:         dbadapters.NumericToFloat64(r.BtErpPrice),
		MarkupPercentage:   dbadapters.NumericToFloat64(r.MarkupPercentage),
		DiscountPercentage: dbadapters.NumericToFloat64(r.DiscountPercentage),
		CurrencyCode:       r.CurrencyCode,
		CurrencyRate:       dbadapters.NumericToFloat64(r.CurrencyRate),
		StoredTotalPrice:   dbadapters.NumericToFloat64Ptr(r.TotalPrice),
	})
}

// NumericParams is Params as it comes out of Postgres. It exists for the query row shapes that
// select the price columns without selecting the whole reservation - dashboard and billing rows -
// so they do not each repeat the same six conversions.
type NumericParams struct {
	PurchasePrice      dbadapters.Numeric
	BrokerErpPrice     dbadapters.Numeric
	BtErpPrice         dbadapters.Numeric
	MarkupPercentage   dbadapters.Numeric
	DiscountPercentage dbadapters.Numeric
	CurrencyCode       string
	CurrencyRate       dbadapters.Numeric
	// TotalPrice is the stored total_price column. Select it whenever the query has it - a NULL
	// numeric falls back to recomputing, which is a worse answer. See StoredTotalPrice.
	TotalPrice dbadapters.Numeric
}

// NewWithNumerics returns the pricing params of a partial reservation row.
func NewWithNumerics(p NumericParams) Params {
	return New(Params{
		PurchasePrice:      dbadapters.NumericToFloat64(p.PurchasePrice),
		BrokerErpPrice:     dbadapters.NumericToFloat64(p.BrokerErpPrice),
		BtErpPrice:         dbadapters.NumericToFloat64(p.BtErpPrice),
		MarkupPercentage:   dbadapters.NumericToFloat64(p.MarkupPercentage),
		DiscountPercentage: dbadapters.NumericToFloat64(p.DiscountPercentage),
		CurrencyCode:       p.CurrencyCode,
		CurrencyRate:       dbadapters.NumericToFloat64(p.CurrencyRate),
		StoredTotalPrice:   dbadapters.NumericToFloat64Ptr(p.TotalPrice),
	})
}

// Amount is one money value of the breakdown, in the reservation's own currency and in shekels.
//
// Value and ILS are exact: rounding money mid-calculation loses agorot the next step would have
// kept, and total_price is stored to six decimals - so anything that feeds further arithmetic,
// the database included, reads these. ValueRounded and ILSRounded are the same numbers settled to
// the agora, for whatever displays or bills them.
type Amount struct {
	Value        float64 `json:"value"`
	ILS          float64 `json:"ils"`
	ValueRounded float64 `json:"valueRounded"`
	ILSRounded   float64 `json:"ilsRounded"`
}

// Details is the complete money breakdown of a single reservation - cost, markup, discount and
// profit - and is the only place these numbers are defined.
//
// The chain runs: car cost and broker ERP cost are marked up, the discount comes off both, and
// our own BT ERP charge is added after the discount. Profit is whatever is left of the price the
// customer actually pays once the broker has been paid, so a discount reduces profit rather than
// revenue-with-a-footnote.
type Details struct {
	CurrencyCode       string  `json:"currencyCode"`
	CurrencyRate       float64 `json:"currencyRate"`
	MarkupPercentage   float64 `json:"markupPercentage"`
	DiscountPercentage float64 `json:"discountPercentage"`

	// What we owe the broker.
	CarCost       Amount `json:"carCost"`
	BrokerErpCost Amount `json:"brokerErpCost"`
	TotalCost     Amount `json:"totalCost"`

	// Our own ERP charge, which has no cost behind it.
	BtErpPrice Amount `json:"btErpPrice"`

	// What we sell for, before the discount. ErpFullPrice is the two ERP charges together as the
	// customer sees them on one line - the broker's marked up, ours as is.
	CarWithMarkup              Amount `json:"carWithMarkup"`
	BrokerErpWithMarkup        Amount `json:"brokerErpWithMarkup"`
	ErpFullPrice               Amount `json:"erpFullPrice"`
	CarWithBrokerErpWithMarkup Amount `json:"carWithBrokerErpWithMarkup"`
	PriceBeforeDiscount        Amount `json:"priceBeforeDiscount"`

	// What the customer pays. TotalPrice is PriceBeforeDiscount less TotalDiscount, and is the
	// number stored on the reservation.
	TotalDiscount Amount `json:"totalDiscount"`
	TotalPrice    Amount `json:"totalPrice"`

	// What we keep. CarProfit and ErpProfit are both net of the discount and add up to
	// TotalProfit, which is also TotalPrice less TotalCost.
	CarProfit   Amount `json:"carProfit"`
	ErpProfit   Amount `json:"erpProfit"`
	TotalProfit Amount `json:"totalProfit"`
}

// Details computes the full breakdown of the reservation these params describe.
func (p Params) Details() Details {
	carWithMarkup := ApplyMarkup(p.PurchasePrice, p.MarkupPercentage)
	brokerErpWithMarkup := ApplyMarkup(p.BrokerErpPrice, p.MarkupPercentage)

	// The car and the broker ERP are discounted at the same rate, so the discount can be split
	// between them and the car's profit stated on its own.
	carAfterDiscount := CalculateDiscountedPrice(carWithMarkup, p.DiscountPercentage)

	carWithBrokerErpWithMarkup := carWithMarkup + brokerErpWithMarkup
	priceBeforeDiscount := carWithBrokerErpWithMarkup + p.BtErpPrice

	// Discounted as one sum rather than as the two halves above, so that this is the same
	// float64, bit for bit, as the CalculateTotalPrice that writes a reservation's total_price.
	// Once that column exists it is the price, rounded to the agora by Postgres - see
	// StoredTotalPrice.
	totalPrice := CalculateDiscountedPrice(carWithBrokerErpWithMarkup, p.DiscountPercentage) + p.BtErpPrice
	if p.StoredTotalPrice != nil {
		totalPrice = *p.StoredTotalPrice
	}

	// Grouped as total less cost, the same way the billing breakdown has always computed it.
	totalProfit := totalPrice - (p.PurchasePrice + p.BrokerErpPrice)
	carProfit := carAfterDiscount - p.PurchasePrice

	rate := p.CurrencyRate

	return Details{
		CurrencyCode:       p.CurrencyCode,
		CurrencyRate:       rate,
		MarkupPercentage:   p.MarkupPercentage,
		DiscountPercentage: p.DiscountPercentage,

		CarCost:       newAmount(p.PurchasePrice, rate),
		BrokerErpCost: newAmount(p.BrokerErpPrice, rate),
		TotalCost:     newAmount(p.PurchasePrice+p.BrokerErpPrice, rate),

		BtErpPrice: newAmount(p.BtErpPrice, rate),

		CarWithMarkup:              newAmount(carWithMarkup, rate),
		BrokerErpWithMarkup:        newAmount(brokerErpWithMarkup, rate),
		ErpFullPrice:               newAmount(brokerErpWithMarkup+p.BtErpPrice, rate),
		CarWithBrokerErpWithMarkup: newAmount(carWithBrokerErpWithMarkup, rate),
		PriceBeforeDiscount:        newAmount(priceBeforeDiscount, rate),

		TotalDiscount: newAmount(priceBeforeDiscount-totalPrice, rate),
		TotalPrice:    newAmount(totalPrice, rate),

		CarProfit: newAmount(carProfit, rate),
		// The ERP side takes the residual rather than being derived on its own, so that the two
		// halves always add up to the whole: when the total came from the column it carries
		// Postgres's rounding, and that half-agora has to land somewhere.
		ErpProfit:   newAmount(totalProfit-carProfit, rate),
		TotalProfit: newAmount(totalProfit, rate),
	}
}

// newAmount states one value in both currencies, exactly and rounded to the agora.
func newAmount(value, currencyRate float64) Amount {
	ils := value * currencyRate
	return Amount{
		Value:        value,
		ILS:          ils,
		ValueRounded: roundMoney(value),
		ILSRounded:   roundMoney(ils),
	}
}

// roundMoney rounds a price to 2 decimal places.
func roundMoney(price float64) float64 {
	return math.Round(price*100) / 100
}
