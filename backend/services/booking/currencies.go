package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// --- Request / Response types ---

type CurrencyResponse struct {
	ID              int32   `json:"id"`
	CurrencyCode    string  `json:"currencyCode"`
	CurrencyISOName string  `json:"currencyISOName"`
	Rate            float64 `json:"rate"`
}

type ListCurrenciesResponse struct {
	Currencies []CurrencyResponse `json:"currencies"`
}

type CreateCurrencyRequest struct {
	CurrencyCode    string  `json:"currencyCode" validate:"required,notblank"`
	CurrencyISOName string  `json:"currencyISOName" validate:"required,notblank"`
	Rate            float64 `json:"rate" validate:"required,gt=0"`
}

func (p CreateCurrencyRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type UpdateCurrencyRequest struct {
	CurrencyCode    *string  `json:"currencyCode" validate:"omitempty,notblank" encore:"optional"`
	CurrencyISOName *string  `json:"currencyISOName" validate:"omitempty,notblank" encore:"optional"`
	Rate            *float64 `json:"rate" validate:"omitempty,gt=0" encore:"optional"`
}

func (p UpdateCurrencyRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// --- Helpers ---

func toCurrencyResponse(c db.Currency) CurrencyResponse {
	return CurrencyResponse{
		ID:              c.ID,
		CurrencyCode:    c.CurrencyCode,
		CurrencyISOName: c.CurrencyIsoName,
		Rate:            db.NumericToFloat64(c.Rate),
	}
}

// --- Endpoints ---

// ListCurrencies lists all currencies.
//
//encore:api auth method=GET path=/currencies tag:admin
func (s *Service) ListCurrencies(ctx context.Context) (*ListCurrenciesResponse, error) {
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

// CreateCurrency creates a new currency.
//
//encore:api auth method=POST path=/currencies tag:admin
func (s *Service) CreateCurrency(ctx context.Context, params CreateCurrencyRequest) (*CurrencyResponse, error) {
	row, err := s.query.CreateCurrency(ctx, db.CreateCurrencyParams{
		CurrencyCode:    params.CurrencyCode,
		CurrencyIsoName: params.CurrencyISOName,
		Rate:            db.NumericFromFloat64(params.Rate),
	})
	if err != nil {
		rlog.Error("failed to create currency", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toCurrencyResponse(row)
	return &resp, nil
}

// UpdateCurrency updates an existing currency.
//
//encore:api auth method=PUT path=/currencies/:id tag:admin
func (s *Service) UpdateCurrency(ctx context.Context, id int32, params UpdateCurrencyRequest) (*CurrencyResponse, error) {
	dbParams := db.UpdateCurrencyParams{
		ID:              id,
		CurrencyCode:    params.CurrencyCode,
		CurrencyIsoName: params.CurrencyISOName,
	}
	if params.Rate != nil {
		dbParams.Rate = db.NumericFromFloat64(*params.Rate)
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

// DeleteCurrency deletes a currency by its ID.
//
//encore:api auth method=DELETE path=/currencies/:id tag:admin
func (s *Service) DeleteCurrency(ctx context.Context, id int32) error {
	err := s.query.DeleteCurrency(ctx, id)
	if err != nil {
		rlog.Error("failed to delete currency", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
