package currency

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
)

// DeleteCurrency removes a currency from the database by its ID.
func (s *CurrencyService) DeleteCurrency(ctx context.Context, id int64) error {
	err := s.query.DeleteCurrency(ctx, id)
	if err != nil {
		rlog.Error("failed to delete currency", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
