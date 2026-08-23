// Package reservation_pricing carries the shapes the accountant billing workflows pass around.
// The money on them is computed by internal/pricing, which is where the reservation price
// breakdown is defined.
package reservation_pricing

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
