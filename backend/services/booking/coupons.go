package booking

import (
	"context"
	"errors"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/rlog"
)

// --- Request / Response types ---

type CouponResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Discount  int32  `json:"discount"`
	IsEnabled bool   `json:"isEnabled"`
}

type ListCouponsResponse struct {
	Coupons []CouponResponse `json:"coupons"`
}

type CreateCouponRequest struct {
	Name      string `json:"name" validate:"required,notblank"`
	Code      string `json:"code" validate:"required,notblank"`
	Discount  int32  `json:"discount" validate:"required,gte=1,lte=100"`
	IsEnabled bool   `json:"isEnabled"`
}

func (p CreateCouponRequest) Validate() error {
	return validation.ValidateStruct(p)
}

type UpdateCouponRequest struct {
	Name      *string `json:"name" validate:"omitempty,notblank" encore:"optional"`
	Code      *string `json:"code" validate:"omitempty,notblank" encore:"optional"`
	Discount  *int32  `json:"discount" validate:"omitempty,gte=1,lte=100" encore:"optional"`
	IsEnabled *bool   `json:"isEnabled" encore:"optional"`
}

func (p UpdateCouponRequest) Validate() error {
	return validation.ValidateStruct(p)
}

// --- Helpers ---

func toCouponResponse(c db.Coupon) CouponResponse {
	return CouponResponse{
		ID:        c.ID,
		Name:      c.Name,
		Code:      c.Code,
		Discount:  c.Discount,
		IsEnabled: c.IsEnabled,
	}
}

// --- Endpoints ---

// ListCoupons lists all coupons.
//
//encore:api auth method=GET path=/coupons tag:admin
func (s *Service) ListCoupons(ctx context.Context) (*ListCouponsResponse, error) {
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

// CreateCoupon creates a new coupon.
//
//encore:api auth method=POST path=/coupons tag:admin
func (s *Service) CreateCoupon(ctx context.Context, params CreateCouponRequest) (*CouponResponse, error) {
	row, err := s.query.CreateCoupon(ctx, db.CreateCouponParams{
		Name:      params.Name,
		Code:      params.Code,
		Discount:  params.Discount,
		IsEnabled: params.IsEnabled,
	})
	if err != nil {
		rlog.Error("failed to create coupon", "error", err)
		return nil, api_errors.ErrInternalError
	}

	resp := toCouponResponse(row)
	return &resp, nil
}

// UpdateCoupon updates an existing coupon.
//
//encore:api auth method=PUT path=/coupons/:id tag:admin
func (s *Service) UpdateCoupon(ctx context.Context, id int32, params UpdateCouponRequest) (*CouponResponse, error) {
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

// DeleteCoupon deletes a coupon by its ID.
//
//encore:api auth method=DELETE path=/coupons/:id tag:admin
func (s *Service) DeleteCoupon(ctx context.Context, id int32) error {
	err := s.query.DeleteCoupon(ctx, id)
	if err != nil {
		rlog.Error("failed to delete coupon", "error", err, "id", id)
		return api_errors.ErrInternalError
	}
	return nil
}
