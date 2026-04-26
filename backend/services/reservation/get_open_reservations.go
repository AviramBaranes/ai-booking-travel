package reservation

import (
	"context"
	"fmt"

	"encore.app/internal/pricing"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type OpenReservation struct {
	ID                  int64
	BrokerReservationID string
	AgentID             int32
	CreatedAt           string
	VoucheredAt         string
	VoucherNumber       string
	DriverName          string
	PickupDate          string
	DropoffDate         string
	RentalDays          int
	CountryCode         string
	CurrencyCode        string
	CarPrice            float64
	ERPPrice            float64
	TotalPrice          float64
}

type GetOpenReservationsResponse struct {
	Reservations []OpenReservation
}

// encore:api private
func (s *Service) GetOpenReservations(ctx context.Context) (*GetOpenReservationsResponse, error) {
	rows, err := s.query.GetPaymentPendingReservations(ctx)
	if err != nil {
		rlog.Error("failed to get open reservations", "error", err)
		return nil, err
	}

	return &GetOpenReservationsResponse{
		Reservations: mapRowsToOpenReservations(rows),
	}, nil
}

// mapRowsToOpenReservations converts a slice of database rows representing payment pending reservations into a slice of OpenReservation structs.
func mapRowsToOpenReservations(rows []db.GetPaymentPendingReservationsRow) []OpenReservation {
	reservations := make([]OpenReservation, len(rows))
	for i, row := range rows {
		cp, erpP, tp := getReservationPrices(row)

		reservations[i] = OpenReservation{
			ID:                  row.ID,
			BrokerReservationID: row.BrokerReservationID,
			AgentID:             row.UserID,
			CreatedAt:           db.TimestamptzToString(row.CreatedAt),
			VoucheredAt:         db.TimestamptzToString(row.VoucheredAt),
			VoucherNumber:       ptrToString(row.VoucherNumber),
			DriverName:          fmt.Sprintf("%s %s %s", row.DriverTitle, row.DriverFirstName, row.DriverLastName),
			PickupDate:          db.DateToString(row.PickupDate),
			DropoffDate:         db.DateToString(row.ReturnDate),
			RentalDays:          int(row.RentalDays),
			CountryCode:         row.CountryCode,
			CurrencyCode:        row.CurrencyCode,
			CarPrice:            cp,
			ERPPrice:            erpP,
			TotalPrice:          tp,
		}
	}

	return reservations
}

// getReservationPrices calculates the car price, ERP price, and total price for a reservation based on the provided database row.
func getReservationPrices(row db.GetPaymentPendingReservationsRow) (carPrice, erpPrice, totalPrice float64) {
	pp := db.NumericToFloat64(row.PurchasePrice)
	mp := db.NumericToFloat64(row.MarkupPercentage)
	bErp := db.NumericToFloat64(row.BrokerErpPrice)
	btErp := float64(row.BtErpPrice)

	carFullPrice := pricing.ApplyMarkup(pp, mp)
	erpFullPrice := pricing.ApplyMarkup(bErp, mp) + btErp

	return carFullPrice, erpFullPrice, carFullPrice + erpFullPrice
}

// ptrToString is a helper function that converts a pointer to a string into a string value, returning an empty string if the pointer is nil.
func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
