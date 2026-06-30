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
}

func NewCurrenciesCache(cache cache) *CurrenciesCache {
	return &CurrenciesCache{
		cache: cache,
	}
}

// GetCurrencyRate retrieves the exchange rate for the given currency code from the cache.
// If the rate is not found in the cache, it updates the cache with the latest rate and returns it.
// If any error occurs during the process, it returns an error.
func (c *CurrenciesCache) GetCurrencyRate(ctx context.Context, currencyCode string) (float64, error) {
	rate, err := c.Get(ctx, currencyCode)
	if err != nil {
		if errors.Is(err, ec.Miss) {
			return c.setCurrenciesRates(ctx, currencyCode)
		}
		return 0, err
	}

	return rate, nil
}

func (c *CurrenciesCache) setCurrenciesRates(ctx context.Context, currencyCode string) (float64, error) {
	ic := icount.NewIcount()
	resp, err := ic.FetchCurrencies()
	if err != nil {
		return 0, err
	}

	var cr float64
	for currency, rate := range resp.Rates {
		if err := c.Set(ctx, currency, rate); err != nil {
			return 0, err
		}

		if currency == currencyCode {
			cr = rate
		}
	}

	return cr, nil
}
