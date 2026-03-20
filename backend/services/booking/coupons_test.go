package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/booking/db"
	"encore.dev/beta/errs"
	"go.uber.org/mock/gomock"
)

// --- Helpers ---

func validCreateCouponParams() CreateCouponRequest {
	return CreateCouponRequest{
		Name:      "Summer Sale",
		Code:      "SUMMER2026",
		Discount:  15,
		IsEnabled: true,
	}
}

func validUpdateCouponParams() UpdateCouponRequest {
	name := "Winter Sale"
	code := "WINTER2026"
	discount := int32(25)
	enabled := false
	return UpdateCouponRequest{
		Name:      &name,
		Code:      &code,
		Discount:  &discount,
		IsEnabled: &enabled,
	}
}

func couponInvalidValueErr(field string) error {
	return api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
		Code: api_errors.CodeInvalidValue, Field: field,
	})
}

// createTestCoupon is a shorthand to seed a coupon with a unique code.
func createTestCoupon(t *testing.T, s *Service, code string) *CouponResponse {
	t.Helper()
	p := validCreateCouponParams()
	p.Code = code
	resp, err := s.CreateCoupon(context.Background(), p)
	if err != nil {
		t.Fatalf("failed to seed coupon %s: %v", code, err)
	}
	return resp
}

// --- Tests grouped by endpoint ---

