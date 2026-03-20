package booking

import (
	"context"
	"math"

	"encore.app/services/booking/db"
	"encore.dev/config"
	"github.com/jackc/pgx/v5/pgtype"
)

// MarkupProvider calculates markup for a single vehicle.
// Constructed per-search with all static params already resolved.
type MarkupProvider interface {
	CalculateMarkup(isAgent bool, basePrice float64, carGroup, brand string) float64
}

func applyMarkup(basePrice, markupPct float64) float64 {
	return math.Round(basePrice*(1+markupPct/100)*100) / 100
}

// --- Hertz ---

type hertzMarkupKey struct {
	CarGroup string
	Brand    string
}

type markupRates struct {
	Gross float64
	Net   float64
}

// HertzMarkupProvider fetches all matching markup rates from the DB at construction,
// then does pure map lookups per car.
type HertzMarkupProvider struct {
	rates map[hertzMarkupKey]markupRates
}

func NewHertzMarkupProvider(ctx context.Context, q db.Querier, country, pickupDate string, rentalDays int, brands, carGroups []string) (*HertzMarkupProvider, error) {
	rows, err := q.GetHertzMarkupRates(ctx, db.GetHertzMarkupRatesParams{
		Country:    country,
		PickupDate: pgtype.Date{Time: parseDate(pickupDate).Time, Valid: true},
		RentalDays: int32(rentalDays),
		Brands:     brands,
		CarGroups:  carGroups,
	})
	if err != nil {
		return nil, err
	}

	rates := make(map[hertzMarkupKey]markupRates, len(rows))
	for _, r := range rows {
		rates[hertzMarkupKey{CarGroup: r.CarGroup, Brand: r.Brand}] = markupRates{
			Gross: r.MarkUpGross,
			Net:   r.MarkUpNet,
		}
	}
	return &HertzMarkupProvider{rates: rates}, nil
}

func (h *HertzMarkupProvider) CalculateMarkup(isAgent bool, basePrice float64, carGroup, brand string) float64 {
	r, ok := h.rates[hertzMarkupKey{CarGroup: carGroup, Brand: brand}]
	if !ok {
		return 0
	}
	pct := r.Gross
	if isAgent {
		pct = r.Net
	}
	return applyMarkup(basePrice, pct)
}

// --- Flex ---

type FlexMarkupConfig struct {
	MarkUpGross config.Float64
	MarkUpNet   config.Float64
}

// FlexMarkupProvider uses config-driven markup percentages.
type FlexMarkupProvider struct {
	markUpGross float64
	markUpNet   float64
}

func NewFlexMarkupProvider(cfg *FlexMarkupConfig) *FlexMarkupProvider {
	return &FlexMarkupProvider{
		markUpGross: cfg.MarkUpGross(),
		markUpNet:   cfg.MarkUpNet(),
	}
}

func (f *FlexMarkupProvider) CalculateMarkup(isAgent bool, basePrice float64, _, _ string) float64 {
	pct := f.markUpGross
	if isAgent {
		pct = f.markUpNet
	}
	return applyMarkup(basePrice, pct)
}
