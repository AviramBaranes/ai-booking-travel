package actions

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
)

type ResolvePenaltiesParams struct {
	IDs []int64 `json:"ids" validate:"required,min=1"`
}

// ResolvePenalties marks the given penalties as collected from the customer. Only paid_at is set:
// billing settles a batch of items with a single invoice and receipt, so there is no per-penalty
// document number to record here.
func (s *ActionService) ResolvePenalties(ctx context.Context, p ResolvePenaltiesParams) error {
	err := s.query.ResolvePenaltiesPayment(ctx, p.IDs)
	if err != nil {
		rlog.Error("failed to resolve penalties", "error", err, "penalty_ids", p.IDs)
		return api_errors.ErrInternalError
	}

	return nil
}
