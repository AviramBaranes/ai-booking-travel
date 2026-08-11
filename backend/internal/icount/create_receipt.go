package icount

import (
	"encoding/json"
	"fmt"
)

// CreateInvoice creates an invoice in iCount using the provided parameters and returns the response from iCount, the response might contain error details if the creation was not successful.
func (i *Icount) CreateReceipt(params CreateDocParams) (*ICountCreateDocResponse, error) {
	icountReq := i.createReceiptDocRequest(params)

	body, err := i.DoRequest(createDocEndpoint, icountReq)
	if err != nil {
		return nil, fmt.Errorf("creating receipt: %w", err)
	}

	var result ICountCreateDocResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return &result, nil
}

// createReceiptDocRequest constructs the ICountCreateDocRequest for a receipt document from the provided CreateInvoiceParams and the Icount struct
func (i *Icount) createReceiptDocRequest(p CreateDocParams) ICountCreateDocRequest {
	bankTransfer, payments := p.PaymentMethod.ToRequest()
	if p.DocType == "" {
		p.DocType = "receipt"
	}
	return ICountCreateDocRequest{
		ClientID:     p.ClientID,
		DocType:      p.DocType,
		CurrencyID:   p.CurrencyID,
		Rate:         p.Rate,
		BankTransfer: bankTransfer,
		CC:           payments,
		Items:        p.Items,
		Deductions:   p.Deductions,
	}
}
