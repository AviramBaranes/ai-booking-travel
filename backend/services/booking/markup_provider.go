package booking

import (
	"context"

	"encore.app/services/booking/db"
	"encore.dev/rlog"
	"github.com/jackc/pgx/v5/pgtype"
)

// MarkupProvider calculates markup for a single vehicle.
// Constructed per-search with all static params already resolved.
type MarkupProvider interface {
	GetMarkup(isAgent bool, carGroup, brand string) float64
}

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

// NewHertzMarkupProvider constructs a HertzMarkupProvider by fetching the relevant rates from the DB based on the search parameters.
func NewHertzMarkupProvider(ctx context.Context, q db.Querier, country, pickupDate string, rentalDays int, carGroups []string) (*HertzMarkupProvider, error) {
	rows, err := q.GetHertzMarkupRates(ctx, db.GetHertzMarkupRatesParams{
		Country:    country,
		PickupDate: pgtype.Date{Time: parseDate(pickupDate).Time, Valid: true},
		RentalDays: int32(rentalDays),
		CarGroups:  carGroups,
	})
	if err != nil {
		return nil, err
	}

	// log the rates for debugging:
	rlog.Info("fetched Hertz markup rates", "country", country, "pickupDate", pickupDate, "rentalDays", rentalDays, "carGroups", carGroups, "rates", rows)

	rates := make(map[hertzMarkupKey]markupRates, len(rows))
	for _, r := range rows {
		rates[hertzMarkupKey{CarGroup: r.CarGroup, Brand: r.Brand}] = markupRates{
			Gross: r.MarkUpGross,
			Net:   r.MarkUpNet,
		}
	}
	return &HertzMarkupProvider{rates: rates}, nil
}

// GetMarkup returns the markup percentage for the given car group and brand, or false if no specific rate is found.
func (h *HertzMarkupProvider) GetMarkup(isAgent bool, carGroup, brand string) float64 {
	r, ok := h.rates[hertzMarkupKey{CarGroup: carGroup, Brand: brand}]
	if !ok {
		return 0
	}
	if isAgent {
		return r.Net
	}
	return r.Gross
}

// --- Flex ---

// FlexMarkupProvider uses config-driven markup percentages.
type FlexMarkupProvider struct {
	markUpGross float64
	markUpNet   float64
}

func NewFlexMarkupProvider(MarkUpGross, MarkUpNet float64) *FlexMarkupProvider {
	return &FlexMarkupProvider{
		markUpGross: MarkUpGross,
		markUpNet:   MarkUpNet,
	}
}

func (f *FlexMarkupProvider) GetMarkup(isAgent bool, _, _ string) float64 {
	if isAgent {
		return f.markUpNet
	}
	return f.markUpGross
}
