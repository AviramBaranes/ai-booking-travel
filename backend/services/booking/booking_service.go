package booking

import (
	"encore.app/services/booking/db"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query db.Querier
}

var bookingsDB = sqldb.NewDatabase("bookings", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})
var pgxdb *pgxpool.Pool
var query *db.Queries

func initService() (*Service, error) {
	pgxdb = sqldb.Driver[*pgxpool.Pool](bookingsDB)
	query = db.New(pgxdb)

	return &Service{
		query: query,
	}, nil
}
