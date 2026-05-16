package coupon_handlers

import (
	"encore.app/services/booking/db"
)

type CouponService struct {
	query db.Querier
}

func NewCouponService(query db.Querier) *CouponService {
	return &CouponService{query: query}
}
