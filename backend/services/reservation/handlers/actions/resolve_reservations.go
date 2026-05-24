package actions

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
)

type ResolveReservationsParams struct {
	IDs []int64 `json:"ids" validate:"required,min=1"`
}

func (s *ActionService) ResolveReservations(ctx context.Context, p ResolveReservationsParams) error {
	err := s.query.ResolveReservationsPayment(ctx, p.IDs)
	if err != nil {
		rlog.Error("failed to resolve reservations", "error", err, "reservation_ids", p.IDs)
		return api_errors.ErrInternalError
	}

	return nil
}
