package reservation

import (
	"context"
	"time"

	"encore.app/internal/currency"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	actions "encore.app/services/reservation/handlers/actions"
	queries "encore.app/services/reservation/handlers/queries"
	"encore.app/services/reservation/handlers/reservation_pricing"
	"encore.dev/config"
	"encore.dev/pubsub"
	"encore.dev/storage/cache"
	"encore.dev/storage/sqldb"
	"github.com/jackc/pgx/v5/pgxpool"
	"x.encore.dev/infra/pubsub/outbox"
)

const (
	PaymentStatusUnpaid        = "unpaid"
	PaymentStatusPaid          = "paid"
	PaymentStatusRefunded      = "refunded"
	PaymentStatusRefundPending = "refund_pending"

	ReservationStatusBooked    = "booked"
	ReservationStatusCanceled  = "canceled"
	ReservationStatusVouchered = "vouchered"
)

type GetReservationResponse = queries.GetReservationResponse
type ListReservationsParams = queries.ListReservationsParams
type ReservationSummary = queries.ReservationSummary
type ListReservationsResponse = queries.ListReservationsResponse
type OpenReservation = queries.OpenReservation
type GetOpenReservationsResponse = queries.GetOpenReservationsResponse
type ListOpenReservationsByBillingEntityParams = queries.ListOpenReservationsByBillingEntityParams
type BillingReservation = reservation_pricing.BillingReservation
type ListOpenReservationsByBillingEntityResponse = queries.ListOpenReservationsByBillingEntityResponse
type CurrencyGroup = queries.CurrencyGroup

// --- Actions type aliases ---

type BookingCancellationEvent = actions.BookingCancellationEvent
type CreateReservationParams = actions.CreateReservationParams
type CreateReservationResponse = actions.CreateReservationResponse
type ApplyVoucherParams = actions.ApplyVoucherParams
type ResolveReservationsParams = actions.ResolveReservationsParams
type SeedReservationsResponse = actions.SeedReservationsResponse
type PayAtPickup = actions.PayAtPickup
type SelectedAddon = actions.SelectedAddon
type VoucherReservationAfterPaymentParams = actions.VoucherReservationAfterPaymentParams

// --- Error re-exports ---

var ErrInvalidBillingEntity = queries.ErrInvalidBillingEntity
var ErrOfficeInOrganicOrg = queries.ErrOfficeInOrganicOrg
var ErrOrgIsInorganic = queries.ErrOrgIsInorganic
var ErrCancellationWindowExceeded = actions.ErrCancellationWindowExceeded

// BookingCancellationEvents is a pub/sub topic that publishes events whenever a reservation is canceled.
var BookingCancellationEvents = pubsub.NewTopic[*actions.BookingCancellationEvent]("booking-cancellation-events", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

// currenciesRates is a cache for storing currency rates with a default expiry of 12 hours.
var currenciesRates = cache.NewFloatKeyspace[string](accounts.GlobalCache, cache.KeyspaceConfig{
	KeyPattern:    "reservation-currencies/:key",
	DefaultExpiry: cache.ExpireIn(12 * time.Hour),
})

// encore:service
type Service struct {
	query             db.Querier
	pool              *pgxpool.Pool
	cancellationTopic pubsub.Publisher[*actions.BookingCancellationEvent]
	currencyCache     *currency.CurrenciesCache
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

	cancellationTopic := pubsub.TopicRef[pubsub.Publisher[*actions.BookingCancellationEvent]](BookingCancellationEvents)

	relay := outbox.NewRelay(outbox.SQLDBStore(reservationsDB))
	outbox.RegisterTopic(relay, cancellationTopic)
	go relay.PollForMessages(context.Background(), -1)

	return &Service{
		query:             query,
		pool:              pgxdb,
		cancellationTopic: cancellationTopic,
		currencyCache:     currency.NewCurrenciesCache(currenciesRates),
	}, nil
}
