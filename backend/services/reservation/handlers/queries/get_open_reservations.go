package queries

import (
	"context"
	"fmt"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/pricing"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type OpenReservation struct {
	ID                  int64
	PaymentStatus       string
	BrokerReservationID string
	AgentID             int64
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

func (s *QueryService) GetOpenReservations(ctx context.Context) (*GetOpenReservationsResponse, error) {
	rows, err := s.query.GetPaymentPendingReservations(ctx)
	if err != nil {
		rlog.Error("failed to get open reservations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	return &GetOpenReservationsResponse{
		Reservations: mapRowsToOpenReservations(rows),
	}, nil
}

func mapRowsToOpenReservations(rows []db.GetPaymentPendingReservationsRow) []OpenReservation {
	reservations := make([]OpenReservation, len(rows))
	for i, row := range rows {
		cp, erpP, tp := getReservationPrices(row)

		reservations[i] = OpenReservation{
			ID:                  row.ID,
			PaymentStatus:       string(row.PaymentStatus),
			BrokerReservationID: row.BrokerReservationID,
			AgentID:             row.UserID,
			CreatedAt:           dbadapters.TimestamptzToString(row.CreatedAt),
			VoucheredAt:         dbadapters.TimestamptzToString(row.VoucheredAt),
			VoucherNumber:       ptrToString(row.VoucherNumber),
			DriverName:          fmt.Sprintf("%s %s %s", row.DriverTitle, row.DriverFirstName, row.DriverLastName),
			PickupDate:          dbadapters.DateToString(row.PickupDate),
			DropoffDate:         dbadapters.DateToString(row.DropoffDate),
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

func getReservationPrices(row db.GetPaymentPendingReservationsRow) (carPrice, erpPrice, totalPrice float64) {
	pp := dbadapters.NumericToFloat64(row.PurchasePrice)
	mp := dbadapters.NumericToFloat64(row.MarkupPercentage)
	bErp := dbadapters.NumericToFloat64(row.BrokerErpPrice)
	btErp := dbadapters.NumericToFloat64(row.BtErpPrice)

	carFullPrice := pricing.ApplyMarkup(pp, mp)
	erpFullPrice := pricing.ApplyMarkup(bErp, mp) + btErp

	return carFullPrice, erpFullPrice, carFullPrice + erpFullPrice
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