func TestListCoupons(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("returns coupons successfully", func(t *testing.T) {
		c1 := createTestCoupon(t, s, "LIST-A")
		c2 := createTestCoupon(t, s, "LIST-B")
		c3 := createTestCoupon(t, s, "LIST-C")

		resp, err := s.ListCoupons(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		ids := make(map[int32]bool)
		for _, c := range resp.Coupons {
			ids[c.ID] = true
		}
		for _, want := range []*CouponResponse{c1, c2, c3} {
			if !ids[want.ID] {
				t.Fatalf("expected coupon %d (%s) in list", want.ID, want.Code)
			}
		}
	})

	t.Run("returns empty list when no coupons exist", func(t *testing.T) {
		// Use a mock to simulate an empty database
		q, s := mockService(t)
		q.EXPECT().ListCoupons(gomock.Any()).Return([]db.Coupon{}, nil)

		resp, err := s.ListCoupons(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Coupons) != 0 {
			t.Fatalf("expected 0 coupons, got %d", len(resp.Coupons))
		}
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().ListCoupons(gomock.Any()).Return(nil, errors.New("db error"))

		_, err := s.ListCoupons(ctx)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestCreateCoupon(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects missing name", func(t *testing.T) {
		p := validCreateCouponParams()
		p.Name = ""
		api_errors.AssertApiError(t, couponInvalidValueErr("name"), p.Validate())
	})

	t.Run("validation rejects blank name", func(t *testing.T) {
		p := validCreateCouponParams()
		p.Name = "   "
		api_errors.AssertApiError(t, couponInvalidValueErr("name"), p.Validate())
	})

	t.Run("validation rejects missing code", func(t *testing.T) {
		p := validCreateCouponParams()
		p.Code = ""
		api_errors.AssertApiError(t, couponInvalidValueErr("code"), p.Validate())
	})

	t.Run("validation rejects blank code", func(t *testing.T) {
		p := validCreateCouponParams()
		p.Code = "   "
		api_errors.AssertApiError(t, couponInvalidValueErr("code"), p.Validate())
	})

	t.Run("validation rejects discount below 1", func(t *testing.T) {
		p := validCreateCouponParams()
		p.Discount = 0
		api_errors.AssertApiError(t, couponInvalidValueErr("discount"), p.Validate())
	})

	t.Run("validation rejects discount above 100", func(t *testing.T) {
		p := validCreateCouponParams()
		p.Discount = 101
		api_errors.AssertApiError(t, couponInvalidValueErr("discount"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validCreateCouponParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("creates coupon successfully", func(t *testing.T) {
		resp, err := s.CreateCoupon(ctx, CreateCouponRequest{
			Name:      "Test Coupon",
			Code:      "CREATE-OK",
			Discount:  20,
			IsEnabled: true,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.ID == 0 {
			t.Fatal("expected non-zero ID")
		}
		if resp.Name != "Test Coupon" {
			t.Fatalf("expected name 'Test Coupon', got %q", resp.Name)
		}
		if resp.Code != "CREATE-OK" {
			t.Fatalf("expected code 'CREATE-OK', got %q", resp.Code)
		}
		if resp.Discount != 20 {
			t.Fatalf("expected discount 20, got %d", resp.Discount)
		}
		if !resp.IsEnabled {
			t.Fatal("expected isEnabled to be true")
		}
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().CreateCoupon(gomock.Any(), gomock.Any()).Return(db.Coupon{}, errors.New("db error"))

		_, err := s.CreateCoupon(ctx, validCreateCouponParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestUpdateCoupon(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("validation rejects blank name", func(t *testing.T) {
		p := validUpdateCouponParams()
		blank := "   "
		p.Name = &blank
		api_errors.AssertApiError(t, couponInvalidValueErr("name"), p.Validate())
	})

	t.Run("validation rejects blank code", func(t *testing.T) {
		p := validUpdateCouponParams()
		blank := "   "
		p.Code = &blank
		api_errors.AssertApiError(t, couponInvalidValueErr("code"), p.Validate())
	})

	t.Run("validation rejects discount below 1", func(t *testing.T) {
		p := validUpdateCouponParams()
		d := int32(0)
		p.Discount = &d
		api_errors.AssertApiError(t, couponInvalidValueErr("discount"), p.Validate())
	})

	t.Run("validation rejects discount above 100", func(t *testing.T) {
		p := validUpdateCouponParams()
		d := int32(101)
		p.Discount = &d
		api_errors.AssertApiError(t, couponInvalidValueErr("discount"), p.Validate())
	})

	t.Run("validation accepts valid params", func(t *testing.T) {
		if err := validUpdateCouponParams().Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("validation accepts partial update with only name", func(t *testing.T) {
		name := "Partial"
		p := UpdateCouponRequest{Name: &name}
		if err := p.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("updates coupon successfully", func(t *testing.T) {
		created := createTestCoupon(t, s, "UPDATE-OK")

		resp, err := s.UpdateCoupon(ctx, created.ID, validUpdateCouponParams())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != "Winter Sale" {
			t.Fatalf("expected name 'Winter Sale', got %q", resp.Name)
		}
		if resp.Code != "WINTER2026" {
			t.Fatalf("expected code 'WINTER2026', got %q", resp.Code)
		}
		if resp.Discount != 25 {
			t.Fatalf("expected discount 25, got %d", resp.Discount)
		}
		if resp.IsEnabled {
			t.Fatal("expected isEnabled to be false")
		}
	})

	t.Run("partial update only changes provided fields", func(t *testing.T) {
		created := createTestCoupon(t, s, "PARTIAL-UPD")

		newName := "Updated Name"
		resp, err := s.UpdateCoupon(ctx, created.ID, UpdateCouponRequest{Name: &newName})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != "Updated Name" {
			t.Fatalf("expected name 'Updated Name', got %q", resp.Name)
		}
		// Other fields should remain unchanged
		if resp.Code != "PARTIAL-UPD" {
			t.Fatalf("expected code unchanged 'PARTIAL-UPD', got %q", resp.Code)
		}
		if resp.Discount != created.Discount {
			t.Fatalf("expected discount unchanged %d, got %d", created.Discount, resp.Discount)
		}
	})

	t.Run("returns not found when coupon does not exist", func(t *testing.T) {
		_, err := s.UpdateCoupon(ctx, 999999, validUpdateCouponParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().UpdateCoupon(gomock.Any(), gomock.Any()).Return(db.Coupon{}, errors.New("db error"))

		_, err := s.UpdateCoupon(ctx, 1, validUpdateCouponParams())
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

func TestDeleteCoupon(t *testing.T) {
	ctx := context.Background()
	s := &Service{query: testQuerier()}

	t.Run("deletes coupon successfully", func(t *testing.T) {
		created := createTestCoupon(t, s, "DELETE-OK")

		if err := s.DeleteCoupon(ctx, created.ID); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify it's gone by trying to update it
		_, err := s.UpdateCoupon(ctx, created.ID, validUpdateCouponParams())
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error when db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteCoupon(gomock.Any(), int32(1)).Return(errors.New("db error"))

		api_errors.AssertApiError(t, api_errors.ErrInternalError, s.DeleteCoupon(ctx, 1))
	})
}
