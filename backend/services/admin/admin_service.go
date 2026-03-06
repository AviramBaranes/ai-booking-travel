package admin

// "encore.dev/storage/sqldb"
// "github.com/jackc/pgx/v5/pgxpool"

// encore:service
type Service struct {
	// query db.Querier
}

// var adminDB = sqldb.NewDatabase("admin", sqldb.DatabaseConfig{
// 	Migrations: "./db/migrations/",
// })

func initService() (*Service, error) {
	// pgxdb := sqldb.Driver[*pgxpool.Pool](adminDB)
	// query := db.New(pgxdb)

	return &Service{
		// query: query,
	}, nil
}
