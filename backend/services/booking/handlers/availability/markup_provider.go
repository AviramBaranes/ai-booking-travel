package availability

import (
	"context"

	"encore.app/services/booking/db"
)

// MarkupProvider calculates markup for a single vehicle.
// Constructed per-search with all static params already resolved.
type MarkupProvider interface {
	GetMarkup(isAgent bool) float64
}

type markupRates struct {
	Gross float64
	Net   float64
}

// markupProvider is the type that implements MarkupProvider
type markupProvider struct {
	netRate   float64
	grossRate float64
}

func NewMarkupProviderFromCfg(netRate, grossRate float64) MarkupProvider {
	return &markupProvider{
		netRate:   netRate,
		grossRate: grossRate,
	}
}

// NewHertzMarkupProvider constructs a HertzMarkupProvider by fetching the relevant rates from the DB based on the search parameters.
func NewMarkupProviderFromDB(ctx context.Context, q db.Querier, countryCode, broker string) (MarkupProvider, error) {
	ratesRow, err := q.GetBrokerMarkupRateByCountryCode(ctx, db.GetBrokerMarkupRateByCountryCodeParams{
		CountryCode: countryCode,
		Broker:      db.Broker(broker),
	})
	if err != nil {
		return nil, err
	}

	return &markupProvider{
		netRate:   ratesRow.MarkUpNet,
		grossRate: ratesRow.MarkUpGross,
	}, nil
}

// GetMarkup returns the markup percentage for the given car group and brand, or false if no specific rate is found.
func (h *markupProvider) GetMarkup(isAgent bool) float64 {
	if isAgent {
		return h.netRate
	}
	return h.grossRate
}
