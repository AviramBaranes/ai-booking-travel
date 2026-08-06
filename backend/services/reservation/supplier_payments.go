package reservation

import (
	"context"

	"encore.app/services/reservation/handlers/supplier_payments"
)

// encore:api auth method=GET path=/supplier-payments/unpaid-reservations tag:accountant
func (s *Service) ListUnpaidSupplierReservations(ctx context.Context, p *supplier_payments.ListUnpaidSupplierReservationsParams) (*supplier_payments.ListUnpaidSupplierReservationsResponse, error) {
	sps := supplier_payments.NewSupplierPaymentsService(s.query)
	return sps.ListUnpaidSupplierReservations(ctx, p)
}
