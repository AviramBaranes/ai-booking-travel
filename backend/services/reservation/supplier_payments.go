package reservation

import (
	"context"
	"net/http"

	fileupload "encore.app/internal/file_upload"
	"encore.app/services/reservation/handlers/supplier_payments"
	"encore.dev/beta/errs"
)

// encore:api auth method=GET path=/supplier-payments/unpaid-reservations tag:accountant
func (s *Service) ListUnpaidSupplierReservations(ctx context.Context, p *supplier_payments.ListUnpaidSupplierReservationsParams) (*supplier_payments.ListUnpaidSupplierReservationsResponse, error) {
	sps := supplier_payments.NewSupplierPaymentsService(s.query)
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

	sps := supplier_payments.NewSupplierPaymentsService(s.query)
	sps.ValidateFlexPaymentSummaryRaw(w, req, file)
}
