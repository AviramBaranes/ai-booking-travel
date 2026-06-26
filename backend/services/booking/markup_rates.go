package booking

import (
	"context"

	markup_rate "encore.app/services/booking/handlers/markup_rate"
)

// ListMarkupRates lists  markup rates with pagination, optional filtering, and sorting.
//
//encore:api auth method=GET path=/-markup-rates tag:admin
func (s *Service) ListMarkupRates(ctx context.Context, p markup_rate.ListMarkupRatesParams) (*markup_rate.ListMarkupRatesResponse, error) {
	h := markup_rate.NewMarkupRateService(s.query)
	return h.ListMarkupRates(ctx, p)
}

// CreateMarkupRate creates a new  markup rate.
//
//encore:api auth method=POST path=/-markup-rates tag:admin
func (s *Service) CreateMarkupRate(ctx context.Context, p markup_rate.CreateMarkupRateParams) (*markup_rate.MarkupRateResponse, error) {
	h := markup_rate.NewMarkupRateService(s.query)
	return h.CreateMarkupRate(ctx, p)
}

// UpdateMarkupRate updates an existing  markup rate.
//
//encore:api auth method=PUT path=/-markup-rates/:id tag:admin
func (s *Service) UpdateMarkupRate(ctx context.Context, id int64, p markup_rate.UpdateMarkupRateParams) (*markup_rate.MarkupRateResponse, error) {
	h := markup_rate.NewMarkupRateService(s.query)
	return h.UpdateMarkupRate(ctx, id, p)
}

// DeleteMarkupRate deletes a  markup rate by its ID.
//
//encore:api auth method=DELETE path=/-markup-rates/:id tag:admin
func (s *Service) DeleteMarkupRate(ctx context.Context, id int64) error {
	h := markup_rate.NewMarkupRateService(s.query)
	return h.DeleteMarkupRate(ctx, id)
}
