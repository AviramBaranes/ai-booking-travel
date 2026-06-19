package organization

import (
	"context"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
)

type UpdateOrganizationBalanceDueParams struct {
	ID            int64   `json:"id" validate:"required,gt=0"`
	BalanceChange float64 `json:"balanceChange" validate:"required"`
}

// UpdateOrganizationBalanceDue updates the balance due for an organization work for plus or minus the current balance due.
func (s *OrganizationService) UpdateOrganizationBalanceDue(ctx context.Context, p UpdateOrganizationBalanceDueParams) error {
	if err := s.query.UpdateOrganizationBalanceDue(ctx, db.UpdateOrganizationBalanceDueParams{
		ID:         p.ID,
		BalanceDue: dbadapters.NumericFromFloat64(p.BalanceChange),
	}); err != nil {
		rlog.Error("failed to update organization balance due", "error", err)
		return api_errors.ErrInternalError
	}

	return nil
}
