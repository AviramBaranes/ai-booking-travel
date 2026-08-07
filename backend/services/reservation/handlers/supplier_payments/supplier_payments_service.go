// Package supplier_payments holds the logic for reconciling and settling what we owe
// the suppliers (brokers) for reservations we sold.
package supplier_payments

import (
	"encore.app/internal/icount"
	"encore.app/services/reservation/db"
)

// Config holds configuration values required by the supplier payments service.
type Config struct {
	// FlexSupplierID is the Flex supplier record in iCount.
	FlexSupplierID int
	// ExpenseTypeID is the iCount expense type all supplier payments are recorded under.
	ExpenseTypeID int
}

// ExpenseCreator records an expense in the accounting system. Narrowed to the one call this
// service makes so tests can stand in for iCount.
type ExpenseCreator interface {
	CreateExpense(p icount.CreateExpenseParams) (*icount.ICountCreateExpenseResponse, error)
}

type SupplierPaymentsService struct {
	query    db.Querier
	expenses ExpenseCreator
	cfg      Config
}

func NewSupplierPaymentsService(query db.Querier, expenses ExpenseCreator, cfg Config) *SupplierPaymentsService {
	return &SupplierPaymentsService{query: query, expenses: expenses, cfg: cfg}
}
