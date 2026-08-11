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

// penaltyExpenseDocNumSuffix distinguishes a fee's expense from the rental's, since both are
// recorded under the same booking id.
const penaltyExpenseDocNumSuffix = "-FEE"

// Failure reasons for a reservation we could not settle.
const (
	ReasonUnsupportedBroker     = "unsupported_broker"
	ReasonExpenseCreationFailed = "expense_creation_failed"
	// ReasonMarkPaidFailed means the expense exists in iCount but we failed to link it to the
	// reservation, so the reservation is still outstanding and must be reconciled by hand.
	ReasonMarkPaidFailed = "mark_paid_failed"
)

var ErrNothingToPay = api_errors.NewValidationError("at least one reservation or penalty must be provided")

type PaySupplierParams struct {
	ReservationIDs []int64 `json:"reservationIds"`
	PenaltyIDs     []int64 `json:"penaltyIds"`
	// PaymentDate is the day the supplier was actually paid. It is recorded both on the iCount
	// expense and as the supplier_paid_at of the reservation or fee.
	PaymentDate string `json:"paymentDate" validate:"required,datetime=2006-01-02"`
}

func (p PaySupplierParams) Validate() error {
	if err := validation.ValidateStruct(p); err != nil {
		return err
	}

	// A payment may settle reservations, fees, or both — but not nothing.
	if len(p.ReservationIDs) == 0 && len(p.PenaltyIDs) == 0 {
		return ErrNothingToPay
	}

	return nil
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

// PaidSupplierPenalty is a fee now settled and linked to its iCount expense.
type PaidSupplierPenalty struct {
	PenaltyID     int64   `json:"penaltyId"`
	ReservationID int64   `json:"reservationId"`
	ExpenseID     string  `json:"expenseId"`
	Amount        float64 `json:"amount"`
	CurrencyCode  string  `json:"currencyCode"`
}

// FailedSupplierPenalty is a fee that could not be settled.
type FailedSupplierPenalty struct {
	PenaltyID int64  `json:"penaltyId"`
	Reason    string `json:"reason"`
	// ExpenseID is set when the expense was created but could not be linked to the fee.
	ExpenseID string `json:"expenseId,omitempty"`
}

type PaySupplierResponse struct {
	Paid []PaidSupplierReservation `json:"paid"`
	// Skipped holds the requested ids that are no longer outstanding — unknown, canceled, or
	// already paid to the supplier.
	Skipped       []int64                 `json:"skipped"`
	Failed        []FailedSupplierPayment `json:"failed"`
	PaidPenalties []PaidSupplierPenalty   `json:"paidPenalties"`
	// SkippedPenalties holds the requested fee ids that are unknown or already paid to the supplier.
	SkippedPenalties []int64                 `json:"skippedPenalties"`
	FailedPenalties  []FailedSupplierPenalty `json:"failedPenalties"`
}

// PaySupplier settles the given reservations and fees with their supplier: it re-checks that each
// one is still outstanding, records an iCount expense for it, and stores the expense id and payment
// date against the reservation or fee. Each one is handled independently, so one failure does not
// hold up the rest.
func (s *SupplierPaymentsService) PaySupplier(ctx context.Context, p PaySupplierParams) (*PaySupplierResponse, error) {
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

	penaltyRows, err := s.query.ListUnpaidSupplierPenaltiesByIDs(ctx, p.PenaltyIDs)
	if err != nil {
		rlog.Error("failed to list unpaid penalties by ids", "error", err)
		return nil, api_errors.ErrInternalError
	}

	outstandingPenalties := make(map[int64]db.ListUnpaidSupplierPenaltiesByIDsRow, len(penaltyRows))
	for _, r := range penaltyRows {
		outstandingPenalties[r.ID] = r
	}

	resp := &PaySupplierResponse{
		Paid:             []PaidSupplierReservation{},
		Skipped:          []int64{},
		Failed:           []FailedSupplierPayment{},
		PaidPenalties:    []PaidSupplierPenalty{},
		SkippedPenalties: []int64{},
		FailedPenalties:  []FailedSupplierPenalty{},
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

	for _, id := range p.PenaltyIDs {
		penalty, isOutstanding := outstandingPenalties[id]
		if !isOutstanding {
			resp.SkippedPenalties = append(resp.SkippedPenalties, id)
			continue
		}

		supplierID, err := s.supplierIDForBroker(penalty.Broker)
		if err != nil {
			rlog.Error("cannot pay supplier for penalty", "error", err, "penalty_id", id)
			resp.FailedPenalties = append(resp.FailedPenalties, FailedSupplierPenalty{PenaltyID: id, Reason: ReasonUnsupportedBroker})
			continue
		}

		amount := dbadapters.NumericToFloat64(penalty.Amount)
		expenseID, err := s.createPenaltyExpense(penalty, supplierID, amount, p.PaymentDate)
		if err != nil {
			rlog.Error("failed to create supplier expense for penalty", "error", err, "penalty_id", id)
			resp.FailedPenalties = append(resp.FailedPenalties, FailedSupplierPenalty{PenaltyID: id, Reason: ReasonExpenseCreationFailed})
			continue
		}

		if err := s.markPenaltyPaid(ctx, id, expenseID, paidAt); err != nil {
			// The expense is already in iCount; surface its id so it can be reconciled by hand.
			rlog.Error("created supplier expense but failed to mark penalty as paid",
				"error", err, "penalty_id", id, "expense_id", expenseID)
			resp.FailedPenalties = append(resp.FailedPenalties, FailedSupplierPenalty{
				PenaltyID: id,
				Reason:    ReasonMarkPaidFailed,
				ExpenseID: expenseID,
			})
			continue
		}

		resp.PaidPenalties = append(resp.PaidPenalties, PaidSupplierPenalty{
			PenaltyID:     id,
			ReservationID: penalty.ReservationID,
			ExpenseID:     expenseID,
			Amount:        amount,
			CurrencyCode:  penalty.CurrencyCode,
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

// createPenaltyExpense records the fee's supplier cost in iCount and returns the expense id. The
// document number is suffixed, since the rental of the same booking may already have been expensed
// under the bare booking id.
func (s *SupplierPaymentsService) createPenaltyExpense(p db.ListUnpaidSupplierPenaltiesByIDsRow, supplierID int, amount float64, paymentDate string) (string, error) {
	result, err := s.expenses.CreateExpense(icount.CreateExpenseParams{
		SupplierID:    supplierID,
		ExpenseTypeID: s.cfg.ExpenseTypeID,
		DocType:       supplierExpenseDocType,
		DocNum:        p.BrokerReservationID + penaltyExpenseDocNumSuffix,
		Sum:           amount,
		CurrencyCode:  p.CurrencyCode,
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

func (s *SupplierPaymentsService) markPenaltyPaid(ctx context.Context, penaltyID int64, expenseID string, paidAt time.Time) error {
	return s.query.MarkPenaltiesPaidToSupplier(ctx, db.MarkPenaltiesPaidToSupplierParams{
		Ids:               []int64{penaltyID},
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
