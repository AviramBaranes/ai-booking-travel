package billing

import (
	"fmt"
	"time"

	"encore.app/internal/currency"
	"encore.app/services/accounts"
	"encore.app/services/billing/db"
	"encore.dev/config"
	"encore.dev/storage/cache"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query      db.Querier
	ratesCache *currency.CurrenciesCache
}

// currenciesRates is a cache for storing currency rates with a default expiry of 24 hours.
var currenciesRates = cache.NewFloatKeyspace[string](accounts.GlobalCache, cache.KeyspaceConfig{
	KeyPattern:    "billing-currencies/:key",
	DefaultExpiry: cache.ExpireIn(12 * time.Hour),
})

var billingDB = sqldb.NewDatabase("billing", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})
var pgxdb *pgxpool.Pool
var query *db.Queries

func initService() (*Service, error) {
	pgxdb = sqldb.Driver[*pgxpool.Pool](billingDB)
	query = db.New(pgxdb)

	if currenciesRates == nil {
		return nil, fmt.Errorf("currenciesRates cache is not initialized")
	}
	ratesCache := currency.NewCurrenciesCache(currenciesRates)
	return &Service{
		query:      query,
		ratesCache: ratesCache,
	}, nil
}

type billingConfig struct {
	MonthlyReport monthlyReportConfig
	Icount        icountConfig
	Invoice       invoiceConfig
}

type icountConfig struct {
	AccountID          config.Int
	AgentsPaypageID    config.Int
	CustomersPaypageID config.Int
}

var cfg = config.Load[*billingConfig]()
