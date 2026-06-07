package accounts

import (
	"encore.app/services/accounts/db"
	lookup "encore.app/services/accounts/handlers/lookup"
	office "encore.app/services/accounts/handlers/office"
	organization "encore.app/services/accounts/handlers/organization"
	user "encore.app/services/accounts/handlers/user"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
)

// encore:service
type Service struct {
	query db.Querier
}

// TEMP CHANGE

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
type AccountName = lookup.AccountName
type GetAccountsLookupParams = lookup.GetAccountsLookupParams
type GetAccountsLookupResponse = lookup.GetAccountsLookupResponse
type CreateOrganizationParams = organization.CreateOrganizationParams
type OrganizationResponse = organization.OrganizationResponse
type CreateOfficeParams = office.CreateOfficeParams
type OfficeResponse = office.OfficeResponse
type CreateAgentParams = user.CreateAgentParams
type CreateAgentResponse = user.CreateAgentResponse
type CreateAdminParams = user.CreateAdminParams
type CreateAdminResponse = user.CreateAdminResponse
