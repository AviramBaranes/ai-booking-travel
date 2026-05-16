package booking

import (
	"context"

	"encore.app/services/booking/handlers/currency_handlers"
	"encore.dev/config"
	"encore.dev/cron"
)

type icountConfig struct {
	CID  config.String
	User config.String
}

var icountCfg = config.Load[*icountConfig]()

// ListCurrencies lists all currencies.
//
//encore:api auth method=GET path=/currencies tag:admin
func (s *Service) ListCurrencies(ctx context.Context) (*currency_handlers.ListCurrenciesResponse, error) {
	cs := currency_handlers.NewCurrencyService(s.query)
	return cs.ListCurrencies(ctx)
}

// CreateCurrency creates a new currency.
//
//encore:api auth method=POST path=/currencies tag:admin
func (s *Service) CreateCurrency(ctx context.Context, params currency_handlers.CreateCurrencyParams) (*currency_handlers.CurrencyResponse, error) {
	cs := currency_handlers.NewCurrencyService(s.query)
	return cs.CreateCurrency(ctx, params)
}

// UpdateCurrency updates an existing currency.
//
//encore:api auth method=PUT path=/currencies/:id tag:admin
func (s *Service) UpdateCurrency(ctx context.Context, id int64, params currency_handlers.UpdateCurrencyParams) (*currency_handlers.CurrencyResponse, error) {
	cs := currency_handlers.NewCurrencyService(s.query)
	return cs.UpdateCurrency(ctx, id, params)
}

// DeleteCurrency deletes a currency by its ID.
//
//encore:api auth method=DELETE path=/currencies/:id tag:admin
func (s *Service) DeleteCurrency(ctx context.Context, id int64) error {
	cs := currency_handlers.NewCurrencyService(s.query)
	return cs.DeleteCurrency(ctx, id)
}

// UpdateCurrenciesRates updates the exchange rates for all currencies from iCount, it runs via a cron job.
//
//encore:api private
func (s *Service) UpdateCurrenciesRates(ctx context.Context) error {
	cs := currency_handlers.NewCurrencyService(s.query)
	return cs.UpdateCurrenciesRates(ctx, icountCfg.CID(), icountCfg.User())
}

var _ = cron.NewJob("currencies-sync", cron.JobConfig{
	Title:    "Sync Currencies Rates",
	Schedule: "0 0 * * *", // At 00:00 every day.
	Endpoint: UpdateCurrenciesRates,
})
