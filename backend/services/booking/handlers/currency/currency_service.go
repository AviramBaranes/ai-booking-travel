package currency

import (
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/booking/db"
)

// CurrencyService provides business logic for managing currencies.
type CurrencyService struct {
	query db.Querier
}

// NewCurrencyService creates a new CurrencyService.
func NewCurrencyService(query db.Querier) *CurrencyService {
	return &CurrencyService{query: query}
}

// CurrencyResponse represents a currency returned by the service.
type CurrencyResponse struct {
	ID              int64   `json:"id"`
	CurrencyCode    string  `json:"currencyCode"`
	CurrencyISOName string  `json:"currencyISOName"`
	Rate            float64 `json:"rate"`
}

// toCurrencyResponse converts a database Currency to a CurrencyResponse.
func toCurrencyResponse(c db.Currency) CurrencyResponse {
	return CurrencyResponse{
		ID:              c.ID,
		CurrencyCode:    c.CurrencyCode,
		CurrencyISOName: c.CurrencyIsoName,
		Rate:            dbadapters.NumericToFloat64(c.Rate),
	}
}
