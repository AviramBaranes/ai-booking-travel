package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// getCouponDiscount retrieves the discount percentage for a given coupon code. If the code is empty, it returns 0.
func (s *Service) getCouponDiscount(ctx context.Context, code string) (int, error) {
	if code == "" {
		return 0, nil
	}

	coupon, err := s.query.FindCouponByCode(ctx, code)
	if err != nil && !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to find coupon by code", "error", err, "code", code)
		return 0, api_errors.ErrInternalError
	}

	return int(coupon.Discount), nil
}
