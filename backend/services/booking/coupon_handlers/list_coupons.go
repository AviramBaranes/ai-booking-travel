package coupon_handlers

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.dev/rlog"
)

type ListCouponsResponse struct {
	Coupons []CouponResponse `json:"coupons"`
}

func (s *CouponService) ListCoupons(ctx context.Context) (*ListCouponsResponse, error) {
	rows, err := s.query.ListCoupons(ctx)
	if err != nil {
		rlog.Error("failed to list coupons", "error", err)
		return nil, api_errors.ErrInternalError
	}

	coupons := make([]CouponResponse, 0, len(rows))
	for _, r := range rows {
		coupons = append(coupons, toCouponResponse(r))
	}

	return &ListCouponsResponse{Coupons: coupons}, nil
}
