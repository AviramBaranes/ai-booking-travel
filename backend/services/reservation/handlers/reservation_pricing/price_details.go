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
	CarProfit        float64
	ErpSellingPrice  float64
	TotalPrice       float64
}

// roundPrice rounds a price to 2 decimal places.
func roundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}

// GetReservationPriceDetails computes purchase price, selling price, profit, and ERP price from a db row.
func GetReservationPriceDetails(row db.GetPaymentPendingReservationsByBillingEntityRow) PriceDetails {
	carPurchasePrice := dbadapters.NumericToFloat64(row.PurchasePrice) + dbadapters.NumericToFloat64(row.BrokerErpPrice)
	mp := dbadapters.NumericToFloat64(row.MarkupPercentage)
	carSellingPrice := pricing.ApplyMarkup(carPurchasePrice, mp)

	return PriceDetails{
		CarPurchasePrice: roundPrice(carPurchasePrice),
		CarSellingPrice:  roundPrice(carSellingPrice),
		CarProfit:        roundPrice(carSellingPrice - carPurchasePrice),
		ErpSellingPrice:  dbadapters.NumericToFloat64(row.BtErpPrice),
		TotalPrice:       dbadapters.NumericToFloat64(row.TotalPrice),
	}
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
	ProfitOnCar         float64 `json:"profitOnCar"`
	TotalPrice          float64 `json:"totalPrice"`
	CurrencyCode        string  `json:"currencyCode"`
	CurrencyRate        float64 `json:"currencyRate"`
	CreatedAt           string  `json:"createdAt"`
	PickupDate          string  `json:"pickupDate"`
}
