package hertz_markup_rate_handlers

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type ListHertzMarkupRatesParams struct {
	Country  string `query:"country" validate:"omitempty"`
	Brand    string `query:"brand" validate:"omitempty"`
	CarGroup string `query:"carGroup" validate:"omitempty"`
	SortBy   string `query:"sortBy" validate:"omitempty"`
	SortDir  string `query:"sortDir" validate:"omitempty,oneof=asc desc"`
	Page     int32  `query:"page" validate:"required,gte=1"`
}

func (p ListHertzMarkupRatesParams) Validate() error {
	if err := validation.ValidateStruct(p); err != nil {
		return err
	}
	if p.SortBy != "" && !allowedSortFields[p.SortBy] {
		return api_errors.ErrInvalidValue
	}
	return nil
}

type ListHertzMarkupRatesResponse struct {
	Rates []HertzMarkupRateResponse `json:"rates"`
	Total int64                     `json:"total"`
}

func (s *HertzMarkupRateService) ListHertzMarkupRates(ctx context.Context, p ListHertzMarkupRatesParams) (*ListHertzMarkupRatesResponse, error) {
	offset := (p.Page - 1) * limit

	sortField := p.SortBy
	sortDir := p.SortDir
	if sortField == "" {
		sortField = "country"
		sortDir = "asc"
	}
	if sortDir == "" {
		sortDir = "asc"
	}

	filterParams := db.CountHertzMarkupRatesParams{
		Country:  toStringPtr(p.Country),
		Brand:    toStringPtr(p.Brand),
		CarGroup: toStringPtr(p.CarGroup),
	}

	total, err := s.query.CountHertzMarkupRates(ctx, filterParams)
	if err != nil {
		rlog.Error("failed to count hertz markup rates", "error", err)
		return nil, api_errors.ErrInternalError
	}

	rows, err := s.query.ListHertzMarkupRates(ctx, db.ListHertzMarkupRatesParams{
		Country:     filterParams.Country,
		Brand:       filterParams.Brand,
		CarGroup:    filterParams.CarGroup,
		SortField:   sortField,
		SortDir:     sortDir,
		QueryOffset: offset,
		QueryLimit:  limit,
	})
	if err != nil {
		rlog.Error("failed to list hertz markup rates", "error", err)
		return nil, api_errors.ErrInternalError
	}

	rates := make([]HertzMarkupRateResponse, 0, len(rows))
	for _, r := range rows {
		rates = append(rates, toHertzMarkupRateResponse(r))
	}

	return &ListHertzMarkupRatesResponse{Rates: rates, Total: total}, nil
}
