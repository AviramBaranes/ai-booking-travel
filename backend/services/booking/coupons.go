package booking

import (
	"context"

	"encore.app/services/booking/handlers/coupon_handlers"
)

// ListCoupons lists all coupons.
//
//encore:api auth method=GET path=/coupons tag:admin
func (s *Service) ListCoupons(ctx context.Context) (*coupon_handlers.ListCouponsResponse, error) {
	cs := coupon_handlers.NewCouponService(s.query)
	return cs.ListCoupons(ctx)
}

// CreateCoupon creates a new coupon.
//
//encore:api auth method=POST path=/coupons tag:admin
func (s *Service) CreateCoupon(ctx context.Context, p coupon_handlers.CreateCouponParams) (*coupon_handlers.CouponResponse, error) {
	cs := coupon_handlers.NewCouponService(s.query)
	return cs.CreateCoupon(ctx, p)
}

// UpdateCoupon updates an existing coupon.
//
//encore:api auth method=PUT path=/coupons/:id tag:admin
func (s *Service) UpdateCoupon(ctx context.Context, id int64, params coupon_handlers.UpdateCouponParams) (*coupon_handlers.CouponResponse, error) {
	cs := coupon_handlers.NewCouponService(s.query)
	return cs.UpdateCoupon(ctx, id, params)
}

// DeleteCoupon deletes a coupon by its ID.
//
//encore:api auth method=DELETE path=/coupons/:id tag:admin
func (s *Service) DeleteCoupon(ctx context.Context, id int64) error {
	cs := coupon_handlers.NewCouponService(s.query)
	return cs.DeleteCoupon(ctx, id)
}
