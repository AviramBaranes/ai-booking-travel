package locations

import (
	"encore.app/services/locations/db"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query db.Querier
}

var locationsDB = sqldb.NewDatabase("locations", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})

func initService() (*Service, error) {
	pgxdb := sqldb.Driver[*pgxpool.Pool](locationsDB)
	query := db.New(pgxdb)

	return &Service{
		query: query,
	}, nil
}
