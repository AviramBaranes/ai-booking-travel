package currency_handlers

import (
	"context"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// CreateCurrencyParams contains the parameters for creating a new currency.
type CreateCurrencyParams struct {
	CurrencyCode    string  `json:"currencyCode" validate:"required,notblank"`
	CurrencyISOName string  `json:"currencyISOName" validate:"required,notblank"`
	Rate            float64 `json:"rate" validate:"required,gt=0"`
}

func (p CreateCurrencyParams) Validate() error {
	return validation.ValidateStruct(p)
}

// CreateCurrency creates a new currency in the database.
func (s *CurrencyService) CreateCurrency(ctx context.Context, params CreateCurrencyParams) (*CurrencyResponse, error) {
	row, err := s.query.CreateCurrency(ctx, db.CreateCurrencyParams{
		CurrencyCode:    params.CurrencyCode,
		CurrencyIsoName: params.CurrencyISOName,
		Rate:            dbadapters.NumericFromFloat64(params.Rate),
	})
	if err != nil {
		rlog.Error("failed to create currency", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toCurrencyResponse(row)
	return &resp, nil
}
