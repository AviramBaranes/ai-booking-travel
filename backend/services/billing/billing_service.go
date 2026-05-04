package billing

import "encore.dev/config"

type billingConfig struct {
	MonthlyReport monthlyReportConfig
	Icount        icountConfig
	Invoice       invoiceConfig
}

type monthlyReportConfig struct {
	Headers monthlyReportHeadersConfig
	Styles  monthlyReportStylesConfig
}

type monthlyReportHeadersConfig struct {
	OfficeName           config.String
	AgentName            config.String
	DriverName           config.String
	ReservationCreatedAt config.String
	ReservationID        config.String
	VoucherDate          config.String
	VoucherNumber        config.String
	AgentVoucherNumber   config.String
	PickupDate           config.String
	ReturnDate           config.String
	CountryCode          config.String
	RentalDays           config.String
	Currency             config.String
	NetPrice             config.String
	FullCoverage         config.String
	TotalNetPrice        config.String
}

type monthlyReportStylesConfig struct {
	HeaderBackgroundColor    config.String
	RefundRowBackgroundColor config.String
	TotalRowBackgroundColor  config.String
	BorderColor              config.String
}

type icountConfig struct {
	CID       config.String
	User      config.String
	AccountID config.Int
}

type invoiceConfig struct {
	PurchaseItemDescription     config.String
	ProfitItemDescription       config.String
	ProfitAndErpItemDescription config.String
}

var cfg = config.Load[*billingConfig]()
