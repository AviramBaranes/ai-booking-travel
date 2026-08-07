package supplier_payments

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"

	"encore.app/internal/api_errors"
	"encore.app/internal/broker"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

// Rejection reasons for a summary line that cannot be approved for payment.
const (
	// ReasonNotFound covers both "we have no such reservation" and "we already paid for it",
	// since a settled reservation is no longer part of the outstanding set.
	ReasonNotFound = "not_found_or_already_paid"
	// ReasonInvalidPrice means we found the reservation but the balance does not match what we owe.
	ReasonInvalidPrice = "invalid_price"
)

// priceTolerance absorbs float representation error only — both sides of the comparison are
// stored to two decimals, so any real difference is at least 0.01.
const priceTolerance = 0.005

var ErrInvalidPaymentSummaryFile = api_errors.NewErrorWithDetail(
	errs.InvalidArgument,
	"The uploaded file is not a valid Flex payment summary",
	api_errors.ErrorDetails{Code: api_errors.CodeInvalidPaymentSummaryFile},
)

// ApprovedReservation is a summary line matched to an outstanding reservation at the expected amount.
type ApprovedReservation struct {
	ReservationID       int64   `json:"reservationId"`
	BrokerReservationID string  `json:"brokerReservationId"`
	CurrencyCode        string  `json:"currencyCode"`
	Amount              float64 `json:"amount"`
}

// RejectedReservation is a summary line we will not pay, with the reason why.
type RejectedReservation struct {
	// ReservationID is nil when no reservation matched the summary line.
	ReservationID       *int64 `json:"reservationId,omitempty"`
	BrokerReservationID string `json:"brokerReservationId"`
	Reason              string `json:"reason"`
	// CurrencyCode and Balance are what the summary file asked for.
	CurrencyCode string  `json:"currencyCode"`
	Balance      float64 `json:"balance"`
	// ExpectedAmount and ExpectedCurrencyCode are what we have on record, set only for ReasonInvalidPrice.
	ExpectedAmount       *float64 `json:"expectedAmount,omitempty"`
	ExpectedCurrencyCode string   `json:"expectedCurrencyCode,omitempty"`
}

type ValidateFlexPaymentSummaryResponse struct {
	Approved []ApprovedReservation `json:"approved"`
	Rejected []RejectedReservation `json:"rejected"`
}

// ValidateFlexPaymentSummaryRaw reads the uploaded xlsx out of the multipart request, validates it
// and writes the result as JSON. It exists so the raw endpoint stays a thin shell around
// ValidateFlexPaymentSummary, which is where the logic lives.
func (s *SupplierPaymentsService) ValidateFlexPaymentSummaryRaw(w http.ResponseWriter, req *http.Request, file io.Reader) {
	resp, err := s.ValidateFlexPaymentSummary(req.Context(), file)
	if err != nil {
		errs.HTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		rlog.Error("failed to encode flex payment summary validation response", "error", err)
	}
}

// ValidateFlexPaymentSummary reconciles a Flex payment-required summary against the reservations
// we still owe Flex for. Every line in the file is either approved for payment or rejected with a
// reason; outstanding reservations absent from the file are left untouched.
func (s *SupplierPaymentsService) ValidateFlexPaymentSummary(ctx context.Context, file io.Reader) (*ValidateFlexPaymentSummaryResponse, error) {
	groups, err := broker.NewFlex().ReadPaymentSummary(file)
	if err != nil {
		rlog.Error("failed to read flex payment summary file", "error", err)
		return nil, ErrInvalidPaymentSummaryFile
	}

	rows, err := s.query.ListUnpaidSupplierReservations(ctx, db.BrokerFlex)
	if err != nil {
		rlog.Error("failed to list unpaid flex reservations", "error", err)
		return nil, api_errors.ErrInternalError
	}

	outstanding := make(map[string]db.ListUnpaidSupplierReservationsRow, len(rows))
	for _, r := range rows {
		outstanding[r.BrokerReservationID] = r
	}

	resp := &ValidateFlexPaymentSummaryResponse{
		Approved: []ApprovedReservation{},
		Rejected: []RejectedReservation{},
	}

	for _, group := range groups {
		for _, line := range group.Lines {
			reservation, found := outstanding[line.BookingID]
			if !found {
				resp.Rejected = append(resp.Rejected, RejectedReservation{
					BrokerReservationID: line.BookingID,
					Reason:              ReasonNotFound,
					CurrencyCode:        group.CurrencyCode,
					Balance:             line.Balance,
				})
				continue
			}

			expected := amountOwed(reservation.PurchasePrice, reservation.BrokerErpPrice)
			if reservation.CurrencyCode != group.CurrencyCode || math.Abs(expected-line.Balance) > priceTolerance {
				resp.Rejected = append(resp.Rejected, RejectedReservation{
					ReservationID:        &reservation.ID,
					BrokerReservationID:  line.BookingID,
					Reason:               ReasonInvalidPrice,
					CurrencyCode:         group.CurrencyCode,
					Balance:              line.Balance,
					ExpectedAmount:       &expected,
					ExpectedCurrencyCode: reservation.CurrencyCode,
				})
				continue
			}

			resp.Approved = append(resp.Approved, ApprovedReservation{
				ReservationID:       reservation.ID,
				BrokerReservationID: line.BookingID,
				CurrencyCode:        group.CurrencyCode,
				Amount:              line.Balance,
			})
		}
	}

	return resp, nil
}

// amountOwed is the supplier's own cost for a reservation: the car price plus the broker's ERP
// day charge, both before our markup and discount.
func amountOwed(purchasePrice, brokerErpPrice dbadapters.Numeric) float64 {
	return dbadapters.NumericToFloat64(purchasePrice) + dbadapters.NumericToFloat64(brokerErpPrice)
}
