package icount

import "strconv"

// ------- Create Document -------

type ICountCreateDocRequest struct {
	ClientID     int                        `json:"client_id,omitempty"`
	DocType      string                     `json:"doctype"`
	CurrencyID   int                        `json:"currency_id"`
	Rate         float64                    `json:"rate,omitempty"`
	Items        []ICountInvoiceItem        `json:"items"`
	BankTransfer *ICountBankTransferPayment `json:"banktransfer,omitempty"`
	CC           []ICountCCPayment          `json:"cc,omitempty"`
}

type PaymentMethod interface {
	ToRequest() (*ICountBankTransferPayment, []ICountCCPayment)
}

type ICountBankTransferPayment struct {
	Sum     float64 `json:"sum"`
	Date    string  `json:"date"`
	Account int     `json:"account"`
}

func NewBankTransferPayment(sum float64, date string, account int) *ICountBankTransferPayment {
	return &ICountBankTransferPayment{
		Sum:     sum,
		Date:    date,
		Account: account,
	}
}

func (b *ICountBankTransferPayment) ToRequest() (*ICountBankTransferPayment, []ICountCCPayment) {
	return b, nil
}

type ICountCCPayment struct {
	Sum              float64 `json:"sum"`
	Date             string  `json:"payment_date,omitempty"`
	Installments     string  `json:"installments"`
	CCLast4          string  `json:"cc_last4,omitempty"`
	CardType         string  `json:"card_type,omitempty"`
	AlreadyCharged   bool    `json:"already_charged,omitempty"`
	ConfirmationCode string  `json:"confirmation_code,omitempty"`
}

func (c *ICountCCPayment) ToRequest() (*ICountBankTransferPayment, []ICountCCPayment) {
	return nil, []ICountCCPayment{*c}
}

type ICountInvoiceItem struct {
	Description     string   `json:"description"`
	IsTaxExempt     bool     `json:"tax_exempt"`
	SKU             string   `json:"sku,omitempty"`
	UnitPrice       *float64 `json:"unitprice,omitempty"`
	UnitPriceIncvat *float64 `json:"unitprice_incvat,omitempty"`
	Quantity        int      `json:"quantity"`
	CurrencyID      int      `json:"currency_id,omitempty"`
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

// ------- Generate Paypage Sale -------

type ICountGeneratePaypageSaleRequest struct {
	ClientName   string  `json:"client_name,omitempty"`
	FirstName    string  `json:"first_name,omitempty"`
	LastName     string  `json:"last_name,omitempty"`
	Email        string  `json:"email,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	PaypageID    int     `json:"paypage_id"`
	Sum          float64 `json:"sum"`
	CurrencyCode string  `json:"currency_code"`
	Description  string  `json:"description"`
	IsIframe     bool    `json:"is_iframe"`
	XOrderID     int64   `json:"x_order_id,omitempty"`
	SuccessURL   string  `json:"success_url,omitempty"`
	IpnURL       string  `json:"ipn_url,omitempty"`
}

type ICountGeneratePaypageSaleResponse struct {
	Status           bool     `json:"status"`
	Reason           string   `json:"reason,omitempty"`
	ErrorDescription string   `json:"error_description,omitempty"`
	ErrorDetails     []string `json:"error_details,omitempty"`

	SaleUniqID string `json:"sale_uniqid"`
	SaleSID    string `json:"sale_sid"`
	SaleURL    string `json:"sale_url"`
}

// ------- Get Transaction -------

type ICountGetTransactionRequest struct {
	ConfirmationCode string `json:"confirmation_code"`
}

type ICountGetTransactionResponse struct {
	Status           bool          `json:"status"`
	Reason           string        `json:"reason,omitempty"`
	ErrorDescription string        `json:"error_description,omitempty"`
	ErrorDetails     []string      `json:"error_details,omitempty"`
	ResultsCount     int           `json:"results_count,omitempty"`
	ResultsList      []Transaction `json:"results_list,omitempty"`
}

type Transaction struct {
	CCBillLogID      string `json:"cc_bill_log_id"`
	CCCardNumber     string `json:"cc_cardnumber"`
	CCCardType       string `json:"cc_cardtype"`
	CCChargeDate     string `json:"cc_charge_date"`
	CCHolderName     string `json:"cc_holder_name"`
	CCMonth          string `json:"cc_month"`
	CCNumOfPayments  string `json:"cc_numofpayments"`
	CCPaymentType    string `json:"cc_paymenttype"`
	CCProvider       string `json:"cc_provider"`
	CCYear           string `json:"cc_year"`
	CCFirstPayment   string `json:"ccfirstpayment"`
	CCTotal          string `json:"cctotal"`
	ClientEmail      string `json:"client_email"`
	ClientID         string `json:"client_id"`
	ClientIDNo       string `json:"client_idno"`
	ClientName       string `json:"client_name"`
	ConfirmationCode string `json:"confirmation_code"`
	CurrencyCode     string `json:"currency_code"`
	CurrencyID       int    `json:"currency_id"`
	CurrencySign     string `json:"currency_sign"`
	CustomClientID   string `json:"custom_client_id"`
	DealID           string `json:"deal_id"`
	DocNumber        string `json:"docnumber"`
	DocType          string `json:"doctype"`
	Is3DS            string `json:"is_3ds"`
	IsChargeback     string `json:"is_chargeback"`
	// Rate             string `json:"rate"`
	Refunded string `json:"refunded"`
}

func (r *Transaction) ToCCPayment(isCharged bool) *ICountCCPayment {
	total, _ := strconv.ParseFloat(r.CCTotal, 64)
	return &ICountCCPayment{
		Sum:              total,
		Date:             r.CCChargeDate,
		Installments:     r.CCNumOfPayments,
		CCLast4:          r.CCCardNumber,
		CardType:         r.CCCardType,
		ConfirmationCode: r.ConfirmationCode,
		AlreadyCharged:   isCharged,
	}
}
