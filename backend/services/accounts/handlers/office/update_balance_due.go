package office

import (
	"context"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type UpdateOfficeBalanceDueParams struct {
	ID            int64   `json:"id" validate:"required,gt=0"`
	BalanceChange float64 `json:"balanceChange" validate:"required"`
}

// UpdateOfficeBalanceDue updates the balance due for an office work for plus or minus the current balance due.
func (s *OfficeService) UpdateOfficeBalanceDue(ctx context.Context, p UpdateOfficeBalanceDueParams) error {
	if err := s.query.UpdateOfficeBalanceDue(ctx, db.UpdateOfficeBalanceDueParams{
		ID:    p.ID,
		Delta: dbadapters.NumericFromFloat64(p.BalanceChange),
	}); err != nil {
		rlog.Error("failed to update office balance due", "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
