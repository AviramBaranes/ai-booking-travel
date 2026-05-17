package booking

import (
	"context"
	"errors"
	"testing"

	"encore.app/internal/api_errors"
	"encore.app/services/booking/db"
	"encore.app/services/booking/handlers/location"
	"go.uber.org/mock/gomock"
)

func TestDeleteLocation(t *testing.T) {
	ctx := context.Background()

	t.Run("deletes broker code and keeps location when other broker codes exist", func(t *testing.T) {
		q := testQuerier()
		s := &Service{query: q}

		cityTLV := "Tel Aviv"
		iataTLV := "TLV"
		loc, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "Israel", CountryCode: "IL", City: &cityTLV, Name: "Ben Gurion Airport", Iata: &iataTLV,
			},
			db.BrokerFlex, "flex-tlv-del-keep",
		)

		// Add a second broker code to the same location
		_, err := q.InsertLocationBrokerCode(ctx, db.InsertLocationBrokerCodeParams{
			LocationID:       loc.ID,
			Broker:           db.BrokerHertz,
			BrokerLocationID: "hertz-tlv-del-keep",
		})
		if err != nil {
			t.Fatalf("failed to insert second broker code: %v", err)
		}

		// Delete the first broker code
		err = s.DeleteLocation(ctx, lbc.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Location should still exist
		_, err = q.GetLocationById(ctx, loc.ID)
		if err != nil {
			t.Fatalf("location should still exist, got error: %v", err)
		}
	})

	t.Run("deletes broker code and location when no other broker codes exist", func(t *testing.T) {
		q := testQuerier()
		s := &Service{query: q}

		cityParis := "Paris"
		iataCDG := "CDG"
		loc, lbc := seedLocationWithBrokerCode(t, q,
			db.InsertLocationParams{
				Country: "France", CountryCode: "FR", City: &cityParis, Name: "Charles de Gaulle", Iata: &iataCDG,
			},
			db.BrokerFlex, "flex-cdg-del-orphan",
		)

		err := s.DeleteLocation(ctx, lbc.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Location should be deleted
		_, err = q.GetLocationById(ctx, loc.ID)
		if !errors.Is(err, db.ErrNoRows) {
			t.Fatalf("expected location to be deleted, got err: %v", err)
		}
	})

	t.Run("returns not found for non-existent broker code", func(t *testing.T) {
		q := testQuerier()
		s := &Service{query: q}

		err := s.DeleteLocation(ctx, 999999)
		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("returns error when delete broker code fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteLocationBrokerCode(gomock.Any(), int64(1)).
			Return(int64(0), errors.New("db error"))

		err := s.DeleteLocation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when count fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteLocationBrokerCode(gomock.Any(), int64(1)).
			Return(int64(10), nil)
		q.EXPECT().CountLocationBrokerCodesByLocationID(gomock.Any(), int64(10)).
			Return(int64(0), errors.New("db error"))

		err := s.DeleteLocation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when delete location fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DeleteLocationBrokerCode(gomock.Any(), int64(1)).
			Return(int64(10), nil)
		q.EXPECT().CountLocationBrokerCodesByLocationID(gomock.Any(), int64(10)).
			Return(int64(0), nil)
		q.EXPECT().DeleteLocationByID(gomock.Any(), int64(10)).
			Return(errors.New("db error"))

		err := s.DeleteLocation(ctx, 1)
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}

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

		err := s.ToggleLocation(ctx, lbc.ID, location.ToggleLocationParams{Enabled: false})
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

		err := s.ToggleLocation(ctx, lbc.ID, location.ToggleLocationParams{Enabled: true})
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

		err := s.ToggleLocation(ctx, lbc.ID, location.ToggleLocationParams{Enabled: true})
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

		err := s.ToggleLocation(ctx, lbc.ID, location.ToggleLocationParams{Enabled: false})
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

		err := s.ToggleLocation(ctx, 999, location.ToggleLocationParams{Enabled: true})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})

	t.Run("returns error when disable fails", func(t *testing.T) {
		q, s := mockService(t)
		q.EXPECT().DisableLocationBrokerCode(gomock.Any(), int64(999)).
			Return(errors.New("db error"))

		err := s.ToggleLocation(ctx, 999, location.ToggleLocationParams{Enabled: false})
		api_errors.AssertApiError(t, api_errors.ErrInternalError, err)
	})
}
