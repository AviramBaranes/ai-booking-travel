package icount

import (
	"encoding/json"
	"fmt"
)

// CreateDocParams contains the parameters required to create a doc in iCount.
type CreateDocParams struct {
	ClientID      int
	CurrencyID    int
	Items         []ICountInvoiceItem
	PaymentMethod PaymentMethod
	Rate          float64
	DocType       string
	// Deductions maps a deduction type id to the amount withheld at source, in document currency.
	// Read by createReceiptDocRequest only — iCount ignores deductions on an invoice, which states
	// a debt rather than the payment the tax was withheld from.
	Deductions map[int]float64
}

// CreateInvoice creates an invoice in iCount using the provided parameters and returns the response from iCount, the response might contain error details if the creation was not successful.
func (i *Icount) CreateInvoice(params CreateDocParams) (*ICountCreateDocResponse, error) {
	icountReq := i.createInvoiceDocRequest(params)

	body, err := i.DoRequest(createDocEndpoint, icountReq)
	if err != nil {
		return nil, fmt.Errorf("creating invoice: %w", err)
	}

	var result ICountCreateDocResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return &result, nil
}

// createInvoiceDocRequest constructs the ICountCreateDocRequest from the provided CreateInvoiceParams and the Icount struct
func (i *Icount) createInvoiceDocRequest(p CreateDocParams) ICountCreateDocRequest {
	bankTransfer, payments := p.PaymentMethod.ToRequest()
	if p.DocType == "" {
		p.DocType = "invoice"
	}
	return ICountCreateDocRequest{
		ClientID:     p.ClientID,
		DocType:      p.DocType,
		CurrencyID:   p.CurrencyID,
		Rate:         p.Rate,
		BankTransfer: bankTransfer,
		CC:           payments,
		Items:        p.Items,
	}
}
