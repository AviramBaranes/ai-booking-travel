package booking

import (
	"context"
	"time"

	"encore.app/internal/currency"
	"encore.app/services/accounts"
	"encore.app/services/booking/db"
	"encore.dev/storage/cache"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query db.Querier
	t     *TranslationCache
	c     *currency.CurrenciesCache
}

var bookingsDB = sqldb.NewDatabase("bookings", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})
var pgxdb *pgxpool.Pool
var query *db.Queries

// currenciesRates is a cache for storing currency rates with a default expiry of 12 hours.
var currenciesRates = cache.NewFloatKeyspace[string](accounts.GlobalCache, cache.KeyspaceConfig{
	KeyPattern:    "booking-currencies/:key",
	DefaultExpiry: cache.ExpireIn(12 * time.Hour),
})

// initService initializes the booking service by setting up the database connection and loading translations into the cache.
func initService() (*Service, error) {
	pgxdb = sqldb.Driver[*pgxpool.Pool](bookingsDB)
	query = db.New(pgxdb)

	svc := &Service{
		query: query,
		t: &TranslationCache{
			translations: make(map[string]string),
			known:        make(map[string]struct{}),
		},
		c: currency.NewCurrenciesCache(currenciesRates),
	}

	if err := svc.refreshTranslations(context.Background()); err != nil {
		return nil, err
	}

	go svc.startBackgroundRefresh()

	return svc, nil
}
