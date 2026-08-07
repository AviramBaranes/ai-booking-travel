package supplier_payments

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/icount"
	"encore.app/internal/validation"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

// supplierExpenseDocType is the iCount document type supplier statements are recorded as.
const supplierExpenseDocType = "invoice"

// Failure reasons for a reservation we could not settle.
const (
	ReasonUnsupportedBroker     = "unsupported_broker"
	ReasonExpenseCreationFailed = "expense_creation_failed"
	// ReasonMarkPaidFailed means the expense exists in iCount but we failed to link it to the
	// reservation, so the reservation is still outstanding and must be reconciled by hand.
	ReasonMarkPaidFailed = "mark_paid_failed"
)

type PaySupplierReservationsParams struct {
	ReservationIDs []int64 `json:"reservationIds" validate:"required,min=1"`
	// PaymentDate is the day the supplier was actually paid. It is recorded both on the iCount
	// expense and as the reservation's supplier_paid_at.
	PaymentDate string `json:"paymentDate" validate:"required,datetime=2006-01-02"`
}

func (p PaySupplierReservationsParams) Validate() error {
	return validation.ValidateStruct(p)
}

// PaidSupplierReservation is a reservation now settled and linked to its iCount expense.
type PaidSupplierReservation struct {
	ReservationID int64   `json:"reservationId"`
	ExpenseID     string  `json:"expenseId"`
	Amount        float64 `json:"amount"`
	CurrencyCode  string  `json:"currencyCode"`
}

// FailedSupplierPayment is a reservation that could not be settled.
type FailedSupplierPayment struct {
	ReservationID int64  `json:"reservationId"`
	Reason        string `json:"reason"`
	// ExpenseID is set when the expense was created but could not be linked to the reservation.
	ExpenseID string `json:"expenseId,omitempty"`
}

type PaySupplierReservationsResponse struct {
	Paid []PaidSupplierReservation `json:"paid"`
	// Skipped holds the requested ids that are no longer outstanding — unknown, canceled, or
	// already paid to the supplier.
	Skipped []int64                 `json:"skipped"`
	Failed  []FailedSupplierPayment `json:"failed"`
}

// PaySupplierReservations settles the given reservations with their supplier: it re-checks that
// each one is still outstanding, records an iCount expense for it, and stores the expense id and
// payment date against the reservation. Each reservation is handled independently, so one
// failure does not hold up the rest.
func (s *SupplierPaymentsService) PaySupplierReservations(ctx context.Context, p PaySupplierReservationsParams) (*PaySupplierReservationsResponse, error) {
	paidAt, err := time.Parse(time.DateOnly, p.PaymentDate)
	if err != nil {
		// Validation already enforced the format, so this means the two disagree.
		rlog.Error("failed to parse supplier payment date", "error", err, "payment_date", p.PaymentDate)
		return nil, api_errors.ErrInternalError
	}

	rows, err := s.query.ListUnpaidSupplierReservationsByIDs(ctx, p.ReservationIDs)
	if err != nil {
		rlog.Error("failed to list unpaid reservations by ids", "error", err)
		return nil, api_errors.ErrInternalError
	}

	outstanding := make(map[int64]db.ListUnpaidSupplierReservationsByIDsRow, len(rows))
	for _, r := range rows {
		outstanding[r.ID] = r
	}

	resp := &PaySupplierReservationsResponse{
		Paid:    []PaidSupplierReservation{},
		Skipped: []int64{},
		Failed:  []FailedSupplierPayment{},
	}

	for _, id := range p.ReservationIDs {
		reservation, isOutstanding := outstanding[id]
		if !isOutstanding {
			resp.Skipped = append(resp.Skipped, id)
			continue
		}

		supplierID, err := s.supplierIDForBroker(reservation.Broker)
		if err != nil {
			rlog.Error("cannot pay supplier for reservation", "error", err, "reservation_id", id)
			resp.Failed = append(resp.Failed, FailedSupplierPayment{ReservationID: id, Reason: ReasonUnsupportedBroker})
			continue
		}

		amount := amountOwed(reservation.PurchasePrice, reservation.BrokerErpPrice)
		expenseID, err := s.createExpense(reservation, supplierID, amount, p.PaymentDate)
		if err != nil {
			rlog.Error("failed to create supplier expense", "error", err, "reservation_id", id)
			resp.Failed = append(resp.Failed, FailedSupplierPayment{ReservationID: id, Reason: ReasonExpenseCreationFailed})
			continue
		}

		if err := s.markPaid(ctx, id, expenseID, paidAt); err != nil {
			// The expense is already in iCount; surface its id so it can be reconciled by hand.
			rlog.Error("created supplier expense but failed to mark reservation as paid",
				"error", err, "reservation_id", id, "expense_id", expenseID)
			resp.Failed = append(resp.Failed, FailedSupplierPayment{
				ReservationID: id,
				Reason:        ReasonMarkPaidFailed,
				ExpenseID:     expenseID,
			})
			continue
		}

		resp.Paid = append(resp.Paid, PaidSupplierReservation{
			ReservationID: id,
			ExpenseID:     expenseID,
			Amount:        amount,
			CurrencyCode:  reservation.CurrencyCode,
		})
	}

	return resp, nil
}

// createExpense records the reservation's supplier cost in iCount and returns the expense id.
func (s *SupplierPaymentsService) createExpense(r db.ListUnpaidSupplierReservationsByIDsRow, supplierID int, amount float64, paymentDate string) (string, error) {
	result, err := s.expenses.CreateExpense(icount.CreateExpenseParams{
		SupplierID:    supplierID,
		ExpenseTypeID: s.cfg.ExpenseTypeID,
		DocType:       supplierExpenseDocType,
		DocNum:        r.BrokerReservationID,
		Sum:           amount,
		CurrencyCode:  r.CurrencyCode,
		Rate:          dbadapters.NumericToFloat64(r.CurrencyRate),
		Paid:          true,
		PaidDate:      paymentDate,
	})
	if err != nil {
		return "", err
	}
	if !result.Status {
		return "", fmt.Errorf("icount rejected the expense: %s (%s) %v", result.Reason, result.ErrorDescription, result.ErrorDetails)
	}

	return strconv.Itoa(result.ExpenseID), nil
}

func (s *SupplierPaymentsService) markPaid(ctx context.Context, reservationID int64, expenseID string, paidAt time.Time) error {
	return s.query.MarkReservationsPaidToSupplier(ctx, db.MarkReservationsPaidToSupplierParams{
		Ids:               []int64{reservationID},
		SupplierExpenseID: &expenseID,
		SupplierPaidAt:    dbadapters.DBTime(paidAt),
	})
}

func (s *SupplierPaymentsService) supplierIDForBroker(b db.Broker) (int, error) {
	switch b {
	case db.BrokerFlex:
		return s.cfg.FlexSupplierID, nil
	default:
		return 0, fmt.Errorf("no icount supplier id configured for broker %q", b)
	}
}
