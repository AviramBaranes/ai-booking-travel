package actions

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/services/reservation/db"
	"encore.dev/rlog"
)

type SaveInvoiceDocNumParams struct {
	ID     int64
	DocNum string
}

// SaveInvoiceDocNum saves the invoice document number for a reservation.
func (s *ActionService) SaveInvoiceDocNum(ctx context.Context, p SaveInvoiceDocNumParams) error {
	if err := s.query.SaveInvoiceDocNum(ctx, db.SaveInvoiceDocNumParams{
		ID:            p.ID,
		InvoiceDocNum: &p.DocNum,
	}); err != nil {
		rlog.Error("failed to save invoice doc num", "error", err, "reservation_id", p.ID, "doc_num", p.DocNum)
		return api_errors.ErrInternalError
	}

	return nil
}
