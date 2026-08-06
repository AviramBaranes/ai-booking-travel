// Package supplier_payments holds the logic for reconciling and settling what we owe
// the suppliers (brokers) for reservations we sold.
package supplier_payments

import "encore.app/services/reservation/db"

type SupplierPaymentsService struct {
	query db.Querier
}

func NewSupplierPaymentsService(query db.Querier) *SupplierPaymentsService {
	return &SupplierPaymentsService{query: query}
}
