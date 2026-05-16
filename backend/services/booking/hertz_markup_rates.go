package booking

import (
	"context"

	"encore.app/services/booking/handlers/hertz_markup_rate_handlers"
)

// ListHertzMarkupRates lists hertz markup rates with pagination, optional filtering, and sorting.
//
//encore:api auth method=GET path=/hertz-markup-rates tag:admin
func (s *Service) ListHertzMarkupRates(ctx context.Context, p hertz_markup_rate_handlers.ListHertzMarkupRatesParams) (*hertz_markup_rate_handlers.ListHertzMarkupRatesResponse, error) {
	h := hertz_markup_rate_handlers.NewHertzMarkupRateService(s.query)
	return h.ListHertzMarkupRates(ctx, p)
}

// CreateHertzMarkupRate creates a new hertz markup rate.
//
//encore:api auth method=POST path=/hertz-markup-rates tag:admin
func (s *Service) CreateHertzMarkupRate(ctx context.Context, p hertz_markup_rate_handlers.CreateHertzMarkupRateParams) (*hertz_markup_rate_handlers.HertzMarkupRateResponse, error) {
	h := hertz_markup_rate_handlers.NewHertzMarkupRateService(s.query)
	return h.CreateHertzMarkupRate(ctx, p)
}

// UpdateHertzMarkupRate updates an existing hertz markup rate.
//
//encore:api auth method=PUT path=/hertz-markup-rates/:id tag:admin
func (s *Service) UpdateHertzMarkupRate(ctx context.Context, id int64, p hertz_markup_rate_handlers.UpdateHertzMarkupRateParams) (*hertz_markup_rate_handlers.HertzMarkupRateResponse, error) {
	h := hertz_markup_rate_handlers.NewHertzMarkupRateService(s.query)
	return h.UpdateHertzMarkupRate(ctx, id, p)
}

// DeleteHertzMarkupRate deletes a hertz markup rate by its ID.
//
//encore:api auth method=DELETE path=/hertz-markup-rates/:id tag:admin
func (s *Service) DeleteHertzMarkupRate(ctx context.Context, id int64) error {
	h := hertz_markup_rate_handlers.NewHertzMarkupRateService(s.query)
	return h.DeleteHertzMarkupRate(ctx, id)
}
