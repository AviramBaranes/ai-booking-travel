package reservation

import (
	"context"

	"encore.app/services/reservation/db"
	"encore.dev/pubsub"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
	"x.encore.dev/infra/pubsub/outbox"
)

// encore:service
type Service struct {
	query             db.Querier
	pool              *pgxpool.Pool
	cancellationTopic pubsub.Publisher[*BookingCancellationEvent]
}

var reservationsDB = sqldb.NewDatabase("reservations", sqldb.DatabaseConfig{
	Migrations: "./db/migrations/",
})
var pgxdb *pgxpool.Pool
var query *db.Queries

func initService() (*Service, error) {
	pgxdb = sqldb.Driver[*pgxpool.Pool](reservationsDB)
	query = db.New(pgxdb)

	cancellationTopic := pubsub.TopicRef[pubsub.Publisher[*BookingCancellationEvent]](BookingCancellationEvents)

	relay := outbox.NewRelay(outbox.SQLDBStore(reservationsDB))
	outbox.RegisterTopic(relay, cancellationTopic)
	go relay.PollForMessages(context.Background(), -1)

	return &Service{
		query:             query,
		pool:              pgxdb,
		cancellationTopic: cancellationTopic,
	}, nil
}
