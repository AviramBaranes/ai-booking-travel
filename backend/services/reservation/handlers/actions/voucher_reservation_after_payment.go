package actions

import (
	"context"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
	"encore.app/services/reservation/handlers/reservation_pricing"
	"encore.dev/rlog"
)

type VoucherReservationAfterPaymentParams struct {
	ReservationID           int64
	TransactionRate         float64
	PaymentConfirmationCode string
	PaymentDocNum           string
	UserEmail               string
}

type VoucherReservationAfterPaymentResponse struct {
	BillingReservation reservation_pricing.BillingReservation
}

func (s *ActionService) VoucherReservationAfterPayment(ctx context.Context, p *VoucherReservationAfterPaymentParams) (*VoucherReservationAfterPaymentResponse, error) {
	reservation, err := s.query.VoucherReservationAfterPayment(ctx, db.VoucherReservationAfterPaymentParams{
		ID:                      p.ReservationID,
		PaymentConfirmationCode: &p.PaymentConfirmationCode,
		PaymentDocNum:           &p.PaymentDocNum,
	})
	if err != nil {
		return nil, api_errors.ErrInternalError
	}

	s.updateReservationCurrencyRate(ctx, reservation.ID, reservation.CurrencyCode)
	sendReservationVoucherToUser(ctx, reservation, p.UserEmail, "N/A")

	pd := reservation_pricing.GetReservationPriceDetails(db.GetPaymentPendingReservationsByBillingEntityRow{
		ID:                  reservation.ID,
		BrokerReservationID: reservation.BrokerReservationID,
		PaymentStatus:       reservation.PaymentStatus,
		ReservationStatus:   reservation.ReservationStatus,
		PurchasePrice:       reservation.PurchasePrice,
		MarkupPercentage:    reservation.MarkupPercentage,
		BtErpPrice:          reservation.BtErpPrice,
		BrokerErpPrice:      reservation.BrokerErpPrice,
		TotalPrice:          reservation.TotalPrice,
		CurrencyCode:        reservation.CurrencyCode,
		CreatedAt:           reservation.CreatedAt,
		PickupDate:          reservation.PickupDate,
	})

	return &VoucherReservationAfterPaymentResponse{
		BillingReservation: reservation_pricing.BillingReservation{
			ID:                  reservation.ID,
			BrokerReservationID: reservation.BrokerReservationID,
			PaymentStatus:       string(reservation.PaymentStatus),
			ReservationStatus:   string(reservation.ReservationStatus),
			CarPurchasePrice:    pd.CarPurchasePrice,
			CarSellingPrice:     pd.CarSellingPrice,
			ERPSellingPrice:     pd.ErpSellingPrice,
			ProfitOnCar:         pd.CarProfit,
			TotalPrice:          pd.TotalPrice,
			CurrencyCode:        reservation.CurrencyCode,
			CurrencyRate:        dbadapters.NumericToFloat64(reservation.CurrencyRate),
			CreatedAt:           dbadapters.TimestamptzToString(reservation.CreatedAt),
			PickupDate:          dbadapters.DateToString(reservation.PickupDate),
		}}, nil

}

func (s *ActionService) updateReservationCurrencyRate(ctx context.Context, reservationID int64, currencyCode string) {
	rate, err := s.currencyCache.Get(ctx, currencyCode)
	if err != nil {
		rlog.Error("failed to get currency rate", "error", err, "currency_code", currencyCode)
		return
	}

	err = s.query.UpdateReservationCurrencyRate(ctx, db.UpdateReservationCurrencyRateParams{
		ID:           reservationID,
		CurrencyRate: dbadapters.NumericFromFloat64(rate),
	})
	if err != nil {
		rlog.Error("failed to update reservation currency rate", "error", err, "reservationID", reservationID)
	}
}
