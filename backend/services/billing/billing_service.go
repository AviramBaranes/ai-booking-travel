package billing

import (
	"fmt"
	"time"

	"encore.app/internal/currency"
	"encore.app/services/accounts"
	"encore.dev/config"
	"encore.dev/storage/cache"
)

// encore:service
type Service struct {
	ratesCache *currency.CurrenciesCache
}

// currenciesRates is a cache for storing currency rates with a default expiry of 24 hours.
var currenciesRates = cache.NewFloatKeyspace[string](accounts.GlobalCache, cache.KeyspaceConfig{
	KeyPattern:    "billing-currencies/:key",
	DefaultExpiry: cache.ExpireIn(12 * time.Hour),
})

func initService() (*Service, error) {
	if currenciesRates == nil {
		return nil, fmt.Errorf("currenciesRates cache is not initialized")
	}
	ratesCache := currency.NewCurrenciesCache(currenciesRates)
	return &Service{
		ratesCache: ratesCache,
	}, nil
}

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
	DropoffDate          config.String
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
	AccountID config.Int
	PaypageID config.Int
}

type invoiceConfig struct {
	PurchaseItemDescription     config.String
	ProfitItemDescription       config.String
	ProfitAndErpItemDescription config.String
}

var cfg = config.Load[*billingConfig]()
