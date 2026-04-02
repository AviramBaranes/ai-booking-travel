package accounts

import (
	"encore.app/services/accounts/db"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query db.Querier
}

var usersDB = sqldb.NewDatabase("users", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})

func initService() (*Service, error) {
	pgxdb := sqldb.Driver[*pgxpool.Pool](usersDB)
	query := db.New(pgxdb)

	createFirstAdmin(query)

	return &Service{
		query: query,
	}, nil
}
