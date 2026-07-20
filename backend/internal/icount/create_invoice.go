package icount

import (
	"encoding/json"
	"fmt"
)

// CreateInvoiceParams contains the parameters required to create an invoice in iCount.
type CreateInvoiceParams struct {
	ClientID      int
	CurrencyID    int
	Items         []ICountInvoiceItem
	PaymentMethod PaymentMethod
	Rate          float64
	DocType       string
}

// CreateInvoice creates an invoice in iCount using the provided parameters and returns the response from iCount, the response might contain error details if the creation was not successful.
func (i *Icount) CreateInvoice(params CreateInvoiceParams) (*ICountCreateDocResponse, error) {
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
func (i *Icount) createInvoiceDocRequest(p CreateInvoiceParams) ICountCreateDocRequest {
	bankTransfer, payments := p.PaymentMethod.ToRequest()
	if p.DocType == "" {
		p.DocType = "invrec"
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
