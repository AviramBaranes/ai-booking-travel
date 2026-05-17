package coupon

import (
	"context"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type CreateCouponParams struct {
	Name      string `json:"name" validate:"required,notblank"`
	Code      string `json:"code" validate:"required,notblank"`
	Discount  int32  `json:"discount" validate:"required,gte=1,lte=100"`
	IsEnabled bool   `json:"isEnabled"`
}

func (p CreateCouponParams) Validate() error {
	return validation.ValidateStruct(p)
}

type CouponResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Discount  int32  `json:"discount"`
	IsEnabled bool   `json:"isEnabled"`
}

// CreateCoupon creates a new coupon.
func (s *CouponService) CreateCoupon(ctx context.Context, p CreateCouponParams) (*CouponResponse, error) {
	row, err := s.query.CreateCoupon(ctx, db.CreateCouponParams{
		Name:      p.Name,
		Code:      p.Code,
		Discount:  p.Discount,
		IsEnabled: p.IsEnabled,
	})
	if err != nil {
		rlog.Error("failed to create coupon", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toCouponResponse(row)
	return &resp, nil
}

// toCouponResponse converts a db.Coupon to a CouponResponse.
func toCouponResponse(c db.Coupon) CouponResponse {
	return CouponResponse{
		ID:        c.ID,
		Name:      c.Name,
		Code:      c.Code,
		Discount:  c.Discount,
		IsEnabled: c.IsEnabled,
	}
}
