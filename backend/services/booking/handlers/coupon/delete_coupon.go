package coupon

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
)

func (s *CouponService) DeleteCoupon(ctx context.Context, id int64) error {
	err := s.query.DeleteCoupon(ctx, id)
	if err != nil {
		rlog.Error("failed to delete coupon", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
