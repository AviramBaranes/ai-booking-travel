package booking

import (
	"context"

	markup_rate "encore.app/services/booking/handlers/markup_rate"
)

// ListHertzMarkupRates lists hertz markup rates with pagination, optional filtering, and sorting.
//
//encore:api auth method=GET path=/hertz-markup-rates tag:admin
func (s *Service) ListHertzMarkupRates(ctx context.Context, p markup_rate.ListHertzMarkupRatesParams) (*markup_rate.ListHertzMarkupRatesResponse, error) {
	h := markup_rate.NewHertzMarkupRateService(s.query)
	return h.ListHertzMarkupRates(ctx, p)
}

// CreateHertzMarkupRate creates a new hertz markup rate.
//
//encore:api auth method=POST path=/hertz-markup-rates tag:admin
func (s *Service) CreateHertzMarkupRate(ctx context.Context, p markup_rate.CreateHertzMarkupRateParams) (*markup_rate.HertzMarkupRateResponse, error) {
	h := markup_rate.NewHertzMarkupRateService(s.query)
	return h.CreateHertzMarkupRate(ctx, p)
}

// UpdateHertzMarkupRate updates an existing hertz markup rate.
//
//encore:api auth method=PUT path=/hertz-markup-rates/:id tag:admin
func (s *Service) UpdateHertzMarkupRate(ctx context.Context, id int64, p markup_rate.UpdateHertzMarkupRateParams) (*markup_rate.HertzMarkupRateResponse, error) {
	h := markup_rate.NewHertzMarkupRateService(s.query)
	return h.UpdateHertzMarkupRate(ctx, id, p)
}

// DeleteHertzMarkupRate deletes a hertz markup rate by its ID.
//
//encore:api auth method=DELETE path=/hertz-markup-rates/:id tag:admin
func (s *Service) DeleteHertzMarkupRate(ctx context.Context, id int64) error {
	h := markup_rate.NewHertzMarkupRateService(s.query)
	return h.DeleteHertzMarkupRate(ctx, id)
}
