package currency_handlers

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
)

// ListCurrenciesResponse represents the response containing a list of currencies.
type ListCurrenciesResponse struct {
	Currencies []CurrencyResponse `json:"currencies"`
}

// ListCurrencies retrieves all currencies from the database.
func (s *CurrencyService) ListCurrencies(ctx context.Context) (*ListCurrenciesResponse, error) {
	rows, err := s.query.ListCurrencies(ctx)
	if err != nil {
		rlog.Error("failed to list currencies", "error", err)
		return nil, api_errors.ErrInternalError
	}

	currencies := make([]CurrencyResponse, 0, len(rows))
	for _, r := range rows {
		currencies = append(currencies, toCurrencyResponse(r))
	}

	return &ListCurrenciesResponse{Currencies: currencies}, nil
}
