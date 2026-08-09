package reservation

import (
	"context"
	"time"

	"encore.app/internal/currency"
	"encore.app/services/accounts"
	"encore.app/services/reservation/db"
	actions "encore.app/services/reservation/handlers/actions"
	"encore.app/services/reservation/handlers/penalties"
	queries "encore.app/services/reservation/handlers/queries"
	"encore.app/services/reservation/handlers/reservation_pricing"
	"encore.app/services/reservation/handlers/supplier_payments"
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

	PenaltyTypeCancellation = "cancellation"
	PenaltyTypeNoShow       = "no_show"
)

type GetReservationResponse = queries.GetReservationResponse
type ListReservationsParams = queries.ListReservationsParams
type ReservationSummary = queries.ReservationSummary
type ListReservationsResponse = queries.ListReservationsResponse
type OpenReservation = queries.OpenReservation
type OpenPenalty = queries.OpenPenalty
type GetOpenReservationsResponse = queries.GetOpenReservationsResponse
type ListOpenReservationsByBillingEntityParams = queries.ListOpenReservationsByBillingEntityParams
type BillingReservation = reservation_pricing.BillingReservation
type ListOpenReservationsByBillingEntityResponse = queries.ListOpenReservationsByBillingEntityResponse
type CurrencyGroup = queries.CurrencyGroup

// --- Supplier payments type aliases ---

type ListUnpaidSupplierReservationsParams = supplier_payments.ListUnpaidSupplierReservationsParams
type ListUnpaidSupplierReservationsResponse = supplier_payments.ListUnpaidSupplierReservationsResponse
type UnpaidSupplierReservation = supplier_payments.UnpaidSupplierReservation
type UnpaidSupplierCurrencyGroup = supplier_payments.UnpaidSupplierCurrencyGroup
type ValidateFlexPaymentSummaryResponse = supplier_payments.ValidateFlexPaymentSummaryResponse
type ApprovedReservation = supplier_payments.ApprovedReservation
type RejectedReservation = supplier_payments.RejectedReservation
type PaySupplierReservationsParams = supplier_payments.PaySupplierReservationsParams
type PaySupplierReservationsResponse = supplier_payments.PaySupplierReservationsResponse
type PaidSupplierReservation = supplier_payments.PaidSupplierReservation
type FailedSupplierPayment = supplier_payments.FailedSupplierPayment

var ErrInvalidPaymentSummaryFile = supplier_payments.ErrInvalidPaymentSummaryFile

// --- Actions type aliases ---

type BookingCancellationEvent = actions.BookingCancellationEvent
type CreateReservationParams = actions.CreateReservationParams
type CreateReservationResponse = actions.CreateReservationResponse
type ApplyVoucherParams = actions.ApplyVoucherParams
type ResolveReservationsParams = actions.ResolveReservationsParams
type ResolvePenaltiesParams = actions.ResolvePenaltiesParams
type SeedReservationsResponse = actions.SeedReservationsResponse
type PayAtPickup = actions.PayAtPickup
type SelectedAddon = actions.SelectedAddon
type VoucherReservationAfterPaymentParams = actions.VoucherReservationAfterPaymentParams

// --- Penalties type aliases ---

type CreatePenaltyParams = penalties.CreatePenaltyParams
type CreatePenaltyResponse = penalties.CreatePenaltyResponse
type BillingPenalty = reservation_pricing.BillingPenalty

// --- Error re-exports ---

var ErrInvalidBillingEntity = queries.ErrInvalidBillingEntity
var ErrOfficeInOrganicOrg = queries.ErrOfficeInOrganicOrg
var ErrOrgIsInorganic = queries.ErrOrgIsInorganic
var ErrCancellationWindowExceeded = actions.ErrCancellationWindowExceeded
var ErrReservationNotCanceled = penalties.ErrReservationNotCanceled
var ErrPenaltyAlreadyExists = penalties.ErrPenaltyAlreadyExists

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
	CID         config.String
	User        config.String
	SupplierIDs supplierIDsConfig
	// ExpenseTypeID is the iCount expense type supplier payments are recorded under.
	ExpenseTypeID config.Int
}

type supplierIDsConfig struct {
	Flex config.Int
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
