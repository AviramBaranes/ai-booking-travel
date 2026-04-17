package reservation

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.dev/beta/errs"
	"encore.dev/et"
	"go.uber.org/mock/gomock"
)

func TestApplyVoucher(t *testing.T) {
	t.Run("it validates voucher id exists", func(t *testing.T) {
		r := ApplyVoucherRequest{}
		err := r.Validate()
		if err == nil {
			t.Fatal("expected validation to fail")
		}

		expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
			Code: api_errors.CodeInvalidValue, Field: "voucher",
		})

		api_errors.AssertApiError(t, expectedErr, err)
	})

	t.Run("it returns not found of id doesn't exists", func(t *testing.T) {
		var userID int32 = 123
		var reservationID int64 = 456
		ctx := authContext(userID)
		err := ApplyVoucher(ctx, reservationID, ApplyVoucherRequest{
			Voucher: "123",
		})
		if err == nil {
			t.Fatal("expected applying voucher to fail")
		}

		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("it returns not found if id exists but reservation doesn't belong to user", func(t *testing.T) {
		var authenticatedUserID int32 = 123
		ctx := authContext(authenticatedUserID)
		res, err := CreateReservation(context.Background(), *validCreateReservationParams())

		if err != nil {
			t.Fatalf("creating reservation failed: %v", err)
		}

		err = ApplyVoucher(ctx, res.ID, ApplyVoucherRequest{Voucher: "123"})
		if err == nil {
			t.Fatal("expected applying voucher to fail")
		}

		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("it send internal err if db fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().ApplyVoucher(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("db error"))
		et.MockService[Interface]("reservation", s)
		ctx := authContext(123)

		err := ApplyVoucher(ctx, 1, ApplyVoucherRequest{
			Voucher: "123",
		})

		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("it updates the reservation with the voucher and current timestamp", func(t *testing.T) {
		var authenticatedUserID int32 = 123
		ctx := authContext(authenticatedUserID)
		reservation := validCreateReservationParams()
		reservation.UserID = authenticatedUserID
		res, err := CreateReservation(context.Background(), *reservation)

		if err != nil {
			t.Fatalf("creating reservation failed: %v", err)
		}

		err = ApplyVoucher(ctx, res.ID, ApplyVoucherRequest{Voucher: "123"})
		if err != nil {
			t.Fatalf("expected applying voucher to succeed, got error: %v", err)
		}

		updatedRes, err := GetReservation(ctx, res.ID)
		if err != nil {
			t.Fatalf("expected getting reservation to succeed, got error: %v", err)
		}

		if *updatedRes.Voucher != "123" {
			t.Fatalf("expected voucher to be '123', got: %v", updatedRes.Voucher)
		}

		if updatedRes.VoucheredAt == nil {
			t.Fatal("expected voucheredAt to be set, got nil")
		}
	})
}
