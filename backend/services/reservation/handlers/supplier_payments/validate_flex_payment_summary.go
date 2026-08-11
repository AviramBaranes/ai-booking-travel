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
	// ReasonNotFound means we hold no reservation at all under this booking id.
	ReasonNotFound = "not_found"
	// ReasonAlreadyPaid means we hold the reservation but have already settled it with the supplier.
	ReasonAlreadyPaid = "already_paid"
	// ReasonCanceled means the reservation was canceled, so the line is likely a cancellation or
	// no-show fee we have not recorded yet.
	ReasonCanceled = "canceled"
	// ReasonInvalidPrice means we found the reservation or fee but the balance does not match what we owe.
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

// ApprovedPenalty is a summary line matched to an outstanding fee at the expected amount.
type ApprovedPenalty struct {
	PenaltyID           int64   `json:"penaltyId"`
	ReservationID       int64   `json:"reservationId"`
	BrokerReservationID string  `json:"brokerReservationId"`
	Type                string  `json:"type"`
	CurrencyCode        string  `json:"currencyCode"`
	Amount              float64 `json:"amount"`
}

// RejectedPenalty is a summary line matched to an outstanding fee charged at a different amount.
// A line matching no fee at all is reported as a rejected reservation instead, since we cannot
// tell a fee we never recorded from a booking id we do not know.
type RejectedPenalty struct {
	PenaltyID           int64  `json:"penaltyId"`
	ReservationID       int64  `json:"reservationId"`
	BrokerReservationID string `json:"brokerReservationId"`
	Type                string `json:"type"`
	Reason              string `json:"reason"`
	// CurrencyCode and Balance are what the summary file asked for.
	CurrencyCode string  `json:"currencyCode"`
	Balance      float64 `json:"balance"`
	// ExpectedAmount and ExpectedCurrencyCode are the fee we have on record.
	ExpectedAmount       float64 `json:"expectedAmount"`
	ExpectedCurrencyCode string  `json:"expectedCurrencyCode"`
}

type ValidateFlexPaymentSummaryResponse struct {
	Approved          []ApprovedReservation `json:"approved"`
	Rejected          []RejectedReservation `json:"rejected"`
	ApprovedPenalties []ApprovedPenalty     `json:"approvedPenalties"`
	RejectedPenalties []RejectedPenalty     `json:"rejectedPenalties"`
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

	penaltyRows, err := s.query.ListUnpaidSupplierPenalties(ctx, db.BrokerFlex)
	if err != nil {
		rlog.Error("failed to list unpaid flex penalties", "error", err)
		return nil, api_errors.ErrInternalError
	}

	// A canceled reservation is excluded from the outstanding set, so a booking id resolves to
	// either a rental charge or a fee, never both.
	outstandingPenalties := make(map[string]db.ListUnpaidSupplierPenaltiesRow, len(penaltyRows))
	for _, p := range penaltyRows {
		outstandingPenalties[p.BrokerReservationID] = p
	}

	resp := &ValidateFlexPaymentSummaryResponse{
		Approved:          []ApprovedReservation{},
		Rejected:          []RejectedReservation{},
		ApprovedPenalties: []ApprovedPenalty{},
		RejectedPenalties: []RejectedPenalty{},
	}

	for _, group := range groups {
		for _, line := range group.Lines {
			if reservation, found := outstanding[line.BookingID]; found {
				addReservationResult(resp, reservation, group.CurrencyCode, line)
				continue
			}

			if penalty, found := outstandingPenalties[line.BookingID]; found {
				addPenaltyResult(resp, penalty, group.CurrencyCode, line)
				continue
			}

			resp.Rejected = append(resp.Rejected, RejectedReservation{
				BrokerReservationID: line.BookingID,
				Reason:              ReasonNotFound,
				CurrencyCode:        group.CurrencyCode,
				Balance:             line.Balance,
			})
		}
	}

	if err := s.refineNotFoundReasons(ctx, resp.Rejected); err != nil {
		return nil, err
	}

	return resp, nil
}

