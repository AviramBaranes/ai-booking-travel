package icount

import (
	"encoding/json"
	"fmt"

	"encore.dev/rlog"
)

// CreateInvoiceParams contains the parameters required to create an invoice in iCount.
type CreateInvoiceParams struct {
	ClientID      int
	CurrencyID    int
	Items         []ICountInvoiceItem
	PaymentMethod PaymentMethod
	Rate          float64
}

// CreateInvoice creates an invoice in iCount using the provided parameters and returns the response from iCount, the response might contain error details if the creation was not successful.
func (i *Icount) CreateInvoice(params CreateInvoiceParams) (*ICountCreateDocResponse, error) {
	icountReq := i.createInvoiceDocRequest(params)
	rlog.Info("creating invoice in iCount", "request", icountReq)

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
func (i *Icount) createInvoiceDocRequest(params CreateInvoiceParams) ICountCreateDocRequest {
	bankTransfer, payments := params.PaymentMethod.ToRequest()
	return ICountCreateDocRequest{
		ClientID:     params.ClientID,
		DocType:      "invrec",
		CurrencyID:   params.CurrencyID,
		Rate:         params.Rate,
		BankTransfer: bankTransfer,
		CC:           payments,
		Items:        params.Items,
	}
}
