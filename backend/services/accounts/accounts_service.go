package accounts

import (
	"encore.app/services/accounts/db"
	user_handlers "encore.app/services/accounts/handlers/user_handlers"
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

type GetUserEmailParams = user_handlers.GetUserEmailParams
type GetAgentsByOfficeIDRequest = user_handlers.GetAgentsByOfficeIDParams
type GetAgentsByOrganizationIDRequest = user_handlers.GetAgentsByOrganizationIDParams
type ListAdminsEmailsResponse = user_handlers.ListAdminsEmailsResponse
type GetAgentsResponse = user_handlers.GetAgentsResponse
