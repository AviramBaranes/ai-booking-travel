package currency

import (
	"context"
	"errors"

	"encore.app/internal/icount"
	ec "encore.dev/storage/cache"
)

type cache interface {
	Get(ctx context.Context, currencyCode string) (float64, error)
	Set(ctx context.Context, currencyCode string, rate float64) error
}

type CurrenciesCache struct {
	cache
	ic *icount.Icount
}

func NewCurrenciesCache(cache cache) *CurrenciesCache {
	return &CurrenciesCache{
		cache: cache,
		ic:    icount.NewIcount(),
	}
}

// GetCurrencyRate retrieves the exchange rate for the given currency code from the cache.
// If the rate is not found in the cache, it updates the cache with the latest rate and returns it.
// If any error occurs during the process, it returns an error.
func (c *CurrenciesCache) GetCurrencyRate(ctx context.Context, currencyCode string) (float64, error) {
	rate, err := c.Get(ctx, currencyCode)
	if err != nil {
		if errors.Is(err, ec.Miss) {
			rates, err := c.refreshRates(ctx)
			if err != nil {
				return 0, err
			}
			return rates[currencyCode], nil
		}
		return 0, err
	}

	return rate, nil
}

// GetCurrencyRates retrieves the exchange rates for all the given currency codes.
// Any code missing from the cache triggers a single refresh of the whole rates table,
// so a batch of currencies costs at most one call to the rates provider.
func (c *CurrenciesCache) GetCurrencyRates(ctx context.Context, currencyCodes []string) (map[string]float64, error) {
	rates := make(map[string]float64, len(currencyCodes))

	var missing bool
	for _, code := range currencyCodes {
		if _, ok := rates[code]; ok {
			continue
		}

		rate, err := c.Get(ctx, code)
		if err != nil {
			if errors.Is(err, ec.Miss) {
				missing = true
				continue
			}
			return nil, err
		}
		rates[code] = rate
	}

	if !missing {
		return rates, nil
	}

	refreshed, err := c.refreshRates(ctx)
	if err != nil {
		return nil, err
	}

	for _, code := range currencyCodes {
		if _, ok := rates[code]; !ok {
			rates[code] = refreshed[code]
		}
	}

	return rates, nil
}

// refreshRates fetches the full rates table from the provider and writes every rate to the cache.
func (c *CurrenciesCache) refreshRates(ctx context.Context) (map[string]float64, error) {
	resp, err := c.ic.FetchCurrencies()
	if err != nil {
		return nil, err
	}

	for currency, rate := range resp.Rates {
		if err := c.Set(ctx, currency, rate); err != nil {
			return nil, err
		}
	}

	return resp.Rates, nil
}
