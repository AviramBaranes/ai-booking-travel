package markup_rate

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type ListMarkupRatesParams struct {
	Country string `query:"country" validate:"omitempty"`
	Broker  string `json:"broker" validate:"omitempty,oneof=hertz flex"`
	SortBy  string `query:"sortBy" validate:"omitempty"`
	SortDir string `query:"sortDir" validate:"omitempty,oneof=asc desc"`
	Page    int32  `query:"page" validate:"required,gte=1"`
}

func (p ListMarkupRatesParams) Validate() error {
	if err := validation.ValidateStruct(p); err != nil {
		return err
	}
	if p.SortBy != "" && !allowedSortFields[p.SortBy] {
		return api_errors.ErrInvalidValue
	}
	return nil
}

type ListMarkupRatesResponse struct {
	Rates []MarkupRateResponse `json:"rates"`
	Total int64                `json:"total"`
}

func (s *MarkupRateService) ListMarkupRates(ctx context.Context, p ListMarkupRatesParams) (*ListMarkupRatesResponse, error) {
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

	filterParams := db.CountMarkupRatesParams{
		CountryCode: toStringPtr(p.Country),
		Broker:      toStringPtr(p.Broker),
	}

	total, err := s.query.CountMarkupRates(ctx, filterParams)
	if err != nil {
		rlog.Error("failed to count  markup rates", "error", err)
		return nil, api_errors.ErrInternalError
	}

	rows, err := s.query.ListMarkupRates(ctx, db.ListMarkupRatesParams{
		CountryCode: filterParams.CountryCode,
		Broker:      filterParams.Broker,
		SortField:   sortField,
		SortDir:     sortDir,
		QueryOffset: offset,
		QueryLimit:  limit,
	})
	if err != nil {
		rlog.Error("failed to list  markup rates", "error", err)
		return nil, api_errors.ErrInternalError
	}

	rates := make([]MarkupRateResponse, 0, len(rows))
	for _, r := range rows {
		rates = append(rates, toMarkupRateResponse(r))
	}

	return &ListMarkupRatesResponse{Rates: rates, Total: total}, nil
}
