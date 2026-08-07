package reservation

import (
	"context"
	"net/http"

	fileupload "encore.app/internal/file_upload"
	"encore.app/internal/icount"
	"encore.app/services/reservation/handlers/supplier_payments"
	"encore.dev/beta/errs"
)

func (s *Service) newSupplierPaymentsService() *supplier_payments.SupplierPaymentsService {
	return supplier_payments.NewSupplierPaymentsService(s.query, icount.NewIcount(), supplier_payments.Config{
		FlexSupplierID: cfg.Icount.SupplierIDs.Flex(),
		ExpenseTypeID:  cfg.Icount.ExpenseTypeID(),
	})
}

// encore:api auth method=GET path=/supplier-payments/unpaid-reservations tag:accountant
func (s *Service) ListUnpaidSupplierReservations(ctx context.Context, p *supplier_payments.ListUnpaidSupplierReservationsParams) (*supplier_payments.ListUnpaidSupplierReservationsResponse, error) {
	sps := s.newSupplierPaymentsService()
	return sps.ListUnpaidSupplierReservations(ctx, p)
}

// ValidateFlexPaymentSummary reconciles an uploaded Flex payment-required summary (xlsx, sent as
// multipart/form-data under the "file" field) against the reservations we still owe Flex for.
//
//encore:api auth method=POST path=/supplier-payments/flex/validate tag:accountant raw
func (s *Service) ValidateFlexPaymentSummary(w http.ResponseWriter, req *http.Request) {
	file, err := fileupload.ExtractFile(req)
	if err != nil {
		errs.HTTPError(w, err)
		return
	}
	defer file.Close()

	sps := s.newSupplierPaymentsService()
	sps.ValidateFlexPaymentSummaryRaw(w, req, file)
}

// PaySupplierReservations records an iCount expense for each of the given reservations that is
// still outstanding, and marks it as paid to the supplier.
//
// encore:api auth method=POST path=/supplier-payments tag:accountant
func (s *Service) PaySupplierReservations(ctx context.Context, p supplier_payments.PaySupplierReservationsParams) (*supplier_payments.PaySupplierReservationsResponse, error) {
	sps := s.newSupplierPaymentsService()
	return sps.PaySupplierReservations(ctx, p)
}
