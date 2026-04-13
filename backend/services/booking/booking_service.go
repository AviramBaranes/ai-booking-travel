package booking

import (
	"context"

	"encore.app/services/booking/db"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query db.Querier
	t     *TranslationCache
}

var bookingsDB = sqldb.NewDatabase("bookings", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})
var pgxdb *pgxpool.Pool
var query *db.Queries

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
	}

	if err := svc.refreshTranslations(context.Background()); err != nil {
		return nil, err
	}

	return svc, nil
}
