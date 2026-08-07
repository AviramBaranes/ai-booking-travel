package icount

import (
	"encoding/json"
	"fmt"
)

// CreateExpenseParams contains the parameters required to create an expense in iCount.
type CreateExpenseParams struct {
	SupplierID    int
	ExpenseTypeID int
	// DocType and DocNum identify the document we received from the supplier.
	DocType string
	DocNum  string
	Sum     float64
	// CurrencyCode is an ISO-4217 code; iCount defaults to ILS when empty. Rate is optional —
	// iCount resolves it by date when omitted.
	CurrencyCode string
	Rate         float64
	Paid         bool
	// PaidDate is formatted as iCount expects (YYYY-MM-DD); it defaults to the invoice date.
	PaidDate string
}

// CreateExpense records an expense in iCount and returns the response, which may carry error
// details if the creation was not successful.
func (i *Icount) CreateExpense(params CreateExpenseParams) (*ICountCreateExpenseResponse, error) {
	body, err := i.DoRequest(createExpenseEndpoint, i.createExpenseRequest(params))
	if err != nil {
		return nil, fmt.Errorf("creating expense: %w", err)
	}

	var result ICountCreateExpenseResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return &result, nil
}

func (i *Icount) createExpenseRequest(p CreateExpenseParams) ICountCreateExpenseRequest {
	return ICountCreateExpenseRequest{
		SupplierID:    p.SupplierID,
		ExpenseTypeID: p.ExpenseTypeID,
		DocType:       p.DocType,
		DocNum:        p.DocNum,
		Sum:           p.Sum,
		CurrencyCode:  p.CurrencyCode,
		Rate:          p.Rate,
		Paid:          p.Paid,
		PaidDate:      p.PaidDate,
	}
}
