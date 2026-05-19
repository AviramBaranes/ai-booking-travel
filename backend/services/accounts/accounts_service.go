package accounts

import (
	"encore.app/services/accounts/db"
	user "encore.app/services/accounts/handlers/user"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query db.Querier
}

var accountsDb = sqldb.NewDatabase("accounts", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})

func initService() (*Service, error) {
	pgxdb := sqldb.Driver[*pgxpool.Pool](accountsDb)
	query := db.New(pgxdb)

	createFirstAdmin(query)

	return &Service{
		query: query,
	}, nil
}

type GetUserEmailParams = user.GetUserEmailParams
type ListAdminsEmailsResponse = user.ListAdminsEmailsResponse
type GetAgentsResponse = user.GetAgentsResponse