// addReservationResult approves a line matched to an outstanding reservation, or rejects it when
// the balance is not what we owe the supplier.
func addReservationResult(resp *ValidateFlexPaymentSummaryResponse, reservation db.ListUnpaidSupplierReservationsRow, currencyCode string, line broker.PaymentSummaryLine) {
	expected := amountOwed(reservation.PurchasePrice, reservation.BrokerErpPrice)
	if reservation.CurrencyCode != currencyCode || math.Abs(expected-line.Balance) > priceTolerance {
		resp.Rejected = append(resp.Rejected, RejectedReservation{
			ReservationID:        &reservation.ID,
			BrokerReservationID:  line.BookingID,
			Reason:               ReasonInvalidPrice,
			CurrencyCode:         currencyCode,
			Balance:              line.Balance,
			ExpectedAmount:       &expected,
			ExpectedCurrencyCode: reservation.CurrencyCode,
		})
		return
	}

	resp.Approved = append(resp.Approved, ApprovedReservation{
		ReservationID:       reservation.ID,
		BrokerReservationID: line.BookingID,
		CurrencyCode:        currencyCode,
		Amount:              line.Balance,
	})
}

// addPenaltyResult approves a line matched to an outstanding fee, or rejects it when the balance is
// not the fee we recorded. The supplier charges the fee and we charge the customer the same amount,
// so the fee itself is the figure reconciled.
func addPenaltyResult(resp *ValidateFlexPaymentSummaryResponse, penalty db.ListUnpaidSupplierPenaltiesRow, currencyCode string, line broker.PaymentSummaryLine) {
	expected := dbadapters.NumericToFloat64(penalty.Amount)
	if penalty.CurrencyCode != currencyCode || math.Abs(expected-line.Balance) > priceTolerance {
		resp.RejectedPenalties = append(resp.RejectedPenalties, RejectedPenalty{
			PenaltyID:            penalty.ID,
			ReservationID:        penalty.ReservationID,
			BrokerReservationID:  line.BookingID,
			Type:                 string(penalty.PenaltyType),
			Reason:               ReasonInvalidPrice,
			CurrencyCode:         currencyCode,
			Balance:              line.Balance,
			ExpectedAmount:       expected,
			ExpectedCurrencyCode: penalty.CurrencyCode,
		})
		return
	}

	resp.ApprovedPenalties = append(resp.ApprovedPenalties, ApprovedPenalty{
		PenaltyID:           penalty.ID,
		ReservationID:       penalty.ReservationID,
		BrokerReservationID: line.BookingID,
		Type:                string(penalty.PenaltyType),
		CurrencyCode:        currencyCode,
		Amount:              line.Balance,
	})
}

// refineNotFoundReasons replaces the blanket not-found reason with what the reservation actually
// says, for the lines that matched nothing outstanding. A booking we hold but did not offer for
// payment was either settled already or canceled, and a canceled one is a fee we have yet to
// record. Lines we genuinely know nothing about keep ReasonNotFound.
func (s *SupplierPaymentsService) refineNotFoundReasons(ctx context.Context, rejected []RejectedReservation) error {
	bookingIDs := make([]string, 0, len(rejected))
	for _, r := range rejected {
		if r.Reason == ReasonNotFound {
			bookingIDs = append(bookingIDs, r.BrokerReservationID)
		}
	}

	if len(bookingIDs) == 0 {
		return nil
	}

	rows, err := s.query.ListReservationsByBrokerReservationIDs(ctx, db.ListReservationsByBrokerReservationIDsParams{
		Broker:               db.BrokerFlex,
		BrokerReservationIds: bookingIDs,
	})
	if err != nil {
		rlog.Error("failed to look up flex reservations for rejected summary lines", "error", err)
		return api_errors.ErrInternalError
	}

	known := make(map[string]db.ListReservationsByBrokerReservationIDsRow, len(rows))
	for _, r := range rows {
		known[r.BrokerReservationID] = r
	}

	for i, r := range rejected {
		if r.Reason != ReasonNotFound {
			continue
		}

		row, found := known[r.BrokerReservationID]
		if !found {
			continue
		}

		rejected[i].ReservationID = &row.ID
		// Settlement outranks cancellation: having paid already is the more actionable fact.
		switch {
		case row.SupplierPaidAt.Valid:
			rejected[i].Reason = ReasonAlreadyPaid
		case row.ReservationStatus == db.ReservationStatusCanceled:
			rejected[i].Reason = ReasonCanceled
		}
	}

	return nil
}

// amountOwed is the supplier's own cost for a reservation: the car price plus the broker's ERP
// day charge, both before our markup and discount.
func amountOwed(purchasePrice, brokerErpPrice dbadapters.Numeric) float64 {
	return dbadapters.NumericToFloat64(purchasePrice) + dbadapters.NumericToFloat64(brokerErpPrice)
}
