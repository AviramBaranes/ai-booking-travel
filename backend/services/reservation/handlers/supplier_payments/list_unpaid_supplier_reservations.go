package supplier_payments

import (
	"context"
	"fmt"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

// ListUnpaidSupplierReservationsParams selects the broker (supplier) whose outstanding
// reservations should be listed.
type ListUnpaidSupplierReservationsParams struct {
	Broker string `query:"broker" validate:"required,oneof=flex hertz"`
}

func (p *ListUnpaidSupplierReservationsParams) Validate() error {
	return validation.ValidateStruct(p)
}

// UnpaidSupplierReservation is a reservation we still owe the supplier for.
type UnpaidSupplierReservation struct {
	ID                  int64   `json:"id"`
	BrokerReservationID string  `json:"brokerReservationId"`
	DriverName          string  `json:"driverName"`
	PickupDate          string  `json:"pickupDate"`
	PickupLocationName  string  `json:"pickupLocationName"`
	RentalDays          int32   `json:"rentalDays"`
	// AmountOwed is what we owe the supplier for this reservation: the car price plus the
	// broker's ERP day charge. It is the figure reconciled against the supplier's statement.
	AmountOwed        float64 `json:"amountOwed"`
	CurrencyCode      string  `json:"currencyCode"`
	ReservationStatus string  `json:"reservationStatus"`
	PaymentStatus     string  `json:"paymentStatus"`
}

// UnpaidSupplierPenalty is a cancellation or no-show fee we still owe the supplier. The supplier
// charges the fee and we charge the customer the same amount, so a single amount covers both sides.
type UnpaidSupplierPenalty struct {
	ID                  int64   `json:"id"`
	ReservationID       int64   `json:"reservationId"`
	BrokerReservationID string  `json:"brokerReservationId"`
	Type                string  `json:"type"`
	DriverName          string  `json:"driverName"`
	PickupDate          string  `json:"pickupDate"`
	PickupLocationName  string  `json:"pickupLocationName"`
	AmountOwed          float64 `json:"amountOwed"`
	CurrencyCode        string  `json:"currencyCode"`
}

// ListUnpaidSupplierReservationsResponse lists what is still owed to the supplier, both lists
// ordered by pickup date. Each row carries its own currency code, since a supplier is settled per
// currency but the screen lists them together.
type ListUnpaidSupplierReservationsResponse struct {
	Reservations []UnpaidSupplierReservation `json:"reservations"`
	Penalties    []UnpaidSupplierPenalty     `json:"penalties"`
}

// ListUnpaidSupplierReservations returns every non-canceled reservation of the given broker
// that has not been marked as paid to the supplier yet, grouped by currency.
func (s *SupplierPaymentsService) ListUnpaidSupplierReservations(ctx context.Context, p *ListUnpaidSupplierReservationsParams) (*ListUnpaidSupplierReservationsResponse, error) {
	rows, err := s.query.ListUnpaidSupplierReservations(ctx, db.Broker(p.Broker))
	if err != nil {
		rlog.Error("failed to list unpaid supplier reservations", "error", err, "broker", p.Broker)
		return nil, api_errors.ErrInternalError
	}

	penaltyRows, err := s.query.ListUnpaidSupplierPenalties(ctx, db.Broker(p.Broker))
	if err != nil {
		rlog.Error("failed to list unpaid supplier penalties", "error", err, "broker", p.Broker)
		return nil, api_errors.ErrInternalError
	}

	return &ListUnpaidSupplierReservationsResponse{
		Reservations: toUnpaidSupplierReservations(rows),
		Penalties:    toUnpaidSupplierPenalties(penaltyRows),
	}, nil
}

func toUnpaidSupplierReservations(rows []db.ListUnpaidSupplierReservationsRow) []UnpaidSupplierReservation {
	reservations := make([]UnpaidSupplierReservation, len(rows))
	for i, r := range rows {
		reservations[i] = UnpaidSupplierReservation{
			ID:                  r.ID,
			BrokerReservationID: r.BrokerReservationID,
			DriverName:          fmt.Sprintf("%s %s %s", r.DriverTitle, r.DriverFirstName, r.DriverLastName),
			PickupDate:          dbadapters.DateToString(r.PickupDate),
			PickupLocationName:  r.PickupLocationName,
			RentalDays:          r.RentalDays,
			AmountOwed:          amountOwed(r.PurchasePrice, r.BrokerErpPrice),
			CurrencyCode:        r.CurrencyCode,
			ReservationStatus:   string(r.ReservationStatus),
			PaymentStatus:       string(r.PaymentStatus),
		}
	}

	return reservations
}

// toUnpaidSupplierPenalties maps fee rows to the response type. Every field other than the fee
// itself comes from the reservation the fee was charged on.
func toUnpaidSupplierPenalties(rows []db.ListUnpaidSupplierPenaltiesRow) []UnpaidSupplierPenalty {
	penalties := make([]UnpaidSupplierPenalty, len(rows))
	for i, p := range rows {
		penalties[i] = UnpaidSupplierPenalty{
			ID:                  p.ID,
			ReservationID:       p.ReservationID,
			BrokerReservationID: p.BrokerReservationID,
			Type:                string(p.PenaltyType),
			DriverName:          fmt.Sprintf("%s %s %s", p.DriverTitle, p.DriverFirstName, p.DriverLastName),
			PickupDate:          dbadapters.DateToString(p.PickupDate),
			PickupLocationName:  p.PickupLocationName,
			AmountOwed:          dbadapters.NumericToFloat64(p.Amount),
			CurrencyCode:        p.CurrencyCode,
		}
	}

	return penalties
}
