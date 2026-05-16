package currency_handlers

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// UpdateCurrencyParams contains the parameters for updating a currency.
type UpdateCurrencyParams struct {
	CurrencyCode    *string  `json:"currencyCode" validate:"omitempty,notblank" encore:"optional"`
	CurrencyISOName *string  `json:"currencyISOName" validate:"omitempty,notblank" encore:"optional"`
	Rate            *float64 `json:"rate" validate:"omitempty,gt=0" encore:"optional"`
}

func (p UpdateCurrencyParams) Validate() error {
	return validation.ValidateStruct(p)
}

// UpdateCurrency updates an existing currency.
func (s *CurrencyService) UpdateCurrency(ctx context.Context, id int64, params UpdateCurrencyParams) (*CurrencyResponse, error) {
	dbParams := db.UpdateCurrencyParams{
		ID:              id,
		CurrencyCode:    params.CurrencyCode,
		CurrencyIsoName: params.CurrencyISOName,
	}
	if params.Rate != nil {
		dbParams.Rate = dbadapters.NumericFromFloat64(*params.Rate)
	}

	row, err := s.query.UpdateCurrency(ctx, dbParams)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to update currency", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toCurrencyResponse(row)
	return &resp, nil
}
