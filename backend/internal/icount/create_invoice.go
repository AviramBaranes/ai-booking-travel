package icount

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type icount struct {
	httpClient *http.Client
	cid        string
	user       string
	pass       string
	accountID  int
}

// NewIcount initializes and returns an icount struct with the provided credentials and account information.
func NewIcount(cid, user, pass string, accountID int) icount {
	return icount{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cid:        cid,
		user:       user,
		pass:       pass,
		accountID:  accountID,
	}
}

// CreateInvoiceParams contains the parameters required to create an invoice in iCount.
type CreateInvoiceParams struct {
	ClientID   int
	CurrencyID int
	Sum        float64
	Date       string
	Items      []ICountInvoiceItem
}

// CreateInvoice creates an invoice in iCount using the provided parameters and returns the response from iCount, the response might contain error details if the creation was not successful.
func (i icount) CreateInvoice(params CreateInvoiceParams) (*ICountCreateDocResponse, error) {
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

// createInvoiceDocRequest constructs the ICountCreateDocRequest from the provided CreateInvoiceParams and the icount struct
func (i icount) createInvoiceDocRequest(params CreateInvoiceParams) ICountCreateDocRequest {
	return ICountCreateDocRequest{
		CID:        i.cid,
		User:       i.user,
		Pass:       i.pass,
		ClientID:   params.ClientID,
		DocType:    "invrec",
		CurrencyID: params.CurrencyID,
		BankTransfer: &ICountBankTransferPayment{
			Sum:     params.Sum,
			Date:    params.Date,
			Account: i.accountID,
		},
		Items: params.Items,
	}
}
