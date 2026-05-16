package coupon_handlers

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

type UpdateCouponParams struct {
	Name      *string `json:"name" validate:"omitempty,notblank" encore:"optional"`
	Code      *string `json:"code" validate:"omitempty,notblank" encore:"optional"`
	Discount  *int32  `json:"discount" validate:"omitempty,gte=1,lte=100" encore:"optional"`
	IsEnabled *bool   `json:"isEnabled" encore:"optional"`
}

func (p UpdateCouponParams) Validate() error {
	return validation.ValidateStruct(p)
}

func (s *CouponService) UpdateCoupon(ctx context.Context, id int64, params UpdateCouponParams) (*CouponResponse, error) {
	row, err := s.query.UpdateCoupon(ctx, db.UpdateCouponParams{
		ID:        id,
		Name:      params.Name,
		Code:      params.Code,
		Discount:  params.Discount,
		IsEnabled: params.IsEnabled,
	})
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, api_errors.ErrNotFound
		}
		rlog.Error("failed to update coupon", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toCouponResponse(row)
	return &resp, nil
}
