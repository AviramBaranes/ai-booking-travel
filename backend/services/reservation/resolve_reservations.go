package reservation

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
)

type ResolveReservationsParams struct {
	IDs []int64 `json:"ids" validate:"required,min=1"`
}

// encore:api private
func (s *Service) ResolveReservations(ctx context.Context, params ResolveReservationsParams) error {
	err := s.query.ResolveReservationsPayment(ctx, params.IDs)
	if err != nil {
		rlog.Error("failed to resolve reservations", "error", err, "reservation_ids", params.IDs)
		return api_errors.ErrInternalError
	}

	return nil
}
