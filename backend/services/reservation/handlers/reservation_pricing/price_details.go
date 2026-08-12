package reservation_pricing

import (
	"math"

	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	"encore.app/services/reservation/db"
)

// PriceDetails holds the computed price breakdown for a single reservation.
type PriceDetails struct {
	CarPurchasePrice float64
	CarSellingPrice  float64
	TotalProfit      float64
	ErpSellingPrice  float64
	TotalPrice       float64
}

// roundPrice rounds a price to 2 decimal places.
func roundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}

// PriceInputs are the raw reservation price columns, already converted out of pgtype.
// It exists so the price breakdown can be computed from any query's row shape.
type PriceInputs struct {
	PurchasePrice    float64
	BrokerErpPrice   float64
	MarkupPercentage float64
	BtErpPrice       float64
	TotalPrice       float64
}

// ComputePriceDetails is the single place the reservation price breakdown is defined.
//
// Cost is what we owe the supplier (the car plus the broker's ERP day charge), and profit
// is what is left of the price the customer actually pays after that cost — so any discount
// already baked into total_price reduces the profit, as it should.
func ComputePriceDetails(in PriceInputs) PriceDetails {
	purchasePrice := in.PurchasePrice + in.BrokerErpPrice
	carSellingPrice := pricing.ApplyMarkup(purchasePrice, in.MarkupPercentage)
	profit := in.TotalPrice - purchasePrice

	return PriceDetails{
		CarPurchasePrice: roundPrice(purchasePrice),
		CarSellingPrice:  roundPrice(carSellingPrice),
		TotalProfit:      roundPrice(profit),
		ErpSellingPrice:  in.BtErpPrice,
		TotalPrice:       in.TotalPrice,
	}
}

// GetReservationPriceDetails computes purchase price, selling price, profit, and ERP price from a db row.
func GetReservationPriceDetails(row db.GetPaymentPendingReservationsByBillingEntityRow) PriceDetails {
	return ComputePriceDetails(PriceInputs{
		PurchasePrice:    dbadapters.NumericToFloat64(row.PurchasePrice),
		BrokerErpPrice:   dbadapters.NumericToFloat64(row.BrokerErpPrice),
		MarkupPercentage: dbadapters.NumericToFloat64(row.MarkupPercentage),
		BtErpPrice:       dbadapters.NumericToFloat64(row.BtErpPrice),
		TotalPrice:       dbadapters.NumericToFloat64(row.TotalPrice),
	})
}

// BillingReservation is a reservation summary tailored for accountant billing workflows.
type BillingReservation struct {
	ID                  int64   `json:"id"`
	BrokerReservationID string  `json:"brokerReservationId"`
	PaymentStatus       string  `json:"paymentStatus"`
	ReservationStatus   string  `json:"reservationStatus"`
	CarPurchasePrice    float64 `json:"carPurchasePrice"`
	CarSellingPrice     float64 `json:"carSellingPrice"`
	ERPSellingPrice     float64 `json:"erpSellingPrice"`
	TotalProfit         float64 `json:"totalProfit"`
	TotalPrice          float64 `json:"totalPrice"`
	CurrencyCode        string  `json:"currencyCode"`
	CurrencyRate        float64 `json:"currencyRate"`
	CreatedAt           string  `json:"createdAt"`
	PickupDate          string  `json:"pickupDate"`
	VoucheredAt         string  `json:"voucheredAt"`
}

// BillingPenalty is a cancellation or no-show fee tailored for accountant billing workflows.
// It carries no price breakdown: the supplier's charge is passed on to the customer as is.
type BillingPenalty struct {
	ID                  int64   `json:"id"`
	ReservationID       int64   `json:"reservationId"`
	BrokerReservationID string  `json:"brokerReservationId"`
	Type                string  `json:"type"`
	Amount              float64 `json:"amount"`
	CurrencyCode        string  `json:"currencyCode"`
	CurrencyRate        float64 `json:"currencyRate"`
	CreatedAt           string  `json:"createdAt"`
}
