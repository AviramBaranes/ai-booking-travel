package reservation

import (
	"context"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/internal/validation"
	"encore.app/services/accounts"
	"encore.app/services/accounts/handlers/office"
	"encore.app/services/accounts/handlers/organization"
	"encore.app/services/accounts/handlers/user"
	actions "encore.app/services/reservation/handlers/actions"
	"encore.dev/beta/errs"
	"encore.dev/et"
)

func TestApplyVoucher(t *testing.T) {
	t.Run("it validates voucher id exists", func(t *testing.T) {
		r := ApplyVoucherParams{}
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
		var userID int64 = 123
		var reservationID int64 = 456
		ctx := authContext(userID)
		err := ApplyVoucher(ctx, reservationID, ApplyVoucherParams{
			Voucher: "123",
		})
		if err == nil {
			t.Fatal("expected applying voucher to fail")
		}

		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("it returns not found if id exists but reservation doesn't belong to user", func(t *testing.T) {
		var authenticatedUserID int64 = 123
		ctx := authContext(authenticatedUserID)
		res, err := CreateReservation(context.Background(), *validCreateReservationParams())

		if err != nil {
			t.Fatalf("creating reservation failed: %v", err)
		}

		err = ApplyVoucher(ctx, res.ID, ApplyVoucherParams{Voucher: "123"})
		if err == nil {
			t.Fatal("expected applying voucher to fail")
		}

		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("it updates the reservation with the voucher and current timestamp", func(t *testing.T) {
		var authenticatedUserID int64 = 123
		et.MockEndpoint(accounts.GetUserCredit, func(ctx context.Context) (*user.GetUserCreditResponse, error) {
			return &user.GetUserCreditResponse{
				Obligo:     1000,
				BalanceDue: 0,
			}, nil
		})

		et.MockEndpoint(accounts.UpdateOfficeBalanceDue, func(ctx context.Context, p office.UpdateOfficeBalanceDueParams) error {
			return nil
		})

		authCtx := authContextWithOrgCtx(authenticatedUserID, 1, 1, false)
		ctx := authContext(authenticatedUserID)
		reservation := validCreateReservationParams()
		reservation.UserID = authenticatedUserID
		res, err := CreateReservation(context.Background(), *reservation)

		if err != nil {
			t.Fatalf("creating reservation failed: %v", err)
		}
		if err := currenciesRates.Set(ctx, "USD", 3.65); err != nil {
			t.Fatalf("failed to seed USD rate: %v", err)
		}
		err = ApplyVoucher(authCtx, res.ID, ApplyVoucherParams{Voucher: "123"})
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

	t.Run("it fails if user has no credit", func(t *testing.T) {
		tests := []struct {
			name        string
			resp        *user.GetUserCreditResponse
			expectedErr error
		}{
			{
				name: "zero obligo",
				resp: &user.GetUserCreditResponse{
					Obligo:     0,
					BalanceDue: 0,
				},
				expectedErr: actions.ErrNoObligo,
			},
			{
				name: "balance due exceeds obligo",
				resp: &user.GetUserCreditResponse{
					Obligo:     150,
					BalanceDue: 100,
				},
				expectedErr: actions.ErrNotEnoughCredits,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var authenticatedUserID int64 = 123
				ctx := authContext(authenticatedUserID)
				reservation := validCreateReservationParams()
				reservation.PurchasePrice = 10000 // set a high price to ensure credit check fails
				reservation.UserID = authenticatedUserID
				res, err := CreateReservation(context.Background(), *reservation)

				if err != nil {
					t.Fatalf("creating reservation failed: %v", err)
				}

				et.MockEndpoint(accounts.GetUserCredit, func(ctx context.Context) (*user.GetUserCreditResponse, error) {
					return tt.resp, nil
				})
				if err := currenciesRates.Set(ctx, "USD", 3.65); err != nil {
					t.Fatalf("failed to seed USD rate: %v", err)
				}
				err = ApplyVoucher(ctx, res.ID, ApplyVoucherParams{Voucher: "123"})
				if err == nil {
					t.Fatal("expected applying voucher to fail")
				}

				api_errors.AssertApiError(t, tt.expectedErr, err)
			})

		}
	})

	t.Run("it updates the correct billing entity balance after successful vouchering", func(t *testing.T) {
		tests := []struct {
			name      string
			orgID     int64
			officeID  int64
			isOrganic bool
		}{
			{
				name:      "updates organization balance for organic org",
				orgID:     10,
				officeID:  20,
				isOrganic: true,
			},
			{
				name:      "updates office balance for non-organic org",
				orgID:     10,
				officeID:  20,
				isOrganic: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var authenticatedUserID int64 = 123

				et.MockEndpoint(accounts.GetUserCredit, func(ctx context.Context) (*user.GetUserCreditResponse, error) {
					return &user.GetUserCreditResponse{
						Obligo:     1000,
						BalanceDue: 0,
					}, nil
				})

				orgBalanceUpdated := false
				officeBalanceUpdated := false

				et.MockEndpoint(accounts.UpdateOrganizationBalanceDue, func(ctx context.Context, p organization.UpdateOrganizationBalanceDueParams) error {
					orgBalanceUpdated = true
					return nil
				})

				et.MockEndpoint(accounts.UpdateOfficeBalanceDue, func(ctx context.Context, p office.UpdateOfficeBalanceDueParams) error {
					officeBalanceUpdated = true
					return nil
				})

				authCtx := authContextWithOrgCtx(authenticatedUserID, tt.orgID, tt.officeID, tt.isOrganic)
				reservation := validCreateReservationParams()
				reservation.UserID = authenticatedUserID
				res, err := CreateReservation(context.Background(), *reservation)
				if err != nil {
					t.Fatalf("creating reservation failed: %v", err)
				}

				if err := currenciesRates.Set(authCtx, "USD", 3.65); err != nil {
					t.Fatalf("failed to seed USD rate: %v", err)
				}

				err = ApplyVoucher(authCtx, res.ID, ApplyVoucherParams{Voucher: "123"})
				if err != nil {
					t.Fatalf("expected applying voucher to succeed, got error: %v", err)
				}

				if tt.isOrganic {
					if !orgBalanceUpdated {
						t.Fatal("expected organization balance to be updated for organic org")
					}
					if officeBalanceUpdated {
						t.Fatal("expected office balance NOT to be updated for organic org")
					}
				} else {
					if !officeBalanceUpdated {
						t.Fatal("expected office balance to be updated for non-organic org")
					}
					if orgBalanceUpdated {
						t.Fatal("expected organization balance NOT to be updated for non-organic org")
					}
				}
			})
		}
	})
}
