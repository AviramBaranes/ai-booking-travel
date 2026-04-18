package reservation

import (
	"encore.app/services/reservation/db"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query db.Querier
	pool  *pgxpool.Pool
}

var reservationsDB = sqldb.NewDatabase("reservations", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})
var pgxdb *pgxpool.Pool
var query *db.Queries

func initService() (*Service, error) {
	pgxdb = sqldb.Driver[*pgxpool.Pool](reservationsDB)
	query = db.New(pgxdb)

	return &Service{
		query: query,
		pool:  pgxdb,
	}, nil
}
