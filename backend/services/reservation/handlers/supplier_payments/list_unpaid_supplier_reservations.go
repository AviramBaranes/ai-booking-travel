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
	ReservationStatus string  `json:"reservationStatus"`
	PaymentStatus       string  `json:"paymentStatus"`
}

// UnpaidSupplierCurrencyGroup is a set of unpaid reservations sharing the same currency.
// Suppliers are settled per currency, so the payment summary is reconciled group by group.
type UnpaidSupplierCurrencyGroup struct {
	CurrencyCode string                      `json:"currencyCode"`
	Reservations []UnpaidSupplierReservation `json:"reservations"`
}

type ListUnpaidSupplierReservationsResponse struct {
	CurrencyGroups []UnpaidSupplierCurrencyGroup `json:"currencyGroups"`
}

// ListUnpaidSupplierReservations returns every non-canceled reservation of the given broker
// that has not been marked as paid to the supplier yet, grouped by currency.
func (s *SupplierPaymentsService) ListUnpaidSupplierReservations(ctx context.Context, p *ListUnpaidSupplierReservationsParams) (*ListUnpaidSupplierReservationsResponse, error) {
	rows, err := s.query.ListUnpaidSupplierReservations(ctx, db.Broker(p.Broker))
	if err != nil {
		rlog.Error("failed to list unpaid supplier reservations", "error", err, "broker", p.Broker)
		return nil, api_errors.ErrInternalError
	}

	return &ListUnpaidSupplierReservationsResponse{
		CurrencyGroups: toCurrencyGroups(rows),
	}, nil
}

// toCurrencyGroups splits the rows into per-currency groups. The query orders by currency_code,
// so rows of the same currency are contiguous and a new group starts on every change.
func toCurrencyGroups(rows []db.ListUnpaidSupplierReservationsRow) []UnpaidSupplierCurrencyGroup {
	groups := []UnpaidSupplierCurrencyGroup{}
	for _, r := range rows {
		if len(groups) == 0 || groups[len(groups)-1].CurrencyCode != r.CurrencyCode {
			groups = append(groups, UnpaidSupplierCurrencyGroup{
				CurrencyCode: r.CurrencyCode,
				Reservations: []UnpaidSupplierReservation{},
			})
		}

		group := &groups[len(groups)-1]
		group.Reservations = append(group.Reservations, UnpaidSupplierReservation{
			ID:                  r.ID,
			BrokerReservationID: r.BrokerReservationID,
			DriverName:          fmt.Sprintf("%s %s %s", r.DriverTitle, r.DriverFirstName, r.DriverLastName),
			PickupDate:          dbadapters.DateToString(r.PickupDate),
			PickupLocationName:  r.PickupLocationName,
			RentalDays:          r.RentalDays,
			AmountOwed:          amountOwed(r.PurchasePrice, r.BrokerErpPrice),
			ReservationStatus:   string(r.ReservationStatus),
			PaymentStatus:       string(r.PaymentStatus),
		})
	}
	return groups
}
