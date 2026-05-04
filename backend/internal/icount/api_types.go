package icount

// ------- Create Document -------

type ICountCreateDocRequest struct {
	CID          string                     `json:"cid"`
	User         string                     `json:"user"`
	Pass         string                     `json:"pass"`
	ClientID     int                        `json:"client_id"`
	DocType      string                     `json:"doctype"`
	CurrencyID   int                        `json:"currency_id"`
	Items        []ICountInvoiceItem        `json:"items"`
	BankTransfer *ICountBankTransferPayment `json:"banktransfer,omitempty"`
}

type ICountBankTransferPayment struct {
	Sum     float64 `json:"sum"`
	Date    string  `json:"date"`
	Account int     `json:"account"`
}

type ICountInvoiceItem struct {
	Description string  `json:"description"`
	IsTaxExempt bool    `json:"tax_exempt"`
	SKU         string  `json:"sku,omitempty"`
	UnitPrice   float64 `json:"unitprice"`
	Quantity    int     `json:"quantity"`
}

type ICountCreateDocResponse struct {
	Status           bool     `json:"status"`
	Reason           string   `json:"reason"`
	ErrorDescription string   `json:"error_description"`
	ErrorDetails     []string `json:"error_details"`
	DocNum           string   `json:"docnum"`
	DocURL           string   `json:"doc_url"`
}

// ------- Fetch Currencies -------

type GetCurrenciesRatesResponse struct {
	Status bool               `json:"status"`
	Reason string             `json:"reason"`
	Rates  map[string]float64 `json:"currency_rates,omitempty"`
}

type GetCurrenciesRatesRequest struct {
	CID  string `json:"cid"`
	User string `json:"user"`
	Pass string `json:"pass"`
}
