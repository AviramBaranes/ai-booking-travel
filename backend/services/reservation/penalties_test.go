package reservation

import (
	"context"
	"testing"

	"encore.app/internal/api_errors"
	dbadapters "encore.app/internal/db_adapters"
	"encore.app/services/reservation/db"
)

// --- Helpers ---

func validCreatePenaltyParams(reservationID int64) CreatePenaltyParams {
	return CreatePenaltyParams{
		ReservationID: reservationID,
		Type:          string(db.PenaltyTypeCancellation),
		Amount:        150.50,
	}
}

// seedUSDRate primes the currency cache so the handler resolves the rate from it instead of
// calling the rates provider. The cache is isolated per test, so every test needs its own seed.
func seedUSDRate(t *testing.T, ctx context.Context) float64 {
	t.Helper()

	const rate = 3.65
	if err := currenciesRates.Set(ctx, "USD", rate); err != nil {
		t.Fatalf("failed to seed USD rate: %v", err)
	}
	return rate
}

// seedCanceledReservation creates a reservation owned by userID and cancels it, so a penalty
// can be charged against it.
func seedCanceledReservation(t *testing.T, ctx context.Context, userID int64) int64 {
	t.Helper()

	params := validCreateReservationParams()
	params.UserID = userID
	params.PickupDate = futurePickup()
	res, err := CreateReservation(ctx, *params)
	if err != nil {
		t.Fatalf("failed to create reservation: %v", err)
	}

	if err := CancelReservation(authContext(userID), res.ID); err != nil {
		t.Fatalf("failed to cancel reservation: %v", err)
	}

	return res.ID
}

// --- Tests ---

func TestCreatePenaltyValidation(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(p *CreatePenaltyParams)
		wantErr error
	}{
		{
			name:    "rejects missing reservation id",
			modify:  func(p *CreatePenaltyParams) { p.ReservationID = 0 },
			wantErr: invalidValueErr("reservationId"),
		},
		{
			name:    "rejects negative reservation id",
			modify:  func(p *CreatePenaltyParams) { p.ReservationID = -1 },
			wantErr: invalidValueErr("reservationId"),
		},
		{
			name:    "rejects missing type",
			modify:  func(p *CreatePenaltyParams) { p.Type = "" },
			wantErr: invalidValueErr("type"),
		},
		{
			name:    "rejects unknown type",
			modify:  func(p *CreatePenaltyParams) { p.Type = "late_return" },
			wantErr: invalidValueErr("type"),
		},
		{
			name:    "rejects zero amount",
			modify:  func(p *CreatePenaltyParams) { p.Amount = 0 },
			wantErr: invalidValueErr("amount"),
		},
		{
			name:    "rejects negative amount",
			modify:  func(p *CreatePenaltyParams) { p.Amount = -1 },
			wantErr: invalidValueErr("amount"),
		},
		{
			name:    "accepts a no-show penalty",
			modify:  func(p *CreatePenaltyParams) { p.Type = string(db.PenaltyTypeNoShow) },
			wantErr: nil,
		},
		{
			name:    "accepts valid params",
			modify:  func(p *CreatePenaltyParams) {},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validCreatePenaltyParams(1)
			tt.modify(&p)

			err := p.Validate()
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			api_errors.AssertApiError(t, tt.wantErr, err)
		})
	}
}

func TestCreatePenalty(t *testing.T) {
	ctx := context.Background()

	t.Run("returns 404 for non-existent reservation", func(t *testing.T) {
		_, err := CreatePenalty(authContext(1), validCreatePenaltyParams(99999999))
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("rejects a penalty on a reservation that is not canceled", func(t *testing.T) {
		const userID int64 = 2001
		params := validCreateReservationParams()
		params.UserID = userID
		params.PickupDate = futurePickup()
		res, err := CreateReservation(ctx, *params)
		if err != nil {
			t.Fatalf("failed to create reservation: %v", err)
		}

		_, err = CreatePenalty(authContext(userID), validCreatePenaltyParams(res.ID))
		api_errors.AssertApiError(t, ErrReservationNotCanceled, err)
	})

	t.Run("records the penalty against the canceled reservation", func(t *testing.T) {
		const userID int64 = 2002
		rate := seedUSDRate(t, ctx)
		reservationID := seedCanceledReservation(t, ctx, userID)

		p := validCreatePenaltyParams(reservationID)
		got, err := CreatePenalty(authContext(userID), p)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.ReservationID != reservationID {
			t.Fatalf("expected reservation id %d, got %d", reservationID, got.ReservationID)
		}
		if got.Type != p.Type {
			t.Fatalf("expected type %q, got %q", p.Type, got.Type)
		}
		if got.Amount != p.Amount {
			t.Fatalf("expected amount %v, got %v", p.Amount, got.Amount)
		}
		// The currency is taken from the reservation, and the rate from the cache.
		if got.CurrencyCode != "USD" {
			t.Fatalf("expected currency code %q, got %q", "USD", got.CurrencyCode)
		}
		if got.CurrencyRate != rate {
			t.Fatalf("expected currency rate %v, got %v", rate, got.CurrencyRate)
		}

		row, err := query.GetReservationPenaltyByReservationID(ctx, reservationID)
		if err != nil {
			t.Fatalf("failed to read back penalty: %v", err)
		}
		if row.ID != got.ID {
			t.Fatalf("expected persisted id %d, got %d", got.ID, row.ID)
		}
		if row.PenaltyType != db.PenaltyTypeCancellation {
			t.Fatalf("expected persisted type %q, got %q", db.PenaltyTypeCancellation, row.PenaltyType)
		}
		if amount := dbadapters.NumericToFloat64(row.Amount); amount != p.Amount {
			t.Fatalf("expected persisted amount %v, got %v", p.Amount, amount)
		}
		if row.CreatedByUserID == nil || *row.CreatedByUserID != userID {
			t.Fatalf("expected persisted created_by_user_id %d, got %v", userID, row.CreatedByUserID)
		}
		if row.PaidAt.Valid {
			t.Fatalf("expected paid_at to be NULL, got %v", row.PaidAt.Time)
		}
		if row.SupplierPaidAt.Valid {
			t.Fatalf("expected supplier_paid_at to be NULL, got %v", row.SupplierPaidAt.Time)
		}
		if row.InvoiceDocNum != nil || row.PaymentDocNum != nil {
			t.Fatalf("expected doc numbers to be NULL, got %v and %v", row.InvoiceDocNum, row.PaymentDocNum)
		}
	})

	t.Run("rejects a second penalty on the same reservation", func(t *testing.T) {
		const userID int64 = 2003
		seedUSDRate(t, ctx)
		reservationID := seedCanceledReservation(t, ctx, userID)

		if _, err := CreatePenalty(authContext(userID), validCreatePenaltyParams(reservationID)); err != nil {
			t.Fatalf("failed to create the first penalty: %v", err)
		}

		p := validCreatePenaltyParams(reservationID)
		p.Type = string(db.PenaltyTypeNoShow)
		_, err := CreatePenalty(authContext(userID), p)
		api_errors.AssertApiError(t, ErrPenaltyAlreadyExists, err)
	})
}
