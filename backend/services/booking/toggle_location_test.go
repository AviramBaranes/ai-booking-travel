package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"go.uber.org/mock/gomock"
)

func TestToggleLocation(t *testing.T) {
	ctx := context.Background()
	q := testQuerier()
	s := &Service{query: q}

	// Helper to read back the enabled state of a broker code.
	getEnabled := func(t *testing.T, lbc db.LocationBrokerCode) bool {
		t.Helper()
		row, err := q.GetLocationBrokerCode(ctx, db.GetLocationBrokerCodeParams{
			Broker:           lbc.Broker,
			BrokerLocationID: lbc.BrokerLocationID,
			LocationID:       lbc.LocationID,
		})
		if err != nil {
			t.Fatalf("failed to fetch broker code: %v", err)
		}
		return row.Enabled
	}

	t.Run("disables an enabled location", func(t *testing.T) {
		_, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Toggle Disable Test",
			},
			db.BrokerFlex, "flex-toggle-disable",
		)

		if !getEnabled(t, lbc) {
			t.Fatal("expected location to be enabled by default")
		}

		err := s.ToggleLocation(ctx, lbc.ID, &ToggleLocationRequest{Enabled: false})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if getEnabled(t, lbc) {
			t.Fatal("expected location to be disabled after toggle")
		}
	})

	t.Run("enables a disabled location", func(t *testing.T) {
		_, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Toggle Enable Test",
			},
			db.BrokerFlex, "flex-toggle-enable",
		)

		// Disable first
		if err := q.DisableLocationBrokerCode(ctx, lbc.ID); err != nil {
			t.Fatalf("failed to disable: %v", err)
		}
		if getEnabled(t, lbc) {
			t.Fatal("expected location to be disabled")
		}

		err := s.ToggleLocation(ctx, lbc.ID, &ToggleLocationRequest{Enabled: true})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !getEnabled(t, lbc) {
			t.Fatal("expected location to be enabled after toggle")
		}
	})

	t.Run("enabling an already enabled location is idempotent", func(t *testing.T) {
		_, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Toggle Idempotent Test",
			},
			db.BrokerFlex, "flex-toggle-idempotent",
		)

		if !getEnabled(t, lbc) {
			t.Fatal("expected location to be enabled by default")
		}

		err := s.ToggleLocation(ctx, lbc.ID, &ToggleLocationRequest{Enabled: true})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !getEnabled(t, lbc) {
			t.Fatal("expected location to still be enabled")
		}
	})

	t.Run("disabling an already disabled location is idempotent", func(t *testing.T) {
		_, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", Name: "Toggle Idempotent Disable",
			},
			db.BrokerFlex, "flex-toggle-idempotent-dis",
		)

		if err := q.DisableLocationBrokerCode(ctx, lbc.ID); err != nil {
			t.Fatalf("failed to disable: %v", err)
		}

		err := s.ToggleLocation(ctx, lbc.ID, &ToggleLocationRequest{Enabled: false})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if getEnabled(t, lbc) {
			t.Fatal("expected location to still be disabled")
		}
	})

	t.Run("returns error when enable fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().EnableLocationBrokerCode(gomock.Any(), int64(999)).
			Return(errors.New("db error"))

		err := s.ToggleLocation(ctx, 999, &ToggleLocationRequest{Enabled: true})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when disable fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DisableLocationBrokerCode(gomock.Any(), int64(999)).
			Return(errors.New("db error"))

		err := s.ToggleLocation(ctx, 999, &ToggleLocationRequest{Enabled: false})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
