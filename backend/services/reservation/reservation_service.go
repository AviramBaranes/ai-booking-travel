package reservation

import (
	"context"

	"encore.app/services/reservation/db"
	queries "encore.app/services/reservation/handlers/queries"
	"encore.dev/config"
	"encore.dev/pubsub"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
	"x.encore.dev/infra/pubsub/outbox"
)

type GetReservationResponse = queries.GetReservationResponse
type ListReservationsParams = queries.ListReservationsParams
type ReservationSummary = queries.ReservationSummary
type ListReservationsResponse = queries.ListReservationsResponse
type OpenReservation = queries.OpenReservation
type GetOpenReservationsResponse = queries.GetOpenReservationsResponse
type ListOpenReservationsByBillingEntityParams = queries.ListOpenReservationsByBillingEntityParams
type BillingReservation = queries.BillingReservation
type ListOpenReservationsByBillingEntityResponse = queries.ListOpenReservationsByBillingEntityResponse
type CurrencyGroup = queries.CurrencyGroup

// --- Error re-exports ---

var ErrInvalidBillingEntity = queries.ErrInvalidBillingEntity
var ErrOfficeInOrganicOrg = queries.ErrOfficeInOrganicOrg
var ErrOrgIsInorganic = queries.ErrOrgIsInorganic

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

type ReservationCfg struct {
	VAT config.Float64

	Icount icountConfig
}

type icountConfig struct {
	CID  config.String
	User config.String
}

var cfg = config.Load[*ReservationCfg]()

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
