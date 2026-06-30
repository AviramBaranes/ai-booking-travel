package actions

import (
	"encore.app/internal/currency"
	"encore.app/services/reservation/db"
	"encore.dev/pubsub"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds configuration values required by the action service.
type Config struct {
	VAT float64
}

// ActionService provides all reservation write/action operations.
type ActionService struct {
	query             db.Querier
	pool              *pgxpool.Pool
	cancellationTopic pubsub.Publisher[*BookingCancellationEvent]
	cfg               Config
	currencyCache     *currency.CurrenciesCache
}

// NewActionService creates a new ActionService with the given dependencies.
func NewActionService(
	query db.Querier,
	pool *pgxpool.Pool,
	topic pubsub.Publisher[*BookingCancellationEvent],
	currencyCache *currency.CurrenciesCache,
	cfg Config,
) *ActionService {
	return &ActionService{
		query:             query,
		pool:              pool,
		cancellationTopic: topic,
		cfg:               cfg,
		currencyCache:     currencyCache,
	}
}
